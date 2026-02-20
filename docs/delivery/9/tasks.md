# Tasks for PBI 9: Deploy Flow Integration

This document lists all tasks associated with PBI 9.

**Parent PBI**: [PBI 9: Deploy Flow Integration](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 9-1 | [Deploy Service — Types, State Tracker & Core Logic](./9-1.md) | Proposed | Define deploy types/interfaces; implement DeploymentManager with deploy, teardown, status, and diff logic |
| 9-2 | [Deploy REST Endpoints](./9-2.md) | Proposed | Add POST/DELETE/GET /api/deploy handlers and wire into router |
| 9-3 | [WebSocket Status Hub & Broadcasting](./9-3.md) | Proposed | Multi-client WebSocket hub with health-polling broadcast integration |
| 9-4 | [Frontend API Client & Deploy Store](./9-4.md) | Proposed | Fetch wrappers for deploy endpoints, WebSocket client, Zustand deploy store |
| 9-5 | [Deploy UI — Toolbar Buttons & Node Status Badges](./9-5.md) | Proposed | Deploy/Teardown buttons in toolbar; status indicators on canvas nodes; error feedback |
| 9-6 | [Live Topology Updates](./9-6.md) | Proposed | Incremental deploy/undeploy when diagram changes while deployed |
| 9-7 | [E2E CoS Test](./9-7.md) | Proposed | End-to-end verification of all PBI 9 acceptance criteria |

## Dependency Graph

```
9-1 (Deploy Service Core)
 ├──▶ 9-2 (REST Endpoints)
 │     └──▶ 9-4 (Frontend API Client & Store) ──▶ 9-5 (Deploy UI)
 └──▶ 9-3 (WebSocket Hub)                          │
       └──▶ 9-4                                     ▼
                                               9-6 (Live Topology Updates)
                                                    │
                                                    ▼
                                               9-7 (E2E CoS Test)
```

## Implementation Order

1. **9-1** — Deploy Service Types, State Tracker & Core Logic (foundation — all other tasks depend on these types and the DeploymentManager)
2. **9-2** — Deploy REST Endpoints (exposes deploy service via HTTP; needed before frontend can call anything)
3. **9-3** — WebSocket Status Hub & Broadcasting (completes the backend — real-time status push; can be parallel with 9-2 but ordered here for sequential flow)
4. **9-4** — Frontend API Client & Deploy Store (consumes endpoints from 9-2 and WebSocket from 9-3)
5. **9-5** — Deploy UI — Toolbar Buttons & Node Status Badges (renders state from 9-4's deploy store)
6. **9-6** — Live Topology Updates (builds on deployed state from 9-5 + diff logic from 9-1; adds PUT endpoint)
7. **9-7** — E2E CoS Test (validates everything end-to-end; must be last)

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 9-1     | Complex    | None              |
| 9-2     | Simple     | None              |
| 9-3     | Medium     | None              |
| 9-4     | Medium     | None              |
| 9-5     | Medium     | None              |
| 9-6     | Medium     | None              |
| 9-7     | Complex    | None              |

## External Package Research Required

None — all functionality uses existing project dependencies (React Flow, Zustand, Go stdlib, Docker SDK, gorilla/websocket).
