# Go Orchestration Full-Flow Migration Optimization Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Harden the Go orchestration runtime against architectural drift, improve failure diagnosability and runtime stability, and add efficiency/UX upgrades without reintroducing shell-era behavior.

**Architecture:** Implement optimization in three waves: governance guardrails (must-have), reliability/observability (core runtime quality), and efficiency/UX (cost and developer experience). All changes stay config-driven and centered in `internal/orchestration` + `internal/policy` + `internal/validation`, with explicit CLI contracts for explainability and error semantics.

**Tech Stack:** Go 1.24, YAML config, JSON Schema validation, CLI (`cmd/mp`, `cmd/mp-devx`), tests (`go test ./...`), markdown evidence docs.

---

## Task Status Tracker (MANDATORY UPDATE)

| ID | Task | Status |
|---|---|---|
| O01 | Routing hardcode guardrails | NOT_STARTED |
| O01-S01 | Extend hardcode scanner rules and allowlist | COMPLETED |
| O01-S02 | Add forbidden phase/agent/provider routing pattern tests | COMPLETED |
| O01-S03 | Wire guard into full validation suite and docs | COMPLETED |
| O02 | Config governance + lint gate | NOT_STARTED |
| O02-S01 | Add JSON Schema for orchestration/workflows/providers | NOT_STARTED |
| O02-S02 | Implement `mp-devx lint-config` strict checks | COMPLETED |
| O02-S03 | Enforce pre-commit/CI config lint gate | COMPLETED |
| O03 | Standardized error semantics | NOT_STARTED |
| O03-S01 | Introduce typed orchestration error codes | COMPLETED |
| O03-S02 | Map error codes to CLI exit codes | COMPLETED |
| O03-S03 | Add compatibility tests across workflows/cli/hooks | COMPLETED |
| O04 | Executor retry reliability | NOT_STARTED |
| O04-S01 | Add idempotent retry policy fields to plan/step config | COMPLETED |
| O04-S02 | Implement exponential backoff + jitter retry loop | NOT_STARTED |
| O04-S03 | Add deterministic retry tests for 429/503/timeouts | NOT_STARTED |
| O05 | Traceability and observability | NOT_STARTED |
| O05-S01 | Generate and propagate `trace_id` across runtime | NOT_STARTED |
| O05-S02 | Emit structured logs for step lifecycle/fallback | NOT_STARTED |
| O05-S03 | Add trace correlation tests and log format tests | NOT_STARTED |
| O06 | Merge explainability CLI | NOT_STARTED |
| O06-S01 | Build explain resolver output model (value + source) | NOT_STARTED |
| O06-S02 | Add `mp orchestrate explain` command | NOT_STARTED |
| O06-S03 | Add explain tests and docs | NOT_STARTED |
| O07 | Result caching | NOT_STARTED |
| O07-S01 | Add cache key strategy + storage abstraction | NOT_STARTED |
| O07-S02 | Add orchestration execution read/write cache path | NOT_STARTED |
| O07-S03 | Add TTL/invalidations tests + cache metadata | NOT_STARTED |
| O08 | Planner visualization | NOT_STARTED |
| O08-S01 | Add DAG export from `ExecutionPlan` | NOT_STARTED |
| O08-S02 | Add Mermaid renderer + CLI output mode | NOT_STARTED |
| O08-S03 | Add visualization tests and documentation | NOT_STARTED |
| O09 | Layered regression tests (golden) | NOT_STARTED |
| O09-S01 | Add plan snapshot golden tests | NOT_STARTED |
| O09-S02 | Add synthesis report golden tests | NOT_STARTED |
| O09-S03 | Add degraded/fallback golden tests | NOT_STARTED |
| O10 | Verification evidence and rollout docs | NOT_STARTED |
| O10-S01 | Capture optimization verification evidence | NOT_STARTED |
| O10-S02 | Update architecture/ops docs and rollout notes | NOT_STARTED |

## Mandatory Status Update Rule (No Exceptions)

After **every subtask** completion, perform both actions:

1. Update this document table (`NOT_STARTED -> IN_PROGRESS -> COMPLETED`).
2. Persist state to runtime:

