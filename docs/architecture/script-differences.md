# main vs go 脚本差异与覆盖清单（按域全量）

日期：2026-03-02  
比较分支：`main` vs `go`  
基线提交：`main=f6a815a326ec`，`go=8835e073834f`

## 判定口径

- `需要覆盖`：包含业务逻辑、策略判断、状态流转、路由决策、质量门禁、调度编排，或承载关键行为验证的测试脚本。
- `可不覆盖`：纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，且不包含核心业务决策。

语义状态定义：
- `equivalent`：语义与行为已等价（含从 main 只读保留脚本）。
- `partial`：已在 go 侧有目标文件，但行为/断言尚需补齐。
- `missing`：go 侧目标文件尚未落地或不可定位。
- `intentional-diff`：go 分支新增脚本（main 无对应）。

迁移决策说明（产品约束优先）：
- 本文档不要求把 `main` 所有脚本逐文件迁移为 Go 实现。
- 本文档现有 `处理策略` 列即 `decision` 字段，使用：`MIGRATE_TO_GO`、`COPY_FROM_MAIN`、`KEEP_IN_GO`。
- 对仍为 `missing` 且非当前产品范围必需能力，可采用 `EXCLUDE_WITH_REASON` 或 `DEFER_WITH_CONDITION`（通过补充整改与专项索引显式记录）。
- `COPY_FROM_MAIN` 均采用 `read-only-from-main` 同步策略。

## Evidence Legend

- `E0`：doc-only plan（仅文档规划）
- `E1`：symbol exists（目标符号已存在）
- `E2`：test exists（已有对应测试）
- `E3`：verified output recorded（有验证输出记录）

当前大表默认证据级别口径：
- `equivalent` + `DONE`：至少 `E1`（推荐补到 `E2/E3`）。
- `partial` + `TODO`：至少 `E0`。
- `missing`：`E0` 且必须补充 `decision` 与触发条件/原因。

## 数据基线

| 指标 | 数量 |
|---|---:|
| `main` 脚本总数 (`*.sh`) | 135 |
| `go` 脚本总数 (`*.sh`) | 16 |
| shared | 14 |
| main-only | 121 |
| go-only | 2 |

覆盖策略统计（全量 137 行 = main 135 + go-only 2）：
- `需要覆盖=121`
- `可不覆盖=16`（包含 go-only 2）

语义状态统计（全量 137 行）：
- `equivalent=14`
- `partial=71`
- `missing=50`
- `intentional-diff=2`

## 最新代码复检（go=8835e073834f）

- 已基于最新 `go` 分支重新回放主脚本映射：`tmp/recheck/scripts-refresh.tsv`（main 135 行全覆盖）。
- 覆盖校验结果：`main` 脚本 `missing=0`、`extra=0`（即文档与主分支脚本树一一对应）。
- `MIGRATE_TO_GO` 目标文件存在性：`target_exists=71`，`target_missing=50`，与当前状态统计 `partial=71`、`missing=50` 一致。

## 历史文档校正

本次对旧版 `script-differences.md` 做了基线回放，移除 10 条不在当前 `main` 脚本树中的历史行：

- `hooks/context-reinforcement.sh`
- `hooks/plan-mode-interceptor.sh`
- `hooks/telemetry-webhook.sh`
- `hooks/worktree-setup.sh`
- `hooks/worktree-teardown.sh`
- `tests/test-continuation.sh`
- `tests/test-v8.24.0-perplexity-integration.sh`
- `tests/test-v8.25.0-dark-factory.sh`
- `tests/test-v8.26.0-changelog-integration.sh`
- `tests/test-v8.27.0-superpowers-hardening.sh`

## 按域全量清单（文件级一对一补充）

### Root 安装与部署

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `deploy.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `deploy.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |
| `install.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `install.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |

### Claude 客户端 Hooks

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `.claude/hooks/pre-commit.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `.claude/hooks/pre-commit.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |
| `.claude/hooks/visual-feedback.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/render/banner.go` | `Activated` | `TODO` | 对齐 internal/render/banner.go 的行为与断言；补充回归测试后升级为 equivalent |

