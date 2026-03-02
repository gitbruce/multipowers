# 2026-03-02 Architecture Diff Docs Supplement Design

## Background

Current architecture diff docs were refreshed:
- `docs/architecture/commands_skills_difference.md`
- `docs/architecture/script-differences.md`
- `docs/architecture/other-differences.md`

But they still need consistency hardening and traceability supplementation to satisfy:
- `.multipowers/product-guidelines.md`
- `.multipowers/product.md`

## Design Goal

Make the three architecture diff docs consistent, auditable, and product-aligned without forcing mechanical one-to-one migration of all `main` files/features into `go`.

## Hard Constraints

1. `main` is read-only for delivery; implementation and mapping updates land on `go`.
2. Upstream baseline and `v8.31.1` migration traceability remain explicit.
3. Mapping language must be explicit: `source file -> target file -> target symbol`.
4. This project does **not** require migrating every `main` file/feature to `go`.
5. Migration decisions must satisfy product intent and invocation/hook contracts in `.multipowers/product.md`.

## Scope And Non-Goals

### In Scope

- Fix baseline/statistics/status terminology consistency across the three docs.
- Add unified evidence skeleton for `partial`/`missing` rows.
- Add explicit migration decision fields based on product constraints.

### Non-Goals

- No full implementation migration in this design task.
- No requirement for one-to-one mechanical translation for all legacy files.
- No modifications on `main`.

## Supplement Strategy (A2)

Use "consistency first + evidence skeleton":

1. Fix document consistency defects first (baseline SHA, status dictionary, terminology).
2. Add uniform evidence fields to all `partial`/`missing` mappings.
3. Separate migration decisions into:
   - `MIGRATE_TO_GO` (required by product behavior/contract)
   - `COPY_FROM_MAIN` (read-only continuity wrapper assets)
   - `EXCLUDE_WITH_REASON` (intentionally not migrated; explicit rationale)
   - `DEFER_WITH_CONDITION` (not now; trigger condition stated)

## Migration Decision Gate (Product-Aligned)

A `main` file/feature is `MIGRATE_TO_GO` only when at least one holds:

1. Needed to preserve upstream product intent (multi-provider orchestration, phase gates, consensus, safety).
2. Needed by invocation/call process or hook lifecycle (`SessionStart`, `PreToolUse`, `PostToolUse`, `Stop/SubagentStop`, etc.).
3. Needed to maintain deterministic runtime contract (`status/action/error_code/message/data/remediation`).
4. Needed for policy, security, or quality enforcement in runtime-critical paths.

Otherwise:
- `COPY_FROM_MAIN` for simple wrappers without business logic.
- `EXCLUDE_WITH_REASON` when outside product scope or intentionally retired.
- `DEFER_WITH_CONDITION` when migration depends on future runtime ownership.

## Per-Document Supplements

### 1) commands_skills_difference.md

Add minimal fields for all `partial/missing` rows:
- `runtime_target_symbol`
- `test_or_verification`
- `contract_coverage` (`yes|no|planned`)
- `decision` (`MIGRATE_TO_GO|EXCLUDE_WITH_REASON|DEFER_WITH_CONDITION`)

Priority rows:
- `extract-skill` wrong-route risk
- `octo -> mp` root intent routing parity
- `claw/doctor/schedule/scheduler/sentinel` missing commands/skills

### 2) script-differences.md

Keep domain-full list; add normalized columns:
- `implemented_symbol`
- `planned_symbol`
- `evidence_level` (`E0|E1|E2|E3`)
- `decision_reason`

Mandatory note for all `COPY_FROM_MAIN` rows:
- `sync_policy=read-only-from-main`

Add Hook Lifecycle index section:
- map legacy hook scripts to lifecycle events and Go ownership symbols.

### 3) other-differences.md

Fix baseline `go` SHA to align with the other two docs.

For `partial/missing` rows add:
- `target_symbol_or_contract`
- `decision`
- `decision_reason`
- `test_or_verification`

For `mcp-server/*` and `openclaw/*`:
- no unresolved "floating missing"
- each row must be `MIGRATE_TO_GO`, `EXCLUDE_WITH_REASON`, or `DEFER_WITH_CONDITION`.

## Unified Evidence Legend

Add same legend to all three docs:
- `E0`: doc-only plan
- `E1`: symbol exists
- `E2`: test exists
- `E3`: verified output recorded

Rule:
- every `partial/missing` row must have at least `E0`.

## Acceptance Criteria

1. Three docs use the same baseline pair (`main`/`go`) and same status dictionary.
2. No ambiguous mapping row without `source -> target -> symbol/contract`.
3. No unresolved `missing` for major domains without explicit decision/reason.
4. `COPY_FROM_MAIN` rows include read-only sync statement.
5. Product alignment is explicit: migration decisions are driven by product constraints, not by mechanical one-to-one parity.

## Deliverables

1. Updated:
   - `docs/architecture/commands_skills_difference.md`
   - `docs/architecture/script-differences.md`
   - `docs/architecture/other-differences.md`
2. Verification transcript snippets for baseline and consistency checks.
3. Follow-up implementation plan generated via `writing-plans` workflow.
