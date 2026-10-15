package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"sangpalito/backend/internal/domain/entity"
)

// QueueRepository defines the contract for queue data access.
type QueueRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Queue, error)
	FindActiveByPolyclinic(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]*entity.Queue, error)
	Create(ctx context.Context, q *entity.Queue) error
	CreateTx(ctx context.Context, tx pgx.Tx, q *entity.Queue) error
	UpdateStatus(ctx context.Context, q *entity.Queue) error
	GetNextQueueNumber(ctx context.Context, tx pgx.Tx, polyclinicID uuid.UUID, date time.Time) (int, error)
	CountByPolyclinicAndDate(ctx context.Context, polyclinicID uuid.UUID, date time.Time, status entity.QueueStatus) (int, error)
	FindCompletedDurations(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]float64, error)
	FindActiveQueueDetails(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]*entity.QueueWithPatient, error)
	FindQueueWithPatientByID(ctx context.Context, queueID uuid.UUID) (*entity.QueueWithPatient, error)
}
