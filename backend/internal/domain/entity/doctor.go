package entity

import (
	"time"

	"github.com/google/uuid"
)

// Doctor represents a doctor entity in the facility.
type Doctor struct {
	ID             uuid.UUID  `json:"id"`
	FacilityID     uuid.UUID  `json:"facility_id"`
	UserID         *uuid.UUID `json:"user_id,omitempty"` // ID akun staf (bisa null)
	FullName       string     `json:"full_name"`
	Email          string     `json:"email,omitempty"`
	Phone          string     `json:"phone,omitempty"`
	Specialty      *string    `json:"specialty,omitempty"`
	Specialization string     `json:"specialization,omitempty"`
	STRNo          *string    `json:"str_no,omitempty"`
	SIPNo          string     `json:"sip_no,omitempty"`
	PolyclinicID   *uuid.UUID `json:"polyclinic_id,omitempty"`
	PolyclinicName string     `json:"polyclinic_name,omitempty"`
	IsActive       bool       `json:"is_active"`
	IsFirstLogin   bool       `json:"is_first_login,omitempty"`
	CreatedAt      time.Time  `json:"created_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at,omitempty"`
}
