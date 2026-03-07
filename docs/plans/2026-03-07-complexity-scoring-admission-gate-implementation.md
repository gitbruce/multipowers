# Complexity Scoring Admission Gate Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a unified spec admission gate so high-complexity spec-driven commands require canonical planning artifacts and worktree execution before code changes proceed.

**Architecture:** Keep the current spec command entrypoint in `internal/app/pipeline.go`, but replace the narrow complexity check with a richer admission result. Reuse `internal/tracks` for canonical artifact inspection and interrupted-context persistence, and let `internal/cli/root.go` own track allocation/reuse plus response enrichment. Execute implementation in a dedicated linked worktree because the change spans scoring, validation, pipeline, CLI, and track lifecycle behavior.

**Tech Stack:** Go, existing `internal/tracks`, `internal/validation`, `internal/app`, `internal/cli`, canonical track templates under `.multipowers/tracks/<id>/`, `go test`

---

## Execution Mode Decision

- Complexity Score: 7
- Worktree Required: YES
- Why: This feature changes cross-module runtime behavior across scoring, validation, pipeline, CLI, metadata persistence, and lifecycle tests.
- Enforcement: Implement only inside a linked git worktree; after each task run the listed tests and commit with `track(complexity-scoring-gate-20260307): group <gid> - <title>`.

## Before Starting

- Create and switch to a dedicated worktree from `go`.
- Recommended branch: `complexity-scoring-admission-gate-20260307`
- Recommended worktree path: `.worktrees/complexity-scoring-admission-gate-20260307`
- Confirm the active track is `complexity-scoring-gate-20260307` before touching runtime code.

Example:

```bash
git worktree add .worktrees/complexity-scoring-admission-gate-20260307 -b complexity-scoring-admission-gate-20260307 go
cd .worktrees/complexity-scoring-admission-gate-20260307
```

---

### Task 1: Add Canonical Planning Artifact Inspection

**Why:** The new gate must distinguish “track exists” from “planning is actually complete under the canonical artifact contract”.

**What:** Create a reusable checker that validates the canonical artifact set and rejects old `spec.md` / `plan.md` as sufficient planning evidence.

**How:** Add a small `internal/tracks` helper that inspects `.multipowers/tracks/<id>/` and returns missing artifact names in stable order.

**Key Design:** Planning completeness is a track-level concept; do not mix it into `internal/context/RequiredFiles`.

**Files:**
- Create: `internal/tracks/artifacts.go`
- Create: `internal/tracks/artifacts_test.go`
- Reference: `internal/tracks/coordinator.go`
- Reference: `custom/templates/conductor/track/index.md.tpl`

**Step 1: Write the failing tests**

```go
func TestCanonicalArtifacts_AllPresent(t *testing.T) {}
func TestCanonicalArtifacts_MissingDesignAndPlan(t *testing.T) {}
func TestCanonicalArtifacts_LegacySpecPlanAreNotSufficient(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/tracks -run 'TestCanonicalArtifacts_' -count=1`
Expected: FAIL because the artifact inspection helper does not exist yet.

**Step 3: Write minimal implementation**

