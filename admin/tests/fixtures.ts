import { test as base } from '@playwright/test';
import type { Page } from '@playwright/test';

/**
 * Test fixtures for admin panel E2E tests
 */

export interface AdminUser {
  email: string;
  password: string;
  name: string;
}

export const adminCredentials: AdminUser = {
  email: 'admin@mynute.com',
  password: 'Admin@123456',
  name: 'Admin User',
};

// Define fixture types
type AuthFixtures = {
  authenticatedPage: Page;
};

// Extend base test with custom fixtures
export const test = base.extend<AuthFixtures>({
  // Auto-login fixture
  authenticatedPage: async ({ page }, use) => {
    // Navigate to login page (baseURL is set to http://localhost:4000/admin in playwright.config.ts)
    await page.goto('/');
    
    // Fill login form
    await page.fill('input[type="email"]', adminCredentials.email);
    await page.fill('input[type="password"]', adminCredentials.password);
    
    // Submit form
    await page.click('button[type="submit"]');
    
    // Wait for navigation to dashboard
    await page.waitForURL('/', { timeout: 5000 });
    
    // Use the authenticated page
    await use(page);
  },
});

export { expect } from '@playwright/test';
