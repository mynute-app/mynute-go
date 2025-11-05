# Swagger Documentation Setup

## Current Status

The project has been refactored into microservices, but Swagger documentation is currently centralized in `/docs`. This document outlines the plan for service-specific Swagger documentation.

## Architecture Overview

### Services
- **Core Service** (port 4000): Business logic, companies, branches, employees, clients, appointments
- **Auth Service** (port 4001): Authentication and user management
- **Email Service** (port 4002): Email sending functionality

## Planned Swagger Setup

### 1. Service-Specific Documentation

Each service should have its own Swagger documentation:

```
core/
  docs/
    swagger.json
    swagger.yaml
    docs.go

auth/
  docs/
    swagger.json
    swagger.yaml
    docs.go

email/
  docs/
    swagger.json
    swagger.yaml
    docs.go
```

### 2. Generation Commands

Add to Makefile (or run directly):

```bash
# Core Service
swag init -g cmd/business-service/main.go -o core/docs --parseDependency --parseInternal

# Auth Service  
swag init -g cmd/auth-service/main.go -o auth/docs --parseDependency --parseInternal

# Email Service
swag init -g cmd/email-service/main.go -o email/docs --parseDependency --parseInternal
```

### 3. Swagger UI Access

Once configured, Swagger UI will be available at:

- Core: http://localhost:4000/swagger/index.html
- Auth: http://localhost:4001/swagger/index.html
- Email: http://localhost:4002/swagger/index.html

## Current Issues

### DTO Type References

Some Swagger annotations reference incorrect DTO types that need to be fixed:

1. ✅ **Fixed**: `auth/api/controller/admin.go` - Changed `DTO.AdminUpdate` to `DTO.AdminUpdateRequest`
2. ✅ **Fixed**: `auth/api/controller/client.go` - Changed `DTO.UpdateClient` to `DTO.UpdateClientRequest` (type created)
3. **Pending**: Cross-service DTO references may cause parsing errors

### Email Service Annotations

The email service currently lacks Swagger annotations. Need to add annotations to:

- `email/routes.go` - Main route handlers
- Add DTO types for email requests/responses

Example annotation needed:

```go
// handleSendTemplateEmail godoc
//
//	@Summary		Send template email
//	@Description	Send an email using a template
//	@Tags			Email
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SendTemplateRequest		true	"Email request"
//	@Success		200		{object}	SendEmailResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/v1/emails/send-template [post]
func handleSendTemplateEmail(c *fiber.Ctx) error {
    // ... implementation
}
```

## Recommended Approach

### Phase 1: Clean Separation (Recommended)
1. Each service generates its own Swagger docs independently
2. Only document the routes that service actually handles
3. Use service-specific DTO packages
4. Each service serves its own Swagger UI

### Phase 2: Unified API Gateway (Future)
1. Use an API Gateway (e.g., Traefik, Kong, Ambassador)
2. Aggregate Swagger docs from all services
3. Single entry point for all API documentation
4. Service discovery and routing

## Installation Requirements

```bash
# Install swaggo CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Install Fiber Swagger middleware
go get -u github.com/gofiber/swagger
```

## Integration with Services

Each service needs to import and serve Swagger docs:

```go
import (
    swagger "github.com/gofiber/swagger"
    _ "mynute-go/services/core/docs" // Import generated docs
)

func setupSwagger(app *fiber.App) {
    app.Get("/swagger/*", swagger.HandlerDefault)
}
```

## Next Steps

1. **Fix remaining DTO type references** across all controllers
2. **Add Swagger annotations** to email service routes
3. **Generate service-specific docs** using swag init
4. **Test each service** Swagger UI independently
5. **Update main.go** to show Swagger URLs for each service
6. **Document API** changes in service README files

## References

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [Fiber Swagger Middleware](https://github.com/gofiber/swagger)
- [OpenAPI 2.0 Specification](https://swagger.io/specification/v2/)
