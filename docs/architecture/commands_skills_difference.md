# main vs go 命令与技能差异（内容级重比对）

日期：2026-03-02  
比较分支：`main` vs `go`  
基线提交：`main=f6a815a326ec`，`go=8835e073834f`

## 关键说明

本次不再把“同名文件”直接视为等价，而是按**内容级语义**重比对：
- 同名但仅前缀替换（如 `/octo:` -> `/mp:`）且流程保持：可判 `equivalent`。
- 同名但主体改为运行时委派/薄包装：判 `partial`。
- main 存在而 go 无对等入口：判 `missing`。
- go 新增能力：判 `intentional-diff`。

用户关注样例已验证：
- `.claude/skills/extract-skill.md`（main 231 行完整指南）
- `.claude-plugin/.claude/skills/extract-skill.md`（go 9 行薄包装）
- 该映射应为 `partial`，不应标记为 `equivalent`。

## 结果概览

名称级交集规模：
- commands: main=46, go=41, shared=38, main-only=8, go-only=3
- skills: main=48, go=47, shared=46, main-only=2, go-only=1

内容级状态统计（全映射，含 only 行）：
- commands: `equivalent=32`, `partial=10`, `missing=5`, `intentional-diff=2`
- skills: `partial=46`, `missing=2`, `intentional-diff=1`

结论：
- command 层为“部分等价 + 明显缺口”；
- skill 层为“全面运行时重写（同名但几乎均非等价文本/流程）”；
- 总体判定：`commands/skills = partial parity`。

## 内容差异样例（Commands）

| command | main lines | go lines | add | del | why partial |
|---|---:|---:|---:|---:|---|
| embrace | 202 | 19 | 10 | 193 | main includes multi-step orchestration logic; go delegates to mp runtime wrapper |
| model-config | 343 | 307 | 106 | 142 | high textual delta indicates substantial content divergence |
| sys-setup | 245 | 229 | 6 | 22 | non-trivial content delta on options/contracts |
| multi | 187 | 187 | 13 | 13 | non-trivial content delta on options/contracts |
| km | 63 | 63 | 9 | 9 | non-trivial content delta on options/contracts |
| dev | 62 | 62 | 4 | 4 | non-trivial content delta on options/contracts |

## 内容差异样例（Skills）

| skill | main lines | go lines | add | del | wrapper | why partial |
|---|---:|---:|---:|---:|---:|---|
| skill-parallel-agents | 778 | 9 | 4 | 773 | 1 | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance |
| flow-deliver | 809 | 153 | 71 | 727 | 0 | substantial textual delta with retained same-name skill |
| flow-discover | 786 | 152 | 67 | 701 | 0 | substantial textual delta with retained same-name skill |
| skill-task-management-v2 | 683 | 9 | 4 | 678 | 1 | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance |
| skill-debate-integration | 663 | 9 | 4 | 658 | 1 | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance |
| flow-define | 748 | 227 | 122 | 643 | 0 | substantial textual delta with retained same-name skill |
| flow-develop | 723 | 169 | 80 | 634 | 0 | substantial textual delta with retained same-name skill |
| skill-debate | 597 | 9 | 4 | 592 | 1 | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance |
| skill-content-pipeline | 588 | 9 | 5 | 584 | 1 | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance |
| skill-audit | 572 | 9 | 5 | 568 | 1 | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance |
| flow-parallel | 643 | 245 | 139 | 537 | 0 | substantial textual delta with retained same-name skill |
| skill-validate | 586 | 141 | 91 | 536 | 0 | substantial textual delta with retained same-name skill |
| skill-meta-prompt | 518 | 9 | 5 | 514 | 1 | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance |
| skill-resume-enhanced | 494 | 9 | 4 | 489 | 1 | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance |

## 全量映射：Commands（main -> go）

