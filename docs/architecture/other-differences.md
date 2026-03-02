# main vs go 其他类型文件差异（非 commands/skills/scripts）

日期：2026-03-02  
比较分支：`main` vs `go`  
基线提交：`main=f6a815a326ec`，`go=5484dd8`

## 范围与口径

本文件覆盖所有 **非 command/skill/script** 文件：
- 排除：`*.sh`、`.claude/commands/*.md`、`.claude/skills/*.md`、`.claude-plugin/.claude/commands/*.md`、`.claude-plugin/.claude/skills/*.md`
- 纳入：`go/ts/js/json/toml/yaml/md` 等其他类型文件
- 本文件不替代 `commands_skills_difference.md` 与 `script-differences.md`，只补充其余类型文件差异。

状态定义：
- `equivalent`：语义等价，存在明确文件映射。
- `partial`：存在语义承接但非逐文件等价。
- `missing`：main 文件在 go 侧无等价承接。
- `intentional-diff`：有意的分支差异（新增或不迁移）。

规则：先语义一对一，再补文件一对一；所有 `partial/missing` 给出整改。

迁移决策说明（产品约束优先）：
- 不要求把 `main` 的每个文件机械迁移到 `go`。
- `partial/missing` 必须有显式 `decision`：`MIGRATE_TO_GO`、`COPY_FROM_MAIN`、`EXCLUDE_WITH_REASON`、`DEFER_WITH_CONDITION`。
- 决策依据以 `.multipowers/product-guidelines.md` 与 `.multipowers/product.md` 为准（尤其 invocation/hook/contract 核心路径）。
- 映射口径保持 `source file -> target file/domain -> target symbol/contract`。

## Evidence Legend

- `E0`：doc-only plan（仅文档规划）
- `E1`：symbol exists（目标符号已存在）
- `E2`：test exists（已有对应测试）
- `E3`：verified output recorded（有验证输出记录）

规则：
- 所有 `partial/missing` 至少标注 `E0`。
- 关键路径迁移项建议达到 `E1/E2`。

## 数据基线

| 指标 | 数量 |
|---|---:|
| `main` other 文件总数 | 195 |
| `go` other 文件总数 | 334 |
| shared | 157 |
| main-only | 38 |
| go-only | 177 |

main-only（38）状态统计：
- `equivalent=5`
- `partial=9`
- `missing=23`
- `intentional-diff=1`

go-only（177）状态统计：
- `equivalent=5`
- `partial=1`
- `intentional-diff=171`

## 语义差异主表

| 语义域 | main 形态 | go 形态 | 状态 | 说明 | 整改 |
|---|---|---|---|---|---|
| Claude 工作区文档 | `.claude/*` | `.claude-plugin/.claude/*` | equivalent | 发生路径重构但语义一致 | 无 |
| 插件配置 | `.claude-plugin/settings.json` | `.claude-plugin/custom/config/setup.toml` | partial | 配置模型转换，未见字段级映射文档 | 补充 JSON->TOML 字段映射与转换步骤 |
| MCP 配置入口 | `.mcp.json` | `.dependencies/claude-skills` + go 配置体系 | partial | 入口形式变化，缺显式迁移说明 | 补充 mcp 入口兼容与迁移文档 |
| MCP Server 子项目 | `mcp-server/*` (TS) | `internal/providers/*`（语义承接） | missing | 无逐文件等价，子项目形态消失 | `decision=DEFER_WITH_CONDITION`；`reason=当前产品以 providers 语义承接为主`；`condition=若恢复独立 server 边界则新增 Go 子模块/桥接层` |
| OpenClaw 子项目 | `openclaw/*` (TS) | 无明确域 | missing | go 分支无逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| Legacy benchmark/live 资产 | `tests/benchmark/*`, `tests/live/README.md` | `internal/workflows/*_test.go` + 新文档 | partial | 测试体系迁移但未保留逐文件资产 | 建立旧资产到 go 测试用例的对照表 |
| Go runtime + custom layer | 无 | `cmd/`, `internal/`, `pkg/`, `custom/`, `.multipowers/` | intentional-diff | no-shell runtime 重构新增域 | 无 |

