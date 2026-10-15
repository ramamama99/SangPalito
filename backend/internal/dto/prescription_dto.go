package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePrescriptionItemRequest struct {
	MedicationID   *uuid.UUID `json:"medication_id"`
	MedicationName string     `json:"medication_name" binding:"required"`
	DosageForm     string     `json:"dosage_form" binding:"required"`
	Quantity       int        `json:"quantity" binding:"required,min=1"`
	Unit           string     `json:"unit" binding:"required"`
	Route          string     `json:"route" binding:"required"` // oral_minum, oral_kunyah, topikal_oles, tetes_mata, tetes_telinga, kumur, inhalasi
	Frequency      string     `json:"frequency" binding:"required"`
	Timing         string     `json:"timing" binding:"required"`
	Instruction    string     `json:"instruction" binding:"required"`
	Notes          string     `json:"notes"`
}

type CreatePrescriptionRequest struct {
	AppointmentID *uuid.UUID                      `json:"appointment_id"`
	VisitID       *uuid.UUID                      `json:"visit_id"`
	PatientID     uuid.UUID                       `json:"patient_id" binding:"required"`
	DoctorID      *uuid.UUID                      `json:"doctor_id"`
	DoctorName    string                          `json:"doctor_name"`
	Notes         string                          `json:"notes"`
	Items         []CreatePrescriptionItemRequest `json:"items" binding:"required,min=1"`
}

type PrescriptionItemResponse struct {
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
}

type PrescriptionResponse struct {
	ID             uuid.UUID                  `json:"id"`
	PrescriptionNo string                     `json:"prescription_no"`
	AppointmentID  *uuid.UUID                 `json:"appointment_id,omitempty"`
	VisitID        *uuid.UUID                 `json:"visit_id,omitempty"`
	PatientID      uuid.UUID                  `json:"patient_id"`
	DoctorID       *uuid.UUID                 `json:"doctor_id,omitempty"`
	DoctorName     string                     `json:"doctor_name,omitempty"`
	Status         string                     `json:"status"`
	StatusLabel    string                     `json:"status_label"`
	Notes          string                     `json:"notes,omitempty"`
	DispensedAt    *time.Time                 `json:"dispensed_at,omitempty"`
	CreatedAt      time.Time                  `json:"created_at"`
	Items          []PrescriptionItemResponse `json:"items"`
}
