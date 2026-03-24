#!/usr/bin/env bash
set -euo pipefail

TEMPLATE_MODULE="github.com/savrin-sharif/go-rest-template"

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <new-module-path>" >&2
  exit 1
fi

NEW_MODULE="$1"
DEFAULT_NAME="${NEW_MODULE##*/}"
DEFAULT_DESCRIPTION="Production-ready starter for Go REST APIs following the golang-standards/project-layout conventions."
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Update go.mod module path
if command -v go >/dev/null 2>&1; then
  go mod edit -module "$NEW_MODULE"
else
  echo "Go is required to run this script." >&2
  exit 1
fi

# Align go directive and Dockerfile with locally installed Go version (major.minor).
GO_VERSION_FULL="$(go env GOVERSION)"           # e.g., go1.22.5
GO_VERSION_SHORT="${GO_VERSION_FULL#go}"       # 1.22.5
GO_VERSION_MAJOR_MINOR="$(echo "$GO_VERSION_SHORT" | cut -d. -f1-2)" # 1.22

go mod edit -go "$GO_VERSION_MAJOR_MINOR"

# Update Dockerfile ARG GO_VERSION to match.
perl -pi -e "s/^(ARG\\s+GO_VERSION=).*/\\1${GO_VERSION_MAJOR_MINOR}/" Dockerfile

# Collect project metadata interactively (defaults are safe to accept).
read -r -p "Project name [${DEFAULT_NAME}]: " PROJECT_NAME
PROJECT_NAME=${PROJECT_NAME:-$DEFAULT_NAME}

read -r -p "Project description [${DEFAULT_DESCRIPTION}]: " PROJECT_DESCRIPTION
PROJECT_DESCRIPTION=${PROJECT_DESCRIPTION:-$DEFAULT_DESCRIPTION}

DEFAULT_BADGES="[![Go Version](https://img.shields.io/badge/go-${GO_VERSION_MAJOR_MINOR}-blue)](https://go.dev) [![CI](https://img.shields.io/badge/ci-ready-lightgrey)](https://github.com/${NEW_MODULE}/actions)"
read -r -p "Badges markdown (optional) [default badges]: " PROJECT_BADGES
PROJECT_BADGES=${PROJECT_BADGES:-$DEFAULT_BADGES}
PROJECT_NAME_ESC=${PROJECT_NAME//\//\\/}
PROJECT_DESCRIPTION_ESC=${PROJECT_DESCRIPTION//\//\\/}
PROJECT_DOCKER_TAG=$(echo "$PROJECT_NAME" | tr '[:upper:]' '[:lower:]' | tr ' /_' '-')

# Rewrite import paths in source files
find . -path './.git' -prune -o -type f -print0 | \
  xargs -0 perl -pi -e "s|${TEMPLATE_MODULE}|${NEW_MODULE}|g"

# Update README with provided metadata.
cat > README.md <<EOF
# ${PROJECT_NAME}
${PROJECT_BADGES}

${PROJECT_DESCRIPTION}

## Features
- Chi router with request ID, recovery, CORS, and slog-based request logging
- Health and welcome endpoints out of the box
- Environment-first configuration via Viper with sane defaults
- Graceful shutdown on SIGINT/SIGTERM
- Make targets for build/test/lint/tidy
- Multi-stage Dockerfile and docker-compose for local runs
- init-project.sh to rename the module, rewrite imports, and re-init git

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
docker build -t ${PROJECT_DOCKER_TAG} .
docker run -p 8080:8080 ${PROJECT_DOCKER_TAG}
# or
docker-compose up --build
```

## Re-brand the template
```bash
./init-project.sh github.com/yourname/another-service
```
EOF

# Sync config default name.
perl -0777 -pi -e "s/(?m)^app:\\n\\s+name:.*$/app:\\n  name: ${PROJECT_NAME_ESC}/" configs/config.yaml

# Sync OpenAPI info.
perl -0777 -pi -e "s/^info:\\n  title:.*\\n  version:.*\\n  description:.*$/info:\\n  title: ${PROJECT_NAME_ESC} API\\n  version: 1.0.0\\n  description: ${PROJECT_DESCRIPTION_ESC}/m" api/openapi.yaml

# Refresh dependencies
go mod tidy

# Reset git history
if [[ -d .git ]]; then
  rm -rf .git
fi

git init
git add .
git commit -m "chore: bootstrap project from template"

echo "Project initialized with module '${NEW_MODULE}'."
echo
cat <<'TODO'
Next steps:
  1) Update README.md (project name, description, badges).
  2) Review configs/config.yaml defaults and environment overrides.
  3) Extend api/openapi.yaml with your endpoints.
  4) Set container image name/tag in Dockerfile & docker-compose.yaml.
  5) Add LICENSE and CI workflow (e.g., GitHub Actions for lint/test).
  6) Run: make tidy && make test && make lint
TODO
