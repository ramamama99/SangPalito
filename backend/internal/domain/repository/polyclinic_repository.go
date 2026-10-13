package repository

import (
	"context"

	"sangpalito/backend/internal/domain/entity"

	"github.com/google/uuid"
)

type PolyclinicRepository interface {
	FindAll(ctx context.Context, facilityID uuid.UUID) ([]entity.Polyclinic, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Polyclinic, error)
	Create(ctx context.Context, p *entity.Polyclinic) error
	Update(ctx context.Context, p *entity.Polyclinic) error
}