## 关键缺口决策与契约索引

| 域 | source scope | target file/domain | target symbol/contract | evidence level | decision | decision reason | evidence_upgrade_path | owner_domain |
|---|---|---|---|---|---|---|---|---|
| `mcp-server` | `mcp-server/*` | `internal/providers/*` | provider orchestration contract（`DetectAll` / `RouteIntent` / provider registry） | `E0` | `DEFER_WITH_CONDITION` | 当前产品以 Go providers 语义承接；若恢复独立 server 边界再拆分子模块 | `E0 -> E1: Create adapter interface in internal/providers/mcp_adapter.go` | providers |
| `openclaw` | `openclaw/*` | `N/A` | `N/A` | `E0` | `EXCLUDE_WITH_REASON` | 当前产品范围不包含 openclaw TS 子项目逐文件迁移；后续如恢复产品线再单独立项 | N/A | external |
| `legacy benchmark/live assets` | `tests/benchmark/*` + `tests/live/README.md` | `internal/workflows/*_test.go` + `docs/multipowers/README.md` | regression/benchmark/live test contract | `E0` | `MIGRATE_TO_GO` | 属于行为验证资产，需在 go 测试体系保持可追踪承接 | `E0 -> E2: Add TestBenchmarkRunner, TestLiveTestHarness in internal/workflows/*_test.go` | workflows |
| `.claude/settings.json` | `.claude/settings.json` | `.claude-plugin/.claude/settings.json` | Claude workspace settings | `E0` | `COPY_FROM_MAIN` | Claude 工作区配置需要保留 | `E0 -> E1: Copy file with path migration` | context |
| `.claude-plugin/settings.json` | `.claude-plugin/settings.json` | `.claude-plugin/custom/config/setup.toml` | plugin configuration | `E0` | `MIGRATE_TO_GO` | JSON -> TOML 配置模型转换 | `E0 -> E2: Document field mapping, add conversion test` | config |
| `.mcp.json` | `.mcp.json` | `.dependencies/claude-skills` | MCP dependency entry | `E0` | `MIGRATE_TO_GO` | MCP 入口形式变化 | `E0 -> E1: Document new dependency model` | deps |
| `docs/SCHEDULER.md` | `docs/SCHEDULER.md` | `docs/architecture/script-differences.md` | scheduler documentation | `E0` | `MIGRATE_TO_GO` | 调度能力转入脚本迁移文档 | `E0 -> E1: Add scheduler section to script-differences.md` | docs |
| `agents/personas/openclaw-admin.md` | `agents/personas/openclaw-admin.md` | `.claude-plugin/.claude/commands/persona.md` | persona configuration | `E0` | `DEFER_WITH_CONDITION` | 具体人格文档未逐文件保留 | `E0 -> E1: If persona needed, add to persona lanes config` | personas |

## main-only 文件映射（按域全量）

### 配置与元数据

| main 文件 | 状态 | go 对应文件/域 | 说明 | 整改 |
|---|---|---|---|---|
| `.claude-plugin/settings.json` | `partial` | `.claude-plugin/custom/config/setup.toml` | 插件配置模型从 settings.json 转为 setup.toml 体系 | 补充字段级配置映射说明并提供转换策略 |
| `.mcp.json` | `partial` | `.dependencies/claude-skills` | MCP/依赖入口改为 go 分支依赖目录与配置体系 | 补充 mcp 入口等价关系与迁移步骤 |

### Claude 工作区文档（路径重构）

