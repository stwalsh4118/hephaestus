# PBI-2: Diagram Canvas

[View in Backlog](../backlog.md)

## Overview

Implement the core visual canvas using React Flow where users can drag components from a palette, drop them onto the canvas, and arrange them freely. This is the primary interaction surface for the entire application.

## Problem Statement

Users need a visual workspace to design system architectures. Without a canvas that supports drag-and-drop, positioning, and standard canvas interactions (zoom, pan), no other frontend feature can be built.

## User Stories

- As a user, I want to drag service components from a palette onto a canvas so that I can start designing my system
- As a user, I want to reposition nodes by dragging them so that I can arrange my architecture logically
- As a user, I want to zoom and pan the canvas so that I can work with diagrams of any size

## Technical Approach

- Integrate React Flow as the canvas library
- Create a sidebar/palette with draggable items (generic placeholders — specific service types come in PBI-3)
- Implement drag-from-palette-to-canvas using React Flow's onDrop/onDragOver handlers
- Canvas state managed via React Flow's built-in state or a lightweight store (Zustand)
- Support zoom (scroll wheel), pan (click-drag on background), and fit-to-view

## UX/UI Considerations

- Canvas occupies the main viewport area
- Palette/toolbar on the left or top for dragging components
- Minimap in corner for orientation on large diagrams
- Grid or dot background for alignment reference
- Node selection with visual highlight

## Acceptance Criteria

1. React Flow canvas renders in the main viewport
2. A sidebar palette displays draggable component items
3. Dragging an item from the palette onto the canvas creates a new node
4. Nodes can be repositioned by dragging on the canvas
5. Canvas supports zoom (scroll wheel) and pan (background drag)
6. Canvas state (nodes and positions) persists during the session (not lost on re-render)

## Dependencies

- **Depends on**: PBI-1 (project foundation)
- **External**: React Flow library

## Open Questions

- Zustand vs React Context for canvas state management?
- Should the minimap be included now or deferred?

## Related Tasks

_Tasks will be created when this PBI moves to Agreed via `/plan-pbi 2`._
