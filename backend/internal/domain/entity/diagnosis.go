package entity

import "github.com/google/uuid"

// Master Kode Penyakit Standar WHO
type ICD10Code struct {
	Code        string `json:"code"`        // Contoh: "J06.9", "I10"
	Description string `json:"description"` // Contoh: "Acute upper respiratory infection, unspecified"
}

// Diagnosis yang ditegakkan dokter untuk kunjungan pasien
type Diagnosis struct {
	ID               uuid.UUID  `json:"id"`
	VisitID          uuid.UUID  `json:"visit_id"`
	ICD10Code        string     `json:"icd10_code"`
	IsPrimary        bool       `json:"is_primary"` // Diagnosis utama (wajib 1 per kunjungan)
	Notes            *string    `json:"notes,omitempty"`
	ICD10Description string     `json:"icd10_description,omitempty"`
	ICD10            *ICD10Code `json:"icd10,omitempty"`
}
