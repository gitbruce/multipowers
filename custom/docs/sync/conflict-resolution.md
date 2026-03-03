# Conflict Resolution Runbooks

## First Rule

- Keep `main` unchanged except syncing from `upstream/main`.
- Resolve customization conflicts only in `go`.
- Never run sync by switching current worktree branch.
- Use temporary worktrees rooted at `.worktrees/sync-*` for sync actions.
- Never resolve by revert/reset of local uncommitted files.

## Resolve Rebase Path
1. Check status:
```bash
git status --short --branch
```
2. Resolve conflicted files manually.
3. Mark resolved files:
```bash
git add <resolved-file-1> <resolved-file-2>
```
4. Continue rebase:
```bash
git rebase --continue
```
5. Repeat until complete.

## Abort Rebase Path
Use this when conflict resolution exceeds the working SLA.

```bash
git rebase --abort
git status --short --branch
```

Then recover using sync flow:
```bash
./scripts/sync-upstream-main.sh -dry-run
./scripts/sync-main-to-go.sh -dry-run
./scripts/sync-all.sh -dry-run
```

## Manual Verification Checklist
```bash
git status --short --branch
./scripts/sync-upstream-main.sh -dry-run
./scripts/sync-main-to-go.sh -dry-run
./scripts/sync-all.sh -dry-run
./.claude-plugin/bin/mp persona list
```

Expected outcomes:
- no unresolved merge markers
- sync dry-run commands return success
- persona list returns configured personas

## Conductor Path Canonicalization

Use only:
- `conductor/tracks.md`
- `conductor/tracks/`

Avoid legacy or typo variants like `conductor/track/`.
