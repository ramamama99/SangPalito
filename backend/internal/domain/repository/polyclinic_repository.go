package repository

import (
	"context"

	"github.com/google/uuid"

	"sangpalito/backend/internal/domain/entity"
)

// PolyclinicRepository defines the contract for polyclinic data access.
type PolyclinicRepository interface {
	FindAll(ctx context.Context, facilityID uuid.UUID) ([]entity.Polyclinic, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Polyclinic, error)
	FindByFacilityID(ctx context.Context, facilityID uuid.UUID) ([]*entity.Polyclinic, error)
	FindActive(ctx context.Context, facilityID uuid.UUID) ([]*entity.Polyclinic, error)
	Create(ctx context.Context, p *entity.Polyclinic) error
	Update(ctx context.Context, p *entity.Polyclinic) error
}
