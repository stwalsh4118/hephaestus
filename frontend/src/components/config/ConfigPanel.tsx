"use client";

import { SERVICE_COLORS, SERVICE_ICONS, SERVICE_LABELS } from "@/constants/canvas";
import { useCanvasStore, useSelectedNode } from "@/store/canvas-store";
import {
  ST_API_SERVICE,
  ST_NGINX,
  ST_POSTGRESQL,
  ST_RABBITMQ,
  ST_REDIS,
} from "@/types/canvas";
import type {
  ApiServiceConfig,
  NginxConfig,
  PostgresqlConfig,
  RabbitMQConfig,
  RedisConfig,
  ServiceType,
} from "@/types/canvas";

import { ApiServiceForm } from "./forms/ApiServiceForm";
import { NginxForm } from "./forms/NginxForm";
import { PostgresqlForm } from "./forms/PostgresqlForm";
import { RabbitmqForm } from "./forms/RabbitmqForm";
import { RedisForm } from "./forms/RedisForm";

const PANEL_WIDTH = 320;

function ServiceConfigForm({
  nodeId,
  serviceType,
  config,
}: {
  nodeId: string;
  serviceType: ServiceType;
  config?: unknown;
}) {
  switch (serviceType) {
    case ST_POSTGRESQL:
      return <PostgresqlForm nodeId={nodeId} config={config as PostgresqlConfig | undefined} />;
    case ST_REDIS:
      return <RedisForm nodeId={nodeId} config={config as RedisConfig | undefined} />;
    case ST_NGINX:
      return <NginxForm nodeId={nodeId} config={config as NginxConfig | undefined} />;
    case ST_RABBITMQ:
      return <RabbitmqForm nodeId={nodeId} config={config as RabbitMQConfig | undefined} />;
    case ST_API_SERVICE:
      return <ApiServiceForm nodeId={nodeId} config={config as ApiServiceConfig | undefined} />;
    default:
      return null;
  }
}

export function ConfigPanel() {
  const selectedNode = useSelectedNode();
  const selectNode = useCanvasStore((state) => state.selectNode);
  const updateNodeLabel = useCanvasStore((state) => state.updateNodeLabel);

  const isOpen = selectedNode !== null;
  const serviceType = selectedNode?.data.type;
  const color = serviceType ? SERVICE_COLORS[serviceType] : undefined;
  const icon = serviceType ? SERVICE_ICONS[serviceType] : undefined;
  const label = serviceType ? SERVICE_LABELS[serviceType] : undefined;

  return (
    <aside
      className={`flex h-full shrink-0 flex-col bg-white transition-all duration-200 ease-in-out ${isOpen ? "border-l border-slate-200" : ""}`}
      style={{ width: isOpen ? PANEL_WIDTH : 0, minWidth: isOpen ? PANEL_WIDTH : 0 }}
    >
      {isOpen && selectedNode && (
        <div className="flex h-full flex-col overflow-hidden" style={{ width: PANEL_WIDTH }}>
          <header
            className="flex items-center justify-between px-4 py-3"
            style={{ backgroundColor: color }}
          >
            <div className="flex items-center gap-2">
              <span className="inline-flex h-6 w-6 items-center justify-center rounded text-[10px] font-bold text-white opacity-90">
                {icon}
              </span>
              <span className="text-sm font-semibold text-white">{label}</span>
            </div>
            <button
              type="button"
              onClick={() => selectNode(null)}
              className="rounded p-1 text-white/80 transition-colors hover:bg-white/20 hover:text-white"
              aria-label="Close configuration panel"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <line x1="18" y1="6" x2="6" y2="18" />
                <line x1="6" y1="6" x2="18" y2="18" />
              </svg>
            </button>
          </header>

          <div className="border-b border-slate-200 px-4 py-3">
            <label className="block text-xs font-medium text-slate-500">Node Name</label>
            <input
              type="text"
              value={selectedNode.data.label}
              onChange={(e) => updateNodeLabel(selectedNode.id, e.target.value)}
              className="mt-1 w-full rounded border border-slate-300 px-2 py-1 text-sm text-slate-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>

          <div className="flex-1 overflow-y-auto px-4 py-3">
            <ServiceConfigForm
              nodeId={selectedNode.id}
              serviceType={selectedNode.data.type}
              config={selectedNode.data.config}
            />
          </div>
        </div>
      )}
    </aside>
  );
}