- Add `CanonicalArtifacts()` returning the five required filenames.
- Add `CheckCanonicalArtifacts(projectDir, trackID)` returning `complete + missing[]`.
- Keep missing artifact ordering deterministic for stable test output.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/tracks -run 'TestCanonicalArtifacts_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/tracks/artifacts.go internal/tracks/artifacts_test.go
git commit -m "track(complexity-scoring-gate-20260307): group g1 - add canonical planning artifact inspection"
```

---

### Task 2: Extend Complexity Scoring for Admission Semantics

**Why:** The runtime needs to express “requires planning” separately from the older single-bit “requires worktree” model.

**What:** Extend the scoring contract so admission-time complexity can explain whether planning is required and why.

**How:** Evolve the decision type in `internal/tracks/scoring.go` and cover the new behavior with unit tests instead of encoding policy in string messages.

**Key Design:** Admission scoring stays lightweight and deterministic; plan-time refinement can add richer evidence later.

**Files:**
- Modify: `internal/tracks/scoring.go`
- Modify: `internal/tracks/complexity_test.go`
- Optional Create: `internal/tracks/scoring_admission_test.go`

**Step 1: Write the failing tests**

```go
func TestCalculateComplexity_HighIntentRequiresPlanning(t *testing.T) {}
func TestCalculateComplexity_LowIntentDoesNotRequirePlanning(t *testing.T) {}
func TestCalculateComplexity_ReportsAdmissionSource(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/tracks -run 'TestCalculateComplexity_' -count=1`
Expected: FAIL because the decision type does not yet expose planning-specific semantics.

**Step 3: Write minimal implementation**

- Extend `ComplexityDecision` with admission-facing fields such as `RequiresPlanning` and `Source`.
- Keep existing threshold behavior stable: score `>= 4` still implies worktree-required for high-complexity execution.
- Make rationale explicit enough for later response metadata.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/tracks -run 'TestCalculateComplexity_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/tracks/scoring.go internal/tracks/complexity_test.go internal/tracks/scoring_admission_test.go
git commit -m "track(complexity-scoring-gate-20260307): group g2 - extend admission complexity semantics"
```

---

### Task 3: Introduce a Structured Admission Gate Result

**Why:** `Valid/Reason` is too weak to distinguish missing planning from missing worktree readiness.

**What:** Add a richer admission result and upgrade validation logic to return structured fields: track id, score, missing artifacts, requires planning, requires worktree, remediation.

**How:** Create a dedicated admission helper in `internal/validation`, then let `gates.go` delegate complexity-specific checks to it.

**Key Design:** Keep global context checking (`/mp:init`) separate from track admission checking.

**Files:**
- Create: `internal/validation/admission.go`
- Create: `internal/validation/admission_test.go`
- Modify: `internal/validation/gates.go`
- Modify: `internal/validation/gates_test.go`
- Reference: `internal/context/requirements.go`
- Reference: `internal/tracks/metadata.go`

**Step 1: Write the failing tests**

```go
func TestAdmission_HighComplexityWithoutTrackRequiresPlan(t *testing.T) {}
func TestAdmission_HighComplexityWithMissingArtifactsRequiresPlan(t *testing.T) {}
func TestAdmission_HighComplexityPlannedButOutsideWorktreeRequiresIsolation(t *testing.T) {}
func TestAdmission_LowComplexityAllowsExecution(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/validation -run 'TestAdmission_' -count=1`
Expected: FAIL because the admission gate and result contract do not exist yet.

**Step 3: Write minimal implementation**

- Add `AdmissionResult` with stable fields used by pipeline and CLI.
- Rework complexity gating to:
  - allow low complexity,
  - block high complexity without active track,
  - block high complexity with incomplete canonical artifacts,
  - block high complexity when planning is complete but current checkout is not a linked worktree.
- Keep `/mp:init` failures on the existing context path.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/validation -run 'TestAdmission_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/validation/admission.go internal/validation/admission_test.go internal/validation/gates.go internal/validation/gates_test.go
git commit -m "track(complexity-scoring-gate-20260307): group g3 - add structured spec admission gate"
```

---

### Task 4: Wire Structured Admission into the Spec Pipeline

**Why:** The pipeline must turn admission results into a stable response contract instead of burying semantics inside message strings.

**What:** Update `RunSpecPipeline` so planning-required and worktree-required blocks emit explicit response data and remediation.

**How:** Consume `AdmissionResult` in `internal/app/pipeline.go`, preserving existing `/mp:init` behavior while enriching blocked responses for high-complexity flows.

**Key Design:** `pipeline` is the single runtime entrypoint for spec-driven commands; it should not persist track state itself.

**Files:**
- Modify: `internal/app/pipeline.go`
- Modify: `internal/app/pipeline_test.go`
- Create: `internal/app/spec_admission_pipeline_test.go`
- Reference: `internal/validation/admission.go`

**Step 1: Write the failing tests**

```go
func TestRunSpecPipeline_HighComplexityMissingPlanReturnsStructuredBlock(t *testing.T) {}
func TestRunSpecPipeline_HighComplexityMissingWorktreeReturnsStructuredBlock(t *testing.T) {}
func TestRunSpecPipeline_LowComplexityContinues(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/app -run 'TestRunSpecPipeline_.*(Plan|Worktree|Continues)' -count=1`
Expected: FAIL because the current pipeline does not expose the new admission fields.

**Step 3: Write minimal implementation**

- Replace direct `EnsureComplexityGate` string handling with `AdmissionResult` handling.
- Populate response `Data` with `track_id`, `complexity_score`, `requires_planning`, `requires_worktree`, `missing_artifacts`, and resume hints when present.
- Preserve current blocked semantics for missing global context.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/app -run 'TestRunSpecPipeline_.*(Plan|Worktree|Continues)' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/app/pipeline.go internal/app/pipeline_test.go internal/app/spec_admission_pipeline_test.go
git commit -m "track(complexity-scoring-gate-20260307): group g4 - wire structured admission into pipeline"
```

---

### Task 5: Persist Interrupted Context and Reuse the Same Track for `/mp:plan`

**Why:** The product requirement is not just to block; it must preserve the original task intent and let `/mp:plan` continue the same lifecycle.

**What:** Save interrupted command context on high-complexity block, allocate or reuse a track, and ensure `/mp:plan` reuses that track instead of fragmenting work across multiple ids.

**How:** Centralize track allocation/reuse in `internal/cli/root.go` with small helpers in `internal/tracks` so blocked flows and `/mp:plan` share the same track identity.

**Key Design:** Recovery is explicit and user-driven; do not auto-execute the blocked command after planning finishes.

**Files:**
- Modify: `internal/cli/root.go`
- Modify: `internal/cli/root_test.go`
- Modify: `internal/tracks/metadata.go`
- Modify: `internal/tracks/metadata_test.go`
- Modify: `internal/tracks/coordinator.go`
- Optional Create: `internal/tracks/interrupted_context.go`

**Step 1: Write the failing tests**

```go
func TestHighComplexityBlock_SavesInterruptedContext(t *testing.T) {}
func TestPlanAfterBlock_ReusesSameTrackID(t *testing.T) {}
func TestPlanCompletion_DoesNotAutoReplayBlockedCommand(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/cli ./internal/tracks -run 'Test(HighComplexityBlock|PlanAfterBlock|PlanCompletion)_' -count=1`
Expected: FAIL because interrupted context is not persisted or reused yet.

**Step 3: Write minimal implementation**

- On planning-required block, resolve a track id if one does not already exist.
- Persist `InterruptedContext` with command, subcommand, prompt, and timestamp.
- Ensure `/mp:plan` resolves the same track and updates its canonical artifacts.
- Keep replay manual: only expose resume metadata; never auto-run the blocked command.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/cli ./internal/tracks -run 'Test(HighComplexityBlock|PlanAfterBlock|PlanCompletion)_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/cli/root.go internal/cli/root_test.go internal/tracks/metadata.go internal/tracks/metadata_test.go internal/tracks/coordinator.go internal/tracks/interrupted_context.go
git commit -m "track(complexity-scoring-gate-20260307): group g5 - persist interrupted context and reuse track"
```

---

### Task 6: Add End-to-End Admission Lifecycle Coverage and Update Operator Docs

**Why:** This feature spans runtime behavior, track lifecycle, and user guidance; without E2E tests and docs it will regress quickly.

**What:** Add lifecycle tests for the full block -> `/mp:plan` -> worktree-ready -> continue flow, and update docs to describe the new admission behavior.

**How:** Extend app/CLI lifecycle tests and update user-facing docs to document when `/mp:plan` is mandatory for high-complexity tasks.

**Key Design:** The docs must describe canonical planning artifacts, not the old `spec.md` / `plan.md` shape.

**Files:**
- Create: `internal/app/spec_admission_lifecycle_test.go`
- Modify: `internal/app/spec_track_lifecycle_test.go`
- Modify: `docs/CLI-REFERENCE.md`
- Modify: `custom/docs/target-project/getting-started.md`
- Reference: `docs/plans/2026-03-07-complexity-scoring-admission-gate-design.md`

**Step 1: Write the failing tests**

```go
func TestAdmissionLifecycle_BlockPlanThenContinueInWorktree(t *testing.T) {}
func TestAdmissionLifecycle_LegacySpecPlanDoNotSatisfyPlanning(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/app -run 'TestAdmissionLifecycle_' -count=1`
Expected: FAIL because the end-to-end planning lifecycle is not fully enforced yet.

**Step 3: Write minimal implementation and docs**

- Add the lifecycle coverage.
- Update runtime help/docs to state:
  - high complexity requires `/mp:plan`,
  - planning completeness is based on canonical artifacts,
  - execution still requires a linked worktree.

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/app -run 'TestAdmissionLifecycle_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/app/spec_admission_lifecycle_test.go internal/app/spec_track_lifecycle_test.go docs/CLI-REFERENCE.md custom/docs/target-project/getting-started.md
git commit -m "track(complexity-scoring-gate-20260307): group g6 - add admission lifecycle coverage and docs"
```

---

## Final Verification

After all tasks are complete, run:

```bash
go test ./internal/tracks ./internal/validation ./internal/app ./internal/cli -count=1
go test ./... -count=1
```

Expected:

- high-complexity spec-driven commands block unless canonical planning is complete,
- `/mp:plan` reuses the same track and preserves interrupted context,
- high-complexity execution still blocks outside a linked worktree,
- low-complexity commands continue without mandatory planning,
- no runtime path treats legacy `spec.md` / `plan.md` as sufficient planning evidence.
