# Product Backlog

**PRD**: [View PRD](../prd.md)

## Backlog Items

| ID | Actor | User Story | Status | Conditions of Satisfaction (CoS) |
|----|-------|-----------|--------|----------------------------------|
| 1 | Developer | As a developer, I want the project scaffolded with Next.js frontend, Go backend, and Docker dev environment so that I have a working foundation to build on | Done | Monorepo structure created; Next.js app runs on dev server; Go backend compiles and serves HTTP; Docker Compose dev environment starts all services; basic CI lint/build passes |
| 2 | User | As a user, I want a visual canvas where I can drag, drop, and arrange system components so that I can design architectures visually | Done | React Flow canvas renders; nodes can be dragged from palette onto canvas; nodes can be repositioned; canvas supports zoom and pan; canvas state persists during session |
| 3 | User | As a user, I want a library of service components (API, PostgreSQL, Redis, Nginx, RabbitMQ) with configuration panels so that I can define service-specific settings | Done | 5 service types available in palette with distinct visuals; clicking a node opens config panel; API service config supports endpoint definition (method, path, response schema); config changes persist to diagram state |
| 4 | User | As a user, I want to draw connections between services and export the topology as JSON so that my diagram can be sent to the backend for deployment | Done | Edges can be drawn between nodes; connections have labels; diagram exports to JSON matching the PRD schema; JSON can be re-imported to restore diagram state |
| 5 | Developer | As a developer, I want a Go backend with REST API endpoints for diagram CRUD and WebSocket scaffolding so that the frontend has a backend to communicate with | Done | Go HTTP server starts and listens; diagram CRUD endpoints (POST/GET/PUT) functional; diagrams persisted to storage; WebSocket endpoint established and accepts connections |
| 6 | Developer | As a developer, I want a Docker orchestration engine using Docker SDK so that I can programmatically manage container lifecycles | Done | Docker SDK integrated in Go backend; containers can be created, started, stopped, removed; single shared Docker network created and managed; health checks report container status |
| 7 | Developer | As a developer, I want service-type-to-container-template mappings so that each diagram node type translates to the correct Docker container configuration | Done | Container templates exist for all 5 service types; templates produce correct images/configs; service dependency ordering enforced at startup; port allocation avoids conflicts |
| 8 | Developer | As a developer, I want a mock API system that generates OpenAPI specs from user-defined endpoints and serves them via Prism so that API services return realistic mock data | Done | Endpoint definitions translate to valid OpenAPI 3.0 specs; Prism container reads mounted spec and serves endpoints; mock responses contain realistic fake data; cross-service HTTP calls resolve within Docker network |
| 9 | User | As a user, I want to click Deploy and see my diagram become running containers with real-time status updates so that I can interact with my architecture | Agreed | Deploy button sends topology to backend; backend spins up containers matching diagram; real-time container status shown via WebSocket; live updates work (add/remove node without full teardown); teardown cleans up all containers |
| 10 | User | As a user, I want to generate traffic through my system and configure load profiles so that I can see how my architecture handles requests | Proposed | k6 integrated as traffic generator; load test starts/stops from frontend; traffic routes through diagram connections; configurable request patterns and load profiles |
| 11 | User | As a user, I want real-time metrics (latency, throughput, error rate, resource usage) for each service so that I can evaluate my design's performance | Proposed | Prometheus collects metrics from all containers; per-service latency, throughput, error rate displayed; CPU/memory usage per container shown; metrics refresh at 1-second intervals; metrics stream via WebSocket |

_Items are ordered by priority (highest first)._

## PBI Details

| ID | Title | Detail Document |
|----|-------|----------------|
| 1 | Project Foundation & Dev Environment | [View Details](./1/prd.md) |
| 2 | Diagram Canvas | [View Details](./2/prd.md) |
| 3 | Service Component Library & Configuration | [View Details](./3/prd.md) |
| 4 | Connections & Topology Export | [View Details](./4/prd.md) |
| 5 | Go Backend API | [View Details](./5/prd.md) |
| 6 | Docker Orchestration Engine | [View Details](./6/prd.md) |
| 7 | Service-to-Container Mapping | [View Details](./7/prd.md) |
| 8 | Mock API System | [View Details](./8/prd.md) |
| 9 | Deploy Flow Integration | [View Details](./9/prd.md) |
| 10 | Traffic Generation | [View Details](./10/prd.md) |
| 11 | Observability & Metrics Dashboard | [View Details](./11/prd.md) |

