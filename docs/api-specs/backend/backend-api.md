# Backend API Specification

## Runtime Configuration

```text
PORT (optional): HTTP listen port, defaults to 8080
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
