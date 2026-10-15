package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
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
	query := `INSERT INTO audit_logs (id, user_id, action, entity, entity_id, ip_address, metadata, created_at)
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

// Find retrieves paginated audit logs filtered by entity and/or entity_id.
func (r *auditLogRepo) Find(ctx context.Context, entityName string, entityID *uuid.UUID, page, limit int) ([]*entity.AuditLog, int, error) {
	offset := (page - 1) * limit

	whereClause := "WHERE 1=1"
	args := []any{}
	argIdx := 1

	if entityName != "" {
		whereClause += fmt.Sprintf(" AND entity = $%d", argIdx)
		args = append(args, entityName)
		argIdx++
	}
	if entityID != nil {
		whereClause += fmt.Sprintf(" AND entity_id = $%d", argIdx)
		args = append(args, *entityID)
		argIdx++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs %s", whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("auditLogRepo.Find count: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, action, entity, entity_id, ip_address, metadata, created_at
		FROM audit_logs
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("auditLogRepo.Find query: %w", err)
	}
	defer rows.Close()

	logs := make([]*entity.AuditLog, 0)
	for rows.Next() {
		var l entity.AuditLog
		var metaJSON []byte
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.Entity, &l.EntityID, &l.IPAddress, &metaJSON, &l.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("auditLogRepo.Find scan: %w", err)
		}
		if len(metaJSON) > 0 {
			_ = json.Unmarshal(metaJSON, &l.Meta)
		}
		logs = append(logs, &l)
	}

	return logs, total, nil
}
