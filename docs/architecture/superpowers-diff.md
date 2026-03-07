# Superpowers 与 Multipowers 功能映射与架构差异分析

日期：2026-03-07

## 背景与比较范围

本文比较两个项目：

- `superpowers`：以 `commands/*.md` 与 `skills/*/SKILL.md` 为核心的流程型技能库。
- `multipowers`（本项目）：在插件命令/技能之外，还引入 `persona/workflow`、Go 原生运行时、track 状态、验证 gate、doctor 与 policy 子系统。

比较目标不是“同名文件是否存在”，而是回答：`superpowers` 的每个 `command` / `skill`，在本项目中分别由哪些对象从功能层面承接。

## 分析方法与判定标准

### 判定等级

| 判定 | 含义 |
|---|---|
| `直接对应` | 本项目中存在职责和交互方式都比较接近的承接对象。 |
| `部分对应` | 只覆盖部分职责，或者名字相近但工作方式明显不同。 |
| `间接承接` | 没有单一等价物，但多个对象组合可以承接该能力。 |
| `缺失` | 本项目没有清晰的一等能力承接。 |

### 取证层次

本文使用四层证据：

1. **命令层**：`.claude-plugin/.claude/commands/*`
2. **技能层**：`.claude-plugin/.claude/skills/*`
3. **角色/工作流层**：`agents/personas/*`、`config/workflows.yaml`、`internal/workflows/*`
4. **运行时/治理层**：`internal/tracks/*`、`internal/validation/*`、`internal/doctor/*`、`internal/policy/*`

## 核心架构差异先看

| 维度 | Superpowers | Multipowers |
|---|---|---|
| 主执行哲学 | 以 `SKILL.md` 约束 agent 行为，核心是流程 discipline | 以插件入口 + Go runtime 承接流程，核心是 runtime contract |
| 命令层角色 | 很薄，`command` 基本只是 skill wrapper | 混合：一部分仍是 skill wrapper，一部分已变成 Go runtime proxy |
| 流程强制方式 | 依赖 agent 遵循提示词和 checklist | 既有 prompt 约束，也有 admission gate、track gate、doctor、policy |
| 状态管理 | 主要靠对话上下文与 Markdown 计划/清单 | 显式 `.multipowers/tracks/*`、`metadata.json`、`state.json`、JSONL 日志 |
| 运行时边界 | 基本没有独立 runtime 子系统 | 有 `internal/workflows`、`internal/orchestration`、`internal/policy`、`internal/doctor` |
| 人在环节奏 | 强调先设计、先批准、再计划、再执行 | 更偏可自动编排与可恢复执行，人的审批点不总是第一等能力 |

**结论一句话：** `superpowers` 更像“把优秀工程流程写进 agent 脑子里”，`multipowers` 更像“把部分流程下沉到插件与运行时里”。

## 能力域总览

| 能力域 | Superpowers | Multipowers | 观察 |
|---|---|---|---|
| 设计前置 | `brainstorming` 强制先设计再实现 | `/mp:plan`、`/mp:define`、admission/track artifacts | 本项目有规划工件，但“设计先于实现”的对话 discipline 没有 `superpowers` 那么统一 |
| 实施计划 | `writing-plans` | `skill-writing-plans`、`/mp:plan` | 两边都强，但本项目把“战略计划”和“零上下文实施计划”混在不同入口 |
| 执行编排 | `executing-plans`、`subagent-driven-development` | `flow-develop`、`/mp:embrace`、`skill-parallel-agents`、`internal/orchestration` | 本项目运行时更强，但与“按计划逐批执行”不是完全同一语义 |
| 调试 | `systematic-debugging` | `/mp:debug`、`skill-debug`、`/mp:doctor` | 本项目除了 RCA 还有运行时治理诊断 |
| 评审 | `requesting-code-review`、`receiving-code-review` | `/mp:review`、`skill-code-review`、deliver personas | 本项目“请求评审”更强，“接收反馈的技术纪律”更弱 |
| 验证 | `verification-before-completion` | `skill-verify`、`EnsureTrackExecution`、doctor/report | 本项目既有提示词验证，也有 runtime gate |
| 分支/收尾 | `using-git-worktrees`、`finishing-a-development-branch` | worktree admission、`skill-finish-branch`、`/mp:ship`、`/mp:rollback` | 本项目更系统，但入口更分散 |
| 技能治理 | `using-superpowers`、`writing-skills` | `plugin.json` 注册、trigger blocks、各 skill execution_mode | 本项目有插件资产管理，但缺少 `superpowers` 那种元技能自律层 |
| 状态持久化 / 治理 | 相对弱 | `tracks`、`policy`、`doctor`、`autosync` | 这是本项目最明显的系统化增强 |

