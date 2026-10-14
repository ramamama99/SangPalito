package entity

import (
	"time"

	"github.com/google/uuid"
)

// Doctor represents a doctor in the system.
type Doctor struct {
	ID             uuid.UUID
	FacilityID     uuid.UUID
	UserID         *uuid.UUID
	FullName       string
	Specialization string
	SIPNo          string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
