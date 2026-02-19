# PBI-8: Mock API System

[View in Backlog](../backlog.md)

## Overview

Build the system that takes user-defined API endpoints (from the service configuration panel) and generates valid OpenAPI 3.0 specifications, then serves them via Prism mock server containers that return realistic fake data.

## Problem Statement

API Service nodes in the diagram need to actually respond to HTTP requests with realistic data. Users define endpoints in the UI (method, path, response schema), but there's no mechanism to turn those definitions into a running mock server. This PBI closes that gap.

## User Stories

- As a user, I want my defined API endpoints to be served by a real mock server so that I can send requests and get responses
- As a developer, I want to generate OpenAPI specs from endpoint definitions so that Prism can serve them automatically

## Technical Approach

- Build an OpenAPI 3.0 spec generator in the Go backend that takes endpoint definitions (method, path, response schema) and produces valid specs
- Use JSON Schema Faker concepts to generate realistic response examples (or rely on Prism's built-in dynamic mocking)
- Create a Prism wrapper Docker image (or use the official Prism image) that reads a mounted OpenAPI spec file
- When deploying an API Service node:
  1. Generate OpenAPI spec from the node's endpoint config
  2. Write spec to a temp file
  3. Mount the spec into the Prism container
  4. Start Prism pointing at the spec
- Verify that mock endpoints respond correctly and that cross-service HTTP calls resolve within the Docker network

## UX/UI Considerations

N/A — backend PBI, though the endpoint definitions come from PBI-3's configuration panel.

## Acceptance Criteria

1. Endpoint definitions (method, path, response schema) are translated into valid OpenAPI 3.0 specs
2. Generated specs pass OpenAPI validation
3. Prism container starts with a mounted spec and serves the defined endpoints
4. Mock responses contain data matching the defined response schemas
5. Cross-service HTTP calls resolve via Docker DNS within the shared network (e.g., `http://user-service:8080/users`)

## Dependencies

- **Depends on**: PBI-6 (orchestration engine), PBI-7 (service-to-container mapping for Prism template)
- **External**: Prism (stoplight/prism Docker image)

## Open Questions

- Prism vs WireMock — PRD mentions both. Prism is simpler for OpenAPI-based mocking; recommend Prism for MVP.
- How complex can response schemas be? (nested objects, arrays, references)
- Should we support request validation (Prism can validate requests against the spec)?

## Related Tasks

[View Tasks](./tasks.md)
