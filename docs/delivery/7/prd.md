# PBI-7: Service-to-Container Mapping

[View in Backlog](../backlog.md)

## Overview

Create container configuration templates for each of the 5 MVP service types so that diagram nodes translate into correctly configured Docker containers. Handle service dependency ordering and port allocation.

## Problem Statement

The orchestration engine (PBI-6) can manage generic containers, but it doesn't know how to turn a "PostgreSQL" diagram node into a properly configured PostgreSQL container. Each service type needs a specific image, environment variables, ports, and volumes. Startup order matters (databases before API services).

## User Stories

- As a developer, I want each service type to map to a concrete Docker container configuration so that deploying a diagram produces correctly running services
- As a developer, I want service dependencies to determine startup order so that services start in the right sequence

## Technical Approach

- Define a template/factory pattern: each service type has a corresponding container config builder
- Templates for MVP service types:
  - **API Service**: Prism mock server image, mounted OpenAPI spec, configurable port
  - **PostgreSQL**: Official postgres image, default credentials via env vars, data volume
  - **Redis**: Official redis image, default config, standard port
  - **Nginx**: Official nginx image, generated config for upstream routing, exposed port
  - **RabbitMQ**: Official rabbitmq image with management plugin, standard ports
- Dependency graph: build order from diagram edges (e.g., databases start before services that connect to them)
- Port allocator: assign host ports dynamically to avoid conflicts
- Translate diagram node config into container-specific environment variables and mounts

## UX/UI Considerations

N/A — backend infrastructure PBI.

## Acceptance Criteria

1. Each of the 5 service types has a container template that produces valid Docker configurations
2. Templates use correct official Docker images for each service type
3. Service-specific configuration (env vars, volumes, ports) is correctly applied
4. Dependency ordering ensures databases and caches start before services that depend on them
5. Port allocation assigns unique host ports without conflicts
6. A diagram with all 5 service types can be translated into container configs that the orchestration engine accepts

## Dependencies

- **Depends on**: PBI-6 (Docker orchestration engine)
- **External**: Official Docker images (postgres, redis, nginx, rabbitmq)

## Open Questions

- Should we pull images proactively or lazily on first deploy?
- Default credentials for databases — hardcoded for MVP or user-configurable?

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 7`._
