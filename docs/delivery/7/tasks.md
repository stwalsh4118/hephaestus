# Tasks for PBI 7: Service-to-Container Mapping

This document lists all tasks associated with PBI 7.

**Parent PBI**: [PBI 7: Service-to-Container Mapping](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 7-1 | [Container Template Types and Registry](./7-1.md) | Review | Define ContainerTemplate interface, TemplateRegistry, and image/port constants for all 5 service types |
| 7-2 | [Service Container Templates](./7-2.md) | Review | Implement template builders for all 5 service types (PostgreSQL, Redis, Nginx, RabbitMQ, API Service) |
| 7-3 | [Port Allocator](./7-3.md) | Review | Dynamic host port allocation with conflict avoidance |
| 7-4 | [Dependency Resolver](./7-4.md) | Review | Build dependency graph from diagram edges and produce topological startup order |
| 7-5 | [Diagram-to-Container Translator](./7-5.md) | Review | Orchestrate templates, port allocator, and dependency resolver to translate a Diagram into ordered ContainerConfigs |
| 7-6 | [API Specification Update](./7-6.md) | Review | Update docker-api.md with new types, interfaces, and constants from PBI 7 |
| 7-7 | [E2E CoS Test](./7-7.md) | Review | End-to-end verification of all PBI 7 acceptance criteria |

## Dependency Graph

```
7-1 (Types & Registry)
 │
 ├──► 7-2 (Service Templates)
 │         │
 │         ├──────────────────────► 7-5 (Translator)
 │         │                           │
 │    7-3 (Port Allocator) ────────►   │
 │         │                           │
 │    7-4 (Dependency Resolver) ───►   │
 │                                     │
 │                                     ▼
 │                                 7-6 (API Spec Update)
 │                                     │
 │                                     ▼
 └────────────────────────────────► 7-7 (E2E CoS Test)
```

## Implementation Order

1. **7-1** — Container Template Types and Registry (foundation types required by all subsequent tasks)
2. **7-2** — Service Container Templates (depends on 7-1 types; implements the 5 builders)
3. **7-3** — Port Allocator (independent of 7-2 but needed by 7-5)
4. **7-4** — Dependency Resolver (independent of 7-2/7-3 but needed by 7-5)
5. **7-5** — Diagram-to-Container Translator (depends on 7-2, 7-3, 7-4; ties everything together)
6. **7-6** — API Specification Update (depends on 7-5; documents final interfaces)
7. **7-7** — E2E CoS Test (depends on all above; validates acceptance criteria)

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 7-1 | Simple | None |
| 7-2 | Medium | None |
| 7-3 | Simple | None |
| 7-4 | Medium | None |
| 7-5 | Medium | None |
| 7-6 | Simple | None |
| 7-7 | Complex | None |

## External Package Research Required

None — all tasks use the Go standard library and the existing Docker SDK integration from PBI 6.
