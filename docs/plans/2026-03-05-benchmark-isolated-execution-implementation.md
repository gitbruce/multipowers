# External Command Isolation (Shared) Enhancement Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Enforce shared per-model/task isolated `worktree + branch` execution at runtime whenever external commands may edit files, with benchmark mode as one profile, and integrate critic top-1 with one same-model repair retry on failure.

**Architecture:** Add shared isolation policy/runtime components used by all external-command flows in orchestration. Every eligible candidate runs in a dedicated sandbox branch/worktree. Executor emits model progress events with heartbeats. A sync gate (waitgroup + timeout) advances with completed candidates when timeout occurs. Benchmark becomes one profile on top of shared logic. Critic ranking chooses top-1 for integration branch; failures trigger one same-model repair retry.

**Tech Stack:** Go, YAML config, git worktree CLI, JSONL persistence, existing `internal/orchestration`, `internal/isolation` (new shared package), profile adapters (including benchmark), `go test`.

---

## Mandatory Execution Rules

1. Use TDD for every task: failing test -> minimal implementation -> passing test.
2. Keep commits small: one commit per task.
3. **When executing this plan, update task status in real time:**
- Set task `Status` to `IN_PROGRESS` before any code change.
- Set task `Status` to `DONE` immediately after verification and commit.
- If blocked, set `Status` to `BLOCKED` and record blocker reason under the task.
4. Allowed status values: `TODO`, `IN_PROGRESS`, `BLOCKED`, `DONE`.
5. Keep task table and per-task `Status` field synchronized.
6. No destructive git commands.

## Required Skills and Guardrails

- `@superpowers:using-git-worktrees`
- `@superpowers:test-driven-development`
- `@superpowers:verification-before-completion`
- `@superpowers:executing-plans`

## Task Status Board

| Task ID | Task Name | Status | Last Update |
|---|---|---|---|
| E1 | Config Schema: Shared Execution Isolation + Sync Gate | DONE | 2026-03-05 02:21:40 CST |
| E2 | Shared Isolation Policy Resolver | DONE | 2026-03-05 02:26:44 CST |
| E3 | Shared Git Worktree Runtime Manager | DONE | 2026-03-05 02:30:18 CST |
| E4 | Event-Driven Model Progress + Heartbeats | TODO | - |
| E5 | Sync Gate Collector with Timeout Degradation | TODO | - |
| E6 | Critic Top-1 Deterministic Selection | TODO | - |
| E7 | Integration Branch + Same-Model Repair Retry | TODO | - |
| E8 | JSONL/Docs/E2E Verification | TODO | - |

---

### Task E1: Config Schema: Shared Execution Isolation + Sync Gate

**Status:** DONE

**Why**
- Runtime enforcement requires explicit config to avoid hardcoded behavior.
- Long-running fan-out needs configurable timeout/proceed policy.

**What**
- Add shared `execution_isolation` root config and profile override support.
- Add defaults and validation.

**How**
- Add `ExecutionIsolationConfig` to orchestration types at root scope.
- Add optional profile override parsing for benchmark profile.
- Parse and default fields in config loader.
- Validate enum/range values.

**Key Design**
- Backward compatible: missing section keeps current behavior.
- Safe defaults: isolation off, proceed policy deterministic.

**Files:**
- Modify: `internal/orchestration/types.go`
- Modify: `internal/orchestration/load.go`
- Modify: `internal/orchestration/config_test.go`
- Modify: `config/orchestration.yaml`

**Step 0: Status update**
- Set task `E1` to `IN_PROGRESS` in board and this section.

**Step 1: Write failing tests**

```go
func TestLoadOrchestrationConfig_ExecutionIsolation(t *testing.T) {
    // assert execution_isolation fields are loaded/defaulted
}

func TestLoadOrchestrationConfig_ExecutionIsolationValidation(t *testing.T) {
    // invalid proceed_policy or timeout should fail
}
```

**Step 2: Run failing tests**

Run: `go test ./internal/orchestration -run TestLoadOrchestrationConfig_ExecutionIsolation -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
type ExecutionIsolationConfig struct {
    Enabled                  bool     `yaml:"enabled"`
    CommandWhitelist         []string `yaml:"command_whitelist,omitempty"`
    BranchPrefix             string   `yaml:"branch_prefix,omitempty"`
    WorktreeRoot             string   `yaml:"worktree_root,omitempty"`
    RepairRetryMax           int      `yaml:"repair_retry_max,omitempty"`
    GlobalTimeoutMs          int      `yaml:"global_timeout_ms,omitempty"`
    ProceedPolicy            string   `yaml:"proceed_policy,omitempty"`
    MinCompletedModels       int      `yaml:"min_completed_models,omitempty"`
    HeartbeatIntervalSeconds int      `yaml:"heartbeat_interval_seconds,omitempty"`
    LogsSubdir               string   `yaml:"logs_subdir,omitempty"`
}
```

