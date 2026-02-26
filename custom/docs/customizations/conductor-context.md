# Conductor Context and /octo:init

## What Changed

Spec-driven `/octo` commands now use project context under `.multipowers/` in the target project, not `.claude/session-*` files.

## `/octo:init`

Run:

```text
/octo:init
```

Execution anchor:
- `/octo:init` should pass `--dir "$PWD"` to `scripts/orchestrate.sh` so `.multipowers/` is created in the current target project, not plugin/cache directories.
- Runtime guard blocks spec/init commands if `PROJECT_ROOT` resolves to plugin/cache paths.

This initializes:

- `.multipowers/product.md`
- `.multipowers/product-guidelines.md`
- `.multipowers/tech-stack.md`
- `.multipowers/workflow.md`
- `.multipowers/code_styleguides/`
- `.multipowers/tracks.md`

Templates come from `custom/templates/.multipowers/`.

## Auto-Init Guard

For spec-driven commands (`/octo:plan`, discover/define/develop/deliver, embrace, review, debate, research), orchestrator checks conductor context first.

If missing/incomplete, it auto-runs `/octo:init` interactively before task execution.

## Track Files

Spec-driven runs create or update checkbox-based tracking files in:

- `.multipowers/tracks/`

Canonical paths:

- `.multipowers/tracks.md`
- `.multipowers/tracks/`

Spec-driven outputs are stored at `.multipowers/tracks/<track_id>/plan.md` and `.multipowers/tracks/<track_id>/intent.md`.

## Upstream Borrowing

Setup behavior is adapted from:

- https://github.com/gemini-cli-extensions/conductor

Reference mapping and attribution:

- `custom/references/conductor-upstream/SOURCE-MAP.md`
