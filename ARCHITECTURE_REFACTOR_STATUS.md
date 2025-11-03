# Auth Service Architecture Refactor - Status Report

**Date:** November 3, 2025  
**Status:** 🚧 WORK IN PROGRESS - Does not compile yet  
**Branch:** admin  
**Commit:** 11f912c

## Overview

Major architectural refactoring to separate authentication concerns from business logic by creating a unified User model in the auth service and moving business-specific user data to the core service.

---

## ✅ Completed Work

### 1. Auth Service - Unified User Model
**File:** `auth/model/user.go`

Created a single `User` model to replace separate Admin, Client, Employee models:

```go
type User struct {
    ID        uuid.UUID
    Email     string
    Password  string      // hashed
    Verified  bool
    Type      string      // "admin", "client", "employee"
    Meta      UserMeta    // Only auth-related metadata
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Benefits:**
- Single source of truth for authentication
- Simplified auth logic - one table, one model
- Type field differentiates user roles
- Clean separation from business logic

### 2. Auth Service - Simplified UserMeta
**File:** `auth/model/json/user_meta.go`

Removed business fields, kept only authentication metadata:

```go
type UserMeta struct {
    Login LoginConfig  // validation codes, password reset
}
```

**Removed:**
- DesignConfig (moved to core service)
- Business-specific configurations
- File: `auth/model/json/design.go` (deleted)

### 3. Auth Service - Independent Utilities
Created auth's own utility packages (decoupled from core):

**Created directories:**
- `auth/lib/` - error.go, validator.go, database.go, context.go, time.go, random.go, email.go, send_response.go, dto.go, env.go
- `auth/handler/` - jwt.go, auth.go, gorm.go
- `auth/dto/` - admin.go, client.go, employee.go, auth.go, error.go
- `auth/config/namespace/` - namespace constants

**Key functions:**
- `lib.PrepareEmail()` - Email validation and URL decoding
- `lib.GenerateRandomInt()` - Validation code generation
- `handler.ComparePassword()` - Password verification
- `handler.JWT(c).Encode()` - JWT token generation
- `lib.Session(c)` - Database session management

### 4. Core Service - Updated Models

#### Admin Model (`core/src/config/db/model/admin.go`)
```go
type Admin struct {
    UserID    uuid.UUID      `gorm:"primaryKey"`  // FK to auth.users.id
    Name      string
    Surname   string
    IsActive  bool
    Roles     []RoleAdmin
    Meta      mJSON.UserMeta  // Business metadata (design, etc.)
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Removed:** Email, Password, Verified, HashPassword(), MatchPassword()  
**Changed:** ID → UserID (primary key)

#### Client Model (`core/src/config/db/model/client.go`)
```go
type Client struct {
    UserID    uuid.UUID      `gorm:"primaryKey"`  // FK to auth.users.id
    Name      string
    Surname   string
    Phone     string
    Meta      mJSON.UserMeta  // Business metadata
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Removed:** Email, Password, Verified, ClientMeta wrapper, password methods  
**Changed:** ID → UserID (primary key)

#### Employee Model (`core/src/config/db/model/employee.go`)
```go
type Employee struct {
    UserID              uuid.UUID           `gorm:"primaryKey"`  // FK to auth.users.id
    Name                string
    Surname             string
    Phone               string
    Tags                []string
    SlotTimeDiff        uint
    WorkSchedule        []EmployeeWorkRange
    Appointments        []Appointment
    CompanyID           uuid.UUID
    Branches            []*Branch
    Services            []*Service
    Roles               []*Role
    TimeZone            string
    TotalServiceDensity uint32
    Meta                mJSON.UserMeta  // Business metadata (design, etc.)
    CreatedAt           time.Time
    UpdatedAt           time.Time
}
```

**Removed:** Email, Password, Verified, password methods  
**Changed:** ID → UserID (primary key)  
**Updated:** All references from `e.ID` to `e.UserID`

---

## ⚠️ Partially Completed (HAS ERRORS)

### Auth Controllers
**Files:** `auth/api/controller/auth.go`, `user.go`, `admin.go`, `policy.go`, `endpoint.go`, `authorization.go`

**Status:** Started refactoring but has compile errors

**What was done:**
- Updated `LoginByPassword()` to use unified User model
- Updated `LoginByEmailCode()` to use unified User model
- Removed reflection code - now direct field access
- Query users by email AND type: `WHERE email = ? AND type = ?`

**What needs fixing:**
1. Function signatures changed - callers need updating
2. Some functions still reference old `model.Client`, `model.Employee`, `model.Admin`
3. `GenerateLoginValidationCode()` needs type parameter instead of model instance
4. `SendLoginValidationCodeByEmail()` needs type parameter
5. All Employee/Admin login endpoints need updating
6. Import errors (unused mJSON import)

**Example fixes needed:**
```go
// OLD (broken):
token, err := LoginByPassword(namespace.ClientKey.Name, &model.Client{}, c)

// NEW (correct):
token, err := LoginByPassword(namespace.ClientKey.Name, c)
```

---

## 📋 TODO - Next Steps

### Priority 1: Fix Auth Controllers
- [ ] Update all LoginByPassword calls to remove model parameter
- [ ] Update all LoginByEmailCode calls to remove model parameter
- [ ] Rewrite GenerateLoginValidationCode to use type string
- [ ] Rewrite SendLoginValidationCodeByEmail to use type string
- [ ] Update Employee login endpoints
- [ ] Update Admin login endpoints
- [ ] Remove unused imports (mJSON, reflect)

### Priority 2: Update Auth DTOs
- [ ] Review `auth/dto/employee.go` - remove WorkSchedule, Design fields
- [ ] Review `auth/dto/admin.go` - keep only auth-related fields
- [ ] Review `auth/dto/client.go` - keep only auth-related fields
- [ ] Add Type field to CreateUser, UpdateUser DTOs
- [ ] Update swagger documentation

### Priority 3: Update user.go Controller
**File:** `auth/api/controller/user.go`

Current status: Partially updated, needs more work

**Required changes:**
- [ ] Change CreateUser/UpdateUser to work with unified User model
- [ ] Add type parameter to creation endpoints
- [ ] Remove separate Client/Employee endpoints or make them aliases
- [ ] Update all CRUD operations to use users table

### Priority 4: Database Migration
**Create:** `auth/migrations/YYYYMMDDHHMMSS_unified_users_table.up.sql`

```sql
-- Drop old tables (if exist)
DROP TABLE IF EXISTS public.admin_users;
DROP TABLE IF EXISTS public.client_users;
DROP TABLE IF EXISTS public.employee_users;

-- Create unified users table
CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('admin', 'client', 'employee')),
    meta JSONB DEFAULT '{"login": {}}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON public.users(email);
