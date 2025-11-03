# Auth Service Seeding

This document explains how to seed the auth service with endpoints, resources, and policies from the business service.

## Overview

The auth service needs to know about all endpoints, resources, and policies to perform authorization checks. This information is defined in the business service codebase under `core/src/config/db/seed/` and must be synced to the auth service database.

## Seeding Process

### 1. Start the Auth Service

First, ensure the auth service is running:

```powershell
go run cmd/auth-service/main.go
```

The auth service should be accessible at `http://localhost:4001` (or your configured `AUTH_SERVICE_URL`).

### 2. Run the Seed Command

From the project root, run:

```powershell
go run cmd/seed-auth/main.go
```

This will:
1. Load all resource definitions from `core/src/config/db/seed/resource/`
2. Load all endpoint definitions from `core/src/config/db/seed/endpoint/`
3. Send POST requests to the auth service to create each endpoint
4. Handle duplicates gracefully (409 Conflict responses are ignored)

### 3. Verify Seeding

Check the auth service database:

```sql
-- Count endpoints
SELECT COUNT(*) FROM endpoints;

-- List all endpoints
SELECT method, path, description FROM endpoints ORDER BY path;

-- Check for specific endpoints
SELECT * FROM endpoints WHERE path LIKE '%/admin%';
```

## Seed Data Sources

### Resources
Defined in: `core/src/config/db/seed/resource/resource.go`

```go
var Resources = []*authModel.Resource{
    Appointment,
    Branch,
    Client,
    Company,
    Employee,
    Holiday,
    Role,
    Sector,
    Service,
}
```

### Endpoints
Defined in: `core/src/config/db/seed/endpoint/*.go`

Each file contains endpoint definitions for a specific controller:
- `admin.go` - Admin management endpoints
- `appointment.go` - Appointment endpoints  
- `auth.go` - Authentication endpoints
- `branch.go` - Branch management endpoints
- `client.go` - Client management endpoints
- `company.go` - Company management endpoints
- `employee.go` - Employee management endpoints
- `holiday.go` - Holiday management endpoints
- `sector.go` - Sector management endpoints
- `service.go` - Service management endpoints

### Policies
Defined in: `core/src/config/db/seed/policy/*.go`

Policies are more complex and should be reviewed before seeding:
- `appointment.go` - Appointment access policies
- `branch.go` - Branch access policies
- `client.go` - Client access policies
- `company.go` - Company access policies
- `employee.go` - Employee access policies
- `helpers.go` - Reusable policy conditions
- `holiday.go` - Holiday access policies
- `service.go` - Service access policies

**Note:** Policy seeding is not yet automated. Review each policy definition carefully and create them via the auth service admin panel or API.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `AUTH_SERVICE_URL` | `http://localhost:4001` | Auth service base URL |

## Endpoint Format

When creating endpoints via API, use this format:

```json
{
  "method": "GET",
  "path": "/api/users/:id",
  "description": "Get user by ID",
  "controller_name": "GetUserById",
  "deny_unauthorized": true,
  "resource_id": "uuid-here"
}
```

## Troubleshooting

### Auth Service Not Responding

**Error:** `failed to send request: dial tcp ... connection refused`

**Solution:** Ensure the auth service is running on the expected port:
```powershell
# Check if port 4001 is listening
netstat -an | findstr ":4001"

# Start the auth service
go run cmd/auth-service/main.go
```

### Duplicate Endpoint Errors

**Error:** `unexpected status 409`

**Solution:** This is expected and handled gracefully. The seeder ignores 409 Conflict responses, which indicate the endpoint already exists.

### Missing Resource IDs

**Error:** `resource_id not found`

**Solution:** Resources must be seeded before endpoints that reference them. The current seeder skips resource creation. To manually create resources, use the auth service API or database migration.

## Manual Seeding

To manually seed a single endpoint:

```powershell
$body = @{
    method = "GET"
    path = "/api/test"
    description = "Test endpoint"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:4001/endpoints" -Method Post -Body $body -ContentType "application/json"
```

## Next Steps

1. **Review Policies** - Examine policy definitions in `core/src/config/db/seed/policy/`
2. **Create Policies** - Use auth service admin panel to create policies
3. **Test Authorization** - Use `/authorize/by-method-and-path` to test access checks
4. **Add New Endpoints** - When adding new endpoints to business service, re-run the seeder

## See Also

- [Auth Service Implementation](AUTH_SERVICE_IMPLEMENTATION.md)
- [Auth Service Migration](AUTH_SERVICE_MIGRATION.md)
- [Seeding Documentation](SEEDING.md)
