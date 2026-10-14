package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
)

type appointmentRepo struct {
	pool *pgxpool.Pool
}

// NewAppointmentRepo creates a new PostgreSQL appointment repository.
func NewAppointmentRepo(pool *pgxpool.Pool) *appointmentRepo {
	return &appointmentRepo{pool: pool}
}

const appointmentCols = `id, facility_id, patient_id, polyclinic_id, schedule_id, doctor_id,
	visit_date, appointment_token, qr_token, chief_complaint, status, registered_by, created_at, updated_at`

func scanAppointment(row pgx.Row) (*entity.Appointment, error) {
	a := &entity.Appointment{}
	err := row.Scan(
		&a.ID, &a.FacilityID, &a.PatientID, &a.PolyclinicID, &a.ScheduleID,
		&a.DoctorID, &a.VisitDate, &a.AppointmentToken, &a.QRToken,
		&a.ChiefComplaint, &a.Status, &a.RegisteredBy, &a.CreatedAt, &a.UpdatedAt,
	)
	return a, err
}

// FindByID retrieves an appointment by its UUID.
func (r *appointmentRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Appointment, error) {
	query := fmt.Sprintf("SELECT %s FROM appointments WHERE id = $1", appointmentCols)
	a, err := scanAppointment(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Appointment", err)
		}
		return nil, fmt.Errorf("appointmentRepo.FindByID: %w", err)
	}
	return a, nil
}

// FindByToken retrieves an appointment by its public token.
func (r *appointmentRepo) FindByToken(ctx context.Context, token uuid.UUID) (*entity.Appointment, error) {
	query := fmt.Sprintf("SELECT %s FROM appointments WHERE appointment_token = $1", appointmentCols)
	a, err := scanAppointment(r.pool.QueryRow(ctx, query, token))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Appointment", err)
		}
		return nil, fmt.Errorf("appointmentRepo.FindByToken: %w", err)
	}
	return a, nil
}

// FindByQRToken retrieves an appointment by its QR token.
func (r *appointmentRepo) FindByQRToken(ctx context.Context, qrToken uuid.UUID) (*entity.Appointment, error) {
	query := fmt.Sprintf("SELECT %s FROM appointments WHERE qr_token = $1", appointmentCols)
	a, err := scanAppointment(r.pool.QueryRow(ctx, query, qrToken))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Appointment", err)
		}
		return nil, fmt.Errorf("appointmentRepo.FindByQRToken: %w", err)
	}
	return a, nil
}

// FindByPatientPolyDate checks if a patient already registered for a poli on a specific date.
func (r *appointmentRepo) FindByPatientPolyDate(ctx context.Context, patientID, polyclinicID uuid.UUID, date time.Time) (*entity.Appointment, error) {
	query := fmt.Sprintf("SELECT %s FROM appointments WHERE patient_id = $1 AND polyclinic_id = $2 AND visit_date = $3", appointmentCols)
	a, err := scanAppointment(r.pool.QueryRow(ctx, query, patientID, polyclinicID, date))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found = no duplicate
		}
		return nil, fmt.Errorf("appointmentRepo.FindByPatientPolyDate: %w", err)
	}
	return a, nil
}

// FindByPatientID retrieves all appointments for a patient.
func (r *appointmentRepo) FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Appointment, error) {
	query := fmt.Sprintf("SELECT %s FROM appointments WHERE patient_id = $1 ORDER BY visit_date DESC", appointmentCols)
	rows, err := r.pool.Query(ctx, query, patientID)
	if err != nil {
		return nil, fmt.Errorf("appointmentRepo.FindByPatientID: %w", err)
	}
	defer rows.Close()

	var appointments []*entity.Appointment
	for rows.Next() {
		a := &entity.Appointment{}
		if err := rows.Scan(
			&a.ID, &a.FacilityID, &a.PatientID, &a.PolyclinicID, &a.ScheduleID,
			&a.DoctorID, &a.VisitDate, &a.AppointmentToken, &a.QRToken,
			&a.ChiefComplaint, &a.Status, &a.RegisteredBy, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("appointmentRepo.FindByPatientID scan: %w", err)
		}
		appointments = append(appointments, a)
	}
	return appointments, nil
}

// Create inserts a new appointment.
func (r *appointmentRepo) Create(ctx context.Context, a *entity.Appointment) error {
	return r.createWith(ctx, r.pool, a)
}

// CreateTx inserts a new appointment within a transaction.
func (r *appointmentRepo) CreateTx(ctx context.Context, tx pgx.Tx, a *entity.Appointment) error {
	return r.createWith(ctx, tx, a)
}

func (r *appointmentRepo) createWith(ctx context.Context, exec interface{}, a *entity.Appointment) error {
	query := `INSERT INTO appointments (id, facility_id, patient_id, polyclinic_id, schedule_id, doctor_id,
		visit_date, appointment_token, qr_token, chief_complaint, status, registered_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	now := time.Now()
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.AppointmentToken == uuid.Nil {
		a.AppointmentToken = uuid.New()
	}
	if a.QRToken == uuid.Nil {
		a.QRToken = uuid.New()
	}
	a.CreatedAt = now
	a.UpdatedAt = now

	args := []interface{}{
		a.ID, a.FacilityID, a.PatientID, a.PolyclinicID, a.ScheduleID,
		a.DoctorID, a.VisitDate, a.AppointmentToken, a.QRToken,
		a.ChiefComplaint, a.Status, a.RegisteredBy, a.CreatedAt, a.UpdatedAt,
	}

	var err error
	switch e := exec.(type) {
	case *pgxpool.Pool:
		_, err = e.Exec(ctx, query, args...)
	case pgx.Tx:
		_, err = e.Exec(ctx, query, args...)
	}
	if err != nil {
		return fmt.Errorf("appointmentRepo.Create: %w", err)
	}
	return nil
}

// UpdateStatus updates the status of an appointment.
func (r *appointmentRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE appointments SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("appointmentRepo.UpdateStatus: %w", err)
	}
	return nil
}
