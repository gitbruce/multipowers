# Main vs Go Semantic and File Mapping Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 基于 `main` 与 `go` 分支当前真实文件树，修正命令/技能与脚本差异文档，并新增“其他类型文件”差异文档，形成可审计的一对一语义+文件映射结论。

**Architecture:** 先机器化生成 `main`/`go` 文件清单与分类统计，再按“语义能力一对一优先、非双射时补充文件一对一”的规则更新三份架构文档。三份文档共用同一套状态口径（`equivalent` / `partial` / `missing` / `intentional-diff`）与整改字段，最终通过一致性校验闭环。

**Tech Stack:** Git (`git ls-tree`, `git show`), Bash (`sort`, `comm`, `awk`), Ripgrep (`rg`), Markdown 文档编辑

---

### Task 1: 生成 main/go 基础清单与分类基线

**Files:**
- Create: `tmp/main-files.txt`
- Create: `tmp/go-files.txt`
- Create: `tmp/main-sh.txt`
- Create: `tmp/go-sh.txt`
- Create: `tmp/main-commands.txt`
- Create: `tmp/go-commands.txt`
- Create: `tmp/main-skills.txt`
- Create: `tmp/go-skills.txt`

**Step 1: 生成两分支全量文件清单**

Run:
```bash
mkdir -p tmp
git ls-tree -r --name-only main | sort > tmp/main-files.txt
git ls-tree -r --name-only go | sort > tmp/go-files.txt
```
Expected: 两个文件生成且非空。

**Step 2: 生成脚本/命令/技能子清单**

Run:
```bash
rg '\.sh$' tmp/main-files.txt -N > tmp/main-sh.txt
rg '\.sh$' tmp/go-files.txt -N > tmp/go-sh.txt
rg '^\.claude/commands/.*\.md$' tmp/main-files.txt -N > tmp/main-commands.txt
rg '^\.claude-plugin/\.claude/commands/.*\.md$' tmp/go-files.txt -N > tmp/go-commands.txt
rg '^\.claude/skills/.*\.md$' tmp/main-files.txt -N > tmp/main-skills.txt
rg '^\.claude-plugin/\.claude/skills/.*\.md$' tmp/go-files.txt -N > tmp/go-skills.txt
```
Expected: 各清单文件均生成。

**Step 3: 快速校验数量**

Run:
```bash
wc -l tmp/main-files.txt tmp/go-files.txt tmp/main-sh.txt tmp/go-sh.txt tmp/main-commands.txt tmp/go-commands.txt tmp/main-skills.txt tmp/go-skills.txt
```
Expected: 输出各类基线数量，用于后续文档一致性校验。

**Step 4: Commit**

```bash
git add tmp/main-files.txt tmp/go-files.txt tmp/main-sh.txt tmp/go-sh.txt tmp/main-commands.txt tmp/go-commands.txt tmp/main-skills.txt tmp/go-skills.txt
git commit -m "chore: add branch inventory baselines for parity analysis"
```

### Task 2: 建立命令/技能语义与文件映射数据草表

**Files:**
- Create: `tmp/commands-main-names.txt`
- Create: `tmp/commands-go-names.txt`
- Create: `tmp/skills-main-names.txt`
- Create: `tmp/skills-go-names.txt`
- Create: `tmp/commands-main-only.txt`
- Create: `tmp/commands-go-only.txt`
- Create: `tmp/skills-main-only.txt`
- Create: `tmp/skills-go-only.txt`

**Step 1: 提取命令与技能名（去路径与扩展名）**

Run:
```bash
sed 's#^.*/##; s#\.md$##' tmp/main-commands.txt | sort > tmp/commands-main-names.txt
sed 's#^.*/##; s#\.md$##' tmp/go-commands.txt | sort > tmp/commands-go-names.txt
sed 's#^.*/##; s#\.md$##' tmp/main-skills.txt | sort > tmp/skills-main-names.txt
sed 's#^.*/##; s#\.md$##' tmp/go-skills.txt | sort > tmp/skills-go-names.txt
```
Expected: 输出纯名称列表。

**Step 2: 计算双向差异**

Run:
```bash
comm -23 tmp/commands-main-names.txt tmp/commands-go-names.txt > tmp/commands-main-only.txt
comm -13 tmp/commands-main-names.txt tmp/commands-go-names.txt > tmp/commands-go-only.txt
comm -23 tmp/skills-main-names.txt tmp/skills-go-names.txt > tmp/skills-main-only.txt
comm -13 tmp/skills-main-names.txt tmp/skills-go-names.txt > tmp/skills-go-only.txt
```
Expected: 产出“仅 main/仅 go”清单。

**Step 3: 校验清单可读性**

Run:
```bash
for f in tmp/commands-main-only.txt tmp/commands-go-only.txt tmp/skills-main-only.txt tmp/skills-go-only.txt; do echo "== $f =="; cat "$f"; done
```
Expected: 差异集合清晰可用于文档映射。

**Step 4: Commit**

