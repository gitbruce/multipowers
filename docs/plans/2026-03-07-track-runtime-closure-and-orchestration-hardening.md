# Track Runtime Closure and Orchestration Hardening Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Close the remaining spec-track runtime gaps, finish the still-necessary Go orchestration hardening work, and update stale plan/design docs so they match the current codebase before any new implementation claims are made.

**Architecture:** Start by re-baselining the old spec-track and orchestration optimization docs against current `go` branch behavior. Then fix the spec-track runtime root cause: command preflight is currently overloading `current_group` with command names, so group enforcement is only half-real. After that, finish the orchestration hardening work that still matters for runtime correctness and supportability: retry policy, deterministic retry execution, trace propagation, structured lifecycle logs, regression goldens, and verification evidence. Treat explainability / caching / visualization as follow-up enhancements unless the re-baseline step proves they are already partially implemented and must stay in the same epic.

**Tech Stack:** Go (`internal/tracks`, `internal/cli`, `internal/app`, `internal/validation`, `internal/orchestration`, `internal/policy`), Markdown docs in `docs/plans` and `custom/docs`, JSONL runtime artifacts under `.multipowers`, `go test`, `rg`, `git`.

---

## Preconditions

- Execute this plan in a dedicated git worktree created from the current `go` branch tip.
- Before touching code, run `go test ./... -count=1` and record the green baseline in the evidence doc created by Task 1.
- Do not trust the old status tables in `docs/plans/2026-03-04-go-orchestration-full-flow-migration-optimization-implementation.md`; re-baseline them from code first.
- Do not trust command names stored in `tracks.Metadata.CurrentGroup`; that is a known stale-doc / stale-implementation mismatch this plan must remove.

## Scope Decisions To Freeze In Task 1

1. `spec-track-runtime` is only considered complete when group progression is explicit, machine-callable, and enforced with commit + verification evidence.
2. `go-orchestration-optimization` is only considered complete for this wave when the runtime-hardening items are finished:
   - retry policy + deterministic retry loop,
   - trace propagation + structured logs,
   - golden regression coverage,
   - verification evidence + doc refresh.
3. `mp orchestrate explain`, result caching, and plan visualization should move to a follow-up backlog doc unless Task 1 finds code already landed that makes separating them misleading.

---

### Task 1: Re-baseline stale docs against current code

**Files:**
- Modify: `docs/plans/2026-03-06-spec-track-runtime-design.md`
- Modify: `docs/plans/2026-03-04-go-orchestration-full-flow-migration-optimization-implementation.md`
- Create: `docs/plans/2026-03-07-go-orchestration-ux-followups.md`
- Create: `docs/plans/evidence/track-runtime/2026-03-07-baseline-audit.md`

**Step 1: Capture proof that the old docs are stale**

Run:

```bash
go test ./... -count=1
rg -n 'CurrentGroup|CompletedGroups|LastCommitSHA|trace_id|backoff|jitter|cache_hit|mermaid' \
  internal custom docs -S
```

Expected:
- `go test` passes at current HEAD.
- Grep shows spec-track code still writes command names into `CurrentGroup`.
- Grep shows `trace_id`, retry backoff, cache, and mermaid items are either absent or only mentioned in docs.

**Step 2: Write the evidence snapshot**

Record the exact commands, exit codes, and 3-10 lines of relevant output in:

```markdown
# 2026-03-07 Baseline Audit

- Baseline commit: `<git rev-parse HEAD>`
- `go test ./... -count=1`: PASS
- Spec-track mismatch: `CurrentGroup` is populated by command names in `internal/cli/root.go` and `internal/hooks/post_tool_use.go`
- Optimization mismatch: old plan marks retry / traceability progress that is not present in runtime code
```

**Step 3: Refresh the old docs to match today’s baseline**

- Update `docs/plans/2026-03-06-spec-track-runtime-design.md` so it clearly states:
  - what is already implemented,
  - what is still open,
  - that command preflight and implementation-group lifecycle are different concerns.
