#!/usr/bin/env bash
# Plan evidence consistency tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing plan evidence checker..."

# Test 1: Mixed-language DONE status is recognized and passes with evidence.
echo "[TEST 1] Mixed-language DONE status is supported"
tmp_plan=$(mktemp)
cat > "$tmp_plan" <<'EOF_PLAN'
# Demo Plan

| Task ID | Status | Priority | Owner | Depends On |
|---|---|---|---|---|
| T9-000 | `TODO` | P0 | Demo | - |

### T9-000 Demo task
- **Status**: `DONE`
- **状态**：`DONE`

## 7) Evidence

- **Coverage Task IDs**: `T9-000`
- **Date**: `2026-02-10`
- **Verifier**: `tester`
- **Command(s)**:
  - `echo test`
- **Exit Code**: `0`
- **Key Output**:
  - `ok`
EOF_PLAN

if python3 scripts/check_plan_evidence.py "$tmp_plan" >/tmp/test_plan_evidence_out.txt 2>/tmp/test_plan_evidence_err.txt; then
    echo "  [PASS] Mixed-language DONE status recognized"
else
    echo "  [FAIL] Expected evidence checker to pass"
    cat /tmp/test_plan_evidence_err.txt
    rm -f "$tmp_plan"
    exit 1
fi
rm -f "$tmp_plan"

# Test 2: Missing evidence section should fail.
echo "[TEST 2] Missing evidence section is detected"
tmp_plan=$(mktemp)
cat > "$tmp_plan" <<'EOF_PLAN'
# Demo Plan

### T9-001 Demo task
- **状态**：`DONE`
EOF_PLAN

set +e
python3 scripts/check_plan_evidence.py "$tmp_plan" >/tmp/test_plan_evidence_tmp_out.txt 2>/tmp/test_plan_evidence_tmp_err.txt
status=$?
set -e
rm -f "$tmp_plan"

if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected checker to fail for missing evidence"
    exit 1
fi
if grep -q "no evidence section found" /tmp/test_plan_evidence_tmp_err.txt; then
    echo "  [PASS] Missing evidence detected"
else
    echo "  [FAIL] Expected missing evidence diagnostic"
    cat /tmp/test_plan_evidence_tmp_err.txt
    exit 1
fi

# Test 3: Coverage IDs must include all DONE tasks.
echo "[TEST 3] Coverage ID mismatch is detected"
tmp_plan=$(mktemp)
cat > "$tmp_plan" <<'EOF_PLAN'
# Demo Plan

### T9-003 Demo task
- **Status**: `DONE`

## 7) 执行与验证证据

- **Coverage Task IDs**: `T9-999`
- **Date**: `2026-02-10`
- **Verifier**: `tester`
- **Command(s)**:
  - `echo test`
- **Exit Code**: `0`
- **Key Output**:
  - `ok`
EOF_PLAN

set +e
python3 scripts/check_plan_evidence.py "$tmp_plan" >/tmp/test_plan_evidence_tmp_out.txt 2>/tmp/test_plan_evidence_tmp_err.txt
status=$?
set -e
rm -f "$tmp_plan"

if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected checker to fail for coverage mismatch"
    exit 1
fi
if grep -q "missing coverage IDs" /tmp/test_plan_evidence_tmp_err.txt; then
    echo "  [PASS] Coverage mismatch detected"
else
    echo "  [FAIL] Expected coverage mismatch diagnostic"
    cat /tmp/test_plan_evidence_tmp_err.txt
    exit 1
fi

# Test 4: Status updater updates table + section fields.
echo "[TEST 4] Status updater updates section and board row"
tmp_plan=$(mktemp)
cat > "$tmp_plan" <<'EOF_PLAN'
# Demo Plan

| Task ID | Status | Priority | Owner | Depends On |
|---|---|---|---|---|
| T9-010 | `TODO` | P1 | Demo | - |

### T9-010 Demo task
- **Status**: `TODO`
- **状态**：`TODO`
EOF_PLAN

python3 scripts/update_plan_task_status.py --file "$tmp_plan" --task-id T9-010 --status DONE >/tmp/test_plan_status_out.txt 2>/tmp/test_plan_status_err.txt

if grep -Fq -- '| T9-010 | `DONE` |' "$tmp_plan" \
   && grep -Fq -- '- **Status**: `DONE`' "$tmp_plan" \
   && grep -Fq -- '- **状态**：`DONE`' "$tmp_plan"; then
    echo "  [PASS] Status updater handled table + dual fields"
else
    echo "  [FAIL] Status updater did not update expected fields"
    cat "$tmp_plan"
    rm -f "$tmp_plan"
    exit 1
fi
rm -f "$tmp_plan"

echo ""
echo "Plan evidence tests PASSED"
