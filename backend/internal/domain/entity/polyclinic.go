package entity

import (
	"github.com/google/uuid"
)

type Polyclinic struct {
	ID            uuid.UUID `json:"id"`
	FacilityID    uuid.UUID `json:"facility_id"`
	Name          string    `json:"name"`            // contoh: Poli Umum, Poli Gigi
	QueuePrefix   string    `json:"queue_prefix"`    // contoh: A, B, C
	AvgServiceMin int       `json:"avg_service_min"` // estimasi menit layanan (default: 10)
	IsActive      bool      `json:"is_active"`
}
