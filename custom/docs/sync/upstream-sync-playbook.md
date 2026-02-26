# Upstream Sync Playbook

## Branch Principle

- `main` must remain a clean mirror of `upstream/main`.
- `multipowers` is the only customization branch.
- Sync direction is one-way: `upstream/main -> main -> multipowers`.
- Minimize edits in high-conflict upstream files (`bin/mp`, `.claude-plugin/.claude/*`, `.claude-plugin/*`); prefer `custom/*`.

## Routine Sync Sequence

```bash
git fetch upstream origin
git switch main
git merge --ff-only upstream/main
git switch multipowers
git merge main -m "chore(sync): merge main into multipowers"
./custom/scripts/mp-devx overlay
go test ./...
```

## Pre-Sync Guards

```bash
git status --short --branch
git switch main
git merge --ff-only upstream/main
git switch multipowers
```

Required outcomes:
- working tree is clean before sync
- `main` fast-forwards to `upstream/main`
- sync happens by merging `main` into `multipowers`

## Conflict SLA and Fallback

- Target: resolve sync conflicts within 30 minutes.
- If not resolved within SLA:

```bash
git rebase --abort
# or git merge --abort, depending on operation
```

Then restart from clean baseline using this playbook.

## Example Sync Transcript

See: `custom/docs/sync/verification-transcript.md`

## Expected Result

- `main` matches `upstream/main`
- overlay reapplied successfully
- sync/registration tests pass

## Conductor Source Reference

Conductor-style setup behavior used by `/mp:init` is tracked here:
- `custom/references/conductor-upstream/SOURCE-MAP.md`
