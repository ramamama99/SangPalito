package repository

import (
	"context"

	"sangpalito/backend/internal/domain/entity"
)

// AuditLogRepository defines the contract for audit log data access.
type AuditLogRepository interface {
	Create(ctx context.Context, log *entity.AuditLog) error
}
