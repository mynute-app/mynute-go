# Mynute Go - Development Guide for Claude

## Project Overview

Mynute Go is a microservices-based appointment management system built with Go, Fiber v2, and PostgreSQL. The project uses a **Go workspace** setup with three independent microservices that communicate via HTTP REST APIs.

## Architecture

### Three Microservices

1. **Core Service** (Port 4000)
   - Business logic: appointments, employees, clients, services, companies, branches
   - Multi-tenant architecture with per-company database schemas
   - Stores policy definitions (but doesn't evaluate them)
   - AWS S3 integration for file uploads
   - Admin panel frontend (TypeScript/React)

2. **Auth Service** (Port 4001)
   - Authentication (JWT tokens)
   - Authorization (policy evaluation engine)
   - Admin user management
   - Reads and evaluates policies from Core DB
   - Policy-based access control with condition trees

3. **Email Service** (Port 4002)
   - Email delivery (Resend for production, MailHog for dev)
   - Template-based emails
   - Stateless service (no database)

### Key Architectural Decisions

- **Policy Storage vs Evaluation**: Core service stores policy models (TenantPolicy, ClientPolicy, AdminPolicy) as JSONB condition trees. Auth service reads and evaluates these policies during authorization checks.
- **Service Independence**: Each service has its own `go.mod` for true dependency isolation
- **Multi-tenant Isolation**: Each company gets its own PostgreSQL schema (`company_{uuid}`)
- **No Cross-Service Imports**: Services communicate via HTTP APIs, never by importing each other's code

## Quick Start

### Prerequisites
- Go 1.23.4 or higher
- PostgreSQL 17.5 or higher
- Docker & Docker Compose (optional but recommended)

### Environment Setup

Each service has its own `.env` file:

```bash
# Core Service
cp services/core/.env.example services/core/.env

# Auth Service
cp services/auth/.env.example services/auth/.env

# Email Service
cp services/email/.env.example services/email/.env
```

### Running with Docker (Recommended)

Start all infrastructure (databases, monitoring, etc.):

```bash
# Start all Docker services (PostgreSQL, Prometheus, MailHog, etc.)
go run cmd/docker-dev/main.go up

# View logs
go run cmd/docker-dev/main.go logs

# Stop all services
go run cmd/docker-dev/main.go down
```

Then run the application services:

```bash
# Option 1: Run all services together
go run .

# Option 2: Run services individually
go run ./cmd/business-service  # Core Service
go run ./cmd/auth-service       # Auth Service
go run ./cmd/email-service      # Email Service
```

### Running Without Docker

1. **Setup PostgreSQL databases manually**
   - Create `mynute_prod` database for Core service
   - Create `mynute_auth` database for Auth service
   - Update connection strings in each service's `.env`

2. **Run migrations**
   ```bash
   make migrate-up
   ```

3. **Run services**
   ```bash
   go run .
   ```

## Testing

### Run Tests

```bash
# Test all services
go test ./...

# Test specific service
go test ./services/core/...
go test ./services/auth/...
go test ./services/email/...

# Test with coverage
go test -v -cover ./...

# Test specific package
go test ./services/core/api/controller/...
```

### Integration Tests

```bash
# Auth service integration tests (requires running database)
cd services/auth
go test ./test/e2e/...
```

### Test Database

The application automatically uses test database when `APP_ENV=test` in `.env` files.

## Database Management

### Migrations

```bash
# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create NAME=add_new_feature

# Check migration status
make migrate-version
```

### Seeding

```bash
# Seed system data (roles, policies, resources, endpoints)
make seed

# Build seed binary for production
make seed-build
```

### Smart Migration Tools

```bash
# Generate migration for specific models
make migrate-smart NAME=update_employee_table MODELS=Employee

# Generate comprehensive migration
make migrate-generate NAME=new_feature_migration
```

## Development Workflow

### Project Structure

```
mynute-go/
├── go.work                    # Workspace configuration
├── go.mod                     # Root module (launcher)
├── main.go                    # Multi-service launcher
├── services/
│   ├── core/                  # Core Service
│   │   ├── go.mod            # Independent module
│   │   ├── server.go
│   │   ├── api/
│   │   ├── config/
│   │   │   └── db/
│   │   │       ├── model/    # Data models
│   │   │       │   ├── policy.go     # TenantPolicy, ClientPolicy, AdminPolicy
│   │   │       │   ├── endpoint.go   # EndPoint model
│   │   │       │   └── resource.go   # Resource model
│   │   │       └── seed/
│   │   │           ├── policy/       # Policy seed data by domain
│   │   │           │   ├── tenant_company.go
│   │   │           │   ├── tenant_employee.go
│   │   │           │   ├── client_profile.go
│   │   │           │   ├── helpers_tenant.go
│   │   │           │   └── helpers_client.go
│   │   │           ├── endpoint/     # Endpoint definitions
│   │   │           └── resource/     # Resource definitions
│   │   └── admin/            # Admin panel frontend
│   ├── auth/                  # Auth Service
│   │   ├── go.mod            # Independent module
│   │   ├── server.go
│   │   └── handler/          # Policy evaluation logic
│   └── email/                 # Email Service
│       ├── go.mod            # Independent module
│       └── server.go
├── cmd/                       # Service entry points
├── migrations/                # SQL migrations
├── scripts/                   # Build scripts
└── tools/                     # Dev tools
```

### Adding Dependencies

```bash
# Add to specific service
cd services/core
go get github.com/some/package
go mod tidy

# Sync workspace
cd ../..
go work sync
```

### API Documentation

Each service has Swagger documentation:

- Core: http://localhost:4000/swagger/index.html
- Auth: http://localhost:4001/swagger/index.html
- Email: http://localhost:4002/swagger/index.html

Regenerate after changes:

```bash
# Regenerate all Swagger docs
make swagger-all

# Or individually
make swagger-core
make swagger-auth
make swagger-email
```

## Common Development Tasks

### Adding a New API Endpoint

1. Define endpoint in `services/core/config/db/seed/endpoint/`
2. Create controller in `services/core/api/controller/`
3. Add route in `services/core/api/routes/`
4. Add Swagger annotations to controller
5. Regenerate Swagger: `make swagger-core`

### Adding a New Policy

1. Define condition tree in `services/core/config/db/seed/policy/`
2. Organize by domain (tenant_*.go or client_*.go)
3. Add to aggregation function (all_tenant.go or all_client.go)
4. Seed to database: `make seed`

### Adding a New Database Model

1. Create model in `services/core/config/db/model/`
2. Add to model registry in `index.go`
3. Generate migration: `make migrate-generate NAME=add_new_model`
4. Run migration: `make migrate-up`

## Important Files

### Configuration
- `services/core/.env` - Core service config
- `services/auth/.env` - Auth service config
- `services/email/.env` - Email service config

### Documentation
- `ARCHITECTURE.md` - Detailed architecture overview
- `README.md` - Project readme with setup instructions
- `GO_WORKSPACE.md` - Go workspace setup explanation

### Database
- `migrations/*.sql` - Database migrations
- `services/core/config/db/model/` - GORM models
- `services/core/config/db/seed/` - Seed data

### Policy System
- `services/core/config/db/model/policy.go` - Policy data models
- `services/core/config/db/seed/policy/` - Policy definitions
- Auth service evaluates policies (reads from Core DB)

## Debugging

### View Service Logs

```bash
# All services
go run cmd/docker-dev/main.go logs

# Specific service
docker-compose -f services/core/docker-compose.dev.yml logs -f
```

### Check Database

```bash
# Connect to Core DB
psql -h localhost -p 5432 -U your_user -d mynute_prod

# Connect to Auth DB
psql -h localhost -p 5433 -U your_user -d mynute_auth

# List schemas
\dn

# List tables
\dt

# List company schemas
\dn company_*
```

### Common Issues

1. **Port already in use**: Stop other services or change ports in `.env`
2. **Database connection failed**: Check Docker is running and databases are up
3. **Migration failed**: Check migration syntax, rollback if needed
4. **Swagger not updating**: Run `make swagger-all` to regenerate

## Build & Deploy

### Build Binaries

```bash
# Build all services
go build ./cmd/...

# Build specific service
go build -o bin/mynute-core ./cmd/business-service
go build -o bin/mynute-auth ./cmd/auth-service
go build -o bin/mynute-email ./cmd/email-service
```

### Docker Production

```bash
# Build production images
docker-compose -f services/core/docker-compose.prod.yml build
docker-compose -f services/auth/docker-compose.prod.yml build
docker-compose -f services/email/docker-compose.prod.yml build

# Deploy
docker-compose -f services/core/docker-compose.prod.yml up -d
docker-compose -f services/auth/docker-compose.prod.yml up -d
docker-compose -f services/email/docker-compose.prod.yml up -d
```

## Monitoring

- **Prometheus**: Each service has metrics at `/metrics`
- **Loki**: Each service has its own Loki config (ports 3100, 3101, 3102)
- **Health Checks**: Each service has `/health` endpoint

## Key Conventions

1. **Service Boundaries**: Never import code between services, use HTTP APIs
2. **Policy Architecture**: Core stores, Auth evaluates
3. **Database Schemas**: Use `company_{uuid}` for tenant isolation
4. **Error Handling**: Return structured JSON errors with proper HTTP status codes
5. **Testing**: Write tests for new features, maintain coverage
6. **Documentation**: Update Swagger annotations for API changes
7. **Migrations**: Always create migrations for schema changes
8. **Commits**: Use conventional commit format (feat:, fix:, refactor:, docs:, etc.)

## Resources

- **Fiber Documentation**: https://docs.gofiber.io/
- **GORM Documentation**: https://gorm.io/docs/
- **Swagger/OpenAPI**: https://swagger.io/specification/
- **Go Workspaces**: https://go.dev/doc/tutorial/workspaces

---

**Last Updated**: November 15, 2025
