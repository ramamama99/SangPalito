package entity

import (
	"math"
)

const (
	// DefaultEMASmoothingAlpha is the weight given to the latest completed consultation duration.
	DefaultEMASmoothingAlpha = 0.30

	// MinConsultationDurationMinutes is the lower bound clamp for consultation duration.
	MinConsultationDurationMinutes = 3.0

	// MaxConsultationDurationMinutes is the upper bound clamp to exclude anomalous outliers.
	MaxConsultationDurationMinutes = 45.0
)

// QueueWithPatient encapsulates queue and clinical complexity parameters for wait-time modeling.
type QueueWithPatient struct {
	Queue          *Queue
	PatientAge     int
	Priority       string
	ChiefComplaint string
	Symptoms       []string
	TriageFlag     string
}

// PatientComplexityInput holds patient factors that influence clinical consultation duration.
type PatientComplexityInput struct {
	AgeYears       int
	Priority       string   // normal, urgent, emergency
	TriageFlag     string   // green, yellow, red
	ChiefComplaint string
	Symptoms       []string
}

// CalculateComplexityFactor calculates the clinical duration multiplier (C_p).
// Base factor is 1.0. Returns the factor and human-readable breakdown notes.
func CalculateComplexityFactor(input PatientComplexityInput) (float64, []string) {
	factor := 1.0
	var notes []string

	// 1. Age-based complexity
	if input.AgeYears >= 60 {
		factor += 0.30
		notes = append(notes, "Pasien Geriatri/Lansia (≥60 thn): +30%")
	} else if input.AgeYears < 5 && input.AgeYears >= 0 {
		factor += 0.20
		notes = append(notes, "Pasien Pediatri/Balita (<5 thn): +20%")
	} else if input.AgeYears <= 12 && input.AgeYears >= 5 {
		factor += 0.10
		notes = append(notes, "Pasien Anak (5-12 thn): +10%")
	}

	// 2. Priority & Triage
	if input.Priority == "emergency" || input.TriageFlag == "red" {
		factor += 0.40
		notes = append(notes, "Prioritas Emergency / Triage Merah: +40%")
	} else if input.Priority == "urgent" || input.TriageFlag == "yellow" {
		factor += 0.25
		notes = append(notes, "Prioritas Urgent / Triage Kuning: +25%")
	}

	// 3. Multi-symptoms / Anamnesis complexity
	symptomCount := len(input.Symptoms)
	if symptomCount >= 4 {
		factor += 0.25
		notes = append(notes, "Kompleksitas Gejala Tinggi (≥4 keluhan): +25%")
	} else if symptomCount >= 2 || len(input.ChiefComplaint) > 60 {
		factor += 0.15
		notes = append(notes, "Keluhan Multi-gejala: +15%")
	}

	// Clamp between 0.7 and 2.5
	if factor < 0.7 {
		factor = 0.7
	}
	if factor > 2.5 {
		factor = 2.5
	}

	// Round to 2 decimal places
	factor = math.Round(factor*100) / 100
	return factor, notes
}

// CalculateEMA updates the Exponential Moving Average with actual duration data.
// If actualDurations is empty, returns the default baseline.
func CalculateEMA(baselineMinutes float64, actualDurations []float64, alpha float64) float64 {
	if alpha <= 0 || alpha > 1 {
		alpha = DefaultEMASmoothingAlpha
	}
	if baselineMinutes <= 0 {
		baselineMinutes = 10.0
	}

	ema := baselineMinutes
	for _, dur := range actualDurations {
		// Clamp outliers
		clamped := dur
		if clamped < MinConsultationDurationMinutes {
			clamped = MinConsultationDurationMinutes
		}
		if clamped > MaxConsultationDurationMinutes {
			clamped = MaxConsultationDurationMinutes
		}
		ema = (alpha * clamped) + ((1.0 - alpha) * ema)
	}

	return math.Round(ema*10) / 10
}
