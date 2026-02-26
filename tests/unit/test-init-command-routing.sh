#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

for f in "$ROOT/.claude/commands/init.md" "$ROOT/custom/commands/init.md"; do
  rg -n '"\$\{CLAUDE_PLUGIN_ROOT\}/scripts/orchestrate\.sh" --dir "\$PWD" init' "$f" >/dev/null
  rg -n '\.multipowers/CLAUDE\.md|\.multipowers/FAQ\.md|\.multipowers/context/runtime\.json' "$f" >/dev/null
  rg -n 'Do not use `Write/Edit/Update`' "$f" >/dev/null
  echo "PASS init routing contract: $(basename "$f")"
done

echo "PASS test-init-command-routing"
