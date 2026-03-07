# Workflow Skills

The public workflow surface is intentionally narrow.

## Mainline

The default path is:

`/mp:init → /mp:brainstorm → /mp:design → /mp:plan → /mp:execute`

### `mainline-brainstorm`
- wraps upstream `skills/brainstorming/SKILL.md`
- used by `/mp:brainstorm`
- may fan out across configured models for early exploration

### `mainline-design`
- reuses the same upstream brainstorming skill
- used by `/mp:design`
- turns exploration into a solution direction

### `mainline-plan`
- wraps upstream `skills/writing-plans/SKILL.md`
- used by `/mp:plan`
- produces the executable implementation plan

### `mainline-execute`
- wraps upstream `skills/executing-plans/SKILL.md`
- used by `/mp:execute`
- covers implementation, review checkpoints, and branch finishing

## Special entries

### `mainline-debug`
- wraps upstream `skills/systematic-debugging/SKILL.md`
- used by `/mp:debug`
- bypasses the normal design/planning path when direct debugging is appropriate

### `mainline-debate`
- thin local wrapper for multi-model deliberation
- used by `/mp:debate`
- always fans out to all configured models

## Design principle

The wrapper Markdown stays intentionally small. Same-function workflow text is sourced from upstream `superpowers`, while init gates, model policy, hooks, and Go runtime behavior stay local.
