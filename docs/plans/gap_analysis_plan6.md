# Multipowers Gap Analysis Plan 6 (Intention vs Executable Codebase)

> **For Claude:** REQUIRED SUB-SKILL: Use `superpowers:executing-plans` to implement this plan task-by-task.

**Goal:** Close the highest-impact gaps between the tool-project intention in `conductor/context/*.md` and the current executable codebase.

**Architecture:** Implement a real Router routing/execution layer on top of existing role dispatch (`bin/ask-role`), then complete missing runtime foundations (Claude connector + MCP wiring + plugin runtime + governance automation + observability).

**Tech Stack:** Bash, Python, Node.js, role connectors (`codex` / `gemini` / `claude`), OpenCode plugin runtime.

---

## 1) Analysis Scope and Method

This gap analysis compares:

- **Intention baseline:** `conductor/context/product.md`, `product-guidelines.md`, `product-vision.md`, `workflow.md`, `tech-stack.md`
- **Actual implementation scope:** all code **except** `conductor/`, `docs/`, `templates/`, `outputs/`

Primary implementation files reviewed:

- Runtime CLIs: `bin/multipowers`, `bin/ask-role`
- Role and schema config: `config/roles.default.json`, `config/roles.schema.json`, `config/mcp.default.json`
- Connectors/logging: `connectors/codex.py`, `connectors/gemini.py`, `connectors/utils.py`
- Hooks/commands/skill plumbing: `hooks/session-start.sh`, `commands/*.md`, `lib/skills-core.js`
- Validation and tests: `scripts/*.py`, `tests/opencode/*.sh`, `tests/claude-code/*.sh`

---

## 2) Intention-to-Reality Gap Summary

| Gap ID | Intended Capability | Current Code Evidence | Risk | Priority |
|---|---|---|---|---|
| GA6-01 | Router chooses Fast vs Standard lane programmatically | `bin/multipowers` exposes only `init|doctor|update|track`; no routing command | Workflow intent remains manual, inconsistent execution | P0 |
| GA6-02 | Standard lane is workflow-first and node-role aware | Only command stubs (`commands/*.md`) and skill docs; no executable workflow engine | Role switching logic is not enforceable or observable | P0 |
| GA6-03 | Multi-CLI role execution includes `claude` | `config/roles.schema.json` tool enum excludes `claude`; no `connectors/claude.py` | Model coverage gap vs stated architecture | P0 |
| GA6-04 | MCP is a first-class tool capability | `config/mcp.default.json` exists but has no runtime references | MCP promise is non-functional | P0 |
| GA6-05 | Major-change governance enforces `semgrep` + `biome` + `ruff` + docs sync | `package.json` scripts and test runners do not execute these gates | High-risk changes can ship without required controls | P1 |
| GA6-06 | Plugin runtime is part of deliverable, not test fallback | `.opencode/plugins/superpowers.js` missing; `tests/opencode/setup.sh` silently creates fallback stub | Integration confidence is overstated | P1 |
| GA6-07 | Task status updates by vibe coding are standardized and machine-checkable | `scripts/check_plan_evidence.py` uses narrow status regex and file-scope assumptions | Status/evidence automation is brittle across plan formats | P1 |
| GA6-08 | Routing/workflow decisions are observable end-to-end | Logs mainly capture context + connector execution, not routing/workflow-node events | Hard to audit why a lane/role/workflow was chosen | P2 |

---

## 3) Task Board (Vibe-Coding Status Updatable)

Status values: `TODO` / `IN_PROGRESS` / `BLOCKED` / `DONE`

> Update both `Status` and `状态` fields in each task to keep compatibility with existing evidence checker conventions.

| Task ID | Status | Priority | Owner | Depends On |
|---|---|---|---|---|
| T6-001 | `TODO` | P0 | Router | - |
| T6-002 | `TODO` | P0 | Router + Coder | T6-001 |
| T6-003 | `TODO` | P0 | Coder | - |
| T6-004 | `TODO` | P0 | Architect + Coder | T6-001 |
| T6-005 | `TODO` | P1 | Architect + Coder | T6-002, T6-003 |
| T6-006 | `TODO` | P1 | Coder | - |
| T6-007 | `TODO` | P1 | Architect | T6-005 |
| T6-008 | `TODO` | P2 | Router | T6-001, T6-002 |

---

### T6-001 (P0) Implement executable Router lane routing

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Router
- **Why**: The project intention requires deterministic lane selection, but current routing is manual and undocumented in runtime behavior.
- **What to implement**:
  - Modify: `bin/multipowers`
  - Add: `scripts/route_task.py`
  - Add: `tests/opencode/test-routing-lanes.sh`
  - Update: `README.md`
