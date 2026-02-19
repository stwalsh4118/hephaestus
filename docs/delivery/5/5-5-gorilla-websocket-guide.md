# gorilla/websocket Guide

**Date**: 2026-02-19
**Docs**: https://pkg.go.dev/github.com/gorilla/websocket

## Installation

```bash
go get github.com/gorilla/websocket
```

## Key Patterns for Task 5-5

### Upgrader

```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Validate against allowed origin
        return true
    },
}
```

### Upgrade HTTP to WebSocket

```go
conn, err := upgrader.Upgrade(w, r, nil)
if err != nil {
    log.Println(err)
    return
}
defer conn.Close()
```

### Ping/Pong

Built-in: gorilla/websocket automatically responds to pings with pongs by default.
Custom handler:

```go
conn.SetPongHandler(func(appData string) error {
    log.Println("pong received")
    return nil
})
```

### Read Loop (required for ping/pong to work)

```go
for {
    _, _, err := conn.ReadMessage()
    if err != nil {
        if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
            log.Printf("unexpected close: %v", err)
        }
        break
    }
}
```

### Close

```go
msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))
```

### Concurrency

- One reader goroutine, one writer goroutine max.
- `WriteControl` is safe to call concurrently.
