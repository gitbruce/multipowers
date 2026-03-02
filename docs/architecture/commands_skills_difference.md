# main vs go 命令与技能差异（语义优先 + 文件一对一补充）

日期：2026-03-02  
比较分支：`main` vs `go`  
基线提交：`main=f6a815a326ec`，`go=d01e74d99977`

## 判定口径

状态枚举（全表统一）：
- `equivalent`：语义能力等价，且可追踪到明确文件映射。
- `partial`：存在部分能力或入口，但行为/入口/覆盖不完整。
- `missing`：`main` 能力在 `go` 无对应实现。
- `intentional-diff`：有意的分支差异（新增或策略性不迁移）。

规则：
1. 先做语义能力一对一映射。  
2. 语义非双射时，补充文件级一对一映射。  
3. 所有 `partial/missing` 必须给出整改路径。

## 结果概览

| 维度 | main | go | shared | main-only | go-only |
|---|---:|---:|---:|---:|---:|
| commands (`*.md`) | 46 | 41 | 38 | 8 | 3 |
| skills (`*.md`) | 48 | 47 | 46 | 2 | 1 |

命令映射状态统计（49 行，含 go-only 补充行）：
- `equivalent=38`
- `partial=4`
- `missing=5`
- `intentional-diff=2`

技能映射状态统计（49 行，含 go-only 补充行）：
- `equivalent=46`
- `missing=2`
- `intentional-diff=1`

结论：`go` 在命令/技能维度尚未达到与 `main` 的完全等价（主要缺口集中在 `claw/doctor/schedule/scheduler/sentinel`）。

## 语义差异主表

| domain | main capability | go counterpart | status | evidence | remediation |
|---|---|---|---|---|---|
| command | `/octo:octo` smart intent router | `/mp:mp` root alias | partial | both are root entry commands, but `main` includes intent detection/scoring and `go` alias is thin | add intent routing policy/threshold behavior to go router layer |
| command | `/octo:parallel` command入口 | `flow-parallel` skill (no command) | partial | capability exists as skill name in both branches, but go lacks dedicated command wrapper | add `/mp:parallel` command mapped to existing parallel flow |
| command | `/octo:spec` command入口 | `flow-spec` skill (no command) | partial | spec flow skill shared; command endpoint absent in go | add `/mp:spec` command wrapper |
| command | `/octo:claw` OpenClaw 管理 | none | missing | main has command + skill-claw, go has neither | decide deprecation or port claw command + skill |
| command | `/octo:doctor` diagnostics | none | missing | main has command + skill-doctor, go lacks dedicated equivalent | add `/mp:doctor` plus diagnostics coverage |
| command | `/octo:schedule` jobs CRUD | none | missing | main has scheduler job command, no go counterpart | port scheduler job management interface |
| command | `/octo:scheduler` daemon control | none | missing | main has daemon lifecycle command, no go counterpart | port daemon control interface |
| command | `/octo:sentinel` monitor | none | missing | command exists only in main | port or mark explicitly removed with replacement guidance |
| command | go-only `/mp:init` | none in main | intentional-diff | dedicated `.multipowers` bootstrap exists only in go | keep additive; document as go-specific |
| command | go-only `/mp:persona` | none in main | intentional-diff | dedicated persona command only in go | keep additive; optional backport not required |
| skill | `skill-claw` | none | missing | skill file exists only in main | port or deprecate explicitly |
| skill | `skill-doctor` | none | missing | skill file exists only in main | port or consolidate into other go diagnostics entrypoints |
| skill | go-only `skill-persona` | none in main | intentional-diff | new persona skill in go | keep additive |

## 文件级一对一补充：Commands（全文件）