- **How to implement**:
  1. Add `multipowers route` subcommand with inputs: task text, optional risk hint, optional force lane.
  2. Implement deterministic routing heuristics in `scripts/route_task.py` (Fast vs Standard).
  3. Return machine-readable output (`json` mode) including lane, reason, suggested workflow/role.
  4. Add unit-style shell tests for representative routing scenarios.
- **Acceptance criteria**:
  - [ ] `multipowers route` exists and is documented.
  - [ ] Same input yields same lane and reason.
  - [ ] Tests cover both lanes and override behavior.
- **Verification commands**:
  - `bash tests/opencode/test-routing-lanes.sh`
  - `bash -n bin/multipowers`

### T6-002 (P0) Add standard-lane workflow execution engine

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Router + Coder
- **Why**: Standard lane is currently conceptual only; there is no executable flow that enforces workflow-first and node-level role switching.
- **What to implement**:
  - Add: `config/workflows.default.json`
  - Modify: `bin/multipowers`
  - Add: `scripts/execute_workflow.py`
  - Add: `tests/opencode/test-workflow-engine.sh`
- **How to implement**:
  1. Define workflow spec format: workflow name, default executor role, ordered nodes, node-level role overrides.
  2. Add `multipowers workflow run <workflow> --task <...>` command.
  3. Resolve workflow config (project override then default), then execute nodes via `bin/ask-role`.
  4. Fail fast on missing nodes/roles and surface actionable errors.
  5. Add tests for node override behavior (e.g., Architect-only review node).
- **Acceptance criteria**:
  - [ ] Workflows are configured, not hardcoded.
  - [ ] Node-level role overrides work and are test-covered.
  - [ ] Failure paths include clear diagnostics.
- **Verification commands**:
  - `bash tests/opencode/test-workflow-engine.sh`
  - `python3 -m py_compile scripts/execute_workflow.py`

### T6-003 (P0) Add Claude CLI connector and schema support

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Coder
- **Why**: The intended multi-CLI architecture includes `claude`, but runtime and schema only support `gemini`/`codex`/`system`.
- **What to implement**:
  - Add: `connectors/claude.py`
  - Modify: `bin/ask-role`
  - Modify: `config/roles.schema.json`
  - Modify: `config/roles.default.json` (add at least one valid `claude` role mapping)
  - Add: `tests/opencode/test-claude-connector.sh`
- **How to implement**:
  1. Implement `connectors/claude.py` with argument normalization, error propagation, and structured logs (parity with existing connectors).
  2. Extend role schema tool enum to include `claude`.
  3. Add `claude)` dispatch branch in `bin/ask-role`.
  4. Add connector tests with a stubbed `claude` binary to validate call path and exit-code handling.
- **Acceptance criteria**:
  - [ ] Roles can declare `tool: "claude"` and pass schema validation.
  - [ ] `ask-role` dispatches to `connectors/claude.py` correctly.
  - [ ] Structured logs include correct role/tool metadata.
- **Verification commands**:
  - `bash tests/opencode/test-claude-connector.sh`
  - `python3 -m py_compile connectors/claude.py`

### T6-004 (P0) Wire MCP configuration into runtime behavior

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Architect + Coder
- **Why**: MCP is declared in config but not loaded or validated by any runtime path.
- **What to implement**:
  - Modify: `bin/multipowers`
  - Add: `scripts/validate_mcp.py`
  - Modify: `package.json`
  - Add: `tests/opencode/test-mcp-config.sh`
- **How to implement**:
  1. Define config priority for MCP (project config override + default fallback).
  2. Add `multipowers doctor` checks for MCP config existence/shape.
  3. Add validator script for `mcpServers` schema-level sanity.
  4. Ensure package artifact includes MCP defaults.
  5. Add tests for valid/invalid config and precedence.
- **Acceptance criteria**:
  - [ ] MCP config is validated in doctor flow.
  - [ ] Invalid MCP config fails with clear remediation guidance.
  - [ ] MCP defaults are included in distributable package contents.
- **Verification commands**:
  - `bash tests/opencode/test-mcp-config.sh`
  - `python3 scripts/validate_mcp.py --config config/mcp.default.json`

### T6-005 (P1) Implement major-change governance pipeline

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Architect + Coder
- **Why**: Required governance gates (`semgrep`, `biome`, `ruff`, docs sync) are not enforced by scripts, tests, or npm workflows.
- **What to implement**:
  - Add: `scripts/run_governance_checks.sh`
  - Add: `scripts/check_docs_sync.py`
  - Modify: `package.json`
  - Modify: `tests/opencode/run-tests.sh`
  - Add: `tests/opencode/test-governance-checks.sh`
