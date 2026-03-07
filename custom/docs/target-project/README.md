# Target Project User Guide

Audience: users who install this plugin and run `/mp` in their own project.

## Install (User Scope)

```text
/plugin marketplace add /mnt/f/src/ai/multipowers/.claude-plugin/marketplace.json
/plugin install mp@multipowers-plugins --scope user
```

## Initialize Once Per Target Project

Open your target project directory, then run:

```text
/mp:init
```

Expected behavior:
- uses Conductor setup protocol from `custom/config/setup.toml`
- creates `.multipowers/` in your current target project directory
- generates:
  - `.multipowers/product.md`
  - `.multipowers/product-guidelines.md`
  - `.multipowers/tech-stack.md`
  - `.multipowers/workflow.md`
  - `.multipowers/code_styleguides/`
  - `.multipowers/tracks/tracks.md`
  - `.multipowers/CLAUDE.md` (project working agreement)
  - `.multipowers/FAQ.md` (auto-generated failure avoidance knowledge)
  - `.multipowers/context/runtime.json` (runtime + pre-run hooks; `pre_run.enabled=false` by default)
  - `.multipowers/tracks/<track_id>/` (spec-driven artifacts)

## Use Spec-Driven Commands

Examples:
- `/mp:plan <goal>`
- `/mp:discover <topic>`
- `/mp:define <scope>`
- `/mp:develop <implementation>`
- `/mp:deliver <validation>`
- `/mp:embrace <end-to-end request>`

If `.multipowers/` context is missing or incomplete, spec-driven commands stop with `run_init` guidance; runtime files are never generated silently.

## Explicit Group Lifecycle

After a spec track exists, implementation groups are advanced explicitly:

```text
mp track group-start --track-id <track_id> --group g1 --execution-mode workspace --json
mp track group-complete --track-id <track_id> --group g1 --commit-sha <sha> --json
```

Rules:
- use `.multipowers/tracks/<track_id>/` for all spec-driven artifacts
- finish each active group with commit + verification evidence before the next spec pipeline step
- if a track is marked worktree-required, `group-start` must run from a linked git worktree checkout

## Where Outputs Go

- Tracks registry: `.multipowers/tracks/tracks.md`
- Intent contract: `.multipowers/tracks/<track_id>/intent.md`
- Design doc: `.multipowers/tracks/<track_id>/design.md`
- Implementation plan: `.multipowers/tracks/<track_id>/implementation-plan.md`
- Track metadata: `.multipowers/tracks/<track_id>/metadata.json`
- Track index: `.multipowers/tracks/<track_id>/index.md`
- Runtime hooks config: `.multipowers/context/runtime.json`
- Project rules: `.multipowers/CLAUDE.md`
- Auto FAQ: `.multipowers/FAQ.md`

Legacy note:
- `.multipowers/tracks.md` is no longer read by the runtime. Use `.multipowers/tracks/tracks.md` only.

## Update / Remove

```text
/plugin uninstall mp@multipowers-plugins --scope user
/plugin marketplace remove multipowers-plugins
```

## Troubleshooting

- If Claude uses an old plugin cache version, reinstall in user scope.
- If `/mp` behavior looks outdated, run uninstall/install again.