## `superpowers` Commands 逐项映射

### 总览表

| Superpowers Command | 核心功能 | 本项目主要落点 | 判定 |
|---|---|---|---|
| `commands/brainstorm.md` | 在任何创造性工作前强制进入设计探索 | `/mp:brainstorm`、`/mp:plan`、`/mp:define`、`flow-define`、`EnsureSpecAdmission` | `部分对应` |
| `commands/write-plan.md` | 调起零上下文、可执行的实施计划生成 | `skill-writing-plans`、`/mp:plan` | `部分对应` |
| `commands/execute-plan.md` | 按既有计划分批执行并在批次间汇报 | `/mp:embrace`、`flow-develop`、`skill-parallel-agents`、`internal/orchestration` | `部分对应` |

### `brainstorm.md`

- **核心功能**：命令本身极薄，只负责强制调起 `brainstorming` skill；语义重点不是“创意发散”，而是“先设计、后实现”。
- **本项目对应对象**：
  - 命令层：`.claude-plugin/.claude/commands/brainstorm.md`、`.claude-plugin/.claude/commands/plan.md`、`.claude-plugin/.claude/commands/define.md`
  - 技能层：`.claude-plugin/.claude/skills/skill-thought-partner.md`、`.claude-plugin/.claude/skills/flow-define.md`
  - 运行时层：`internal/validation/admission.go`、`internal/tracks/*`
- **判定：`部分对应`**
- **原因**：
  - 本项目名称最接近的是 `/mp:brainstorm`，但它更偏创意 thought partner，而不是 `superpowers:brainstorming` 的设计 gate。
  - 真正接近“形成规格并留下工件”的，是 `/mp:plan`、`/mp:define` 与 track artifacts（`intent.md`、`design.md`、`implementation-plan.md`）。
  - 本项目用 admission gate 强制高复杂度任务先有计划工件，但不是统一由 `brainstorm` 命令承担。

### `write-plan.md`

- **核心功能**：调起 `writing-plans` skill，产出细粒度、零上下文、可交接的实施计划。
- **本项目对应对象**：
  - 技能层：`.claude-plugin/.claude/skills/skill-writing-plans.md`
  - 命令层：`.claude-plugin/.claude/commands/plan.md`
  - 状态层：`docs/plans/*`、`.multipowers/tracks/<track_id>/implementation-plan.md`
- **判定：`部分对应`**
- **原因**：
  - 从 skill 功能看，`skill-writing-plans` 和 `superpowers:writing-plans` 非常接近，都是“零上下文 + bite-sized tasks”。
  - 但从命令语义看，本项目 `/mp:plan` 更宽：它先做 intent capture，再决定是原生 plan mode 还是 multi-AI orchestration，而不是单纯输出实施计划。

### `execute-plan.md`

- **核心功能**：调起 `executing-plans`，按现成计划分批执行、批间汇报并等待反馈。
- **本项目对应对象**：
  - 命令层：`.claude-plugin/.claude/commands/embrace.md`、`.claude-plugin/.claude/commands/develop.md`
  - 技能层：`.claude-plugin/.claude/skills/skill-parallel-agents.md`、`.claude-plugin/.claude/skills/flow-develop.md`
  - 运行时层：`internal/workflows/*`、`internal/orchestration/*`、`internal/tracks/progress.go`
