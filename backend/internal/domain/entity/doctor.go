package entity

import (
	"github.com/google/uuid"
)

type Doctor struct {
	ID         uuid.UUID  `json:"id"`
	FacilityID uuid.UUID  `json:"facility_id"`
	UserID     *uuid.UUID `json:"user_id,omitempty"` // ID akun staf (bisa null)
	FullName   string     `json:"full_name"`
	Specialty  *string    `json:"specialty,omitempty"`
	STRNo      *string    `json:"str_no,omitempty"`
	IsActive   bool       `json:"is_active"`
}
