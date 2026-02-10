# Multipowers 项目 Gap Analysis（第三版）与优化计划

> 基于 `conductor/context/` 的产品背景与 2026-02-09 当前代码现状，聚焦“仍未闭合”的 gap，并给出任务级落地方案：**为什么改、改什么文件、怎么改**。

## 1) 分析基线

- 背景来源：
  - `conductor/context/product.md`
  - `conductor/context/product-guidelines.md`
  - `conductor/context/workflow.md`
  - `conductor/context/tech-stack.md`
- 现状核查命令（2026-02-09）：
  - `bash -n bin/multipowers` → 通过
  - `bash -n bin/ask-role` → 通过
  - `npm test --silent` → 通过（核心测试 10/10 通过）

## 2) 目标状态（context 约束）

1. `Context-First`：context 注入稳定、可控、有预算与可观测。
2. `Role-Driven`：角色配置一致、路由明确、参数与失败语义准确。
3. `Methodology Enforcement`：Brainstorming → Planning → TDD/SDD → Verification 可执行。
4. `Verification-First`：完成状态必须有可复现验证证据。

## 3) Gap Analysis（当前仍待解决）

| Gap ID | 差距描述 | 当前表现 | 风险 | 优先级 |
|---|---|---|---|---|
| GA3-01 | 主 CLI 脚本不可执行 | `bin/multipowers` 存在 shell 语法错误 | 核心命令不可用 | P0 |
| GA3-02 | 错误退出码路径损坏 | 存在 `exit1` 拼写错误 | 异常路径不可预测 | P0 |
| GA3-03 | Track 状态字段格式不一致 | 模板 `**Status:**` 与命令 `*Status:` 混用 | 状态机不稳定 | P0 |
| GA3-04 | schema 与实现校验脱节 | `ask-role` 读取 schema 但使用硬编码校验 | 规则漂移、假通过 | P1 |
| GA3-05 | 校验输出污染模型输出 | `Schema validation passed` 打到 stdout | 污染下游响应 | P1 |
| GA3-06 | prompt 过度转义 | connectors 对参数数组调用仍做 shell 转义 | prompt 语义被扭曲 | P1 |
| GA3-07 | 失败日志泄露风险 | connector 失败打印完整命令（含 prompt） | 敏感上下文泄露 | P1 |
| GA3-08 | context 裁剪策略与文档不一致 | 文档写“按文件优先级”，实现为字符截断 | 关键 context 可能被误裁 | P1 |
| GA3-09 | 默认测试集覆盖不完整 | 核心新增测试未纳入默认 test 集 | 回归漏检 | P0 |
| GA3-10 | 测试环境路径假设过强 | `tests/opencode/setup.sh` 依赖不存在的 `.opencode/plugins/superpowers.js` | 测试高失败率 | P0 |
| GA3-11 | doctor 未优先检查生效配置 | 仅围绕默认配置检查 | 项目级配置问题漏检 | P1 |
| GA3-12 | 文档与真实状态漂移 | 部分任务标记 DONE，但关键命令仍失败 | 计划可信度下降 | P2 |
| GA3-13 | 文档产物路径仍有漂移 | `workflow.md` 与 `design-system.md` 对设计产物目录约定不一致 | 团队执行歧义 | P1 |
| GA3-14 | 配置 schema 本身与默认配置冲突 | schema `additionalProperties: false`，但默认配置含 `model_alias` | schema 不可用或被绕过 | P1 |

## 4) 优化任务清单（为什么改 / 改什么文件 / 怎么改）

> 状态枚举建议：`TODO` / `IN_PROGRESS` / `BLOCKED` / `DONE`

### T3-001（P0，对应 GA3-01）修复 `bin/multipowers` 语法错误
- **状态**：`DONE`
- **为什么改**：CLI 主入口不可执行，阻断所有流程。
- **改什么文件**：
  - 修改：`bin/multipowers`
- **怎么改**：
  1. 修正 `doctor` 分支内 `if python3 -c ...; then ... else ... fi` 嵌套结构。
  2. 去掉 `then` 分支里无意义的 `if [ $? -ne 0 ]`。
  3. 修正 `for track in conductor/tracks/*.md 2>/dev/null` 的非法写法。
- **成功判定**：
  - [x] `bash -n bin/multipowers` 返回 0。
  - [x] `./bin/multipowers doctor` 可执行。

### T3-002（P0，对应 GA3-02）修复 CLI 异常退出码路径
- **状态**：`DONE`
- **为什么改**：错误分支必须稳定返回非零码，便于脚本集成。
- **改什么文件**：
  - 修改：`bin/multipowers`
- **怎么改**：
  1. 将 `exit1` 改为 `exit 1`。
  2. 保证 unknown subcommand 路径统一输出 usage + 非零退出。
- **成功判定**：
  - [x] `./bin/multipowers unknown` 返回码为 1。

### T3-003（P0，对应 GA3-03）统一 Track 元数据格式
- **状态**：`DONE`
- **为什么改**：状态字段格式不一致导致解析与写回不稳定。
- **改什么文件**：
  - 修改：`templates/conductor/tracks/example-track.md`
  - 修改：`bin/multipowers`
  - 修改：`tests/opencode/test-track-workflow.sh`
  - 修改：`README.md`
