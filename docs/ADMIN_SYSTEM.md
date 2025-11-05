# Admin System Documentation

## Overview

The Admin System provides **system-wide administrator access** to your multi-tenant SaaS platform. Admins can access all tenants and bypass normal tenant-based restrictions, making it ideal for:

- **Platform Support**: Help customers across all tenants
- **System Monitoring**: Audit and monitor all tenant activities
- **Data Management**: Manage data across multiple tenants

---

## 📋 Table of Contents

1. [Architecture](#architecture)
2. [Database Schema](#database-schema)
3. [Setup & Installation](#setup--installation)
4. [Creating Admins](#creating-admins)
5. [Admin Authentication](#admin-authentication)
6. [API Endpoints](#api-endpoints)
7. [Security Considerations](#security-considerations)
8. [Testing](#testing)

---

## Architecture

### Components

```
┌─────────────────────────────────────────────────────────────┐
│                     Admin System                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐    ┌──────────────┐    ┌─────────────┐  │
│  │   Models     │    │  Controllers │    │ Middleware  │  │
│  ├──────────────┤    ├──────────────┤    ├─────────────┤  │
│  │ Admin        │───>│ AdminAuth    │───>│ WhoAreYou   │  │
│  │ RoleAdmin    │    │ Admin        │    │ DenyUn...   │  │
│  └──────────────┘    └──────────────┘    └─────────────┘  │
│                                                             │
│  ┌──────────────┐    ┌──────────────┐                      │
│  │   DTOs       │    │   Handler    │                      │
│  ├──────────────┤    ├──────────────┤                      │
│  │ AdminClaims  │    │ JWT.Admin    │                      │
│  │ AdminLogin   │    │              │                      │
│  └──────────────┘    └──────────────┘                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Key Features

- **Role-Based Access Control**: Admins can have multiple roles (superadmin, support, auditor)
- **Tenant Bypass**: Admins bypass tenant-specific RBAC/ABAC policies
- **JWT Authentication**: Secure token-based authentication
- **Password Security**: bcrypt hashing with automatic validation

---

## Database Schema

### Tables

#### `public.admins`
```sql
CREATE TABLE public.admins (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(100) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    is_active  BOOLEAN DEFAULT true,
    meta       JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
```

#### `public.role_admins`
```sql
CREATE TABLE public.role_admins (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP WITH TIME ZONE
);
```

#### `public.admin_role_admins` (Join Table)
```sql
CREATE TABLE public.admin_role_admins (
    admin_id      UUID NOT NULL REFERENCES public.admins(id) ON DELETE CASCADE,
    role_admin_id UUID NOT NULL REFERENCES public.role_admins(id) ON DELETE CASCADE,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (admin_id, role_admin_id)
);
```

### Default Roles

- **superadmin**: Full access to all tenants and resources
- **support**: Customer support with read-only access to tenant data
- **auditor**: Read-only access for compliance and auditing

---

## Setup & Installation

### 1. Run Database Migration

```powershell
# Apply the admin migration
cd c:\code\mynute-go
.\scripts\migrate.ps1 -action up
```

Or manually run:
```powershell
psql -U <username> -d <database> -f migrations/20251031140524_add_admin_tables.up.sql
```

### 2. Seed Default Admin Data

```powershell
# Run the admin seeder
go run cmd/seed-admin/main.go
```

This will create:
- 3 default admin roles (superadmin, support, auditor)
- 1 default admin user with superadmin role

### 3. Environment Variables (Optional)

Add to your `.env` file:

```env
# Default admin credentials (optional)
DEFAULT_ADMIN_EMAIL=admin@yourcompany.com
DEFAULT_ADMIN_PASSWORD=YourSecurePassword123!
DEFAULT_ADMIN_NAME=System Administrator

# JWT Secret (required)
JWT_SECRET=your-secret-key-here
```

---

## Creating Admins

### Method 1: Using the Seeder (First Admin)

```powershell
go run cmd/seed-admin/main.go
```

**Default Credentials:**
- Email: `admin@mynute.com`
- Password: `Admin@123456`

⚠️ **Change this password immediately after first login!**

### Method 2: Using the API (Subsequent Admins)

Only **superadmin** users can create new admins.

```bash
# Login as superadmin first
curl -X POST http://localhost:3000/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@mynute.com",
    "password": "Admin@123456"
  }'

# Use the returned token to create new admin
curl -X POST http://localhost:3000/admin/create \
  -H "Content-Type: application/json" \
  -H "X-Auth-Token: <your-token>" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123!",
    "is_active": true,
    "roles": ["support"]
  }'
```

---

## Admin Authentication

### Login Flow

1. **Admin sends credentials** to `/admin/auth/login`
2. **System validates** email and password
3. **JWT token is generated** with admin claims
4. **Token is returned** to the client
5. **Client includes token** in `X-Auth-Token` header for subsequent requests

### JWT Token Claims

```json
{
  "data": {
    "id": "uuid-of-admin",
    "name": "Admin Name",
    "email": "admin@example.com",
    "password": "hashed-password",
    "is_admin": true,
    "is_active": true,
    "roles": ["superadmin"],
    "type": "admin"
  },
  "exp": 1234567890
}
```

### Middleware Flow

```
Request → WhoAreYou → DenyUnauthorized → Controller
            ↓              ↓
      Check Admin    Check Admin Bypass
      Token First    (If admin → Allow)
```

---

## API Endpoints

### Authentication Endpoints

#### `POST /admin/auth/login`
Authenticate admin user and receive JWT token.

**Request:**
```json
{
  "email": "admin@example.com",
  "password": "SecurePassword123!"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "admin": {
      "id": "uuid",
      "name": "Admin Name",
      "email": "admin@example.com",
      "is_active": true,
      "roles": ["superadmin"]
    }
  }
}
```

#### `GET /admin/auth/me`
Get current admin's information.

**Headers:**
```
X-Auth-Token: <your-jwt-token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Admin Name",
    "email": "admin@example.com",
    "is_active": true,
    "roles": ["superadmin"]
  }
}
```

#### `POST /admin/auth/refresh`
Refresh JWT token.

**Headers:**
```
X-Auth-Token: <your-jwt-token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "new-jwt-token"
  }
}
```

---

### Admin Management Endpoints

#### `GET /admin/list`
List all admin users. **Requires:** Any admin role.

#### `POST /admin/create`
Create new admin. **Requires:** superadmin role.

#### `PATCH /admin/:id`
Update admin. **Requires:** superadmin role.

#### `DELETE /admin/:id`
Delete admin (soft delete). **Requires:** superadmin role.

---

### Role Management Endpoints

#### `GET /admin/roles`
List all admin roles. **Requires:** Any admin role.

#### `POST /admin/roles`
Create new role. **Requires:** superadmin role.

#### `PATCH /admin/roles/:id`
Update role. **Requires:** superadmin role.

#### `DELETE /admin/roles/:id`
Delete role. **Requires:** superadmin role.

---

## Security Considerations

### Password Requirements

- Minimum 8 characters
- Must include uppercase, lowercase, number, and special character
- Automatically hashed using bcrypt (cost factor: 10)

### Token Security

- Tokens expire after 90 days
- Password changes invalidate existing tokens
- Tokens include hashed password for validation

### Best Practices

1. **Change default password immediately** after seeding
2. **Use environment variables** for production credentials
3. **Limit superadmin access** - only assign to trusted users
4. **Audit admin actions** - implement logging for admin activities
5. **Use HTTPS** in production to protect tokens
6. **Rotate JWT secrets** periodically
7. **Implement rate limiting** on login endpoints

### Admin Bypass Logic

Admins bypass normal tenant restrictions:
- ✅ Superadmins have **full access** to all endpoints
- ✅ Other admin roles can access **all tenant data**
- ✅ No CompanyID validation required for admins
- ⚠️ Admin actions should be logged for audit trails

---

## Testing

### Manual Testing with curl

1. **Login as admin:**
```bash
curl -X POST http://localhost:3000/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@mynute.com","password":"Admin@123456"}'
```

2. **Get admin info:**
```bash
curl http://localhost:3000/admin/auth/me \
  -H "X-Auth-Token: <token-from-login>"
```

3. **List all admins:**
```bash
curl http://localhost:3000/admin/list \
  -H "X-Auth-Token: <token>"
```

4. **Create new admin:**
```bash
curl -X POST http://localhost:3000/admin/create \
  -H "Content-Type: application/json" \
  -H "X-Auth-Token: <token>" \
  -d '{
    "name": "Support User",
    "email": "support@example.com",
    "password": "Support123!",
    "is_active": true,
    "roles": ["support"]
  }'
```

### Automated Testing

Create unit tests in `test/src/admin_test.go`:

```go
package test

import (
    "testing"
    "mynute-go/services/core/src/config/db/model"
)

func TestAdminLogin(t *testing.T) {
    // Test admin login flow
}

func TestAdminBypass(t *testing.T) {
    // Test that admins bypass tenant restrictions
}

func TestSuperAdminAccess(t *testing.T) {
    // Test superadmin can create other admins
}
```

---

## Troubleshooting

### Common Issues

**Issue: "Admin account is inactive"**
- Solution: Update `is_active = true` in database or via API

**Issue: "Invalid credentials"**
- Solution: Verify email/password, check password hasn't been changed

**Issue: "Superadmin role required"**
- Solution: Ensure admin has `superadmin` role assigned

**Issue: Token expired**
- Solution: Login again or use `/admin/auth/refresh` endpoint

### Database Queries

Check admin roles:
```sql
SELECT a.email, r.name as role
FROM public.admins a
JOIN public.admin_role_admins ara ON a.id = ara.admin_id
JOIN public.role_admins r ON ara.role_admin_id = r.id
WHERE a.email = 'admin@mynute.com';
```

Reset admin password manually:
```sql
-- Generate bcrypt hash of "NewPassword123!" first
UPDATE public.admins
SET password = '$2a$10$...'  -- bcrypt hash here
WHERE email = 'admin@mynute.com';
```

---

## Migration Rollback

To remove admin tables:

```powershell
.\scripts\migrate.ps1 -action down
```

Or manually:
```sql
DROP TABLE IF EXISTS public.admin_role_admins CASCADE;
DROP TABLE IF EXISTS public.admins CASCADE;
DROP TABLE IF EXISTS public.role_admins CASCADE;
```

---

## Summary

The Admin System provides:
- ✅ System-wide admin access across all tenants
- ✅ Role-based permissions (superadmin, support, auditor)
- ✅ Secure JWT authentication
- ✅ Automatic tenant bypass for admins
- ✅ RESTful API endpoints for admin management
- ✅ Easy seeding for initial setup

For questions or issues, refer to the main project documentation or contact the development team.