- **判定：`部分对应`**
- **原因**：
  - 本项目已经有比 `executing-plans` 更强的执行引擎，但主要围绕 workflow orchestration，而不是“读取实施计划文件后按 3 个任务一批执行”。
  - 因此能力更强，但语义不完全同构。

## `superpowers` Skills 逐项映射

### 总览表

| Superpowers Skill | 本项目主要落点 | 判定 |
|---|---|---|
| `brainstorming` | `/mp:plan`、`/mp:define`、`/mp:brainstorm`、track/admission | `部分对应` |
| `dispatching-parallel-agents` | `skill-parallel-agents`、`flow-*`、`internal/orchestration` | `直接对应` |
| `executing-plans` | `/mp:embrace`、`flow-develop`、`internal/orchestration`、tracks | `部分对应` |
| `finishing-a-development-branch` | `skill-finish-branch`、`/mp:ship`、`/mp:rollback` | `直接对应` |
| `receiving-code-review` | 零散承接于 review/issue loop，无独立 skill | `缺失` |
| `requesting-code-review` | `/mp:review`、`skill-code-review`、deliver personas | `直接对应` |
| `systematic-debugging` | `/mp:debug`、`skill-debug`、`/mp:doctor` | `直接对应` |
| `subagent-driven-development` | `flow-develop`、`skill-parallel-agents`、`internal/orchestration` | `部分对应` |
| `test-driven-development` | `/mp:tdd`、`skill-tdd`、`tdd-orchestrator` persona | `直接对应` |
| `using-git-worktrees` | `EnsureSpecAdmission`、`IsLinkedWorktreeCheckout`、track group lifecycle | `间接承接` |
| `using-superpowers` | `plugin.json` 技能注册、trigger blocks、命令自动装载 skill | `部分对应` |
| `verification-before-completion` | `skill-verify`、track 完成字段、doctor/report | `直接对应` |
| `writing-plans` | `skill-writing-plans`、`/mp:plan`、track plan artifacts | `直接对应` |
| `writing-skills` | 无直接对应；`extract-skill` 用途不同 | `缺失` |

### 1. `brainstorming`

- **核心功能**：在任何实现前先完成上下文探索、单问题澄清、方案比较、分段审批，并产出设计文档。
- **本项目落点**：
  - 命令：`.claude-plugin/.claude/commands/plan.md`、`.claude-plugin/.claude/commands/define.md`、`.claude-plugin/.claude/commands/brainstorm.md`
  - 技能：`flow-define.md`、`skill-thought-partner.md`
  - 运行时：`internal/validation/admission.go`、`internal/tracks/template_renderer.go`
- **判定：`部分对应`**
- **说明**：
  - 本项目有计划工件和复杂度 gate，但没有一个与 `superpowers:brainstorming` 等价的“先设计、再批准、再进入实现”的统一元流程。
  - `/mp:brainstorm` 在本项目是创意探索工具，不是设计审批 gate。

### 2. `dispatching-parallel-agents`

- **核心功能**：遇到多个独立问题时并行分发 agent，避免串行调查浪费时间。
- **本项目落点**：
  - 技能：`.claude-plugin/.claude/skills/skill-parallel-agents.md`
  - 工作流：`flow-discover.md`、`flow-define.md`、`flow-develop.md`、`flow-deliver.md`
  - 运行时：`internal/workflows/adapter_helper.go`、`internal/orchestration/*`
  - 路由：`internal/providers/router_intent.go`
- **判定：`直接对应`**
- **说明**：
  - 本项目不仅能做并行 agent dispatch，还把它扩展成多 provider、多 phase 的 orchestration runtime。
  - 相比 `superpowers`，本项目的并行能力更系统，但也更重。

### 3. `executing-plans`

- **核心功能**：加载现成计划、批量执行、批间汇报、等待架构师反馈。
- **本项目落点**：
  - 命令：`/mp:embrace`、`/mp:develop`
  - 技能：`flow-develop.md`、`skill-parallel-agents.md`
  - 运行时：`internal/orchestration/*`、`internal/tracks/progress.go`
- **判定：`部分对应`**
- **说明**：
  - 本项目更偏 workflow runtime，而不是“严格读取某个实施计划并按 batch 执行”。
  - 如果按执行引擎能力看，本项目更强；如果按“计划驱动的人在环节奏”看，`superpowers` 更明确。

