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
	"sangpalito/backend/internal/domain/repository"
)

type appointmentRepo struct {
	pool *pgxpool.Pool
}

func NewAppointmentRepo(pool *pgxpool.Pool) repository.AppointmentRepository {
	return &appointmentRepo{pool: pool}
}

const appointmentCols = `id, facility_id, patient_id, polyclinic_id, schedule_id,
	visit_date, source, priority, appointment_token, qr_token, created_by, created_at`

func scanAppointment(row pgx.Row) (*entity.Appointment, error) {
	a := &entity.Appointment{}
	var qrTokenStr string
	var schedID *uuid.UUID
	err := row.Scan(
		&a.ID, &a.FacilityID, &a.PatientID, &a.PolyclinicID, &schedID,
		&a.VisitDate, &a.Source, &a.Priority, &a.AppointmentToken, &qrTokenStr,
		&a.RegisteredBy, &a.CreatedAt,
	)
	if err == nil {
		a.QRToken, _ = uuid.Parse(qrTokenStr)
		if schedID != nil {
			a.ScheduleID = *schedID
		}
	}
	return a, err
}

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

func (r *appointmentRepo) FindByQRToken(ctx context.Context, qrToken uuid.UUID) (*entity.Appointment, error) {
	query := fmt.Sprintf("SELECT %s FROM appointments WHERE qr_token = $1", appointmentCols)
	a, err := scanAppointment(r.pool.QueryRow(ctx, query, qrToken.String()))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Appointment", err)
		}
		return nil, fmt.Errorf("appointmentRepo.FindByQRToken: %w", err)
	}
	return a, nil
}

func (r *appointmentRepo) FindByPatientPolyDate(ctx context.Context, patientID, polyclinicID uuid.UUID, date time.Time) (*entity.Appointment, error) {
	query := fmt.Sprintf("SELECT %s FROM appointments WHERE patient_id = $1 AND polyclinic_id = $2 AND visit_date = $3", appointmentCols)
	a, err := scanAppointment(r.pool.QueryRow(ctx, query, patientID, polyclinicID, date))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("appointmentRepo.FindByPatientPolyDate: %w", err)
	}
	return a, nil
}

func (r *appointmentRepo) FindByPatientID(ctx context.Context, patientID uuid.UUID) ([]*entity.Appointment, error) {
	query := fmt.Sprintf("SELECT %s FROM appointments WHERE patient_id = $1 ORDER BY visit_date DESC", appointmentCols)
	rows, err := r.pool.Query(ctx, query, patientID)
	if err != nil {
		return nil, fmt.Errorf("appointmentRepo.FindByPatientID: %w", err)
	}
	defer rows.Close()

	var appointments []*entity.Appointment
	for rows.Next() {
		a, err := scanAppointment(rows)
		if err != nil {
			return nil, fmt.Errorf("appointmentRepo.FindByPatientID scan: %w", err)
		}
		appointments = append(appointments, a)
	}
	return appointments, nil
}

func (r *appointmentRepo) FindToday(ctx context.Context, facilityID uuid.UUID, polyclinicID *uuid.UUID, date string) ([]entity.Appointment, error) {
	query := `
		SELECT 
			a.id, a.facility_id, a.patient_id, a.polyclinic_id, a.schedule_id, 
			a.visit_date, a.source, a.priority, a.appointment_token, a.qr_token, a.created_by, a.created_at,
			p.full_name, poly.name
		FROM appointments a
		JOIN patients p ON p.id = a.patient_id
		JOIN polyclinics poly ON poly.id = a.polyclinic_id
		WHERE a.facility_id = $1 AND a.visit_date = $2::date
	`
	args := []interface{}{facilityID, date}
	if polyclinicID != nil {
		query += " AND a.polyclinic_id = $3"
		args = append(args, *polyclinicID)
	}
	query += " ORDER BY a.created_at ASC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("appointmentRepo.FindToday: %w", err)
	}
	defer rows.Close()

	var list []entity.Appointment
	for rows.Next() {
		var a entity.Appointment
		var qrStr string
		var schedID *uuid.UUID
		err := rows.Scan(
			&a.ID, &a.FacilityID, &a.PatientID, &a.PolyclinicID, &schedID,
			&a.VisitDate, &a.Source, &a.Priority, &a.AppointmentToken, &qrStr, &a.RegisteredBy, &a.CreatedAt,
			&a.PatientFullName, &a.PolyclinicName,
		)
		if err != nil {
			return nil, fmt.Errorf("appointmentRepo.FindToday scan: %w", err)
		}
		if schedID != nil {
			a.ScheduleID = *schedID
		}
		a.QRToken, _ = uuid.Parse(qrStr)
		list = append(list, a)
	}
	return list, nil
}

func (r *appointmentRepo) Create(ctx context.Context, a *entity.Appointment) error {
	return r.createWith(ctx, r.pool, a)
}

func (r *appointmentRepo) CreateTx(ctx context.Context, tx pgx.Tx, a *entity.Appointment) error {
	return r.createWith(ctx, tx, a)
}

func (r *appointmentRepo) createWith(ctx context.Context, exec interface{}, a *entity.Appointment) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.AppointmentToken == uuid.Nil {
		a.AppointmentToken = uuid.New()
	}
	if a.QRToken == uuid.Nil {
		a.QRToken = uuid.New()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}

	source := "online"
	if a.Source != "" {
		source = a.Source
	}
	priority := "normal"
	if a.Priority != "" {
		priority = a.Priority
	}

	var schedID *uuid.UUID
	if a.ScheduleID != uuid.Nil {
		schedID = &a.ScheduleID
	}

	query := `
		INSERT INTO appointments (
			id, facility_id, patient_id, polyclinic_id, schedule_id, visit_date, 
			source, priority, appointment_token, qr_token, created_by, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	var err error
	switch e := exec.(type) {
	case *pgxpool.Pool:
		_, err = e.Exec(ctx, query,
			a.ID, a.FacilityID, a.PatientID, a.PolyclinicID, schedID, a.VisitDate,
			source, priority, a.AppointmentToken, a.QRToken.String(), a.RegisteredBy, a.CreatedAt,
		)
	case pgx.Tx:
		_, err = e.Exec(ctx, query,
			a.ID, a.FacilityID, a.PatientID, a.PolyclinicID, schedID, a.VisitDate,
			source, priority, a.AppointmentToken, a.QRToken.String(), a.RegisteredBy, a.CreatedAt,
		)
	default:
		return fmt.Errorf("appointmentRepo.createWith: unsupported executor")
	}

	if err != nil {
		return fmt.Errorf("appointmentRepo.createWith: %w", err)
	}
	return nil
}

func (r *appointmentRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE queues SET status = $2 WHERE appointment_id = $1`
	_, err := r.pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("appointmentRepo.UpdateStatus: %w", err)
	}
	return nil
}

func (r *appointmentRepo) Cancel(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE queues 
		SET status = 'cancelled' 
		WHERE appointment_id = $1 AND status NOT IN ('done', 'cancelled')
	`
	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("appointmentRepo.Cancel: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return apperror.BadRequest("Kunjungan sudah selesai atau sudah dibatalkan", nil)
	}
	return nil
}
