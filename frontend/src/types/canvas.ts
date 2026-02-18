import type { Edge, Node, Viewport } from "@xyflow/react";

export type NodeType = "service" | "worker" | "database" | "queue";

export interface PaletteItem {
  id: NodeType;
  label: string;
  icon: string;
  description: string;
}

export interface CanvasNodeData extends Record<string, unknown> {
  label: string;
  type: NodeType;
  description: string;
}

export type CanvasNode = Node<CanvasNodeData>;
export type CanvasEdge = Edge;
export type CanvasViewport = Viewport;
