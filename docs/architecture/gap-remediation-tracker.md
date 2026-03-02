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

*To be populated from Task 10 analysis.*

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
