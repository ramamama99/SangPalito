package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"sangpalito/backend/internal/domain/entity"
)

// ScheduleRepository defines the contract for schedule data access.
type ScheduleRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error)
	FindByPolyclinic(ctx context.Context, polyclinicID uuid.UUID) ([]*entity.Schedule, error)
	FindByPolyclinicAndDay(ctx context.Context, polyclinicID uuid.UUID, dayOfWeek int) ([]*entity.Schedule, error)
	GetRemainingQuota(ctx context.Context, tx pgx.Tx, scheduleID uuid.UUID, date time.Time) (int, error)
}
