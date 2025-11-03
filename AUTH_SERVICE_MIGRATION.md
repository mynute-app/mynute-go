# Auth Service Migration Guide

This document explains how to migrate authentication and user management from the monolithic business service to the separate auth service.

## Architecture Overview

### Before (Monolithic)
```
Business Service (:4000)
├── Authentication (JWT validation)
├── User Management (CRUD)
├── Business Logic
└── Single Database
```

### After (Microservices)
```
Auth Service (:4001)              Business Service (:4000)
├── Authentication (login, JWT)   ├── Business Logic
├── User Management (CRUD)        ├── Authorization (policies)
├── Token Validation              └── Calls Auth Service for validation
└── Auth Database                 └── Main Database
```

## Migration Strategies

### Strategy 1: Keep Existing Routes (Recommended for Gradual Migration)

The existing routes in the business service can continue to work by using the auth service as a backend.

#### Option A: Proxy Pattern
Business service proxies authentication requests to auth service:

```go
// In business service routes
app.Post("/api/client/login", middleware.ProxyAuthServiceLogin(namespace.ClientKey.Name))
app.Post("/api/employee/login", middleware.ProxyAuthServiceLogin(namespace.EmployeeKey.Name))
app.Post("/api/admin/login", middleware.ProxyAuthServiceLogin(namespace.AdminKey.Name))
```

**Pros:**
- No client changes needed
- Gradual migration
- Backward compatible

**Cons:**
- Extra network hop
- Slight performance overhead

#### Option B: Direct Migration
Update clients to call auth service directly:

```go
// Clients call auth service directly
POST http://localhost:4001/auth/client/login
POST http://localhost:4001/auth/employee/login
POST http://localhost:4001/auth/admin/login
```

**Pros:**
- Better performance (direct call)
- Clean separation

**Cons:**
- Requires client updates
- Coordinated deployment

### Strategy 2: Update Middleware for Authorization

The business service middleware needs to validate tokens. Two options:

#### Option A: Local JWT Validation (Current Approach)
Business service validates JWTs locally using shared secret:

```go
// Already implemented in core/src/middleware/auth.go
func DenyUnauthorized(c *fiber.Ctx) error {
    // Parse JWT locally
    claims, err := handler.JWT(c).WhoAreYou()
    // Validate against database
    // Check policies
}
```

**Pros:**
- Fast (no network call)
- Works offline

**Cons:**
- Shared secret between services
- Token revocation requires database check

#### Option B: Remote Token Validation
Business service calls auth service to validate tokens:

```go
// New middleware in core/src/middleware/auth_service_proxy.go
func DenyUnauthorizedViaAuthService(c *fiber.Ctx) error {
    authClient := NewAuthServiceClient()
    claims, err := authClient.ValidateToken(token)
    c.Locals(namespace.RequestKey.Auth_Claims, claims)
}
```

**Pros:**
- Single source of truth
- Centralized token management
- Easy revocation

**Cons:**
- Network latency
- Auth service becomes critical dependency

## Recommended Migration Path

### Phase 1: Run Both Services (Current State)
1. ✅ Auth service running on :4001
2. ✅ Business service running on :4000
3. Both services share same databases (for now)
4. Business service continues using local JWT validation

### Phase 2: Migrate Authentication Endpoints
1. Update frontend/mobile apps to call auth service for login:
   - `POST /auth/client/login` → Auth Service :4001
   - `POST /auth/employee/login` → Auth Service :4001
   - `POST /auth/admin/login` → Auth Service :4001

2. Keep business service routes as proxies for backward compatibility:
   ```go
   // In core/src/config/api/routes/index.go
   apiRoutes.Post("/client/login", middleware.ProxyAuthServiceLogin(namespace.ClientKey.Name))
   ```

3. Add deprecation warnings to proxied endpoints

### Phase 3: Migrate User Management (Optional)
1. Move user CRUD operations to auth service
2. Update business service to call auth service APIs
3. Remove user management code from business service

### Phase 4: Separate Databases (Future)
1. Auth service uses only AuthDB
2. Business service uses only MainDB
3. Inter-service communication via HTTP/gRPC

## Environment Variables

### Auth Service
```env
# Auth service specific
AUTH_SERVICE_PORT=4001

# Auth database connection
POSTGRES_AUTH_HOST=localhost
POSTGRES_AUTH_PORT=5432
POSTGRES_AUTH_DB=mynute_auth
POSTGRES_AUTH_USER=postgres
POSTGRES_AUTH_PASSWORD=postgres

# JWT secret (must match business service for now)
JWT_SECRET=your-secret-key
```

### Business Service
```env
# Business service
APP_PORT=4000

# Main database connection
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=mynute_main
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres

# Auth service URL (for token validation)
AUTH_SERVICE_URL=http://localhost:4001

# JWT secret (must match auth service for now)
JWT_SECRET=your-secret-key
```

## Testing the Setup

### 1. Start Auth Service
```powershell
cd c:\code\mynute-go
go run cmd/auth-service/main.go
```

Expected output:
```
Server is starting at http://localhost:4001
```

### 2. Start Business Service
```powershell
cd c:\code\mynute-go
go run cmd/business-service/main.go
```

Expected output:
```
Server is starting at http://localhost:4000
```

