# Go Workspace Setup

This project uses **Go workspaces** to manage multiple microservices in a monorepo.

## Structure

```
mynute-go/
├── go.work                      # Workspace configuration
├── go.work.sum                  # Workspace checksums
├── go.mod                       # Root module (main launcher)
├── main.go                      # Multi-service launcher
│
├── services/core/
│   ├── go.mod                   # Core service module
│   ├── go.sum                   # Core dependencies
│   └── ...
│
├── services/auth/
│   ├── go.mod                   # Auth service module
│   ├── go.sum                   # Auth dependencies
│   └── ...
│
└── services/email/
    ├── go.mod                   # Email service module
    ├── go.sum                   # Email dependencies
    └── ...
```

## Commands

### Running Services

```bash
# Run all services together
go run .

# Run individual services
go run ./cmd/business-service   # Core service
go run ./cmd/auth-service        # Auth service
go run ./cmd/email-service       # Email service
```

### Managing Dependencies

```bash
# Sync workspace (after adding dependencies to any service)
go work sync

# Add dependency to a specific service
cd services/core
go get github.com/some/package
go mod tidy

# Update all services
go work sync
```

### Testing

```bash
# Test all services
go test ./...

# Test specific service
go test ./services/core/...
go test ./services/auth/...
go test ./services/email/...
```

### Building

```bash
# Build all services
go build ./cmd/...

# Build specific service
go build ./cmd/business-service
go build ./cmd/auth-service
go build ./cmd/email-service
```

## Benefits

1. **True Independence**: Each service has its own go.mod with specific dependencies
2. **Isolated Updates**: Upgrade dependencies per service without affecting others
3. **Smaller Builds**: Each service only includes what it needs
4. **Docker Optimization**: Better layer caching per service
5. **No Conflicts**: Prevents issues like the Swagger double-registration we had

## Workspace Commands Reference

```bash
# Initialize workspace (already done)
go work init
go work use . ./services/core ./services/auth ./services/email

# Add a new service module
go work use ./services/new-service

# Edit workspace
go work edit

# Sync workspace dependencies
go work sync
```

## Migration Notes

We migrated from a single go.mod to this workspace structure to:
- Resolve Swagger registration conflicts between services
- Enable true microservices independence
- Improve build times and Docker caching
- Follow Go best practices for monorepos

Each service can now be developed, tested, and deployed completely independently.
