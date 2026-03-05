# 从 main 基线版本（`v8.23.1`）到 upstream `v8.34.0` 的变更

- main 基线版本：`v8.23.1`
- 目标 upstream 版本：`v8.34.0`
- 数据来源：`upstream` 分支 `CHANGELOG.md`（`refs/tmp/upstream-main`）

## 按版本变更摘要

### `v8.25.0`（2026-02-25）
- 增加 Dark Factory 模式（`/octo:factory`），支持从规范到交付的自动化流水线。
- 增加 holdout 场景测试、加权满意度评分、失败重试与修复建议闭环。

### `v8.26.0`（2026-02-26）
- 增加 Claude Code v2.1.46-v2.1.59 功能旗标与 changelog 集成。
- 增加 worktree create/remove 生命周期 hooks，并扩展 hook 事件类型。
- 增加更多默认配置项与 doctor 检查覆盖。
- 扩展隔离能力与 native auto-memory 行为。

### `v8.27.0`（2026-02-26）
- 强化上下文压缩后（context compaction）的规则保持能力。
- 重写部分 skill 描述，降低“浅触发”导致的执行偏差。
- 引入 `<HARD-GATE>` 强约束标记与 expensive skills 的 human-only 调用控制。
- 增加 staged review 流程与 Plan 模式拦截 hook。

### `v8.30.0`（2026-02-28）
- 增加迭代流程中的 continuation/resume 能力。

### `v8.31.0`（2026-02-28）
- 增强多模型协作智能。

### `v8.31.1`（2026-03-01）
- 增加 `/octo:batch` 别名。
- 强化并行流程默认质量策略。

### `v8.32.0`（2026-03-04）
- 增加 marketing / finance / legal 人设。
- 增加 IDE 集成能力。

### `v8.33.0`（2026-03-04）
- 增加 UI/UX 设计工作流与 BM25 设计智能能力。

### `v8.34.0`（2026-03-04）
- 增加 recurrence detection（重复问题识别）。
- 增加 issue categorization（问题分类）。
- 增加 JSONL 决策日志。
- 集成 CodeRabbit 相关流程增强。

## 备注

- upstream changelog 中未单列 `v8.24.0`、`v8.28.0`、`v8.29.0` 条目。
- 因此从 `v8.23.1` 到 `v8.34.0` 的主要可追踪变化集中在 `v8.25.0`、`v8.26.0`、`v8.27.0`、`v8.30.0`、`v8.31.x`、`v8.32.0`、`v8.33.0`、`v8.34.0`。
