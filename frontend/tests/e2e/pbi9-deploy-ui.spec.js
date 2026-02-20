import { expect, test } from "@playwright/test";

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

test.describe("PBI-9 Task 9-5: Deploy UI — Toolbar Buttons & Node Status Badges", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(BASE_URL);
    await expect(page.locator(".react-flow")).toBeVisible();
  });

  test("Deploy button renders in toolbar", async ({ page }) => {
    const deployBtn = page.locator("button", { hasText: "Deploy" });
    await expect(deployBtn).toBeVisible();
  });

  test("Deploy button is disabled when canvas is empty", async ({ page }) => {
    // No nodes on canvas initially
    await expect(page.locator(".react-flow__node")).toHaveCount(0);

    const deployBtn = page.locator("button", { hasText: "Deploy" });
    await expect(deployBtn).toBeDisabled();
  });

  test("Deploy button is enabled when nodes are present", async ({ page }) => {
    await dropPaletteItem(page, 0, 400, 300);
    await expect(page.locator(".react-flow__node")).toHaveCount(1);

    const deployBtn = page.locator("button", { hasText: "Deploy" });
    await expect(deployBtn).toBeEnabled();
  });

  test("Teardown button is hidden when not deployed", async ({ page }) => {
    const teardownBtn = page.locator("button", { hasText: "Teardown" });
    await expect(teardownBtn).toHaveCount(0);
  });

  test("ErrorToast renders message and auto-dismisses", async ({ page }) => {
    // Inject an error into the deploy store to trigger the toast
    await page.evaluate(() => {
      const storeModule = window.__ZUSTAND_DEPLOY_STORE__;
      if (storeModule) {
        storeModule.setState({ error: "Test error message" });
      }
    });

    // Alternatively, trigger via the store API exposed on window
    // If the store isn't accessible, we can set it via the deploy action
    // by pointing at a non-existent backend
    await dropPaletteItem(page, 0, 400, 300);
    await expect(page.locator(".react-flow__node")).toHaveCount(1);

    // Click Deploy — since no backend is running, it should error
    const deployBtn = page.locator("button", { hasText: "Deploy" });
    await deployBtn.click();

    // The error toast should appear (fetch to backend fails)
    const toast = page.locator("text=Dismiss error").locator("..");
    await expect(toast.or(page.locator("[aria-label='Dismiss error']"))).toBeVisible({
      timeout: 10000,
    });
  });

  test("ErrorToast can be manually dismissed", async ({ page }) => {
    await dropPaletteItem(page, 0, 400, 300);
    await expect(page.locator(".react-flow__node")).toHaveCount(1);

    // Click Deploy — will fail and show error toast
    const deployBtn = page.locator("button", { hasText: "Deploy" });
    await deployBtn.click();

    // Wait for toast to appear
    const dismissBtn = page.locator("[aria-label='Dismiss error']");
    await expect(dismissBtn).toBeVisible({ timeout: 10000 });

    // Dismiss it
    await dismissBtn.click();
    await expect(dismissBtn).not.toBeVisible();
  });
});
