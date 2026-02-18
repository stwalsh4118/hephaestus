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

### Node Data

```typescript
interface CanvasNodeData extends Record<string, unknown> {
  label: string;
  type: ServiceType;
  description: string;
  config?: ServiceConfig;
}

type CanvasNode = Node<CanvasNodeData>;
type CanvasEdge = Edge;
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
