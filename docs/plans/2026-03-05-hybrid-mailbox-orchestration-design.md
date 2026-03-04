# Hybrid Mailbox Orchestration Design (Boundary-First + Reviewer-Led Abort)

## Scope

Design a non-blocking orchestration pipeline where the implementer maintains throughput across multi-task runs while safely pivoting to rework when feedback invalidates prior work. The design extends current isolated worktree execution with a filesystem-backed mailbox and deterministic gate/abort behavior.

This design is additive to existing shared execution isolation and integration-retry behavior.

## Problem Statement

Traditional orchestration trade-off:

1. Sequential review blocks implementer throughput.
2. Blind parallelization risks stacking tasks on flawed upstream logic.

Goal:

1. Keep implementer velocity high via boundary-first progression.
2. Prevent poisoned downstream work via targeted immediate aborts.
3. Preserve deterministic recovery and cleanup in long multi-task runs.

## Product Decisions (Confirmed)

1. **Primary scheduling model:** non-preemptive, boundary-first gate checks.
2. **Immediate abort policy:** hybrid two-tier strategy:
- reviewer semantic invalidation (`invalidate_descendants`)
- orchestrator structural invalidation (file-overlap conflict)
3. **Resume model:** explicit states:
- `RESUME_IN_PLACE`
- `RESTART_FROM_SCRATCH`
4. **IPC medium:** filesystem mailbox with artifact pointers.
5. **Cleanup model:** accepted/integrated cleanup + aborted tombstoning + run-end sweep.
6. **Resource guardrail:** cap active isolated worktrees with queue backpressure.

## Architecture

### Core Layers

1. **Implementer Loop (M1):** executes task steps in isolated worktrees and checks gate at each boundary.
2. **Reviewer Loop (M2):** evaluates submitted artifacts in provided worktree and emits structured verdicts.
3. **Orchestrator Control Plane:** dependency-aware abort controller, conflict monitor, mailbox watcher, and lifecycle manager.

### Execution Model

1. Each task attempt executes in a dedicated git worktree/branch.
2. M1 submits immutable artifact pointers (worktree + commit + context).
3. M1 does not block on review by default; it proceeds to next step/task until gate or control event requires pivot.

### Abort Authority Split

1. **Semantic tier (reviewer-led):** `invalidate_descendants` verdict invalidates dependency descendants.
2. **Structural tier (orchestrator-led):** changed-file overlap invalidates active stale-base task attempts.

### Real-Time Watcher Tier

A background `MailboxWatcher` goroutine monitors reviewer/orchestrator inbox events in near real time.

1. On high-priority semantic or structural invalidation, watcher triggers immediate task abort signaling.
2. Control event is recorded and handed to executor state transition engine for deterministic requeue.

### Finality Layer (Artifact Lifecycle & Cleanup)

1. **Promotion:** `accepted` artifact is integrated into run integration branch; sandbox marked for cleanup.
2. **Tombstoning:** aborted attempt is removed immediately (`git worktree remove --force` + branch delete).
3. **Run-end sweep:** remove leftover temporary sandboxes for the run.

## Data Models / Mailbox Schema

### Runtime Layout

```text
~/.claude-octopus/runs/<run_id>/
  mailbox/
    inbox-m1/
    inbox-m2/
    inbox-orchestrator/
    processed/
    tmp/
  manifests/
    task-<id>-changed-files.json
    task-<id>-state-snapshot.json
  state/
    active-tasks.json
```

### Message Envelope

```json
{
  "message_id": "uuid",
  "run_id": "run-20260305-123",
  "type": "artifact_ready|review_verdict|control_abort|task_requeue|task_accepted|cleanup_tombstone",
  "from": "m1|m2|orchestrator",
  "to": "m1|m2|orchestrator",
  "priority": "normal|high",
  "created_at": "2026-03-05T10:00:00Z",
  "task_id": "task-2",
  "attempt_id": "task-2-attempt-1",
  "step_id": "step-1",
  "payload": {}
}
```

### Artifact Pointer Payload

1. `worktree_path`
2. `branch`
3. `commit_sha`
4. `base_sha`
5. `changed_files_manifest`
6. `task_context`
7. `depends_on`

### Review Verdict Payload

1. `verdict` = `accepted|soft_rework|hard_rework|invalidate_descendants`
2. `reason_code`
3. `evidence_ref`
4. `target_commit_sha`

### Abort/Requeue Payload

1. `abort_reason` = `semantic_invalidate|structural_overlap|operator_cancel`
2. `overlap_files`
3. `requeue_base_sha`
4. `resume_mode` = `RESUME_IN_PLACE|RESTART_FROM_SCRATCH`
5. `resume_from` = `task_start|step_boundary`
6. `stale_artifact_id`

### Reliability Contract

1. **Atomic handover rule:** writers must never write directly into `inbox-*`.
2. Writers must create JSON in `mailbox/tmp/` on the same filesystem, `fsync` as needed, then atomically `os.Rename` into destination inbox.
3. Consumers process one file per message and move it to `processed/` after successful handling.
4. Consumption must be idempotent by `message_id`.
5. Processing order is `(created_at, message_id)`.

## Control Flow + Gate Algorithm

### Boundary-First Implementer Loop

