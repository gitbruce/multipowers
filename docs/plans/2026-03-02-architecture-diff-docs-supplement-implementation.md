# Architecture Diff Docs Supplement Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make architecture diff documentation consistent, product-aligned, and explicitly decision-driven without forcing one-to-one migration of all `main` files/features.

**Architecture:** Apply consistency-first updates to `commands_skills_difference.md`, `script-differences.md`, and `other-differences.md`; then add a shared evidence/decision schema (`source -> target -> symbol/contract -> evidence -> decision`) for all `partial`/`missing` mappings. Migration decisions are governed by `.multipowers/product-guidelines.md` and `.multipowers/product.md`, not by mechanical parity.

**Tech Stack:** Markdown docs, bash/rg consistency checks, git.

---

### Task 1: Baseline And Vocabulary Drift Detection

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`
- Modify: `docs/architecture/script-differences.md`
- Modify: `docs/architecture/other-differences.md`
- Test: `docs/architecture/commands_skills_difference.md`

**Step 1: Write the failing test (consistency check command)**

```bash
bash -lc 'set -euo pipefail
rg -o "go=[0-9a-f]{7,}" docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md | awk -F: "{print \$NF}" | sort -u'
```

**Step 2: Run test to verify it fails**

Run:
```bash
bash -lc 'set -euo pipefail
vals=$(rg -o "go=[0-9a-f]{7,}" docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md | awk -F: "{print \$NF}" | sort -u | wc -l)
[ "$vals" -eq 1 ]'
```

Expected: FAIL (currently go baseline hashes are not fully aligned).

**Step 3: Record expected status dictionary**

Add or normalize this set in each doc (exact tokens):

```md
`equivalent` / `partial` / `missing` / `intentional-diff`
```

**Step 4: Run check to verify vocabulary is present**

Run:
```bash
rg -n "equivalent|partial|missing|intentional-diff" docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md
```

Expected: PASS (all 3 docs contain the same four status terms).

**Step 5: Commit**

```bash
git add docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md
git commit -m "docs(architecture): normalize baseline and status vocabulary"
```

### Task 2: Fix `other-differences.md` Baseline And Scope Contract

**Files:**
- Modify: `docs/architecture/other-differences.md`
- Test: `docs/architecture/other-differences.md`

**Step 1: Write the failing test**

```bash
rg -n "go=d01e74d99977" docs/architecture/other-differences.md
```

**Step 2: Run test to verify it fails**

Run:
```bash
test "$(rg -c "go=d01e74d99977" docs/architecture/other-differences.md)" -eq 0
```

Expected: FAIL (old go baseline still exists).

**Step 3: Write minimal implementation**

Update header baseline line to the current aligned pair used by the other two docs:

```md
基线提交：`main=<same-main-sha>`，`go=<same-go-sha-as-other-2-docs>`
```

Add scope boundary text:

```md
本文件仅覆盖 non commands/skills/scripts 差异，不替代另外两份差异文档。
```

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
test "$(rg -c "go=d01e74d99977" docs/architecture/other-differences.md)" -eq 0
vals=$(rg -o "go=[0-9a-f]{7,}" docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md | awk -F: "{print \$NF}" | sort -u | wc -l)
[ "$vals" -eq 1 ]'
```

Expected: PASS.

**Step 5: Commit**

```bash
git add docs/architecture/other-differences.md
git commit -m "docs(architecture): align other-differences baseline and scope"
```

### Task 3: Add Unified Evidence Legend (`E0-E3`) To All Three Docs

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`
- Modify: `docs/architecture/script-differences.md`
- Modify: `docs/architecture/other-differences.md`
- Test: `docs/architecture/script-differences.md`

**Step 1: Write the failing test**

```bash
bash -lc 'set -euo pipefail
for f in docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md; do
  rg -q "E0|E1|E2|E3" "$f" || { echo "missing evidence legend in $f"; exit 1; }
done'
```

**Step 2: Run test to verify it fails**

Run:
```bash
bash -lc 'set -euo pipefail
for f in docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md; do
  rg -q "E0|E1|E2|E3" "$f" || exit 1
done'
```

Expected: FAIL for at least one file.

**Step 3: Write minimal implementation**

Insert this block in each file:

```md
## Evidence Legend
- `E0`: doc-only plan
- `E1`: symbol exists
- `E2`: test exists
- `E3`: verified output recorded
```

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
for f in docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md; do
  rg -q "E0" "$f"
  rg -q "E1" "$f"
  rg -q "E2" "$f"
  rg -q "E3" "$f"
done'
```

Expected: PASS.

**Step 5: Commit**

```bash
git add docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md
git commit -m "docs(architecture): add unified evidence legend for mapping rows"
```

### Task 4: Add Decision Columns For `partial/missing` Rows (Product-Gated)

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`
- Modify: `docs/architecture/script-differences.md`
- Modify: `docs/architecture/other-differences.md`
- Test: `docs/architecture/commands_skills_difference.md`

**Step 1: Write the failing test**

```bash
bash -lc 'set -euo pipefail
for f in docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md; do
  rg -q "decision" "$f" || { echo "missing decision field in $f"; exit 1; }
done'
```

**Step 2: Run test to verify it fails**

Run:
```bash
bash -lc 'set -euo pipefail
for f in docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md; do
  rg -q "MIGRATE_TO_GO|COPY_FROM_MAIN|EXCLUDE_WITH_REASON|DEFER_WITH_CONDITION" "$f" || exit 1
