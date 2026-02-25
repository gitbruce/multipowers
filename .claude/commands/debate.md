---
command: debate
description: Structured multi-LLM debate (Codex + Gemini + Claude) via orchestrate grapple workflow
---

# /octo:debate

This command MUST execute the multi-LLM debate workflow. Do not return a single-model opinion.

## Mandatory Behavior

1. Build prompt text from user arguments.
2. Execute debate workflow:

```bash
"${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh" --dir "$PWD" grapple "<user-prompt>"
```

3. If command exits non-zero:
- Stop and report error.
- Do not produce a synthetic single-model "debate" answer.

4. If command succeeds:
- Summarize the debate result from generated artifacts/output.
- Keep attribution clear: include Codex/Gemini/Claude perspectives when present.

## Prohibited

- Do not answer with one model only.
- Do not skip debate execution when user explicitly requests `/octo:debate`.
