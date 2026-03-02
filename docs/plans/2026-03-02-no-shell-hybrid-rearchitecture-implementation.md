# No-Shell Hybrid Runtime Re-architecture Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** 将当前 Go 黑盒 + Markdown 空壳改造成「Go 原子引擎 + Markdown 推理编排」架构，并一次性完成 Big-Bang 切换。

**Architecture:** Go 只暴露可组合的原子命令（state/validate/hook/route/test/coverage/status）；Markdown Skills 负责步骤化推理与分支决策；两者通过统一 JSON 契约连接。保留旧高层命令作为兼容 facade，但内部走原子层。

**Tech Stack:** Go (`internal/cli`, `internal/validation`, `internal/hooks`, `internal/providers`, `internal/workflows`, `internal/tracks`), Markdown skills (`.claude-plugin/.claude/skills`), docs (`docs/*`), Go tests (`go test ./...`).

---

## 0. 使用说明（给初级工程师）

- 你只需要按任务顺序执行，不需要自己做架构设计。
- 每个任务都有固定模板：`Why / What / Key Design / How`。
- 每个任务执行完就更新状态表并提交一次 commit。
- 任何一步失败，先修复当前任务，不要跳到下个任务。

### 状态枚举

- `TODO`: 未开始
- `IN_PROGRESS`: 进行中（同一时间只允许一个）
- `DONE`: 已完成且测试通过
- `BLOCKED`: 被阻塞，需记录原因

### 状态总表（执行时持续更新）

| Task ID | Name | Status | Owner | Last Update |
|---|---|---|---|---|
| T01 | Baseline Snapshot | DONE | Claude | 2026-03-02 |
| T02 | Atomic CLI Contract Tests | TODO | - | - |
| T03 | State Atomic Commands | TODO | - | - |
| T04 | Typed Validation | TODO | - | - |
| T05 | Test/Coverage Atomic Commands | TODO | - | - |
| T06 | Route Atomic Command | TODO | - | - |
| T07 | Hook Contract Normalization | TODO | - | - |
| T08 | De-stub Status Command | TODO | - | - |
| T09 | Rehydrate Markdown Skills | TODO | - | - |
| T10 | Legacy Command Facades | TODO | - | - |
| T11 | Documentation Alignment | TODO | - | - |
| T12 | End-to-End Verification | TODO | - | - |
| T13 | Final Review Gate | TODO | - | - |

---

## T01 - Baseline Snapshot

**Why**  
先记录现状，避免后续回归时不知道是“新问题”还是“原本就存在的问题”。

**What**  
采集测试基线、CLI 现有行为、工作区变更基线。

**Key Design**  
基线必须可复现：固定命令、固定输出位置、固定记录格式。

**How**
1. 运行：`go test ./...`。
2. 运行：`go test ./internal/cli -v`。
3. 运行：`./.claude-plugin/bin/mp status --dir . --json`。
4. 运行：`git status --short`。
5. 在本文件 `T01` 下追加“Baseline Result”小节，粘贴关键结果摘要。
6. 更新状态表：`T01 -> DONE`。
7. 提交：
```bash
git add docs/plans/2026-03-02-no-shell-hybrid-rearchitecture-implementation.md
git commit -m "chore(plan): record baseline before hybrid runtime refactor"
```

### Baseline Result (2026-03-02)

**Go Test Results:**
```
?   	github.com/gitbruce/claude-octopus/cmd/mp	[no test files]
?   	github.com/gitbruce/claude-octopus/cmd/mp-devx	[no test files]
ok  	github.com/gitbruce/claude-octopus/internal/app	(cached)
ok  	github.com/gitbruce/claude-octopus/internal/cli	0.004s
ok  	github.com/gitbruce/claude-octopus/internal/context	(cached)
ok  	github.com/gitbruce/claude-octopus/internal/devx	(cached)
ok  	github.com/gitbruce/claude-octopus/internal/faq	(cached)
ok  	github.com/gitbruce/claude-octopus/internal/fsboundary	(cached)
ok  	github.com/gitbruce/claude-octopus/internal/hooks	0.004s
ok  	github.com/gitbruce/claude-octopus/internal/modelroute	(cached)
ok  	github.com/gitbruce/claude-octopus/internal/providers	(cached)
ok  	github.com/gitbruce/claude-octopus/internal/tracks	(cached)
ok  	github.com/gitbruce/claude-octopus/internal/validation	(cached)
ok  	github.com/gitbruce/claude-octopus/internal/workflows	(cached)
ok  	github.com/gitbruce/claude-octopus/pkg/api	(cached)
```

