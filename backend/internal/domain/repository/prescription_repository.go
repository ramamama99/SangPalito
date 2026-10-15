package repository

import (
	"context"

	"github.com/google/uuid"

	"sangpalito/backend/internal/domain/entity"
)

type PrescriptionRepository interface {
	Create(ctx context.Context, p *entity.Prescription) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Prescription, error)
	FindByPrescriptionNo(ctx context.Context, no string) (*entity.Prescription, error)
	FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Prescription, error)
	FindByVisitID(ctx context.Context, visitID uuid.UUID) (*entity.Prescription, error)
	FindByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*entity.Prescription, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, dispensedBy *uuid.UUID) error
}