```bash
go run ./cmd/mp state set --dir . --key plan.go_orchestration_optimization.<TASK_ID> --value <IN_PROGRESS|COMPLETED> --json
```

Example:

```bash
go run ./cmd/mp state set --dir . --key plan.go_orchestration_optimization.O04-S02 --value COMPLETED --json
```

No subtask may be claimed complete without both updates.

---

### Task O01: Routing hardcode guardrails

#### Subtask O01-S01: Extend hardcode scanner rules and allowlist
**Why:** Current guard only focuses on model strings and cannot prevent phase/agent/provider drift logic from leaking outside orchestration boundaries.
**What:** Expand guard to scan for forbidden routing patterns (phase switch, agent name lists, provider-specific if/else) in disallowed packages.
**How:** Refactor current scanner in `internal/validation/model_hardcode_guard_test.go` into reusable rule groups and package-level allowlist.
**Key Design:** Allow routing logic only in `internal/orchestration/**` and `internal/policy/**`; everywhere else must consume resolved contracts.

**Files:**
- Modify: `internal/validation/model_hardcode_guard_test.go`
- Create: `internal/validation/routing_hardcode_guard_test.go` (if cleaner split is preferred)

**Steps:**
1. Write failing tests for forbidden patterns in `internal/workflows`, `internal/cli`, `internal/hooks`.
2. Run: `go test ./internal/validation -run HardcodeGuard -count=1` (expect FAIL).
3. Implement rule extensions + clear allowlist constants.
4. Re-run same test command (expect PASS).
5. Update status (`O01-S01`) and commit.

#### Subtask O01-S02: Add forbidden routing pattern tests
**Why:** Guard without realistic negative fixtures can regress silently.
**What:** Add test fixtures or inline test cases for phase alias hardcode, agent array hardcode, provider branch hardcode.
**How:** Build table-driven tests with path + content + expected violation reason.
**Key Design:** Violations should report actionable message: forbidden token + file path + suggested config-driven alternative.

**Files:**
- Modify: `internal/validation/model_hardcode_guard_test.go`
- Create/Modify: `internal/validation/testdata/*` (if fixture-based)

**Steps:** RED test -> verify fail -> implement detection -> verify pass -> status update + commit.

#### Subtask O01-S03: Wire guard into full validation suite and docs
**Why:** Guards are useless if not part of normal validation workflow.
**What:** Ensure guard is included in default `go test ./...` path and referenced in validation docs.
**How:** Add mention in validation readme/docs and ensure no skip tags bypass it.
**Key Design:** No opt-out flag in CI.

**Files:**
- Modify: `docs/architecture/*` relevant validation section
- Modify: `docs/CLI-REFERENCE.md` (if validation commands are listed)

---

### Task O02: Config governance + lint gate

#### Subtask O02-S01: Add JSON Schema for orchestration/workflows/providers
**Why:** Human-edited YAML is the highest-risk drift vector.
**What:** Define strict schema (required fields, enum values, array/object shapes, regex field checks).
**How:** Add schema files and schema load tests against valid/invalid fixtures.
**Key Design:** `workflows.yaml` holds flow overrides; globals remain in `orchestration.yaml`.

**Files:**
- Create: `config/schema/orchestration.schema.json`
- Create: `config/schema/workflows.schema.json`
- Create: `config/schema/providers.schema.json`
- Modify: `internal/policy/*load*` tests to validate schema path

**Steps:** fixture RED tests -> schema implementation -> schema validation pass -> status + commit.

#### Subtask O02-S02: Implement `mp-devx lint-config`
**Why:** Schema exists but must be executable via tooling for developers and CI.
**What:** Add `lint-config` action to parse YAML + validate schema + cross-reference models/providers/fallbacks.
**How:** Extend devx action router and return machine-readable errors.
**Key Design:** Fail-fast with line/file context.

**Files:**
- Modify: `cmd/mp-devx/main.go`
- Modify: `internal/devx/*` action runner files
- Create/Modify tests in `cmd/mp-devx/main_test.go` and `internal/devx/*_test.go`

