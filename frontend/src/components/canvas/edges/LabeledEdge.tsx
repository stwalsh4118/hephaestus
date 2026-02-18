"use client";

import {
  BaseEdge,
  EdgeLabelRenderer,
  type EdgeProps,
  getSmoothStepPath,
} from "@xyflow/react";
import { useCallback, useRef, useState } from "react";

import { useCanvasStore } from "@/store/canvas-store";
import type { CanvasEdge } from "@/types/canvas";

export function LabeledEdge({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  markerEnd,
  data,
}: EdgeProps<CanvasEdge>) {
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState("");
  const cancelledRef = useRef(false);

  const [edgePath, labelX, labelY] = getSmoothStepPath({
    sourceX,
    sourceY,
    targetX,
    targetY,
    sourcePosition,
    targetPosition,
  });

  const label = data?.label ?? "";
  const updateEdgeLabel = useCanvasStore((state) => state.updateEdgeLabel);

  const commitLabel = useCallback(
    (value: string) => {
      updateEdgeLabel(id, value);
      setIsEditing(false);
    },
    [id, updateEdgeLabel]
  );

  const handleDoubleClick = useCallback(() => {
    cancelledRef.current = false;
    setEditValue(label);
    setIsEditing(true);
  }, [label]);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === "Enter") {
        commitLabel(editValue);
      }
      if (e.key === "Escape") {
        cancelledRef.current = true;
        setIsEditing(false);
      }
    },
    [commitLabel, editValue]
  );

  const handleBlur = useCallback(() => {
    if (!cancelledRef.current) {
      commitLabel(editValue);
    }
    cancelledRef.current = false;
  }, [commitLabel, editValue]);

  const inputRefCallback = useCallback((el: HTMLInputElement | null) => {
    el?.focus();
  }, []);

  return (
    <>
      <BaseEdge
        id={id}
        path={edgePath}
        markerEnd={markerEnd}
      />
      <EdgeLabelRenderer>
        {/* nodrag nopan: prevent label interactions from moving the canvas */}
        <div
          className="nodrag nopan absolute"
          style={{
            transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            pointerEvents: "all",
          }}
          onDoubleClick={handleDoubleClick}
        >
          {isEditing ? (
            <input
              ref={inputRefCallback}
              className="rounded border border-blue-400 bg-white px-1 py-0.5 text-xs outline-none"
              value={editValue}
              onChange={(e) => setEditValue(e.target.value)}
              onKeyDown={handleKeyDown}
              onBlur={handleBlur}
            />
          ) : label ? (
            <span className="rounded bg-white px-1.5 py-0.5 text-xs text-slate-700 shadow-sm">
              {label}
            </span>
          ) : (
            <span className="rounded px-3 py-1.5 text-xs text-slate-400 hover:bg-slate-100 hover:text-slate-500">
              +
            </span>
          )}
        </div>
      </EdgeLabelRenderer>
    </>
  );
}
