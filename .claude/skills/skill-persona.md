---
name: skill-persona
aliases:
  - persona
  - octo:persona
description: Run a specific configured persona, or list available personas
---

# Persona Skill

Routes persona requests to the orchestrator persona subcommand.

## Usage

```bash
/octo:persona list
/octo:persona <persona-name> <prompt>
```

## Execution

```bash
${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh persona "$ARGUMENTS"
```

## Notes

- `list` prints personas from `agents/config.yaml`.
- Persona execution prints the selected model lane before running (for example `codex:gpt-5.3-codex`).
