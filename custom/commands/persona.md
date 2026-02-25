---
command: persona
description: Run a specific pre-configured persona, or list available personas
skill: skill-persona
---

# Persona - Explicit Persona Runner

## INSTRUCTIONS FOR CLAUDE

When the user invokes this command (e.g., `/octo:persona <arguments>`):

### REQUIRED EXECUTION PATH

Always execute via orchestrate CLI:

```bash
${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh persona "$ARGUMENTS"
```

### PROHIBITED

Do NOT run persona requests using Claude Code Task tool subagents.

### REQUIRED OUTPUT

Surface the orchestrator output, including the explicit model lane line:

- `Using: <provider>:<model>`

# Persona Command

Run an explicit persona from the Claude Octopus persona catalog.

## Usage

```bash
/octo:persona list
/octo:persona <persona-name> <prompt>
```

## Behavior

- `list`: Shows all personas configured in `agents/config.yaml`.
- `<persona-name> <prompt>`: Runs the prompt with that persona and prints the execution lane.

## Implementation

```bash
${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh persona "$ARGUMENTS"
```