| main command | main file | go file | status | notes/remediation |
|---|---|---|---|---|
| brainstorm | `.claude/commands/brainstorm.md` | `.claude-plugin/.claude/commands/brainstorm.md` | equivalent | same command name; runtime moved to mp binary |
| claw | `.claude/commands/claw.md` | `N/A` | missing | no OpenClaw admin command/skill parity in go; remediation: port command + capability or declare deprecation |
| debate | `.claude/commands/debate.md` | `.claude-plugin/.claude/commands/debate.md` | equivalent | same command name; runtime moved to mp binary |
| debug | `.claude/commands/debug.md` | `.claude-plugin/.claude/commands/debug.md` | equivalent | same command name; runtime moved to mp binary |
| deck | `.claude/commands/deck.md` | `.claude-plugin/.claude/commands/deck.md` | equivalent | same command name; runtime moved to mp binary |
| define | `.claude/commands/define.md` | `.claude-plugin/.claude/commands/define.md` | equivalent | same command name; runtime moved to mp binary |
| deliver | `.claude/commands/deliver.md` | `.claude-plugin/.claude/commands/deliver.md` | equivalent | same command name; runtime moved to mp binary |
| dev | `.claude/commands/dev.md` | `.claude-plugin/.claude/commands/dev.md` | equivalent | same command name; runtime moved to mp binary |
| develop | `.claude/commands/develop.md` | `.claude-plugin/.claude/commands/develop.md` | equivalent | same command name; runtime moved to mp binary |
| discover | `.claude/commands/discover.md` | `.claude-plugin/.claude/commands/discover.md` | equivalent | same command name; runtime moved to mp binary |
| docs | `.claude/commands/docs.md` | `.claude-plugin/.claude/commands/docs.md` | equivalent | same command name; runtime moved to mp binary |
| doctor | `.claude/commands/doctor.md` | `N/A` | missing | no dedicated diagnostics command parity; remediation: add /mp:doctor with internal checks |
| embrace | `.claude/commands/embrace.md` | `.claude-plugin/.claude/commands/embrace.md` | equivalent | same lifecycle capability; markdown logic moved into Go runtime |
| extract | `.claude/commands/extract.md` | `.claude-plugin/.claude/commands/extract.md` | equivalent | same command name; runtime moved to mp binary |
| grasp | `.claude/commands/grasp.md` | `.claude-plugin/.claude/commands/grasp.md` | equivalent | same command name; runtime moved to mp binary |
| ink | `.claude/commands/ink.md` | `.claude-plugin/.claude/commands/ink.md` | equivalent | same command name; runtime moved to mp binary |
| issues | `.claude/commands/issues.md` | `.claude-plugin/.claude/commands/issues.md` | equivalent | same command name; runtime moved to mp binary |
| km | `.claude/commands/km.md` | `.claude-plugin/.claude/commands/km.md` | equivalent | same command name; runtime moved to mp binary |
| loop | `.claude/commands/loop.md` | `.claude-plugin/.claude/commands/loop.md` | equivalent | same command name; runtime moved to mp binary |
| meta-prompt | `.claude/commands/meta-prompt.md` | `.claude-plugin/.claude/commands/meta-prompt.md` | equivalent | same command name; runtime moved to mp binary |
| model-config | `.claude/commands/model-config.md` | `.claude-plugin/.claude/commands/model-config.md` | equivalent | same command name; runtime moved to mp binary |
| multi | `.claude/commands/multi.md` | `.claude-plugin/.claude/commands/multi.md` | equivalent | same command name; runtime moved to mp binary |
| octo | `.claude/commands/octo.md` | `.claude-plugin/.claude/commands/mp.md` | partial | root routing exists as /mp, but smart intent router behavior is reduced; remediation: add intent-routing parity in /mp or dedicated router command |
| parallel | `.claude/commands/parallel.md` | `.claude-plugin/.claude/skills/flow-parallel.md` | partial | parallel capability exists at skill layer but command entry /mp:parallel missing; remediation: add command wrapper |
| pipeline | `.claude/commands/pipeline.md` | `.claude-plugin/.claude/commands/pipeline.md` | equivalent | same command name; runtime moved to mp binary |
| plan | `.claude/commands/plan.md` | `.claude-plugin/.claude/commands/plan.md` | equivalent | same command name; runtime moved to mp binary |
| prd | `.claude/commands/prd.md` | `.claude-plugin/.claude/commands/prd.md` | equivalent | same command name; runtime moved to mp binary |
| prd-score | `.claude/commands/prd-score.md` | `.claude-plugin/.claude/commands/prd-score.md` | equivalent | same command name; runtime moved to mp binary |
| probe | `.claude/commands/probe.md` | `.claude-plugin/.claude/commands/probe.md` | equivalent | same command name; runtime moved to mp binary |
| quick | `.claude/commands/quick.md` | `.claude-plugin/.claude/commands/quick.md` | equivalent | same command name; runtime moved to mp binary |
| research | `.claude/commands/research.md` | `.claude-plugin/.claude/commands/research.md` | equivalent | same command name; runtime moved to mp binary |
| resume | `.claude/commands/resume.md` | `.claude-plugin/.claude/commands/resume.md` | equivalent | same command name; runtime moved to mp binary |
| review | `.claude/commands/review.md` | `.claude-plugin/.claude/commands/review.md` | equivalent | same command name; runtime moved to mp binary |
| rollback | `.claude/commands/rollback.md` | `.claude-plugin/.claude/commands/rollback.md` | equivalent | same command name; runtime moved to mp binary |
| schedule | `.claude/commands/schedule.md` | `N/A` | missing | no schedule command parity; remediation: add scheduler job management command |
| scheduler | `.claude/commands/scheduler.md` | `N/A` | missing | no scheduler daemon management parity; remediation: add scheduler control command |
| security | `.claude/commands/security.md` | `.claude-plugin/.claude/commands/security.md` | equivalent | same command name; runtime moved to mp binary |
| sentinel | `.claude/commands/sentinel.md` | `N/A` | missing | no sentinel monitoring command parity; remediation: port sentinel workflow |
| setup | `.claude/commands/setup.md` | `.claude-plugin/.claude/commands/setup.md` | equivalent | same command name; runtime moved to mp binary |
| ship | `.claude/commands/ship.md` | `.claude-plugin/.claude/commands/ship.md` | equivalent | same command name; runtime moved to mp binary |
| spec | `.claude/commands/spec.md` | `.claude-plugin/.claude/skills/flow-spec.md` | partial | spec capability exists at skill layer but /mp:spec command missing; remediation: add command wrapper |
| status | `.claude/commands/status.md` | `.claude-plugin/.claude/commands/status.md` | equivalent | same command name; runtime moved to mp binary |
| sys-setup | `.claude/commands/sys-setup.md` | `.claude-plugin/.claude/commands/sys-setup.md` | equivalent | same command name; runtime moved to mp binary |
| tangle | `.claude/commands/tangle.md` | `.claude-plugin/.claude/commands/tangle.md` | equivalent | same command name; runtime moved to mp binary |
| tdd | `.claude/commands/tdd.md` | `.claude-plugin/.claude/commands/tdd.md` | equivalent | same command name; runtime moved to mp binary |
| validate | `.claude/commands/validate.md` | `.claude-plugin/.claude/commands/validate.md` | equivalent | same command name; runtime moved to mp binary |
| (go-only) init | `N/A` | `.claude-plugin/.claude/commands/init.md` | intentional-diff | new go-only guided bootstrap for .multipowers |
| (go-only) mp | `.claude/commands/octo.md` | `.claude-plugin/.claude/commands/mp.md` | partial | acts as root alias equivalent surface, but less routing logic than main octo |
| (go-only) persona | `N/A` | `.claude-plugin/.claude/commands/persona.md` | intentional-diff | new go-only persona launcher; not present as dedicated command in main |

