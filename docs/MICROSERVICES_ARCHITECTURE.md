# Microservices Architecture

This project has been refactored into a microservices architecture with separate services for different concerns.

## Services

### 1. Core/Business Service
**Location:** `/core`  
**Port:** `4000`  
**Entry Point:** `cmd/business-service/main.go`

The main business logic service handling:
- Business operations
- Main database interactions
- Static file serving
- Admin panel
- API routes for business entities

**Run Locally:**
```bash
go run cmd/business-service/main.go
```

**Docker Development:**
```bash
docker-compose -p mynute-go -f core/docker-compose.dev.yml up -d
```

**Docker Production:**
```bash
docker-compose -f core/docker-compose.prod.yml up -d --build
```

**Configuration Files:**
- `core/Dockerfile` - Container configuration
- `core/docker-compose.dev.yml` - Development environment
- `core/docker-compose.prod.yml` - Production environment
- `core/server.go` - Server initialization and configuration
- `core/.env` - Environment variables

---

### 2. Auth Service
**Location:** `/auth`  
**Port:** `4001`  
**Entry Point:** `cmd/auth-service/main.go`

Dedicated authentication and authorization service handling:
- User authentication
- JWT token management
- Role-based access control (RBAC)
- Auth database management
- Endpoint/policy management

**Run Locally:**
```bash
go run cmd/auth-service/main.go
```

**Docker Development:**
```bash
docker-compose -p mynute-go-auth -f auth/docker-compose.dev.yml up -d
```

**Docker Production:**
```bash
docker-compose -f auth/docker-compose.prod.yml up -d --build
```

**Configuration Files:**
- `auth/Dockerfile` - Container configuration
- `auth/docker-compose.dev.yml` - Development environment
- `auth/docker-compose.prod.yml` - Production environment
- `auth/server.go` - Server initialization and configuration
- `auth/.env` - Environment variables

---

### 3. Email Service
**Location:** `/email`  
**Port:** `4002`  
**Entry Point:** `cmd/email-service/main.go`

Dedicated email microservice handling:
- Send single emails
- Send template-based emails
- Send bulk emails
- Email provider abstraction (SMTP, Resend, MailHog)
- Email queue management

**Run Locally:**
```bash
go run cmd/email-service/main.go
```

**Docker Development:**
```bash
docker-compose -p mynute-go-email -f email/docker-compose.dev.yml up -d
```

**Docker Production:**
```bash
docker-compose -f email/docker-compose.prod.yml up -d --build
```

**Configuration Files:**
- `email/Dockerfile` - Container configuration
- `email/docker-compose.dev.yml` - Development environment
- `email/docker-compose.prod.yml` - Production environment
- `email/server.go` - Server initialization and configuration
- `email/.env` - Environment variables

**API Endpoints:**
- `POST /api/v1/emails/send` - Send a single email
- `POST /api/v1/emails/send-template` - Send template email
- `POST /api/v1/emails/send-bulk` - Send bulk emails
- `GET /health` - Health check

---

## Service Communication

Each service runs independently and manages its own:
- Server configuration
- Routes and middleware
- Dependencies

**Ports:**
- Business Service: `4000`
- Auth Service: `4001`
- Email Service: `4002`

**Communication Pattern:**
Services can communicate via HTTP/REST APIs. The email service exposes endpoints that other services can call to send emails.

**Environment Variables:**

Each service has its own `.env` file:
- `core/.env` - Business service configuration
- `auth/.env` - Auth service configuration
- `email/.env` - Email service configuration

---

## Running Multiple Services

### Development

You have several options to run the services:

**Option 1: Run all services together (Recommended for local development)**
```bash
go run main.go
```
This will start all three services concurrently:
- Business Service (port 4000)
- Auth Service (port 4001)
- Email Service (port 4002)

Press `Ctrl+C` to gracefully shutdown all services.

**Option 2: Run services individually**
```bash
# Terminal 1 - Business Service
go run cmd/business-service/main.go

# Terminal 2 - Auth Service
go run cmd/auth-service/main.go

# Terminal 3 - Email Service
go run cmd/email-service/main.go
```

**Option 3: Using Docker Compose**
```bash
# Start all services with their dependencies
docker-compose -f core/docker-compose.dev.yml up -d
docker-compose -f auth/docker-compose.dev.yml up -d
docker-compose -f email/docker-compose.dev.yml up -d
```

### Production

Deploy each service independently based on your infrastructure needs:

```bash
# Deploy Business Service
cd core && docker-compose -f docker-compose.prod.yml up -d

# Deploy Auth Service
cd auth && docker-compose -f docker-compose.prod.yml up -d

# Deploy Email Service
cd email && docker-compose -f docker-compose.prod.yml up -d
```

---

## Migration from Monolith

The project was refactored from a monolithic structure where:
- All Docker files were in the root directory
- `main.go` started a single unified server
- All functionality was bundled together

**Changes Made:**

1. **Separated Docker Configurations:**
   - Moved `Dockerfile`, `docker-compose.dev.yml`, `docker-compose.prod.yml` to `/core`
   - Created new Docker files in `/auth` for the auth service
   - Created new Docker files in `/email` for the email service

2. **Created Independent Servers:**
   - `core/server.go` - Business service server
   - `auth/server.go` - Auth service server
   - `email/server.go` - Email service server

3. **Updated Entry Points:**
   - `cmd/business-service/main.go` - Business service entry point
   - `cmd/auth-service/main.go` - Auth service entry point
   - `cmd/email-service/main.go` - Email service entry point
   - `main.go` - Runs all services concurrently

4. **Separated Email Functionality:**
   - Moved email library from `core/src/lib/email` to `email/lib`
   - Moved email templates from `core/static/email` to `email/static/email`
   - Moved email translations from `core/translation/email` to `email/translation/email`
   - Created REST API for email operations

5. **Updated Documentation:**
   - `README.md` - Updated with new Docker commands
   - `docs/MICROSERVICES_ARCHITECTURE.md` - This architecture document
   - `docs/EMAIL_MIGRATION.md` - Email service migration guide
   - `email/README.md` - Email service documentation

---

## Benefits

✅ **Separation of Concerns:** Auth and email logic isolated from business logic  
✅ **Independent Scaling:** Scale services independently based on load  
✅ **Independent Deployment:** Deploy services without affecting others  
✅ **Better Organization:** Clear boundaries between service responsibilities  
✅ **Easier Testing:** Test services in isolation  
✅ **Flexible Infrastructure:** Deploy services on different servers/containers  
✅ **Email Provider Abstraction:** Switch email providers without affecting other services  
✅ **Reduced Coupling:** Services communicate via well-defined HTTP APIs

---

## Next Steps

Consider these enhancements for your microservices architecture:

- [ ] Add API Gateway for unified entry point
- [ ] Implement service-to-service authentication
- [ ] Add health check endpoints for each service
- [ ] Set up centralized logging
- [ ] Implement distributed tracing
- [ ] Add service discovery mechanism
- [ ] Create shared libraries for common functionality
- [ ] Set up CI/CD pipelines for each service