- Update `docs/plans/2026-03-04-go-orchestration-full-flow-migration-optimization-implementation.md` so parent task statuses match the code, not the old assumptions.
- Create `docs/plans/2026-03-07-go-orchestration-ux-followups.md` to hold O06/O07/O08 if Task 1 confirms they are not part of the runtime-closure critical path.

**Step 4: Verify the refreshed docs are internally consistent**

Run:

```bash
rg -n 'CurrentGroup|trace_id|backoff|cache|mermaid|NOT_STARTED|COMPLETED' \
  docs/plans/2026-03-06-spec-track-runtime-design.md \
  docs/plans/2026-03-04-go-orchestration-full-flow-migration-optimization-implementation.md \
  docs/plans/2026-03-07-go-orchestration-ux-followups.md -S
```

Expected: every status/claim in those docs matches the code snapshot captured in Step 1.

**Step 5: Commit**

```bash
git add \
  docs/plans/2026-03-06-spec-track-runtime-design.md \
  docs/plans/2026-03-04-go-orchestration-full-flow-migration-optimization-implementation.md \
  docs/plans/2026-03-07-go-orchestration-ux-followups.md \
  docs/plans/evidence/track-runtime/2026-03-07-baseline-audit.md
git commit -m "docs(plan): rebaseline track runtime and orchestration hardening"
```

---

### Task 2: Add explicit track progress state instead of overloading command names

**Files:**
- Create: `internal/tracks/progress.go`
- Create: `internal/tracks/progress_test.go`
- Modify: `internal/tracks/metadata.go`
- Modify: `internal/tracks/metadata_test.go`

**Step 1: Write the failing tests**

```go
func TestRecordCommandTouchDoesNotMutateCurrentGroup(t *testing.T) {}

func TestStartGroupClearsPreviousEvidenceAndMarksInProgress(t *testing.T) {}

func TestCompleteGroupRequiresMatchingGroupAndCommitSHA(t *testing.T) {}
```

The tests should prove the new invariants:
- `last_command` / `last_command_at` track spec-command touches.
- `current_group` is reserved for actual implementation groups like `g1`, `g2`, ...
- starting a group clears prior `last_commit_sha` / `last_verified_at` evidence.
- completing a group requires an explicit commit SHA.

**Step 2: Run the tests to verify failure**

Run:

```bash
go test ./internal/tracks -run 'RecordCommandTouch|StartGroup|CompleteGroup|Metadata' -count=1
```

Expected: FAIL because `progress.go` does not exist and metadata lacks the new fields/helpers.

**Step 3: Write the minimal implementation**

Add a small domain API instead of sprinkling ad-hoc metadata writes around the codebase:

```go
type GroupStatus string

const (
    GroupStatusInProgress GroupStatus = "in_progress"
    GroupStatusCompleted  GroupStatus = "completed"
)

func RecordCommandTouch(projectDir, trackID, command string) error
func StartGroup(projectDir, trackID, groupID, executionMode string, worktreeRequired bool) error
func CompleteGroup(projectDir, trackID, groupID, commitSHA, verifiedAt string) error
```

Extend `Metadata` with only the fields required to separate command activity from group lifecycle:

```go
LastCommand   string      `json:"last_command,omitempty"`
LastCommandAt string      `json:"last_command_at,omitempty"`
GroupStatus   GroupStatus `json:"group_status,omitempty"`
```

Do **not** add a giant nested state machine. Keep the API tiny and deterministic.

**Step 4: Run the tests to verify they pass**

Run the same command from Step 2.

Expected: PASS.

**Step 5: Commit**

```bash
git add \
  internal/tracks/progress.go \
  internal/tracks/progress_test.go \
  internal/tracks/metadata.go \
  internal/tracks/metadata_test.go
git commit -m "feat(tracks): add explicit track progress state"
```

---

