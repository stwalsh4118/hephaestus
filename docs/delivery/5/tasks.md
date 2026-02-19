# Tasks for PBI 5: Go Backend API

This document lists all tasks associated with PBI 5.

**Parent PBI**: [PBI 5: Go Backend API](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 5-1 | [Diagram model types and validation](./5-1.md) | Proposed | Define Go structs for the diagram JSON schema and implement validation |
| 5-2 | [File-based diagram storage](./5-2.md) | Proposed | Implement a file-based JSON storage layer for diagram CRUD persistence |
| 5-3 | [CORS middleware](./5-3.md) | Proposed | Add CORS middleware to allow frontend dev server origin |
| 5-4 | [Diagram CRUD HTTP handlers](./5-4.md) | Proposed | Implement REST handlers for POST/GET/PUT /api/diagrams and wire routing |
| 5-5 | [WebSocket endpoint scaffolding](./5-5.md) | Proposed | Scaffold /ws/status WebSocket endpoint with ping/pong support |
| 5-6 | [E2E CoS Test](./5-6.md) | Proposed | End-to-end verification of all PBI-5 acceptance criteria |

## Dependency Graph

```
5-1 (Model & Validation)
 │
 ▼
5-2 (Storage)        5-3 (CORS)
 │                    │
 └────────┬───────────┘
          ▼
       5-4 (Handlers)     5-5 (WebSocket)
          │                  │
          └────────┬─────────┘
                   ▼
              5-6 (E2E Test)
```

## Implementation Order

1. **5-1** — Diagram model types and validation (no dependencies; foundational types used by all subsequent tasks)
2. **5-2** — File-based diagram storage (depends on 5-1 for diagram types)
3. **5-3** — CORS middleware (independent of 5-1/5-2; can be parallelized with 5-2 if desired)
4. **5-4** — Diagram CRUD HTTP handlers (depends on 5-1, 5-2, 5-3; integrates all three)
5. **5-5** — WebSocket endpoint scaffolding (independent; can be parallelized with 5-4)
6. **5-6** — E2E CoS Test (depends on all tasks; validates complete feature)

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 5-1 | Simple | None |
| 5-2 | Medium | `google/uuid` |
| 5-3 | Simple | None |
| 5-4 | Medium | None |
| 5-5 | Medium | `gorilla/websocket` |
| 5-6 | Complex | None |

## External Package Research Required

| Package | Task | Guide Document |
|---------|------|----------------|
| `gorilla/websocket` | 5-5 | `5-5-gorilla-websocket-guide.md` |
