package entity

import (
	"time"

	"github.com/google/uuid"
)

// Medication represents a cataloged medicine/drug.
type Medication struct {
	ID                     uuid.UUID `json:"id"`
	Code                   string    `json:"code"`
	Name                   string    `json:"name"`
	Strength               string    `json:"strength,omitempty"`
	GenericName            string    `json:"generic_name"`
	Category               string    `json:"category"`
	DosageForm             string    `json:"dosage_form"`
	Route                  string    `json:"route"`
	DefaultFrequency       string    `json:"default_frequency"` // 3 x 1 tablet sehari
	DefaultTiming          string    `json:"default_timing"`    // Sesudah makan, Sebelum makan

	Indication             string    `json:"indication"`        // Indikasi
	Unit                   string    `json:"unit"`
	IsPrescriptionRequired bool      `json:"is_prescription_required"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
