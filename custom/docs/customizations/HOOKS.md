# HOOKS 运维说明

本文面向插件运维/维护人员，说明当前 Claude Code 生命周期 hooks 的注册位置、触发时机、决策逻辑、落盘副作用与排障方式。

> 说明：这里的 hooks 指 **Claude Code lifecycle hooks**，不是 git hooks。
> 当前文档以仓库当前分支的代码实现为准；统一入口为 `mp hook --event ... --json`。

---

## 1. 注册位置与执行入口

### 1.1 Claude Code hook 注册表

文件：`.claude-plugin/hooks.json`

当前注册了 9 个事件：

- `SessionStart`
- `EnterPlanMode`
- `UserPromptSubmit`
- `PreToolUse`
- `PostToolUse`
- `WorktreeCreate`
- `WorktreeRemove`
- `Stop`
- `SubagentStop`

每个事件都调用同一个 Go 入口：

```text
${CLAUDE_PLUGIN_ROOT}/bin/mp hook --event <EventName> --dir "${CLAUDE_PROJECT_DIR}" --json
```

这意味着：

- hook 调度是集中式的，不是每个事件一个独立脚本
- 行为统一由 `internal/hooks/handler.go` 分发
- 返回值统一是 JSON hook result，核心字段是：
  - `decision`: `allow` / `block`
  - `reason`
  - `remediation`
  - `metadata`

### 1.2 统一分发入口

文件：`internal/hooks/handler.go`

分发关系：

- `SessionStart` -> `SessionStartData(...)`
- `EnterPlanMode` -> `handleEnterPlanMode(...)`
- `UserPromptSubmit` -> 内联判定 + `resolveModelRouting(...)`
- `PreToolUse` -> `PreToolUse(...)`
- `PostToolUse` -> `PostToolUse(...)`
- `WorktreeCreate` / `WorktreeRemove` -> `appendWorktreeEvent(...)`
- `Stop` / `SubagentStop` -> `StopDecision(...)`

所有 hook 事件都会先写一条 autosync raw event，用于运行期观测；其中 `PreToolUse`、`PostToolUse`、`Stop` 还会额外写对应的细分 raw event。

---

## 2. Hook 一览表

| Hook | 触发时机 | 默认动作 | 主要用途 | 可能 block |
|---|---|---|---|---|
| `SessionStart` | 会话开始 | `allow` | 注入项目/策略摘要元数据 | 否 |
| `EnterPlanMode` | 进入 Plan 模式前 | `allow` / `block` | 强制要求显式 `/mp:plan` 意图 | 是 |
| `UserPromptSubmit` | 用户提交 prompt 时 | `allow` / `block` | spec 命令上下文检查 + 路由元数据注入 | 是 |
| `PreToolUse` | 工具执行前 | `allow` / `block` | 写入边界校验 + 安全策略拦截 | 是 |
| `PostToolUse` | 工具执行后 | `allow` / `block` | FAQ 学习循环 + track 后处理 | 是 |
| `WorktreeCreate` | 创建 worktree 时 | `allow` | 记录 worktree 生命周期事件 | 否 |
| `WorktreeRemove` | 删除 worktree 时 | `allow` | 记录 worktree 生命周期事件 | 否 |
| `Stop` | 主会话结束前 | `allow` / `block` | 阻止未完成必需流程时直接退出 | 是 |
| `SubagentStop` | 子代理结束前 | `allow` / `block` | 同上，防止子代理提前退出 | 是 |

---

## 3. 各 Hook 详细行为

### 3.1 `SessionStart`

**触发场景**

- Claude Code 新会话开始时

**实现位置**

- 注册：`.claude-plugin/hooks.json`
- 处理：`internal/hooks/handler.go`
- 数据构造：`internal/hooks/session_start.go`

**做什么**

- 从 `.multipowers/` 读取并摘要以下上下文文件：
  - `product.md`
  - `product-guidelines.md`
  - `tech-stack.md`
  - `workflow.md`
  - `CLAUDE.md`
- 如果 policy 可加载，附带：
  - `policy_version`
  - `policy_checksum`
  - `workflows_configured`
  - `agents_configured`
- 返回 `track_status=unknown` 占位信息

**运维意义**

- 给会话提供“热启动”上下文
- 出现“为什么会话没拿到项目上下文”时，优先检查 `.multipowers/` 文件是否完整
- 出现“policy 元数据缺失”时，优先检查 policy 编译/加载链路

