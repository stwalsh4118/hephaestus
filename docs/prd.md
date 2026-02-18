# Visual System Design Simulator - Product Requirements Document

## Overview

A visual infrastructure-as-code platform that bridges the gap between system design theory and practice. Users draw system architectures using a drag-and-drop interface, and the platform translates those diagrams into real running containers with mock APIs, allowing users to test data flow, monitor performance, and simulate scaling scenarios.

## Problem Statement

System design learning is largely theoretical—whiteboard diagrams don't run, so learners never see how their designs actually perform under load. This platform makes diagrams executable, providing immediate feedback on architectural decisions.

## Target Users

- Software engineers preparing for system design interviews
- Developers learning distributed systems concepts
- Teams prototyping architectures before implementation

## Key Differentiators

1. **Real Infrastructure**: Unlike draw.io diagrams, these actually run
2. **Learning Focus**: Built specifically for system design education
3. **Immediate Feedback**: See how designs perform under load instantly
4. **Cost Awareness**: Show real resource usage (stretch goal)

---

## Technical Architecture

### Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Next.js + React Flow |
| Backend | Go |
| Container Runtime | Docker (via Docker SDK) |
| Mock Server | Prism or WireMock |
| Metrics | Prometheus |
| Load Testing | k6 or Artillery |
| Networking | Single Docker network (simplified) |

### Core Components

#### 1. Diagram Engine (Frontend)

**Purpose**: Visual editor and diagram parser

**Features**:
- Drag-and-drop canvas using React Flow
- Component library with pre-built service types
- Connection drawing (arrows) for data flow
- Service configuration panel (name, endpoints, ports)
- JSON export of diagram topology

**Service Types for MVP**:
- API Service (generic, user-defined endpoints)
- Database (PostgreSQL)
- Cache (Redis)
- Load Balancer (Nginx)
- Message Queue (RabbitMQ)

**Extended Service Types** (post-MVP):
- Object Storage
- CDN
- Search (Elasticsearch)
- Serverless Functions
- Auth Service
- Kafka Streams
- Pub/Sub

#### 2. Infrastructure Translator (Backend)

**Purpose**: Convert JSON diagrams into running containers

**Features**:
- Parse JSON diagram into Go structs
- Map service types to container configurations
- Generate Docker network configuration
- Handle service dependencies and startup order

**Implementation Details**:
- Use Docker SDK for Go (not Docker Compose) to enable live updates
- Single shared Docker network for all containers (simplified networking)
- Mount mock server configs as volumes

#### 3. Mock API Generator

**Purpose**: Enable user-defined endpoints without business logic

**Features**:
- Users define endpoints via UI (method, path, response schema)
- Backend generates OpenAPI specs from definitions
- Prism/WireMock container reads spec and serves mock responses
- JSON Schema Faker generates realistic response data

**Example Flow**:
```
User defines: POST /users → { "id": "uuid", "name": "string" }
    ↓
Backend generates OpenAPI spec
    ↓
Spec mounted into Prism container
    ↓
Prism serves endpoint with fake data
```

#### 4. Orchestration Engine

**Purpose**: Deploy and manage containers with live updates

**Features**:
- Spin up containers on demand via Docker SDK
- Live updates: add/remove/modify containers without teardown
- Health checks for all services
- Service discovery via Docker DNS

**Key Decision**: Use Docker SDK instead of Docker Compose to enable:
- Adding new containers without restarting everything
- Modifying individual services in isolation
- Real-time topology changes from the frontend

#### 5. Data Flow Engine

**Purpose**: Route real traffic through the system

**Features**:
- Traffic generator using k6 or Artillery
- Request routing based on diagram connections
- Support for different data formats (JSON primarily)
- Configurable request patterns and load profiles

#### 6. Observability Stack

**Purpose**: Monitor and visualize system behavior

**Features**:
- Prometheus for metrics collection
- Real-time metrics displayed in frontend
- Per-service latency, throughput, error rates
- Resource usage (CPU, memory) per container

**Stretch Goals**:
- Grafana dashboards
- Distributed tracing (Jaeger)
- Alerting on anomalies

#### 7. Scaling Simulator (Post-MVP)

**Purpose**: Simulate real-world scaling scenarios

**Features**:
- Horizontal scaling (spin up multiple instances)
- Load simulation (CPU, memory, I/O stress)
- Chaos engineering (inject failures)
- Auto-scaling based on metrics

---

## User Workflow

### 1. Design Phase
```
User opens canvas
    ↓
Drags "API Service" onto canvas, names it "User Service"
    ↓
Opens config panel, defines endpoints:
  - GET /users
  - POST /users
  - GET /users/:id
    ↓
Drags "Database" (PostgreSQL) onto canvas
    ↓
Draws arrow from User Service → Database
    ↓
Drags "Cache" (Redis) onto canvas
    ↓
Draws arrow from User Service → Cache
```

### 2. Deploy Phase
```
User clicks "Deploy" (or changes happen live)
    ↓
Frontend sends JSON topology to Go backend
    ↓
Backend spins up containers via Docker SDK:
  - user-service (Prism with OpenAPI spec)
  - postgres (PostgreSQL container)
  - redis (Redis container)
    ↓
All containers join shared Docker network
    ↓
Frontend shows "Running" status for each service
```

### 3. Test Phase
```
User clicks "Run Traffic"
    ↓
Backend triggers k6 load test
    ↓
Requests flow: k6 → User Service → Database/Cache
    ↓
Prometheus collects metrics
    ↓
Frontend displays real-time graphs:
  - Requests/sec
  - Latency (p50, p95, p99)
  - Error rate
```

