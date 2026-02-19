# Backend API Specification

## Runtime Configuration

```text
PORT (optional): HTTP listen port, defaults to 8080
CORS_ORIGIN (optional): Allowed CORS origin, defaults to http://localhost:3000
```

## REST Endpoints

### Health Check

```http
GET /health
```

Response `200 OK`:

```json
{
  "status": "ok"
}
```

### Create Diagram

```http
POST /api/diagrams
Content-Type: application/json
```

Request body: `Diagram` JSON (see schema below). The `id` field in the request is ignored; a UUID is generated server-side.

Response `201 Created`:

```json
{
  "id": "<uuid>"
}
```

Error `400 Bad Request`:

```json
{
  "error": "validation failed: ..."
}
```

### Get Diagram

```http
GET /api/diagrams/{id}
```

Response `200 OK`: Full `Diagram` JSON.

Error `404 Not Found`:

```json
{
  "error": "diagram not found"
}
```

### Update Diagram

```http
PUT /api/diagrams/{id}
Content-Type: application/json
```

Request body: `Diagram` JSON. The `id` in the URL takes precedence.

Response `200 OK`: Updated `Diagram` JSON.

Error `404 Not Found`:

```json
{
  "error": "diagram not found"
}
```

Error `400 Bad Request`:

```json
{
  "error": "validation failed: ..."
}
```

## WebSocket Endpoints

### Status Stream

```
/ws/status
```

Upgrades HTTP connection to WebSocket. Scaffold only — no business logic messages yet.

- **Origin check**: Must match `CORS_ORIGIN` (or be empty)
- **Keep-alive**: Server sends periodic pings; client must respond with pongs
- **Purpose**: Future real-time status updates (PBI-9, PBI-11)
- **Non-WebSocket requests**: Returns `400 Bad Request`

## Diagram Schema

```go
type Diagram struct {
    ID    string        `json:"id"`
    Name  string        `json:"name"`
    Nodes []DiagramNode `json:"nodes"`
    Edges []DiagramEdge `json:"edges"`
}

type DiagramNode struct {
    ID          string          `json:"id"`
    Type        string          `json:"type"`       // api-service | postgresql | redis | nginx | rabbitmq
    Name        string          `json:"name"`
    Description string          `json:"description"`
    Position    *Position       `json:"position"`
    Config      json.RawMessage `json:"config,omitempty"`
}

type DiagramEdge struct {
    ID     string `json:"id"`
    Source string `json:"source"`
    Target string `json:"target"`
    Label  string `json:"label"`
}

type Position struct {
    X float64 `json:"x"`
    Y float64 `json:"y"`
}
```

## Storage

Diagrams are persisted as individual JSON files in `./data/diagrams/<id>.json`.
