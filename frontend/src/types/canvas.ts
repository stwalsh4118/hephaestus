import type { Edge, Node, Viewport } from "@xyflow/react";

export const ST_API_SERVICE = "api-service" as const;
export const ST_POSTGRESQL = "postgresql" as const;
export const ST_REDIS = "redis" as const;
export const ST_NGINX = "nginx" as const;
export const ST_RABBITMQ = "rabbitmq" as const;

export type ServiceType =
  | typeof ST_API_SERVICE
  | typeof ST_POSTGRESQL
  | typeof ST_REDIS
  | typeof ST_NGINX
  | typeof ST_RABBITMQ;

export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

export interface Endpoint {
  method: HttpMethod;
  path: string;
  responseSchema: string;
}

export interface ApiServiceConfig {
  type: typeof ST_API_SERVICE;
  endpoints: Endpoint[];
  port: number;
}

export interface PostgresqlConfig {
  type: typeof ST_POSTGRESQL;
  engine: string;
  version: string;
}

export interface RedisConfig {
  type: typeof ST_REDIS;
  maxMemory: string;
  evictionPolicy: string;
}

export interface NginxConfig {
  type: typeof ST_NGINX;
  upstreamServers: string[];
}

export interface RabbitMQConfig {
  type: typeof ST_RABBITMQ;
  vhost: string;
}

export type ServiceConfig =
  | ApiServiceConfig
  | PostgresqlConfig
  | RedisConfig
  | NginxConfig
  | RabbitMQConfig;

export interface PaletteItem {
  id: ServiceType;
  label: string;
  icon: string;
  description: string;
}

export interface CanvasNodeData extends Record<string, unknown> {
  label: string;
  type: ServiceType;
  description: string;
  config?: ServiceConfig;
}

export interface CanvasEdgeData extends Record<string, unknown> {
  label: string;
}

export type CanvasNode = Node<CanvasNodeData>;
export type CanvasEdge = Edge<CanvasEdgeData>;
export type CanvasViewport = Viewport;
