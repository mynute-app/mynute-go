import { test, expect } from './fixtures';

/**
 * Tests for routing with basepath="/admin"
 * 
 * These tests cover critical bugs that were fixed:
 * 1. Login should redirect to /admin, not /
 * 2. Navigation should work within /admin scope
 * 3. F5 reload should stay on /admin, not redirect to /
 * 4. All route() calls should use relative paths with basepath
 */
test.describe('Router Basepath Functionality', () => {
  test('should use /admin basepath for all routes', async ({ authenticatedPage: page }) => {
    // Dashboard
    await expect(page).toHaveURL(/\/admin\/?$/);
    
    // Navigate to users
    await page.click('a:has-text("Admin Users")');
    await page.waitForURL('/admin/users', { timeout: 5000 });
    expect(page.url()).toContain('/admin/users');
    
    // Navigate to companies
    await page.click('a:has-text("Companies")');
    await page.waitForURL('/admin/companies', { timeout: 5000 });
    expect(page.url()).toContain('/admin/companies');
    
    // Navigate to clients
    await page.click('a:has-text("Clients")');
    await page.waitForURL('/admin/clients', { timeout: 5000 });
    expect(page.url()).toContain('/admin/clients');
  });

  test('should not double /admin prefix in URLs', async ({ authenticatedPage: page }) => {
    // Click through all navigation
    const links = ['Companies', 'Clients', 'Admin Users', 'Dashboard'];
    
    for (const linkText of links) {
      await page.click(`a:has-text("${linkText}")`);
      await page.waitForTimeout(500);
      
      // URL should never have /admin/admin/
      const url = page.url();
      expect(url).not.toContain('/admin/admin/');
    }
  });

  test('should stay on /admin when reloading dashboard', async ({ authenticatedPage: page }) => {
    // Ensure we're on dashboard
    await expect(page).toHaveURL(/\/admin\/?$/);
    await expect(page.locator('h1')).toContainText('Dashboard');
    
    // Reload
    await page.reload();
    
    // Should still be on /admin dashboard
    await expect(page).toHaveURL(/\/admin\/?$/);
    await expect(page.locator('h1')).toContainText('Dashboard');
    
    // Should NOT be on root /
    const url = page.url();
    expect(url).not.toMatch(/^https?:\/\/[^\/]+\/$/);
  });

  test('should navigate correctly from dashboard stat cards', async ({ authenticatedPage: page }) => {
    await expect(page.locator('h1')).toContainText('Dashboard');
    
    // Click on "Total Companies" stat
    const companiesStat = page.locator('text=Total Companies');
    if (await companiesStat.isVisible()) {
      const statCard = companiesStat.locator('..');
      await statCard.click();
      await page.waitForTimeout(500);
      
      // Should go to /admin/companies
      expect(page.url()).toContain('/admin/companies');
      expect(page.url()).not.toContain('/admin/admin/companies');
      await expect(page.locator('h1')).toContainText('Companies');
    }
  });

  test('should navigate correctly from dashboard quick actions', async ({ authenticatedPage: page }) => {
    // Test "View All Companies" button
    const companiesButton = page.locator('button:has-text("View All Companies")');
    if (await companiesButton.isVisible()) {
      await companiesButton.click();
      await page.waitForTimeout(500);
      
      expect(page.url()).toContain('/admin/companies');
      expect(page.url()).not.toContain('/admin/admin/');
      await expect(page.locator('h1')).toContainText('Companies');
    }
    
    // Go back to dashboard
    await page.click('a:has-text("Dashboard")');
    await page.waitForURL('/admin/', { timeout: 5000 });
    
    // Test "View All Clients" button
    const clientsButton = page.locator('button:has-text("View All Clients")');
    if (await clientsButton.isVisible()) {
      await clientsButton.click();
      await page.waitForTimeout(500);
      
      expect(page.url()).toContain('/admin/clients');
      expect(page.url()).not.toContain('/admin/admin/');
      await expect(page.locator('h1')).toContainText('Clients');
    }
  });

  test('should navigate to company detail with correct basepath', async ({ authenticatedPage: page }) => {
    // Navigate to companies
    await page.click('a:has-text("Companies")');
    await page.waitForURL('/admin/companies', { timeout: 5000 });
    
    // Wait for companies to load
    await page.waitForTimeout(1000);
    
    // Click "View Details" if available
    const viewButton = page.locator('button:has-text("View Details")').first();
    if (await viewButton.isVisible()) {
      await viewButton.click();
      await page.waitForTimeout(500);
      
      // Should navigate to /admin/companies/:id, not /admin/admin/companies/:id
      expect(page.url()).toMatch(/\/admin\/companies\/\d+/);
      expect(page.url()).not.toContain('/admin/admin/');
    }
  });

  test('should navigate back from company detail correctly', async ({ authenticatedPage: page }) => {
    // Navigate to companies
    await page.click('a:has-text("Companies")');
    await page.waitForURL('/admin/companies', { timeout: 5000 });
    
    await page.waitForTimeout(1000);
    
    // Click "View Details" if available
    const viewButton = page.locator('button:has-text("View Details")').first();
    if (await viewButton.isVisible()) {
      await viewButton.click();
      await page.waitForTimeout(500);
      
      // Click "Back to Companies" button
      const backButton = page.locator('button:has-text("Back to Companies")');
      if (await backButton.isVisible()) {
        await backButton.click();
        await page.waitForTimeout(500);
        
        // Should go back to /admin/companies
        expect(page.url()).toContain('/admin/companies');
        expect(page.url()).not.toMatch(/\/companies\/\d+/);
        expect(page.url()).not.toContain('/admin/admin/');
      }
    }
  });

  test('should handle direct URL access to /admin routes', async ({ authenticatedPage: page }) => {
    // Directly navigate to /admin/users
    await page.goto('/admin/users');
    await expect(page.locator('h1')).toContainText('Admin Users');
    
    // Directly navigate to /admin/companies
    await page.goto('/admin/companies');
    await expect(page.locator('h1')).toContainText('Companies');
    
    // Directly navigate to /admin/clients
    await page.goto('/admin/clients');
    await expect(page.locator('h1')).toContainText('Clients');
    
    // All should work correctly
    expect(page.url()).toContain('/admin/');
  });

  test('should reload at any /admin/* route without issues', async ({ authenticatedPage: page }) => {
    const routes = [
      { url: '/admin/', heading: 'Dashboard' },
      { url: '/admin/users', heading: 'Admin Users' },
      { url: '/admin/companies', heading: 'Companies' },
      { url: '/admin/clients', heading: 'Clients' },
    ];

    for (const route of routes) {
      // Navigate to route
      await page.goto(route.url);
      await expect(page.locator('h1')).toContainText(route.heading);
      
      // Reload
      await page.reload();
      
      // Should stay on same route
      await expect(page).toHaveURL(route.url);
      await expect(page.locator('h1')).toContainText(route.heading);
      
      // Should not redirect to login
      await expect(page.locator('input[type="email"]')).not.toBeVisible();
    }
  });

  test('should handle browser back/forward with basepath correctly', async ({ authenticatedPage: page }) => {
    // Navigate through several pages
    await expect(page.locator('h1')).toContainText('Dashboard');
    
    await page.click('a:has-text("Companies")');
    await page.waitForURL('/admin/companies', { timeout: 5000 });
    
    await page.click('a:has-text("Clients")');
    await page.waitForURL('/admin/clients', { timeout: 5000 });
    
    // Go back
    await page.goBack();
    await expect(page).toHaveURL('/admin/companies');
    await expect(page.locator('h1')).toContainText('Companies');
    
    // Go back again
    await page.goBack();
    await expect(page).toHaveURL(/\/admin\/?$/);
    await expect(page.locator('h1')).toContainText('Dashboard');
    
    // Go forward
    await page.goForward();
    await expect(page).toHaveURL('/admin/companies');
    await expect(page.locator('h1')).toContainText('Companies');
  });
});

