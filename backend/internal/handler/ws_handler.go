package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sangpalito/backend/internal/apperror"
	"sangpalito/backend/internal/dto"
	"sangpalito/backend/internal/hub"
	"sangpalito/backend/pkg/response"
)

// WSHandler manages WebSocket connections and public display screen polling endpoints.
type WSHandler struct {
	hub  *hub.Hub
	pool *pgxpool.Pool
}

// NewWSHandler creates a new WSHandler.
func NewWSHandler(hub *hub.Hub, pool *pgxpool.Pool) *WSHandler {
	return &WSHandler{
		hub:  hub,
		pool: pool,
	}
}

// ServeWS upgrades the HTTP connection to a WebSocket stream.
// Supports both global listener (/ws/queues) and polyclinic room (/ws/queues/:id).
func (h *WSHandler) ServeWS(c *gin.Context) {
	polyclinicID := c.Param("id")
	if polyclinicID == "" {
		polyclinicID = c.Param("polyclinic_id")
	}

	if polyclinicID != "" {
		if _, err := uuid.Parse(polyclinicID); err != nil {
			_ = c.Error(apperror.BadRequest("polyclinic_id harus berupa UUID yang valid", err))
			return
		}
	}

	hub.ServeWS(h.hub, c.Writer, c.Request, polyclinicID)
}

// GetPolyclinicDisplay provides polling fallback data for TV queue displays.
// GET /api/v1/public/queues/:id/display (or /api/v1/public/queues/:polyclinic_id/display)
func (h *WSHandler) GetPolyclinicDisplay(c *gin.Context) {
	polyclinicIDStr := c.Param("id")
	if polyclinicIDStr == "" {
		polyclinicIDStr = c.Param("polyclinic_id")
	}

	polyclinicID, err := uuid.Parse(polyclinicIDStr)
	if err != nil {
		_ = c.Error(apperror.BadRequest("polyclinic_id harus berupa UUID yang valid", err))
		return
	}

	ctx := c.Request.Context()

	// 1. Fetch polyclinic info
	var polyName, polyCode string
	err = h.pool.QueryRow(ctx, "SELECT name, queue_prefix FROM polyclinics WHERE id = $1", polyclinicID).Scan(&polyName, &polyCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_ = c.Error(apperror.NotFound("Poliklinik", nil))
			return
		}
		_ = c.Error(apperror.Internal(fmt.Errorf("WSHandler.GetPolyclinicDisplay: %w", err)))
		return
	}

	// 2. Fetch active/current queue today ('in_service' prioritised over 'called')
	activeQuery := `
		SELECT id, queue_number, queue_code, status, called_at, service_end
		FROM queues
		WHERE polyclinic_id = $1 
		  AND queue_date = CURRENT_DATE 
		  AND status IN ('called', 'in_service')
		ORDER BY CASE status WHEN 'in_service' THEN 1 ELSE 2 END, called_at DESC NULLS LAST
		LIMIT 1
	`
	var currentQueue *dto.QueueDisplayItem
	var qID uuid.UUID
	var qNum int
	var qCode, qStatus string
	var qCalledAt, qServiceEnd *time.Time

	err = h.pool.QueryRow(ctx, activeQuery, polyclinicID).Scan(&qID, &qNum, &qCode, &qStatus, &qCalledAt, &qServiceEnd)
	if err == nil {
		currentQueue = &dto.QueueDisplayItem{
			ID:          qID.String(),
			QueueNumber: qNum,
			QueueCode:   qCode,
			Status:      qStatus,
			CalledAt:    qCalledAt,
			ServiceEnd:  qServiceEnd,
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		_ = c.Error(apperror.Internal(fmt.Errorf("WSHandler.GetPolyclinicDisplay query active: %w", err)))
		return
	}

	// 3. Fetch recent 5 done queues today
	doneQuery := `
		SELECT id, queue_number, queue_code, status, called_at, service_end
		FROM queues
		WHERE polyclinic_id = $1 
		  AND queue_date = CURRENT_DATE 
		  AND status = 'done'
		ORDER BY service_end DESC NULLS LAST
		LIMIT 5
	`
	rows, err := h.pool.Query(ctx, doneQuery, polyclinicID)
	if err != nil {
		_ = c.Error(apperror.Internal(fmt.Errorf("WSHandler.GetPolyclinicDisplay query recent done: %w", err)))
		return
	}
	defer rows.Close()

	recentDone := make([]dto.QueueDisplayItem, 0)
	for rows.Next() {
		var item dto.QueueDisplayItem
		var itemID uuid.UUID
		if err := rows.Scan(&itemID, &item.QueueNumber, &item.QueueCode, &item.Status, &item.CalledAt, &item.ServiceEnd); err != nil {
			_ = c.Error(apperror.Internal(fmt.Errorf("WSHandler.GetPolyclinicDisplay scan recent done: %w", err)))
			return
		}
		item.ID = itemID.String()
		recentDone = append(recentDone, item)
	}

	// 4. Count waiting and checked_in queues today
	var waitingCount int
	err = h.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM queues
		WHERE polyclinic_id = $1 
		  AND queue_date = CURRENT_DATE 
		  AND status IN ('waiting', 'checked_in')
	`, polyclinicID).Scan(&waitingCount)
	if err != nil {
		_ = c.Error(apperror.Internal(fmt.Errorf("WSHandler.GetPolyclinicDisplay count waiting: %w", err)))
		return
	}

	// 5. Count total queues today
	var totalToday int
	err = h.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM queues
		WHERE polyclinic_id = $1 
		  AND queue_date = CURRENT_DATE
	`, polyclinicID).Scan(&totalToday)
	if err != nil {
		_ = c.Error(apperror.Internal(fmt.Errorf("WSHandler.GetPolyclinicDisplay count total: %w", err)))
		return
	}

	resp := dto.PolyclinicDisplayResponse{
		PolyclinicID:   polyclinicID.String(),
		PolyclinicName: polyName,
		PolyclinicCode: polyCode,
		CurrentQueue:   currentQueue,
		RecentDone:     recentDone,
		WaitingCount:   waitingCount,
		TotalToday:     totalToday,
		ServerTime:     time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response.OK(resp))
}
