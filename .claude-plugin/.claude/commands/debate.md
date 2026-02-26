---
command: debate
description: Thin wrapper that delegates to Go runtime (octo)
---

# /mp:debate

Use Go runtime only.

Actions:
1. Ensure `${CLAUDE_PLUGIN_ROOT}/bin/mp` exists.
2. Execute:
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" debate --dir "$PWD" --prompt "<user-prompt>" --json
```
3. Parse JSON response.
4. If `status` is `error` or `blocked`, stop immediately.

Do not implement command logic in markdown.
