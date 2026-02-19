# Tasks for PBI 6: Docker Orchestration Engine

This document lists all tasks associated with PBI 6.

**Parent PBI**: [PBI 6: Docker Orchestration Engine](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 6-1 | [Docker SDK Integration & Client Wrapper](./6-1.md) | Done | Add Docker SDK dependency, create research guide, and implement client wrapper with connection handling |
| 6-2 | [Container Types, Constants & Orchestrator Interface](./6-2.md) | Done | Define container models, status types, naming constants, and the orchestrator interface contract |
| 6-3 | [Docker Network Management](./6-3.md) | Done | Create and manage a shared Docker bridge network for all simulator containers |
| 6-4 | [Container Lifecycle Management](./6-4.md) | Done | Implement container create, start, stop, remove, list, and inspect operations |
| 6-5 | [Health Check & Status Reporting](./6-5.md) | Review | Implement periodic container health check polling and status reporting |
| 6-6 | [Teardown & Graceful Cleanup](./6-6.md) | Review | Implement full teardown of managed containers and network, wired into server shutdown |
| 6-7 | [E2E CoS Test](./6-7.md) | Review | End-to-end tests verifying all PBI 6 acceptance criteria |

## Dependency Graph

```
6-1 (SDK + Client)
 │
 ├──► 6-2 (Types & Interface)
 │     │
 │     ├──► 6-3 (Network Management)
 │     │     │
 │     │     └──► 6-4 (Container Lifecycle)
 │     │           │
 │     │           ├──► 6-5 (Health Check)
 │     │           │
 │     │           └──► 6-6 (Teardown)
 │     │
 │     └────────────────► 6-7 (E2E CoS Test)
 │                          ▲
 (all tasks feed into 6-7)
```

## Implementation Order

1. **6-1** — Foundation: Docker SDK dependency and client wrapper (no other tasks can proceed without this)
2. **6-2** — Types and interface contract (defines the shapes all implementations use)
3. **6-3** — Network management (containers need a network to join during creation)
4. **6-4** — Container lifecycle (depends on network being available, uses types from 6-2)
5. **6-5** — Health checks (depends on containers existing to inspect)
6. **6-6** — Teardown (depends on container lifecycle + network management to clean up)
7. **6-7** — E2E CoS Test (validates all prior tasks against real Docker daemon)

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 6-1 | Medium | Docker SDK for Go |
| 6-2 | Simple | None |
| 6-3 | Medium | None |
| 6-4 | Complex | None |
| 6-5 | Medium | None |
| 6-6 | Medium | None |
| 6-7 | Complex | None |

## External Package Research Required

| Package | Guide Document | Required By |
|---------|---------------|-------------|
| Docker SDK for Go (`github.com/docker/docker`) | `6-1-docker-sdk-guide.md` | Task 6-1 |
