---
command: execute
description: "Execute the implementation plan with upstream execution guidance"
skill: mainline-execute
---

# /mp:execute

REQUIRES /mp:init before entering this flow.

Thin wrapper role: `executor`.

Runtime bridge:

`${CLAUDE_PLUGIN_ROOT}/bin/mp execute --dir "$PWD" --prompt "$ARGUMENTS" --json`

## Upstream Workflow

---
description: Execute plan in batches with review checkpoints
disable-model-invocation: true
---

Invoke the superpowers:executing-plans skill and follow it exactly as presented to you

