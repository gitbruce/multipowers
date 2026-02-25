---
command: init
description: Initialize project context using Conductor setup protocol adapted for /octo:init
---

# Init - Conductor Setup (Octo-Compatible Protocol)

When user runs `/octo:init`, do not use adapted instructions.
Load and execute the protocol from:

- `custom/config/setup.toml`

Rules:
- Follow `custom/config/setup.toml` as the source of truth.
- This file is based on upstream Conductor setup protocol with `/octo:init` compatibility adjustments.
- Keep generated project context under `conductor/` in the target project.
- Use interactive setup flow driven by `AskUserQuestion` as defined in `custom/config/setup.toml`.

Prohibited:
- Do NOT shortcut initialization by directly running:
  - `\${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh --dir "$PWD" init`
- Do NOT continue to planning/execution if the setup wizard has not completed.

Reference:
- Upstream source: `https://github.com/gemini-cli-extensions/conductor/blob/main/commands/conductor/setup.toml`
