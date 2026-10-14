package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"sangpalito/backend/internal/domain/entity"
)

type auditLogRepo struct {
	pool *pgxpool.Pool
}

// NewAuditLogRepo creates a new PostgreSQL audit log repository.
func NewAuditLogRepo(pool *pgxpool.Pool) *auditLogRepo {
	return &auditLogRepo{pool: pool}
}

// Create inserts a new audit log entry.
func (r *auditLogRepo) Create(ctx context.Context, log *entity.AuditLog) error {
	query := `INSERT INTO audit_logs (id, user_id, action, entity, entity_id, ip_address, meta, created_at)
	           VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	metaJSON, _ := json.Marshal(log.Meta)

	_, err := r.pool.Exec(ctx, query,
		log.ID, log.UserID, log.Action, log.Entity, log.EntityID,
		log.IPAddress, metaJSON, log.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("auditLogRepo.Create: %w", err)
	}
	return nil
}
