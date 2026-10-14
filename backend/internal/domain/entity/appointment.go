package entity

import (
	"time"

	"github.com/google/uuid"
)

// Appointment represents a patient visit registration/appointment.
type Appointment struct {
	ID               uuid.UUID  `json:"id"`
	FacilityID       uuid.UUID  `json:"facility_id"`
	PatientID        uuid.UUID  `json:"patient_id"`
	PolyclinicID     uuid.UUID  `json:"polyclinic_id"`
	ScheduleID       uuid.UUID  `json:"schedule_id"`
	DoctorID         uuid.UUID  `json:"doctor_id"`
	VisitDate        time.Time  `json:"visit_date"`
	AppointmentToken uuid.UUID  `json:"appointment_token"`
	QRToken          uuid.UUID  `json:"qr_token"`
	ChiefComplaint   string     `json:"chief_complaint"`
	Status           string     `json:"status"` // "registered", "checked_in", "completed", "cancelled"
	RegisteredBy     *uuid.UUID `json:"registered_by,omitempty"`
	Source           string     `json:"source,omitempty"`   // "online" | "walk_in"
	Priority         string     `json:"priority,omitempty"` // "normal" | "elderly" | "pregnant" | "disabled" | "emergency"
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Relasi data tambahan
	PatientFullName string        `json:"patient_full_name,omitempty"`
	PolyclinicName  string        `json:"polyclinic_name,omitempty"`
	PreAnamnesis    *PreAnamnesis `json:"pre_anamnesis,omitempty"`
}
