# Conductor Context and /mp:init

## What Changed

Spec-driven `/mp` commands now use project context under `conductor/` in the target project, not `.claude/session-*` files.

## `/mp:init`

Run:

```text
/mp:init
```

Execution anchor:
- `/mp:init` should pass `--dir "$PWD"` to `bin/mp` so `conductor/` is created in the current target project, not plugin/cache directories.
- Runtime guard blocks spec/init commands if `PROJECT_ROOT` resolves to plugin/cache paths.

This initializes:

- `conductor/product.md`
- `conductor/product-guidelines.md`
- `conductor/tech-stack.md`
- `conductor/workflow.md`
- `conductor/code_styleguides/`
- `conductor/tracks.md`

Templates come from `custom/templates/conductor/`.

## Auto-Init Guard

For spec-driven commands (`/mp:plan`, discover/define/develop/deliver, embrace, review, debate, research), orchestrator checks conductor context first.

If missing/incomplete, it auto-runs `/mp:init` interactively before task execution.

## Track Files

Spec-driven runs create or update checkbox-based tracking files in:

- `conductor/tracks/`

Canonical paths:

- `conductor/tracks.md`
- `conductor/tracks/`

Spec-driven outputs are stored at `conductor/tracks/<track_id>/plan.md` and `conductor/tracks/<track_id>/intent.md`.

## Upstream Borrowing

Setup behavior is adapted from:

- https://github.com/gemini-cli-extensions/conductor

Reference mapping and attribution:

- `custom/references/conductor-upstream/SOURCE-MAP.md`
