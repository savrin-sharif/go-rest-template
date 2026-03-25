#!/usr/bin/env bash
set -euo pipefail

TEMPLATE_MODULE="github.com/savrin-sharif/go-rest-template"
TEMPLATE_GIT_URL="https://${TEMPLATE_MODULE}.git"
NEW_MODULE="$TEMPLATE_MODULE"

START_DIR="$(pwd)"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORK_DIR="$SCRIPT_DIR"

# If not run from a template checkout, fetch it into a temp dir (similar to the Flutter script pattern).
if [[ ! -f "$WORK_DIR/go.mod" ]]; then
  if ! command -v git >/dev/null 2>&1; then
    echo "git is required to fetch the template (${TEMPLATE_GIT_URL})." >&2
    exit 1
  fi
  echo "Template not found locally; cloning ${TEMPLATE_GIT_URL} ..." >&2
  WORK_DIR="$(mktemp -d)"
  git clone --depth 1 "$TEMPLATE_GIT_URL" "$WORK_DIR" >&2
fi

cd "$WORK_DIR"

DEFAULT_NAME="${NEW_MODULE##*/}"
DEFAULT_DESCRIPTION="Production-ready starter for Go REST APIs following the golang-standards/project-layout conventions."

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

# Collect project metadata (prompt with defaults; env vars override prompts).
if [[ -z "${PROJECT_NAME:-}" ]]; then
  read -r -p "Project name [${DEFAULT_NAME}]: " PROJECT_NAME_INPUT
  PROJECT_NAME="${PROJECT_NAME_INPUT:-$DEFAULT_NAME}"
fi
PROJECT_NAME="${PROJECT_NAME:-$DEFAULT_NAME}"

if [[ -z "${PROJECT_DESCRIPTION:-}" ]]; then
  read -r -p "Project description [${DEFAULT_DESCRIPTION}]: " PROJECT_DESCRIPTION_INPUT
  PROJECT_DESCRIPTION="${PROJECT_DESCRIPTION_INPUT:-$DEFAULT_DESCRIPTION}"
fi
PROJECT_DESCRIPTION="${PROJECT_DESCRIPTION:-$DEFAULT_DESCRIPTION}"

DEFAULT_BADGES="[![Go Version](https://img.shields.io/badge/go-${GO_VERSION_MAJOR_MINOR}-blue)](https://go.dev) [![CI](https://img.shields.io/badge/ci-ready-lightgrey)](https://github.com/${NEW_MODULE}/actions)"
PROJECT_BADGES="${PROJECT_BADGES:-$DEFAULT_BADGES}"
PROJECT_NAME_ESC=${PROJECT_NAME//\//\\/}
PROJECT_DESCRIPTION_ESC=${PROJECT_DESCRIPTION//\//\\/}

# Rewrite import paths in source files
find . -path './.git' -prune -o -type f -print0 | \
  xargs -0 perl -pi -e "s|${TEMPLATE_MODULE}|${NEW_MODULE}|g"

# Sync config default name.
perl -0777 -pi -e "s/(?m)^app:\\n\\s+name:.*$/app:\\n  name: ${PROJECT_NAME_ESC}/" configs/config.yaml

# Sync OpenAPI info.
perl -0777 -pi -e "s/^info:\\n  title:.*\\n  version:.*\\n  description:.*$/info:\\n  title: ${PROJECT_NAME_ESC} API\\n  version: 1.0.0\\n  description: ${PROJECT_DESCRIPTION_ESC}/m" api/openapi.yaml

# Copy scaffold into a folder named after the project in the starting directory.
TARGET_DIR="${START_DIR%/}/${PROJECT_NAME}"
mkdir -p "$TARGET_DIR"
rsync -a --exclude '.git' "$WORK_DIR"/ "$TARGET_DIR"/

# Remove initializer from the generated project.
rm -f "$TARGET_DIR/init-project.sh" "$TARGET_DIR/init-go.sh"

echo "Project initialized with module '${NEW_MODULE}' in ${TARGET_DIR}"
cat <<'TODO'
Next steps:
  1) cd into the project folder.
  2) Update configs/config.yaml and api/openapi.yaml to match your service.
  3) Run: go mod tidy && go test ./...
  4) Initialize git and create your first commit if desired.
TODO
