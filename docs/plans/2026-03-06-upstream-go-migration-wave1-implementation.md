# Upstream Governance Migration Wave 1 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Deliver a Go-native migration wave that includes `/mp:doctor` (16 checks), decision JSONL + recurrence intelligence, issue categorization normalization, plan-mode/worktree lifecycle hard hooks, and CodeRabbit governance in the `/mp` + `.multipowers` architecture.

**Architecture:** Implement a single doctor engine in Go under `internal/doctor`, expose it via `mp-devx --action doctor`, and make `mp doctor` a pure proxy for code reuse. Extend runtime governance by adding typed decision logging (`.multipowers/decisions/decisions.jsonl`), recurrence analysis, plan/worktree hook hard gates, and category-normalized issue tracking guidance aligned to Go-first responsibility boundaries.

**Tech Stack:** Go (`context`, `time`, `os/exec`, JSON), existing `internal/cli`, `cmd/mp-devx`, `internal/hooks`, `.claude-plugin` command/skill assets, YAML/JSON config, `go test`.

---

### Task 1: Create Doctor Domain Contract And Registry (16 checks)

**Files:**
- Create: `internal/doctor/types.go`
- Create: `internal/doctor/registry.go`
- Create: `internal/doctor/registry_test.go`

**Step 1: Write the failing test**

```go
func TestRegistry_Has16ChecksSortedByID(t *testing.T) {
    checks := DefaultRegistry()
    if len(checks) != 16 {
        t.Fatalf("len=%d want 16", len(checks))
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/doctor -run TestRegistry_Has16ChecksSortedByID -v`
Expected: FAIL (package/files missing)

**Step 3: Write minimal implementation**

- Define `CheckSpec`, `CheckResult`, `RunReport`, `Status`.
- Register 16 `check_id`s.
- Mark `fail_capable=true` only for:
  - `auth`, `command-boundary`, `config`, `hooks`, `multipowers-boundary`, `no-shell-runtime`, `policy-freshness`, `providers`, `skills`

**Step 4: Run test to verify it passes**

Run: `go test ./internal/doctor -run TestRegistry_Has16ChecksSortedByID -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/doctor/types.go internal/doctor/registry.go internal/doctor/registry_test.go
git commit -m "feat(doctor): add 16-check registry and domain contract"
```

### Task 2: Implement Concurrent Runner With Timeout, Single Check, List, Save

**Files:**
- Create: `internal/doctor/runner.go`
- Create: `internal/doctor/output.go`
- Create: `internal/doctor/runner_test.go`

**Step 1: Write failing tests**

```go
func TestRunner_InvalidCheckIDReturnsError(t *testing.T) {}
func TestRunner_DefaultTimeout_AllVsSingle(t *testing.T) {}
func TestRunner_TimeoutMarksWarnTimedOut(t *testing.T) {}
func TestRunner_ListReturnsIDPurposeFailCapable(t *testing.T) {}
func TestRunner_SaveWritesExpectedPath(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/doctor -run TestRunner_ -v`
Expected: FAIL

**Step 3: Write minimal implementation**

- Concurrency: one goroutine per selected check.
- Deterministic output: always sort by `check_id`.
- Timeout policy:
  - all checks default `30s`
  - `--check-id` default `45s`
  - explicit `--timeout` overrides both.
- Timeout behavior: only current check canceled; result `warn + timed_out=true`.
- `--list`: columns `check_id/purpose/fail_capable`.
- `--save`: write JSON to `.multipowers/doctor/reports/doctor[-<check_id>]-YYYYMMDD-HHMMSS.json`.
- Exit code semantics at integration layer: non-zero only if any `fail`.

**Step 4: Run tests to verify pass**

Run: `go test ./internal/doctor -run TestRunner_ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/doctor/runner.go internal/doctor/output.go internal/doctor/runner_test.go
git commit -m "feat(doctor): add concurrent runner with timeout/list/save contracts"
```

