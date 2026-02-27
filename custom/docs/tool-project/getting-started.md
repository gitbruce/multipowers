# Getting Started (Tool Project Maintainers)

## Naming Baseline

- Slash command namespace: `/mp:*`
- Plugin id: `mp`
- Marketplace id: `multipowers-plugins`

## Branch Discipline

1. Keep `main` synced to `upstream/main` only.
2. Do all custom development on `multipowers`.
3. Never merge `multipowers` back into `main`.
4. Periodically merge `main` into `multipowers` and sync single-source command docs.

## Quick Start

```bash
git switch main
git fetch upstream
git merge --ff-only upstream/main
git switch multipowers
./scripts/mp-devx overlay
./scripts/mp persona list
```

## Install / Uninstall (User Scope)

Install:

```text
/plugin marketplace add /mnt/f/src/ai/claude-octopus/.claude-plugin/marketplace.json
/plugin install mp@multipowers-plugins --scope user
```

Uninstall:

```text
/plugin uninstall mp@multipowers-plugins --scope user
/plugin marketplace remove multipowers-plugins
```

## Daily Maintainer Workflow

1. Sync upstream with `./scripts/mp-devx sync`
2. Reapply overlay
3. Validate with `go test ./...`
