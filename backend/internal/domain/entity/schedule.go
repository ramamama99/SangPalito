package entity

import (
	"time"

	"github.com/google/uuid"
)

// Schedule represents a doctor's clinic schedule.
type Schedule struct {
	ID            uuid.UUID `json:"id"`
	FacilityID    uuid.UUID `json:"facility_id,omitempty"`
	PolyclinicID  uuid.UUID `json:"polyclinic_id"`
	DoctorID      uuid.UUID `json:"doctor_id"`
	DayOfWeek     int       `json:"day_of_week"`    // 1=Senin ... 7=Minggu
	StartTime     string    `json:"start_time"`     // format "08:00:00"
	EndTime       string    `json:"end_time"`       // format "12:00:00"
	Quota         int       `json:"quota"`          // batas kuota pendaftaran
	RemainingSlot int       `json:"remaining_slot"` // sisa kuota hari ini
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`

	// Joined fields
	DoctorName     string `json:"doctor_name,omitempty"`
	PolyclinicName string `json:"polyclinic_name,omitempty"`
}
