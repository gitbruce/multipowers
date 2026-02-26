# Conflict Resolution Runbooks

## First Rule

- Keep `main` unchanged except syncing from `upstream/main`.
- Resolve customization conflicts only in `multipowers`.

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
./custom/scripts/octo-devx sync
```

## Manual Verification Checklist
```bash
git status --short --branch
./custom/scripts/octo-devx sync
./custom/scripts/octo-devx overlay
./bin/octo persona list
```

Expected outcomes:
- no unresolved merge markers
- overlay apply succeeds
- persona list returns configured personas

## Conductor Path Canonicalization

Use only:
- `conductor/tracks.md`
- `conductor/tracks/`

Avoid legacy or typo variants like `conductor/track/`.
