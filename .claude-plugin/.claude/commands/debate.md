---
command: debate
description: "Run multi-model debate using all configured models"
skill: mainline-debate
---

# /mp:debate

REQUIRES /mp:init before entering this flow.

Thin wrapper role: `debater`.

Runtime bridge:

`${CLAUDE_PLUGIN_ROOT}/bin/mp debate --dir "$PWD" --prompt "$ARGUMENTS" --json`