test.describe('Sidebar Navigation with Basepath', () => {
  test('should highlight correct active item with basepath', async ({ authenticatedPage: page }) => {
    // Dashboard should be active
    const dashboardLink = page.locator('a:has-text("Dashboard")');
    await expect(dashboardLink).toHaveClass(/bg-gray-800/);
    
    // Navigate to companies
    await page.click('a:has-text("Companies")');
    await page.waitForURL('/admin/companies', { timeout: 5000 });
    
    // Companies should now be active
    const companiesLink = page.locator('a:has-text("Companies")');
    await expect(companiesLink).toHaveClass(/bg-gray-800/);
    
    // Dashboard should not be active
    await expect(dashboardLink).not.toHaveClass(/bg-gray-800/);
  });

  test('should maintain active state after reload', async ({ authenticatedPage: page }) => {
    // Navigate to companies
    await page.click('a:has-text("Companies")');
    await page.waitForURL('/admin/companies', { timeout: 5000 });
    
    // Companies should be active
    const companiesLink = page.locator('a:has-text("Companies")');
    await expect(companiesLink).toHaveClass(/bg-gray-800/);
    
    // Reload
    await page.reload();
    
    // Companies should still be active
    await expect(companiesLink).toHaveClass(/bg-gray-800/);
  });

  test('should activate dashboard for root /admin route', async ({ authenticatedPage: page }) => {
    // Go to exact /admin
    await page.goto('/admin');
    
    // Dashboard link should be active
    const dashboardLink = page.locator('a:has-text("Dashboard")');
    await expect(dashboardLink).toHaveClass(/bg-gray-800/);
  });
});
