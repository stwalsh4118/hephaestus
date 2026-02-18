import { expect, test } from "@playwright/test";
import fs from "node:fs";
import path from "node:path";

const BASE_URL = "http://127.0.0.1:3100";

const dropPaletteItem = async (page, itemIndex, x, y) => {
  const paletteItems = page.locator('aside button[draggable="true"]');
  const source = paletteItems.nth(itemIndex);
  const canvasPane = page.locator(".react-flow__pane").first();
  const dataTransfer = await page.evaluateHandle(() => new DataTransfer());

  await source.dispatchEvent("dragstart", { dataTransfer });
  await canvasPane.dispatchEvent("dragover", { dataTransfer });
  await canvasPane.dispatchEvent("drop", {
    dataTransfer,
    clientX: x,
    clientY: y,
  });
};

const connectNodes = async (page, sourceNode, targetNode) => {
  const sourceHandle = sourceNode.locator(".react-flow__handle-right");
  const targetHandle = targetNode.locator(".react-flow__handle-left");

  const sourceBox = await sourceHandle.boundingBox();
  const targetBox = await targetHandle.boundingBox();
  if (!sourceBox || !targetBox) {
    throw new Error("Handle bounding box unavailable");
  }

  await page.mouse.move(
    sourceBox.x + sourceBox.width / 2,
    sourceBox.y + sourceBox.height / 2
  );
  await page.mouse.down();
  await page.mouse.move(
    targetBox.x + targetBox.width / 2,
    targetBox.y + targetBox.height / 2,
    { steps: 10 }
  );
  await page.mouse.up();
};

const importDiagramJson = async (page, diagramJson) => {
  const tmpDir = path.join(process.cwd(), "test-results");
  fs.mkdirSync(tmpDir, { recursive: true });
  const filePath = path.join(tmpDir, `import-${Date.now()}.json`);
  fs.writeFileSync(filePath, JSON.stringify(diagramJson));

  const fileChooserPromise = page.waitForEvent("filechooser");
  await page.locator("button", { hasText: "Import JSON" }).click();
  const fileChooser = await fileChooserPromise;
  await fileChooser.setFiles(filePath);

  // Wait for import to process before cleaning up the temp file
  await page.waitForTimeout(100);
  fs.unlinkSync(filePath);
};

