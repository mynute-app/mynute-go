# Etapa 1: build da aplicação
FROM golang:1.23.11-alpine AS builder

RUN apk update && apk upgrade

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Compila o binário da aplicação
RUN go build -o app .

# Etapa 2: imagem final
FROM alpine:latest

WORKDIR /app

# Copia o binário
COPY --from=builder /app/app .

# Copia o restante do projeto (arquivos externos usados no runtime)
COPY --from=builder /app ./

EXPOSE 3000

CMD ["./app"]