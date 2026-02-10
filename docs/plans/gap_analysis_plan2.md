# Multipowers 项目 Gap Analysis（第二版）与优化计划

> 目标：基于 `conductor/context/` 的项目背景，聚焦“当前仍存在的 gap”，并给出可执行的优化任务（为什么改、改什么文件、怎么改）。

## 1) 分析基线

- 分析日期：2026-02-09
- 背景依据：
  - `conductor/context/product.md`
  - `conductor/context/product-guidelines.md`
  - `conductor/context/workflow.md`
  - `conductor/context/tech-stack.md`
- 现状抽样：
  - `bin/ask-role`
  - `bin/multipowers`
  - `connectors/codex.py`
  - `connectors/gemini.py`
  - `connectors/utils.py`
  - `config/roles.default.json`
  - `config/roles.schema.json`
  - `tests/opencode/*`
  - `package.json`

### 1.1 快速验证结果（当前状态）

- `bash -n bin/multipowers`：通过
- `bash -n bin/ask-role`：通过
- `npm test --silent`：通过（核心测试 10/10 通过）

> 结论：核心流程稳定性与测试可信度已达标，后续以持续回归与证据化更新为主。

## 2) 目标状态（来自 context）

1. **Context-First**：上下文注入稳定且可控（预算、裁剪、可追踪）。
2. **Role-Driven**：角色配置一致、路由明确、失败可传播。
3. **Methodology Enforcement**：流程必须可执行，不是仅文档化。
4. **Verification-First**：完成状态必须由可运行验证支撑。

## 3) Gap Analysis（仍待解决）

| Gap ID | 差距 | 当前表现 | 风险 | 优先级 |
|---|---|---|---|---|
| GA2-01 | `multipowers` 脚本不可执行 | `bin/multipowers` 存在语法错误（doctor 分支 if/fi 结构问题） | CLI 主入口不可用 | P0 |
| GA2-02 | CLI 默认分支存在运行错误 | 末尾存在 `exit1`（应为 `exit 1`） | 异常路径返回码失效 | P0 |
| GA2-03 | Track 解析格式不稳 | 模板使用 `**Status:**`，命令逻辑依赖 `*Status:` | 状态读取/变更不一致 | P0 |
| GA2-04 | `ask-role` schema 校验实现与 schema 脱钩 | 读取 schema 文件但未按 schema 通用校验，内联脚本硬编码规则 | 规则漂移、维护成本高 | P1 |
| GA2-05 | schema 校验输出污染主输出 | `ask-role` 会输出 `Schema validation passed` 到标准输出 | 污染模型返回内容 | P1 |
| GA2-06 | Prompt 仍被过度转义 | connectors 使用 `sanitize_prompt()`，但调用是参数数组非 shell 拼接 | prompt 内容被污染，影响模型理解 | P1 |
| GA2-07 | 错误日志可能泄露完整 prompt | connector 失败日志打印完整命令（含 prompt） | 泄露敏感上下文 | P1 |
| GA2-08 | Context 预算实现与文档声明不一致 | 文档声明“按文件优先级裁剪”，实现为整段字符截断 | 裁剪不可控，关键信息可能丢失 | P1 |
| GA2-09 | 测试入口覆盖不足 | `run-tests.sh` 默认仅跑少量用例，新增关键测试未纳入默认集 | 回归未被及时发现 | P0 |
| GA2-10 | 测试环境假设过强 | `tests/opencode/setup.sh` 依赖 `.opencode/plugins/superpowers.js`，仓库不存在 | CI/本地测试高失败率 | P0 |
| GA2-11 | `doctor` 只校验默认配置 | 未优先校验 `conductor/config/roles.json`（项目配置） | 真实项目配置问题被漏检 | P1 |
| GA2-12 | 状态文档与真实代码可能漂移 | 任务状态标记“DONE”但关键命令仍失败 | 决策依据失真 | P2 |

## 4) 优化任务清单（为什么改 / 改什么文件 / 怎么改）

> 状态建议：`TODO` / `IN_PROGRESS` / `BLOCKED` / `DONE`

### T2-001（P0，对应 GA2-01）修复 `bin/multipowers` 语法错误
- **状态**：`DONE`
- **为什么改**：主入口脚本语法失败会阻断全部命令。
- **改什么文件**：
  - 修改：`bin/multipowers`
- **怎么改**：
  1. 修正 `doctor` 分支中的 `if python3 -c ...; then ... else ... fi` 嵌套。
  2. 删除无效的 `if [ $? -ne 0 ]`（位于 then 分支内逻辑错误）。
  3. 修正 `for track in conductor/tracks/*.md 2>/dev/null` 的语法（改为先判空或 `nullglob`）。
- **成功判定**：
  - [x] `bash -n bin/multipowers` 返回 0。
  - [x] `./bin/multipowers doctor` 可运行并返回预期码。

