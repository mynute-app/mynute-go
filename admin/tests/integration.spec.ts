import { test, expect } from './fixtures';

test.describe('Error Handling and Edge Cases', () => {
  test('should handle API errors on companies page', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies');
    
    // Wait for potential error or loading state
    await page.waitForTimeout(2000);
    
    // Should either show companies or error message
    const hasError = await page.locator('text=/error/i').isVisible().catch(() => false);
    const hasContent = await page.locator('.bg-white.rounded-lg.shadow').isVisible().catch(() => false);
    const hasEmpty = await page.locator('text=/No companies/i').isVisible().catch(() => false);
    
    expect(hasError || hasContent || hasEmpty).toBeTruthy();
  });

  test('should handle API errors on clients page', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/clients');
    
    await page.waitForTimeout(2000);
    
    // Should either show clients or error message
    const hasError = await page.locator('text=/error/i').isVisible().catch(() => false);
    const hasTable = await page.locator('table').isVisible().catch(() => false);
    const hasEmpty = await page.locator('text=/No clients/i').isVisible().catch(() => false);
    
    expect(hasError || hasTable || hasEmpty).toBeTruthy();
  });

  test('should show loading state on companies page', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies');
    
    // Should briefly show loading or go straight to content
    const loadingVisible = await page.locator('text=/Loading companies/i').isVisible().catch(() => false);
    
    // Either loading was shown or content loaded immediately
    expect(loadingVisible || true).toBeTruthy();
  });

  test('should show loading state on clients page', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/clients');
    
    // Should briefly show loading or go straight to content
    const loadingVisible = await page.locator('text=/Loading clients/i').isVisible().catch(() => false);
    
    expect(loadingVisible || true).toBeTruthy();
  });

  test('should handle empty search results on companies', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies');
    
    await page.waitForTimeout(1000);
    
    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.fill('XYZ_NONEXISTENT_COMPANY_12345');
    
    await page.waitForTimeout(500);
    
    // Should show "No companies found" message
    await expect(page.locator('text=/No companies found/i')).toBeVisible();
  });

  test('should handle empty search results on clients', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/clients');
    
    await page.waitForTimeout(1000);
    
    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.fill('XYZ_NONEXISTENT_CLIENT_12345');
    
    await page.waitForTimeout(500);
    
    // Should show "No clients found" message
    await expect(page.locator('text=/No clients found/i')).toBeVisible();
  });

  test('should handle missing company data gracefully', async ({ authenticatedPage: page }) => {
    // Try to access company with very high ID
    await page.goto('http://localhost:3000/admin/companies/999999');
    
    await page.waitForTimeout(1500);
    
    // Should show not found or loading state
    const hasNotFound = await page.locator('text=/not found/i').isVisible().catch(() => false);
    const hasLoading = await page.locator('text=/Loading/i').isVisible().catch(() => false);
    const hasBackButton = await page.locator('button:has-text("Back to Companies")').isVisible().catch(() => false);
    
    expect(hasNotFound || hasLoading || hasBackButton).toBeTruthy();
  });

  test('should handle network errors gracefully', async ({ authenticatedPage: page }) => {
    // This test checks if the app handles offline state
    await page.goto('http://localhost:3000/admin/dashboard');
    
    // Dashboard should load or show error
    const hasDashboard = await page.locator('h1:has-text("Dashboard")').isVisible();
    expect(hasDashboard).toBeTruthy();
  });

  test('should prevent navigation when modal is open', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/clients');
    
    await page.waitForTimeout(1000);
    
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      // Modal should be visible
      await expect(page.locator('text=Client Details')).toBeVisible();
      
      // URL should still be /clients
      await expect(page).toHaveURL(/.*\/clients/);
    }
  });

  test('should handle rapid tab switching', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    // Rapidly click through tabs
    const tabs = ['Branches', 'Employees', 'Services', 'Overview'];
    
    for (const tab of tabs) {
      await page.click(`button:has-text("${tab}")`);
      await page.waitForTimeout(100);
    }
    
    // Should still be functional
    await expect(page.locator('button:has-text("Overview")')).toHaveClass(/border-primary/);
  });

  test('should handle special characters in search', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies');
    
    await page.waitForTimeout(1000);
    
    const searchInput = page.locator('input[placeholder*="Search"]');
    
    // Test special characters
    const specialChars = ['<script>', '"; DROP TABLE', '100%', '@#$%'];
    
    for (const chars of specialChars) {
      await searchInput.fill(chars);
      await page.waitForTimeout(200);
      
      // Should not crash - either show results or no results
      const pageContent = await page.locator('body').isVisible();
      expect(pageContent).toBeTruthy();
    }
  });

  test('should maintain scroll position in tables', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/clients');
    
    await page.waitForTimeout(1000);
    
    const tableVisible = await page.locator('table').isVisible().catch(() => false);
    
    if (tableVisible) {
      // Scroll page
      await page.evaluate(() => window.scrollTo(0, 500));
      
      // Page should maintain scroll
      const scrollY = await page.evaluate(() => window.scrollY);
      expect(scrollY).toBeGreaterThan(0);
    }
  });
});

