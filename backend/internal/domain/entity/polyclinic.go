package entity

import (
	"time"

	"github.com/google/uuid"
)

// Polyclinic represents a polyclinic unit within a facility.
type Polyclinic struct {
	ID            uuid.UUID `json:"id"`
	FacilityID    uuid.UUID `json:"facility_id"`
	Name          string    `json:"name"`            // contoh: Poli Umum, Poli Gigi
	Code          string    `json:"code,omitempty"`  // alias
	QueuePrefix   string    `json:"queue_prefix"`    // contoh: A, B, C
	AvgServiceMin int       `json:"avg_service_min"` // estimasi menit layanan (default: 10)
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}
