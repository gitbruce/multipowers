# Architecture Diff Gap Remediation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Close the remaining documentation gaps after the latest architecture-diff updates, so migration decisions are product-gated, evidence-graded, and operationally verifiable for junior contributors.

**Architecture:** Keep the three diff docs as source-of-truth and add a thin governance layer: (1) explicit gap tracker, (2) per-domain decision/evidence completion, (3) reproducible verification script. Do not force one-to-one mechanical migration from `main`; drive decisions from `.multipowers/product-guidelines.md` and `.multipowers/product.md`.

**Tech Stack:** Markdown, Bash, `rg`, `awk`, `sed`, Git.

---

## Scope, Why, and Success Criteria

**Why this plan exists**
- Current docs already introduced `decision` + `E0-E3`, but many entries are still at `E0` and not explicitly mapped to concrete symbol/test evidence.
- `script-differences.md` changed P0 to "classify decision first", but missing rows are not yet systematically grouped for execution.
- `other-differences.md` fixed baseline and added `mcp-server/openclaw` decisions, but remaining `partial/missing` items still need a structured gap-closure workflow.

**What will be delivered**
- A single gap tracker for unresolved rows.
- Three docs updated with explicit "what still missing + how to close" guidance.
- A deterministic verification script to prevent regression.

**How success is measured**
- Every unresolved gap is tracked with owner, decision, evidence level, and next action.
- High-risk command/skill rows include concrete target symbol + test reference (or explicit defer/exclude condition).
- Script and other-diff gaps are classified by domain, not left as free-form TODO text.
- One command verifies consistency across all three docs.

**Key Design**
- Product constraints override parity pressure: not every `main` file must migrate.
- Keep edits incremental and reviewable; avoid table-wide mechanical rewrites when domain-index can express intent clearly.
- Prefer "decision matrix + verification automation" over ad-hoc text additions.

---

### Task 1: Create A Single Gap Tracker

**Why:** Junior contributors need one place to see what is unresolved; currently gaps are spread across three docs.

**What:** Add a tracker doc that lists remaining gaps and links back to source rows.

**How:** Create a dedicated tracker with strict columns and initial placeholders.

**Key Design:** Tracker is execution-oriented, not narrative: each row must be actionable.

**Files:**
- Create: `docs/architecture/gap-remediation-tracker.md`
- Test: `docs/architecture/gap-remediation-tracker.md`

**Step 1: Write the failing test**

```bash
test -f docs/architecture/gap-remediation-tracker.md
```

**Step 2: Run test to verify it fails**

Run:
```bash
test -f docs/architecture/gap-remediation-tracker.md
```

Expected: FAIL (`No such file or directory`).

**Step 3: Write minimal implementation**

Create file with this template:

```md
# Architecture Diff Gap Remediation Tracker

| gap_id | source_doc | source_anchor | gap_type | current_state | target_state | decision | evidence_level | owner | next_action |
|---|---|---|---|---|---|---|---|---|---|

## Commands/Skills High-Risk
## Script Missing Decision Classification
## Other-Differences Partial/Missing Contracts
```

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
test -f docs/architecture/gap-remediation-tracker.md
rg -q "gap_id" docs/architecture/gap-remediation-tracker.md
rg -q "Commands/Skills High-Risk" docs/architecture/gap-remediation-tracker.md'
```

Expected: PASS.

**Step 5: Commit**

```bash
git add docs/architecture/gap-remediation-tracker.md
git commit -m "docs(architecture): add gap remediation tracker"
```

### Task 2: Close `commands_skills_difference` High-Risk E0 Gaps

**Why:** High-risk rows (`extract-skill`, `octo->mp`, missing command/skill entries) still need concrete closure paths for implementation and review.

**What:** Upgrade high-risk index from "decision only" to "decision + target symbol + test reference + closure condition".

**How:** Expand high-risk table rows in `commands_skills_difference.md` and add tracker entries.

**Key Design:** Keep `DEFER_WITH_CONDITION` explicit (trigger condition must be concrete and testable).

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`
- Modify: `docs/architecture/gap-remediation-tracker.md`
- Test: `docs/architecture/commands_skills_difference.md`

**Step 1: Write the failing test**

```bash
bash -lc 'set -euo pipefail
rg -n "决策与证据索引（高风险项）" docs/architecture/commands_skills_difference.md
rg -n "internal/cli/root.go|internal/providers/router_intent.go|internal/hooks/handler.go" docs/architecture/commands_skills_difference.md'
```

**Step 2: Run test to verify it fails**

