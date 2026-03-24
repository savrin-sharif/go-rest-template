# Go REST API Template

Production-ready starter for Go (1.21+) REST APIs following the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) conventions.

## Features
- Chi router with request ID, recovery, CORS, and slog-based request logging
- Health and welcome endpoints out of the box
- Environment-first configuration via Viper with sane defaults
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
make tidy
make run
# visit http://localhost:8080 and http://localhost:8080/health
```

## Docker
```bash
docker build -t go-rest-template .
docker run -p 8080:8080 go-rest-template
# or
docker-compose up --build
```

## Re-brand the template
```bash
./init-project.sh github.com/yourname/awesome-service
```

