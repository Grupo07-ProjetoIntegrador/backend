# Estágio de Compilação
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instala dependências do sistema necessárias para compilação (se houver)
RUN apk add --no-cache git

# Copia os arquivos de dependência do Go
COPY go.mod go.sum ./
RUN go mod download

# Copia o código-fonte restante
COPY . .

# Compila o binário otimizado para produção
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api/main.go

# Estágio de Execução
FROM alpine:latest

WORKDIR /app

# Adiciona certificados CA para conexões HTTPS externas seguras
RUN apk add --no-cache ca-certificates tzdata

# Copia o binário compilado do estágio anterior
COPY --from=builder /app/main .

# O binário lê variáveis de ambiente diretamente do sistema ou arquivo .env se presente
EXPOSE 8080

CMD ["./main"]
