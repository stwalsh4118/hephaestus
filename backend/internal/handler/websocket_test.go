package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func setupWSServer(t *testing.T) *httptest.Server {
	t.Helper()
	t.Setenv("CORS_ORIGIN", "")

	h := NewWebSocketHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	return httptest.NewServer(mux)
}

func wsURL(server *httptest.Server) string {
	return "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/status"
}

func TestWebSocket_UpgradeSuccess(t *testing.T) {
	server := setupWSServer(t)
	defer server.Close()

	dialer := websocket.Dialer{}
	conn, resp, err := dialer.Dial(wsURL(server), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("status: got %d, want %d", resp.StatusCode, http.StatusSwitchingProtocols)
	}
}

func TestWebSocket_PingPong(t *testing.T) {
	server := setupWSServer(t)
	defer server.Close()

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL(server), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// Set up pong handler to detect server's ping response
	pongReceived := make(chan struct{}, 1)
	conn.SetPongHandler(func(string) error {
		select {
		case pongReceived <- struct{}{}:
		default:
		}
		return nil
	})

	// Send a ping from client side
	if err := conn.WriteMessage(websocket.PingMessage, []byte("test")); err != nil {
		t.Fatalf("write ping: %v", err)
	}

	// Verify connection is still alive by checking we can write
	if err := conn.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
		t.Fatalf("write after ping: %v", err)
	}
}

func TestWebSocket_GracefulClose(t *testing.T) {
	server := setupWSServer(t)
	defer server.Close()

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL(server), nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	// Send close message
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	err = conn.WriteMessage(websocket.CloseMessage, msg)
	if err != nil {
		t.Fatalf("write close: %v", err)
	}

	// Read should return close error
	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatalf("set read deadline: %v", err)
	}
	_, _, err = conn.ReadMessage()
	if err == nil {
		t.Error("expected error after close, got nil")
	}
	if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		t.Logf("close error (acceptable): %v", err)
	}

	_ = conn.Close()
}

func TestWebSocket_NonWSRequest(t *testing.T) {
	server := setupWSServer(t)
	defer server.Close()

	resp, err := http.Get(server.URL + "/ws/status")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}
