# Benchmark Isolated Execution Enhancement Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Enforce per-model isolated `worktree + branch` execution at runtime when benchmark mode is enabled, code intent is true, and command whitelist matches; then integrate critic top-1 with one same-model repair retry on failure.

**Architecture:** Add an isolation policy gate and git-isolation runtime manager inside benchmark/orchestration flow. Every eligible model candidate runs in a dedicated sandbox branch/worktree. Critic ranking chooses top-1 for integration branch. If top-1 integration fails, run one same-model repair attempt and retry once.

**Tech Stack:** Go, YAML config, git worktree CLI, JSONL persistence, existing `internal/orchestration`, `internal/benchmark`, `go test`.

---

## Required Skills and Guardrails

- `@superpowers:using-git-worktrees`
- `@superpowers:test-driven-development`
- `@superpowers:verification-before-completion`
- `@superpowers:executing-plans`

Rules:
1. TDD only: failing test first, then minimal code.
2. One commit per task.
3. Update status board immediately after each task.
4. No destructive git commands.
5. Isolation runtime must not mutate the active developer worktree branch.

## Task Status Board

| Task ID | Task Name | Status | Last Update |
|---|---|---|---|
| E1 | Config Schema: Execution Isolation | TODO | - |
| E2 | Isolation Policy Resolver | TODO | - |
| E3 | Git Worktree Runtime Manager | TODO | - |
| E4 | Orchestration Integration for Isolated Candidate Runs | TODO | - |
| E5 | Critic Top-1 Deterministic Selection | TODO | - |
| E6 | Integration Branch + Same-Model Repair Retry | TODO | - |
| E7 | JSONL Schema/Store Enhancements | TODO | - |
| E8 | E2E + Docs + Full Verification | TODO | - |

---

### Task E1: Config Schema: Execution Isolation

**Files:**
- Modify: `internal/orchestration/types.go`
- Modify: `internal/orchestration/load.go`
- Modify: `internal/orchestration/config_test.go`
- Modify: `config/orchestration.yaml`

**Step 1: Write failing test**

```go
func TestLoadOrchestrationConfig_ExecutionIsolation(t *testing.T) {
    // load yaml with benchmark_mode.execution_isolation and assert fields
}
```

**Step 2: Run failing test**

Run: `go test ./internal/orchestration -run TestLoadOrchestrationConfig_ExecutionIsolation -v`
Expected: FAIL with missing struct fields.

**Step 3: Minimal implementation**

```go
type ExecutionIsolationConfig struct {
    Enabled          bool     `yaml:"enabled"`
    CommandWhitelist []string `yaml:"command_whitelist,omitempty"`
    BranchPrefix     string   `yaml:"branch_prefix,omitempty"`
    WorktreeRoot     string   `yaml:"worktree_root,omitempty"`
    RepairRetryMax   int      `yaml:"repair_retry_max,omitempty"`
}
```

Add defaults:
- `branch_prefix=bench`
- `worktree_root=.worktrees/bench`
- `repair_retry_max=1`

**Step 4: Re-run test**

Run: `go test ./internal/orchestration -run TestLoadOrchestrationConfig_ExecutionIsolation -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/orchestration/types.go internal/orchestration/load.go internal/orchestration/config_test.go config/orchestration.yaml
git commit -m "feat(config): add execution isolation schema for benchmark"
```

**Step 6: Status update**
- Mark `E1` as `DONE`.

---

### Task E2: Isolation Policy Resolver

**Files:**
- Create: `internal/benchmark/isolation_policy.go`
- Create: `internal/benchmark/isolation_policy_test.go`
- Modify: `internal/hooks/handler.go`

**Step 1: Write failing tests**

```go
func TestResolveIsolationPolicy(t *testing.T) {
    // enforce true only when benchmark enabled + code related + whitelist match + isolation enabled
}
```

**Step 2: Run failing tests**

Run: `go test ./internal/benchmark -run TestResolveIsolationPolicy -v`
Expected: FAIL with undefined resolver.

**Step 3: Minimal implementation**

```go
type IsolationPolicyInput struct {
    BenchmarkEnabled bool
    IsolationEnabled bool
    CodeRelated      bool
    Command          string
    Whitelist        []string
}

func ResolveIsolationPolicy(in IsolationPolicyInput) IsolationPolicyDecision {
    // deterministic boolean logic + reason fields
}
```

