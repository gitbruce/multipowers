# 脚本差异与覆盖清单（按域全量，v8.31.1 基线）

日期：2026-03-02  
基线来源：`upstream` (`https://github.com/nyldn/claude-octopus`)  
基线标签：`v8.31.1`  
baseline commit：`4b38832403ac6ff5bc716b5d0ab62be43c983dfa`

## 判定口径

- `需要覆盖`：包含业务逻辑、策略判断、状态流转、路由决策、质量门禁、调度编排，或承载关键行为验证的测试脚本。
- `可不覆盖`：纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，且不包含核心业务决策。
- 本清单要求“尽量全覆盖”，仅对“无逻辑 OS wrapper”给出 `可不覆盖`。

## 结果概览

- `upstream/v8.31.1` `.sh` 总数：`145`
- `需要覆盖`：`131`
- `可不覆盖`：`14`
- `COPY_FROM_MAIN`：`14`
- `MIGRATE_TO_GO`：`131`
- 迁移状态：`DONE=14`，`TODO=131`
- 当前 `go` 工作树现存 `.sh`：`14`

## 字段说明

- `处理策略`：`COPY_FROM_MAIN` 或 `MIGRATE_TO_GO`
- `目标位置`：目标 Go 文件路径或保留脚本路径
- `目标符号`：函数/方法名；测试脚本允许映射到统一测试套件符号
- `迁移状态`：`TODO`/`IN_PROGRESS`/`DONE`/`N/A`

## 按域全量清单

### Root 安装与部署

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `deploy.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `deploy.sh` | `N/A` | `DONE` |
| `install.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `install.sh` | `N/A` | `DONE` |

### Claude 客户端 Hooks

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `.claude/hooks/pre-commit.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `.claude/hooks/pre-commit.sh` | `N/A` | `DONE` |
| `.claude/hooks/visual-feedback.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/render/banner.go` | `Activated` | `TODO` |

### Runtime Hooks

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `hooks/agent-teams-phase-gate.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleAgentTeamsPhaseGate (planned)` | `TODO` |
| `hooks/architecture-gate.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleArchitectureGate (planned)` | `TODO` |
| `hooks/budget-gate.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleBudgetGate (planned)` | `TODO` |
| `hooks/code-quality-gate.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleCodeQualityGate (planned)` | `TODO` |
| `hooks/config-change-handler.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleConfigChangeHandler (planned)` | `TODO` |
| `hooks/context-reinforcement.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleContextReinforcement (planned)` | `TODO` |
| `hooks/frontend-gate.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleFrontendGate (planned)` | `TODO` |
| `hooks/octopus-statusline.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/status.go` | `statusData` | `TODO` |
| `hooks/perf-gate.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandlePerfGate (planned)` | `TODO` |
| `hooks/plan-mode-interceptor.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandlePlanModeInterceptor (planned)` | `TODO` |
| `hooks/provider-routing-validator.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent.go` | `RouteIntent` | `TODO` |
| `hooks/quality-gate.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleQualityGate (planned)` | `TODO` |
| `hooks/scheduler-security-gate.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleSchedulerSecurityGate (planned)` | `TODO` |
| `hooks/security-gate.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleSecurityGate (planned)` | `TODO` |
| `hooks/session-sync.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/state.go` | `ReadState/WriteState` | `TODO` |
| `hooks/sysadmin-safety-gate.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleSysadminSafetyGate (planned)` | `TODO` |
| `hooks/task-completed-transition.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleTaskCompletedTransition (planned)` | `TODO` |
| `hooks/task-completion-checkpoint.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleTaskCompletionCheckpoint (planned)` | `TODO` |
| `hooks/task-dependency-validator.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleTaskDependencyValidator (planned)` | `TODO` |
| `hooks/teammate-idle-dispatch.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleTeammateIdleDispatch (planned)` | `TODO` |
| `hooks/telemetry-webhook.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleTelemetryWebhook (planned)` | `TODO` |
| `hooks/worktree-setup.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleWorktreeSetup (planned)` | `TODO` |
| `hooks/worktree-teardown.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler.go` | `HandleWorktreeTeardown (planned)` | `TODO` |

### State 管理

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/octo-state.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/state_kv.go` | `KVGetAll/KVUpdateFromJSON` | `TODO` |
| `scripts/state-manager.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/state_kv.go` | `KVGet/KVSet/KVUpdate` | `TODO` |

