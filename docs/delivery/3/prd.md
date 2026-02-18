# PBI-3: Service Component Library & Configuration

[View in Backlog](../backlog.md)

## Overview

Create the five MVP service type nodes (API Service, PostgreSQL, Redis, Nginx, RabbitMQ) with distinct visual representations and a configuration side panel where users can define service-specific settings, including endpoint definitions for API services.

## Problem Statement

The canvas (PBI-2) supports generic nodes, but users need specific service types with distinct visuals and meaningful configuration options. Without typed components and a config panel, users can't define what their services actually do.

## User Stories

- As a user, I want to see distinct visual representations for different service types (database, cache, load balancer, etc.) so that my diagram is immediately readable
- As a user, I want to click a node and configure its settings (name, endpoints, ports) in a side panel so that I can define service behavior
- As a user, I want to define API endpoints (method, path, response schema) for API Service nodes so that mock APIs can be generated later

## Technical Approach

- Create custom React Flow node components for each of the 5 service types
- Each node type has a distinct icon, colour scheme, and shape to be visually distinguishable
- Implement a configuration side panel that appears when a node is selected
- Panel fields are dynamic based on service type (e.g., API Service shows endpoint editor, Database shows engine/version)
- Endpoint definition UI: list of {method, path, responseSchema} entries with add/remove
- Config changes update the diagram state immediately

## UX/UI Considerations

- Nodes should be visually distinct at a glance (colour + icon per type)
- Config panel slides in from the right when a node is selected
- Endpoint editor supports adding multiple endpoints with method dropdown, path input, and JSON schema editor
- Clear visual feedback when config is saved/applied

## Acceptance Criteria

1. 5 service types available in the palette: API Service, PostgreSQL, Redis, Nginx, RabbitMQ
2. Each service type has a distinct visual appearance (icon and colour)
3. Clicking a node opens a configuration side panel
4. API Service configuration includes an endpoint editor (method, path, response schema)
5. Database configuration includes engine and version selection
6. Configuration changes persist to the diagram state

## Dependencies

- **Depends on**: PBI-2 (diagram canvas)
- **External**: Icon library (e.g., Lucide, Heroicons)

## Open Questions

- How detailed should the JSON response schema editor be? (free-form JSON vs structured builder)
- Should we support custom node names or auto-generate from type?

## Related Tasks

[View Tasks](./tasks.md)