CREATE INDEX idx_users_type ON public.users(type);
CREATE INDEX idx_users_deleted_at ON public.users(deleted_at);
```

**Down migration:**
```sql
DROP TABLE IF EXISTS public.users;
```

### Priority 5: Update Core Service
- [ ] Update core migrations to add UserID foreign keys to admins/clients/employees
- [ ] Update core controllers to call auth service APIs for user creation
- [ ] Implement auth service client in core (HTTP calls to validate tokens)
- [ ] Update core DTOs to exclude auth fields

### Priority 6: Integration & Testing
- [ ] Test auth service compiles independently
- [ ] Test login endpoints with type parameter
- [ ] Verify zero core dependencies in auth package
- [ ] Test token generation for all user types
- [ ] Test validation code flow
- [ ] Integration test: Create user in auth → Create profile in core

### Priority 7: Documentation
- [ ] Update ADMIN_SYSTEM.md with new architecture
- [ ] Document auth service API endpoints
- [ ] Document core service integration with auth
- [ ] Update README with service responsibilities
- [ ] Add sequence diagrams for user creation flow
- [ ] Document migration process from old to new structure

---

## Architecture Design

### Service Responsibilities

#### Auth Service (Port 4001)
**Responsibilities:**
- User authentication (login, logout)
- JWT token generation and validation
- Password hashing and verification
- Validation code generation and storage
- Email validation
- Authorization rules (policies, endpoints, resources)

**Does NOT handle:**
- Email sending (returns validation code to caller)
- Image uploads
- Business logic
- Appointments
- Company management

**Database:**
- Single `users` table with `type` field
- Authorization tables (policies, endpoints, resources)

#### Core Service (Port 4000)
**Responsibilities:**
- Business logic (appointments, services, schedules)
- User profiles (admins, clients, employees with business fields)
- Email sending (receives validation code from auth)
- Image uploads (profile, logo, banner)
- Company/branch management
- Payment processing
- Notifications

**Database:**
- `admins` table (UserID FK to auth.users.id)
- `clients` table (UserID FK to auth.users.id)
- `employees` table (UserID FK to auth.users.id)
- All business tables (appointments, services, companies, etc.)

### Data Flow Examples

#### User Registration
```
1. Client → POST /auth/users
   Body: { email, password, type: "client" }
   
2. Auth Service:
   - Validates email/password
   - Hashes password
   - Creates user in users table
   - Returns user_id and token
   
3. Client → POST /core/clients
   Headers: { X-Auth-Token: <token> }
   Body: { name, surname, phone }
   
4. Core Service:
   - Validates token with auth service
   - Extracts user_id from token
   - Creates client profile with user_id as PK
   - Returns client profile
