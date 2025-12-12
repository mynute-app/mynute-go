# Etapa 1: build da aplicação
FROM golang:1.23-alpine AS builder

WORKDIR /mynute-go

RUN apk update && apk add --no-cache git curl

# Install Atlas CLI
RUN curl -sSf https://atlasgo.sh | sh

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . .

# Build main application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./mynute-backend-app

# Build migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o ./migrate-tool migrate/main.go

# Build seed tool
RUN CGO_ENABLED=0 GOOS=linux go build -o ./seed-tool cmd/seed/main.go

# Etapa 2: imagem final
FROM alpine:latest

WORKDIR /mynute-go

# Install netcat for database health check and tzdata for timezone support
RUN apk add --no-cache netcat-openbsd tzdata

# Copy Atlas CLI from builder
COPY --from=builder /usr/local/bin/atlas /usr/local/bin/atlas

# Copy application binary
COPY --from=builder /mynute-go/mynute-backend-app .

# Copy migration and seed tools
COPY --from=builder /mynute-go/migrate-tool .
COPY --from=builder /mynute-go/seed-tool .

# Copy migration files and Atlas config
COPY --from=builder /mynute-go/migrations ./migrations
COPY --from=builder /mynute-go/atlas.hcl ./atlas.hcl

# Copy static files
COPY --from=builder /mynute-go/static ./static

EXPOSE 4000

# Just run the app - migrations/seeding are separate manual steps
CMD ["./mynute-backend-app"]