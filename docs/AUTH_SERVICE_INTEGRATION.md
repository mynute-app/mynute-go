# Auth Service API Integration

## Overview

The core service now fetches all endpoint definitions from the auth service API instead of directly accessing the auth database. This enforces proper microservices architecture where the core service communicates with the auth service through HTTP APIs.

## Architecture

```
Core Service (Port 4000)
    ↓ HTTP GET
Auth Service API (Port 4001)
    ↓ Database Query
Auth Database (PostgreSQL)
```

## How It Works

1. **Core Service Startup**:
   - Creates an `AuthClient` to communicate with auth service
   - Checks if auth service is available at `http://localhost:4001/health`
   - Fetches endpoints from `GET http://localhost:4001/endpoints`

2. **Endpoint Registration**:
   - Receives array of `EndPoint` objects from auth service
   - Dynamically registers routes using the endpoint middleware
   - Each endpoint includes: `controller_name`, `method`, `path`, `deny_unauthorized`, `needs_company_id`

3. **Graceful Degradation**:
   - If auth service is unavailable, core service logs a warning
   - No routes are registered (test mode behavior)
   - In production, auth service should always be available

## Configuration

Set the auth service URL via environment variable (defaults to `http://localhost:4001`):

```bash
export AUTH_SERVICE_URL=http://localhost:4001
```

## Testing

### Prerequisites

1. **Start Auth Service Database**:
   ```bash
   docker-compose -p mynute-go-auth -f services/auth/docker-compose.dev.yml up -d
   ```

2. **Run Auth Service Migrations**:
   ```bash
   # From project root
   cd services/auth
   go run ../../cmd/migrate/main.go -action=up -path=../../migrations
   ```

3. **Seed Auth Service**:
   ```bash
   # From project root
   go run ./cmd/seed-auth
   ```
   This creates all endpoints, resources, policies, and roles.

4. **Start Auth Service**:
   ```bash
   go run ./cmd/auth-service
   ```
   Should be running on port 4001.

5. **Start Core Service Database**:
   ```bash
   docker-compose -p mynute-go-core -f services/core/docker-compose.dev.yml up -d
   ```

6. **Run Core Service Migrations**:
   ```bash
   cd services/core
   go run ../../cmd/migrate/main.go -action=up -path=../../migrations
   ```

7. **Start Core Service**:
   ```bash
   go run ./cmd/business-service
   ```
   Should be running on port 4000.
   You should see logs like:
   ```
   Successfully fetched 25 endpoints from auth service
   Registered route: POST /api/client -> CreateClient
   Registered route: GET /api/client/:id -> GetClientById
   ...
   Routes build finished! Registered 25 endpoints
   ```

### Run E2E Tests

```bash
cd services/core
go test ./test/e2e/...
```

Tests should now pass because routes are dynamically registered from the auth service.

## Troubleshooting

### "Auth service is not available"

- Make sure auth service is running on port 4001
- Check if auth service database is accessible
- Verify `AUTH_SERVICE_URL` environment variable

### "No routes registered"

- Auth service must be seeded with endpoint data
- Run `go run ./cmd/seed-auth` to create endpoints
- Check auth service logs for errors

### "Controller not found"

- Endpoint's `controller_name` doesn't match any registered controller
- Check that controllers are registered in `services/core/api/controller/*`
- Controller names must match exactly (e.g., "CreateClient")

## Implementation Details

### Auth Client (`services/core/api/lib/auth_client/`)

- `NewAuthClient()`: Creates HTTP client with 10s timeout
- `FetchEndpoints()`: GET `/endpoints` from auth service
- `IsAvailable()`: Checks `/health` endpoint

### Endpoint Middleware (`services/core/api/middleware/endpoint.go`)

- `Build()`: Original method (deprecated, queries auth DB directly)
- `BuildFromAPI()`: New method, accepts endpoints from API
- Builds middleware chain: Session → Authorization → Schema → Controller

### Routes (`services/core/api/routes/index.go`)

- Creates auth client
- Checks availability
- Fetches endpoints
- Registers routes dynamically

## Migration Path

**Before**: Core service directly connected to auth database
```go
Database.AuthDB = connectAuthDB()
routes.Build(db.Gorm, db.AuthDB, app)
```

**After**: Core service uses auth service API
```go
// No AuthDB connection
authClient := auth_client.NewAuthClient()
endpoints, _ := authClient.FetchEndpoints()
routes.Build(db.Gorm, app)
```

## Benefits

✅ Proper microservices separation
✅ Auth service owns its data completely  
✅ Core service can't accidentally corrupt auth data
✅ Easier to scale services independently
✅ Clear API contract between services
✅ Follows single responsibility principle
