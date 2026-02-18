"use client";

import type { DragEvent } from "react";

import { CANVAS_DROP_DATA_KEY, PALETTE_ITEMS } from "@/constants/canvas";
import type { NodeType } from "@/types/canvas";

const handleDragStart = (event: DragEvent<HTMLButtonElement>, type: NodeType) => {
  event.dataTransfer.setData(CANVAS_DROP_DATA_KEY, type);
  event.dataTransfer.effectAllowed = "copy";
};

export function ComponentPalette() {
  return (
    <aside className="flex h-full w-72 shrink-0 flex-col border-r border-slate-200 bg-slate-50">
      <header className="border-b border-slate-200 px-4 py-3">
        <h2 className="text-sm font-semibold uppercase tracking-wide text-slate-700">Components</h2>
        <p className="mt-1 text-xs text-slate-500">Drag items onto the canvas.</p>
      </header>

      <div className="flex-1 space-y-3 overflow-y-auto p-4">
        {PALETTE_ITEMS.map((item) => (
          <button
            key={item.id}
            type="button"
            draggable
            onDragStart={(event) => handleDragStart(event, item.id)}
            className="group w-full cursor-grab rounded-lg border border-slate-300 bg-white p-3 text-left shadow-sm transition hover:border-slate-400 hover:shadow active:cursor-grabbing"
          >
            <div className="flex items-start gap-3">
              <span className="inline-flex h-8 w-8 items-center justify-center rounded-md bg-slate-900 text-xs font-bold tracking-wide text-white">
                {item.icon}
              </span>

              <span className="min-w-0">
                <span className="block text-sm font-medium text-slate-900">{item.label}</span>
                <span className="mt-1 block text-xs text-slate-600">{item.description}</span>
              </span>
            </div>
          </button>
        ))}
      </div>
    </aside>
  );
}
