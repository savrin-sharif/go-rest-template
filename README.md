# Go REST API Template

Production-ready starter for Go (1.21+) REST APIs following the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) conventions.

## Features
- Chi router with request ID, recovery, CORS, and slog-based request logging
- Health and welcome endpoints out of the box
- Environment-first configuration via Viper with sane defaults
- Built-in database config with single-source DSN (`APP_DATABASE_URL`)
- Graceful shutdown on `SIGINT`/`SIGTERM`
- Make targets for build/test/lint/tidy
- Multi-stage Dockerfile and docker-compose for local runs
- `init-project.sh` to rename the module, rewrite imports, and re-init git

## Layout
```
cmd/server          # Application entry point
internal/config     # Configuration + logger construction
internal/server     # HTTP server & router wiring
internal/handler    # HTTP handlers
internal/middleware # Custom middleware (logging, recovery)
pkg/httputil        # Reusable helpers safe to import by other projects
configs/            # Default config
api/                # OpenAPI specification
```

## Quickstart
```bash
cp .env.example .env
make tidy
make run
# visit http://localhost:8080 and http://localhost:8080/health
```

The app auto-loads `.env` on startup and requires `APP_DATABASE_URL`.

## Docker
```bash
docker build -t go-rest-template .
docker run -p 8080:8080 \
  -e APP_DATABASE_URL="postgres://postgres:postgres@host.docker.internal:5432/app_db?sslmode=disable" \
  go-rest-template
# or
docker-compose up --build
```

`docker-compose.yaml` includes a ready-to-use PostgreSQL service and wires app
database setup from a single env source: `APP_DATABASE_URL` (required at runtime).

## Re-brand the template
```bash
./init-project.sh github.com/yourname/awesome-service
```

The initializer also offers an optional DB setup prompt and writes your chosen
`APP_DATABASE_URL` into `.env` and `.env.example` in the generated project.
