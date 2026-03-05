# MP/MP-DEVX Command Boundary Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 统一命令入口边界：运行时能力只保留在 `mp`，工程/运维能力统一收敛到 `mp-devx`，并提供可回归验证与迁移文档。

**Architecture:** 采用“两可执行 + 共享 internal 包”的模式。`cmd/mp` 只承接用户运行时路径（workflow/orchestrate/hook/state/context）；`cmd/mp-devx` 承接构建、校验、性能、兼容性检查。通过 CLI 合约测试锁定边界，避免后续命令漂移。

**Tech Stack:** Go (`flag`, `testing`), 现有 CLI 框架 (`internal/cli`, `cmd/mp-devx`), 文档 (`docs/architecture`, `docs/plans`).

---

## Command Ownership Checklist (Target)

### Keep in `mp` (Runtime)
- `init`
- `context guard`
- `state get|set|update`
- `plan|discover|research|define|develop|deliver|review|embrace|debate|persona`
- `orchestrate select-agent`
- `loop`
- `hook`
- `status`
- `route`
- `extract`
- `checkpoint save|get|delete`
- `config get|show-model-routing`

### Move to `mp-devx` (Ops/Devx)
- `mp test run` -> `mp-devx --action suite --suite <unit|integration|all>`
- `mp coverage check` -> `mp-devx --action coverage` (新增)
- `mp validate --type no-shell` / `mp validate --strict-no-shell` -> `mp-devx --action validate-runtime` (新增，复用 validation 扫描)
- `mp cost report` -> `mp-devx --action cost-report` (新增，消费 `.multipowers/metrics`)

### Keep in both (明确职责差异)
- `mp cost estimate` 保留在 runtime（执行前预算提示）
- `mp-devx --action bench` 保留在 devx（性能门禁）

---

## Migration Script Checklist (mp -> mp-devx)

1. 在 `mp-devx` 增加 `coverage`/`validate-runtime`/`cost-report` 三个 action。
2. 在 `internal/cli/root.go` 将 `test`、`coverage`、`validate(no-shell)` 标记为 deprecated 并输出迁移提示（第一阶段不立刻删除）。
3. 第二阶段删除 `mp` 中对应分支，只保留错误提示和迁移文案。
4. 更新 `CLAUDE.md` / `docs/architecture/*` 的命令示例。
5. 在 CI 增加“命令归属漂移检查”（grep case 分支 + 快照文件比对）。

---

### Task 1: 固化命令归属文档与快照

**Files:**
- Create: `docs/architecture/command-ownership.md`
- Create: `docs/plans/evidence/command-boundary/commands-snapshot-2026-03-05.md`
- Modify: `CLAUDE.md`

**Step 1: Write the failing test**

```go
func TestCommandOwnershipSnapshotExists(t *testing.T) {
    if _, err := os.Stat("docs/plans/evidence/command-boundary/commands-snapshot-2026-03-05.md"); err != nil {
        t.Fatalf("snapshot missing: %v", err)
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/validation -run TestCommandOwnershipSnapshotExists -v`  
Expected: FAIL with "snapshot missing"

**Step 3: Write minimal implementation**

