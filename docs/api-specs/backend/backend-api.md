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

### Deploy Diagram

```http
POST /api/deploy
Content-Type: application/json
```

Request body: Full `Diagram` JSON (same schema as diagram CRUD).

Response `202 Accepted`:

```json
{
  "status": "deploying"
}
```

Error `400 Bad Request`:

```json
{
  "error": "validation failed: ..."
}
```

Error `409 Conflict` (already deploying/deployed):

```json
{
  "error": "deployment already in progress"
}
```

### Update Deployment (Live Topology Update)

```http
PUT /api/deploy
Content-Type: application/json
```

Request body: Full `Diagram` JSON with updated topology. The backend diffs the new nodes against the last-deployed diagram and applies only the changes (creates new containers for added nodes, removes containers for removed nodes).

Response `200 OK`:

```json
{
  "deployStatus": "deployed",
  "nodeStatuses": [
    {
      "nodeId": "abc",
      "containerId": "xyz",
      "status": "running"
    }
  ]
}
```

Error `400 Bad Request`:

```json
{
  "error": "validation failed: ..."
}
```

Error `409 Conflict` (not deployed):

```json
{
  "error": "no active deployment"
}
```

### Teardown Deployment

```http
DELETE /api/deploy
```

Response `200 OK`:

```json
{
  "status": "idle"
}
```

Error `409 Conflict` (not deployed):

```json
{
  "error": "no active deployment"
}
```

### Get Deploy Status

```http
GET /api/deploy/status
```

Response `200 OK`:

```json
{
  "deployStatus": "deployed",
  "nodeStatuses": [
    {
      "nodeId": "abc",
      "containerId": "xyz",
      "status": "running"
    }
  ]
}
```

`deployStatus` values: `idle`, `deploying`, `deployed`, `tearing_down`, `error`.

`status` values (per node): `created`, `running`, `stopped`, `error`, `healthy`, `unhealthy`.

## WebSocket Endpoints

### Status Stream

```
/ws/status
```

Upgrades HTTP connection to WebSocket. Broadcasts real-time container status updates to all connected clients.

- **Origin check**: Must match `CORS_ORIGIN` (or be empty)
- **Keep-alive**: Server sends periodic pings; client must respond with pongs
- **Broadcast messages**: JSON `StatusMessage` pushed to all clients when container health changes
- **Non-WebSocket requests**: Returns `400 Bad Request`

Message format:

```json
{
  "type": "status_update",
  "deployStatus": "deployed",
  "nodeStatuses": [
    { "nodeId": "abc", "containerId": "xyz", "status": "running" }
  ]
}
```

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
