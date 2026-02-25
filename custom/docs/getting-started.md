# Getting Started

## Branch Discipline (Required)

1. Keep `main` synced to `upstream/main` only.
2. Do all custom development on `multipowers`.
3. Never merge `multipowers` back into `main`.
4. Periodically merge `main` into `multipowers` and reapply overlay.

Quick checks:

```bash
git switch main
git fetch upstream
git merge --ff-only upstream/main
git switch multipowers
```

1. Run `./custom/scripts/apply-custom-overlay.sh`
2. Verify commands: `./scripts/orchestrate.sh persona list`
3. Review customization docs in `customizations/`

## Initialize Project Context

Run once per target project:

```text
/octo:init
```

This generates `conductor/` context files from `custom/templates/conductor/`.

Spec-driven commands auto-check context and auto-run init when missing:
- `/octo:plan`
- `/octo:discover`, `/octo:define`, `/octo:develop`, `/octo:deliver`
- `/octo:embrace`, `/octo:review`, `/octo:debate`, `/octo:research`

## Local Marketplace Lifecycle

### Install marketplace (local folder)

From Claude Code chat:

```text
/plugin marketplace add /mnt/f/src/ai/claude-octopus
```

### Install plugin from that marketplace

```text
/plugin install octo@nyldn-plugins --scope user
```

### Uninstall plugin

```text
/plugin uninstall octo@nyldn-plugins --scope user
```

### Remove local marketplace

```text
/plugin marketplace remove nyldn-plugins
```

## Daily Workflow

1. Sync upstream with `./custom/scripts/sync-upstream.sh`
2. Reapply overlay automatically via sync script
3. Validate with `bash tests/integration/test-sync-overlay.sh`

## If Claude Still Uses Old/Broken Plugin Cache

Claude may run a cached plugin version under `~/.claude/plugins/cache/...`, not your current working tree.

Refresh user-scope install:

```text
/plugin uninstall octo@nyldn-plugins --scope user
/plugin install octo@nyldn-plugins --scope user
```

Then verify:

```text
/octo:persona list
```
