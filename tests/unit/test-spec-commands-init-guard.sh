#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

for f in \
  "$ROOT/.claude/commands/plan.md" \
  "$ROOT/.claude/commands/discover.md" \
  "$ROOT/.claude/commands/define.md" \
  "$ROOT/.claude/commands/develop.md" \
  "$ROOT/.claude/commands/deliver.md" \
  "$ROOT/.claude/commands/embrace.md" \
  "$ROOT/.claude/commands/research.md" \
  "$ROOT/.claude/commands/review.md" \
  "$ROOT/.claude/commands/debate.md"; do
  [[ -f "$f" ]] || continue
  rg -n 'CLAUDE\.md' "$f" >/dev/null
  rg -n '"\$\{CLAUDE_PLUGIN_ROOT\}/scripts/orchestrate\.sh" --dir "\$PWD" init' "$f" >/dev/null
  rg -n 'hard-stop|fail-fast|stop with an initialization failure|STOP' "$f" >/dev/null
  echo "PASS command guard: $(basename "$f")"
done

echo "PASS test-spec-commands-init-guard"
