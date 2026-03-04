# Hybrid Mailbox Orchestration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement boundary-first mailbox orchestration with reviewer-led semantic aborts, orchestrator structural abort safety, deterministic resume/requeue rules, and active worktree cap backpressure.

**Architecture:** Add a filesystem mailbox module (atomic `tmp -> rename`, ordered/idempotent processing), connect it to orchestration via real-time watcher + deterministic gate logic, and extend runtime lifecycle for tombstoning/cleanup with resource-cap scheduling.

**Tech Stack:** Go, YAML, git worktree, filesystem JSON IPC, existing `internal/orchestration`, `internal/isolation`, and `go test`.

---

## Mandatory Status-Tracking Contract

1. Before touching code for a task, set that task status to `IN_PROGRESS` in the **Task Status Board**.
2. Each task has subtasks. Before each subtask starts, set subtask status to `IN_PROGRESS` in its **Subtask Status Board**.
3. After each subtask verification passes, set subtask status to `DONE` and add timestamp + commit hash (or `N/A` if no commit yet).
4. After the task-level verification and commit pass, set task status to `DONE`.
5. Allowed status values only: `TODO`, `IN_PROGRESS`, `BLOCKED`, `DONE`.
6. If blocked, set status to `BLOCKED` immediately and record blocker reason.

## Mandatory Design Metadata Contract

Every task and every subtask in this plan includes:

1. `Why`
2. `What`
3. `How`
4. `Key Design`

Do not execute a task/subtask unless those fields stay explicit and up to date.

## Required Skills and Guardrails

- `@superpowers:using-git-worktrees`
- `@superpowers:test-driven-development`
- `@superpowers:verification-before-completion`
- `@superpowers:executing-plans`

## Task Status Board

| Task ID | Task Name | Status | Last Update | Commit |
|---|---|---|---|---|
| T1 | Config Schema: Mailbox + Worktree Cap | DONE | 2026-03-05 03:23 CST | pending |
| T2 | Plan Metadata: Dependency Graph + Resume Modes | DONE | 2026-03-05 03:27 CST | pending |
| T3 | Mailbox Atomic Writer | TODO | - | - |
| T4 | Mailbox Reader + Idempotent Processor | TODO | - | - |
| T5 | Structural Conflict Monitor | TODO | - | - |
| T6 | Mailbox Watcher + Control Events | TODO | - | - |
| T7 | Deterministic Gate Engine | TODO | - | - |
| T8 | Worktree Slot Scheduler (Cap Backpressure) | TODO | - | - |
| T9 | Lifecycle Manager (Promotion/Tombstone/Sweep) | TODO | - | - |
| T10 | E2E Validation + Docs Sync | TODO | - | - |

---

### Task T1: Config Schema: Mailbox + Worktree Cap

**Status:** DONE

**Why:** Runtime orchestration needs explicit configuration for mailbox root/polling and active-worktree cap to enforce predictable safety and throughput.

**What:** Extend `execution_isolation` config schema with mailbox and cap fields, defaults, and validation.

**How:** Add new fields to config types, update loader defaults/validation, update config tests, and add sample config values.

**Key Design:** Keep backward compatibility; missing fields must default safely without enabling risky behavior.

**Files:**
- Modify: `internal/orchestration/types.go`
- Modify: `internal/orchestration/load.go`
- Modify: `internal/orchestration/config_test.go`
- Modify: `config/orchestration.yaml`

**Subtask Status Board**

