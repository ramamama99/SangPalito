package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/domain/repository"
)

type polyclinicRepo struct {
	pool *pgxpool.Pool
}

func NewPolyclinicRepo(pool *pgxpool.Pool) repository.PolyclinicRepository {
	return &polyclinicRepo{pool: pool}
}

func (r *polyclinicRepo) FindAll(ctx context.Context, facilityID uuid.UUID) ([]entity.Polyclinic, error) {
	query := `
		SELECT id, facility_id, name, queue_prefix, avg_service_min, is_active
		FROM polyclinics
		WHERE facility_id = $1 AND is_active = true
		ORDER BY name ASC
	`
	rows, err := r.pool.Query(ctx, query, facilityID)
	if err != nil {
		return nil, fmt.Errorf("polyclinicRepo.FindAll: %w", err)
	}
	defer rows.Close()

	var list []entity.Polyclinic
	for rows.Next() {
		var p entity.Polyclinic
		if err := rows.Scan(&p.ID, &p.FacilityID, &p.Name, &p.QueuePrefix, &p.AvgServiceMin, &p.IsActive); err != nil {
			return nil, fmt.Errorf("polyclinicRepo.FindAll scan: %w", err)
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *polyclinicRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Polyclinic, error) {
	query := `
		SELECT id, facility_id, name, queue_prefix, avg_service_min, is_active
		FROM polyclinics
		WHERE id = $1
	`
	var p entity.Polyclinic
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.FacilityID, &p.Name, &p.QueuePrefix, &p.AvgServiceMin, &p.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Poliklinik", err)
		}
		return nil, fmt.Errorf("polyclinicRepo.FindByID: %w", err)
	}
	return &p, nil
}

func (r *polyclinicRepo) Create(ctx context.Context, p *entity.Polyclinic) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	query := `
		INSERT INTO polyclinics (id, facility_id, name, queue_prefix, avg_service_min, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query, p.ID, p.FacilityID, p.Name, p.QueuePrefix, p.AvgServiceMin, p.IsActive)
	if err != nil {
		return fmt.Errorf("polyclinicRepo.Create: %w", err)
	}
	return nil
}

func (r *polyclinicRepo) Update(ctx context.Context, p *entity.Polyclinic) error {
	query := `
		UPDATE polyclinics
		SET name = $2, queue_prefix = $3, avg_service_min = $4, is_active = $5
		WHERE id = $1
	`
	cmd, err := r.pool.Exec(ctx, query, p.ID, p.Name, p.QueuePrefix, p.AvgServiceMin, p.IsActive)
	if err != nil {
		return fmt.Errorf("polyclinicRepo.Update: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return apperror.NotFound("Poliklinik", nil)
	}
	return nil
}