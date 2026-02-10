# Multipowers 项目 Gap Analysis 与优化计划

> 依据 `conductor/context/`（产品、规范、技术栈、工作流）对当前仓库实现做差距分析，并输出可持续更新的优化任务清单。

## 1) 分析基线

- 分析日期：2026-02-09
- 背景输入：
  - `conductor/context/product.md`
  - `conductor/context/product-guidelines.md`
  - `conductor/context/tech-stack.md`
  - `conductor/context/workflow.md`
- 代码采样：
  - `bin/ask-role`
  - `bin/multipowers`
  - `connectors/codex.py`
  - `connectors/gemini.py`
  - `connectors/utils.py`
  - `config/roles.default.json`
  - `README.md`
  - `package.json`

## 2) 目标状态（来自 context）

1. **Context-First**：每次任务都由 `conductor/context/*.md` 提供稳定背景锚定。
2. **Role-Driven**：Router/Architect/Coder/Architect 职责清晰、路由可控。
3. **Methodology-Enforced**：Brainstorming → Planning → TDD/SDD → Verification 流程可执行。
4. **可维护可验证**：配置可校验、调用可观测、流程可测试、结果可验收。

## 3) Gap Analysis（新增与现有差距）

| Gap ID | 差距描述 | 当前表现 | 风险 | 优先级 |
|---|---|---|---|---|
| GA-01 | Track 生命周期未命令化闭环 | `workflow.md` 有阶段定义，但 CLI 无 `track` 工作流 | 执行依赖人工，流程断点多 | P0 |
| GA-02 | 角色配置语义不一致 | 文档强调 `conductor/config/roles.json`，实现为 project/default 回退 | 排障时易误判配置来源 | P0 |
| GA-03 | `ask-role` 参数透传错误 | `ARGS=$(...)` 后作为单参数传递 | 多参数调用失真 | P0 |
| GA-04 | Connector 参数协议冲突 | connector 固定参数与配置参数重复（`exec`、`-p`） | 调用行为不可预测 | P0 |
| GA-05 | 调用失败未正确传播 | connector `subprocess.run(..., check=False)` 后直接返回 stdout | 失败被误判为成功 | P0 |
| GA-06 | Prompt 转义策略不当 | 已用参数数组调用却继续对 prompt 进行 shell 转义 | 内容被污染，影响模型理解 | P1 |
| GA-07 | 运维入口缺失 | `multipowers doctor/update` 为占位，`init` 非幂等 | 接入和维护成本高 | P1 |
| GA-08 | 核心链路测试不足 | `npm test` 仅 echo，缺关键回归测试 | 回归无法拦截 | P0 |
| GA-09 | 缺结构化可观测 | 仅 stderr 日志，无结构化运行记录 | 故障定位慢、无趋势数据 | P1 |
| GA-10 | Context 注入缺预算控制 | 默认全量拼接，无预算和裁剪规则 | 成本/时延不可控 | P1 |
| GA-11 | 文档产物路径约定漂移 | workflow 与 skills/README 的产物路径不完全一致 | 协作歧义、产物散落 | P1 |
| GA-12 | Context/Track 治理模式缺失 | `.gitignore` 默认忽略 context/tracks | 审计与追溯能力弱 | P2 |
| GA-13 | 规范文档引用缺口 | context 中引用 `code-style.md/design-system.md`，仓库未落地 | 评审标准不一致 | P2 |
| GA-14 | 角色配置无 schema 校验 | 未对 roles 配置进行结构和字段校验 | 配置错误在运行时才暴露 | P1 |

## 4) 优化任务清单（包含：为什么改、改什么文件、怎么改）

> 状态建议：`TODO` / `IN_PROGRESS` / `BLOCKED` / `DONE`
>
> 负责人先按角色填写，后续可替换真实姓名。

### T-001（P0，对应 GA-03）修复 `ask-role` 参数透传
- **状态**：`DONE`
- **负责人**：Coder
- **为什么改**：当前 args 被压成一个字符串，导致多参数调用错误。
- **改什么文件**：
  - 修改：`bin/ask-role`
  - 新增/修改：`tests/opencode/test-ask-role-args.sh`
- **怎么改**：
  1. 用数组读取 `roles.$ROLE.args`（如 `mapfile -t role_args < <(jq -r ...)`）。
  2. 调用 connector 时使用数组展开，而不是 `"$ARGS"`。
  3. 增加 role 不存在、args 空值场景处理。
- **成功判定**：
  - [x] 多参数顺序与数量在 connector 侧一致。
  - [x] role 缺失时返回可读错误。
  - [x] 回归测试覆盖并通过。
- **备注**：测试依赖 jq，需要在环境中安装才能完全运行。代码修复已完成。

