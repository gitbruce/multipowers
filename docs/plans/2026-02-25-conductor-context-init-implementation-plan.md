# Conductor Context Auto-Init Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use executing-plans to implement this plan task-by-task.

**Goal:** Replace `.claude/session-*` spec artifacts with Conductor-style project context under `conductor/`, with auto-init and checkbox tracking for all spec-driven `/octo` commands.

**Architecture:** Keep upstream-impact minimal: implement most behavior in `custom/*` (templates/libs/docs/tests), and add only thin hooks in high-churn upstream files. Use a central command-level guard in `bin/octo` for spec-driven commands. If context is missing/incomplete, auto-run interactive `/octo:init`, render templates from `custom/templates/`, then execute with mandatory conductor-context loading and `conductor/tracks/` checkbox updates.

**Tech Stack:** Bash (`bin/octo`, command markdown contracts), custom overlay libs/scripts, markdown templates, shell tests.

---

## Task 0: Branch Strategy Guardrails (Main-Upstream Principle)

**Files:**
- Modify: `custom/scripts/octo-devx sync`
- Create: `tests/unit/test-main-upstream-discipline.sh`
- Modify: `custom/docs/sync/upstream-sync-playbook.md`

**Subtasks:**
- [x] T0.1 Document invariant: `main` is upstream mirror only (no custom commits)
- [x] T0.2 Enforce all customization work on `multipowers` only
- [x] T0.3 Add conflict-budget rule: minimize edits in `.claude/*`, `.claude-plugin/*`, `bin/octo`
- [x] T0.4 Add verification command to show customization footprint is mostly `custom/*`

**Verification:**
- [x] `git log main --not upstream/main --oneline` is empty
- [x] `git diff --name-only upstream/main...multipowers | awk -F/ '{print $1}' | sort | uniq -c`

---

## Task 1: Add Conductor Template Assets

**Files:**
- Create: `custom/templates/conductor/product.md`
- Create: `custom/templates/conductor/product-guidelines.md`
- Create: `custom/templates/conductor/tech-stack.md`
- Create: `custom/templates/conductor/workflow.md`
- Create: `custom/templates/conductor/tracks.md`
- Create: `custom/templates/conductor/code_styleguides/README.md`

**Subtasks:**
- [x] T1.1 Create template directory structure
- [x] T1.2 Add variable placeholders for init answers
- [x] T1.3 Add template metadata/version comments
- [x] T1.4 Add minimal default checkbox sections for tracks index

**Verification:**
- [x] `find custom/templates/conductor -type f | sort`
- [x] `rg -n "\{\{[A-Z0-9_]+\}\}" custom/templates/conductor`

---

## Task 1.5: Import Conductor Setup References (Explicit Source Borrowing)

**Files:**
- Create: `custom/references/conductor-upstream/README.md`
- Create: `custom/references/conductor-upstream/SOURCE-MAP.md`
- Modify: `custom/docs/customizations/conductor-context.md` (or create if missing)
- Modify: `custom/docs/README.md`

**Subtasks:**
- [x] T1.5.1 Pull setup-flow reference from `https://github.com/gemini-cli-extensions/conductor`
- [x] T1.5.2 Document exactly which source files/behaviors are borrowed
- [x] T1.5.3 Add license/attribution note and update policy
- [x] T1.5.4 Record what is copied verbatim vs adapted

**Verification:**
- [x] `custom/references/conductor-upstream/SOURCE-MAP.md` lists upstream commit/ref
- [x] docs clearly state source attribution and sync/update process

---

## Task 2: Implement Conductor Context Library

**Files:**
- Create: `custom/lib/conductor-context.sh`
- Modify: `custom/lib/overlay-loader.sh`

**Subtasks:**
- [x] T2.1 Implement `is_spec_driven_command`
- [x] T2.2 Implement context existence/completeness checks
- [x] T2.3 Implement `ensure_conductor_context`
- [x] T2.4 Implement `load_conductor_context_for_prompt`
- [x] T2.5 Source this library from overlay loader safely

