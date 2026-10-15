package postgres

import (
	"context"
	"encoding/json"
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

const queueBaseCols = `id, appointment_id, polyclinic_id, queue_number, queue_code, status,
	checked_in_at, call_count, called_at, service_start, service_end, queue_date`

func scanQueue(row pgx.Row) (*entity.Queue, error) {
	q := &entity.Queue{}
	err := row.Scan(
		&q.ID, &q.AppointmentID, &q.PolyclinicID, &q.QueueNumber, &q.QueueCode,
		&q.Status, &q.CheckedInAt, &q.CallCount, &q.CalledAt, &q.ServiceStart, &q.ServiceEnd,
		&q.QueueDate,
	)
	return q, err
}

// FindByID retrieves a queue entry by its UUID.
func (r *queueRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Queue, error) {
	query := fmt.Sprintf("SELECT %s FROM queues WHERE id = $1", queueBaseCols)
	q, err := scanQueue(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Antrean", err)
		}
		return nil, fmt.Errorf("queueRepo.FindByID: %w", err)
	}
	return q, nil
}

// FindActiveByPolyclinic retrieves all active queues for a polyclinic on a given date.
func (r *queueRepo) FindActiveByPolyclinic(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]*entity.Queue, error) {
	query := fmt.Sprintf("SELECT %s FROM queues WHERE polyclinic_id = $1 AND queue_date = $2 ORDER BY queue_number ASC", queueBaseCols)
	rows, err := r.pool.Query(ctx, query, polyclinicID, date)
	if err != nil {
		return nil, fmt.Errorf("queueRepo.FindActiveByPolyclinic: %w", err)
	}
	defer rows.Close()

	var queues []*entity.Queue
	for rows.Next() {
		q := &entity.Queue{}
		if err := rows.Scan(
			&q.ID, &q.AppointmentID, &q.PolyclinicID, &q.QueueNumber, &q.QueueCode,
			&q.Status, &q.CheckedInAt, &q.CallCount, &q.CalledAt, &q.ServiceStart, &q.ServiceEnd,
			&q.QueueDate,
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
	query := `INSERT INTO queues (id, appointment_id, polyclinic_id, queue_number, queue_code,
	           status, queue_date)
	           VALUES ($1, $2, $3, $4, $5, $6, $7)`

	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	if q.QueueCode == "" {
		q.QueueCode = fmt.Sprintf("A-%03d", q.QueueNumber)
	}

	args := []interface{}{
		q.ID, q.AppointmentID, q.PolyclinicID,
		q.QueueNumber, q.QueueCode, q.Status, q.QueueDate,
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

// UpdateStatus updates the status and service timing of a queue entry.
func (r *queueRepo) UpdateStatus(ctx context.Context, q *entity.Queue) error {
	query := `UPDATE queues SET status = $1, call_count = $2, called_at = $3,
	           checked_in_at = $4, service_start = $5, service_end = $6
	           WHERE id = $7`

	_, err := r.pool.Exec(ctx, query,
		q.Status, q.CallCount, q.CalledAt, q.CheckedInAt, q.ServiceStart, q.ServiceEnd, q.ID,
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

// FindCompletedDurations retrieves actual consultation durations (in minutes) for completed visits today.
func (r *queueRepo) FindCompletedDurations(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]float64, error) {
	query := `SELECT EXTRACT(EPOCH FROM (service_end - service_start)) / 60.0 AS duration_min
	           FROM queues
	           WHERE polyclinic_id = $1 AND queue_date = $2 AND status = 'done'
	             AND service_start IS NOT NULL AND service_end IS NOT NULL AND service_end > service_start
	           ORDER BY service_end ASC`

	rows, err := r.pool.Query(ctx, query, polyclinicID, date)
	if err != nil {
		return nil, fmt.Errorf("queueRepo.FindCompletedDurations: %w", err)
	}
	defer rows.Close()

	var durations []float64
	for rows.Next() {
		var dur float64
		if err := rows.Scan(&dur); err != nil {
			return nil, fmt.Errorf("queueRepo.FindCompletedDurations scan: %w", err)
		}
		durations = append(durations, dur)
	}
	return durations, nil
}

const queueWithPatientQuery = `SELECT q.id, q.appointment_id, q.polyclinic_id, q.queue_number, q.queue_code, q.status,
	q.checked_in_at, q.call_count, q.called_at, q.service_start, q.service_end, q.queue_date,
	COALESCE(EXTRACT(YEAR FROM age(CURRENT_DATE, p.birth_date)), 30)::int as patient_age,
	COALESCE(a.priority::text, 'normal') as priority,
	COALESCE(pa.chief_complaint, '') as chief_complaint,
	COALESCE(pa.triage_flag, '') as triage_flag,
	COALESCE(pa.symptoms, '[]'::jsonb) as symptoms
	FROM queues q
	JOIN appointments a ON q.appointment_id = a.id
	JOIN patients p ON a.patient_id = p.id
	LEFT JOIN pre_anamneses pa ON pa.appointment_id = a.id`

func scanQueueWithPatient(row pgx.Row) (*entity.QueueWithPatient, error) {
	q := &entity.Queue{}
	qp := &entity.QueueWithPatient{Queue: q}
	var symptomsJSON []byte

	err := row.Scan(
		&q.ID, &q.AppointmentID, &q.PolyclinicID, &q.QueueNumber, &q.QueueCode, &q.Status,
		&q.CheckedInAt, &q.CallCount, &q.CalledAt, &q.ServiceStart, &q.ServiceEnd, &q.QueueDate,
		&qp.PatientAge, &qp.Priority, &qp.ChiefComplaint, &qp.TriageFlag, &symptomsJSON,
	)
	if err != nil {
		return nil, err
	}

	if len(symptomsJSON) > 0 {
		var symptoms []string
		_ = json.Unmarshal(symptomsJSON, &symptoms)
		qp.Symptoms = symptoms
	}

	return qp, nil
}

// FindActiveQueueDetails retrieves active queues with full patient clinical context.
func (r *queueRepo) FindActiveQueueDetails(ctx context.Context, polyclinicID uuid.UUID, date time.Time) ([]*entity.QueueWithPatient, error) {
	query := fmt.Sprintf("%s WHERE q.polyclinic_id = $1 AND q.queue_date = $2 ORDER BY q.queue_number ASC", queueWithPatientQuery)

	rows, err := r.pool.Query(ctx, query, polyclinicID, date)
	if err != nil {
		return nil, fmt.Errorf("queueRepo.FindActiveQueueDetails: %w", err)
	}
	defer rows.Close()

	var results []*entity.QueueWithPatient
	for rows.Next() {
		qp, err := scanQueueWithPatient(rows)
		if err != nil {
			return nil, fmt.Errorf("queueRepo.FindActiveQueueDetails scan: %w", err)
		}
		results = append(results, qp)
	}
	return results, nil
}

// FindQueueWithPatientByID retrieves a single queue entry with patient clinical context.
func (r *queueRepo) FindQueueWithPatientByID(ctx context.Context, queueID uuid.UUID) (*entity.QueueWithPatient, error) {
	query := fmt.Sprintf("%s WHERE q.id = $1", queueWithPatientQuery)
	qp, err := scanQueueWithPatient(r.pool.QueryRow(ctx, query, queueID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Antrean", err)
		}
		return nil, fmt.Errorf("queueRepo.FindQueueWithPatientByID: %w", err)
	}
	return qp, nil
}
