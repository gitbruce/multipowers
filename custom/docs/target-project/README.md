# Target Project User Guide

Audience: users who install this plugin and run `/octo` in their own project.

## Install (User Scope)

```text
/plugin marketplace add /mnt/f/src/ai/claude-octopus
/plugin install octo@nyldn-plugins --scope user
```

## Initialize Once Per Target Project

Open your target project directory, then run:

```text
/octo:init
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
  - `.multipowers/tracks.md`
  - `.multipowers/CLAUDE.md` (project working agreement)
  - `.multipowers/FAQ.md` (auto-generated failure avoidance knowledge)
  - `.multipowers/context/runtime.json` (runtime + pre-run hooks)

## Use Spec-Driven Commands

Examples:
- `/octo:plan <goal>`
- `/octo:discover <topic>`
- `/octo:define <scope>`
- `/octo:develop <implementation>`
- `/octo:deliver <validation>`
- `/octo:embrace <end-to-end request>`

If `.multipowers/` context is missing, spec-driven commands auto-run `/octo:init`.

## Where Outputs Go

- Track plan: `.multipowers/tracks/<track_id>/plan.md`
- Intent contract: `.multipowers/tracks/<track_id>/intent.md`
- Track metadata: `.multipowers/tracks/<track_id>/metadata.json`
- Runtime hooks config: `.multipowers/context/runtime.json`
- Project rules: `.multipowers/CLAUDE.md`
- Auto FAQ: `.multipowers/FAQ.md`

## Update / Remove

```text
/plugin uninstall octo@nyldn-plugins --scope user
/plugin marketplace remove nyldn-plugins
```

## Troubleshooting

- If Claude uses an old plugin cache version, reinstall in user scope.
- If `/octo` behavior looks outdated, run uninstall/install again.
