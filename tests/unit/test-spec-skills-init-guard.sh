#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

for f in \
  "$ROOT/.claude/skills/flow-discover.md" \
  "$ROOT/.claude/skills/flow-define.md" \
  "$ROOT/.claude/skills/flow-develop.md" \
  "$ROOT/.claude/skills/flow-deliver.md" \
  "$ROOT/.claude/skills/skill-code-review.md" \
  "$ROOT/.claude/skills/skill-debate.md"; do
  [[ -f "$f" ]] || continue
  rg -n "\\.multipowers" "$f" >/dev/null
  rg -n 'CLAUDE\.md' "$f" >/dev/null
  rg -n '"\$\{CLAUDE_PLUGIN_ROOT\}/scripts/orchestrate\.sh" --dir "\$PWD" init' "$f" >/dev/null
  rg -n 'continue without init|Do not run .*Write.*Edit.*Update' "$f" >/dev/null
  echo "PASS guard markers: $(basename "$f")"
done

echo "PASS test-spec-skills-init-guard"
