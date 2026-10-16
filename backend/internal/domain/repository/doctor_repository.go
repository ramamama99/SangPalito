package repository

import (
	"context"

	"sangpalito/backend/internal/domain/entity"

	"github.com/google/uuid"
)

type DoctorRepository interface {
	FindAll(ctx context.Context, facilityID uuid.UUID) ([]entity.Doctor, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Doctor, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.Doctor, error)
	Create(ctx context.Context, d *entity.Doctor) error
	Update(ctx context.Context, d *entity.Doctor) error
	Delete(ctx context.Context, id uuid.UUID) error
}

