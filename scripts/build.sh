#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if ! command -v go >/dev/null 2>&1; then
  echo "ERROR: go is required to build mp binaries." >&2
  exit 127
fi

mkdir -p "${ROOT_DIR}/bin"
go build -o "${ROOT_DIR}/bin/mp" "${ROOT_DIR}/cmd/octo"
go build -o "${ROOT_DIR}/bin/mp-devx" "${ROOT_DIR}/cmd/octo-devx"

echo "Built: ${ROOT_DIR}/bin/mp"
echo "Built: ${ROOT_DIR}/bin/mp-devx"
