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

type userRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo creates a new PostgreSQL user repository.
func NewUserRepo(pool *pgxpool.Pool) *userRepo {
	return &userRepo{pool: pool}
}

const userSelectQuery = `
	SELECT u.id, u.facility_id, u.email, u.password_hash, u.full_name,
	       COALESCE(r.name, 'registrar') AS role,
	       u.is_active, u.last_login_at, u.refresh_token_hash, u.created_at, u.created_at AS updated_at,
	       u.phone, u.polyclinic_id, COALESCE(p.name, '') AS polyclinic_name
	FROM users u
	LEFT JOIN user_roles ur ON ur.user_id = u.id
	LEFT JOIN roles r ON r.id = ur.role_id
	LEFT JOIN polyclinics p ON p.id = u.polyclinic_id
`

// FindByID retrieves a user by their UUID.
func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := userSelectQuery + ` WHERE u.id = $1`

	u := &entity.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.FacilityID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role,
		&u.IsActive, &u.LastLoginAt, &u.RefreshTokenHash, &u.CreatedAt, &u.UpdatedAt,
		&u.Phone, &u.PolyclinicID, &u.PolyclinicName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("User", err)
		}
		return nil, fmt.Errorf("userRepo.FindByID: %w", err)
	}
	return u, nil
}

// FindByEmail retrieves a user by email within a facility.
func (r *userRepo) FindByEmail(ctx context.Context, email string, facilityID uuid.UUID) (*entity.User, error) {
	query := userSelectQuery + ` WHERE u.email = $1 AND u.facility_id = $2`

	u := &entity.User{}
	err := r.pool.QueryRow(ctx, query, email, facilityID).Scan(
		&u.ID, &u.FacilityID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role,
		&u.IsActive, &u.LastLoginAt, &u.RefreshTokenHash, &u.CreatedAt, &u.UpdatedAt,
		&u.Phone, &u.PolyclinicID, &u.PolyclinicName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("User", err)
		}
		return nil, fmt.Errorf("userRepo.FindByEmail: %w", err)
	}
	return u, nil
}

// FindByFacilityID retrieves paginated users belonging to a facility.
func (r *userRepo) FindByFacilityID(ctx context.Context, facilityID uuid.UUID, page, limit int) ([]*entity.User, int, error) {
	countQuery := `SELECT COUNT(*) FROM users WHERE facility_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, facilityID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("userRepo.FindByFacilityID count: %w", err)
	}

	offset := (page - 1) * limit
	query := userSelectQuery + ` WHERE u.facility_id = $1 ORDER BY u.created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, facilityID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("userRepo.FindByFacilityID: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		u := &entity.User{}
		if err := rows.Scan(
			&u.ID, &u.FacilityID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role,
			&u.IsActive, &u.LastLoginAt, &u.RefreshTokenHash, &u.CreatedAt, &u.UpdatedAt,
			&u.Phone, &u.PolyclinicID, &u.PolyclinicName,
		); err != nil {
			return nil, 0, fmt.Errorf("userRepo.FindByFacilityID scan: %w", err)
		}
		users = append(users, u)
	}
	return users, total, nil
}

// Create inserts a new user into the database.
func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (id, facility_id, email, password_hash, full_name, phone, polyclinic_id, is_active, created_at)
	           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	now := time.Now()
	user.ID = uuid.New()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.pool.Exec(ctx, query,
		user.ID, user.FacilityID, user.Email, user.PasswordHash, user.FullName,
		user.Phone, user.PolyclinicID, user.IsActive, user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("userRepo.Create: %w", err)
	}

	if user.Role != "" {
		roleQuery := `INSERT INTO user_roles (user_id, role_id)
		              SELECT $1, id FROM roles WHERE name = $2
		              ON CONFLICT DO NOTHING`
		_, _ = r.pool.Exec(ctx, roleQuery, user.ID, user.Role)
	}

	return nil
}

// Update modifies an existing user's data.
func (r *userRepo) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET full_name = $1, phone = $2, polyclinic_id = $3, is_active = $4 WHERE id = $5`
	user.UpdatedAt = time.Now()

	_, err := r.pool.Exec(ctx, query, user.FullName, user.Phone, user.PolyclinicID, user.IsActive, user.ID)
	if err != nil {
		return fmt.Errorf("userRepo.Update: %w", err)
	}

	if user.Role != "" {
		_, _ = r.pool.Exec(ctx, `DELETE FROM user_roles WHERE user_id = $1`, user.ID)
		roleQuery := `INSERT INTO user_roles (user_id, role_id)
		              SELECT $1, id FROM roles WHERE name = $2
		              ON CONFLICT DO NOTHING`
		_, _ = r.pool.Exec(ctx, roleQuery, user.ID, user.Role)
	}

	return nil
}

// Delete removes a user and their role assignments.
func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, _ = r.pool.Exec(ctx, `DELETE FROM user_roles WHERE user_id = $1`, id)
	cmd, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("userRepo.Delete: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return apperror.NotFound("User", nil)
	}
	return nil
}

// FindByRefreshTokenHash retrieves an active user matching the given refresh token SHA-256 hash.
func (r *userRepo) FindByRefreshTokenHash(ctx context.Context, tokenHash string) (*entity.User, error) {
	query := userSelectQuery + ` WHERE u.refresh_token_hash = $1 AND u.is_active = true`

	u := &entity.User{}
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&u.ID, &u.FacilityID, &u.Email, &u.PasswordHash, &u.FullName, &u.Role,
		&u.IsActive, &u.LastLoginAt, &u.RefreshTokenHash, &u.CreatedAt, &u.UpdatedAt,
		&u.Phone, &u.PolyclinicID, &u.PolyclinicName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NotFound("User", err)
		}
		return nil, fmt.Errorf("userRepo.FindByRefreshTokenHash: %w", err)
	}
	return u, nil
}

// UpdateRefreshToken updates the refresh token hash for a user (nil clears it on logout).
func (r *userRepo) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash *string) error {
	query := `UPDATE users SET refresh_token_hash = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, tokenHash, userID)
	if err != nil {
		return fmt.Errorf("userRepo.UpdateRefreshToken: %w", err)
	}
	return nil
}

// UpdateLastLogin sets last_login_at to now.
func (r *userRepo) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET last_login_at = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("userRepo.UpdateLastLogin: %w", err)
	}
	return nil
}

// UpdatePassword updates a user's password hash.
func (r *userRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("userRepo.UpdatePassword: %w", err)
	}
	return nil
}
