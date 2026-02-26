#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")/../.." && pwd)"
project="${1:-$root}"
outdir="$root/docs/plans/evidence/go-big-bang/perf"
mkdir -p "$outdir"
raw="$outdir/preflight-samples.txt"
: > "$raw"
for i in $(seq 1 30); do
  start=$(date +%s%3N)
  "$root/bin/octo" context guard --dir "$project" --json >/dev/null || true
  end=$(date +%s%3N)
  echo $((end-start)) >> "$raw"
done
sort -n "$raw" > "$raw.sorted"
p95_index=$(( (30*95 + 99)/100 ))
p95=$(sed -n "${p95_index}p" "$raw.sorted")
cat > "$outdir/report.md" <<R
# Preflight Benchmark

samples: 30
p95_ms: ${p95}

Targets:
- hot path <= 50ms
- cold start <= 120ms
R