### Task 3: Rewire templates, coordinator, CLI, and hooks to the new progress model

**Files:**
- Modify: `internal/tracks/coordinator.go`
- Modify: `internal/tracks/template_renderer.go`
- Modify: `internal/tracks/template_renderer_test.go`
- Modify: `internal/cli/root.go`
- Modify: `internal/cli/root_test.go`
- Modify: `internal/hooks/post_tool_use.go`
- Modify: `internal/hooks/handler_test.go`
- Modify: `internal/app/spec_track_lifecycle_test.go`
- Modify: `custom/templates/conductor/track/intent.md.tpl`
- Modify: `custom/templates/conductor/track/design.md.tpl`
- Modify: `custom/templates/conductor/track/implementation-plan.md.tpl`
- Modify: `custom/templates/conductor/track/index.md.tpl`
- Modify: `custom/templates/conductor/track/metadata.json.tpl`

**Step 1: Write the failing tests**

```go
func TestTemplateRendererAllowsNoActiveGroup(t *testing.T) {}

func TestSpecCommandsWriteLastCommandWithoutFakeGroupProgress(t *testing.T) {}

func TestPostToolUseWritesLastCommandWithoutCompletingAGroup(t *testing.T) {}
```

The spec lifecycle integration test should assert:
- `plan` / `develop` preflight still creates artifacts,
- `metadata.last_command` is updated,
- `metadata.current_group` remains empty until an explicit group-start call,
- templates render “no active group” cleanly instead of faking a group name.

**Step 2: Run the tests to verify failure**

Run:

```bash
go test ./internal/tracks ./internal/cli ./internal/hooks ./internal/app \
  -run 'TemplateRenderer|SpecCommands|PostToolUse|SpecTrackLifecycle' -count=1
```

Expected: FAIL because the current code still writes command names into `CurrentGroup` and templates require a group unconditionally.

**Step 3: Write the minimal implementation**

- `TrackCoordinator.DefaultArtifactValues` should stop pretending that `command == current_group`.
- `prepareSpecTrack` and `PostToolUse` should call `RecordCommandTouch(...)`.
- Template rendering should treat `CurrentGroup` as optional and surface `LastCommand` separately.
- `implementation-plan.md.tpl` should still render `Why / What / How / Key Design / Execution Mode Decision`, but it must no longer claim a CLI command is a task group.

Example template direction:

```md
## Execution State
- Last Command: {{.LastCommand}}
- Current Group: {{if .CurrentGroup}}{{.CurrentGroup}}{{else}}(not started){{end}}
```

**Step 4: Run the tests to verify they pass**

Run the same command from Step 2.

Expected: PASS.

**Step 5: Commit**

```bash
git add \
  internal/tracks/coordinator.go \
  internal/tracks/template_renderer.go \
  internal/tracks/template_renderer_test.go \
  internal/cli/root.go \
  internal/cli/root_test.go \
  internal/hooks/post_tool_use.go \
  internal/hooks/handler_test.go \
  internal/app/spec_track_lifecycle_test.go \
  custom/templates/conductor/track/intent.md.tpl \
  custom/templates/conductor/track/design.md.tpl \
  custom/templates/conductor/track/implementation-plan.md.tpl \
  custom/templates/conductor/track/index.md.tpl \
  custom/templates/conductor/track/metadata.json.tpl
git commit -m "feat(tracks): separate command activity from group progress"
```

---

### Task 4: Add machine-callable group lifecycle commands and make enforcement real

**Files:**
- Create: `internal/tracks/worktree_check.go`
- Create: `internal/tracks/worktree_check_test.go`
- Modify: `internal/cli/root.go`
- Modify: `internal/cli/root_test.go`
- Modify: `internal/validation/gates.go`
- Modify: `internal/validation/group_enforcement_test.go`
- Modify: `internal/app/pipeline.go`
- Modify: `internal/app/pipeline_test.go`

**Step 1: Write the failing tests**