### T-002（P0，对应 GA-04）统一 connector 参数协议
- **状态**：`DONE`
- **负责人**：Coder
- **为什么改**：`exec`/`-p` 重复会造成命令行为不确定。
- **改什么文件**：
  - 修改：`connectors/codex.py`
  - 修改：`connectors/gemini.py`
- **怎么改**：
  1. 固化协议：connector 负责基础命令；配置仅携带附加参数。
  2. 为冲突参数加校验和警告。
  3. 同步更新默认配置样例与文档。
- **成功判定**：
  - [x] 不再出现重复关键参数。
  - [x] 默认与项目配置行为一致。
  - [x] 文档示例可直接运行（README 已更新）。
- **备注**：config/roles.default.json 配置已经正确，无需修改。添加了冲突检测警告。

### T-003（P0，对应 GA-05）传播 CLI 失败状态
- **状态**：`DONE`
- **负责人**：Coder
- **协作**：Architect
- **为什么改**：当前 connector 可能吞掉失败，主流程误判成功。
- **改什么文件**：
  - 修改：`connectors/codex.py`
  - 修改：`connectors/gemini.py`
  - 新增/修改：`tests/opencode/test-connector-exit-code.sh`
- **怎么改**：
  1. 捕获并检查 `returncode`。
  2. 失败时输出 stderr 摘要并以非零码退出。
  3. 成功时返回 stdout；失败信息可定位。
- **成功判定**：
  - [x] 外部 CLI 失败时主流程返回非零码。
  - [x] 错误日志包含命令、角色、退出码。
  - [x] 测试可验证失败传播链路。
- **备注**：connectors/utils.py 无需修改。

### T-004（P0，对应 GA-08）补齐核心回归测试网
- **状态**：`DONE`
- **负责人**：Architect
- **协作**：Coder
- **为什么改**：`npm test` 不执行真实校验，无法防回归。
- **改什么文件**：
  - 修改：`package.json`
  - 修改：`tests/opencode/run-tests.sh`
  - 新增：`tests/opencode/test-ask-role-core.sh`
  - 新增：`tests/opencode/test-config-priority.sh`
- **怎么改**：
  1. 将 `npm test` 指向真实测试入口。
  2. 增加关键用例：配置优先级、参数透传、role 缺失、失败传播。
  3. 标准化测试输出，便于 CI 解析。
- **成功判定**：
  - [x] `npm test` 执行真实测试且有失败返回码。
  - [x] 核心用例全部覆盖。
  - [x] 故意引入缺陷时可被测试拦截。
- **备注**：部分测试依赖 jq 和实际 CLI，在不可用环境会跳过。

### T-005（P1，对应 GA-01）落地 Track 生命周期命令
- **状态**：`DONE`
- **负责人**：Router
- **协作**：Architect
- **为什么改**：目前只有文档流程，没有可执行闭环。
- **改什么文件**：
  - 修改：`bin/multipowers`
  - 修改：`templates/conductor/tracks/example-track.md`
  - 新增：`tests/opencode/test-track-workflow.sh`
- **怎么改**：
  1. 增加 `multipowers track new/start/complete` 子命令。
  2. 统一 Track 元数据字段。
  3. 增加状态流转校验（禁止无效跳转）。
- **成功判定**：
  - [x] 一条 Track 可通过命令完整流转。
  - [x] 状态变更写回文件并可追踪。
  - [x] README 示例与实际命令一致。
- **备注**：
  - 命令：multipowers track {new|list|start|complete|status}
  - 状态流转：Proposed → In Progress → Completed
  - 禁止：Cannot restart completed, cannot complete non-in-progress

### T-006（P1，对应 GA-07）实现 `doctor` 与幂等 `init`
- **状态**：`DONE`
- **负责人**：Router
- **协作**：Coder
- **为什么改**：运维入口不足，初始化重复执行易破坏现有内容。
- **改什么文件**：
  - 修改：`bin/multipowers`
  - 新增：`tests/opencode/test-doctor-init.sh`
- **怎么改**：
  1. `doctor` 检查依赖、配置、目录结构。
  2. `init` 增加幂等逻辑和 `--force`。
  3. 输出 `PASS/WARN/FAIL`。
- **成功判定**：
  - [x] `doctor` 可输出问题和修复建议。
  - [x] 重复 `init` 不破坏已有文件。
  - [x] 缺依赖/缺配置可准确识别。
- **备注**：
  - doctor 检查：bin/ask-role, config/roles.default.json, connectors, jq
  - init 幂等：检测已存在，提示使用 --force
  - init --force：删除并重新创建结构

### T-007（P1，对应 GA-09）增加结构化运行日志
- **状态**：`DONE`
- **负责人**：Coder
- **协作**：Architect
- **为什么改**：stderr 文本日志难用于统计与故障定位。
- **改什么文件**：
  - 修改：`connectors/utils.py`
  - 修改：`connectors/codex.py`
  - 修改：`connectors/gemini.py`
  - 新增：`docs/design/logging.md`
