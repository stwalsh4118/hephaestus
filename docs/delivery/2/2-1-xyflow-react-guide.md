# 2-1 External Package Guide: @xyflow/react

- Research date: 2026-02-18
- Package: `@xyflow/react@12.10.0`
- Docs consulted:
  - https://reactflow.dev/learn
  - https://reactflow.dev/api-reference/react-flow
  - https://reactflow.dev/api-reference/hooks/use-react-flow
  - https://reactflow.dev/examples/interaction/drag-and-drop

## Verified API Usage (PBI-2 Scope)

### 1) Core canvas component and required CSS

```tsx
import { ReactFlow } from "@xyflow/react";
import "@xyflow/react/dist/style.css";

<div className="h-full w-full">
  <ReactFlow nodes={nodes} edges={edges} />
</div>;
```

Notes:
- React Flow requires the stylesheet import to render correctly.
- The parent wrapper must have explicit width/height.

### 2) Controlled state handlers

```tsx
import {
  ReactFlow,
  applyNodeChanges,
  applyEdgeChanges,
  type NodeChange,
  type EdgeChange,
} from "@xyflow/react";

const onNodesChange = (changes: NodeChange[]) => {
  setNodes((current) => applyNodeChanges(changes, current));
};

const onEdgesChange = (changes: EdgeChange[]) => {
  setEdges((current) => applyEdgeChanges(changes, current));
};
```

Notes:
- `applyNodeChanges` and `applyEdgeChanges` are the documented utilities for controlled flows.

### 3) Built-in viewport UX components

```tsx
import { Background, Controls, MiniMap, ReactFlow } from "@xyflow/react";

<ReactFlow nodes={nodes} edges={edges}>
  <Background gap={16} size={1} />
  <MiniMap />
  <Controls />
</ReactFlow>;
```

Notes:
- `Background`, `MiniMap`, and `Controls` are first-class React Flow components.

### 4) Drag-and-drop coordinate conversion

```tsx
import { useReactFlow } from "@xyflow/react";

const { screenToFlowPosition } = useReactFlow();

const onDrop = (event: React.DragEvent) => {
  event.preventDefault();

  const position = screenToFlowPosition({
    x: event.clientX,
    y: event.clientY,
  });

  addNodeAt(position);
};
```

Notes:
- `screenToFlowPosition` is the documented way to map pointer coordinates to canvas coordinates.

### 5) Provider requirements for hooks

```tsx
import { ReactFlowProvider } from "@xyflow/react";

<ReactFlowProvider>
  <DiagramCanvas />
</ReactFlowProvider>;
```

Notes:
- `useReactFlow` must be used inside `ReactFlowProvider` or within a `ReactFlow` child tree.
