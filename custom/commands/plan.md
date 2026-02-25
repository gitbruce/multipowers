---
command: plan
description: "Create plan artifacts under conductor/tracks via orchestrate.sh (guarded by /octo:init context check)"
aliases:
  - build-plan
  - intent
---

# /octo:plan

This command is intentionally thin and delegates execution to the shell runtime so guard logic is deterministic.

## Mandatory Behavior

1. Build prompt text from user arguments.
2. Execute:

```bash
"${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh" --dir "$PWD" plan "<user-prompt>"
```

3. If command exits non-zero:
- Stop immediately.
- Report the error.
- Do not continue with interactive intent questions.

4. If command succeeds:
- Return the generated track path and files from command output.

## Notes

- `orchestrate.sh plan` enforces spec-driven conductor context checks.
- If required context files are missing, it must run `/octo:init` first and only continue when context is complete.
