/** Overall deployment state. */
export type DeployStatus =
  | "idle"
  | "deploying"
  | "deployed"
  | "tearing_down"
  | "error";

/** Per-container status matching the backend docker.ContainerStatus values. */
export type ContainerStatus =
  | "created"
  | "running"
  | "stopped"
  | "error"
  | "healthy"
  | "unhealthy";

/** Status of a single deployed node. */
export interface NodeStatus {
  nodeId: string;
  containerId: string;
  status: ContainerStatus;
}

/** WebSocket message format for real-time status updates. */
export interface StatusMessage {
  type: "status_update";
  deployStatus: DeployStatus;
  nodeStatuses: NodeStatus[];
}

/** HTTP response from GET /api/deploy/status. */
export interface DeployStatusResponse {
  deployStatus: DeployStatus;
  nodeStatuses: NodeStatus[];
}

/** HTTP response from POST /api/deploy and DELETE /api/deploy. */
export interface DeployResponse {
  status: string;
}
