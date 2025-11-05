# Mynute Go - Appointment Management System

A modern, scalable appointment management system built with Go, Fiber web framework, and PostgreSQL. Mynute enables businesses to manage appointments, employees, clients, and services with multi-tenant architecture support.

## 🚀 Features

- **Multi-tenant Architecture**: Support for multiple companies with isolated data schemas
- **Appointment Management**: Full CRUD operations for appointments with scheduling capabilities
- **Employee Management**: Handle employee profiles, roles, and availability
- **Client Management**: Manage client information and appointment history  
- **Service Management**: Define and organize business services
- **Branch Management**: Support for multiple business locations
- **Authentication**: OAuth integration with popular providers
- **Email System**: Multi-language email templates with multiple provider support (Resend, MailHog)
- **Database Migrations**: Version-controlled schema management with golang-migrate
- **API Documentation**: Comprehensive Swagger/OpenAPI documentation
- **Monitoring**: Prometheus metrics integration
- **Cloud Storage**: AWS S3 integration for file uploads
- **Docker Support**: Full containerization with development and production configurations

## 🏗️ Architecture

### Microservices Architecture

Mynute Go follows a microservices architecture with independent, deployable services:

- **Core Service** (Port 4000): Business logic, appointments, employees, clients, companies
- **Auth Service** (Port 4001): Authentication, authorization, access control policies
- **Email Service** (Port 4002): Email sending, templates, bulk emails

Each service has:
- Independent deployment with Docker
- Service-specific Swagger documentation at `/swagger/index.html`
- Isolated configuration and environment variables
- Own database connection management

### Tech Stack

- **Backend**: Go 1.23.4 with Fiber v2 web framework
- **Database**: PostgreSQL 17.5 with GORM ORM
- **Authentication**: JWT tokens with OAuth providers (Goth)
- **Email**: Resend for production, MailHog for testing
- **Cloud Storage**: AWS S3
- **Monitoring**: Prometheus
- **Documentation**: Swagger/OpenAPI 2.0
- **Containerization**: Docker & Docker Compose
- **Migrations**: golang-migrate

### Project Structure

```
mynute-go/
├── services/                  # Microservices
│   ├── core/                 # Core/Business Service (Port 4000)
│   │   ├── server.go         # Server initialization
│   │   ├── docs/             # Swagger documentation
│   │   ├── admin/            # Admin panel frontend
│   │   ├── loki-config.yaml  # Loki logging config
│   │   ├── prometheus.yml    # Prometheus metrics config
│   │   ├── docker-compose.*.yml
│   │   ├── Dockerfile
│   │   └── src/
│   │       ├── api/          # API routes, controllers, DTOs
│   │       ├── config/       # Database models, configs
│   │       ├── lib/          # Shared utilities
│   │       ├── middleware/   # HTTP middleware
│   │       └── service/      # Business logic
│   ├── auth/                 # Auth Service (Port 4001)
│   │   ├── server.go         # Auth server initialization
│   │   ├── docs/             # Swagger documentation
│   │   ├── loki-config.yaml  # Loki logging config
│   │   ├── prometheus.yml    # Prometheus metrics config
│   │   ├── docker-compose.*.yml
│   │   ├── Dockerfile
│   │   ├── api/              # Auth API routes & controllers
│   │   ├── config/           # Auth DTOs and models
│   │   └── handler/          # JWT, auth logic, access control
│   └── email/                # Email Service (Port 4002)
│       ├── server.go         # Email server initialization
│       ├── docs/             # Swagger documentation
│       ├── loki-config.yaml  # Loki logging config
│       ├── prometheus.yml    # Prometheus metrics config
│       ├── docker-compose.*.yml
│       ├── Dockerfile
│       ├── controller/       # Email HTTP handlers
│       ├── dto/              # Email request/response types
│       ├── lib/              # Email providers (Resend, MailHog)
│       └── routes.go         # Email API routes
├── cmd/                      # Service entry points
│   ├── business-service/     # Core service main.go
│   ├── auth-service/         # Auth service main.go
│   ├── email-service/        # Email service main.go
│   ├── migrate/              # Migration tool
│   ├── job/                  # Job scripts (e.g., create random companies)
│   ├── seed/                 # Database seeding tools
│   ├── seed-admin/           # Admin seeding
│   └── seed-auth/            # Auth seeding
├── migrations/               # Database migration files
├── scripts/                  # Build and deployment scripts
├── tools/                    # Development tools
│   ├── generate-migration/   # Migration generators
│   └── smart-migration/      # Smart migration tools
├── prometheus-all-services.yml  # Unified Prometheus config
└── main.go                   # Multi-service launcher
```

## 🛠️ Quick Start

### Prerequisites

- Go 1.23.4 or higher
- PostgreSQL 17.5 or higher
- Docker & Docker Compose (optional)

