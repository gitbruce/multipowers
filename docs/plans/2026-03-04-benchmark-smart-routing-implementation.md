# Benchmark + Smart Routing Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build async benchmark fan-out + judge scoring + daily JSONL history + optional smart-routing override for `/mp:*` code-related requests.

**Architecture:** Add an async benchmark pipeline beside orchestration. Main execution path only emits events and never waits for storage/scoring/routing analytics. Background workers classify intent, persist JSONL, score outputs, and compute history-based routing overrides.

**Tech Stack:** Go, YAML config, JSONL file IO, existing `internal/orchestration`, `internal/modelroute`, `go test`.

---

## Mandatory Execution Rules

1. Use TDD for every task: failing test -> minimal implementation -> passing test.
2. Keep commits small: one commit per task.
3. **After each task is completed, update its status immediately in this file before starting the next task.**
4. Allowed status values: `TODO`, `IN_PROGRESS`, `BLOCKED`, `DONE`.
5. If blocked, record blocker reason under that task.

## Task Status Board (Must Keep Updated)

| Task ID | Task Name | Status | Last Update |
|---|---|---|---|
| T1 | Config Schema for Benchmark/Smart Routing | DONE | 2026-03-04 21:55:09 CST |
| T2 | Code-Intent Classification Contract | DONE | 2026-03-04 21:59:00 CST |
| T3 | Force-All-Models Routing Override | DONE | 2026-03-04 22:01:49 CST |
| T4 | Async Queue (Non-Blocking) | DONE | 2026-03-04 22:04:15 CST |
| T5 | Daily JSONL Store | DONE | 2026-03-04 22:06:58 CST |
| T6 | Judge Scoring Worker | TODO | - |
| T7 | Smart Routing from History | TODO | - |
| T8 | E2E Failure Isolation + Docs | TODO | - |

## Shared Context for Junior Developers

1. Main command handling currently enters through `/mp:*` route and orchestration modules under `internal/orchestration`.
2. Existing model selection exists in `internal/modelroute`.
3. New code should live under `internal/benchmark` to keep concerns isolated.
4. “Do not affect main execution performance” means no blocking writes and no hard dependency on benchmark workers.
5. Any benchmark failure is recoverable and must not fail user-facing run results.

### Task T1: Config Schema for Benchmark/Smart Routing

**Status:** DONE

**Why**
- We need config-first behavior switches (`benchmark_mode`, `smart_routing`) to avoid hardcoded runtime logic.

**What**
- Extend orchestration config structs and YAML loading.
- Add default values and validation for sample gate and toggles.

**How**
- Add new config structs in orchestration config types.
- Update config loader tests with benchmark-enabled fixture.
- Add minimal validation for `min_samples_per_model >= 1`.

**Key Design**
- Keep backward compatibility: missing new config sections should not break existing workflows.

**Files:**
- Modify: `internal/orchestration/types.go`
- Modify: `config/orchestration.yaml`
- Create/Modify: `internal/orchestration/config_test.go`
- Create: `internal/orchestration/testdata/orchestration_with_benchmark.yaml`

**Step 1: Write failing test**
Run: `go test ./internal/orchestration -run TestLoadOrchestrationConfig_BenchmarkAndSmartRouting -v`
Expected: FAIL.

**Step 2: Minimal implementation**
- Add `BenchmarkModeConfig` and `SmartRoutingConfig` fields in root config.

**Step 3: Re-run test**
Run: `go test ./internal/orchestration -run TestLoadOrchestrationConfig_BenchmarkAndSmartRouting -v`
Expected: PASS.

**Step 4: Commit**
```bash
git add internal/orchestration/types.go internal/orchestration/config_test.go internal/orchestration/testdata/orchestration_with_benchmark.yaml config/orchestration.yaml
git commit -m "feat(config): add benchmark and smart-routing schema"
```

**Step 5: REQUIRED status update**
- Update board row `T1` to `DONE`.
- Set `Last Update` timestamp.

### Task T2: Code-Intent Classification Contract

**Status:** DONE

**Why**
- Benchmark fan-out should happen only for code-related requests.

**What**
- Implement classification API combining whitelist hits and LLM semantic decision.
- Enforce product rule: LLM final decision priority.

**How**
- Create intent request/decision structs.
- Implement deterministic decision function.
- Add tests for: whitelist hit + LLM true/false conflict cases.

**Key Design**
- Keep classifier pure and deterministic for easy testing; external LLM call wrapped outside pure function.

**Files:**
- Create: `internal/benchmark/intent.go`
- Create: `internal/benchmark/intent_test.go`
- Modify: `internal/hooks/handler.go`

