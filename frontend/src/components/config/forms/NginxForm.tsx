"use client";

import { useCanvasStore } from "@/store/canvas-store";
import { ST_NGINX } from "@/types/canvas";
import type { NginxConfig } from "@/types/canvas";

interface NginxFormProps {
  nodeId: string;
  config?: NginxConfig;
}

export function NginxForm({ nodeId, config }: NginxFormProps) {
  const updateNodeConfig = useCanvasStore((state) => state.updateNodeConfig);

  const servers = config?.upstreamServers ?? [];

  const updateServers = (newServers: string[]) => {
    updateNodeConfig(nodeId, {
      type: ST_NGINX,
      upstreamServers: newServers,
    });
  };

  const handleAdd = () => {
    updateServers([...servers, ""]);
  };

  const handleRemove = (index: number) => {
    updateServers(servers.filter((_, i) => i !== index));
  };

  const handleChange = (index: number, value: string) => {
    updateServers(servers.map((s, i) => (i === index ? value : s)));
  };

  return (
    <div className="space-y-4">
      <div>
        <label className="block text-xs font-medium text-slate-600">Upstream Servers</label>
        <div className="mt-1 space-y-2">
          {servers.map((server, index) => (
            <div key={index} className="flex items-center gap-2">
              <input
                type="text"
                value={server}
                onChange={(e) => handleChange(index, e.target.value)}
                placeholder="e.g. 127.0.0.1:3000"
                className="flex-1 rounded border border-slate-300 px-2 py-1.5 text-sm text-slate-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
              />
              <button
                type="button"
                onClick={() => handleRemove(index)}
                className="rounded p-1 text-slate-400 transition-colors hover:bg-red-50 hover:text-red-500"
                aria-label="Remove server"
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
          ))}
        </div>
        <button
          type="button"
          onClick={handleAdd}
          className="mt-2 rounded border border-dashed border-slate-300 px-3 py-1.5 text-xs text-slate-600 transition-colors hover:border-slate-400 hover:text-slate-800"
        >
          + Add Server
        </button>
      </div>
    </div>
  );
}
