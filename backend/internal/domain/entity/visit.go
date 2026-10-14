package entity

import (
	"time"

	"github.com/google/uuid"
)

// Visit represents a clinical visit record.
type Visit struct {
	ID            uuid.UUID
	FacilityID    uuid.UUID
	AppointmentID uuid.UUID
	PatientID     uuid.UUID
	DoctorID      uuid.UUID
	PolyclinicID  uuid.UUID
	VisitDate     time.Time
	Diagnosis     string
	ICDCode       string
	TreatmentPlan string
	Prescription  string
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
