package entity

import (
	"math"
	"time"

	"github.com/google/uuid"
)

// Patient represents a registered patient.
type Patient struct {
	ID                uuid.UUID
	FacilityID        uuid.UUID
	MedicalRecordNo   string
	NIKHash           string
	NIKEncrypted      []byte
	FullName          string
	BirthDate         time.Time
	Gender            string // "male" or "female"
	Phone             string
	Address           string
	BloodType         string
	Guarantor         string // "bpjs", "umum", "asuransi", "lainnya"
	BPJSNo            string
	GuardianName      *string
	GuardianNIKHash   *string
	GuardianNIKEncrypted *[]byte
	GuardianPhone     *string
	GuardianRelation  *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

// Age calculates the patient's age in years from their birth date.
// Age is computed at runtime, not stored in the database.
func (p *Patient) Age() int {
	now := time.Now()
	years := now.Year() - p.BirthDate.Year()
	// Adjust if birthday hasn't occurred yet this year
	if now.YearDay() < p.BirthDate.YearDay() {
		years--
	}
	return int(math.Max(0, float64(years)))
}

// IsMinor returns true if the patient is under 17 years old.
func (p *Patient) IsMinor() bool {
	return p.Age() < 17
}

// RequiresGuardian returns true if the patient is a minor who needs guardian info.
func (p *Patient) RequiresGuardian() bool {
	return p.IsMinor()
}
