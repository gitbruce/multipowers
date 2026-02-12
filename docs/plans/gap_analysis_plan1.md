# Multipowers Gap Analysis Plan 1

> **Goal:** Based on `conductor/context` (tool project background), identify remaining gaps between project intention and current executable codebase, then provide an implementation plan with vibe-coding updatable task statuses.

## 1) Analysis Baseline

### 1.1 Intention Sources (Tool Project Background)
- `conductor/context/product.md` - Core boundaries and objectives
- `conductor/context/product-vision.md` - Strategic positioning and design principles
- `conductor/context/product-guidelines.md` - Routing, workflows, and roles
- `conductor/context/workflow.md` - Delivery and execution routing
- `conductor/context/tech-stack.md` - Languages and core components

### 1.2 Codebase Scope (for gap analysis)
Analyzed implementation scope excludes:
- `/conductor` - Template/context definitions (stable)
- `/docs` - Documentation outputs
- `/templates` - User project bootstrap templates
- `/outputs` - Generated artifacts

Included key implementation paths:
- `bin/*` - CLI entrypoints (`multipowers`, `ask-role`)
- `scripts/*` - Python validation and execution scripts
- `connectors/*` - Model CLI wrappers (codex, gemini, claude)
- `config/*` - Default configuration schemas
- `lib/*` - Core JavaScript utilities (skills-core.js)
- `skills/*` - Superpowers workflow skill definitions
- `tests/*` - Integration and unit test suites

---

## 2) Intention-to-Reality Gap Summary

| Gap ID | Intended Capability | Current Code Evidence | Risk | Priority |
|---|---|---|---|---|
| GA1-01 | Router should provide one-command execution (route + execute) | `multipowers run` already exists; documentation gap remains | Users may not know about `run` command | P2 |
| GA1-02 | Skills should be auto-discovered and loadable | `lib/skills-core.js` implements skill discovery, but no CLI integration for skill listing/loading | Skill system exists but disconnected from main CLI | P0 |
| GA1-03 | Role dispatch should validate role existence before external CLI call | `ask-role` validates roles via schema but only after parsing config | Late validation, potential unclear errors | P1 |
| GA1-04 | Context files should be mandatory with strict mode enforcement | `ask-role` implements strict/lenient modes, but no global enforcement mechanism | Context quality depends on manual discipline | P1 |
| GA1-05 | Workflow execution should support checkpoint/resume | `execute_workflow.py` runs nodes sequentially but no checkpoint/state persistence | Long-running workflows can't recover from failure | P1 |
| GA1-06 | Governance checks should integrate with workflow execution | `run_governance_checks.sh` exists but not integrated into workflow nodes | Governance requires manual invocation | P1 |
| GA1-07 | Track lifecycle should enforce governance gates | `track complete` has governance hooks but bypassable via `--skip-governance` | "Done" may bypass policy without trace | P1 |
| GA1-08 | Plan evidence validation should support custom rulesets | `check_plan_evidence.py` has hardcoded field requirements | Not extensible for project-specific evidence | P1 |
| GA1-09 | Connector retry logic should be configurable | Connectors have basic error propagation but no retry/backoff config | Transient failures may abort valid work | P1 |
| GA1-10 | Skills should have role-contract enforcement | Skills define role contracts (e.g., subagent-driven-development) but no validation | Skills may be used outside intended role context | P1 |
| GA1-13 | Test coverage should meet governance thresholds | Tests exist but no coverage reporting or enforcement | Quality signal is weak | P0 |
| GA1-14 | Template sync should be automated from conductor | No automation for syncing conductor/ improvements to templates/ | `workflow.md` §5 requirement unaddressed | P0 |
| GA1-15 | Evidence validation should check content, not just structure | `check_plan_evidence.py` checks field existence only | Evidence may pass without proving completion | P1 |

---

## 3) Task Board (Vibe-Coding Status Updatable)

Status values: `TODO` / `IN_PROGRESS` / `BLOCKED` / `DONE`

Update task status:
```bash
python3 scripts/update_plan_task_status.py --file docs/plans/gap_analysis_plan1.md --task-id T1-001 --status IN_PROGRESS
```

