#!/usr/bin/env bash
set -euo pipefail

TEMPLATE_MODULE="github.com/savrin-sharif/go-rest-template"
TEMPLATE_GIT_URL="${TEMPLATE_GIT_URL:-https://${TEMPLATE_MODULE}.git}"

START_DIR="$(pwd)"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

WORK_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$WORK_DIR"
}
trap cleanup EXIT

# If script is run from the template checkout, use it offline.
if [[ -f "$SCRIPT_DIR/go.mod" ]] && grep -Eq "^module[[:space:]]+${TEMPLATE_MODULE}$" "$SCRIPT_DIR/go.mod"; then
  rsync -a --exclude '.git' "$SCRIPT_DIR"/ "$WORK_DIR"/
else
  if ! command -v git >/dev/null 2>&1; then
    echo "git is required to fetch the template (${TEMPLATE_GIT_URL})." >&2
    exit 1
  fi
  echo "Template files not found next to script; cloning ${TEMPLATE_GIT_URL} ..." >&2
  if ! git clone --depth 1 "$TEMPLATE_GIT_URL" "$WORK_DIR" >&2; then
    echo "Failed to clone template. Check internet access or run from a local template clone." >&2
    exit 1
  fi
fi

cd "$WORK_DIR"

NEW_MODULE="${1:-${NEW_MODULE:-}}"
DEFAULT_NAME="my-service"
if [[ -n "$NEW_MODULE" ]]; then
  DEFAULT_NAME="${NEW_MODULE##*/}"
fi

# Collect project metadata (prompt with defaults; env vars override prompts).
if [[ -z "${PROJECT_NAME:-}" ]]; then
  read -r -p "Project folder name [${DEFAULT_NAME}]: " PROJECT_NAME_INPUT
  PROJECT_NAME="${PROJECT_NAME_INPUT:-$DEFAULT_NAME}"
fi
PROJECT_NAME="${PROJECT_NAME:-$DEFAULT_NAME}"

if [[ -z "$NEW_MODULE" ]]; then
  DEFAULT_MODULE="github.com/yourname/${PROJECT_NAME}"
  read -r -p "Go module path [${DEFAULT_MODULE}]: " NEW_MODULE_INPUT
  NEW_MODULE="${NEW_MODULE_INPUT:-$DEFAULT_MODULE}"
fi