**是否会 block**

- 不会，固定 `allow`

---

### 3.2 `EnterPlanMode`

**触发场景**

- 会话准备进入 Plan 模式时

**实现位置**

- `internal/hooks/handler.go`

**做什么**

- 读取本次事件的 `tool_input.prompt`
- 只有 prompt 明确以 `/mp:plan` 开头时才允许进入 Plan 模式
- 否则返回 `block`

**block 条件**

- 当前意图不是 `/mp:plan`

**返回结果**

- `allow`: `plan-mode intent confirmed via /mp:plan`
- `block`: `plan mode requires explicit /mp:plan intent`
- `metadata.required_command=/mp:plan`

**运维意义**

- 防止用户/代理在没有显式规划意图时误进入 Plan 模式
- 如果用户反馈“进不了 Plan 模式”，先看 prompt 是否真的是 `/mp:plan ...`

---

### 3.3 `UserPromptSubmit`

**触发场景**

- 用户每次提交 prompt 时

**实现位置**

- `internal/hooks/handler.go`

**识别范围**

只对 spec-driven 命令做特殊处理：

- `/mp:plan`
- `/mp:discover`
- `/mp:define`
- `/mp:develop`
- `/mp:deliver`
- `/mp:embrace`
- `/mp:review`
- `/mp:research`
- `/mp:debate`

**做什么**

#### A. 上下文完整性校验

如果是 spec 命令，但 `.multipowers` 上下文不完整：

- 直接 `block`
- 返回：
  - `action_required=run_init`
  - `recommended_command=/mp:init`
  - `resume_command=<原始 prompt>`
  - `missing_files=<缺失文件列表>`

这意味着：

- runtime 不会偷偷自动生成上下文
- 运维上应把它理解为“显式初始化门禁”

当前强依赖的上下文文件来自 `internal/context/RequiredFiles`：

- `.multipowers/product.md`
- `.multipowers/product-guidelines.md`
- `.multipowers/tech-stack.md`
- `.multipowers/workflow.md`
- `.multipowers/tracks/tracks.md`
- `.multipowers/CLAUDE.md`
- `.multipowers/context/runtime.json`

#### B. 路由与执行隔离元数据注入

如果上下文完整，则 `allow` 并返回：

- `model_routing`
  - 解析出的 workflow
  - requested model
  - executor kind/provider
  - executor profile
  - enforcement
  - fallback target
  - source ref
- `benchmark_code_intent`
  - `code_related`
  - `source`
  - `whitelist`
- `execution_isolation`
  - 是否启用隔离
  - 触发原因
  - whitelist match 情况
  - `branch_prefix`
  - `worktree_root`

**运维意义**

- 这是 spec 命令的第一层运行前诊断点
- 如果模型选择异常，先看这里返回的 `model_routing`
- 如果 external command 是否该进 worktree 有争议，先看这里的 `execution_isolation`

---

### 3.4 `PreToolUse`

**触发场景**

- Claude Code 即将执行工具前

**实现位置**

- `internal/hooks/pre_tool_use.go`

**做什么**

#### A. 安全策略硬拦截

如果 `tool_input.safety_block=true`：

- 直接 `block`
- 原因：`blocked by safety-critical policy`

#### B. 写入边界校验

仅对以下写类工具检查文件路径：

- `Write`
- `Edit`
- `MultiEdit`

检查逻辑：

- 读取 `tool_input.file_path`
- 调用 `fsboundary.ValidateWritePath(...)`
- 如果路径越界：
  - 追加一条 quality gate 记录
  - 直接 `block`

这是当前 hooks 中最关键的“防误写”门禁。

#### C. Autosync 高置信 proposal 警告

如果存在 **高置信但未解决** 的 autosync proposal：

- 不 block
- 继续 `allow`
- 但 metadata 中加：
  - `autosync_warning=unresolved_high_confidence_proposal`

proposal 判断来源：

- 文件：`.multipowers/policy/autosync/proposals.jsonl`
- 规则：
  - `confidence >= 0.95`
  - `status` 不属于已解决集合：
    - `auto-applied`
    - `manual-required`
    - `ignored`
    - `revoked`
    - `rolled-back`
    - `expired`

**运维意义**

- 如果出现“工具没执行，直接被拦”，先查是否是写入越界或 safety policy
- 如果工具还能执行但 UI 有 autosync warning，说明是高置信策略建议还没闭环

---

### 3.5 `PostToolUse`

