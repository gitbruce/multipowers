# Multi Mainline Architecture

## Goal

This branch reduces Multipowers to one public development path:

`init → brainstorm → design → plan → execute`

with `debug` and `debate` as the only direct-entry exceptions.

## Upstream relationship

Workflow text for same-function behaviors is sourced from upstream `superpowers`:

- command Markdown sync manifest: `custom/config/superpowers-sync.yaml`
- synced upstream files: `custom/references/superpowers-upstream/`
- generated wrappers: `custom/templates/mainline-wrapper/`

## Thin wrapper model

The generated command/skill Markdown should only describe:

- init prerequisite behavior
- runtime bridge to the Go CLI
- model/debate specifics that are unique to Multipowers

Everything else should come from upstream bodies.

## Fixed-role runtime

Internal routing uses fixed roles instead of selectable personas:

- `initializer`
- `facilitator`
- `planner`
- `executor`
- `reviewer`
- `debugger`
- `debater`

These roles are declared in `config/roles.yaml` and resolved via `internal/roles`.

## Model policy

- `brainstorm`, `design`, and `plan` may use multiple configured models in parallel
- `execute` and `debug` use phase-specific runtime policy
- `debate` always uses the full configured model/provider set

## Init gate

`/mp:init` is a hard prerequisite for all mainline and special commands. Missing init returns a blocked response, points the user to `/mp:init`, and preserves resume metadata.
