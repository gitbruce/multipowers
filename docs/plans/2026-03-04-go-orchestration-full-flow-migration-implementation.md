# Go Orchestration Full-Flow Migration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement Go-native advanced orchestration flows for `discover/define/develop/deliver/debate/embrace` with process-level parity to main shell behavior, using config-driven planning and provider dispatch.

**Architecture:** Build a unified orchestration runtime with planner, executor, and synthesizer modules. Planner resolves merged config from `workflows.yaml` flow-specific overrides and `orchestration.yaml` global defaults. Executor performs parallel step dispatch using policy resolver/dispatcher with one-hop fallback. Synthesizer supports progressive and final synthesis for each flow.

**Tech Stack:** Go (`internal/orchestration`, `internal/workflows`, `internal/cli`, `internal/policy`), YAML config (`config/workflows.yaml`, `config/orchestration.yaml`, `config/providers.yaml`, `config/agents.yaml`), Go tests (`go test ./...`).

---

## Task Status Tracker (MANDATORY UPDATE)

| ID | Task | Status |
|---|---|---|
| T01 | Finalize config contracts and schema alignment | COMPLETED |
| T01-S01 | Add orchestration config schema/types | COMPLETED |
| T01-S02 | Add workflow override config schema | COMPLETED |
| T01-S03 | Add provider model-pattern inference schema support | COMPLETED |
| T02 | Implement config loading and merge resolution | NOT_STARTED |
| T02-S01 | Load `orchestration.yaml` globals | COMPLETED |
| T02-S02 | Load workflow flow-specific override sections | COMPLETED |
| T02-S03 | Implement precedence merge (`task > workflow > global`) | COMPLETED |
| T03 | Implement planner for all flows | NOT_STARTED |
| T03-S01 | Build phase graph and perspective plan objects | NOT_STARTED |
| T03-S02 | Build flow/task-specific plan builder API | NOT_STARTED |
| T03-S03 | Add planner tests for all 6 flows | NOT_STARTED |
| T04 | Implement parallel executor and result collection | NOT_STARTED |
| T04-S01 | Build worker-pool step executor with context cancel | NOT_STARTED |
| T04-S02 | Add fallback-aware step result model | NOT_STARTED |
| T04-S03 | Add synchronization/progress event stream | NOT_STARTED |
| T05 | Implement progressive and final synthesis | NOT_STARTED |
| T05-S01 | Progressive trigger engine (`min_completed`, `min_bytes`) | NOT_STARTED |
| T05-S02 | Final synthesis aggregator and report model | NOT_STARTED |
| T05-S03 | Synthesis failure/degraded behavior contract | NOT_STARTED |
| T06 | Migrate workflow entrypoints to orchestration engine | NOT_STARTED |
| T06-S01 | Replace `internal/workflows/discover.go` facade | NOT_STARTED |
| T06-S02 | Replace `define/develop/deliver/debate/embrace` facades | NOT_STARTED |
| T06-S03 | Ensure returned metadata parity and structured report output | NOT_STARTED |
| T07 | Integrate orchestration commands into CLI | NOT_STARTED |
| T07-S01 | Add `mp orchestrate select-agent` final behavior | NOT_STARTED |
| T07-S02 | Add `mp loop` final ralph-wiggum behavior | NOT_STARTED |
| T07-S03 | Add CLI tests for new commands and invalid config cases | NOT_STARTED |
| T08 | Remove legacy shell-era assumptions in Go paths | NOT_STARTED |
| T08-S01 | Remove remaining phase hardcode assumptions from non-orchestration paths | NOT_STARTED |
| T08-S02 | Ensure no compatibility fallback path remains | NOT_STARTED |
| T09 | End-to-end verification for flow equivalence | NOT_STARTED |
| T09-S01 | Build E2E tests per flow (happy path) | NOT_STARTED |
| T09-S02 | Build E2E tests per flow (failure/degraded path) | NOT_STARTED |
| T09-S03 | Verify progressive synthesis behavior with deterministic fixtures | NOT_STARTED |
| T10 | Docs and migration evidence | NOT_STARTED |
| T10-S01 | Update architecture docs for orchestration runtime | NOT_STARTED |
| T10-S02 | Capture verification evidence and commands | NOT_STARTED |

