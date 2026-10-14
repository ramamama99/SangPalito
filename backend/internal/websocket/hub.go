package websocket

import (
	"net/http"

	"sangpalito/backend/internal/hub"
)

// Hub is an alias to hub.Hub for backward compatibility.
type Hub = hub.Hub

// NewHub delegates to hub.NewHub.
func NewHub() *Hub {
	return hub.NewHub()
}

// ServeWS delegates to hub.ServeWS.
func ServeWS(h *Hub, w http.ResponseWriter, r *http.Request) {
	hub.ServeWS(h, w, r, "")
}
