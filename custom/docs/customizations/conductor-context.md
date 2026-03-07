# Conductor Context and /mp:init

## What Changed

Spec-driven `/mp` commands now use project context under `.multipowers/` in the target project, not `.claude/session-*` files.

## `/mp:init`

Run:

```text
/mp:init
```

Execution anchor:
- `/mp:init` should pass `--dir "$PWD"` to `.claude-plugin/bin/mp` so `.multipowers/` is created in the current target project, not plugin/cache directories.
- Runtime guard blocks spec/init commands if `PROJECT_ROOT` resolves to plugin/cache paths.

This initializes:

- `.multipowers/product.md`
- `.multipowers/product-guidelines.md`
- `.multipowers/tech-stack.md`
- `.multipowers/workflow.md`
- `.multipowers/code_styleguides/`
- `.multipowers/context/runtime.json`
- `.multipowers/tracks/tracks.md`

Templates come from `custom/templates/conductor/`.

## Auto-Init Guard

For spec-driven commands (`/mp:plan`, discover/define/develop/deliver, embrace, review, debate, research), orchestrator checks `.multipowers/` context first.

If missing/incomplete, it returns `run_init` guidance before task execution; context files are never generated silently.

## Group Lifecycle Enforcement

Spec-command preflight and implementation-group progress are separate concerns:

- spec commands update `last_command` / `last_command_at` in track metadata
- implementation work starts with `mp track group-start ...`
- implementation work completes with `mp track group-complete ... --commit-sha <sha>`
- while `group_status=in_progress`, the next spec pipeline call is blocked until commit and verification evidence exist

## Track Files

Spec-driven runs create or update canonical track artifacts in:

- `.multipowers/tracks/<track_id>/`

Canonical paths:

- `.multipowers/tracks/tracks.md`
- `.multipowers/tracks/<track_id>/intent.md`
- `.multipowers/tracks/<track_id>/design.md`
- `.multipowers/tracks/<track_id>/implementation-plan.md`
- `.multipowers/tracks/<track_id>/metadata.json`
- `.multipowers/tracks/<track_id>/index.md`

Legacy note:

- `.multipowers/tracks.md` is not compatible with the current runtime and is not read.

## Upstream Borrowing

Setup behavior is adapted from:

- https://github.com/gemini-cli-extensions/conductor

Reference mapping and attribution:

- `custom/references/conductor-upstream/SOURCE-MAP.md`
