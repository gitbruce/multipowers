# Main vs Go Commands & Skills One-to-One Mapping

日期：2026-03-05  
比较分支：`main` vs `go` (current)

本文件详细列出了 `main` 分支中所有命令（Commands）与技能（Skills）在 `go` 分支中的对应关系及定制化差异。

## 映射判定口径

| 状态 | 定义 | 备注 |
|---|---|---|
| `Equivalent` | 内容 100% 一致 | 基础能力，直接迁移。 |
| `Go-Customized (Model)` | 针对模型路由/能力进行了优化 | 修改了 prompt 中的模型建议或能力权重。 |
| `Go-Customized (Hook)` | 针对 Go 运行时钩子进行了适配 | 增加了对 `internal/hooks` 中质量门禁或策略的引用。 |
| `Re-written` | 内容针对 Go 引擎进行了重构 | 逻辑完全改变，以适配 Go 原子引擎。 |
| `Missing` | 仅在 `main` 分支中存在 | 尚未迁移或计划弃用。 |
| `New (Go-Only)` | 仅在 `go` 分支中存在 | 新增的原子能力或治理命令。 |

---

## Commands 映射清单 (46 total in main)

| Main Command | Go Path | Status | Mapping / Customization Details |
|---|---|---|---|
| `arch.md` | `.../commands/arch.md` | `Equivalent` | - |
| `bug.md` | `.../commands/bug.md` | `Go-Customized (Hook)` | 集成了 `quality-gate` 钩子以确调试信息完整性。 |
| `chat.md` | `.../commands/chat.md` | `Equivalent` | - |
| `claw.md` | `N/A` | `Missing` | 待产品范围确认。 |
| `commit.md` | `.../commands/commit.md` | `Go-Customized (Model)` | 优化了针对 `gpt-4o` 的提交信息生成策略。 |
| `config.md` | `.../commands/config.md` | `Re-written` | 完全重写，以映射到 Go 侧的 `internal/settings`。 |
| `context.md` | `.../commands/context.md` | `Go-Customized (Hook)` | 适配了 Go 侧的原子上下文边界检查。 |
| `debate.md` | `.../commands/debate.md` | `Equivalent` | - |
| `deliver.md` | `.../commands/deliver.md` | `Go-Customized (Hook)` | 强制执行交付前的 P95 性能基准测试校验。 |
| `deploy.md` | `.../commands/deploy.md` | `Equivalent` | - |
| `develop.md` | `.../commands/develop.md` | `Go-Customized (Model)` | 针对 Go 源码开发优化了思考路径。 |
| `diff.md` | `.../commands/diff.md` | `Equivalent` | - |
| `doctor.md` | `N/A` | `Missing` | 已被 `mp-devx` 的内部校验取代。 |
| `edit.md` | `.../commands/edit.md` | `Equivalent` | - |
| `embrace.md` | `.../commands/embrace.md` | `Re-written` | 从单纯的脚本包装转为 Go 运行时的内部工作流委派。 |
| `extract.md` | `.../commands/extract.md` | `Re-written` | 迁移至 `internal/extract` 逻辑实现。 |
| `feedback.md` | `.../commands/feedback.md` | `Equivalent` | - |
| `help.md` | `.../commands/help.md` | `Go-Customized (Model)` | 更新了支持的 Go 侧 Action 列表。 |
| `init.md` (shared) | `.../commands/init.md` | `Re-written` | 适配 Go 版项目的脚手架初始化逻辑。 |
| `install.md` | `.../commands/install.md` | `Equivalent` | - |
| `learn.md` | `.../commands/learn.md` | `Equivalent` | - |
| `log.md` | `.../commands/log.md` | `Equivalent` | - |
| `mcp.md` | `.../commands/mcp.md` | `Equivalent` | - |
| `migrate.md` | `.../commands/migrate.md` | `Equivalent` | - |
| `octo.md` | `.../commands/mp.md` | `Re-written` | **重点**：重命名并重写为 `mp` 命令作为统一入口。 |
| `parallel.md` | `N/A` | `Missing` | 待 Go 侧并行执行器稳定后再合并。 |
| `plan.md` | `.../commands/plan.md` | `Go-Customized (Hook)` | 引入了音轨（Tracks）计划的结构校验钩子。 |
| `policy.md` | `.../commands/policy.md` | `Equivalent` | - |
| `preview.md` | `.../commands/preview.md` | `Equivalent` | - |
| `publish.md` | `.../commands/publish.md` | `Equivalent` | - |
| `quality.md` | `.../commands/quality.md` | `Go-Customized (Hook)` | 直接调用 Go 侧的 `internal/validation` 门禁。 |
| `recheck.md` | `.../commands/recheck.md` | `Equivalent` | - |
| `refactor.md` | `.../commands/refactor.md` | `Equivalent` | - |
| `release.md` | `.../commands/release.md` | `Equivalent` | - |
| `review.md` | `.../commands/review.md` | `Equivalent` | - |
| `run.md` | `.../commands/run.md` | `Equivalent` | - |
| `schedule.md` | `N/A` | `Missing` | 待调度域落地。 |
| `scheduler.md` | `N/A` | `Missing` | 待调度域落地。 |
| `security.md` | `.../commands/security.md` | `Equivalent` | - |
| `sentinel.md` | `N/A` | `Missing` | 已被 Go 侧的物理门禁（Isolation）机制取代。 |
| `spec.md` | `N/A` | `Missing` | 待产品定义最终确认。 |
| `status.md` | `.../commands/status.md` | `Equivalent` | - |
| `summarize.md` | `.../commands/summarize.md` | `Equivalent` | - |
| `task.md` | `.../commands/task.md` | `Equivalent` | - |
| `test.md` | `.../commands/test.md` | `Go-Customized (Hook)` | 运行时入口弱化，测试/覆盖率执行收敛至 `mp-devx` 运维命令。 |
| `trace.md` | `.../commands/trace.md` | `Equivalent` | - |

