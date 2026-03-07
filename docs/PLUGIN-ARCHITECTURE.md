# Plugin Architecture

## Public surface

The published plugin surface is generated from `custom/config/mainline-surface.yaml` and contains only the approved mainline commands and wrapper skills.

## Source of truth

Multipowers separates prompt content from runtime glue:

- upstream workflow Markdown is synced into `custom/references/superpowers-upstream/`
- thin local templates live in `custom/templates/mainline-wrapper/`
- generated assets are written to `.claude-plugin/.claude/`
- the final public contract is recorded in `.claude-plugin/plugin.json`

## Runtime layers

1. Wrapper Markdown exposes the public command and skill surface.
2. Go runtime enforces init gating, tracks, hooks, policy, and doctor behavior.
3. Policy resolves models/providers for each phase.
4. Debate forces all configured providers into the run.

## Fixed roles

Persona selection is no longer part of the public API. Internally, the branch uses a fixed role set:

- `initializer`
- `facilitator`
- `planner`
- `executor`
- `reviewer`
- `debugger`
- `debater`

## Build flow

```bash
go run ./cmd/mp-devx --action sync-superpowers
go run ./cmd/mp-devx --action build-runtime
```
