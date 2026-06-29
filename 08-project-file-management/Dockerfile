# Stage 1: Build
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o file-server cmd/server/main.go

# Stage 2: Runtime
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/file-server .
COPY .env.example .env

EXPOSE 8080

CMD ["./file-server"]
