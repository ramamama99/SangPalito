package entity

import (
	"time"

	"github.com/google/uuid"
)

// Facility represents a healthcare facility.
type Facility struct {
	ID        uuid.UUID
	Name      string
	Code      string
	Address   string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
