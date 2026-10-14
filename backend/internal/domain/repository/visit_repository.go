package repository

import (
	"context"
	"sangpalito/backend/internal/domain/entity"

	"github.com/google/uuid"
)

type VisitRepository interface {
	Create(ctx context.Context, visit *entity.Visit, diagnoses []entity.Diagnosis) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Visit, []entity.Diagnosis, error)
	FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Visit, error)
	Update(ctx context.Context, visit *entity.Visit) error
}
