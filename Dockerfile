# Etapa 1: build da aplicação
FROM golang:1.23.11-alpine AS builder

WORKDIR /mynute-go

RUN apk update && apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

# Etapa 2: imagem final
FROM alpine:latest

WORKDIR /mynute-go

COPY --from=builder /mynute-go .

EXPOSE 3000

CMD ["./app"]