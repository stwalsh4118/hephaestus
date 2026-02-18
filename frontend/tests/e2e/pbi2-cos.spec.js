import { expect, test } from "@playwright/test";

const parseMatrix = (matrix) => {
  const values = matrix
    .replace("matrix(", "")
    .replace(")", "")
    .split(",")
    .map((value) => Number.parseFloat(value.trim()));

  return {
    scale: values[0] ?? 1,
    translateX: values[4] ?? 0,
    translateY: values[5] ?? 0,
  };
};

test("PBI-2 acceptance criteria pass end-to-end", async ({ page }) => {
  await page.goto("http://127.0.0.1:3100");

  await expect(page.locator(".react-flow")).toBeVisible();
  await expect(page.locator(".react-flow__background")).toBeVisible();
  await expect(page.locator(".react-flow__controls")).toBeVisible();
  await expect(page.locator(".react-flow__minimap")).toBeVisible();

  const paletteItems = page.locator('aside button[draggable="true"]');
  await expect(paletteItems).toHaveCount(5);
  await expect(paletteItems.nth(0)).toContainText("API Service");
  await expect(paletteItems.nth(1)).toContainText("PostgreSQL");

  const canvasPane = page.locator(".react-flow__pane").first();

  const dropPaletteItem = async (itemIndex, x, y) => {
    const source = paletteItems.nth(itemIndex);
    const dataTransfer = await page.evaluateHandle(() => new DataTransfer());

    await source.dispatchEvent("dragstart", { dataTransfer });
    await canvasPane.dispatchEvent("dragover", { dataTransfer });
    await canvasPane.dispatchEvent("drop", {
      dataTransfer,
      clientX: x,
      clientY: y,
    });
  };

  await dropPaletteItem(0, 560, 220);
  await dropPaletteItem(1, 720, 320);
  await expect(page.locator(".react-flow__node")).toHaveCount(2);

  const firstNode = page.locator(".react-flow__node").first();
  const initialNodeTransform = await firstNode.evaluate((node) => node.style.transform);
  const firstNodeBox = await firstNode.boundingBox();
  if (!firstNodeBox) {
    throw new Error("Node bounding box unavailable");
  }

  await page.mouse.move(
    firstNodeBox.x + firstNodeBox.width / 2,
    firstNodeBox.y + firstNodeBox.height / 2
  );
  await page.mouse.down();
  await page.mouse.move(
    firstNodeBox.x + firstNodeBox.width / 2 + 120,
    firstNodeBox.y + firstNodeBox.height / 2 + 80,
    { steps: 10 }
  );
  await page.mouse.up();

  const movedNodeTransform = await firstNode.evaluate((node) => node.style.transform);
  expect(movedNodeTransform).not.toBe(initialNodeTransform);

  const nodeCountBeforeResize = await page.locator(".react-flow__node").count();
  await page.setViewportSize({ width: 1220, height: 760 });
  await expect(page.locator(".react-flow__node")).toHaveCount(nodeCountBeforeResize);

  const nodeTransformAfterResize = await firstNode.evaluate((node) => node.style.transform);
  expect(nodeTransformAfterResize).toBe(movedNodeTransform);

  await page.mouse.move(900, 420);
  const viewport = page.locator(".react-flow__viewport");

  const initialViewportTransform = parseMatrix(
    await viewport.evaluate((node) => getComputedStyle(node).transform)
  );

  await page.mouse.wheel(0, -600);
  await page.waitForTimeout(200);

  const zoomedViewportTransform = parseMatrix(
    await viewport.evaluate((node) => getComputedStyle(node).transform)
  );
  expect(zoomedViewportTransform.scale).toBeGreaterThan(initialViewportTransform.scale);

  await page.mouse.move(940, 420);
  await page.mouse.down();
  await page.mouse.move(860, 360, { steps: 12 });
  await page.mouse.up();

  const pannedViewportTransform = parseMatrix(
    await viewport.evaluate((node) => getComputedStyle(node).transform)
  );
  expect(pannedViewportTransform.translateX).not.toBe(zoomedViewportTransform.translateX);
  expect(pannedViewportTransform.translateY).not.toBe(zoomedViewportTransform.translateY);

  await page.screenshot({ path: "test-results/pbi2-cos-final.png", fullPage: true });
});
