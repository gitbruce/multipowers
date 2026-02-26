---
command: persona
description: Run a specific pre-configured persona, or list available personas
skill: skill-persona
---

# Persona - Explicit Persona Runner

## INSTRUCTIONS FOR CLAUDE

When the user invokes this command (e.g., `/octo:persona <arguments>`):

### REQUIRED EXECUTION PATH

```bash
${CLAUDE_PLUGIN_ROOT}/bin/octo persona --dir "$PWD" --prompt "$ARGUMENTS" --json
```

### PROHIBITED

Do NOT run persona requests using Claude Code Task tool subagents.

## Usage

```bash
/octo:persona list
/octo:persona <persona-name> <prompt>
```

## Output Contract

- `list`: one-line table output in `name | description | model` format
- run: includes `using_line` with `Using: <provider>:<model>`
