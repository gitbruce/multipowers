---
command: resume
description: "Resume the last interrupted mainline flow"
---

# /mp:resume

REQUIRES /mp:init before entering this flow.

Thin wrapper role: `reviewer`.

Runtime bridge:

`${CLAUDE_PLUGIN_ROOT}/bin/mp resume --dir "$PWD" --prompt "$ARGUMENTS" --json`
