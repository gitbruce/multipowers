# Conflict Resolution Runbooks

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
git fetch upstream origin
git switch main
git merge --ff-only upstream/main
git switch multipowers
./custom/scripts/sync-upstream.sh
```

## Manual Verification Checklist
```bash
git status --short --branch
./custom/scripts/sync-upstream.sh
./custom/scripts/apply-custom-overlay.sh
./scripts/orchestrate.sh persona list
```

Expected outcomes:
- no unresolved merge markers
- overlay apply succeeds
- persona list returns configured personas