- **How to implement**:
  1. Build changed-file driven checker pipeline (`git diff --name-only` input supported).
  2. Run `semgrep`, `biome`, `ruff` conditionally by file type.
  3. Add docs-sync policy checks (behavioral/code changes require matching docs update).
  4. Integrate governance script as an explicit npm script and test gate.
- **Acceptance criteria**:
  - [ ] Governance checks run locally via one command.
  - [ ] Check failures are actionable and non-silent.
  - [ ] Docs-sync violations are detectable.
- **Verification commands**:
  - `bash tests/opencode/test-governance-checks.sh`
  - `npm run governance:check`

### T6-006 (P1) Restore real plugin runtime artifact and remove silent fallback masking

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Coder
- **Why**: Current tests can pass with a generated fallback plugin stub, which masks missing runtime implementation.
- **What to implement**:
  - Add: `.opencode/plugins/superpowers.js`
  - Modify: `tests/opencode/setup.sh`
  - Modify: `tests/opencode/test-plugin-loading.sh`
- **How to implement**:
  1. Commit the real plugin runtime file into repository.
  2. Change setup to fail if plugin source is missing (or make fallback explicit-fail mode in CI).
  3. Extend plugin-loading test to assert runtime capabilities, not just file presence.
- **Acceptance criteria**:
  - [ ] Repository includes real plugin runtime source.
  - [ ] Tests fail when runtime source is missing.
  - [ ] Plugin loading test validates meaningful runtime behavior.
- **Verification commands**:
  - `bash tests/opencode/test-plugin-loading.sh`
  - `node --check .opencode/plugins/superpowers.js`

### T6-007 (P1) Standardize task-status/evidence automation for vibe coding

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Architect
- **Why**: Plan status/evidence automation is brittle across language/style variations and hinders reliable status updates by agents.
- **What to implement**:
  - Modify: `scripts/check_plan_evidence.py`
  - Add: `scripts/update_plan_task_status.py`
  - Modify: `tests/opencode/test-plan-evidence.sh`
  - Update: `docs/plans/status-evidence-template.md`
- **How to implement**:
  1. Expand parser support for both `Status` and `状态` labels.
  2. Add task status updater utility by Task ID.
  3. Validate required evidence fields for `DONE` tasks regardless of language label.
  4. Add regression tests for mixed-language status headers.
- **Acceptance criteria**:
  - [ ] Status updates can be done by tool using task ID + target status.
  - [ ] Evidence check correctly recognizes `DONE` tasks in supported formats.
  - [ ] Backward compatibility with existing gap plan documents is preserved.
- **Verification commands**:
  - `bash tests/opencode/test-plan-evidence.sh`
  - `python3 scripts/update_plan_task_status.py --help`

### T6-008 (P2) Add routing/workflow structured observability

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Router
- **Why**: Current logs do not explain lane/workflow/node decisions, reducing debuggability and governance traceability.
- **What to implement**:
  - Modify: `connectors/utils.py`
  - Modify: `bin/multipowers`
  - Add: `tests/opencode/test-routing-observability.sh`
  - Update: `README.md`
- **How to implement**:
  1. Add structured events for `lane_selected`, `workflow_started`, `workflow_node_executed`, `workflow_finished`.
  2. Carry `request_id`/`track_id` through lane + workflow events.
  3. Assert event presence and field completeness in tests.
- **Acceptance criteria**:
  - [ ] Routing decisions are queryable from JSONL logs.
  - [ ] Workflow node timeline is reconstructable per request.
  - [ ] Tests verify event schema and sequence.
- **Verification commands**:
  - `bash tests/opencode/test-routing-observability.sh`
  - `python3 -m py_compile connectors/utils.py`

---

## 4) Recommended Execution Sequence

1. **Foundation**: T6-001 → T6-002 → T6-003 → T6-004
2. **Governance + Runtime Integrity**: T6-005 → T6-006 → T6-007
3. **Observability Hardening**: T6-008

---

## 5) Definition of Done (Plan-Level)

1. Fast/Standard routing is executable, deterministic, and test-covered.
2. Standard lane workflows are configurable and enforce node-level role routing.
3. `claude`, `codex`, and `gemini` are all valid connector targets.
4. MCP config is validated and wired into runtime expectations.
5. Major-change governance checks run automatically and gate completion claims.
6. Plugin runtime is present in-repo and not masked by silent fallback stubs.
7. Task status updates and evidence checks are reliable for vibe-coding workflows.
8. Routing/workflow decisions are visible in structured logs.

---

## 6) Evidence Section (to fill when tasks are DONE)

- **Coverage Task IDs**: ``
- **Date**: ``
- **Verifier**: ``
- **Command(s)**:
  - ``
- **Exit Code**: ``
- **Key Output**:
  - ``