## 文件级一对一补充：Skills（全文件）

| main skill | main file | go file | status | notes/remediation |
|---|---|---|---|---|
| extract-skill | `.claude/skills/extract-skill.md` | `.claude-plugin/.claude/skills/extract-skill.md` | equivalent | same skill name and intent |
| flow-define | `.claude/skills/flow-define.md` | `.claude-plugin/.claude/skills/flow-define.md` | equivalent | same skill name and intent |
| flow-deliver | `.claude/skills/flow-deliver.md` | `.claude-plugin/.claude/skills/flow-deliver.md` | equivalent | same skill name and intent |
| flow-develop | `.claude/skills/flow-develop.md` | `.claude-plugin/.claude/skills/flow-develop.md` | equivalent | same skill name and intent |
| flow-discover | `.claude/skills/flow-discover.md` | `.claude-plugin/.claude/skills/flow-discover.md` | equivalent | same skill name and intent |
| flow-parallel | `.claude/skills/flow-parallel.md` | `.claude-plugin/.claude/skills/flow-parallel.md` | equivalent | same skill name and intent |
| flow-spec | `.claude/skills/flow-spec.md` | `.claude-plugin/.claude/skills/flow-spec.md` | equivalent | same skill name and intent |
| skill-adversarial-security | `.claude/skills/skill-adversarial-security.md` | `.claude-plugin/.claude/skills/skill-adversarial-security.md` | equivalent | same skill name and intent |
| skill-architecture | `.claude/skills/skill-architecture.md` | `.claude-plugin/.claude/skills/skill-architecture.md` | equivalent | same skill name and intent |
| skill-audit | `.claude/skills/skill-audit.md` | `.claude-plugin/.claude/skills/skill-audit.md` | equivalent | same skill name and intent |
| skill-claw | `.claude/skills/skill-claw.md` | `N/A` | missing | no OpenClaw admin skill parity in go; remediation: port skill with go runtime contract |
| skill-code-review | `.claude/skills/skill-code-review.md` | `.claude-plugin/.claude/skills/skill-code-review.md` | equivalent | same skill name and intent |
| skill-content-pipeline | `.claude/skills/skill-content-pipeline.md` | `.claude-plugin/.claude/skills/skill-content-pipeline.md` | equivalent | same skill name and intent |
| skill-context-detection | `.claude/skills/skill-context-detection.md` | `.claude-plugin/.claude/skills/skill-context-detection.md` | equivalent | same skill name and intent |
| skill-debate | `.claude/skills/skill-debate.md` | `.claude-plugin/.claude/skills/skill-debate.md` | equivalent | same skill name and intent |
| skill-debate-integration | `.claude/skills/skill-debate-integration.md` | `.claude-plugin/.claude/skills/skill-debate-integration.md` | equivalent | same skill name and intent |
| skill-debug | `.claude/skills/skill-debug.md` | `.claude-plugin/.claude/skills/skill-debug.md` | equivalent | same skill name and intent |
| skill-decision-support | `.claude/skills/skill-decision-support.md` | `.claude-plugin/.claude/skills/skill-decision-support.md` | equivalent | same skill name and intent |
| skill-deck | `.claude/skills/skill-deck.md` | `.claude-plugin/.claude/skills/skill-deck.md` | equivalent | same skill name and intent |
| skill-deep-research | `.claude/skills/skill-deep-research.md` | `.claude-plugin/.claude/skills/skill-deep-research.md` | equivalent | same skill name and intent |
| skill-doc-delivery | `.claude/skills/skill-doc-delivery.md` | `.claude-plugin/.claude/skills/skill-doc-delivery.md` | equivalent | same skill name and intent |
| skill-doctor | `.claude/skills/skill-doctor.md` | `N/A` | missing | no doctor diagnostics skill parity in go; remediation: add doctor skill or merge into /mp:sys-setup + /mp:status with equivalent checks |
| skill-finish-branch | `.claude/skills/skill-finish-branch.md` | `.claude-plugin/.claude/skills/skill-finish-branch.md` | equivalent | same skill name and intent |
| skill-intent-contract | `.claude/skills/skill-intent-contract.md` | `.claude-plugin/.claude/skills/skill-intent-contract.md` | equivalent | same skill name and intent |
| skill-issues | `.claude/skills/skill-issues.md` | `.claude-plugin/.claude/skills/skill-issues.md` | equivalent | same skill name and intent |
| skill-iterative-loop | `.claude/skills/skill-iterative-loop.md` | `.claude-plugin/.claude/skills/skill-iterative-loop.md` | equivalent | same skill name and intent |
| skill-knowledge-work | `.claude/skills/skill-knowledge-work.md` | `.claude-plugin/.claude/skills/skill-knowledge-work.md` | equivalent | same skill name and intent |
| skill-meta-prompt | `.claude/skills/skill-meta-prompt.md` | `.claude-plugin/.claude/skills/skill-meta-prompt.md` | equivalent | same skill name and intent |
| skill-parallel-agents | `.claude/skills/skill-parallel-agents.md` | `.claude-plugin/.claude/skills/skill-parallel-agents.md` | equivalent | same skill name and intent |
| skill-prd | `.claude/skills/skill-prd.md` | `.claude-plugin/.claude/skills/skill-prd.md` | equivalent | same skill name and intent |
| skill-quick | `.claude/skills/skill-quick.md` | `.claude-plugin/.claude/skills/skill-quick.md` | equivalent | same skill name and intent |
| skill-quick-review | `.claude/skills/skill-quick-review.md` | `.claude-plugin/.claude/skills/skill-quick-review.md` | equivalent | same skill name and intent |
| skill-resume | `.claude/skills/skill-resume.md` | `.claude-plugin/.claude/skills/skill-resume.md` | equivalent | same skill name and intent |
| skill-resume-enhanced | `.claude/skills/skill-resume-enhanced.md` | `.claude-plugin/.claude/skills/skill-resume-enhanced.md` | equivalent | same skill name and intent |
| skill-rollback | `.claude/skills/skill-rollback.md` | `.claude-plugin/.claude/skills/skill-rollback.md` | equivalent | same skill name and intent |
| skill-security-audit | `.claude/skills/skill-security-audit.md` | `.claude-plugin/.claude/skills/skill-security-audit.md` | equivalent | same skill name and intent |
| skill-security-framing | `.claude/skills/skill-security-framing.md` | `.claude-plugin/.claude/skills/skill-security-framing.md` | equivalent | same skill name and intent |
| skill-ship | `.claude/skills/skill-ship.md` | `.claude-plugin/.claude/skills/skill-ship.md` | equivalent | same skill name and intent |
| skill-status | `.claude/skills/skill-status.md` | `.claude-plugin/.claude/skills/skill-status.md` | equivalent | same skill name and intent |
| skill-task-management | `.claude/skills/skill-task-management.md` | `.claude-plugin/.claude/skills/skill-task-management.md` | equivalent | same skill name and intent |
| skill-task-management-v2 | `.claude/skills/skill-task-management-v2.md` | `.claude-plugin/.claude/skills/skill-task-management-v2.md` | equivalent | same skill name and intent |
| skill-tdd | `.claude/skills/skill-tdd.md` | `.claude-plugin/.claude/skills/skill-tdd.md` | equivalent | same skill name and intent |
| skill-thought-partner | `.claude/skills/skill-thought-partner.md` | `.claude-plugin/.claude/skills/skill-thought-partner.md` | equivalent | same skill name and intent |
| skill-validate | `.claude/skills/skill-validate.md` | `.claude-plugin/.claude/skills/skill-validate.md` | equivalent | same skill name and intent |
| skill-verify | `.claude/skills/skill-verify.md` | `.claude-plugin/.claude/skills/skill-verify.md` | equivalent | same skill name and intent |
| skill-visual-feedback | `.claude/skills/skill-visual-feedback.md` | `.claude-plugin/.claude/skills/skill-visual-feedback.md` | equivalent | same skill name and intent |
| skill-writing-plans | `.claude/skills/skill-writing-plans.md` | `.claude-plugin/.claude/skills/skill-writing-plans.md` | equivalent | same skill name and intent |
| sys-configure | `.claude/skills/sys-configure.md` | `.claude-plugin/.claude/skills/sys-configure.md` | equivalent | same skill name and intent |
| (go-only) skill-persona | `N/A` | `.claude-plugin/.claude/skills/skill-persona.md` | intentional-diff | go-only additive persona capability |

## 缺口整改建议（按优先级）

1. `P0`: 补齐 `doctor` 与 `scheduler` 相关命令入口（`/mp:doctor`, `/mp:schedule`, `/mp:scheduler`），确保运维可用性不回退。  
2. `P0`: 明确 `sentinel` 去留；若保留能力，迁移到 go 运行时并提供命令入口。  
3. `P1`: 为已存在 skill 但缺命令入口的能力补齐命令包装（`/mp:parallel`, `/mp:spec`）。  
4. `P1`: 评估 `claw/skill-claw` 是否迁移；若不迁移，在文档中显式声明退役与替代方案。

## Parity 结论

- 语义层面：大部分 commands/skills 已对齐，但 `main` 仍有关键能力未在 `go` 达到同等级入口与行为。  
- 文件层面：共享能力已建立可追踪的一对一文件映射；非双射部分已补充显式映射与整改路径。  
- 当前判定：`commands/skills = partial parity`。
