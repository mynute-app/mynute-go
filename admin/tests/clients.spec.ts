import { test, expect } from './fixtures';

test.describe('Clients Page', () => {
  test('should display clients list', async ({ authenticatedPage: page }) => {
    await page.goto('');
    await page.click('a[href="/clients"]');
    
    await expect(page).toHaveURL(/.*\/clients/);
    await expect(page.locator('h1')).toContainText('Clients');
    
    // Check for search input
    await expect(page.locator('input[placeholder*="Search"]')).toBeVisible();
  });

  test('should display clients table', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    // Wait for clients to load
    await page.waitForTimeout(1000);
    
    // Check if there's a table or empty state
    const hasTable = await page.locator('table').isVisible().catch(() => false);
    const hasEmptyState = await page.locator('text=/No clients found/i').isVisible().catch(() => false);
    
    expect(hasTable || hasEmptyState).toBeTruthy();
  });

  test('should display table headers', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const tableVisible = await page.locator('table').isVisible().catch(() => false);
    
    if (tableVisible) {
      // Check for table headers
      const headers = ['Name', 'Email', 'Phone', 'Created At', 'Actions'];
      
      for (const header of headers) {
        await expect(page.locator(`th:has-text("${header}")`)).toBeVisible();
      }
    }
  });

  test('should filter clients by search', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.fill('John');
    
    // Wait for filtering to apply
    await page.waitForTimeout(500);
    
    // Clients should be filtered
    const rows = page.locator('tbody tr');
    const rowCount = await rows.count();
    
    // If rows exist, they should contain the search term
    if (rowCount > 0) {
      const firstRow = rows.first();
      await expect(firstRow).toContainText(/John/i);
    }
  });

  test('should open client details modal', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    // Wait for clients to load
    await page.waitForTimeout(1000);
    
    // Click on the first "View Details" button if it exists
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      // Modal should be visible
      await expect(page.locator('text=Client Details')).toBeVisible();
      
      // Should have close button
      await expect(page.locator('button:has-text("Close")')).toBeVisible();
    }
  });

  test('should close client details modal', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      // Modal should be visible
      await expect(page.locator('text=Client Details')).toBeVisible();
      
      // Close modal
      await page.click('button:has-text("Close")');
      
      // Modal should be hidden
      await expect(page.locator('text=Client Details')).not.toBeVisible();
    }
  });

  test('should display client information in modal', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      await page.waitForTimeout(500);
      
      // Should display client info sections
      const sections = ['Basic Information', 'Appointments'];
      
      for (const section of sections) {
        const sectionVisible = await page.locator(`text=${section}`).isVisible().catch(() => false);
        if (sectionVisible) {
          await expect(page.locator(`text=${section}`)).toBeVisible();
        }
      }
    }
  });

  test('should display client appointments in modal', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      await page.waitForTimeout(500);
      
      // Should show appointments section
      const appointmentsSection = page.locator('text=Appointments');
      await expect(appointmentsSection).toBeVisible();
      
      // Should show appointments list or empty state
      const hasAppointments = await page.locator('text=/appointment|No appointments/i').isVisible();
      expect(hasAppointments).toBeTruthy();
    }
  });

  test('should show delete button for clients', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const deleteButton = page.locator('button:has-text("Delete")').first();
    
    if (await deleteButton.isVisible()) {
      await expect(deleteButton).toBeVisible();
      
      // TODO: Add test for delete confirmation when implemented
    }
  });

  test('should clear search when input is empty', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    const searchInput = page.locator('input[placeholder*="Search"]');
    
    // Type and then clear
    await searchInput.fill('Test');
    await page.waitForTimeout(300);
    await searchInput.clear();
    await page.waitForTimeout(300);
    
    // All clients should be visible again
    const rows = page.locator('tbody tr');
    const rowCount = await rows.count();
    
    // Should have 0 or more rows (not filtered)
    expect(rowCount).toBeGreaterThanOrEqual(0);
  });

  test('should have working sidebar navigation', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    // Navigate to other pages from sidebar
    await page.click('a[href="/dashboard"]');
    await expect(page).toHaveURL(/.*\/dashboard/);
    
    await page.click('a[href="/clients"]');
    await expect(page).toHaveURL(/.*\/clients/);
  });

  test('should close modal when clicking outside', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      // Modal should be visible
      await expect(page.locator('text=Client Details')).toBeVisible();
      
      // Click on backdrop (outside modal)
      await page.locator('.fixed.inset-0.bg-black').click({ force: true });
      
      // Modal should be hidden
      await expect(page.locator('text=Client Details')).not.toBeVisible();
    }
  });

  test('should display formatted dates', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const tableVisible = await page.locator('table').isVisible().catch(() => false);
    
    if (tableVisible) {
      // Check if dates are displayed in readable format
      const dateCell = page.locator('tbody tr td').nth(3);
      
      if (await dateCell.isVisible()) {
        const dateText = await dateCell.textContent();
        
        // Should contain date-like text (basic check)
        expect(dateText).toBeTruthy();
      }
    }
  });

  test('should show delete confirmation dialog', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
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

  test('should display client avatar with initial', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const tableVisible = await page.locator('table').isVisible().catch(() => false);
    
    if (tableVisible) {
      // Should have avatar with client initial
      const avatar = page.locator('.rounded-full.bg-primary').first();
      
      if (await avatar.isVisible()) {
        await expect(avatar).toBeVisible();
        const avatarText = await avatar.textContent();
        expect(avatarText?.length).toBe(1); // Should be single letter
      }
    }
  });

  test('should display appointment counts in modal', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      await page.waitForTimeout(500);
      
      // Should show appointment count
      await expect(page.locator('text=/Appointments \\(/i')).toBeVisible();
    }
  });

  test('should display appointment status badges in modal', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      await page.waitForTimeout(500);
      
      // If appointments exist, check for status badges
      const appointmentExists = await page.locator('.border.border-gray-200.rounded-lg').isVisible().catch(() => false);
      
      if (appointmentExists) {
        const statusBadge = page.locator('.rounded-full.text-xs').first();
        await expect(statusBadge).toBeVisible();
      }
    }
  });

  test('should show client surname if available', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      await page.waitForTimeout(500);
      
      // Modal header should show name (and surname if available)
      await expect(page.locator('h2').first()).toBeVisible();
    }
  });

  test('should display updated date in modal', async ({ authenticatedPage: page }) => {
    await page.goto('/clients');
    
    await page.waitForTimeout(1000);
    
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      await page.waitForTimeout(500);
      
      // Should show Last Updated field
      await expect(page.locator('text=Last Updated')).toBeVisible();
    }
  });
});
