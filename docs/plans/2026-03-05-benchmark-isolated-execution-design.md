# Benchmark Isolated Execution Design (Enhancement)

## Scope

Enhance the already completed benchmark + smart-routing capability into a shared external-command isolation framework, so any workflow that invokes external commands and may edit files can enforce isolated git state per model/task.

This design is an additive enhancement over:
- `docs/plans/2026-03-04-benchmark-smart-routing-design.md`
- `docs/plans/2026-03-04-benchmark-smart-routing-implementation.md`

This enhancement is also generalized as shared runtime logic for any workflow path that invokes external commands and may edit files.

## Goals

1. Enforce `worktree + branch` isolation for each model/task in eligible external-command runs.
2. Keep existing fault-tolerance guarantees (no critical-path blocking by side effects).
3. Introduce deterministic critic-based winner selection and integration behavior.
4. Standardize repair behavior when winner integration fails.

## Product Decisions (Confirmed)

1. **Enforcement type:** runtime enforcement (not documentation-only, not soft warning).
2. **Activation scope:** configurable command whitelist.
3. **Default winner policy:** critic rank #1 merges into integration branch directly.
4. **Failure policy for rank #1 merge/gate failure:** create repair task for the same model and retry once.
5. **Applicability:** shared isolation logic is not limited to `benchmark_mode`.

## Runtime Enforcement Policy

Isolation enforcement is active only when all are true:

1. `execution_isolation.enabled = true`
2. `external_command_involved = true`
3. `may_edit_files = true`
4. optional profile gate matches (for benchmark profile: code-intent + profile whitelist)

If any condition is false, run existing non-isolated behavior.

## Shared External Command Isolation Policy

The same isolation/sync/progress design should be reused whenever external command execution is involved.

General policy input:
1. `external_command_involved`
2. `may_edit_files`
3. feature-specific toggle/scope (benchmark whitelist is one profile)

General enforcement rule:
- If external command execution may change workspace state, execute in isolated worktree/branch using the same runtime manager and status telemetry contract.

Benchmark mode is one policy profile built on top of this shared logic; non-benchmark external-command flows call the same policy/runtime APIs.

## Architecture

### Main Path Additions

1. Resolve execution context (`external_command_involved`, `may_edit_files`, optional profile signals like code intent).
2. Evaluate `IsolationPolicy` (new).
3. If enforced:
   - create integration branch for run context.
   - create per-model branch/worktree sandbox.
   - dispatch each model execution in its own worktree.
4. Collect outputs and execution metadata.
5. Critic ranks candidates.
6. Pick top-1 candidate and integrate into run integration branch.
7. If merge/gate fails, issue one same-model repair retry.
8. Emit model progress events continuously so long-running jobs are observable.

### New Components

1. `internal/isolation/policy.go`
- Shared policy resolver for external-command isolation.
- Supports profile-specific gates (benchmark and future profiles).

2. `internal/isolation/runtime.go`
- Shared worktree/branch lifecycle manager.
- Safe create/cleanup contract.

3. `internal/isolation/critic_selection.go`
- Deterministic ranking resolution with tie-break rules.

4. `internal/isolation/integration_flow.go`
- Integration branch handling.
- Top-1 merge + single repair retry flow.

5. `internal/orchestration/*` integration points
- Inject isolation preflight and model execution context.

6. `internal/orchestration/events.go` progress payload contract
- Reuse `EventTypeStepStart`, `EventTypeStepProgress`, `EventTypeStepEnd`.
- Standardize isolation progress payload in `Event.Data`:
  - `run_id`
  - `model`
  - `status` (`queued|sandbox_ready|running|completed|failed|timeout|repair_retry`)
  - `percent`
  - `branch`
  - `worktree_path`
  - `heartbeat_at`

7. shared isolation adapter for external command paths
- Any new external command flow can call the same isolation runtime/policy APIs.

### Sync Gate for Long-Running Model Fan-Out

At the end of model fan-out, orchestrator waits with a bounded gate:

1. Use `sync.WaitGroup` to wait candidate goroutines.
2. Wrap wait with timeout (`context.WithTimeout` from `execution_isolation` config).
3. If all candidates finish before timeout: proceed normally to critic ranking.
4. If timeout hits:
- proceed with completed candidates only.
- mark unfinished candidates as `timeout`.
- persist timeout metadata to benchmark records.

This gate blocks candidate collection/ranking only; async persistence remains best-effort and non-blocking.

## Configuration Extension

Add shared config at root level:

