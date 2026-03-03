#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if ! command -v go >/dev/null 2>&1; then
  echo "ERROR: go is required to build mp binaries." >&2
  exit 127
fi

# Create output directories
mkdir -p "${ROOT_DIR}/.claude-plugin/bin"
mkdir -p "${ROOT_DIR}/.claude-plugin/runtime"

# Build policy first (fail fast if config is invalid)
echo "Building policy..."
(
  cd "${ROOT_DIR}"
  go run ./cmd/mp-devx -action build-policy -config-dir config -output-dir .claude-plugin/runtime
)

# Build binaries
echo "Building binaries..."
(
  cd "${ROOT_DIR}"
  go build -o "${ROOT_DIR}/.claude-plugin/bin/mp" ./cmd/mp
  go build -o "${ROOT_DIR}/.claude-plugin/bin/mp-devx" ./cmd/mp-devx
)

echo ""
echo "Build complete:"
echo "  - ${ROOT_DIR}/.claude-plugin/runtime/policy.json"
echo "  - ${ROOT_DIR}/.claude-plugin/bin/mp"
echo "  - ${ROOT_DIR}/.claude-plugin/bin/mp-devx"
