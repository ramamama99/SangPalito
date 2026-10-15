package entity

import (
	"time"

	"github.com/google/uuid"
)

type Vitals struct {
	BPSys           *int     `json:"bp_sys,omitempty"`           // Tekanan Darah Sistolik (misal: 120)
	BPDia           *int     `json:"bp_dia,omitempty"`           // Tekanan Darah Diastolik (misal: 80)
	BloodPressure   string   `json:"blood_pressure,omitempty"`   // e.g. "120/80"
	HR              *int     `json:"hr,omitempty"`               // Heart Rate / Denyut Nadi (bpm)
	HeartRate       *int     `json:"heart_rate,omitempty"`       // Alias bpm
	RespiratoryRate *int     `json:"respiratory_rate,omitempty"` // x/min
	Temp            *float64 `json:"temp,omitempty"`             // Suhu Tubuh (°C, misal: 36.5)
	Temperature     *float64 `json:"temperature,omitempty"`      // Alias Celcius
	SpO2            *int     `json:"spo2,omitempty"`             // Saturasi Oksigen (%)
	Weight          *float64 `json:"weight,omitempty"`           // Berat Badan (kg)
	Height          *float64 `json:"height,omitempty"`           // Tinggi Badan (cm)
}

// DiagnosisItem represents an ICD-10 coded diagnosis associated with a visit.
type DiagnosisItem struct {
	ID          uuid.UUID `json:"id"`
	VisitID     uuid.UUID `json:"visit_id"`
	ICD10Code   string    `json:"icd10_code"`
	Description string    `json:"description,omitempty"`
	IsPrimary   bool      `json:"is_primary"`
	Notes       string    `json:"notes,omitempty"`
}

type Visit struct {
	ID              uuid.UUID       `json:"id"`
	AppointmentID   uuid.UUID       `json:"appointment_id"`
	PatientID       uuid.UUID       `json:"patient_id"`
	DoctorID        *uuid.UUID      `json:"doctor_id,omitempty"`
	DoctorName      string          `json:"doctor_name,omitempty"`
	DoctorSpecialty string          `json:"doctor_specialty,omitempty"`
	PolyclinicID    uuid.UUID       `json:"polyclinic_id,omitempty"`
	PolyclinicName  string          `json:"polyclinic_name,omitempty"`
	VisitDate       time.Time       `json:"visit_date,omitempty"`
	Vitals          *Vitals         `json:"vitals,omitempty"`
	Subjective      string          `json:"subjective"` // S: Keluhan pasien yang diceritakan
	Objective       string          `json:"objective"`  // O: Hasil pengamatan fisik dokter
	Assessment      string          `json:"assessment"` // A: Kesimpulan klinis awal dokter
	Plan            string          `json:"plan"`       // P: Rencana terapi / obat / rujukan
	Prescription    string          `json:"prescription,omitempty"`
	CreatedBy       *uuid.UUID      `json:"created_by,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`

	// Daftar diagnosis ICD-10 untuk kunjungan ini
	Diagnoses []Diagnosis `json:"diagnoses,omitempty"`

	// Field kompatibilitas riwayat rekam medis pasien (patient_service)
	Diagnosis     string `json:"diagnosis,omitempty"`
	ICDCode       string `json:"icd_code,omitempty"`
	TreatmentPlan string `json:"treatment_plan,omitempty"`
}