```yaml
execution_isolation:
  enabled: false
  command_whitelist:
    - develop
    - review
    - embrace
  branch_prefix: "bench"
  worktree_root: ".worktrees/bench"
  repair_retry_max: 1
  global_timeout_ms: 120000
  proceed_policy: "all_or_timeout"
  min_completed_models: 1
  heartbeat_interval_seconds: 30
  logs_subdir: "logs"

benchmark_mode:
  execution_profile:
    enabled: true
    require_code_intent: true
    command_whitelist:
      - develop
      - review
      - embrace
```

### Backward Compatibility

- Missing `execution_isolation` keeps current behavior.
- Defaults preserve non-isolation unless explicitly enabled.
- Non-benchmark external command flows can adopt this logic incrementally without behavior changes until explicitly enabled.

## Critic Ranking Contract

1. Rank by weighted score descending.
2. Tie-breakers:
- fewer execution failures
- shorter duration
- lexical model name
3. Select rank #1 directly for integration.

## Integration + Retry Contract

1. Target branch: `bench/<run_id>/integration`
2. Merge/cherry-pick rank #1 candidate.
3. On merge conflict or gate failure:
- generate repair prompt for same model branch
- execute repair once (`repair_retry_max=1`)
- retry merge/gate once
4. If retry still fails:
- record `failed_after_retry`
- emit structured error events
- do not auto-fallback to rank #2

## Data Model Enhancements

### New stream

`isolation_runs.YYYY-MM-DD.jsonl`:
- `run_id`
- `enforced`
- `reason`
- `command`
- `code_intent_final`
- `whitelist_match`
- `models`
- `worktree_root`
- `branch_prefix`

### Extended existing streams

1. `model_outputs.*.jsonl`
- `exec_branch`
- `exec_worktree`
- `exec_head_sha`
- `files_changed_count`

2. `judge_scores.*.jsonl`
- `candidate_branch`
- `candidate_worktree`
- `rank`

3. `route_overrides.*.jsonl`
- `selection_mode = critic_top1_direct`
- `integration_branch`
- `integration_status`
- `repair_retry_used`

4. `errors.*.jsonl`
- stage additions:
  - `isolation_setup`
  - `integration_merge`
  - `integration_gate`
  - `repair_retry`

5. `async_jobs.*.jsonl`
- add/update fields:
  - `run_id`
  - `model`
  - `stage`
  - `heartbeat_at`
  - `attempt`
  - `status`

### Optional Runtime Logs

For long-running and debugging scenarios, each isolated model sandbox may stream execution logs to:

`.worktrees/bench/<run_id>/<model>/logs/`

Log write failures must never fail the main orchestration result.

## Failure Isolation

1. Isolation support failures remain best-effort and must not deadlock orchestrator main path.
2. Isolation setup failure is explicitly logged and surfaced in isolation/benchmark records.
3. File/metrics write failures are bounded-retry then safe-drop.

## Acceptance Criteria (Enhancement)

1. Eligible external-command runs always execute in isolated worktrees and branches.
2. Isolation metadata is persisted for every candidate output.
3. Critic top-1 is selected deterministically and integrated first.
4. Merge/gate failure triggers exactly one same-model repair retry.
5. Retry outcomes are persisted with explicit statuses.
6. Existing no-block guarantee remains true.
7. Long-running model tasks surface continuous progress/heartbeat updates.
8. With timeout policy enabled, critic ranking proceeds with completed candidates instead of failing the whole run.
9. Isolation runtime/policy components are reusable by non-benchmark external command paths with shared semantics.

## Implementation Notes (2026-03-05)

1. Shared isolation entrypoint is exposed via `ResolveExternalCommandIsolation(...)` for benchmark and non-benchmark external-command flows.
2. Runtime isolation lifecycle (`CreateModelSandbox`, `CleanupModelSandbox`) is implemented in `internal/isolation/runtime.go`.
3. Progress telemetry uses `EventTypeStepProgress` with `ModelProgressData` payload and heartbeat timestamps.
4. Sync gate timeout degradation is implemented as shared `internal/isolation/sync_gate.go` and consumed by orchestration helpers.
5. Critic deterministic top-1 selection and same-model single repair retry integration flow are implemented in `internal/isolation/*`.
6. Benchmark JSONL schema is extended with isolation/progress/integration metadata fields and `isolation_runs` stream contract.

## Out of Scope

1. Multi-retry cascading across rank #2/#3.
2. Cross-run automated conflict resolution heuristics.
3. Distributed/remote worktree orchestration.
