# Multipowers Gap Analysis Plan 8 (Intention → Enforced Execution)

> Goal: Based on `conductor/context` (tool project background), identify remaining gaps between project intention and current executable codebase, then provide an implementation plan with vibe-coding updatable task statuses.

## 1) Analysis Baseline

### 1.1 Intention Sources (Tool Project Background)
- `conductor/context/product.md`
- `conductor/context/product-vision.md`
- `conductor/context/product-guidelines.md`
- `conductor/context/workflow.md`
- `conductor/context/tech-stack.md`

### 1.2 Codebase Scope (for gap analysis)
Analyzed implementation scope excludes:
- `/conductor`
- `/docs`
- `/templates`
- `/outputs`

Included key implementation paths:
- `bin/*`
- `scripts/*`
- `connectors/*`
- `config/*`
- `.opencode/plugins/*`
- `lib/*`
- `tests/opencode/*`
- `README.md`, `package.json`

---

## 2) Intention-to-Reality Gap Summary

| Gap ID | Intended Capability | Current Code Evidence | Risk | Priority |
|---|---|---|---|---|
| GA8-01 | Router should choose lane and execute path | `multipowers route` returns decision only; execution requires separate commands | Routing remains partially manual | P0 |
| GA8-02 | Execution should be track-first | `route/workflow` accept optional `--track-id`, no active-track enforcement | Work can run outside track lifecycle | P0 |
| GA8-03 | Standard lane should be workflow-first with valid role graph | `execute_workflow.py` validates structure but not role existence against effective roles config | Late runtime failures (`ask-role` fails mid-workflow) | P0 |
| GA8-04 | Major-change governance should be enforced, not optional | `run_governance_checks.sh` defaults to warn/skip when tools unavailable | Governance signal can be bypassed silently | P1 |
| GA8-05 | Governance should leave machine-readable artifacts | Governance script emits terminal output only; no artifact record | Hard to audit/trace compliance | P1 |
| GA8-06 | Completion should require governance artifacts | `track complete` has no governance gate | “Done” may bypass policy | P1 |
| GA8-07 | Evidence checks should verify governance evidence for major work | `check_plan_evidence.py` checks generic fields, not governance proof | Evidence quality gap for major changes | P1 |
| GA8-08 | Docs sync should be component-aware | `check_docs_sync.py` only checks doc-vs-non-doc presence | Minimal docs update can pass without meaningful alignment | P1 |
| GA8-09 | Routing observability should cover full lifecycle | Existing events cover lane/workflow but weak coverage for fast-lane execution and governance stage | Partial timeline visibility | P2 |
| GA8-10 | Operational CLI should avoid placeholder behavior | `multipowers update` still placeholder | Tool lifecycle maintenance incomplete | P2 |
| GA8-11 | Template sync rule should have automation support | No checker for `conductor`-to-`templates` sync candidates | Drift accumulates without reminders | P2 |

---

## 3) Task Board (Vibe-Coding Status Updatable)

Status values: `TODO` / `IN_PROGRESS` / `BLOCKED` / `DONE`

Use status updater script:
- `python3 scripts/update_plan_task_status.py --file docs/plans/gap_analysis_plan8.md --task-id T8-001 --status IN_PROGRESS`

| Task ID | Status | Priority | Owner | Depends On |
|---|---|---|---|---|
| T8-001 | `DONE` | P0 | Integrator | - |
| T8-002 | `DONE` | P0 | Integrator | T8-001 |
| T8-003 | `DONE` | P0 | Integrator | T8-001 |
| T8-004 | `DONE` | P1 | Integrator | T8-001 |
| T8-005 | `DONE` | P1 | Integrator | T8-004 |
| T8-006 | `DONE` | P1 | Integrator | T8-004, T8-005 |
| T8-007 | `DONE` | P1 | Integrator | T8-005 |
| T8-008 | `DONE` | P1 | Integrator | T8-004 |
| T8-009 | `DONE` | P2 | Integrator | T8-001, T8-005 |
| T8-010 | `DONE` | P2 | Integrator | - |
| T8-011 | `DONE` | P2 | Integrator | T8-008 |

---

### T8-001 (P0) Add one-command router execution entry (`multipowers run`)

- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator
- **Why**:
  - Intention says router selects a lane and drives execution path.
  - Current implementation splits decision (`route`) and execution (`workflow run` / `ask-role`).
- **What to implement**:
  - Modify: `bin/multipowers`
  - Add: `tests/opencode/test-router-run-command.sh`
- **How to implement**:
  1. Add new subcommand `multipowers run --task <text> [--risk-hint ...] [--track-id ...] [--json]`.
  2. Internally call route logic, then:
     - fast lane → dispatch `ask-role` with selected role
     - standard lane → invoke `workflow run <selected_workflow>`
  3. Ensure `request_id` is generated once and propagated end-to-end.
  4. Return structured JSON in `--json` mode including lane + execution result.
- **Acceptance criteria**:
  - [ ] One command completes route + execution path.
  - [ ] Fast and standard lanes both supported.
  - [ ] `request_id` is consistent across generated logs.
- **Verification commands**:
  - `bash tests/opencode/test-router-run-command.sh`
  - `./bin/multipowers run --task "fix typo in readme" --json`

---

### T8-002 (P0) Enforce active-track context for execution commands

- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator
- **Why**:
  - Workflow baseline requires setup + track lifecycle before execution.
  - Current commands allow untracked execution by default.
- **What to implement**:
  - Modify: `bin/multipowers`
  - Add: `tests/opencode/test-active-track-enforcement.sh`
- **How to implement**:
  1. Persist active track marker when `track start` succeeds (e.g., `conductor/.active_track`).
  2. Require active track for `run`/`workflow run` unless `--allow-untracked` is provided.
  3. Add actionable error message with next steps (`track new/start`).
- **Acceptance criteria**:
  - [ ] Execution commands fail clearly when no active track is present.
  - [ ] `--allow-untracked` enables explicit bypass.
  - [ ] Track start updates active track marker deterministically.
- **Verification commands**:
  - `bash tests/opencode/test-active-track-enforcement.sh`

---

### T8-003 (P0) Add workflow role preflight validation against effective roles config

- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator
- **Why**:
  - Standard lane should use valid role graph; invalid node roles should fail before execution.
- **What to implement**:
  - Modify: `scripts/execute_workflow.py`
  - Modify: `bin/multipowers`
  - Add: `tests/opencode/test-workflow-role-preflight.sh`
- **How to implement**:
  1. Resolve effective roles config (`conductor/config/roles.json` fallback to default).
  2. Before executing nodes, ensure every role exists in roles config.
  3. Emit clear diagnostics (`unknown role`, `workflow`, `node id`).
  4. Expose this check from `multipowers workflow validate`.
- **Acceptance criteria**:
  - [ ] Unknown node role fails before any node execution.
  - [ ] Error message includes workflow and node identifiers.
  - [ ] Valid workflow still runs unchanged.
- **Verification commands**:
  - `bash tests/opencode/test-workflow-role-preflight.sh`
  - `./bin/multipowers workflow validate`

---

### T8-004 (P1) Strengthen governance execution policy (strict by design)

- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator
- **Why**:
  - Major-change policy requires semgrep/biome/ruff checks; warn-only behavior weakens enforcement.
- **What to implement**:
  - Modify: `scripts/run_governance_checks.sh`
  - Modify: `package.json`
  - Add: `tests/opencode/test-governance-strictness.sh`
- **How to implement**:
  1. Introduce explicit modes:
     - `--mode strict` (default in CI / governance command)
     - `--mode advisory` (local opt-in)
  2. In strict mode, missing tool or execution failure returns non-zero.
  3. Keep actionable hints for missing tools and remediation.
- **Acceptance criteria**:
  - [ ] Strict mode blocks on missing/failed tools.
  - [ ] Advisory mode preserves developer velocity when needed.
  - [ ] `npm run governance` uses strict mode by default.
- **Verification commands**:
  - `bash tests/opencode/test-governance-strictness.sh`
  - `npm run governance -- --help`

---

### T8-005 (P1) Emit governance artifacts for audit and traceability

- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator
- **Why**:
  - Policy requires evidence; terminal-only output is hard to audit.
