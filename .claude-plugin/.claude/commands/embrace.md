---
command: embrace
description: Thin wrapper that delegates to Go runtime (mp)
---

# /mp:embrace

Use Go runtime only.

Actions:
1. Resolve mp runtime path (plugin first, workspace fallback).
2. Execute:
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

"$MP_BIN" embrace --dir "$PWD" --prompt "<user-prompt>" --json
```
3. Parse JSON response.
4. If `status` is `error` or `blocked`, stop immediately.

Do not implement workflow logic in markdown.