**Step 4: Re-run tests**

Run: `go test ./internal/orchestration -run TestLoadOrchestrationConfig_ExecutionIsolation -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/orchestration/types.go internal/orchestration/load.go internal/orchestration/config_test.go config/orchestration.yaml
git commit -m "feat(config): add shared execution isolation and sync gate schema"
```

**Step 6: Status update**
- Set task `E1` to `DONE` with timestamp.

---

### Task E2: Shared Isolation Policy Resolver

**Status:** DONE

**Why**
- Enforcement must be deterministic and auditable.

**What**
- Implement shared policy decision for isolation activation in any external-command flow.

**How**
- Build pure resolver from isolation toggle + external command involvement + file-edit intent + optional profile gates.
- Return decision with reason fields.

**Key Design**
- No side effects in policy resolver.
- Normalized command matching for predictable behavior.

**Files:**
- Create: `internal/isolation/policy.go`
- Create: `internal/isolation/policy_test.go`
- Create: `internal/isolation/profile_benchmark.go`
- Modify: `internal/hooks/handler.go`

**Step 0: Status update**
- Set task `E2` to `IN_PROGRESS`.

**Step 1: Write failing tests**

```go
func TestResolveIsolationPolicy(t *testing.T) {
    // enforce true only on full condition match
}
```

**Step 2: Run failing tests**

Run: `go test ./internal/isolation -run TestResolveIsolationPolicy -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
type IsolationPolicyInput struct {
    IsolationEnabled bool
    ExternalCommand  bool
    MayEditFiles     bool
    CodeRelated      bool
    Command          string
    Whitelist        []string
}

func ResolveIsolationPolicy(in IsolationPolicyInput) IsolationPolicyDecision {
    // deterministic result + reason
}
```

**Step 4: Re-run tests**

Run: `go test ./internal/isolation -run TestResolveIsolationPolicy -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/isolation/policy.go internal/isolation/policy_test.go internal/isolation/profile_benchmark.go internal/hooks/handler.go
git commit -m "feat(isolation): add shared runtime isolation policy resolver"
```

**Step 6: Status update**
- Set task `E2` to `DONE` with timestamp.

---

### Task E3: Shared Git Worktree Runtime Manager

**Status:** DONE

**Why**
- Hard isolation requires per-model sandbox lifecycle control.

**What**
- Create shared runtime manager to create/cleanup command sandboxes.

**How**
- Add git runner abstraction.
- Create branch/worktree naming conventions.
- Support logs directory initialization.

**Key Design**
- Never mutate active developer worktree branch.
- Always use bounded cleanup and error surfaces.

**Files:**
- Create: `internal/isolation/runtime.go`
- Create: `internal/isolation/runtime_test.go`

**Step 0: Status update**
- Set task `E3` to `IN_PROGRESS`.

**Step 1: Write failing tests**

```go
func TestIsolationRuntime_CreateModelSandbox(t *testing.T) {}
func TestIsolationRuntime_CleanupModelSandbox(t *testing.T) {}
```

**Step 2: Run failing tests**

Run: `go test ./internal/isolation -run TestIsolationRuntime -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
type ModelSandbox struct {
    Model        string
    Branch       string
    WorktreePath string
    LogsPath     string
}

func (r RuntimeManager) CreateModelSandbox(runID, model, baseRef string) (ModelSandbox, error) {
    // git worktree add -b <branch> <path> <baseRef>
}
```

**Step 4: Re-run tests**

Run: `go test ./internal/isolation -run TestIsolationRuntime -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/isolation/runtime.go internal/isolation/runtime_test.go
git commit -m "feat(isolation): add shared worktree runtime manager"
```

**Step 6: Status update**
- Set task `E3` to `DONE` with timestamp.

---

### Task E4: Event-Driven Model Progress + Heartbeats

**Status:** TODO

**Why**
- Long-running model tasks need continuous user-visible status.

**What**
- Emit model progress events from executor.
- Add periodic heartbeat updates.

**How**
- Reuse `EventTypeStepProgress` with structured `Event.Data` payload.
- Emit `queued/sandbox_ready/running/completed/failed/timeout/repair_retry` transitions.
- Emit heartbeat every configured interval.

**Key Design**
- Event emission is non-blocking and drop-tolerant.
- Progress events should not affect execution correctness.

