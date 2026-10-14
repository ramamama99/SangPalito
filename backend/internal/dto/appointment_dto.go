package dto

import "github.com/google/uuid"

type CreateAppointmentRequest struct {
	PatientID    uuid.UUID                  `json:"patient_id" validate:"required"`
	PolyclinicID uuid.UUID                  `json:"polyclinic_id" validate:"required"`
	ScheduleID   *uuid.UUID                 `json:"schedule_id,omitempty"`
	VisitDate    string                     `json:"visit_date" validate:"required"` // YYYY-MM-DD
	Priority     string                     `json:"priority"`                       // normal | elderly | pregnant | disabled | emergency
	PreAnamnesis *CreatePreAnamnesisRequest `json:"pre_anamnesis,omitempty"`
}

type CreatePreAnamnesisRequest struct {
	ChiefComplaint string   `json:"chief_complaint" validate:"required"`
	Symptoms       []string `json:"symptoms"`
	DurationDays   *int     `json:"duration_days,omitempty"`
	Allergies      *string  `json:"allergies,omitempty"`
}

type UpdatePreAnamnesisRequest struct {
	ChiefComplaint string   `json:"chief_complaint" validate:"required"`
	Symptoms       []string `json:"symptoms"`
	DurationDays   *int     `json:"duration_days,omitempty"`
	Allergies      *string  `json:"allergies,omitempty"`
}
