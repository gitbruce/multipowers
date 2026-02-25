#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

rg -n 'ensure_spec_command_context "\$COMMAND" "\$\*"' "$ROOT/scripts/orchestrate.sh" >/dev/null
rg -n "Missing project conductor context" "$ROOT/scripts/orchestrate.sh" >/dev/null

echo "PASS test-spec-commands-auto-init"