- **What to implement**:
  - Modify: `scripts/run_governance_checks.sh`
  - Add: `scripts/record_governance_artifact.py`
  - Add: `tests/opencode/test-governance-artifact.sh`
- **How to implement**:
  1. Add `--artifact <path>` option to governance script.
  2. Record JSON artifact including:
     - timestamp, mode, changed_files
     - tool execution result per tool
     - overall exit_code and summary
  3. Include optional `request_id`/`track_id` fields when provided.
- **Acceptance criteria**:
  - [ ] Governance run can output machine-readable artifact.
  - [ ] Artifact includes per-tool pass/fail status.
  - [ ] Artifact path is deterministic and script-friendly.
- **Verification commands**:
  - `bash tests/opencode/test-governance-artifact.sh`

---

### T8-006 (P1) Gate `track complete` with governance evidence

- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator
- **Why**:
  - Completion should not bypass major-change governance expectations.
- **What to implement**:
  - Modify: `bin/multipowers`
  - Modify: `scripts/run_governance_checks.sh`
  - Add: `tests/opencode/test-track-complete-governance-gate.sh`
- **How to implement**:
  1. Before `track complete`, require recent governance artifact or run governance checks inline.
  2. Add explicit bypass flag `--skip-governance` with warning log.
  3. Persist completion metadata (track file footer) including governance status.
- **Acceptance criteria**:
  - [ ] Track completion is blocked without governance evidence by default.
  - [ ] Bypass is explicit and auditable.
  - [ ] Completion metadata includes governance outcome.
- **Verification commands**:
  - `bash tests/opencode/test-track-complete-governance-gate.sh`

---

### T8-007 (P1) Extend plan evidence checker with governance evidence rule

- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator
- **Why**:
  - Generic evidence fields are necessary but insufficient for major-change policy compliance.
- **What to implement**:
  - Modify: `scripts/check_plan_evidence.py`
  - Modify: `docs/plans/status-evidence-template.md`
  - Add: `tests/opencode/test-plan-governance-evidence.sh`
- **How to implement**:
  1. Add optional flag `--require-governance-evidence`.
  2. Under this flag, require one of:
     - governance command in evidence (`run_governance_checks.sh`, `semgrep`, `biome`, `ruff`)
     - governance artifact reference.
  3. Keep backward-compatible default behavior for legacy plans.
- **Acceptance criteria**:
  - [ ] Governance evidence can be enforced for selected plan sets.
  - [ ] Legacy plans still validate under default mode.
  - [ ] Diagnostics clearly indicate missing governance proof.
- **Verification commands**:
  - `bash tests/opencode/test-plan-governance-evidence.sh`
  - `python3 scripts/check_plan_evidence.py --help`

---

### T8-008 (P1) Upgrade docs-sync checker to component-aware mapping

- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator
- **Why**:
  - “Any docs changed” is too weak for meaningful intention-code consistency.
- **What to implement**:
  - Modify: `scripts/check_docs_sync.py`
  - Add: `config/docs-sync-rules.json`
  - Add: `tests/opencode/test-docs-sync-mapping.sh`
- **How to implement**:
  1. Define prefix-based rules (e.g., `bin/` => `README.md`, `connectors/` => logging/connector docs).
  2. Require mapped docs updates when code paths change.
  3. Keep ignores for test-only and non-behavioral changes.
- **Acceptance criteria**:
  - [ ] Missing mapped docs update fails with actionable output.
  - [ ] Correct docs updates pass.
  - [ ] Rule config is easy to extend.
- **Verification commands**:
  - `bash tests/opencode/test-docs-sync-mapping.sh`

---

### T8-009 (P2) Complete runtime observability for fast-lane + governance lifecycle

- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator
- **Why**:
  - Current events are strong for workflow path but weaker for fast-lane execution and governance stage visibility.
- **What to implement**:
  - Modify: `bin/multipowers`
  - Modify: `connectors/utils.py`
  - Add: `tests/opencode/test-observability-lifecycle.sh`
- **How to implement**:
  1. Add events: `fast_lane_dispatched`, `fast_lane_finished`, `governance_started`, `governance_finished`.
  2. Require `request_id` carry-through for these events.
  3. Validate event order and mandatory fields in tests.
