// Content managed by Project Forge, see [projectforge.md] for details.
import {defineConfig, devices, PlaywrightTestOptions, PlaywrightWorkerOptions, Project} from "@playwright/test";

const prj = (
    name: string, device: string, channel: string, width: number, height: number,
    dark: boolean = false, reduceMotion: boolean = false, jsDisabled: boolean = false,
): Project<PlaywrightTestOptions & PlaywrightWorkerOptions> => {
  return {
    name,
    use: {
      ...devices[device],
      channel,
      colorScheme: (dark ? "dark" : "light"),
      contextOptions: {reducedMotion: reduceMotion ? "reduce" : "no-preference"},
      javaScriptEnabled: !jsDisabled,
      trace: "on",
      viewport: {width, height},
      video: {mode: "on", size: {width, height}}
    },
  };
}

export default defineConfig({
  testDir: ".",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ["html"],
    ["json", { outputFile: "playwright-report/results.json"}],
    ["list"]
  ],
  use: {
    baseURL: process.env.TEST_URL || "http://127.0.0.1:40000",
  },
  projects: [
    prj("chrome", "Desktop Chrome", "chrome", 1280, 720, false, false),
    prj("chrome.nojs", "Desktop Chrome", "chrome", 1280, 720, false, false, true),
    prj("chrome.nomotion", "Desktop Chrome", "chrome", 1280, 720, false, true),
    prj("chrome.dark", "Desktop Chrome", "chrome", 1280, 720, true, false),
    prj("chrome.dark.nojs", "Desktop Chrome", "chrome", 1280, 720, true, false, true),
    prj("chrome.dark.nomotion", "Desktop Chrome", "chrome", 1280, 720, true, true),
    prj("chrome.mobile", "Pixel 5", "", 393, 727, false, false),
    prj("edge", "Desktop Edge", "msedge", 1280, 720, false, false),
    prj("firefox", "Desktop Firefox", "", 1280, 720, false, false),
    prj("safari", "Desktop Safari", "", 1280, 720, false, false),
    prj("safari.mobile", "iPhone 12", "", 390, 664, false, false),
    prj("safari.mobile.nojs", "iPhone 12", "", 390, 664, false, false, true),
    prj("safari.mobile.nomotion", "iPhone 12", "", 390, 664, false, true),
    prj("safari.mobile.dark", "iPhone 12", "", 390, 664, true, false),
    prj("safari.mobile.landscape", "iPhone 12 landscape", "", 750, 340, false, false),
    prj("safari.mobile.landscape.dark", "iPhone 12 landscape", "", 750, 340, true, false),
  ],
  webServer: {
    command: "../../build/release/projectforge",
    url: "http://127.0.0.1:40000",
    reuseExistingServer: !process.env.CI,
    stdout: "pipe",
    stderr: "pipe",
    timeout: 60000
  },
});
