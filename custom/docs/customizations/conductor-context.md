# Conductor Context and /octo:init

## What Changed

Spec-driven `/octo` commands now use project context under `conductor/` in the target project, not `.claude/session-*` files.

## `/octo:init`

Run:

```text
/octo:init
```

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

No active spec-driven writes should go to `.claude/session-plan.md` or `.claude/session-intent.md`.

## Upstream Borrowing

Setup behavior is adapted from:

- https://github.com/gemini-cli-extensions/conductor

Reference mapping and attribution:

- `custom/references/conductor-upstream/SOURCE-MAP.md`
