# Main (SH) -> Go (Go Runtime) Script Capability Audit

日期：2026-03-05  
比较分支：`main` vs `go` (current)

本文从**功能级别**核对 `main` 的 Shell 架构迁移到 `go` 分支后的承接情况，重点回答：
1. 关键能力是否已迁移；
2. 未 1:1 迁移时是否有更优替代；
3. 哪些能力仍存在迁移缺口。

---

## 审计口径与证据来源

- 只统计 `main` 中可执行 Shell 能力（不把 `scripts/token-extraction/*` 这类 TS/文档资产混入 Shell 核心能力统计）。
- `main` 基线统计（通过 `git ls-tree main`）：
  - `scripts/*.sh`：35
  - `hooks/*.sh`：18（另有 `hooks/pre-push` 1）
  - `tests/*.sh`：78
- `go` 现状（工作区）：
  - `scripts/*.sh`：12（以构建/发布/同步包装为主）
  - `hooks/`：7（文档 + `pre-push` + `octopus-hud.mjs`）
  - `tests/*.sh`：2（`tests/run-all*.sh`）
  - Go 测试：`internal|cmd|pkg` 下 `_test.go` 约 80 个

判定状态：
- `Migrated`：Go 有等价功能与可定位证据。
- `Replaced-Better`：非 1:1，但 Go 方案更稳健/可测。
- `Gap-Needs-Migration`：关键能力在 Go 侧无等价实现。
- `Deferred-By-Product`：能力被产品边界暂缓，不作为当前发布阻塞。

---

## 总览：能力域差异结论

| 能力域 | 结论 | 说明 |
|---|---|---|
| 编排引擎（`orchestrate.sh`） | `Migrated + Replaced-Better` | 主流程已落到 `internal/orchestration/*`，并引入更强并发/隔离/可测试结构。 |
| 状态与上下文 | `Migrated` | `internal/tracks/*` + `internal/context/*` 替代 `state-manager.sh/context-manager.sh` 核心职责。 |
| Provider 路由 | `Migrated (Core)` | `internal/providers/router_intent.go` 与 `internal/policy/*` 承接主能力。 |
| 权限与边界控制 | `Replaced-Better` | 从 Shell 字符串校验升级为 `fsboundary` + hook 入口强制校验。 |
| 生命周期 Hook | `Migrated (Runtime)` | `internal/hooks/*` 统一事件入口，替代多数 `hooks/*.sh`。 |
| 调度器（`scripts/scheduler/*.sh`） | `Deferred-By-Product` | Go 侧 `internal/scheduler` 尚未落地；现阶段并未承接 v8.15 的调度域。 |
| 提取能力（`scripts/extract/core-extractor.sh`） | `Gap-Needs-Migration` | 当前无 `internal/extract` 与 `mp extract` 命令实现。 |
| Shell 测试迁移 | `Gap-Needs-Migration` | Go 测试体系已建立，但大量原 Shell 回归脚本尚未形成等价 Go 用例。 |

---

## 重点深度拆解：`scripts/orchestrate.sh`（功能级）

> `main` 文件行数：17,797（`git show main:scripts/orchestrate.sh | wc -l`）

### 已迁移 / 更优替代能力

