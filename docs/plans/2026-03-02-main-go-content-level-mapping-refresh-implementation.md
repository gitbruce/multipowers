# Main vs Go Content-Level Mapping Refresh Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 基于最新 `go` 分支代码，重做 `commands/skills/scripts` 的主分支到 go 分支映射，按“语义能力 + 内容差异”修正 `docs/architecture/commands_skills_difference.md` 与 `docs/architecture/script-differences.md`。

**Architecture:** 先生成可复现的分支快照与内容 diff 指标，再按“同名不等于等价”的规则回填映射状态。`commands/skills` 增加内容级分层（文本近似、包装化、显著重写）；`scripts` 重新用最新 go 源码存在性与语义承接路径校准 `partial/missing/equivalent`。最终用覆盖率与计数一致性校验收尾。

**Tech Stack:** Git (`git ls-tree`, `git show`, `git diff --no-index`), Bash (`awk`, `sed`, `comm`), Ripgrep (`rg`), Markdown 文档编辑

---

### Task 1: 生成最新分支快照与可复现差异基线

**Files:**
- Create: `tmp/recheck/main-command-names.txt`
- Create: `tmp/recheck/go-command-names.txt`
- Create: `tmp/recheck/shared-command-names.txt`
- Create: `tmp/recheck/main-skill-names.txt`
- Create: `tmp/recheck/go-skill-names.txt`
- Create: `tmp/recheck/shared-skill-names.txt`
- Create: `tmp/recheck/main-sh.txt`
- Create: `tmp/recheck/go-sh.txt`

**Step 1: 生成命令/技能名称快照**

Run:
```bash
mkdir -p tmp/recheck
git ls-tree -r --name-only main | rg '^\\.claude/commands/.*\\.md$' -N | sed 's#.*/##; s#\\.md$##' | sort > tmp/recheck/main-command-names.txt
git ls-tree -r --name-only go | rg '^\\.claude-plugin/\\.claude/commands/.*\\.md$' -N | sed 's#.*/##; s#\\.md$##' | sort > tmp/recheck/go-command-names.txt
comm -12 tmp/recheck/main-command-names.txt tmp/recheck/go-command-names.txt > tmp/recheck/shared-command-names.txt

git ls-tree -r --name-only main | rg '^\\.claude/skills/.*\\.md$' -N | sed 's#.*/##; s#\\.md$##' | sort > tmp/recheck/main-skill-names.txt
git ls-tree -r --name-only go | rg '^\\.claude-plugin/\\.claude/skills/.*\\.md$' -N | sed 's#.*/##; s#\\.md$##' | sort > tmp/recheck/go-skill-names.txt
comm -12 tmp/recheck/main-skill-names.txt tmp/recheck/go-skill-names.txt > tmp/recheck/shared-skill-names.txt
```
Expected: shared commands=38, shared skills=46。

**Step 2: 生成脚本快照**

Run:
```bash
git ls-tree -r --name-only main | rg '\\.sh$' -N | sort > tmp/recheck/main-sh.txt
git ls-tree -r --name-only go | rg '\\.sh$' -N | sort > tmp/recheck/go-sh.txt
```
Expected: `main-sh=135`, `go-sh=16`（若变化，以当前输出为准并回写文档）。

**Step 3: 校验快照计数**

Run:
```bash
wc -l tmp/recheck/main-command-names.txt tmp/recheck/go-command-names.txt tmp/recheck/shared-command-names.txt tmp/recheck/main-skill-names.txt tmp/recheck/go-skill-names.txt tmp/recheck/shared-skill-names.txt tmp/recheck/main-sh.txt tmp/recheck/go-sh.txt
```
Expected: 命令/技能/脚本计数可复现。

**Step 4: Commit**

```bash
git add tmp/recheck/main-command-names.txt tmp/recheck/go-command-names.txt tmp/recheck/shared-command-names.txt tmp/recheck/main-skill-names.txt tmp/recheck/go-skill-names.txt tmp/recheck/shared-skill-names.txt tmp/recheck/main-sh.txt tmp/recheck/go-sh.txt
git commit -m "chore: add refreshed main-go snapshot baselines for mapping review"
```

### Task 2: 生成 commands 内容级映射指标

**Files:**
- Create: `tmp/recheck/commands-diff-metrics.tsv`
- Create: `tmp/recheck/commands-content-mapping.tsv`

**Step 1: 生成每个 shared command 的内容差异指标**

Run:
```bash
# 输出字段: name, main_lines, go_lines, add, del, norm_equal, go_wrapper
bash -lc '<commands-diff-script>'
```
Expected: `commands-diff-metrics.tsv` 含 38 行 shared command 指标。

**Step 2: 生成状态映射**

规则：
- `norm_equal=1` => `equivalent`
- `go_wrapper=1` 或关键逻辑显著下沉 => `partial`
- main-only => `missing`
- go-only => `intentional-diff`

Run:
```bash
# 生成 commands-content-mapping.tsv: main_name, main_file, go_file_or_target, status, evidence, remediation
bash -lc '<commands-mapping-script>'
```
Expected: 38 shared + main-only + go-only 全量映射行。

**Step 3: 针对关键样例人工复核**

至少复核：
- `octo -> mp`
- `embrace`（wrapper 化）
- `extract`（确认是否同构）
- `plan`（路径与内容是否仅命名替换）

Run:
```bash
git show main:.claude/commands/embrace.md | sed -n '1,120p'
sed -n '1,120p' .claude-plugin/.claude/commands/embrace.md
```
Expected: 样例结论与映射状态一致。