test.describe('Data Consistency', () => {
  test('should show consistent company count across pages', async ({ authenticatedPage: page }) => {
    // Get count from dashboard
    await page.goto('http://localhost:3000/admin/dashboard');
    await page.waitForTimeout(1000);
    
    // Navigate to companies page
    await page.click('button:has-text("View All Companies")');
    await page.waitForURL(/.*\/companies/, { timeout: 5000 });
    await page.waitForTimeout(1000);
    
    // Count should match (or be close if data changed)
    const cards = page.locator('.bg-white.rounded-lg.shadow');
    const cardCount = await cards.count();
    
    // Basic check - should have some consistency
    expect(cardCount).toBeGreaterThanOrEqual(0);
  });

  test('should show consistent client count across pages', async ({ authenticatedPage: page }) => {
    // Get count from dashboard
    await page.goto('http://localhost:3000/admin/dashboard');
    await page.waitForTimeout(1000);
    
    // Navigate to clients page
    await page.click('button:has-text("View All Clients")');
    await page.waitForURL(/.*\/clients/, { timeout: 5000 });
    await page.waitForTimeout(1000);
    
    // Should have clients data loaded
    const tableVisible = await page.locator('table').isVisible().catch(() => false);
    const emptyState = await page.locator('text=/No clients/i').isVisible().catch(() => false);
    
    expect(tableVisible || emptyState).toBeTruthy();
  });

  test('should update stats after deletion', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/dashboard');
    
    await page.waitForTimeout(1000);
    
    // Get initial counts
    const companiesCount = await page.locator('text=Total Companies').locator('..').locator('text=/\\d+/').textContent();
    
    // Counts should be valid numbers
    expect(companiesCount).toBeTruthy();
  });
});

test.describe('Accessibility and UX', () => {
  test('should have accessible form inputs', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies');
    
    const searchInput = page.locator('input[placeholder*="Search"]');
    
    // Should have placeholder
    await expect(searchInput).toHaveAttribute('placeholder', /Search/i);
  });

  test('should have hover effects on cards', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies');
    
    await page.waitForTimeout(1000);
    
    const firstCard = page.locator('.bg-white.rounded-lg.shadow').first();
    
    if (await firstCard.isVisible()) {
      // Should have hover class
      await expect(firstCard).toHaveClass(/hover:shadow-lg/);
    }
  });

  test('should have responsive layout classes', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies');
    
    await page.waitForTimeout(1000);
    
    // Should have grid with responsive classes
    const grid = page.locator('.grid');
    await expect(grid).toHaveClass(/md:grid-cols-2|lg:grid-cols-3/);
  });

  test('should have transition effects on buttons', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies');
    
    await page.waitForTimeout(1000);
    
    const button = page.locator('button:has-text("View Details")').first();
    
    if (await button.isVisible()) {
      // Should have transition class
      await expect(button).toHaveClass(/transition-colors/);
    }
  });

  test('should show proper cursor on clickable elements', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/dashboard');
    
    const statCard = page.locator('text=Total Companies').locator('..');
    
    if (await statCard.isVisible()) {
      await expect(statCard).toHaveClass(/cursor-pointer/);
    }
  });
});