done'
```

Expected: FAIL for one or more files.

**Step 3: Write minimal implementation**

Add or extend table columns for all `partial/missing` rows:

```md
| ... | target_symbol_or_contract | evidence_level | decision | decision_reason |
```

Allowed `decision` values:

```md
`MIGRATE_TO_GO` | `COPY_FROM_MAIN` | `EXCLUDE_WITH_REASON` | `DEFER_WITH_CONDITION`
```

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
for f in docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md; do
  rg -q "MIGRATE_TO_GO" "$f"
  rg -q "COPY_FROM_MAIN" "$f"
done'
```

Expected: PASS (at least these two core decisions are present across docs).

**Step 5: Commit**

```bash
git add docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md
git commit -m "docs(architecture): add product-gated migration decision fields"
```

### Task 5: Resolve Floating `missing` In `other-differences` High-Risk Domains

**Files:**
- Modify: `docs/architecture/other-differences.md`
- Test: `docs/architecture/other-differences.md`

**Step 1: Write the failing test**

```bash
rg -n "mcp-server/|openclaw/" docs/architecture/other-differences.md
```

**Step 2: Run test to verify it fails**

Run:
```bash
bash -lc 'set -euo pipefail
rows=$(rg -n "mcp-server/|openclaw/" docs/architecture/other-differences.md | wc -l)
[ "$rows" -gt 0 ]
rg -n "mcp-server/|openclaw/" docs/architecture/other-differences.md | rg -vq "EXCLUDE_WITH_REASON|MIGRATE_TO_GO|DEFER_WITH_CONDITION"'
```

Expected: FAIL (currently rows exist without explicit final decision class).

**Step 3: Write minimal implementation**

For each `mcp-server/*` and `openclaw/*` row, fill:

```md
decision: EXCLUDE_WITH_REASON | MIGRATE_TO_GO | DEFER_WITH_CONDITION
decision_reason: <product-aligned rationale>
target_symbol_or_contract: <if migrate/defer>
```

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
rg -n "mcp-server/|openclaw/" docs/architecture/other-differences.md | rg -q "EXCLUDE_WITH_REASON|MIGRATE_TO_GO|DEFER_WITH_CONDITION"'
```

Expected: PASS.

**Step 5: Commit**

```bash
git add docs/architecture/other-differences.md
git commit -m "docs(architecture): resolve floating missing decisions for mcp-server/openclaw domains"
```

### Task 6: Add Hook Lifecycle Index In `script-differences`

**Files:**
- Modify: `docs/architecture/script-differences.md`
- Test: `docs/architecture/script-differences.md`

**Step 1: Write the failing test**

```bash
bash -lc 'set -euo pipefail
for evt in SessionStart PreToolUse PostToolUse Stop SubagentStop; do
  rg -q "$evt" docs/architecture/script-differences.md || { echo "missing $evt"; exit 1; }
done'
```

**Step 2: Run test to verify it fails**

Run:
```bash
bash -lc 'set -euo pipefail
for evt in SessionStart PreToolUse PostToolUse Stop SubagentStop; do
  rg -q "$evt" docs/architecture/script-differences.md || exit 1
done'
```

Expected: FAIL if lifecycle index is absent/incomplete.

**Step 3: Write minimal implementation**

Add section:

```md
## Hook Lifecycle Alignment Index
| lifecycle_event | legacy_script_source | go_target_symbol | evidence_level | note |
|---|---|---|---|---|
```

Populate key rows with existing mapped hooks.

**Step 4: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
for evt in SessionStart PreToolUse PostToolUse Stop SubagentStop; do
  rg -q "$evt" docs/architecture/script-differences.md
done'
```

Expected: PASS.

**Step 5: Commit**

```bash
git add docs/architecture/script-differences.md
git commit -m "docs(architecture): add hook lifecycle alignment index for script mapping"
```

### Task 7: Final Verification And Handoff

**Files:**
- Modify: `docs/architecture/commands_skills_difference.md`
- Modify: `docs/architecture/script-differences.md`
- Modify: `docs/architecture/other-differences.md`
- Test: `docs/architecture/other-differences.md`

**Step 1: Write the failing test checklist command**

```bash
bash -lc 'set -euo pipefail
vals=$(rg -o "go=[0-9a-f]{7,}" docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md | awk -F: "{print \$NF}" | sort -u | wc -l)
[ "$vals" -eq 1 ]
for f in docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md; do
  rg -q "E0" "$f"
  rg -q "decision" "$f"
done'
```

**Step 2: Run test to verify it passes**

Run:
```bash
bash -lc 'set -euo pipefail
vals=$(rg -o "go=[0-9a-f]{7,}" docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md | awk -F: "{print \$NF}" | sort -u | wc -l)
[ "$vals" -eq 1 ]
for f in docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md; do
  rg -q "E0" "$f"
  rg -q "E1" "$f"
  rg -q "E2" "$f"
  rg -q "E3" "$f"
done'
```

Expected: PASS.

**Step 3: Final documentation quality check**

Run:
```bash
git diff -- docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md
```

Expected: all changes are doc-structure/evidence/decision focused; no scope drift.

**Step 4: Commit**

```bash
git add docs/architecture/commands_skills_difference.md docs/architecture/script-differences.md docs/architecture/other-differences.md
git commit -m "docs(architecture): complete consistency-first supplement and product-gated decisions"
```

**Step 5: Pre-push verification skill**

Run `@superpowers:verification-before-completion` and capture verification output in the execution log before push.

