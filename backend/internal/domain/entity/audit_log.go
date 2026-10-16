package entity

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an audit trail entry.
type AuditLog struct {
	ID        int64                  `json:"id"`
	UserID    *uuid.UUID             `json:"user_id,omitempty"`
	UserName  string                 `json:"user_name,omitempty"`
	UserEmail string                 `json:"user_email,omitempty"`
	UserRole  string                 `json:"user_role,omitempty"`
	Action    string                 `json:"action"` // "create", "read", "update", "delete", "login", "logout", "exam_patient"
	Entity    string                 `json:"entity"` // "patient", "queue", "visit", "user", "auth"
	EntityID  *uuid.UUID             `json:"entity_id,omitempty"`
	IPAddress string                 `json:"ip_address,omitempty"`
	Meta      map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}