### Task 3: Implement Upstream 9 Check Handlers In Go

**Files:**
- Create: `internal/doctor/checks_upstream.go`
- Create: `internal/doctor/checks_upstream_test.go`
- Reuse: `internal/validation/no_shell_runtime.go`, `config/*`, `.claude-plugin/*`

**Step 1: Write failing tests**

```go
func TestChecks_ConfigFailsWhenCoderabbitMissing(t *testing.T) {}
func TestChecks_AuthFailsWhenNoProviderAuth(t *testing.T) {}
func TestChecks_HooksValidatesCommandTargets(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/doctor -run TestChecks_ -v`
Expected: FAIL

**Step 3: Write minimal implementation**

Implement check IDs:
- `providers`, `auth`, `config`, `state`, `hooks`, `skills`, `conflicts`, `agents`, `recurrence`

Key rules:
- `config`: `.coderabbit.yaml` missing => `fail`.
- `recurrence`: `48h>=3 || 7d>=5` warn, plus source concentration `>=3` warn.
- `hooks/skills`: verify referenced files exist.

**Step 4: Run tests to verify pass**

Run: `go test ./internal/doctor -run TestChecks_ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/doctor/checks_upstream.go internal/doctor/checks_upstream_test.go
git commit -m "feat(doctor): implement upstream-compatible 9 checks"
```

### Task 4: Implement Local 7 Check Handlers

**Files:**
- Create: `internal/doctor/checks_local.go`
- Create: `internal/doctor/checks_local_test.go`

**Step 1: Write failing tests**

```go
func TestLocalChecks_CommandBoundaryDetectsDrift(t *testing.T) {}
func TestLocalChecks_NoShellRuntimeUsesValidator(t *testing.T) {}
func TestLocalChecks_PolicyFreshnessDetectsMissingCompiledPolicy(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/doctor -run TestLocalChecks_ -v`
Expected: FAIL

**Step 3: Write minimal implementation**

Implement local check IDs:
- `command-boundary`
- `no-shell-runtime`
- `multipowers-boundary`
- `namespace-drift`
- `policy-freshness`
- `checkpoint-health`
- `runtime-status-consistency`

**Step 4: Run tests to verify pass**

Run: `go test ./internal/doctor -run TestLocalChecks_ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/doctor/checks_local.go internal/doctor/checks_local_test.go
git commit -m "feat(doctor): add 7 local governance checks"
```

### Task 5: Integrate mp-devx Doctor Action

**Files:**
- Modify: `cmd/mp-devx/main.go`
- Modify: `cmd/mp-devx/main_test.go`

**Step 1: Write failing test**

```go
func TestRun_ActionDoctor(t *testing.T) {
    rc := run([]string{"-action", "doctor", "-list"}, &out, io.Discard)
    if rc != 0 { t.Fatalf("expected rc=0") }
}
```

**Step 2: Run test to verify failure**

Run: `go test ./cmd/mp-devx -run TestRun_ActionDoctor -v`
Expected: FAIL (`unknown action`)

**Step 3: Write minimal implementation**

Add action `doctor` with flags:
- `-check-id`
- `-timeout`
- `-list`
- `-save`
- `-verbose`
- `-json`

Exit code: non-zero only when report contains `fail`.

**Step 4: Run test to verify pass**

Run: `go test ./cmd/mp-devx -run TestRun_ActionDoctor -v`
Expected: PASS

**Step 5: Commit**

```bash
git add cmd/mp-devx/main.go cmd/mp-devx/main_test.go
git commit -m "feat(mp-devx): add doctor action backed by internal doctor engine"
```

### Task 6: Integrate mp Doctor As Pure Proxy To mp-devx

**Files:**
- Modify: `internal/cli/root.go`
- Modify: `internal/cli/root_test.go`

**Step 1: Write failing tests**

```go
func TestDoctor_UnknownCheckIDReturnsError(t *testing.T) {}
func TestDoctor_ProxyMissingMPDevxReturnsError(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/cli -run TestDoctor_ -v`
Expected: FAIL

