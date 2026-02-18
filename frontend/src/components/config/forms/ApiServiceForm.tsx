"use client";

import { useCanvasStore } from "@/store/canvas-store";
import type { ApiServiceConfig, Endpoint, HttpMethod } from "@/types/canvas";

const HTTP_METHODS: HttpMethod[] = ["GET", "POST", "PUT", "DELETE", "PATCH"];
const DEFAULT_PORT = 8080;

const createEmptyEndpoint = (): Endpoint => ({
  method: "GET",
  path: "",
  responseSchema: "",
});

interface ApiServiceFormProps {
  nodeId: string;
  config?: ApiServiceConfig;
}

export function ApiServiceForm({ nodeId, config }: ApiServiceFormProps) {
  const updateNodeConfig = useCanvasStore((state) => state.updateNodeConfig);

  const port = config?.port ?? DEFAULT_PORT;
  const endpoints = config?.endpoints ?? [];

  const pushConfig = (newPort: number, newEndpoints: Endpoint[]) => {
    updateNodeConfig(nodeId, {
      type: "api-service",
      port: newPort,
      endpoints: newEndpoints,
    });
  };

  const handlePortChange = (value: string) => {
    const parsed = Number.parseInt(value, 10);
    if (!Number.isNaN(parsed)) {
      pushConfig(parsed, endpoints);
    }
  };

  const handleAddEndpoint = () => {
    pushConfig(port, [...endpoints, createEmptyEndpoint()]);
  };

  const handleRemoveEndpoint = (index: number) => {
    pushConfig(
      port,
      endpoints.filter((_, i) => i !== index)
    );
  };

  const handleEndpointChange = (index: number, field: keyof Endpoint, value: string) => {
    pushConfig(
      port,
      endpoints.map((ep, i) => (i === index ? { ...ep, [field]: value } : ep))
    );
  };

  return (
    <div className="space-y-4">
      <div>
        <label className="block text-xs font-medium text-slate-600">Port</label>
        <input
          type="number"
          value={port}
          onChange={(e) => handlePortChange(e.target.value)}
          className="mt-1 w-full rounded border border-slate-300 px-2 py-1.5 text-sm text-slate-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
        />
      </div>

      <div>
        <label className="block text-xs font-medium text-slate-600">Endpoints</label>
        <div className="mt-2 space-y-3">
          {endpoints.map((endpoint, index) => (
            <div key={index} className="rounded-lg border border-slate-200 bg-slate-50 p-3">
              <div className="mb-2 flex items-center justify-between">
                <span className="text-[10px] font-medium uppercase tracking-wide text-slate-400">
                  Endpoint {index + 1}
                </span>
                <button
                  type="button"
                  onClick={() => handleRemoveEndpoint(index)}
                  className="rounded p-1 text-slate-400 transition-colors hover:bg-red-50 hover:text-red-500"
                  aria-label={`Remove endpoint ${index + 1}`}
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="14"
                    height="14"
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
              </div>

              <div className="flex gap-2">
                <select
                  value={endpoint.method}
                  onChange={(e) => handleEndpointChange(index, "method", e.target.value)}
                  className="w-24 rounded border border-slate-300 bg-white px-2 py-1.5 text-sm text-slate-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                >
                  {HTTP_METHODS.map((m) => (
                    <option key={m} value={m}>
                      {m}
                    </option>
                  ))}
                </select>
                <input
                  type="text"
                  value={endpoint.path}
                  onChange={(e) => handleEndpointChange(index, "path", e.target.value)}
                  placeholder="/users"
                  className="flex-1 rounded border border-slate-300 px-2 py-1.5 text-sm text-slate-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                />
              </div>

              <div className="mt-2">
                <label className="block text-[10px] font-medium text-slate-500">
                  Response Schema (JSON)
                </label>
                <textarea
                  value={endpoint.responseSchema}
                  onChange={(e) => handleEndpointChange(index, "responseSchema", e.target.value)}
                  placeholder='{ "id": 1, "name": "example" }'
                  rows={3}
                  className="mt-1 w-full rounded border border-slate-300 px-2 py-1.5 font-mono text-xs text-slate-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
                />
              </div>
            </div>
          ))}
        </div>

        <button
          type="button"
          onClick={handleAddEndpoint}
          className="mt-2 rounded border border-dashed border-slate-300 px-3 py-1.5 text-xs text-slate-600 transition-colors hover:border-slate-400 hover:text-slate-800"
        >
          + Add Endpoint
        </button>
      </div>
    </div>
  );
}
