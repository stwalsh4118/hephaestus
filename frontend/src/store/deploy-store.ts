import { create } from "zustand";

import {
  deployDiagram,
  teardownDiagram,
  updateDeploy as updateDeployApi,
} from "@/lib/deploy-api";
import { type DiagramJson } from "@/lib/topology";
import { connectStatusWs, disconnectStatusWs } from "@/lib/ws-client";
import type {
  ContainerStatus,
  DeployStatus,
  StatusMessage,
} from "@/types/deploy";

interface DeployStore {
  /** Overall deployment status. */
  deployStatus: DeployStatus;
  /** Per-node container statuses, keyed by canvas node ID. */
  nodeStatuses: Map<string, ContainerStatus>;
  /** Error message from the last failed operation. */
  error: string | null;
  /** The diagram snapshot at the time of last deploy. */
  lastDeployedDiagram: DiagramJson | null;
  /** Whether an update deploy request is in-flight. */
  isUpdating: boolean;

  /** Deploy the given diagram. */
  deploy: (diagram: DiagramJson) => Promise<void>;
  /** Teardown all deployed containers. */
  teardown: () => Promise<void>;
  /** Update deployed topology incrementally. */
  updateDeploy: (diagram: DiagramJson) => Promise<void>;
  /** Update state from a WebSocket status message. */
  updateFromStatusMessage: (msg: StatusMessage) => void;
  /** Reset all state. */
  reset: () => void;
  /** Dismiss the current error. */
  clearError: () => void;
}

export const useDeployStore = create<DeployStore>()((set, get) => ({
  deployStatus: "idle",
  nodeStatuses: new Map(),
  error: null,
  lastDeployedDiagram: null,
  isUpdating: false,

  deploy: async (diagram: DiagramJson) => {
    set({ deployStatus: "deploying", error: null });

    // Connect WebSocket for real-time updates.
    connectStatusWs((msg) => {
      get().updateFromStatusMessage(msg);
    });

    try {
      await deployDiagram(diagram);
      set({ lastDeployedDiagram: diagram });
    } catch (err) {
      const message = err instanceof Error ? err.message : "Deploy failed";
      set({ deployStatus: "error", error: message });
      disconnectStatusWs();
    }
  },

  teardown: async () => {
    set({ deployStatus: "tearing_down", error: null });

    try {
      await teardownDiagram();
      disconnectStatusWs();
      set({
        deployStatus: "idle",
        nodeStatuses: new Map(),
        lastDeployedDiagram: null,
      });
    } catch (err) {
      const message = err instanceof Error ? err.message : "Teardown failed";
      set({ deployStatus: "error", error: message });
    }
  },

  updateDeploy: async (diagram: DiagramJson) => {
    if (get().isUpdating) return;
    set({ isUpdating: true, error: null });

    try {
      await updateDeployApi(diagram);
      set({ lastDeployedDiagram: diagram, isUpdating: false });
    } catch (err) {
      const message = err instanceof Error ? err.message : "Update deploy failed";
      set({ error: message, isUpdating: false });
    }
  },

  updateFromStatusMessage: (msg: StatusMessage) => {
    const nodeStatuses = new Map<string, ContainerStatus>();
    for (const ns of msg.nodeStatuses) {
      nodeStatuses.set(ns.nodeId, ns.status);
    }
    set({ deployStatus: msg.deployStatus, nodeStatuses });
  },

  reset: () => {
    disconnectStatusWs();
    set({
      deployStatus: "idle",
      nodeStatuses: new Map(),
      error: null,
      lastDeployedDiagram: null,
      isUpdating: false,
    });
  },

  clearError: () => {
    set({ error: null });
  },
}));
