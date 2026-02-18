import {
  applyEdgeChanges,
  applyNodeChanges,
  type EdgeChange,
  type NodeChange,
  type NodeRemoveChange,
  type XYPosition,
} from "@xyflow/react";
import { create } from "zustand";

import { NODE_TYPE_SERVICE } from "@/components/canvas/nodes";
import { DEFAULT_CANVAS_VIEWPORT } from "@/constants/canvas";
import type {
  CanvasEdge,
  CanvasNode,
  CanvasNodeData,
  CanvasViewport,
  ServiceConfig,
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
  selectedNodeId: string | null;
  onNodesChange: (changes: NodeChange<CanvasNode>[]) => void;
  onEdgesChange: (changes: EdgeChange<CanvasEdge>[]) => void;
  onViewportChange: (viewport: CanvasViewport) => void;
  addNode: (input: AddNodeInput) => void;
  selectNode: (nodeId: string | null) => void;
  updateNodeLabel: (nodeId: string, label: string) => void;
  updateNodeConfig: (nodeId: string, config: ServiceConfig) => void;
}

export const useCanvasStore = create<CanvasStore>()((set) => ({
  nodes: INITIAL_NODES,
  edges: INITIAL_EDGES,
  viewport: INITIAL_VIEWPORT,
  selectedNodeId: null,
  onNodesChange: (changes) => {
    set((state) => {
      const newNodes = applyNodeChanges(changes, state.nodes);
      const removedIds = changes
        .filter((c): c is NodeRemoveChange => c.type === "remove")
        .map((c) => c.id);
      const shouldClearSelection =
        state.selectedNodeId !== null && removedIds.includes(state.selectedNodeId);
      return {
        nodes: newNodes,
        ...(shouldClearSelection ? { selectedNodeId: null } : {}),
      };
    });
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
          type: NODE_TYPE_SERVICE,
        },
      ],
    }));
  },
  selectNode: (nodeId) => {
    set({ selectedNodeId: nodeId });
  },
  updateNodeLabel: (nodeId, label) => {
    set((state) => ({
      nodes: state.nodes.map((node) =>
        node.id === nodeId ? { ...node, data: { ...node.data, label } } : node
      ),
    }));
  },
  updateNodeConfig: (nodeId, config) => {
    set((state) => ({
      nodes: state.nodes.map((node) =>
        node.id === nodeId ? { ...node, data: { ...node.data, config } } : node
      ),
    }));
  },
}));

export const useSelectedNode = (): CanvasNode | null => {
  return useCanvasStore((state) => {
    if (state.selectedNodeId === null) return null;
    return state.nodes.find((n) => n.id === state.selectedNodeId) ?? null;
  });
};
