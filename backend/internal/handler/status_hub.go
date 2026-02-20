package handler

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// StatusHub manages connected WebSocket clients and broadcasts messages to all.
type StatusHub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]struct{}
}

// NewStatusHub creates a new StatusHub.
func NewStatusHub() *StatusHub {
	return &StatusHub{
		clients: make(map[*websocket.Conn]struct{}),
	}
}

// Register adds a WebSocket connection to the hub.
func (h *StatusHub) Register(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = struct{}{}
}

// Unregister removes a WebSocket connection from the hub.
func (h *StatusHub) Unregister(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

// Broadcast sends a JSON message to all connected clients.
// Clients that fail to receive the message are removed from the hub.
func (h *StatusHub) Broadcast(msg any) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("status hub: marshal broadcast message: %v", err)
		return
	}

	h.mu.RLock()
	clients := make([]*websocket.Conn, 0, len(h.clients))
	for c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	var failed []*websocket.Conn
	for _, c := range clients {
		if err := c.SetWriteDeadline(time.Now().Add(wsWriteWait)); err != nil {
			failed = append(failed, c)
			continue
		}
		if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
			failed = append(failed, c)
		}
	}

	if len(failed) > 0 {
		h.mu.Lock()
		for _, c := range failed {
			delete(h.clients, c)
		}
		h.mu.Unlock()
	}
}

// ClientCount returns the number of connected clients.
func (h *StatusHub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
