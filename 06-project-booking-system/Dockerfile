# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/server/main.go

# Run stage
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/.env.example ./.env

EXPOSE 8080

CMD ["./server"]