### T2-002（P0，对应 GA2-02）修复默认分支退出码错误
- **状态**：`DONE`
- **为什么改**：错误分支应稳定返回非零码。
- **改什么文件**：
  - 修改：`bin/multipowers`
- **怎么改**：
  1. 将 `exit1` 修正为 `exit 1`。
  2. 为未知子命令补充统一错误提示。
- **成功判定**：
  - [x] `./bin/multipowers unknown` 输出 usage 且返回码为 1。

### T2-003（P0，对应 GA2-03）统一 Track 元数据格式与解析
- **状态**：`DONE`
- **为什么改**：状态字段格式前后不一致会导致状态机逻辑不稳。
- **改什么文件**：
  - 修改：`templates/conductor/tracks/example-track.md`
  - 修改：`bin/multipowers`
  - 修改：`README.md`
  - 修改：`tests/opencode/test-track-workflow.sh`
- **怎么改**：
  1. 统一采用一种格式（推荐 `**Status:**` / `**Updated At:**` / `**Owner:**`）。
  2. `track new/start/complete/status` 统一按该格式读写。
  3. 回归测试覆盖新建→开始→完成→状态查询全链路。
- **成功判定**：
  - [x] 不同命令对同一 Track 的状态读取一致。
  - [x] `test-track-workflow.sh` 全通过。

### T2-004（P1，对应 GA2-04/GA2-05）重构 `ask-role` 配置校验
- **状态**：`DONE`
- **为什么改**：当前校验逻辑与 schema 文件耦合弱，且污染标准输出。
- **改什么文件**：
  - 修改：`bin/ask-role`
  - 新增：`scripts/validate_roles.py`（建议）
  - 修改：`bin/multipowers`（doctor 复用同一校验器）
- **怎么改**：
  1. 提取统一校验脚本（输入 config path + role）。
  2. 所有校验输出写入 stderr，stdout 仅保留模型结果。
  3. 错误信息包含字段路径与修复建议。
- **成功判定**：
  - [x] `ask-role` stdout 不再出现校验提示。
  - [x] 非法配置可在执行前被阻断并指明字段。

### T2-005（P1，对应 GA2-06）去除不必要 prompt 转义
- **状态**：`DONE`
- **为什么改**：`subprocess.run(cmd_list)` 不需要 shell 转义，转义会扭曲 prompt。
- **改什么文件**：
  - 修改：`connectors/codex.py`
  - 修改：`connectors/gemini.py`
  - 修改：`connectors/utils.py`
- **怎么改**：
  1. connector 直接传原始 prompt 参数。
  2. 保留安全边界在“命令数组构造”，不在内容上做 shell 风格 escape。
  3. 为包含 `$`, `"`, `(`, `)` 的 prompt 增加回归测试。
- **成功判定**：
  - [x] 特殊字符 prompt 在目标 CLI 中语义不变。
  - [x] 不再出现反斜杠污染。

### T2-006（P1，对应 GA2-07）收敛失败日志的敏感暴露
- **状态**：`DONE`
- **为什么改**：失败日志打印完整命令会泄露 prompt/context。
- **改什么文件**：
  - 修改：`connectors/codex.py`
  - 修改：`connectors/gemini.py`
  - 修改：`docs/design/logging.md`
- **怎么改**：
  1. 失败日志仅打印命令骨架与参数数量，不打印完整 prompt。
  2. 结构化日志增加 `error_class`/`error_summary`，避免原文泄露。
- **成功判定**：
  - [x] 错误日志不含完整 prompt。
  - [x] 排障信息仍可定位问题。

### T2-007（P1，对应 GA2-08）实现“按文件优先级”的 context 裁剪
- **状态**：`DONE`
- **为什么改**：当前截断是纯字符级，和文档声明不一致。
- **改什么文件**：
  - 修改：`bin/ask-role`
  - 修改：`connectors/utils.py`
  - 修改：`conductor/context/tech-stack.md`
  - 新增：`tests/opencode/test-context-budget-priority.sh`
- **怎么改**：
  1. 按文件优先级加载并逐步累计 token。
  2. 超限时从低优先级文件开始裁剪/摘要。
  3. 在日志中记录被裁剪文件列表。
- **成功判定**：
  - [x] 超限时高优先级文件始终保留。
  - [x] 日志可见裁剪决策。

### T2-008（P0，对应 GA2-09）修复默认测试入口覆盖
- **状态**：`DONE`
- **为什么改**：关键测试未进入默认集，回归可能漏检。
- **改什么文件**：
  - 修改：`tests/opencode/run-tests.sh`
  - 修改：`package.json`
- **怎么改**：
  1. 将 `test-ask-role-args.sh`、`test-connector-exit-code.sh`、`test-roles-schema.sh`、`test-track-workflow.sh` 加入默认集。
  2. 区分“必跑核心集”和“可选集成集”。
- **成功判定**：
  - [x] `npm test` 覆盖核心链路关键测试。
  - [x] 默认失败能及时发现关键回归。