- **Acceptance criteria**:
  - [ ] Full lifecycle trace available in JSONL for both lanes.
  - [ ] Event order is deterministic and testable.
  - [ ] Governance stage is visible in runtime logs.
- **Verification commands**:
  - `bash tests/opencode/test-observability-lifecycle.sh`

---

### T8-010 (P2) Implement safe `multipowers update` workflow

- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator
- **Why**:
  - CLI advertises `update` but currently returns placeholder.
- **What to implement**:
  - Modify: `bin/multipowers`
  - Add: `scripts/check_update_state.sh`
  - Add: `tests/opencode/test-update-command.sh`
- **How to implement**:
  1. Implement non-destructive update check (`status`, `behind`, `dirty working tree`).
  2. Support `update --check` and `update --apply` (with explicit confirmation).
  3. Keep clear guidance for manual fallback paths.
- **Acceptance criteria**:
  - [ ] `update --check` reports actionable status.
  - [ ] `update --apply` is explicit and safe.
  - [ ] Dirty tree handling is clear and non-destructive.
- **Verification commands**:
  - `bash tests/opencode/test-update-command.sh`

---

### T8-011 (P2) Add template-sync candidate checker for maintainer workflow

- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator
- **Why**:
  - Context workflow requires evaluating safe sync from maintainer conventions into templates.
- **What to implement**:
  - Add: `scripts/check_template_sync_candidates.py`
  - Modify: `bin/multipowers` (doctor warn-only integration)
  - Add: `tests/opencode/test-template-sync-candidates.sh`
- **How to implement**:
  1. Compare selected maintainer policy files with template equivalents using configurable mapping.
  2. Report drift summary and suggested sync targets.
  3. Add doctor warning output (non-blocking) when drift exceeds threshold.
- **Acceptance criteria**:
  - [ ] Script reports candidate sync diffs with file-level granularity.
  - [ ] Doctor surfaces drift status as warning, not hard fail.
  - [ ] Tests cover no-drift and drift scenarios.
- **Verification commands**:
  - `bash tests/opencode/test-template-sync-candidates.sh`
  - `python3 scripts/check_template_sync_candidates.py --help`

---

## 4) Recommended Execution Sequence

### Phase A — Routing + Workflow Enforcement
1. T8-001
2. T8-002
3. T8-003

### Phase B — Governance as Gate
4. T8-004
5. T8-005
6. T8-006
7. T8-007
8. T8-008

### Phase C — Ops + Observability + Maintenance
9. T8-009
10. T8-010
11. T8-011

---

## 5) Definition of Done (Plan-Level)

1. Router supports one-command lane decision + execution.
2. Standard-lane workflows fail fast on invalid role graphs.
3. Track completion is governance-gated by default.
4. Governance produces structured artifacts and is enforceable.
5. Evidence validation supports governance-proof checks.
6. Docs sync is component-aware, not only file-type-aware.
7. Observability covers fast lane, standard lane, and governance stages.
8. `update` command has safe, non-placeholder behavior.
9. Template sync candidates are auditable by script.

---

## 6) Evidence Section (fill when tasks are DONE)

- **Coverage Task IDs**: `T8-001, T8-002, T8-003, T8-004, T8-005, T8-006, T8-007, T8-008, T8-009, T8-010, T8-011`
- **Date**: `2026-02-10`
- **Verifier**: `Codex CLI`
- **Command(s)**:
  - `bash tests/opencode/run-tests.sh`
  - `bash tests/opencode/test-router-run-command.sh`
  - `bash tests/opencode/test-track-complete-governance-gate.sh`
  - `python3 scripts/check_plan_evidence.py --require-governance-evidence docs/plans/gap_analysis_plan8.md`
  - `python3 scripts/check_docs_sync.py --changed-file bin/multipowers --changed-file README.md`
  - `bash scripts/run_governance_checks.sh --mode advisory --changed-file bin/multipowers --changed-file README.md`
- **Exit Code**: `0`
- **Key Output**:
  - `STATUS: PASSED (34 passed, 0 failed)`
  - `[PLAN-EVIDENCE] PASS (1 files)`
