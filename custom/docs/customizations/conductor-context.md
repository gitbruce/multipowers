# Conductor Context and /octo:init

## What Changed

Spec-driven `/octo` commands now use project context under `conductor/` in the target project, not `.claude/session-*` files.

## `/octo:init`

Run:

```text
/octo:init
```

Execution anchor:
- `/octo:init` should pass `--dir "$PWD"` to `scripts/orchestrate.sh` so `conductor/` is created in the current target project, not plugin/cache directories.

This initializes:

- `conductor/product.md`
- `conductor/product-guidelines.md`
- `conductor/tech-stack.md`
- `conductor/workflow.md`
- `conductor/code_styleguides/`
- `conductor/tracks.md`

Templates come from `custom/templates/conductor/`.

## Auto-Init Guard

For spec-driven commands (`/octo:plan`, discover/define/develop/deliver, embrace, review, debate, research), orchestrator checks conductor context first.

If missing/incomplete, it auto-runs `/octo:init` interactively before task execution.

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
