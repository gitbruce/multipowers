# main vs go 脚本差异与覆盖清单（最新状态）

日期：2026-03-03  
比较分支：`main` vs `go`  
基线提交：`main=f6a815a326ec`，`go=cf865fa764fe`

## 判定口径

- `需要覆盖`：包含业务逻辑、策略判断、状态流转、路由决策、质量门禁、调度编排，或关键行为验证。
- `可不覆盖`：安装/发布/构建/清理/一次性迁移/聚合 wrapper，且不含核心业务决策。

状态定义：`equivalent` / `partial` / `missing` / `intentional-diff`

决策（decision）取值：`MIGRATE_TO_GO`、`KEEP_IN_GO`、`DEFER_WITH_CONDITION`、`EXCLUDE_WITH_REASON`

## Evidence Legend

- `E0`：doc-only plan（仅文档规划）
- `E1`：symbol exists（目标符号已存在）
- `E2`：test exists（已有对应测试）
- `E3`：verified output recorded（有验证输出记录）

## 最新统计

| 指标 | 数量 |
|---|---:|
| `main` 脚本总数 (`*.sh`) | 135 |
| `go` 脚本总数 (`*.sh`) | 15 |
| shared | 14 |
| main-only | 121 |
| go-only | 1 |

覆盖策略统计（全量 136 行 = main 135 + go-only 1）：
- `需要覆盖=121`
- `可不覆盖=15`（包含 go-only 1）

语义状态统计（全量 136 行）：
- `equivalent=14`
- `partial=71`
- `missing=50`
- `intentional-diff=1`

## 最新 go-only 脚本

- `scripts/auto-push.sh`
- `scripts/build.sh`

## `tmp/compare.md` 补充漏项（2026-03-03）

- `.claude/hooks/visual-feedback.sh`：`main` 中存在、`go` 中已移除；当前应标记为 `partial`，决策为 `MIGRATE_TO_GO`。
  - 语义承接：`.claude/skills/skill-visual-feedback.md` 与 `internal/render/banner.go`。

## 核心域迁移快照（最新）

| 域 | source pattern | decision | target (go) | 状态 |
|---|---|---|---|---|
| State/Session/Task | `scripts/{state-manager,octo-state,session-manager,task-manager}.sh` | `MIGRATE_TO_GO` | `internal/tracks/*` | `partial` |
| Context/Validation | `scripts/{context-manager,validate-*.sh}` | `MIGRATE_TO_GO` | `internal/context/*`, `internal/validation/*` | `partial` |
| Provider/Routing | `scripts/{provider-router,mcp-provider-detection}.sh`, `scripts/lib/{intelligence,personas,routing}.sh` | `MIGRATE_TO_GO` | `internal/providers/*`, `internal/workflows/persona.go` | `partial` |
| Orchestrate | `scripts/orchestrate.sh` | `MIGRATE_TO_GO` | `internal/workflows/*` | `partial` |
| Hooks | `hooks/*.sh` | `MIGRATE_TO_GO` | `internal/hooks/*` | `partial` |
| Scheduler | `scripts/scheduler/*.sh` | `DEFER_WITH_CONDITION` | `internal/scheduler/*` | `missing` |
| Metrics/Permissions/Agent/Extract | `scripts/{metrics-tracker,permissions-manager,agent-teams-bridge}.sh`, `scripts/extract/core-extractor.sh` | `MIGRATE_TO_GO` | `internal/{metrics,permissions,teams,extract}/*` | `missing/partial` |
| 测试脚本集 | `tests/**/*.sh`（除 run-all wrappers） | `MIGRATE_TO_GO` | `internal/*_test.go` | `missing/partial` |

## Hook Lifecycle Alignment Index（最新）

| lifecycle_event | source | go target | response_contract_fields | evidence | decision |
|---|---|---|---|---|---|
| `SessionStart` | `hooks/session-sync.sh` | `internal/hooks/*` + `internal/tracks/state.go` | `status`, `action`, `message`, `data` | `E1` | `MIGRATE_TO_GO` |
| `UserPromptSubmit` | `scripts/context-manager.sh` + `hooks/architecture-gate.sh` | `internal/context/checker.go` + `internal/hooks/handler.go` | `status`, `action`, `error_code`, `message`, `remediation` | `E2` | `MIGRATE_TO_GO` |
| `PreToolUse` | `hooks/security-gate.sh` + `hooks/provider-routing-validator.sh` | `internal/hooks/pre_tool_use.go` + `internal/providers/router_intent.go` | `status`, `action`, `error_code`, `message`, `data` | `E1` | `MIGRATE_TO_GO` |
| `PostToolUse` | `hooks/quality-gate.sh` + `hooks/task-completion-checkpoint.sh` | `internal/hooks/post_tool_use.go` + `internal/hooks/handler.go` | `status`, `action`, `error_code`, `message`, `data` | `E2` | `MIGRATE_TO_GO` |
| `Stop` | `hooks/task-completed-transition.sh` | `internal/hooks/stop.go` | `status`, `action`, `message`, `data` | `E2` | `MIGRATE_TO_GO` |
| `SubagentStop` | `hooks/teammate-idle-dispatch.sh` + `hooks/task-dependency-validator.sh` | `internal/hooks/stop.go` + `internal/hooks/handler.go` | `status`, `action`, `message`, `data` | `E2` | `MIGRATE_TO_GO` |