**CLI Tests (verbose):**
- `TestInitRequiresExplicitPrompt` - PASS
- `TestInitWithPromptCreatesContext` - PASS

**mp status output:**
```json
{"status":"ok","data":{"context_complete":false}}
```

**Git status:** 2 untracked files (docs/architecture/difference.md, docs/plans/2026-03-02-no-shell-hybrid-rearchitecture-implementation.md)

---

## T02 - Atomic CLI Contract Tests

**Why**  
先写契约测试可以约束后续实现，避免命令接口反复变动。

**What**  
为新原子命令写失败测试，约束输出 JSON 最低字段：`status/action/error_code/message/data/remediation`。

**Key Design**  
测试优先验证“契约存在且稳定”，再验证业务值。

**How**
1. 修改文件：`internal/cli/root_test.go`。
2. 新增测试用例（table-driven）覆盖命令：
- `mp state get`
- `mp state set`
- `mp state update`
- `mp validate --type tdd-env`
- `mp hook --event PreToolUse`
- `mp route --intent develop`
- `mp test run`
- `mp coverage check`
3. 每个测试至少断言：退出码、JSON 可解析、`status` 非空。
4. 运行：`go test ./internal/cli -run Atomic -v`，确认失败（expected fail）。
5. 更新状态表：`T02 -> DONE`。
6. 提交：
```bash
git add internal/cli/root_test.go
git commit -m "test(cli): add atomic command contract tests"
```

---

## T03 - State Atomic Commands (`state get/set/update`)

**Why**  
Skill 编排需要稳定状态读写，否则无法做阶段推进和恢复。

**What**  
实现 `mp state get/set/update`，支持键值查询与批量更新。

**Key Design**  
复用 `internal/tracks/state.go` 文件锁机制，保证并发安全。

**How**
1. 新建：`internal/tracks/state_kv.go`，实现 KV 读写辅助函数。
2. 新建：`internal/tracks/state_kv_test.go`，覆盖：空状态、set 后 get、update 原子合并。
3. 修改：`internal/cli/root.go`，添加 `state` 子命令分发。
4. 修改：`internal/cli/root_test.go`，补充 state 命令测试通过断言。
5. 运行：
- `go test ./internal/tracks -v`
- `go test ./internal/cli -run State -v`
6. 更新状态表：`T03 -> DONE`。
7. 提交：
```bash
git add internal/tracks/state_kv.go internal/tracks/state_kv_test.go internal/cli/root.go internal/cli/root_test.go
git commit -m "feat(state): add atomic state get/set/update commands"
```

---

## T04 - Typed Validation (`validate --type`)

**Why**  
不同质量门禁必须可独立调用，Skill 才能做精细化分支。

**What**  
实现 `--type` 分发：`workspace/no-shell/tdd-env/test-run/coverage`。

**Key Design**  
验证逻辑集中在 validation 层，CLI 只做参数解析与响应封装。

**How**
1. 新建：`internal/validation/types.go`（类型调度入口）。
2. 新建：`internal/validation/types_test.go`（每种 type 的 happy path + invalid type）。
3. 修改：`internal/validation/gates.go`（必要复用逻辑抽取）。
4. 修改：`internal/cli/root.go`（接入 `validate --type`）。
5. 修改：`internal/cli/root_test.go`（validate 子命令测试）。
6. 运行：
- `go test ./internal/validation -v`
- `go test ./internal/cli -run Validate -v`
7. 更新状态表：`T04 -> DONE`。
8. 提交：
```bash
git add internal/validation/types.go internal/validation/types_test.go internal/validation/gates.go internal/cli/root.go internal/cli/root_test.go
git commit -m "feat(validate): add typed validation command contract"
```

---

## T05 - Atomic Test/Coverage Commands

**Why**  
TDD skill 需要原子测试结果和覆盖率结果作为分支依据。

**What**  
实现：`mp test run`、`mp coverage check`。

**Key Design**  
返回结构化结果（如 `passed`, `failed_tests`, `coverage_pct`），不要返回仅自然语言。

**How**
1. 新建：`internal/workflows/test_run.go`、`internal/workflows/test_run_test.go`。
2. 新建：`internal/workflows/coverage.go`、`internal/workflows/coverage_test.go`。
3. 修改：`internal/cli/root.go` 接入 `test` 与 `coverage` 子命令。
4. 修改：`internal/cli/root_test.go` 补充对应测试。
5. 运行：
- `go test ./internal/workflows -v`
- `go test ./internal/cli -run "Test|Coverage" -v`
6. 更新状态表：`T05 -> DONE`。
7. 提交：
```bash
git add internal/workflows/test_run.go internal/workflows/test_run_test.go internal/workflows/coverage.go internal/workflows/coverage_test.go internal/cli/root.go internal/cli/root_test.go
git commit -m "feat(cli): add atomic test and coverage commands"
```

