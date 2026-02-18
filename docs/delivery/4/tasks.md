# Tasks for PBI 4: Connections & Topology Export

This document lists all tasks associated with PBI 4.

**Parent PBI**: [PBI 4: Connections & Topology Export](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 4-1 | [Edge Types, Constants & Store Extensions](./4-1.md) | Done | Define edge data types, edge-related constants, and extend the canvas store with onConnect handler, duplicate edge validation, and edge label update action |
| 4-2 | [Custom Labeled Edge Component](./4-2.md) | Done | Create a custom React Flow edge component with directional arrows and editable labels (double-click to edit) |
| 4-3 | [Canvas Edge Integration](./4-3.md) | Done | Wire edge creation, custom edge types, and connection handling into DiagramCanvas |
| 4-4 | [Topology JSON Export](./4-4.md) | Done | Implement export function mapping React Flow state to PRD diagram JSON schema, with toolbar button and file download |
| 4-5 | [Topology JSON Import](./4-5.md) | Done | Implement import function parsing PRD JSON schema to restore React Flow state, with toolbar button and file picker |
| 4-6 | [E2E CoS Test](./4-6.md) | Done | End-to-end Playwright tests verifying all PBI-4 acceptance criteria |

## Dependency Graph

```
4-1 (Types, Constants, Store)
 ├──► 4-2 (Labeled Edge Component)
 │     └──► 4-3 (Canvas Integration)
 │           ├──► 4-4 (Export)
 │           └──► 4-5 (Import)
 │                 └──► 4-6 (E2E CoS Test)
 └──────────────────────► 4-4 (Export)
```

## Implementation Order

1. **4-1** — Foundation: types, constants, and store logic must exist before any edge UI or feature.
2. **4-2** — Edge component depends on types/constants from 4-1.
3. **4-3** — Canvas wiring depends on edge component (4-2) and store (4-1).
4. **4-4** — Export depends on working edges (4-3) and store state (4-1).
5. **4-5** — Import depends on export schema (4-4) and canvas integration (4-3).
6. **4-6** — E2E tests verify all features are working together; depends on all above.

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 4-1 | Simple | None |
| 4-2 | Medium | None |
| 4-3 | Simple | None |
| 4-4 | Medium | None |
| 4-5 | Medium | None |
| 4-6 | Medium | None |

## External Package Research Required

None — all functionality is achievable with the existing `@xyflow/react` and standard browser APIs.