| main 文件 | 状态 | go 对应文件/域 | 说明 | 整改 |
|---|---|---|---|---|
| `.claude/DEVELOPMENT.md` | `equivalent` | `.claude-plugin/.claude/DEVELOPMENT.md` | 路径重构：.claude -> .claude-plugin/.claude | 无 |
| `.claude/claude-octopus.local.md` | `equivalent` | `.claude-plugin/.claude/claude-octopus.local.md` | 路径重构：.claude -> .claude-plugin/.claude | 无 |
| `.claude/references/stub-detection.md` | `equivalent` | `.claude-plugin/.claude/references/stub-detection.md` | 路径重构：.claude -> .claude-plugin/.claude | 无 |
| `.claude/references/validation-gates.md` | `equivalent` | `.claude-plugin/.claude/references/validation-gates.md` | 路径重构：.claude -> .claude-plugin/.claude | 无 |
| `.claude/settings.json` | `missing` | `.claude-plugin/.claude/settings.json` | 预期路径重构目标不存在 | 补齐迁移文件或修正文档映射 |
| `.claude/state/state-manager.md` | `equivalent` | `.claude-plugin/.claude/state/state-manager.md` | 路径重构：.claude -> .claude-plugin/.claude | 无 |

### Legacy 文档

| main 文件 | 状态 | go 对应文件/域 | 说明 | 整改 |
|---|---|---|---|---|
| `STEELMAN.md` | `intentional-diff` | `N/A` | 策略文档未在 go 分支保留 | 无 |
| `docs/SCHEDULER.md` | `partial` | `docs/architecture/script-differences.md` | 调度能力转入脚本迁移与架构文档，缺独立 scheduler 专篇 | 补充 scheduler 专项文档或在现有文档中补充完整章节 |

### Legacy Personas

| main 文件 | 状态 | go 对应文件/域 | 说明 | 整改 |
|---|---|---|---|---|
| `agents/personas/openclaw-admin.md` | `partial` | `.claude-plugin/.claude/commands/persona.md` | 人格入口存在于 go，但该具体 persona 文档未逐文件保留 | 若需保留，迁移 persona 文档并接入 lane 配置 |

### MCP Server 子项目（TS）

| main 文件 | 状态 | go 对应文件/域 | 说明 | 整改 |
|---|---|---|---|---|
| `mcp-server/dist/index.d.ts` | `missing` | `internal/providers/*` | 独立 mcp-server TS 子项目在 go 分支无逐文件等价 | `decision=DEFER_WITH_CONDITION`；`reason=当前产品以 providers 语义承接为主`；`condition=若恢复独立 server 边界则新增 Go 子模块或桥接层` |
| `mcp-server/dist/index.js` | `missing` | `internal/providers/*` | 独立 mcp-server TS 子项目在 go 分支无逐文件等价 | `decision=DEFER_WITH_CONDITION`；`reason=当前产品以 providers 语义承接为主`；`condition=若恢复独立 server 边界则新增 Go 子模块或桥接层` |
| `mcp-server/dist/index.js.map` | `missing` | `internal/providers/*` | 独立 mcp-server TS 子项目在 go 分支无逐文件等价 | `decision=DEFER_WITH_CONDITION`；`reason=当前产品以 providers 语义承接为主`；`condition=若恢复独立 server 边界则新增 Go 子模块或桥接层` |
| `mcp-server/package.json` | `missing` | `internal/providers/*` | 独立 mcp-server TS 子项目在 go 分支无逐文件等价 | `decision=DEFER_WITH_CONDITION`；`reason=当前产品以 providers 语义承接为主`；`condition=若恢复独立 server 边界则新增 Go 子模块或桥接层` |
| `mcp-server/src/index.ts` | `missing` | `internal/providers/*` | 独立 mcp-server TS 子项目在 go 分支无逐文件等价 | `decision=DEFER_WITH_CONDITION`；`reason=当前产品以 providers 语义承接为主`；`condition=若恢复独立 server 边界则新增 Go 子模块或桥接层` |
| `mcp-server/src/schema/skill-schema.json` | `missing` | `internal/providers/*` | 独立 mcp-server TS 子项目在 go 分支无逐文件等价 | `decision=DEFER_WITH_CONDITION`；`reason=当前产品以 providers 语义承接为主`；`condition=若恢复独立 server 边界则新增 Go 子模块或桥接层` |
| `mcp-server/tsconfig.json` | `missing` | `internal/providers/*` | 独立 mcp-server TS 子项目在 go 分支无逐文件等价 | `decision=DEFER_WITH_CONDITION`；`reason=当前产品以 providers 语义承接为主`；`condition=若恢复独立 server 边界则新增 Go 子模块或桥接层` |