```go
func TestTrackGroupStartRequiresWorktreeWhenTrackDemandsIt(t *testing.T) {}

func TestTrackGroupCompleteRecordsCommitAndVerificationEvidence(t *testing.T) {}

func TestPipelineBlocksWhileActiveGroupIsMissingEvidence(t *testing.T) {}
```

Add CLI contract tests for a new atomic interface:

```text
mp track group-start --track-id <id> --group g1 --execution-mode worktree --json
mp track group-complete --track-id <id> --group g1 --commit-sha abc1234 --json
```

**Step 2: Run the tests to verify failure**

Run:

```bash
go test ./internal/tracks ./internal/validation ./internal/app ./internal/cli \
  -run 'TrackGroup|Worktree|Pipeline' -count=1
```

Expected: FAIL because there is no `track` atomic command and enforcement still depends on the `g[0-9]+` regex hack.

**Step 3: Write the minimal implementation**

- Add a `track` command family to `internal/cli/root.go` with `group-start` and `group-complete` subcommands.
- Use `StartGroup(...)` / `CompleteGroup(...)` from Task 2 instead of open-coding metadata edits.
- Add a small helper in `internal/tracks/worktree_check.go` that determines whether the current repo path is itself a git worktree checkout; use that helper only when `worktree_required=true`.
- Rewrite `EnsureTrackExecution(...)` so it checks explicit group state (`group_status == in_progress`) rather than inferring from a regex on `current_group`.

Minimal gate direction:

```go
if meta.GroupStatus == tracks.GroupStatusInProgress {
    if meta.LastCommitSHA == "" || meta.LastVerifiedAt == "" {
        return blocked(...)
    }
}
```

**Step 4: Run the tests to verify they pass**

Run the same command from Step 2.

Expected: PASS.

**Step 5: Commit**

```bash
git add \
  internal/tracks/worktree_check.go \
  internal/tracks/worktree_check_test.go \
  internal/cli/root.go \
  internal/cli/root_test.go \
  internal/validation/gates.go \
  internal/validation/group_enforcement_test.go \
  internal/app/pipeline.go \
  internal/app/pipeline_test.go
git commit -m "feat(tracks): enforce explicit group lifecycle evidence"
```

---

### Task 5: Add end-to-end closure coverage for spec-track runtime and refresh target-project docs

**Files:**
- Create: `internal/app/spec_track_group_lifecycle_test.go`
- Modify: `internal/app/spec_track_lifecycle_test.go`
- Modify: `custom/docs/target-project/README.md`
- Modify: `custom/docs/target-project/getting-started.md`
- Modify: `custom/docs/customizations/conductor-context.md`

**Step 1: Write the failing integration test**

```go
func TestSpecTrackGroupLifecycleRequiresCompletionEvidenceBetweenGroups(t *testing.T) {}
```

The scenario should be:
1. run init,
2. run `plan` preflight,
3. run `track group-start g1`,
4. verify the next spec pipeline call is blocked,
5. run `track group-complete g1 --commit-sha ...`,
6. verify the next spec pipeline call succeeds.

**Step 2: Run the tests to verify failure**

Run:

```bash
go test ./internal/app -run 'SpecTrack' -count=1
```

Expected: FAIL until the new command + gate behavior is fully wired.

**Step 3: Write the minimal implementation/docs**

- Fix any remaining orchestration between `RunSpecPipeline`, the new track commands, and metadata helpers.
- Update target-project docs so they describe:
  - canonical paths,
  - `runtime.json`,
  - explicit group lifecycle commands,
  - the fact that old `.multipowers/tracks.md` is not read.

**Step 4: Run the tests and doc checks**

Run:

```bash
go test ./internal/app -run 'SpecTrack' -count=1
rg -n '\.multipowers/tracks\.md|group-start|group-complete|runtime.json' custom/docs -S
```

Expected:
- integration tests PASS,
- docs mention only the canonical registry path,
- docs mention the new group lifecycle interface.

