package entity

import (
	"time"

	"github.com/google/uuid"
)

// Schedule represents a doctor's practice schedule.
type Schedule struct {
	ID           uuid.UUID
	FacilityID   uuid.UUID
	PolyclinicID uuid.UUID
	DoctorID     uuid.UUID
	DayOfWeek    int // 0=Sunday, 1=Monday, ..., 6=Saturday
	StartTime    string
	EndTime      string
	Quota        int
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Joined fields (populated by certain queries)
	DoctorName     string
	PolyclinicName string
}
