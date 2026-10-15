package entity

import (
	"time"

	"github.com/google/uuid"
)

// Prescription represents a digital medical prescription issued by a doctor.
type Prescription struct {
	ID             uuid.UUID          `json:"id"`
	PrescriptionNo string             `json:"prescription_no"`
	AppointmentID  *uuid.UUID         `json:"appointment_id,omitempty"`
	VisitID        *uuid.UUID         `json:"visit_id,omitempty"`
	PatientID      uuid.UUID          `json:"patient_id"`
	DoctorID       *uuid.UUID         `json:"doctor_id,omitempty"`
	DoctorName     string             `json:"doctor_name,omitempty"`
	Status         string             `json:"status"` // pending_dispense, dispensed, cancelled
	Notes          string             `json:"notes,omitempty"`
	DispensedAt    *time.Time         `json:"dispensed_at,omitempty"`
	DispensedBy    *uuid.UUID         `json:"dispensed_by,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	Items          []PrescriptionItem `json:"items,omitempty"`
}

// PrescriptionItem represents an individual medication entry in a prescription.
type PrescriptionItem struct {
	ID             uuid.UUID  `json:"id"`
	PrescriptionID uuid.UUID  `json:"prescription_id"`
	MedicationID   *uuid.UUID `json:"medication_id,omitempty"`
	MedicationName string     `json:"medication_name"`
	DosageForm     string     `json:"dosage_form"`
	Quantity       int        `json:"quantity"`
	Unit           string     `json:"unit"`
	Route          string     `json:"route"`
	Frequency      string     `json:"frequency"`
	Timing         string     `json:"timing"`
	Instruction    string     `json:"instruction"`
	Notes          string     `json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
