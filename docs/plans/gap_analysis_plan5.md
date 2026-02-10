# Multipowers 项目 Gap Analysis（第五版）与优化计划

> 基于 `conductor/context` 方法论（Context-First / Role-Driven / Verification-First）与 2026-02-10 当前仓库状态，识别在 `plan4` 收口后仍存在的执行与治理 gap，并给出可落地优化任务（为什么改、改什么文件、怎么改）。

## 1) 分析基线

- 背景依据：
  - `conductor/context/*.md`（当前工作区缺失）
  - `templates/conductor/context/*.md`
  - `README.md`
  - `docs/design/logging.md`
- 现状核查（2026-02-10）：
  - `./bin/multipowers doctor` → 退出码 `1`（`context directory missing: conductor/context`）
  - `npm test --silent` → 通过（`Passed: 12, Failed: 0`）
  - `python3 scripts/check_plan_evidence.py` → 通过（`PASS (4 files)`）
- 关键观察：
  1. 本地仓库默认（private mode）下 `conductor/context` 不在版本控制，导致 doctor 与测试结果的“体感一致性”仍有偏差。
  2. `ask-role` 与 `doctor` 对 context 缺失的处理强度不一致（一个阻断、一个告警）。

## 2) 目标状态（延续 context 约束）

1. Context 缺失时，所有核心入口给出一致且可执行的处理策略。
2. 初始化与修复流程默认安全，不以破坏性操作换取“简单”。
3. 日志与证据具备请求级可追踪性，而非仅事件级记录。
4. 测试与真实上手路径一致，能暴露“文档可行但工程不可行”的偏差。

## 3) Gap Analysis（当前仍待解决）

| Gap ID | 差距描述 | 当前表现 | 风险 | 优先级 |
|---|---|---|---|---|
| GA5-01 | context 执行策略不一致 | doctor 对 context 缺失直接失败；ask-role 仅告警后继续运行 | 运行语义不稳定、团队认知分裂 | P0 |
| GA5-02 | init 修复路径偏破坏性 | `init --force` 通过整体删除 `conductor/` 重建 | 易误删本地 track/config/context | P0 |
| GA5-03 | context 模板虽齐但缺“完成度校验” | 模板可初始化，但占位内容可直接进入运行链路 | 低质量 context 进入生产流程 | P1 |
| GA5-04 | 根目录上手与测试通过存在落差 | 核心测试全绿，但根目录 doctor 默认失败（缺 context） | 新用户误判“项目已可直接使用” | P1 |
| GA5-05 | 证据检查规则仍偏宽松 | `check_plan_evidence.py` 主要按 Task ID 文本包含判断 | 证据“伪覆盖”风险仍在 | P1 |
| GA5-06 | 结构化日志缺请求关联键 | `context_prepared` 与 connector 记录缺统一 request_id | 难以还原单次调用全链路 | P1 |
| GA5-07 | Track 歧义匹配提示仍可改进 | 多匹配时 helper 会列出候选，但调用层仍补充“not found”类提示 | 用户排障路径不清晰 | P2 |
| GA5-08 | CI 仅覆盖核心门禁，缺治理/集成分层策略 | 当前 workflow 仅 core；集成与治理检查未分层触发 | 回归防线粒度不足 | P2 |

## 4) 优化点整理（按优先级，可更新任务列表）

> 状态建议：`TODO` / `IN_PROGRESS` / `BLOCKED` / `DONE`

### T5-001（P0，对应 GA5-01）统一 context 缺失时的执行策略
- **状态**：`DONE`
- **负责人**：Router
- **为什么改**：核心入口（doctor/ask-role）必须保持一致语义，否则无法形成稳定心智模型。
- **改什么文件**：
  - 修改：`bin/ask-role`
  - 修改：`bin/multipowers`
  - 修改：`README.md`
  - 修改：`tests/opencode/test-ask-role-core.sh`
- **怎么改**：
  1. 定义单一策略：`strict`（缺 context 阻断）或 `lenient`（告警继续），并通过显式参数/环境变量控制。
  2. doctor 与 ask-role 共用同一判断函数或同一规则文案。
  3. README 明确默认策略与切换方式。
