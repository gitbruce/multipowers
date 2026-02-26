#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
found=0
while IFS= read -r f; do
  lines=$(wc -l < "$f")
  if [ "$lines" -gt 500 ]; then
    echo "WARN: $f has $lines lines (>500)"
    found=1
  fi
done < <(cd "$root" && find . -type f -name '*.go' -not -path './vendor/*' | sort)
exit 0
