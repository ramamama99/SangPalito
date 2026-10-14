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

type scheduleRepo struct {
	pool *pgxpool.Pool
}

func NewScheduleRepo(pool *pgxpool.Pool) repository.ScheduleRepository {
	return &scheduleRepo{pool: pool}
}

func (r *scheduleRepo) FindWithRemainingSlot(ctx context.Context, polyclinicID uuid.UUID, date time.Time, dayOfWeek int) ([]entity.Schedule, error) {
	query := `
		SELECT 
			s.id, s.polyclinic_id, s.doctor_id, s.day_of_week,
			s.start_time::text, s.end_time::text, s.quota,
			GREATEST(0, s.quota - COUNT(q.id) FILTER (WHERE q.status != 'cancelled'))::int AS remaining_slot
		FROM schedules s
		LEFT JOIN appointments a ON a.schedule_id = s.id AND a.visit_date = $2
		LEFT JOIN queues q ON q.appointment_id = a.id
		WHERE s.polyclinic_id = $1 AND s.day_of_week = $3
		GROUP BY s.id
		ORDER BY s.start_time ASC
	`
	rows, err := r.pool.Query(ctx, query, polyclinicID, date, dayOfWeek)
	if err != nil {
		return nil, fmt.Errorf("scheduleRepo.FindWithRemainingSlot: %w", err)
	}
	defer rows.Close()

	var list []entity.Schedule
	for rows.Next() {
		var s entity.Schedule
		err := rows.Scan(
			&s.ID, &s.PolyclinicID, &s.DoctorID, &s.DayOfWeek,
			&s.StartTime, &s.EndTime, &s.Quota, &s.RemainingSlot,
		)
		if err != nil {
			return nil, fmt.Errorf("scheduleRepo.FindWithRemainingSlot scan: %w", err)
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *scheduleRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error) {
	query := `
		SELECT s.id, s.polyclinic_id, s.doctor_id, s.day_of_week, s.start_time::text, s.end_time::text, s.quota,
		       d.full_name AS doctor_name, p.name AS polyclinic_name
		FROM schedules s
		JOIN doctors d ON d.id = s.doctor_id
		JOIN polyclinics p ON p.id = s.polyclinic_id
		WHERE s.id = $1
	`
	var s entity.Schedule
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.PolyclinicID, &s.DoctorID, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.Quota,
		&s.DoctorName, &s.PolyclinicName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Jadwal", err)
		}
		return nil, fmt.Errorf("scheduleRepo.FindByID: %w", err)
	}
	return &s, nil
}

func (r *scheduleRepo) FindByPolyclinic(ctx context.Context, polyclinicID uuid.UUID) ([]*entity.Schedule, error) {
	query := `
		SELECT s.id, s.polyclinic_id, s.doctor_id, s.day_of_week, s.start_time::text, s.end_time::text, s.quota,
		       d.full_name AS doctor_name, p.name AS polyclinic_name
		FROM schedules s
		JOIN doctors d ON d.id = s.doctor_id
		JOIN polyclinics p ON p.id = s.polyclinic_id
		WHERE s.polyclinic_id = $1
		ORDER BY s.day_of_week, s.start_time
	`
	rows, err := r.pool.Query(ctx, query, polyclinicID)
	if err != nil {
		return nil, fmt.Errorf("scheduleRepo.FindByPolyclinic: %w", err)
	}
	defer rows.Close()

	var list []*entity.Schedule
	for rows.Next() {
		sc := &entity.Schedule{}
		if err := rows.Scan(
			&sc.ID, &sc.PolyclinicID, &sc.DoctorID, &sc.DayOfWeek,
			&sc.StartTime, &sc.EndTime, &sc.Quota,
			&sc.DoctorName, &sc.PolyclinicName,
		); err != nil {
			return nil, fmt.Errorf("scheduleRepo.FindByPolyclinic scan: %w", err)
		}
		list = append(list, sc)
	}
	return list, nil
}

func (r *scheduleRepo) FindByPolyclinicAndDay(ctx context.Context, polyclinicID uuid.UUID, dayOfWeek int) ([]*entity.Schedule, error) {
	// Support 1=Mon..7=Sun or 0=Sun..6=Sat mapping
	mappedDay := dayOfWeek
	if mappedDay == 0 {
		mappedDay = 7
	}

	query := `
		SELECT s.id, s.polyclinic_id, s.doctor_id, s.day_of_week, s.start_time::text, s.end_time::text, s.quota,
		       d.full_name AS doctor_name, p.name AS polyclinic_name
		FROM schedules s
		JOIN doctors d ON d.id = s.doctor_id
		JOIN polyclinics p ON p.id = s.polyclinic_id
		WHERE s.polyclinic_id = $1 AND (s.day_of_week = $2 OR s.day_of_week = $3)
		ORDER BY s.start_time
	`
	rows, err := r.pool.Query(ctx, query, polyclinicID, dayOfWeek, mappedDay)
	if err != nil {
		return nil, fmt.Errorf("scheduleRepo.FindByPolyclinicAndDay: %w", err)
	}
	defer rows.Close()

	var list []*entity.Schedule
	for rows.Next() {
		sc := &entity.Schedule{}
		if err := rows.Scan(
			&sc.ID, &sc.PolyclinicID, &sc.DoctorID, &sc.DayOfWeek,
			&sc.StartTime, &sc.EndTime, &sc.Quota,
			&sc.DoctorName, &sc.PolyclinicName,
		); err != nil {
			return nil, fmt.Errorf("scheduleRepo.FindByPolyclinicAndDay scan: %w", err)
		}
		list = append(list, sc)
	}
	return list, nil
}

func (r *scheduleRepo) GetRemainingQuota(ctx context.Context, tx pgx.Tx, scheduleID uuid.UUID, date time.Time) (int, error) {
	var quota int
	quotaQuery := `SELECT quota FROM schedules WHERE id = $1 FOR UPDATE`
	if err := tx.QueryRow(ctx, quotaQuery, scheduleID).Scan(&quota); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, apperror.NotFound("Jadwal", err)
		}
		return 0, fmt.Errorf("scheduleRepo.GetRemainingQuota (quota): %w", err)
	}

	var booked int
	countQuery := `
		SELECT COUNT(*) 
		FROM appointments a
		LEFT JOIN queues q ON q.appointment_id = a.id
		WHERE a.schedule_id = $1 AND a.visit_date = $2 AND (q.status IS NULL OR q.status != 'cancelled')
	`
	if err := tx.QueryRow(ctx, countQuery, scheduleID, date).Scan(&booked); err != nil {
		return 0, fmt.Errorf("scheduleRepo.GetRemainingQuota (count): %w", err)
	}

	remaining := quota - booked
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

func (r *scheduleRepo) Create(ctx context.Context, s *entity.Schedule) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	query := `
		INSERT INTO schedules (id, polyclinic_id, doctor_id, day_of_week, start_time, end_time, quota)
		VALUES ($1, $2, $3, $4, $5::time, $6::time, $7)
	`
	_, err := r.pool.Exec(ctx, query, s.ID, s.PolyclinicID, s.DoctorID, s.DayOfWeek, s.StartTime, s.EndTime, s.Quota)
	if err != nil {
		return fmt.Errorf("scheduleRepo.Create: %w", err)
	}
	return nil
}

func (r *scheduleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schedules WHERE id = $1`
	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("scheduleRepo.Delete: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return apperror.NotFound("Jadwal", nil)
	}
	return nil
}