| Task ID | Status | Priority | Owner | Depends On |
|---|---|---|---|---|
| T1-001 | `TODO` | P0 | Integrator | - |
| T1-002 | `TODO` | P0 | Integrator | T1-001 |
| T1-013 | `TODO` | P0 | Integrator | - |
| T1-011 | `TODO` | P0 | Integrator | - |
| T1-014 | `TODO` | P0 | Integrator | - |
| T1-003 | `TODO` | P1 | Integrator | T1-001 |
| T1-004 | `TODO` | P1 | Integrator | T1-002 |
| T1-007 | `TODO` | P1 | Integrator | T1-004, T1-011 |
| T1-006 | `TODO` | P1 | Integrator | T1-004 |
| T1-015 | `TODO` | P1 | Integrator | T1-013 |
| T1-008 | `TODO` | P1 | Integrator | - |
| T1-009 | `TODO` | P1 | Integrator | T1-003 |
| T1-005 | `TODO` | P2 | Integrator | T1-003 |

---

### T1-001 (P0) Implement CLI skill system integration

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - Intention: Skills should be the primary method for workflow execution
  - Reality: `lib/skills-core.js` implements discovery but no CLI commands exist to list/load skills
  - Gap: Users cannot discover available skills from CLI
- **What to implement**:
  - Add: `bin/multipowers skill` subcommand with `list`, `info`, `validate` actions
  - Modify: `lib/skills-core.js` to support CLI invocation
  - Add: `tests/opencode/test-skill-cli-integration.sh`
- **How to implement**:
  1. Add `skill_command()` function to `bin/multipowers` handling:
     - `multipowers skill list` - Enumerate all skills from `skills/` and `.skills/`
     - `multipowers skill info <skill-name>` - Show skill description and role contract
     - `multipowers skill validate` - Validate all skill frontmatter
  2. Reuse `findSkillsInDir()` from `lib/skills-core.js`
  3. Add test coverage for skill discovery edge cases (missing frontmatter, duplicate names)
- **Acceptance criteria**:
  - [ ] `multipowers skill list` shows all 13 superpowers skills
  - [ ] `multipowers skill info subagent-driven-development` shows role contract
  - [ ] `multipowers skill validate` passes for valid skills, fails for invalid frontmatter
- **Verification commands**:
  - `bash tests/opencode/test-skill-cli-integration.sh`
  - `./bin/multipowers skill list | grep brainstorming`

---

### T1-002 (P0) Add role-contract validation to skill invocation

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - Skills document role contracts (e.g., "Main Role: router", "Step Roles: coder, architect")
  - No validation enforces these contracts during execution
  - Risk: Skills may be invoked with wrong model/role combination
- **What to implement**:
  - Modify: `bin/ask-role` to accept `--expected-role` parameter
  - Add: `scripts/validate_skill_role_contract.py`
  - Add: `tests/opencode/test-skill-role-contract.sh`
- **How to implement**:
  1. Extend skill frontmatter schema to include required role contract fields:
     ```yaml
     roles:
       main: router
       steps:
         - role: coder
           stage: implementation
         - role: architect
           stage: review
     ```
  2. In `ask-role`, add `--validate-role-contract <skill-name>` option
  3. Validation script parses skill frontmatter, extracts role requirements, validates against `roles.json`
  4. Return error if requested role doesn't match skill contract
- **Acceptance criteria**:
  - [ ] Skill with main:router validates when invoked with router role
  - [ ] Skill invocation with mismatched role fails with clear error
  - [ ] Validation is optional (via flag) for backward compatibility
- **Verification commands**:
  - `bash tests/opencode/test-skill-role-contract.sh`
  - `./bin/ask-role coder "test" --validate-role-contract subagent-driven-development; echo $?` # should fail

---

### T1-013 (P0) Add test coverage enforcement

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - Tests exist but no coverage measurement
  - Quality signal is weak without coverage data
  - Intention: "Evidence before claims" - coverage is foundational
  - **Moved to Phase 0**: Without coverage measurement, other evidence is hollow
- **What to implement**:
  - Modify: `scripts/run_governance_checks.sh` - Add coverage check
  - Add: `config/coverage-rules.json` for thresholds
  - Add: `tests/opencode/test-coverage-enforcement.sh`
