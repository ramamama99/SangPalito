package repository

import (
	"context"

	"sangpalito/backend/internal/domain/entity"

	"github.com/google/uuid"
)

type FacilityRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Facility, error)
	GetFirst(ctx context.Context) (*entity.Facility, error) // untuk default faskes di MVP
	Create(ctx context.Context, f *entity.Facility) error
}
