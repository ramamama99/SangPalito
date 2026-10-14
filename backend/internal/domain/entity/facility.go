package entity

import (
	"time"

	"github.com/google/uuid"
)

// Facility represents a healthcare facility (Puskesmas / Klinik).
type Facility struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // puskesmas | klinik
	Code      *string   `json:"code,omitempty"`
	Address   *string   `json:"address,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
