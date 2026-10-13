package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type scheduleRepo struct {
	pool *pgxpool.Pool
}

func NewScheduleRepo(pool *pgxpool.Pool) repository.ScheduleRepository {
	return &scheduleRepo{pool: pool}
}

func (r *scheduleRepo) FindWithRemainingSlot(ctx context.Context, polyclinicID uuid.UUID, date time.Time, dayOfWeek int) ([]entity.Schedule, error) {
	// Menghitung sisa kuota hari itu: Quota dikurangi jumlah antrean yang belum dibatalkan
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
		SELECT id, polyclinic_id, doctor_id, day_of_week, start_time::text, end_time::text, quota
		FROM schedules
		WHERE id = $1
	`
	var s entity.Schedule
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.PolyclinicID, &s.DoctorID, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.Quota,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Jadwal", err)
		}
		return nil, fmt.Errorf("scheduleRepo.FindByID: %w", err)
	}
	return &s, nil
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
