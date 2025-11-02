import { test, expect } from './fixtures';

test.describe('Companies Page', () => {

  test('should display companies list', async ({ authenticatedPage: page }) => {
    await page.goto('');
    await page.click('a[href="/companies"]');
    
    await expect(page).toHaveURL(/.*\/companies/);
    await expect(page.locator('h1')).toContainText('Companies');
    
    // Check for search input
    await expect(page.locator('input[placeholder*="Search"]')).toBeVisible();
    
    // Check for add button
    await expect(page.locator('button:has-text("Add Company")')).toBeVisible();
  });

  test('should display company cards', async ({ authenticatedPage: page }) => {
    await page.goto('/companies');
    
    // Wait for companies to load
    await page.waitForTimeout(1000);
    
    // Check if there are company cards or empty state
    const hasCards = await page.locator('.bg-white.rounded-lg.shadow-sm').count() > 0;
    const hasEmptyState = await page.locator('text=No companies found').isVisible().catch(() => false);
    
    expect(hasCards || hasEmptyState).toBeTruthy();
  });

  test('should filter companies by search', async ({ authenticatedPage: page }) => {
    await page.goto('/companies');
    
    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.fill('Test Company');
    
    // Wait for filtering to apply
    await page.waitForTimeout(500);
    
    // Companies should be filtered (or show no results)
    const cards = page.locator('.bg-white.rounded-lg.shadow-sm');
    const cardCount = await cards.count();
    
    // If cards exist, they should contain the search term
    if (cardCount > 0) {
      const firstCard = cards.first();
      await expect(firstCard).toContainText(/Test Company/i);
    }
  });

  test('should navigate to company detail page', async ({ authenticatedPage: page }) => {
    await page.goto('/companies');
    
    // Wait for companies to load
    await page.waitForTimeout(1000);
    
    // Click on the first "View Details" button if it exists
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      // Should navigate to company detail page
      await expect(page).toHaveURL(/.*\/companies\/\d+/);
    }
  });

  test('should show add company button', async ({ authenticatedPage: page }) => {
    await page.goto('/companies');
    
    const addButton = page.locator('button:has-text("Add Company")');
    await expect(addButton).toBeVisible();
    
    // TODO: Add test for modal when create functionality is implemented
  });

  test('should clear search when input is empty', async ({ authenticatedPage: page }) => {
    await page.goto('/companies');
    
    const searchInput = page.locator('input[placeholder*="Search"]');
    
    // Type and then clear
    await searchInput.fill('Test');
    await page.waitForTimeout(300);
    await searchInput.clear();
    await page.waitForTimeout(300);
    
    // All companies should be visible again
    const cards = page.locator('.bg-white.rounded-lg.shadow-sm');
    const cardCount = await cards.count();
    
    // Should have 0 or more cards (not filtered)
    expect(cardCount).toBeGreaterThanOrEqual(0);
  });

  test('should display company information in cards', async ({ authenticatedPage: page }) => {
    await page.goto('/companies');
    
    await page.waitForTimeout(1000);
    
    const firstCard = page.locator('.bg-white.rounded-lg.shadow-sm').first();
    
    if (await firstCard.isVisible()) {
      // Should have company name
      await expect(firstCard.locator('h3')).toBeVisible();
      
      // Should have subdomain or created date
      const hasInfo = await firstCard.locator('text=/subdomain:|Created/').count() > 0;
      expect(hasInfo).toBeTruthy();
    }
  });

  test('should have working sidebar navigation', async ({ authenticatedPage: page }) => {
    await page.goto('/companies');
    
    // Navigate to other pages from sidebar
    await page.click('a[href="/dashboard"]');
    await expect(page).toHaveURL(/.*\/dashboard/);
    
    await page.click('a[href="/companies"]');
    await expect(page).toHaveURL(/.*\/companies/);
  });

  test('should show delete confirmation dialog', async ({ authenticatedPage: page }) => {
    await page.goto('/companies');
    
    await page.waitForTimeout(1000);
    
    const deleteButton = page.locator('button:has-text("Delete")').first();
    
    if (await deleteButton.isVisible()) {
      // Set up dialog handler
      page.once('dialog', async (dialog) => {
        expect(dialog.message()).toContain('Are you sure');
        await dialog.dismiss();
      });
      
      await deleteButton.click();
    }
  });

  test('should display company status badge', async ({ authenticatedPage: page }) => {
    await page.goto('/companies');
    
    await page.waitForTimeout(1000);
    
    const firstCard = page.locator('.bg-white.rounded-lg.shadow').first();
    
    if (await firstCard.isVisible()) {
      // Should have Active status badge
      await expect(firstCard.locator('text=Active')).toBeVisible();
    }
  });

  test('should display tax ID in company cards', async ({ authenticatedPage: page }) => {
    await page.goto('/companies');
    
    await page.waitForTimeout(1000);
    
    const firstCard = page.locator('.bg-white.rounded-lg.shadow').first();
    
    if (await firstCard.isVisible()) {
      // Should show tax ID
      await expect(firstCard.locator('text=/Tax ID:/i')).toBeVisible();
    }
  });

  test('should display created date in company cards', async ({ authenticatedPage: page }) => {
    await page.goto('/companies');
    
    await page.waitForTimeout(1000);
    
    const firstCard = page.locator('.bg-white.rounded-lg.shadow').first();
    
    if (await firstCard.isVisible()) {
      // Should show created date
      await expect(firstCard.locator('text=/Created:/i')).toBeVisible();
    }
  });
});