**Step 4: Re-run tests**

Run: `go test ./internal/benchmark -run TestResolveIsolationPolicy -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/benchmark/isolation_policy.go internal/benchmark/isolation_policy_test.go internal/hooks/handler.go
git commit -m "feat(benchmark): add isolation policy resolver"
```

**Step 6: Status update**
- Mark `E2` as `DONE`.

---

### Task E3: Git Worktree Runtime Manager

**Files:**
- Create: `internal/benchmark/isolation_runtime.go`
- Create: `internal/benchmark/isolation_runtime_test.go`

**Step 1: Write failing tests**

```go
func TestIsolationRuntime_CreateModelSandbox(t *testing.T) {
    // verifies branch/worktree paths and git commands
}

func TestIsolationRuntime_Cleanup(t *testing.T) {
    // verifies safe removal behavior
}
```

**Step 2: Run failing tests**

Run: `go test ./internal/benchmark -run TestIsolationRuntime -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
type GitRunner interface {
    Run(workdir string, args ...string) (string, error)
}

type ModelSandbox struct {
    Model        string
    Branch       string
    WorktreePath string
}

func (r RuntimeManager) CreateModelSandbox(runID, model, baseBranch string) (ModelSandbox, error) {
    // git worktree add -b <branch> <path> <base>
}
```

**Step 4: Re-run tests**

Run: `go test ./internal/benchmark -run TestIsolationRuntime -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/benchmark/isolation_runtime.go internal/benchmark/isolation_runtime_test.go
git commit -m "feat(benchmark): add git worktree runtime manager for model isolation"
```

**Step 6: Status update**
- Mark `E3` as `DONE`.

---

### Task E4: Orchestration Integration for Isolated Candidate Runs

**Files:**
- Modify: `internal/orchestration/executor.go`
- Modify: `internal/orchestration/result_types.go`
- Modify: `internal/orchestration/executor_test.go`

**Step 1: Write failing tests**

```go
func TestExecutor_UsesIsolatedWorktreeWhenPolicyEnforced(t *testing.T) {
    // each model candidate carries worktree/branch metadata
}
```

**Step 2: Run failing test**

Run: `go test ./internal/orchestration -run TestExecutor_UsesIsolatedWorktreeWhenPolicyEnforced -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
type DispatchContext struct {
    ExecBranch   string
    ExecWorktree string
    ExecHeadSHA  string
}
```

Wire isolation policy + runtime manager into candidate dispatch path.

**Step 4: Re-run test**

Run: `go test ./internal/orchestration -run TestExecutor_UsesIsolatedWorktreeWhenPolicyEnforced -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/orchestration/executor.go internal/orchestration/result_types.go internal/orchestration/executor_test.go
git commit -m "feat(orchestration): run benchmark candidates in isolated worktrees"
```

**Step 6: Status update**
- Mark `E4` as `DONE`.

---

### Task E5: Critic Top-1 Deterministic Selection

**Files:**
- Create: `internal/benchmark/critic_selection.go`
- Create: `internal/benchmark/critic_selection_test.go`

**Step 1: Write failing tests**

```go
func TestSelectTopCandidate(t *testing.T) {
    // weighted score desc + deterministic tie-breakers
}
```

**Step 2: Run failing tests**

Run: `go test ./internal/benchmark -run TestSelectTopCandidate -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
func SelectTopCandidate(candidates []CandidateScore) (CandidateScore, error) {
    // score desc, failures asc, duration asc, model lexical asc
}
```

**Step 4: Re-run tests**

Run: `go test ./internal/benchmark -run TestSelectTopCandidate -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/benchmark/critic_selection.go internal/benchmark/critic_selection_test.go
git commit -m "feat(benchmark): add deterministic critic top-1 selection"
```

**Step 6: Status update**
- Mark `E5` as `DONE`.

---

### Task E6: Integration Branch + Same-Model Repair Retry

**Files:**
- Create: `internal/benchmark/integration_flow.go`
- Create: `internal/benchmark/integration_flow_test.go`
- Modify: `internal/orchestration/synthesis_final_test.go`

**Step 1: Write failing tests**

