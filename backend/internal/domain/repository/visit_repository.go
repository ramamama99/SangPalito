package repository

import (
	"context"

	"github.com/google/uuid"

	"sangpalito/backend/internal/domain/entity"
)

// VisitRepository defines the contract for visit data access.
type VisitRepository interface {
	FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Visit, error)
	Create(ctx context.Context, v *entity.Visit) error
}