| Subtask | Status | Why | What | How | Key Design | Verification |
|---|---|---|---|---|---|---|
| T1.1 Failing tests | DONE | Prevent schema drift | Add tests for new fields | Write load/validation tests first | TDD gate for config behavior | `go test ./internal/orchestration -run "TestLoadOrchestrationConfig_MailboxAndCap|TestLoadOrchestrationConfig_MailboxAndCapValidation" -v` |
| T1.2 Type fields | DONE | Config must parse fields | Add struct fields | Extend `ExecutionIsolationConfig` | Keep YAML names stable | compile + tests |
| T1.3 Defaults | DONE | Avoid nil/zero ambiguity | Apply defaults | Update `applyExecutionIsolationDefaults` | Safe defaults for V1 | tests pass |
| T1.4 Validation | DONE | Reject invalid config early | Validate cap/polling/root | Update `validateExecutionIsolationConfig` | Actionable error fields | validation test |
| T1.5 Config sample | DONE | Keep docs/runtime aligned | Add YAML keys | Update `config/orchestration.yaml` | No behavior flip unless enabled | loader tests |
| T1.6 Commit + status update | DONE | Preserve traceability | Commit and board update | `git add` + `git commit` | One task, one commit | clean staged diff |

---

### Task T2: Plan Metadata: Dependency Graph + Resume Modes

**Status:** DONE

**Why:** Abort/requeue logic needs deterministic dependency traversal and explicit resume semantics.

**What:** Add `dependency_graph` and resume metadata to execution plan structures.

**How:** Extend plan types, derive dependency graph in planner, add tests for deterministic graph output.

**Key Design:** V1 graph is static at plan build; keep model deterministic and loggable.

**Files:**
- Modify: `internal/orchestration/plan_types.go`
- Modify: `internal/orchestration/planner.go`
- Modify: `internal/orchestration/planner_test.go`

**Subtask Status Board**

| Subtask | Status | Why | What | How | Key Design | Verification |
|---|---|---|---|---|---|---|
| T2.1 Failing tests | DONE | Lock behavior before code | Add dependency/resume tests | New planner tests | Deterministic output checks | `go test ./internal/orchestration -run "TestBuildPlan_DependencyGraph|TestBuildPlan_TaskSnapshotDefaults" -v` |
| T2.2 Plan types | DONE | Represent runtime decisions | Add graph + resume structs | Update `plan_types.go` | Explicit resume enums | compile |
| T2.3 Graph builder | DONE | Resolve descendants fast | Build parent/descendant maps | Extend planner | Stable ordering in maps/slices | planner tests |
| T2.4 Planner wiring | DONE | Make graph available runtime | Attach graph to plan metadata | Update `BuildPlan` | Preserve old behavior if no deps | planner tests |
| T2.5 Commit + status update | DONE | Keep history clear | Commit and update board | standard git flow | one logical change | clean diff |

---

### Task T3: Mailbox Atomic Writer

**Status:** TODO

**Why:** Filesystem IPC fails if readers see partial JSON writes.

**What:** Implement mailbox writer using same-filesystem `tmp -> rename` atomic handover.

**How:** Add mailbox envelope types and atomic writer API with tests.

**Key Design:** Never write directly to `inbox-*`; enforce `tmp` path and atomic move.

**Files:**
- Create: `internal/mailbox/types.go`
- Create: `internal/mailbox/writer.go`
- Create: `internal/mailbox/writer_test.go`

**Subtask Status Board**

| Subtask | Status | Why | What | How | Key Design | Verification |
|---|---|---|---|---|---|---|
| T3.1 Failing tests | TODO | Prevent non-atomic IPC | Add writer tests first | Test tmp and final paths | Assert no partial inbox file | `go test ./internal/mailbox -run "TestWriteMessageAtomic" -v` |
| T3.2 Envelope types | TODO | Keep schema stable | Add envelope structs | Define fields from design doc | JSON tags match mailbox spec | compile |
| T3.3 Atomic write impl | TODO | Guarantee read safety | Write tmp then rename | `os.Rename` on same FS | Fail if tmp/inbox not same root policy | writer tests |
| T3.4 Error paths | TODO | Fail clearly and recoverably | Return typed errors | validate dirs and write failures | No silent drops | writer tests |
| T3.5 Commit + status update | TODO | Keep incremental trace | Commit + board update | standard git flow | one commit per task | clean diff |

---

### Task T4: Mailbox Reader + Idempotent Processor

**Status:** TODO

**Why:** Deterministic ordering and idempotency are required to avoid duplicate or out-of-order control actions.

