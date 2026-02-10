#!/usr/bin/env bash
# Plan evidence consistency tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing plan evidence checker..."

# Test 1: Current gap plans should pass

echo "[TEST 1] Existing gap plans satisfy evidence coverage"
if python3 scripts/check_plan_evidence.py >/tmp/test_plan_evidence_out.txt 2>/tmp/test_plan_evidence_err.txt; then
    echo "  [PASS] Existing gap plans pass evidence check"
else
    echo "  [FAIL] Existing gap plans should pass"
    cat /tmp/test_plan_evidence_err.txt
    exit 1
fi

# Test 2: Missing evidence section should fail

echo "[TEST 2] Missing evidence section is detected"
tmp_plan=$(mktemp)
cat > "$tmp_plan" <<'EOF_PLAN'
# Demo Plan

### T9-001 Demo task
- **状态**：`DONE`
- **成功判定**：
  - [x] done
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

# Test 3: Missing required evidence fields should fail

echo "[TEST 3] Missing required fields are detected"
tmp_plan=$(mktemp)
cat > "$tmp_plan" <<'EOF_PLAN'
# Demo Plan

### T9-002 Demo task
- **状态**：`DONE`

## 7) 执行与验证证据

- **Coverage Task IDs**: `T9-002`
EOF_PLAN

set +e
python3 scripts/check_plan_evidence.py "$tmp_plan" >/tmp/test_plan_evidence_tmp_out.txt 2>/tmp/test_plan_evidence_tmp_err.txt
status=$?
set -e
rm -f "$tmp_plan"

if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected checker to fail for missing required fields"
    exit 1
fi
if grep -q "missing required field" /tmp/test_plan_evidence_tmp_err.txt; then
    echo "  [PASS] Required field checks enforced"
else
    echo "  [FAIL] Expected required field diagnostic"
    cat /tmp/test_plan_evidence_tmp_err.txt
    exit 1
fi

# Test 4: Coverage IDs must include all DONE tasks

echo "[TEST 4] Coverage ID mismatch is detected"
tmp_plan=$(mktemp)
cat > "$tmp_plan" <<'EOF_PLAN'
# Demo Plan

### T9-003 Demo task
- **状态**：`DONE`

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

echo ""
echo "Plan evidence tests PASSED"
