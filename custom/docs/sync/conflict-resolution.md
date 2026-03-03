# Conflict Resolution Runbooks

## Safety Baseline

- Do not revert local uncommitted files.
- Do not run destructive reset commands.
- Resolve sync issues only in isolated sync worktrees.

## A. Git Conflict During Sync

1. Inspect state:
```bash
git status --short --branch
```
2. Resolve conflicted files manually.
3. Mark resolved files:
```bash
git add <resolved-files>
```
4. Continue the operation (`rebase --continue` or merge continue path).
5. Re-run dry-run checks:
```bash
./scripts/sync-upstream-main.sh -dry-run
./scripts/sync-main-to-go.sh -dry-run
./scripts/sync-all.sh -dry-run
```

## B. Structure Parity Validation Failure

1. Reproduce:
```bash
./scripts/validate-claude-structure.sh -dry-run
```
2. For unexpected drift, sync from `main` again and re-check.
3. For intentional divergence, update explicit allow/ignore entries in:
- `config/sync/claude-structure-rules.json`
4. Re-run:
```bash
./scripts/validate-claude-structure.sh -dry-run
go test ./internal/devx ./cmd/mp-devx -v
```

## C. Abort and Retry Path

Use this when conflict resolution exceeds SLA.

```bash
git rebase --abort || true
git status --short --branch
./scripts/sync-all.sh -dry-run
./scripts/validate-claude-structure.sh -dry-run
```

## Manual Verification Checklist

```bash
git status --short --branch
./scripts/sync-upstream-main.sh -dry-run
./scripts/sync-main-to-go.sh -dry-run
./scripts/sync-all.sh -dry-run
./scripts/validate-claude-structure.sh -dry-run
```
