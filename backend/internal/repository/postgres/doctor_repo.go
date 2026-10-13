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

type doctorRepo struct {
	pool *pgxpool.Pool
}

func NewDoctorRepo(pool *pgxpool.Pool) repository.DoctorRepository {
	return &doctorRepo{pool: pool}
}

func (r *doctorRepo) FindAll(ctx context.Context, facilityID uuid.UUID) ([]entity.Doctor, error) {
	query := `
		SELECT id, facility_id, user_id, full_name, specialty, str_no, is_active
		FROM doctors
		WHERE facility_id = $1 AND is_active = true
		ORDER BY full_name ASC
	`
	rows, err := r.pool.Query(ctx, query, facilityID)
	if err != nil {
		return nil, fmt.Errorf("doctorRepo.FindAll: %w", err)
	}
	defer rows.Close()

	var list []entity.Doctor
	for rows.Next() {
		var d entity.Doctor
		if err := rows.Scan(&d.ID, &d.FacilityID, &d.UserID, &d.FullName, &d.Specialty, &d.STRNo, &d.IsActive); err != nil {
			return nil, fmt.Errorf("doctorRepo.FindAll scan: %w", err)
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *doctorRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Doctor, error) {
	query := `
		SELECT id, facility_id, user_id, full_name, specialty, str_no, is_active
		FROM doctors
		WHERE id = $1
	`
	var d entity.Doctor
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.FacilityID, &d.UserID, &d.FullName, &d.Specialty, &d.STRNo, &d.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Dokter", err)
		}
		return nil, fmt.Errorf("doctorRepo.FindByID: %w", err)
	}
	return &d, nil
}

func (r *doctorRepo) Create(ctx context.Context, d *entity.Doctor) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	query := `
		INSERT INTO doctors (id, facility_id, user_id, full_name, specialty, str_no, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(ctx, query, d.ID, d.FacilityID, d.UserID, d.FullName, d.Specialty, d.STRNo, d.IsActive)
	if err != nil {
		return fmt.Errorf("doctorRepo.Create: %w", err)
	}
	return nil
}