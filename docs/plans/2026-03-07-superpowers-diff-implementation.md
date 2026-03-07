# Superpowers Diff Documentation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 产出一份从代码层面分析 `superpowers` 与本项目功能映射和优缺点的架构文档，并保存到 `docs/architecture/superpowers-diff.md`。

**Architecture:** 先以 `superpowers` 的每个 `command` / `skill` 为主键提炼其功能职责，再映射到本项目的 `commands`、`skills`、`personas/workflows`、Go 运行时与治理模块。最终文档采用“总览 + 逐项映射 + 优缺点分析”的结构，保证既能快速扫描，也能落到证据层。

**Tech Stack:** Markdown、`rg`、`sed`、仓库内现有架构文档、`superpowers` 上游源码快照

---

## Prerequisites

- [x] 本仓库结构、`docs/AGENTS.md` 和相关架构文档已勘察
- [x] `superpowers` 上游仓库已拉取到 `/tmp/superpowers-obra`
- [x] 用户已确认：逐项功能映射、范围包含 runtime 层、优缺点分析以前者为主

### Task 1: 盘点 `superpowers` 入口面

**Files:**
- Read: `/tmp/superpowers-obra/commands/*.md`
- Read: `/tmp/superpowers-obra/skills/*/SKILL.md`
- Read: `/tmp/superpowers-obra/README.md`

**Step 1: 列出全部 `commands` 与 `skills`**

Run:

```bash
cd /tmp/superpowers-obra && find commands skills -maxdepth 2 -type f | sort
```

Expected: 能看到 3 个 `commands` 与全部 `skills/*/SKILL.md`

**Step 2: 提炼每项功能语义**

Run:

```bash
cd /tmp/superpowers-obra && for f in commands/*.md skills/*/SKILL.md; do echo "=== $f ==="; sed -n '1,80p' "$f"; done
```

Expected: 能提炼出每项的目标任务、触发条件与流程约束

### Task 2: 盘点本项目承接层

**Files:**
- Read: `.claude-plugin/plugin.json`
- Read: `.claude-plugin/.claude/commands/*.md`
- Read: `.claude-plugin/.claude/skills/*.md`
- Read: `config/workflows.yaml`
- Read: `config/agents.yaml`
- Read: `internal/workflows/*.go`
- Read: `internal/tracks/*.go`
- Read: `internal/validation/*.go`
- Read: `internal/doctor/*.go`
- Read: `internal/policy/*.go`

**Step 1: 识别本项目命令与技能表层入口**

Run:

```bash
cd /mnt/f/src/ai/multipowers && find .claude-plugin/.claude/commands .claude-plugin/.claude/skills -maxdepth 1 -type f | sort
```

Expected: 得到本项目所有 `command` / `skill` 资产清单

**Step 2: 识别 workflow、persona 与 runtime 承接点**

Run:

```bash
cd /mnt/f/src/ai/multipowers && rg -n "runWorkflowHelper|TrackCoordinator|EnsureTrackExecution|EnsureSpecAdmission|doctor|policy|DispatchWithFallback|GenerateReport|worktree" internal config README.md docs
```

Expected: 得到 workflow、track、验证、doctor、policy 等关键落点

### Task 3: 建立能力域总览表

**Files:**
- Modify: `docs/architecture/superpowers-diff.md`

**Step 1: 汇总能力域**

将能力域整理为：

- 设计前置
- 实施计划
- 执行编排
- 调试
- 评审
- 验证
- 分支/收尾
- 技能治理
- 状态持久化 / 运行时治理

**Step 2: 写入总览表**

Expected: 读者在进入逐项映射前，先理解两项目的心智模型差异

### Task 4: 写 `commands` 逐项映射

**Files:**
- Modify: `docs/architecture/superpowers-diff.md`

**Step 1: 覆盖 `brainstorm` / `write-plan` / `execute-plan` 三项**

每项必须包含：

- 核心功能
- 本项目对应对象
- 映射等级
- 差异说明
- 代码证据路径

**Step 2: 复核无漏项**

Run:

```bash
cd /mnt/f/src/ai/multipowers && rg -n "brainstorm\.md|write-plan\.md|execute-plan\.md" docs/architecture/superpowers-diff.md
```

Expected: 三项都出现在最终文档中

### Task 5: 写 `skills` 逐项映射

**Files:**
- Modify: `docs/architecture/superpowers-diff.md`

**Step 1: 覆盖全部 `superpowers skills`**

至少覆盖以下项：

- `brainstorming`
- `dispatching-parallel-agents`
- `executing-plans`
- `finishing-a-development-branch`
- `receiving-code-review`
- `requesting-code-review`
- `systematic-debugging`
- `subagent-driven-development`
- `test-driven-development`
- `using-git-worktrees`
- `using-superpowers`
- `verification-before-completion`
- `writing-plans`
- `writing-skills`

**Step 2: 用统一模板落文档**

每项至少写出：

- `command` 落点
- `skill` 落点
- `workflow/persona` 落点
- `runtime/governance` 落点
- 判定与理由

**Step 3: 复核无漏项**

Run:

```bash
cd /mnt/f/src/ai/multipowers && rg -n "brainstorming|dispatching-parallel-agents|executing-plans|finishing-a-development-branch|receiving-code-review|requesting-code-review|systematic-debugging|subagent-driven-development|test-driven-development|using-git-worktrees|using-superpowers|verification-before-completion|writing-plans|writing-skills" docs/architecture/superpowers-diff.md
```

Expected: 全部技能均在文档中出现

### Task 6: 写优缺点分析与结论

**Files:**
- Modify: `docs/architecture/superpowers-diff.md`

**Step 1: 写架构与工程化优缺点**

重点比较：

- 强约束是在 prompt 里还是 runtime 里
- 状态是隐式还是显式
- 执行是认知遵循还是系统 gate
- 扩展能力是加 Markdown 还是加 Go/配置/runtime

**Step 2: 写用户体验与提示设计优缺点**

重点比较：

- 触发自然度
- 审批节奏
- 可发现性
- 认知负担
- 名称与语义一致性

**Step 3: 写结论与建议**

Expected: 结论必须回扣前面逐项映射，不做无证据判断

### Task 7: 最终验证

**Files:**
- Verify: `docs/architecture/superpowers-diff.md`
- Verify: `docs/plans/2026-03-07-superpowers-diff-design.md`
- Verify: `docs/plans/2026-03-07-superpowers-diff-implementation.md`

**Step 1: 检查文件存在**

Run:

```bash
cd /mnt/f/src/ai/multipowers && ls docs/plans/2026-03-07-superpowers-diff-design.md docs/plans/2026-03-07-superpowers-diff-implementation.md docs/architecture/superpowers-diff.md
```

Expected: 三个文件都存在

**Step 2: 检查文档覆盖度**

Run:

```bash
cd /mnt/f/src/ai/multipowers && rg -n "直接对应|部分对应|间接承接|缺失|优点|缺点|结论" docs/architecture/superpowers-diff.md
```

Expected: 能看到映射判定与优缺点分析段落

**Step 3: 检查最终文档可读性**

Run:

```bash
cd /mnt/f/src/ai/multipowers && sed -n '1,260p' docs/architecture/superpowers-diff.md
```

Expected: 结构完整、章节顺序正确、无明显残缺或占位符

## 备注

- 本次工作以文档分析为主，不涉及代码逻辑变更。
- 若后续需要把缺失项转为实现 backlog，可在结论后追加“差距到实现项”跟踪表。
- 如需 git 提交，应由用户显式要求后再执行。