## Mandatory Status Update Rule

After each subtask completion, do both actions:

1. Update this table status (`NOT_STARTED -> IN_PROGRESS -> COMPLETED`).
2. Persist status in runtime state:

```bash
go run ./cmd/mp state set --dir . --key plan.go_orchestration_full_flow.<TASK_ID> --value <IN_PROGRESS|COMPLETED> --json
```

Example:

```bash
go run ./cmd/mp state set --dir . --key plan.go_orchestration_full_flow.T03-S02 --value COMPLETED --json
```

No subtask may be marked done without this status update.

---

### Task T01: Finalize config contracts and schema alignment

**Why:** All runtime behavior depends on strict and predictable config contracts. Without finalized schema, implementation becomes unstable and breaks across flows.

**What:** Align `orchestration/workflows/providers/agents` schema roles and enforce correct field ownership.

**How:** Write schema tests first, then implement minimal type/load validation to make tests pass.

**Key Design:**
- `workflows.yaml` contains only flow-specific overrides.
- `orchestration.yaml` contains global defaults.
- `providers.yaml` resolves model-to-provider execution behavior.

#### Subtask T01-S01: Add orchestration config schema/types

**Files:**
- Modify: `internal/orchestration/types.go`
- Modify: `internal/orchestration/load.go`
- Test: `internal/orchestration/load_test.go`

**Step 1 (RED):** Add failing tests for required orchestration sections and default values.

**Step 2 (RED verify):**
Run: `go test ./internal/orchestration -run TestLoadConfigFromProjectDir -count=1`
Expected: fail for missing parsing/default merge.

**Step 3 (GREEN):** Implement schema/type load support.

**Step 4 (GREEN verify):** rerun command and expect pass.

**Step 5 (Status):** Update tracker + `state set` for `T01-S01`.

#### Subtask T01-S02: Add workflow override config schema

**Files:**
- Modify: `internal/policy/types.go`
- Modify: `internal/policy/schema_test.go`

**Why:** Workflow behavior must be override-based, not global-default duplicated.

**What:** Ensure workflow schema supports override sections (phase/perspective/parallel/synthesis optional).

**How:** Add tests for optional override blocks and no hard requirement for executor profile.

**Key Design:** Missing override means fallback to orchestration global.

#### Subtask T01-S03: Add provider model-pattern inference schema support

**Files:**
- Modify: `internal/policy/types.go`
- Modify: `internal/policy/runtime_policy.go`
- Modify: `internal/policy/compile.go`
- Test: `internal/policy/resolve_test.go`

**Why:** Workflow should declare models, provider is inferred.

**What:** Add and persist `model_patterns` for providers.

**How:** Test-first for inference, then compile/runtime fields.

**Key Design:** Explicit profile wins; otherwise deterministic inference by patterns.

---

### Task T02: Implement config loading and merge resolution

**Why:** Engine needs one merged runtime plan source from global + flow overrides.

**What:** Implement resolver that merges configs by precedence.

**How:** Add pure-function merge unit tests and implement minimal merge logic.

**Key Design:** `task > workflow.default > orchestration.global`.

#### Subtask T02-S01: Load `orchestration.yaml` globals

**Files:**
- Modify: `internal/orchestration/load.go`
- Test: `internal/orchestration/load_test.go`

**Why:** Global defaults must be independently validated.

**What:** Parse phase defaults, loop settings, skill triggers with strict defaults.

**How:** Add table tests for malformed/empty configs.

**Key Design:** Hard-fail on invalid regex patterns or required missing fields.

#### Subtask T02-S02: Load workflow flow-specific override sections

**Files:**
- Modify: `internal/policy/load.go`
- Test: `internal/policy/load_test.go`

**Why:** Overrides must be discoverable without duplicating global config.

**What:** Parse optional orchestration override nodes in workflow/task.

**How:** Add tests with minimal and full override examples.

**Key Design:** No override node should change default behavior.

#### Subtask T02-S03: Implement precedence merge

**Files:**
- Create: `internal/orchestration/merge.go`
- Create: `internal/orchestration/merge_test.go`

**Why:** Centralized precedence avoids divergence across flows.

**What:** Merge global/workflow/task configs into one immutable resolved plan config.