### Runtime Hooks

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `hooks/agent-teams-phase-gate.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleAgentTeamsPhaseGate (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/architecture-gate.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleArchitectureGate (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/budget-gate.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleBudgetGate (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/code-quality-gate.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleCodeQualityGate (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/config-change-handler.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleConfigChangeHandler (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/frontend-gate.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleFrontendGate (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/octopus-statusline.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/status.go` | `statusData` | `TODO` | 对齐 internal/cli/status.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/perf-gate.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandlePerfGate (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/provider-routing-validator.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent.go` | `RouteIntent` | `TODO` | 对齐 internal/providers/router_intent.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/quality-gate.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleQualityGate (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/scheduler-security-gate.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleSchedulerSecurityGate (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/security-gate.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleSecurityGate (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/session-sync.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/state.go` | `ReadState/WriteState` | `TODO` | 对齐 internal/tracks/state.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/sysadmin-safety-gate.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleSysadminSafetyGate (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/task-completed-transition.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleTaskCompletedTransition (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/task-completion-checkpoint.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleTaskCompletionCheckpoint (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/task-dependency-validator.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleTaskDependencyValidator (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |
| `hooks/teammate-idle-dispatch.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleTeammateIdleDispatch (planned)` | `TODO` | 对齐 internal/hooks/handler.go 的行为与断言；补充回归测试后升级为 equivalent |

### State 管理

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/octo-state.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/state_kv.go` | `KVGetAll/KVUpdateFromJSON` | `TODO` | 对齐 internal/tracks/state_kv.go 的行为与断言；补充回归测试后升级为 equivalent |
| `scripts/state-manager.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/state_kv.go` | `KVGet/KVSet/KVUpdate` | `TODO` | 对齐 internal/tracks/state_kv.go 的行为与断言；补充回归测试后升级为 equivalent |

### Context 管理

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/context-manager.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/context/checker.go` | `Missing/Complete` | `TODO` | 对齐 internal/context/checker.go 的行为与断言；补充回归测试后升级为 equivalent |

### Validation

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/validate-no-hardcoded-paths.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/validation/no_shell_runtime.go` | `ScanNoShellRuntime` | `TODO` | 对齐 internal/validation/no_shell_runtime.go 的行为与断言；补充回归测试后升级为 equivalent |
| `scripts/validate-release.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/validation/gates.go` | `EnsureTargetWorkspace + ValidateReleaseArtifacts (planned)` | `TODO` | 对齐 internal/validation/gates.go 的行为与断言；补充回归测试后升级为 equivalent |

### Provider 与 Intelligence

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/lib/common.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/lib/common.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |
| `scripts/lib/intelligence.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent.go` | `buildRoutingReason` | `TODO` | 对齐 internal/providers/router_intent.go 的行为与断言；补充回归测试后升级为 equivalent |
| `scripts/lib/personas.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/persona.go` | `RunPersona/RenderPersonaList` | `TODO` | 对齐 internal/workflows/persona.go 的行为与断言；补充回归测试后升级为 equivalent |
| `scripts/lib/routing.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent.go` | `IsValidIntent/AllValidIntents` | `TODO` | 对齐 internal/providers/router_intent.go 的行为与断言；补充回归测试后升级为 equivalent |
| `scripts/mcp-provider-detection.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/detector.go` | `DetectAll` | `TODO` | 对齐 internal/providers/detector.go 的行为与断言；补充回归测试后升级为 equivalent |
| `scripts/provider-router.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent.go` | `RouteIntent` | `TODO` | 对齐 internal/providers/router_intent.go 的行为与断言；补充回归测试后升级为 equivalent |

### Workflow 编排

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/orchestrate.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/discover.go` | `Discover/Define/Develop/Deliver facades` | `TODO` | 对齐 internal/workflows/discover.go 的行为与断言；补充回归测试后升级为 equivalent |

### Session 与 Task

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/session-manager.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/state.go` | `ReadState/WriteState` | `TODO` | 对齐 internal/tracks/state.go 的行为与断言；补充回归测试后升级为 equivalent |
| `scripts/task-manager.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/checkbox.go` | `WriteTracking + TaskState helpers (planned)` | `TODO` | 对齐 internal/tracks/checkbox.go 的行为与断言；补充回归测试后升级为 equivalent |

### Metrics 与成本

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/metrics-tracker.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/metrics/tracker.go` | `TrackTokens/TrackCost/TrackDuration (planned)` | `TODO` | 在 internal/metrics/tracker.go 落地实现（按域拆分）并补充测试覆盖 |

### Scheduler

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/scheduler/cron.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/cron.go` | `Cron (planned)` | `TODO` | 在 internal/scheduler/cron.go 落地实现（按域拆分）并补充测试覆盖 |
| `scripts/scheduler/daemon.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/daemon.go` | `Daemon (planned)` | `TODO` | 在 internal/scheduler/daemon.go 落地实现（按域拆分）并补充测试覆盖 |
| `scripts/scheduler/octopus-scheduler.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/octopus-scheduler.go` | `OctopusScheduler (planned)` | `TODO` | 在 internal/scheduler/octopus-scheduler.go 落地实现（按域拆分）并补充测试覆盖 |
| `scripts/scheduler/policy.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/policy.go` | `Policy (planned)` | `TODO` | 在 internal/scheduler/policy.go 落地实现（按域拆分）并补充测试覆盖 |
| `scripts/scheduler/runner.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/runner.go` | `Runner (planned)` | `TODO` | 在 internal/scheduler/runner.go 落地实现（按域拆分）并补充测试覆盖 |
| `scripts/scheduler/store.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/store.go` | `Store (planned)` | `TODO` | 在 internal/scheduler/store.go 落地实现（按域拆分）并补充测试覆盖 |

### Agent Teams

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/agent-teams-bridge.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/teams/bridge.go` | `SyncTaskLedger/DispatchTeammate (planned)` | `TODO` | 在 internal/teams/bridge.go 落地实现（按域拆分）并补充测试覆盖 |

### Permissions

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/permissions-manager.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/permissions/manager.go` | `EvaluateConsent/RequireApproval (planned)` | `TODO` | 在 internal/permissions/manager.go 落地实现（按域拆分）并补充测试覆盖 |

### Extract

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/extract/core-extractor.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/extract/core.go` | `ExtractCore (planned)` | `TODO` | 在 internal/extract/core.go 落地实现（按域拆分）并补充测试覆盖 |

### Build/Release 运维

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/build-openclaw.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/build-openclaw.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |
| `scripts/clean-deployment.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/clean-deployment.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |
| `scripts/deploy.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/deploy.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |
| `scripts/install-hooks.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/install-hooks.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |
| `scripts/release.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/release.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |

### Async/Tmux

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/async-tmux-features.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/async.go` | `RunAsyncTask/TrackTmuxSession (planned)` | `TODO` | 在 internal/workflows/async.go 落地实现（按域拆分）并补充测试覆盖 |

### 一次性迁移与修复

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/fix-command-frontmatter.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/fix-command-frontmatter.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |
| `scripts/integrate-v2.1.20-features.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/integrate-v2.1.20-features.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |
| `scripts/migrate-todos.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/migrate-todos.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |

### Legacy 测试入口

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/test-claude-octopus.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestClaudeOctopusEndToEnd (planned)` | `TODO` | 在 internal/workflows/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `scripts/test-v7.13.0-features.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestV713FeatureParity (planned)` | `TODO` | 在 internal/workflows/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |

### Tests Runner

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `tests/run-all-tests.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `tests/run-all-tests.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |
| `tests/run-all.sh` | **可不覆盖** | `equivalent` | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `tests/run-all.sh` | `N/A` | `DONE` | 保持从 main 只读同步拷贝；无需 Go 迁移 |

### Tests Helpers

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `tests/helpers/generate-coverage-report.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/devx/helpers_test.go` | `TestGenerateCoverageReport (planned)` | `TODO` | 在 internal/devx/helpers_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/helpers/live-test-harness.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/devx/helpers_test.go` | `TestLiveTestHarness (planned)` | `TODO` | 在 internal/devx/helpers_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/helpers/mock-helpers.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/devx/helpers_test.go` | `TestMockHelpers (planned)` | `TODO` | 在 internal/devx/helpers_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/helpers/test-framework.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/devx/helpers_test.go` | `TestFramework (planned)` | `TODO` | 在 internal/devx/helpers_test.go 落地实现（按域拆分）并补充测试覆盖 |

### Tests Benchmark

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `tests/benchmark/manual-test.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/benchmark_test.go` | `TestManualTest (planned)` | `TODO` | 在 internal/workflows/benchmark_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/benchmark/run-benchmark.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/benchmark_test.go` | `TestRunBenchmark (planned)` | `TODO` | 在 internal/workflows/benchmark_test.go 落地实现（按域拆分）并补充测试覆盖 |

### Tests Integration

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `tests/integration/test-plugin-expert-review.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestPluginExpertReview (planned)` | `TODO` | 在 internal/workflows/integration_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/integration/test-plugin-lifecycle.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestPluginLifecycle (planned)` | `TODO` | 在 internal/workflows/integration_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/integration/test-probe-workflow.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestProbeWorkflow (planned)` | `TODO` | 在 internal/workflows/integration_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/integration/test-readme-compliance.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestReadmeCompliance (planned)` | `TODO` | 在 internal/workflows/integration_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/integration/test-scheduler-lifecycle.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestSchedulerLifecycle (planned)` | `TODO` | 在 internal/workflows/integration_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/integration/test-value-proposition.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestValueProposition (planned)` | `TODO` | 在 internal/workflows/integration_test.go 落地实现（按域拆分）并补充测试覆盖 |

### Tests Live

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `tests/live/fix-loop.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/live_test.go` | `TestFixLoop (planned)` | `TODO` | 在 internal/workflows/live_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/live/test-prd-skill.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/live_test.go` | `TestPrdSkill (planned)` | `TODO` | 在 internal/workflows/live_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/live/test-skill-loading.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/live_test.go` | `TestSkillLoading (planned)` | `TODO` | 在 internal/workflows/live_test.go 落地实现（按域拆分）并补充测试覆盖 |

### Tests Smoke

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `tests/smoke/test-debate-skill.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestDebateSkill (planned)` | `TODO` | 在 internal/cli/smoke_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/smoke/test-detect-providers.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestDetectProviders (planned)` | `TODO` | 在 internal/cli/smoke_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/smoke/test-dry-run-all.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestDryRunAll (planned)` | `TODO` | 在 internal/cli/smoke_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/smoke/test-help-commands.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestHelpCommands (planned)` | `TODO` | 在 internal/cli/smoke_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/smoke/test-packaging-integrity.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestPackagingIntegrity (planned)` | `TODO` | 在 internal/cli/smoke_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/smoke/test-sentinel-command.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestSentinelCommand (planned)` | `TODO` | 在 internal/cli/smoke_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/smoke/test-syntax.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestSyntax (planned)` | `TODO` | 在 internal/cli/smoke_test.go 落地实现（按域拆分）并补充测试覆盖 |

### Tests Functional/Regression

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `tests/test-command-registration.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestCommandRegistration (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/test-debug-mode-simple.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestDebugModeSimple (planned)` | `TODO` | 在 internal/cli/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-enforcement-pattern.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestEnforcementPattern (planned)` | `TODO` | 在 internal/cli/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-intent-contract-skill.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestIntentContractSkill (planned)` | `TODO` | 在 internal/cli/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-intent-questions.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestIntentQuestions (planned)` | `TODO` | 在 internal/cli/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-lifecycle-commands.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestLifecycleCommands (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/test-model-config-simple.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestModelConfigSimple (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/test-multi-command.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestMultiCommand (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/test-octo-state.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/state_kv_test.go` | `TestOctoState (planned)` | `TODO` | 对齐 internal/tracks/state_kv_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/test-pdf-pages.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestPdfPages (planned)` | `TODO` | 在 internal/cli/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-phases-1-2-3.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestPhases123 (planned)` | `TODO` | 在 internal/workflows/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-plan-command.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestPlanCommand (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/test-ux-features-v7.16.0.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestUxFeaturesV7160 (planned)` | `TODO` | 在 internal/cli/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-v2.1.12-integration.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV2112Integration (planned)` | `TODO` | 在 internal/cli/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-v7.19.0-performance-fixes.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV7190PerformanceFixes (planned)` | `TODO` | 在 internal/cli/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-v8.0.0-opus-integration.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestV800OpusIntegration (planned)` | `TODO` | 在 internal/workflows/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-v8.1.0-feature-detection.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV810FeatureDetection (planned)` | `TODO` | 在 internal/cli/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-v8.2.0-agent-fields.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV820AgentFields (planned)` | `TODO` | 在 internal/cli/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-v8.2.0-sonnet-workflows.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestV820SonnetWorkflows (planned)` | `TODO` | 在 internal/workflows/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-v8.5.0-strategic-features.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestV850StrategicFeatures (planned)` | `TODO` | 在 internal/workflows/regression_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/test-version-check.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestVersionCheck (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/test-version-consistency.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestVersionConsistency (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/validate-openclaw.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/validation/gates_test.go` | `TestValidateOpenClaw (planned)` | `TODO` | 对齐 internal/validation/gates_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/validate-plugin-name.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/validation/gates_test.go` | `TestValidatePluginName (planned)` | `TODO` | 对齐 internal/validation/gates_test.go 的行为与断言；补充回归测试后升级为 equivalent |

### Tests Unit

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `tests/unit/test-ceremonies.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestCeremonies (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-claude-2114-features.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestClaude2114Features (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-command-frontmatter.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestCommandFrontmatter (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-crash-recovery.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestCrashRecovery (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-cron-parser.sh` | **需要覆盖** | `missing` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/cron_test.go` | `TestCronParser (planned)` | `TODO` | 在 internal/scheduler/cron_test.go 落地实现（按域拆分）并补充测试覆盖 |
| `tests/unit/test-cross-model-review.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestCrossModelReview (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-debate-routing.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent_test.go` | `TestDebateRouting (planned)` | `TODO` | 对齐 internal/providers/router_intent_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-dependency-wbs.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestDependencyWbs (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-docs-sync.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestDocsSync (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-earned-skills.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestEarnedSkills (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-error-learning.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestErrorLearning (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-gate-thresholds.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler_test.go` | `TestGateThresholds (planned)` | `TODO` | 对齐 internal/hooks/handler_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-heartbeat-timeout.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestHeartbeatTimeout (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-intelligence.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/parity_test.go` | `TestIntelligence (planned)` | `TODO` | 对齐 internal/providers/parity_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-knowledge-routing.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent_test.go` | `TestKnowledgeRouting (planned)` | `TODO` | 对齐 internal/providers/router_intent_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-lockout-integration.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestLockoutIntegration (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-lockout-protocol.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestLockoutProtocol (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-mode-switching.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestModeSwitching (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-observation-importance.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestObservationImportance (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-openclaw-compat.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestOpenclawCompat (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-persona-packs.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/persona_test.go` | `TestPersonaPacks (planned)` | `TODO` | 对齐 internal/workflows/persona_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-provider-history.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent_test.go` | `TestProviderHistory (planned)` | `TODO` | 对齐 internal/providers/router_intent_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-response-mode.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestResponseMode (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-routing-rules.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent_test.go` | `TestRoutingRules (planned)` | `TODO` | 对齐 internal/providers/router_intent_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-routing.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent_test.go` | `TestRouting (planned)` | `TODO` | 对齐 internal/providers/router_intent_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-security-functions.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler_test.go` | `TestSecurityFunctions (planned)` | `TODO` | 对齐 internal/hooks/handler_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-sentinel.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler_test.go` | `TestSentinel (planned)` | `TODO` | 对齐 internal/hooks/handler_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-skill-frontmatter.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestSkillFrontmatter (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-structured-decisions.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestStructuredDecisions (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |
| `tests/unit/test-tool-policy.sh` | **需要覆盖** | `partial` | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestToolPolicy (planned)` | `TODO` | 对齐 internal/cli/root_test.go 的行为与断言；补充回归测试后升级为 equivalent |

### go 分支新增脚本（main 无对应）

| 文件 | 判定 | 语义状态 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 | 补充整改 |
|---|---|---|---|---|---|---|---|---|
| `scripts/auto-push.sh` | **可不覆盖** | `intentional-diff` | go 分支 CI/协作辅助 wrapper，main 无此文件 | `KEEP_IN_GO` | `scripts/auto-push.sh` | `N/A` | `N/A` | 保留为 go 分支附加能力，无需回迁 main |
| `scripts/build.sh` | **可不覆盖** | `intentional-diff` | go 分支构建 wrapper，main 无此文件 | `KEEP_IN_GO` | `scripts/build.sh` | `N/A` | `N/A` | 保留为 go 分支附加能力，无需回迁 main |


## Hook Lifecycle Alignment Index

| lifecycle_event | legacy_script_source | go_target_symbol | response_contract_fields | test_reference | evidence_level | decision | note |
|---|---|---|---|---|---|---|---|
| `SessionStart` | `hooks/session-sync.sh` | `internal/hooks/session_start.go` + `internal/tracks/state.go` | `status`, `action`, `message`, `data.session_id` | `internal/hooks/stop_test.go:TestHookStopAndSubagentStop` (covers lifecycle) | `E1` | `MIGRATE_TO_GO` | 初始化与状态同步属于运行时控制面；test planned for session lifecycle |
| `UserPromptSubmit` | `scripts/context-manager.sh` + `hooks/architecture-gate.sh` | `internal/context/checker.go` + `internal/hooks/handler.go` | `status`, `action`, `error_code`, `message`, `remediation` | `internal/hooks/handler_test.go:TestUserPromptSubmitBlocksMissingContext` | `E2` | `MIGRATE_TO_GO` | 提示词提交前需做上下文契约与门禁检查 |
| `PreToolUse` | `hooks/security-gate.sh` + `hooks/provider-routing-validator.sh` | `internal/hooks/pre_tool_use.go` + `internal/providers/router_intent.go` | `status`, `action`, `error_code`, `message`, `data.tool_name`, `remediation` | `internal/hooks/handler_test.go` (coverage via integration) | `E1` | `MIGRATE_TO_GO` | 高成本工具调用前的策略校验；dedicated test planned |
| `PostToolUse` | `hooks/quality-gate.sh` + `hooks/task-completion-checkpoint.sh` | `internal/hooks/post_tool_use.go` + `internal/hooks/handler.go` | `status`, `action`, `error_code`, `message`, `data.result_quality` | `internal/hooks/handler_test.go:TestPostToolUseWritesFaqAndTrack` | `E2` | `MIGRATE_TO_GO` | 工具调用后质量门禁与任务检查点 |
| `Stop` | `hooks/task-completed-transition.sh` | `internal/hooks/stop.go` | `status`, `action`, `message`, `data.pending_tasks`, `data.completed_tasks` | `internal/hooks/stop_test.go:TestHookStopAndSubagentStop` | `E2` | `MIGRATE_TO_GO` | 会话停止阶段的收口与状态落盘 |
| `SubagentStop` | `hooks/teammate-idle-dispatch.sh` + `hooks/task-dependency-validator.sh` | `internal/hooks/stop.go` + `internal/hooks/handler.go` | `status`, `action`, `message`, `data.subagent_id`, `data.idle_reason` | `internal/hooks/stop_test.go:TestHookStopAndSubagentStop` | `E2` | `MIGRATE_TO_GO` | 子代理生命周期与依赖转移收口 |

**Response Contract Field Definitions:**
- `status`: `"ok"` | `"blocked"` | `"error"` - Hook execution outcome
- `action`: `"continue"` | `"block"` | `"remediate"` - Recommended next step
- `error_code`: `string` (optional) - Machine-readable error identifier
- `message`: `string` - Human-readable explanation
- `data`: `object` - Event-specific payload (session_id, tool_name, etc.)
- `remediation`: `string` (optional) - Suggested fix when blocked

## Missing Decision Classification Matrix

The following matrix provides default decision classification for all `missing` scripts grouped by domain pattern.

| source_pattern | count | default_decision | decision_reason | closure_path | exception_rules |
|----------------|-------|------------------|-----------------|--------------|-----------------|
| `scripts/scheduler/*.sh` | 6 | `DEFER_WITH_CONDITION` | Scheduler domain contract undefined in `.multipowers/product.md` | `internal/scheduler/*_test.go` after domain contract defined | None |
| `scripts/extract/*.sh` | 1 | `MIGRATE_TO_GO` | Core extraction workflow required | `internal/extract/core.go` + `internal/extract/core_test.go` | None |
| `tests/smoke/*.sh` | 7 | `MIGRATE_TO_GO` | Critical validation tests for CLI surface | `internal/cli/smoke_test.go` | None |
| `tests/live/*.sh` | 3 | `DEFER_WITH_CONDITION` | Live tests require external service availability | `internal/workflows/live_test.go` | Enable when test environment ready |
| `tests/benchmark/*.sh` | 2 | `MIGRATE_TO_GO` | Performance regression guard | `internal/workflows/benchmark_test.go` | None |
| `tests/integration/*.sh` | 6 | `MIGRATE_TO_GO` | Integration tests for plugin lifecycle | `internal/workflows/integration_test.go` | None |
| `tests/helpers/*.sh` | 4 | `MIGRATE_TO_GO` | Test infrastructure utilities | `internal/devx/helpers_test.go` | None |
| `scripts/metrics-tracker.sh` | 1 | `MIGRATE_TO_GO` | Cost tracking is core observability | `internal/metrics/tracker.go` + `internal/metrics/tracker_test.go` | None |
| `scripts/permissions-manager.sh` | 1 | `MIGRATE_TO_GO` | Consent management is governance capability | `internal/permissions/manager.go` + `internal/permissions/manager_test.go` | None |
| `scripts/agent-teams-bridge.sh` | 1 | `MIGRATE_TO_GO` | Team coordination required for subagent workflows | `internal/teams/bridge.go` + `internal/teams/bridge_test.go` | None |
| `scripts/async-tmux-features.sh` | 1 | `MIGRATE_TO_GO` | Async execution is core workflow feature | `internal/workflows/async.go` + `internal/workflows/async_test.go` | None |
| `tests/test-*.sh` (regression) | 21 | `MIGRATE_TO_GO` | Feature regression tests ensure parity | `internal/cli/regression_test.go` + `internal/workflows/regression_test.go` | None |
| `tests/unit/test-cron-parser.sh` | 1 | `MIGRATE_TO_GO` | Scheduler dependency requires cron parsing | `internal/scheduler/cron_test.go` | None |

**Summary:**
- `MIGRATE_TO_GO`: 48 scripts (scheduler-independent tests, core features)
- `DEFER_WITH_CONDITION`: 10 scripts (scheduler domain, live tests)
- `EXCLUDE_WITH_REASON`: 0 scripts (all missing items have migration path)

## 整改优先级

1. `P0`：先对 `missing=50` 完成逐项 `decision` 分类（`MIGRATE_TO_GO` / `EXCLUDE_WITH_REASON` / `DEFER_WITH_CONDITION`），不再默认要求”全部迁移”。
2. `P0`：对 `internal/hooks/handler.go`、`internal/providers/router_intent.go`、`internal/workflows/*_test.go` 的 `partial` 项补行为断言与回归测试。
3. `P1`：将所有 `MIGRATE_TO_GO + TODO` 项转为可验证的 `*_test.go` 对应测试，并记录通过证据。
4. `P1`：维持 `COPY_FROM_MAIN` 项只读同步策略，避免 main/go wrapper 漂移。

## Parity 结论

- main 脚本清单已实现 100% 文件级覆盖登记（`missing_from_doc=0`）。
- go 侧对 main 脚本能力的语义迁移仍为 `partial parity`（`partial+missing=121`）。
- 当前可认定结论：脚本域未达到与 main 完全等价；需按上述优先级持续收敛。
