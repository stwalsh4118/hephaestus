# PBI-4: Connections & Topology Export

[View in Backlog](../backlog.md)

## Overview

Enable users to draw labelled connections (edges) between service nodes and export the complete diagram topology as a JSON document conforming to the PRD schema. Support re-importing JSON to restore a diagram.

## Problem Statement

A system architecture diagram is not just nodes — the connections between services define data flow and dependencies. Users need to draw edges, label them, and produce a machine-readable JSON representation that the backend can consume for deployment.

## User Stories

- As a user, I want to draw connections between services by dragging from one node to another so that I can define data flow
- As a user, I want to label connections (e.g., "reads/writes", "caches") so that the relationship is clear
- As a user, I want to export my diagram as JSON so that it can be sent to the backend for deployment
- As a user, I want to import a JSON topology to restore a previously saved diagram

## Technical Approach

- Use React Flow's built-in edge/connection system for drawing connections
- Custom edge component with label support (editable label on click)
- Connection validation: prevent invalid connections (e.g., duplicate edges between same nodes)
- JSON export produces a document matching the PRD's Diagram JSON Schema (nodes, edges, positions, configs)
- JSON import parses the schema and reconstructs the React Flow state
- Export/import accessible via toolbar buttons

## UX/UI Considerations

- Edges drawn by dragging from a node handle to another node
- Edge labels displayed along the connection line, editable on double-click
- Animated or styled edges to indicate direction (arrows)
- Export button downloads JSON file; import button opens file picker

## Acceptance Criteria

1. Edges can be drawn between any two nodes by dragging from source to target handle
2. Edges display an arrow indicating direction
3. Edge labels can be added and edited
4. Duplicate edges between the same source-target pair are prevented
5. Export produces valid JSON matching the PRD Diagram JSON Schema
6. Importing a previously exported JSON restores the diagram accurately (nodes, positions, edges, configs)

## Dependencies

- **Depends on**: PBI-2 (diagram canvas), PBI-3 (service components with config data)
- **External**: None

## Open Questions

- Should we validate that certain connection types make sense (e.g., prevent connecting two databases directly)?
- Should export include diagram metadata (name, created date)?

## Related Tasks

[View Tasks](./tasks.md)
