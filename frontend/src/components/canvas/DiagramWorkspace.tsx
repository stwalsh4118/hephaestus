"use client";

import { ReactFlowProvider } from "@xyflow/react";

import { ComponentPalette } from "@/components/palette/ComponentPalette";

import { DiagramCanvas } from "./DiagramCanvas";

export function DiagramWorkspace() {
  return (
    <ReactFlowProvider>
      <main className="flex h-screen min-h-screen w-full bg-slate-100 text-slate-900">
        <ComponentPalette />

        <section className="flex flex-1 border-l border-slate-200 bg-white">
          <DiagramCanvas />
        </section>
      </main>
    </ReactFlowProvider>
  );
}