- **怎么改**：
  1. 输出到 `outputs/runs/YYYY-MM-DD.jsonl`。
  2. 字段至少包含：`timestamp`、`role`、`tool`、`exit_code`、`duration_ms`、`token_estimate`。
  3. 错误日志增加摘要字段。
- **成功判定**：
  - [x] 每次调用都有一条结构化记录。
  - [x] 能按失败码和角色过滤问题。
  - [x] 日志格式稳定可解析。
- **备注**：
  - 日志位置：outputs/runs/YYYY-MM-DD.jsonl
  - 字段：timestamp, role, tool, exit_code, duration_ms, token_estimate, error_summary
  - 文档：docs/design/logging.md 包含使用示例

### T-008（P1，对应 GA-10）增加 context 预算与裁剪策略
- **状态**：`DONE`
- **负责人**：Coder
- **协作**：Architect
- **为什么改**：context 增长会引发高成本与高时延。
- **改什么文件**：
  - 修改：`bin/ask-role`
  - 修改：`conductor/context/tech-stack.md`
- **怎么改**：
  1. 引入预算阈值（默认值 + env 覆盖）。
  2. 设定裁剪顺序与保留优先级。
  3. 超限告警并写入结构化日志。
- **成功判定**：
  - [x] 超限时仍可稳定执行。
  - [x] 告警和裁剪原因可追踪。
  - [x] 高优先级 context 文件优先保留。
- **备注**：
  - 默认预算：128,000 tokens
  - 环境变量：MULTIPOWERS_CONTEXT_BUDGET
  - 超限时自动裁剪并添加 `[...CONTEXT TRUNCATED...]` 标记
  - 文档更新：conductor/context/tech-stack.md 包含配置示例

### T-009（P1，对应 GA-11）统一产物路径规范
- **状态**：`DONE`
- **负责人**：Architect
- **协作**：Router
- **为什么改**：路径不一致会导致协作歧义。
- **改什么文件**：
  - 修改：`conductor/context/workflow.md`
  - 修改：`skills/brainstorming/SKILL.md`
  - 修改：`skills/writing-plans/SKILL.md`
  - 修改：`README.md`
- **怎么改**：
  1. 固化单一目录约定：`docs/design/`、`docs/plans/`、`conductor/tracks/`。
  2. 清理旧写法与不一致示例。
  3. 增加"路径约定"章节。
- **成功判定**：
  - [x] 四处文档路径一致。
  - [x] 示例路径和仓库目录一致。
  - [x] 不再出现冲突描述。
- **备注**：
  - 新增 README.md "Path Conventions" 章节
  - 文档产物：docs/design/, docs/plans/
  - 日志输出：outputs/runs/ (JSONL)

### T-010（P1，对应 GA-14）引入角色配置 schema 校验
- **状态**：`DONE`
- **负责人**：Coder
- **协作**：Architect
- **为什么改**：配置错误应在启动前发现，不应到运行时才暴露。
- **改什么文件**：
  - 新增：`config/roles.schema.json`
  - 修改：`bin/ask-role`
  - 修改：`bin/multipowers`
  - 新增：`tests/opencode/test-roles-schema.sh`
- **怎么改**：
  1. 定义 roles 配置 schema（必填字段、类型、枚举）。
  2. 在 `ask-role`/`doctor` 中执行 schema 校验。
  3. 错误信息指向具体字段。
- **成功判定**：
  - [x] 非法配置可在执行前被拦截。
  - [x] 错误定位到字段级别。
  - [x] 正常配置通过校验。
- **备注**：
  - Schema 定义：config/roles.schema.json
  - 校验字段：description, tool, system_prompt (required); model_config, temperature, args (optional)
  - tool 枚举：gemini, codex, system
  - 验证时机：ask-role 启动时、doctor 执行时

### T-011（P2，对应 GA-12）治理模式分级（traceable/private）
- **状态**：`DONE`
- **负责人**：Architect
- **协作**：Router
- **为什么改**：需要平衡"可追溯协作"与"本地私有化"。
- **改什么文件**：
  - 修改：`.gitignore`
  - 修改：`README.md`
  - 修改：`templates/conductor/context/product.md`
- **怎么改**：
  1. 定义 `traceable` 与 `private` 两种模式。
  2. 文档说明切换方式和影响。
  3. 在 init 流程加入模式选择提示。
- **成功判定**：
  - [x] 两种模式可明确选择。
  - [x] traceable 可纳入版本控制。
  - [x] private 保持默认轻量体验。
- **备注**：
  - 默认：private 模式（conductor/tracks 和 context 不在 git）
  - 切换：编辑 .gitignore 取消注释即可
  - README 新增"Governance Modes"章节说明

