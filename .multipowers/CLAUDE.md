# claude-octopus-go - Claude Working Agreement

## Project Context
- Summary: No-shell hybrid runtime with deterministic Go atomic engine plus Markdown reasoning orchestration.
- Target users: `go` branch maintainers and migration implementers.
- Primary goal: Stable command contracts and full upstream script mapping traceability.
- Constraints: Keep `main` untouched, sync upstream regularly, preserve normalized response schema.

## Runtime Profile
- Runtime: Go 1.21+
- Framework: standard library + internal domain packages
- Database: repository-managed file/state artifacts
- Deployment: commit and push to `origin/go` after verification

## Execution Rules
- Keep project context artifacts under `.multipowers/`.
- Do not implement migration work on `main` branch.
- Before major updates: `git fetch upstream --prune` and confirm baseline refs.
- Validate before completion claims; include concrete command evidence.
- Keep mapping rows explicit: strategy, target path, target symbol, status.

## Upstream Sync Policy
- `main` is read-only in local workflow: no feature commits, no direct edits, no migration work.
- Periodically sync upstream into `main` to keep baseline current (for example before each major planning/implementation wave).
- Use `main` only as canonical upstream baseline reference; active development happens on `go`.
- Sync flow is: `upstream/main -> main -> go`.
- Never run sync by switching current worktree branch.
- Run sync branch mutations only in `.worktrees/sync-*`.
- Never resolve sync with revert/reset of local uncommitted files.
- Preferred entrypoints: `./scripts/sync-upstream-main.sh`, `./scripts/sync-main-to-go.sh`, `./scripts/sync-all.sh`.

## Workflow
Discover -> Define -> Develop -> Deliver, with contract-first verification and parity checks at delivery time.
