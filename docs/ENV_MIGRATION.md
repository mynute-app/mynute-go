# Environment Variables Migration

## ⚠️ DEPRECATED: Root `.env` and `.env.example`

As of the microservices refactoring, the root `.env` and `.env.example` files are **deprecated** and should no longer be used.

## New Structure

Each service now manages its own environment variables:

### Core/Business Service
- **Location:** `core/.env` and `core/.env.example`
- **Contains:** Business service configuration, cloud storage, Grafana, MinIO, etc.

### Auth Service
- **Location:** `auth/.env` and `auth/.env.example`
- **Contains:** Auth service configuration, JWT settings, email/SMTP for auth purposes

## Migration Guide

If you have an existing root `.env` file, you need to split it into service-specific files:

### 1. Create Core Service Environment

```bash
cp core/.env.example core/.env
```

Then copy these variables from your root `.env` to `core/.env`:
- `APP_ENV`
- `APP_PORT`
- `APP_HOST`
- `POSTGRES_*` (all PostgreSQL variables)
- `LOKI_HOST`
- `PGADMIN_*`
- `GF_SECURITY_*` (Grafana)
- `GRAFANA_*`
- `GOOGLE_*` (OAuth)
- `R2_*` (Cloudflare R2)
- `MINIO_*`
- `S3_*` (AWS S3)
- `STORAGE_DRIVER`
- `SWAGGER_*`
- `RESEND_*` (if used by core)
- `MAILHOG_*` (if used by core)

### 2. Create Auth Service Environment

```bash
cp auth/.env.example auth/.env
```

Then copy these variables from your root `.env` to `auth/.env`:
- `APP_ENV`
- `AUTH_SERVICE_PORT`
- `POSTGRES_*` (all PostgreSQL variables)
- `JWT_SECRET`
- `EMAIL_FROM`
- `SMTP_*` (all SMTP variables for auth emails)
- `RESEND_*` (if used by auth)
- `MAILHOG_*` (if used by auth)

### 3. Remove Root Environment Files

After migrating, you can safely remove:
```bash
rm .env
rm .env.example
```

## Docker Compose Updates

The docker-compose files have been updated to reference local `.env` files:

- `core/docker-compose.dev.yml` → uses `core/.env`
- `core/docker-compose.prod.yml` → uses `core/.env`
- `auth/docker-compose.dev.yml` → uses `auth/.env`
- `auth/docker-compose.prod.yml` → uses `auth/.env`

## Running Services Locally

When running services with `go run`, the services will automatically load their respective `.env` files:

```bash
# Business Service loads core/.env
go run cmd/business-service/main.go

# Auth Service loads auth/.env
go run cmd/auth-service/main.go

# Running both together
go run main.go  # Loads both core/.env and auth/.env
```

## Shared Variables

Some variables like `POSTGRES_*` are duplicated across both services because both need database access. This is intentional and allows each service to:
- Use different databases if needed
- Scale independently
- Have different connection settings
- Be deployed separately

## Environment-Specific Configuration

Remember to update `APP_ENV` appropriately:
- `dev` - Development with auto-migrations and seeding
- `test` - Testing environment
- `prod` - Production (manual migrations required)
