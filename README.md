# Uptime Monitor

A self-hosted distributed uptime monitoring application. The backend is written in Go; the frontend (React) is planned for a later phase.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Go](https://go.dev/dl/) 1.26+ (to run the API locally)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) (optional — only needed if you want to run migrations from the Makefile instead of Docker)

## Local setup

### 1. Clone the repository

```bash
git clone https://github.com/publiosilva/uptime-monitor.git
cd uptime-monitor
```

### 2. Configure environment variables

Copy the example env file:

```bash
cp .env.example .env
```

Default values work out of the box for local development:

| Variable | Default |
|----------|---------|
| `POSTGRES_USER` | `uptime` |
| `POSTGRES_PASSWORD` | `uptime` |
| `POSTGRES_DB` | `uptime_monitor` |
| `POSTGRES_PORT` | `5432` |
| `DATABASE_URL` | `postgres://uptime:uptime@localhost:5432/uptime_monitor?sslmode=disable` |

### 3. Start PostgreSQL and run migrations

```bash
docker compose up -d
```

This starts Postgres and automatically runs pending migrations via the `migrate` service. On first run you should see `1/u init_schema` in the logs:

```bash
docker compose logs migrate
```

Re-running `docker compose up -d` is safe — only new migrations are applied.

### 4. Verify the database

List tables:

```bash
docker compose exec postgres psql -U uptime -d uptime_monitor -c '\dt'
```

Expected tables: `users`, `monitors`, `heartbeats`, `schema_migrations`.

To open an interactive shell:

```bash
docker compose exec postgres psql -U uptime -d uptime_monitor
```

Use `\dt` to list tables and `\q` to quit.

### 5. Run the API

```bash
cd backend
go run ./cmd/api
```

The API listens on `http://localhost:3333`.

## Installing golang-migrate (optional)

Migrations run automatically through Docker Compose, so you do **not** need the CLI for a basic local setup. Install it only if you want to use the Makefile targets in `backend/` (e.g. creating new migrations or rolling back).

### macOS (Homebrew)

```bash
brew install golang-migrate
```

### Linux

```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
```

### Go install

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1
```

Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is on your `PATH`.

## Database migrations

Migration files live in `backend/db/migrations/` and follow [golang-migrate](https://github.com/golang-migrate/migrate) naming (`NNNNNN_name.up.sql` / `.down.sql`).

### Via Docker Compose (default)

```bash
docker compose up -d
```

### Via Makefile (requires golang-migrate CLI)

From the `backend/` directory:

```bash
make migrate-up        # apply pending migrations
make migrate-down      # roll back the last migration
make migrate-version   # show current version
make migrate-create    # create a new migration pair
```

## Stopping services

```bash
docker compose down
```

To remove the database volume as well (full reset):

```bash
docker compose down -v
```

After a volume reset, run `docker compose up -d` again to recreate the schema from scratch.