### OpenClaw 子项目（TS）

| main 文件 | 状态 | go 对应文件/域 | 说明 | 整改 |
|---|---|---|---|---|
| `openclaw/dist/index.d.ts` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/dist/index.js` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/dist/index.js.map` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/dist/skill-loader.d.ts` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/dist/skill-loader.js` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/dist/skill-loader.js.map` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/dist/tools/index.d.ts` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/dist/tools/index.js` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/dist/tools/index.js.map` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/openclaw.plugin.json` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/package.json` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/src/index.ts` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/src/skill-loader.ts` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/src/tools/index.ts` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |
| `openclaw/tsconfig.json` | `missing` | `N/A` | openclaw 子项目在 go 分支未保留逐文件对应 | `decision=EXCLUDE_WITH_REASON`；`reason=当前产品范围不包含 openclaw TS 子项目逐文件迁移`；`condition=若恢复 openclaw 产品线再单独立项` |

### Legacy 测试资产

| main 文件 | 状态 | go 对应文件/域 | 说明 | 整改 |
|---|---|---|---|---|
| `tests/benchmark/MANUAL-TEST-GUIDE.md` | `partial` | `internal/workflows/test_run_test.go` | 基准能力迁移到 Go 测试体系，原 shell/fixture 结构未逐文件保留 | 补充 shell benchmark -> go benchmark 的对照用例表 |
| `tests/benchmark/README.md` | `partial` | `internal/workflows/test_run_test.go` | 基准能力迁移到 Go 测试体系，原 shell/fixture 结构未逐文件保留 | 补充 shell benchmark -> go benchmark 的对照用例表 |
| `tests/benchmark/test-cases/vulnerable/sql-injection-login/code.py` | `partial` | `internal/workflows/test_run_test.go` | 基准能力迁移到 Go 测试体系，原 shell/fixture 结构未逐文件保留 | 补充 shell benchmark -> go benchmark 的对照用例表 |
| `tests/benchmark/test-cases/vulnerable/sql-injection-login/ground-truth.json` | `partial` | `internal/workflows/test_run_test.go` | 基准能力迁移到 Go 测试体系，原 shell/fixture 结构未逐文件保留 | 补充 shell benchmark -> go benchmark 的对照用例表 |
| `tests/live/README.md` | `partial` | `docs/multipowers/README.md` | live 测试说明迁移到新文档体系，原路径未保留 | 补充 live test 新入口与旧路径映射 |


## go-only 文件映射（按域全量）

### Plugin 封装与内置资产

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `.claude-plugin/.claude/DEVELOPMENT.md` | `equivalent` | `.claude/DEVELOPMENT.md` | 路径重构：.claude -> .claude-plugin/.claude |
| `.claude-plugin/.claude/claude-octopus.local.md` | `equivalent` | `.claude/claude-octopus.local.md` | 路径重构：.claude -> .claude-plugin/.claude |
| `.claude-plugin/.claude/references/stub-detection.md` | `equivalent` | `.claude/references/stub-detection.md` | 路径重构：.claude -> .claude-plugin/.claude |
| `.claude-plugin/.claude/references/validation-gates.md` | `equivalent` | `.claude/references/validation-gates.md` | 路径重构：.claude -> .claude-plugin/.claude |
| `.claude-plugin/.claude/state/state-manager.md` | `equivalent` | `.claude/state/state-manager.md` | 路径重构：.claude -> .claude-plugin/.claude |
| `.claude-plugin/bin/mp` | `intentional-diff` | `cmd/mp/main.go` | go 构建产物（二进制） |
| `.claude-plugin/bin/mp-devx` | `intentional-diff` | `cmd/mp-devx/main.go` | go 构建产物（二进制） |
| `.claude-plugin/custom/config/models.json` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |
| `.claude-plugin/custom/config/setup.toml` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |
| `.claude-plugin/scripts/mp` | `intentional-diff` | `scripts/mp` | 插件内入口脚本（wrapper） |
| `.claude-plugin/scripts/mp-devx` | `intentional-diff` | `scripts/mp-devx` | 插件内入口脚本（wrapper） |

### Multipowers 上下文

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `.multipowers/CLAUDE.md` | `intentional-diff` | `N/A` | go 分支初始化与产品上下文文档 |
| `.multipowers/FAQ.md` | `intentional-diff` | `N/A` | go 分支初始化与产品上下文文档 |
| `.multipowers/product-guidelines.md` | `intentional-diff` | `N/A` | go 分支初始化与产品上下文文档 |
| `.multipowers/product.md` | `intentional-diff` | `N/A` | go 分支初始化与产品上下文文档 |
| `.multipowers/setup_state.json` | `intentional-diff` | `N/A` | go 分支初始化与产品上下文文档 |
| `.multipowers/tech-stack.md` | `intentional-diff` | `N/A` | go 分支初始化与产品上下文文档 |
| `.multipowers/tracks.md` | `intentional-diff` | `N/A` | go 分支初始化与产品上下文文档 |
| `.multipowers/workflow.md` | `intentional-diff` | `N/A` | go 分支初始化与产品上下文文档 |

### Custom 定制层

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `custom/README.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/config/models.json` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/config/persona-lanes.json` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/config/proxy.json` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/config/setup.toml` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/customizations/conductor-context.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/customizations/models-and-lanes.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/customizations/persona-command.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/customizations/proxy-routing.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/reference/compatibility.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/reference/config-schema.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/reference/faq.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/sync/conflict-resolution.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/sync/go-upstream-diff-discipline.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/sync/upstream-sync-playbook.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/sync/verification-transcript.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/target-project/README.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/target-project/getting-started.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/target-project/troubleshooting.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/tool-project/README.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/tool-project/getting-started.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/tool-project/product-vision.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/tool-project/product.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/tool-project/tech-stack.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/docs/tool-project/troubleshooting.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/references/conductor-upstream/README.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/references/conductor-upstream/SOURCE-MAP.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/CLAUDE.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/FAQ.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/conductor/code_styleguides/cpp.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/conductor/code_styleguides/csharp.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/conductor/code_styleguides/dart.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/conductor/code_styleguides/general.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/conductor/code_styleguides/go.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/conductor/code_styleguides/html-css.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/conductor/code_styleguides/javascript.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/conductor/code_styleguides/python.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/conductor/code_styleguides/typescript.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |
| `custom/templates/conductor/workflow.md` | `intentional-diff` | `N/A` | 定制化层与模板资产 |

