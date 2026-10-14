package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents a system user (staff member).
type User struct {
	ID               uuid.UUID
	FacilityID       uuid.UUID
	Email            string
	PasswordHash     string
	FullName         string
	Role             string // "admin", "registrar", "medic", "head"
	IsActive         bool
	LastLoginAt      *time.Time
	RefreshTokenHash *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
