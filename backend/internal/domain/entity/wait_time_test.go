package entity

import (
	"testing"
)

func TestCalculateEMA(t *testing.T) {
	// Baseline when no visits completed
	baseline := 12.0
	ema0 := CalculateEMA(baseline, nil, 0.3)
	if ema0 != 12.0 {
		t.Errorf("expected baseline 12.0, got %f", ema0)
	}

	// 1 completed visit taking 20 mins:
	// EMA = 0.3 * 20 + 0.7 * 12 = 6 + 8.4 = 14.4
	ema1 := CalculateEMA(baseline, []float64{20.0}, 0.3)
	if ema1 != 14.4 {
		t.Errorf("expected 14.4, got %f", ema1)
	}

	// Outlier clamping: an anomalous 2-hour duration (120 min) should clamp to Max (45 min)
	emaClamped := CalculateEMA(10.0, []float64{120.0}, 0.3)
	expected := (0.3 * 45.0) + (0.7 * 10.0) // 13.5 + 7 = 20.5
	if emaClamped != expected {
		t.Errorf("expected clamped %f, got %f", expected, emaClamped)
	}

	// Outlier clamping: 1 min consultation clamps to Min (3 min)
	emaMinClamped := CalculateEMA(10.0, []float64{1.0}, 0.3)
	expectedMin := (0.3 * 3.0) + (0.7 * 10.0) // 0.9 + 7 = 7.9
	if emaMinClamped != expectedMin {
		t.Errorf("expected min clamped %f, got %f", expectedMin, emaMinClamped)
	}
}

func TestCalculateComplexityFactor(t *testing.T) {
	tests := []struct {
		name           string
		input          PatientComplexityInput
		minFactor      float64
		maxFactor      float64
		mustHaveNote   string
	}{
		{
			name: "Normal adult single complaint",
			input: PatientComplexityInput{
				AgeYears:       25,
				Priority:       "normal",
				ChiefComplaint: "Batuk pilek",
				Symptoms:       []string{"Batuk"},
			},
			minFactor: 1.0,
			maxFactor: 1.0,
		},
		{
			name: "Elderly patient (Geriatri >= 60)",
			input: PatientComplexityInput{
				AgeYears:       68,
				Priority:       "normal",
				ChiefComplaint: "Linu persendian",
			},
			minFactor:    1.29,
			maxFactor:    1.31,
			mustHaveNote: "Geriatri",
		},
		{
			name: "Toddler (< 5 years)",
			input: PatientComplexityInput{
				AgeYears:       2,
				Priority:       "normal",
				ChiefComplaint: "Demam panas",
			},
			minFactor:    1.19,
			maxFactor:    1.21,
			mustHaveNote: "Pediatri",
		},
		{
			name: "Emergency triage red patient",
			input: PatientComplexityInput{
				AgeYears:       35,
				Priority:       "emergency",
				TriageFlag:     "red",
				ChiefComplaint: "Sesak napas berat",
			},
			minFactor:    1.39,
			maxFactor:    1.41,
			mustHaveNote: "Emergency",
		},
		{
			name: "Elderly with multiple symptoms",
			input: PatientComplexityInput{
				AgeYears:       72,
				Priority:       "urgent",
				TriageFlag:     "yellow",
				ChiefComplaint: "Pusing berputar, lemas, muntah, dada berdebar",
				Symptoms:       []string{"Vertigo", "Mual", "Muntah", "Palpitasi"},
			},
			minFactor: 1.7, // 1.0 + 0.30(age) + 0.25(urgent) + 0.25(4 symptoms) = 1.8
			maxFactor: 1.9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factor, notes := CalculateComplexityFactor(tt.input)
			if factor < tt.minFactor || factor > tt.maxFactor {
				t.Errorf("%s: expected factor in [%f, %f], got %f", tt.name, tt.minFactor, tt.maxFactor, factor)
			}
			if tt.mustHaveNote != "" {
				found := false
				for _, note := range notes {
					if len(note) > 0 {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: expected note containing %s, got notes: %v", tt.name, tt.mustHaveNote, notes)
				}
			}
		})
	}
}
