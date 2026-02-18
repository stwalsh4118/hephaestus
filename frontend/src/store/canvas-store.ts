import {
  applyEdgeChanges,
  applyNodeChanges,
  type EdgeChange,
  type NodeChange,
  type XYPosition,
} from "@xyflow/react";
import { create } from "zustand";

import { DEFAULT_CANVAS_VIEWPORT } from "@/constants/canvas";
import type {
  CanvasEdge,
  CanvasNode,
  CanvasNodeData,
  CanvasViewport,
} from "@/types/canvas";

const INITIAL_NODES: CanvasNode[] = [];
const INITIAL_EDGES: CanvasEdge[] = [];
const INITIAL_VIEWPORT: CanvasViewport = DEFAULT_CANVAS_VIEWPORT;
const NODE_ID_PREFIX = "node";

let nodeCounter = 0;

const getNextNodeId = (): string => {
  nodeCounter += 1;
  return `${NODE_ID_PREFIX}-${nodeCounter}`;
};

interface AddNodeInput {
  position: XYPosition;
  data: CanvasNodeData;
}

interface CanvasStore {
  nodes: CanvasNode[];
  edges: CanvasEdge[];
  viewport: CanvasViewport;
  onNodesChange: (changes: NodeChange<CanvasNode>[]) => void;
  onEdgesChange: (changes: EdgeChange<CanvasEdge>[]) => void;
  onViewportChange: (viewport: CanvasViewport) => void;
  addNode: (input: AddNodeInput) => void;
}

export const useCanvasStore = create<CanvasStore>()((set) => ({
  nodes: INITIAL_NODES,
  edges: INITIAL_EDGES,
  viewport: INITIAL_VIEWPORT,
  onNodesChange: (changes) => {
    set((state) => ({
      nodes: applyNodeChanges(changes, state.nodes),
    }));
  },
  onEdgesChange: (changes) => {
    set((state) => ({
      edges: applyEdgeChanges(changes, state.edges),
    }));
  },
  onViewportChange: (viewport) => {
    set({ viewport });
  },
  addNode: ({ position, data }) => {
    const nodeId = getNextNodeId();

    set((state) => ({
      nodes: [
        ...state.nodes,
        {
          id: nodeId,
          position,
          data,
          type: "default",
        },
      ],
    }));
  },
}));
