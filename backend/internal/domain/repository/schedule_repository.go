package repository

import (
	"context"
	"time"

	"sangpalito/backend/internal/domain/entity"

	"github.com/google/uuid"
)

type ScheduleRepository interface {
	// FindWithRemainingSlot mencari jadwal dokter dan otomatis menghitung sisa kuota hari itu
	FindWithRemainingSlot(ctx context.Context, polyclinicID uuid.UUID, date time.Time, dayOfWeek int) ([]entity.Schedule, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error)
	Create(ctx context.Context, s *entity.Schedule) error
	Delete(ctx context.Context, id uuid.UUID) error
}
