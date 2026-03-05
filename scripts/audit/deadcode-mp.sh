#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

OUT_TXT="docs/plans/evidence/command-boundary/deadcode-baseline-2026-03-05.txt"
OUT_MD="tmp/unused_go.md"
mkdir -p "$(dirname "$OUT_TXT")" tmp

TMP_RAW="$(mktemp)"
trap 'rm -f "$TMP_RAW" "$TMP_FILES"' EXIT

# Static reachability from runtime/devx command roots (excluding tests)
go run golang.org/x/tools/cmd/deadcode@latest -test=false ./cmd/mp ./cmd/mp-devx > "$TMP_RAW"
cp "$TMP_RAW" "$OUT_TXT"

TMP_FILES="$(mktemp)"
awk -F: '{print $1}' "$TMP_RAW" | sort -u > "$TMP_FILES"

METHOD_COUNT="$(wc -l < "$TMP_RAW" | tr -d ' ')"
FILE_COUNT="$(wc -l < "$TMP_FILES" | tr -d ' ')"

{
  echo '# Go files/methods without mp/mp-devx entry (static reachability)'
  echo
  echo "Date: $(date +%F)"
  echo
  echo 'Method:'
  echo '- Command roots: `./cmd/mp`, `./cmd/mp-devx`'
  echo '- Tool: `go run golang.org/x/tools/cmd/deadcode@latest -test=false ./cmd/mp ./cmd/mp-devx`'
  echo '- Interpretation: listed funcs are not reachable from runtime/devx command roots (tests excluded).'
  echo
  echo '## Summary'
  echo "- Unreachable methods: ${METHOD_COUNT}"
  echo "- Files containing unreachable methods: ${FILE_COUNT}"
  echo
  echo '## Files'
  sed 's#^#- `#; s#$#`#' "$TMP_FILES"
  echo
  echo '## Methods'
  sed 's#^#- `#; s#$#`#' "$TMP_RAW"
  echo
  echo '## Notes'
  echo '- This is a static call-graph result, not a direct delete list.'
  echo '- Some methods may be intentionally reserved for future wiring or only used by tests/tools.'
} > "$OUT_MD"

echo "wrote $OUT_TXT"
echo "wrote $OUT_MD"
