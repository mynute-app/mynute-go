import { test, expect } from './fixtures';

test.describe('Login Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/admin');
  });

  test('should display login form', async ({ page }) => {
    // Check for main heading
    await expect(page.locator('h1')).toContainText('Mynute Admin');
    
    // Check for email input
    await expect(page.locator('input[type="email"]')).toBeVisible();
    
    // Check for password input
    await expect(page.locator('input[type="password"]')).toBeVisible();
    
    // Check for login button
    await expect(page.locator('button[type="submit"]')).toContainText('Login');
  });

  test('should show validation for empty fields', async ({ page }) => {
    // Click login without filling fields
    await page.click('button[type="submit"]');
    
    // HTML5 validation should prevent submission
    const emailInput = page.locator('input[type="email"]');
    await expect(emailInput).toHaveAttribute('required', '');
  });

  test('should show error for invalid credentials', async ({ page }) => {
    // Fill with invalid credentials
    await page.fill('input[type="email"]', 'invalid@example.com');
    await page.fill('input[type="password"]', 'wrongpassword');
    
    // Submit form
    await page.click('button[type="submit"]');
    
    // Wait for error message
    await expect(page.locator('.bg-red-50')).toBeVisible({ timeout: 5000 });
  });

  test('should login successfully with valid credentials', async ({ page }) => {
    // Fill login form
    await page.fill('input[type="email"]', 'admin@mynute.com');
    await page.fill('input[type="password"]', 'admin123');
    
    // Submit form
    await page.click('button[type="submit"]');
    
    // Should redirect to dashboard
    await page.waitForURL('/admin/', { timeout: 5000 });
    
    // Check for dashboard heading
    await expect(page.locator('h1')).toContainText('Dashboard');
  });

  test('should show loading state during login', async ({ page }) => {
    // Fill login form
    await page.fill('input[type="email"]', 'admin@mynute.com');
    await page.fill('input[type="password"]', 'admin123');
    
    // Click login button
    const loginButton = page.locator('button[type="submit"]');
    await loginButton.click();
    
    // Button should show loading text
    await expect(loginButton).toContainText(/Logging in/);
  });
});
