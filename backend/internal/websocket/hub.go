package websocket

import (
	"encoding/json"
	"log/slog"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients    map[chan []byte]bool
	broadcast  chan []byte
	register   chan chan []byte
	unregister chan chan []byte
	mu         sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[chan []byte]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan chan []byte),
		unregister: make(chan chan []byte),
	}
}

// Run starts the hub loop in a goroutine with panic recovery.
func (h *Hub) Run() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("websocket hub panic recovered", "panic", r)
		}
	}()

	for {
		select {
		case ch := <-h.register:
			h.mu.Lock()
			h.clients[ch] = true
			h.mu.Unlock()
		case ch := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[ch]; ok {
				delete(h.clients, ch)
				close(ch)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for ch := range h.clients {
				select {
				case ch <- message:
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast serializes payload to JSON and dispatches to connected clients.
func (h *Hub) Broadcast(payload any) {
	if h == nil {
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("websocket marshal error", "err", err)
		return
	}
	select {
	case h.broadcast <- data:
	default:
		slog.Warn("websocket broadcast channel full, dropping message")
	}
}
