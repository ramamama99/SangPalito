package entity

import (
	"time"

	"github.com/google/uuid"
)

// QueueStatus represents the status of a queue entry.
type QueueStatus string

const (
	StatusWaiting   QueueStatus = "waiting"
	StatusCheckedIn QueueStatus = "checked_in"
	StatusCalled    QueueStatus = "called"
	StatusInService QueueStatus = "in_service"
	StatusDone      QueueStatus = "done"
	StatusCancelled QueueStatus = "cancelled"
	StatusNoShow    QueueStatus = "no_show"
)

// ValidTransitions defines the allowed state machine transitions for queue status.
var ValidTransitions = map[QueueStatus][]QueueStatus{
	StatusWaiting:   {StatusCheckedIn, StatusCancelled},
	StatusCheckedIn: {StatusCalled, StatusCancelled},
	StatusCalled:    {StatusInService, StatusNoShow, StatusCancelled},
	StatusInService: {StatusDone},
	// done, cancelled, no_show are terminal states — no transitions allowed
}

// CanTransitionTo checks if the current status can transition to the target status.
func (s QueueStatus) CanTransitionTo(target QueueStatus) bool {
	allowed, exists := ValidTransitions[s]
	if !exists {
		return false
	}
	for _, a := range allowed {
		if a == target {
			return true
		}
	}
	return false
}

// IsTerminal returns true if this status is a terminal (final) state.
func (s QueueStatus) IsTerminal() bool {
	return s == StatusDone || s == StatusCancelled || s == StatusNoShow
}

// Queue represents a queue entry for a patient appointment.
type Queue struct {
	ID            uuid.UUID
	FacilityID    uuid.UUID
	AppointmentID uuid.UUID
	PolyclinicID  uuid.UUID
	QueueNumber   int
	QueueCode     string
	Status        QueueStatus
	CheckedInAt   *time.Time
	CallCount     int
	CalledAt      *time.Time
	ServiceStart  *time.Time
	ServiceEnd    *time.Time
	QueueDate     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