### Environment Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/mynute-app/mynute-go.git
   cd mynute-go
   ```

2. **Create environment files for each service**
   
   Each service manages its own configuration:
   
   ```bash
   # Core/Business Service
   cp services/core/.env.example services/core/.env
   
   # Auth Service
   cp services/auth/.env.example services/auth/.env
   
   # Email Service
   cp services/email/.env.example services/email/.env
   ```

3. **Configure environment variables** in each service's `.env`:
   
   Core Service (`services/core/.env`):
   ```env
   # Application
   APP_ENV=dev
   APP_PORT=4000
   
   # Database
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5432
   POSTGRES_USER=your_user
   POSTGRES_PASSWORD=your_password
   POSTGRES_DB=mynute_prod
   
   # AWS S3
   AWS_ACCESS_KEY_ID=your-access-key
   AWS_SECRET_ACCESS_KEY=your-secret-key
   AWS_REGION=us-east-1
   AWS_S3_BUCKET=your-bucket-name
   ```
   
   Auth Service (`services/auth/.env`):
   ```env
   # Application
   APP_ENV=dev
   AUTH_SERVICE_PORT=4001
   
   # Database
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5432
   POSTGRES_USER=your_user
   POSTGRES_PASSWORD=your_password
   POSTGRES_AUTH_DB=mynute_auth
   
   # JWT
   JWT_SECRET=your-jwt-secret
   JWT_EXPIRATION=24h
   ```
   
   Email Service (`services/email/.env`):
   ```env
   # Application
   APP_ENV=dev
   EMAIL_SERVICE_PORT=4002
   
   # Resend (Production)
   RESEND_API_KEY=your-resend-api-key
   RESEND_DEFAULT_FROM=noreply@yourdomain.com
   
   # MailHog (Development)
   MAILHOG_HOST=localhost
   MAILHOG_PORT=1025
   ```

### Development with Docker (Recommended)

1. **Start development environment**
   
   Each service has its own Docker Compose configuration:
   
   ```bash
   # Core/Business Service (Port 4000)
   docker-compose -p mynute-core -f services/core/docker-compose.dev.yml up -d
   
   # Auth Service (Port 4001)
   docker-compose -p mynute-auth -f services/auth/docker-compose.dev.yml up -d
   
   # Email Service (Port 4002)
   docker-compose -p mynute-email -f services/email/docker-compose.dev.yml up -d
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the application**
   
   **Option A: Run all services together (Recommended)**
   ```bash
   go run main.go
   ```
   
   **Option B: Run services individually**
   
   ```bash
   # Core/Business Service
   go run cmd/business-service/main.go
   
   # Auth Service
   go run cmd/auth-service/main.go
   
   # Email Service
   go run cmd/email-service/main.go
   ```

### Access Points

Once running, access each service:

- **Core Service**: http://localhost:4000
  - Swagger UI: http://localhost:4000/swagger/index.html
  - Admin Panel: http://localhost:4000/admin
  
- **Auth Service**: http://localhost:4001
  - Swagger UI: http://localhost:4001/swagger/index.html
  
- **Email Service**: http://localhost:4002
  - Swagger UI: http://localhost:4002/swagger/index.html

### Manual Setup

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Setup PostgreSQL database**
   - Create databases for development and testing
   - Update connection strings in `.env`

3. **Run database migrations**
   ```bash
   make migrate-up
   ```

4. **Start the application**
   ```bash
   go run main.go
   ```

The application will be available at `http://localhost:4000`

## 📚 API Documentation

### Swagger UI

Each microservice has its own interactive API documentation:

- **Core Service**: http://localhost:4000/swagger/index.html
- **Auth Service**: http://localhost:4001/swagger/index.html
- **Email Service**: http://localhost:4002/swagger/index.html

### Regenerate Swagger Documentation

```bash
# Generate docs for all services
make swagger-all

# Or generate for specific service
make swagger-core    # Core service
make swagger-auth    # Auth service
make swagger-email   # Email service
```

### Key API Endpoints

**Core Service (Port 4000)**
- `/api/v1/appointments/*` - Appointment management
- `/api/v1/employees/*` - Employee management
- `/api/v1/clients/*` - Client management
- `/api/v1/services/*` - Service configuration
- `/api/v1/companies/*` - Company management
- `/api/v1/branches/*` - Branch management

**Auth Service (Port 4001)**
- `/api/v1/auth/*` - Authentication & login
- `/api/v1/admin/*` - Admin management
- `/api/v1/policies/*` - Access control policies
- `/api/v1/resources/*` - Resource management
- `/api/v1/endpoints/*` - Endpoint permissions

**Email Service (Port 4002)**
- `/api/v1/emails/send` - Send email(s) to single or multiple recipients
- `/api/v1/emails/send-template-merge` - Send email using template HTML from Core service
- `/health` - Health check

## 🗄️ Database Management

### Migrations

The project uses a custom migration tool for database schema management:

```bash
# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create NAME=add_new_feature

# Check migration status
make migrate-version

# Rollback N migrations
make migrate-down-n STEPS=3
```

**Manual Migration (if needed):**
```bash
go run cmd/migrate/main.go -action=up -path=./migrations
```

### Seeding

Seed system data (roles, resources, endpoints, policies):

```bash
# Run seeding in development
make seed

# Build seed binary for production
make seed-build

# Show seeding help
make seed-help
```