- **How to implement**:
  1. Add coverage tools to governance:
     - Python: `coverage run -m pytest` + `coverage report`
     - Bash: `kcov` or `bashcov` (if available)
  2. Define coverage thresholds by file type:
     ```json
     {
       "thresholds": {
         "python": 80,
         "bash": 60
       }
     }
     ```
  3. In governance script, run coverage if test files changed
  4. Fail if coverage below threshold
  5. Generate coverage report to `outputs/coverage/`
- **Acceptance criteria**:
  - [ ] Coverage measured for Python test files
  - [ ] Governance fails if coverage below threshold
  - [ ] Coverage reports generated
  - [ ] Thresholds configurable by file type
- **Verification commands**:
  - `bash tests/opencode/test-coverage-enforcement.sh`

---

### T1-011 (P0) Add context quality enforcement mechanism

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - `check_context_quality.py` exists but only warns
  - No enforcement mechanism prevents low-quality context
  - Intention: "Keep stable project background in context docs"
  - **Moved to Phase 0**: Bad context inputs make all other work wasteful
- **What to implement**:
  - Modify: `scripts/check_context_quality.py` - Add `--enforce` flag
  - Modify: `bin/multipowers doctor` - Add context quality blocking
  - Add: `config/context-quality-rules.json` for quality thresholds
- **How to implement**:
  1. Define quality rules:
     ```json
     {
       "rules": {
         "min_content_length": 200,
         "forbidden_patterns": ["TODO", "FIXME", "PLACEHOLDER"],
         "required_sections": ["## Overview", "## Usage"]
       }
     }
     ```
  2. In `check_context_quality.py`, add `--enforce` mode:
     - Fail if files don't meet quality thresholds
     - Report specific violations
  3. In `doctor` command, call with `--enforce` if `MULTIPOWERS_CONTEXT_MODE=strict`
- **Acceptance criteria**:
  - [ ] Context files with TODO/PLACEHOLDER fail in enforce mode
  - [ ] Minimum content length enforced
  - [ ] Required sections validated
  - [ ] Doctor fails context checks in strict mode
- **Verification commands**:
  - `python3 scripts/check_context_quality.py --enforce --context-dir conductor/context`

---

### T1-014 (P0) Implement template sync automation

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - `workflow.md` §5 explicitly states: "evaluate and sync safe parts into templates/conductor/"
  - Current: No automation for template syncing
  - Gap: conductor/ and templates/ will diverge over time
  - Impact: User projects bootstrap with outdated patterns
- **What to implement**:
  - Add: `scripts/sync_templates.py` for automated sync
  - Add: `config/template-sync-rules.json` for sync mapping
  - Add: `tests/opencode/test-template-sync-automation.sh`
- **How to implement**:
  1. Define sync mapping:
     ```json
     {
       "sync_map": {
         "conductor/config/roles.json": "templates/conductor/config/roles.json",
         "conductor/context/product.md": "templates/conductor/context/product.md",
         "conductor/context/workflow.md": "templates/conductor/context/workflow.md"
       }
     }
     ```
  2. Implement sync script:
     - Compare files between conductor/ and templates/
     - Report differences (missing, outdated, conflicts)
     - Support `--dry-run`, `--apply`, `--interactive` modes
     - Create backup before applying changes
  3. Add `multipowers template-sync` subcommand:
     - `multipowers template-sync --check` - Show drift
     - `multipowers template-sync --apply` - Sync changes
  4. Integrate with `doctor` to warn about drift
- **Acceptance criteria**:
  - [ ] Script reports differences between conductor/ and templates/
  - [ ] `--apply` safely updates template files
  - [ ] Doctor warns when template drift exceeds threshold
  - [ ] Backup created before modifications
- **Verification commands**:
  - `bash tests/opencode/test-template-sync-automation.sh`
  - `python3 scripts/sync_templates.py --check`

---

### T1-003 (P1) Implement workflow checkpoint/resume mechanism

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - Long-running workflows (multi-node) can fail mid-execution
  - Current `execute_workflow.py` has no state persistence
  - Users must re-run entire workflow after failure
- **What to implement**:
  - Modify: `scripts/execute_workflow.py`
  - Add: `.workflows/state/` directory for checkpoint files
  - Add: `tests/opencode/test-workflow-checkpoint-resume.sh`
