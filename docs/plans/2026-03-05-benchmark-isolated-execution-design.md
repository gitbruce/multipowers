# Benchmark Isolated Execution Design (Enhancement)

## Scope

Enhance the already completed benchmark + smart-routing capability so that when benchmark mode is enabled and a request is both code-related and whitelist-matched, multi-model execution is enforced through isolated git state per model.

This design is an additive enhancement over:
- `docs/plans/2026-03-04-benchmark-smart-routing-design.md`
- `docs/plans/2026-03-04-benchmark-smart-routing-implementation.md`

## Goals

1. Enforce `worktree + branch` isolation for each model in eligible benchmark runs.
2. Keep existing benchmark fault-tolerance guarantees (no critical-path blocking by benchmark side effects).
3. Introduce deterministic critic-based winner selection and integration behavior.
4. Standardize repair behavior when winner integration fails.

## Product Decisions (Confirmed)

1. **Enforcement type:** runtime enforcement (not documentation-only, not soft warning).
2. **Activation scope:** configurable command whitelist.
3. **Default winner policy:** critic rank #1 merges into integration branch directly.
4. **Failure policy for rank #1 merge/gate failure:** create repair task for the same model and retry once.

## Runtime Enforcement Policy

Isolation enforcement is active only when all are true:

1. `benchmark_mode.enabled = true`
2. `code_intent_final = true`
3. `command in benchmark_mode.execution_isolation.command_whitelist`
4. `benchmark_mode.execution_isolation.enabled = true`

If any condition is false, run existing non-isolated benchmark behavior.

## Architecture

### Main Path Additions

1. Resolve benchmark config and code-intent (existing path).
2. Evaluate `IsolationPolicy` (new).
3. If enforced:
   - create integration branch for run context.
   - create per-model branch/worktree sandbox.
   - dispatch each model execution in its own worktree.
4. Collect outputs and execution metadata.
5. Critic ranks candidates.
6. Pick top-1 candidate and integrate into run integration branch.
7. If merge/gate fails, issue one same-model repair retry.

### New Components

1. `internal/benchmark/isolation_policy.go`
- Deterministic policy resolver based on benchmark toggle + code intent + whitelist.

2. `internal/benchmark/isolation_runtime.go`
- Worktree/branch lifecycle manager.
- Safe create/cleanup contract.

3. `internal/benchmark/critic_selection.go`
- Deterministic ranking resolution with tie-break rules.

4. `internal/benchmark/integration_flow.go`
- Integration branch handling.
- Top-1 merge + single repair retry flow.

5. `internal/orchestration/*` integration points
- Inject isolation preflight and model execution context.

## Configuration Extension

Add under `benchmark_mode`:

```yaml
benchmark_mode:
  execution_isolation:
    enabled: false
    command_whitelist:
      - develop
      - review
      - embrace
    branch_prefix: "bench"
    worktree_root: ".worktrees/bench"
    repair_retry_max: 1
```

### Backward Compatibility

- Missing `execution_isolation` keeps current behavior.
- Defaults preserve non-isolation unless explicitly enabled.

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

## Failure Isolation

1. Benchmark support failures remain best-effort and must not deadlock orchestrator main path.
2. Isolation setup failure is explicitly logged and surfaced in benchmark records.
3. File/metrics write failures are bounded-retry then safe-drop.

## Acceptance Criteria (Enhancement)

1. Eligible benchmark runs always execute models in isolated worktrees and branches.
2. Isolation metadata is persisted for every candidate output.
3. Critic top-1 is selected deterministically and integrated first.
4. Merge/gate failure triggers exactly one same-model repair retry.
5. Retry outcomes are persisted with explicit statuses.
6. Existing benchmark no-block guarantee remains true.

## Out of Scope

1. Multi-retry cascading across rank #2/#3.
2. Cross-run automated conflict resolution heuristics.
3. Distributed/remote worktree orchestration.