test.describe("PBI-4: Connections & Topology Export", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(BASE_URL);
    await expect(page.locator(".react-flow")).toBeVisible();
  });

  test("AC1: Edges can be drawn between nodes by dragging from source to target handle", async ({
    page,
  }) => {
    await dropPaletteItem(page, 0, 400, 300);
    await dropPaletteItem(page, 1, 700, 300);
    await expect(page.locator(".react-flow__node")).toHaveCount(2);

    const nodeA = page.locator(".react-flow__node").nth(0);
    const nodeB = page.locator(".react-flow__node").nth(1);

    await connectNodes(page, nodeA, nodeB);

    await expect(page.locator(".react-flow__edge")).toHaveCount(1);
  });

  test("AC2: Edges display an arrow indicating direction", async ({
    page,
  }) => {
    await dropPaletteItem(page, 0, 400, 300);
    await dropPaletteItem(page, 1, 700, 300);
    await expect(page.locator(".react-flow__node")).toHaveCount(2);

    const nodeA = page.locator(".react-flow__node").nth(0);
    const nodeB = page.locator(".react-flow__node").nth(1);
    await connectNodes(page, nodeA, nodeB);

    await expect(page.locator(".react-flow__edge")).toHaveCount(1);

    // Verify arrowhead marker definition exists in the SVG
    const markerDef = page.locator("marker[id*='arrowclosed']");
    await expect(markerDef.first()).toBeAttached();

    // Verify the edge path references the marker
    const markerEnd = await page
      .locator(".react-flow__edge-path")
      .first()
      .getAttribute("marker-end");
    expect(markerEnd).toContain("arrowclosed");
  });

  test("AC3: Edge labels can be added and edited", async ({ page }) => {
    // Import a diagram with a labeled edge to test editing
    const diagramWithLabel = {
      id: "test-label",
      name: "Label Test",
      nodes: [
        {
          id: "n1",
          type: "api-service",
          name: "API",
          description: "",
          position: { x: 200, y: 200 },
        },
        {
          id: "n2",
          type: "postgresql",
          name: "DB",
          description: "",
          position: { x: 500, y: 200 },
        },
      ],
      edges: [
        { id: "e1", source: "n1", target: "n2", label: "reads" },
      ],
    };

    await importDiagramJson(page, diagramWithLabel);

    await expect(page.locator(".react-flow__node")).toHaveCount(2);
    await expect(page.locator(".react-flow__edge")).toHaveCount(1);

    // Verify the label text is displayed
    const edgeLabelRenderer = page.locator(".react-flow__edgelabel-renderer");
    await expect(edgeLabelRenderer).toContainText("reads");

    // Double-click the label to enter edit mode
    const labelSpan = edgeLabelRenderer.locator("span").first();
    await labelSpan.dblclick();

    // Verify input appears
    const labelInput = edgeLabelRenderer.locator("input");
    await expect(labelInput).toBeVisible();

    // Clear and type new label
    await labelInput.fill("reads/writes");
    await labelInput.press("Enter");

    // Verify the updated label text is displayed
    await expect(edgeLabelRenderer).toContainText("reads/writes");

    // Verify via export that the label persisted
    const downloadPromise = page.waitForEvent("download");
    await page.locator("button", { hasText: "Export JSON" }).click();
    const download = await downloadPromise;
    const downloadPath = await download.path();
    const exported = JSON.parse(fs.readFileSync(downloadPath, "utf-8"));
    expect(exported.edges[0].label).toBe("reads/writes");
  });

  test("AC4: Duplicate edges between same source-target pair are prevented", async ({
    page,
  }) => {
    await dropPaletteItem(page, 0, 400, 300);
    await dropPaletteItem(page, 1, 700, 300);
    await expect(page.locator(".react-flow__node")).toHaveCount(2);

    const nodeA = page.locator(".react-flow__node").nth(0);
    const nodeB = page.locator(".react-flow__node").nth(1);

    // Draw first edge
    await connectNodes(page, nodeA, nodeB);
    await expect(page.locator(".react-flow__edge")).toHaveCount(1);

    // Attempt to draw duplicate edge
    await connectNodes(page, nodeA, nodeB);

    // Should still be only 1 edge
    await expect(page.locator(".react-flow__edge")).toHaveCount(1);
  });

  test("AC5: Export produces valid JSON matching PRD schema", async ({
    page,
  }) => {
    await dropPaletteItem(page, 0, 400, 300); // API Service
    await dropPaletteItem(page, 1, 700, 300); // PostgreSQL
    await expect(page.locator(".react-flow__node")).toHaveCount(2);

    const nodeA = page.locator(".react-flow__node").nth(0);
    const nodeB = page.locator(".react-flow__node").nth(1);
    await connectNodes(page, nodeA, nodeB);
    await expect(page.locator(".react-flow__edge")).toHaveCount(1);

    // Click Export and intercept the download
    const downloadPromise = page.waitForEvent("download");
    await page.locator("button", { hasText: "Export JSON" }).click();
    const download = await downloadPromise;

    // Read the downloaded file
    const downloadPath = await download.path();
    const jsonContent = fs.readFileSync(downloadPath, "utf-8");
    const diagram = JSON.parse(jsonContent);

    // Verify top-level schema
    expect(diagram).toHaveProperty("id");
    expect(diagram).toHaveProperty("name");
    expect(diagram).toHaveProperty("nodes");
    expect(diagram).toHaveProperty("edges");
    expect(typeof diagram.id).toBe("string");
    expect(typeof diagram.name).toBe("string");
    expect(Array.isArray(diagram.nodes)).toBe(true);
    expect(Array.isArray(diagram.edges)).toBe(true);

    // Verify nodes
    expect(diagram.nodes).toHaveLength(2);
    for (const node of diagram.nodes) {
      expect(node).toHaveProperty("id");
      expect(node).toHaveProperty("type");
      expect(node).toHaveProperty("name");
      expect(node).toHaveProperty("position");
      expect(typeof node.position.x).toBe("number");
      expect(typeof node.position.y).toBe("number");
    }

    // Verify node types are service types, not React Flow types
    const types = diagram.nodes.map((n) => n.type);
    expect(types).toContain("api-service");
    expect(types).toContain("postgresql");

    // Verify edges
    expect(diagram.edges).toHaveLength(1);
    const edge = diagram.edges[0];
    expect(edge).toHaveProperty("id");
    expect(edge).toHaveProperty("source");
    expect(edge).toHaveProperty("target");
    expect(edge).toHaveProperty("label");
    expect(typeof edge.label).toBe("string");
  });

  test("AC6: Importing exported JSON restores diagram accurately", async ({
    page,
  }) => {
    await dropPaletteItem(page, 0, 400, 300); // API Service
    await dropPaletteItem(page, 1, 700, 300); // PostgreSQL
    await expect(page.locator(".react-flow__node")).toHaveCount(2);

    const nodeA = page.locator(".react-flow__node").nth(0);
    const nodeB = page.locator(".react-flow__node").nth(1);
    await connectNodes(page, nodeA, nodeB);
    await expect(page.locator(".react-flow__edge")).toHaveCount(1);

    // Export the diagram
    const downloadPromise = page.waitForEvent("download");
    await page.locator("button", { hasText: "Export JSON" }).click();
    const download = await downloadPromise;
    const downloadPath = await download.path();
    const originalJson = fs.readFileSync(downloadPath, "utf-8");
    const originalDiagram = JSON.parse(originalJson);

    // Reload to clear state
    await page.reload();
    await expect(page.locator(".react-flow")).toBeVisible();
    await expect(page.locator(".react-flow__node")).toHaveCount(0);
    await expect(page.locator(".react-flow__edge")).toHaveCount(0);

    // Import the previously exported diagram
    await importDiagramJson(page, originalDiagram);

    // Wait for nodes and edges to appear
    await expect(page.locator(".react-flow__node")).toHaveCount(2);
    await expect(page.locator(".react-flow__edge")).toHaveCount(1);

    // Re-export and compare
    const reExportPromise = page.waitForEvent("download");
    await page.locator("button", { hasText: "Export JSON" }).click();
    const reExportDownload = await reExportPromise;
    const reExportPath = await reExportDownload.path();
    const reExportedDiagram = JSON.parse(
      fs.readFileSync(reExportPath, "utf-8")
    );

    // Compare nodes (ignoring top-level diagram id which is regenerated)
    expect(reExportedDiagram.nodes).toHaveLength(originalDiagram.nodes.length);
    for (const orig of originalDiagram.nodes) {
      const reimported = reExportedDiagram.nodes.find((n) => n.id === orig.id);
      expect(reimported).toBeDefined();
      expect(reimported.type).toBe(orig.type);
      expect(reimported.name).toBe(orig.name);
      expect(reimported.position.x).toBeCloseTo(orig.position.x, 0);
      expect(reimported.position.y).toBeCloseTo(orig.position.y, 0);
    }

    // Compare edges
    expect(reExportedDiagram.edges).toHaveLength(originalDiagram.edges.length);
    for (const orig of originalDiagram.edges) {
      const reimported = reExportedDiagram.edges.find((e) => e.id === orig.id);
      expect(reimported).toBeDefined();
      expect(reimported.source).toBe(orig.source);
      expect(reimported.target).toBe(orig.target);
      expect(reimported.label).toBe(orig.label);
    }
  });
});
