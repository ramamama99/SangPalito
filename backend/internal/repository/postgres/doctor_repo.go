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

const doctorSelectQuery = `
	SELECT d.id, d.facility_id, d.user_id, d.full_name, d.specialty, d.str_no, d.is_active,
	       COALESCE(d.email, u.email, '') AS email,
	       COALESCE(d.phone, u.phone, '') AS phone,
	       COALESCE(d.polyclinic_id, u.polyclinic_id) AS polyclinic_id,
	       COALESCE(p.name, up.name, '') AS polyclinic_name,
	       COALESCE(u.is_first_login, false) AS is_first_login
	FROM doctors d
	LEFT JOIN users u ON u.id = d.user_id
	LEFT JOIN polyclinics p ON p.id = d.polyclinic_id
	LEFT JOIN polyclinics up ON up.id = u.polyclinic_id
`

func (r *doctorRepo) FindAll(ctx context.Context, facilityID uuid.UUID) ([]entity.Doctor, error) {
	query := doctorSelectQuery + `
		WHERE d.facility_id = $1 AND d.is_active = true
		ORDER BY d.full_name ASC
	`
	rows, err := r.pool.Query(ctx, query, facilityID)
	if err != nil {
		return nil, fmt.Errorf("doctorRepo.FindAll: %w", err)
	}
	defer rows.Close()

	var list []entity.Doctor
	for rows.Next() {
		var d entity.Doctor
		if err := rows.Scan(
			&d.ID, &d.FacilityID, &d.UserID, &d.FullName, &d.Specialty, &d.STRNo, &d.IsActive,
			&d.Email, &d.Phone, &d.PolyclinicID, &d.PolyclinicName, &d.IsFirstLogin,
		); err != nil {
			return nil, fmt.Errorf("doctorRepo.FindAll scan: %w", err)
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *doctorRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Doctor, error) {
	query := doctorSelectQuery + ` WHERE d.id = $1`
	var d entity.Doctor
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.FacilityID, &d.UserID, &d.FullName, &d.Specialty, &d.STRNo, &d.IsActive,
		&d.Email, &d.Phone, &d.PolyclinicID, &d.PolyclinicName, &d.IsFirstLogin,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("Dokter", err)
		}
		return nil, fmt.Errorf("doctorRepo.FindByID: %w", err)
	}
	return &d, nil
}

func (r *doctorRepo) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.Doctor, error) {
	query := doctorSelectQuery + ` WHERE d.user_id = $1`
	var d entity.Doctor
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&d.ID, &d.FacilityID, &d.UserID, &d.FullName, &d.Specialty, &d.STRNo, &d.IsActive,
		&d.Email, &d.Phone, &d.PolyclinicID, &d.PolyclinicName, &d.IsFirstLogin,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found is ok
		}
		return nil, fmt.Errorf("doctorRepo.FindByUserID: %w", err)
	}
	return &d, nil
}

func (r *doctorRepo) Create(ctx context.Context, d *entity.Doctor) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	query := `
		INSERT INTO doctors (id, facility_id, user_id, full_name, specialty, str_no, is_active, polyclinic_id, phone, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query, d.ID, d.FacilityID, d.UserID, d.FullName, d.Specialty, d.STRNo, d.IsActive, d.PolyclinicID, d.Phone, d.Email)
	if err != nil {
		return fmt.Errorf("doctorRepo.Create: %w", err)
	}
	return nil
}

func (r *doctorRepo) Update(ctx context.Context, d *entity.Doctor) error {
	query := `
		UPDATE doctors
		SET full_name = $1, specialty = $2, str_no = $3, is_active = $4,
		    polyclinic_id = $5, phone = $6, email = $7
		WHERE id = $8
	`
	_, err := r.pool.Exec(ctx, query, d.FullName, d.Specialty, d.STRNo, d.IsActive, d.PolyclinicID, d.Phone, d.Email, d.ID)
	if err != nil {
		return fmt.Errorf("doctorRepo.Update: %w", err)
	}
	return nil
}

func (r *doctorRepo) Delete(ctx context.Context, id uuid.UUID) error {
	// First unassign any schedules
	_, _ = r.pool.Exec(ctx, `DELETE FROM schedules WHERE doctor_id = $1`, id)
	cmd, err := r.pool.Exec(ctx, `DELETE FROM doctors WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("doctorRepo.Delete: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return apperror.NotFound("Doctor", nil)
	}
	return nil
}