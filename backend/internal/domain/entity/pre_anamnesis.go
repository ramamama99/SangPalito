package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type PreAnamnesis struct {
	ID             uuid.UUID `json:"id"`
	AppointmentID  uuid.UUID `json:"appointment_id"`
	ChiefComplaint string    `json:"chief_complaint"` // Keluhan utama (teks bebas)
	Symptoms       []string  `json:"symptoms"`        // ["Demam", "Batuk", "Sesak napas"]
	DurationDays   *int      `json:"duration_days,omitempty"`
	Allergies      *string   `json:"allergies,omitempty"` // Riwayat alergi ("tidak ada", dll)
	TriageFlag     string    `json:"triage_flag"`         // "normal" | "attention" | "urgent"
	CreatedAt      time.Time `json:"created_at"`
}

// DetermineTriage menentukan otomatis tingkat keparahan pasien
func DetermineTriage(symptoms []string, durationDays int) string {
	urgentKeywords := []string{
		"sesak napas", "sesak berat", "nyeri dada", "tidak sadar", "kejang", "pendarahan",
	}

	for _, sym := range symptoms {
		s := strings.ToLower(strings.TrimSpace(sym))
		for _, urgent := range urgentKeywords {
			if strings.Contains(s, urgent) {
				return "urgent" // Merah / Darurat!
			}
		}
	}

	// Jika sakit sudah lebih dari 7 hari atau gejala lebih dari 4 macam
	if durationDays > 7 || len(symptoms) >= 5 {
		return "attention" // Kuning / Butuh perhatian khusus
	}

	return "normal" 
}