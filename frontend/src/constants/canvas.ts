import {
  ST_API_SERVICE,
  ST_NGINX,
  ST_POSTGRESQL,
  ST_RABBITMQ,
  ST_REDIS,
} from "@/types/canvas";
import type { CanvasViewport, PaletteItem, ServiceType } from "@/types/canvas";

export const DEFAULT_CANVAS_VIEWPORT: CanvasViewport = {
  x: 0,
  y: 0,
  zoom: 1,
};

export const CANVAS_GRID_SPACING = 16;
export const CANVAS_GRID_STROKE_WIDTH = 1;
export const CANVAS_MIN_ZOOM = 0.5;
export const CANVAS_MAX_ZOOM = 2;
export const CANVAS_DROP_DATA_KEY = "application/hephaestus-node-type";

export const EDGE_TYPE_LABELED = "labeled-edge";
export const EDGE_DEFAULT_LABEL = "";

export const SERVICE_COLORS: Record<ServiceType, string> = {
  [ST_API_SERVICE]: "#3b82f6",
  [ST_POSTGRESQL]: "#336791",
  [ST_REDIS]: "#dc2626",
  [ST_NGINX]: "#009639",
  [ST_RABBITMQ]: "#ff6600",
};

export const SERVICE_ICONS: Record<ServiceType, string> = {
  [ST_API_SERVICE]: "API",
  [ST_POSTGRESQL]: "PG",
  [ST_REDIS]: "RD",
  [ST_NGINX]: "NX",
  [ST_RABBITMQ]: "MQ",
};

export const SERVICE_LABELS: Record<ServiceType, string> = {
  [ST_API_SERVICE]: "API Service",
  [ST_POSTGRESQL]: "PostgreSQL",
  [ST_REDIS]: "Redis",
  [ST_NGINX]: "Nginx",
  [ST_RABBITMQ]: "RabbitMQ",
};

export const PALETTE_ITEMS: PaletteItem[] = [
  {
    id: ST_API_SERVICE,
    label: SERVICE_LABELS[ST_API_SERVICE],
    icon: SERVICE_ICONS[ST_API_SERVICE],
    description: "RESTful API service with configurable endpoints",
  },
  {
    id: ST_POSTGRESQL,
    label: SERVICE_LABELS[ST_POSTGRESQL],
    icon: SERVICE_ICONS[ST_POSTGRESQL],
    description: "Relational database with SQL support",
  },
  {
    id: ST_REDIS,
    label: SERVICE_LABELS[ST_REDIS],
    icon: SERVICE_ICONS[ST_REDIS],
    description: "In-memory data store and cache",
  },
  {
    id: ST_NGINX,
    label: SERVICE_LABELS[ST_NGINX],
    icon: SERVICE_ICONS[ST_NGINX],
    description: "Reverse proxy and load balancer",
  },
  {
    id: ST_RABBITMQ,
    label: SERVICE_LABELS[ST_RABBITMQ],
    icon: SERVICE_ICONS[ST_RABBITMQ],
    description: "Message broker for async communication",
  },
];