**How:** Write merge tests first, then implement.

**Key Design:** Merge result must be deterministic and pure (same input, same output).

---

### Task T03: Implement planner for all flows

**Why:** Planner turns configs into executable plan structures; without it flow logic remains hardcoded.

**What:** Build generic planner that supports all six flows.

**How:** Build plan structs + per-flow templates + tests.

**Key Design:** One planner API for all workflows.

#### Subtask T03-S01: Build phase graph and perspective plan objects

**Files:**
- Create: `internal/orchestration/planner.go`
- Create: `internal/orchestration/plan_types.go`
- Test: `internal/orchestration/planner_test.go`

**Why:** Need typed plan objects for executor and synthesizer.

**What:** Define `ExecutionPlan`, `PhasePlan`, `StepPlan`, `SynthesisPlan`.

**How:** Add failing tests for discover/develop plan generation shape.

**Key Design:** Steps are immutable after plan creation.

#### Subtask T03-S02: Build flow/task-specific plan builder API

**Files:**
- Modify: `internal/orchestration/planner.go`
- Test: `internal/orchestration/planner_test.go`

**Why:** Task-level overrides must alter flow behavior predictably.

**What:** Add `BuildPlan(workflow, task, prompt, dir)` API.

**How:** Add tests for task-specific perspective/parallel override.

**Key Design:** Planner must attach resolved source refs for traceability.

#### Subtask T03-S03: Add planner tests for all 6 flows

**Files:**
- Modify: `internal/orchestration/planner_test.go`

**Why:** Avoid hidden regressions in non-discover flows.

**What:** Add test cases for `discover/define/develop/deliver/debate/embrace`.

**How:** Use fixture table with expected phases and step count.

**Key Design:** Keep tests data-driven.

---

### Task T04: Implement parallel executor and result collection

**Why:** Flow-level parity depends on concurrent multi-perspective execution.

**What:** Execute planned steps in parallel with reliable synchronization.

**How:** Implement worker pool with context cancellation and structured result collection.

**Key Design:** No goroutine leaks; deterministic join.

#### Subtask T04-S01: Build worker-pool step executor with context cancel

**Files:**
- Create: `internal/orchestration/executor.go`
- Test: `internal/orchestration/executor_test.go`

**Why:** Need bounded concurrency and cancellation safety.

**What:** Add configurable `max_workers`, timeout, cancellation behavior.

**How:** Write tests that assert step count completion and cancel propagation.

**Key Design:** Use `context.Context` and `sync.WaitGroup`.

#### Subtask T04-S02: Add fallback-aware step result model

**Files:**
- Create: `internal/orchestration/result_types.go`
- Modify: `internal/orchestration/executor.go`
- Test: `internal/orchestration/executor_test.go`

**Why:** Must expose degraded/fallback metadata for auditability.

**What:** Persist fallback fields for each step.

**How:** Mock dispatch returning degraded and validate serialization fields.

**Key Design:** Preserve raw dispatch metadata.

#### Subtask T04-S03: Add synchronization/progress event stream

**Files:**
- Create: `internal/orchestration/events.go`
- Modify: `internal/orchestration/executor.go`
- Test: `internal/orchestration/executor_test.go`

**Why:** Progressive synthesis requires completion events.

**What:** Emit step lifecycle events over channel.

**How:** Add tests asserting event order and completion counts.

**Key Design:** Non-blocking event publish with bounded buffer.

---

### Task T05: Implement progressive and final synthesis

**Why:** Shell parity requires both early synthesis and final report generation.

**What:** Add synthesis engine that works incrementally and finalizes at completion.

**How:** Implement trigger engine + final aggregator with tests.

**Key Design:** Synthesis should be robust under partial failures.

#### Subtask T05-S01: Progressive trigger engine

**Files:**
- Create: `internal/orchestration/synthesis_progressive.go`
- Test: `internal/orchestration/synthesis_progressive_test.go`

**Why:** Reduce user wait by synthesizing early valid results.

**What:** Trigger only when `min_completed` and `min_bytes` are met.

**How:** Add deterministic fixture tests for threshold conditions.

**Key Design:** Deduplicate repeated trigger windows.

#### Subtask T05-S02: Final synthesis aggregator and report model

