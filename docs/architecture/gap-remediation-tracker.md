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

| gap_id | pattern | count | default_decision | decision_reason | closure_path | status |
|--------|---------|-------|------------------|-----------------|--------------|--------|
| SCR-001 | `scripts/scheduler/*.sh` | 6 | `DEFER_WITH_CONDITION` | Scheduler domain contract undefined | `internal/scheduler/*_test.go` | deferred |
| SCR-002 | `scripts/extract/*.sh` | 1 | `MIGRATE_TO_GO` | Core extraction workflow | `internal/extract/core_test.go` | pending |
| SCR-003 | `tests/smoke/*.sh` | 7 | `MIGRATE_TO_GO` | CLI surface validation | `internal/cli/smoke_test.go` | pending |
| SCR-004 | `tests/live/*.sh` | 3 | `DEFER_WITH_CONDITION` | External service dependency | `internal/workflows/live_test.go` | deferred |
| SCR-005 | `tests/benchmark/*.sh` | 2 | `MIGRATE_TO_GO` | Performance regression guard | `internal/workflows/benchmark_test.go` | pending |
| SCR-006 | `tests/integration/*.sh` | 6 | `MIGRATE_TO_GO` | Plugin lifecycle tests | `internal/workflows/integration_test.go` | pending |
| SCR-007 | `tests/helpers/*.sh` | 4 | `MIGRATE_TO_GO` | Test infrastructure | `internal/devx/helpers_test.go` | pending |
| SCR-008 | `scripts/metrics-tracker.sh` | 1 | `MIGRATE_TO_GO` | Cost observability | `internal/metrics/tracker_test.go` | pending |
| SCR-009 | `scripts/permissions-manager.sh` | 1 | `MIGRATE_TO_GO` | Consent governance | `internal/permissions/manager_test.go` | pending |
| SCR-010 | `scripts/agent-teams-bridge.sh` | 1 | `MIGRATE_TO_GO` | Team coordination | `internal/teams/bridge_test.go` | pending |
| SCR-011 | `scripts/async-tmux-features.sh` | 1 | `MIGRATE_TO_GO` | Async execution | `internal/workflows/async_test.go` | pending |
| SCR-012 | `tests/test-*.sh` (regression) | 21 | `MIGRATE_TO_GO` | Feature parity tests | `internal/*/regression_test.go` | pending |
| SCR-013 | `tests/unit/test-cron-parser.sh` | 1 | `MIGRATE_TO_GO` | Scheduler dependency | `internal/scheduler/cron_test.go` | pending |

**Summary:** 48 `MIGRATE_TO_GO` + 10 `DEFER_WITH_CONDITION` = 58 total missing script gaps

**Source:** `docs/architecture/script-differences.md` § Missing Decision Classification Matrix

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
