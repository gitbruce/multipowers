# Architecture

This document describes the current runtime architecture on the `go` branch.

## High-Level Layers

1. **Context + Runtime Bootstrapping**
   - `internal/context` owns `/mp:init` materialization.
   - initialization creates `.multipowers/context/runtime.json` and the baseline governance files.
   - `runtime.json` is the single runtime/pre-run contract; `pre_run.enabled=false` by default.

2. **Spec Track Runtime**
   - `internal/tracks.TrackCoordinator` resolves the active track, allocates or reuses IDs, renders artifacts, and refreshes the canonical registry.
   - the only runtime registry path is `.multipowers/tracks/tracks.md`.
   - all spec artifacts live under `.multipowers/tracks/<track_id>/`.

3. **Spec Pipeline Enforcement**
   - `internal/app.RunSpecPipeline` checks context completeness, optional pre-run hooks, and track execution gates before command execution.
   - `internal/validation.EnsureTrackExecution` blocks the next spec command while an implementation group is still `in_progress` without completion evidence.

4. **Orchestration Runtime**
   - `internal/orchestration` builds immutable execution plans, executes step graphs with bounded workers, emits lifecycle events, and generates reports.
   - `internal/policy` remains the routing / execution-contract layer.

## Spec Track Runtime Model

### Canonical Files

- `.multipowers/context/runtime.json`
- `.multipowers/tracks/tracks.md`
- `.multipowers/tracks/<track_id>/intent.md`
- `.multipowers/tracks/<track_id>/design.md`
- `.multipowers/tracks/<track_id>/implementation-plan.md`
- `.multipowers/tracks/<track_id>/metadata.json`
- `.multipowers/tracks/<track_id>/index.md`

Legacy note:

- `.multipowers/tracks.md` is intentionally ignored by the runtime.

### Command Activity vs Group Lifecycle

Track metadata now separates two different concerns:

- **command activity**: `last_command`, `last_command_at`
- **implementation lifecycle**: `current_group`, `group_status`, `completed_groups`, `last_commit_sha`, `last_verified_at`

This prevents spec commands like `plan` or `develop` from pretending to be implementation groups.

### Explicit Group Lifecycle

Implementation progress is advanced through machine-callable commands:

```bash
mp track group-start --track-id <id> --group g1 --execution-mode workspace|worktree --json
mp track group-complete --track-id <id> --group g1 --commit-sha <sha> --json
```

Enforcement rules:

- `group-start` marks `group_status=in_progress`
- `group-complete` records commit + verification evidence and marks the group completed
- tracks marked `worktree_required=true` can only start groups from a linked git worktree checkout
- after every group transition, the canonical registry is refreshed

## Orchestration Hardening

### Retry Policy

Step plans carry explicit retry policy:

- `idempotent`
- `max_retries`
- `backoff_ms`
- `jitter_ratio`
- `retryable_codes`

Only idempotent steps are eligible for automatic retry.

### Deterministic Retry Controller

The executor wraps dispatch with a deterministic retry controller:

- retries only retryable failures
- uses bounded backoff
- respects `ctx.Done()` immediately
- records attempt count and terminal error metadata on `StepResult`

### Trace Propagation

Each orchestration run gets a stable `trace_id`:

- generated during plan build
- copied into step plans and step results
- attached to execution / phase / step events
- reused by structured lifecycle logs

### Structured Lifecycle Logs

The executor persists JSONL lifecycle records under:

- `.multipowers/<logs_subdir>/orchestration-<trace_id>.jsonl`

The log payload is machine-oriented and does not include raw prompt bodies by default.

### Regression Protection

`internal/orchestration/testdata/golden/` contains golden snapshots for:

- resolved plans
- generated reports
- degraded / fallback outputs

These protect the hardened runtime contracts from silent drift.

## Follow-Up Scope

The current hardening wave intentionally defers the UX / efficiency backlog tracked in:

- `docs/plans/2026-03-07-go-orchestration-ux-followups.md`

That follow-up document currently covers explainability CLI work, caching, and planner visualization.
