# Mynute Go - Architecture Documentation

## Overview

Mynute Go is built as a microservices architecture with three independent services that communicate via HTTP APIs. This design provides scalability, independent deployment, and clear separation of concerns.

## Microservices

### 1. Core Service (Port 4000)

**Purpose**: Business logic and domain management

**Responsibilities**:
- Appointment scheduling and management
- Employee management and availability
- Client/customer management
- Service catalog and pricing
- Company and branch management
- Multi-tenant schema management
- File uploads (AWS S3)

**Key Technologies**:
- Fiber v2 web framework
- GORM for database operations
- PostgreSQL with multi-tenant schemas
- AWS S3 for file storage
- Admin panel (TypeScript/React)

**Database**: 
- Main DB: Business data with per-company schemas
- Policy storage: TenantPolicy, ClientPolicy, AdminPolicy models
- Auth DB: Shared authentication data (read-only access from Auth service)

**API Endpoints**: `/api/v1/*`

### 2. Auth Service (Port 4001)

**Purpose**: Authentication and authorization

**Responsibilities**:
- User authentication (JWT tokens)
- Admin user management
- Policy-based access control (PBAC)
- Policy evaluation and validation
- Resource and endpoint permissions
- Condition tree evaluation for authorization
- OAuth integration (future)

**Key Technologies**:
- Fiber v2 web framework
- GORM for database operations
- JWT for token management
- Policy-based access control

**Database**:
- Auth DB: Authentication data, users, roles, policies

**API Endpoints**: `/api/v1/auth/*`, `/api/v1/admin/*`, `/api/v1/policies/*`

### 3. Email Service (Port 4002)

**Purpose**: Email delivery and template management

**Responsibilities**:
- Send single emails
- Send template-based emails
- Send bulk emails
- Multi-provider support (Resend, MailHog)
- Email logging and tracking

**Key Technologies**:
- Fiber v2 web framework
- Resend API (production)
- MailHog (development/testing)
- Template rendering

**Database**: None (stateless service)

**API Endpoints**: `/api/v1/emails/*`

## Service Communication

### Communication Pattern

Services communicate via **HTTP REST APIs**. Each service exposes endpoints that other services can call.

```
┌─────────────┐     HTTP      ┌──────────────┐
│ Core Service│ ──────────────▶│ Auth Service │
│  (Port 4000)│                │  (Port 4001) │
└─────────────┘                └──────────────┘
       │                              
       │ HTTP                         
       │                              
       ▼                              
┌──────────────┐                      
│Email Service │                      
│  (Port 4002) │                      
└──────────────┘                      
```

### Example Flows

#### User Login Flow
1. Client → Core Service: POST `/api/v1/auth/login`
2. Core Service → Auth Service: POST `/api/v1/auth/verify`
3. Auth Service → Core Service: JWT token
4. Core Service → Client: Login response with token

#### Send Appointment Confirmation
1. Core Service: Create appointment
2. Core Service → Email Service: POST `/api/v1/emails/send-template`
3. Email Service → Resend API: Send email
4. Email Service → Core Service: Success response

## Database Architecture

### Multi-Database Setup

```
┌──────────────────┐
│  PostgreSQL      │
├──────────────────┤
│                  │
│  ┌────────────┐  │
│  │  Auth DB   │  │  ← Auth Service (R/W)
│  │            │  │  ← Core Service (R only)
│  └────────────┘  │
│                  │
│  ┌────────────┐  │
│  │  Main DB   │  │  ← Core Service (R/W)
│  │            │  │
│  │  • public  │  │  (Companies, system tables)
│  │  • company_│  │  (Per-tenant schemas)
│  │    {uuid}  │  │
│  └────────────┘  │
└──────────────────┘
```

### Schema Isolation

**Auth Database**: Single schema for all authentication data
- Tables: admins, admin_roles, resources, endpoints

**Policy Database**: Stored in Core service, evaluated by Auth service
- Tables: tenant_policies, client_policies, admin_policies
- Policy types separated by user context (tenant/client/admin)
- Each policy contains JSONB condition trees for authorization logic

**Main Database**: Multi-tenant with schema per company
- Public schema: companies, sectors, system data, policies (tenant_policies, client_policies, admin_policies), resources, endpoints
- Company schemas: company_{uuid} containing appointments, employees, clients, services, branches
- Policy models store condition trees as JSONB for flexible authorization rules

## Project Structure