- **怎么改**：
  1. 固定一种元数据格式（建议 `**Status:**` / `**Updated At:**` / `**Owner:**`）。
  2. `track new/start/complete/status` 统一按该格式读写。
  3. 测试覆盖完整状态流转。
- **成功判定**：
  - [x] Track 生命周期测试全绿。
  - [x] 无格式混用导致的状态误判。

### T3-004（P1，对应 GA3-04/GA3-14）统一配置校验机制（基于 schema）
- **状态**：`DONE`
- **为什么改**：当前为硬编码校验，且与 schema/默认配置冲突。
- **改什么文件**：
  - 修改：`config/roles.schema.json`
  - 修改：`bin/ask-role`
  - 修改：`bin/multipowers`
  - 新增：`scripts/validate_roles.py`（建议）
- **怎么改**：
  1. 先对齐 schema 与 `config/roles.default.json`（包含 `model_alias` 兼容策略）。
  2. 抽离统一校验器，`ask-role` 和 `doctor` 共用。
  3. 错误信息输出字段路径与修复建议。
- **成功判定**：
  - [x] 非法配置可在执行前被阻断。
  - [x] schema 与默认配置不冲突。

### T3-005（P1，对应 GA3-05）清理校验日志通道
- **状态**：`DONE`
- **为什么改**：校验提示进入 stdout 会污染模型输出。
- **改什么文件**：
  - 修改：`bin/ask-role`
  - 修改：`scripts/validate_roles.py`（若新增）
- **怎么改**：
  1. 所有诊断与校验日志统一写 stderr。
  2. stdout 仅输出角色调用结果。
- **成功判定**：
  - [x] `ask-role` stdout 不包含“schema validation passed”等运维文本。

### T3-006（P1，对应 GA3-06）移除不必要 prompt 转义
- **状态**：`DONE`
- **为什么改**：参数数组调用不需要 shell escaping，现有转义会扭曲 prompt。
- **改什么文件**：
  - 修改：`connectors/codex.py`
  - 修改：`connectors/gemini.py`
  - 修改：`connectors/utils.py`
  - 新增：`tests/opencode/test-prompt-preserve.sh`
- **怎么改**：
  1. 直接传原始 prompt 给 subprocess 参数数组。
  2. 保留安全边界在命令数组构造，不在文本内容层强行 escape。
  3. 增加特殊字符回归测试。
- **成功判定**：
  - [x] 特殊字符 prompt 语义保持不变。

### T3-007（P1，对应 GA3-07）收敛错误日志敏感信息
- **状态**：`DONE`
- **为什么改**：完整命令日志可能泄露 prompt/context。
- **改什么文件**：
  - 修改：`connectors/codex.py`
  - 修改：`connectors/gemini.py`
  - 修改：`docs/design/logging.md`
- **怎么改**：
  1. 错误日志改为“命令骨架 + 参数数量 + exit code”。
  2. 结构化日志只保留摘要，不写原始 prompt。
- **成功判定**：
  - [x] 失败日志中不出现完整 prompt 内容。

### T3-008（P1，对应 GA3-08）按文件优先级实现 context 裁剪
- **状态**：`DONE`
- **为什么改**：当前字符截断与文档承诺不一致。
- **改什么文件**：
  - 修改：`bin/ask-role`
  - 修改：`connectors/utils.py`
  - 修改：`conductor/context/tech-stack.md`
  - 新增：`tests/opencode/test-context-budget-priority.sh`
- **怎么改**：
  1. 按文件优先级逐个累加 context。
  2. 超限时从低优先级文件开始裁剪或摘要。
  3. 记录被裁剪文件清单到结构化日志。
- **成功判定**：
  - [x] 超限时核心 context 文件优先保留。
  - [x] 裁剪决策可追踪。

### T3-009（P0，对应 GA3-09）扩大默认测试入口覆盖
- **状态**：`DONE`
- **为什么改**：关键回归测试未默认运行，无法形成有效护栏。
- **改什么文件**：
  - 修改：`tests/opencode/run-tests.sh`
  - 修改：`package.json`
- **怎么改**：
  1. 把关键测试加入默认集（args/exit-code/schema/track）。
  2. 保留 integration 测试为可选。
- **成功判定**：
  - [x] `npm test` 覆盖核心链路测试并返回准确状态码。

### T3-010（P0，对应 GA3-10）修复测试环境资源依赖
- **状态**：`DONE`
- **为什么改**：测试依赖不存在文件导致大量假失败。
- **改什么文件**：
  - 修改：`tests/opencode/setup.sh`
  - 修改：`tests/opencode/test-plugin-loading.sh`
  - 修改：`tests/opencode/test-skills-core.sh`
- **怎么改**：
  1. 以仓库真实存在资源为基准构建测试环境。
  2. 将缺失可选资源改为 skip，不阻断核心测试。
- **成功判定**：
  - [x] 干净环境下 `npm test` 不因路径不存在而直接失败。