### 4. Iterate Phase
```
User adds "Load Balancer" to diagram
    ↓
Backend spins up Nginx container (no teardown)
    ↓
User reconfigures traffic flow
    ↓
System updates routing live
```

---

## MVP Scope (8 weeks)

### Week 1-2: Diagram Engine
- [ ] Set up Next.js project with React Flow
- [ ] Implement drag-and-drop canvas
- [ ] Create component library (5 service types)
- [ ] Build connection drawing
- [ ] Export diagram as JSON

### Week 3-4: Infrastructure Translator
- [ ] Set up Go backend with Docker SDK
- [ ] Parse JSON diagram to Go structs
- [ ] Create container templates for each service type
- [ ] Implement single-network topology
- [ ] Build deploy endpoint

### Week 5-6: Mock API System
- [ ] Build endpoint definition UI in frontend
- [ ] Generate OpenAPI specs from definitions
- [ ] Integrate Prism as mock server
- [ ] Mount specs into containers
- [ ] Verify cross-service communication

### Week 7-8: Observability & Traffic
- [ ] Integrate Prometheus
- [ ] Build metrics display in frontend
- [ ] Set up k6 for traffic generation
- [ ] Create basic load test profiles
- [ ] End-to-end integration testing

---

## API Endpoints

### Backend REST API

```
POST   /api/diagrams              Create/save diagram
GET    /api/diagrams/:id          Get diagram
PUT    /api/diagrams/:id          Update diagram

POST   /api/deploy                Deploy diagram to containers
DELETE /api/deploy                Tear down all containers
GET    /api/deploy/status         Get status of all containers

POST   /api/traffic/start         Start load test
POST   /api/traffic/stop          Stop load test
GET    /api/traffic/config        Get traffic config

GET    /api/metrics               Get Prometheus metrics (proxied)

WebSocket /ws/status              Real-time container status updates
WebSocket /ws/metrics             Real-time metrics stream
```

---

## Data Models

### Diagram JSON Schema

```json
{
  "id": "uuid",
  "name": "My System Design",
  "nodes": [
    {
      "id": "node-1",
      "type": "api-service",
      "name": "User Service",
      "position": { "x": 100, "y": 200 },
      "config": {
        "endpoints": [
          { "method": "GET", "path": "/users", "responseSchema": {...} },
          { "method": "POST", "path": "/users", "responseSchema": {...} }
        ],
        "port": 8080
      }
    },
    {
      "id": "node-2",
      "type": "database",
      "name": "PostgreSQL",
      "position": { "x": 300, "y": 200 },
      "config": {
        "engine": "postgresql",
        "version": "15"
      }
    }
  ],
  "edges": [
    {
      "id": "edge-1",
      "source": "node-1",
      "target": "node-2",
      "label": "reads/writes"
    }
  ]
}
```

### Container State

```go
type ContainerState struct {
    ID          string            `json:"id"`
    NodeID      string            `json:"nodeId"`
    Name        string            `json:"name"`
    Status      string            `json:"status"` // running, stopped, error
    Health      string            `json:"health"` // healthy, unhealthy, starting
    Ports       map[string]string `json:"ports"`
    CreatedAt   time.Time         `json:"createdAt"`
}
```

---

## Non-Functional Requirements

### Performance
- Container startup: < 10 seconds per service
- UI responsiveness: < 100ms for drag operations
- Metrics refresh: 1-second intervals

### Scalability
- Support diagrams with up to 20 services
- Handle 1000 RPS in load tests

### Reliability
- Graceful handling of container failures
- Auto-recovery for crashed services
- Clear error messages in UI

---

## Future Enhancements (Post-MVP)

1. **Saved Templates**: Pre-built architectures (e.g., "Twitter clone", "E-commerce")
2. **Collaborative Editing**: Multiple users designing together
3. **Cost Estimation**: Show estimated cloud costs for designs
4. **Kubernetes Mode**: Deploy to real k8s cluster
5. **Chaos Engineering**: Built-in failure injection
6. **Performance Benchmarks**: Compare against known architectures
7. **Export to IaC**: Generate Terraform/Pulumi from diagrams

---

## Success Metrics

- User can go from blank canvas to running system in < 5 minutes
- Load test results display within 30 seconds of starting
- Zero-downtime updates when modifying topology
- Clear correlation between diagram changes and performance impact

---

## Open Questions

1. Should we support custom Docker images for advanced users?
2. How do we handle database persistence between sessions?
3. Should traffic patterns be configurable or fixed profiles?
4. Do we need authentication for the MVP?

---

## Repository Structure

```
visual-system-design-simulator/
├── frontend/                 # Next.js app
│   ├── components/
│   │   ├── Canvas/          # React Flow wrapper
│   │   ├── ServicePanel/    # Service configuration
│   │   ├── MetricsPanel/    # Real-time metrics
│   │   └── Toolbar/         # Actions (deploy, traffic)
│   ├── lib/
│   │   ├── api.ts           # Backend API client
│   │   └── diagram.ts       # Diagram state management
│   └── pages/
│       └── index.tsx        # Main canvas page
│
├── backend/                  # Go backend
│   ├── cmd/
│   │   └── server/          # Main entrypoint
│   ├── internal/
│   │   ├── api/             # HTTP handlers
│   │   ├── docker/          # Docker SDK wrapper
│   │   ├── translator/      # Diagram → Container logic
│   │   ├── mock/            # OpenAPI spec generation
│   │   └── metrics/         # Prometheus integration
│   └── templates/           # Container config templates
│
├── docker/                   # Base container images
│   ├── mock-server/         # Prism wrapper image
│   └── load-generator/      # k6 wrapper image
│
└── docs/
    └── prd.md               # This document
```
