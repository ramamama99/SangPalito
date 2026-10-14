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

type queueRepo struct {
	pool *pgxpool.Pool
}

// NewQueueRepo creates a new PostgreSQL queue repository.
func NewQueueRepo(pool *pgxpool.Pool) *queueRepo {
	return &queueRepo{pool: pool}
}

// FindByID retrieves a queue entry by its UUID.
func (r *queueRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Queue, error) {
	query := `SELECT id, facility_id, appointment_id, polyclinic_id, queue_number, status,
	           call_count, called_at, service_start, service_end, queue_date, created_at, updated_at
	           FROM queues WHERE id = $1`

	q := &entity.Queue{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&q.ID, &q.FacilityID, &q.AppointmentID, &q.PolyclinicID, &q.QueueNumber,
		&q.Status, &q.CallCount, &q.CalledAt, &q.ServiceStart, &q.ServiceEnd,
		&q.QueueDate, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Antrean", err)
		}
		return nil, fmt.Errorf("queueRepo.FindByID: %w", err)
	}
	return q, nil
}

// FindActiveByPolyclinic retrieves all active (non-terminal) queues for a polyclinic on a given date.
func (r *queueRepo) FindActiveByPolyclinic(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]*entity.Queue, error) {
	query := `SELECT id, facility_id, appointment_id, polyclinic_id, queue_number, status,
	           call_count, called_at, service_start, service_end, queue_date, created_at, updated_at
	           FROM queues
	           WHERE polyclinic_id = $1 AND queue_date = $2
	           ORDER BY queue_number ASC`

	rows, err := r.pool.Query(ctx, query, polyclinicID, date)
	if err != nil {
		return nil, fmt.Errorf("queueRepo.FindActiveByPolyclinic: %w", err)
	}
	defer rows.Close()

	var queues []*entity.Queue
	for rows.Next() {
		q := &entity.Queue{}
		if err := rows.Scan(
			&q.ID, &q.FacilityID, &q.AppointmentID, &q.PolyclinicID, &q.QueueNumber,
			&q.Status, &q.CallCount, &q.CalledAt, &q.ServiceStart, &q.ServiceEnd,
			&q.QueueDate, &q.CreatedAt, &q.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("queueRepo.FindActiveByPolyclinic scan: %w", err)
		}
		queues = append(queues, q)
	}
	return queues, nil
}

// Create inserts a new queue entry.
func (r *queueRepo) Create(ctx context.Context, q *entity.Queue) error {
	return r.createWith(ctx, r.pool, q)
}

// CreateTx inserts a new queue entry within a transaction.
func (r *queueRepo) CreateTx(ctx context.Context, tx pgx.Tx, q *entity.Queue) error {
	return r.createWith(ctx, tx, q)
}

func (r *queueRepo) createWith(ctx context.Context, exec interface{}, q *entity.Queue) error {
	query := `INSERT INTO queues (id, facility_id, appointment_id, polyclinic_id, queue_number,
	           status, queue_date, created_at, updated_at)
	           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	now := time.Now()
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	q.CreatedAt = now
	q.UpdatedAt = now

	args := []interface{}{
		q.ID, q.FacilityID, q.AppointmentID, q.PolyclinicID,
		q.QueueNumber, q.Status, q.QueueDate, q.CreatedAt, q.UpdatedAt,
	}

	var err error
	switch e := exec.(type) {
	case *pgxpool.Pool:
		_, err = e.Exec(ctx, query, args...)
	case pgx.Tx:
		_, err = e.Exec(ctx, query, args...)
	}
	if err != nil {
		return fmt.Errorf("queueRepo.Create: %w", err)
	}
	return nil
}

// UpdateStatus updates the status and related fields of a queue entry.
func (r *queueRepo) UpdateStatus(ctx context.Context, q *entity.Queue) error {
	query := `UPDATE queues SET status = $1, call_count = $2, called_at = $3,
	           service_start = $4, service_end = $5, updated_at = $6
	           WHERE id = $7`

	q.UpdatedAt = time.Now()
	_, err := r.pool.Exec(ctx, query,
		q.Status, q.CallCount, q.CalledAt, q.ServiceStart, q.ServiceEnd, q.UpdatedAt, q.ID,
	)
	if err != nil {
		return fmt.Errorf("queueRepo.UpdateStatus: %w", err)
	}
	return nil
}

// GetNextQueueNumber returns the next queue number for a polyclinic on a date using SELECT FOR UPDATE.
func (r *queueRepo) GetNextQueueNumber(ctx context.Context, tx pgx.Tx, polyclinicID uuid.UUID, date time.Time) (int, error) {
	query := `SELECT COALESCE(MAX(queue_number), 0) FROM queues
	           WHERE polyclinic_id = $1 AND queue_date = $2 FOR UPDATE`

	var maxNum int
	err := tx.QueryRow(ctx, query, polyclinicID, date).Scan(&maxNum)
	if err != nil {
		return 0, fmt.Errorf("queueRepo.GetNextQueueNumber: %w", err)
	}
	return maxNum + 1, nil
}

// CountByPolyclinicAndDate counts queues by polyclinic, date, and optionally status.
func (r *queueRepo) CountByPolyclinicAndDate(ctx context.Context, polyclinicID uuid.UUID, date time.Time, status entity.QueueStatus) (int, error) {
	var count int
	var err error

	if status == "" {
		query := `SELECT COUNT(*) FROM queues WHERE polyclinic_id = $1 AND queue_date = $2`
		err = r.pool.QueryRow(ctx, query, polyclinicID, date).Scan(&count)
	} else {
		query := `SELECT COUNT(*) FROM queues WHERE polyclinic_id = $1 AND queue_date = $2 AND status = $3`
		err = r.pool.QueryRow(ctx, query, polyclinicID, date, status).Scan(&count)
	}

	if err != nil {
		return 0, fmt.Errorf("queueRepo.CountByPolyclinicAndDate: %w", err)
	}
	return count, nil
}
