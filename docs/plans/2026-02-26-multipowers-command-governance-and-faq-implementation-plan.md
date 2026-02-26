# Multipowers Command Governance & Auto-FAQ Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use executing-plans to implement this plan task-by-task.

**Goal:** Implement standardized spec-command preflight governance, runtime fail-fast contract, and fully automatic `CLAUDE.md` + `FAQ.md` knowledge loop for target projects.

**Architecture:** Keep upstream conflict surface minimal by placing new logic in `custom/*` and using thin wiring in `scripts/orchestrate.sh`. Centralize spec-driven preflight and FAQ synthesis in reusable shell helpers; ensure all target-project writes stay under `/.multipowers/`.

**Tech Stack:** Bash (`scripts/orchestrate.sh`, `custom/lib/*.sh`), Markdown templates (`custom/templates/*`), shell tests (`tests/unit/*.sh`, `tests/integration/*.sh`).

---

## Execution Tracker (Tick as you go)

- [x] Phase 0 - Baseline and safety checkpoint
- [x] Phase 1 - Template assets for target-project `CLAUDE.md`/`FAQ.md`
- [x] Phase 2 - Auto-init generation wiring
- [x] Phase 3 - Spec-command Step 0 preflight hardening
- [x] Phase 4 - Runtime preconditions fail-fast enforcement (all execution paths)
- [x] Phase 5 - Debate provider quorum + fallback enforcement
- [x] Phase 6 - Auto-FAQ synthesis, dedup, rewrite
- [x] Phase 7 - Docs update (tool-project + target-project)
- [x] Phase 8 - Verification and release handoff

---

### Phase 0: Baseline and Safety Checkpoint

**Files:**
- Modify: none

**Steps:**
- [x] P0.1 Confirm branch and working tree status
  - Run: `git rev-parse --abbrev-ref HEAD && git status --short`
  - Expected: on `multipowers`; existing unrelated changes are visible and untouched.
- [x] P0.2 Capture current head for rollback reference
  - Run: `git rev-parse --short HEAD`
  - Record in notes.
- [x] P0.3 Verify existing syntax baseline
  - Run: `bash -n scripts/orchestrate.sh`
  - Expected: no syntax error.

**Commit checkpoint:** none (informational baseline)

---

### Phase 1: Add Target Project Templates (`CLAUDE.md`, `FAQ.md`)

**Files:**
- Create: `custom/templates/CLAUDE.md`
- Create: `custom/templates/FAQ.md`
- Modify: `custom/docs/reference/config-schema.md` (template section)

**Steps:**
- [x] P1.1 Write `custom/templates/CLAUDE.md` as language-agnostic contract
  - Include sections: scope/priority, security baseline, testing gates, runtime fail-fast, FAQ reference.
- [x] P1.2 Write `custom/templates/FAQ.md` scaffold with error-type sections
  - Sections only by error type (not by command).
- [x] P1.3 Add schema docs for new generated files
  - Update paths and ownership: tool template -> target generated artifact.
- [x] P1.4 Verify templates are discoverable
  - Run: `find custom/templates -maxdepth 2 -type f | sort`
  - Expected: includes `custom/templates/CLAUDE.md` and `custom/templates/FAQ.md`.

**Commit checkpoint:**
- [x] `git add custom/templates/CLAUDE.md custom/templates/FAQ.md custom/docs/reference/config-schema.md && git commit -m "feat(templates): add CLAUDE and FAQ target-project templates"`

---

### Phase 2: Wire `/octo:init` to Generate Both Files in Target `/.multipowers/`

**Files:**
- Modify: `scripts/orchestrate.sh`
- Modify: `custom/config/setup.toml`
- Modify: `custom/lib/conductor-context.sh`
- Test: `tests/unit/test-octo-init-render.sh`
- Create: `tests/unit/test-claude-faq-init-render.sh`

**Steps:**
- [x] P2.1 Extend init render path to copy/render `CLAUDE.md` to `"$PROJECT_ROOT/.multipowers/CLAUDE.md"`
- [x] P2.2 Extend init render path to copy/render `FAQ.md` to `"$PROJECT_ROOT/.multipowers/FAQ.md"`
- [x] P2.3 Ensure no writes happen outside target project root
  - Guard with absolute path checks before write.
- [x] P2.4 Add unit test: init creates both files in target cwd
  - New test asserts presence under `/.multipowers/` and absence under tool project and `$HOME`.
- [x] P2.5 Run unit tests
  - Run: `bash tests/unit/test-octo-init-render.sh && bash tests/unit/test-claude-faq-init-render.sh`
  - Expected: PASS both.

**Commit checkpoint:**
- [x] `git add scripts/orchestrate.sh custom/config/setup.toml custom/lib/conductor-context.sh tests/unit/test-octo-init-render.sh tests/unit/test-claude-faq-init-render.sh && git commit -m "feat(init): generate CLAUDE and FAQ in target .multipowers"`

