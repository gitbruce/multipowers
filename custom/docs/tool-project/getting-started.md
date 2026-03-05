# Getting Started (Tool Project Maintainers)

## Naming Baseline

- Slash command namespace: `/mp:*`
- Plugin id: `mp`
- Marketplace id: `multipowers-plugins`

## Branch Discipline (Go Development)

1. Keep `main` synced to `upstream/main` only (as a clean mirror).
2. All **Go-based no-shell runtime development** occurs on the `go` branch.
3. Sync upstream with `./scripts/sync-all.sh`.
4. Run validation with `./scripts/validate-claude-structure.sh`.

## Maintenance Commands

```bash
# Verify the no-shell runtime and logic
go test ./...

# Rebuild the main CLI
go build -o .claude-plugin/bin/mp ./cmd/mp

# Build the DevX helper
go build -o .claude-plugin/bin/mp-devx ./cmd/mp-devx

# Run an atomic command directly for testing
./.claude-plugin/bin/mp status
```

## Install / Uninstall (User Scope)

Install:

```text
/plugin marketplace add /mnt/f/src/ai/multipowers/.claude-plugin/marketplace.json
/plugin install mp@multipowers-plugins --scope user
```

Uninstall:

```text
/plugin uninstall mp@multipowers-plugins --scope user
/plugin marketplace remove multipowers-plugins
```

## Daily Maintainer Workflow

1. Sync upstream changes via `scripts/sync-upstream-main.sh` and `scripts/sync-main-to-go.sh`.
2. Implement Go packages in `internal/` or CLI commands in `cmd/`.
3. Validate parity and architecture with `scripts/verify-architecture-diff-docs.sh`.
4. Ensure all Go tests pass: `go test ./... -v`.
