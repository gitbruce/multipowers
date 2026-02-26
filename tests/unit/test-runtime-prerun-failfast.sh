#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

# Runtime pre-run hooks must fail-fast unless explicitly configured otherwise.
rg -n 'apply_pre_run_context "\$agent_type" "\$\{phase:-unknown\}" "\$\{role:-none\}" \|\| return \$\?' "$ROOT/scripts/orchestrate.sh" >/dev/null
rg -n 'if \[\[ "\$on_fail" != "continue" \]\]; then' "$ROOT/scripts/orchestrate.sh" >/dev/null

echo "PASS test-runtime-prerun-failfast"
