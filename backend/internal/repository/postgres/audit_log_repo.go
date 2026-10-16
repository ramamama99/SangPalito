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
	query := `INSERT INTO audit_logs (user_id, action, entity, entity_id, ip_address, metadata, created_at)
	           VALUES ($1, $2, $3, $4, $5, $6, $7)`

	metaJSON, _ := json.Marshal(log.Meta)

	var ipParam any
	if log.IPAddress != "" {
		ipParam = log.IPAddress
	}

	_, err := r.pool.Exec(ctx, query,
		log.UserID, log.Action, log.Entity, log.EntityID,
		ipParam, metaJSON, log.CreatedAt,
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
		whereClause += fmt.Sprintf(" AND a.entity = $%d", argIdx)
		args = append(args, entityName)
		argIdx++
	}
	if entityID != nil {
		whereClause += fmt.Sprintf(" AND a.entity_id = $%d", argIdx)
		args = append(args, *entityID)
		argIdx++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs a %s", whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("auditLogRepo.Find count: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT a.id, a.user_id, a.action, a.entity, a.entity_id, a.ip_address, a.metadata, a.created_at,
		       COALESCE(u.full_name, 'Sistem / Anonim') AS user_name,
		       COALESCE(u.email, '-') AS user_email,
		       COALESCE(r.name, '-') AS user_role
		FROM audit_logs a
		LEFT JOIN users u ON u.id = a.user_id
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN roles r ON r.id = ur.role_id
		%s
		ORDER BY a.created_at DESC
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
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.Action, &l.Entity, &l.EntityID, &l.IPAddress, &metaJSON, &l.CreatedAt,
			&l.UserName, &l.UserEmail, &l.UserRole,
		); err != nil {
			return nil, 0, fmt.Errorf("auditLogRepo.Find scan: %w", err)
		}
		if len(metaJSON) > 0 {
			_ = json.Unmarshal(metaJSON, &l.Meta)
		}
		logs = append(logs, &l)
	}

	return logs, total, nil
}
