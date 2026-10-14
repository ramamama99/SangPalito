package entity

import (
	"time"

	"github.com/google/uuid"
)

// Appointment represents a patient visit registration/appointment.
type Appointment struct {
	ID               uuid.UUID
	FacilityID       uuid.UUID
	PatientID        uuid.UUID
	PolyclinicID     uuid.UUID
	ScheduleID       uuid.UUID
	DoctorID         uuid.UUID
	VisitDate        time.Time
	AppointmentToken uuid.UUID
	QRToken          uuid.UUID
	ChiefComplaint   string
	Status           string // "registered", "checked_in", "completed", "cancelled"
	RegisteredBy     *uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