**Files:**
- Create: `internal/orchestration/synthesis_final.go`
- Create: `internal/orchestration/report.go`
- Test: `internal/orchestration/synthesis_final_test.go`

**Why:** Need consistent final report structure across flows.

**What:** Build final synthesis prompt/context from valid results.

**How:** Test report structure and section completeness.

**Key Design:** Exclude invalid/empty step outputs.

#### Subtask T05-S03: Synthesis failure/degraded behavior contract

**Files:**
- Modify: `internal/orchestration/report.go`
- Test: `internal/orchestration/synthesis_final_test.go`

**Why:** Failures should return structured degraded result, not panic.

**What:** Add `status`, `error_class`, partial context fields.

**How:** Simulate synthesis dispatch failure and assert degraded response shape.

**Key Design:** Preserve useful intermediate outputs.

---

### Task T06: Migrate workflow entrypoints to orchestration engine

**Why:** Existing facades do not execute advanced business flow logic.

**What:** Replace workflow wrappers with orchestration execution path.

**How:** Rewrite one flow at a time with tests, then remove old facade-only behavior.

**Key Design:** Uniform return envelope across all flows.

#### Subtask T06-S01: Replace discover facade

**Files:**
- Modify: `internal/workflows/discover.go`
- Test: `internal/workflows/discover_test.go` (create if missing)

**Why:** Discover has highest orchestration complexity and parity requirement.

**What:** Use planner->executor->synthesizer pipeline.

**How:** Add tests for perspectives, parallel execution metadata, synthesis output presence.

**Key Design:** Discover should drive progressive synthesis by config.

#### Subtask T06-S02: Replace define/develop/deliver/debate/embrace facades

**Files:**
- Modify: `internal/workflows/define.go`
- Modify: `internal/workflows/develop.go`
- Modify: `internal/workflows/deliver.go`
- Modify: `internal/workflows/debate.go`
- Modify: `internal/workflows/embrace.go`
- Test: corresponding `*_test.go` files

**Why:** Full-flow migration required, not discover-only.

**What:** Route all flows through same orchestration runtime.

**How:** Add table-driven tests to confirm flow-specific plan and report outputs.

**Key Design:** Keep flow wrappers thin and deterministic.

#### Subtask T06-S03: Ensure metadata parity and structured report output

**Files:**
- Modify: `internal/workflows/*.go`
- Modify: `internal/hooks/handler.go` (if metadata shape changes)
- Test: `internal/hooks/handler_test.go`

**Why:** Existing consumers expect normalized metadata fields.

**What:** Preserve `model_routing`, fallback, status, and route fields.

**How:** Add compatibility assertions in tests.

**Key Design:** No hidden implicit fields; all explicit in contracts.

---

### Task T07: Integrate orchestration commands into CLI

**Why:** Operators need direct orchestration controls for selection and loops.

**What:** Finalize CLI contract for `orchestrate select-agent` and `loop`.

**How:** Add command tests before refining handlers.

**Key Design:** CLI errors must be structured and actionable.

#### Subtask T07-S01: Add final select-agent behavior

**Files:**
- Modify: `internal/cli/root.go`
- Test: `internal/cli/root_test.go`

**Why:** Validate phase->agent selection against merged config.

**What:** Return selected agent, reason, and candidate list.

**How:** Add JSON mode tests and invalid phase tests.

**Key Design:** Canonicalize phase aliases (`discover->probe`, etc.).

#### Subtask T07-S02: Add final loop behavior

**Files:**
- Modify: `internal/cli/root.go`
- Test: `internal/cli/root_test.go`

**Why:** Implement ralph-style iterative loop in Go path.

**What:** Execute loop with completion promise and max iteration controls.

**How:** Add deterministic mock tests for completion and max-iteration stop.

**Key Design:** No shell fallback path.

#### Subtask T07-S03: Add CLI invalid config tests

**Files:**
- Modify: `internal/cli/root_test.go`

**Why:** One-way cutover requires strong failure visibility.

**What:** Test missing/invalid orchestration config behavior.

**How:** Assert blocked/error responses with clear remediation messages.

**Key Design:** Fail fast before execution starts.

---

### Task T08: Remove legacy assumptions and compatibility paths

