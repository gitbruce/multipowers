# Multipowers 项目 Gap Analysis（第四版）与优化计划

> 基于 `conductor/context` 的方法论约束（Context-First / Role-Driven / Verification-First）与当前代码现状（2026-02-10）进行复盘，聚焦“已稳定后仍未闭合”的工程 gap，并给出可更新任务清单（为什么改、改什么文件、怎么改）。

## 1) 分析基线

- 背景来源（当前可用）：
  - `templates/conductor/context/product.md`
  - `templates/conductor/context/tech-stack.md`
  - `README.md`
  - `docs/design/logging.md`
- 现状核查（本地）：
  - `conductor/context/` 当前目录缺失（工作区无项目背景文件）。
  - `templates/conductor/context/` 仅有 `product.md`、`tech-stack.md` 两个模板，且内容为占位符。
  - `npm test --silent`：通过（核心测试 `10/10`）。
- 结论：当前主要风险已从“命令不可用”转移为“上下文基线、可观测精度、治理自动化”不足。

## 2) 目标状态（结合项目背景）

1. **Context-First**：初始化后即具备完整、可编辑、可校验的 `conductor/context` 基线。
2. **Role-Driven**：日志与运行证据能准确反映真实角色与调用路径。
3. **Verification-First**：测试不污染工作区，状态更新有自动化证据约束。
4. **可持续治理**：文档、模板、CLI 与测试行为保持一致并可长期维护。

## 3) Gap Analysis（当前仍存在）

| Gap ID | 差距描述 | 当前表现 | 风险 | 优先级 |
|---|---|---|---|---|
| GA4-01 | 项目 context 基线缺失 | `conductor/context/` 在当前仓库状态下不存在 | 无法基于真实背景执行角色流程 | P0 |
| GA4-02 | init 模板与文档约定不一致 | 模板仅 2 个 context 文件，README/ask-role 约定至少 4 个高优先级文件 | 新项目初始化后信息缺口大 | P0 |
| GA4-03 | doctor 未校验 context 健康度 | doctor 只聚焦 roles/config/connector，不检查关键 context 文件 | 背景缺失无法前置暴露 | P1 |
| GA4-04 | 测试污染工作区状态 | `test-doctor-init.sh` 直接删改仓库内 `conductor/`，执行后不恢复 | 本地分析与实际状态漂移 | P0 |
| GA4-05 | 结构化日志角色归属不精确 | connectors 中角色字段仍是工具侧常量，不是调用时真实 role | 运行数据失真，排障困难 | P1 |
| GA4-06 | context 裁剪仅 stderr 可见 | 裁剪决策未沉淀到结构化 JSONL | 无法做历史审计与趋势分析 | P1 |
| GA4-07 | Track 名称 slug 对非 ASCII 不友好 | `track new` 清洗后可能得到空/低可读 slug | 中文等场景易生成不稳定文件名 | P2 |
| GA4-08 | 缺少持续集成门禁 | 当前仅本地执行 `npm test`，无仓库级自动门禁 | 回归风险依赖人工执行 | P1 |
| GA4-09 | 计划状态证据缺自动检查 | 证据模板已存在，但未自动校验 `DONE` 合规性 | 仍可能出现“文档完成但事实未完成” | P2 |

## 4) 优化点整理（按优先级，可更新任务列表）

> 状态枚举建议：`TODO` / `IN_PROGRESS` / `BLOCKED` / `DONE`

### T4-001（P0，对应 GA4-01/GA4-02）补齐并固定 context 初始化基线
- **状态**：`DONE`
- **负责人**：Router
- **为什么改**：没有完整 context 基线会直接削弱 Context-First 能力。
- **改什么文件**：
  - 修改：`templates/conductor/context/product.md`
  - 修改：`templates/conductor/context/tech-stack.md`
  - 新增：`templates/conductor/context/product-guidelines.md`
  - 新增：`templates/conductor/context/workflow.md`
  - 新增（可选）：`templates/conductor/context/code-style.md`
  - 新增（可选）：`templates/conductor/context/design-system.md`
- **怎么改**：
  1. 将模板从“占位符文本”升级为可执行骨架（明确必填项与示例）。
  2. 与 `ask-role` 的优先级文件约定对齐（至少 4 个核心文件）。
  3. 在 init 后提供“待补全清单”提示。
- **成功判定**：
  - [x] `./bin/multipowers init --force` 后 `conductor/context/` 至少生成 4 个核心文件。
  - [x] 新模板不再只有占位符句子。