- **成功判定**：
  - [x] doctor 与 ask-role 在“缺 context”场景输出一致级别（FAIL/WARN）与一致修复指令。
  - [x] 回归测试覆盖两条入口语义一致性。

### T5-002（P0，对应 GA5-02）提供非破坏性 context 修复模式
- **状态**：`DONE`
- **负责人**：Coder
- **协作**：Architect
- **为什么改**：`init --force` 目前会重建整个 `conductor/`，对已有本地信息风险过高。
- **改什么文件**：
  - 修改：`bin/multipowers`
  - 修改：`README.md`
  - 修改：`tests/opencode/test-doctor-init.sh`
- **怎么改**：
  1. 新增 `init --repair`（仅补齐缺失文件，不覆盖已有文件）。
  2. `init --force` 增加显式确认（或 dry-run 提示）以防误删。
  3. 测试覆盖“保留 tracks/config + 补齐 context”场景。
- **成功判定**：
  - [x] 能在不删除现有 `conductor/tracks` 与 `conductor/config` 的情况下补齐 context。
  - [x] `--force` 行为可预期且具备防误操作提示。

### T5-003（P1，对应 GA5-03）新增 context 完成度校验器
- **状态**：`DONE`
- **负责人**：Architect
- **协作**：Router
- **为什么改**：模板齐全不等于内容可用，需要防止占位文本直接上线。
- **改什么文件**：
  - 新增：`scripts/check_context_quality.py`
  - 修改：`bin/multipowers`
  - 修改：`templates/conductor/context/*.md`
  - 新增：`tests/opencode/test-context-quality.sh`
- **怎么改**：
  1. 校验占位符（如 `[Project Name]`、`[... ]`）是否仍存在。
  2. 校验核心章节最小信息量（如非空 bullet 数）。
  3. doctor 增加 quality 检查结果与修复建议。
- **成功判定**：
  - [x] 含占位符的 context 会被明确标红或告警。
  - [x] 完整 context 可通过 quality 检查。

### T5-004（P1，对应 GA5-04）增加“从新克隆到可运行”烟雾测试
- **状态**：`DONE`
- **负责人**：Architect
- **为什么改**：需要确保测试通过与真实上手路径一致，避免“测试绿、上手红”。
- **改什么文件**：
  - 新增：`tests/opencode/test-onboarding-smoke.sh`
  - 修改：`tests/opencode/run-tests.sh`
  - 修改：`README.md`
- **怎么改**：
  1. 在临时目录模拟 fresh clone 最小环境。
  2. 验证流程：doctor（预期失败）→ init/repair → doctor（预期通过）。
  3. 将该测试纳入核心或 pre-core 快速集。
- **成功判定**：
  - [x] 可稳定复现并验证“首次接入”完整流程。
  - [x] 文档步骤与自动化烟雾测试严格一致。

### T5-005（P1，对应 GA5-05）强化计划证据校验规则
- **状态**：`DONE`
- **负责人**：Architect
- **协作**：Architect
- **为什么改**：当前规则偏文本匹配，存在误判“已覆盖”的窗口。
- **改什么文件**：
  - 修改：`scripts/check_plan_evidence.py`
  - 修改：`docs/plans/status-evidence-template.md`
  - 修改：`tests/opencode/test-plan-evidence.sh`
- **怎么改**：
  1. 引入结构化解析：DONE 任务必须被 `Coverage Task IDs` 显式覆盖。
  2. 校验 evidence block 必填字段（Date/Verifier/Command/Exit Code/Key Output）。
  3. 扩展检查范围到 `docs/plans/*.md`（可配置排除列表）。
- **成功判定**：
  - [x] 缺字段或伪覆盖会被可靠拦截。
  - [x] 合规计划文档可稳定通过检查。

### T5-006（P1，对应 GA5-06）为单次调用增加 request_id 级链路追踪
- **状态**：`DONE`
- **负责人**：Coder
- **协作**：Architect
- **为什么改**：没有 request_id，context 日志与 connector 日志难以归并。
- **改什么文件**：
  - 修改：`bin/ask-role`
  - 修改：`connectors/utils.py`
  - 修改：`connectors/codex.py`
  - 修改：`connectors/gemini.py`
  - 修改：`docs/design/logging.md`
  - 修改：`tests/opencode/test-context-budget-priority.sh`
