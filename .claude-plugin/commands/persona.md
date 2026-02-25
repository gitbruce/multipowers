---
command: persona
description: Run a specific pre-configured persona, or list available personas
skill: skill-persona
---

# Persona - Explicit Persona Runner

## INSTRUCTIONS FOR CLAUDE

When the user invokes this command (e.g., `/octo:persona <arguments>`):

### ✓ REQUIRED EXECUTION PATH

Always execute via orchestrate CLI:

```bash
${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh persona "$ARGUMENTS"
```

### ✗ PROHIBITED

Do NOT run persona requests using Claude Code Task tool subagents (for example `octo:personas:*`), because that path hides the provider:model lane and breaks `/octo:persona` contract.

### REQUIRED OUTPUT

Surface the orchestrator output, including the explicit model lane line:

- `Using: <provider>:<model>`
- If fallback occurs in Claude Code nested session protection:
  - `Configured: <provider>:<model>`
  - `Using: <provider>:<model>`

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
