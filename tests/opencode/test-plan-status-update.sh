#!/usr/bin/env bash
# Plan status updater tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing plan status updater..."

tmp_plan=$(mktemp)
cleanup() {
    rm -f "$tmp_plan"
}
trap cleanup EXIT

cat > "$tmp_plan" <<'PLAN_EOF'
# Demo Plan

| Task ID | Status | Priority | Owner | Depends On |
|---|---|---|---|---|
| T9-001 | `TODO` | P0 | Integrator | - |

### T9-001 Demo task
- **Status**: `TODO`
- **状态**：`TODO`
- **Owner**: Integrator

## Evidence
- **Coverage Task IDs**: `T9-001`
- **Date**: `2026-02-10`
- **Verifier**: `tester`
- **Command(s)**:
  - `echo test`
- **Exit Code**: `0`
- **Key Output**:
  - `ok`
PLAN_EOF

# Test 1: IN_PROGRESS updates task section + task board.
echo "[TEST 1] Update to IN_PROGRESS updates both section and table"
python3 scripts/update_plan_task_status.py --file "$tmp_plan" --task-id T9-001 --status IN_PROGRESS >/tmp/test_plan_status_update_out.txt
if grep -q '\| T9-001 \| `IN_PROGRESS` \|' "$tmp_plan" && \
   grep -q '\*\*Status\*\*: `IN_PROGRESS`' "$tmp_plan" && \
   grep -q '\*\*状态\*\*：`IN_PROGRESS`' "$tmp_plan"; then
    echo "  [PASS] Section and table statuses updated"
else
    echo "  [FAIL] Status update did not modify expected fields"
    cat "$tmp_plan"
    exit 1
fi

# Test 2: DONE status recognized by evidence checker.
echo "[TEST 2] DONE status is recognized by evidence checker"
python3 scripts/update_plan_task_status.py --file "$tmp_plan" --task-id T9-001 --status DONE >/tmp/test_plan_status_update_out2.txt
if python3 scripts/check_plan_evidence.py "$tmp_plan" >/tmp/test_plan_status_update_chk_out.txt 2>/tmp/test_plan_status_update_chk_err.txt; then
    echo "  [PASS] Evidence checker recognizes updated DONE status"
else
    echo "  [FAIL] Evidence checker should pass for updated DONE task"
    cat /tmp/test_plan_status_update_chk_err.txt
    exit 1
fi

echo ""
echo "Plan status updater tests PASSED"
