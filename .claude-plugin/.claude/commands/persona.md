---
command: persona
description: Run a specific pre-configured persona, or list available personas
---

# Persona - Explicit Persona Runner

## INSTRUCTIONS FOR CLAUDE

When the user invokes this command (e.g., `/mp:persona <arguments>`):

Preferred invocation is `/mp:persona`.
If this appears as workspace-local `/persona`, treat it as a local alias and still execute with the same runtime path resolution.

### REQUIRED EXECUTION PATH

```bash
MP_BIN=""
if [[ -n "${CLAUDE_PLUGIN_ROOT:-}" ]] && [[ -x "${CLAUDE_PLUGIN_ROOT}/bin/mp" ]]; then
  MP_BIN="${CLAUDE_PLUGIN_ROOT}/bin/mp"
else
  CACHE_MP_BIN="$(ls -td "$HOME/.claude/plugins/cache/multipowers-plugins/mp/"*/bin/mp 2>/dev/null | head -1)"
  if [[ -n "${CACHE_MP_BIN:-}" ]] && [[ -x "${CACHE_MP_BIN}" ]]; then
    MP_BIN="${CACHE_MP_BIN}"
  else
    LEGACY_CACHE_MP_BIN="$(ls -td "$HOME/.claude/plugins/cache/multipowers-plugins/multipowers/"*/.claude-plugin/bin/mp 2>/dev/null | head -1)"
    if [[ -n "${LEGACY_CACHE_MP_BIN:-}" ]] && [[ -x "${LEGACY_CACHE_MP_BIN}" ]]; then
      MP_BIN="${LEGACY_CACHE_MP_BIN}"
    fi
  fi
fi

if [[ -z "${MP_BIN}" ]] && [[ -x "$PWD/.claude-plugin/bin/mp" ]]; then
  MP_BIN="$PWD/.claude-plugin/bin/mp"
elif [[ -z "${MP_BIN}" ]] && [[ -x "./.claude-plugin/bin/mp" ]]; then
  MP_BIN="./.claude-plugin/bin/mp"
fi

if [[ -z "${MP_BIN}" ]]; then
  echo "mp binary not found. Reinstall plugin, restart Claude Code, then retry /mp:* commands." >&2
  echo "Troubleshooting: /plugin install mp@multipowers-plugins --scope user" >&2
  exit 1
fi

"$MP_BIN" persona --dir "$PWD" --prompt "$ARGUMENTS" --json
```

### PROHIBITED

Do NOT run persona requests using Claude Code Task tool subagents.

## Usage

```bash
/mp:persona list
/mp:persona <persona-name> <prompt>
```

## Output Contract

- `list`: one-line table output in `name | description | model` format
- run: includes `using_line` with `Using: <provider>:<model>`