### 4. `finishing-a-development-branch`

- **核心功能**：实现完成后跑验证、给出 merge / PR / keep / discard 选项，并处理 worktree 收尾。
- **本项目落点**：
  - 技能：`.claude-plugin/.claude/skills/skill-finish-branch.md`
  - 命令：`.claude-plugin/.claude/commands/ship.md`、`.claude-plugin/.claude/commands/rollback.md`
  - 验证：`.claude-plugin/.claude/skills/skill-verify.md`
  - worktree/runtime：`internal/tracks/worktree_check.go`
- **判定：`直接对应`**
- **说明**：
  - 两边目标几乎一致，本项目还补上了 checkpoint/rollback 能力。

### 5. `receiving-code-review`

- **核心功能**：收到 review feedback 后，先验证反馈是否适用于当前代码库，再决定接受、澄清或技术性反驳。
- **本项目落点**：
  - 零散存在于 `/mp:review`、`skill-code-review`、问题跟踪/循环命令周边
- **判定：`缺失`**
- **说明**：
  - 本项目有很强的“发起评审”能力，但缺少一个与 `receiving-code-review` 对应的显式元技能，来规范“如何接收反馈、如何避免表演式赞同、如何技术性 push back”。
  - 这是本项目相对 `superpowers` 的一个真实缺口。

### 6. `requesting-code-review`

- **核心功能**：在关键节点主动请求代码评审，及时暴露问题。
- **本项目落点**：
  - 命令：`.claude-plugin/.claude/commands/review.md`、`.claude-plugin/.claude/commands/deliver.md`
  - 技能：`.claude-plugin/.claude/skills/skill-code-review.md`
  - persona：`agents/personas/code-reviewer.md`、`agents/personas/security-auditor.md`、`agents/personas/performance-engineer.md`
  - 运行时：`internal/workflows/deliver.go`
- **判定：`直接对应`**
- **说明**：
  - 本项目在这一项上是增强版：不仅能 review，还能并发做质量、安全、性能审查。

### 7. `systematic-debugging`

- **核心功能**：按阶段化 RCA 流程排障，而不是凭直觉乱改。
- **本项目落点**：
  - 命令：`.claude-plugin/.claude/commands/debug.md`
  - 技能：`.claude-plugin/.claude/skills/skill-debug.md`
  - 辅助运行时：`internal/workflows/test_run.go`、`internal/doctor/*`
- **判定：`直接对应`**
- **说明**：
  - `skill-debug` 与 `systematic-debugging` 的精神高度一致。
  - 额外的 `/mp:doctor` 不是同义替代，而是本项目增加的运行时/环境诊断面。

### 8. `subagent-driven-development`

- **核心功能**：按计划逐任务派发实现 agent，并在每个任务后做 spec review + code quality review。
- **本项目落点**：
  - 技能：`skill-parallel-agents.md`、`flow-develop.md`
  - 工作流：`internal/workflows/develop.go`
  - 运行时：`internal/orchestration/*`
- **判定：`部分对应`**
- **说明**：
  - 本项目能调度多 agent / 多 provider，但没有把“每任务 fresh subagent + 两阶段 review loop”做成与 `superpowers` 同样显式的固定套路。
  - 它更像广义 orchestration，而不是严格的 per-task subagent discipline。

### 9. `test-driven-development`

- **核心功能**：强制 RED → GREEN → REFACTOR，禁止先写生产代码。
- **本项目落点**：
  - 命令：`.claude-plugin/.claude/commands/tdd.md`
  - 技能：`.claude-plugin/.claude/skills/skill-tdd.md`
  - persona：`agents/personas/tdd-orchestrator.md`
  - 辅助：`internal/workflows/test_run.go`
- **判定：`直接对应`**
- **说明**：
  - 这一项几乎是明确继承并扩展：不仅保留 TDD discipline，还给了 persona 与测试执行辅助。

### 10. `using-git-worktrees`

