# Tasks for PBI 3: Service Component Library & Configuration

This document lists all tasks associated with PBI 3.

**Parent PBI**: [PBI 3: Service Component Library & Configuration](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 3-1 | [Service type definitions & constants](./3-1.md) | Proposed | Define 5 MVP service types, config interfaces, colour/icon constants, and updated palette items |
| 3-2 | [Custom React Flow node components](./3-2.md) | Proposed | Create visually distinct custom node components for each service type and register with React Flow |
| 3-3 | [Extend store with node selection & config](./3-3.md) | Proposed | Add selected node tracking, node config update action, and config data model to Zustand store |
| 3-4 | [Configuration side panel](./3-4.md) | Proposed | Create sliding panel UI that opens on node selection and displays service info |
| 3-5 | [Service-specific config forms](./3-5.md) | Proposed | Build dynamic config forms for PostgreSQL, Redis, Nginx, and RabbitMQ service types |
| 3-6 | [API endpoint editor](./3-6.md) | Proposed | Build endpoint definition UI with add/remove rows, method dropdown, path input, and JSON schema editor |
| 3-7 | [E2E CoS Test](./3-7.md) | Proposed | Automated end-to-end tests verifying all PBI-3 acceptance criteria |

## Dependency Graph

```
3-1 (Types & Constants)
 ├──► 3-2 (Custom Node Components)
 └──► 3-3 (Store: Selection & Config)
       ├──► 3-4 (Config Side Panel)
       │     └──► 3-5 (Service Config Forms)
       │           └──► 3-6 (API Endpoint Editor)
       └──────────────────┘
                          └──► 3-7 (E2E CoS Test)
```

## Implementation Order

1. **3-1** — Foundation: types, interfaces, constants needed by all subsequent tasks
2. **3-2** — Custom nodes: depends on 3-1 for service type definitions and colour constants
3. **3-3** — Store extension: depends on 3-1 for config type definitions
4. **3-4** — Config panel shell: depends on 3-3 for selected node state
5. **3-5** — Service forms: depends on 3-3 (updateNodeConfig) and 3-4 (panel container)
6. **3-6** — API endpoint editor: depends on 3-5 (form pattern and panel integration)
7. **3-7** — E2E test: depends on all prior tasks being complete

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 3-1 | Simple | None |
| 3-2 | Medium | None |
| 3-3 | Medium | None |
| 3-4 | Medium | None |
| 3-5 | Medium | None |
| 3-6 | Complex | None |
| 3-7 | Medium | None |

## External Package Research Required

None. All tasks use existing dependencies (React Flow, Zustand, Tailwind CSS, Playwright).
