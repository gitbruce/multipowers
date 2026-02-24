---
command: persona
description: Run a specific pre-configured persona, or list available personas
skill: skill-persona
---

# Persona Command

Run an explicit persona from the Claude Octopus persona catalog.

## Usage

```bash
/octo:persona list
/octo:persona <persona-name> <prompt>
```

## Behavior

- `list`: Shows all personas configured in `agents/config.yaml`.
- `<persona-name> <prompt>`: Runs the prompt with that persona and prints the execution lane in verbose form (for example: `codex:gpt-5.3-codex`) before execution.

## Implementation

```bash
${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh persona "$ARGUMENTS"
```
