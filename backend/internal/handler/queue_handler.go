package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/domain/entity"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/internal/service"
	"sangpalito/backend/pkg/response"
)

// QueueHandler handles queue-related HTTP requests.
type QueueHandler struct {
	queueService service.QueueService
	validate     *validator.Validate
}

// NewQueueHandler creates a new QueueHandler.
func NewQueueHandler(queueService service.QueueService, validate *validator.Validate) *QueueHandler {
	return &QueueHandler{queueService: queueService, validate: validate}
}

// List handles GET /api/v1/queues.
func (h *QueueHandler) List(c *gin.Context) {
	polyclinicIDStr := c.Query("polyclinic_id")
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	polyclinicID, err := uuid.Parse(polyclinicIDStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("polyclinic_id harus UUID yang valid", err))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("Format tanggal tidak valid (YYYY-MM-DD)", err))
		return
	}

	results, err := h.queueService.GetActiveQueues(c.Request.Context(), polyclinicID, date)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(results))
}

// UpdateStatus handles PATCH /api/v1/queues/:id/status.
func (h *QueueHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(apperror.BadRequest("ID tidak valid", err))
		return
	}

	var req dto.UpdateQueueStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.BadRequest("Request body tidak valid", err))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		_ = c.Error(apperror.ValidationError(err.Error()))
		return
	}

	result, err := h.queueService.AdvanceQueueStatus(c.Request.Context(), id, entity.QueueStatus(req.Status))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.OK(result))
}