### 3. Test Auth Service Directly
```powershell
# Test health check
curl http://localhost:4001/health

# Test client login
curl -X POST http://localhost:4001/auth/client/login `
  -H "Content-Type: application/json" `
  -d '{"email":"client@example.com","password":"password123"}'

# Response includes X-Auth-Token header
```

### 4. Test Token Validation
```powershell
# Get token from login response, then validate
curl -X POST http://localhost:4001/auth/validate `
  -H "X-Auth-Token: <your-token-here>"

# Returns user claims
```

### 5. Test Business Service with Auth
```powershell
# Business service validates token locally
curl http://localhost:4000/api/some-protected-endpoint `
  -H "X-Auth-Token: <your-token-here>"
```

## Code Organization

### Auth Service Files
```
cmd/auth-service/main.go          # Entry point
auth/api/
  ├── controller/
  │   ├── auth.go                 # Login, token validation
  │   ├── user.go                 # Client, employee CRUD
  │   └── admin.go                # Admin management
  └── routes/
      └── routes.go               # Route registration
```

### Business Service Files
```
cmd/business-service/main.go      # Entry point (wraps core.NewServer)
core/server.go                    # Existing server logic
core/src/middleware/
  ├── auth.go                     # Local JWT validation (existing)
  └── auth_service_proxy.go       # Auth service client (new)
```

## API Endpoints

### Auth Service (:4001)

#### Authentication
- `POST /auth/client/login` - Client login with password
- `POST /auth/client/login-with-code` - Client login with email code
- `POST /auth/client/send-login-code/email/:email` - Send login code
- `POST /auth/employee/login` - Employee login
- `POST /auth/employee/login-with-code` - Employee login with code
- `POST /auth/employee/send-login-code/email/:email` - Send login code
- `POST /auth/admin/login` - Admin login
- `POST /auth/validate` - Validate user token
- `POST /auth/validate-admin` - Validate admin token

#### User Management
- `POST /users/client` - Create client
- `GET /users/client/:id` - Get client by ID
- `GET /users/client/email/:email` - Get client by email
- `PATCH /users/client/:id` - Update client
- `DELETE /users/client/:id` - Delete client
- `POST /users/employee` - Create employee
- `GET /users/employee/:id` - Get employee by ID
- `GET /users/employee/email/:email` - Get employee by email
- `PATCH /users/employee/:id` - Update employee
- `DELETE /users/employee/:id` - Delete employee

#### Admin Management
- `GET /users/admin/are_there_any_superadmin` - Check if superadmin exists
- `POST /users/admin/first_superadmin` - Create first superadmin
- `GET /users/admin` - List all admins
- `POST /users/admin` - Create admin
- `GET /users/admin/:id` - Get admin by ID
- `PATCH /users/admin/:id` - Update admin
- `DELETE /users/admin/:id` - Delete admin

### Business Service (:4000)

#### Existing Business Endpoints (unchanged)
- `/api/company/*` - Company management
- `/api/branch/*` - Branch management
- `/api/appointment/*` - Appointment management
- `/api/service/*` - Service management
- etc.

#### Auth Proxies (optional, for backward compatibility)
- `POST /api/client/login` → Proxies to auth service
- `POST /api/employee/login` → Proxies to auth service
- `POST /api/admin/login` → Proxies to auth service

## Troubleshooting

### Auth service not reachable
Check `AUTH_SERVICE_URL` in business service:
```env
AUTH_SERVICE_URL=http://localhost:4001
```

### Token validation fails
Ensure JWT_SECRET is the same in both services:
```env
# Both .env files must have:
JWT_SECRET=same-secret-for-both-services
```

### Database connection errors
Check database environment variables are set correctly for each service.

### CORS issues
If calling auth service from browser, add CORS middleware in auth service main.go.

## Next Steps

1. ✅ Implement auth controllers
2. ✅ Wire up routes
3. ⏳ Update business service middleware
4. ⏳ Test both services
5. ⏳ Update frontend to call auth service
6. ⏳ Deploy both services
7. ⏳ Monitor and optimize

## Security Considerations

1. **JWT Secret Sharing**: Currently both services share JWT_SECRET. Consider:
   - Using asymmetric keys (RS256) instead of symmetric (HS256)
   - Auth service signs with private key
   - Business service validates with public key

2. **Token Revocation**: With local JWT validation:
   - Tokens can't be revoked until expiry
   - Consider shorter token lifetimes
   - Implement refresh tokens

3. **Service-to-Service Auth**: When business service calls auth service:
   - Use API keys or service tokens
   - Implement rate limiting
   - Use HTTPS in production

4. **Database Security**:
   - Use separate database users with minimal permissions
   - AuthDB user: only auth tables
   - MainDB user: only business tables

## Performance Considerations

1. **Token Validation**:
   - Local validation: ~0.1ms
   - Remote validation: ~10-50ms (network latency)
   - For high-traffic endpoints, use local validation

2. **Caching**:
   - Cache validated tokens (short TTL)
   - Cache user data to reduce database queries

3. **Database Connections**:
   - Use connection pooling
   - Monitor connection counts

4. **Monitoring**:
   - Add metrics for auth service calls
   - Monitor token validation latency
   - Track authentication success/failure rates