```
mynute-go/
├── services/                    # Microservices directory
│   ├── core/                   # Core Service
│   │   ├── server.go           # Server initialization
│   │   ├── docs/               # Swagger docs (auto-generated)
│   │   ├── admin/              # Admin panel frontend
│   │   ├── docker-compose.*.yml
│   │   ├── Dockerfile
│   │   ├── api/                # Routes, controllers, DTOs
│   │   │   ├── middleware/     # Auth middleware (calls Auth service)
│   │   │   └── ...
│   │   ├── config/             # DB models, configs
│   │   │   └── db/
│   │   │       ├── model/      # Data models (policy.go, endpoint.go, resource.go)
│   │   │       └── seed/       # Seed data
│   │   │           ├── policy/ # Policy definitions by domain
│   │   │           │   ├── tenant_*.go  # Tenant-specific policies
│   │   │           │   ├── client_*.go  # Client-specific policies
│   │   │           │   └── helpers_*.go # Reusable condition checks
│   │   │           ├── endpoint/        # Endpoint definitions
│   │   │           └── resource/        # Resource definitions
│   │   └── lib/                # Utilities, email, auth
│   │
│   ├── auth/                   # Auth Service
│   │   ├── server.go           # Auth server
│   │   ├── docs/               # Swagger docs (auto-generated)
│   │   ├── docker-compose.*.yml
│   │   ├── Dockerfile
│   │   ├── api/                # Auth routes & controllers
│   │   │   ├── routes/
│   │   │   └── controller/
│   │   ├── config/             # DTOs, models
│   │   │   ├── dto/
│   │   │   └── db/
│   │   └── handler/            # JWT, auth logic
│   │
│   └── email/                  # Email Service
│       ├── server.go           # Email server
│       ├── routes.go           # Email routes
│       ├── docs/               # Swagger docs (auto-generated)
│       ├── docker-compose.*.yml
│       ├── Dockerfile
│       ├── controller/         # HTTP handlers
│       ├── dto/                # Request/response types
│       └── lib/                # Email providers
│
├── cmd/                        # Service entry points
│   ├── business-service/       # Core service main.go
│   ├── auth-service/           # Auth service main.go
│   ├── email-service/          # Email service main.go
│   └── seed/                   # Database seeding
│
├── migrations/                 # Database migrations
│   └── *.sql                   # Migration files
│
├── tools/                      # Development tools
│   ├── generate-migration/     # Migration generators
│   └── smart-migration/        # Smart migration tools
│
├── scripts/                    # Build & deployment scripts
│   ├── migrate.sh/ps1          # Migration scripts
│   └── seed.sh/ps1             # Seeding scripts
│
├── main.go                     # Multi-service launcher
├── Makefile                    # Build automation
└── go.mod                      # Go dependencies
```

## API Documentation

Each service maintains its own Swagger/OpenAPI documentation:

- **Auto-generated**: Run `make swagger-all` to regenerate
- **Location**: `services/{service}/docs/`
- **UI Access**: `http://localhost:{port}/swagger/index.html`

### Swagger Generation

```bash
# Generate all service docs
make swagger-all

# Individual services
make swagger-core    # Core service
make swagger-auth    # Auth service
make swagger-email   # Email service
```

## Configuration Management

### Service-Specific Configuration

Each service has its own `.env` file:

```
services/
├── core/.env          # Core service config
├── auth/.env          # Auth service config
└── email/.env         # Email service config
```

### Shared Configuration

Common settings in root `.env` (deprecated, being phased out):
- Database credentials
- Infrastructure configs

## Deployment Strategy

### Independent Deployment

Each service can be deployed independently:

```bash
# Deploy individual services
docker-compose -f services/core/docker-compose.prod.yml up -d
docker-compose -f services/auth/docker-compose.prod.yml up -d
docker-compose -f services/email/docker-compose.prod.yml up -d
```

### Service Discovery

Currently using **environment variables** for service URLs:
- Core Service: `http://localhost:4000`
- Auth Service: `http://localhost:4001`
- Email Service: `http://localhost:4002`

Future: Consider adding service discovery (Consul, etcd) for production.

## Scaling Considerations

### Horizontal Scaling

Services can be scaled independently:

**Stateless Services** (can scale freely):
- Email Service: No database, pure computation
- Auth Service: Stateless JWT validation

**Stateful Services** (require coordination):
- Core Service: Database connections, file uploads

### Load Balancing

Use reverse proxy (nginx, Traefik) to distribute traffic:

```
                    ┌──────────────┐
                    │ Load Balancer│
                    └──────┬───────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐       ┌─────▼────┐      ┌─────▼────┐
   │ Core #1 │       │ Core #2  │      │ Core #3  │
   └─────────┘       └──────────┘      └──────────┘
```

## Authorization Architecture

### Policy-Based Access Control (PBAC)

Mynute uses a sophisticated policy-based authorization system with clear separation of concerns:

#### Policy Types

**TenantPolicy**: Tenant/company user authorization
- Contains `tenant_id` for multi-tenant isolation
- Used for employees, managers, company operations
- Examples: company owners, branch managers, employees

**ClientPolicy**: Client/customer authorization
- No tenant association (clients access multiple companies)
- Used for appointment booking, profile management
- Examples: view own appointments, update profile

**AdminPolicy**: System administrator authorization
- Platform-wide access control
- Used for system management, super admin operations
- Examples: manage companies, system configuration

#### Policy Structure

Policies are stored as **JSONB condition trees**:

```json
{
  "description": "Allow company owner to update company",
  "logic_type": "AND",
  "children": [
    {
      "leaf": {
        "attribute": "subject.roles[*].id",
        "operator": "Contains",
        "value": "owner-role-uuid"
      }
    },
    {
      "leaf": {
        "attribute": "subject.company_id",
        "operator": "Equals",
        "resource_attribute": "resource.id"
      }
    }
  ]
}
```