**Files:**
- Modify: `internal/orchestration/events.go`
- Modify: `internal/orchestration/executor.go`
- Modify: `internal/orchestration/executor_test.go`

**Step 0: Status update**
- Set task `E4` to `IN_PROGRESS`.

**Step 1: Write failing tests**

```go
func TestExecutor_EmitsModelProgressEvents(t *testing.T) {
    // expects step_progress events with model payload and heartbeat
}
```

**Step 2: Run failing tests**

Run: `go test ./internal/orchestration -run TestExecutor_EmitsModelProgressEvents -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
type BenchmarkProgress struct {
    RunID       string
    Model       string
    Status      string
    Percent     int
    Branch      string
    Worktree    string
    HeartbeatAt time.Time
}
```

**Step 4: Re-run tests**

Run: `go test ./internal/orchestration -run TestExecutor_EmitsModelProgressEvents -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/orchestration/events.go internal/orchestration/executor.go internal/orchestration/executor_test.go
git commit -m "feat(orchestration): emit model progress and heartbeat events for isolated external commands"
```

**Step 6: Status update**
- Set task `E4` to `DONE` with timestamp.

---

### Task E5: Sync Gate Collector with Timeout Degradation

**Status:** TODO

**Why**
- One slow model must not block the whole command indefinitely.

**What**
- Add waitgroup-based sync gate with timeout.
- Proceed with completed candidates based on policy.

**How**
- Use `sync.WaitGroup` + `context.WithTimeout`.
- Mark unfinished candidates as `timeout`.
- Support proceed policies (`all_done`, `all_or_timeout`, `majority_or_timeout`).

**Key Design**
- Graceful degradation: use available candidates rather than fail full run.
- Gate applies only to candidate collection/ranking boundary.

**Files:**
- Modify: `internal/orchestration/executor.go`
- Create: `internal/isolation/sync_gate.go`
- Create: `internal/isolation/sync_gate_test.go`

**Step 0: Status update**
- Set task `E5` to `IN_PROGRESS`.

**Step 1: Write failing tests**

```go
func TestSyncGate_ProceedWithCompletedCandidatesOnTimeout(t *testing.T) {}
func TestSyncGate_MajorityPolicy(t *testing.T) {}
```

**Step 2: Run failing tests**

Run: `go test ./internal/isolation -run TestSyncGate -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
func WaitForCandidates(ctx context.Context, in SyncGateInput) SyncGateResult {
    // waitgroup + timeout + timeout marking + policy decision
}
```

**Step 4: Re-run tests**

Run: `go test ./internal/isolation -run TestSyncGate -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/isolation/sync_gate.go internal/isolation/sync_gate_test.go internal/orchestration/executor.go
git commit -m "feat(isolation): add shared sync gate timeout degradation policies"
```

**Step 6: Status update**
- Set task `E5` to `DONE` with timestamp.

---

### Task E6: Critic Top-1 Deterministic Selection

**Status:** TODO

**Why**
- Integration choice must be deterministic and explainable.

**What**
- Implement top-1 selector with tie-break rules.

**How**
- Sort by weighted score desc, failures asc, duration asc, model lexical asc.

**Key Design**
- Deterministic output for same input set.
- No side effects or IO in selector.

**Files:**
- Create: `internal/isolation/critic_selection.go`
- Create: `internal/isolation/critic_selection_test.go`

**Step 0: Status update**
- Set task `E6` to `IN_PROGRESS`.

**Step 1: Write failing tests**

```go
func TestSelectTopCandidate(t *testing.T) {
    // validates ordering and tie-break determinism
}
```

**Step 2: Run failing tests**

Run: `go test ./internal/isolation -run TestSelectTopCandidate -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
func SelectTopCandidate(candidates []CandidateScore) (CandidateScore, error) {
    // deterministic sort and top-1 selection
}
```

**Step 4: Re-run tests**

Run: `go test ./internal/isolation -run TestSelectTopCandidate -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/isolation/critic_selection.go internal/isolation/critic_selection_test.go
git commit -m "feat(isolation): add deterministic critic top-1 selection"
```

**Step 6: Status update**
- Set task `E6` to `DONE` with timestamp.

---

### Task E7: Integration Branch + Same-Model Repair Retry

**Status:** TODO

**Why**
- Top-1 merge/gate failures need deterministic automated recovery.

**What**
- Implement integration flow with one same-model repair retry.

**How**
- Create integration branch.
- Try merge/cherry-pick top-1.
- On failure, run same-model repair prompt once and retry once.

