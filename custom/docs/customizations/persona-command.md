# Persona Command

## What Changed From Upstream
Introduced overlay-managed `/octo:persona` command source and lane display.

## Why This Exists
To keep persona behavior explicit and stable across upstream syncs.

## How To Use
Use `/octo:persona list` and `/octo:persona <name> <prompt>`.

### Persona List Table (with Model IDs)

When listing personas, show lane details using this mapping:

| Persona | Provider Lane | Model ID |
|---|---|---|
| `backend-architect` | `codex` | `gpt-5.3-codex` |
| `code-reviewer` | `codex` | `gpt-5.3-codex` |
| `business-analyst` | `gemini` | `gemini-3-pro-preview` |
| `docs-architect` | `claude_light` | `claude-sonnet` |
| _all others (fallback)_ | `claude_light` | `claude-sonnet` |

## Operational Impact
Persona execution shows `Using: <provider>:<model>` before running.

## Rollback Path
Remove `persona.md` from overlay and command registration.
