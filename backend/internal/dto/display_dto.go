package dto

import "time"

// QueueDisplayItem represents a minimal queue representation on the TV display screen.
type QueueDisplayItem struct {
	ID          string     `json:"id"`
	QueueNumber int        `json:"queue_number"`
	QueueCode   string     `json:"queue_code"`
	Status      string     `json:"status"`
	CalledAt    *time.Time `json:"called_at,omitempty"`
	ServiceEnd  *time.Time `json:"service_end,omitempty"`
}

// PolyclinicDisplayResponse is the payload for GET /api/v1/public/queues/:polyclinic_id/display.
type PolyclinicDisplayResponse struct {
	PolyclinicID   string             `json:"polyclinic_id"`
	PolyclinicName string             `json:"polyclinic_name"`
	PolyclinicCode string             `json:"polyclinic_code"`
	CurrentQueue   *QueueDisplayItem  `json:"current_queue"`
	RecentDone     []QueueDisplayItem `json:"recent_done"`
	WaitingCount   int                `json:"waiting_count"`
	TotalToday     int                `json:"total_today"`
	ServerTime     string             `json:"server_time"`
}

// QueueEvent is the JSON envelope sent over WebSocket connections.
type QueueEvent struct {
	Event string `json:"event"` // e.g. "queue_updated", "queue_called"
	Data  any    `json:"data"`
}
