# Tasks for PBI 8: Mock API System

This document lists all tasks associated with PBI 8.

**Parent PBI**: [PBI 8: Mock API System](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 8-1 | [Add Cmd field to ContainerConfig](./8-1.md) | Review | Add command/args support to ContainerConfig and update Docker engine to pass Cmd when creating containers |
| 8-2 | [OpenAPI 3.0 Spec Generator](./8-2.md) | Review | Create a Go package that converts endpoint definitions into valid OpenAPI 3.0.0 JSON specs |
| 8-3 | [Integrate Spec Generation into API Service Template](./8-3.md) | Review | Update APIServiceTemplate to parse endpoint config, generate OpenAPI spec, write to disk, and mount into Prism container |
| 8-4 | [E2E CoS Test](./8-4.md) | Review | End-to-end verification of all PBI-8 acceptance criteria with running Prism containers |

## Dependency Graph

```
8-1 (Cmd field) ──┐
                   ├──► 8-3 (Template Integration) ──► 8-4 (E2E CoS Test)
8-2 (Spec Gen) ───┘
```

## Implementation Order

1. **8-1** — Add Cmd field to ContainerConfig (no dependencies; foundational for Prism args)
2. **8-2** — OpenAPI 3.0 Spec Generator (no dependencies; can be done in parallel with 8-1)
3. **8-3** — Integrate Spec Generation into API Service Template (depends on 8-1 and 8-2)
4. **8-4** — E2E CoS Test (depends on 8-3; verifies all acceptance criteria end-to-end)

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 8-1 | Simple | None |
| 8-2 | Medium | None |
| 8-3 | Medium | None |
| 8-4 | Complex | None |

## External Package Research Required

None. PBI-8 uses the existing `stoplight/prism:latest` Docker image (already configured in PBI-7) and Go standard library for JSON/OpenAPI generation.