**Verification:**
- [x] `bash -n custom/lib/conductor-context.sh`
- [x] `bash -n custom/lib/overlay-loader.sh`

---

## Task 3: Add `/octo:init` Interactive Setup

**Files:**
- Create: `custom/commands/init.md`
- Modify: `.claude-plugin/plugin.json`
- Modify: `bin/octo`
- Modify: `custom/scripts/octo-devx overlay`

**Subtasks:**
- [x] T3.1 Add command contract for `/octo:init` under `custom/commands/`
- [x] T3.2 Register `init` command in plugin command list
- [x] T3.3 Add orchestrator command handler for `init`
- [x] T3.4 Implement interactive Q&A flow (conductor-style)
- [x] T3.5 Render templates into target `conductor/`
- [x] T3.6 Implement overwrite decision per existing file
- [x] T3.7 Sync `custom/commands/init.md` into runtime command path via overlay script

**Verification:**
- [x] `bash -n bin/octo`
- [x] `/octo:init` creates all required conductor files

---

## Task 4: Route Plan/Intent Outputs from `.claude` to `conductor/tracks`

**Files:**
- Modify: `bin/octo`
- Modify: `custom/commands/plan.md` (if override is required)

**Subtasks:**
- [x] T4.1 Route plan output to `conductor/tracks/<track_id>/plan.md`
- [x] T4.2 Route intent contract to `conductor/tracks/<track_id>/intent.md`
- [x] T4.3 Define track file naming in `conductor/tracks/`
- [x] T4.4 Ensure track file includes checkbox status blocks
- [x] T4.5 Keep upstream command/skill file edits to minimum; prefer override through `custom/commands/*`

**Verification:**
- [x] `/octo:plan` writes only under `conductor/tracks/`
- [x] `rg -n "conductor/tracks/.*/(plan|intent)\.md" scripts` confirms conductor track write paths

---

## Task 5: Add Command-Level Guard to All Spec-Driven Commands

**Files:**
- Modify: `bin/octo`
- Modify: `custom/lib/conductor-context.sh`
- Modify: `docs/COMMAND-REFERENCE.md`

**Subtasks:**
- [x] T5.1 Add preflight call to `ensure_conductor_context`
- [x] T5.2 If missing context, auto-call `/octo:init`
- [x] T5.3 Re-check context after init before proceeding
- [x] T5.4 Exclude non-spec commands from this guard
- [x] T5.5 Implement spec-driven command allowlist in one central place
- [x] T5.6 Keep upstream hook footprint minimal (single dispatch point in `bin/octo`)

**Verification:**
- [x] Missing conductor context triggers init automatically
- [x] Existing complete context skips init

---

## Task 6: Mandatory Context Read Before Spec-Driven Execution

**Files:**
- Modify: `bin/octo`
- Modify: `custom/lib/conductor-context.sh`

**Subtasks:**
- [x] T6.1 Load `conductor/*.md` into execution context
- [x] T6.2 Ensure context loading occurs before agent/task dispatch
- [x] T6.3 Add explicit log line indicating context load success

**Verification:**
- [x] Spec-driven command logs conductor context load before execution
- [x] Command behavior changes when context content changes

---

## Task 7: Checkbox Tracking in `conductor/tracks/`

**Files:**
- Modify: `bin/octo`
- Create/Modify: `custom/templates/conductor/tracks.md`

**Subtasks:**
- [x] T7.1 Create per-run track files with standardized sections
- [x] T7.2 Write tasks/subtasks as checkboxes
- [x] T7.3 Update checkbox state through workflow lifecycle
- [x] T7.4 Record final status and unresolved items
- [x] T7.5 Enforce canonical path `conductor/tracks/` (no `conductor/track/` variants)