| main name | main file | go target | status | evidence | remediation |
|---|---|---|---|---|---|
| brainstorm | `.claude/commands/brainstorm.md` | `.claude-plugin/.claude/commands/brainstorm.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| debate | `.claude/commands/debate.md` | `.claude-plugin/.claude/commands/debate.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| debug | `.claude/commands/debug.md` | `.claude-plugin/.claude/commands/debug.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| deck | `.claude/commands/deck.md` | `.claude-plugin/.claude/commands/deck.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| define | `.claude/commands/define.md` | `.claude-plugin/.claude/commands/define.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| deliver | `.claude/commands/deliver.md` | `.claude-plugin/.claude/commands/deliver.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| dev | `.claude/commands/dev.md` | `.claude-plugin/.claude/commands/dev.md` | `partial` | non-trivial content delta on options/contracts | validate option-level parity against main command contract |
| develop | `.claude/commands/develop.md` | `.claude-plugin/.claude/commands/develop.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| discover | `.claude/commands/discover.md` | `.claude-plugin/.claude/commands/discover.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| docs | `.claude/commands/docs.md` | `.claude-plugin/.claude/commands/docs.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| embrace | `.claude/commands/embrace.md` | `.claude-plugin/.claude/commands/embrace.md` | `partial` | main includes multi-step orchestration logic; go delegates to mp runtime wrapper | capture main behavior in runtime tests and docs |
| extract | `.claude/commands/extract.md` | `.claude-plugin/.claude/commands/extract.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| grasp | `.claude/commands/grasp.md` | `.claude-plugin/.claude/commands/grasp.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| ink | `.claude/commands/ink.md` | `.claude-plugin/.claude/commands/ink.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| issues | `.claude/commands/issues.md` | `.claude-plugin/.claude/commands/issues.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| km | `.claude/commands/km.md` | `.claude-plugin/.claude/commands/km.md` | `partial` | non-trivial content delta on options/contracts | validate option-level parity against main command contract |
| loop | `.claude/commands/loop.md` | `.claude-plugin/.claude/commands/loop.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| meta-prompt | `.claude/commands/meta-prompt.md` | `.claude-plugin/.claude/commands/meta-prompt.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| model-config | `.claude/commands/model-config.md` | `.claude-plugin/.claude/commands/model-config.md` | `partial` | high textual delta indicates substantial content divergence | review behavioral parity test coverage for this command |
| multi | `.claude/commands/multi.md` | `.claude-plugin/.claude/commands/multi.md` | `partial` | non-trivial content delta on options/contracts | validate option-level parity against main command contract |
| pipeline | `.claude/commands/pipeline.md` | `.claude-plugin/.claude/commands/pipeline.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| plan | `.claude/commands/plan.md` | `.claude-plugin/.claude/commands/plan.md` | `equivalent` | same planning structure; key differences are /octo->/mp and .claude->.multipowers path updates | none |
| prd | `.claude/commands/prd.md` | `.claude-plugin/.claude/commands/prd.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| prd-score | `.claude/commands/prd-score.md` | `.claude-plugin/.claude/commands/prd-score.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| probe | `.claude/commands/probe.md` | `.claude-plugin/.claude/commands/probe.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| quick | `.claude/commands/quick.md` | `.claude-plugin/.claude/commands/quick.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| research | `.claude/commands/research.md` | `.claude-plugin/.claude/commands/research.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| resume | `.claude/commands/resume.md` | `.claude-plugin/.claude/commands/resume.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| review | `.claude/commands/review.md` | `.claude-plugin/.claude/commands/review.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| rollback | `.claude/commands/rollback.md` | `.claude-plugin/.claude/commands/rollback.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| security | `.claude/commands/security.md` | `.claude-plugin/.claude/commands/security.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| setup | `.claude/commands/setup.md` | `.claude-plugin/.claude/commands/setup.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| ship | `.claude/commands/ship.md` | `.claude-plugin/.claude/commands/ship.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| status | `.claude/commands/status.md` | `.claude-plugin/.claude/commands/status.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| sys-setup | `.claude/commands/sys-setup.md` | `.claude-plugin/.claude/commands/sys-setup.md` | `partial` | non-trivial content delta on options/contracts | validate option-level parity against main command contract |
| tangle | `.claude/commands/tangle.md` | `.claude-plugin/.claude/commands/tangle.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| tdd | `.claude/commands/tdd.md` | `.claude-plugin/.claude/commands/tdd.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| validate | `.claude/commands/validate.md` | `.claude-plugin/.claude/commands/validate.md` | `equivalent` | prefix/path migration with equivalent behavior surface | none |
| claw | `.claude/commands/claw.md` | `N/A` | `missing` | main command has no same-name go command | port command or explicitly deprecate with replacement mapping |
| doctor | `.claude/commands/doctor.md` | `N/A` | `missing` | main command has no same-name go command | port command or explicitly deprecate with replacement mapping |
| octo | `.claude/commands/octo.md` | `.claude-plugin/.claude/commands/mp.md` | `partial` | root alias exists but smart router logic reduced | restore intent routing logic in go root command |
| parallel | `.claude/commands/parallel.md` | `.claude-plugin/.claude/skills/flow-parallel.md` | `partial` | flow skill exists but command entry missing | add /mp:parallel command wrapper |
| schedule | `.claude/commands/schedule.md` | `N/A` | `missing` | main command has no same-name go command | port command or explicitly deprecate with replacement mapping |
| scheduler | `.claude/commands/scheduler.md` | `N/A` | `missing` | main command has no same-name go command | port command or explicitly deprecate with replacement mapping |
| sentinel | `.claude/commands/sentinel.md` | `N/A` | `missing` | main command has no same-name go command | port command or explicitly deprecate with replacement mapping |
| spec | `.claude/commands/spec.md` | `.claude-plugin/.claude/skills/flow-spec.md` | `partial` | flow skill exists but command entry missing | add /mp:spec command wrapper |
| (go-only) init | `N/A` | `.claude-plugin/.claude/commands/init.md` | `intentional-diff` | go-only additive command | none |
| (go-only) mp | `.claude/commands/octo.md` | `.claude-plugin/.claude/commands/mp.md` | `partial` | root command counterpart exists but behavior differs | document reduced routing logic |
| (go-only) persona | `N/A` | `.claude-plugin/.claude/commands/persona.md` | `intentional-diff` | go-only additive command | none |