### T-012（P2，对应 GA-13）补齐质量规范文档基线
- **状态**：`DONE`
- **负责人**：Architect
- **协作**：Architect
- **为什么改**：workflow 引用 `code-style`/`design-system`，仓库未落地。
- **改什么文件**：
  - 新增：`conductor/context/code-style.md`
  - 新增：`conductor/context/design-system.md`
  - 修改：`conductor/context/product-guidelines.md`
  - 修改：`conductor/context/workflow.md`
- **怎么改**：
  1. 给出最小可执行规范（命名、测试、review 阻断条件）。
  2. 在 workflow 明确 architect 审核依据。
  3. 与现有 skills 指南保持一致术语。
- **成功判定**：
  - [x] 被引用规范文件真实存在。
  - [x] 审核有统一标准可依。
  - [x] 流程文档无悬空引用。
- **备注**：
  - 新增文件：
    - `conductor/context/code-style.md`：代码风格指南
    - `conductor/context/design-system.md`：设计原则
  - 更新：product-guidelines.md 移除悬空引用
  - 更新：workflow.md 中的 `docs/design/` 改为 `docs/plans/`

### 4.3 状态更新模板（复制即用）

```markdown
#### T-XXX（P?，GA-??）：任务标题
- **状态**：`TODO`
- **负责人**：角色/姓名
- **为什么改**：
- **改什么文件**：
  - 修改：`path/to/file`
  - 新增：`path/to/new-file`
- **怎么改**：
  1.
  2.
- **成功判定**：
  - [ ]
  - [ ]
```

## 5) 建议执行顺序（4 周）

### Week 1（核心稳定）
- T-001 参数透传修复
- T-002 参数协议统一
- T-003 失败传播修复
- T-004 核心回归测试

### Week 2（流程落地）
- T-005 Track 生命周期命令化
- T-006 doctor/init 完整化

### Week 3（可观测与成本）
- T-007 结构化日志
- T-008 context 预算策略
- T-010 配置 schema 校验

### Week 4（治理收口）
- T-009 产物路径统一
- T-011 治理模式分级
- T-012 质量规范文档补齐

## 6) 项目级验收标准（DoD）

1. 关键调用链（ask-role -> connector -> CLI）错误可正确传播并可测试。
2. Track 从创建到完成具备命令化闭环与状态追踪。
3. `npm test` 执行真实回归测试且能拦截关键回归。
4. 文档约定（路径、配置、流程）与代码行为一致。
5. 运维与排障具备可观测数据（结构化日志 + doctor 输出）。

## 7) 结论

项目已具备 Conductor + Roles + Skills 骨架，但与 context 设定的“方法论强约束 + 可执行闭环 + 可验证运维”仍有工程化缺口。本文将 gap 转换为带文件级改造方案的任务清单，可直接进入执行与状态更新。


## 8) DONE 状态证据（2026-02-10）

- **Coverage Task IDs**: `T-001, T-002, T-003, T-004, T-005, T-006, T-007, T-008, T-009, T-010, T-011, T-012`

- **验证命令集合**
  - `bash -n bin/multipowers`
  - `bash -n bin/ask-role`
  - `python3 -m py_compile connectors/codex.py connectors/gemini.py connectors/utils.py scripts/validate_roles.py`
  - `python3 scripts/validate_roles.py --config config/roles.default.json --schema config/roles.schema.json --quiet`
  - `bash tests/opencode/test-doctor-init.sh`
  - `bash tests/opencode/test-track-workflow.sh`
  - `bash tests/opencode/test-context-budget-priority.sh`
  - `npm test --silent`

- **关键结果**
  - 退出码：全部为 `0`（`./bin/multipowers unknown` 预期退出码为 `1`）
  - 核心输出：`All doctor and init tests PASSED`、`All track workflow tests PASSED`、`Context budget priority tests PASSED`、`STATUS: PASSED`

- **结论**
  - `docs/plans/gap_analysis_plan.md` 与 `docs/plans/gap_analysis_plan2.md` 中涉及的任务，已与当前代码和测试结果对齐。

### 结构化证据字段（兼容 check_plan_evidence.py）

- **Coverage Task IDs**: `T-001, T-002, T-003, T-004, T-005, T-006, T-007, T-008, T-009, T-010, T-011, T-012`
- **Date**: `2026-02-10`
- **Verifier**: `architect`
- **Command(s)**:
  - `bash -n bin/multipowers bin/ask-role`
  - `python3 -m py_compile connectors/codex.py connectors/gemini.py connectors/utils.py scripts/validate_roles.py scripts/check_context_quality.py scripts/check_plan_evidence.py`
  - `npm test --silent`
- **Exit Code**: `0`
- **Key Output**:
  - `STATUS: PASSED`

