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

### Tech Stack

- **Backend**: Go 1.23.4 with Fiber v2 web framework
- **Database**: PostgreSQL 17.5 with GORM ORM
- **Authentication**: JWT tokens with OAuth providers (Goth)
- **Email**: Resend for production, MailHog for testing
- **Cloud Storage**: AWS S3
- **Monitoring**: Prometheus
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose
- **Migrations**: golang-migrate

### Project Structure

```
mynute-go/
├── core/                     # Core application logic
│   ├── server.go            # Server initialization
│   └── src/
│       ├── config/          # Configuration (API routes, DB, cloud)
│       ├── controller/      # HTTP controllers
│       ├── handler/         # Business logic handlers
│       ├── lib/             # Utilities and libraries
│       ├── middleware/      # HTTP middleware
│       └── service/         # Business services
├── docs/                    # API documentation
├── migrations/              # Database migration files
├── static/                  # Static assets (email templates, pages)
├── test/                    # Test files
├── tools/                   # Development tools
└── scripts/                 # Build and deployment scripts
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

2. **Create environment file**
   ```bash
   cp .env.example .env
   ```

3. **Configure environment variables** in `.env`:
   ```env
   # Application
   APP_ENV=dev
   PORT=4000
   
   # Database
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5432
   POSTGRES_USER=your_user
   POSTGRES_PASSWORD=your_password
   POSTGRES_DB_PROD=mynute_prod
   POSTGRES_DB_DEV=mynute_dev
   POSTGRES_DB_TEST=mynute_test
   
   # Authentication
   JWT_SECRET=your-jwt-secret
   
   # Email (Resend)
   RESEND_API_KEY=your-resend-api-key
   RESEND_DEFAULT_FROM=noreply@yourdomain.com
   
   # AWS S3
   AWS_ACCESS_KEY_ID=your-access-key
   AWS_SECRET_ACCESS_KEY=your-secret-key
   AWS_REGION=us-east-1
   AWS_S3_BUCKET=your-bucket-name
   ```

### Development with Docker (Recommended)

1. **Start development environment**
   ```bash
   docker-compose -p mynute-go -f docker-compose.dev.yml up -d --force-recreate
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

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
Access the interactive API documentation at: `http://localhost:4000/swagger/`

### Key Endpoints

- **Authentication**: `/auth/*` - OAuth login/logout
- **Appointments**: `/appointments/*` - CRUD operations
- **Employees**: `/employees/*` - Employee management
- **Clients**: `/clients/*` - Client management
- **Services**: `/services/*` - Service configuration
- **Companies**: `/companies/*` - Multi-tenant company management
- **Branches**: `/branches/*` - Branch location management

## 🗄️ Database Management

### Important: Database Configuration

⚠️ **Migration tools use `POSTGRES_DB_PROD` environment variable to determine the target database.**

Set this explicitly in your `.env` file:
- Development: `POSTGRES_DB_PROD=devdb`
- Production: `POSTGRES_DB_PROD=maindb`

See [docs/MIGRATION_DATABASE_CONFIG.md](docs/MIGRATION_DATABASE_CONFIG.md) for details.

### Migrations

The project uses golang-migrate for database schema management:

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

See [docs/SEEDING.md](docs/SEEDING.md) for detailed production seeding guide.

### Smart Migration Tools

Generate migrations automatically based on model changes:

```bash
# Generate migration for specific models
make migrate-smart NAME=update_employee_table MODELS=Employee

# Generate comprehensive migration
make migrate-generate NAME=new_feature_migration
```

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific test package
go test ./core/src/service/...
```

### Test Database
The application automatically uses the test database (`POSTGRES_DB_TEST`) when `APP_ENV=test`.

## 🚀 Deployment

### Production with Docker

1. **Build production image**
   ```bash
   docker-compose -f docker-compose.prod.yml build
   ```

2. **Deploy with production configuration**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Manual Production Deployment

1. **Build binary**
   ```bash
   CGO_ENABLED=0 GOOS=linux go build -o mynute-backend-app
   ```

2. **Run migrations manually**
   ```bash
   make migrate-up
   ```

3. **Start application**
   ```bash
   ./mynute-backend-app
   ```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Application environment (dev/test/prod) | `dev` |
| `PORT` | Server port | `4000` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `JWT_SECRET` | JWT signing secret | Required |
| `RESEND_API_KEY` | Resend email API key | Required for email |
| `AWS_ACCESS_KEY_ID` | AWS access key | Required for S3 |

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

- [Database Migrations Guide](docs/MIGRATIONS.md)
- [Migration Workflow](docs/MIGRATION_WORKFLOW.md)
- [Email System Documentation](core/src/lib/email/README.md)
- [API Documentation](docs/swagger.yaml)

## 📄 License

This project is licensed under the Apache 2.0 License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the GitHub repository
- Contact: fiber@swagger.io
- Documentation: Available at `/swagger/` endpoint

## 🔗 Related Projects

- Frontend Application: [mynute-frontend](https://github.com/mynute-app/mynute-frontend)
- Mobile Application: [mynute-mobile](https://github.com/mynute-app/mynute-mobile)

---

**Mynute** - Simplifying appointment management for businesses worldwide.