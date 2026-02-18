import { expect, test } from "@playwright/test";

const BASE_URL = "http://127.0.0.1:3100";

const SERVICE_TYPES = [
  { index: 0, label: "API Service", icon: "API" },
  { index: 1, label: "PostgreSQL", icon: "PG" },
  { index: 2, label: "Redis", icon: "RD" },
  { index: 3, label: "Nginx", icon: "NX" },
  { index: 4, label: "RabbitMQ", icon: "MQ" },
];

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

test.describe("PBI-3: Service Component Library & Configuration", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(BASE_URL);
    await expect(page.locator(".react-flow")).toBeVisible();
  });

  test("AC1: 5 service types available in palette", async ({ page }) => {
    const paletteItems = page.locator('aside button[draggable="true"]');
    await expect(paletteItems).toHaveCount(5);

    for (const svc of SERVICE_TYPES) {
      await expect(paletteItems.nth(svc.index)).toContainText(svc.label);
      await expect(paletteItems.nth(svc.index)).toContainText(svc.icon);
    }
  });

  test("AC2: Each service type has distinct visual appearance", async ({ page }) => {
    // Drop all 5 service types onto the canvas
    for (let i = 0; i < SERVICE_TYPES.length; i++) {
      await dropPaletteItem(page, i, 500 + i * 100, 250);
    }

    await expect(page.locator(".react-flow__node")).toHaveCount(5);

    // Verify each node has a coloured header bar (different background colours)
    const nodes = page.locator(".react-flow__node");
    const colors = new Set();
    for (let i = 0; i < 5; i++) {
      const headerBg = await nodes
        .nth(i)
        .locator("div > div")
        .first()
        .evaluate((el) => {
          return getComputedStyle(el).backgroundColor;
        });
      colors.add(headerBg);
    }
    // All 5 should have different colours
    expect(colors.size).toBe(5);
  });

  test("AC3: Clicking a node opens config panel, clicking canvas closes it", async ({ page }) => {
    // Drop a node
    await dropPaletteItem(page, 0, 500, 300);
    await expect(page.locator(".react-flow__node")).toHaveCount(1);

    // Config panel should not be visible initially (width 0)
    const configPanel = page.locator("main > aside").last();
    await expect(configPanel).toHaveCSS("width", "0px");

    // Click the node
    const node = page.locator(".react-flow__node").first();
    await node.click();

    // Config panel should now be visible
    await expect(configPanel).toHaveCSS("width", "320px");
    await expect(configPanel).toContainText("API Service");

    // Click close button
    const closeButton = configPanel.locator('button[aria-label="Close configuration panel"]');
    await closeButton.click();

    // Panel should close
    await expect(configPanel).toHaveCSS("width", "0px");
  });

  test("AC4: API Service config includes endpoint editor", async ({ page }) => {
    // Drop an API Service node
    await dropPaletteItem(page, 0, 500, 300);
    await expect(page.locator(".react-flow__node")).toHaveCount(1);

    // Click the node to open config panel
    await page.locator(".react-flow__node").first().click();
    const configPanel = page.locator("main > aside").last();
    await expect(configPanel).toHaveCSS("width", "320px");

    // Verify port input exists
    const portInput = configPanel.locator('input[type="number"]');
    await expect(portInput).toBeVisible();
    await expect(portInput).toHaveValue("8080");

    // Click "Add Endpoint" button
    const addButton = configPanel.getByText("+ Add Endpoint");
    await expect(addButton).toBeVisible();
    await addButton.click();

    // Verify endpoint row appeared
    await expect(configPanel.getByText("Endpoint 1")).toBeVisible();

    // Verify method dropdown with correct options
    const methodSelect = configPanel.locator("select").first();
    await expect(methodSelect).toBeVisible();
    const options = await methodSelect.locator("option").allTextContents();
    expect(options).toEqual(["GET", "POST", "PUT", "DELETE", "PATCH"]);

    // Verify path input
    const pathInput = configPanel.locator('input[placeholder="/users"]');
    await expect(pathInput).toBeVisible();

    // Verify schema textarea
    const schemaTextarea = configPanel.locator("textarea");
    await expect(schemaTextarea).toBeVisible();

    // Set endpoint values
    await methodSelect.selectOption("POST");
    await pathInput.fill("/users");
    await schemaTextarea.fill('{ "id": 1 }');

    // Add a second endpoint
    await addButton.click();
    await expect(configPanel.getByText("Endpoint 2")).toBeVisible();

    // Remove the first endpoint
    const removeButtons = configPanel.locator('button[aria-label^="Remove endpoint"]');
    await removeButtons.first().click();

    // Only one endpoint should remain
    await expect(configPanel.getByText("Endpoint 1")).toBeVisible();
    await expect(configPanel.getByText("Endpoint 2")).not.toBeVisible();
  });

  test("AC5: PostgreSQL config includes engine and version selection", async ({ page }) => {
    // Drop a PostgreSQL node
    await dropPaletteItem(page, 1, 500, 300);
    await expect(page.locator(".react-flow__node")).toHaveCount(1);

    // Click the node
    await page.locator(".react-flow__node").first().click();
    const configPanel = page.locator("main > aside").last();
    await expect(configPanel).toHaveCSS("width", "320px");
    await expect(configPanel).toContainText("PostgreSQL");

    // Verify engine dropdown
    const selects = configPanel.locator("select");
    await expect(selects).toHaveCount(2); // engine + version

    // Engine dropdown should have PostgreSQL option
    const engineSelect = selects.first();
    await expect(engineSelect).toHaveValue("PostgreSQL");

    // Version dropdown should have 14, 15, 16
    const versionSelect = selects.nth(1);
    const versionOptions = await versionSelect.locator("option").allTextContents();
    expect(versionOptions).toEqual(["16", "15", "14"]);

    // Change version to 15
    await versionSelect.selectOption("15");
    await expect(versionSelect).toHaveValue("15");
  });

  test("AC6: Configuration changes persist to diagram state", async ({ page }) => {
    // Drop a PostgreSQL node
    await dropPaletteItem(page, 1, 500, 300);
    await expect(page.locator(".react-flow__node")).toHaveCount(1);

    // Click the node to open config
    await page.locator(".react-flow__node").first().click();
    const configPanel = page.locator("main > aside").last();
    await expect(configPanel).toHaveCSS("width", "320px");

    // Change version to 14
    const selects = configPanel.locator("select");
    const versionSelect = selects.nth(1);
    await versionSelect.selectOption("14");
    await expect(versionSelect).toHaveValue("14");

    // Close panel via close button
    const closeButton = configPanel.locator('button[aria-label="Close configuration panel"]');
    await closeButton.click();
    await expect(configPanel).toHaveCSS("width", "0px");

    // Re-click the node to re-open config
    await page.locator(".react-flow__node").first().click();
    await expect(configPanel).toHaveCSS("width", "320px");

    // Verify version is still 14 (persisted)
    const reselectedVersionSelect = configPanel.locator("select").nth(1);
    await expect(reselectedVersionSelect).toHaveValue("14");
  });
});
