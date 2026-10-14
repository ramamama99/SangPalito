package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"sangpalito/backend/internal/domain/entity"
)

// PatientRepository defines the contract for patient data access.
type PatientRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Patient, error)
	FindByNIKHash(ctx context.Context, nikHash string, facilityID uuid.UUID) (*entity.Patient, error)
	Search(ctx context.Context, facilityID uuid.UUID, query string, page, limit int) ([]*entity.Patient, int, error)
	Create(ctx context.Context, p *entity.Patient) error
	CreateTx(ctx context.Context, tx pgx.Tx, p *entity.Patient) error
	Update(ctx context.Context, p *entity.Patient) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	GenerateMedicalRecordNo(ctx context.Context, date time.Time) (string, error)
	GenerateMedicalRecordNoTx(ctx context.Context, tx pgx.Tx, date time.Time) (string, error)
}
