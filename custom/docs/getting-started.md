# Getting Started

1. Run `./custom/scripts/apply-custom-overlay.sh`
2. Verify commands: `./scripts/orchestrate.sh persona list`
3. Review customization docs in `customizations/`

## Local Marketplace Lifecycle

### Install marketplace (local folder)

From Claude Code chat:

```text
/plugin marketplace add /mnt/f/src/ai/claude-octopus
```

### Install plugin from that marketplace

```text
/plugin install claude-octopus@nyldn-plugins --scope user
```

### Uninstall plugin

```text
/plugin uninstall claude-octopus@nyldn-plugins --scope user
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
/plugin uninstall claude-octopus@nyldn-plugins --scope user
/plugin install claude-octopus@nyldn-plugins --scope user
```

Then verify:

```text
/octo:persona list
```
