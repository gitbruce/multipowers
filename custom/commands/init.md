---
command: init
description: Initialize project conductor context files in the current target project
---

# Init - Project Context Bootstrap

When user runs `/octo:init`, always execute:

```bash
${CLAUDE_PLUGIN_ROOT}/scripts/orchestrate.sh init
```

Behavior:
- Interactive setup wizard (conductor-style)
- Creates/updates project files under `conductor/`
- Uses templates maintained in `custom/templates/conductor/`

Artifacts:
- `conductor/product.md`
- `conductor/product-guidelines.md`
- `conductor/tech-stack.md`
- `conductor/workflow.md`
- `conductor/code_styleguides/`
- `conductor/tracks.md`