**Step 5: Commit**

```bash
git add \
  internal/app/spec_track_group_lifecycle_test.go \
  internal/app/spec_track_lifecycle_test.go \
  custom/docs/target-project/README.md \
  custom/docs/target-project/getting-started.md \
  custom/docs/customizations/conductor-context.md
git commit -m "test(tracks): add spec-track closure lifecycle coverage"
```

---

### Task 6: Add real retry policy fields to orchestration config and plan structures

**Files:**
- Modify: `internal/orchestration/plan_types.go`
- Modify: `internal/orchestration/merge.go`
- Modify: `internal/orchestration/load.go`
- Modify: `internal/orchestration/merge_test.go`
- Modify: `internal/orchestration/load_test.go`

**Step 1: Write the failing tests**

```go
func TestLoadConfigParsesRetryPolicyFields(t *testing.T) {}

func TestMergePreservesTaskLevelRetryOverrides(t *testing.T) {}
```

The tests should prove the plan really carries:
- `idempotent`
- `max_retries`
- `backoff_ms`
- `jitter_ratio`
- `retryable_codes`

**Step 2: Run the tests to verify failure**

Run:

```bash
go test ./internal/orchestration -run 'LoadConfig|Merge.*Retry' -count=1
```

Expected: FAIL because those fields are not actually present in the current structs/merge path.

**Step 3: Write the minimal implementation**

Extend `StepPlan` (or a nested retry policy sub-struct if cleaner) and thread the values through config load + merge.

Minimal direction:

```go
type RetryPolicy struct {
    Idempotent     bool
    MaxRetries     int
    BackoffMs      int
    JitterRatio    float64
    RetryableCodes []string
}
```

Keep defaults conservative: non-idempotent steps do not auto-retry.

**Step 4: Run the tests to verify they pass**

Run the same command from Step 2.

Expected: PASS.

**Step 5: Commit**

```bash
git add \
  internal/orchestration/plan_types.go \
  internal/orchestration/merge.go \
  internal/orchestration/load.go \
  internal/orchestration/merge_test.go \
  internal/orchestration/load_test.go
git commit -m "feat(orchestration): add retry policy fields to plans"
```

---

### Task 7: Implement deterministic retry execution in the executor

**Files:**
- Create: `internal/orchestration/retry.go`
- Modify: `internal/orchestration/executor.go`
- Modify: `internal/orchestration/result_types.go`
- Modify: `internal/orchestration/executor_test.go`

**Step 1: Write the failing tests**

```go
func TestExecutor_RetrySucceedsAfterTransientFailure(t *testing.T) {}

func TestExecutor_RetryStopsOnNonRetryableFailure(t *testing.T) {}

func TestExecutor_RetryHonorsContextCancellation(t *testing.T) {}
```

Use an injectable sleeper / clock so the tests are deterministic and do not rely on wall-clock sleeps.

**Step 2: Run the tests to verify failure**

Run:

```bash
go test ./internal/orchestration -run 'Executor_.*Retry' -count=1
```

Expected: FAIL because the executor currently dispatches once and returns.

**Step 3: Write the minimal implementation**

- Add a tiny retry controller in `retry.go`.
- Retry only if the step is idempotent and the error/status is retryable.
- Record attempt count and terminal failure reason in `StepResult`.
- Respect `ctx.Done()` immediately.

Sketch:

```go
type AttemptInfo struct {
    Count int
    LastError string
}
```

**Step 4: Run the tests to verify they pass**

Run the same command from Step 2.

Expected: PASS.

**Step 5: Commit**

```bash
git add \
  internal/orchestration/retry.go \
  internal/orchestration/executor.go \
  internal/orchestration/result_types.go \
  internal/orchestration/executor_test.go
git commit -m "feat(orchestration): add deterministic retry controller"
```

---

### Task 8: Add trace propagation and structured lifecycle logs

