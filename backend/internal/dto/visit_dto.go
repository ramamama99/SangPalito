package dto

import (
	"time"

	"sangpalito/backend/internal/domain/entity"

	"github.com/google/uuid"
)

type CreateDiagnosisRequest struct {
	ICD10Code string `json:"icd10_code" binding:"required"`
	IsPrimary bool   `json:"is_primary"`
	Notes     string `json:"notes"`
}

type CreateVisitRequest struct {
	AppointmentID   uuid.UUID                       `json:"appointment_id" binding:"required"`
	PatientID       uuid.UUID                       `json:"patient_id" binding:"required"`
	DoctorID        uuid.UUID                       `json:"doctor_id" binding:"required"`
	Vitals          entity.Vitals                   `json:"vitals"`
	Subjective      string                          `json:"subjective" binding:"required"`
	Objective       string                          `json:"objective"`
	Assessment      string                          `json:"assessment" binding:"required"`
	Plan            string                          `json:"plan" binding:"required"`
	MedicationPlan  string                          `json:"medication_plan"`
	PrescriptionNo  string                          `json:"prescription_no"`
	Diagnoses       []CreateDiagnosisRequest        `json:"diagnoses"`
	Prescription    *CreatePrescriptionRequest      `json:"prescription,omitempty"`
}

type UpdateVisitRequest struct {
	Vitals         *entity.Vitals `json:"vitals"`
	Subjective     string         `json:"subjective"`
	Objective      string         `json:"objective"`
	Assessment     string         `json:"assessment"`
	Plan           string         `json:"plan"`
	MedicationPlan string         `json:"medication_plan"`
}

type DiagnosisResponse struct {
	ID          uuid.UUID `json:"id"`
	ICD10Code   string    `json:"icd10_code"`
	Description string    `json:"description,omitempty"`
	IsPrimary   bool      `json:"is_primary"`
	Notes       string    `json:"notes"`
}

type VisitResponse struct {
	ID             uuid.UUID             `json:"id"`
	AppointmentID  uuid.UUID             `json:"appointment_id"`
	PatientID      uuid.UUID             `json:"patient_id"`
	DoctorID       uuid.UUID             `json:"doctor_id"`
	Vitals         entity.Vitals         `json:"vitals"`
	Subjective     string                `json:"subjective"`
	Objective      string                `json:"objective"`
	Assessment     string                `json:"assessment"`
	Plan           string                `json:"plan"`
	MedicationPlan string                `json:"medication_plan,omitempty"`
	PrescriptionNo string                `json:"prescription_no,omitempty"`
	Diagnoses      []DiagnosisResponse   `json:"diagnoses"`
	Prescription   *PrescriptionResponse `json:"prescription,omitempty"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
}

type ICD10Response struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type LogVisitAuditRequest struct {
	VisitID        string                 `json:"visit_id"`
	AppointmentID  string                 `json:"appointment_id"`
	PatientID      string                 `json:"patient_id"`
	PatientName    string                 `json:"patient_name"`
	Polyclinic     string                 `json:"polyclinic"`
	DoctorName     string                 `json:"doctor_name"`
	Diagnosis      string                 `json:"diagnosis"`
	PrescriptionNo string                 `json:"prescription_no"`
	Vitals         map[string]interface{} `json:"vitals"`
	Subjective     string                 `json:"subjective"`
	Objective      string                 `json:"objective"`
	Plan           string                 `json:"plan"`
	MedicationPlan string                 `json:"medication_plan"`
}