```go
func TestIntegrateTopCandidate_RepairRetryOnce(t *testing.T) {
    // merge fail -> same model repair -> retry once
}
```

**Step 2: Run failing tests**

Run: `go test ./internal/benchmark -run TestIntegrateTopCandidate_RepairRetryOnce -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
type IntegrationResult struct {
    Status          string
    RepairRetryUsed int
}

func IntegrateTopCandidate(in Input) IntegrationResult {
    // top1 merge attempt, one same-model repair retry on failure
}
```

**Step 4: Re-run tests**

Run: `go test ./internal/benchmark -run TestIntegrateTopCandidate_RepairRetryOnce -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/benchmark/integration_flow.go internal/benchmark/integration_flow_test.go internal/orchestration/synthesis_final_test.go
git commit -m "feat(benchmark): add top-1 integration flow with single repair retry"
```

**Step 6: Status update**
- Mark `E6` as `DONE`.

---

### Task E7: JSONL Schema/Store Enhancements

**Files:**
- Modify: `internal/benchmark/records.go`
- Modify: `internal/benchmark/store_jsonl.go`
- Modify: `internal/benchmark/store_jsonl_test.go`
- Create: `internal/benchmark/isolation_records_test.go`

**Step 1: Write failing tests**

```go
func TestJSONLStore_IsolationStreams(t *testing.T) {
    // writes isolation_runs and extended fields
}
```

**Step 2: Run failing tests**

Run: `go test ./internal/benchmark -run TestJSONLStore_IsolationStreams -v`
Expected: FAIL.

**Step 3: Minimal implementation**

```go
type IsolationRunRecord struct {
    RunID          string   `json:"run_id"`
    Enforced       bool     `json:"enforced"`
    Command        string   `json:"command"`
    Models         []string `json:"models"`
    WorktreeRoot   string   `json:"worktree_root"`
    BranchPrefix   string   `json:"branch_prefix"`
}
```

Extend existing records with exec/integration fields from design.

**Step 4: Re-run tests**

Run: `go test ./internal/benchmark -run TestJSONLStore_IsolationStreams -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/benchmark/records.go internal/benchmark/store_jsonl.go internal/benchmark/store_jsonl_test.go internal/benchmark/isolation_records_test.go
git commit -m "feat(storage): persist isolation and integration benchmark metadata"
```

**Step 6: Status update**
- Mark `E7` as `DONE`.

---

### Task E8: E2E + Docs + Full Verification

**Files:**
- Modify: `internal/benchmark/integration_test.go`
- Modify: `README.md`
- Modify: `docs/TRIGGERS.md`
- Modify: `docs/plans/2026-03-05-benchmark-isolated-execution-design.md`

**Step 1: Write failing end-to-end test**

```go
func TestBenchmarkIsolationEnforcedForWhitelistedCommand(t *testing.T) {
    // benchmark on + code intent + whitelist => isolated runs observed
}
```

**Step 2: Run failing test**

Run: `go test ./internal/benchmark -run TestBenchmarkIsolationEnforcedForWhitelistedCommand -v`
Expected: FAIL.

**Step 3: Minimal implementation & docs wiring**
- Final e2e wiring.
- README add isolation policy section.
- TRIGGERS add whitelist/isolation behavior.

**Step 4: Re-run target tests**

Run: `go test ./internal/benchmark ./internal/orchestration ./internal/modelroute -v`
Expected: PASS.

**Step 5: Full verification**

Run: `go test ./...`
Expected: PASS.

**Step 6: Commit**

```bash
git add internal/benchmark internal/orchestration README.md docs/TRIGGERS.md docs/plans/2026-03-05-benchmark-isolated-execution-design.md
git commit -m "feat(benchmark): enforce isolated model execution and top-1 integration retry policy"
```

**Step 7: Status update**
- Mark `E8` as `DONE`.
- Confirm all rows are `DONE`.

---

## Completion Checklist

- [ ] Isolation policy is runtime-enforced with whitelist scope.
- [ ] Every eligible model run uses unique branch/worktree.
- [ ] Critic top-1 direct integration works deterministically.
- [ ] Same-model repair retry executes exactly once on merge/gate failure.
- [ ] JSONL records include isolation/integration metadata.
- [ ] All tests pass.

