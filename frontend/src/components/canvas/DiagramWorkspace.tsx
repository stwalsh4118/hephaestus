"use client";

import { useCallback, useMemo, useRef } from "react";
import type { ChangeEvent } from "react";

import { ReactFlowProvider } from "@xyflow/react";

import { ConfigPanel } from "@/components/config/ConfigPanel";
import { ComponentPalette } from "@/components/palette/ComponentPalette";
import { ErrorToast } from "@/components/ui/ErrorToast";
import { hasDiagramChanges } from "@/lib/diagram-diff";
import {
  downloadDiagram,
  exportDiagram,
  importDiagram,
} from "@/lib/topology";
import { useCanvasStore } from "@/store/canvas-store";
import { useDeployStore } from "@/store/deploy-store";

import { DiagramCanvas } from "./DiagramCanvas";

export function DiagramWorkspace() {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const nodes = useCanvasStore((s) => s.nodes);
  const deployStatus = useDeployStore((s) => s.deployStatus);
  const deployError = useDeployStore((s) => s.error);
  const clearError = useDeployStore((s) => s.clearError);

  const lastDeployedDiagram = useDeployStore((s) => s.lastDeployedDiagram);
  const isUpdating = useDeployStore((s) => s.isUpdating);

  const isDeploying = deployStatus === "deploying";
  const isDeployed = deployStatus === "deployed";
  const isTearingDown = deployStatus === "tearing_down";
  const canDeploy = nodes.length > 0 && !isDeploying && !isDeployed && !isTearingDown;
  const isError = deployStatus === "error";
  const canTeardown = (isDeployed || isError) && !isTearingDown;

  const hasChanges = useMemo(
    () => isDeployed && hasDiagramChanges(nodes, lastDeployedDiagram),
    [isDeployed, lastDeployedDiagram, nodes]
  );

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

  const handleDeploy = useCallback(() => {
    const { nodes, edges } = useCanvasStore.getState();
    const diagram = exportDiagram(nodes, edges);
    useDeployStore.getState().deploy(diagram);
  }, []);

  const handleTeardown = useCallback(() => {
    const confirmed = window.confirm(
      "Are you sure you want to tear down all deployed containers?"
    );
    if (confirmed) {
      useDeployStore.getState().teardown();
    }
  }, []);

  const handleUpdateDeploy = useCallback(() => {
    const { nodes, edges } = useCanvasStore.getState();
    const diagram = exportDiagram(nodes, edges);
    useDeployStore.getState().updateDeploy(diagram);
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

            <div className="mx-1 h-5 w-px bg-slate-300" />

            <button
              type="button"
              onClick={handleDeploy}
              disabled={!canDeploy}
              className="inline-flex items-center gap-1.5 rounded bg-emerald-600 px-3 py-1.5 text-xs font-medium text-white hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {isDeploying && (
                <span className="inline-block h-3 w-3 animate-spin rounded-full border-2 border-white border-t-transparent" />
              )}
              {isDeploying ? "Deploying..." : "Deploy"}
            </button>

            {hasChanges && (
              <button
                type="button"
                onClick={handleUpdateDeploy}
                disabled={isUpdating}
                className="inline-flex items-center gap-1.5 rounded bg-amber-500 px-3 py-1.5 text-xs font-medium text-white hover:bg-amber-600 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {isUpdating && (
                  <span className="inline-block h-3 w-3 animate-spin rounded-full border-2 border-white border-t-transparent" />
                )}
                {isUpdating ? "Updating..." : "Update Deploy"}
              </button>
            )}

            {(isDeployed || isTearingDown || isError) && (
              <button
                type="button"
                onClick={handleTeardown}
                disabled={!canTeardown}
                className="inline-flex items-center gap-1.5 rounded bg-red-600 px-3 py-1.5 text-xs font-medium text-white hover:bg-red-700 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {isTearingDown && (
                  <span className="inline-block h-3 w-3 animate-spin rounded-full border-2 border-white border-t-transparent" />
                )}
                {isTearingDown ? "Tearing down..." : "Teardown"}
              </button>
            )}
          </div>
          <div className="flex-1">
            <DiagramCanvas />
          </div>
        </section>

        <ConfigPanel />
      </main>

      {deployError && (
        <ErrorToast message={deployError} onDismiss={clearError} />
      )}
    </ReactFlowProvider>
  );
}
