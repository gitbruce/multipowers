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
- creates `conductor/` in your current target project directory
- generates:
  - `conductor/product.md`
  - `conductor/product-guidelines.md`
  - `conductor/tech-stack.md`
  - `conductor/workflow.md`
  - `conductor/code_styleguides/`
  - `conductor/tracks.md`

## Use Spec-Driven Commands

Examples:
- `/octo:plan <goal>`
- `/octo:discover <topic>`
- `/octo:define <scope>`
- `/octo:develop <implementation>`
- `/octo:deliver <validation>`
- `/octo:embrace <end-to-end request>`

If `conductor/` context is missing, spec-driven commands auto-run `/octo:init`.

## Where Outputs Go

- Track plan: `conductor/tracks/<track_id>/plan.md`
- Intent contract: `conductor/tracks/<track_id>/intent.md`
- Track metadata: `conductor/tracks/<track_id>/metadata.json`

## Update / Remove

```text
/plugin uninstall octo@nyldn-plugins --scope user
/plugin marketplace remove nyldn-plugins
```

## Troubleshooting

- If Claude uses an old plugin cache version, reinstall in user scope.
- If `/octo` behavior looks outdated, run uninstall/install again.
