package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func setupHubWSServer(t *testing.T) (*httptest.Server, *StatusHub) {
	t.Helper()
	t.Setenv("CORS_ORIGIN", "")

	hub := NewStatusHub()
	h := NewWebSocketHandler(hub)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	return httptest.NewServer(mux), hub
}

func hubWSURL(server *httptest.Server) string {
	return "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/status"
}

func TestStatusHub_RegisterAndBroadcast(t *testing.T) {
	server, hub := setupHubWSServer(t)
	defer server.Close()

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(hubWSURL(server), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// Wait for registration.
	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount() != 1 {
		t.Errorf("expected 1 client, got %d", hub.ClientCount())
	}

	// Broadcast a message.
	msg := map[string]string{"type": "test", "data": "hello"}
	hub.Broadcast(msg)

	// Read the message.
	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatalf("set read deadline: %v", err)
	}
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read message: %v", err)
	}

	var received map[string]string
	if err := json.Unmarshal(data, &received); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if received["type"] != "test" || received["data"] != "hello" {
		t.Errorf("unexpected message: %v", received)
	}
}

func TestStatusHub_MultipleClients(t *testing.T) {
	server, hub := setupHubWSServer(t)
	defer server.Close()

	dialer := websocket.Dialer{}
	conn1, _, err := dialer.Dial(hubWSURL(server), nil)
	if err != nil {
		t.Fatalf("dial 1: %v", err)
	}
	defer func() { _ = conn1.Close() }()

	conn2, _, err := dialer.Dial(hubWSURL(server), nil)
	if err != nil {
		t.Fatalf("dial 2: %v", err)
	}
	defer func() { _ = conn2.Close() }()

	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount() != 2 {
		t.Errorf("expected 2 clients, got %d", hub.ClientCount())
	}

	msg := map[string]string{"type": "broadcast"}
	hub.Broadcast(msg)

	for i, conn := range []*websocket.Conn{conn1, conn2} {
		if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
			t.Fatalf("set read deadline %d: %v", i, err)
		}
		_, data, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read message %d: %v", i, err)
		}
		var received map[string]string
		if err := json.Unmarshal(data, &received); err != nil {
			t.Fatalf("unmarshal %d: %v", i, err)
		}
		if received["type"] != "broadcast" {
			t.Errorf("client %d: unexpected message: %v", i, received)
		}
	}
}

func TestStatusHub_UnregisterOnDisconnect(t *testing.T) {
	server, hub := setupHubWSServer(t)
	defer server.Close()

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(hubWSURL(server), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if hub.ClientCount() != 1 {
		t.Errorf("expected 1 client, got %d", hub.ClientCount())
	}

	// Close connection
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	_ = conn.WriteMessage(websocket.CloseMessage, msg)
	_ = conn.Close()

	// Wait for server to process close.
	time.Sleep(100 * time.Millisecond)

	if hub.ClientCount() != 0 {
		t.Errorf("expected 0 clients after disconnect, got %d", hub.ClientCount())
	}
}