**Step 3: Write minimal implementation**

- Add `doctor` command in `mp`.
- Proxy execution to `mp-devx --action doctor ...`.
- stdout/stderr/exit code passthrough.
- If `mp-devx` missing: return `error` + remediation (no auto-build).

**Step 4: Run tests to verify pass**

Run: `go test ./internal/cli -run TestDoctor_ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/root.go internal/cli/root_test.go
git commit -m "feat(mp): add doctor proxy to mp-devx action doctor"
```

### Task 7: Add Decision JSONL Domain And Recurrence Input Pipeline

**Files:**
- Create: `internal/decisions/store.go`
- Create: `internal/decisions/store_test.go`
- Modify: `internal/hooks/pre_tool_use.go`
- Modify: `internal/hooks/stop.go`

**Step 1: Write failing tests**

```go
func TestDecisionStore_AppendsJSONL(t *testing.T) {}
func TestDecisionStore_DefaultPathUnderMultipowers(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/decisions -v`
Expected: FAIL

**Step 3: Write minimal implementation**

- Append-only JSONL at `.multipowers/decisions/decisions.jsonl`.
- Fields: `id,type,timestamp,source,summary,scope,confidence,importance`.
- Hook integration:
  - boundary block writes `quality-gate` decision.
  - stop block writes `quality-gate` decision.

**Step 4: Run tests to verify pass**

Run: `go test ./internal/decisions ./internal/hooks -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/decisions/store.go internal/decisions/store_test.go internal/hooks/pre_tool_use.go internal/hooks/stop.go
git commit -m "feat(decisions): add append-only decision JSONL and hook failure logging"
```

### Task 8: Migrate Issue Categorization To Normalized Enum

**Files:**
- Create: `internal/issues/category.go`
- Create: `internal/issues/category_test.go`
- Modify: `internal/faq/classify.go`
- Modify: `.claude-plugin/.claude/skills/skill-issues.md`

**Step 1: Write failing tests**

```go
func TestNormalizeCategory_AllowsKnownValues(t *testing.T) {}
func TestNormalizeCategory_RejectsUnknown(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/issues -v`
Expected: FAIL

**Step 3: Write minimal implementation**

- Enum categories:
  - `logic-error`, `integration`, `quality-gate`, `security`, `performance`, `ux`, `architecture`
- Implement normalize/validate helpers.
- Update FAQ classification to map tool/runtime signals into normalized categories.
- Update `skill-issues.md` from `.octo` to `.multipowers` and enforce category selection list.

**Step 4: Run tests to verify pass**

Run: `go test ./internal/issues ./internal/faq -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/issues/category.go internal/issues/category_test.go internal/faq/classify.go .claude-plugin/.claude/skills/skill-issues.md
git commit -m "feat(issues): add normalized category enum and migrate issue skill paths"
```

### Task 9: Add Plan-Mode Interceptor And Worktree Lifecycle Hooks

**Files:**
- Modify: `.claude-plugin/hooks.json`
- Modify: `internal/hooks/handler.go`
- Modify: `internal/hooks/handler_test.go`
- Create: `internal/hooks/worktree_events.go`
- Create: `internal/hooks/worktree_events_test.go`

**Step 1: Write failing tests**

```go
func TestEnterPlanMode_BlocksWithoutPlanIntent(t *testing.T) {}
func TestWorktreeEvents_ArePersisted(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/hooks -run 'Test(EnterPlanMode|WorktreeEvents)' -v`
Expected: FAIL

**Step 3: Write minimal implementation**

- Add hook events in plugin hooks config:
  - `EnterPlanMode`
  - `WorktreeCreate`
  - `WorktreeRemove`
- `EnterPlanMode`: enforce `/mp:plan`-aligned remediation path.
- Worktree events: append event records to `.multipowers/temp/worktree-events.jsonl`.

**Step 4: Run tests to verify pass**