---

### Phase 3: Harden Spec-Driven Step 0 Preflight (Single Contract)

**Files:**
- Modify: `scripts/orchestrate.sh`
- Modify: `custom/lib/conductor-context.sh`
- Modify: `.claude/commands/plan.md`
- Modify: `.claude/commands/discover.md`
- Modify: `.claude/commands/define.md`
- Modify: `.claude/commands/develop.md`
- Modify: `.claude/commands/deliver.md`
- Modify: `.claude/commands/embrace.md`
- Modify: `.claude/commands/research.md`
- Test: `tests/unit/test-conductor-context-guard.sh`
- Test: `tests/integration/test-spec-commands-auto-init.sh`

**Steps:**
- [x] P3.1 Ensure Step 0 always runs before command-specific logic in orchestrator
- [x] P3.2 Keep command markdown guard text minimal and front-loaded (no duplicated tail sections)
- [x] P3.3 Align required context file list with `.multipowers` contract
- [x] P3.4 Update auto-init integration test expected messages to `.multipowers`
- [x] P3.5 Run tests
  - Run: `bash tests/unit/test-conductor-context-guard.sh && bash tests/integration/test-spec-commands-auto-init.sh`
  - Expected: PASS.

**Commit checkpoint:**
- [x] `git add scripts/orchestrate.sh custom/lib/conductor-context.sh .claude/commands/plan.md .claude/commands/discover.md .claude/commands/define.md .claude/commands/develop.md .claude/commands/deliver.md .claude/commands/embrace.md .claude/commands/research.md tests/unit/test-conductor-context-guard.sh tests/integration/test-spec-commands-auto-init.sh && git commit -m "fix(spec-commands): enforce unified Step-0 multipowers preflight"`

---

### Phase 4: Enforce Runtime Preconditions (`fail-fast`) in All Execution Paths

**Files:**
- Modify: `scripts/orchestrate.sh`
- Modify: `custom/docs/target-project/getting-started.md`
- Modify: `custom/docs/target-project/troubleshooting.md`
- Create: `tests/unit/test-runtime-prerun-failfast.sh`

**Steps:**
- [x] P4.1 Audit all command/provider execution entrypoints for missing `apply_pre_run_context`
- [x] P4.2 Ensure every path reads `/.multipowers/context/runtime.json`
- [x] P4.3 Enforce fail-fast semantics on pre-run command failures
- [x] P4.4 Add unit test with intentional failing pre-run command
  - Expected: command exits non-zero before provider call.
- [x] P4.5 Run tests
  - Run: `bash tests/unit/test-runtime-prerun-failfast.sh`
  - Expected: PASS.

**Commit checkpoint:**
- [x] `git add scripts/orchestrate.sh custom/docs/target-project/getting-started.md custom/docs/target-project/troubleshooting.md tests/unit/test-runtime-prerun-failfast.sh && git commit -m "feat(runtime): enforce fail-fast pre-run contract across execution paths"`

---

### Phase 5: Debate Provider Quorum/Fallback Guarantees

**Files:**
- Modify: `scripts/orchestrate.sh`
- Modify: `tests/unit/test-debate-routing.sh`
- Modify: `tests/integration/test-debate-integration.sh`
- Modify: `custom/docs/customizations/proxy-routing.md`

**Steps:**
- [x] P5.1 Validate debate pipeline rule: up to 3 providers, continue with 2, fail under 2
- [x] P5.2 Ensure synthesis fallback chain is deterministic and logged
- [x] P5.3 Ensure Codex/Gemini proxy env is injected consistently in debate paths
- [x] P5.4 Extend tests for two scenarios:
  - one provider unavailable -> workflow still completes;
  - two providers unavailable -> explicit fail with remediation.
- [x] P5.5 Run tests
  - Run: `bash tests/unit/test-debate-routing.sh && bash tests/integration/test-debate-integration.sh`
  - Expected: PASS.

**Commit checkpoint:**
- [x] `git add scripts/orchestrate.sh tests/unit/test-debate-routing.sh tests/integration/test-debate-integration.sh custom/docs/customizations/proxy-routing.md && git commit -m "fix(debate): enforce >=2 provider quorum with consistent fallback and proxy env"`

---

### Phase 6: Implement Auto-FAQ Synthesis (No Manual Maintenance)

**Files:**
- Create: `custom/lib/faq-synthesizer.sh`
- Modify: `scripts/orchestrate.sh`
- Create: `tests/unit/test-faq-synthesizer.sh`
- Create: `tests/integration/test-faq-auto-update.sh`

**Steps:**
- [x] P6.1 Define structured failure-event schema in shell-friendly format
  - fields: `timestamp`, `command`, `error_type`, `root_cause`, `suggested_fix`, `provider`, `exit_code`.
