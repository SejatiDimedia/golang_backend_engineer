# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Salin go.mod dan go.sum terlebih dahulu untuk caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Salin semua file kode sumber
COPY . .

# Build binary aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/server/main.go

# Stage 2: Runtime
FROM alpine:3.18

WORKDIR /app

# Salin binary dari stage builder
COPY --from=builder /app/url-shortener .
# Salin .env.example sebagai fallback jika .env tidak dipasang
COPY --from=builder /app/.env.example .env

EXPOSE 8080

CMD ["./url-shortener"]
