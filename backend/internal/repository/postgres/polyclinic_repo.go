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
)

type polyclinicRepo struct {
	pool *pgxpool.Pool
}

// NewPolyclinicRepo creates a new PostgreSQL polyclinic repository.
func NewPolyclinicRepo(pool *pgxpool.Pool) *polyclinicRepo {
	return &polyclinicRepo{pool: pool}
}

// FindByID retrieves a polyclinic by its UUID.
func (r *polyclinicRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Polyclinic, error) {
	query := `SELECT id, facility_id, name, code, is_active, created_at, updated_at
	           FROM polyclinics WHERE id = $1`

	p := &entity.Polyclinic{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.FacilityID, &p.Name, &p.Code, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Poliklinik", err)
		}
		return nil, fmt.Errorf("polyclinicRepo.FindByID: %w", err)
	}
	return p, nil
}

// FindByFacilityID retrieves all polyclinics for a facility.
func (r *polyclinicRepo) FindByFacilityID(ctx context.Context, facilityID uuid.UUID) ([]*entity.Polyclinic, error) {
	query := `SELECT id, facility_id, name, code, is_active, created_at, updated_at
	           FROM polyclinics WHERE facility_id = $1 ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query, facilityID)
	if err != nil {
		return nil, fmt.Errorf("polyclinicRepo.FindByFacilityID: %w", err)
	}
	defer rows.Close()

	var polyclinics []*entity.Polyclinic
	for rows.Next() {
		p := &entity.Polyclinic{}
		if err := rows.Scan(&p.ID, &p.FacilityID, &p.Name, &p.Code, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("polyclinicRepo.FindByFacilityID scan: %w", err)
		}
		polyclinics = append(polyclinics, p)
	}
	return polyclinics, nil
}

// FindActive retrieves all active polyclinics for a facility.
func (r *polyclinicRepo) FindActive(ctx context.Context, facilityID uuid.UUID) ([]*entity.Polyclinic, error) {
	query := `SELECT id, facility_id, name, code, is_active, created_at, updated_at
	           FROM polyclinics WHERE facility_id = $1 AND is_active = true ORDER BY name ASC`

	rows, err := r.pool.Query(ctx, query, facilityID)
	if err != nil {
		return nil, fmt.Errorf("polyclinicRepo.FindActive: %w", err)
	}
	defer rows.Close()

	var polyclinics []*entity.Polyclinic
	for rows.Next() {
		p := &entity.Polyclinic{}
		if err := rows.Scan(&p.ID, &p.FacilityID, &p.Name, &p.Code, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("polyclinicRepo.FindActive scan: %w", err)
		}
		polyclinics = append(polyclinics, p)
	}
	return polyclinics, nil
}