```

#### Login Flow
```
1. Client → POST /auth/client/login
   Body: { email, password }
   
2. Auth Service:
   - Queries: SELECT * FROM users WHERE email = ? AND type = 'client'
   - Verifies password
   - Generates JWT token
   - Returns token in X-Auth-Token header
   
3. Client uses token for all subsequent requests
```

#### Validation Code Flow
```
1. Client → POST /auth/client/send-login-code/email/{email}

2. Auth Service:
   - Generates 6-digit code
   - Stores in users.meta.login.validation_code
   - Returns: { validation_code: "123456", email: "user@example.com" }
   
3. Core Service (or external email service):
   - Receives validation code from response
   - Sends email to user with code
   
4. Client → POST /auth/client/login-with-code
   Body: { email, code: "123456" }
   
5. Auth Service:
   - Validates code and expiry
   - Clears code
   - Generates JWT token
   - Returns token
```

---

## File Changes Summary

### Files Created (38 new files)
```
auth/config/namespace/index.go
auth/dto/admin.go
auth/dto/auth.go
auth/dto/client.go
auth/dto/employee.go
auth/dto/error.go
auth/handler/auth.go
auth/handler/gorm.go
auth/handler/jwt.go
auth/lib/context.go
auth/lib/database.go
auth/lib/dto.go
auth/lib/email.go
auth/lib/env.go
auth/lib/error.go
auth/lib/random.go
auth/lib/send_response.go
auth/lib/time.go
auth/lib/validator.go
auth/model/json/login.go
auth/model/json/user_meta.go
auth/model/role.go
auth/model/user.go
shared/handler/auth.go
shared/handler/gorm.go
shared/handler/jwt.go
shared/lib/database.go
shared/lib/error.go
```

### Files Modified
```
auth/api/controller/admin.go
auth/api/controller/auth.go
auth/api/controller/authorization.go
auth/api/controller/endpoint.go
auth/api/controller/policy.go
auth/api/controller/user.go
auth/model/base_model.go
core/src/config/db/model/admin.go
core/src/config/db/model/client.go
core/src/config/db/model/employee.go
```

### Files Deleted
```
auth/model/json/design.go  (business logic, moved to core)
```

### Files Still in Auth (To Review)
```
auth/model/admin.go     # Should be deleted
auth/model/client.go    # Should be deleted
auth/model/employee.go  # Should be deleted
```

---

## Known Issues

1. **Compile Errors in auth/api/controller/auth.go**
   - Function signature mismatches
   - Unused imports
   - References to deleted models

2. **Shared Directory**
   - Created `shared/` directory with duplicated code
   - Should be deleted - each service has its own utilities

3. **Auth Model Files**
   - Old admin.go, client.go, employee.go still exist in auth/model/
   - Should be deleted since we now have unified User model

4. **Core Dependencies**
   - Core service needs HTTP client to call auth APIs
   - Need to implement token validation flow

5. **Missing Migrations**
   - No migration file for unified users table
   - No migration to update core tables with UserID FK

---

## How to Continue

### Step 1: Clean up auth controllers (Priority)
```bash
# Fix compile errors in auth.go
# Focus on these functions:
- LoginClientByPassword
- LoginClientByEmailCode
- SendClientLoginValidationCodeByEmail
- LoginEmployeeByPassword
- LoginEmployeeByEmailCode  
- SendEmployeeLoginValidationCodeByEmail
- LoginAdminByPassword
- GetAdminClaims
```

### Step 2: Remove old models from auth
```bash
rm auth/model/admin.go
rm auth/model/client.go
rm auth/model/employee.go
rm -rf shared/
```

### Step 3: Create migration
```bash
# Create migration file
# Run migration in auth database
# Verify users table created correctly
```

### Step 4: Test
```bash
# Try to compile auth service
cd auth
go build ./...

# Fix any remaining compile errors
# Test each login endpoint
```

---

## Questions & Decisions Needed

1. **Authentication Flow:** Should core service call auth HTTP APIs or share database?
   - Recommendation: HTTP APIs for true service separation
   
2. **Email Sending:** Who triggers email sending?
   - Current design: Auth returns code, core sends email
   - Alternative: Message queue pattern
   
3. **User Creation:** Two-step or single-step process?
   - Current: Create in auth, then create in core
   - Alternative: Core service creates both via auth API

4. **Migration Strategy:** How to migrate existing users?
   - Need data migration script
   - Map existing admin/client/employee records to unified users table

---

## Related Documentation

- `/docs/ADMIN_SYSTEM.md` - Admin system overview (needs update)
- `/docs/MIGRATIONS_QUICKSTART.md` - Migration guide
- `/ADMIN_IMPLEMENTATION.md` - Implementation details

---

**Last Updated:** November 3, 2025  
**Next Review:** After fixing auth controller compile errors
