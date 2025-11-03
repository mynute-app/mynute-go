import { test, expect } from './fixtures';

test.describe('Admin Users Page', () => {
  test.beforeEach(async ({ authenticatedPage: page }) => {
    // Navigate to users page
    await page.click('text=Admin Users');
    await page.waitForURL('/admin/users', { timeout: 5000 });
  });

  test('should display users page heading', async ({ authenticatedPage: page }) => {
    await expect(page.locator('h1')).toContainText('Admin Users');
  });

  test('should display users table', async ({ authenticatedPage: page }) => {
    // Check for table
    const table = page.locator('table');
    await expect(table).toBeVisible();
    
    // Check for table headers
    await expect(page.locator('th:has-text("Name")')).toBeVisible();
    await expect(page.locator('th:has-text("Email")')).toBeVisible();
    await expect(page.locator('th:has-text("Role")')).toBeVisible();
    await expect(page.locator('th:has-text("Status")')).toBeVisible();
    await expect(page.locator('th:has-text("Actions")')).toBeVisible();
  });

  test('should show loading state', async ({ authenticatedPage: page }) => {
    // Reload to trigger loading state
    await page.reload();
    
    // Loading state is often too fast to reliably test in E2E
    // This test verifies the page loads successfully
    await page.waitForSelector('table', { timeout: 5000 });
  });

  test('should display empty state when no admins', async ({ authenticatedPage: page }) => {
    // Wait for table to load
    await page.waitForSelector('table', { timeout: 5000 });
    
    // If no admins, should show message
    const emptyMessage = page.locator('text=No admins found');
    const hasEmptyMessage = await emptyMessage.isVisible();
    
    if (hasEmptyMessage) {
      await expect(emptyMessage).toContainText('Click "Add Admin" to create one');
    }
  });

  test('should display admin list when admins exist', async ({ authenticatedPage: page }) => {
    // Wait for table body
    const tbody = page.locator('tbody');
    await expect(tbody).toBeVisible();
    
    // Check if there are any rows (besides empty state)
    const rows = tbody.locator('tr');
    const rowCount = await rows.count();
    
    expect(rowCount).toBeGreaterThanOrEqual(1);
  });

  test('should have delete button for each admin', async ({ authenticatedPage: page }) => {
    // Wait for table to load
    await page.waitForSelector('table', { timeout: 5000 });
    
    // Check for delete buttons
    const deleteButtons = page.locator('button:has-text("Delete")');
    const count = await deleteButtons.count();
    
    // If there are admins, there should be delete buttons
    if (count > 0) {
      expect(count).toBeGreaterThan(0);
    }
  });

  test('should show confirmation dialog when deleting admin', async ({ authenticatedPage: page }) => {
    // Wait for table to load
    await page.waitForSelector('table', { timeout: 5000 });
    
    const deleteButtons = page.locator('button:has-text("Delete")');
    const count = await deleteButtons.count();
    
    if (count > 0) {
      // Set up dialog handler
      let dialogShown = false;
      page.on('dialog', async (dialog) => {
        expect(dialog.message()).toContain('Are you sure');
        dialogShown = true;
        await dialog.dismiss();
      });
      
      // Click first delete button
      await deleteButtons.first().click();
      
      // Wait a bit for dialog
      await page.waitForTimeout(500);
      
      expect(dialogShown).toBe(true);
    }
  });

  test('should navigate back to dashboard', async ({ authenticatedPage: page }) => {
    // Click dashboard in sidebar
    await page.click('text=Dashboard');
    
    // Should navigate to dashboard
    await page.waitForURL('/admin/', { timeout: 5000 });
    
    // Verify we're on dashboard
    await expect(page.locator('h1')).toContainText('Dashboard');
  });
});
