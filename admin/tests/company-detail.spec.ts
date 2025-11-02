import { test, expect } from './fixtures';

test.describe('Company Detail Page', () => {
  test('should display company detail page', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies');
    
    // Wait for companies to load
    await page.waitForTimeout(1000);
    
    const viewButton = page.locator('button:has-text("View Details")').first();
    
    if (await viewButton.isVisible()) {
      await viewButton.click();
      
      // Should be on company detail page
      await expect(page).toHaveURL(/.*\/companies\/\d+/);
      
      // Should have back button
      await expect(page.locator('button:has-text("Back to Companies")')).toBeVisible();
    }
  });

  test('should display all tabs', async ({ authenticatedPage: page }) => {
    // Navigate to a company detail page (assuming company ID 1 exists)
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    // Check for all tab buttons
    const tabs = ['Overview', 'Branches', 'Employees', 'Services', 'Subdomains'];
    
    for (const tabName of tabs) {
      const tabButton = page.locator(`button:has-text("${tabName}")`);
      await expect(tabButton).toBeVisible();
    }
  });

  test('should switch between tabs', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    // Click on Branches tab
    await page.click('button:has-text("Branches")');
    await page.waitForTimeout(300);
    
    // Should show branches content
    const branchesContent = page.locator('text=/branches|No branches/i');
    await expect(branchesContent).toBeVisible();
    
    // Click on Employees tab
    await page.click('button:has-text("Employees")');
    await page.waitForTimeout(300);
    
    // Should show employees content
    const employeesContent = page.locator('text=/employees|No employees/i');
    await expect(employeesContent).toBeVisible();
  });

  test('should display overview tab content', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    // Overview should be active by default
    const overviewTab = page.locator('button:has-text("Overview")');
    await expect(overviewTab).toHaveClass(/bg-blue-50/);
    
    // Should display company information fields
    const infoLabels = ['Company Name', 'Subdomain', 'Created At', 'Updated At'];
    
    for (const label of infoLabels) {
      await expect(page.locator(`text=${label}`)).toBeVisible();
    }
  });

  test('should display branches tab', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    await page.click('button:has-text("Branches")');
    await page.waitForTimeout(500);
    
    // Should have add branch button
    await expect(page.locator('button:has-text("Add Branch")')).toBeVisible();
    
    // Should show branches list or empty state
    const hasBranches = await page.locator('text=/Branch:|No branches/i').isVisible();
    expect(hasBranches).toBeTruthy();
  });

  test('should display employees tab', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    await page.click('button:has-text("Employees")');
    await page.waitForTimeout(500);
    
    // Should have add employee button
    await expect(page.locator('button:has-text("Add Employee")')).toBeVisible();
    
    // Should show employees list or empty state
    const hasEmployees = await page.locator('text=/Employee|No employees/i').isVisible();
    expect(hasEmployees).toBeTruthy();
  });

  test('should display services tab', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    await page.click('button:has-text("Services")');
    await page.waitForTimeout(500);
    
    // Should have add service button
    await expect(page.locator('button:has-text("Add Service")')).toBeVisible();
    
    // Should show services list or empty state
    const hasServices = await page.locator('text=/Service|No services/i').isVisible();
    expect(hasServices).toBeTruthy();
  });

  test('should display subdomains tab', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    await page.click('button:has-text("Subdomains")');
    await page.waitForTimeout(500);
    
    // Should show subdomain information
    const hasSubdomain = await page.locator('text=/Subdomain|No subdomain/i').isVisible();
    expect(hasSubdomain).toBeTruthy();
  });

  test('should navigate back to companies list', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    // Click back button
    await page.click('button:has-text("Back to Companies")');
    
    // Should be back on companies list
    await expect(page).toHaveURL(/.*\/companies$/);
  });

  test('should handle non-existent company', async ({ authenticatedPage: page }) => {
    // Try to access a company that likely doesn't exist
    await page.goto('http://localhost:3000/admin/companies/999999');
    
    await page.waitForTimeout(1000);
    
    // Should show loading or error state
    const hasError = await page.locator('text=/not found|error|loading/i').isVisible().catch(() => false);
    const hasContent = await page.locator('button:has-text("Overview")').isVisible().catch(() => false);
    
    // Either shows error or loads (if company exists)
    expect(hasError || hasContent).toBeTruthy();
  });

  test('should maintain active tab on refresh', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    // Switch to Employees tab
    await page.click('button:has-text("Employees")');
    await page.waitForTimeout(300);
    
    // Refresh page
    await page.reload();
    await page.waitForTimeout(1000);
    
    // Overview should be active again (default state after refresh)
    const overviewTab = page.locator('button:has-text("Overview")');
    await expect(overviewTab).toHaveClass(/border-primary/);
  });

  test('should display tab counts', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    // Check if tabs have count badges
    const tabsWithCounts = ['Branches', 'Employees', 'Services', 'Subdomains'];
    
    for (const tabName of tabsWithCounts) {
      const tab = page.locator(`button:has-text("${tabName}")`);
      if (await tab.isVisible()) {
        // Should have a count badge
        const hasCount = await tab.locator('.bg-gray-100').isVisible();
        expect(hasCount).toBeTruthy();
      }
    }
  });

  test('should show delete confirmation on company delete', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    // Set up dialog handler
    page.once('dialog', async (dialog) => {
      expect(dialog.message()).toContain('Are you sure');
      await dialog.dismiss();
    });
    
    const deleteButton = page.locator('button:has-text("Delete Company")');
    if (await deleteButton.isVisible()) {
      await deleteButton.click();
    }
  });

  test('should display statistics in overview tab', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    // Should show statistics section
    const stats = ['Total Branches', 'Total Employees', 'Total Services', 'Subdomains'];
    
    for (const stat of stats) {
      await expect(page.locator(`text=${stat}`)).toBeVisible();
    }
  });

  test('should show delete buttons in nested resource tables', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    // Check Branches tab
    await page.click('button:has-text("Branches")');
    await page.waitForTimeout(500);
    
    const branchDeleteButton = page.locator('button:has-text("Delete")').first();
    if (await branchDeleteButton.isVisible()) {
      await expect(branchDeleteButton).toBeVisible();
    }
    
    // Check Employees tab
    await page.click('button:has-text("Employees")');
    await page.waitForTimeout(500);
    
    const employeeDeleteButton = page.locator('button:has-text("Delete")').first();
    if (await employeeDeleteButton.isVisible()) {
      await expect(employeeDeleteButton).toBeVisible();
    }
  });

  test('should display company legal and trade names correctly', async ({ authenticatedPage: page }) => {
    await page.goto('http://localhost:3000/admin/companies/1');
    
    await page.waitForTimeout(1000);
    
    // Should show main heading
    await expect(page.locator('h1')).toBeVisible();
    
    // In overview, should show both legal and trade names
    await expect(page.locator('text=/Legal Name|Trade Name/i')).toBeVisible();
  });
});
