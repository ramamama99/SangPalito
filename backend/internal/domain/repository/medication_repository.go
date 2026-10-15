package repository

import (
	"context"

	"github.com/google/uuid"

	"sangpalito/backend/internal/domain/entity"
)

type MedicationRepository interface {
	FindAll(ctx context.Context, search, category, route string) ([]*entity.Medication, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Medication, error)
	Create(ctx context.Context, m *entity.Medication) error
	Update(ctx context.Context, m *entity.Medication) error
	Delete(ctx context.Context, id uuid.UUID) error
}

