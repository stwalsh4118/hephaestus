# Tasks for PBI 2: Diagram Canvas

This document lists all tasks associated with PBI 2.

**Parent PBI**: [PBI 2: Diagram Canvas](./prd.md)

## Task Summary

| Task ID | Name | Status | Description |
| :------ | :--- | :----- | :---------- |
| 2-1 | [Install React Flow and Zustand Dependencies](./2-1.md) | Proposed | Install @xyflow/react and zustand packages; create research guides with verified API usage |
| 2-2 | [Define Canvas Types and Constants](./2-2.md) | Proposed | Define TypeScript interfaces for node data, palette items, and canvas constants |
| 2-3 | [Create Zustand Canvas State Store](./2-3.md) | Proposed | Implement Zustand store for managing nodes, edges, and viewport state |
| 2-4 | [Build Base React Flow Canvas Component](./2-4.md) | Proposed | Create React Flow canvas with background grid, minimap, zoom/pan controls |
| 2-5 | [Build Component Palette Sidebar](./2-5.md) | Proposed | Create sidebar with draggable generic component items |
| 2-6 | [Integrate Drag-and-Drop from Palette to Canvas](./2-6.md) | Proposed | Wire onDrop/onDragOver handlers to create nodes from palette drops |
| 2-7 | [E2E CoS Test](./2-7.md) | Proposed | Verify all PBI-2 acceptance criteria end-to-end |

## Dependency Graph

```
2-1 (Install Dependencies)
 │
 ▼
2-2 (Types & Constants)
 │
 ├──────────┬──────────┐
 ▼          ▼          ▼
2-3 (Store) 2-5 (Palette)
 │          │          │
 ├──────────┘          │
 ▼                     │
2-4 (Canvas)           │
 │                     │
 ├─────────────────────┘
 ▼
2-6 (DnD Integration)
 │
 ▼
2-7 (E2E CoS Test)
```

## Implementation Order

1. **2-1**: Install React Flow and Zustand Dependencies — no dependencies, foundational
2. **2-2**: Define Canvas Types and Constants — depends on 2-1 for React Flow type imports
3. **2-3**: Create Zustand Canvas State Store — depends on 2-2 for type definitions
4. **2-5**: Build Component Palette Sidebar — depends on 2-2 for palette item types/constants (can run in parallel with 2-3)
5. **2-4**: Build Base React Flow Canvas Component — depends on 2-2 (types) and 2-3 (store)
6. **2-6**: Integrate Drag-and-Drop from Palette to Canvas — depends on 2-3, 2-4, and 2-5
7. **2-7**: E2E CoS Test — depends on all prior tasks (2-6 transitively)

## Complexity Ratings

| Task ID | Complexity | External Packages |
|---------|------------|-------------------|
| 2-1 | Simple | @xyflow/react, zustand |
| 2-2 | Simple | — |
| 2-3 | Medium | — |
| 2-4 | Medium | — |
| 2-5 | Simple | — |
| 2-6 | Medium | — |
| 2-7 | Medium | — |

## External Package Research Required

| Package | Guide Document | Required By |
|---------|---------------|-------------|
| @xyflow/react | `2-1-xyflow-react-guide.md` | Task 2-1 |
| zustand | `2-1-zustand-guide.md` | Task 2-1 |
