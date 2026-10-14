package entity

import (
	"time"

	"github.com/google/uuid"
)

type Vitals struct {
	BPSys  *int     `json:"bp_sys,omitempty"` // Tekanan Darah Sistolik (misal: 120)
	BPDia  *int     `json:"bp_dia,omitempty"` // Tekanan Darah Diastolik (misal: 80)
	HR     *int     `json:"hr,omitempty"`     // Heart Rate / Denyut Nadi (bpm)
	Temp   *float64 `json:"temp,omitempty"`   // Suhu Tubuh (°C, misal: 36.5)
	SpO2   *int     `json:"spo2,omitempty"`   // Saturasi Oksigen (%)
	Weight *float64 `json:"weight,omitempty"` // Berat Badan (kg)
	Height *float64 `json:"height,omitempty"` // Tinggi Badan (cm)
}

type Visit struct {
	ID            uuid.UUID  `json:"id"`
	AppointmentID uuid.UUID  `json:"appointment_id"`
	PatientID     uuid.UUID  `json:"patient_id"`
	DoctorID      *uuid.UUID `json:"doctor_id,omitempty"`
	Vitals        *Vitals    `json:"vitals,omitempty"`
	Subjective    string     `json:"subjective"` // S: Keluhan pasien yang diceritakan
	Objective     string     `json:"objective"`  // O: Hasil pengamatan fisik dokter
	Assessment    string     `json:"assessment"` // A: Kesimpulan klinis awal dokter
	Plan          string     `json:"plan"`       // P: Rencana terapi / obat / rujukan
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Daftar diagnosis ICD-10 untuk kunjungan ini
	Diagnoses []Diagnosis `json:"diagnoses,omitempty"`

	// Field kompatibilitas riwayat rekam medis pasien (patient_service)
	VisitDate     time.Time `json:"visit_date,omitempty"`
	Diagnosis     string    `json:"diagnosis,omitempty"`
	ICDCode       string    `json:"icd_code,omitempty"`
	TreatmentPlan string    `json:"treatment_plan,omitempty"`
}
