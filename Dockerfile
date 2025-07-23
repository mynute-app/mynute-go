# Etapa 1: build da aplicação
FROM golang:1.23.11-alpine AS builder

RUN apk update && apk upgrade

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

# Etapa 2: imagem final
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 3000

CMD ["./app"]