package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/internal/worker"
)

// QueueBroadcaster defines an abstraction for realtime message broadcasting.
type QueueBroadcaster interface {
	Broadcast(payload any)
}

// QueueNotifier defines an abstraction for asynchronous notification dispatch.
type QueueNotifier interface {
	Enqueue(job worker.NotificationJob)
}

// QueueService defines the contract for queue management operations.
type QueueService interface {
	GetActiveQueues(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]*dto.QueueResponse, error)
	AdvanceQueueStatus(ctx context.Context, queueID uuid.UUID, newStatus entity.QueueStatus) (*dto.QueueResponse, error)
}

type queueService struct {
	queueRepo    repository.QueueRepository
	auditLogRepo repository.AuditLogRepository
	broadcaster  QueueBroadcaster
	notifier     QueueNotifier
}

// NewQueueService creates a new queue service with optional realtime and notification hooks.
func NewQueueService(
	queueRepo repository.QueueRepository,
	auditLogRepo repository.AuditLogRepository,
	broadcaster QueueBroadcaster,
	notifier QueueNotifier,
) QueueService {
	return &queueService{
		queueRepo:    queueRepo,
		auditLogRepo: auditLogRepo,
		broadcaster:  broadcaster,
		notifier:     notifier,
	}
}

// GetActiveQueues retrieves all queues for a polyclinic on a given date.
func (s *queueService) GetActiveQueues(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]*dto.QueueResponse, error) {
	queues, err := s.queueRepo.FindActiveByPolyclinic(ctx, polyclinicID, date)
	if err != nil {
		return nil, apperror.Internal(fmt.Errorf("queueService.GetActiveQueues: %w", err))
	}

	var results []*dto.QueueResponse
	for _, q := range queues {
		results = append(results, toQueueResponse(q))
	}
	return results, nil
}

// AdvanceQueueStatus changes the status of a queue entry, enforcing the state machine.
func (s *queueService) AdvanceQueueStatus(ctx context.Context, queueID uuid.UUID, newStatus entity.QueueStatus) (*dto.QueueResponse, error) {
	queue, err := s.queueRepo.FindByID(ctx, queueID)
	if err != nil {
		return nil, err
	}

	// Validate state transition
	if !queue.Status.CanTransitionTo(newStatus) {
		return nil, apperror.UnprocessableEntity(
			fmt.Sprintf("Tidak bisa mengubah status dari '%s' ke '%s'", queue.Status, newStatus), nil,
		)
	}

	now := time.Now()

	// Apply side effects based on new status
	switch newStatus {
	case entity.StatusCalled:
		queue.CallCount++
		queue.CalledAt = &now
	case entity.StatusInService:
		queue.ServiceStart = &now
	case entity.StatusDone:
		queue.ServiceEnd = &now
	case entity.StatusNoShow:
		if queue.CallCount < 3 {
			return nil, apperror.BadRequest(
				fmt.Sprintf("Pasien baru dipanggil %d kali. Minimal 3 kali sebelum no_show.", queue.CallCount), nil,
			)
		}
	}

	queue.Status = newStatus

	if err := s.queueRepo.UpdateStatus(ctx, queue); err != nil {
		return nil, apperror.Internal(fmt.Errorf("queueService.AdvanceQueueStatus: %w", err))
	}

	resp := toQueueResponse(queue)

	// Broadcast status update to realtime WebSocket clients
	if s.broadcaster != nil {
		s.broadcaster.Broadcast(map[string]any{
			"event": "queue_updated",
			"data":  resp,
		})
	}

	// Trigger async notification when patient is called
	if newStatus == entity.StatusCalled && s.notifier != nil {
		s.notifier.Enqueue(worker.NotificationJob{
			Type:        worker.NotifyQueueCalled,
			QueueID:     queue.ID,
			QueueNumber: queue.QueueNumber,
			VisitDate:   queue.QueueDate,
		})
	}

	// Record audit log
	if s.auditLogRepo != nil {
		entry := &entity.AuditLog{
			ID:        uuid.New(),
			Action:    "update_status",
			Entity:    "queue",
			EntityID:  &queue.ID,
			Meta:      map[string]interface{}{"new_status": string(newStatus), "queue_number": queue.QueueNumber},
			CreatedAt: time.Now(),
		}
		_ = s.auditLogRepo.Create(ctx, entry)
	}

	slog.Info("queue status changed",
		"queue_id", queueID,
		"new_status", newStatus,
		"queue_number", queue.QueueNumber,
	)

	return resp, nil
}

func toQueueResponse(q *entity.Queue) *dto.QueueResponse {
	return &dto.QueueResponse{
		ID:            q.ID.String(),
		FacilityID:    q.FacilityID.String(),
		AppointmentID: q.AppointmentID.String(),
		PolyclinicID:  q.PolyclinicID.String(),
		QueueNumber:   q.QueueNumber,
		Status:        string(q.Status),
		CallCount:     q.CallCount,
		CalledAt:      q.CalledAt,
		ServiceStart:  q.ServiceStart,
		ServiceEnd:    q.ServiceEnd,
		QueueDate:     q.QueueDate.Format("2006-01-02"),
		CreatedAt:     q.CreatedAt,
		UpdatedAt:     q.UpdatedAt,
	}
}
