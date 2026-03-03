# main vs go 命令与技能差异（最新状态）

日期：2026-03-03  
比较分支：`main` vs `go`  
基线提交：`main=f6a815a326ec`，`go=cf865fa764fe`

## 范围与判定口径

- 仅关注：
  - `main:.claude/commands/*` vs `go:.claude/commands/*`
  - `main:.claude/skills/*` vs `go:.claude/skills/*`
- 状态定义：`equivalent` / `partial` / `missing` / `intentional-diff`
- 决策（decision）取值：`MIGRATE_TO_GO`、`COPY_FROM_MAIN`、`EXCLUDE_WITH_REASON`、`DEFER_WITH_CONDITION`
- 决策依据：`.multipowers/product-guidelines.md`、`.multipowers/product.md`

## Evidence Legend

- `E0`：doc-only plan（仅文档规划）
- `E1`：symbol exists（目标符号已存在）
- `E2`：test exists（已有对应测试）
- `E3`：verified output recorded（有验证输出记录）

## 最新统计

名称级交集规模：
- commands：`main=46`，`go=41`，`shared=38`，`main-only=8`，`go-only=3`
- skills：`main=48`，`go=47`，`shared=46`，`main-only=2`，`go-only=1`

内容级状态统计（最新）：
- commands：`equivalent=32`，`partial=10`，`missing=5`，`intentional-diff=2`
- skills：`equivalent=46`，`partial=0`，`missing=2`，`intentional-diff=1`

当前结论：`commands/skills = partial parity`

## 与 `tmp/compare.md` 的对齐校验（2026-03-03）

- `commands` 与 `skills` 的关键统计项复核一致：
  - `commands (main-only)=8`
  - `skills (main-only)=2`
- 本轮无需调整差异清单与决策快照，仅补充该校验结论。

## 目录结构治理（当前生效）

规则文件：`config/sync/claude-structure-rules.json`

- `MUST_HOMOMORPHIC`：
  - `.claude/commands` -> `.claude/commands`（共享子集）
  - `.claude/skills` -> `.claude/skills`（共享子集）
  - `.claude/references` -> `.claude/references`
  - `.claude/state` -> `.claude/state`
- `ALLOW_FORK`：
  - `.claude/commands/init.md`
  - `.claude/commands/mp.md`
  - `.claude/commands/persona.md`
  - `.claude/skills/skill-persona.md`
  - 以及规则内显式 ignore 列表

校验入口：`./scripts/validate-claude-structure.sh -dry-run`

## 当前差异清单（仅最新）

main-only commands：
- `claw.md`
- `doctor.md`
- `octo.md`
- `parallel.md`
- `schedule.md`
- `scheduler.md`
- `sentinel.md`
- `spec.md`

go-only commands：
- `init.md`
- `mp.md`
- `persona.md`

main-only skills：
- `skill-claw.md`
- `skill-doctor.md`

go-only skills：
- `skill-persona.md`

## 高风险差异决策快照

| source | target | decision | evidence | 最新状态 |
|---|---|---|---|---|
| `.claude/commands/octo.md` | `.claude/commands/mp.md` | `MIGRATE_TO_GO` | `E0` | 仍为 `partial` |
| `.claude/commands/sentinel.md` | `internal/hooks/*`（治理门禁域） | `MIGRATE_TO_GO` | `E0` | 仍为 `missing` |
| `.claude/commands/schedule.md` + `.claude/commands/scheduler.md` | `internal/scheduler/*` | `DEFER_WITH_CONDITION` | `E0` | 待 scheduler 域契约落地 |
| `.claude/commands/claw.md` | `N/A` | `DEFER_WITH_CONDITION` | `E0` | 待产品范围确认 |
| `.claude/commands/doctor.md` | `sys-configure` 能力域 | `EXCLUDE_WITH_REASON` | `E0` | 维持排除 |
| `.claude/skills/skill-claw.md` + `.claude/skills/skill-doctor.md` | `N/A` | `EXCLUDE_WITH_REASON` | `E0` | 维持排除 |
