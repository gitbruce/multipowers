---
command: status
description: "Inspect current mainline state"
---

# /mp:status

REQUIRES /mp:init before entering this flow.

Thin wrapper role: `reviewer`.

Runtime bridge:

`${CLAUDE_PLUGIN_ROOT}/bin/mp status --dir "$PWD" --prompt "$ARGUMENTS" --json`