**Files:**
- Create: `internal/orchestration/log_writer.go`
- Create: `internal/orchestration/log_writer_test.go`
- Modify: `internal/orchestration/plan_types.go`
- Modify: `internal/orchestration/result_types.go`
- Modify: `internal/orchestration/planner.go`
- Modify: `internal/orchestration/executor.go`
- Modify: `internal/orchestration/events.go`
- Modify: `internal/orchestration/executor_test.go`
- Modify: `internal/orchestration/report_test.go`
- Modify: `internal/policy/dispatch.go`

**Step 1: Write the failing tests**

```go
func TestBuildPlanGeneratesStableTraceIDPerRun(t *testing.T) {}

func TestExecutorPropagatesTraceIDToStepsAndEvents(t *testing.T) {}

func TestLogWriterWritesStructuredLifecycleEvents(t *testing.T) {}
```

**Step 2: Run the tests to verify failure**

Run:

```bash
go test ./internal/orchestration -run 'TraceID|LogWriter|StructuredLifecycle' -count=1
```

Expected: FAIL because `trace_id` is not present in plan/result/event structs and events are not persisted as structured logs.

**Step 3: Write the minimal implementation**

- Add `TraceID` to plan metadata, execution result, step result, and event payloads.
- Generate it once per orchestration run.
- Persist JSONL step lifecycle events under `.multipowers/logs/` using the existing `logs_subdir` config.
- Never log raw prompt bodies; log hashes and lengths only.

Minimal event shape:

```go
type Event struct {
    TraceID string
    StepID  string
    Type    EventType
    Status  string
    Data    any
}
```

**Step 4: Run the tests to verify they pass**

Run the same command from Step 2.

Expected: PASS.

**Step 5: Commit**

```bash
git add \
  internal/orchestration/log_writer.go \
  internal/orchestration/log_writer_test.go \
  internal/orchestration/plan_types.go \
  internal/orchestration/result_types.go \
  internal/orchestration/planner.go \
  internal/orchestration/executor.go \
  internal/orchestration/events.go \
  internal/orchestration/executor_test.go \
  internal/orchestration/report_test.go \
  internal/policy/dispatch.go
git commit -m "feat(orchestration): add trace ids and structured lifecycle logs"
```

---

### Task 9: Add golden regressions for planner, report, and degraded paths

**Files:**
- Create: `internal/orchestration/testdata/golden/plan/discover.golden.json`
- Create: `internal/orchestration/testdata/golden/plan/define.golden.json`
- Create: `internal/orchestration/testdata/golden/plan/develop.golden.json`
- Create: `internal/orchestration/testdata/golden/plan/deliver.golden.json`
- Create: `internal/orchestration/testdata/golden/plan/debate.golden.json`
- Create: `internal/orchestration/testdata/golden/plan/embrace.golden.json`
- Create: `internal/orchestration/testdata/golden/report/*.golden.json`
- Create: `internal/orchestration/testdata/golden/degraded/*.golden.json`
- Modify: `internal/orchestration/planner_test.go`
- Modify: `internal/orchestration/synthesis_final_test.go`
- Modify: `internal/orchestration/e2e_test.go`

**Step 1: Write the failing golden harness**

```go
func TestPlannerGoldenSnapshots(t *testing.T) {}

func TestSynthesisReportGoldenSnapshots(t *testing.T) {}

func TestDegradedFallbackGoldenSnapshots(t *testing.T) {}
```

Require an explicit update flag/environment variable before rewriting golden files.

**Step 2: Run the tests to verify failure**

Run:

```bash
go test ./internal/orchestration -run 'Golden' -count=1
```

Expected: FAIL until the goldens and normalization helpers exist.

**Step 3: Write the minimal implementation**

- Normalize nondeterministic fields before compare.
- Snapshot:
  - full resolved plans,
  - synthesis reports,
  - degraded/fallback outputs.
- Reuse the same golden-file style already used by `internal/policy/compile_test.go`.

