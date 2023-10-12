import { test, expect } from '@playwright/test';

test.describe("pages", () => {
  test('home', async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveTitle("{{{ .Title }}}");
  });

  test('about', async ({ page }) => {
    await page.goto("/about");
  });
});