---

## T06 - Atomic Route Command

**Why**  
Skill 需要可观察的 provider 路由决策，不应隐藏在黑盒流程里。

**What**  
实现 `mp route --intent <x> [--provider-policy <y>]`。

**Key Design**  
路由结果必须可解释：输出 `provider`, `model`, `role`, `reason`。

**How**
1. 新建：`internal/providers/router_intent.go`。
2. 新建：`internal/providers/router_intent_test.go`（discover/develop/deliver/fallback）。
3. 修改：`internal/providers/router.go`（必要公共逻辑复用）。
4. 修改：`internal/cli/root.go` 与 `internal/cli/root_test.go`。
5. 运行：
- `go test ./internal/providers -v`
- `go test ./internal/cli -run Route -v`
6. 更新状态表：`T06 -> DONE`。
7. 提交：
```bash
git add internal/providers/router_intent.go internal/providers/router_intent_test.go internal/providers/router.go internal/cli/root.go internal/cli/root_test.go
git commit -m "feat(route): expose atomic provider routing command"
```

---

## T07 - Hook Contract Normalization

**Why**  
Hook 决策要被 Skill 可靠消费，字段不一致会直接破坏编排。

**What**  
统一 `hook` CLI 输出契约，确保 allow/block 都有可处理字段。

**Key Design**  
`HookResult -> api.Response` 做一层稳定映射，避免直接泄漏内部结构。

**How**
1. 修改：`internal/hooks/handler_test.go`，先写失败用例覆盖 allow/block。
2. 修改：`internal/hooks/handler.go`，补齐 metadata 与 remediation。
3. 修改：`internal/cli/root.go`，规范 `mp hook` 的 JSON 响应。
4. 运行：
- `go test ./internal/hooks -v`
- `go test ./internal/cli -run Hook -v`
5. 更新状态表：`T07 -> DONE`。
6. 提交：
```bash
git add internal/hooks/handler.go internal/hooks/handler_test.go internal/cli/root.go internal/cli/root_test.go
git commit -m "feat(hook): normalize hook command contract for skill orchestration"
```

---

## T08 - De-stub `mp status`

**Why**  
当前 status 是占位符，会误导 Skill 和用户的决策。

**What**  
让 `mp status` 返回真实运行态：context/provider/validation/hook readiness。

**Key Design**  
状态聚合器独立成文件，便于测试和后续扩展。

**How**
1. 新建：`internal/cli/status.go`、`internal/cli/status_test.go`。
2. 修改：`internal/cli/root.go` 使用状态聚合器。
3. 修改：`internal/cli/root_test.go` 增加 status 合约断言。
4. 运行：`go test ./internal/cli -run Status -v`。
5. 更新状态表：`T08 -> DONE`。
6. 提交：
```bash
git add internal/cli/status.go internal/cli/status_test.go internal/cli/root.go internal/cli/root_test.go
git commit -m "feat(status): replace placeholder with runtime health data"
```

---

## T09 - Rehydrate Markdown Skills

**Why**  
目前 `flow-*` 与部分 skill 是空壳，缺少 LLM 推理链。

**What**  
将核心 skills 改为“指导推理 + 调用原子命令 + 基于 JSON 分支”。

**Key Design**  
Skill 里保留推理；Go 只提供确定性结果，不在 Go 内决定所有下一步。

**How**
1. 修改：
- `.claude-plugin/.claude/skills/flow-discover.md`
- `.claude-plugin/.claude/skills/flow-define.md`
- `.claude-plugin/.claude/skills/flow-develop.md`
- `.claude-plugin/.claude/skills/flow-deliver.md`
- `.claude-plugin/.claude/skills/skill-tdd.md`
- `.claude-plugin/.claude/skills/skill-validate.md`
- `.claude-plugin/.claude/skills/skill-status.md`
2. 每个文件都加入：阶段目标、命令调用、JSON 分支条件、失败补救。
3. 运行：
- `go test ./internal/validation -run NoShell -v`
- `rg -n "Thin wrapper skill|\.sh" .claude-plugin/.claude/skills`
4. 更新状态表：`T09 -> DONE`。
5. 提交：
```bash
git add .claude-plugin/.claude/skills/flow-discover.md .claude-plugin/.claude/skills/flow-define.md .claude-plugin/.claude/skills/flow-develop.md .claude-plugin/.claude/skills/flow-deliver.md .claude-plugin/.claude/skills/skill-tdd.md .claude-plugin/.claude/skills/skill-validate.md .claude-plugin/.claude/skills/skill-status.md
git commit -m "feat(skills): restore reasoning-driven skills with atomic mp commands"
```

