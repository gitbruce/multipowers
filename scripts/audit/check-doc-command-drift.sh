#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

PATTERN='mp test run|mp coverage check|mp validate --type no-shell|mp validate --strict-no-shell|mp cost report'

matches="$(rg -n "$PATTERN" docs/architecture CLAUDE.md SAFEGUARDS.md 2>/dev/null || true)"
if [[ -z "$matches" ]]; then
  echo "doc command drift check: ok"
  exit 0
fi

bad="$(printf '%s\n' "$matches" | rg -v 'command-ownership.md|mp-devx|迁移|deprecated|moved|Move|Planned Migrations|Compatibility Policy' || true)"
if [[ -n "$bad" ]]; then
  echo "doc command drift check: failed"
  echo "$bad"
  exit 1
fi

echo "doc command drift check: ok"
