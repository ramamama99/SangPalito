package dto

import "time"

// UpdateQueueStatusRequest is the payload for PATCH /api/v1/queues/:id/status.
type UpdateQueueStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=waiting checked_in called in_service done cancelled no_show"`
}

// QueueResponse is the public representation of a queue entry.
type QueueResponse struct {
	ID            string     `json:"id"`
	FacilityID    string     `json:"facility_id"`
	AppointmentID string     `json:"appointment_id"`
	PolyclinicID  string     `json:"polyclinic_id"`
	QueueNumber   int        `json:"queue_number"`
	QueueCode     string     `json:"queue_code,omitempty"`
	Status           string     `json:"status"`
	AlreadyCheckedIn bool       `json:"already_checked_in"`
	CheckedInAt      *time.Time `json:"checked_in_at,omitempty"`
	CallCount     int        `json:"call_count"`
	CalledAt      *time.Time `json:"called_at,omitempty"`
	ServiceStart  *time.Time `json:"service_start,omitempty"`
	ServiceEnd    *time.Time `json:"service_end,omitempty"`
	QueueDate     string     `json:"queue_date"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
