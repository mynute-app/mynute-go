# Integration Tests for Admin API

This directory contains integration tests for the admin management API endpoints. These tests verify the complete request-response cycle including HTTP routing, authentication, database operations, and response formatting.

## Test Files

- **`admin_test.go`** - Unit tests for admin logic (68 tests)
  - Password hashing and validation
  - Email validation
  - DTO structure validation
  - Edge cases for admin creation
  
- **`admin_integration_test.go`** - Integration tests for admin API endpoints (8 test suites, 30+ tests)
  - Full HTTP request/response cycle
  - Database CRUD operations
  - Authentication and authorization
  - Complete workflow testing

## Running Tests

### Run All Tests (Unit + Integration)

```powershell
# From auth directory
cd auth
go test ./... -v

# From controller directory
cd auth/api/controller
go test -v
```

### Run Only Unit Tests (Fast)

```powershell
# Skip integration tests using -short flag
go test ./... -v -short
```

### Run Only Integration Tests

```powershell
# Run tests matching "Integration" pattern
go test -v -run Integration
```

### Run Specific Test

```powershell
# Run a single test
go test -v -run TestCreateFirstAdminIntegration

# Run tests matching a pattern
go test -v -run "TestAdmin.*Integration"
```

### With Coverage

```powershell
go test ./... -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Integration Test Setup

Integration tests require a PostgreSQL test database. They will automatically **skip** if the database is not configured.

### Quick Setup

1. **Create a test database:**
   ```sql
   CREATE DATABASE auth_test;
   ```

2. **Set environment variable:**
   ```powershell
   # PowerShell
   $env:TEST_DATABASE_URL="postgresql://user:password@localhost:5432/auth_test?sslmode=disable"
   
   # CMD
   set TEST_DATABASE_URL=postgresql://user:password@localhost:5432/auth_test?sslmode=disable
   ```

3. **Run tests:**
   ```powershell
   go test -v -run Integration
   ```

### Docker PostgreSQL Setup (Recommended)

```powershell
# Start PostgreSQL in Docker
docker run --name auth-test-db `
  -e POSTGRES_USER=testuser `
  -e POSTGRES_PASSWORD=testpass `
  -e POSTGRES_DB=auth_test `
  -p 5433:5432 `
  -d postgres:15

# Set environment variable
$env:TEST_DATABASE_URL="postgresql://testuser:testpass@localhost:5433/auth_test?sslmode=disable"

# Run tests
go test -v -run Integration

# Cleanup
docker stop auth-test-db
docker rm auth-test-db
```

### Using Existing Database

If you already have the development database running:

```powershell
# Create a separate test database
docker exec -it mynute-postgres psql -U mynute -c "CREATE DATABASE auth_test;"

# Set environment variable
$env:TEST_DATABASE_URL="postgresql://mynute:mynute@localhost:5432/auth_test?sslmode=disable"

# Run tests
go test -v -run Integration
```

## Integration Test Coverage

### Endpoints Tested

| Endpoint | Method | Test | Description |
|----------|--------|------|-------------|
| `/users/admin/are_there_any_superadmin` | GET | `TestAreThereAnyAdminIntegration` | Check if any admin exists |
| `/users/admin/first_superadmin` | POST | `TestCreateFirstAdminIntegration` | Create first admin (no auth) |
| `/users/admin` | POST | `TestCreateAdminIntegration` | Create admin (requires auth) |
| `/users/admin/:id` | GET | `TestGetAdminByIdIntegration` | Get admin by ID |
| `/users/admin/:id` | PATCH | `TestUpdateAdminByIdIntegration` | Update admin |
| `/users/admin/:id` | DELETE | `TestDeleteAdminByIdIntegration` | Soft delete admin |
| `/users/admin` | GET | `TestListAdminsIntegration` | List all admins |

### Test Scenarios

#### `TestAreThereAnyAdminIntegration`
- ✅ Returns false when no admin exists
- ✅ Returns true when admin exists

#### `TestCreateFirstAdminIntegration`
- ✅ Creates first admin without authentication
- ✅ Rejects when admin already exists
- ✅ Validates required fields
- ✅ Rejects invalid email format
- ✅ Rejects weak passwords

#### `TestCreateAdminIntegration`
- ✅ Creates admin with valid superadmin token
- ✅ Rejects without authentication
- ✅ Rejects with invalid token
- ✅ Prevents duplicate emails

#### `TestGetAdminByIdIntegration`
- ✅ Retrieves admin with valid token
- ✅ Returns 404 for non-existent admin
- ✅ Rejects without authentication

#### `TestUpdateAdminByIdIntegration`
- ✅ Updates admin with valid token
- ✅ Validates email format on update

#### `TestDeleteAdminByIdIntegration`
- ✅ Soft deletes admin with valid token
- ✅ Rejects without authentication

#### `TestListAdminsIntegration`
- ✅ Lists all admins with valid token
- ✅ Excludes soft-deleted admins
- ✅ Rejects without authentication

