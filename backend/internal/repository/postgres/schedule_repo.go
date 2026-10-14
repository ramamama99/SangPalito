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

type scheduleRepo struct {
	pool *pgxpool.Pool
}

// NewScheduleRepo creates a new PostgreSQL schedule repository.
func NewScheduleRepo(pool *pgxpool.Pool) *scheduleRepo {
	return &scheduleRepo{pool: pool}
}

// FindByID retrieves a schedule by its UUID, joining doctor and polyclinic names.
func (r *scheduleRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error) {
	query := `SELECT s.id, s.facility_id, s.polyclinic_id, s.doctor_id, s.day_of_week,
	           s.start_time, s.end_time, s.quota, s.is_active, s.created_at, s.updated_at,
	           d.full_name AS doctor_name, p.name AS polyclinic_name
	           FROM schedules s
	           JOIN doctors d ON d.id = s.doctor_id
	           JOIN polyclinics p ON p.id = s.polyclinic_id
	           WHERE s.id = $1`

	sc := &entity.Schedule{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sc.ID, &sc.FacilityID, &sc.PolyclinicID, &sc.DoctorID, &sc.DayOfWeek,
		&sc.StartTime, &sc.EndTime, &sc.Quota, &sc.IsActive, &sc.CreatedAt, &sc.UpdatedAt,
		&sc.DoctorName, &sc.PolyclinicName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Jadwal", err)
		}
		return nil, fmt.Errorf("scheduleRepo.FindByID: %w", err)
	}
	return sc, nil
}

// FindByPolyclinic retrieves all active schedules for a polyclinic.
func (r *scheduleRepo) FindByPolyclinic(ctx context.Context, polyclinicID uuid.UUID) ([]*entity.Schedule, error) {
	query := `SELECT s.id, s.facility_id, s.polyclinic_id, s.doctor_id, s.day_of_week,
	           s.start_time, s.end_time, s.quota, s.is_active, s.created_at, s.updated_at,
	           d.full_name AS doctor_name, p.name AS polyclinic_name
	           FROM schedules s
	           JOIN doctors d ON d.id = s.doctor_id
	           JOIN polyclinics p ON p.id = s.polyclinic_id
	           WHERE s.polyclinic_id = $1 AND s.is_active = true
	           ORDER BY s.day_of_week, s.start_time`

	return r.scanSchedules(ctx, query, polyclinicID)
}

// FindByPolyclinicAndDay retrieves schedules for a polyclinic on a specific day of week.
func (r *scheduleRepo) FindByPolyclinicAndDay(ctx context.Context, polyclinicID uuid.UUID, dayOfWeek int) ([]*entity.Schedule, error) {
	query := `SELECT s.id, s.facility_id, s.polyclinic_id, s.doctor_id, s.day_of_week,
	           s.start_time, s.end_time, s.quota, s.is_active, s.created_at, s.updated_at,
	           d.full_name AS doctor_name, p.name AS polyclinic_name
	           FROM schedules s
	           JOIN doctors d ON d.id = s.doctor_id
	           JOIN polyclinics p ON p.id = s.polyclinic_id
	           WHERE s.polyclinic_id = $1 AND s.day_of_week = $2 AND s.is_active = true
	           ORDER BY s.start_time`

	rows, err := r.pool.Query(ctx, query, polyclinicID, dayOfWeek)
	if err != nil {
		return nil, fmt.Errorf("scheduleRepo.FindByPolyclinicAndDay: %w", err)
	}
	defer rows.Close()

	return r.collectSchedules(rows)
}

// GetRemainingQuota calculates how many slots are still available for a schedule on a date.
// Uses SELECT ... FOR UPDATE to prevent race conditions during booking.
func (r *scheduleRepo) GetRemainingQuota(ctx context.Context, tx pgx.Tx, scheduleID uuid.UUID, date time.Time) (int, error) {
	// First get the schedule's total quota (lock the row)
	var quota int
	quotaQuery := `SELECT quota FROM schedules WHERE id = $1 FOR UPDATE`
	if err := tx.QueryRow(ctx, quotaQuery, scheduleID).Scan(&quota); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, apperror.NotFound("Jadwal", err)
		}
		return 0, fmt.Errorf("scheduleRepo.GetRemainingQuota (quota): %w", err)
	}

	// Count existing appointments for that schedule on that date
	var booked int
	countQuery := `SELECT COUNT(*) FROM appointments WHERE schedule_id = $1 AND visit_date = $2 AND status != 'cancelled'`
	if err := tx.QueryRow(ctx, countQuery, scheduleID, date).Scan(&booked); err != nil {
		return 0, fmt.Errorf("scheduleRepo.GetRemainingQuota (count): %w", err)
	}

	remaining := quota - booked
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

func (r *scheduleRepo) scanSchedules(ctx context.Context, query string, args ...interface{}) ([]*entity.Schedule, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("scheduleRepo.scanSchedules: %w", err)
	}
	defer rows.Close()
	return r.collectSchedules(rows)
}

func (r *scheduleRepo) collectSchedules(rows pgx.Rows) ([]*entity.Schedule, error) {
	var schedules []*entity.Schedule
	for rows.Next() {
		sc := &entity.Schedule{}
		if err := rows.Scan(
			&sc.ID, &sc.FacilityID, &sc.PolyclinicID, &sc.DoctorID, &sc.DayOfWeek,
			&sc.StartTime, &sc.EndTime, &sc.Quota, &sc.IsActive, &sc.CreatedAt, &sc.UpdatedAt,
			&sc.DoctorName, &sc.PolyclinicName,
		); err != nil {
			return nil, fmt.Errorf("scheduleRepo.collectSchedules scan: %w", err)
		}
		schedules = append(schedules, sc)
	}
	return schedules, nil
}
