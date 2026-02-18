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

export const SERVICE_COLORS: Record<ServiceType, string> = {
  "api-service": "#3b82f6",
  postgresql: "#336791",
  redis: "#dc2626",
  nginx: "#009639",
  rabbitmq: "#ff6600",
};

export const SERVICE_ICONS: Record<ServiceType, string> = {
  "api-service": "API",
  postgresql: "PG",
  redis: "RD",
  nginx: "NX",
  rabbitmq: "MQ",
};

export const SERVICE_LABELS: Record<ServiceType, string> = {
  "api-service": "API Service",
  postgresql: "PostgreSQL",
  redis: "Redis",
  nginx: "Nginx",
  rabbitmq: "RabbitMQ",
};

export const PALETTE_ITEMS: PaletteItem[] = [
  {
    id: "api-service",
    label: SERVICE_LABELS["api-service"],
    icon: SERVICE_ICONS["api-service"],
    description: "RESTful API service with configurable endpoints",
  },
  {
    id: "postgresql",
    label: SERVICE_LABELS["postgresql"],
    icon: SERVICE_ICONS["postgresql"],
    description: "Relational database with SQL support",
  },
  {
    id: "redis",
    label: SERVICE_LABELS["redis"],
    icon: SERVICE_ICONS["redis"],
    description: "In-memory data store and cache",
  },
  {
    id: "nginx",
    label: SERVICE_LABELS["nginx"],
    icon: SERVICE_ICONS["nginx"],
    description: "Reverse proxy and load balancer",
  },
  {
    id: "rabbitmq",
    label: SERVICE_LABELS["rabbitmq"],
    icon: SERVICE_ICONS["rabbitmq"],
    description: "Message broker for async communication",
  },
];