---

## T10 - Legacy Command Facades

**Why**  
Big-Bang 后仍要保证旧命令可用，避免用户入口断裂。

**What**  
让 `discover/define/develop/deliver` 成为兼容 facade，内部调用原子层。

**Key Design**  
Facade 只编排，不复制逻辑，不回退到 shell。

**How**
1. 修改：`internal/cli/root.go`。
2. 修改：
- `internal/workflows/discover.go`
- `internal/workflows/define.go`
- `internal/workflows/develop.go`
- `internal/workflows/deliver.go`
3. 修改：`internal/cli/root_test.go`（旧命令行为回归测试）。
4. 运行：
- `go test ./internal/cli -v`
- `go test ./internal/workflows -v`
5. 更新状态表：`T10 -> DONE`。
6. 提交：
```bash
git add internal/cli/root.go internal/workflows/discover.go internal/workflows/define.go internal/workflows/develop.go internal/workflows/deliver.go internal/cli/root_test.go
git commit -m "refactor(workflows): keep legacy commands via atomic facades"
```

---

## T11 - Documentation Alignment

**Why**  
文档若仍是 shell/旧行为，会导致后续错误使用。

**What**  
更新 CLI、命令参考、workflow skills 文档到新架构。

**Key Design**  
文档以“原子命令 + JSON 契约 + skill 分支”为主线。

**How**
1. 修改：
- `docs/CLI-REFERENCE.md`
- `docs/COMMAND-REFERENCE.md`
- `docs/WORKFLOW-SKILLS.md`
- `RELEASE_NOTES_NO_SHELL_RUNTIME.md`
2. 删除过时 shell 调用示例，换成 `mp` 原子命令。
3. 运行：
- `rg -n "octo-state\.sh|doctor\.sh|quality-gate\.sh|Thin wrapper skill" docs .claude-plugin/.claude/skills`
4. 更新状态表：`T11 -> DONE`。
5. 提交：
```bash
git add docs/CLI-REFERENCE.md docs/COMMAND-REFERENCE.md docs/WORKFLOW-SKILLS.md RELEASE_NOTES_NO_SHELL_RUNTIME.md
git commit -m "docs: align docs with hybrid atomic runtime architecture"
```

---

## T12 - End-to-End Verification

**Why**  
在宣称完成前必须给出可重复验证证据。

**What**  
跑全量测试 + 关键命令烟测 + 禁用模式扫描。

**Key Design**  
验证命令与预期结果写入计划文档，形成审计轨迹。

**How**
1. 运行全量：`go test ./...`。
2. 运行关键命令：
- `./.claude-plugin/bin/mp state get --dir . --json`
- `./.claude-plugin/bin/mp validate --type no-shell --dir . --json`
- `./.claude-plugin/bin/mp route --intent develop --dir . --json`
- `./.claude-plugin/bin/mp test run --dir . --json`
- `./.claude-plugin/bin/mp coverage check --dir . --json`
- `./.claude-plugin/bin/mp status --dir . --json`
3. 运行扫描：`rg -n "scripts/(mp|mp-devx)|\.sh" .claude-plugin/.claude/skills internal docs`。
4. 在本文件追加“Verification Result”小节记录结果摘要。
5. 更新状态表：`T12 -> DONE`。
6. 提交：
```bash
git add docs/plans/2026-03-02-no-shell-hybrid-rearchitecture-implementation.md
git commit -m "chore: record verification evidence for hybrid runtime migration"
```

---

## T13 - Final Review Gate

**Why**  
Big-Bang 改动面大，必须进行最终结构化 code review。

**What**  
做变更审计并请求正式评审。

**Key Design**  
评审标准聚焦：契约稳定、No-Shell 合规、skill 推理质量、兼容性。

**How**
1. 运行：
- `git log --oneline --decorate -n 30`
- `git diff --stat <base-branch>...HEAD`
2. 对照设计文档检查范围漂移。
3. 使用 `@requesting-code-review` 发起评审。
4. 如有缺陷，修复后回到 `T12` 重跑验证。
5. 更新状态表：`T13 -> DONE`。

---

## 完成定义（Definition of Done）

- 所有任务状态为 `DONE`。
- `go test ./...` 通过。
- 关键 `mp` 原子命令可用且返回统一 JSON 契约。
- 核心 skills 不再是 thin wrapper，具备推理步骤与分支逻辑。
- 文档不含陈旧 shell runtime 路径。