Run:
```bash
bash -lc 'set -euo pipefail
rg -q "internal/cli/root.go" docs/architecture/commands_skills_difference.md
rg -q "internal/cli/root_test.go" docs/architecture/commands_skills_difference.md
rg -q "internal/providers/router_intent_test.go" docs/architecture/commands_skills_difference.md'
```

Expected: FAIL (at least one concrete symbol/test reference is missing).

**Step 3: Write minimal implementation**

For each high-risk row, add:

```md
- target_symbol/contract: <file + symbol or "planned symbol">
- test_reference: <existing test path or "planned test path">
- closure_condition: <explicit condition to move E0 -> E1/E2>
```

Minimum rows to complete:
- `extract-skill`
- `octo -> mp`
- `claw`, `doctor`, `schedule`, `scheduler`, `sentinel`
- `skill-claw`, `skill-doctor`

Also add matching rows to `gap-remediation-tracker.md`.

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
for token in "extract-skill" "octo" "sentinel"; do rg -q "$token" docs/architecture/commands_skills_difference.md; done
for token in "target symbol" "test_reference" "closure_condition"; do rg -q "$token" docs/architecture/commands_skills_difference.md; done
rg -q "commands_skills_difference.md" docs/architecture/gap-remediation-tracker.md'
```

Expected: PASS.

**Step 5: Commit**

```bash
git add docs/architecture/commands_skills_difference.md docs/architecture/gap-remediation-tracker.md
git commit -m "docs(architecture): add closure conditions for high-risk command/skill gaps"
```

### Task 3: Convert Script `missing` Into Domain Decision Matrix

**Why:** `script-differences.md` changed P0 strategy, but execution still lacks a domain-level matrix junior developers can apply consistently.

**What:** Add `Missing Decision Classification Matrix` with glob patterns and mandatory decision outcomes.

**How:** Group unresolved script rows by domain/pattern and define default decision + exception rule.

**Key Design:** Pattern-first classification reduces manual error across large tables.

**Files:**
- Modify: `docs/architecture/script-differences.md`
- Modify: `docs/architecture/gap-remediation-tracker.md`
- Test: `docs/architecture/script-differences.md`

**Step 1: Write the failing test**

```bash
rg -n "Missing Decision Classification Matrix" docs/architecture/script-differences.md
```

**Step 2: Run test to verify it fails**

Run:
```bash
rg -q "Missing Decision Classification Matrix" docs/architecture/script-differences.md
```

Expected: FAIL.

**Step 3: Write minimal implementation**

Add section with table:

```md
## Missing Decision Classification Matrix
| source_pattern | default_decision | decision_reason | closure_path |
|---|---|---|---|
| `scripts/scheduler/*.sh` | `DEFER_WITH_CONDITION` | ... | ... |
| `scripts/extract/*.sh` | `MIGRATE_TO_GO` | ... | ... |
| `tests/smoke/*.sh` | `MIGRATE_TO_GO` | ... | ... |
| `tests/live/*.sh` | `DEFER_WITH_CONDITION` | ... | ... |
```

Link each pattern to target package/test family (e.g. `internal/hooks/*`, `internal/workflows/*_test.go`, `internal/cli/*`).

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
rg -q "Missing Decision Classification Matrix" docs/architecture/script-differences.md
for p in "scripts/scheduler" "scripts/extract" "tests/smoke" "tests/live"; do rg -q "$p" docs/architecture/script-differences.md; done
rg -q "script-differences.md" docs/architecture/gap-remediation-tracker.md'
```

Expected: PASS.

**Step 5: Commit**

```bash
git add docs/architecture/script-differences.md docs/architecture/gap-remediation-tracker.md
git commit -m "docs(architecture): add missing decision classification matrix for scripts"
```

### Task 4: Strengthen Hook Lifecycle Index With Contract/Test Links

**Why:** Lifecycle index exists, but junior implementers still need exact contract fields and test files for completion.

**What:** Add `response_contract_fields` and `test_reference` columns.

**How:** Extend existing lifecycle table entries with concrete file/test links.

**Key Design:** Hook mapping is only actionable when contract + test are both explicit.

**Files:**
- Modify: `docs/architecture/script-differences.md`
- Test: `docs/architecture/script-differences.md`

**Step 1: Write the failing test**

```bash
bash -lc 'set -euo pipefail
rg -q "Hook Lifecycle Alignment Index" docs/architecture/script-differences.md
rg -q "response_contract_fields" docs/architecture/script-differences.md
rg -q "test_reference" docs/architecture/script-differences.md'
```

**Step 2: Run test to verify it fails**

Run:
```bash
bash -lc 'set -euo pipefail
rg -q "response_contract_fields" docs/architecture/script-differences.md
rg -q "test_reference" docs/architecture/script-differences.md'
```

Expected: FAIL.

**Step 3: Write minimal implementation**

Extend table to include:
- `response_contract_fields` (must mention `status/action/error_code/message/data/remediation`)
- `test_reference` (examples: `internal/hooks/handler_test.go`, `internal/hooks/stop_test.go`, `internal/cli/status_test.go`)

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
for evt in SessionStart UserPromptSubmit PreToolUse PostToolUse Stop SubagentStop; do
  rg -q "$evt" docs/architecture/script-differences.md
done
for token in "status" "action" "error_code" "remediation" "internal/hooks/handler_test.go"; do
  rg -q "$token" docs/architecture/script-differences.md
done'
```

Expected: PASS.

**Step 5: Commit**

```bash
git add docs/architecture/script-differences.md
git commit -m "docs(architecture): add contract and test links to hook lifecycle index"
```

### Task 5: Fill `other-differences` Remaining Partial/Missing Closure Paths

**Why:** `mcp-server/openclaw` decisions are explicit, but other partial/missing rows still lack uniform closure metadata.

**What:** Add `target_symbol_or_contract`, `evidence_upgrade_path`, and `owner_domain` for remaining high-impact rows.

**How:** Extend "关键缺口决策与契约索引" to include `.claude/settings.json`, `.mcp.json`, `docs/SCHEDULER.md`, benchmark/live docs.

**Key Design:** Keep `EXCLUDE/DEFER` auditable by linking each row to owner and trigger condition.

**Files:**
- Modify: `docs/architecture/other-differences.md`
- Modify: `docs/architecture/gap-remediation-tracker.md`
- Test: `docs/architecture/other-differences.md`

**Step 1: Write the failing test**

```bash
bash -lc 'set -euo pipefail
for token in ".claude/settings.json" ".mcp.json" "docs/SCHEDULER.md" "tests/live/README.md"; do
  rg -q "$token" docs/architecture/other-differences.md || exit 1
done
rg -q "owner_domain" docs/architecture/other-differences.md'
```

**Step 2: Run test to verify it fails**

Run:
```bash
rg -q "owner_domain" docs/architecture/other-differences.md
```

Expected: FAIL.

**Step 3: Write minimal implementation**

In `other-differences.md` add/extend closure index fields:

```md
| ... | target_symbol_or_contract | evidence_upgrade_path | owner_domain |
```

Populate at least for:
- `.claude/settings.json`
- `.mcp.json`
- `docs/SCHEDULER.md`
- `tests/benchmark/*`
- `tests/live/README.md`

Sync corresponding entries in `gap-remediation-tracker.md`.

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
for token in "target_symbol_or_contract" "evidence_upgrade_path" "owner_domain"; do
  rg -q "$token" docs/architecture/other-differences.md
done
rg -q "other-differences.md" docs/architecture/gap-remediation-tracker.md'
```

Expected: PASS.

**Step 5: Commit**

```bash
git add docs/architecture/other-differences.md docs/architecture/gap-remediation-tracker.md
git commit -m "docs(architecture): add closure metadata for remaining other-file gaps"
```

### Task 6: Add Deterministic Verification Script

**Why:** Manual checks are error-prone; future edits need one repeatable command.

**What:** Add `scripts/verify-architecture-diff-docs.sh`.

**How:** Script validates baseline consistency, legend presence, decision tokens, hook lifecycle events, and mcp/openclaw decision tags.

**Key Design:** Fail-fast with explicit error messages for junior maintainers.

**Files:**
- Create: `scripts/verify-architecture-diff-docs.sh`
- Modify: `docs/architecture/gap-remediation-tracker.md`
- Test: `scripts/verify-architecture-diff-docs.sh`

**Step 1: Write the failing test**

```bash
test -x scripts/verify-architecture-diff-docs.sh
```

**Step 2: Run test to verify it fails**

Run:
```bash
test -x scripts/verify-architecture-diff-docs.sh
```

Expected: FAIL.

**Step 3: Write minimal implementation**

Create script:

```bash
#!/usr/bin/env bash
set -euo pipefail

DOC1="docs/architecture/commands_skills_difference.md"
DOC2="docs/architecture/script-differences.md"
DOC3="docs/architecture/other-differences.md"

err() { echo "ERROR: $*" >&2; exit 1; }

vals=$(rg -o "go=[0-9a-f]{7,}" "$DOC1" "$DOC2" "$DOC3" | awk -F: '{print $NF}' | sort -u | wc -l | tr -d ' ')
[ "$vals" = "1" ] || err "go baseline hash mismatch across docs"

for f in "$DOC1" "$DOC2" "$DOC3"; do
  rg -q "E0" "$f" || err "missing E0 legend in $f"
  rg -q "E1" "$f" || err "missing E1 legend in $f"
  rg -q "E2" "$f" || err "missing E2 legend in $f"
  rg -q "E3" "$f" || err "missing E3 legend in $f"
  rg -q "decision" "$f" || err "missing decision token in $f"
done

for evt in SessionStart UserPromptSubmit PreToolUse PostToolUse Stop SubagentStop; do
  rg -q "$evt" "$DOC2" || err "missing lifecycle event $evt in script-differences.md"
done

rg -n "mcp-server/|openclaw/" "$DOC3" | rg -q "EXCLUDE_WITH_REASON|MIGRATE_TO_GO|DEFER_WITH_CONDITION" \
  || err "mcp/openclaw rows missing explicit decision"

echo "verify-architecture-diff-docs: PASS"
```

Then:
```bash
chmod +x scripts/verify-architecture-diff-docs.sh
```

**Step 4: Run test to verify it passes**

Run:
```bash
scripts/verify-architecture-diff-docs.sh
```

Expected: PASS with `verify-architecture-diff-docs: PASS`.

**Step 5: Commit**

```bash
git add scripts/verify-architecture-diff-docs.sh docs/architecture/gap-remediation-tracker.md
git commit -m "chore(architecture): add deterministic diff-doc verification script"
```

### Task 7: Final Cross-Doc Consistency Sweep And Evidence Upgrade Plan

**Why:** Remaining E0 rows need clear next-hop to E1/E2 so execution can continue without ambiguity.

**What:** Add final "E0 Upgrade Queue" section and verify no floating unresolved items without owner/action.

**How:** Update tracker and run verification commands + git diff review.

**Key Design:** No unresolved row is acceptable without `owner + next_action + decision`.

**Files:**
- Modify: `docs/architecture/gap-remediation-tracker.md`
- Test: `docs/architecture/gap-remediation-tracker.md`

**Step 1: Write the failing test**

```bash
bash -lc 'set -euo pipefail
rg -q "E0 Upgrade Queue" docs/architecture/gap-remediation-tracker.md
rg -q "| owner | next_action |" docs/architecture/gap-remediation-tracker.md'
```

**Step 2: Run test to verify it fails**

Run:
```bash
rg -q "E0 Upgrade Queue" docs/architecture/gap-remediation-tracker.md
```

Expected: FAIL.

**Step 3: Write minimal implementation**

Add:

```md
## E0 Upgrade Queue
| gap_id | current_evidence | target_evidence | owner | next_action | due |
|---|---|---|---|---|---|
```

Populate all still-`E0` high-risk rows.

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
rg -q "E0 Upgrade Queue" docs/architecture/gap-remediation-tracker.md
scripts/verify-architecture-diff-docs.sh'
```

Expected: PASS.

**Step 5: Commit**

```bash
git add docs/architecture/gap-remediation-tracker.md
git commit -m "docs(architecture): add E0 upgrade queue for remaining gaps"
```

### Task 8: Completion Verification And Review Gate

**Why:** Plan outputs are governance artifacts; completion must be evidence-backed.

**What:** Run full verification, request review, and prepare handoff summary.

**How:** Use verification + review skills and attach outputs.

**Key Design:** No "done" claim without command evidence.

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`
- Modify: `docs/architecture/script-differences.md`
- Modify: `docs/architecture/other-differences.md`
- Modify: `docs/architecture/gap-remediation-tracker.md`

**Step 1: Write the failing test**

```bash
bash -lc 'set -euo pipefail
scripts/verify-architecture-diff-docs.sh
rg -q "E0 Upgrade Queue" docs/architecture/gap-remediation-tracker.md'
```

**Step 2: Run test to verify it fails**

Run:
```bash
bash -lc 'set -euo pipefail
scripts/verify-architecture-diff-docs.sh
rg -q "E0 Upgrade Queue" docs/architecture/gap-remediation-tracker.md'
```

Expected: FAIL until all previous tasks are complete.

**Step 3: Write minimal implementation**

Complete remaining edits from Tasks 1-7 and rerun checks.

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
scripts/verify-architecture-diff-docs.sh
git diff -- docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md docs/architecture/gap-remediation-tracker.md | sed -n "1,200p"'
```

Expected: PASS; diff only contains gap-remediation scope changes.

**Step 5: Commit**

```bash
git add docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md docs/architecture/gap-remediation-tracker.md scripts/verify-architecture-diff-docs.sh
git commit -m "docs(architecture): complete gap remediation with tracker and verification guard"
```

---

## Required Skills During Execution

- `@superpowers:executing-plans` (required orchestrator)
- `@superpowers:verification-before-completion` (must run before claiming done)
- `@superpowers:requesting-code-review` (quality gate before merge/push)

