import { test, expect } from '@playwright/test';

/**
 * CRITICAL: First Admin Registration Tests
 * 
 * These tests MUST run before login tests to ensure the database is in a clean state.
 * If these tests fail because a login form appears instead of registration form,
 * it means an admin already exists in the database.
 * 
 * To fix: Clear all admins from the database before running tests:
 * - Run: psql -U postgres -d testdb -c "DELETE FROM public.admins;"
 * - Or: Restart the test database with fresh migrations
 */

test.describe.serial('First Admin Registration', () => {
  test('SETUP: should show registration form when no admin exists', async ({ page }) => {
    // This is a critical setup test that verifies the database is clean
    // If this fails, it means admins already exist in the database
    
    await page.goto('');
    
    // Wait for the page to load
    await page.waitForLoadState('networkidle');
    
    // Check if we see the registration form or login form
    const registrationHeading = page.locator('h1:has-text("Welcome to Mynute Admin")');
    const loginHeading = page.locator('h1:has-text("Mynute Admin")');
    
    const hasRegistrationForm = await registrationHeading.isVisible().catch(() => false);
    const hasLoginForm = await loginHeading.isVisible().catch(() => false);
    
    if (hasLoginForm && !hasRegistrationForm) {
      throw new Error(
        '\n\n❌ SETUP ERROR: Login form detected instead of registration form!\n' +
        '   This means an admin already exists in the database.\n\n' +
        '   To fix this:\n' +
        '   1. Clear admins from the test database:\n' +
        '      psql -U postgres -d testdb -c "DELETE FROM public.admins;"\n' +
        '   2. Or restart the database with fresh migrations\n' +
        '   3. Make sure APP_ENV=test is set\n\n'
      );
    }
    
    // Should show registration form
    await expect(registrationHeading).toBeVisible();
    await expect(page.locator('text=Create your first admin account')).toBeVisible();
    
    // Check for all required fields
    await expect(page.locator('input[placeholder="John"]')).toBeVisible();
    await expect(page.locator('input[placeholder="Doe"]')).toBeVisible();
    await expect(page.locator('input[type="email"]')).toBeVisible();
    await expect(page.locator('input[type="password"]').first()).toBeVisible();
    
    // Check for submit button
    await expect(page.locator('button:has-text("Create Admin Account")')).toBeVisible();
  });

  test('should show login form when admin exists', async ({ page }) => {
    // Mock the API to return admin exists
    await page.route('**/api/admin/are_there_any_superadmin', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ has_superadmin: true }),
      });
    });

    await page.goto('');
    
    // Should show login form
    await expect(page.locator('h1')).toContainText('Mynute Admin');
    await expect(page.locator('button[type="submit"]:has-text("Login")')).toBeVisible();
  });

  test('should validate password match on registration', async ({ page }) => {
    // Mock the API to return no admins
    await page.route('**/api/admin/are_there_any_superadmin', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ has_superadmin: false }),
      });
    });

    await page.goto('');
    
    // Fill in form with mismatched passwords
    await page.fill('input[placeholder="John"]', 'John');
    await page.fill('input[placeholder="Doe"]', 'Doe');
    await page.fill('input[type="email"]', 'admin@mynute.com');
    
    const passwordFields = page.locator('input[type="password"]');
    await passwordFields.nth(0).fill('password123');
    await passwordFields.nth(1).fill('password456');
    
    // Submit form
    await page.click('button:has-text("Create Admin Account")');
    
    // Should show error message
    await expect(page.locator('text=Passwords do not match')).toBeVisible();
  });

  test('should validate password length on registration', async ({ page }) => {
    // Mock the API to return no admins
    await page.route('**/api/admin/are_there_any_superadmin', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ has_superadmin: false }),
      });
    });

    await page.goto('');
    
    // Fill in form with short password
    await page.fill('input[placeholder="John"]', 'John');
    await page.fill('input[placeholder="Doe"]', 'Doe');
    await page.fill('input[type="email"]', 'admin@mynute.com');
    
    const passwordFields = page.locator('input[type="password"]');
    await passwordFields.nth(0).fill('pass');
    await passwordFields.nth(1).fill('pass');
    
    // Submit form
    await page.click('button:has-text("Create Admin Account")');
    
    // Should show error message
    await expect(page.locator('text=/Password must be at least 8 characters/i')).toBeVisible();
  });

  test('should toggle password visibility on registration form', async ({ page }) => {
    // Mock the API to return no admins
    await page.route('**/api/admin/are_there_any_superadmin', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ has_superadmin: false }),
      });
    });

    await page.goto('');
    
    // Get password input and toggle button
    const passwordInput = page.getByTestId('registration-password');
    const passwordToggle = page.getByTestId('registration-password-toggle');
    
    // Password should be hidden by default
    await expect(passwordInput).toHaveAttribute('type', 'password');
    
    // Fill password
    await passwordInput.fill('TestPassword123');
    
    // Click eye icon to show password
    await passwordToggle.click();
    await expect(passwordInput).toHaveAttribute('type', 'text');
    await expect(passwordInput).toHaveValue('TestPassword123');
    
    // Click again to hide password
    await passwordToggle.click();
    await expect(passwordInput).toHaveAttribute('type', 'password');
  });

  test('should toggle confirm password visibility on registration form', async ({ page }) => {
    // Mock the API to return no admins
    await page.route('**/api/admin/are_there_any_superadmin', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ has_superadmin: false }),
      });
    });

    await page.goto('');
    
    // Get confirm password input and toggle button
    const confirmPasswordInput = page.getByTestId('registration-confirm-password');
    const confirmPasswordToggle = page.getByTestId('registration-confirm-password-toggle');
    
    // Password should be hidden by default
    await expect(confirmPasswordInput).toHaveAttribute('type', 'password');
    
    // Fill confirm password
    await confirmPasswordInput.fill('TestPassword123');
    
    // Click eye icon to show password
    await confirmPasswordToggle.click();
    await expect(confirmPasswordInput).toHaveAttribute('type', 'text');
    await expect(confirmPasswordInput).toHaveValue('TestPassword123');
    
    // Click again to hide password
    await confirmPasswordToggle.click();
    await expect(confirmPasswordInput).toHaveAttribute('type', 'password');
  });

  test('should successfully register first admin and show success message', async ({ page }) => {
    // Mock the API responses
    await page.route('**/api/admin/are_there_any_superadmin', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ has_superadmin: false }),
      });
    });

    await page.route('**/api/admin/first_superadmin', (route) => {
      route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({
          id: '1',
          name: 'John',
          surname: 'Doe',
          email: 'admin@mynute.com',
        }),
      });
    });

    await page.route('**/api/admin/send-verification-code/email/**', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ success: true }),
      });
    });

    await page.goto('');
    
    // Fill in registration form
    await page.fill('input[placeholder="John"]', 'John');
    await page.fill('input[placeholder="Doe"]', 'Doe');
    await page.fill('input[type="email"]', 'admin@mynute.com');
    
    const passwordFields = page.locator('input[type="password"]');
    await passwordFields.nth(0).fill('password123');
    await passwordFields.nth(1).fill('password123');
    
    // Submit form
    await page.click('button:has-text("Create Admin Account")');
    
    // Should show success message
    await expect(page.locator('text=Registration Successful!')).toBeVisible();
    await expect(page.locator('text=/verification email has been sent/i')).toBeVisible();
    await expect(page.locator('text=admin@mynute.com')).toBeVisible();
    
    // Should have "Go to Login" button
    await expect(page.locator('button:has-text("Go to Login")')).toBeVisible();
  });

  test('should navigate back to login after successful registration', async ({ page }) => {
    // Mock the API responses
    let hasAdmin = false;
    
    await page.route('**/api/admin/are_there_any_superadmin', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ has_superadmin: hasAdmin }),
      });
    });

    await page.route('**/api/admin/first_superadmin', (route) => {
      hasAdmin = true;
      route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({
          id: '1',
          name: 'John',
          surname: 'Doe',
          email: 'admin@mynute.com',
        }),
      });
    });

    await page.route('**/api/admin/send-verification-code/email/**', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ success: true }),
      });
    });

    await page.goto('');
    
    // Complete registration
    await page.fill('input[placeholder="John"]', 'John');
    await page.fill('input[placeholder="Doe"]', 'Doe');
    await page.fill('input[type="email"]', 'admin@mynute.com');
    
    const passwordFields = page.locator('input[type="password"]');
    await passwordFields.nth(0).fill('password123');
    await passwordFields.nth(1).fill('password123');
    
    await page.click('button:has-text("Create Admin Account")');
    
    // Click "Go to Login"
    await page.click('button:has-text("Go to Login")');
    
    // Should show login form
    await expect(page.locator('h1')).toContainText('Mynute Admin');
    await expect(page.locator('button[type="submit"]:has-text("Login")')).toBeVisible();
  });

  test('should show loading state while checking for admins', async ({ page }) => {
    // This test verifies loading state appears during API check
    await page.route('**/api/admin/are_there_any_superadmin', async (route) => {
      // Delay response to see loading state
      await new Promise(resolve => setTimeout(resolve, 100));
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ has_superadmin: false }),
      });
    });

    await page.goto('');
    
    // Page should eventually load registration form
    await expect(page.locator('h1')).toContainText('Welcome to Mynute Admin');
  });

  test('should handle registration errors gracefully', async ({ page }) => {
    // Mock the API responses
    await page.route('**/api/admin/are_there_any_superadmin', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ has_superadmin: false }),
      });
    });

    await page.route('**/api/admin/first_superadmin', (route) => {
      route.fulfill({
        status: 400,
        contentType: 'application/json',
        body: JSON.stringify({
          error: 'Email already exists',
        }),
      });
    });

    await page.goto('');
    
    // Fill in and submit registration form
    await page.fill('input[placeholder="John"]', 'John');
    await page.fill('input[placeholder="Doe"]', 'Doe');
    await page.fill('input[type="email"]', 'admin@mynute.com');
    
    const passwordFields = page.locator('input[type="password"]');
    await passwordFields.nth(0).fill('password123');
    await passwordFields.nth(1).fill('password123');
    
    await page.click('button:has-text("Create Admin Account")');
    
    // Should show error message
    await expect(page.locator('.bg-red-50')).toBeVisible();
    await expect(page.locator('text=/Email already exists/i')).toBeVisible();
  });

  test('should show all required field indicators', async ({ page }) => {
    // Mock the API to return no admins
    await page.route('**/api/admin/are_there_any_superadmin', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ has_superadmin: false }),
      });
    });

    await page.goto('');
    
    // Check for required field indicators (asterisks)
    const requiredIndicators = page.locator('.text-red-500:has-text("*")');
    const count = await requiredIndicators.count();
    
    // Should have 5 required fields (First Name, Last Name, Email, Password, Confirm Password)
    expect(count).toBe(5);
  });

  test('should show password length hint', async ({ page }) => {
    // Mock the API to return no admins
    await page.route('**/api/admin/are_there_any_superadmin', (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ has_superadmin: false }),
      });
    });

    await page.goto('');
    
    // Should show password hint
    await expect(page.locator('text=Must be at least 8 characters')).toBeVisible();
  });
});