- **核心功能**：在进入实现前使用隔离 worktree，避免主工作区污染，并验证基线。
- **本项目落点**：
  - 运行时：`internal/validation/admission.go`、`internal/tracks/worktree_check.go`
  - 状态层：`internal/tracks/progress.go`
  - 文档/CLI：`docs/CLI-REFERENCE.md` 中的 `mp track group-start` / `group-complete`
- **判定：`间接承接`**
- **说明**：
  - 本项目没有等价的显式“使用 worktree 技能”，而是把这件事下沉成 runtime gate：高复杂度任务若要求 worktree，会直接阻断非 linked worktree 执行。
  - 就系统强度来说本项目更强；就可理解性和用户引导来说 `superpowers` 更直接。

### 11. `using-superpowers`

- **核心功能**：要求 agent 在任何响应前先检查并调用适用 skill，是一层元流程治理。
- **本项目落点**：
  - 资产注册：`.claude-plugin/plugin.json`
  - 各 skill 的 `trigger:`、`execution_mode:`、`validation_gates:`
  - 各命令中的 “Auto-loads the skill” 模式
- **判定：`部分对应`**
- **说明**：
  - 本项目有插件注册和 trigger 机制，但没有一个与 `using-superpowers` 等价的“元技能自律协议”。
  - 换句话说，本项目更像平台分发与命令路由，`superpowers` 更像 agent 行为宪法。

### 12. `verification-before-completion`

- **核心功能**：在声称完成前必须运行并阅读验证命令输出。
- **本项目落点**：
  - 技能：`.claude-plugin/.claude/skills/skill-verify.md`
  - track gate：`internal/validation/gates.go`
  - group lifecycle：`internal/tracks/progress.go`
  - doctor/report：`internal/doctor/*`
- **判定：`直接对应`**
- **说明**：
  - 本项目不止保留了“先验证再宣称完成”的提示词纪律，还把部分验证前提写进了 track 元数据和 runtime gate。

### 13. `writing-plans`

- **核心功能**：面向零上下文执行者，输出可逐步实现的详细计划。
- **本项目落点**：
  - 技能：`.claude-plugin/.claude/skills/skill-writing-plans.md`
  - 命令：`.claude-plugin/.claude/commands/plan.md`
  - 工件：`.multipowers/tracks/<track_id>/implementation-plan.md`
- **判定：`直接对应`**
- **说明**：
  - 本项目明确拥有同类能力，而且保留了“bite-sized tasks、零上下文、如何验证”的写法。
  - 不同点在于，本项目额外加入了 intent capture、workflow routing 与 track artifacts。

### 14. `writing-skills`

- **核心功能**：把“写 skill”本身流程化，要求用 TDD 思路验证技能文档是否真正改变 agent 行为。
- **本项目落点**：
  - 最近似但不等价的是 `.claude-plugin/.claude/skills/extract-skill.md`，其目标是做设计/产品逆向提炼，而不是编写技能系统本身。
- **判定：`缺失`**
- **说明**：
  - 本项目缺少一个与 `writing-skills` 对等的“技能工程学”文档与测试流程。

## 本项目中超出 `superpowers` 的能力

这些能力在 `superpowers` 中没有明显一等对应物，但构成了本项目的重要差异：

- **Track runtime**：`internal/tracks/*` + `.multipowers/tracks/<track_id>/*`
- **Spec admission gate**：`internal/validation/admission.go`
- **Execution gate**：`internal/validation/gates.go`
- **Doctor diagnostics**：`internal/doctor/*` + `/mp:doctor`
- **Policy routing / fallback**：`internal/policy/*` + `config/providers.yaml` / `config/workflows.yaml`
- **Persona runtime**：`/mp:persona` + `agents/personas/*`
- **Workflow adapter + orchestration report**：`internal/workflows/*`、`internal/orchestration/*`
- **Structured lifecycle logs**：`docs/ARCHITECTURE.md` 描述的 `trace_id` 与 JSONL 日志

这说明本项目不是单纯“把 `superpowers` 换个名字”，而是把一部分 skill-driven discipline 下沉成了 runtime system。

## 两个项目各自的优缺点

