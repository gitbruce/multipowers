#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

rg -n "main.*upstream mirror|multipowers" "$ROOT/custom/docs/README.md" "$ROOT/custom/docs/sync/upstream-sync-playbook.md" >/dev/null

echo "PASS test-main-upstream-discipline"