| 功能点（main） | main 代表函数/机制 | Go 证据 | 结论 | 说明 |
|---|---|---|---|---|
| 全流程状态机 | `case $PHASE`，discover/develop/deliver | `internal/orchestration/workflow_adapter.go` (`RunWorkflow/RunDiscover/RunDevelop/RunDeliver`) | `Migrated` | 工作流入口已模块化。 |
| 计划构建 | `jq` 组装 phase/step 计划 | `internal/orchestration/planner.go` (`BuildPlan`) + `planner_test.go` | `Replaced-Better` | 从文本拼接变成类型化 plan。 |
| 并发执行 | `parallel_execute/wait/flock` | `internal/orchestration/executor.go` (`ExecutePlan/ExecutePhase`) | `Replaced-Better` | 并行执行和结果聚合由 Go 并发模型管理。 |
| 背压控制 | Shell 锁与进程竞争 | `internal/orchestration/worktree_slots.go` (`Acquire/Release`) + `worktree_slots_test.go` | `Replaced-Better` | 明确容量上限，避免无界并发。 |
| 门禁与重试 | `evaluate_quality_branch/retry` | `internal/orchestration/gate.go` (`EvaluateGate`) + `loop.go` (`RunLoop`) | `Migrated` | 门禁与循环重试已拆分为可测组件。 |
| 异步信号监听 | 目录轮询 / inotify 变体 | `internal/orchestration/mailbox_watcher.go` + `mailbox_watcher_test.go` | `Replaced-Better` | 统一 watcher 抽象，去掉 shell 轮询噪音。 |
| 冲突检测 | `git diff` 重叠检测 | `internal/orchestration/conflict_monitor.go` (`HasOverlap`) + `conflict_monitor_test.go` | `Migrated` | 保留核心目标（变更重叠检测）。 |
| 报告合成 | markdown 拼接 | `internal/orchestration/report.go` (`GenerateReport/ToMarkdown`) + `synthesis_progressive.go` | `Migrated` | 渐进式+最终报告链路齐全。 |
| 生命周期收敛 | run 结束清理 | `internal/orchestration/lifecycle.go` (`OnAccepted/OnAborted/SweepRun`) + `lifecycle_test.go` | `Migrated` | 接受/中止/清扫具备独立契约。 |
| 实时事件与心跳 | 进度日志 + heartbeat | `internal/orchestration/events.go` + `executor.go` (`emitModelProgress`) + `executor_test.go` | `Migrated` | 可观测性从日志字符串升级为结构化事件。 |
| 运行时边界校验 | 路径校验函数 | `internal/hooks/pre_tool_use.go` + `internal/fsboundary/policy.go` | `Replaced-Better` | 在 Hook 入口做写入边界硬拦截。 |

### 未 1:1 迁移能力（需判断是否迁移）

| main 能力 | main 代表函数 | Go 现状 | 判定 | 建议 |
|---|---|---|---|---|
| 成本估算/计费提示 | `get_model_pricing` `estimate_workflow_cost` `generate_usage_report` | 未发现等价用户侧成本估算命令；现有 `internal/benchmark/*` 记录 tokens/耗时用于评测 | `Gap-Needs-Migration` | 若产品仍强调“执行前成本可见”，建议补 `mp cost estimate/report`（P1）。 |
| Provider 锁定与历史回避 | `lock_provider` `read_provider_history` | 有 fallback/policy 与 `providers.Degrade`，但无 lockout history 语义 | `Replaced-Better (Behavior Changed)` | 当前策略更简单可测；若需“失败 provider 冷却”再增补。 |
| 语义缓存/收敛裁剪 | `check_cache_semantic` `save_to_cache_semantic` `check_convergence` | 有 `synthesis_progressive` 去重，但无语义缓存层 | `Gap-Needs-Migration` | 高负载场景建议补轻量缓存（P2）。 |
| Agent checkpoint 恢复 | `save_agent_checkpoint` `load_agent_checkpoint` | 未发现等价恢复机制 | `Gap-Needs-Migration` | 长流程可靠性要求高时需补（P1/P2）。 |
| 外部 URL 安全包装 | `validate_external_url` `wrap_untrusted_content` | 未发现对应运行时实现 | `Gap-Needs-Migration` | 若仍支持外部抓取输入，应补输入净化链路（P1）。 |

> 注：以上“未发现”均基于当前树检索：`rg` 在 `internal/ cmd/ pkg/ scripts/ .claude-plugin/` 下未命中同名能力关键字。

---

## 关键非 orchestrate 脚本：功能承接结论