### 一、架构与工程化

#### `superpowers` 的优点

- **流程 discipline 非常清晰**：`brainstorming → writing-plans → executing-plans/subagent-driven-development → finishing-a-development-branch` 形成了完整主线。
- **元治理显式**：`using-superpowers`、`verification-before-completion`、`receiving-code-review` 直接约束 agent 的行为方式，而不是只约束任务结果。
- **可移植性强**：技能主要是 Markdown，可跨代理环境迁移，修改成本低。
- **人类可审查性高**：工作流规范本身就是文档，读文档就能理解系统的“宪法”。

#### `superpowers` 的缺点

- **强制力主要依赖 agent 自觉**：是否真的遵守流程，仍取决于模型是否按指令行动。
- **状态与执行证据较弱**：缺少类似 track metadata、doctor report、policy artifacts 的系统化状态层。
- **复杂执行场景支撑有限**：并行、worktree、路由、重试、报告等工程能力更多停留在流程建议层。
- **上下文开销高**：大型 skill 文档会持续占用 prompt 上下文。

#### `multipowers` 的优点

- **运行时承接更强**：`internal/orchestration`、`internal/tracks`、`internal/validation`、`internal/doctor` 把流程从“建议”变成了“contract”。
- **状态显式且可恢复**：track artifacts、group lifecycle、`state.json`、JSONL 日志使长流程更可追踪。
- **治理能力强**：doctor、policy validation、fallback routing、admission gate 使系统更适合复杂项目和团队化约束。
- **扩展面更宽**：既能加 prompt skill，也能加 workflow、persona、Go runtime、policy、doctor check。

#### `multipowers` 的缺点

- **系统复杂度更高**：命令、技能、persona、workflow、runtime、policy、doctor、tracks 多层并存，理解门槛显著上升。
- **语义分散**：同一个能力可能分散在多个入口，例如“计划”既可能来自 `/mp:plan`，也可能来自 `skill-writing-plans` 与 track artifacts。
- **部分工程 discipline 没有一等表达**：例如 `receiving-code-review`、`using-superpowers`、`writing-skills` 这类元流程能力较弱。
- **可移植性较弱**：Go runtime、插件命令与配置编译链使它比 `superpowers` 更依赖特定平台和工程环境。

### 二、用户体验与提示设计

#### `superpowers` 的优点

- **主线明确**：先问清楚、再设计、再计划、再执行，用户不容易迷路。
- **审批点清楚**：设计分段审批、执行分批汇报，让“人在环”体验很强。
- **语言一致性好**：command 与 skill 基本同义，名字和行为关系清晰。

#### `superpowers` 的缺点

- **可能显得过于严格或慢**：对于小任务，持续澄清与审批会增加交互回合。
- **命令面较小**：更多依赖 skill 自驱而不是显式工具面板，不如插件型系统那样“命令即能力”。

#### `multipowers` 的优点

- **命令面丰富**：`/mp:plan`、`/mp:review`、`/mp:debug`、`/mp:doctor`、`/mp:persona`、`/mp:status` 等具备很强的可操作性。
- **自动路由更强**：可基于 intent、workflow、provider policy 做更复杂的分派。
- **可观测性更好**：status、doctor、track artifacts、report 让用户更容易看见系统状态。

#### `multipowers` 的缺点

- **可发现性与认知负担并存**：命令很多，但语义边界不总是直观。
- **同名不等义问题更明显**：例如 `/mp:brainstorm` 与 `superpowers:brainstorming` 看似同名，实际职责差异很大。
- **隐藏逻辑更多**：一些关键约束在 Go 代码与 runtime 里，用户不如读 `superpowers` skill 那样一眼看懂。

## 结论与建议

### 结论

1. **本项目不是对 `superpowers` 的一对一镜像，而是“提示词工作流 + Go runtime”混合系统。**
2. **可直接映射的强项**主要集中在：并行编排、TDD、调试、请求评审、完成前验证、分支收尾。
3. **明显弱于 `superpowers` 的地方**主要是元流程 discipline：
   - 没有等价的 `receiving-code-review`
   - 没有等价的 `writing-skills`
   - 没有和 `using-superpowers` 同等级别的“先查 skill 再行动”元治理
   - 没有由单一入口统一承担的 `brainstorming` 设计 gate
