"use client";

import { useCanvasStore } from "@/store/canvas-store";
import { ST_REDIS } from "@/types/canvas";
import type { RedisConfig } from "@/types/canvas";

const EVICTION_POLICIES = ["noeviction", "allkeys-lru", "volatile-lru", "allkeys-random"] as const;
const DEFAULT_MAX_MEMORY = "256mb";
const DEFAULT_EVICTION_POLICY = "noeviction";

interface RedisFormProps {
  nodeId: string;
  config?: RedisConfig;
}

export function RedisForm({ nodeId, config }: RedisFormProps) {
  const updateNodeConfig = useCanvasStore((state) => state.updateNodeConfig);

  const maxMemory = config?.maxMemory ?? DEFAULT_MAX_MEMORY;
  const evictionPolicy = config?.evictionPolicy ?? DEFAULT_EVICTION_POLICY;

  const handleChange = (field: keyof Omit<RedisConfig, "type">, value: string) => {
    updateNodeConfig(nodeId, {
      type: ST_REDIS,
      maxMemory: field === "maxMemory" ? value : maxMemory,
      evictionPolicy: field === "evictionPolicy" ? value : evictionPolicy,
    });
  };

  return (
    <div className="space-y-4">
      <div>
        <label className="block text-xs font-medium text-slate-600">Max Memory</label>
        <input
          type="text"
          value={maxMemory}
          onChange={(e) => handleChange("maxMemory", e.target.value)}
          placeholder="e.g. 256mb"
          className="mt-1 w-full rounded border border-slate-300 px-2 py-1.5 text-sm text-slate-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
        />
      </div>

      <div>
        <label className="block text-xs font-medium text-slate-600">Eviction Policy</label>
        <select
          value={evictionPolicy}
          onChange={(e) => handleChange("evictionPolicy", e.target.value)}
          className="mt-1 w-full rounded border border-slate-300 bg-white px-2 py-1.5 text-sm text-slate-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
        >
          {EVICTION_POLICIES.map((p) => (
            <option key={p} value={p}>
              {p}
            </option>
          ))}
        </select>
      </div>
    </div>
  );
}
