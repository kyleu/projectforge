import {test, expect, Page, TestInfo} from "@playwright/test";

const pageTest = async(page: Page, testInfo: TestInfo, browserName: string, key: string, path: string) => {
  await page.goto(path);
  const fn = `${key}/${browserName}`;
  const ss = await page.screenshot({fullPage: true});
  await testInfo.attach(fn, { body: ss, contentType: "image/png" });
}

test.describe("pages", () => {
  test('home', async ({ page, browserName }, testInfo) => {
    await pageTest(page, testInfo, browserName, "home", "/");
    await expect(page).toHaveTitle(/{{{ .Title }}}/);
  });

  test('about', async ({ page, browserName }, testInfo) => {
    await pageTest(page, testInfo, browserName, "about", "/about");
    await expect(page).toHaveTitle(/{{{ .Title }}}/);
  });
});