import type { CanvasViewport, PaletteItem } from "@/types/canvas";

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

export const PALETTE_ITEMS: PaletteItem[] = [
  {
    id: "service",
    label: "Service",
    icon: "SVC",
    description: "Generic stateless compute service",
  },
  {
    id: "worker",
    label: "Worker",
    icon: "WRK",
    description: "Background processing component",
  },
  {
    id: "database",
    label: "Database",
    icon: "DB",
    description: "Persistent data storage component",
  },
  {
    id: "queue",
    label: "Queue",
    icon: "Q",
    description: "Asynchronous message buffer",
  },
];
