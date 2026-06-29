<!--
TEMPLATE: SETUP.md
-->

# Setup: [Project Name]

## Prerequisites

- Go [version]
- Docker & Docker Compose
- [Other: PostgreSQL client, Redis CLI, etc. — only if needed outside Docker]

## Environment variables

| Variable | Description | Example |
|---|---|---|
| `PORT` | | `8080` |
| `DATABASE_URL` | | `postgres://user:pass@localhost:5432/dbname` |
| | | |

Copy `.env.example` to `.env` and fill in values before running.

## Local setup (Docker — recommended)

```bash
docker-compose up -d
```

This starts: [list services — app, postgres, redis, etc.]

## Local setup (without Docker)

```bash
# 1. Install dependencies
go mod download

# 2. Start dependent services (Postgres, Redis, etc.) manually or via docker-compose up -d postgres

# 3. Run migrations
[migration command]

# 4. Run the service
go run cmd/api/main.go
```

## Verifying it's running

```bash
curl http://localhost:PORT/health
```

Expected response: `{"status": "ok"}`

## Running tests

```bash
go test ./...
```

See `TESTING.md` for strategy and coverage expectations.

## Troubleshooting

| Issue | Likely cause | Fix |
|---|---|---|
| | | |