**触发场景**

- 工具执行完成后

**实现位置**

- `internal/hooks/post_tool_use.go`

**做什么**

#### A. FAQ 学习循环

- 校验 `.multipowers/FAQ.md` 路径合法
- 根据 `evt.ToolName` 生成 FAQ 事件
- 去重后写入 `.multipowers/FAQ.md`

作用：

- 把工具失败/工具类型沉淀为后续规避知识

#### B. post-tool track 后处理

- 通过 `TrackCoordinator` 解析/复用一个 `post-tool` track
- 确保 canonical track artifacts 存在
- 更新 track metadata：
  - `title=Post Tool Track`（若为空）
  - `status=in_progress`
  - `execution_mode=hook`（若为空）
- 记录：
  - `last_command=<ToolName>`
  - `last_command_at=<UTC 时间>`
- 刷新 `.multipowers/tracks/tracks.md`

**关键点**

当前实现已经修正为：

- **只记录 command touch**
- **不再把工具名伪装成 implementation group**

也就是说，`PostToolUse` 不会再把：

- `current_group=post-tool`
- `completed_groups=[post-tool]`

这种“假进度”写进 runtime 状态。

**运维意义**

- 如果 FAQ 没更新，先查 `.multipowers/FAQ.md` 路径和权限
- 如果 track registry 没更新，查 `TrackCoordinator.ResolveTrack` / `EnsureArtifacts` / `UpdateRegistry`
- 如果用户说“工具执行把 group 状态弄乱了”，现在理论上不应再发生；应检查是否是旧版本插件缓存

---

### 3.6 `WorktreeCreate`

**触发场景**

- Claude Code / runtime 创建 worktree 时

**实现位置**

- `internal/hooks/handler.go`
- `internal/hooks/worktree_events.go`

**做什么**

- 把事件追加写入：
  - `.multipowers/temp/worktree-events.jsonl`
- 记录字段：
  - `id`
  - `event`
  - `timestamp`
  - `session_id`
  - `project_dir`
  - `metadata`（来自 `evt.ToolInput`）

**返回**

- 成功：`allow` + `worktree_events_log=<路径>`
- 失败：仍然 `allow`，但 metadata 会附带 `worktree_event_error`

**运维意义**

- 这是 worktree 生命周期审计日志
- 编排层如果要重建“谁什么时候建了 worktree”，看这个文件最直接

---

### 3.7 `WorktreeRemove`

**触发场景**

- worktree 被移除时

**行为**

- 与 `WorktreeCreate` 相同，只是 `event=WorktreeRemove`
- 同样写入 `.multipowers/temp/worktree-events.jsonl`

**运维意义**

- 用于核对 worktree 是否被异常清理、谁触发了清理、是否存在 dangling worktree 现象

---

### 3.8 `Stop`

**触发场景**

- 主会话结束前

**实现位置**

- `internal/hooks/stop.go`

**做什么**

- 先写 autosync raw event：`hook.stop`
- 根据 `ctxpkg.Complete(projectDir)` 判断是否允许停止

**allow 条件**

- `ctxpkg.Complete(projectDir)=true`
- 也就是必需上下文文件齐全，没有 mandatory checkpoint pending

**block 条件**

- `ctxpkg.Complete(projectDir)=false`
- 常见缺项包括 `.multipowers/context/runtime.json`、`.multipowers/tracks/tracks.md` 或其他必需上下文文件
- 会追加一条 quality gate：
  - source = 当前 stop 来源
  - message = `mandatory checkpoint pending`
  - class = `session-stop`

**返回**

- `allow`: `no mandatory checkpoint pending`
- `block`: `mandatory checkpoint pending`
- remediation: `finish required init/context workflow before stop`

**运维意义**

- 如果用户抱怨“结束不了会话”，通常说明初始化或上下文闭环还没完成
- 这是防止中途放弃导致 runtime 状态半残的最后门禁

---

### 3.9 `SubagentStop`

**触发场景**

- 子代理结束前

**行为**

- 与 `Stop` 完全相同
- 区别仅在于来源字段不同，便于区分是主会话还是 subagent 触发

**运维意义**

- 防止 subagent 在必需 checkpoint 未完成时提前退出
- 排查多代理流程异常终止时，要同时查 `Stop` 和 `SubagentStop`

---

## 4. 关键落盘文件

hook 会直接或间接影响这些文件：

- `.multipowers/FAQ.md`
  - `PostToolUse` FAQ 学习循环写入
