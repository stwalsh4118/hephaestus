# PBI-6: Docker Orchestration Engine

[View in Backlog](../backlog.md)

## Overview

Integrate the Docker SDK for Go to programmatically manage container lifecycles — creating, starting, stopping, and removing containers. Establish a single shared Docker network and implement health checking for all managed containers.

## Problem Statement

The core value proposition of the platform is turning diagrams into running containers. The orchestration engine is the component that interacts with Docker to make this happen. Without it, diagrams remain static.

## User Stories

- As a developer, I want to programmatically create and start Docker containers so that diagram nodes can become running services
- As a developer, I want to manage a shared Docker network so that all containers can communicate
- As a developer, I want health checks for containers so that the system knows which services are healthy

## Technical Approach

- Use the Docker SDK for Go (`github.com/docker/docker/client`)
- Implement a container manager service with methods: Create, Start, Stop, Remove, List, Inspect
- Create and manage a single Docker bridge network for all simulator containers
- Containers join the shared network on creation with their service name as hostname
- Health check polling: periodically inspect container state and report status (running, stopped, error, healthy, unhealthy)
- Graceful cleanup: teardown method removes all managed containers and the network
- Container naming convention to avoid conflicts with user's other Docker containers

## UX/UI Considerations

N/A — backend infrastructure PBI.

## Acceptance Criteria

1. Docker SDK client initialises and connects to the Docker daemon
2. Containers can be created from specified images with configuration (env vars, ports, volumes)
3. A shared Docker network is created and all managed containers join it
4. Containers can be started, stopped, and removed individually
5. Health check reports accurate container status (running/stopped/error)
6. Teardown removes all managed containers and the shared network cleanly
7. Container names are prefixed/namespaced to avoid conflicts

## Dependencies

- **Depends on**: PBI-1 (project foundation), PBI-5 (Go backend to host the orchestration service)
- **External**: Docker Engine running on host, Docker SDK for Go

## Open Questions

- Should we support resource limits (CPU/memory) on containers from the start?
- How to handle Docker daemon not running (error messaging)?

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 6`._
