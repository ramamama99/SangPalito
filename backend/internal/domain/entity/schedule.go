package entity

import (
	"github.com/google/uuid"
)

type Schedule struct {
	ID            uuid.UUID `json:"id"`
	PolyclinicID  uuid.UUID `json:"polyclinic_id"`
	DoctorID      uuid.UUID `json:"doctor_id"`
	DayOfWeek     int       `json:"day_of_week"`    // 1=Senin ... 7=Minggu
	StartTime     string    `json:"start_time"`     // format "08:00:00"
	EndTime       string    `json:"end_time"`       // format "12:00:00"
	Quota         int       `json:"quota"`          // batas kuota pendaftaran
	RemainingSlot int       `json:"remaining_slot"` // sisa kuota hari ini
}
