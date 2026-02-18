"use client";

import { Handle, Position } from "@xyflow/react";
import type { NodeProps } from "@xyflow/react";

import { SERVICE_COLORS, SERVICE_ICONS, SERVICE_LABELS } from "@/constants/canvas";
import type { CanvasNode } from "@/types/canvas";

const FALLBACK_COLOR = "#64748b";
const FALLBACK_ICON = "?";
const FALLBACK_LABEL = "Unknown";

export function ServiceNode({ data, selected }: NodeProps<CanvasNode>) {
  const color = SERVICE_COLORS[data.type] ?? FALLBACK_COLOR;
  const icon = SERVICE_ICONS[data.type] ?? FALLBACK_ICON;
  const defaultLabel = SERVICE_LABELS[data.type] ?? FALLBACK_LABEL;

  return (
    <div
      className={`min-w-[140px] rounded-lg border-2 bg-white shadow-md transition-shadow ${
        selected ? "shadow-lg ring-2 ring-blue-400" : ""
      }`}
      style={{ borderColor: color }}
    >
      <div
        className="flex items-center gap-2 rounded-t-md px-3 py-2"
        style={{ backgroundColor: color }}
      >
        <span className="inline-flex h-6 w-6 items-center justify-center rounded text-[10px] font-bold text-white opacity-90">
          {icon}
        </span>
        <span className="truncate text-xs font-semibold text-white">
          {data.label || defaultLabel}
        </span>
      </div>

      <div className="px-3 py-2">
        <span className="text-[10px] text-slate-500">{defaultLabel}</span>
      </div>

      <Handle
        type="target"
        position={Position.Left}
        className="!h-3 !w-3 !border-2 !bg-white"
        style={{ borderColor: color }}
      />
      <Handle
        type="source"
        position={Position.Right}
        className="!h-3 !w-3 !border-2 !bg-white"
        style={{ borderColor: color }}
      />
    </div>
  );
}
