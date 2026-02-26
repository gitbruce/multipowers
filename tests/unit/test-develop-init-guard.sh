#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

rg -n 'Step 0: Enforce Conductor Context Guard' "$ROOT/.claude/commands/develop.md" >/dev/null
rg -n 'Do \*\*not\*\* perform any implementation action' "$ROOT/.claude/commands/develop.md" >/dev/null
rg -n '^### STEP 0: Enforce `\.multipowers` Context Guard \(MANDATORY\)$' "$ROOT/.claude/skills/flow-develop.md" >/dev/null
rg -n 'Do not run `Write`, `Edit`, `Update`' "$ROOT/.claude/skills/flow-develop.md" >/dev/null

echo "PASS test-develop-init-guard"
