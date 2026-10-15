package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"sangpalito/backend/internal/domain/entity"
)

type medicationRepo struct {
	pool *pgxpool.Pool
}

func NewMedicationRepo(pool *pgxpool.Pool) *medicationRepo {
	return &medicationRepo{pool: pool}
}

func (r *medicationRepo) FindAll(ctx context.Context, search, category, route string) ([]*entity.Medication, error) {
	query := `
		SELECT id, code, name, generic_name, category, dosage_form, route,
		       default_frequency, default_timing, indication,
		       unit, is_prescription_required, created_at, updated_at
		FROM medications
		WHERE 1=1
	`
	args := []any{}
	argIdx := 1

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR generic_name ILIKE $%d OR indication ILIKE $%d)", argIdx, argIdx, argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}

	if route != "" {
		query += fmt.Sprintf(" AND route = $%d", argIdx)
		args = append(args, route)
		argIdx++
	}

	query += " ORDER BY name ASC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("medicationRepo.FindAll: %w", err)
	}
	defer rows.Close()

	meds := make([]*entity.Medication, 0)
	for rows.Next() {
		var m entity.Medication
		if err := rows.Scan(
			&m.ID, &m.Code, &m.Name, &m.GenericName, &m.Category, &m.DosageForm, &m.Route,
			&m.DefaultFrequency, &m.DefaultTiming, &m.Indication,
			&m.Unit, &m.IsPrescriptionRequired, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("medicationRepo.FindAll scan: %w", err)
		}
		meds = append(meds, &m)
	}

	return meds, nil
}

func (r *medicationRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Medication, error) {
	query := `
		SELECT id, code, name, generic_name, category, dosage_form, route,
		       default_frequency, default_timing, indication,
		       unit, is_prescription_required, created_at, updated_at
		FROM medications
		WHERE id = $1
	`
	var m entity.Medication
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.Code, &m.Name, &m.GenericName, &m.Category, &m.DosageForm, &m.Route,
		&m.DefaultFrequency, &m.DefaultTiming, &m.UsageInstruction, &m.Indication,
		&m.Stock, &m.Unit, &m.IsPrescriptionRequired, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("medicationRepo.FindByID: %w", err)
	}
	return &m, nil
}

func (r *medicationRepo) Create(ctx context.Context, m *entity.Medication) error {
	query := `
		INSERT INTO medications (
			id, code, name, generic_name, category, dosage_form, route,
			default_frequency, default_timing, usage_instruction, indication,
			stock, unit, is_prescription_required, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11,
			$12, $13, $14, $15, $16
		)
	`
	now := time.Now()
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	m.CreatedAt = now
	m.UpdatedAt = now

	_, err := r.pool.Exec(ctx, query,
		m.ID, m.Code, m.Name, m.GenericName, m.Category, m.DosageForm, m.Route,
		m.DefaultFrequency, m.DefaultTiming, m.UsageInstruction, m.Indication,
		m.Stock, m.Unit, m.IsPrescriptionRequired, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("medicationRepo.Create: %w", err)
	}
	return nil
}

func (r *medicationRepo) Update(ctx context.Context, m *entity.Medication) error {
	query := `
		UPDATE medications SET
			code = $1, name = $2, generic_name = $3, category = $4,
			dosage_form = $5, route = $6, default_frequency = $7,
			default_timing = $8, usage_instruction = $9, indication = $10,
			stock = $11, unit = $12, is_prescription_required = $13,
			updated_at = $14
		WHERE id = $15
	`
	m.UpdatedAt = time.Now()
	_, err := r.pool.Exec(ctx, query,
		m.Code, m.Name, m.GenericName, m.Category,
		m.DosageForm, m.Route, m.DefaultFrequency,
		m.DefaultTiming, m.UsageInstruction, m.Indication,
		m.Stock, m.Unit, m.IsPrescriptionRequired,
		m.UpdatedAt, m.ID,
	)
	if err != nil {
		return fmt.Errorf("medicationRepo.Update: %w", err)
	}
	return nil
}

func (r *medicationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM medications WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("medicationRepo.Delete: %w", err)
	}
	return nil
}
