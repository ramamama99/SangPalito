package repository

import (
	"context"

	"github.com/google/uuid"

	"sangpalito/backend/internal/domain/entity"
)

// UserRepository defines the contract for user data access.
type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string, facilityID uuid.UUID) (*entity.User, error)
	FindByFacilityID(ctx context.Context, facilityID uuid.UUID, page, limit int) ([]*entity.User, int, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	FindByRefreshTokenHash(ctx context.Context, tokenHash string) (*entity.User, error)
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash *string) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
}
