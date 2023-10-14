// Content managed by Project Forge, see [projectforge.md] for details.
import {defineConfig, devices} from '@playwright/test';

export default defineConfig({
  testDir: '.',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: process.env.TEST_URL || 'http://127.0.0.1:40000',
    trace: 'on-first-retry',
  },
  projects: [
    {name: 'chrome', use: {...devices['Desktop Chrome'], channel: 'chrome'}},
    {name: 'chrome.mobile', use: {...devices['Pixel 5']}},
    {name: 'chromium', use: {...devices['Desktop Chrome']}},
    {name: 'edge', use: {...devices['Desktop Edge'], channel: 'msedge'}},
    {name: 'firefox', use: {...devices['Desktop Firefox']}},
    {name: 'safari', use: {...devices['Desktop Safari']}},
    {name: 'safari.mobile', use: {...devices['iPhone 12']}},
  ],
  webServer: {
    command: "../../build/release/projectforge",
    url: "http://127.0.0.1:40000",
    reuseExistingServer: !process.env.CI,
    stdout: 'pipe',
    stderr: 'pipe',
    timeout: 5000
  },
});
