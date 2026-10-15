package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"sangpalito/backend/internal/domain/entity"
)

// AppointmentRepository defines the contract for appointment data access.
type AppointmentRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Appointment, error)
	FindByToken(ctx context.Context, token uuid.UUID) (*entity.Appointment, error)
	FindByQRToken(ctx context.Context, qrToken string) (*entity.Appointment, error)
	FindByPatientPolyDate(ctx context.Context, patientID, polyclinicID uuid.UUID, date time.Time) (*entity.Appointment, error)
	FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Appointment, error)
	FindToday(ctx context.Context, facilityID uuid.UUID, polyclinicID *uuid.UUID, date string) ([]entity.Appointment, error)
	Create(ctx context.Context, a *entity.Appointment) error
	CreateTx(ctx context.Context, tx pgx.Tx, a *entity.Appointment) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Cancel(ctx context.Context, id uuid.UUID) error
}