**Step 4: Commit**

```bash
git add tmp/recheck/commands-diff-metrics.tsv tmp/recheck/commands-content-mapping.tsv
git commit -m "chore: add content-level command mapping metrics"
```

### Task 3: 生成 skills 内容级映射指标（重点处理 extract-skill）

**Files:**
- Create: `tmp/recheck/skills-diff-metrics.tsv`
- Create: `tmp/recheck/skills-content-mapping.tsv`

**Step 1: 生成每个 shared skill 的内容差异指标**

Run:
```bash
# 输出字段: name, main_lines, go_lines, add, del, norm_equal, go_thin_wrapper
bash -lc '<skills-diff-script>'
```
Expected: `skills-diff-metrics.tsv` 含 46 行 shared skill 指标。

**Step 2: 生成状态映射**

规则：
- go 文件为 `Thin wrapper skill` => `partial`（语义承接但实现形态已大改）
- 文本与流程保留完整 => `equivalent`
- main-only => `missing`
- go-only => `intentional-diff`

Run:
```bash
# 生成 skills-content-mapping.tsv: main_name, main_file, go_file, status, evidence, remediation
bash -lc '<skills-mapping-script>'
```
Expected: 明确标出 `extract-skill` 为内容显著重写（`partial`）。

**Step 3: 样例复核（用户关注点）**

Run:
```bash
git show main:.claude/skills/extract-skill.md | sed -n '1,220p'
sed -n '1,120p' .claude-plugin/.claude/skills/extract-skill.md
```
Expected: 在映射中有证据说明“主分支为完整实现指南，go 为薄包装/运行时委派”。

**Step 4: Commit**

```bash
git add tmp/recheck/skills-diff-metrics.tsv tmp/recheck/skills-content-mapping.tsv
git commit -m "chore: add content-level skill mapping metrics including extract-skill"
```

### Task 4: 更新 `commands_skills_difference.md`（从名称映射升级到内容映射）

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`

**Step 1: 重写结果概览与状态统计**

要求：
- 拆分“名称级 shared”与“内容级 equivalent/partial/missing/intentional-diff”
- 引入 `extract-skill` 等关键差异案例

**Step 2: 写入全量映射表**

输入来源：
- `tmp/recheck/commands-content-mapping.tsv`
- `tmp/recheck/skills-content-mapping.tsv`

输出要求：
- 每行包含 `main file -> go file/target`
- 每个 `partial/missing` 给出 remediation

**Step 3: 文档一致性校验**

Run:
```bash
rg -n "equivalent|partial|missing|intentional-diff|extract-skill" docs/architecture/commands_skills_difference.md
```
Expected: 状态口径完整且包含 `extract-skill` 证据段。

**Step 4: Commit**

```bash
git add docs/architecture/commands_skills_difference.md
git commit -m "docs: refresh command-skill mapping with content-level parity analysis"
```

### Task 5: 更新 `script-differences.md`（按最新 go 代码回归映射）

**Files:**
- Modify: `docs/architecture/script-differences.md`
- Create: `tmp/recheck/scripts-refresh.tsv`

**Step 1: 以最新 go 源码重新计算脚本承接状态**

Run:
```bash
# 基于 main 135 脚本，重算每条映射 target 是否存在、语义承接状态、整改项
bash -lc '<scripts-refresh-script>'
```
Expected: `scripts-refresh.tsv` 覆盖 main 脚本全集。

**Step 2: 回写按域全量清单**

要求：
- 保持“按域 + 全文件”
- 每条主脚本都给出 `main -> go 目标文件/方法`
- `partial/missing` 必须有整改路径

**Step 3: 覆盖率校验**

Run:
```bash
# 确保 main 脚本 100% 出现在文档表格中
awk '/^\\| `/{print}' docs/architecture/script-differences.md | sed -E 's/^\\| `([^`]+)`.*/\\1/' | sort > tmp/recheck/script-doc-listed.txt
comm -23 tmp/recheck/main-sh.txt tmp/recheck/script-doc-listed.txt
```
Expected: 输出为空（0 漏项）。

**Step 4: Commit**

```bash
git add docs/architecture/script-differences.md tmp/recheck/scripts-refresh.tsv
git commit -m "docs: refresh script mapping against latest go codebase"
```

### Task 6: 交叉一致性校验与最终提交

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`
- Modify: `docs/architecture/script-differences.md`

**Step 1: 校验三类计数一致性**

Run:
```bash
wc -l tmp/recheck/main-command-names.txt tmp/recheck/go-command-names.txt tmp/recheck/shared-command-names.txt tmp/recheck/main-skill-names.txt tmp/recheck/go-skill-names.txt tmp/recheck/shared-skill-names.txt tmp/recheck/main-sh.txt tmp/recheck/go-sh.txt
```
Expected: 文档 summary 与快照一致。

**Step 2: 校验状态分布可复现**

Run:
```bash
rg -n "equivalent=|partial=|missing=|intentional-diff=" docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md
```
Expected: 各文档都有显式统计并与映射表行数一致。

**Step 3: 最终检查并提交**

Run:
```bash
git status --short
git diff -- docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md | sed -n '1,260p'
```
Expected: 仅包含目标文档和证据文件改动。

**Step 4: Final Commit**

```bash
git add docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md tmp/recheck
git commit -m "docs: complete refreshed main-go mapping with content-aware parity"
```