- [x] P6.2 Emit events from key failure points (init/context/runtime/provider/debate)
- [x] P6.3 Build synthesizer pipeline:
  - ingest events -> classify by error type -> dedup key (`error_type + normalized_root_cause + normalized_fix`) -> rewrite FAQ.
- [x] P6.4 Enforce FAQ max entries cap (e.g. 120 by recency/frequency)
- [x] P6.5 Ensure no backup/archive files are created
- [x] P6.6 Add deterministic tests for dedup and rewrite behavior
- [x] P6.7 Run tests
  - Run: `bash tests/unit/test-faq-synthesizer.sh && bash tests/integration/test-faq-auto-update.sh`
  - Expected: PASS.

**Commit checkpoint:**
- [x] `git add custom/lib/faq-synthesizer.sh scripts/orchestrate.sh tests/unit/test-faq-synthesizer.sh tests/integration/test-faq-auto-update.sh && git commit -m "feat(faq): auto-generate deduped multipowers FAQ by error type"`

---

### Phase 7: Update Documentation (Tool + Target Views)

**Files:**
- Modify: `custom/docs/tool-project/README.md`
- Modify: `custom/docs/tool-project/getting-started.md`
- Modify: `custom/docs/tool-project/troubleshooting.md`
- Modify: `custom/docs/target-project/README.md`
- Modify: `custom/docs/target-project/getting-started.md`
- Modify: `custom/docs/target-project/troubleshooting.md`
- Modify: `custom/docs/customizations/conductor-context.md`
- Modify: `custom/docs/reference/faq.md`

**Steps:**
- [x] P7.1 Document `CLAUDE.md` template ownership and target generation behavior
- [x] P7.2 Document FAQ auto-update lifecycle and error-type categories
- [x] P7.3 Document no-manual-maintenance policy for FAQ
- [x] P7.4 Document strict target-output boundaries and fail-fast runtime contract
- [x] P7.5 Run docs consistency check
  - Run: `bash tests/unit/test-docs-sync.sh`
  - Expected: PASS.

**Commit checkpoint:**
- [x] `git add custom/docs/tool-project/README.md custom/docs/tool-project/getting-started.md custom/docs/tool-project/troubleshooting.md custom/docs/target-project/README.md custom/docs/target-project/getting-started.md custom/docs/target-project/troubleshooting.md custom/docs/customizations/conductor-context.md custom/docs/reference/faq.md && git commit -m "docs: add CLAUDE/FAQ governance and auto-learning workflow"`

---

### Phase 8: Final Verification and Push

**Files:**
- Modify: none

**Steps:**
- [x] P8.1 Run syntax checks
  - Run: `bash -n scripts/orchestrate.sh && bash tests/smoke/test-syntax.sh`
- [x] P8.2 Run focused test suite
  - Run: `bash tests/unit/test-conductor-context-guard.sh`
  - Run: `bash tests/unit/test-octo-init-render.sh`
  - Run: `bash tests/unit/test-runtime-prerun-failfast.sh`
  - Run: `bash tests/unit/test-faq-synthesizer.sh`
  - Run: `bash tests/unit/test-debate-routing.sh`
  - Run: `bash tests/integration/test-spec-commands-auto-init.sh`
  - Run: `bash tests/integration/test-faq-auto-update.sh`
  - Run: `bash tests/integration/test-debate-integration.sh`
- [x] P8.3 Ensure no forbidden output paths were touched during test runs
  - Run: `test ! -d "$HOME/.claude-octopus"`
  - Run: `rg -n "\.claude-octopus|/home/.*/\.claude-octopus" . || true`
- [x] P8.4 Push branch
  - Run: `git push origin multipowers`

**Completion criteria:**
- [x] All phase checkboxes completed
- [x] All listed tests pass
- [x] Target-project output boundary requirements satisfied
- [x] Docs updated for both audiences

Implementation note:
- `tests/unit/test-docs-sync.sh` reports existing baseline documentation/version issues unrelated to this change set.
- `tests/unit/test-debate-routing.sh` depends on a missing helper path in current repo layout (`tests/helpers/test-framework.sh`).

---

## Risks & Mitigations
- Risk: upstream drift in `.claude/commands/*`
  - Mitigation: keep only minimal Step 0 delta; avoid broad rewrites.
- Risk: FAQ growth/noise
  - Mitigation: strict dedup key + bounded entry cap + rewrite mode.
- Risk: hidden execution path bypasses runtime pre-run
  - Mitigation: explicit execution-path audit + fail-fast test coverage.

## Rollback Plan
- Revert by phase commit if regression appears.
- If runtime pre-run blocks valid flows unexpectedly, temporarily disable via controlled flag in `runtime.json` while keeping fail-fast default documented.
