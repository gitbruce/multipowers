---
command: plan
description: "Build an implementation plan from upstream planning docs"
skill: mainline-plan
---

# /mp:plan

REQUIRES /mp:init before entering this flow.

Thin wrapper role: `planner`.

Runtime bridge:

`${CLAUDE_PLUGIN_ROOT}/bin/mp plan --dir "$PWD" --prompt "$ARGUMENTS" --json`

## Upstream Workflow

---
description: Create detailed implementation plan with bite-sized tasks
disable-model-invocation: true
---

Invoke the superpowers:writing-plans skill and follow it exactly as presented to you

