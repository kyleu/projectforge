import {test, expect, Page, TestInfo} from "@playwright/test";

const pageTest = async(page: Page, testInfo: TestInfo, browserName: string, key: string, path: string) => {
  await page.goto(path);
  const fn = `${key}/${browserName}`;
  const ss = await page.screenshot({path: `playwright-assets/screenshots/${fn}.png`, fullPage: true});
  await testInfo.attach(fn, { body: ss, contentType: "image/png" });
}

test.describe("pages", () => {
  test('home', async ({ page, browserName }, testInfo) => {
    await pageTest(page, testInfo, browserName, "home", "/");
    await expect(page).toHaveTitle(/Project Forge/);
  });

  test('about', async ({ page, browserName }, testInfo) => {
    await pageTest(page, testInfo, browserName, "about", "/about");
    await expect(page).toHaveTitle(/Project Forge/);
  });
});
