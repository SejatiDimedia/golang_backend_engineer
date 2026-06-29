# Stage 1: Build
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o prompt-server cmd/server/main.go

# Stage 2: Runtime
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/prompt-server .
COPY .env.example .env

# Buat folder certs dummy (untuk nanti dimount ke public.key asli)
RUN mkdir -p certs

EXPOSE 8082

CMD ["./prompt-server"]
