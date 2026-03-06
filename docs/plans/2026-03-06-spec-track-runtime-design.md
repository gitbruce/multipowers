# Spec-Driven Track 产物统一化与 runtime 初始化修复设计

## 1. 背景与目标

本次优化聚焦三件事：

1. 修复 `/mp:init` 未生成 `.multipowers/context/runtime.json`，导致 pre-run 配置读取链路不完整。
2. 规范 spec-driven 产物目录：统一落盘到 `.multipowers/tracks/<id>/`，并在 `.multipowers/tracks/tracks.md` 做统计跟踪。
3. 建立统一 `<id>` 生命周期：`/mp:plan` 创建并激活；其他 spec-driven 命令默认复用 active track；若无 active track 则隐式创建。

补充约束（用户确认）：
- `tracks.md` 采用硬迁移到 `.multipowers/tracks/tracks.md`，旧路径不兼容读取。
- `implementation-plan.md` 必须包含最细颗粒任务的 `why/what/how/key design`。
- 每完成一个任务组必须更新状态并提交；复杂任务需使用 worktree，且由 LLM 自动计算并写入计划。

## 2. 方案选择

采用中心化 `TrackCoordinator`（推荐方案 2）：

- 统一处理 track 解析/创建、模板渲染、registry 统计更新、状态持久化。
- 所有 spec-driven 命令入口通过同一协调层，避免分散逻辑导致行为漂移。
- 命令和 skills 仅表达意图，不直接负责文件结构与命名规范。

## 3. 目录与文件规范

### 3.1 规范目录

- `.multipowers/context/runtime.json`
- `.multipowers/tracks/tracks.md`
- `.multipowers/tracks/<id>/intent.md`
- `.multipowers/tracks/<id>/design.md`
- `.multipowers/tracks/<id>/implementation-plan.md`
- `.multipowers/tracks/<id>/metadata.json`
- `.multipowers/tracks/<id>/index.md`

### 3.2 模板来源与职责

模板统一放到 `custom/templates/conductor/track/`：

- `intent.md.tpl`: 目标、成功标准、边界、约束。
- `design.md.tpl`: 技术架构为主，含 Scope Mapping（intent -> design decisions）。
- `implementation-plan.md.tpl`: 任务组级执行契约与细粒度任务模板。
- `metadata.json.tpl`: 机器可读执行状态与复杂度判定结果。
- `index.md.tpl`: track 内文档索引。

## 4. 数据流

### 4.1 `/mp:init`

1. 创建 `.multipowers/context/runtime.json` 默认骨架（`pre_run.enabled=false`）。
2. 初始化 `.multipowers/tracks/tracks.md`。
3. 若检测到旧路径 `.multipowers/tracks.md`：执行一次硬迁移；若新旧冲突则阻断并要求人工处理。

### 4.2 spec-driven 命令

命令范围：`plan/discover/define/develop/deliver/review/research/debate/embrace`。

1. 进入 `TrackCoordinator.ResolveTrack`：
   - `/mp:plan`：创建显式 track，设为 active。
   - 其他命令：优先复用 active；无 active 时隐式创建。
2. `TrackCoordinator.EnsureArtifacts`：按模板补齐缺失工件。
3. `TrackCoordinator.UpdateRegistry`：刷新 `tracks/tracks.md` 统计（总数、状态计数、active、最近更新时间）。

## 5. Implementation Plan 强约束

`implementation-plan.md` 必须包含：

- 任务组级字段：`Why/What/How/Key Design/Depends On/Verification/Done When/Rollback`。
- 每组完成后强制动作：
  1. 更新 `metadata.json` 执行态。
  2. 更新 `tracks/tracks.md` 统计。
  3. 执行 commit（规范消息：`track(<id>): group <gid> - <title>`）。

### 5.1 复杂度自动判定与 worktree 决策

LLM 在生成计划时自动计算 `complexity_score`：

- `+2` changed files >= 8
- `+2` touched modules >= 2
- `+2` migration/security critical path
- `+1` external integration
- `+1` groups >= 4
- `+1` estimated time > 2h

判定：`score >= 4 => worktree_required=true`。

`implementation-plan.md` 固定包含 `Execution Mode Decision` 段：
- Complexity Score
- Worktree Required (YES/NO)
- Rationale
- Enforcement（组完成后 commit + 状态更新）

## 6. 错误处理

- `worktree_required=true` 且当前不在 worktree：`blocked`。
- 任务组完成后未 commit：`blocked`。
- metadata 或 tracks 统计更新失败：`error`，禁止进入下一组。
- 模板缺失或渲染失败：`error`，返回缺失模板路径。

## 7. 测试策略

1. `internal/context`：`runtime.json` 生成与 init 回滚测试。
2. `internal/tracks/coordinator`：
   - active/implicit `<id>` 分配。
   - 工件模板渲染完整性。
   - registry 统计更新。
3. `internal/cli` + `internal/app`：spec-driven 入口是否统一进入 coordinator。
4. `internal/doctor`/`internal/validation`：只检查新路径 `tracks/tracks.md`。
5. 集成测试：任务组完成后的三件套校验（状态更新 + registry 更新 + commit 证据）。

## 8. 非目标

- 不在本次内实现旧路径长期双读兼容。
- 不在本次内重写全部旧 skills 文本；仅做最小必要对齐。
