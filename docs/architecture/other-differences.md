# main vs go 其他类型文件差异（最新状态）

日期：2026-03-03  
比较分支：`main` vs `go`  
基线提交：`main=f6a815a326ec`，`go=cf865fa764fe`

## 范围与判定口径

- 本文覆盖：非 `commands/skills` 且非 `*.sh` 的文件。
- 状态定义：`equivalent` / `partial` / `missing` / `intentional-diff`
- 决策（decision）取值：`MIGRATE_TO_GO`、`COPY_FROM_MAIN`、`EXCLUDE_WITH_REASON`、`DEFER_WITH_CONDITION`
- 决策依据：`.multipowers/product-guidelines.md`、`.multipowers/product.md`

## Evidence Legend

- `E0`：doc-only plan（仅文档规划）
- `E1`：symbol exists（目标符号已存在）
- `E2`：test exists（已有对应测试）
- `E3`：verified output recorded（有验证输出记录）

## 最新统计

| 指标 | 数量 |
|---|---:|
| `main` other 文件总数 | 195 |
| `go` other 文件总数 | 371 |
| shared | 157 |
| main-only | 38 |
| go-only | 214 |

main-only（38）状态统计：
- `equivalent=5`
- `partial=9`
- `missing=23`
- `intentional-diff=1`

go-only（214）状态统计：
- `equivalent=5`
- `partial=1`
- `intentional-diff=208`

## 当前主差异域（最新）

| 语义域 | main | go | 状态 | decision |
|---|---|---|---|---|
| Claude 工作区文档 | `.claude/*` | `.claude-plugin/.claude/*` | `equivalent` | `COPY_FROM_MAIN` |
| 插件配置入口 | `.claude-plugin/settings.json` | `.claude-plugin/custom/config/setup.toml` | `partial` | `MIGRATE_TO_GO` |
| MCP 配置入口 | `.mcp.json` | `.dependencies/claude-skills` + go 配置体系 | `partial` | `MIGRATE_TO_GO` |
| MCP Server 子项目 | `mcp-server/*` | `internal/providers/*` 语义承接 | `missing` | `DEFER_WITH_CONDITION` |
| OpenClaw 子项目 | `openclaw/*` | `N/A` | `missing` | `EXCLUDE_WITH_REASON` |
| benchmark/live 资产 | `tests/benchmark/*`, `tests/live/*` | `internal/workflows/*_test.go` + 新文档 | `partial` | `MIGRATE_TO_GO` |
| no-shell runtime + custom 层 | `N/A` | `cmd/`, `internal/`, `pkg/`, `custom/`, `.multipowers/` | `intentional-diff` | `KEEP_IN_GO` |

## 关键决策快照（最新）

| source scope | target scope | decision | evidence | 最新状态 |
|---|---|---|---|---|
| `mcp-server/*` | `internal/providers/*` | `DEFER_WITH_CONDITION` | `E0` | 待独立 server 边界需求 |
| `openclaw/*` | `N/A` | `EXCLUDE_WITH_REASON` | `E0` | 当前产品范围排除 |
| `tests/benchmark/*` + `tests/live/*` | `internal/workflows/*_test.go` | `MIGRATE_TO_GO` | `E0` | 持续迁移中 |
| `.claude/settings.json` | `.claude-plugin/.claude/settings.json` | `COPY_FROM_MAIN` | `E0` | 需保持 workspace 配置承接 |
| `.claude-plugin/settings.json` | `.claude-plugin/custom/config/setup.toml` | `MIGRATE_TO_GO` | `E0` | 字段映射文档待补全 |

## go-only 增量说明（当前）

`go-only` 的新增文件以 no-shell runtime、sync 自动化、结构治理和自定义层为主，属于当前分支预期 `intentional-diff`，不需要回迁 `main`。

