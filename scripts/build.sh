#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if ! command -v go >/dev/null 2>&1; then
  echo "ERROR: go is required to build mp binaries." >&2
  exit 127
fi

mkdir -p "${ROOT_DIR}/.claude-plugin/bin"
(
  cd "${ROOT_DIR}"
  go build -o "${ROOT_DIR}/.claude-plugin/bin/mp" ./cmd/mp
  go build -o "${ROOT_DIR}/.claude-plugin/bin/mp-devx" ./cmd/mp-devx
)

echo "Built: ${ROOT_DIR}/.claude-plugin/bin/mp"
echo "Built: ${ROOT_DIR}/.claude-plugin/bin/mp-devx"
