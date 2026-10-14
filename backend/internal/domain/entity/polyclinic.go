package entity

import (
	"time"

	"github.com/google/uuid"
)

// Polyclinic represents a clinic/polyclinic within a facility.
type Polyclinic struct {
	ID         uuid.UUID
	FacilityID uuid.UUID
	Name       string
	Code       string
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
