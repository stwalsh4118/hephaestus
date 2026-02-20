# Frontend Canvas API Specification

## Types (`frontend/src/types/canvas.ts`)

### Service Types

```typescript
type ServiceType = "api-service" | "postgresql" | "redis" | "nginx" | "rabbitmq";
type HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
```

### Service Configuration Interfaces

```typescript
interface Endpoint {
  method: HttpMethod;
  path: string;
  responseSchema: string;
}

interface ApiServiceConfig {
  type: "api-service";
  endpoints: Endpoint[];
  port: number; // default 8080
}

interface PostgresqlConfig {
  type: "postgresql";
  engine: string;       // "PostgreSQL"
  version: string;      // "14" | "15" | "16"
}

interface RedisConfig {
  type: "redis";
  maxMemory: string;        // e.g. "256mb"
  evictionPolicy: string;   // "noeviction" | "allkeys-lru" | "volatile-lru" | "allkeys-random"
}

interface NginxConfig {
  type: "nginx";
  upstreamServers: string[];
}

interface RabbitMQConfig {
  type: "rabbitmq";
  vhost: string; // default "/"
}

type ServiceConfig = ApiServiceConfig | PostgresqlConfig | RedisConfig | NginxConfig | RabbitMQConfig;
```

### Node & Edge Data

```typescript
interface CanvasNodeData extends Record<string, unknown> {
  label: string;
  type: ServiceType;
  description: string;
  config?: ServiceConfig;
}

interface CanvasEdgeData extends Record<string, unknown> {
  label: string;
}

type CanvasNode = Node<CanvasNodeData>;
type CanvasEdge = Edge<CanvasEdgeData>;
type CanvasViewport = Viewport;
```

### Palette

```typescript
interface PaletteItem {
  id: ServiceType;
  label: string;
  icon: string;
  description: string;
}
```

## Constants (`frontend/src/constants/canvas.ts`)

```typescript
const SERVICE_COLORS: Record<ServiceType, string>;  // hex colour per service
const SERVICE_ICONS: Record<ServiceType, string>;    // 2-3 char abbreviation per service
const SERVICE_LABELS: Record<ServiceType, string>;   // display name per service
const PALETTE_ITEMS: PaletteItem[];                  // 5 items: api-service, postgresql, redis, nginx, rabbitmq
const CANVAS_DROP_DATA_KEY = "application/hephaestus-node-type";
const EDGE_TYPE_LABELED = "labeled-edge";
const EDGE_DEFAULT_LABEL = "";
```

## Canvas Store (`frontend/src/store/canvas-store.ts`)

### State

```typescript
interface CanvasStore {
  nodes: CanvasNode[];
  edges: CanvasEdge[];
  viewport: CanvasViewport;
  selectedNodeId: string | null;
}
```

### Actions

```typescript
onNodesChange(changes: NodeChange<CanvasNode>[]): void   // applies changes, clears selection on delete
onEdgesChange(changes: EdgeChange<CanvasEdge>[]): void
onViewportChange(viewport: CanvasViewport): void
addNode(input: { position: XYPosition; data: CanvasNodeData }): void
selectNode(nodeId: string | null): void
updateNodeLabel(nodeId: string, label: string): void
updateNodeConfig(nodeId: string, config: ServiceConfig): void  // full replacement, not merge
onConnect(connection: Connection): void              // creates edge, rejects duplicates (same source+target)
updateEdgeLabel(edgeId: string, label: string): void
removeEdge(edgeId: string): void
loadDiagram(nodes: CanvasNode[], edges: CanvasEdge[]): void  // replaces entire diagram, resets selection
```

### Hooks

```typescript
useCanvasStore: UseBoundStore<StoreApi<CanvasStore>>
useSelectedNode(): CanvasNode | null  // derived selector
```

## Custom Node Component (`frontend/src/components/canvas/nodes/`)

```typescript
const NODE_TYPE_SERVICE = "service-node";
const nodeTypes: NodeTypes = { [NODE_TYPE_SERVICE]: ServiceNode };
```

`ServiceNode` renders a coloured header (from `SERVICE_COLORS`), icon abbreviation, label, and source/target handles.

## Custom Edge Component (`frontend/src/components/canvas/edges/`)

```typescript
const edgeTypes: EdgeTypes = { [EDGE_TYPE_LABELED]: LabeledEdge };
```

