package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/repository"
	"sangpalito/backend/pkg/response"
)

// AuditLogHandler handles audit log queries for administrators.
type AuditLogHandler struct {
	repo repository.AuditLogRepository
}

// NewAuditLogHandler creates a new AuditLogHandler.
func NewAuditLogHandler(repo repository.AuditLogRepository) *AuditLogHandler {
	return &AuditLogHandler{repo: repo}
}

// List handles GET /api/v1/audit-logs
func (h *AuditLogHandler) List(c *gin.Context) {
	entityName := c.Query("entity")
	entityIDStr := c.Query("entity_id")
	var entityID *uuid.UUID
	if entityIDStr != "" {
		parsed, err := uuid.Parse(entityIDStr)
		if err != nil {
			_ = c.Error(apperror.BadRequest("entity_id harus berupa UUID yang valid", err))
			return
		}
		entityID = &parsed
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	logs, total, err := h.repo.Find(c.Request.Context(), entityName, entityID, page, limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OKWithMeta(logs, response.Meta{
		Page:  page,
		Limit: limit,
		Total: total,
	}))
}