**Step 1: Write failing tests**
Run: `go test ./internal/benchmark -run TestClassifyCodeIntent -v`
Expected: FAIL.

**Step 2: Minimal implementation**
- Implement `ClassifyCodeIntent(req IntentRequest) IntentDecision`.

**Step 3: Re-run tests**
Run: `go test ./internal/benchmark -run TestClassifyCodeIntent -v`
Expected: PASS.

**Step 4: Commit**
```bash
git add internal/benchmark/intent.go internal/benchmark/intent_test.go internal/hooks/handler.go
git commit -m "feat(benchmark): add intent classification contract"
```

**Step 5: REQUIRED status update**
- Update board row `T2` to `DONE`.

### Task T3: Force-All-Models Routing Override

**Status:** DONE

**Why**
- When benchmark mode is on and intent is code-related, all available models must run regardless of default route.

**What**
- Add override resolver that returns full model list.
- Integrate resolver at routing decision point.

**How**
- Implement `ResolveForcedCandidates(...)`.
- Inject into orchestration candidate selection path.
- Add tests for enabled/disabled and code/non-code branches.

**Key Design**
- Override must be scoped to current request only; no global mutable state.

**Files:**
- Create: `internal/benchmark/override.go`
- Create: `internal/benchmark/override_test.go`
- Modify: `internal/orchestration/select.go`

**Step 1: Write failing tests**
Run: `go test ./internal/benchmark -run TestResolveForcedCandidates -v`
Expected: FAIL.

**Step 2: Minimal implementation**
- Return all models only when both config and intent condition are true.

**Step 3: Re-run tests**
Run: `go test ./internal/benchmark -run TestResolveForcedCandidates -v`
Expected: PASS.

**Step 4: Commit**
```bash
git add internal/benchmark/override.go internal/benchmark/override_test.go internal/orchestration/select.go
git commit -m "feat(routing): add force-all-models benchmark override"
```

**Step 5: REQUIRED status update**
- Update board row `T3` to `DONE`.

### Task T4: Async Queue (Non-Blocking Main Path)

**Status:** DONE

**Why**
- Benchmark should be fully async and must not increase critical-path latency.

**What**
- Create non-blocking queue and worker loop.
- Drop or defer events safely when overloaded.

**How**
- Add `TryEnqueue` API with `select default`.
- Add queue metrics counters (`enqueued`, `dropped`).
- Add tests for full queue non-blocking behavior.

**Key Design**
- Main path never calls blocking enqueue.

**Files:**
- Create: `internal/benchmark/queue.go`
- Create: `internal/benchmark/queue_test.go`
- Modify: `internal/orchestration/executor.go` (or actual run entrypoint)

**Step 1: Write failing tests**
Run: `go test ./internal/benchmark -run TestQueueNonBlocking -v`
Expected: FAIL.

**Step 2: Minimal implementation**
- Implement bounded channel queue and `TryEnqueue`.

**Step 3: Re-run tests**
Run: `go test ./internal/benchmark -run TestQueueNonBlocking -v`
Expected: PASS.

**Step 4: Commit**
```bash
git add internal/benchmark/queue.go internal/benchmark/queue_test.go internal/orchestration/executor.go
git commit -m "feat(benchmark): add bounded non-blocking async queue"
```

**Step 5: REQUIRED status update**
- Update board row `T4` to `DONE`.

### Task T5: Daily JSONL Store

**Status:** DONE

**Why**
- Product requirement changed from DB to local JSONL cross-project storage.

**What**
- Persist benchmark streams into `~/.claude-octopus/metrics/<stream>.<YYYY-MM-DD>.jsonl`.
- Define all event schemas.

**How**
- Add atomic append writer with file lock.
- Add schema structs for `runs/model_outputs/task_fingerprints/judge_scores/route_overrides/async_jobs/errors`.
- Add tests for file naming and valid JSON lines.

**Key Design**
- One JSON object per line; never multi-line records; recoverable by streaming readers.

**Files:**
- Create: `internal/benchmark/store_jsonl.go`
- Create: `internal/benchmark/store_jsonl_test.go`
- Create: `internal/benchmark/records.go`

**Step 1: Write failing tests**
Run: `go test ./internal/benchmark -run TestJSONLStore -v`
Expected: FAIL.

**Step 2: Minimal implementation**
- Implement daily partition path and append record.

**Step 3: Re-run tests**
Run: `go test ./internal/benchmark -run TestJSONLStore -v`
Expected: PASS.

**Step 4: Commit**
```bash
git add internal/benchmark/store_jsonl.go internal/benchmark/store_jsonl_test.go internal/benchmark/records.go
git commit -m "feat(storage): add daily JSONL benchmark store"
```