### T3-011（P1，对应 GA3-11）`doctor` 对齐真实配置优先级
- **状态**：`DONE`
- **为什么改**：运行优先项目配置，doctor 不对齐会漏检。
- **改什么文件**：
  - 修改：`bin/multipowers`
  - 修改：`README.md`
- **怎么改**：
  1. doctor 先检测 `conductor/config/roles.json`，再回退默认配置。
  2. 报告中打印当前生效配置路径。
- **成功判定**：
  - [x] doctor 输出与 `ask-role` 实际生效配置一致。

### T3-012（P1，对应 GA3-13）统一设计产物路径约定
- **状态**：`DONE`
- **为什么改**：`workflow.md` 与设计系统文档产物目录不一致。
- **改什么文件**：
  - 修改：`conductor/context/workflow.md`
  - 修改：`conductor/context/design-system.md`
  - 修改：`README.md`
  - 修改：`skills/brainstorming/SKILL.md`
- **怎么改**：
  1. 固化唯一设计产物目录（建议 `docs/design/`）。
  2. 清理所有文档中冲突写法。
- **成功判定**：
  - [x] 文档中的设计产物路径 100% 一致。

### T3-013（P2，对应 GA3-12）建立状态变更证据模板
- **状态**：`DONE`
- **为什么改**：避免“文档 DONE 但代码未可运行”的状态漂移。
- **改什么文件**：
  - 新增：`docs/plans/status-evidence-template.md`
  - 修改：`docs/plans/gap_analysis_plan.md`
  - 修改：`docs/plans/gap_analysis_plan2.md`
- **怎么改**：
  1. 定义状态更新必填项：命令、退出码、关键输出、日期。
  2. 任何 `DONE` 必须附证据块。
- **成功判定**：
  - [x] 所有 `DONE` 条目都附可复现证据。

## 5) 推荐执行顺序（4 周）

### Week 1（P0 止血）
- T3-001, T3-002, T3-003

### Week 2（测试可靠性）
- T3-009, T3-010, T3-011

### Week 3（调用链质量）
- T3-004, T3-005, T3-006, T3-007

### Week 4（一致性治理）
- T3-008, T3-012, T3-013

## 6) 完成定义（DoD）

1. `bin/multipowers` 与 `bin/ask-role` 均通过语法检查并可执行。
2. `npm test` 在标准环境可稳定运行，核心链路测试全绿。
3. 配置校验与 schema 一致，且不污染 stdout。
4. context 裁剪策略与文档一致，可观测且可解释。
5. 计划状态更新都有验证证据，状态可信。



## 7) 执行与验证证据（2026-02-10）

- **E3-001（T3-001/T3-002）CLI 可执行与退出码**
  - 命令：`bash -n bin/multipowers && ./bin/multipowers unknown`
  - 退出码：`0`（语法检查），`1`（unknown 子命令）
  - 关键输出：`Usage: multipowers {init|doctor|update|track}`

- **E3-002（T3-003/T3-011）Track 与 doctor 行为**
  - 命令：`bash tests/opencode/test-track-workflow.sh`、`./bin/multipowers doctor`
  - 退出码：`0`
  - 关键输出：`All track workflow tests PASSED`；`effective roles config: config/roles.default.json`

- **E3-003（T3-004/T3-005）Schema 校验与输出通道**
  - 命令：`python3 scripts/validate_roles.py --config config/roles.default.json --schema config/roles.schema.json --quiet`
  - 退出码：`0`
  - 关键输出：成功场景无 stdout 污染，错误信息输出到 stderr

- **E3-004（T3-006/T3-007）prompt 保真与日志脱敏**
  - 命令：`bash tests/opencode/test-prompt-preserve.sh`、`bash tests/opencode/test-connector-exit-code.sh`
  - 退出码：`0`
  - 关键输出：`Prompt preservation tests PASSED`；`Prompt masked in error output`

- **E3-005（T3-008）context 预算优先级裁剪**
  - 命令：`bash tests/opencode/test-context-budget-priority.sh`
  - 退出码：`0`
  - 关键输出：`High-priority files are preserved`；`Truncation warning includes trimmed files`

- **E3-006（T3-009/T3-010/T3-012/T3-013）默认回归入口与文档一致性**
  - 命令：`npm test --silent`
  - 退出码：`0`
  - 关键输出：`Passed: 10`、`Failed: 0`、`STATUS: PASSED`

### 结构化证据字段（兼容 check_plan_evidence.py）

- **Coverage Task IDs**: `T3-001, T3-002, T3-003, T3-004, T3-005, T3-006, T3-007, T3-008, T3-009, T3-010, T3-011, T3-012, T3-013`
- **Date**: `2026-02-10`
- **Verifier**: `architect`
- **Command(s)**:
  - `bash -n bin/multipowers bin/ask-role`
  - `python3 -m py_compile connectors/codex.py connectors/gemini.py connectors/utils.py scripts/validate_roles.py scripts/check_context_quality.py scripts/check_plan_evidence.py`
  - `npm test --silent`
- **Exit Code**: `0`
- **Key Output**:
  - `STATUS: PASSED`

