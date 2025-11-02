import { test, expect } from './fixtures';

test.describe('Dashboard Page', () => {
  test('should display dashboard with stats', async ({ authenticatedPage: page }) => {
    // Check main heading
    await expect(page.locator('h1')).toContainText('Dashboard');
    
    // Check for stats grid
    const statsGrid = page.locator('.grid');
    await expect(statsGrid).toBeVisible();
    
    // Check for stat cards (should have at least 3)
    const statCards = page.locator('.bg-white.rounded-lg.shadow-md.p-6');
    const cardCount = await statCards.count();
    expect(cardCount).toBeGreaterThanOrEqual(3);
  });

  test('should display stat labels correctly', async ({ authenticatedPage: page }) => {
    // Check for expected stat labels
    await expect(page.locator('text=Total Companies')).toBeVisible();
    await expect(page.locator('text=Total Clients')).toBeVisible();
    await expect(page.locator('text=Admin Users')).toBeVisible();
  });

  test('should display quick actions section', async ({ authenticatedPage: page }) => {
    // Check for Quick Actions section
    await expect(page.locator('h2:has-text("Quick Actions")')).toBeVisible();
    
    // Check for action buttons
    await expect(page.locator('button:has-text("Manage Companies")')).toBeVisible();
    await expect(page.locator('button:has-text("View Clients")')).toBeVisible();
    await expect(page.locator('button:has-text("Manage Admin Users")')).toBeVisible();
  });

  test('should navigate to companies from quick action', async ({ authenticatedPage: page }) => {
    // Click on Manage Companies button
    await page.click('button:has-text("Manage Companies")');
    
    // Should navigate to companies page
    await expect(page).toHaveURL(/.*\/companies/);
    await expect(page.locator('h1')).toContainText('Companies');
  });

  test('should navigate to clients from quick action', async ({ authenticatedPage: page }) => {
    // Click on View Clients button
    await page.click('button:has-text("View Clients")');
    
    // Should navigate to clients page
    await expect(page).toHaveURL(/.*\/clients/);
    await expect(page.locator('h1')).toContainText('Clients');
  });

  test('should navigate to admin users from quick action', async ({ authenticatedPage: page }) => {
    // Click on Manage Admin Users button
    await page.click('button:has-text("Manage Admin Users")');
    
    // Should navigate to users page
    await expect(page).toHaveURL(/.*\/users/);
    await expect(page.locator('h1')).toContainText('Admin Users');
  });

  test('should display system information section', async ({ authenticatedPage: page }) => {
    // Check for System Information section
    await expect(page.locator('h2:has-text("System Information")')).toBeVisible();
    
    // Should show system status
    const systemStatus = page.locator('text=/Status:|Version:|Last Backup:/');
    const statusCount = await systemStatus.count();
    expect(statusCount).toBeGreaterThan(0);
  });

  test('should display clickable stat cards', async ({ authenticatedPage: page }) => {
    // Stat cards should have cursor-pointer class
    const companyStat = page.locator('text=Total Companies').locator('..');
    
    if (await companyStat.isVisible()) {
      await expect(companyStat).toHaveClass(/cursor-pointer/);
    }
  });

  test('should show loading state initially', async ({ page }) => {
    await page.goto('http://localhost:3000/admin');
    
    // Should handle authentication
    await page.fill('input[type="email"]', 'admin@mynute.com');
    await page.fill('input[type="password"]', 'admin123');
    await page.click('button[type="submit"]');
    
    await page.waitForURL('/admin/', { timeout: 5000 });
    
    // Dashboard should load
    await expect(page.locator('h1')).toContainText('Dashboard');
  });

  test('should display system status as Active', async ({ authenticatedPage: page }) => {
    // Should show Active status
    await expect(page.locator('text=System Status')).toBeVisible();
    await expect(page.locator('text=Active')).toBeVisible();
  });

  test('should display database status in system info', async ({ authenticatedPage: page }) => {
    await expect(page.locator('text=Database Status')).toBeVisible();
    await expect(page.locator('text=Connected')).toBeVisible();
  });

  test('should display environment information', async ({ authenticatedPage: page }) => {
    await expect(page.locator('text=Environment')).toBeVisible();
    await expect(page.locator('text=Production')).toBeVisible();
  });

  test('should show uptime percentage', async ({ authenticatedPage: page }) => {
    await expect(page.locator('text=Uptime')).toBeVisible();
    await expect(page.locator('text=/\\d+\\.\\d+%/')).toBeVisible();
  });

  test('should have working sidebar navigation', async ({ authenticatedPage: page }) => {
    // Click on Admin Users in sidebar
    await page.click('text=Admin Users');
    
    // Should navigate to users page
    await page.waitForURL('/admin/users', { timeout: 5000 });
    
    // Verify we're on the users page
    await expect(page.locator('h1')).toContainText('Admin Users');
  });

  test('should display user info in header', async ({ authenticatedPage: page }) => {
    // Check header has user info
    const header = page.locator('header');
    await expect(header).toBeVisible();
    
    // Should show welcome message
    await expect(header.locator('text=/Welcome back/')).toBeVisible();
  });

  test('should have working logout button', async ({ authenticatedPage: page }) => {
    // Click logout button
    await page.click('button:has-text("Logout")');
    
    // Confirm logout (if confirmation dialog appears)
    page.on('dialog', (dialog) => dialog.accept());
    
    // Should redirect to login page
    await page.waitForURL('/admin', { timeout: 5000 });
    
    // Should show login form
    await expect(page.locator('h1')).toContainText('Mynute Admin');
  });
});