**What:** Implement ordered inbox reading and safe move-to-processed processing.

**How:** Create reader/processor with sorting by `(created_at, message_id)` and idempotent handling by `message_id`.

**Key Design:** Processor must be re-runnable without reapplying side effects.

**Files:**
- Create: `internal/mailbox/reader.go`
- Create: `internal/mailbox/reader_test.go`
- Create: `internal/mailbox/processor.go`
- Create: `internal/mailbox/processor_test.go`

**Subtask Status Board**

| Subtask | Status | Why | What | How | Key Design | Verification |
|---|---|---|---|---|---|---|
| T4.1 Failing tests | TODO | Define ordering contract | Add order/idempotency tests | Seed multiple message files | Include tie on timestamp | `go test ./internal/mailbox -run "TestListInboxMessages|TestProcessOneMessage" -v` |
| T4.2 Ordered listing | TODO | Deterministic gate behavior | Implement message listing | Parse + sort by timestamp/id | Ignore filename ordering | reader tests |
| T4.3 Idempotent processing | TODO | Prevent duplicate effects | Implement processor | Move to `processed/` after success | no-op on duplicate message id | processor tests |
| T4.4 Failure handling | TODO | Avoid message loss | Retry-safe error behavior | preserve message on handler failure | crash-safe semantics | tests |
| T4.5 Commit + status update | TODO | Preserve progress visibility | Commit and board update | standard git flow | one task one commit | clean diff |

---

### Task T5: Structural Conflict Monitor

**Status:** TODO

**Why:** Active tasks can become stale if accepted upstream artifacts touch overlapping files.

**What:** Add file-overlap detector for structural abort decisions.

**How:** Normalize paths and compute intersection between active-touched files and accepted-changed files.

**Key Design:** Keep detector pure and deterministic so abort decisions are explainable.

**Files:**
- Create: `internal/orchestration/conflict_monitor.go`
- Create: `internal/orchestration/conflict_monitor_test.go`

**Subtask Status Board**

| Subtask | Status | Why | What | How | Key Design | Verification |
|---|---|---|---|---|---|---|
| T5.1 Failing tests | TODO | Lock conflict rules | Add overlap/no-overlap tests | Unit tests with normalized paths | deterministic output ordering | `go test ./internal/orchestration -run "TestConflictMonitor" -v` |
| T5.2 Monitor impl | TODO | Produce overlap reason data | Add monitor API | path normalization + set intersection | sorted overlap list | unit tests |
| T5.3 Integration utility | TODO | Reuse across watcher/gate | Add helper signatures | keep package-local minimal API | no side effects | compile |
| T5.4 Commit + status update | TODO | keep audit trail clean | commit and board update | standard git flow | one commit task | clean diff |

---

### Task T6: Mailbox Watcher + Control Events

**Status:** TODO

**Why:** Immediate aborts require real-time control signal extraction, not only boundary polling.

**What:** Add mailbox watcher goroutine and typed control events for semantic/structural aborts.

**How:** Poll reviewer inbox, decode high-priority messages, emit control events, ack to processed on success.

**Key Design:** Watcher emits events only; executor remains sole state-transition authority.

**Files:**
- Create: `internal/orchestration/control_events.go`
- Create: `internal/orchestration/mailbox_watcher.go`
- Create: `internal/orchestration/mailbox_watcher_test.go`
- Modify: `internal/orchestration/events.go`

**Subtask Status Board**

| Subtask | Status | Why | What | How | Key Design | Verification |
|---|---|---|---|---|---|---|
| T6.1 Failing tests | TODO | Protect real-time control behavior | Add watcher tests | semantic and structural cases | assert immediate control emission | `go test ./internal/orchestration -run "TestMailboxWatcher" -v` |
| T6.2 Event contracts | TODO | Keep control API explicit | Add control event structs/enums | create typed event payloads | stable reason codes | compile |
| T6.3 Watcher runtime | TODO | Connect mailbox to control stream | implement poll loop | reader + processor + channel emit | ack only after emit success | watcher tests |
| T6.4 Poll config wiring | TODO | tune performance safely | use config poll interval | pass interval via config | portable polling V1 | tests |
| T6.5 Commit + status update | TODO | keep implementation traceable | commit and board update | standard git flow | one logical commit | clean diff |

