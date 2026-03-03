---
command: init
description: Guided wizard setup for .multipowers context
---

# /mp:init

Run a guided wizard first, then call Go runtime.

Actions:
1. Ask one batched `AskUserQuestion` with these fields:
   - `project_name`
   - `summary`
   - `target_users`
   - `primary_goal`
   - `non_goals`
   - `constraints`
   - `runtime`
   - `framework`
   - `database`
   - `deployment`
   - `workflow`
   - `track_name`
   - `track_objective`
2. Build a single-line JSON object from answers and pass it as `--prompt`.
3. Resolve mp runtime path (plugin first, workspace fallback).
4. Execute:
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

"$MP_BIN" init --dir "$PWD" --prompt "$INIT_WIZARD_JSON" --json
```
5. Parse JSON response.
6. If `status` is `error` or `blocked`, stop immediately.

Do not implement command logic in markdown.
