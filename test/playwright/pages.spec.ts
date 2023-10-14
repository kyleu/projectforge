// Content managed by Project Forge, see [projectforge.md] for details.
import {test, expect, Page, TestInfo} from "@playwright/test";

const pageTest = async(page: Page, testInfo: TestInfo, key: string, path: string) => {
  await page.goto(path);
  const ss = await page.screenshot({path: key+".png", fullPage: true});
  await testInfo.attach(key, { body: ss, contentType: "image/png" });
}

test.describe("pages", () => {
  test('home', async ({ page }, testInfo) => {
    await pageTest(page, testInfo, "home", "/");
    await expect(page).toHaveTitle(/Project Forge/);
  });

  test('about', async ({ page }, testInfo) => {
    await pageTest(page, testInfo, "about", "/about");
    await expect(page).toHaveTitle(/Project Forge/);
  });
});