### T2-009（P0，对应 GA2-10）修复测试环境资源依赖
- **状态**：`DONE`
- **为什么改**：当前测试依赖不存在的插件文件路径，导致稳定失败。
- **改什么文件**：
  - 修改：`tests/opencode/setup.sh`
  - 修改：`tests/opencode/test-plugin-loading.sh`
  - 修改：`tests/opencode/test-skills-core.sh`
- **怎么改**：
  1. 统一以仓库现有路径作为源（如 `lib/`、`skills/`、`hooks/`、`commands/`）。
  2. 若某资产是可选项，测试应降级为 skip 而非 fail。
  3. 明确“最低可运行测试环境”说明。
- **成功判定**：
  - [x] 在干净环境执行 `npm test` 不因路径缺失直接失败。
  - [x] 失败仅由真实逻辑问题导致。

### T2-010（P1，对应 GA2-11）`doctor` 优先校验项目配置
- **状态**：`DONE`
- **为什么改**：真实执行优先使用项目级 roles，doctor 只查默认配置会漏检。
- **改什么文件**：
  - 修改：`bin/multipowers`
  - 修改：`README.md`
- **怎么改**：
  1. `doctor` 按 `conductor/config/roles.json` > `config/roles.default.json` 顺序检查。
  2. 输出“当前生效配置路径”。
- **成功判定**：
  - [x] doctor 报告能准确指出正在生效的配置文件。

### T2-011（P2，对应 GA2-12）建立“状态可信度”规则
- **状态**：`DONE`
- **为什么改**：计划状态与代码真实可运行性可能脱节。
- **改什么文件**：
  - 修改：`docs/plans/gap_analysis_plan.md`
  - 新增：`docs/plans/status-evidence-template.md`
- **怎么改**：
  1. 为任务 `DONE` 增加必须证据：命令、退出码、关键输出。
  2. 每次状态变更附最小验证记录。
- **成功判定**：
  - [x] 文档中的 `DONE` 均有可追溯验证证据。

## 5) 建议执行顺序（4 周）

### Week 1（P0 稳定性止血）
- T2-001, T2-002, T2-003

### Week 2（测试与可靠性）
- T2-008, T2-009, T2-010

### Week 3（执行链质量）
- T2-004, T2-005, T2-006

### Week 4（一致性与治理）
- T2-007, T2-011

## 6) 完成定义（DoD）

1. `bin/multipowers` / `bin/ask-role` 均通过语法检查并可执行。
2. `npm test` 可在标准环境稳定运行，核心链路测试全绿。
3. 角色配置校验统一、输出不污染模型结果。
4. context 预算裁剪可解释、可观测、与文档一致。
5. 任务状态更新有验证证据，不再出现“文档完成但代码不可运行”。



## 7) 状态验证证据（2026-02-10）

- **E2-001（T2-001/T2-002）CLI 基础健康**
  - 命令：`bash -n bin/multipowers`、`./bin/multipowers unknown`
  - 退出码：`0` / `1`
  - 关键输出：`Usage: multipowers {init|doctor|update|track}`

- **E2-002（T2-003/T2-010）Track 与 doctor 优先级**
  - 命令：`bash tests/opencode/test-track-workflow.sh`、`./bin/multipowers doctor`
  - 退出码：`0`
  - 关键输出：`All track workflow tests PASSED`；`effective roles config: conductor/config/roles.json`（有项目配置时）

- **E2-003（T2-004/T2-005/T2-006/T2-007）调用链与裁剪策略**
  - 命令：`bash tests/opencode/test-ask-role-args.sh`、`bash tests/opencode/test-prompt-preserve.sh`、`bash tests/opencode/test-connector-exit-code.sh`、`bash tests/opencode/test-context-budget-priority.sh`
  - 退出码：`0`
  - 关键输出：参数透传正确、prompt 保真、日志不泄露完整 prompt、裁剪决策可追踪

- **E2-004（T2-008/T2-009/T2-011）测试入口与状态可信度**
  - 命令：`npm test --silent`
  - 退出码：`0`
  - 关键输出：`Passed: 10`、`Failed: 0`、`STATUS: PASSED`

### 结构化证据字段（兼容 check_plan_evidence.py）

- **Coverage Task IDs**: `T2-001, T2-002, T2-003, T2-004, T2-005, T2-006, T2-007, T2-008, T2-009, T2-010, T2-011`
- **Date**: `2026-02-10`
- **Verifier**: `oracle`
- **Command(s)**:
  - `bash -n bin/multipowers bin/ask-role`
  - `python3 -m py_compile connectors/codex.py connectors/gemini.py connectors/utils.py scripts/validate_roles.py scripts/check_context_quality.py scripts/check_plan_evidence.py`
  - `npm test --silent`
- **Exit Code**: `0`
- **Key Output**:
  - `STATUS: PASSED`

