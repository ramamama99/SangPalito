package repository

import (
	"context"

	"github.com/google/uuid"

	"sangpalito/backend/internal/domain/entity"
)

// PolyclinicRepository defines the contract for polyclinic data access.
type PolyclinicRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Polyclinic, error)
	FindByFacilityID(ctx context.Context, facilityID uuid.UUID) ([]*entity.Polyclinic, error)
	FindActive(ctx context.Context, facilityID uuid.UUID) ([]*entity.Polyclinic, error)
}