`LabeledEdge` renders a smooth step path with:
- Directional arrowhead marker (`MarkerType.ArrowClosed`)
- Editable label at midpoint (double-click to edit, Enter to commit, Escape to cancel)
- Label hidden when empty

## Topology Export/Import (`frontend/src/lib/topology.ts`)

### DiagramJson Schema (PRD-conforming)

```typescript
interface DiagramJson {
  id: string;           // UUID, generated on export
  name: string;         // "Untitled Diagram"
  nodes: DiagramJsonNode[];
  edges: DiagramJsonEdge[];
}

interface DiagramJsonNode {
  id: string;
  type: string;         // ServiceType string (e.g., "api-service")
  name: string;         // maps from CanvasNodeData.label
  description: string;
  position: { x: number; y: number };
  config?: ServiceConfig;
}

interface DiagramJsonEdge {
  id: string;
  source: string;
  target: string;
  label: string;
}
```

### Functions

```typescript
exportDiagram(nodes: CanvasNode[], edges: CanvasEdge[]): DiagramJson
downloadDiagram(diagram: DiagramJson): void  // triggers browser file download
importDiagram(json: unknown): { nodes: CanvasNode[]; edges: CanvasEdge[] }  // validates and maps
```

## Config Panel (`frontend/src/components/config/ConfigPanel.tsx`)

- Opens when `selectedNodeId !== null` (width: 320px)
- Closed state: width 0px, no border
- Header: service colour background, icon, label, close button
- Body: editable node name input, service-specific config form

### Config Forms (`frontend/src/components/config/forms/`)

| Service Type | Component | Key Fields |
|---|---|---|
| `api-service` | `ApiServiceForm` | Port input, endpoint list (method, path, schema) with add/remove |
| `postgresql` | `PostgresqlForm` | Engine dropdown, version dropdown (16, 15, 14) |
| `redis` | `RedisForm` | Max memory input, eviction policy dropdown |
| `nginx` | `NginxForm` | Upstream servers list with add/remove |
| `rabbitmq` | `RabbitmqForm` | Virtual host input |

All forms call `updateNodeConfig` on every change for immediate persistence.

---

## Deploy Types (`frontend/src/types/deploy.ts`)

```typescript
type DeployStatus = "idle" | "deploying" | "deployed" | "tearing_down" | "error";
type ContainerStatus = "created" | "running" | "stopped" | "error" | "healthy" | "unhealthy";

interface NodeStatus { nodeId: string; containerId: string; status: ContainerStatus; }
interface StatusMessage { type: string; deployStatus: DeployStatus; nodeStatuses: NodeStatus[]; }
interface DeployResponse { status: string; }
interface DeployStatusResponse { deployStatus: DeployStatus; nodeStatuses: NodeStatus[]; }
```

## Deploy API Client (`frontend/src/lib/deploy-api.ts`)

```typescript
deployDiagram(diagram: DiagramJson): Promise<DeployResponse>     // POST /api/deploy
teardownDiagram(): Promise<DeployResponse>                        // DELETE /api/deploy
getDeployStatus(): Promise<DeployStatusResponse>                  // GET /api/deploy/status
updateDeploy(diagram: DiagramJson): Promise<DeployStatusResponse> // PUT /api/deploy
```

## WebSocket Client (`frontend/src/lib/ws-client.ts`)

```typescript
connectStatusWs(onMessage: (msg: StatusMessage) => void): WebSocket  // ws://host/ws/status
disconnectStatusWs(): void
```

## Deploy Store (`frontend/src/store/deploy-store.ts`)

Zustand store managing deployment lifecycle.

```typescript
interface DeployStore {
  deployStatus: DeployStatus;
  nodeStatuses: Map<string, ContainerStatus>;
  error: string | null;
  lastDeployedDiagram: DiagramJson | null;
  isUpdating: boolean;
  deploy(diagram: DiagramJson): Promise<void>;
  teardown(): Promise<void>;
  updateDeploy(diagram: DiagramJson): Promise<void>;
  clearError(): void;
  reset(): void;
}
```

## Diagram Diff (`frontend/src/lib/diagram-diff.ts`)

```typescript
hasDiagramChanges(nodes: CanvasNode[], lastDeployed: DiagramJson | null): boolean
```