| Main 脚本 | Go 承接 | 结论 | 备注 |
|---|---|---|---|
| `scripts/state-manager.sh` | `internal/tracks/state.go` + `state_kv.go` + `mp state get/set/update` | `Migrated` | 原子写 + 文件锁保留。 |
| `scripts/octo-state.sh` | `internal/tracks/state_kv.go` | `Migrated` | KV 读写齐全。 |
| `scripts/context-manager.sh` | `internal/context/init_runner.go` | `Replaced-Better` | 从通用 CRUD 脚本改为 init 驱动的结构化上下文产物。 |
| `scripts/task-manager.sh` | `internal/tracks/checkbox.go` + `internal/hooks/post_tool_use.go` | `Replaced-Better` | 不再依赖会话 task-state JSON，改为 track 文档化。 |
| `scripts/provider-router.sh` + `scripts/lib/routing.sh` | `internal/providers/router_intent.go` + `internal/policy/resolve.go` + `internal/hooks/handler.go` | `Migrated` | 路由与模型契约已类型化。 |
| `scripts/mcp-provider-detection.sh` | 当前仅见 `internal/providers/detector_test.go`，无对应实现文件 | `Gap-Needs-Migration` | 功能未闭环。 |
| `scripts/permissions-manager.sh` | `internal/fsboundary/policy.go` + `internal/isolation/*` | `Replaced-Better` | 运行时围栏强于 shell 校验。 |
| `scripts/metrics-tracker.sh` | `internal/benchmark/store_jsonl.go` + `records.go` | `Replaced-Better` | 指标落盘从 session JSON 转为 append-only JSONL。 |
| `scripts/install-hooks.sh` | `scripts/install-hooks.sh`（仍在） | `Equivalent` | 依然用于安装 git pre-push hook。 |
| `scripts/extract/core-extractor.sh` | 未发现 `internal/extract` / `mp extract` | `Gap-Needs-Migration` | 属核心能力缺口。 |
| `scripts/scheduler/*.sh` | 未发现 `internal/scheduler` | `Deferred-By-Product` | 调度域待产品契约明确。 |

---

## Hooks 迁移核对（main `hooks/*.sh` -> Go runtime）

`main` 的分散 Shell hook 逻辑，已集中到：
- `internal/hooks/handler.go`（事件入口）
- `internal/hooks/pre_tool_use.go`（写入边界拦截）
- `internal/hooks/post_tool_use.go`（FAQ/track 后处理）
- `internal/hooks/stop.go`（停止门禁）

结论：
- Hook 框架层面：`Migrated/Replaced-Better`。
- 但并非所有 `main` gate 都做了同名 1:1 重建（这是架构选择，不是文档遗漏）。

---

## 测试迁移现状（防止“功能已迁移但无回归保护”）

- Go 测试已形成体系（约 80 个 `_test.go`），覆盖编排核心链路。  
- 但 `main` 中大量 `tests/*.sh` 尚未形成逐项等价替代；当前仅保留 `tests/run-all.sh`、`tests/run-all-tests.sh` 两个 shell 入口。

结论：`Gap-Needs-Migration`（测试等价性维度）。

---

## 迁移缺口决策清单（需迁移 vs 可不迁移）

### 建议迁移（避免关键能力缺失）

1. `extract` 能力闭环（`mp extract` + `internal/extract` + tests）  
2. 成本可见性命令（若产品仍要求执行前预算提示）  
3. 外部输入安全包装（若继续支持 URL 抓取分析）  
4. 长流程 checkpoint 恢复（可靠性要求高时）  

### 可暂不迁移（Go 已有更优方案或产品暂缓）

1. Scheduler 全域（待产品契约后再迁）  
2. Provider lockout/history 的 shell 实现（现由 policy/fallback 简化替代）  
3. 大量 shell 回归脚本 1:1 复刻（优先迁移高风险路径，逐步转 Go tests）

---

## 本次更新后的结论

`go` 分支在**核心编排能力**上已完成迁移，并在并发、隔离、类型安全和测试可维护性上优于 `main` Shell 架构。  
当前风险主要不在主编排，而在**域能力缺口**（`extract`、部分安全包装、成本可见性、checkpoint）与**测试等价性缺口**。

