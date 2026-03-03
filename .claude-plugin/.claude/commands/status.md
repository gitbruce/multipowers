---
command: status
description: "Show Claude Octopus workflow and provider status"
skill: skill-status
---

# Status

Display current Claude Octopus state, active agents, and provider readiness.

Run:

```bash
MP_BIN=""
if [[ -n "${CLAUDE_PLUGIN_ROOT:-}" ]] && [[ -x "${CLAUDE_PLUGIN_ROOT}/bin/mp" ]]; then
  MP_BIN="${CLAUDE_PLUGIN_ROOT}/bin/mp"
elif [[ -x "$PWD/.claude-plugin/bin/mp" ]]; then
  MP_BIN="$PWD/.claude-plugin/bin/mp"
elif [[ -x "./.claude-plugin/bin/mp" ]]; then
  MP_BIN="./.claude-plugin/bin/mp"
else
  echo "mp binary not found. Restart Claude Code and use /mp:* commands, or build via scripts/build.sh." >&2
  exit 1
fi

"$MP_BIN" status
```

Then summarize:
- Current mode (dev/knowledge/auto)
- Provider readiness
- Active agents and results availability
- Recommended next command
