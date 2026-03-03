# Upstream Sync Playbook

## Branch Principle

- `main` must remain a clean mirror of `upstream/main`.
- `go` is the only customization branch.
- Sync direction is one-way: `upstream/main -> main -> go`.
- Never run sync by switching current worktree branch.
- Run branch mutation only in temporary worktrees under `.worktrees/sync-*`.
- Never resolve sync by revert/reset of local uncommitted files.
- Minimize edits in high-conflict upstream files (`.claude-plugin/bin/mp`, `.claude-plugin/.claude/*`, `.claude-plugin/*`); prefer `custom/*`.

## Routine Sync Sequence

```bash
./scripts/sync-upstream-main.sh -dry-run
./scripts/sync-main-to-go.sh -dry-run
./scripts/sync-all.sh -dry-run
```

## Pre-Sync Guards

```bash
git status --short --branch
git fetch upstream origin --prune
```

Required outcomes:
- no branch switch in the active developer worktree
- sync branch mutations execute only in `.worktrees/sync-*`
- `main` fast-forwards from `upstream/main`
- `COPY_FROM_MAIN` rules are applied from `main` into `go`

## Conflict SLA and Fallback

- Target: resolve sync conflicts within 30 minutes.
- If not resolved within SLA:

```bash
git worktree list
# remove failed temp worktree and rerun scripts
```

Then restart from clean baseline using this playbook.

## Example Sync Transcript

See: `custom/docs/sync/verification-transcript.md`

## Expected Result

- `main` matches `upstream/main`
- `go` receives only allowed shared-file sync via rules
- dry-run and automation checks pass

## Conductor Source Reference

Conductor-style setup behavior used by `/mp:init` is tracked here:
- `custom/references/conductor-upstream/SOURCE-MAP.md`
