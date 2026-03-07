# Multi Branch Superpowers Mainline Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rebuild the `multi` branch around the single public mainline `init → brainstorm → design → plan → execute`, keep `debug` and `debate` as the only special direct-entry commands, and make `superpowers` upstream command/skill Markdown the source of truth for same-function workflow text.

**Architecture:** Keep `go` as the upstream runtime baseline and implement the product reset only on `multi`. Sync selected upstream `superpowers` command/skill Markdown into repository references, generate final plugin Markdown from thin wrappers plus upstream bodies, replace persona-market entrypoints with fixed internal roles, and rewire runtime policy/hooks so `/mp:init` is a hard prerequisite, phase model policy is user-configurable, and `debate` always fans out to all configured models.

**Tech Stack:** Go, `.claude-plugin` plugin assets, `internal/devx`, `internal/cli`, `internal/hooks`, `internal/policy`, `internal/providers`, `internal/workflows`, `config/workflows.yaml`, `config/providers.yaml`, fixed-role config, upstream Markdown references, `go test`

---

## Execution Mode Decision

- Complexity Score: 8
- Worktree Required: YES
- Why: This change resets the public product surface, runtime routing, wrapper generation, model policy, persona removal, and user docs across plugin assets and Go runtime layers.
- Enforcement: Author the plan on branch `multi`, but execute implementation from a dedicated linked worktree rooted at `multi`. After each task, run the listed tests and commit with `multi(mainline-20260308): group <gid> - <title>`.

## Before Starting

- This plan is authored on branch `multi`, created from `go`.
- Before implementation, create a dedicated linked worktree from `multi`.
- Recommended branch: `multi-mainline-20260308`
- Recommended worktree path: `.worktrees/multi-mainline-20260308`
- Confirm the worktree starts from the current `multi` HEAD before editing files.

Example:

```bash
git worktree add .worktrees/multi-mainline-20260308 -b multi-mainline-20260308 multi
cd .worktrees/multi-mainline-20260308
go test ./internal/devx ./internal/policy ./internal/providers ./internal/workflows ./internal/hooks ./internal/cli -count=1
```

---

### Task 1: Sync Selected Upstream Superpowers Markdown into Repository References

**Why:** Same-function command/skill text must come from upstream `superpowers` so the `multi` branch stops carrying a competing local workflow narrative.

**What:** Add a sync mechanism that pulls only the selected upstream `commands` and `skills` into a stable repository location used by the wrapper generator.

**How:** Introduce a small `mp-devx` sync action plus a sync manifest listing the exact upstream files required by the new public surface and internal sub-skills.

**Key Design:** Sync only the selected upstream assets needed by the new mainline; do not mirror the entire `superpowers` repository.

**Files:**
- Create: `custom/config/superpowers-sync.yaml`
- Create: `custom/references/superpowers-upstream/README.md`
- Create: `internal/devx/superpowers_sync.go`
- Create: `internal/devx/superpowers_sync_test.go`
- Modify: `internal/devx/runner.go`
- Modify: `cmd/mp-devx/main.go`

**Step 1: Write the failing tests**

```go
func TestSyncSuperpowers_WritesSelectedCommandsAndSkills(t *testing.T) {}
func TestSyncSuperpowers_RejectsUnexpectedSelections(t *testing.T) {}
func TestSyncSuperpowers_PreservesStableOutputPaths(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run 'TestSyncSuperpowers_' -count=1`
Expected: FAIL because the sync action and manifest do not exist yet.

**Step 3: Write minimal implementation**

- Add `superpowers-sync.yaml` listing:
  - `commands/brainstorm.md`
  - `commands/write-plan.md`
  - `commands/execute-plan.md`
  - `skills/brainstorming/SKILL.md`
  - `skills/writing-plans/SKILL.md`
  - `skills/executing-plans/SKILL.md`
  - `skills/systematic-debugging/SKILL.md`
  - `skills/verification-before-completion/SKILL.md`
  - `skills/finishing-a-development-branch/SKILL.md`
  - `skills/requesting-code-review/SKILL.md`
  - `skills/receiving-code-review/SKILL.md`
  - `skills/using-git-worktrees/SKILL.md`
  - `skills/subagent-driven-development/SKILL.md`
