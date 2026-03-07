# Multipowers

Multipowers is a Go-native Claude Code plugin centered on one public development mainline:

`/mp:init → /mp:brainstorm → /mp:design → /mp:plan → /mp:execute`

Two special direct-entry commands remain available:

- `/mp:debug`
- `/mp:debate`

## Why this branch

The `multi` branch aggressively narrows the product surface so the plugin tells exactly one story:

- `init` establishes required `.multipowers` artifacts
- `brainstorm` and `design` reuse upstream `superpowers` brainstorming guidance
- `plan` reuses upstream planning guidance
- `execute` owns implementation plus finishing the branch
- `debug` is the direct root-cause-first exception path
- `debate` is the multi-model exception path and fans out to all configured models

## Quick Start

1. Install the plugin and build the runtime.
2. Run `/mp:init` in the target project.
3. Move through the mainline:
   - `/mp:brainstorm`
   - `/mp:design`
   - `/mp:plan`
   - `/mp:execute`
4. Use `/mp:debug` only for direct debugging.
5. Use `/mp:debate` when you want all configured models to deliberate together.

## Runtime model

Multipowers keeps Markdown wrappers thin and shifts runtime-specific behavior into Go:

- thin command/skill wrappers under `.claude-plugin/.claude/`
- upstream source-of-truth Markdown synced into `custom/references/superpowers-upstream/`
- generated plugin assets built from local templates plus upstream bodies
- hooks, tracks, policy, and doctor logic implemented in Go

## Build and verify

```bash
go run ./cmd/mp-devx --action sync-superpowers
go run ./cmd/mp-devx --action build-runtime
go test ./internal/devx ./internal/roles ./internal/policy ./internal/providers ./internal/workflows ./internal/hooks ./internal/cli ./internal/validation ./cmd/mp-devx -count=1
```

## Key docs

- `docs/COMMAND-REFERENCE.md`
- `docs/WORKFLOW-SKILLS.md`
- `docs/PLUGIN-ARCHITECTURE.md`
- `docs/CLI-REFERENCE.md`
- `docs/architecture/multi-mainline.md`
- `docs/architecture/superpowers-diff.md`
