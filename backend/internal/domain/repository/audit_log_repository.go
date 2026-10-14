package repository

import (
	"context"

	"github.com/google/uuid"
	"sangpalito/backend/internal/domain/entity"
)

// AuditLogRepository defines the contract for audit log data access.
type AuditLogRepository interface {
	Create(ctx context.Context, log *entity.AuditLog) error
	Find(ctx context.Context, entityName string, entityID *uuid.UUID, page, limit int) ([]*entity.AuditLog, int, error)
}