#### Service Responsibilities

**Core Service** (`/core`):
- Stores policy data models (TenantPolicy, ClientPolicy, AdminPolicy)
- Defines policy seed data organized by domain
- Contains only basic field definitions
- No policy validation or evaluation logic

**Auth Service** (`/auth`):
- Evaluates policies during authorization checks
- Validates condition tree structure
- Executes condition logic (operators, comparisons)
- Returns allow/deny decisions

#### Policy Organization

Policies are organized by domain for maintainability:

```
services/core/config/db/seed/policy/
├── helpers_tenant.go      # Reusable tenant condition checks
├── helpers_client.go      # Reusable client condition checks
├── tenant_company.go      # Company-level policies
├── tenant_employee.go     # Employee management policies
├── tenant_branch.go       # Branch operation policies
├── tenant_service.go      # Service management policies
├── tenant_holiday.go      # Holiday management policies
├── tenant_appointment.go  # Tenant appointment policies
├── client_profile.go      # Client profile policies
├── client_appointment.go  # Client appointment policies
├── all_tenant.go          # Aggregates all tenant policies
└── all_client.go          # Aggregates all client policies
```

## Security

### Authentication Flow

1. User logs in via Core Service
2. Core Service validates credentials with Auth Service
3. Auth Service issues JWT token
4. User includes JWT in subsequent requests
5. Services validate JWT locally (stateless)

### Inter-Service Communication

**Current**: HTTP with no authentication (services trust each other)

**Future Enhancements**:
- Service-to-service authentication
- Mutual TLS (mTLS)
- API Gateway with service mesh

## Monitoring & Observability

### Health Checks

Each service exposes a health endpoint:
- Core: `GET /health`
- Auth: `GET /health`
- Email: `GET /health`

### Metrics

Prometheus integration available:
- Endpoint: `/metrics`
- Metrics: Request count, latency, errors

### Logging

Structured JSON logging with slog:
```json
{"time":"2024-11-05T10:30:00Z","level":"INFO","msg":"Request completed","path":"/api/v1/appointments","method":"GET","status":200}
```

## Development Workflow

### Running Locally

**With Docker (Recommended)**

Start all infrastructure services (databases, monitoring) with one command:

```bash
# Start all Docker services
go run cmd/docker-dev/main.go up

# Stop all Docker services
go run cmd/docker-dev/main.go down

# Restart all Docker services
go run cmd/docker-dev/main.go restart

# View logs from all services
go run cmd/docker-dev/main.go logs
```

This starts:
- **Core Service Infrastructure**: PostgreSQL (5432), Prometheus (9090), Grafana, Loki (3100), MinIO
- **Auth Service Infrastructure**: PostgreSQL (5433)
- **Email Service Infrastructure**: MailHog SMTP (1025), MailHog UI (8025)

**Run Application Services**

After starting infrastructure, run the application:

```bash
# Option 1: Run all services together
go run .

# Option 2: Run individually
go run ./cmd/business-service  # Core Service (Port 4000)
go run ./cmd/auth-service      # Auth Service (Port 4001)
go run ./cmd/email-service     # Email Service (Port 4002)
```

### Testing

```bash
# Run all tests
go test ./...

# Test specific service
go test ./services/core/...
go test ./services/auth/...
go test ./services/email/...
```

### Hot Reload (Development)

Use Air for automatic reload during development:
```bash
air
```

## Future Enhancements

### Short Term
- [ ] Add API Gateway (Kong, Traefik)
- [ ] Implement circuit breakers
- [ ] Add request tracing (OpenTelemetry)

### Medium Term
- [ ] Event-driven communication (NATS, RabbitMQ)
- [ ] Service mesh (Istio, Linkerd)
- [ ] Centralized logging (ELK, Loki)

### Long Term
- [ ] GraphQL federation
- [ ] gRPC for inter-service communication
- [ ] Kubernetes deployment
- [ ] Service discovery (Consul)

## Best Practices

### Adding a New Service

1. Create service directory in `services/`
2. Add `server.go` with Fiber setup
3. Create `cmd/{service-name}/main.go` entry point
4. Add Docker Compose files (dev/prod)
5. Add Swagger annotations to controllers
6. Generate Swagger docs: `swag init`
7. Add Makefile targets
8. Update this documentation

### Service Communication

- Use DTOs for all API requests/responses
- Handle errors gracefully with proper HTTP status codes
- Implement timeouts on all HTTP calls
- Log all inter-service communications
- Use structured error responses

### Database Access

- **Core Service**: 
  - Read/Write to Main DB (business data, policies)
  - Stores policy definitions (TenantPolicy, ClientPolicy, AdminPolicy)
  - No policy evaluation logic
  
- **Auth Service**: 
  - Reads policies from Core DB for evaluation
  - Validates and executes policy condition trees
  - Returns authorization decisions
  
- **Email Service**: No database access (stateless)

---

**Last Updated**: November 5, 2025
