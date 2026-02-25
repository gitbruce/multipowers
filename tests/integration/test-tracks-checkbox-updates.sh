#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

rg -n "ensure_conductor_track_file|mark_track_checkbox|finalize_spec_tracking" "$ROOT/scripts/orchestrate.sh" "$ROOT/custom/lib/conductor-context.sh" >/dev/null
rg -n "conductor/tracks/" "$ROOT/custom/docs/customizations/conductor-context.md" >/dev/null

echo "PASS test-tracks-checkbox-updates"