- Add a `Runner.SyncSuperpowersAssets(...)` implementation that writes the selected files under `custom/references/superpowers-upstream/...`.
- Add a new `mp-devx --action sync-superpowers` entrypoint.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run 'TestSyncSuperpowers_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add custom/config/superpowers-sync.yaml custom/references/superpowers-upstream/README.md internal/devx/superpowers_sync.go internal/devx/superpowers_sync_test.go internal/devx/runner.go cmd/mp-devx/main.go
git commit -m "multi(mainline-20260308): group g1 - add selected superpowers sync"
```

---

### Task 2: Add a Single Mainline Surface Manifest and Fixed Role Definitions

**Why:** The branch needs one authoritative definition of the public command surface and one authoritative definition of the fixed internal roles that replace personas.

**What:** Add a manifest for public commands/skills and a fixed role manifest for internal runtime use.

**How:** Keep public surface and internal role definitions explicit in config instead of scattering them across `plugin.json`, `agents/personas/*`, and old Double Diamond docs.

**Key Design:** Roles are internal implementation detail only; users never select them directly.

**Files:**
- Create: `custom/config/mainline-surface.yaml`
- Create: `config/roles.yaml`
- Create: `agents/roles/initializer.md`
- Create: `agents/roles/facilitator.md`
- Create: `agents/roles/planner.md`
- Create: `agents/roles/executor.md`
- Create: `agents/roles/reviewer.md`
- Create: `agents/roles/debugger.md`
- Create: `agents/roles/debater.md`
- Create: `internal/roles/manifest.go`
- Create: `internal/roles/manifest_test.go`

**Step 1: Write the failing tests**

```go
func TestLoadRoles_HasExactFixedRoleSet(t *testing.T) {}
func TestLoadSurface_PublicCommandsAreMainlineOnly(t *testing.T) {}
func TestLoadSurface_MapsDesignToBrainstormingUpstream(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/roles -run 'TestLoad(Roles|Surface)_' -count=1`
Expected: FAIL because the fixed-role and mainline-surface manifests do not exist yet.

**Step 3: Write minimal implementation**

- Define public commands exactly as:
  - `init`, `model-config`, `brainstorm`, `design`, `plan`, `execute`, `debug`, `debate`, `status`, `doctor`, `resume`, `setup`
- Define fixed internal roles exactly as:
  - `initializer`, `facilitator`, `planner`, `executor`, `reviewer`, `debugger`, `debater`
- Map:
  - `brainstorm` → upstream `brainstorming`
  - `design` → upstream `brainstorming`
  - `plan` → upstream `writing-plans`
  - `execute` → upstream `executing-plans`
  - `debug` → upstream `systematic-debugging`
  - `debate` → local runtime debate handler

**Step 4: Run test to verify it passes**

Run: `go test ./internal/roles -run 'TestLoad(Roles|Surface)_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add custom/config/mainline-surface.yaml config/roles.yaml agents/roles internal/roles/manifest.go internal/roles/manifest_test.go
git commit -m "multi(mainline-20260308): group g2 - add mainline surface and fixed roles"
```

---

### Task 3: Add Mainline Workflow Entry Points in Go Runtime

**Why:** The public runtime must speak the new workflow vocabulary (`brainstorm`, `design`, `execute`) instead of only the old Double Diamond vocabulary.

**What:** Add explicit workflow entry points and adapter routes for `brainstorm`, `design`, and `execute`.

**How:** Mirror the existing lightweight workflow wrappers, but target the new mainline surface and keep `plan`, `debug`, and `debate` as first-class runtime commands.

**Key Design:** Preserve thin runtime helpers; the new semantic split belongs in command routing and wrapper generation, not in large prompt bodies.

**Files:**
- Create: `internal/workflows/brainstorm.go`
- Create: `internal/workflows/design.go`
- Create: `internal/workflows/execute.go`
- Create: `internal/workflows/mainline_test.go`
- Modify: `internal/workflows/adapter_helper.go`
- Modify: `internal/providers/router_intent.go`
- Modify: `internal/cli/root.go`

**Step 1: Write the failing tests**

```go
func TestMainlineWorkflow_BrainstormEntryPoint(t *testing.T) {}
func TestMainlineWorkflow_DesignEntryPoint(t *testing.T) {}
func TestMainlineWorkflow_ExecuteEntryPoint(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/workflows ./internal/cli -run 'TestMainlineWorkflow_' -count=1`
Expected: FAIL because the new workflow entry points and CLI routing do not exist yet.

**Step 3: Write minimal implementation**

- Add `Brainstorm(prompt string)`, `Design(prompt string)`, and `Execute(prompt string)` wrappers under `internal/workflows`.
- Extend the shared workflow adapter helper and CLI command switch to recognize `brainstorm`, `design`, and `execute`.
- Keep `plan` and `debug` independent runtime commands.
- Remove new-command dependence on legacy `discover/define/develop/deliver` naming.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/workflows ./internal/cli -run 'TestMainlineWorkflow_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/workflows/brainstorm.go internal/workflows/design.go internal/workflows/execute.go internal/workflows/mainline_test.go internal/workflows/adapter_helper.go internal/providers/router_intent.go internal/cli/root.go
git commit -m "multi(mainline-20260308): group g3 - add mainline runtime entrypoints"
```

---

### Task 4: Build a Thin Wrapper Generator for Commands and Skills

**Why:** The branch should not hand-maintain a second full set of command/skill Markdown when upstream already provides the workflow bodies.

**What:** Add a generator that combines upstream Markdown, local thin wrapper metadata, and local runtime-specific injected sections into final plugin assets.

**How:** Use templates plus the surface manifest to build `.claude-plugin/.claude/commands/*` and `.claude-plugin/.claude/skills/*` deterministically during `mp-devx build-runtime`.

**Key Design:** Upstream body stays intact; local wrapper content is additive and limited to init/model policy/runtime glue.

**Files:**
- Create: `custom/templates/mainline-wrapper/command.md.tpl`
- Create: `custom/templates/mainline-wrapper/skill.md.tpl`
- Create: `internal/devx/wrapper_builder.go`
- Create: `internal/devx/wrapper_builder_test.go`
- Modify: `internal/devx/runner.go`
- Modify: `cmd/mp-devx/main.go`

**Step 1: Write the failing tests**

```go
func TestBuildMainlineAssets_WritesPublicCommands(t *testing.T) {}
func TestBuildMainlineAssets_EmbedsSelectedUpstreamBodies(t *testing.T) {}
func TestBuildMainlineAssets_InjectsRuntimeSections(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/devx -run 'TestBuildMainlineAssets_' -count=1`
Expected: FAIL because the wrapper builder and templates do not exist yet.

**Step 3: Write minimal implementation**

- Add a wrapper builder that reads:
  - `custom/config/mainline-surface.yaml`
  - `custom/references/superpowers-upstream/...`
  - local command/skill injection sections
- Generate final assets for:
  - commands: `init`, `model-config`, `brainstorm`, `design`, `plan`, `execute`, `debug`, `debate`, `status`, `doctor`, `resume`, `setup`
  - skills: `mainline-brainstorm`, `mainline-design`, `mainline-plan`, `mainline-execute`, `mainline-debug`, `mainline-debate`
- Wire the generator into `mp-devx build-runtime`.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/devx -run 'TestBuildMainlineAssets_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add custom/templates/mainline-wrapper/command.md.tpl custom/templates/mainline-wrapper/skill.md.tpl internal/devx/wrapper_builder.go internal/devx/wrapper_builder_test.go internal/devx/runner.go cmd/mp-devx/main.go
git commit -m "multi(mainline-20260308): group g4 - add mainline wrapper generator"
```

---

### Task 5: Generate the New Public Plugin Surface and Trim `plugin.json`

**Why:** The public plugin assets must expose only the approved mainline commands and thin wrapper skills.

**What:** Regenerate the plugin assets and reduce `plugin.json` to the new surface.

**How:** Use the wrapper generator to overwrite the final plugin assets, then tighten `plugin.json` to the retained commands/skills only.

**Key Design:** `plugin.json` is the canonical public surface contract; if a command/skill is not in the retained list, it should not remain published by accident.

**Files:**
- Modify: `.claude-plugin/plugin.json`
- Modify: `.claude-plugin/.claude/commands/init.md`
- Modify: `.claude-plugin/.claude/commands/model-config.md`
- Modify: `.claude-plugin/.claude/commands/brainstorm.md`
- Create: `.claude-plugin/.claude/commands/design.md`
- Modify: `.claude-plugin/.claude/commands/plan.md`
- Create: `.claude-plugin/.claude/commands/execute.md`
- Modify: `.claude-plugin/.claude/commands/debug.md`
- Modify: `.claude-plugin/.claude/commands/debate.md`
- Modify: `.claude-plugin/.claude/commands/status.md`
- Modify: `.claude-plugin/.claude/commands/doctor.md`
- Modify: `.claude-plugin/.claude/commands/resume.md`
- Modify: `.claude-plugin/.claude/commands/setup.md`
- Create: `.claude-plugin/.claude/skills/mainline-brainstorm.md`
- Create: `.claude-plugin/.claude/skills/mainline-design.md`
- Create: `.claude-plugin/.claude/skills/mainline-plan.md`
- Create: `.claude-plugin/.claude/skills/mainline-execute.md`
- Create: `.claude-plugin/.claude/skills/mainline-debug.md`
- Create: `.claude-plugin/.claude/skills/mainline-debate.md`

**Step 1: Write the failing tests**

```go
func TestPluginSurface_ContainsOnlyMainlineCommands(t *testing.T) {}
func TestPluginSurface_ContainsOnlyWrapperSkills(t *testing.T) {}
func TestPluginSurface_ExposesDesignAndExecute(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./cmd/mp-devx -run 'TestPluginSurface_' -count=1`
Expected: FAIL because `plugin.json` and generated assets still expose the old public surface.

**Step 3: Write minimal implementation**

- Regenerate the retained command and skill files.
- Trim `plugin.json` to exactly:
  - commands: `init`, `model-config`, `brainstorm`, `design`, `plan`, `execute`, `debug`, `debate`, `status`, `doctor`, `resume`, `setup`
  - skills: `mainline-brainstorm`, `mainline-design`, `mainline-plan`, `mainline-execute`, `mainline-debug`, `mainline-debate`
- Ensure command wrappers reference the generated wrapper skills or direct runtime actions only.

**Step 4: Run test to verify it passes**

Run: `go test ./cmd/mp-devx -run 'TestPluginSurface_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add .claude-plugin/plugin.json .claude-plugin/.claude/commands .claude-plugin/.claude/skills
git commit -m "multi(mainline-20260308): group g5 - publish mainline plugin surface"
```

---

### Task 6: Make `/mp:init` a Hard Prerequisite for Every Mainline Entry

**Why:** The new product surface depends on a consistent initialization contract; users should not guess whether `init` is required.

**What:** Enforce `/mp:init` artifacts as a hard prerequisite for `brainstorm`, `design`, `plan`, `execute`, `debug`, and `debate`, with redirect-and-resume guidance.

**How:** Reuse the existing context-completeness checks, but update the hook/runtime command vocabulary and response metadata to the new mainline commands.

**Key Design:** Missing init blocks should be uniform and should preserve enough metadata to send the user to `/mp:init` and then back to the original command.

**Files:**
- Modify: `internal/hooks/handler.go`
- Create: `internal/hooks/mainline_init_test.go`
- Modify: `internal/cli/root.go`
- Modify: `internal/cli/root_test.go`
- Reference: `internal/context/*`

**Step 1: Write the failing tests**

```go
func TestMainlineInitGate_BlocksBrainstormWithoutContext(t *testing.T) {}
func TestMainlineInitGate_BlocksExecuteWithoutContext(t *testing.T) {}
func TestMainlineInitGate_PreservesResumeMetadata(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/hooks ./internal/cli -run 'TestMainlineInitGate_' -count=1`
Expected: FAIL because hook routing still keys off old command names and does not cover the new mainline surface.

**Step 3: Write minimal implementation**

- Replace old `discover/define/develop/deliver/embrace/review/research` prompt detection with the retained surface.
- Ensure missing init guidance points to `/mp:init`.
- Keep `resume_command` / `resume_prompt` metadata stable so the wrapper can resume the original flow.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/hooks ./internal/cli -run 'TestMainlineInitGate_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/hooks/handler.go internal/hooks/mainline_init_test.go internal/cli/root.go internal/cli/root_test.go
git commit -m "multi(mainline-20260308): group g6 - gate mainline commands on init"
```

---

### Task 7: Add Phase Model Policy and Force Debate to Use All Configured Models

**Why:** The new branch keeps user-configurable model selection, but the policy must be phase-based instead of persona-based, and `debate` must always use every configured model.

**What:** Rewrite workflow model policy around the new phase names and implement all-model fanout for `debate`.

**How:** Keep using policy/provider configuration, but change the vocabulary from old workflow names and persona routing to phase model policy and debate-specific fanout.

**Key Design:** `brainstorm`, `design`, and `plan` may run multiple configured models in parallel; `debate` always uses the full configured model set.

**Files:**
- Modify: `config/workflows.yaml`
- Modify: `config/providers.yaml`
- Modify: `internal/policy/types.go`
- Modify: `internal/policy/validate.go`
- Modify: `internal/policy/validate_test.go`
- Modify: `internal/providers/provider.go`
- Modify: `internal/providers/quorum_test.go`
- Modify: `internal/workflows/debate.go`
- Create: `internal/workflows/model_policy_test.go`

**Step 1: Write the failing tests**

```go
func TestValidateWorkflows_MainlinePhasePolicies(t *testing.T) {}
func TestDebate_UsesAllConfiguredProviders(t *testing.T) {}
func TestBrainstorm_UsesConfiguredParallelProviders(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/policy ./internal/providers ./internal/workflows -run 'Test(ValidateWorkflows_MainlinePhasePolicies|Debate_UsesAllConfiguredProviders|Brainstorm_UsesConfiguredParallelProviders)' -count=1`
Expected: FAIL because the policy vocabulary and debate fanout logic still follow the old surface.

**Step 3: Write minimal implementation**

- Rewrite `config/workflows.yaml` to phase names:
  - `brainstorm`, `design`, `plan`, `execute`, `debug`, `debate`
- Keep user-configurable model selection by phase.
- Make `debate` fan out to all configured/active providers instead of a fixed subset.
- Keep validation deterministic and fail fast on unknown phase names or empty debate provider sets.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/policy ./internal/providers ./internal/workflows -run 'Test(ValidateWorkflows_MainlinePhasePolicies|Debate_UsesAllConfiguredProviders|Brainstorm_UsesConfiguredParallelProviders)' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add config/workflows.yaml config/providers.yaml internal/policy/types.go internal/policy/validate.go internal/policy/validate_test.go internal/providers/provider.go internal/providers/quorum_test.go internal/workflows/debate.go internal/workflows/model_policy_test.go
git commit -m "multi(mainline-20260308): group g7 - add mainline phase model policy"
```

---

### Task 8: Replace Persona Routing with Fixed Roles and Remove the Public Persona Path

**Why:** The new branch keeps internal roles but must stop exposing or depending on persona-market workflows.

**What:** Remove public persona routing, add a minimal internal role dispatcher, and migrate any retained internal calls away from `persona` nomenclature.

**How:** Keep the role set small and fixed; do not preserve user-facing persona enumeration or selection.

**Key Design:** Role choice is an internal runtime detail, not a public API.

**Files:**
- Create: `internal/roles/dispatch.go`
- Create: `internal/roles/dispatch_test.go`
- Delete: `internal/workflows/persona.go`
- Delete: `.claude-plugin/.claude/commands/persona.md`
- Modify: `internal/cli/root.go`
- Modify: `cmd/mp-devx/main_test.go`

**Step 1: Write the failing tests**

```go
func TestRoleDispatch_MapsMainlineRolesOnly(t *testing.T) {}
func TestPublicSurface_DoesNotExposePersonaCommand(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/roles ./cmd/mp-devx -run 'Test(RoleDispatch_|PublicSurface_DoesNotExposePersonaCommand)' -count=1`
Expected: FAIL because persona routing and persona command exposure still exist.

**Step 3: Write minimal implementation**

- Introduce a fixed internal role dispatcher.
- Remove public persona command exposure from CLI/plugin surface.
- Delete legacy persona workflow code once the fixed role dispatcher covers retained runtime needs.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/roles ./cmd/mp-devx -run 'Test(RoleDispatch_|PublicSurface_DoesNotExposePersonaCommand)' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/roles/dispatch.go internal/roles/dispatch_test.go internal/cli/root.go cmd/mp-devx/main_test.go
git rm internal/workflows/persona.go .claude-plugin/.claude/commands/persona.md
git commit -m "multi(mainline-20260308): group g8 - replace personas with fixed roles"
```

---

### Task 9: Remove Legacy Public Commands and Skills from the Plugin Surface

**Why:** The public branch cannot remain a hybrid of the old Double Diamond product and the new mainline product.

**What:** Delete obsolete plugin command/skill assets and add a regression test that bans them from returning.

**How:** Remove the old public files, keep only the retained surface, and lock that surface with validation tests.

**Key Design:** Deletion is a feature here; do not silently keep dead assets around “just in case.”

**Files:**
- Delete: `.claude-plugin/.claude/commands/discover.md`
- Delete: `.claude-plugin/.claude/commands/research.md`
- Delete: `.claude-plugin/.claude/commands/define.md`
- Delete: `.claude-plugin/.claude/commands/develop.md`
- Delete: `.claude-plugin/.claude/commands/deliver.md`
- Delete: `.claude-plugin/.claude/commands/embrace.md`
- Delete: `.claude-plugin/.claude/commands/review.md`
- Delete: `.claude-plugin/.claude/commands/validate.md`
- Delete: `.claude-plugin/.claude/commands/tdd.md`
- Delete: `.claude-plugin/.claude/commands/ship.md`
- Delete: `.claude-plugin/.claude/commands/multi.md`
- Delete: `.claude-plugin/.claude/commands/probe.md`
- Delete: `.claude-plugin/.claude/commands/grasp.md`
- Delete: `.claude-plugin/.claude/commands/tangle.md`
- Delete: `.claude-plugin/.claude/commands/ink.md`
- Delete: `.claude-plugin/.claude/commands/dev.md`
- Delete: `.claude-plugin/.claude/commands/quick.md`
- Delete: `.claude-plugin/.claude/commands/security.md`
- Delete: `.claude-plugin/.claude/commands/issues.md`
- Delete: `.claude-plugin/.claude/commands/rollback.md`
- Delete: `.claude-plugin/.claude/commands/pipeline.md`
- Delete: `.claude-plugin/.claude/commands/prd.md`
- Delete: `.claude-plugin/.claude/commands/prd-score.md`
- Delete: `.claude-plugin/.claude/commands/meta-prompt.md`
- Delete: `.claude-plugin/.claude/commands/deck.md`
- Delete: `.claude-plugin/.claude/commands/docs.md`
- Delete: `.claude-plugin/.claude/commands/extract.md`
- Delete: `.claude-plugin/.claude/commands/km.md`
- Delete: `.claude-plugin/.claude/commands/loop.md`
- Delete: `.claude-plugin/.claude/commands/sys-setup.md`
- Delete: `.claude-plugin/.claude/skills/flow-discover.md`
- Delete: `.claude-plugin/.claude/skills/flow-define.md`
- Delete: `.claude-plugin/.claude/skills/flow-develop.md`
- Delete: `.claude-plugin/.claude/skills/flow-deliver.md`
- Delete: `.claude-plugin/.claude/skills/skill-writing-plans.md`
- Delete: `.claude-plugin/.claude/skills/skill-tdd.md`
- Delete: `.claude-plugin/.claude/skills/skill-debug.md`
- Delete: `.claude-plugin/.claude/skills/skill-verify.md`
- Delete: `.claude-plugin/.claude/skills/skill-finish-branch.md`
- Delete: `.claude-plugin/.claude/skills/skill-code-review.md`
- Delete: `.claude-plugin/.claude/skills/skill-parallel-agents.md`
- Delete: `.claude-plugin/.claude/skills/skill-deep-research.md`
- Delete: `.claude-plugin/.claude/skills/skill-thought-partner.md`
- Delete: `.claude-plugin/.claude/skills/skill-validate.md`
- Delete: `.claude-plugin/.claude/skills/skill-ship.md`
- Delete: `.claude-plugin/.claude/skills/skill-rollback.md`
- Delete: `.claude-plugin/.claude/skills/skill-quick-review.md`
- Delete: `.claude-plugin/.claude/skills/skill-quick.md`
- Delete: `.claude-plugin/.claude/skills/skill-context-detection.md`
- Delete: `.claude-plugin/.claude/skills/skill-visual-feedback.md`
- Delete: `.claude-plugin/.claude/skills/skill-intent-contract.md`
- Delete: `.claude-plugin/.claude/skills/skill-task-management.md`
- Delete: `.claude-plugin/.claude/skills/skill-task-management-v2.md`
- Delete: `.claude-plugin/.claude/skills/skill-iterative-loop.md`
- Delete: `.claude-plugin/.claude/skills/skill-status.md`
- Delete: `.claude-plugin/.claude/skills/skill-issues.md`
- Delete: `.claude-plugin/.claude/skills/skill-resume.md`
- Delete: `.claude-plugin/.claude/skills/skill-resume-enhanced.md`
- Delete: `.claude-plugin/.claude/skills/skill-architecture.md`
- Delete: `.claude-plugin/.claude/skills/skill-security-audit.md`
- Delete: `.claude-plugin/.claude/skills/skill-security-framing.md`
- Delete: `.claude-plugin/.claude/skills/skill-audit.md`
- Delete: `.claude-plugin/.claude/skills/skill-doc-delivery.md`
- Delete: `.claude-plugin/.claude/skills/skill-prd.md`
- Delete: `.claude-plugin/.claude/skills/skill-content-pipeline.md`
- Delete: `.claude-plugin/.claude/skills/skill-knowledge-work.md`
- Delete: `.claude-plugin/.claude/skills/skill-meta-prompt.md`
- Delete: `.claude-plugin/.claude/skills/skill-deck.md`
- Delete: `.claude-plugin/.claude/skills/extract-skill.md`
- Create: `internal/validation/public_surface_test.go`

**Step 1: Write the failing tests**

```go
func TestPublicSurface_AllowsOnlyRetainedCommands(t *testing.T) {}
func TestPublicSurface_AllowsOnlyRetainedSkills(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/validation -run 'TestPublicSurface_' -count=1`
Expected: FAIL because the old files and surface entries still exist.

**Step 3: Write minimal implementation**

- Delete the obsolete public command and skill files listed above.
- Add a regression test that scans `.claude-plugin/plugin.json` and bans the removed public surface.
- Keep only retained commands and generated wrapper skills on the public plugin surface.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/validation -run 'TestPublicSurface_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/validation/public_surface_test.go .claude-plugin/plugin.json
for f in .claude-plugin/.claude/commands/discover.md .claude-plugin/.claude/commands/research.md .claude-plugin/.claude/commands/define.md .claude-plugin/.claude/commands/develop.md .claude-plugin/.claude/commands/deliver.md .claude-plugin/.claude/commands/embrace.md .claude-plugin/.claude/commands/review.md .claude-plugin/.claude/commands/validate.md .claude-plugin/.claude/commands/tdd.md .claude-plugin/.claude/commands/ship.md .claude-plugin/.claude/commands/multi.md .claude-plugin/.claude/commands/probe.md .claude-plugin/.claude/commands/grasp.md .claude-plugin/.claude/commands/tangle.md .claude-plugin/.claude/commands/ink.md .claude-plugin/.claude/commands/dev.md .claude-plugin/.claude/commands/quick.md .claude-plugin/.claude/commands/security.md .claude-plugin/.claude/commands/issues.md .claude-plugin/.claude/commands/rollback.md .claude-plugin/.claude/commands/pipeline.md .claude-plugin/.claude/commands/prd.md .claude-plugin/.claude/commands/prd-score.md .claude-plugin/.claude/commands/meta-prompt.md .claude-plugin/.claude/commands/deck.md .claude-plugin/.claude/commands/docs.md .claude-plugin/.claude/commands/extract.md .claude-plugin/.claude/commands/km.md .claude-plugin/.claude/commands/loop.md .claude-plugin/.claude/commands/sys-setup.md .claude-plugin/.claude/skills/flow-discover.md .claude-plugin/.claude/skills/flow-define.md .claude-plugin/.claude/skills/flow-develop.md .claude-plugin/.claude/skills/flow-deliver.md .claude-plugin/.claude/skills/skill-writing-plans.md .claude-plugin/.claude/skills/skill-tdd.md .claude-plugin/.claude/skills/skill-debug.md .claude-plugin/.claude/skills/skill-verify.md .claude-plugin/.claude/skills/skill-finish-branch.md .claude-plugin/.claude/skills/skill-code-review.md .claude-plugin/.claude/skills/skill-parallel-agents.md .claude-plugin/.claude/skills/skill-deep-research.md .claude-plugin/.claude/skills/skill-thought-partner.md .claude-plugin/.claude/skills/skill-validate.md .claude-plugin/.claude/skills/skill-ship.md .claude-plugin/.claude/skills/skill-rollback.md .claude-plugin/.claude/skills/skill-quick-review.md .claude-plugin/.claude/skills/skill-quick.md .claude-plugin/.claude/skills/skill-context-detection.md .claude-plugin/.claude/skills/skill-visual-feedback.md .claude-plugin/.claude/skills/skill-intent-contract.md .claude-plugin/.claude/skills/skill-task-management.md .claude-plugin/.claude/skills/skill-task-management-v2.md .claude-plugin/.claude/skills/skill-iterative-loop.md .claude-plugin/.claude/skills/skill-status.md .claude-plugin/.claude/skills/skill-issues.md .claude-plugin/.claude/skills/skill-resume.md .claude-plugin/.claude/skills/skill-resume-enhanced.md .claude-plugin/.claude/skills/skill-architecture.md .claude-plugin/.claude/skills/skill-security-audit.md .claude-plugin/.claude/skills/skill-security-framing.md .claude-plugin/.claude/skills/skill-audit.md .claude-plugin/.claude/skills/skill-doc-delivery.md .claude-plugin/.claude/skills/skill-prd.md .claude-plugin/.claude/skills/skill-content-pipeline.md .claude-plugin/.claude/skills/skill-knowledge-work.md .claude-plugin/.claude/skills/skill-meta-prompt.md .claude-plugin/.claude/skills/skill-deck.md .claude-plugin/.claude/skills/extract-skill.md; do git rm "$f"; done
git commit -m "multi(mainline-20260308): group g9 - remove legacy public surface"
```

---

### Task 10: Rewrite Mainline Documentation and Archive Obsolete User Docs

**Why:** The codebase must stop advertising the old product vocabulary once the surface is trimmed.

**What:** Rewrite the primary docs to match the new surface and archive or remove old user-facing docs bound to deleted commands, skills, and personas.

**How:** Keep architecture/runtime docs that still describe valid Go subsystems, but rewrite the README and user command docs around the new mainline.

**Key Design:** Documentation should describe exactly one public workflow story.

**Files:**
- Modify: `README.md`
- Modify: `docs/WORKFLOW-SKILLS.md`
- Modify: `docs/COMMAND-REFERENCE.md`
- Modify: `docs/PLUGIN-ARCHITECTURE.md`
- Modify: `docs/CLI-REFERENCE.md`
- Create: `docs/architecture/multi-mainline.md`
- Create: `docs/archive/README.md`
- Move or Delete: user docs bound to removed commands/persona-market usage

**Step 1: Write the failing tests**

```go
func TestDocs_PublicMainlineVocabularyOnly(t *testing.T) {}
func TestDocs_ReadmeQuickStartUsesMainlineCommands(t *testing.T) {}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/validation -run 'TestDocs_' -count=1`
Expected: FAIL because user docs still describe the old public surface.

**Step 3: Write minimal implementation**

- Rewrite README quick start to center on:
  - `/mp:init`
  - `/mp:brainstorm`
  - `/mp:design`
  - `/mp:plan`
  - `/mp:execute`
  - `/mp:debug`
  - `/mp:debate`
- Rewrite workflow docs so they no longer market `discover/define/develop/deliver` as public commands.
- Add a dedicated `docs/architecture/multi-mainline.md` explaining:
  - upstream superpowers source-of-truth
  - thin wrappers
  - fixed roles
  - phase model policy
  - debate all-model fanout
- Archive or remove stale user-facing docs that would mislead new users.

**Step 4: Run test to verify it passes**

Run: `go test ./internal/validation -run 'TestDocs_' -count=1`
Expected: PASS

**Step 5: Commit**

```bash
git add README.md docs/WORKFLOW-SKILLS.md docs/COMMAND-REFERENCE.md docs/PLUGIN-ARCHITECTURE.md docs/CLI-REFERENCE.md docs/architecture/multi-mainline.md docs/archive/README.md internal/validation
git commit -m "multi(mainline-20260308): group g10 - rewrite docs for mainline surface"
```

---

### Task 11: Run Full Verification and Record the New Surface Contract

**Why:** This branch intentionally deletes public surface area and rewires runtime behavior; it needs strong evidence before any merge discussion.

**What:** Run the focused package tests, rebuild runtime assets, and confirm the published public surface matches the approved command list.

**How:** Verify the generator, routing, policy, hooks, and plugin manifest together before claiming the branch is ready.

**Key Design:** Evidence before claims; do not rely on partial package passes.

**Files:**
- Verify: `.claude-plugin/plugin.json`
- Verify: `.claude-plugin/.claude/commands/*`
- Verify: `.claude-plugin/.claude/skills/*`
- Verify: `config/workflows.yaml`
- Verify: `config/roles.yaml`
- Verify: `docs/architecture/multi-mainline.md`
- Create: `docs/plans/evidence/multi-mainline/2026-03-08-verification.md`

**Step 1: Run focused verification**

Run: `go test ./internal/devx ./internal/roles ./internal/policy ./internal/providers ./internal/workflows ./internal/hooks ./internal/cli ./internal/validation ./cmd/mp-devx -count=1`
Expected: PASS

**Step 2: Rebuild runtime assets**

Run: `go run ./cmd/mp-devx --action build-runtime`
Expected: PASS and regenerated plugin assets with no missing file errors.

**Step 3: Run broader verification**

Run: `go test ./... -count=1`
Expected: PASS, or if unrelated pre-existing failures exist, capture them explicitly in the evidence note and stop before making success claims.

**Step 4: Record verification evidence**

- Save commands, exit status, and any notable output under `docs/plans/evidence/multi-mainline/2026-03-08-verification.md`.
- Include the final allowed public command list and final wrapper skill list.

**Step 5: Commit**

```bash
git add docs/plans/evidence/multi-mainline/2026-03-08-verification.md .claude-plugin/plugin.json .claude-plugin/.claude/commands .claude-plugin/.claude/skills config/workflows.yaml config/roles.yaml docs/architecture/multi-mainline.md
git commit -m "multi(mainline-20260308): group g11 - verify mainline branch reset"
```
