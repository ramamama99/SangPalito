package hub

import (
	"encoding/json"
	"log/slog"
	"sync"

	"sangpalito/backend/internal/dto"
)

// BroadcastMessage encapsulates a room-targeted payload.
type BroadcastMessage struct {
	PolyclinicID string
	Payload      []byte // JSON
}

// Hub maintains the set of active clients across rooms (polyclinics) and broadcasts messages.
type Hub struct {
	// rooms: key = polyclinic_id, value = set of clients
	// an empty key "" holds clients subscribing to global/all polyclinic events
	rooms      map[string]map[*Client]bool
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		broadcast:  make(chan BroadcastMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run executes the main event loop for the Hub.
func (h *Hub) Run() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("websocket hub Run panic recovered", "panic", r)
		}
	}()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			room := client.polyclinicID
			if h.rooms[room] == nil {
				h.rooms[room] = make(map[*Client]bool)
			}
			h.rooms[room][client] = true
			h.mu.Unlock()
			slog.Debug("client registered to websocket room", "room", room)

		case client := <-h.unregister:
			h.mu.Lock()
			room := client.polyclinicID
			if clients, ok := h.rooms[room]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 && room != "" {
						delete(h.rooms, room)
					}
				}
			}
			h.mu.Unlock()
			slog.Debug("client unregistered from websocket room", "room", room)

		case msg := <-h.broadcast:
			h.mu.RLock()
			if msg.PolyclinicID != "" {
				// 1. Deliver to clients in this specific polyclinic room
				if clients, ok := h.rooms[msg.PolyclinicID]; ok {
					for client := range clients {
						select {
						case client.send <- msg.Payload:
						default:
							close(client.send)
							delete(clients, client)
						}
					}
				}
				// 2. Also deliver to global subscribers (room "")
				if globalClients, ok := h.rooms[""]; ok {
					for client := range globalClients {
						select {
						case client.send <- msg.Payload:
						default:
							close(client.send)
							delete(globalClients, client)
						}
					}
				}
			} else {
				// Broadcast to all clients across all rooms
				for _, clients := range h.rooms {
					for client := range clients {
						select {
						case client.send <- msg.Payload:
						default:
							close(client.send)
							delete(clients, client)
						}
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast dispatches a message conforming to service.QueueBroadcaster interface.
func (h *Hub) Broadcast(payload any) {
	if h == nil {
		return
	}

	var polyclinicID string
	var rawData []byte

	switch v := payload.(type) {
	case BroadcastMessage:
		select {
		case h.broadcast <- v:
		default:
			slog.Warn("websocket broadcast channel full, dropping message")
		}
		return
	case *BroadcastMessage:
		select {
		case h.broadcast <- *v:
		default:
			slog.Warn("websocket broadcast channel full, dropping message")
		}
		return
	case map[string]any:
		if data, ok := v["data"]; ok {
			if resp, ok := data.(*dto.QueueResponse); ok && resp != nil {
				polyclinicID = resp.PolyclinicID
			} else if m, ok := data.(map[string]any); ok {
				if pid, ok := m["polyclinic_id"].(string); ok {
					polyclinicID = pid
				}
			}
		}
	}

	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("websocket marshal error", "err", err)
		return
	}
	rawData = data

	select {
	case h.broadcast <- BroadcastMessage{PolyclinicID: polyclinicID, Payload: rawData}:
	default:
		slog.Warn("websocket broadcast channel full, dropping message")
	}
}

// BroadcastToPolyclinic sends a formatted event specifically to clients of a polyclinic.
func (h *Hub) BroadcastToPolyclinic(polyclinicID string, event string, data any) {
	if h == nil {
		return
	}
	evt := dto.QueueEvent{
		Event: event,
		Data:  data,
	}
	bytes, err := json.Marshal(evt)
	if err != nil {
		slog.Error("websocket BroadcastToPolyclinic marshal error", "err", err)
		return
	}
	select {
	case h.broadcast <- BroadcastMessage{PolyclinicID: polyclinicID, Payload: bytes}:
	default:
		slog.Warn("websocket broadcast channel full, dropping message")
	}
}
