# E2E Tests for Mynute Admin

Playwright-based end-to-end tests for the admin panel.

## Running Tests

```bash
# Run all tests
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

## Test Structure

- `tests/fixtures.ts` - Custom fixtures and test utilities
- `tests/login.spec.ts` - Login page tests
- `tests/dashboard.spec.ts` - Dashboard tests
- `tests/users.spec.ts` - Admin users page tests
- `tests/navigation.spec.ts` - Navigation and layout tests

## Writing Tests

```typescript
import { test, expect } from './fixtures';

test('my test', async ({ authenticatedPage: page }) => {
  // Test code here
  await expect(page.locator('h1')).toContainText('Dashboard');
});
```

## Fixtures

### `authenticatedPage`

Automatically logs in before the test runs:

```typescript
test('authenticated test', async ({ authenticatedPage: page }) => {
  // Already logged in
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