#### Subtask O02-S03: Enforce pre-commit/CI lint gate
**Why:** Optional lint command will be skipped under delivery pressure.
**What:** Add hook/CI enforcement before merge.
**How:** Wire `mp-devx lint-config` into pre-commit and workflow checks.
**Key Design:** One command for both local and CI parity.

**Files:**
- Modify: `.github/workflows/*` relevant pipeline
- Modify: `.claude/hooks/pre-commit.sh` (or equivalent existing hook entry)
- Modify docs under `docs/architecture/*` for contributor workflow

---

### Task O03: Standardized error semantics

#### Subtask O03-S01: Introduce typed orchestration error codes
**Why:** Free-form strings create inconsistent behavior across runtime layers.
**What:** Replace ad-hoc status strings with typed constants (`E_DISPATCH_TIMEOUT`, `E_PROVIDER_UNAVAILABLE`, etc.).
**How:** Add enum-like code definitions and migrate result structs to use them.
**Key Design:** Backward-compatible JSON fields: preserve `message`, add canonical `code`.

**Files:**
- Modify: `internal/orchestration/result_types.go`
- Modify: `internal/orchestration/report.go`
- Modify related tests in `internal/orchestration/*_test.go`

#### Subtask O03-S02: Map error codes to CLI exit codes
**Why:** Operators need deterministic shell-level automation behavior.
**What:** Define mapping table from orchestration code -> process exit code.
**How:** Implement centralized mapper used by CLI command handlers.
**Key Design:** Same failure class always returns same exit code.

**Files:**
- Modify: `internal/cli/root.go`
- Create/Modify: `internal/cli/error_codes.go`
- Modify tests: `internal/cli/root_test.go`

#### Subtask O03-S03: Add compatibility tests across workflows/cli/hooks
**Why:** Error semantics must remain consistent across entrypoints.
**What:** Add integration-style tests ensuring workflow return code aligns with CLI and hooks.
**How:** Use mocked dispatch failures for deterministic assertions.
**Key Design:** One source of truth for code mapping.

**Files:**
- Modify: `internal/workflows/*_test.go`
- Modify: `internal/hooks/*_test.go`
- Modify: `internal/cli/root_test.go`

---

### Task O04: Executor retry reliability

#### Subtask O04-S01: Add idempotent retry policy config
**Why:** Retry safety depends on idempotency and explicit policy.
**What:** Add step-level fields: `idempotent`, `max_retries`, `backoff_ms`, `jitter_ratio`, `retryable_codes`.
**How:** Extend plan structs + config merge rules with defaults.
**Key Design:** Non-idempotent steps must never auto-retry.

**Files:**
- Modify: `internal/orchestration/plan_types.go`
- Modify: `internal/orchestration/merge.go`
- Modify: `internal/orchestration/load.go`
- Tests in `internal/orchestration/merge_test.go` and `load_test.go`

#### Subtask O04-S02: Implement retry loop with exponential backoff
**Why:** 429/503 transient failures currently degrade too quickly.
**What:** Retry only eligible failures with bounded attempts and context-aware sleep.
**How:** Refactor `executeStep` path to wrap dispatch call in retry controller.
**Key Design:** Respect context cancel/deadline immediately.

**Files:**
- Modify: `internal/orchestration/executor.go`
- Modify: `internal/orchestration/result_types.go` (attempt metadata)
- Tests: `internal/orchestration/executor_test.go`

#### Subtask O04-S03: Deterministic retry test matrix
**Why:** Retry logic is easy to break during refactor.
**What:** Add tests for success-after-retry, max retry exhausted, non-retryable immediate fail, canceled context.
**How:** Use fake clock or injectable sleeper for deterministic timing.
**Key Design:** No flaky sleep-based tests.

**Files:**
- Modify: `internal/orchestration/executor_test.go`

---

### Task O05: Traceability and observability

#### Subtask O05-S01: Generate and propagate `trace_id`
**Why:** Parallel execution is hard to diagnose without correlation IDs.
**What:** Create `trace_id` once per orchestration run and propagate to every step/synthesis dispatch.
**How:** Extend plan/result/context structs to carry trace metadata.
**Key Design:** `trace_id` immutable for a run; `step_id` unique per step.

