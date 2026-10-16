package dto

import "github.com/google/uuid"

// Polyclinic DTO
type CreatePolyclinicRequest struct {
	Name          string `json:"name" validate:"required,min=2,max=80"`
	QueuePrefix   string `json:"queue_prefix" validate:"required,len=1,uppercase"` // Contoh: "A", "B"
	AvgServiceMin int    `json:"avg_service_min" validate:"required,min=1,max=60"`
}

type UpdatePolyclinicRequest struct {
	Name          string `json:"name" validate:"required,min=2,max=80"`
	QueuePrefix   string `json:"queue_prefix" validate:"required,len=1,uppercase"`
	AvgServiceMin int    `json:"avg_service_min" validate:"required,min=1,max=60"`
	IsActive      bool   `json:"is_active"`
}

// Doctor DTO
type CreateDoctorRequest struct {
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	FullName     string     `json:"full_name" validate:"required,min=2,max=120"`
	Email        string     `json:"email" validate:"required,email,max=100"`
	Phone        string     `json:"phone" validate:"required,min=9,max=25"`
	Specialty    *string    `json:"specialty,omitempty"`
	STRNo        *string    `json:"str_no,omitempty"`
	PolyclinicID *uuid.UUID `json:"polyclinic_id,omitempty"`
	Password     string     `json:"password,omitempty"`
}

type UpdateDoctorRequest struct {
	FullName     string     `json:"full_name" validate:"required,min=2,max=120"`
	Email        string     `json:"email,omitempty"`
	Phone        string     `json:"phone,omitempty"`
	Specialty    *string    `json:"specialty,omitempty"`
	STRNo        *string    `json:"str_no,omitempty"`
	PolyclinicID *uuid.UUID `json:"polyclinic_id,omitempty"`
	IsActive     bool       `json:"is_active"`
}

// Schedule DTO
type CreateScheduleRequest struct {
	PolyclinicID uuid.UUID `json:"polyclinic_id" validate:"required"`
	DoctorID     uuid.UUID `json:"doctor_id" validate:"required"`
	DayOfWeek    int       `json:"day_of_week" validate:"required,min=1,max=7"` // 1=Senin ... 7=Minggu
	StartTime    string    `json:"start_time" validate:"required"`              // Contoh: "08:00:00"
	EndTime      string    `json:"end_time" validate:"required"`                // Contoh: "12:00:00"
	Quota        int       `json:"quota" validate:"required,min=1"`
}

type SaveShiftAssignmentRequest struct {
	PolyclinicID uuid.UUID `json:"polyclinic_id" validate:"required"`
	DoctorID     uuid.UUID `json:"doctor_id" validate:"required"`
	DayOfWeek    int       `json:"day_of_week" validate:"required,min=1,max=7"`
	ShiftID      string    `json:"shift_id" validate:"required"` // "pagi", "siang", "sore"
	StartTime    string    `json:"start_time"`
	EndTime      string    `json:"end_time"`
	Quota        int       `json:"quota"`
}
