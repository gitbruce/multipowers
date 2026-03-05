# Main vs Go Folder Mapping

日期：2026-03-05  
比较分支：`main` vs `go`

该文档描述了 `main` 分支（Bash/JS 混合版）与 `go` 分支（Go 原子引擎版）之间的目录级映射关系。它旨在帮助维护者理解在迁移过程中，各项能力与资产在目录层面的流转路径。

## 核心能力映射 (Core Capability Folders)

| Source Folder (`main`) | Target Folder (`go`) | Parity Status | Description |
|---|---|---|---|
| `.claude/` | `.claude-plugin/.claude/` | `Moved` | 存放命令、技能和参考资料。在 `go` 分支中，这些资产已移动到插件标准的隐藏目录下。 |
| `scripts/` | `internal/` + `scripts/` | `Transformed` | `main` 中的 Bash 业务脚本已大量重构为 `internal/` 下的 Go 领域包。`scripts/` 仅保留少数包装脚本。 |
| `hooks/` | `internal/hooks/` | `Transformed` | 生命周期钩子从独立的 Shell 脚本转变为 `mp-devx` 驱动的 Go 内置钩子处理器。 |
| `templates/` | `templates/` | `Equivalent` | 项目初始化和状态模板基本保持同构。 |

## 基础设施与运行时 (Infrastructure & Runtime)

| Source Folder (`main`) | Target Folder (`go`) | Parity Status | Description |
|---|---|---|---|
| `N/A` | `cmd/mp/` | `New` | `go` 分支的核心入口，取代了 `main` 中的多个独立脚本入口。 |
| `N/A` | `cmd/mp-devx/` | `New` | 开发工具链入口，包含测试、基准测试和策略构建功能。 |
| `config/` | `config/` | `Updated` | 从纯 YAML 配置升级为包含 Provider、Workflow 等强类型定义的配置目录。 |
| `N/A` | `pkg/api/` | `New` | 定义 Go 版 Octopus 的外部 API 契约。 |

## 文档与治理 (Docs & Governance)

| Source Folder (`main`) | Target Folder (`go`) | Parity Status | Description |
|---|---|---|---|
| `docs/` | `docs/` | `Equivalent` | 架构文档和用户指南。`go` 分支在此目录下增加了 `plans/`，并在其下维护 `plans/evidence/` 以追踪迁移进度。 |
| `N/A` | `.multipowers/` | `New` | 存放 `go` 分支特有的工作协议（CLAUDE.md）、技术栈定义和任务追踪器。 |
| `N/A` | `conductor/` | `New` | `conductor` 扩展使用的音轨计划和元数据目录。 |
| `custom/` | `custom/` | `Updated` | 存放用户自定义的文档、模板和配置，结构保持对齐但内容更丰富。 |

## 测试与校验 (Tests & Validation)

| Source Folder (`main`) | Target Folder (`go`) | Parity Status | Description |
|---|---|---|---|
| `tests/` | `tests/` | `Updated` | 从 Shell 测试脚本过渡为调用 Go 单元测试和集成测试的入口。 |
| `N/A` | `.dependencies/` | `New` | 存放本地编译和开发所需的第三方依赖快照。 |

## 总结结论

`go` 分支对 `main` 分支的目录结构进行了显著的**模块化重构**：
1.  **能力下沉**：原有的顶级脚本逻辑下沉到了 `internal/` 目录中，增强了类型安全和可测试性。
2.  **插件规范化**：将 Prompt 资产约束在 `.claude-plugin/` 目录下，符合插件化部署规范。
3.  **治理隔离**：通过 `.multipowers/` 目录隔离了项目的管理元数据与业务逻辑代码。
4.  **文档追踪结构化**：迁移证据并非独立的 `docs/evidence/`，而是归档在 `docs/plans/evidence/`，与计划文档形成闭环。