if [[ "$NEW_MODULE" != */* ]] || [[ "$NEW_MODULE" == *" "* ]]; then
  echo "Invalid module path: '${NEW_MODULE}'" >&2
  echo "Example: github.com/yourname/${PROJECT_NAME}" >&2
  exit 1
fi

if [[ "$PROJECT_NAME" == *"/"* ]] || [[ "$PROJECT_NAME" == *" "* ]]; then
  echo "Invalid project folder name: '${PROJECT_NAME}'" >&2
  echo "Use a simple folder name like: ${NEW_MODULE##*/}" >&2
  exit 1
fi

DEFAULT_DESCRIPTION="Production-ready starter for Go REST APIs following the golang-standards/project-layout conventions."
DEFAULT_DB_URL="$(awk -F= '/^APP_DATABASE_URL=/{sub(/^APP_DATABASE_URL=/, "", $0); print $0; exit}' .env.example 2>/dev/null || true)"
DEFAULT_DB_URL="${DEFAULT_DB_URL:-postgres://postgres:postgres@localhost:5432/app_db?sslmode=disable}"

# Update go.mod module path
if command -v go >/dev/null 2>&1; then
  go mod edit -module "$NEW_MODULE"
else
  echo "Go is required to run this script." >&2
  exit 1
fi

# Align go directive with locally installed Go version (major.minor).
GO_VERSION_FULL="$(go env GOVERSION)"           # e.g., go1.22.5
GO_VERSION_SHORT="${GO_VERSION_FULL#go}"       # 1.22.5
GO_VERSION_MAJOR_MINOR="$(echo "$GO_VERSION_SHORT" | cut -d. -f1-2)" # 1.22

go mod edit -go "$GO_VERSION_MAJOR_MINOR"

if [[ -z "${PROJECT_DESCRIPTION:-}" ]]; then
  read -r -p "Project description [${DEFAULT_DESCRIPTION}]: " PROJECT_DESCRIPTION_INPUT
  PROJECT_DESCRIPTION="${PROJECT_DESCRIPTION_INPUT:-$DEFAULT_DESCRIPTION}"
fi
PROJECT_DESCRIPTION="${PROJECT_DESCRIPTION:-$DEFAULT_DESCRIPTION}"

PROJECT_NAME_ESC=${PROJECT_NAME//\//\\/}
PROJECT_DESCRIPTION_ESC=${PROJECT_DESCRIPTION//\//\\/}

# Optional DB setup. If skipped, use the template example URL.
if [[ -z "${APP_DATABASE_URL:-}" ]]; then
  read -r -p "Configure database URL now? [y/N]: " CONFIGURE_DB_INPUT
  CONFIGURE_DB_INPUT="${CONFIGURE_DB_INPUT:-N}"
  CONFIGURE_DB_INPUT="$(printf '%s' "$CONFIGURE_DB_INPUT" | tr '[:upper:]' '[:lower:]')"

  case "$CONFIGURE_DB_INPUT" in
    y|yes)
      read -r -p "APP_DATABASE_URL [${DEFAULT_DB_URL}]: " APP_DATABASE_URL_INPUT
      APP_DATABASE_URL="${APP_DATABASE_URL_INPUT:-$DEFAULT_DB_URL}"
      ;;
    *)
      APP_DATABASE_URL="$DEFAULT_DB_URL"
      ;;
  esac
fi
APP_DATABASE_URL="${APP_DATABASE_URL:-$DEFAULT_DB_URL}"

# Rewrite import paths in source files
find . -path './.git' -prune -o -type f -print0 | \
  xargs -0 perl -pi -e "s|${TEMPLATE_MODULE}|${NEW_MODULE}|g"

# Sync config default name.
perl -0777 -pi -e "s/(?m)^app:\\n\\s+name:.*$/app:\\n  name: ${PROJECT_NAME_ESC}/" configs/config.yaml

# Sync OpenAPI info.
perl -0777 -pi -e "s/^info:\\n  title:.*\\n  version:.*\\n  description:.*$/info:\\n  title: ${PROJECT_NAME_ESC} API\\n  version: 1.0.0\\n  description: ${PROJECT_DESCRIPTION_ESC}/m" api/openapi.yaml

# Copy scaffold into a folder named after the project in the starting directory.
TARGET_DIR="${START_DIR%/}/${PROJECT_NAME}"
if [[ -e "$TARGET_DIR" ]]; then
  echo "Target directory already exists: ${TARGET_DIR}" >&2
  echo "Choose another project folder name or remove the existing folder first." >&2
  exit 1
fi
mkdir -p "$TARGET_DIR"
rsync -a --exclude '.git' "$WORK_DIR"/ "$TARGET_DIR"/

# Write project-local env files with selected DB URL.
printf 'APP_DATABASE_URL=%s\n' "$APP_DATABASE_URL" > "$TARGET_DIR/.env.example"
printf 'APP_DATABASE_URL=%s\n' "$APP_DATABASE_URL" > "$TARGET_DIR/.env"

# Remove initializer from the generated project.
rm -f "$TARGET_DIR/init-project.sh" "$TARGET_DIR/init-go.sh"

echo "Project initialized with module '${NEW_MODULE}' in ${TARGET_DIR}"
cat <<'TODO'
Next steps:
  1) cd into the project folder.
  2) Review .env (or .env.example) and adjust APP_DATABASE_URL if needed.
  3) Update configs/config.yaml and api/openapi.yaml to match your service.
  4) Run: go mod tidy && go test ./...
TODO
