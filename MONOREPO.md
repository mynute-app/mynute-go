# Monorepo Structure - Auth & Business Services

This project now uses a monorepo structure with two separate services that can be run independently.

## Architecture

```
mynute-go/
├── cmd/
│   ├── auth-service/        # Auth Service entry point (port 4001)
│   └── business-service/    # Business Service entry point (port 4000)
├── auth/                    # Auth logic & API
│   ├── api/
│   │   └── routes/         # Auth service routes
│   ├── model/              # Auth models (endpoints, policies, resources)
│   └── access_controller.go # Policy evaluation engine
├── core/                   # Business logic
│   └── src/
│       ├── controller/     # Business controllers
│       ├── config/
│       │   └── db/
│       │       ├── database.go  # Dual DB connections
│       │       ├── model/       # Business models
│       │       └── seed/        # Seed data
│       └── ...
└── go.mod                  # Single module
```

## Services

### Auth Service (Port 4001)
**Purpose:** Authentication, authorization, user management, policy evaluation

**Database:** `mynute_auth` (or `mynute_auth_dev`, `mynute_auth_test`)

**Contains:**
- User authentication (clients, employees, admins)
- Token generation and validation
- Policy and endpoint management
- Role management (system roles)

**Routes:**
- `POST /auth/login/*` - Login endpoints
- `POST /auth/validate` - Token validation (for business service)
- `GET/POST/PUT/DELETE /users/*` - User management
- `GET/POST/PUT/DELETE /policies/*` - Policy management
- `GET/POST/PUT/DELETE /endpoints/*` - Endpoint registry
- `GET/POST/PUT/DELETE /roles/*` - Role management

### Business Service (Port 4000)
**Purpose:** Core business operations (appointments, companies, branches, etc.)

**Database:** Your existing main database

**Contains:**
- Company management
- Branch operations
- Appointment scheduling
- Service management
- Employee profiles (business data)
- Client interactions

## Running the Services

### Development Mode

**Start Auth Service:**
```bash
# From project root
go run cmd/auth-service/main.go
```

**Start Business Service:**
```bash
# From project root
go run cmd/business-service/main.go
```

**Or run both with Make (if you create Makefile targets):**
```bash
make run-auth      # Runs auth service
make run-business  # Runs business service
make run-all       # Runs both in parallel
```

### Environment Variables

Add these to your `.env`:

```env
# General
APP_ENV=dev
APP_PORT=4000  # Business service default

# Auth Service
AUTH_SERVICE_PORT=4001

# Auth Database
POSTGRES_AUTH_HOST=localhost
POSTGRES_AUTH_PORT=5432
POSTGRES_AUTH_USER=postgres
POSTGRES_AUTH_PASSWORD=your_password
POSTGRES_AUTH_DB=mynute_auth_dev

# Main Database (existing)
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
POSTGRES_DB=mynute_main_dev
```

## Database Setup

Create the auth database:

```sql
-- Development
CREATE DATABASE mynute_auth_dev;

-- Test
CREATE DATABASE mynute_auth_test;

-- Production
CREATE DATABASE mynute_auth;
```

The auth service will automatically:
1. Connect to the auth database
2. Run migrations (in dev/test mode)
3. Seed initial data (endpoints, policies, roles)

## Communication Between Services

The business service can communicate with the auth service in two ways:

### Option 1: HTTP API Calls
Business service calls auth service's `/auth/validate` endpoint for each request.

### Option 2: JWT Validation (Recommended)
1. Auth service issues JWT tokens
2. Business service validates JWT signature locally (faster)
3. Business service only calls auth service for complex policy evaluations

## Next Steps

1. ✅ Create cmd/ structure
2. ✅ Create auth/api/routes
3. 🔄 Implement auth controllers
4. 🔄 Extract existing auth routes from core
5. 🔄 Update business service middleware
6. 🔄 Create Docker configurations
7. 🔄 Test inter-service communication

## Benefits

- **Independent Scaling**: Scale auth and business services separately
- **Security Isolation**: Auth credentials isolated in separate database
- **Clear Boundaries**: Auth vs business logic clearly separated
- **Parallel Development**: Teams can work independently
- **Shared Code**: Common utilities and types still shared via monorepo
- **Single Deployment Unit**: Can still deploy together if needed

## Migration from Current Setup

The business service currently uses `core.NewServer()` which maintains backward compatibility. The auth service is new and additive. Both can run simultaneously without breaking existing functionality.