---

### Task T7: Deterministic Gate Engine

**Status:** TODO

**Why:** Gate must apply mailbox/control inputs in strict priority to prevent race-condition behavior drift.

**What:** Implement gate evaluator with fixed decision order and resume-mode resolver.

**How:** Build pure gate function that processes input sets in sequence: abort -> invalidate -> overlap -> requeue -> continue.

**Key Design:** `stale_artifact_id` must always force `RESTART_FROM_SCRATCH`.

**Files:**
- Create: `internal/orchestration/gate.go`
- Create: `internal/orchestration/gate_test.go`

**Subtask Status Board**

| Subtask | Status | Why | What | How | Key Design | Verification |
|---|---|---|---|---|---|---|
| T7.1 Failing tests | TODO | lock priority behavior | add gate-order tests | construct mixed event inputs | deterministic action assertions | `go test ./internal/orchestration -run "TestGateDecision" -v` |
| T7.2 Decision model | TODO | standardize control outputs | add decision/result structs | include action/reason/resume fields | explicit enums over strings where possible | compile |
| T7.3 Evaluator impl | TODO | enforce deterministic order | implement EvaluateGate | process inputs in fixed sequence | no hidden side effects | gate tests |
| T7.4 Resume resolver | TODO | prevent stale merges | implement stale handling | stale artifact => restart | explicit resume mode | gate tests |
| T7.5 Commit + status update | TODO | maintain delivery quality | commit and board update | standard git flow | one commit task | clean diff |

---

### Task T8: Worktree Slot Scheduler (Cap Backpressure)

**Status:** TODO

**Why:** Unlimited isolated worktrees can exhaust disk/inodes and destabilize long runs.

**What:** Add `ActiveWorktreeCap` scheduler that blocks new task pulls when cap is reached.

**How:** Implement slot manager with blocking acquire/release and wire it into executor scheduling flow.

**Key Design:** Cap pressure pauses queue intake, not review/integration/cleanup progress.

**Files:**
- Create: `internal/orchestration/worktree_slots.go`
- Create: `internal/orchestration/worktree_slots_test.go`
- Modify: `internal/orchestration/executor.go`
- Modify: `internal/orchestration/executor_test.go`

**Subtask Status Board**

| Subtask | Status | Why | What | How | Key Design | Verification |
|---|---|---|---|---|---|---|
| T8.1 Failing tests | TODO | lock cap behavior | add cap-blocking tests | unit + executor tests | assert no pull-next while full | `go test ./internal/orchestration -run "TestWorktreeSlots|TestExecutor_DoesNotPullNextTaskWhenCapReached" -v` |
| T8.2 Slot manager | TODO | enforce bounded concurrency | implement Acquire/Release | condvar + ctx cancel | no deadlocks | slot tests |
| T8.3 Executor wiring | TODO | enforce cap at runtime | gate pull-next path | acquire before new attempt, release on cleanup | avoid starving cleanup path | executor tests |
| T8.4 Event telemetry | TODO | observe pressure state | emit cap reached/freed events | reuse event emitter | drop-tolerant metrics | tests |
| T8.5 Commit + status update | TODO | keep reviewability high | commit and board update | standard git flow | one commit task | clean diff |

---

### Task T9: Lifecycle Manager (Promotion/Tombstone/Sweep)

**Status:** TODO

**Why:** Delayed cleanup leaks worktrees and stale branches, causing drift and resource waste.

**What:** Add lifecycle manager for accepted promotion cleanup, abort tombstoning, and run-end sweep.

**How:** Reuse isolation runtime cleanup APIs and add orchestration lifecycle hooks.

**Key Design:** Aborted attempts are tombstoned immediately; run-end sweep is best-effort with explicit error recording.