- 增加命令归属文档和 snapshot，记录 `mp`/`mp-devx` 当前命令矩阵与目标矩阵。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/validation -run TestCommandOwnershipSnapshotExists -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add docs/architecture/command-ownership.md docs/plans/evidence/command-boundary/commands-snapshot-2026-03-05.md CLAUDE.md
git commit -m "docs: add mp/mp-devx command ownership baseline"
```

### Task 2: 为 mp 命令迁移添加合约测试（先红）

**Files:**
- Modify: `internal/cli/root_test.go`
- Test: `internal/cli/root_test.go`

**Step 1: Write the failing test**

```go
func TestMPShowsDeprecationForTestCoverageAndNoShellValidate(t *testing.T) {
    code := Run([]string{"test", "run", "--json"})
    if code == 0 {
        t.Fatalf("expected non-zero with deprecation guidance")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run TestMPShowsDeprecationForTestCoverageAndNoShellValidate -v`  
Expected: FAIL because current command still executes.

**Step 3: Write minimal implementation**

- 在 `internal/cli/root.go` 的 `test`/`coverage`/`validate(no-shell)` 分支先返回迁移提示（`use mp-devx --action ...`）。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -run TestMPShowsDeprecationForTestCoverageAndNoShellValidate -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/root.go internal/cli/root_test.go
git commit -m "refactor: deprecate ops commands in mp with migration hints"
```

### Task 3: 在 mp-devx 增加运维 action（先红后绿）

**Files:**
- Modify: `cmd/mp-devx/main.go`
- Modify: `cmd/mp-devx/main_test.go`
- Modify: `internal/devx/runner.go`
- Modify: `internal/devx/runner_test.go`

**Step 1: Write the failing test**

```go
func TestRunSupportsCoverageValidateRuntimeAndCostReport(t *testing.T) {
    code := run([]string{"--action", "coverage"}, io.Discard, io.Discard)
    if code != 0 {
        t.Fatalf("expected coverage action supported")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cmd/mp-devx -run TestRunSupportsCoverageValidateRuntimeAndCostReport -v`  
Expected: FAIL with "unknown action"

**Step 3: Write minimal implementation**

- `main.go` 新增 action:
  - `coverage` -> 调用 `workflows.CoverageCheck`
  - `validate-runtime` -> 调用 `validation.ScanNoShellRuntime`
  - `cost-report` -> 调用 `cost.BuildReport`
- 抽象到 `internal/devx/Runner`，避免 `cmd` 层堆业务逻辑。

**Step 4: Run test to verify it passes**

Run: `go test ./cmd/mp-devx ./internal/devx -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add cmd/mp-devx/main.go cmd/mp-devx/main_test.go internal/devx/runner.go internal/devx/runner_test.go
git commit -m "feat: add coverage/validate-runtime/cost-report actions to mp-devx"
```

### Task 4: 清理 mp 中运维命令入口（第二阶段）

**Files:**
- Modify: `internal/cli/root.go`
- Modify: `internal/cli/root_test.go`

**Step 1: Write the failing test**

```go
func TestMPOpsCommandsAreBlockedWithStableErrorCode(t *testing.T) {
    code := Run([]string{"coverage", "check", "--json"})
    if code == 0 {
        t.Fatalf("expected blocked/error response")
    }
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli -run TestMPOpsCommandsAreBlockedWithStableErrorCode -v`  
Expected: FAIL if command still executes.

**Step 3: Write minimal implementation**

- 删除 `mp` 中 `test`/`coverage` 的执行逻辑，保留稳定错误码 + 迁移提示。

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli -v`  
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/root.go internal/cli/root_test.go
git commit -m "refactor: remove ops execution paths from mp"
```

### Task 5: 未挂入口 Go 方法审计自动化

**Files:**
- Create: `scripts/audit/deadcode-mp.sh`
- Create: `tmp/unused_go.md` (generated)
- Create: `docs/plans/evidence/command-boundary/deadcode-baseline-2026-03-05.txt`

**Step 1: Write the failing test**

```bash
#!/usr/bin/env bash
set -euo pipefail
test -f tmp/unused_go.md
```

**Step 2: Run test to verify it fails**

Run: `bash scripts/audit/check-unused-go-report.sh`  
Expected: FAIL when report missing.

**Step 3: Write minimal implementation**

- `deadcode-mp.sh` 生成 deadcode 原始输出和 markdown 报告；
- 报告固定输出到 `tmp/unused_go.md`。

**Step 4: Run test to verify it passes**

Run: `bash scripts/audit/deadcode-mp.sh && bash scripts/audit/check-unused-go-report.sh`  
Expected: PASS

**Step 5: Commit**

```bash
git add scripts/audit/deadcode-mp.sh scripts/audit/check-unused-go-report.sh docs/plans/evidence/command-boundary/deadcode-baseline-2026-03-05.txt tmp/unused_go.md
git commit -m "chore: add deadcode audit for mp/mp-devx entry reachability"
```

### Task 6: 全量验证与发布说明

**Files:**
- Modify: `docs/architecture/script-differences.md`
- Modify: `docs/architecture/commands_skills_difference.md`
- Modify: `SAFEGUARDS.md`

**Step 1: Write the failing test**

```bash
#!/usr/bin/env bash
set -euo pipefail
rg -n "mp test run|mp coverage check" docs/architecture | grep -v "deprecated"
```

**Step 2: Run test to verify it fails**

Run: `bash scripts/audit/check-doc-command-drift.sh`  
Expected: FAIL if docs still claim old ownership.

**Step 3: Write minimal implementation**

- 全面更新文档中的命令入口和示例。

**Step 4: Run test to verify it passes**

Run: `go test ./... && bash scripts/audit/check-doc-command-drift.sh`  
Expected: PASS

**Step 5: Commit**

```bash
git add docs/architecture/script-differences.md docs/architecture/commands_skills_difference.md SAFEGUARDS.md scripts/audit/check-doc-command-drift.sh
git commit -m "docs: align command ownership after mp/mp-devx split"
```

---

## Verification Gate (must pass before merge)

- `go test ./...`
- `go test ./cmd/mp-devx ./internal/devx ./internal/cli -v`
- `bash scripts/audit/deadcode-mp.sh`
- `bash scripts/audit/check-doc-command-drift.sh`

## Rollout Strategy

1. Release N: `mp` 输出迁移提示，不直接删除执行路径。  
2. Release N+1: 删除 `mp` 运维执行路径，保留稳定错误码与迁移提示。  
3. Release N+2: 清理兼容分支与旧文档。
