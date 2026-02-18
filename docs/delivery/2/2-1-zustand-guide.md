# 2-1 External Package Guide: zustand

- Research date: 2026-02-18
- Package: `zustand@5.0.11`
- Docs consulted:
  - https://raw.githubusercontent.com/pmndrs/zustand/main/README.md
  - https://raw.githubusercontent.com/pmndrs/zustand/main/docs/apis/create.md
  - https://raw.githubusercontent.com/pmndrs/zustand/main/docs/guides/beginner-typescript.md
  - https://raw.githubusercontent.com/pmndrs/zustand/main/docs/integrations/persisting-store-data.md

## Verified API Usage (PBI-2 Scope)

### 1) Create a typed store

```ts
import { create } from "zustand";

type CanvasStore = {
  nodes: { id: string; label: string }[];
  addNode: (node: { id: string; label: string }) => void;
};

export const useCanvasStore = create<CanvasStore>()((set) => ({
  nodes: [],
  addNode: (node) => set((state) => ({ nodes: [...state.nodes, node] })),
}));
```

Notes:
- `create<T>()((set, get) => state)` is the documented TypeScript pattern.

### 2) Subscribe with selectors in React

```tsx
const nodes = useCanvasStore((state) => state.nodes);
const addNode = useCanvasStore((state) => state.addNode);
```

Notes:
- Selector usage limits re-renders to selected slices.

### 3) State update patterns

```ts
set((state) => ({ nodes: [...state.nodes, nextNode] }));
set({ nodes: [] });
```

Notes:
- `set` supports functional updates and shallow merge updates.
- Treat state as immutable when updating arrays/objects.

### 4) Optional session persistence middleware

```ts
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

export const useSessionStore = create<{ value: number }>()(
  persist(
    () => ({ value: 0 }),
    {
      name: "hephaestus-session",
      storage: createJSONStorage(() => sessionStorage),
    },
  ),
);
```

Notes:
- `persist` middleware is available if session-level persistence is needed in follow-up work.
- For PBI-2, in-memory store persistence across re-renders is sufficient.