## 全量映射：Skills（main -> go）

| main name | main file | go target | status | evidence | remediation |
|---|---|---|---|---|---|
| extract-skill | `.claude/skills/extract-skill.md` | `.claude-plugin/.claude/skills/extract-skill.md` | `partial` | main is full reverse-engineering guide; go is thin wrapper calling mp status | map extract workflow to concrete mp subcommand/runtime path, not status stub |
| flow-define | `.claude/skills/flow-define.md` | `.claude-plugin/.claude/skills/flow-define.md` | `partial` | substantial textual delta with retained same-name skill | verify runtime contracts for this skill |
| flow-deliver | `.claude/skills/flow-deliver.md` | `.claude-plugin/.claude/skills/flow-deliver.md` | `partial` | substantial textual delta with retained same-name skill | verify runtime contracts for this skill |
| flow-develop | `.claude/skills/flow-develop.md` | `.claude-plugin/.claude/skills/flow-develop.md` | `partial` | substantial textual delta with retained same-name skill | verify runtime contracts for this skill |
| flow-discover | `.claude/skills/flow-discover.md` | `.claude-plugin/.claude/skills/flow-discover.md` | `partial` | substantial textual delta with retained same-name skill | verify runtime contracts for this skill |
| flow-parallel | `.claude/skills/flow-parallel.md` | `.claude-plugin/.claude/skills/flow-parallel.md` | `partial` | substantial textual delta with retained same-name skill | verify runtime contracts for this skill |
| flow-spec | `.claude/skills/flow-spec.md` | `.claude-plugin/.claude/skills/flow-spec.md` | `partial` | substantial textual delta with retained same-name skill | verify runtime contracts for this skill |
| skill-adversarial-security | `.claude/skills/skill-adversarial-security.md` | `.claude-plugin/.claude/skills/skill-adversarial-security.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-architecture | `.claude/skills/skill-architecture.md` | `.claude-plugin/.claude/skills/skill-architecture.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-audit | `.claude/skills/skill-audit.md` | `.claude-plugin/.claude/skills/skill-audit.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-code-review | `.claude/skills/skill-code-review.md` | `.claude-plugin/.claude/skills/skill-code-review.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-content-pipeline | `.claude/skills/skill-content-pipeline.md` | `.claude-plugin/.claude/skills/skill-content-pipeline.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-context-detection | `.claude/skills/skill-context-detection.md` | `.claude-plugin/.claude/skills/skill-context-detection.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-debate | `.claude/skills/skill-debate.md` | `.claude-plugin/.claude/skills/skill-debate.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-debate-integration | `.claude/skills/skill-debate-integration.md` | `.claude-plugin/.claude/skills/skill-debate-integration.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-debug | `.claude/skills/skill-debug.md` | `.claude-plugin/.claude/skills/skill-debug.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-decision-support | `.claude/skills/skill-decision-support.md` | `.claude-plugin/.claude/skills/skill-decision-support.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-deck | `.claude/skills/skill-deck.md` | `.claude-plugin/.claude/skills/skill-deck.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-deep-research | `.claude/skills/skill-deep-research.md` | `.claude-plugin/.claude/skills/skill-deep-research.md` | `partial` | substantial textual delta with retained same-name skill | verify runtime contracts for this skill |
| skill-doc-delivery | `.claude/skills/skill-doc-delivery.md` | `.claude-plugin/.claude/skills/skill-doc-delivery.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-finish-branch | `.claude/skills/skill-finish-branch.md` | `.claude-plugin/.claude/skills/skill-finish-branch.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-intent-contract | `.claude/skills/skill-intent-contract.md` | `.claude-plugin/.claude/skills/skill-intent-contract.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-issues | `.claude/skills/skill-issues.md` | `.claude-plugin/.claude/skills/skill-issues.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-iterative-loop | `.claude/skills/skill-iterative-loop.md` | `.claude-plugin/.claude/skills/skill-iterative-loop.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-knowledge-work | `.claude/skills/skill-knowledge-work.md` | `.claude-plugin/.claude/skills/skill-knowledge-work.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-meta-prompt | `.claude/skills/skill-meta-prompt.md` | `.claude-plugin/.claude/skills/skill-meta-prompt.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-parallel-agents | `.claude/skills/skill-parallel-agents.md` | `.claude-plugin/.claude/skills/skill-parallel-agents.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-prd | `.claude/skills/skill-prd.md` | `.claude-plugin/.claude/skills/skill-prd.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-quick | `.claude/skills/skill-quick.md` | `.claude-plugin/.claude/skills/skill-quick.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-quick-review | `.claude/skills/skill-quick-review.md` | `.claude-plugin/.claude/skills/skill-quick-review.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-resume | `.claude/skills/skill-resume.md` | `.claude-plugin/.claude/skills/skill-resume.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-resume-enhanced | `.claude/skills/skill-resume-enhanced.md` | `.claude-plugin/.claude/skills/skill-resume-enhanced.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-rollback | `.claude/skills/skill-rollback.md` | `.claude-plugin/.claude/skills/skill-rollback.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-security-audit | `.claude/skills/skill-security-audit.md` | `.claude-plugin/.claude/skills/skill-security-audit.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-security-framing | `.claude/skills/skill-security-framing.md` | `.claude-plugin/.claude/skills/skill-security-framing.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-ship | `.claude/skills/skill-ship.md` | `.claude-plugin/.claude/skills/skill-ship.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-status | `.claude/skills/skill-status.md` | `.claude-plugin/.claude/skills/skill-status.md` | `partial` | go skill re-authored as runtime health contract, not main dashboard flow text | document feature-level equivalence and missing dashboard steps |
| skill-task-management | `.claude/skills/skill-task-management.md` | `.claude-plugin/.claude/skills/skill-task-management.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-task-management-v2 | `.claude/skills/skill-task-management-v2.md` | `.claude-plugin/.claude/skills/skill-task-management-v2.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-tdd | `.claude/skills/skill-tdd.md` | `.claude-plugin/.claude/skills/skill-tdd.md` | `partial` | substantial textual delta with retained same-name skill | verify runtime contracts for this skill |
| skill-thought-partner | `.claude/skills/skill-thought-partner.md` | `.claude-plugin/.claude/skills/skill-thought-partner.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-validate | `.claude/skills/skill-validate.md` | `.claude-plugin/.claude/skills/skill-validate.md` | `partial` | substantial textual delta with retained same-name skill | verify runtime contracts for this skill |
| skill-verify | `.claude/skills/skill-verify.md` | `.claude-plugin/.claude/skills/skill-verify.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-visual-feedback | `.claude/skills/skill-visual-feedback.md` | `.claude-plugin/.claude/skills/skill-visual-feedback.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-writing-plans | `.claude/skills/skill-writing-plans.md` | `.claude-plugin/.claude/skills/skill-writing-plans.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| sys-configure | `.claude/skills/sys-configure.md` | `.claude-plugin/.claude/skills/sys-configure.md` | `partial` | go skill is thin wrapper to mp runtime; main skill contains detailed procedural guidance | trace wrapper to concrete runtime methods/tests and maintain behavior parity evidence |
| skill-claw | `.claude/skills/skill-claw.md` | `N/A` | `missing` | main skill has no same-name go counterpart | port skill or declare deprecation |
| skill-doctor | `.claude/skills/skill-doctor.md` | `N/A` | `missing` | main skill has no same-name go counterpart | port skill or declare deprecation |
| (go-only) skill-persona | `N/A` | `.claude-plugin/.claude/skills/skill-persona.md` | `intentional-diff` | go-only additive skill | none |

## 重点整改

1. `P0`：修正薄包装错路由。`extract-skill` 当前包装到 `mp status`，应映射到真正 extract 运行时入口。  
2. `P0`：补齐 main-only command 能力缺口（`claw/doctor/schedule/scheduler/sentinel`）或明确退役声明。  
3. `P1`：为 `octo -> mp` 补齐意图路由能力（当前为弱化版根命令）。  
4. `P1`：对所有 `partial` skills 增加“运行时方法/测试用例”证据链接，避免文档层面等价误判。

## Parity 结论

- “同名=等价”的旧判定在当前 go 代码上不成立。  
- main 到 go 的有效映射应以**语义承接 + 内容证据**为准。  
- 当前状态：`commands/skills` 仅达 `partial parity`。
