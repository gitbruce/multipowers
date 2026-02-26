---
command: init
description: Initialize .multipowers context via orchestrate.sh interactive wizard (no direct file writing path)
---

# /mp:init

This command MUST invoke the orchestrator wizard. Do not manually generate context files in chat.

## Mandatory Contract

1. Target path is `$PWD/.multipowers`.
2. Execute init only through:

```bash
"${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh" --dir "$PWD" init
```

3. If command exits non-zero:
- Stop immediately and report init failure.
- Do not proceed to any spec-driven task.

4. After success, required files must exist:
- `.multipowers/product.md`
- `.multipowers/product-guidelines.md`
- `.multipowers/tech-stack.md`
- `.multipowers/workflow.md`
- `.multipowers/tracks.md`
- `.multipowers/CLAUDE.md`
- `.multipowers/FAQ.md`
- `.multipowers/context/runtime.json`

## Prohibited

- Do not use `Write/Edit/Update` to create these files directly from command text.
- Do not bypass `orchestrate.sh init`.
- Do not continue to `/mp:plan` or other spec-driven commands when init failed.
