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

  test('should stay within /admin scope when navigating', async ({ authenticatedPage: page }) => {
    // All navigation should stay within /admin basepath
    const routes = [
      { text: 'Dashboard', url: '/admin/' },
      { text: 'Companies', url: '/admin/companies' },
      { text: 'Clients', url: '/admin/clients' },
      { text: 'Admin Users', url: '/admin/users' },
    ];

    for (const route of routes) {
      await page.click(`a:has-text("${route.text}")`);
      await page.waitForURL(route.url, { timeout: 5000 });
      
      // Verify URL contains /admin
      const url = page.url();
      expect(url).toContain('/admin');
    }
  });

  test('should navigate correctly with router basepath', async ({ authenticatedPage: page }) => {
    // Test that clicking dashboard stat cards navigates correctly
    await expect(page.locator('h1')).toContainText('Dashboard');
    
    // Click on "Total Companies" stat if visible
    const companiesStat = page.locator('text=Total Companies').locator('..');
    if (await companiesStat.isVisible()) {
      await companiesStat.click();
      await page.waitForTimeout(500);
      
      // Should navigate to /admin/companies, not /admin/admin/companies
      const url = page.url();
      expect(url).toContain('/admin/companies');
      expect(url).not.toContain('/admin/admin/');
    }
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

  test('should persist login at /admin after F5 reload', async ({ authenticatedPage: page }) => {
    // Ensure we're on dashboard
    await expect(page.locator('h1')).toContainText('Dashboard');
    
    // Press F5 to reload
    await page.reload();
    
    // Should NOT show login form
    await expect(page.locator('input[type="email"]')).not.toBeVisible();
    
    // Should still show dashboard
    await expect(page.locator('h1')).toContainText('Dashboard');
    
    // Should stay on /admin URL
    await expect(page).toHaveURL(/\/admin\/?$/);
  });

  test('should persist user data and token in localStorage across reloads', async ({ authenticatedPage: page }) => {
    // Check localStorage has auth data
    const token = await page.evaluate(() => localStorage.getItem('admin_token'));
    const userData = await page.evaluate(() => localStorage.getItem('admin_user'));
    
    expect(token).toBeTruthy();
    expect(userData).toBeTruthy();
    
    // Reload page
    await page.reload();
    
    // LocalStorage should still have the data
    const tokenAfter = await page.evaluate(() => localStorage.getItem('admin_token'));
    const userDataAfter = await page.evaluate(() => localStorage.getItem('admin_user'));
    
    expect(tokenAfter).toBe(token);
    expect(userDataAfter).toBe(userData);
    
    // Should still be authenticated
    await expect(page.locator('h1')).toContainText('Dashboard');
  });

  test('should not make API calls on every page reload when already authenticated', async ({ authenticatedPage: page }) => {
    // Set up request interception to count API calls
    const apiCalls: string[] = [];
    
    page.on('request', (request) => {
      if (request.url().includes('/api/')) {
        apiCalls.push(request.url());
      }
    });
    
    // Reload the page
    await page.reload();
    
    // Wait for page to be ready
    await expect(page.locator('h1')).toContainText('Dashboard');
    
    // Should NOT call /api/admin/email/* on reload (data is in localStorage)
    const emailAPICalls = apiCalls.filter(url => url.includes('/api/admin/email/'));
    expect(emailAPICalls.length).toBe(0);
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

  test('should not show 403 errors on page load when not logged in', async ({ page }) => {
    // Clear storage to simulate not logged in
    await page.goto('/admin');
    await page.evaluate(() => localStorage.clear());
    
    // Reload to trigger checkAuth
    await page.reload();
    
    // Wait a bit for any API calls
    await page.waitForTimeout(1000);
    
    // Should show login form, not error
    await expect(page.locator('h1')).toContainText('Mynute Admin');
    await expect(page.locator('input[type="email"]')).toBeVisible();
    
    // Should NOT show any 403 or error messages
    const errorText = page.locator('text=/403|forbidden|unauthorized/i');
    await expect(errorText).not.toBeVisible();
  });
});
