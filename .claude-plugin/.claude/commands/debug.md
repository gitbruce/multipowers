---
command: debug
description: "Enter direct debugging flow"
skill: mainline-debug
---

# /mp:debug

REQUIRES /mp:init before entering this flow.

Thin wrapper role: `debugger`.

Runtime bridge:

`${CLAUDE_PLUGIN_ROOT}/bin/mp debug --dir "$PWD" --prompt "$ARGUMENTS" --json`
