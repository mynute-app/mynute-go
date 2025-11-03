import { defineConfig, devices } from '@playwright/test';
import * as dotenv from 'dotenv';
import * as path from 'path';

// Load environment variables from parent directory's .env file
dotenv.config({ path: path.resolve(__dirname, '../.env') });

// Helper to check if running in CI
const isCI = process.env.CI === 'true';

// CRITICAL: Tests must only run in test environment
const APP_ENV = process.env.APP_ENV || 'dev';

if (APP_ENV !== 'test') {
  console.error('\n❌ ERROR: Frontend tests can only run when APP_ENV=test');
  console.error(`Current APP_ENV: ${APP_ENV}`);
  console.error('\nTo run tests, set APP_ENV=test in your environment or .env file\n');
  process.exit(1);
}

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: './tests', // Test files are numbered (01-, 02-, etc.) to enforce execution order
  /* Run tests in files in parallel */
  fullyParallel: false, // Changed to false to ensure proper test order
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: isCI,
  /* Retry on CI only */
  retries: isCI ? 2 : 0,
  /* Opt out of parallel tests on CI. */
  workers: 1, // Force single worker to ensure test order
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: 'html',
  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: 'http://localhost:4000/admin',
    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },

  /* Configure projects for major browsers */
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },

    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },

    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },

    /* Test against mobile viewports. */
    // {
    //   name: 'Mobile Chrome',
    //   use: { ...devices['Pixel 5'] },
    // },
    // {
    //   name: 'Mobile Safari',
    //   use: { ...devices['iPhone 12'] },
    // },
  ],

  /* Run your local dev server before starting the tests */
  webServer: {
    command: 'cd .. && go run main.go',
    url: 'http://localhost:4000',
    reuseExistingServer: !isCI,
    timeout: 120 * 1000,
  },
});
