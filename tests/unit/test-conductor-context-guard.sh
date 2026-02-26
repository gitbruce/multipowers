#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

bash -n "$ROOT/custom/lib/conductor-context.sh"
rg -n "is_spec_driven_command|conductor_context_complete|ensure_conductor_context" "$ROOT/custom/lib/conductor-context.sh" >/dev/null
rg -n 'CLAUDE\.md' "$ROOT/custom/lib/conductor-context.sh" >/dev/null
rg -n "ensure_spec_command_context|spec_prompt|run_octo_init_interactive" "$ROOT/scripts/orchestrate.sh" >/dev/null

echo "PASS test-conductor-context-guard"