**Why:** Migration is complete cutover without backward compatibility runtime toggles.

**What:** Remove or neutralize code that implies legacy fallback behavior.

**How:** Search-driven cleanup + guard tests.

**Key Design:** Runtime path must be single-source-of-truth.

#### Subtask T08-S01: Remove remaining phase hardcode assumptions

**Files:**
- Modify: `internal/workflows/persona.go` (if stale parsing behavior remains)
- Modify: `internal/hooks/*` as needed
- Test: `internal/validation/*`

**Why:** Prevent hidden drift from old shell semantics.

**What:** Remove dead legacy branch logic in Go runtime.

**How:** Add targeted tests to ensure no fallback to removed assumptions.

**Key Design:** Keep behavior explicit via config only.

#### Subtask T08-S02: Ensure no compatibility fallback path remains

**Files:**
- Modify: affected runtime files
- Test: `internal/validation/no_shell_runtime_test.go` and related tests

**Why:** Requirement explicitly forbids short-term compatibility switch.

**What:** Ensure no env-based fallback remains.

**How:** Add tests/search guards for forbidden patterns.

**Key Design:** Hard fail if required config missing.

---

### Task T09: End-to-end verification for flow equivalence

**Why:** Need confidence that orchestration behavior matches process-level expectations.

**What:** Add flow-level E2E tests for success and degraded paths.

**How:** Use test fixtures and deterministic dispatch mocks.

**Key Design:** One happy path + one failure path per flow minimum.

#### Subtask T09-S01: E2E happy path for all 6 flows

**Files:**
- Create/Modify: `internal/orchestration/e2e_flow_test.go`

**Why:** Verify all flows run through orchestration runtime.

**What:** Assert plan/execution/synthesis success and report completeness.

**How:** Table-driven test across six flows.

**Key Design:** Keep fixtures small but representative.

#### Subtask T09-S02: E2E degraded path for all 6 flows

**Files:**
- Modify: `internal/orchestration/e2e_flow_test.go`

**Why:** Must verify fallback and degraded contracts.

**What:** Simulate provider failure and assert one-hop fallback metadata.

**How:** Inject failing dispatch mock per flow.

**Key Design:** Ensure deterministic failure triggers.

#### Subtask T09-S03: Progressive synthesis E2E verification

**Files:**
- Modify: `internal/orchestration/e2e_flow_test.go`

**Why:** Progressive synthesis is a core parity behavior.

**What:** Assert trigger threshold behavior and intermediate synthesis outputs.

**How:** Use controlled result sizes and completion ordering.

**Key Design:** Validate both trigger and non-trigger cases.

---

### Task T10: Docs and migration evidence

**Why:** Future maintainers need exact runtime contracts and verification evidence.

**What:** Update architecture docs and capture command evidence.

**How:** Write docs after tests pass.

**Key Design:** Docs must map config fields to runtime behavior.

#### Subtask T10-S01: Update architecture docs

**Files:**
- Modify: `docs/architecture/*` relevant files
- Modify: `README.md` orchestration sections

**Why:** Remove shell-era ambiguity in docs.

**What:** Document planner/executor/synthesizer and config precedence.

**How:** Add explicit field-level reference sections.

**Key Design:** Keep docs implementation-accurate and concise.

#### Subtask T10-S02: Capture verification evidence

**Files:**
- Create: `docs/plans/evidence/model-routing/2026-03-04-go-orchestration-full-flow-verification.md`

**Why:** Completion claims require fresh evidence.

**What:** Record exact commands and outcomes.

**How:** Include:
- `go test ./internal/orchestration ./internal/workflows ./internal/cli ./internal/policy`
- `go test ./...`

**Key Design:** Evidence file must be reproducible and timestamped.

---

## Mandatory Verification Gate (Before completion claim)

Run all commands below successfully before any completion claim:

```bash
go test ./internal/orchestration ./internal/policy ./internal/workflows ./internal/cli -count=1
go test ./...
```

If any fail, do not mark corresponding tasks complete.

## Commit Strategy (Required)

Use frequent commits per major task (at minimum T01-T03, T04-T06, T07-T10 groups), with commit messages tied to task IDs.

Example:

```bash
git commit -m "feat(orchestration): implement planner and config merge (T02,T03)"
```