- `.multipowers/tracks/tracks.md`
  - `PostToolUse` 会刷新 canonical tracks registry
- `.multipowers/tracks/<track_id>/metadata.json`
  - `PostToolUse` 会更新 track metadata / last_command
- `.multipowers/temp/worktree-events.jsonl`
  - `WorktreeCreate` / `WorktreeRemove` 事件日志
- `.multipowers/policy/autosync/proposals.jsonl`
  - `PreToolUse` 会读取，用于 unresolved high-confidence warning
- `.multipowers/policy/autosync/events.raw.<YYYY-MM-DD>.jsonl`
  - `handler.go` 会为所有 hooks 写通用 raw event，`PreToolUse` / `PostToolUse` / `Stop` 还会写细分 raw event
- `.multipowers/context/runtime.json`
  - `UserPromptSubmit` / `Stop` 会把它作为必需上下文的一部分进行校验
- `.multipowers/decisions/...`
  - `PreToolUse` / `Stop` 会写 quality gate 相关记录

---

## 5. 运维排障顺序

### 场景 A：spec 命令一提交就被拦截

先看：

1. `UserPromptSubmit` 返回的 `decision`
2. `metadata.action_required`
3. `metadata.missing_files`

如果是：

- `action_required=run_init`

说明缺的是 `.multipowers` 初始化，不是模型/路由问题。

### 场景 B：进入不了 Plan 模式

先看：

1. `EnterPlanMode` 返回值
2. prompt 是否以 `/mp:plan` 开头

### 场景 C：写文件类工具被拦截

先看：

1. `PreToolUse` 返回 reason
2. `tool_name` 是否是 `Write/Edit/MultiEdit`
3. `file_path` 是否越界
4. `.multipowers/decisions` 中是否有 `write-path` quality gate

### 场景 D：工具执行后 FAQ / track 没更新

先看：

1. `PostToolUse` 返回 reason
2. `.multipowers/FAQ.md` 是否存在且可写
3. `.multipowers/tracks/tracks.md` 是否存在且可写
4. `.multipowers/tracks/<track_id>/metadata.json` 是否更新了 `last_command`

### 场景 E：会话结束不了

先看：

1. `Stop` / `SubagentStop` 的返回 reason
2. `.multipowers` context 是否完整
3. 是否存在 mandatory checkpoint pending

### 场景 F：worktree 生命周期异常

先看：

1. `.multipowers/temp/worktree-events.jsonl`
2. 是否同时存在 `WorktreeCreate` 和对应 `WorktreeRemove`
3. `metadata` 中的路径/会话信息是否匹配

---

## 6. 手工调试方式

可以直接调用统一入口调试某个事件：

```bash
./.claude-plugin/bin/mp hook --event SessionStart --dir "$PWD" --json
./.claude-plugin/bin/mp hook --event EnterPlanMode --dir "$PWD" --json --prompt "/mp:plan fix runtime"
./.claude-plugin/bin/mp hook --event UserPromptSubmit --dir "$PWD" --json --prompt "/mp:develop add retry"
```

说明：

- `PreToolUse` / `PostToolUse` 这类事件通常依赖 `ToolName` / `ToolInput`，更适合通过测试或宿主环境重现
- 如果只是看当前注册情况，直接检查 `.claude-plugin/hooks.json`

---

## 7. 维护原则

运维上请遵守下面几条：

1. 不要把 lifecycle hook 当成 git hook 处理。
2. 不要为单个事件新增散落 shell 脚本；优先走 `mp hook --event ...` 统一入口。
3. `PreToolUse` 里的写入边界校验属于硬门禁，变更时必须审查安全影响。
4. `PostToolUse` 只能记录 command activity，不要重新引入“工具名伪装 group 状态”的写法。
5. `UserPromptSubmit` 的 `run_init` block 是产品契约，不要改回“静默自动初始化”。
6. worktree 生命周期事件必须继续保持可审计落盘。

---

## 8. 代码定位速查

- hook 注册表：`.claude-plugin/hooks.json`
- 分发入口：`internal/hooks/handler.go`
- SessionStart：`internal/hooks/session_start.go`
- PreToolUse：`internal/hooks/pre_tool_use.go`
- PostToolUse：`internal/hooks/post_tool_use.go`
- Stop/SubagentStop：`internal/hooks/stop.go`
- Worktree 事件落盘：`internal/hooks/worktree_events.go`