4. **明显强于 `superpowers` 的地方**是 runtime/gov 层：track、admission、doctor、policy、structured logs、persona runtime。

### 建议

- **建议补一个“review feedback handling”技能**：把 `receiving-code-review` 的纪律补成一等能力。
- **建议补一个“meta skill governance”层**：显式定义何时必须优先进入 `plan/define/debug/tdd/verify`，减少入口歧义。
- **建议重命名或重写 `/mp:brainstorm` 的说明**：避免它与 `superpowers:brainstorming` 在语义上产生误导。
- **建议把 `skill-writing-plans` 与 `/mp:plan` 的边界写得更清楚**：一个负责战略 intent routing，一个负责零上下文实施计划。

## 主要证据路径

### `superpowers`

- `/tmp/superpowers-obra/commands/brainstorm.md`
- `/tmp/superpowers-obra/commands/write-plan.md`
- `/tmp/superpowers-obra/commands/execute-plan.md`
- `/tmp/superpowers-obra/skills/brainstorming/SKILL.md`
- `/tmp/superpowers-obra/skills/dispatching-parallel-agents/SKILL.md`
- `/tmp/superpowers-obra/skills/executing-plans/SKILL.md`
- `/tmp/superpowers-obra/skills/receiving-code-review/SKILL.md`
- `/tmp/superpowers-obra/skills/requesting-code-review/SKILL.md`
- `/tmp/superpowers-obra/skills/systematic-debugging/SKILL.md`
- `/tmp/superpowers-obra/skills/test-driven-development/SKILL.md`
- `/tmp/superpowers-obra/skills/using-git-worktrees/SKILL.md`
- `/tmp/superpowers-obra/skills/using-superpowers/SKILL.md`
- `/tmp/superpowers-obra/skills/verification-before-completion/SKILL.md`
- `/tmp/superpowers-obra/skills/writing-plans/SKILL.md`
- `/tmp/superpowers-obra/skills/writing-skills/SKILL.md`

### 本项目

- `.claude-plugin/plugin.json`
- `.claude-plugin/.claude/commands/brainstorm.md`
- `.claude-plugin/.claude/commands/plan.md`
- `.claude-plugin/.claude/commands/develop.md`
- `.claude-plugin/.claude/commands/debug.md`
- `.claude-plugin/.claude/commands/tdd.md`
- `.claude-plugin/.claude/commands/review.md`
- `.claude-plugin/.claude/commands/embrace.md`
- `.claude-plugin/.claude/commands/doctor.md`
- `.claude-plugin/.claude/commands/persona.md`
- `.claude-plugin/.claude/skills/skill-writing-plans.md`
- `.claude-plugin/.claude/skills/skill-tdd.md`
- `.claude-plugin/.claude/skills/skill-debug.md`
- `.claude-plugin/.claude/skills/skill-code-review.md`
- `.claude-plugin/.claude/skills/skill-verify.md`
- `.claude-plugin/.claude/skills/skill-finish-branch.md`
- `.claude-plugin/.claude/skills/skill-parallel-agents.md`
- `.claude-plugin/.claude/skills/skill-thought-partner.md`
- `config/workflows.yaml`
- `config/agents.yaml`
- `internal/workflows/adapter_helper.go`
- `internal/workflows/discover.go`
- `internal/workflows/define.go`
- `internal/workflows/develop.go`
- `internal/workflows/deliver.go`
- `internal/workflows/debate.go`
- `internal/workflows/embrace.go`
- `internal/tracks/worktree_check.go`
- `internal/tracks/template_renderer.go`
- `internal/tracks/progress.go`
- `internal/validation/admission.go`
- `internal/validation/gates.go`
- `internal/doctor/runner.go`
- `internal/policy/validate.go`
- `docs/ARCHITECTURE.md`
- `docs/CLI-REFERENCE.md`
- `docs/PLUGIN-ARCHITECTURE.md`
