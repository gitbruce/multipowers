# Architecture Diff Gap Remediation Tracker

> **Purpose:** Single source-of-truth for unresolved gaps across architecture diff documents.
> **Last Updated:** 2026-03-02
> **Go Baseline:** See `go=` hash in individual diff docs for current baseline reference.

| gap_id | source_doc | source_anchor | gap_type | current_state | target_state | decision | evidence_level | owner | next_action |
|--------|------------|---------------|----------|---------------|--------------|----------|----------------|-------|-------------|

## Commands/Skills High-Risk

High-risk command and skill gaps requiring explicit closure paths.

| gap_id | item | risk_reason | target_symbol/contract | test_reference | closure_condition | status |
|--------|------|-------------|------------------------|----------------|-------------------|--------|
| CMD-001 | extract-skill | Core workflow misroutes to `mp status` instead of extract | `internal/cli/extract.go:ExtractSkill` | `internal/cli/extract_test.go` | `mp extract` command exists with test coverage ≥80% | pending |
| CMD-002 | octo -> mp | Root intent routing logic reduced in go | `internal/cli/root.go:RouteIntent` | `internal/cli/root_test.go` | `mp route --intent` returns valid routing for all registered intents | pending |
| CMD-003 | claw | External system integration not in current scope | `internal/external/claw/adapter.go` (planned) | `internal/external/claw/adapter_test.go` (planned) | Product requirement explicitly requests claw integration | deferred |
| CMD-004 | doctor | Diagnostic capability replaced by sys-configure | N/A (replaced) | `internal/cli/sys_configure_test.go` | N/A - excluded with reason | closed |
| CMD-005 | schedule/scheduler | Scheduler domain contract undefined | `internal/scheduler/scheduler.go` (planned) | `internal/scheduler/scheduler_test.go` (planned) | Scheduler domain contract defined in `.multipowers/product.md` | deferred |
| CMD-006 | sentinel | Security gate capability required | `internal/governance/sentinel.go` (planned) | `internal/governance/sentinel_test.go` (planned) | Sentinel gate blocks invalid states with test coverage | pending |
| CMD-007 | skill-claw/skill-doctor | Product scope excludes these skills | N/A | N/A | N/A - excluded with reason | closed |
| CMD-008 | parallel command | Command wrapper missing for flow-parallel skill | `internal/cli/parallel.go` | `internal/cli/parallel_test.go` | `/mp:parallel` command invokes flow-parallel skill | pending |
| CMD-009 | spec command | Command wrapper missing for flow-spec skill | `internal/cli/spec.go` | `internal/cli/spec_test.go` | `/mp:spec` command invokes flow-spec skill | pending |

**Source:** `docs/architecture/commands_skills_difference.md` § 决策与证据索引（高风险项）

## Script Missing Decision Classification

Unresolved script rows grouped by domain/pattern.

| gap_id | pattern | default_decision | decision_reason | closure_path | status |
|--------|---------|------------------|-----------------|--------------|--------|

*To be populated from Task 11 analysis.*

## Other-Differences Partial/Missing Contracts

High-impact configuration and documentation gaps.

| gap_id | item | target_symbol_or_contract | evidence_upgrade_path | owner_domain | status |
|--------|------|---------------------------|----------------------|--------------|--------|

*To be populated from Task 13 analysis.*

---

## Legend

**Decision Tokens:**
- `MIGRATE_TO_GO` - Implement in Go atomic commands
- `COPY_FROM_MAIN` - Copy directly from main branch
- `EXCLUDE_WITH_REASON` - Not needed, documented rationale
- `DEFER_WITH_CONDITION` - Postpone until trigger condition met

**Evidence Levels:**
- `E0` - Documentation only, no implementation
- `E1` - Target symbol/contract exists
- `E2` - Test coverage exists
- `E3` - Verified output matches expected

**Gap Types:**
- `missing_command` - Command exists in main, not in go
- `missing_skill` - Skill exists in main, not in go
- `implementation_diff` - Different implementation approach
- `contract_gap` - Response contract mismatch
- `test_gap` - Missing test coverage

---

## Verification

Run: `scripts/verify-architecture-diff-docs.sh`
