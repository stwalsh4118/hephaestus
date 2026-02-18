import type { Edge, Node, Viewport } from "@xyflow/react";

export type ServiceType = "api-service" | "postgresql" | "redis" | "nginx" | "rabbitmq";

export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

export interface Endpoint {
  method: HttpMethod;
  path: string;
  responseSchema: string;
}

export interface ApiServiceConfig {
  type: "api-service";
  endpoints: Endpoint[];
  port: number;
}

export interface PostgresqlConfig {
  type: "postgresql";
  engine: string;
  version: string;
}

export interface RedisConfig {
  type: "redis";
  maxMemory: string;
  evictionPolicy: string;
}

export interface NginxConfig {
  type: "nginx";
  upstreamServers: string[];
}

export interface RabbitMQConfig {
  type: "rabbitmq";
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

export type CanvasNode = Node<CanvasNodeData>;
export type CanvasEdge = Edge;
export type CanvasViewport = Viewport;