**Files:**
- Modify: `internal/orchestration/planner.go`
- Modify: `internal/orchestration/executor.go`
- Modify: `internal/policy/dispatch.go` (request metadata passthrough if needed)

#### Subtask O05-S02: Structured step lifecycle logs
**Why:** Raw logs are not machine-analyzable during incidents.
**What:** Emit `started/completed/failed/fallback` events with JSON fields.
**How:** Reuse/extend `events.go` and add writer hook to `.claude-octopus/logs`.
**Key Design:** Never log prompt raw body by default; include hashes/lengths only.

**Files:**
- Modify: `internal/orchestration/events.go`
- Modify: `internal/orchestration/executor.go`
- Modify docs for log schema

#### Subtask O05-S03: Correlation tests
**Why:** Trace propagation can silently drop across layers.
**What:** Assert same `trace_id` in all step results and event entries.
**How:** Add tests around executor + synthesis path.
**Key Design:** Test both success and fallback paths.

**Files:**
- Modify: `internal/orchestration/executor_test.go`
- Modify: `internal/orchestration/report_test.go`

---

### Task O06: Merge explainability CLI

#### Subtask O06-S01: Build explain resolver model
**Why:** Users cannot debug config precedence from opaque final routing.
**What:** Build model that returns each effective field with source (`global`, `workflow`, `task`).
**How:** Extend merge output to include source refs alongside values.
**Key Design:** Explain path is read-only and side-effect free.

**Files:**
- Modify: `internal/orchestration/merge.go`
- Create: `internal/orchestration/explain.go`
- Tests: `internal/orchestration/merge_test.go`, `explain_test.go`

#### Subtask O06-S02: Add `mp orchestrate explain`
**Why:** Explain data must be user-facing.
**What:** New CLI command accepting `--flow`, optional `--task`, optional `--json`.
**How:** Hook into CLI router and render explain model.
**Key Design:** Non-JSON output should still show value + source clearly.

**Files:**
- Modify: `internal/cli/root.go`
- Modify: `internal/cli/root_test.go`

#### Subtask O06-S03: Explain docs and negative tests
**Why:** Wrong/missing flow/task should fail with remediation guidance.
**What:** Add invalid input tests and CLI reference docs.
**How:** Include examples for common debug questions.
**Key Design:** Errors include nearest valid flow suggestions.

**Files:**
- Modify: `docs/CLI-REFERENCE.md`
- Modify: `internal/cli/root_test.go`

---

### Task O07: Result caching

#### Subtask O07-S01: Cache key and store abstraction
**Why:** Repeated runs with same inputs waste tokens and time.
**What:** Define cache key = `flow + task + prompt_hash + config_checksum (+ model hint)`.
**How:** Add minimal cache interface and filesystem-backed implementation.
**Key Design:** Cache payload includes timestamp + schema version.

**Files:**
- Create: `internal/orchestration/cache.go`
- Create: `internal/orchestration/cache_fs.go`
- Tests: `internal/orchestration/cache_test.go`

#### Subtask O07-S02: Integrate read/write cache in runtime
**Why:** Cache has no value unless integrated at orchestration entrypoint.
**What:** Before execute, try cache hit; after success/degraded, persist result.
**How:** Wire caching into workflow adapter/orchestration facade.
**Key Design:** Include `cache_hit` metadata in output.

**Files:**
- Modify: `internal/orchestration/workflow_adapter.go`
- Modify: `internal/workflows/*.go`
- Tests: `internal/orchestration/e2e_test.go`

#### Subtask O07-S03: TTL and invalidation tests
**Why:** Stale cache can cause incorrect behavior.
**What:** Add TTL config (default 10m) and invalidation on config checksum mismatch.
**How:** Add deterministic tests with fake time.
**Key Design:** Expired cache should behave as miss with audit reason.

**Files:**
- Modify: `internal/orchestration/cache.go`
- Modify: `internal/orchestration/cache_test.go`

---

### Task O08: Planner visualization

#### Subtask O08-S01: DAG export from ExecutionPlan
**Why:** Complex plan structure is hard to review textually.
**What:** Convert plan phases/steps/dependencies into DAG model.
**How:** Add pure transform helper and unit tests.
**Key Design:** Renderer-independent graph model.

