# Etapa 1: build da aplicação
FROM golang:1.23.11-alpine AS builder

WORKDIR /mynute-go

RUN apk update && apk add --no-cache git

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./mynute-backend-app

# Etapa 2: imagem final
FROM alpine:latest

WORKDIR /mynute-go

COPY --from=builder /mynute-go/mynute-backend-app .

EXPOSE 4000

CMD ["./mynute-backend-app"]