---

## Skills 映射清单 (46 total in main)

| Main Skill | Go Path | Status | Mapping / Customization Details |
|---|---|---|---|
| `skill-arch.md` | `.../skills/skill-arch.md` | `Equivalent` | - |
| `skill-bug.md` | `.../skills/skill-bug.md` | `Equivalent` | - |
| `skill-chat.md` | `.../skills/skill-chat.md` | `Equivalent` | - |
| `skill-claw.md` | `N/A` | `Missing` | 暂不提供 Claw 能力。 |
| `skill-commit.md` | `.../skills/skill-commit.md` | `Equivalent` | - |
| `skill-config.md` | `.../skills/skill-config.md` | `Go-Customized (Model)` | 针对 Go 配置项 schema 的理解进行了优化。 |
| `skill-context.md` | `.../skills/skill-context.md` | `Go-Customized (Hook)` | 指导如何与 Go 侧的 Context Guard 交互。 |
| `skill-debate.md` | `.../skills/skill-debate.md` | `Equivalent` | - |
| `skill-deliver.md` | `.../skills/skill-deliver.md` | `Equivalent` | - |
| `skill-deploy.md` | `.../skills/skill-deploy.md` | `Equivalent` | - |
| `skill-develop.md` | `.../skills/skill-develop.md` | `Equivalent` | - |
| `skill-diff.md` | `.../skills/skill-diff.md` | `Equivalent` | - |
| `skill-doctor.md` | `N/A` | `Missing` | 暂不提供该技能。 |
| `skill-edit.md` | `.../skills/skill-edit.md` | `Equivalent` | - |
| `skill-embrace.md` | `.../skills/skill-embrace.md` | `Equivalent` | - |
| `skill-extract.md` | `.../skills/skill-extract.md` | `Equivalent` | - |
| `skill-feedback.md` | `.../skills/skill-feedback.md` | `Equivalent` | - |
| `skill-help.md` | `.../skills/skill-help.md` | `Equivalent` | - |
| `skill-init.md` | `.../skills/skill-init.md` | `Equivalent` | - |
| `skill-install.md` | `.../skills/skill-install.md` | `Equivalent` | - |
| `skill-learn.md` | `.../skills/skill-learn.md` | `Equivalent` | - |
| `skill-log.md` | `.../skills/skill-log.md` | `Equivalent` | - |
| `skill-mcp.md` | `.../skills/skill-mcp.md` | `Equivalent` | - |
| `skill-migrate.md` | `.../skills/skill-migrate.md` | `Equivalent` | - |
| `skill-octo.md` | `.../skills/skill-mp.md` | `Equivalent` | 更名为 `skill-mp`。 |
| `skill-parallel.md` | `.../skills/skill-parallel.md` | `Equivalent` | - |
| `skill-plan.md` | `.../skills/skill-plan.md` | `Equivalent` | - |
| `skill-policy.md` | `.../skills/skill-policy.md` | `Equivalent` | - |
| `skill-preview.md` | `.../skills/skill-preview.md` | `Equivalent` | - |
| `skill-publish.md` | `.../skills/skill-publish.md` | `Equivalent` | - |
| `skill-quality.md` | `.../skills/skill-quality.md` | `Equivalent` | - |
| `skill-recheck.md` | `.../skills/skill-recheck.md` | `Equivalent` | - |
| `skill-refactor.md` | `.../skills/skill-refactor.md` | `Equivalent` | - |
| `skill-release.md` | `.../skills/skill-release.md` | `Equivalent` | - |
| `skill-review.md` | `.../skills/skill-review.md` | `Equivalent` | - |
| `skill-run.md` | `.../skills/skill-run.md` | `Equivalent` | - |
| `skill-schedule.md` | `.../skills/skill-schedule.md` | `Equivalent` | - |
| `skill-scheduler.md` | `.../skills/skill-scheduler.md` | `Equivalent` | - |
| `skill-security.md` | `.../skills/skill-security.md` | `Equivalent` | - |
| `skill-sentinel.md` | `.../skills/skill-sentinel.md` | `Equivalent` | - |
| `skill-spec.md` | `.../skills/skill-spec.md` | `Equivalent` | - |
| `skill-status.md` | `.../skills/skill-status.md` | `Equivalent` | - |
| `skill-summarize.md` | `.../skills/skill-summarize.md` | `Equivalent` | - |
| `skill-task.md` | `.../skills/skill-task.md` | `Equivalent` | - |
| `skill-test.md` | `.../skills/skill-test.md` | `Equivalent` | - |
| `skill-trace.md` | `.../skills/skill-trace.md` | `Equivalent` | - |

---

## 总结结论

1.  **高对齐度**：`Skills` 侧整体保持高一致性（44/46 已映射到 Go；2 项为 `Missing`）。
2.  **深度定制**：`Commands` 侧共有 14 项发生显著变化（`Go-Customized` 9 项 + `Re-written` 5 项），核心目的是适配 Go 原子运行时与 Hook 质量门禁。
3.  **治理收敛**：`main-only` 命令目前为 7 项，主要是待迁移项或已被 Go 侧底层机制替代的治理项。
4.  **技能层改动更克制**：`Skills` 侧仅 2 项发生定制（`skill-config`、`skill-context`），其余以等价迁移为主。
