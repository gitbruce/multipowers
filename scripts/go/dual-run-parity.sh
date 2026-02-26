#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")/../.." && pwd)"
project="${1:-$root}"
outdir="$root/docs/plans/evidence/go-big-bang/dual-run"
mkdir -p "$outdir"
for cmd in plan debate embrace; do
  "$root/bin/octo" "$cmd" --dir "$project" --prompt "smoke parity" --json > "$outdir/go-${cmd}.json" || true
  "$root/scripts/orchestrate.sh" --dir "$project" "$cmd" "smoke parity" > "$outdir/sh-${cmd}.txt" 2>&1 || true
done
cat > "$outdir/report.md" <<R
# Dual-Run Parity Report

Compared commands:
- plan
- debate
- embrace

Generated artifacts:
- go-*.json
- sh-*.txt

Note: compare structure/flow semantics, not byte-identical markdown output.
R
