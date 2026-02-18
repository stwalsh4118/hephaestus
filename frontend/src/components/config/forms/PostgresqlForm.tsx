"use client";

import { useCanvasStore } from "@/store/canvas-store";
import { ST_POSTGRESQL } from "@/types/canvas";
import type { PostgresqlConfig } from "@/types/canvas";

const POSTGRESQL_VERSIONS = ["16", "15", "14"] as const;
const DEFAULT_ENGINE = "PostgreSQL";
const DEFAULT_VERSION = "16";

interface PostgresqlFormProps {
  nodeId: string;
  config?: PostgresqlConfig;
}

export function PostgresqlForm({ nodeId, config }: PostgresqlFormProps) {
  const updateNodeConfig = useCanvasStore((state) => state.updateNodeConfig);

  const engine = config?.engine ?? DEFAULT_ENGINE;
  const version = config?.version ?? DEFAULT_VERSION;

  const handleChange = (field: keyof Omit<PostgresqlConfig, "type">, value: string) => {
    updateNodeConfig(nodeId, {
      type: ST_POSTGRESQL,
      engine: field === "engine" ? value : engine,
      version: field === "version" ? value : version,
    });
  };

  return (
    <div className="space-y-4">
      <div>
        <label className="block text-xs font-medium text-slate-600">Engine</label>
        <select
          value={engine}
          onChange={(e) => handleChange("engine", e.target.value)}
          className="mt-1 w-full rounded border border-slate-300 bg-white px-2 py-1.5 text-sm text-slate-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
        >
          <option value="PostgreSQL">PostgreSQL</option>
        </select>
      </div>

      <div>
        <label className="block text-xs font-medium text-slate-600">Version</label>
        <select
          value={version}
          onChange={(e) => handleChange("version", e.target.value)}
          className="mt-1 w-full rounded border border-slate-300 bg-white px-2 py-1.5 text-sm text-slate-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
        >
          {POSTGRESQL_VERSIONS.map((v) => (
            <option key={v} value={v}>
              {v}
            </option>
          ))}
        </select>
      </div>
    </div>
  );
}