## History

| Timestamp | PBI_ID | Event_Type | Details | User |
|-----------|--------|------------|---------|------|
| 20260218-005152 | ALL | Created | Initial backlog created from PRD decomposition | AI_Agent |
| 20260218-011452 | 1 | Status Change | Proposed → Agreed. User approved for planning. | User |
| 20260218-025856 | 1 | Status Change | Agreed → InReview. Tasks 1-1 through 1-5 completed; task 1-6 remains Blocked pending Docker-enabled AC4 verification. | AI_Agent |
| 20260218-060325 | 1 | Status Change | InReview → Done. Task 1-6 completed after end-to-end verification, including Docker Compose startup and service health checks. | AI_Agent |
| 20260218-073326 | 2 | Status Change | Proposed → Agreed. Auto-approved for planning. | AI_Agent |
| 2026-02-18 07:40:51 | 2 | Status Change | Agreed → InProgress. Started implementation. | AI_Agent |
| 2026-02-18 08:56:00 | 2 | Status Change | InProgress → Done. All PBI-2 tasks completed and verified. | AI_Agent |
| 2026-02-18 09:47:52 | 3 | Status Change | Proposed → Agreed. Auto-approved for planning. | AI_Agent |
| 2026-02-18 09:54:38 | 3 | Status Change | Agreed → InProgress. Started implementation. | AI_Agent |
| 2026-02-18 10:27:34 | 3 | Status Change | InProgress → Done. All PBI-3 tasks completed and verified. | AI_Agent |
| 2026-02-18 11:49:11 | 4 | Status Change | Proposed → Agreed. Auto-approved for planning. | AI_Agent |
| 2026-02-18 11:53:18 | 4 | Status Change | Agreed → InProgress. Started implementation. | AI_Agent |
| 2026-02-18 13:39:47 | 4 | Status Change | InProgress → Done. All PBI-4 tasks completed and verified. | AI_Agent |
| 2026-02-19 04:38:28 | 5 | Status Change | Proposed → Agreed. Auto-approved for planning. | AI_Agent |
| 2026-02-19 04:41:13 | 5 | Status Change | Agreed → InProgress. Started implementation. | AI_Agent |
| 2026-02-19 05:32:09 | 5 | Status Change | InProgress → Done. All PBI-5 tasks completed and verified. | AI_Agent |
| 2026-02-19 06:18:07 | 6 | Status Change | Proposed → Agreed. Auto-approved for planning. | AI_Agent |
| 2026-02-19 06:53:01 | 6 | Status Change | Agreed → InProgress. Started implementation. | AI_Agent |
| 2026-02-19 07:32:50 | 6 | Status Change | InProgress → Done. All 7 tasks completed and verified. | AI_Agent |
| 2026-02-19 08:11:26 | 7 | Status Change | Proposed → Agreed. Auto-approved for planning. | AI_Agent |
| 2026-02-19 08:16:59 | 7 | Status Change | Agreed → InProgress. Started implementation. | AI_Agent |
| 2026-02-19 08:45:58 | 7 | Status Change | InProgress → Done. All 7 tasks completed and verified. | AI_Agent |
| 2026-02-19 10:19:37 | 8 | Status Change | Proposed → Agreed. Auto-approved for planning. | AI_Agent |
| 2026-02-19 13:36:02 | 8 | Status Change | Agreed → InProgress. Started implementation. | AI_Agent |
| 2026-02-20 05:15:09 | 8 | Status Change | InProgress → Done. All 4 tasks completed and verified. | AI_Agent |
| 2026-02-20 08:04:21 | 9 | Status Change | Proposed → Agreed. Auto-approved for planning. | AI_Agent |
