---
command: plan
description: "Create plan artifacts under conductor/tracks via orchestrate.sh (guarded by /octo:init context check)"
aliases:
  - build-plan
  - intent
---

# /octo:plan

## Mandatory Behavior

1. Build prompt text from user arguments.
2. Before planning, verify required conductor context files exist under `$PWD/conductor/`:
- `product.md`
- `product-guidelines.md`
- `tech-stack.md`
- `workflow.md`
- `tracks.md`
3. If any file is missing:
- Invoke `Skill(octo:init)` first (interactive wizard).
- Do not call `orchestrate.sh plan` until init wizard has completed successfully.
- Re-check the same files; if still missing, stop with error.
- If still missing, output only an initialization failure message and EXIT. Do not ask Goal/Knowledge/Clarity questions.
4. Once context is complete, execute:

```bash
"${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh" --dir "$PWD" plan "<user-prompt>"
```

5. If command exits non-zero:
- Stop immediately and report the error.

6. If command succeeds:
- Return the generated track path and files from command output.
