#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")/../.." && pwd)"
target="${1:-$(mktemp -d /tmp/octo-e2e-XXXX)}"
mkdir -p "$target"
outdir="$root/docs/plans/evidence/go-big-bang/e2e"
mkdir -p "$outdir"
{
  echo "target=$target"
  "$root/bin/octo" init --dir "$target" --json
  "$root/bin/octo" context guard --dir "$target" --json
  "$root/bin/octo" plan --dir "$target" --prompt "e2e plan" --json
  "$root/bin/octo" develop --dir "$target" --prompt "e2e develop" --json
  "$root/bin/octo" debate --dir "$target" --prompt "e2e debate" --json || true
  "$root/bin/octo" validate --dir "$target" --json
} > "$outdir/transcript.txt" 2>&1
