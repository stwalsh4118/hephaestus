# PBI-11: Observability & Metrics Dashboard

[View in Backlog](../backlog.md)

## Overview

Integrate Prometheus for metrics collection from all deployed containers and build a real-time metrics dashboard in the frontend displaying per-service latency, throughput, error rates, and resource usage (CPU, memory).

## Problem Statement

Users can deploy containers and send traffic, but they can't see what's happening inside. The educational value depends on observability — seeing how latency increases under load, how error rates spike when a service is overwhelmed, and how resource usage correlates with traffic patterns.

## User Stories

- As a user, I want to see real-time latency, throughput, and error rates per service so that I can evaluate my architecture's performance
- As a user, I want to see CPU and memory usage per container so that I can understand resource consumption
- As a user, I want metrics to refresh continuously so that I can observe changes as traffic flows

## Technical Approach

- Deploy a Prometheus container on the shared Docker network, configured to scrape all service containers
- Expose container-level metrics via cAdvisor or Docker stats API (CPU, memory, network I/O)
- For Prism/mock services: collect HTTP request metrics (request count, latency histograms, error counts)
- Backend proxies Prometheus queries via `GET /api/metrics` endpoint
- WebSocket `/ws/metrics`: stream aggregated metrics to the frontend at 1-second intervals
- Frontend metrics panel: per-service cards or graphs showing latency (p50, p95, p99), throughput (RPS), error rate (%), CPU %, memory usage
- Metrics panel can be a dedicated sidebar/bottom panel or overlay on the canvas

## UX/UI Considerations

- Metrics panel accessible via toolbar toggle or always visible in a collapsible bottom panel
- Per-service metric cards with sparkline graphs
- Colour coding: green = healthy ranges, yellow = elevated, red = critical
- Time range selector (last 30s, 1m, 5m)
- Ability to click a node on the canvas to see its detailed metrics

## Acceptance Criteria

1. Prometheus container deploys automatically alongside user containers and scrapes metrics
2. Per-service latency (p50, p95, p99) is collected and displayed
3. Per-service throughput (requests per second) is displayed
4. Per-service error rate (percentage) is displayed
5. Per-container CPU and memory usage is collected and displayed
6. Metrics refresh at 1-second intervals via WebSocket
7. Metrics panel is accessible from the frontend UI

## Dependencies

- **Depends on**: PBI-9 (deploy flow — containers must be deployed), PBI-10 (traffic generation — need traffic for meaningful metrics)
- **External**: Prometheus Docker image, cAdvisor or Docker stats API

## Open Questions

- cAdvisor vs Docker stats API for resource metrics? (cAdvisor is more comprehensive but adds another container)
- Should we include Grafana as an optional advanced view? (PRD lists it as a stretch goal)
- How to handle metric storage between deploy sessions? (ephemeral or persist?)

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 11`._
