import {expect, Page, test, TestInfo} from "@playwright/test";

const screenshot = async(key: string, page: Page, testInfo: TestInfo, browserName: string) => {
  const ss = await page.screenshot({fullPage: true});
  await testInfo.attach(`${key}/${browserName}`, { body: ss, contentType: "image/png" });
};

test.describe("pages", () => {
  test('about', async ({ page, viewport, browserName }, testInfo) => {
    await page.goto(`http://127.0.0.1:40000/about`)
    await screenshot("about", page, testInfo, browserName);

    await expect(page).toHaveTitle(/Project Forge/);

    await page.locator("#search-input").focus();
    await page.waitForTimeout(1000);
    await screenshot("about/search-focus", page, testInfo, browserName);

    await page.locator("#search-input").blur();
    await page.waitForTimeout(500);

    if (viewport.width < 800) {
      await page.locator(".menu-toggle").click();
      await page.waitForTimeout(1000);
      await screenshot("about/menu-toggle", page, testInfo, browserName);

      await page.locator(".menu-toggle").click();
      await page.waitForTimeout(500);
    }
  });
});