### T4-002（P0，对应 GA4-04）修复测试对仓库状态的污染
- **状态**：`DONE`
- **负责人**：Architect
- **协作**：Coder
- **为什么改**：测试必须可重复且不破坏开发者工作区。
- **改什么文件**：
  - 修改：`tests/opencode/test-doctor-init.sh`
  - 修改：`tests/opencode/setup.sh`
  - 修改：`tests/opencode/run-tests.sh`
- **怎么改**：
  1. 将 `doctor/init` 测试迁移到临时目录执行（复制最小必需文件）。
  2. 若仍需在仓库内执行，增加 trap 备份/恢复 `conductor/`。
  3. 在测试结束后校验“仓库状态未被污染”。
- **成功判定**：
  - [x] 连续两次执行 `npm test --silent` 后，`conductor/` 结构保持不变。
  - [x] `test-doctor-init.sh` 不再直接 `rm -rf` 仓库内长期目录。

### T4-003（P1，对应 GA4-03）为 doctor 增加 context 健康检查
- **状态**：`DONE`
- **负责人**：Router
- **为什么改**：doctor 应覆盖运行前关键前置条件，而不仅是二进制与配置。
- **改什么文件**：
  - 修改：`bin/multipowers`
  - 修改：`README.md`
  - 修改：`tests/opencode/test-doctor-init.sh`
- **怎么改**：
  1. doctor 检查 `conductor/context` 是否存在。
  2. 检查优先级核心文件是否齐全；缺失时给出 `WARN/FAIL + 修复命令`。
  3. 输出“当前 context 文件数量与缺失项”。
- **成功判定**：
  - [x] context 缺失时 doctor 给出明确告警与修复建议。
  - [x] context 完整时 doctor 报告全绿。

### T4-004（P1，对应 GA4-05）修正结构化日志中的真实角色归属
- **状态**：`DONE`
- **负责人**：Coder
- **协作**：Architect
- **为什么改**：角色是审计与统计维度，写死会导致 observability 失真。
- **改什么文件**：
  - 修改：`bin/ask-role`
  - 修改：`connectors/codex.py`
  - 修改：`connectors/gemini.py`
  - 修改：`tests/opencode/test-connector-exit-code.sh`
- **怎么改**：
  1. 在 `ask-role` 调用 connector 时传入 `ROLE`（参数或环境变量）。
  2. connector 使用该 role 写 `log_structured`。
  3. 回归测试验证 `architect/librarian` 调用时日志角色正确。
- **成功判定**：
  - [x] `outputs/runs/*.jsonl` 的 `role` 字段与调用角色一致。

### T4-005（P1，对应 GA4-06）将 context 裁剪决策写入结构化日志
- **状态**：`DONE`
- **负责人**：Coder
- **协作**：Architect
- **为什么改**：仅 stderr 可见会丢失历史证据，无法后续排障与复盘。
- **改什么文件**：
  - 修改：`bin/ask-role`
  - 修改：`connectors/utils.py`
  - 修改：`docs/design/logging.md`
  - 修改：`tests/opencode/test-context-budget-priority.sh`
- **怎么改**：
  1. 增加 context 相关日志字段（如 `context_budget`, `context_token_estimate`, `truncated_files`）。
  2. 裁剪发生时落盘 JSONL，避免只在 stderr 短暂输出。
  3. 测试断言日志字段存在且内容正确。
- **成功判定**：
  - [x] 超预算请求在 JSONL 中可追踪到被裁剪文件清单。

### T4-006（P2，对应 GA4-07）增强 Track slug 对多语言输入的兼容
- **状态**：`DONE`
- **负责人**：Coder
- **为什么改**：中文/特殊字符 feature 名可能生成空 slug，影响可用性。
- **改什么文件**：
  - 修改：`bin/multipowers`
  - 修改：`tests/opencode/test-track-workflow.sh`
  - 修改：`README.md`
- **怎么改**：
  1. slug 清洗后若为空，回退到 `track-<date>-task` 或可读 hash。
  2. 增加中文输入场景测试。
  3. README 明确命名与回退规则。
- **成功判定**：
  - [x] `track new "中文 功能"` 稳定生成合法 track 文件。

### T4-007（P1，对应 GA4-08）建立最小 CI 测试门禁
- **状态**：`DONE`
- **负责人**：Architect
- **为什么改**：仅靠本地执行无法保证每次变更都通过核心回归。
- **改什么文件**：
  - 新增：`.github/workflows/core-tests.yml`
  - 修改：`README.md`
