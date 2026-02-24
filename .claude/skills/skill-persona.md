---
name: skill-persona
aliases:
  - persona
  - octo:persona
description: Run a specific configured persona, or list available personas
---

# Persona Skill

Routes persona requests to the orchestrator persona subcommand.

## Execution Contract (Mandatory)

When invoked by `/octo:persona`, you MUST:

1. Execute only:

```bash
${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh persona "$ARGUMENTS"
```

2. Do not use Claude Code Task tool subagents (`octo:personas:*`) for this command.
3. Return the orchestrator's explicit execution lane lines, especially:
   - `Using: <provider>:<model>`
   - and when applicable, `Configured: <provider>:<model>`

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
