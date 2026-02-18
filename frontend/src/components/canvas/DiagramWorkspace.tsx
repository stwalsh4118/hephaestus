"use client";

import { useCallback, useRef } from "react";
import type { ChangeEvent } from "react";

import { ReactFlowProvider } from "@xyflow/react";

import { ConfigPanel } from "@/components/config/ConfigPanel";
import { ComponentPalette } from "@/components/palette/ComponentPalette";
import {
  downloadDiagram,
  exportDiagram,
  importDiagram,
} from "@/lib/topology";
import { useCanvasStore } from "@/store/canvas-store";

import { DiagramCanvas } from "./DiagramCanvas";

export function DiagramWorkspace() {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleExport = useCallback(() => {
    const { nodes, edges } = useCanvasStore.getState();
    const diagram = exportDiagram(nodes, edges);
    downloadDiagram(diagram);
  }, []);

  const handleImportClick = useCallback(() => {
    fileInputRef.current?.click();
  }, []);

  const handleFileChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      try {
        const json: unknown = JSON.parse(event.target?.result as string);
        const { nodes, edges } = importDiagram(json);
        useCanvasStore.getState().loadDiagram(nodes, edges);
      } catch (err) {
        console.error("Failed to import diagram:", err);
      }
    };
    reader.readAsText(file);

    // Reset input so the same file can be re-imported
    e.target.value = "";
  }, []);

  return (
    <ReactFlowProvider>
      <main className="flex h-screen min-h-screen w-full bg-slate-100 text-slate-900">
        <ComponentPalette />

        <section className="relative flex flex-1 flex-col border-l border-slate-200 bg-white">
          <div className="flex items-center gap-2 border-b border-slate-200 px-4 py-2">
            <button
              type="button"
              onClick={handleExport}
              className="rounded bg-slate-700 px-3 py-1.5 text-xs font-medium text-white hover:bg-slate-800"
            >
              Export JSON
            </button>
            <button
              type="button"
              onClick={handleImportClick}
              className="rounded bg-slate-700 px-3 py-1.5 text-xs font-medium text-white hover:bg-slate-800"
            >
              Import JSON
            </button>
            <input
              ref={fileInputRef}
              type="file"
              accept=".json"
              onChange={handleFileChange}
              className="hidden"
            />
          </div>
          <div className="flex-1">
            <DiagramCanvas />
          </div>
        </section>

        <ConfigPanel />
      </main>
    </ReactFlowProvider>
  );
}