### Context 管理

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/context-manager.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/context/checker.go` | `Missing/Complete` | `TODO` |

### Validation

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/validate-no-hardcoded-paths.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/validation/no_shell_runtime.go` | `ScanNoShellRuntime` | `TODO` |
| `scripts/validate-release.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/validation/gates.go` | `EnsureTargetWorkspace + ValidateReleaseArtifacts (planned)` | `TODO` |

### Provider 与 Intelligence

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/lib/common.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/lib/common.sh` | `N/A` | `DONE` |
| `scripts/lib/intelligence.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent.go` | `buildRoutingReason` | `TODO` |
| `scripts/lib/personas.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/persona.go` | `RunPersona/RenderPersonaList` | `TODO` |
| `scripts/lib/routing.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent.go` | `IsValidIntent/AllValidIntents` | `TODO` |
| `scripts/mcp-provider-detection.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/detector.go` | `DetectAll` | `TODO` |
| `scripts/provider-router.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent.go` | `RouteIntent` | `TODO` |

### Workflow 编排

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/orchestrate.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/discover.go` | `Discover/Define/Develop/Deliver facades` | `TODO` |

### Session 与 Task

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/session-manager.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/state.go` | `ReadState/WriteState` | `TODO` |
| `scripts/task-manager.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/checkbox.go` | `WriteTracking + TaskState helpers (planned)` | `TODO` |

### Metrics 与成本

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/metrics-tracker.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/metrics/tracker.go` | `TrackTokens/TrackCost/TrackDuration (planned)` | `TODO` |

### Scheduler

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/scheduler/cron.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/cron.go` | `Cron (planned)` | `TODO` |
| `scripts/scheduler/daemon.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/daemon.go` | `Daemon (planned)` | `TODO` |
| `scripts/scheduler/octopus-scheduler.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/octopus-scheduler.go` | `OctopusScheduler (planned)` | `TODO` |
| `scripts/scheduler/policy.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/policy.go` | `Policy (planned)` | `TODO` |
| `scripts/scheduler/runner.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/runner.go` | `Runner (planned)` | `TODO` |
| `scripts/scheduler/store.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/store.go` | `Store (planned)` | `TODO` |

### Agent Teams

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/agent-teams-bridge.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/teams/bridge.go` | `SyncTaskLedger/DispatchTeammate (planned)` | `TODO` |

### Permissions

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/permissions-manager.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/permissions/manager.go` | `EvaluateConsent/RequireApproval (planned)` | `TODO` |

### Extract

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/extract/core-extractor.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/extract/core.go` | `ExtractCore (planned)` | `TODO` |

### Build/Release 运维

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/build-openclaw.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/build-openclaw.sh` | `N/A` | `DONE` |
| `scripts/clean-deployment.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/clean-deployment.sh` | `N/A` | `DONE` |
| `scripts/deploy.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/deploy.sh` | `N/A` | `DONE` |
| `scripts/install-hooks.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/install-hooks.sh` | `N/A` | `DONE` |
| `scripts/release.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/release.sh` | `N/A` | `DONE` |

### Async/Tmux

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/async-tmux-features.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/async.go` | `RunAsyncTask/TrackTmuxSession (planned)` | `TODO` |

### 一次性迁移与修复

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/fix-command-frontmatter.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/fix-command-frontmatter.sh` | `N/A` | `DONE` |
| `scripts/integrate-v2.1.20-features.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/integrate-v2.1.20-features.sh` | `N/A` | `DONE` |
| `scripts/migrate-todos.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `scripts/migrate-todos.sh` | `N/A` | `DONE` |

### Legacy 测试入口

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `scripts/test-claude-octopus.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestClaudeOctopusEndToEnd (planned)` | `TODO` |
| `scripts/test-v7.13.0-features.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestV713FeatureParity (planned)` | `TODO` |

