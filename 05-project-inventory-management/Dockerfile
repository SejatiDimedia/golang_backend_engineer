# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-api ./cmd/server/main.go

# Stage 2: Runtime
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/inventory-api .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

CMD ["./inventory-api"]
