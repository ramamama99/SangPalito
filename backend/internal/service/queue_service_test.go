package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/worker"
)

type mockQueueRepo struct {
	queues map[uuid.UUID]*entity.Queue
}

func newMockQueueRepo() *mockQueueRepo {
	return &mockQueueRepo{queues: make(map[uuid.UUID]*entity.Queue)}
}

func (m *mockQueueRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Queue, error) {
	if q, ok := m.queues[id]; ok {
		return q, nil
	}
	return nil, apperror.NotFound("Queue", nil)
}

func (m *mockQueueRepo) FindActiveByPolyclinic(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]*entity.Queue, error) {
	var list []*entity.Queue
	for _, q := range m.queues {
		if q.PolyclinicID == polyclinicID {
			list = append(list, q)
		}
	}
	return list, nil
}

func (m *mockQueueRepo) Create(ctx context.Context, q *entity.Queue) error {
	m.queues[q.ID] = q
	return nil
}

func (m *mockQueueRepo) CreateTx(ctx context.Context, tx pgx.Tx, q *entity.Queue) error {
	return m.Create(ctx, q)
}

func (m *mockQueueRepo) UpdateStatus(ctx context.Context, q *entity.Queue) error {
	m.queues[q.ID] = q
	return nil
}

func (m *mockQueueRepo) GetNextQueueNumber(ctx context.Context, tx pgx.Tx, polyclinicID uuid.UUID, date time.Time) (int, error) {
	return len(m.queues) + 1, nil
}

func (m *mockQueueRepo) CountByPolyclinicAndDate(ctx context.Context, polyclinicID uuid.UUID, date time.Time, status entity.QueueStatus) (int, error) {
	count := 0
	for _, q := range m.queues {
		if q.PolyclinicID == polyclinicID && q.Status == status {
			count++
		}
	}
	return count, nil
}

func (m *mockQueueRepo) FindCompletedDurations(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]float64, error) {
	return nil, nil
}

func (m *mockQueueRepo) FindActiveQueueDetails(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]*entity.QueueWithPatient, error) {
	var list []*entity.QueueWithPatient
	for _, q := range m.queues {
		if q.PolyclinicID == polyclinicID {
			list = append(list, &entity.QueueWithPatient{
				Queue:      q,
				PatientAge: 30,
				Priority:   "normal",
			})
		}
	}
	return list, nil
}

func (m *mockQueueRepo) FindQueueWithPatientByID(ctx context.Context, queueID uuid.UUID) (*entity.QueueWithPatient, error) {
	if q, ok := m.queues[queueID]; ok {
		return &entity.QueueWithPatient{
			Queue:      q,
			PatientAge: 30,
			Priority:   "normal",
		}, nil
	}
	return nil, apperror.NotFound("Queue", nil)
}

// mockBroadcaster records broadcasts for inspection in unit tests.
type mockBroadcaster struct {
	payloads []any
}

func (m *mockBroadcaster) Broadcast(payload any) {
	m.payloads = append(m.payloads, payload)
}

// mockNotifier records notification jobs for unit tests.
type mockNotifier struct {
	jobs []worker.NotificationJob
}

func (m *mockNotifier) Enqueue(job worker.NotificationJob) {
	m.jobs = append(m.jobs, job)
}

func TestQueueService_AdvanceQueueStatus(t *testing.T) {
	ctx := context.Background()
	polyclinicID := uuid.New()
	facilityID := uuid.New()

	setupQueue := func(status entity.QueueStatus, callCount int) (*queueService, *mockQueueRepo, *mockBroadcaster, *mockNotifier, uuid.UUID) {
		repo := newMockQueueRepo()
		broadcaster := &mockBroadcaster{}
		notifier := &mockNotifier{}
		auditRepo := &mockAuditLogRepo{}

		qID := uuid.New()
		repo.queues[qID] = &entity.Queue{
			ID:           qID,
			FacilityID:   facilityID,
			PolyclinicID: polyclinicID,
			QueueNumber:  5,
			Status:       status,
			CallCount:    callCount,
			QueueDate:    time.Now(),
		}

		svc := &queueService{
			queueRepo:    repo,
			auditLogRepo: auditRepo,
			broadcaster:  broadcaster,
			notifier:     notifier,
		}

		return svc, repo, broadcaster, notifier, qID
	}

	t.Run("Valid Transition Waiting to CheckedIn", func(t *testing.T) {
		svc, _, broadcaster, _, qID := setupQueue(entity.StatusWaiting, 0)

		resp, err := svc.AdvanceQueueStatus(ctx, qID, entity.StatusCheckedIn)
		if err != nil {
			t.Fatalf("expected valid transition, got: %v", err)
		}
		if resp.Status != string(entity.StatusCheckedIn) {
			t.Errorf("expected status checked_in, got %s", resp.Status)
		}
		if len(broadcaster.payloads) == 0 {
			t.Error("expected broadcast to be sent")
		}
	})

	t.Run("Invalid Transition Waiting to Done Fails", func(t *testing.T) {
		svc, _, _, _, qID := setupQueue(entity.StatusWaiting, 0)

		_, err := svc.AdvanceQueueStatus(ctx, qID, entity.StatusDone)
		if err == nil {
			t.Fatal("expected error on invalid transition waiting -> done, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) || appErr.Code != "UNPROCESSABLE_ENTITY" {
			t.Errorf("expected UNPROCESSABLE_ENTITY error, got %v", err)
		}
	})

	t.Run("Transition to Called Increments CallCount and Triggers Notification", func(t *testing.T) {
		svc, repo, _, notifier, qID := setupQueue(entity.StatusCheckedIn, 0)

		resp, err := svc.AdvanceQueueStatus(ctx, qID, entity.StatusCalled)
		if err != nil {
			t.Fatalf("expected valid transition, got: %v", err)
		}
		if resp.CallCount != 1 {
			t.Errorf("expected call count 1, got %d", resp.CallCount)
		}
		if repo.queues[qID].CalledAt == nil {
			t.Error("expected called_at to be set")
		}
		if len(notifier.jobs) != 1 || notifier.jobs[0].Type != worker.NotifyQueueCalled {
			t.Errorf("expected async notification job to be enqueued, got %v", notifier.jobs)
		}
	})

	t.Run("Transition to NoShow Fails if CallCount < 3", func(t *testing.T) {
		svc, _, _, _, qID := setupQueue(entity.StatusCalled, 2) // baru 2 kali dipanggil

		_, err := svc.AdvanceQueueStatus(ctx, qID, entity.StatusNoShow)
		if err == nil {
			t.Fatal("expected error when no_show with call_count < 3, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) || appErr.Code != "BAD_REQUEST" {
			t.Errorf("expected BAD_REQUEST error, got %v", err)
		}
	})

	t.Run("Transition to NoShow Succeeds if CallCount >= 3", func(t *testing.T) {
		svc, _, _, _, qID := setupQueue(entity.StatusCalled, 3) // sudah 3 kali dipanggil

		resp, err := svc.AdvanceQueueStatus(ctx, qID, entity.StatusNoShow)
		if err != nil {
			t.Fatalf("expected no_show to succeed after 3 calls, got: %v", err)
		}
		if resp.Status != string(entity.StatusNoShow) {
			t.Errorf("expected status no_show, got %s", resp.Status)
		}
	})
}
