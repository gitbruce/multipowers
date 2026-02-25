---
command: init
description: Initialize project context using exact Conductor setup protocol from custom/config/setup.toml
---

# Init - Conductor Setup (Exact Upstream Protocol)

When user runs `/octo:init`, do not use adapted instructions.
Load and execute the protocol from:

- `custom/config/setup.toml`

Rules:
- Follow `custom/config/setup.toml` verbatim as the source of truth.
- Treat this file as the exact upstream Conductor setup logic.
- Keep generated project context under `conductor/` in the target project.
- Execute orchestration with explicit target directory:
  - `\${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh --dir "$PWD" init`

Reference:
- Upstream source: `https://github.com/gemini-cli-extensions/conductor/blob/main/commands/conductor/setup.toml`
