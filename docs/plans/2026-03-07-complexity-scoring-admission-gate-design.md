# Complexity Scoring Admission Gate 设计

## 0. 2026-03-07 基线评估

基于当前 `go` 分支代码，复杂度门禁已经有一个可工作的雏形，但离目标闭环还有明显差距：

- `internal/validation/gates.go` 目前只实现了“高复杂度 => 需要 active track + worktree”的判定。
- `internal/tracks/scoring.go` 的入口评分仍以 `prompt` 启发式为主，尚未形成“准入评分 + 计划内细化评分”的两阶段模型。
- track 产物规范已经迁移到 canonical 结构（`intent.md`、`design.md`、`implementation-plan.md`、`metadata.json`、`index.md`），但旧草稿 track 仍保留 `spec.md` / `plan.md` 形态，不能直接作为新的 planning completeness 依据。
- `internal/tracks/metadata.go` 已经有 `interrupted_context` 字段，但运行时尚未把“高复杂度被拦截的原始任务意图”稳定写入并恢复。

因此，本设计的核心不是重写复杂度评分器本身，而是将 **复杂度评分、planning 文档完备性、worktree 执行约束、上下文恢复** 收敛为一个统一的 admission gate。

---

## 1. 目标与非目标

### 1.1 目标

本次设计要实现以下闭环：

1. 对所有 spec-driven `mp` 命令入口执行 complexity admission 检查：
   - `discover`
   - `define`
   - `develop`
   - `deliver`
   - `embrace`
   - `review`
   - `research`
   - `debate`
2. 当任务被判定为高复杂度时，不仅要求在 worktree 中执行，还要求当前 active track 已具备完整 canonical planning 文档。
3. 若高复杂度但 planning 文档不完整，必须阻断当前命令并强制跳转到 `/mp:plan`。
4. 被阻断时保存 `interrupted_context`；在 `/mp:plan` 完成后，允许用户显式继续原任务，但不自动偷偷执行原命令。
5. 低复杂度命令继续走现有轻量路径，不被强制要求 `/mp:plan`。

### 1.2 非目标

本次不做以下事情：

- 不把该门禁直接扩展到通用 `Write/Edit/MultiEdit` hook。
- 不自动 merge worktree 回主分支。
- 不恢复旧的 `spec.md` / `plan.md` 双轨兼容；planning 完整性只认 canonical artifacts。
- 不在本次内引入完整 LLM 复杂度评估；入口阶段先基于 `command + prompt + active track metadata` 做准入评分。

---

## 2. 方案选择

本设计采用“**中心化 Spec Admission Gate**”方案。

### 2.1 推荐方案

在现有 `RunSpecPipeline` 前半段增加一个统一 admission gate，按以下顺序做判定：

1. 全局上下文是否完整（已有 `/mp:init` 语义，优先级最高）。
2. 当前命令是否属于受控 spec-driven 命令。
3. admission complexity 是否达到高复杂度阈值。
4. 若为高复杂度，active track 是否存在、canonical planning artifacts 是否完整、metadata 是否已有明确执行结论。
5. 若 planning 完整，再检查当前是否处于 linked worktree。

### 2.2 不选 hook-first 的原因

虽然将 gate 放到 hook 层看起来更“硬”，但当前 hook 层缺少完整任务语义：

- `PreToolUse` 更多看到的是工具调用，而不是“这是不是一个高复杂度 spec-driven 开发任务”。
- 直接在 hook 层拦写文件，容易误伤低复杂度的小修改。
- 当前 `pipeline` 已经是 spec-driven 命令的统一入口，职责更适合承接 admission gate。

---

## 3. 准入模型

### 3.1 两阶段复杂度模型

复杂度判断分为两个阶段：

#### 阶段 A：Admission Complexity

发生在命令刚进入 runtime 时，使用以下输入：

- `command`
- `prompt`
- active track metadata（若存在）

目标不是精确估算全部工作量，而是回答：

> 这个命令现在是否必须先进入显式 planning 流程？

#### 阶段 B：Plan Refinement

发生在 `/mp:plan` 内部或其后续补齐阶段，补充更细粒度信息：

- 预估变更文件数
- 触达模块数
- 任务组数量
- 预计工时
- 外部集成风险
- 迁移/安全关键路径

其产出会写回 track metadata，作为后续执行期 worktree / group enforcement 的依据。

### 3.2 高复杂度的定义

高复杂度并不只等价于“需要 worktree”，而是表示：

- 需要显式 planning artifacts
- 需要 metadata 中存在复杂度与执行模式结论
- 需要在 worktree 中执行真正的修改阶段

也就是说，高复杂度任务的门禁语义是：

> 先计划，再执行；执行时还必须在隔离环境中进行。

---

## 4. Planning Completeness 语义

### 4.1 Canonical artifacts

高复杂度任务必须具备以下完整集合：

- `.multipowers/tracks/<id>/intent.md`
- `.multipowers/tracks/<id>/design.md`
- `.multipowers/tracks/<id>/implementation-plan.md`
- `.multipowers/tracks/<id>/metadata.json`
- `.multipowers/tracks/<id>/index.md`

### 4.2 旧文档不再构成合规 planning

旧草稿格式中的：

- `spec.md`
- `plan.md`

即使存在，也不能被认定为“planning 已完成”。

原因是当前 runtime、模板和 registry 都已经围绕 canonical artifact 体系工作；若继续把旧格式也视作合规，会导致 runtime contract 不一致。

### 4.3 Metadata 要求

对于高复杂度任务，仅有文档文件存在还不够，metadata 还必须至少满足：

- `complexity_score > 0`
- `worktree_required = true`

这样 admission gate 才能区分两种 block：

