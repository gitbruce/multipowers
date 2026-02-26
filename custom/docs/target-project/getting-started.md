# Getting Started (Target Project Users)

## Install Plugin (User Scope)

```text
/plugin marketplace add /mnt/f/src/ai/claude-octopus
/plugin install octo@nyldn-plugins --scope user
```

## Initialize in Your Project

In your target project directory:

```text
/octo:init
```

Expected:
- creates `.multipowers/` in your target project
- initializes project context files and tracks registry
- creates `.multipowers/CLAUDE.md` and `.multipowers/FAQ.md`
- creates `.multipowers/context/runtime.json` for runtime/pre-run settings

## Optional: Configure Pre-Run Hooks

During `/octo:init`, you can configure pre-run hooks (for any runtime/toolchain).

Example use cases:
- activate an environment
- source a shell profile
- export required env vars

Hooks are stored in:
- `.multipowers/context/runtime.json`

All external provider executions (`codex`, `gemini`, `claude`) read this file before running.

## Run Spec-Driven Commands

- `/octo:plan`
- `/octo:discover`, `/octo:define`, `/octo:develop`, `/octo:deliver`
- `/octo:embrace`, `/octo:review`, `/octo:debate`, `/octo:research`

If context is missing, `/octo:init` is auto-triggered.

## FAQ Learning Loop

- `.multipowers/FAQ.md` is auto-generated and auto-refined from observed failures.
- Entries are categorized by error type and deduplicated.
- No manual maintenance, backup, or archive files are required.