**Files:**
- Create: `internal/orchestration/lifecycle.go`
- Create: `internal/orchestration/lifecycle_test.go`
- Modify: `internal/isolation/runtime.go`
- Modify: `internal/isolation/runtime_test.go`
- Modify: `internal/orchestration/workflow_adapter.go`

**Subtask Status Board**

| Subtask | Status | Why | What | How | Key Design | Verification |
|---|---|---|---|---|---|---|
| T9.1 Failing tests | TODO | lock cleanup semantics | add lifecycle tests | accepted/aborted/sweep scenarios | deterministic cleanup actions | `go test ./internal/orchestration ./internal/isolation -run "TestLifecycle|TestIsolationRuntime" -v` |
| T9.2 Lifecycle API | TODO | centralize cleanup logic | add manager methods | OnAccepted/OnAborted/SweepRun | clear ownership boundary | compile |
| T9.3 Runtime extension | TODO | support sweep + tombstone | update runtime helper methods | remove worktree/branch safely | keep force cleanup explicit | runtime tests |
| T9.4 Adapter integration | TODO | trigger final sweep on run end | wire lifecycle calls | hook into adapter close/finish path | no blocking main success path | orchestration tests |
| T9.5 Commit + status update | TODO | preserve change isolation | commit and board update | standard git flow | one commit task | clean diff |

---

### Task T10: E2E Validation + Docs Sync

**Status:** TODO

**Why:** Complex async control systems require scenario-level verification and docs parity to avoid regressions.

**What:** Add E2E test for users/orders/products scenario and update architecture/trigger docs.

**How:** Implement failing E2E first, wire remaining glue, run full verification suite, then update docs.

**Key Design:** Validate both modes: boundary-first pivot and immediate abort path.

**Files:**
- Modify: `internal/orchestration/e2e_test.go`
- Modify: `docs/TRIGGERS.md`
- Modify: `docs/ARCHITECTURE.md`
- Modify: `docs/plans/2026-03-05-hybrid-mailbox-orchestration-design.md`

**Subtask Status Board**

| Subtask | Status | Why | What | How | Key Design | Verification |
|---|---|---|---|---|---|---|
| T10.1 Failing E2E | TODO | prove scenario semantics | add API migration E2E test | users->orders->products with invalidation/overlap branches | assert deterministic transitions | `go test ./internal/orchestration -run TestE2E_HybridMailboxBoundaryAndAbortFlow -v` |
| T10.2 Runtime glue | TODO | make E2E pass | connect watcher/gate/lifecycle/slots | minimal integration changes | no behavior change for non-mailbox flows | targeted tests |
| T10.3 Full verification | TODO | validate no regressions | run package and full test suites | execute listed commands | evidence before completion | `go test ./internal/mailbox ./internal/orchestration ./internal/isolation ./internal/benchmark -count=1` + `go test ./... -count=1` |
| T10.4 Docs sync | TODO | keep ops/dev alignment | update docs files | reflect final runtime behavior | no stale design claims | docs review |
| T10.5 Commit + status update | TODO | close implementation cleanly | commit + final board update | standard git flow | traceable final commit | clean diff |

---

## Standard Command Sequence Per Task

1. Update task/subtask status to `IN_PROGRESS`.
2. Write failing tests.
3. Run targeted tests and confirm failure.
4. Implement minimal code.
5. Run targeted tests and confirm pass.
6. Run neighboring package tests for regression check.
7. Commit task.
8. Update subtask and task statuses to `DONE` with timestamp/commit.

## Final Verification Checklist

1. Atomic mailbox writes are always `tmp -> rename` on same filesystem.
2. Gate decision order is deterministic.
3. `stale_artifact_id` forces `RESTART_FROM_SCRATCH`.
4. `ActiveWorktreeCap` is enforced with queue backpressure.
5. Semantic and structural immediate abort flows are tested.
6. Aborted tasks tombstone immediately.
7. Run-end sweep removes orphaned run worktrees.
8. Full test verification evidence is captured before completion claims.
