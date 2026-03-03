# Upstream Sync Playbook

## Core Principles

- Sync direction is fixed: `upstream/main -> main -> go`.
- `main` stays a clean mirror of `upstream/main`; implementation continues on `go`.
- No overlay mechanism: common content enters `go` only through `COPY_FROM_MAIN` rules.
- All sync mutations run in isolated `.worktrees/sync-*` worktrees, never by switching the current working branch.

## Rules Contracts

- File copy contract: `config/sync/main-to-go-rules.json`
- Structure parity contract: `config/sync/claude-structure-rules.json`
- Validation entrypoint: `./scripts/validate-claude-structure.sh`

## Optional Proxy (If Git Is Slow)

```bash
host_ip=$(ip route show | grep -i default | awk '{print $3}')
export http_proxy="http://$host_ip:7890"
export https_proxy="http://$host_ip:7890"
```

## Dry-Run Sequence

```bash
./scripts/sync-upstream-main.sh -dry-run
./scripts/sync-main-to-go.sh -dry-run
./scripts/sync-all.sh -dry-run
./scripts/validate-claude-structure.sh -dry-run
```

## Apply Sequence

```bash
./scripts/sync-upstream-main.sh
./scripts/sync-main-to-go.sh
./scripts/validate-claude-structure.sh
go test ./internal/devx ./cmd/mp-devx -v
```

## Required Outcomes

- `main` fast-forwards to `upstream/main` in isolated worktree.
- `go` receives only rules-allowed shared files.
- `.claude-plugin/.claude` structure checks pass for `MUST_HOMOMORPHIC` scopes.
- No local uncommitted user edits are reverted.

## Failure Handling

- If sync or validation fails, follow: `custom/docs/sync/conflict-resolution.md`
- Keep evidence in: `custom/docs/sync/verification-transcript.md`