- **How to implement**:
  1. After each successful node execution, write checkpoint:
     - Path: `.workflows/state/<request-id>.json`
     - Content: completed nodes, current index, timestamp
  2. Add `--resume <request-id>` flag to workflow execution
  3. On resume, skip completed nodes, start from last index
  4. Add `--force` flag to ignore existing checkpoints
  5. Clean up checkpoints older than 7 days (configurable)
- **Acceptance criteria**:
  - [ ] Workflow writes checkpoint after each successful node
  - [ ] `--resume` flag continues from last checkpoint
  - [ ] Resume skips completed nodes, logs which nodes were skipped
  - [ ] Checkpoint cleanup removes stale state files
- **Verification commands**:
  - `bash tests/opencode/test-workflow-checkpoint-resume.sh`
  - `ls -la .workflows/state/`

---

### T1-004 (P1) Integrate governance checks into workflow nodes

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - `run_governance_checks.sh` exists but requires manual invocation
  - Intention: "Major-change governance should be enforced, not optional"
  - Gap: No automatic governance execution during workflow
- **What to implement**:
  - Modify: `config/workflows.default.json` - Add governance node to relevant workflows
  - Modify: `scripts/execute_workflow.py` - Support special "governance" node type
  - Add: `tests/opencode/test-workflow-governance-integration.sh`
- **How to implement**:
  1. Define new workflow node type:
     ```json
     {
       "id": "governance",
       "type": "governance",
       "mode": "strict"
     }
     ```
  2. In `execute_workflow.py`, detect governance node type
  3. Invoke `run_governance_checks.sh` with current `git diff --name-only`
  4. Propagate exit code: fail workflow if governance fails
  5. Add optional "governance" node to `subagent-driven-development` workflow (after "implement", before "review")
- **Acceptance criteria**:
  - [ ] Governance node executes checks on changed files
  - [ ] Workflow fails if governance checks fail
  - [ ] Governance artifacts written to standard location
  - [ ] Existing workflows without governance nodes unchanged
- **Verification commands**:
  - `bash tests/opencode/test-workflow-governance-integration.sh`

---

### T1-007 (P1) Enforce governance metadata on track completion

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - Plan8 already implemented governance gating, but metadata writing needs verification
  - Need to ensure "Done" means actually done with governance
  - Gap: Governance may pass but track metadata not updated properly
- **What to implement**:
  - Modify: `bin/multipowers` track complete
  - Verify: Governance metadata written to track file
  - Add: `tests/opencode/test-track-governance-metadata.sh`
- **How to implement**:
  1. After governance checks pass, update track metadata:
     - Add `**Governance:** passed` with timestamp
     - Add `**Governance Artifact:** <path>` for audit trail
  2. Read track file before and after completion to verify metadata exists
  3. Fail completion if metadata update fails
  4. Ensure governance artifact path is absolute and readable
- **Acceptance criteria**:
  - [ ] Track completion writes governance status to track file
  - [ ] Governance artifact path is recorded in metadata
  - [ ] Completion fails if metadata write fails
  - [ ] Track status command shows governance information
- **Verification commands**:
  - `bash tests/opencode/test-track-governance-metadata.sh`
  - `grep "Governance:" conductor/tracks/track-*.md`

---

### T1-006 (P1) Add custom evidence ruleset support

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - `check_plan_evidence.py` has hardcoded required fields
  - Different projects may need different evidence requirements
  - Current implementation not extensible
- **What to implement**:
  - Modify: `scripts/check_plan_evidence.py`
  - Add: `config/evidence-rules.json` for custom requirements
  - Add: `tests/opencode/test-custom-evidence-rules.sh`
- **How to implement**:
  1. Define evidence rules schema:
     ```json
     {
       "rulesets": {
         "default": {
           "required_fields": ["Date", "Verifier", "Command(s)", "Exit Code", "Key Output"],
           "optional_fields": ["Coverage Task IDs"]
         },
         "strict": {
           "required_fields": ["Date", "Verifier", "Command(s)", "Exit Code", "Key Output", "Coverage Task IDs", "Governance Artifact"],
           "require_governance_evidence": true
         }
       }
     }
     ```
  2. Add `--ruleset <name>` flag to `check_plan_evidence.py`
  3. Load rules from config if provided, else use default
  4. Validate evidence sections against selected ruleset
