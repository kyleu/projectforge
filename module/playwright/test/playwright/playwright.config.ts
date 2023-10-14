import {defineConfig, devices} from "@playwright/test";

export default defineConfig({
  testDir: ".",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: "html",
  use: {
    baseURL: process.env.TEST_URL || "http://127.0.0.1:{{{ .Port }}}",
    trace: "on",
    video: {
      mode: "on",
    },
  },
  projects: [
    {name: "chrome", use: {...devices["Desktop Chrome"], channel: "chrome"}},
    {name: "chrome.nomotion", use: {...devices["Desktop Chrome"], channel: "chrome", contextOptions: {reducedMotion: "reduce"}}},
    {name: "chrome.dark", use: {...devices["Desktop Chrome"], channel: "chrome", colorScheme: "dark"}},
    {name: "chrome.dark.nomotion", use: {...devices["Desktop Chrome"], channel: "chrome", colorScheme: "dark", contextOptions: {reducedMotion: "reduce"}}},
    {name: "chrome.mobile", use: {...devices["Pixel 5"]}},
    {name: "edge", use: {...devices["Desktop Edge"], channel: "msedge"}},
    {name: "firefox", use: {...devices["Desktop Firefox"]}},
    {name: "safari", use: {...devices["Desktop Safari"]}},
    {name: "safari.mobile", use: {...devices["iPhone 12"]}},
  ],
  webServer: {
    command: "../../build/release/{{{ .Exec }}}",
    url: "http://127.0.0.1:{{{ .Port }}}",
    reuseExistingServer: !process.env.CI,
    stdout: "pipe",
    stderr: "pipe",
    timeout: 5000
  },
});
