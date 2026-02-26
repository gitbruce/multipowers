#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

# /octo:init must not silently scaffold defaults in non-interactive mode.
rg -n '/octo:init requires interactive wizard mode; non-interactive fallback is disabled' "$ROOT/scripts/orchestrate.sh" >/dev/null
rg -n 'ERROR: /octo:init must run as an interactive wizard to create project context\.' "$ROOT/scripts/orchestrate.sh" >/dev/null

if rg -n '/octo:init running in non-interactive mode; applying defaults for context scaffolding' "$ROOT/scripts/orchestrate.sh" >/dev/null; then
  echo "FAIL: found legacy non-interactive default scaffolding path"
  exit 1
fi

echo "PASS test-init-wizard-required"
