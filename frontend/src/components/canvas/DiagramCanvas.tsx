"use client";

import { useCallback } from "react";
import type { DragEvent } from "react";

import {
  Background,
  BackgroundVariant,
  Controls,
  MiniMap,
  ReactFlow,
  useReactFlow,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";

import {
  CANVAS_DROP_DATA_KEY,
  CANVAS_GRID_SPACING,
  CANVAS_GRID_STROKE_WIDTH,
  CANVAS_MAX_ZOOM,
  CANVAS_MIN_ZOOM,
  DEFAULT_CANVAS_VIEWPORT,
  EDGE_DEFAULT_LABEL,
  EDGE_TYPE_LABELED,
  PALETTE_ITEMS,
} from "@/constants/canvas";
import { useCanvasStore } from "@/store/canvas-store";
import type { ServiceType } from "@/types/canvas";

import { edgeTypes } from "./edges";
import { nodeTypes } from "./nodes";

export function DiagramCanvas() {
  const { screenToFlowPosition } = useReactFlow();
  const nodes = useCanvasStore((state) => state.nodes);
  const edges = useCanvasStore((state) => state.edges);
  const viewport = useCanvasStore((state) => state.viewport);
  const onNodesChange = useCanvasStore((state) => state.onNodesChange);
  const onEdgesChange = useCanvasStore((state) => state.onEdgesChange);
  const onViewportChange = useCanvasStore((state) => state.onViewportChange);
  const addNode = useCanvasStore((state) => state.addNode);
  const selectNode = useCanvasStore((state) => state.selectNode);
  const onConnect = useCanvasStore((state) => state.onConnect);

  const onDragOver = useCallback((event: DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = "copy";
  }, []);

  const onDrop = useCallback(
    (event: DragEvent<HTMLDivElement>) => {
      event.preventDefault();

      const droppedType = event.dataTransfer.getData(CANVAS_DROP_DATA_KEY) as ServiceType;
      if (!droppedType) {
        return;
      }

      const paletteItem = PALETTE_ITEMS.find((item) => item.id === droppedType);
      if (!paletteItem) {
        return;
      }

      const position = screenToFlowPosition({
        x: event.clientX,
        y: event.clientY,
      });

      addNode({
        position,
        data: {
          label: paletteItem.label,
          type: paletteItem.id,
          description: paletteItem.description,
        },
      });
    },
    [addNode, screenToFlowPosition]
  );

  return (
    <div className="h-full w-full">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        onConnect={onConnect}
        defaultEdgeOptions={{
          type: EDGE_TYPE_LABELED,
          data: { label: EDGE_DEFAULT_LABEL },
        }}
        onDrop={onDrop}
        onDragOver={onDragOver}
        onNodeClick={(_, node) => selectNode(node.id)}
        onPaneClick={() => selectNode(null)}
        viewport={viewport}
        defaultViewport={DEFAULT_CANVAS_VIEWPORT}
        onViewportChange={onViewportChange}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        minZoom={CANVAS_MIN_ZOOM}
        maxZoom={CANVAS_MAX_ZOOM}
      >
        <Background
          variant={BackgroundVariant.Dots}
          gap={CANVAS_GRID_SPACING}
          size={CANVAS_GRID_STROKE_WIDTH}
          color="#94a3b8"
        />
        <MiniMap pannable zoomable />
        <Controls showInteractive />
      </ReactFlow>
    </div>
  );
}
