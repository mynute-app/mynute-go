# Etapa 1: build da aplicação
FROM golang:1.23.3-alpine3.21 AS builder

WORKDIR /mynute-go

RUN apk update && apk add --no-cache git ca-certificates && \
    apk upgrade --no-cache

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . .

# Build the auth service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./mynute-auth-app ./cmd/auth-service

# Etapa 2: imagem final - usando distroless para maior segurança
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /mynute-go

COPY --from=builder /mynute-go/mynute-auth-app .

EXPOSE 4001

USER nonroot:nonroot

CMD ["./mynute-auth-app"]
