# PBI-10: Traffic Generation

[View in Backlog](../backlog.md)

## Overview

Integrate k6 as a traffic generation tool so users can send configurable load through their deployed system architecture and observe how services handle requests under pressure.

## Problem Statement

Deployed containers are running but idle. The core educational value of the platform comes from seeing how architectures perform under load. Without a traffic generator, users can't test their designs or learn about bottlenecks, latency, and throughput.

## User Stories

- As a user, I want to start a load test against my deployed system so that I can see how it handles traffic
- As a user, I want to configure traffic patterns (requests per second, duration) so that I can simulate different scenarios
- As a user, I want to stop a running load test so that I can make changes and re-test

## Technical Approach

- Create a k6 wrapper Docker image with pre-configured scripts that can be parameterised
- Backend generates k6 test scripts dynamically based on diagram topology (which endpoints to hit, in what order)
- Traffic routes follow diagram connections: k6 → entry point service → downstream services
- REST endpoints: `POST /api/traffic/start`, `POST /api/traffic/stop`, `GET /api/traffic/config`
- Frontend controls: Start/Stop traffic buttons, load profile selector (low/medium/high or custom RPS + duration)
- k6 container runs on the shared Docker network so it can reach all services

## UX/UI Considerations

- Traffic controls in the toolbar or a dedicated panel
- Load profile presets (e.g., "Gentle: 10 RPS for 30s", "Moderate: 100 RPS for 60s", "Stress: 500 RPS for 120s")
- Custom configuration for advanced users (RPS, duration, ramp-up)
- Visual indicator that traffic is flowing (e.g., animated edges on the canvas)
- Running test status displayed (elapsed time, total requests sent)

## Acceptance Criteria

1. k6 container can be started with a dynamically generated test script
2. Traffic routes through the deployed services following diagram connections
3. `POST /api/traffic/start` initiates a load test with specified parameters
4. `POST /api/traffic/stop` terminates a running load test
5. Frontend provides start/stop controls and load profile selection
6. At least 3 preset load profiles are available (low, medium, high)

## Dependencies

- **Depends on**: PBI-9 (deploy flow — containers must be running to receive traffic)
- **External**: k6 Docker image

## Open Questions

- Should traffic patterns be configurable or fixed profiles? (PRD open question — recommend presets + custom)
- How to determine the entry point service for traffic? (user selection or auto-detect from diagram topology)
- Should we show k6 output (summary stats) directly or rely on Prometheus metrics (PBI-11)?

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 10`._
