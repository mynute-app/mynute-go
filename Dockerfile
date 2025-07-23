# Etapa 1: build da aplicação
FROM golang:1.23.11-alpine AS builder

WORKDIR /build

RUN apk update && apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

# Etapa 2: imagem final
FROM alpine:latest

WORKDIR /app/mynute-go

# Copia tudo pro diretório que o código espera
COPY --from=builder /build ./ 

EXPOSE 3000

CMD ["sh"]