### Go Entrypoints

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `cmd/mp-devx/main.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `cmd/mp/main.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |

### Go Runtime Internal

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `internal/app/architecture_layout_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/app/errors.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/app/pipeline.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/app/pipeline_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/cli/root.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/cli/root_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/cli/status.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/cli/status_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/checker.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/checker_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/init_policy.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/init_policy.json` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/init_policy_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/init_runner.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/init_runner_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/loader.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/parity_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/requirements.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/setup_contract.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/setup_contract_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/summarizer.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/context/summarizer_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/devx/.keep` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/devx/runner.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/devx/runner_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/execx/result.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/execx/run.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/faq/classify.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/faq/dedup.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/faq/events.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/faq/render.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/faq/render_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/fsboundary/policy.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/fsboundary/policy_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/hooks/handler.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/hooks/handler_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/hooks/post_tool_use.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/hooks/pre_tool_use.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/hooks/session_start.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/hooks/stop.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/hooks/stop_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/modelroute/route.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/modelroute/route_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/claude.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/codex.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/degrade.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/degrade_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/detector.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/gemini.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/parity_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/provider.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/proxy.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/quorum.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/quorum_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/registry.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/router.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/router_intent.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/providers/router_intent_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/render/banner.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/runtime/config.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/runtime/prerun.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/tracks/checkbox.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/tracks/files.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/tracks/state.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/tracks/state_kv.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/tracks/state_kv_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/tracks/state_lock_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/tracks/state_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/tracks/track_id.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/util/jsonio.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/util/paths.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/validation/gates.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/validation/gates_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/validation/no_shell_runtime.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/validation/no_shell_runtime_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/validation/types.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/validation/types_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/coverage.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/coverage_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/debate.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/define.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/deliver.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/develop.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/discover.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/embrace.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/persona.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/persona_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/test_run.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `internal/workflows/test_run_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |

### Go Public API

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `pkg/api/README.md` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `pkg/api/jsonschema.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `pkg/api/types.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |
| `pkg/api/types_test.go` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |

### 文档与计划

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `docs/INDEX.md` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |
| `docs/MIGRATION-7.13.0.md` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |
| `docs/RELEASE-v7.17.0.md` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |
| `docs/USER_RESEARCH_PERSONAS.md` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |
| `docs/architecture/script-differences.md` | `partial` | `docs/SCHEDULER.md` | 脚本迁移分析文档承接了部分 scheduler 文档语义 |
| `docs/multipowers/README.md` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |
| `docs/plans/2026-03-02-main-go-semantic-and-file-mapping-design.md` | `intentional-diff` | `N/A` | go 分支迁移设计/实施文档 |
| `docs/plans/2026-03-02-no-shell-hybrid-rearchitecture-design.md` | `intentional-diff` | `N/A` | go 分支迁移设计/实施文档 |
| `docs/plans/2026-03-02-no-shell-hybrid-rearchitecture-implementation.md` | `intentional-diff` | `N/A` | go 分支迁移设计/实施文档 |
| `docs/plans/2026-03-02-upstream-v8.31.1-script-migration-mapping-design.md` | `intentional-diff` | `N/A` | go 分支迁移设计/实施文档 |
| `docs/plans/2026-03-02-v8.31.1-script-migration-mapping-implementation.md` | `intentional-diff` | `N/A` | go 分支迁移设计/实施文档 |

### Go Dev 脚本

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `scripts/check-go-file-length.go` | `intentional-diff` | `N/A` | go 分支开发辅助工具 |
| `scripts/mp` | `intentional-diff` | `N/A` | go 入口 wrapper |
| `scripts/mp-devx` | `intentional-diff` | `N/A` | go 入口 wrapper |

### 依赖与子模块

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `.dependencies/claude-skills` | `intentional-diff` | `N/A` | 依赖子模块与引用管理 |
| `.gitmodules` | `intentional-diff` | `N/A` | 依赖子模块与引用管理 |

### CI/CD

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `.github/workflows/go-ci.yml` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |

### Go Module

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `go.mod` | `intentional-diff` | `N/A` | go 运行时源码与模块定义 |

### 迁移与发布记录

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `IMPLEMENTATION_SUMMARY.md` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |
| `MIGRATION-7.23.0.md` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |
| `RELEASE_NOTES_7.25.0.md` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |
| `RELEASE_NOTES_GO_MIGRATION.md` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |
| `RELEASE_NOTES_NO_SHELL_RUNTIME.md` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |

### 运行时标记文件

| go 文件 | 状态 | main 对应文件/域 | 说明 |
|---|---|---|---|
| `plugin` | `intentional-diff` | `N/A` | go 分支新增文件（no-shell runtime 重构产物） |


## Parity 结论

- 在 other 文件域，`go` 与 `main` 的关系是“**结构性重构**”而非路径平移。
- 已确认等价的是文档路径重构类（`.claude -> .claude-plugin/.claude`）。
- 主要缺口集中在 `mcp-server/*` 与 `openclaw/*` 两个 main 子项目（当前判定 `missing`）。
- 当前判定：`other files = partial parity`。
