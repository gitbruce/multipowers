#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

REPORT="tmp/unused_go.md"

if [[ ! -f "$REPORT" ]]; then
  echo "missing report: $REPORT" >&2
  exit 1
fi

rg -n "^# Go files/methods without mp/mp-devx entry" "$REPORT" >/dev/null
rg -n "^## Summary$" "$REPORT" >/dev/null
rg -n "^## Files$" "$REPORT" >/dev/null
rg -n "^## Methods$" "$REPORT" >/dev/null

echo "unused go report check: ok"
