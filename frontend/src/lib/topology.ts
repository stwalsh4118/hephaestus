import { MarkerType } from "@xyflow/react";

import { NODE_TYPE_SERVICE } from "@/components/canvas/nodes";
import { EDGE_TYPE_LABELED } from "@/constants/canvas";
import type {
  CanvasEdge,
  CanvasNode,
  CanvasNodeData,
  ServiceConfig,
  ServiceType,
} from "@/types/canvas";

const DEFAULT_DIAGRAM_NAME = "Untitled Diagram";

export interface DiagramJsonNode {
  id: string;
  type: string;
  name: string;
  description: string;
  position: { x: number; y: number };
  config?: ServiceConfig;
}

export interface DiagramJsonEdge {
  id: string;
  source: string;
  target: string;
  label: string;
}

export interface DiagramJson {
  id: string;
  name: string;
  nodes: DiagramJsonNode[];
  edges: DiagramJsonEdge[];
}

export function exportDiagram(
  nodes: CanvasNode[],
  edges: CanvasEdge[]
): DiagramJson {
  return {
    id: crypto.randomUUID(),
    name: DEFAULT_DIAGRAM_NAME,
    nodes: nodes.map((node) => ({
      id: node.id,
      type: node.data.type,
      name: node.data.label,
      description: node.data.description,
      position: { x: node.position.x, y: node.position.y },
      ...(node.data.config ? { config: node.data.config } : {}),
    })),
    edges: edges.map((edge) => ({
      id: edge.id,
      source: edge.source,
      target: edge.target,
      label: edge.data?.label ?? "",
    })),
  };
}

export function downloadDiagram(diagram: DiagramJson): void {
  const json = JSON.stringify(diagram, null, 2);
  const blob = new Blob([json], { type: "application/json" });
  const url = URL.createObjectURL(blob);

  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = "diagram.json";
  anchor.click();

  URL.revokeObjectURL(url);
}

function validateNode(node: unknown): node is DiagramJsonNode {
  if (!node || typeof node !== "object") return false;
  const n = node as Record<string, unknown>;
  return (
    typeof n.id === "string" &&
    typeof n.type === "string" &&
    typeof n.name === "string" &&
    n.position !== null &&
    typeof n.position === "object" &&
    typeof (n.position as Record<string, unknown>).x === "number" &&
    typeof (n.position as Record<string, unknown>).y === "number"
  );
}

function validateEdge(edge: unknown): edge is DiagramJsonEdge {
  if (!edge || typeof edge !== "object") return false;
  const e = edge as Record<string, unknown>;
  return (
    typeof e.id === "string" &&
    typeof e.source === "string" &&
    typeof e.target === "string"
  );
}

export function importDiagram(json: unknown): {
  nodes: CanvasNode[];
  edges: CanvasEdge[];
} {
  if (!json || typeof json !== "object") {
    throw new Error("Invalid diagram JSON: expected an object");
  }

  const diagram = json as Record<string, unknown>;

  if (!Array.isArray(diagram.nodes) || !Array.isArray(diagram.edges)) {
    throw new Error("Invalid diagram JSON: missing nodes or edges arrays");
  }

  for (const node of diagram.nodes) {
    if (!validateNode(node)) {
      throw new Error(
        "Invalid diagram JSON: node missing required fields (id, type, name, position)"
      );
    }
  }

  for (const edge of diagram.edges) {
    if (!validateEdge(edge)) {
      throw new Error(
        "Invalid diagram JSON: edge missing required fields (id, source, target)"
      );
    }
  }

  const validNodes = diagram.nodes as DiagramJsonNode[];
  const validEdges = diagram.edges as DiagramJsonEdge[];

  const nodes: CanvasNode[] = validNodes.map((node) => ({
    id: node.id,
    type: NODE_TYPE_SERVICE,
    position: { x: node.position.x, y: node.position.y },
    data: {
      label: node.name,
      type: node.type as ServiceType,
      description: node.description ?? "",
      ...(node.config ? { config: node.config } : {}),
    } as CanvasNodeData,
  }));

  const edges: CanvasEdge[] = validEdges.map((edge) => ({
    id: edge.id,
    source: edge.source,
    target: edge.target,
    type: EDGE_TYPE_LABELED,
    markerEnd: { type: MarkerType.ArrowClosed },
    data: { label: edge.label ?? "" },
  }));

  return { nodes, edges };
}