**Step 5: REQUIRED status update**
- Update board row `T5` to `DONE`.

### Task T6: Judge Scoring Worker

**Status:** TODO

**Why**
- Need consistent multi-dimensional quality scoring for model outputs.

**What**
- Add scoring worker that records 1-5 scores and weighted aggregate.
- Support configurable dimensions and weights.

**How**
- Implement score aggregation utility.
- Add judge output parser + validation (dimension value 1..5).
- Add tests for weighted score and invalid score handling.

**Key Design**
- Keep scoring math independent from LLM integration so it can be unit-tested.

**Files:**
- Create: `internal/benchmark/judge.go`
- Create: `internal/benchmark/judge_test.go`
- Modify: `internal/orchestration/result_types.go`

**Step 1: Write failing tests**
Run: `go test ./internal/benchmark -run TestComputeWeightedScore -v`
Expected: FAIL.

**Step 2: Minimal implementation**
- Implement `ComputeWeightedScore` and range validation.

**Step 3: Re-run tests**
Run: `go test ./internal/benchmark -run TestComputeWeightedScore -v`
Expected: PASS.

**Step 4: Commit**
```bash
git add internal/benchmark/judge.go internal/benchmark/judge_test.go internal/orchestration/result_types.go
git commit -m "feat(judge): add benchmark score worker and aggregation"
```

**Step 5: REQUIRED status update**
- Update board row `T6` to `DONE`.

### Task T7: Smart Routing from History

**Status:** TODO

**Why**
- Need automatic model override for similar scenarios when enabled.

**What**
- Read historical judge scores, group by similarity signature, apply sample gate (`N=10`), pick best average score.

**How**
- Build signature from `task_type+tech_features+framework+language`.
- Implement history reducer.
- Integrate into route selection only when `smart_routing.enabled=true`.

**Key Design**
- No history candidate meeting sample gate => no override.

**Files:**
- Create: `internal/benchmark/routing_history.go`
- Create: `internal/benchmark/routing_history_test.go`
- Modify: `internal/modelroute/route.go`

**Step 1: Write failing tests**
Run: `go test ./internal/benchmark -run TestSelectBestModelByHistory -v`
Expected: FAIL.

**Step 2: Minimal implementation**
- Implement selection with sample threshold check.

**Step 3: Re-run tests**
Run: `go test ./internal/benchmark -run TestSelectBestModelByHistory -v`
Expected: PASS.

**Step 4: Commit**
```bash
git add internal/benchmark/routing_history.go internal/benchmark/routing_history_test.go internal/modelroute/route.go
git commit -m "feat(routing): add history-based smart routing override"
```

**Step 5: REQUIRED status update**
- Update board row `T7` to `DONE`.

### Task T8: End-to-End Failure Isolation + Documentation

**Status:** TODO

**Why**
- Must guarantee benchmark subsystem errors never break user-facing execution.

**What**
- Add integration tests for worker failures and queue/store failures.
- Document new toggles and behavior.

**How**
- Inject failures in async components and assert main execution still succeeds.
- Update README and trigger docs with benchmark and smart routing semantics.

**Key Design**
- “Best effort only” behavior is test-enforced contract, not just documentation.

**Files:**
- Create: `internal/benchmark/integration_test.go`
- Modify: `internal/orchestration/synthesis_final_test.go`
- Modify: `README.md`
- Modify: `docs/TRIGGERS.md`

**Step 1: Write failing tests**
Run: `go test ./internal/benchmark -run TestBenchmarkFailureDoesNotFailMainFlow -v`
Expected: FAIL.

**Step 2: Minimal implementation**
- Wrap benchmark emit path with swallow-and-log error handling.

**Step 3: Re-run tests**
Run: `go test ./internal/benchmark -run TestBenchmarkFailureDoesNotFailMainFlow -v`
Expected: PASS.

**Step 4: Full verification**
Run: `go test ./internal/benchmark ./internal/orchestration ./internal/modelroute -v`
Expected: PASS.

**Step 5: Commit**
```bash
git add internal/benchmark internal/orchestration internal/modelroute README.md docs/TRIGGERS.md
git commit -m "feat(benchmark): complete async benchmark pipeline with failure isolation"
```

**Step 6: REQUIRED status update**
- Update board row `T8` to `DONE`.
- Confirm all rows are `DONE`.

## Completion Checklist

- [ ] All task statuses updated in board.
- [ ] All commits created task-by-task.
- [ ] All listed test commands passing.
- [ ] Docs updated for new config toggles.
- [ ] No benchmark failure path can fail main run.