- **怎么改**：
  1. ask-role 为每次调用生成 `request_id`（UUID/时间戳+随机串）。
  2. 通过环境变量将 `request_id` 传递给 connector 与 context-log。
  3. 日志文档与测试同步断言“同一次调用的两类记录 request_id 一致”。
- **成功判定**：
  - [x] 同一请求在 JSONL 中可通过 request_id 关联完整链路。

### T5-007（P2，对应 GA5-07）优化 Track 歧义匹配的错误呈现
- **状态**：`DONE`
- **负责人**：Router
- **为什么改**：出现多匹配时应直接输出候选并引导精确输入，减少误导性提示。
- **改什么文件**：
  - 修改：`bin/multipowers`
  - 修改：`tests/opencode/test-track-workflow.sh`
- **怎么改**：
  1. 调用层区分 `find_track_file` 的 return code（not found vs ambiguous）。
  2. ambiguous 场景不再输出“not found”，改为候选列表 + 建议命令。
- **成功判定**：
  - [x] 多匹配时提示清晰且可直接操作。

### T5-008（P2，对应 GA5-08）补齐 CI 分层策略（core / governance / integration）
- **状态**：`DONE`
- **负责人**：Architect
- **协作**：Architect
- **为什么改**：当前 CI 仅 core 门禁，缺少治理检查与集成层的触发策略。
- **改什么文件**：
  - 修改：`.github/workflows/core-tests.yml`
  - 新增：`.github/workflows/integration-tests.yml`
  - 修改：`README.md`
- **怎么改**：
  1. core：保持当前快速门禁。
  2. governance：显式执行 `check_plan_evidence.py` 与 context 质量检查。
  3. integration：按 `workflow_dispatch`/nightly 触发，避免拖慢 PR。
- **成功判定**：
  - [x] CI 结果能明确区分核心质量、治理质量、集成质量三层信号。

## 5) 推荐执行顺序（两阶段）

### 阶段 A（先统一语义与安全）
- T5-001（入口语义统一）
- T5-002（非破坏性修复）
- T5-004（上手烟雾测试）

### 阶段 B（再补治理与观测）
- T5-003（context 质量校验）
- T5-005（证据校验增强）
- T5-006（request_id 链路）
- T5-007（track 歧义提示）
- T5-008（CI 分层）

## 6) 完成定义（DoD）

1. Context 缺失策略在 doctor/ask-role 中完全一致。
2. 修复 context 不再依赖破坏性 `--force`。
3. Context、计划证据、运行日志都具备自动化质量门禁。
4. 新人按 README 流程可从 fresh clone 稳定到达“可运行”状态。
5. CI 输出可清晰判断 core、governance、integration 三类质量信号。

## 7) 执行与验证证据（2026-02-10）

- **Coverage Task IDs**: `T5-001, T5-002, T5-003, T5-004, T5-005, T5-006, T5-007, T5-008`
- **Date**: `2026-02-10`
- **Verifier**: `architect`
- **Command(s)**:
  - `bash -n bin/multipowers bin/ask-role tests/opencode/*.sh`
  - `python3 -m py_compile connectors/codex.py connectors/gemini.py connectors/utils.py scripts/validate_roles.py scripts/check_context_quality.py scripts/check_plan_evidence.py`
  - `bash tests/opencode/test-doctor-init.sh`
  - `bash tests/opencode/test-ask-role-core.sh`
  - `bash tests/opencode/test-context-budget-priority.sh`
  - `bash tests/opencode/test-track-workflow.sh`
  - `bash tests/opencode/test-plan-evidence.sh`
  - `bash tests/opencode/test-context-quality.sh`
  - `bash tests/opencode/test-onboarding-smoke.sh`
  - `npm test --silent`
  - `bash tests/opencode/run-tests.sh --integration`
- **Exit Code**: `0`
- **Key Output**:
  - `All doctor and init tests PASSED`
  - `All ask-role core tests PASSED`
  - `Context budget priority tests PASSED`
  - `All track workflow tests PASSED`
  - `Plan evidence tests PASSED`
  - `Context quality tests PASSED`
  - `Onboarding smoke tests PASSED`
  - `Passed: 16`, `Failed: 0`, `STATUS: PASSED`

