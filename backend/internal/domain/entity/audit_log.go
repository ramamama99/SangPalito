package entity

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an audit trail entry.
type AuditLog struct {
	ID        uuid.UUID
	UserID    *uuid.UUID
	Action    string // "create", "read", "update", "delete", "login"
	Entity    string // "patient", "queue", "visit", "user"
	EntityID  *uuid.UUID
	IPAddress string
	Meta      map[string]interface{}
	CreatedAt time.Time
}
