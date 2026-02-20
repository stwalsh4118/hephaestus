package handler

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsReadBufferSize  = 1024
	wsWriteBufferSize = 1024
	wsPongWait        = 60 * time.Second
	wsPingInterval    = (wsPongWait * 9) / 10
	wsWriteWait       = 10 * time.Second
	wsCORSOriginEnv   = "CORS_ORIGIN"
	wsDefaultOrigin   = "http://localhost:3000"
)

// WebSocketHandler handles WebSocket connections at /ws/status.
type WebSocketHandler struct {
	upgrader websocket.Upgrader
	hub      *StatusHub
}

// NewWebSocketHandler creates a WebSocketHandler with origin checking.
// If hub is nil, connections are accepted but no broadcasts are sent.
func NewWebSocketHandler(hub *StatusHub) *WebSocketHandler {
	origin := os.Getenv(wsCORSOriginEnv)
	if origin == "" {
		origin = wsDefaultOrigin
	}

	return &WebSocketHandler{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  wsReadBufferSize,
			WriteBufferSize: wsWriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return r.Header.Get("Origin") == "" || r.Header.Get("Origin") == origin
			},
		},
	}
}

// RegisterRoutes registers the WebSocket route on the given mux.
func (h *WebSocketHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ws/status", h.Handle)
}

// Handle upgrades an HTTP connection to WebSocket and maintains the connection.
func (h *WebSocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	log.Printf("websocket connection opened: %s", r.RemoteAddr)

	// Register with hub for broadcasts.
	if h.hub != nil {
		h.hub.Register(conn)
	}

	defer func() {
		if h.hub != nil {
			h.hub.Unregister(conn)
		}
		if err := conn.Close(); err != nil {
			log.Printf("websocket close error: %v", err)
		}
		log.Printf("websocket connection closed: %s", r.RemoteAddr)
	}()

	if err := conn.SetReadDeadline(time.Now().Add(wsPongWait)); err != nil {
		log.Printf("set read deadline: %v", err)
		return
	}
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(wsPongWait))
	})

	// Start a goroutine to send periodic pings.
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(wsPingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.SetWriteDeadline(time.Now().Add(wsWriteWait)); err != nil {
					return
				}
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Read loop: required for ping/pong handling. No business logic.
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("websocket unexpected close: %v", err)
			}
			break
		}
	}

	close(done)
}