1. Execute current step of active task.
2. Persist step snapshot in `task-<id>-state-snapshot.json`.
3. Enter opportunistic gate before next step.
4. Apply control decisions in deterministic order.
5. Continue or pivot based on gate outcome.

### Deterministic Gate Order

1. Apply `control_abort` for active attempt.
2. Apply semantic `invalidate_descendants` verdicts.
3. Apply structural overlap decisions.
4. Apply pending `task_requeue` instructions.
5. Resume active work or pull next task.

### Semantic Abort Path

1. Reviewer emits `invalidate_descendants` for parent task `K`.
2. Orchestrator resolves descendants using `dependency_graph`.
3. Active descendants are aborted immediately.
4. Requeue emits new `requeue_base_sha` and `stale_artifact_id`.

### Structural Abort Path

1. On accepted/integrated artifact, orchestrator computes changed files.
2. Compare against active attempt touched files.
3. Any overlap triggers immediate structural abort.
4. Aborted task is requeued on fresh base from current integration head.

### Resume Semantics

1. `RESUME_IN_PLACE`:
- boundary pause only
- no stale base
- resume from saved boundary snapshot

2. `RESTART_FROM_SCRATCH`:
- any immediate abort
- any `stale_artifact_id`
- allocate fresh worktree and restart from `task_start`

### Watcher + Executor Coordination

1. Watcher detects and emits high-priority control events.
2. Executor remains sole owner of state transitions and attempt lifecycle writes.
3. Duplicate abort events are debounced per `attempt_id`.

## Lifecycle Management

### Promotion and Integration

1. `accepted` artifacts are integrated into `bench/<run_id>/integration`.
2. Integration metadata is persisted for traceability.

### Tombstoning and Cleanup

1. Aborted attempt is tombstoned immediately.
2. Accepted integrated sandbox is removed after integration confirmation.
3. Run-end sweep removes unreferenced run sandboxes.
4. Cleanup failures are recorded as deferred cleanup tasks.

### Resource Guardrail

1. Enforce `ActiveWorktreeCap` (configurable, recommended default `10-20`).
2. If cap reached, orchestrator pauses "pull next task" and keeps processing review/integration/cleanup events.
3. New task dispatch resumes when an `accepted` or `cleanup_tombstone` event frees a slot.

## Error Handling

### Error Classes

1. `mailbox_io_error`
2. `abort_delivery_error`
3. `state_transition_error`
4. `integration_error`
5. `cleanup_error`
6. `resource_cap_error`

### Failure Policy

1. Mailbox I/O uses bounded retry.
2. If mailbox read fails after retry, fail open for boundary progression and emit structured error.
3. If abort signal delivery fails, escalate to process-group kill path.
4. Invalid resume transitions downgrade to `RESTART_FROM_SCRATCH`.
5. Cap-exceeded condition is backpressure, not fatal run failure.

## Invariants

1. Exactly one active attempt per `task_id`.
2. Every `task_requeue` has explicit `resume_mode` and base metadata.
3. `RESTART_FROM_SCRATCH` must use a new worktree path.
4. Processed mailbox messages are immutable.
5. Abort reason is always enumerated.
6. Active attempts count must never exceed `ActiveWorktreeCap`.

## Observability

Emit structured events with keys: `run_id`, `task_id`, `attempt_id`, `message_id`.

Core events:

1. `watcher_abort_issued`
2. `abort_applied`
3. `task_requeued`
4. `resume_mode_selected`
5. `cleanup_tombstoned`
6. `cleanup_deferred`
7. `worktree_cap_reached`
8. `worktree_slot_freed`

Recommended counters:

1. abort count by reason
2. requeue count
3. resume mode distribution
4. overlap conflicts avoided
5. estimated token waste prevented
6. cap backpressure wait time

## Testing Strategy

### Unit Tests

1. Envelope/schema validation.
2. Atomic handover contract (`tmp -> rename`) and partial-file safety.
3. Idempotent consumption by `message_id`.
4. Gate ordering and resume-mode resolver.
5. `stale_artifact_id` invalidation behavior.
6. active-worktree-cap scheduler gating.

### Integration Tests

1. Semantic invalidation aborts active descendants.
2. Structural overlap triggers immediate abort + requeue on fresh base.
3. Boundary pause with non-breaking rework uses `RESUME_IN_PLACE`.
4. Immediate abort paths force `RESTART_FROM_SCRATCH`.
5. Mailbox transient failures respect retry/fail-open policy.
6. Cap reached pauses new pulls until slot is freed.

### End-to-End Scenario

Use API migration sequence (`/users`, `/orders`, `/products`):

1. Task 1 submitted while Task 2 progresses.
2. Reviewer invalidates Task 1 or overlap is detected.
3. Task 2 is aborted/requeued correctly.
4. Task 3 starts only after gate and cap conditions permit.
5. Run exits with no orphaned run worktrees.

## Out of Scope

1. Cross-machine distributed mailbox transport.
2. Multi-run global fair scheduling.
3. Automatic fallback to rank #2 integration on rank #1 failure.

## Next Step

Create a concrete implementation plan covering:

1. config/schema additions (`dependency_graph`, `active_worktree_cap`, mailbox paths)
2. mailbox persistence API
3. watcher + gate integration points
4. abort/requeue state machine updates
5. cleanup/cap scheduler changes
6. verification matrix and test rollout
