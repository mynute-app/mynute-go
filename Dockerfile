# Etapa 1: build da aplicação
FROM golang:1.23.11-alpine AS builder

WORKDIR /app/mynute-go

RUN apk update && apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

# Etapa 2: imagem final (runtime enxuto)
FROM alpine:latest

WORKDIR /app/mynute-go

# Copia a pasta inteira do projeto + binário
COPY --from=builder /app/mynute-go ./

EXPOSE 3000

CMD ["./app"]