# PBI-12: Deploy Robustness & Hardening

[View in Backlog](../backlog.md)

## Overview

Harden the deploy flow by fixing port conflicts on update deploy, translating raw Docker errors into user-friendly messages, and adding real-Docker integration tests to catch issues that mock-based tests miss.

## Problem Statement

Manual testing of PBI 9 revealed three categories of issues:

1. **Port conflicts on update deploy**: When adding a node to a running deployment, `ApplyDiff` creates a fresh `PortAllocator` starting at port 10000 with no awareness of ports already bound by existing containers. This causes "port already allocated" errors from Docker.
2. **Raw error messages**: Docker daemon errors (port binding failures, image pull errors, container start failures) are passed verbatim to the frontend toast. Users see messages like `start container "3f86f7be...": Error response from daemon: failed to set up container networking...Bind for :::10000 failed: port is already allocated` instead of actionable guidance.
3. **Insufficient test coverage**: All existing deploy tests use mock orchestrators. Real Docker issues (port conflicts, image availability, network setup, container startup failures) are invisible to the test suite.

## User Stories

- As a user, I want deploy operations to handle common issues gracefully and show me clear error messages so that I can understand and fix problems without reading Docker internals
- As a user, I want to add or remove services from a running deployment without encountering port conflicts so that live topology updates work reliably
- As a developer, I want integration tests that exercise the real Docker daemon so that I catch deployment issues before users do

## Technical Approach

### Port Conflict Fix
- Track allocated host ports in `DeploymentManager` (persist the port assignments alongside node→container mappings)
- When `ApplyDiff` creates a new `Translator`, pre-seed the `PortAllocator` with ports already in use by running containers
- Add a method to `PortAllocator` to mark ports as used (e.g., `Reserve(ports []string)`)

### User-Facing Error Messages
- Add an error translation layer in the backend deploy handler that maps known Docker error patterns to user-friendly messages
- Pattern matching on error strings (port already allocated, image not found, network failure, timeout, etc.)
- Preserve the original error in server logs for debugging; return only the friendly message in the API response
- Fallback: if no pattern matches, return a generic "Deploy failed — check server logs for details" rather than the raw error

### Integration Tests
- Gate behind `//go:build docker_integration` build tag
- Test scenarios: multi-service initial deploy, update deploy with added/removed nodes, duplicate service types with different names, full teardown, port allocation across deploy cycles
- Use real Docker daemon — tests create and clean up their own containers
- 30s timeout per test to prevent hangs

## UX/UI Considerations

- Error toast messages become shorter and actionable (e.g., "Port conflict: another service is already using this port. Try removing and re-deploying." instead of raw Docker output)
- No UI changes needed — the existing `ErrorToast` component displays whatever message the deploy store provides

## Acceptance Criteria

1. `PUT /api/deploy` (update deploy) correctly allocates ports that don't conflict with already-running containers from the same deployment
2. Deploy, update, and teardown errors return user-friendly messages in the API `error` field (no raw Docker daemon error strings)
3. Original technical errors are logged server-side for debugging
4. Integration tests exist covering: multi-service deploy, update add node, update remove node, deploy with duplicate service types, full teardown, port allocation across deploy→update cycles
5. Integration tests are gated behind `//go:build docker_integration` and pass when Docker is available
6. Existing unit tests continue to pass unchanged

## Dependencies

- **Depends on**: PBI 9 (Deploy Flow Integration) — must be complete first
- **Blocks**: None
- **External**: Docker daemon required for integration tests

## Open Questions

None — technical approach confirmed during analysis.

## Related Tasks

_Tasks will be created when this PBI is approved via `/plan-pbi 12`._
