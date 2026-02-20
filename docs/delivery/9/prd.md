# PBI-9: Deploy Flow Integration

[View in Backlog](../backlog.md)

## Overview

Wire together the complete deploy pipeline: a Deploy button in the frontend sends the diagram topology to the backend, the backend orchestrates container creation, and real-time container status streams back to the frontend via WebSocket. Support live topology updates (adding/removing nodes without full teardown).

## Problem Statement

Individual backend pieces exist (API, orchestration, templates, mock generation) but there's no end-to-end flow from the frontend. The user has no way to trigger a deployment, see container status, or make live changes. This PBI integrates everything into a cohesive deploy experience.

## User Stories

- As a user, I want to click a Deploy button and see my diagram become running containers so that I can interact with my architecture
- As a user, I want to see real-time status for each service (starting, running, error) so that I know what's happening
- As a user, I want to add or remove services from my diagram while it's deployed so that I can iterate without restarting everything
- As a user, I want to tear down all containers when I'm done so that resources are freed

## Technical Approach

- Frontend: Deploy/Teardown buttons in the toolbar; per-node status badges (colour-coded)
- `POST /api/deploy` endpoint: receives diagram JSON, calls orchestration engine to deploy all containers
- `DELETE /api/deploy` endpoint: tears down all containers
- `GET /api/deploy/status` endpoint: returns current state of all containers
- WebSocket `/ws/status`: streams container status changes in real-time to the frontend
- Live update logic: compare current deployed state with new diagram, compute diff (added/removed/modified nodes), apply changes incrementally
- Error handling: surface container failures to the UI with clear messages

## UX/UI Considerations

- Deploy button prominent in the toolbar, changes to "Deploying..." during deployment
- Each node on the canvas shows a status indicator (green = running, yellow = starting, red = error, grey = stopped)
- Teardown button with confirmation
- Error messages displayed as toast notifications or inline on the affected node
- Status transitions animated for visual feedback

## Acceptance Criteria

1. Deploy button sends the current diagram topology to `POST /api/deploy`
2. Backend creates containers for all nodes and reports progress
3. Real-time container status updates appear on canvas nodes via WebSocket
4. Adding a new node to a deployed diagram deploys only the new container (no full teardown)
5. Removing a node from a deployed diagram removes only that container
6. Teardown button removes all containers and resets status indicators
7. Container errors are surfaced in the UI with actionable messages

## Dependencies

- **Depends on**: PBI-4 (topology export), PBI-5 (backend API + WebSocket), PBI-6 (orchestration), PBI-7 (container templates), PBI-8 (mock API for API service nodes)
- **External**: None

## Open Questions

- Should deploy be manual (button click) or automatic (live as diagram changes)?
- How to handle partial deploy failures (some containers start, others fail)?

## Related Tasks

[View Tasks](./tasks.md)