1. 缺 planning 文档或执行结论，需要 `/mp:plan`
2. planning 已完成，但当前不在 worktree 中

---

## 5. 运行时状态流转

### 5.1 状态定义

建议将高复杂度任务在 admission 层抽象为四个状态：

- `unplanned`
- `planning_required`
- `planned`
- `execution_ready`

含义如下：

- `unplanned`：存在任务意图，但还没有可执行的 planning 载体。
- `planning_required`：高复杂度且缺 active track / 缺 canonical artifacts / 缺 metadata 执行结论。
- `planned`：文档和 metadata 已齐，但尚未进入合规 worktree。
- `execution_ready`：planning 完整，且当前处于合规 worktree 中。

### 5.2 典型流转

#### 场景 A：高复杂度直接执行 `develop`

1. 用户执行 `mp develop --prompt ...`
2. admission complexity 判定为高复杂度
3. 若无 active track，则分配/创建一个 track 作为承载者
4. 保存 `interrupted_context`
5. 返回 `blocked`，明确 remediation 为 `/mp:plan`

#### 场景 B：用户完成 `/mp:plan`

1. `/mp:plan` 复用当前 track
2. 补齐 canonical artifacts
3. 更新 metadata 中的复杂度与 execution mode 决策
4. track 进入 `planned`

#### 场景 C：再次执行 `develop`

1. admission gate 发现 planning 完整
2. 若未在 worktree 中，则返回“必须切换到 worktree”
3. 若已在 worktree 中，则进入 `execution_ready`，放行执行

---

## 6. 组件拆分

### 6.1 `internal/tracks`

新增或补齐以下能力：

- canonical artifact 完整性检查
- 对 active track planning 状态的机器判断
- 复用 `InterruptedContext` 作为被阻断任务的持久化容器

这里不建议把 planning 完整性放进 `context.Missing()`，因为它不是项目初始化上下文，而是 track 级别状态。

### 6.2 `internal/tracks/scoring.go`

保留现有评分器，但扩展其输出表达能力，使其不再只表达 `WorktreeRequired`，还可以表达：

- 是否需要先 planning
- 评分来源（admission / refined）
- 评分 rationale

### 6.3 `internal/validation`

将当前简单的 `EnsureComplexityGate` 升级为 richer admission result：

- 是否允许执行
- block 原因
- 推荐动作
- track id
- complexity score
- requires planning / requires worktree
- missing artifacts
- resume 是否可用

### 6.4 `internal/app/pipeline.go`

仍然作为 spec-driven 命令统一入口，但只负责消费 admission result，不再自己推断复杂度 gate 的具体语义。

### 6.5 `internal/cli/root.go`

负责：

- 在 block 时写入 `interrupted_context`
- 在没有 active track 的情况下分配 track 并绑定被拦截上下文
- `/mp:plan` 复用当前 track
- 将 structured data 返回给上层调用者

---

## 7. Response Contract

### 7.1 高复杂度但缺 planning

返回：

- `status = blocked`
- `action = ask_user_questions`
- `message = High complexity detected. Planning artifacts are required before proceeding.`
- `remediation = Run /mp:plan to complete design and implementation planning for this track.`

`data` 至少包含：

- `track_id`
- `complexity_score`
- `requires_planning = true`
- `requires_worktree = true`
- `missing_artifacts`
- `resume_command`
- `resume_prompt`
- `interrupted_context_saved = true`

### 7.2 Planning 已完整但不在 worktree

返回：

- `status = blocked`
- `requires_planning = false`
- `requires_worktree = true`
- `missing_artifacts = []`

这样上层能够明确知道：

- 不是没计划
- 而是执行环境还不符合要求

### 7.3 低复杂度

低复杂度命令直接放行，不附加 planning block 语义。

---

## 8. `/mp:plan` 后的恢复语义

用户确认的恢复模型是：

- 高复杂度任务被拦截后，保存上下文
- `/mp:plan` 补齐 planning
- 之后允许“继续之前任务”，但不自动直接执行原命令

因此建议语义为：

1. `interrupted_context` 中保存：命令、子命令、prompt、时间戳
2. `/mp:plan` 完成后不自动重放原任务
3. 用户再次执行 `develop/review/...` 时：
   - 若与被中断任务匹配，则返回 `resume_available=true`
   - runtime 正常继续执行当前显式命令
4. 任务成功进入真正执行后，可清空或标记 `interrupted_context` 为已消费

---

## 9. 测试策略

### 9.1 `internal/tracks`

- canonical artifact 完整性检查
- 旧 `spec.md/plan.md` 不应被视为 planning 完整
- metadata 中复杂度/执行模式字段的判定

### 9.2 `internal/validation`

- 高复杂度无 active track -> block `/mp:plan`
- 高复杂度有 track 但缺 `design.md` -> block `/mp:plan`
- 高复杂度 planning 完整但不在 worktree -> block worktree
- 高复杂度 planning 完整且在 worktree -> allow
- 低复杂度 -> allow

### 9.3 `internal/app`

- structured admission result 是否被正确转为 response contract
- `action` / `remediation` / `data` 是否稳定输出

### 9.4 `internal/cli`

- block 时保存 `interrupted_context`
- `/mp:plan` 是否复用同一个 track
- planning 完成后再次执行原命令是否能进入 resume 路径
- runtime 不应偷偷自动执行旧命令

---

## 10. 兼容性与迁移结论

- 当前 `.multipowers/tracks/complexity-scoring-gate-20260307/` 下的 `spec.md` 与 `plan.md` 只能保留为历史讨论稿。
- 新 admission gate 落地后，真正的 planning completeness 只认 canonical artifacts。
- 这与最新 track 模板和 runtime contract 保持一致，也避免在 runtime 里继续维护两套 planning 语义。