**Production Seeding:**
```bash
# 1. Build the binary
make seed-build

# 2. Deploy to server and run
./bin/seed
```

### Smart Migration Tools

Generate migrations automatically based on model changes:

```bash
# Generate migration for specific models
make migrate-smart NAME=update_employee_table MODELS=Employee

# Generate comprehensive migration
make migrate-generate NAME=new_feature_migration
```

## 🔍 Monitoring & Observability

### Prometheus Metrics

Each service has its own Prometheus configuration:

- **Core Service**: `services/core/prometheus.yml`
- **Auth Service**: `services/auth/prometheus.yml`
- **Email Service**: `services/email/prometheus.yml`

**Unified Monitoring** (all services):
```bash
prometheus --config.file=prometheus-all-services.yml
```

### Loki Logging

Each service has isolated logging configuration:

- **Core Service**: `services/core/loki-config.yaml` - Port 3100
- **Auth Service**: `services/auth/loki-config.yaml` - Port 3101
- **Email Service**: `services/email/loki-config.yaml` - Port 3102

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific service tests
go test ./services/core/...
go test ./services/auth/...
go test ./services/email/...
```

### Test Database
The application automatically uses the test database when `APP_ENV=test`.

## 🚀 Deployment

### Production with Docker

Each service can be deployed independently:

1. **Build production images**
   
   ```bash
   # Core/Business Service
   docker-compose -f services/core/docker-compose.prod.yml build
   
   # Auth Service
   docker-compose -f services/auth/docker-compose.prod.yml build
   
   # Email Service
   docker-compose -f services/email/docker-compose.prod.yml build
   ```

2. **Deploy with production configuration**
   
   ```bash
   # Core/Business Service
   docker-compose -p mynute-core-prod -f services/core/docker-compose.prod.yml up -d
   
   # Auth Service
   docker-compose -p mynute-auth-prod -f services/auth/docker-compose.prod.yml up -d
   
   # Email Service
   docker-compose -p mynute-email-prod -f services/email/docker-compose.prod.yml up -d
   ```

### Manual Production Deployment

1. **Build service binaries**
   
   ```bash
   # Core/Business Service
   CGO_ENABLED=0 GOOS=linux go build -o bin/mynute-core ./cmd/business-service
   
   # Auth Service
   CGO_ENABLED=0 GOOS=linux go build -o bin/mynute-auth ./cmd/auth-service
   
   # Email Service
   CGO_ENABLED=0 GOOS=linux go build -o bin/mynute-email ./cmd/email-service
   ```

2. **Run migrations manually** (Core and Auth services only)
   ```bash
   make migrate-up
   ```

3. **Start services**
   ```bash
   # Start each service on its designated port
   ./bin/mynute-core   # Port 4000
   ./bin/mynute-auth   # Port 4001
   ./bin/mynute-email  # Port 4002
   ```

## 🔧 Configuration

### Environment Variables

**Core Service**

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Application environment (dev/test/prod) | `dev` |
| `APP_PORT` | Server port | `4000` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `POSTGRES_DB` | Main database name | Required |
| `AWS_ACCESS_KEY_ID` | AWS access key | Required for S3 |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | Required for S3 |

**Auth Service**

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Application environment | `dev` |
| `AUTH_SERVICE_PORT` | Auth service port | `4001` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_AUTH_DB` | Auth database name | Required |
| `JWT_SECRET` | JWT signing secret | Required |
| `JWT_EXPIRATION` | Token expiration time | `24h` |

**Email Service**

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Application environment | `dev` |
| `EMAIL_SERVICE_PORT` | Email service port | `4002` |
| `RESEND_API_KEY` | Resend email API key | Required (prod) |
| `RESEND_DEFAULT_FROM` | Default sender email | Required |
| `MAILHOG_HOST` | MailHog host (dev) | `localhost` |
| `MAILHOG_PORT` | MailHog port (dev) | `1025` |

### Multi-tenant Configuration

Each company gets its own database schema for data isolation:
- Public schema: `companies`, `users`, `system_data`
- Company schema: `company_{uuid}` containing appointments, employees, clients, etc.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Write tests for new features
- Update documentation for API changes
- Run migrations for database schema changes
- Use semantic commit messages

## 📖 Additional Documentation

- **Architecture**:
  - [Architecture Overview](ARCHITECTURE.md)

- **Microservices**:
  - [Email Service Documentation](services/email/README.md)
  
- **Development**:
  - [Admin Panel Documentation](services/core/admin/README.md)
  - [Integration Tests](services/auth/api/controller/INTEGRATION_TESTS.md)

## 📄 License

This project is licensed under the Apache 2.0 License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the GitHub repository
- Email: support@mynute.com
- Documentation: Available at each service's `/swagger/index.html` endpoint

## 🔗 Related Projects

- Frontend Application: [mynute-frontend](https://github.com/mynute-app/mynute-frontend)
- Mobile Application: [mynute-mobile](https://github.com/mynute-app/mynute-mobile)

---

**Mynute** - Simplifying appointment management for businesses worldwide.