```bash
git add tmp/commands-main-names.txt tmp/commands-go-names.txt tmp/skills-main-names.txt tmp/skills-go-names.txt tmp/commands-main-only.txt tmp/commands-go-only.txt tmp/skills-main-only.txt tmp/skills-go-only.txt
git commit -m "chore: add command and skill name diff artifacts"
```

### Task 3: 更新命令/技能差异文档

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`

**Step 1: 写入文档结构与统一状态口径**

包含：
- 基线与统计（日期、分支、计数）
- 语义一对一主表（commands, skills）
- 非双射文件级补充映射表
- `partial/missing` 的整改路径
- 最终 parity 结论

**Step 2: 用 Task 2 产物填充差异数据**

Run:
```bash
wc -l tmp/main-commands.txt tmp/go-commands.txt tmp/main-skills.txt tmp/go-skills.txt
```
Expected: 与文档 summary 一致。

**Step 3: 文档最小格式校验**

Run:
```bash
rg -n "equivalent|partial|missing|intentional-diff" docs/architecture/commands_skills_difference.md
```
Expected: 至少命中状态字段定义与若干映射行。

**Step 4: Commit**

```bash
git add docs/architecture/commands_skills_difference.md
git commit -m "docs: refresh main vs go command and skill parity mapping"
```

### Task 4: 校验并修正脚本差异文档

**Files:**
- Modify: `docs/architecture/script-differences.md`
- Create: `tmp/script-main-only.txt`
- Create: `tmp/script-go-only.txt`

**Step 1: 生成脚本双向差异**

Run:
```bash
comm -23 tmp/main-sh.txt tmp/go-sh.txt > tmp/script-main-only.txt
comm -13 tmp/main-sh.txt tmp/go-sh.txt > tmp/script-go-only.txt
```
Expected: 产出脚本差异清单。

**Step 2: 逐域核对脚本文档是否覆盖 main 侧脚本**

Run:
```bash
wc -l tmp/main-sh.txt tmp/script-main-only.txt tmp/script-go-only.txt
```
Expected: 文档中“总数/覆盖数/缺口数”可回溯到该输出。

**Step 3: 修正文档使其满足新规则**

要求：
- 语义映射优先；
- 语义非双射时补充文件一对一；
- 每个 `partial/missing` 行有整改目标；
- 明确哪些脚本是 `intentional-diff`（纯 wrapper）。

**Step 4: Commit**

```bash
git add docs/architecture/script-differences.md tmp/script-main-only.txt tmp/script-go-only.txt
git commit -m "docs: validate and correct script parity mapping with file supplements"
```

### Task 5: 新增其他类型差异文档

**Files:**
- Create: `docs/architecture/other-differences.md`
- Create: `tmp/main-other.txt`
- Create: `tmp/go-other.txt`
- Create: `tmp/other-main-only.txt`
- Create: `tmp/other-go-only.txt`

**Step 1: 生成“非命令/技能/脚本”文件集合**

Run:
```bash
grep -vE '^\\.claude/commands/|^\\.claude/skills/|\\.sh$' tmp/main-files.txt > tmp/main-other.txt
grep -vE '^\\.claude-plugin/\\.claude/commands/|^\\.claude-plugin/\\.claude/skills/|\\.sh$' tmp/go-files.txt > tmp/go-other.txt
```
Expected: 形成 other 类目全集。

**Step 2: 生成 other 双向差异**

Run:
```bash
comm -23 tmp/main-other.txt tmp/go-other.txt > tmp/other-main-only.txt
comm -13 tmp/main-other.txt tmp/go-other.txt > tmp/other-go-only.txt
```
Expected: 获取非脚本类差异清单。

**Step 3: 撰写 `other-differences.md`**

包含：
- 按域分组（配置、文档、Go 源码、测试、元数据等）
- 语义映射主表
- 文件级一对一补充表（非双射场景）
- 缺口整改与结论

**Step 4: Commit**

```bash
git add docs/architecture/other-differences.md tmp/main-other.txt tmp/go-other.txt tmp/other-main-only.txt tmp/other-go-only.txt
git commit -m "docs: add main vs go parity analysis for non command-skill-script files"
```

### Task 6: 三文档一致性校验与收尾

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`
- Modify: `docs/architecture/script-differences.md`
- Modify: `docs/architecture/other-differences.md`

**Step 1: 校验三文档基线一致性**

Run:
```bash
rg -n "基线|Date|分支|main|go|equivalent|partial|missing|intentional-diff" docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md
```
Expected: 三文档均具备统一状态口径与基线说明。

**Step 2: 校验“main 侧能力覆盖”结论一致**

人工检查：
- 命令/技能：main 每个能力均在映射表出现；
- 脚本：main 每个脚本在域清单出现；
- 其他：关键类型均有结论与整改项。

**Step 3: 运行最小 Git 检查**

Run:
```bash
git status --short
git diff -- docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md | sed -n '1,220p'
```
Expected: 只出现本次目标变更，且内容与任务目标一致。

**Step 4: Final Commit**

```bash
git add docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md
git commit -m "docs: complete semantic and file-level parity mapping across main vs go"
```