**Files:**
- Create: `internal/orchestration/plan_graph.go`
- Tests: `internal/orchestration/planner_test.go`

#### Subtask O08-S02: Mermaid renderer + CLI
**Why:** Mermaid is easy to embed in docs/PR comments.
**What:** Add `--format mermaid` output option for plan visualization.
**How:** Implement renderer + wire command path.
**Key Design:** Keep node IDs stable for diffability.

**Files:**
- Create: `internal/orchestration/plan_graph_mermaid.go`
- Modify: `internal/cli/root.go`
- Tests: `internal/cli/root_test.go`, `internal/orchestration/*_test.go`

#### Subtask O08-S03: Visualization docs
**Why:** Feature adoption requires usage examples.
**What:** Add CLI examples and generated sample graph.
**How:** Update architecture and CLI docs with copy-paste snippets.
**Key Design:** One simple example per flow class.

**Files:**
- Modify: `docs/CLI-REFERENCE.md`
- Modify: `docs/architecture/*` orchestration sections

---

### Task O09: Layered regression tests (golden)

#### Subtask O09-S01: Plan snapshot golden tests
**Why:** Planner refactors risk silent shape drift.
**What:** Snapshot resolved plan JSON for all 6 flows.
**How:** Build golden harness with explicit update switch.
**Key Design:** Golden updates require deliberate flag.

**Files:**
- Create: `internal/orchestration/testdata/golden/plan/*.golden.json`
- Modify: `internal/orchestration/planner_test.go`

#### Subtask O09-S02: Synthesis report golden tests
**Why:** Report format drift breaks downstream consumers.
**What:** Snapshot synthesis report outputs under deterministic mocks.
**How:** Table tests for normal and partial results.
**Key Design:** Strip nondeterministic fields (timestamps, IDs) before compare.

**Files:**
- Create: `internal/orchestration/testdata/golden/report/*.golden.json`
- Modify: `internal/orchestration/synthesis_final_test.go`

#### Subtask O09-S03: Degraded/fallback golden tests
**Why:** Degraded mode contracts are most likely to regress.
**What:** Snapshot fallback chains, error codes, partial outputs.
**How:** Mock dispatch failures by step/provider.
**Key Design:** Ensure one-hop fallback semantics remain intact.

**Files:**
- Create: `internal/orchestration/testdata/golden/degraded/*.golden.json`
- Modify: `internal/orchestration/e2e_test.go`

---

### Task O10: Verification evidence and rollout docs

#### Subtask O10-S01: Capture optimization verification evidence
**Why:** Completion claims need reproducible command evidence.
**What:** Record commands, exit codes, timestamps, first output lines.
**How:** Run required suite and save into evidence markdown.
**Key Design:** Evidence generated from clean current HEAD.

**Files:**
- Create: `docs/plans/evidence/model-routing/2026-03-04-go-orchestration-optimization-verification.md`

**Required commands:**
```bash
go test ./internal/validation ./internal/orchestration ./internal/policy ./internal/workflows ./internal/cli -count=1
go test ./...
```

#### Subtask O10-S02: Update architecture/ops docs
**Why:** New contracts (error code map, explain, retry, trace, cache) must be documented.
**What:** Add field-level docs and operator guidance.
**How:** Update architecture and CLI references with examples.
**Key Design:** Docs must describe behavior, precedence, and failure handling with concrete examples.

**Files:**
- Modify: `docs/architecture/*` relevant orchestration docs
- Modify: `docs/CLI-REFERENCE.md`
- Modify: `README.md` orchestration section

---

## Mandatory Verification Gate (Before completion claim)

Run and pass:

```bash
go test ./internal/validation ./internal/orchestration ./internal/policy ./internal/workflows ./internal/cli -count=1
go test ./...
```

If any fails, related task/subtask cannot be marked `COMPLETED`.

## Commit Strategy (Required)

- One commit per subtask or tightly-coupled pair only.
- Commit message must include subtask ID, for example:

```bash
git commit -m "feat(orchestration): add idempotent retry controller (O04-S02)"
```

- After each commit, immediately update status table + `mp state set` key.