**Step 4: Run the tests to verify they pass**

Run the same command from Step 2.

Expected: PASS.

**Step 5: Commit**

```bash
git add \
  internal/orchestration/testdata/golden/plan \
  internal/orchestration/testdata/golden/report \
  internal/orchestration/testdata/golden/degraded \
  internal/orchestration/planner_test.go \
  internal/orchestration/synthesis_final_test.go \
  internal/orchestration/e2e_test.go
git commit -m "test(orchestration): add golden regressions for planner and degraded flows"
```

---

### Task 10: Capture verification evidence and refresh architecture / CLI docs

**Files:**
- Modify: `docs/plans/2026-03-04-go-orchestration-full-flow-migration-optimization-implementation.md`
- Create: `docs/plans/evidence/model-routing/2026-03-07-go-orchestration-hardening-verification.md`
- Modify: `docs/CLI-REFERENCE.md`
- Modify: `docs/ARCHITECTURE.md`

**Step 1: Capture the required verification commands**

Run:

```bash
go test ./internal/tracks ./internal/cli ./internal/hooks ./internal/app -count=1
go test ./internal/orchestration ./internal/policy ./internal/validation -count=1
go test ./... -count=1
```

Expected: all PASS.

**Step 2: Write the evidence doc**

Create `docs/plans/evidence/model-routing/2026-03-07-go-orchestration-hardening-verification.md` with:
- commit SHA,
- exact commands,
- exit codes,
- timestamps,
- first output lines,
- notes on spec-track closure behavior.

**Step 3: Refresh the user-facing docs**

- Update `docs/CLI-REFERENCE.md` for the new `mp track group-start` / `group-complete` commands and any orchestration log output contracts.
- Update `docs/ARCHITECTURE.md` to describe:
  - explicit track group lifecycle,
  - retry controller,
  - trace propagation,
  - structured orchestration logs,
  - the fact that O06/O07/O08 live in the follow-up backlog doc if Task 1 moved them out.

**Step 4: Verify docs and evidence are aligned**

Run:

```bash
rg -n 'group-start|group-complete|trace_id|structured log|golden|follow-up' \
  docs/CLI-REFERENCE.md \
  docs/ARCHITECTURE.md \
  docs/plans/2026-03-04-go-orchestration-full-flow-migration-optimization-implementation.md \
  docs/plans/evidence/model-routing/2026-03-07-go-orchestration-hardening-verification.md -S
```

Expected: every new runtime behavior has matching docs/evidence.

**Step 5: Commit**

```bash
git add \
  docs/plans/2026-03-04-go-orchestration-full-flow-migration-optimization-implementation.md \
  docs/plans/evidence/model-routing/2026-03-07-go-orchestration-hardening-verification.md \
  docs/CLI-REFERENCE.md \
  docs/ARCHITECTURE.md
git commit -m "docs: capture track runtime closure and orchestration hardening evidence"
```

---

## Final Verification Checklist

- `go test ./internal/tracks ./internal/cli ./internal/hooks ./internal/app -count=1`
- `go test ./internal/orchestration ./internal/policy ./internal/validation -count=1`
- `go test ./... -count=1`
- `rg -n '\.multipowers/tracks\.md' internal custom docs -S` shows only historical/non-runtime notes.
- `rg -n 'CurrentGroup = command|post-tool.*CurrentGroup' internal -S` returns no runtime writes that overload command names as group IDs.
- `rg -n 'trace_id|group-start|group-complete|golden' docs -S` returns the refreshed docs/evidence.

## Handoff Notes

- If Task 1 shows any part of O06/O07/O08 already landed in code, do **not** delete it; instead, update the follow-up doc so it reflects the true residual work.
- Keep the implementation DRY and boring. The point of this plan is to remove stale-doc drift and make the runtime contracts explicit, not to introduce a new orchestration framework.
- Do not batch multiple tasks into one commit.