### Tests Runner

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `tests/run-all-tests.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `tests/run-all-tests.sh` | `N/A` | `DONE` |
| `tests/run-all.sh` | **可不覆盖** | 纯安装/发布/构建/清理/一次性迁移/聚合执行 wrapper，无核心业务决策 | `COPY_FROM_MAIN` | `tests/run-all.sh` | `N/A` | `DONE` |

### Tests Helpers

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `tests/helpers/generate-coverage-report.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/devx/helpers_test.go` | `TestGenerateCoverageReport (planned)` | `TODO` |
| `tests/helpers/live-test-harness.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/devx/helpers_test.go` | `TestLiveTestHarness (planned)` | `TODO` |
| `tests/helpers/mock-helpers.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/devx/helpers_test.go` | `TestMockHelpers (planned)` | `TODO` |
| `tests/helpers/test-framework.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/devx/helpers_test.go` | `TestFramework (planned)` | `TODO` |

### Tests Benchmark

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `tests/benchmark/manual-test.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/benchmark_test.go` | `TestManualTest (planned)` | `TODO` |
| `tests/benchmark/run-benchmark.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/benchmark_test.go` | `TestRunBenchmark (planned)` | `TODO` |

### Tests Integration

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `tests/integration/test-plugin-expert-review.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestPluginExpertReview (planned)` | `TODO` |
| `tests/integration/test-plugin-lifecycle.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestPluginLifecycle (planned)` | `TODO` |
| `tests/integration/test-probe-workflow.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestProbeWorkflow (planned)` | `TODO` |
| `tests/integration/test-readme-compliance.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestReadmeCompliance (planned)` | `TODO` |
| `tests/integration/test-scheduler-lifecycle.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestSchedulerLifecycle (planned)` | `TODO` |
| `tests/integration/test-value-proposition.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/integration_test.go` | `TestValueProposition (planned)` | `TODO` |

### Tests Live

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `tests/live/fix-loop.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/live_test.go` | `TestFixLoop (planned)` | `TODO` |
| `tests/live/test-prd-skill.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/live_test.go` | `TestPrdSkill (planned)` | `TODO` |
| `tests/live/test-skill-loading.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/live_test.go` | `TestSkillLoading (planned)` | `TODO` |

### Tests Smoke

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `tests/smoke/test-debate-skill.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestDebateSkill (planned)` | `TODO` |
| `tests/smoke/test-detect-providers.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestDetectProviders (planned)` | `TODO` |
| `tests/smoke/test-dry-run-all.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestDryRunAll (planned)` | `TODO` |
| `tests/smoke/test-help-commands.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestHelpCommands (planned)` | `TODO` |
| `tests/smoke/test-packaging-integrity.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestPackagingIntegrity (planned)` | `TODO` |
| `tests/smoke/test-sentinel-command.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestSentinelCommand (planned)` | `TODO` |
| `tests/smoke/test-syntax.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/smoke_test.go` | `TestSyntax (planned)` | `TODO` |

### Tests Functional/Regression

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `tests/test-command-registration.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestCommandRegistration (planned)` | `TODO` |
| `tests/test-continuation.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestContinuation (planned)` | `TODO` |
| `tests/test-debug-mode-simple.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestDebugModeSimple (planned)` | `TODO` |
| `tests/test-enforcement-pattern.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestEnforcementPattern (planned)` | `TODO` |
| `tests/test-intent-contract-skill.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestIntentContractSkill (planned)` | `TODO` |
| `tests/test-intent-questions.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestIntentQuestions (planned)` | `TODO` |
| `tests/test-lifecycle-commands.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestLifecycleCommands (planned)` | `TODO` |
| `tests/test-model-config-simple.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestModelConfigSimple (planned)` | `TODO` |
| `tests/test-multi-command.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestMultiCommand (planned)` | `TODO` |
| `tests/test-octo-state.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/tracks/state_kv_test.go` | `TestOctoState (planned)` | `TODO` |
| `tests/test-pdf-pages.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestPdfPages (planned)` | `TODO` |
| `tests/test-phases-1-2-3.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestPhases123 (planned)` | `TODO` |
| `tests/test-plan-command.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestPlanCommand (planned)` | `TODO` |
| `tests/test-ux-features-v7.16.0.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestUxFeaturesV7160 (planned)` | `TODO` |
| `tests/test-v2.1.12-integration.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV2112Integration (planned)` | `TODO` |
| `tests/test-v7.19.0-performance-fixes.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV7190PerformanceFixes (planned)` | `TODO` |
| `tests/test-v8.0.0-opus-integration.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestV800OpusIntegration (planned)` | `TODO` |
| `tests/test-v8.1.0-feature-detection.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV810FeatureDetection (planned)` | `TODO` |
| `tests/test-v8.2.0-agent-fields.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV820AgentFields (planned)` | `TODO` |
| `tests/test-v8.2.0-sonnet-workflows.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestV820SonnetWorkflows (planned)` | `TODO` |
| `tests/test-v8.24.0-perplexity-integration.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV8240PerplexityIntegration (planned)` | `TODO` |
| `tests/test-v8.25.0-dark-factory.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV8250DarkFactory (planned)` | `TODO` |
| `tests/test-v8.26.0-changelog-integration.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV8260ChangelogIntegration (planned)` | `TODO` |
| `tests/test-v8.27.0-superpowers-hardening.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/regression_test.go` | `TestV8270SuperpowersHardening (planned)` | `TODO` |
| `tests/test-v8.5.0-strategic-features.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/regression_test.go` | `TestV850StrategicFeatures (planned)` | `TODO` |
| `tests/test-version-check.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestVersionCheck (planned)` | `TODO` |
| `tests/test-version-consistency.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestVersionConsistency (planned)` | `TODO` |

### Tests Unit

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `tests/unit/test-ceremonies.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestCeremonies (planned)` | `TODO` |
| `tests/unit/test-claude-2114-features.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestClaude2114Features (planned)` | `TODO` |
| `tests/unit/test-command-frontmatter.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestCommandFrontmatter (planned)` | `TODO` |
| `tests/unit/test-crash-recovery.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestCrashRecovery (planned)` | `TODO` |
| `tests/unit/test-cron-parser.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/scheduler/cron_test.go` | `TestCronParser (planned)` | `TODO` |
| `tests/unit/test-cross-model-review.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestCrossModelReview (planned)` | `TODO` |
| `tests/unit/test-debate-routing.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent_test.go` | `TestDebateRouting (planned)` | `TODO` |
| `tests/unit/test-dependency-wbs.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestDependencyWbs (planned)` | `TODO` |
| `tests/unit/test-docs-sync.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestDocsSync (planned)` | `TODO` |
| `tests/unit/test-earned-skills.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestEarnedSkills (planned)` | `TODO` |
| `tests/unit/test-error-learning.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestErrorLearning (planned)` | `TODO` |
| `tests/unit/test-gate-thresholds.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler_test.go` | `TestGateThresholds (planned)` | `TODO` |
| `tests/unit/test-heartbeat-timeout.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestHeartbeatTimeout (planned)` | `TODO` |
| `tests/unit/test-intelligence.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/parity_test.go` | `TestIntelligence (planned)` | `TODO` |
| `tests/unit/test-knowledge-routing.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent_test.go` | `TestKnowledgeRouting (planned)` | `TODO` |
| `tests/unit/test-lockout-integration.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestLockoutIntegration (planned)` | `TODO` |
| `tests/unit/test-lockout-protocol.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestLockoutProtocol (planned)` | `TODO` |
| `tests/unit/test-mode-switching.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestModeSwitching (planned)` | `TODO` |
| `tests/unit/test-observation-importance.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestObservationImportance (planned)` | `TODO` |
| `tests/unit/test-openclaw-compat.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestOpenclawCompat (planned)` | `TODO` |
| `tests/unit/test-persona-packs.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/workflows/persona_test.go` | `TestPersonaPacks (planned)` | `TODO` |
| `tests/unit/test-provider-history.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent_test.go` | `TestProviderHistory (planned)` | `TODO` |
| `tests/unit/test-response-mode.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestResponseMode (planned)` | `TODO` |
| `tests/unit/test-routing-rules.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent_test.go` | `TestRoutingRules (planned)` | `TODO` |
| `tests/unit/test-routing.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/providers/router_intent_test.go` | `TestRouting (planned)` | `TODO` |
| `tests/unit/test-security-functions.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler_test.go` | `TestSecurityFunctions (planned)` | `TODO` |
| `tests/unit/test-sentinel.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/hooks/handler_test.go` | `TestSentinel (planned)` | `TODO` |
| `tests/unit/test-skill-frontmatter.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestSkillFrontmatter (planned)` | `TODO` |
| `tests/unit/test-structured-decisions.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestStructuredDecisions (planned)` | `TODO` |
| `tests/unit/test-tool-policy.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/cli/root_test.go` | `TestToolPolicy (planned)` | `TODO` |

### Tests Validation

| 文件 | 判定 | 说明 | 处理策略 | 目标位置 | 目标符号 | 迁移状态 |
|---|---|---|---|---|---|---|
| `tests/validate-openclaw.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/validation/gates_test.go` | `TestValidateOpenClaw (planned)` | `TODO` |
| `tests/validate-plugin-name.sh` | **需要覆盖** | 含业务/策略/流程逻辑，需迁移到 Go 原子能力或等价测试覆盖 | `MIGRATE_TO_GO` | `internal/validation/gates_test.go` | `TestValidatePluginName (planned)` | `TODO` |
