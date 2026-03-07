# Superpowers 差异分析设计

日期：2026-03-07

## 目标

产出一份基于代码与配置实物的对比文档，回答下面三个问题：

1. `superpowers` 的每个 `command` / `skill` 在本项目中分别由哪些对象承接。
2. 这些承接发生在什么层：`command`、`skill`、`persona/workflow`、还是 Go 原生运行时与治理模块。
3. 两个项目在架构与工程化、以及用户体验与提示设计上各自的优势、代价与短板是什么。

最终交付文件为 `docs/architecture/superpowers-diff.md`。

## 已确认范围

- 比较对象一：`superpowers` 仓库中的 `commands/` 与 `skills/`
- 比较对象二：本项目中的
  - `.claude-plugin/.claude/commands/*`
  - `.claude-plugin/.claude/skills/*`
  - `agents/personas/*`
  - `config/workflows.yaml`、`config/agents.yaml`
  - `internal/workflows/*`
  - `internal/tracks/*`
  - `internal/validation/*`
  - `internal/doctor/*`
  - `internal/policy/*`
  - 相关架构与 CLI 文档

## 判定口径

文档使用以下映射等级：

- `直接对应`：本项目中存在职责和交互方式都比较接近的承接对象。
- `部分对应`：只覆盖部分职责，或者名字相似但执行边界明显不同。
- `间接承接`：没有单一等价物，但多个对象组合后可以承接同类功能。
- `缺失`：本项目没有清晰的一等能力来承接该项。

为了避免“同名即对应”的误判，所有映射都按功能维度拆解：

- 目标任务
- 触发方式
- 交互模式
- 强制约束 / gate
- 执行主体
- 状态持久化
- 验证与治理机制

## 文档结构

`docs/architecture/superpowers-diff.md` 将按以下结构组织：

1. 背景与比较范围
2. 分析方法与判定标准
3. 核心架构差异总览
4. 能力域总览表
5. `superpowers commands` 逐项映射
6. `superpowers skills` 逐项映射
7. 两项目优缺点分析
8. 结论与建议

## 关键设计判断

### 1. 以 `superpowers` 为主键，而不是以本项目为主键

原因是用户希望回答“`superpowers` 每个 command / skill 在本项目中如何映射”，因此主叙事必须从 `superpowers` 出发，再向本项目的多层能力落点展开。

### 2. 不做文件名对文件名映射

本项目存在大量“语义下沉”现象：一些 `superpowers` 中的提示词技能，在本项目已经拆成：

- 提示词入口层（`commands` / `skills`）
- 配置层（`config/workflows.yaml`, `config/agents.yaml`）
- 运行时层（`internal/workflows`, `internal/orchestration`）
- 治理层（`internal/validation`, `internal/doctor`, `internal/policy`）
- 状态层（`internal/tracks` 与 `.multipowers/*`）

因此最终文档会把这些分层承接关系明确写出来。

### 3. 优缺点分析以“架构与工程化”为主，以“用户体验与提示设计”为辅

重点考察：

- 可组合性
- 可执行性
- 可验证性
- 治理能力
- 状态持久化
- 可扩展性与维护成本

同时补充：

- 命令可发现性
- 认知负担
- 交互节奏
- 人在环审批体验
- 语义一致性

## 取证来源

### `superpowers`

- `commands/brainstorm.md`
- `commands/write-plan.md`
- `commands/execute-plan.md`
- `skills/*/SKILL.md`
- `README.md`

### 本项目

- `.claude-plugin/plugin.json`
- `.claude-plugin/.claude/commands/*`
- `.claude-plugin/.claude/skills/*`
- `config/workflows.yaml`
- `config/agents.yaml`
- `internal/workflows/*`
- `internal/tracks/*`
- `internal/validation/*`
- `internal/doctor/*`
- `internal/policy/*`
- `README.md`
- `docs/ARCHITECTURE.md`
- `docs/CLI-REFERENCE.md`
- `docs/PLUGIN-ARCHITECTURE.md`

## 风险与缓解

### 风险 1：名称相似但职责不同

例如本项目的 `/mp:brainstorm` 更像创意 thought partner，而不是 `superpowers:brainstorming` 那种“设计先行、必须审批后才能进入实现”的硬流程。

缓解方式：在每一项映射里显式区分“名称相近”和“功能相近”。

### 风险 2：本项目的能力被拆散在多个层次

例如 `using-git-worktrees` 在 `superpowers` 中是独立技能，而本项目更多通过 `internal/validation/admission.go`、`internal/tracks/worktree_check.go` 和 track/group CLI 进行运行时强制。

缓解方式：每一项同时列出 `command`、`skill`、`workflow/persona`、`runtime` 四类落点。

### 风险 3：优缺点分析容易流于抽象

缓解方式：所有优缺点都回扣到具体文件和机制，不做空泛结论。

## 约束说明

`superpowers:brainstorming` 建议在保存设计后提交 git；当前会话遵循仓库/工具上层约束，不会自动执行 `git commit`，除非用户显式要求。