#### `TestAdminWorkflowIntegration`
Complete end-to-end workflow:
1. Check no admin exists
2. Create first admin
3. Verify admin exists
4. Generate JWT token
5. Create second admin with auth
6. List all admins
7. Update admin
8. Delete admin
9. Verify final state

## How Integration Tests Work

### Test Flow

```
┌─────────────────────────────────────┐
│  1. Setup Test Database             │
│     - Connect to PostgreSQL         │
│     - Run migrations                │
│     - Clean existing data           │
└─────────────────────────────────────┘
                ↓
┌─────────────────────────────────────┐
│  2. Setup Fiber App                 │
│     - Create Fiber instance         │
│     - Register routes               │
│     - Add middleware                │
│     - Attach database to context    │
└─────────────────────────────────────┘
                ↓
┌─────────────────────────────────────┐
│  3. Create HTTP Request             │
│     - Set method (GET/POST/etc)     │
│     - Set path (/users/admin)       │
│     - Set headers (Content-Type)    │
│     - Set body (JSON)               │
│     - Add auth token if needed      │
└─────────────────────────────────────┘
                ↓
┌─────────────────────────────────────┐
│  4. Execute Request                 │
│     - app.Test(request)             │
│     - Routes through Fiber          │
│     - Calls actual controller       │
│     - Performs database operations  │
│     - Returns HTTP response         │
└─────────────────────────────────────┘
                ↓
┌─────────────────────────────────────┐
│  5. Verify Response                 │
│     - Check status code             │
│     - Parse response body           │
│     - Validate JSON structure       │
│     - Check error messages          │
└─────────────────────────────────────┘
                ↓
┌─────────────────────────────────────┐
│  6. Verify Database                 │
│     - Query database directly       │
│     - Check records created/updated │
│     - Verify soft deletes           │
│     - Validate constraints          │
└─────────────────────────────────────┘
                ↓
┌─────────────────────────────────────┐
│  7. Cleanup                         │
│     - Delete test data              │
│     - Close connections             │
└─────────────────────────────────────┘
```

### Helper Functions

#### `setupIntegrationDB(t *testing.T) *gorm.DB`
Creates a connection to the test database and runs migrations.

#### `setupIntegrationApp(db *gorm.DB) *fiber.App`
Creates a Fiber app with all admin routes registered.

#### `generateJWTTokenForAdmin(admin *model.User) (string, error)`
Generates a valid JWT token for authentication in tests.

#### `createTestAdminUser(t *testing.T, db *gorm.DB, email string) *model.User`
Creates an admin user in the database for testing.

#### `createJSONRequest(method, path string, body interface{}) *http.Request`
Creates an HTTP request with JSON body.

#### `cleanupIntegrationDB(t *testing.T, db *gorm.DB)`
Removes all test data from the database.

## Best Practices

### ✅ Do
- Run integration tests before committing major changes
- Use a separate test database (never production!)
- Clean up test data after each test
- Use descriptive test names
- Test both success and failure cases
- Verify database state, not just HTTP responses

### ❌ Don't
- Don't run integration tests on production database
- Don't skip test cleanup
- Don't make tests dependent on each other
- Don't hard-code test data that could conflict
- Don't commit database credentials

## Troubleshooting

### Tests are Skipped

```
--- SKIP: TestAreThereAnyAdminIntegration (0.00s)
    admin_integration_test.go:45: TEST_DATABASE_URL not set
```

**Solution:** Set the `TEST_DATABASE_URL` environment variable.

### Connection Failed

```
Failed to connect to test database: connection refused
```

**Solution:** Ensure PostgreSQL is running and accessible.

### Migration Failed

```
Failed to migrate test database: relation already exists
```

**Solution:** The database may already have tables. Drop and recreate the test database.

```sql
DROP DATABASE IF EXISTS auth_test;
CREATE DATABASE auth_test;
```

### Tests Fail with "already exists"

```
duplicate key value violates unique constraint
```

**Solution:** Previous test didn't clean up properly. Manually clean:

```sql
TRUNCATE TABLE users CASCADE;
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: testuser
          POSTGRES_PASSWORD: testpass
          POSTGRES_DB: auth_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Run Integration Tests
        env:
          TEST_DATABASE_URL: postgresql://testuser:testpass@localhost:5432/auth_test?sslmode=disable
        run: |
          cd auth
          go test -v -run Integration
```

## Performance

- **Unit Tests**: ~2s for 68 tests
- **Integration Tests**: ~5-10s for 30+ tests (depends on database)
- **Total**: ~7-12s for all tests

## Next Steps

- [ ] Add performance benchmarks
- [ ] Add concurrent request testing
- [ ] Add rate limiting tests
- [ ] Add role permission tests
- [ ] Add password reset flow tests
- [ ] Add email verification tests
- [ ] Add test for bulk operations