**Key Design**
- No fallback to rank #2 in this phase.
- Explicit status model (`merged`, `repair_retry`, `failed_after_retry`).

**Files:**
- Create: `internal/isolation/integration_flow.go`
- Create: `internal/isolation/integration_flow_test.go`
- Modify: `internal/orchestration/synthesis_final_test.go`

**Step 0: Status update**
- Set task `E7` to `IN_PROGRESS`.

**Step 1: Write failing tests**

```go
func TestIntegrateTopCandidate_RepairRetryOnce(t *testing.T) {}
```

**Step 2: Run failing tests**

Run: `go test ./internal/isolation -run TestIntegrateTopCandidate_RepairRetryOnce -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
type IntegrationResult struct {
    Status          string
    RepairRetryUsed int
}

func IntegrateTopCandidate(in Input) IntegrationResult {
    // one retry on same model only
}
```

**Step 4: Re-run tests**

Run: `go test ./internal/isolation -run TestIntegrateTopCandidate_RepairRetryOnce -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/isolation/integration_flow.go internal/isolation/integration_flow_test.go internal/orchestration/synthesis_final_test.go
git commit -m "feat(isolation): add top-1 integration with same-model single retry"
```

**Step 6: Status update**
- Set task `E7` to `DONE` with timestamp.

---

### Task E8: JSONL/Docs/E2E Verification

**Status:** TODO

**Why**
- Enhancement must be observable, documented, and verifiable end-to-end.

**What**
- Extend JSONL records for isolation/progress/integration metadata.
- Update user docs and add e2e tests.

**How**
- Add `isolation_runs` stream and extend existing streams.
- Add e2e tests for benchmark profile and at least one non-benchmark external-command profile.
- Update README/TRIGGERS and finalize verification.

**Key Design**
- Persistence remains best-effort and non-blocking.
- Verification command output is required before completion claims.

**Files:**
- Modify: `internal/benchmark/records.go`
- Modify: `internal/benchmark/store_jsonl.go`
- Modify: `internal/benchmark/store_jsonl_test.go`
- Modify: `internal/benchmark/integration_test.go`
- Create/Modify: `internal/isolation/e2e_profile_test.go`
- Modify: `README.md`
- Modify: `docs/TRIGGERS.md`
- Modify: `docs/plans/2026-03-05-benchmark-isolated-execution-design.md`

**Step 0: Status update**
- Set task `E8` to `IN_PROGRESS`.

**Step 1: Write failing tests**

```go
func TestJSONLStore_IsolationAndProgressStreams(t *testing.T) {}
func TestBenchmarkIsolationEnforcedForWhitelistedCommand(t *testing.T) {}
func TestExternalCommandIsolationNonBenchmarkProfile(t *testing.T) {}
```

**Step 2: Run failing tests**

Run: `go test ./internal/benchmark ./internal/isolation -run 'TestJSONLStore_IsolationAndProgressStreams|TestBenchmarkIsolationEnforcedForWhitelistedCommand|TestExternalCommandIsolationNonBenchmarkProfile' -v`
Expected: FAIL.

**Step 3: Minimal implementation & docs wiring**
- Implement schema/store updates.
- Update README and TRIGGERS.

**Step 4: Re-run targeted suites**

Run: `go test ./internal/benchmark ./internal/orchestration ./internal/modelroute -v`
Expected: PASS.

**Step 5: Full verification**

Run: `go test ./...`
Expected: PASS.

**Step 6: Commit**

```bash
git add internal/isolation internal/benchmark internal/orchestration internal/modelroute README.md docs/TRIGGERS.md docs/plans/2026-03-05-benchmark-isolated-execution-design.md
git commit -m "feat(isolation): enforce shared isolated execution with progress sync gate and retry integration"
```

**Step 7: Status update**
- Set task `E8` to `DONE` with timestamp.
- Confirm all rows are `DONE`.

---

## Completion Checklist

- [ ] Task statuses are updated in real time during execution (`IN_PROGRESS` -> `DONE`/`BLOCKED`).
- [ ] Isolation policy is runtime-enforced for external commands that may edit files.
- [ ] Every eligible model run uses unique branch/worktree.
- [ ] Model progress + heartbeat feedback is visible via events.
- [ ] Sync gate timeout proceeds with completed candidates (graceful degradation).
- [ ] Critic top-1 direct integration works deterministically.
- [ ] Same-model repair retry executes exactly once on merge/gate failure.
- [ ] JSONL records include isolation/progress/integration metadata.
- [ ] Benchmark and non-benchmark profiles both use the same shared isolation runtime/policy.
- [ ] All tests pass.
