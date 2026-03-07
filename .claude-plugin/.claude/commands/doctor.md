---
command: doctor
description: "Diagnose mainline environment readiness"
---

# /mp:doctor

REQUIRES /mp:init before entering this flow.

Thin wrapper role: `reviewer`.

Runtime bridge:

`${CLAUDE_PLUGIN_ROOT}/bin/mp doctor --dir "$PWD" --prompt "$ARGUMENTS" --json`
