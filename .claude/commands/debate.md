---
command: debate
description: Structured multi-LLM debate (Codex + Gemini + Claude) via orchestrate grapple workflow
---

# /octo:debate

This command MUST execute the multi-LLM debate workflow. Do not return a single-model opinion.

## Mandatory Behavior

1. Before debate, verify required context exists under `$PWD/.multipowers/`:
   - `product.md`
   - `product-guidelines.md`
   - `tech-stack.md`
   - `workflow.md`
   - `tracks.md`
   - `CLAUDE.md`
   If any file is missing, you MUST execute:
```bash
"${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh" --dir "$PWD" init
```
   Re-check all required files and continue only when all are present; otherwise fail-fast.
2. Build prompt text from user arguments.
3. Execute debate workflow:

```bash
"${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh" --dir "$PWD" grapple "<user-prompt>"
```

4. If command exits non-zero:
- Stop and report error.
- Do not produce a synthetic single-model "debate" answer.

5. If command succeeds:
- Summarize the debate result from generated artifacts/output.
- Keep attribution clear: include Codex/Gemini/Claude perspectives when present.

## Prohibited

- Do not answer with one model only.
- Do not skip debate execution when user explicitly requests `/octo:debate`.
