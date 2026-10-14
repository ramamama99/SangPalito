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

type facilityRepo struct {
	pool *pgxpool.Pool
}

func NewFacilityRepo(pool *pgxpool.Pool) repository.FacilityRepository {
	return &facilityRepo{pool: pool}
}

func (r *facilityRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Facility, error) {
	query := `
		SELECT id, name, type, code, address, phone, created_at
		FROM facilities
		WHERE id = $1
	`
	var f entity.Facility
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&f.ID, &f.Name, &f.Type, &f.Code, &f.Address, &f.Phone, &f.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Fasilitas", err)
		}
		return nil, fmt.Errorf("facilityRepo.FindByID: %w", err)
	}
	return &f, nil
}

func (r *facilityRepo) GetFirst(ctx context.Context) (*entity.Facility, error) {
	query := `
		SELECT id, name, type, code, address, phone, created_at
		FROM facilities
		ORDER BY created_at ASC
		LIMIT 1
	`
	var f entity.Facility
	err := r.pool.QueryRow(ctx, query).Scan(
		&f.ID, &f.Name, &f.Type, &f.Code, &f.Address, &f.Phone, &f.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Fasilitas default", err)
		}
		return nil, fmt.Errorf("facilityRepo.GetFirst: %w", err)
	}
	return &f, nil
}

func (r *facilityRepo) Create(ctx context.Context, f *entity.Facility) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	query := `
		INSERT INTO facilities (id, name, type, code, address, phone)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query, f.ID, f.Name, f.Type, f.Code, f.Address, f.Phone)
	if err != nil {
		return fmt.Errorf("facilityRepo.Create: %w", err)
	}
	return nil
}