- **怎么改**：
  1. CI 执行 `bash -n` + `python3 -m py_compile` + `npm test --silent`。
  2. 集成测试继续保留为可选 job/手动触发。
  3. 在 README 增加“本地与 CI 一致命令”。
- **成功判定**：
  - [x] PR/推送自动触发核心测试并给出明确状态。

### T4-008（P2，对应 GA4-09）增加计划状态证据自动校验
- **状态**：`DONE`
- **负责人**：Architect
- **协作**：Architect
- **为什么改**：模板存在但无自动检查，状态漂移仍可能复发。
- **改什么文件**：
  - 新增：`scripts/check_plan_evidence.py`
  - 修改：`docs/plans/status-evidence-template.md`
  - 修改：`tests/opencode/run-tests.sh`
- **怎么改**：
  1. 扫描 `docs/plans/gap_analysis_plan*.md` 的 `DONE` 条目是否含证据区块。
  2. 违规时返回非零退出码并输出缺失项。
  3. 将该检查纳入核心测试入口（可作为轻量文档门禁）。
- **成功判定**：
  - [x] 缺失证据的 `DONE` 状态会被自动拦截。

## 5) 推荐执行顺序（两阶段）

### 阶段 A：先稳住基线（P0）
- T4-001（context 基线）
- T4-002（测试隔离）

### 阶段 B：提升可观测与治理（P1/P2）
- T4-003（doctor context 健康检查）
- T4-004（真实角色日志）
- T4-005（裁剪落盘）
- T4-007（CI 门禁）
- T4-006（track slug 多语言）
- T4-008（状态证据自动校验）

## 6) 完成定义（DoD）

1. 初始化后可直接得到完整 context 骨架并通过 doctor 校验。
2. 核心测试执行不再污染工作区。
3. 结构化日志能准确反映角色与 context 裁剪决策。
4. 关键回归由 CI 自动门禁。
5. 计划状态变更具备自动化证据约束。


## 7) 执行与验证证据（2026-02-10）

- **E4-001（T4-001/T4-003）context 基线与 doctor 健康检查**
  - 命令：`bash tests/opencode/test-doctor-init.sh`
  - 退出码：`0`
  - 关键输出：`Doctor reports missing context`、`Init creates required context files`、`Doctor validates context files`

- **E4-002（T4-002）测试隔离与无污染保障**
  - 命令：`bash tests/opencode/test-doctor-init.sh`
  - 退出码：`0`
  - 关键输出：`Repository workspace remains unchanged`、`No repository conductor pollution`

- **E4-003（T4-004/T4-005）角色归属与 context 裁剪结构化日志**
  - 命令：`bash tests/opencode/test-connector-exit-code.sh && bash tests/opencode/test-context-budget-priority.sh`
  - 退出码：`0`
  - 关键输出：`Structured log role matches ask-role caller`、`Structured log contains context truncation details`

- **E4-004（T4-006）Track slug 多语言回退**
  - 命令：`bash tests/opencode/test-track-workflow.sh`
  - 退出码：`0`
  - 关键输出：`Non-ASCII name generated fallback slug`

- **E4-005（T4-007）CI 门禁落地**
  - 命令：`test -f .github/workflows/core-tests.yml && rg -n "npm test --silent|bash -n|py_compile" .github/workflows/core-tests.yml`
  - 退出码：`0`
  - 关键输出：workflow 包含语法检查、Python 编译检查与核心测试执行

- **E4-006（T4-008）计划证据自动校验**
  - 命令：`bash tests/opencode/test-plan-evidence.sh`
  - 退出码：`0`
  - 关键输出：`Existing gap plans pass evidence check`、`Missing evidence is detected`

- **E4-007（全集成回归）**
  - 命令：`npm test --silent`
  - 退出码：`0`
  - 关键输出：`Passed: 12`、`Failed: 0`、`STATUS: PASSED`

### 结构化证据字段（兼容 check_plan_evidence.py）

- **Coverage Task IDs**: `T4-001, T4-002, T4-003, T4-004, T4-005, T4-006, T4-007, T4-008`
- **Date**: `2026-02-10`
- **Verifier**: `architect`
- **Command(s)**:
  - `bash -n bin/multipowers bin/ask-role`
  - `python3 -m py_compile connectors/codex.py connectors/gemini.py connectors/utils.py scripts/validate_roles.py scripts/check_context_quality.py scripts/check_plan_evidence.py`
  - `npm test --silent`
- **Exit Code**: `0`
- **Key Output**:
  - `STATUS: PASSED`

