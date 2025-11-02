import { test, expect } from './fixtures';

test.describe('Navigation', () => {
  test('should have functional sidebar navigation', async ({ authenticatedPage: page }) => {
    // Start at dashboard
    await expect(page.locator('h1')).toContainText('Dashboard');
    
    // Navigate to Users
    await page.click('a:has-text("Admin Users")');
    await page.waitForURL('/admin/users', { timeout: 5000 });
    await expect(page.locator('h1')).toContainText('Admin Users');
    
    // Navigate to Companies
    await page.click('a:has-text("Companies")');
    await page.waitForURL('/admin/companies', { timeout: 5000 });
    await expect(page.locator('h1')).toContainText('Companies');
    
    // Navigate to Clients
    await page.click('a:has-text("Clients")');
    await page.waitForURL('/admin/clients', { timeout: 5000 });
    await expect(page.locator('h1')).toContainText('Clients');
    
    // Navigate back to Dashboard
    await page.click('a:has-text("Dashboard")');
    await page.waitForURL('/admin/', { timeout: 5000 });
    await expect(page.locator('h1')).toContainText('Dashboard');
  });

  test('should highlight active navigation item', async ({ authenticatedPage: page }) => {
    // Dashboard should be active
    const dashboardLink = page.locator('a:has-text("Dashboard")');
    await expect(dashboardLink).toHaveClass(/bg-gray-800/);
    
    // Navigate to users
    await page.click('a:has-text("Admin Users")');
    await page.waitForURL('/admin/users', { timeout: 5000 });
    
    // Users link should now be active
    const usersLink = page.locator('a:has-text("Admin Users")');
    await expect(usersLink).toHaveClass(/bg-gray-800/);
  });

  test('should navigate using sidebar links', async ({ authenticatedPage: page }) => {
    // Test each navigation link
    const routes = [
      { text: 'Dashboard', url: '/admin/', heading: 'Dashboard' },
      { text: 'Companies', url: '/admin/companies', heading: 'Companies' },
      { text: 'Clients', url: '/admin/clients', heading: 'Clients' },
      { text: 'Admin Users', url: '/admin/users', heading: 'Admin Users' },
    ];

    for (const route of routes) {
      await page.click(`a:has-text("${route.text}")`);
      await page.waitForURL(route.url, { timeout: 5000 });
      await expect(page.locator('h1')).toContainText(route.heading);
    }
  });

  test('should maintain navigation state across page reloads', async ({ authenticatedPage: page }) => {
    // Navigate to companies
    await page.click('a:has-text("Companies")');
    await page.waitForURL('/admin/companies', { timeout: 5000 });
    
    // Reload page
    await page.reload();
    
    // Should still be on companies page
    await expect(page).toHaveURL(/.*\/companies/);
    await expect(page.locator('h1')).toContainText('Companies');
  });
});

test.describe('Responsive Layout', () => {
  test('should display sidebar on desktop', async ({ authenticatedPage: page }) => {
    await page.setViewportSize({ width: 1280, height: 720 });
    
    const sidebar = page.locator('.bg-gray-900');
    await expect(sidebar).toBeVisible();
  });

  test('should display header on all viewports', async ({ authenticatedPage: page }) => {
    // Desktop
    await page.setViewportSize({ width: 1280, height: 720 });
    await expect(page.locator('header')).toBeVisible();
    
    // Tablet
    await page.setViewportSize({ width: 768, height: 1024 });
    await expect(page.locator('header')).toBeVisible();
    
    // Mobile
    await page.setViewportSize({ width: 375, height: 667 });
    await expect(page.locator('header')).toBeVisible();
  });
});

test.describe('Authentication State', () => {
  test('should persist login across page reloads', async ({ authenticatedPage: page }) => {
    // Reload page
    await page.reload();
    
    // Should still be logged in (on dashboard)
    await expect(page.locator('h1')).toContainText('Dashboard');
  });

  test('should redirect to login when not authenticated', async ({ page }) => {
    // Clear local storage (logout)
    await page.goto('/admin');
    await page.evaluate(() => localStorage.clear());
    
    // Try to navigate to dashboard
    await page.goto('/admin/');
    
    // Should show login page
    await expect(page.locator('h1')).toContainText('Mynute Admin');
    await expect(page.locator('input[type="email"]')).toBeVisible();
  });
});
