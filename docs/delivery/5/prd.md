# PBI-5: Go Backend API

[View in Backlog](../backlog.md)

## Overview

Build the Go backend HTTP server with REST API endpoints for diagram CRUD operations, persistent storage, and WebSocket scaffolding for real-time communication. This is the backend foundation that all server-side features depend on.

## Problem Statement

The frontend needs a backend to save/load diagrams and later to trigger deployments and stream status updates. Without REST endpoints and WebSocket support, the frontend operates in isolation with no persistence or server-side capabilities.

## User Stories

- As a user, I want to save my diagram so that I can return to it later
- As a user, I want to load a previously saved diagram so that I can continue working on it
- As a developer, I want WebSocket support scaffolded so that real-time features (status, metrics) can be built in later PBIs

## Technical Approach

- Go HTTP server using a lightweight router (chi or standard library)
- REST endpoints per PRD: `POST/GET/PUT /api/diagrams`, `GET /api/diagrams/:id`
- Storage: file-based JSON storage initially (simple, no external DB dependency for MVP)
- Diagram validation: ensure incoming JSON matches expected schema
- WebSocket endpoint scaffolded at `/ws/status` — accepts connections, basic ping/pong, no business logic yet
- CORS configuration for frontend dev server
- Structured logging

## UX/UI Considerations

N/A — backend PBI.

## Acceptance Criteria

1. Go HTTP server starts and listens on a configurable port
2. `POST /api/diagrams` creates a new diagram and returns its ID
3. `GET /api/diagrams/:id` retrieves a saved diagram
4. `PUT /api/diagrams/:id` updates an existing diagram
5. Diagrams are persisted to storage and survive server restarts
6. WebSocket endpoint at `/ws/status` accepts connections and responds to ping
7. CORS configured to allow frontend origin

## Dependencies

- **Depends on**: PBI-1 (project foundation)
- **External**: None beyond Go standard library / chosen router

## Open Questions

- chi router vs standard library ServeMux (Go 1.22+ has improved routing)?
- File-based storage vs SQLite for diagram persistence?
- Should diagram listing (`GET /api/diagrams`) be included now or deferred?

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 5`._
