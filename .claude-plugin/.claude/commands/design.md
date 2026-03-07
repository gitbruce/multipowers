---
command: design
description: "Reuse upstream brainstorming skill for solution design"
skill: mainline-design
---

# /mp:design

REQUIRES /mp:init before entering this flow.

Thin wrapper role: `facilitator`.

Runtime bridge:

`${CLAUDE_PLUGIN_ROOT}/bin/mp design --dir "$PWD" --prompt "$ARGUMENTS" --json`

## Upstream Workflow

---
description: "You MUST use this before any creative work - creating features, building components, adding functionality, or modifying behavior. Explores requirements and design before implementation."
disable-model-invocation: true
---

Invoke the superpowers:brainstorming skill and follow it exactly as presented to you