Run: `go test ./internal/hooks -run 'Test(EnterPlanMode|WorktreeEvents)' -v`
Expected: PASS

**Step 5: Commit**

```bash
git add .claude-plugin/hooks.json internal/hooks/handler.go internal/hooks/handler_test.go internal/hooks/worktree_events.go internal/hooks/worktree_events_test.go
git commit -m "feat(hooks): add plan-mode interceptor and worktree lifecycle event persistence"
```

### Task 10: Add /mp Doctor Command Asset And Plugin Registration

**Files:**
- Create: `.claude-plugin/.claude/commands/doctor.md`
- Modify: `.claude-plugin/plugin.json`

**Step 1: Write failing test**

```go
func TestPluginJSON_ContainsDoctorCommand(t *testing.T) {}
```

**Step 2: Run test to verify failure**

Run: `go test ./internal/doctor -run TestPluginJSON_ContainsDoctorCommand -v`
Expected: FAIL

**Step 3: Write minimal implementation**

- Add `/mp:doctor` command doc that resolves `mp` binary and executes:
  - `mp doctor [--list|--check-id ...|--timeout ...|--json|--verbose|--save]`
- Register command path in plugin manifest.

**Step 4: Run test to verify pass**

Run: `go test ./internal/doctor -run TestPluginJSON_ContainsDoctorCommand -v`
Expected: PASS

**Step 5: Commit**

```bash
git add .claude-plugin/.claude/commands/doctor.md .claude-plugin/plugin.json
git commit -m "feat(plugin): register /mp:doctor command asset"
```

### Task 11: Add CodeRabbit Governance Config (Go-first)

**Files:**
- Create: `.coderabbit.yaml`

**Step 1: Write failing test**

```go
func TestConfigCheck_FailsWhenCodeRabbitMissing(t *testing.T) {}
```

**Step 2: Run test to verify failure**

Run: `go test ./internal/doctor -run TestConfigCheck_FailsWhenCodeRabbitMissing -v`
Expected: FAIL

**Step 3: Write minimal implementation**

- Add `.coderabbit.yaml` with Go-first path instructions:
  - `internal/**`, `cmd/**`, `.claude-plugin/**`, `config/**`, `docs/plans/**`
- Keep strict review gates enabled.

**Step 4: Run test to verify pass**

Run: `go test ./internal/doctor -run TestConfigCheck_FailsWhenCodeRabbitMissing -v`
Expected: PASS

**Step 5: Commit**

```bash
git add .coderabbit.yaml
git commit -m "chore(governance): add coderabbit config for go-first architecture"
```

### Task 12: Full Verification And Evidence

**Files:**
- Modify: `docs/CLI-REFERENCE.md`
- Modify: `docs/COMMAND-REFERENCE.md`
- Create: `docs/plans/evidence/doctor/2026-03-06-doctor-verification.md`

**Step 1: Run targeted tests**

Run:
- `go test ./internal/doctor ./internal/issues ./internal/decisions ./internal/hooks -v`
- `go test ./cmd/mp-devx ./internal/cli -v`

Expected: PASS

**Step 2: Run smoke CLI checks**

Run:
- `go run ./cmd/mp-devx --action doctor --list`
- `go run ./cmd/mp-devx --action doctor --check-id config --json`
- `go run ./cmd/mp doctor --check-id config --timeout 10s --json`

Expected:
- `--list` shows `check_id/purpose/fail_capable`
- JSON includes `timed_out/elapsed_ms/timeout_ms/fail_capable`
- `mp doctor` output/exit mirrors `mp-devx`

**Step 3: Document evidence**

- Save command output snippets and test summary to evidence markdown.
- Update CLI/command docs for new doctor command and flags.

**Step 4: Final commit**

```bash
git add docs/CLI-REFERENCE.md docs/COMMAND-REFERENCE.md docs/plans/evidence/doctor/2026-03-06-doctor-verification.md
git commit -m "docs: add doctor command reference and verification evidence"
```

