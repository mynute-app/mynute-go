# E2E Tests for Mynute Admin

Playwright-based end-to-end tests for the admin panel.

## ⚠️ CRITICAL: Test Environment Setup

### Environment Requirements

**Tests ONLY run when `APP_ENV=test`**

```bash
# Set in .env file
APP_ENV=test

# Or set as environment variable
export APP_ENV=test  # Linux/Mac
set APP_ENV=test     # Windows CMD
$env:APP_ENV="test"  # Windows PowerShell
```

If `APP_ENV` is not set to `test`, the test suite will immediately exit with an error.

### Database Setup

**Before running tests, ensure the test database is clean:**

1. **Clear existing admins** (if tests fail with "admin already exists"):
   ```bash
   psql -U postgres -d testdb -c "DELETE FROM public.admins;"
   ```

2. **Or reset the entire database**:
   ```bash
   # Drop and recreate
   psql -U postgres -c "DROP DATABASE IF EXISTS testdb;"
   psql -U postgres -c "CREATE DATABASE testdb;"
   
   # Run migrations with APP_ENV=test
   APP_ENV=test go run main.go
   ```

## Test Execution Order

Tests run in a **specific order** to ensure proper setup:

1. **`first-admin-registration.spec.ts`** - MUST run first
   - Verifies database is clean (no admins exist)
   - Creates the first admin account
   - If this fails, it means admins already exist in the database

2. **`login.spec.ts`** - Authentication tests
3. **`navigation.spec.ts`** - Navigation and routing
4. **`dashboard.spec.ts`** - Dashboard functionality
5. **Other test files** - Feature-specific tests

## Running Tests

```bash
# Run all tests in correct order
npm test

# Run tests in UI mode (interactive)
npm run test:ui

# Run tests in headed mode (see browser)
npm run test:headed

# Debug tests
npm run test:debug

# View test report
npm run test:report
```

## Troubleshooting

### "Login form detected instead of registration form"

This error in `first-admin-registration.spec.ts` means:
- An admin already exists in the database
- The test database was not properly cleared

**Solution:**
```bash
# Clear admins from database
psql -U postgres -d testdb -c "DELETE FROM public.admins;"

# Re-run tests
npm test
```

### "Tests can only run when APP_ENV=test"

**Solution:**
```bash
# Set environment variable
export APP_ENV=test  # Linux/Mac
$env:APP_ENV="test"  # Windows PowerShell

# Or add to .env file
echo "APP_ENV=test" >> .env

# Then run tests
npm test
```

### Tests failing with wrong credentials

The default admin credentials created by seeding are:
- Email: `admin@mynute.com`
- Password: `Admin@123456`

Make sure these match in `fixtures.ts`.

## Test Structure

- `tests/fixtures.ts` - Custom fixtures and test utilities
- `tests/first-admin-registration.spec.ts` - ⚠️ MUST RUN FIRST
- `tests/login.spec.ts` - Login page tests
- `tests/dashboard.spec.ts` - Dashboard tests
- `tests/users.spec.ts` - Admin users page tests
- `tests/navigation.spec.ts` - Navigation and layout tests
- `tests/routing-basepath.spec.ts` - Router basepath tests
- `tests/companies.spec.ts` - Companies page tests
- `tests/company-detail.spec.ts` - Company detail page tests
- `tests/clients.spec.ts` - Clients page tests
- `tests/integration.spec.ts` - Integration and edge case tests

## Writing Tests

```typescript
import { test, expect } from './fixtures';

test('my test', async ({ authenticatedPage: page }) => {
  // Test code here - already logged in
  await expect(page.locator('h1')).toContainText('Dashboard');
});
```

## Fixtures

### `authenticatedPage`

Automatically logs in before the test runs:

```typescript
test('authenticated test', async ({ authenticatedPage: page }) => {
  // Already logged in with admin@mynute.com
});
```

## Configuration

Edit `playwright.config.ts` to:
- Change base URL
- Add/remove browsers
- Configure timeouts
- Set up CI settings

## CI/CD

Tests run automatically on CI with the `webServer` configuration starting your Go backend.

### Required CI Environment Variables

```yaml
env:
  APP_ENV: test
  DATABASE_URL: postgresql://user:pass@localhost:5432/testdb
```