**Verification:**
- [x] `conductor/tracks/<date>-<slug>.md` is created on spec-driven runs
- [x] Checkboxes transition from `[ ]` to `[x]` during execution
- [x] `rg -n "conductor/track/" . custom docs scripts` returns no active references

---

## Task 8: Tests for Guard, Init, Context, and Tracking

**Files:**
- Create: `tests/unit/test-conductor-context-guard.sh`
- Create: `tests/unit/test-octo-init-render.sh`
- Create: `tests/integration/test-spec-commands-auto-init.sh`
- Create: `tests/integration/test-spec-commands-context-read.sh`
- Create: `tests/integration/test-tracks-checkbox-updates.sh`

**Subtasks:**
- [x] T8.1 Unit tests for completeness checker and allowlist
- [x] T8.2 Unit tests for template render output paths
- [x] T8.3 Integration: missing context auto-inits
- [x] T8.4 Integration: spec commands read context before execution
- [x] T8.5 Integration: no `.claude/session-*` writes remain

**Verification:**
- [x] All new tests pass
- [x] Existing core smoke tests still pass

---

## Task 9: Documentation Updates

**Files:**
- Modify: `custom/docs/getting-started.md`
- Modify: `custom/docs/customizations/persona-command.md` (if command behavior references old flow)
- Create: `custom/docs/customizations/conductor-context.md`
- Modify: `docs/COMMAND-REFERENCE.md`
- Modify: `custom/docs/sync/upstream-sync-playbook.md`

**Subtasks:**
- [x] T9.1 Document `/octo:init` usage and prompts
- [x] T9.2 Document spec-driven command guard behavior
- [x] T9.3 Document conductor file structure and track workflow
- [x] T9.4 Add migration note from `.claude/session-*` to `conductor/tracks/*`
- [x] T9.5 Add branch strategy note: `main` stays upstream mirror; all custom behavior lives on `multipowers`
- [x] T9.6 Explicitly document canonical folder names: `conductor/tracks.md` + `conductor/tracks/`
- [x] T9.7 Document borrowed Conductor setup source map and attribution

**Verification:**
- [x] `rg -n "\.claude/session-(plan|intent)\.md" custom/docs docs` only appears in migration/deprecation notes

---

## Task 10: Final Validation and Release Commit

**Files:**
- Modify: `CHANGELOG.md` (if applicable)

**Subtasks:**
- [x] T10.1 Run syntax checks on changed shell scripts
- [x] T10.2 Run targeted unit/integration tests
- [x] T10.3 Manual sanity pass for `/octo:init` and `/octo:plan`
- [x] T10.4 Commit in logical slices (core guard/init, then docs/tests)
- [x] T10.5 Push to `origin/multipowers`
- [x] T10.6 Verify customization impact is minimized in upstream-heavy paths

**Verification Commands:**
- [x] `bash -n bin/octo custom/lib/*.sh custom/scripts/*.sh`
- [x] `bash tests/unit/test-conductor-context-guard.sh`
- [x] `bash tests/integration/test-spec-commands-auto-init.sh`
- [x] `bash tests/integration/test-tracks-checkbox-updates.sh`
- [x] `git diff --name-only upstream/main...HEAD | rg -v '^(custom/|docs/plans/|tests/(unit|integration)/test-conductor)'`

---

## Rollback Plan
- Revert guard wiring in command files and `bin/octo`.
- Keep `/octo:init` command isolated behind feature flag if needed.
- Restore prior plan behavior only if conductor flow cannot be stabilized.

## Done Criteria
- [x] Spec-driven `/octo` commands write plan/intent under `conductor/tracks/<track_id>/`.
- [x] Missing conductor context triggers interactive `/octo:init` automatically.
- [x] Spec-driven commands read `conductor/*.md` before execution.
- [x] Track files in `conductor/tracks/` contain checkbox status and are updated during run.
- [x] Tests cover guard/init/context/tracking contracts.
