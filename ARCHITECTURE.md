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
- Auth DB: Shared authentication data (read-only access)

**API Endpoints**: `/api/v1/*`

### 2. Auth Service (Port 4001)

**Purpose**: Authentication and authorization

**Responsibilities**:
- User authentication (JWT tokens)
- Admin user management
- Access control policies (RBAC)
- Resource and endpoint permissions
- Policy evaluation
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
- Tables: admins, admin_roles, policies, resources, endpoints

**Main Database**: Multi-tenant with schema per company
- Public schema: companies, sectors, system data
- Company schemas: company_{uuid} containing appointments, employees, clients, services, branches

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
│   │   └── src/
│   │       ├── api/            # Routes, controllers, DTOs
│   │       ├── config/         # DB models, configs
│   │       ├── lib/            # Utilities, email, auth
│   │       ├── middleware/     # HTTP middleware
│   │       └── service/        # Business logic
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

- Core Service: Read/Write to Main DB, Read-only to Auth DB
- Auth Service: Read/Write to Auth DB only
- Email Service: No database access

---

**Last Updated**: November 5, 2025
