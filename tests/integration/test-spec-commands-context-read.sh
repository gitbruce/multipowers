#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

rg -n "apply_conductor_context_to_prompt|load_conductor_context_for_prompt" "$ROOT/scripts/orchestrate.sh" "$ROOT/custom/lib/conductor-context.sh" >/dev/null

echo "PASS test-spec-commands-context-read"