- **Acceptance criteria**:
  - [ ] Default ruleset matches current hardcoded behavior
  - [ ] Custom ruleset can be specified via flag
  - [ ] Strict ruleset enforces governance evidence
  - [ ] Invalid ruleset name returns helpful error
- **Verification commands**:
  - `bash tests/opencode/test-custom-evidence-rules.sh`
  - `python3 scripts/check_plan_evidence.py --ruleset strict docs/plans/gap_analysis_plan1.md`

---

### T1-015 (P1) Add evidence content validation

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - `check_plan_evidence.py` checks field existence only
  - Evidence may pass without proving completion (e.g., "Exit Code: 0" doesn't prove right tests ran)
  - Gap: Structure validation ≠ content validation
- **What to implement**:
  - Modify: `scripts/check_plan_evidence.py`
  - Add: `config/evidence-content-rules.json` for content validation
  - Add: `tests/opencode/test-evidence-content-validation.sh`
- **How to implement**:
  1. Define content validation rules:
     ```json
     {
       "content_rules": {
         "command_patterns": {
           "required": ["pytest", "npm test", "coverage"],
           "forbidden": ["echo", "cat", "ls"]
         },
         "exit_code_meaning": {
           "0": "must have test_output",
           "nonzero": "must have explanation"
         }
       }
     }
     ```
  2. Add `--validate-content` flag to evidence checker
  3. Validate that commands actually test/verify the work
  4. Validate that exit codes have meaningful context
  5. Check that key outputs reference actual test results
- **Acceptance criteria**:
  - [ ] Commands validated to be test/verification commands
  - [ ] Exit codes have explanatory context
  - [ ] Key outputs reference actual results
  - [ ] Content-only evidence fails validation
- **Verification commands**:
  - `bash tests/opencode/test-evidence-content-validation.sh`
  - `python3 scripts/check_plan_evidence.py --validate-content docs/plans/gap_analysis_plan1.md`

---

### T1-008 (P1) Implement connector retry/backoff configuration

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - External CLI calls (codex, gemini) may fail transiently
  - Current connectors have no retry logic
  - Transient failures abort valid work
- **What to implement**:
  - Modify: `connectors/codex.py`, `connectors/gemini.py`, `connectors/claude.py`
  - Add: `config/connector-retry.json` for retry policies
  - Add: `tests/opencode/test-connector-retry-logic.sh`
- **How to implement**:
  1. Define retry policy schema:
     ```json
     {
       "default": {
         "max_retries": 3,
         "backoff": "exponential",
         "initial_delay_ms": 1000,
         "max_delay_ms": 10000,
         "retryable_codes": [503, 502, 429]
       }
     }
     ```
  2. Add retry wrapper in each connector:
     - Catch subprocess.CalledProcessError
     - Check if exit code/signals retryable
     - Sleep with backoff
     - Retry up to max
  3. Log retry attempts to stderr for observability
- **Acceptance criteria**:
  - [ ] Connector retries on transient failures (configurable)
  - [ ] Exponential backoff implemented between retries
  - [ ] Non-retryable errors fail immediately
  - [ ] Retry attempts logged to stderr
- **Verification commands**:
  - `bash tests/opencode/test-connector-retry-logic.sh`

---

### T1-009 (P1) Add per-role context budget configuration

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - Research roles (architect, librarian) need more context than executor roles
  - Current `MULTIPOWERS_CONTEXT_BUDGET` is global
  - One-size-fits-all doesn't match role needs
- **What to implement**:
  - Modify: `config/roles.default.json` - Add `context_budget` field per role
  - Modify: `bin/ask-role` - Read role-specific budget from config
  - Add: `tests/opencode/test-role-context-budgets.sh`
- **How to implement**:
  1. Extend role schema:
     ```json
     {
       "roles": {
         "architect": {
           "context_budget": 256000,
           "tool": "gemini",
           ...
         },
         "coder": {
           "context_budget": 128000,
           "tool": "codex",
           ...
         }
       }
     }
     ```
  2. In `ask-role`, after resolving role, look up `context_budget` field
  3. Use role-specific budget, fallback to global `MULTIPOWERS_CONTEXT_BUDGET`
  4. Log effective budget to stderr
- **Acceptance criteria**:
  - [ ] Architect role gets 256k tokens budget (configurable)
  - [ ] Coder role gets 128k tokens budget (configurable)
  - [ ] Global env var overrides all role budgets
  - [ ] Budget logged during context injection
- **Verification commands**:
  - `bash tests/opencode/test-role-context-budgets.sh`
  - `MULTIPOWERS_CONTEXT_BUDGET=64000 ./bin/ask-role architect "test" 2>&1 | grep "Context size"`

---

### T1-005 (P2) Strengthen track completion governance bypass logging

- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator
- **Why**:
  - Current `track complete` has `--skip-governance` bypass
  - Intention: "Completion should require governance artifacts"
  - Gap: Bypass is too easy, not audited well
- **What to implement**:
  - Modify: `bin/multipowers` track complete command
  - Add: `conductor/tracks/.bypass_log` for audit trail
  - Add: `tests/opencode/test-track-governance-bypass-audit.sh`
- **How to implement**:
  1. When `--skip-governance` is used:
     - Write log entry with timestamp, user, track_id, reason
     - Emit warning: "Governance bypass logged to .bypass_log"
  2. Add env var `MULTIPOWERS_ENFORCE_GOVERNANCE=true` to block bypass entirely
  3. Default: bypass allowed but logged
  4. If enforcement env var set: exit with error and remediation steps
- **Acceptance criteria**:
  - [ ] Governance bypass is logged to `.bypass_log`
  - [ ] Log entry includes timestamp, track_id, requesting user
  - [ ] `MULTIPOWERS_ENFORCE_GOVERNANCE=true` blocks bypass
  - [ ] Clear error message when enforcement blocks bypass
- **Verification commands**:
  - `bash tests/opencode/test-track-governance-bypass-audit.sh`
  - `cat conductor/tracks/.bypass_log`

---

## 4) Recommended Execution Sequence

### Phase 0 — Foundation
1. **T1-013** - Test coverage enforcement (evidence basis)
2. **T1-011** - Context quality enforcement (input quality)
3. **T1-014** - Template sync automation (workflow.md §5 requirement)

### Phase A — Skill System Foundation
4. T1-001 - CLI skill integration (enables skill discovery)
5. T1-002 - Role-contract validation (ensures correct skill usage)

### Phase B — Workflow Resilience
6. T1-003 - Workflow checkpoint/resume
7. T1-004 - Governance integration into workflows

### Phase C — Evidence & Quality
8. T1-007 - Governance metadata on completion
9. T1-006 - Custom evidence rulesets
10. T1-015 - Evidence content validation

### Phase D — Operations & Reliability
11. T1-008 - Connector retry logic
12. T1-009 - Per-role context budgets
13. T1-005 - Bypass logging (nice-to-have)

---

## 5) Definition of Done (Plan-Level)

1. Test coverage is measured and enforced
2. Context quality is enforceable with blocking mode
3. Template sync automation exists (addresses workflow.md §5)
4. Skills are discoverable and validated via CLI
5. Skill role contracts are enforced during invocation
6. Workflows support checkpoint/resume for resilience
7. Governance is integrated into workflow execution
8. Governance metadata is verified on track completion
9. Evidence validation supports custom rulesets
10. Evidence content validation ensures substantive proof
11. Connectors retry on transient failures
12. Context budgets are configurable per role
13. Governance bypass is fully audited

---

## 6) Evidence Section (fill when tasks are DONE)

- **Coverage Task IDs**: `T1-001, T1-002, T1-003, T1-004, T1-005, T1-006, T1-007, T1-008, T1-009, T1-011, T1-013, T1-014, T1-015`
- **Date**: `[FILL WHEN COMPLETE]`
- **Verifier**: `[FILL WHEN COMPLETE]`
- **Command(s)**:
  - `bash tests/opencode/run-tests.sh`
  - `python3 scripts/check_plan_evidence.py --require-governance-evidence docs/plans/gap_analysis_plan1.md`
  - `bash scripts/run_governance_checks.sh --mode strict`
- **Exit Code**: `[FILL WHEN COMPLETE]`
- **Key Output**:
  - `[FILL WHEN COMPLETE]`
