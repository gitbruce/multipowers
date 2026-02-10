#!/usr/bin/env bash
# Plan governance evidence rule tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing plan governance evidence rule..."

tmp_plan=$(mktemp)
cleanup() {
    rm -f "$tmp_plan"
}
trap cleanup EXIT

cat > "$tmp_plan" <<'PLAN_EOF'
# Demo Plan

### T8-999 Demo task
- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator

## Evidence
- **Coverage Task IDs**: `T8-999`
- **Date**: `2026-02-10`
- **Verifier**: `tester`
- **Command(s)**:
  - `bash tests/opencode/run-tests.sh`
- **Exit Code**: `0`
- **Key Output**:
  - `all pass`
PLAN_EOF

# Test 1: require-governance-evidence fails without governance proof.
echo "[TEST 1] governance evidence can be required"
set +e
output=$(python3 scripts/check_plan_evidence.py --require-governance-evidence "$tmp_plan" 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] expected missing governance evidence to fail"
    exit 1
fi
if echo "$output" | grep -q "governance evidence"; then
    echo "  [PASS] missing governance evidence is detected"
else
    echo "  [FAIL] missing governance evidence diagnostics absent"
    echo "  Output: $output"
    exit 1
fi

# Test 2: governance reference passes strict evidence mode.
echo "[TEST 2] governance reference satisfies requirement"
cat > "$tmp_plan" <<'PLAN_EOF'
# Demo Plan

### T8-999 Demo task
- **Status**: `DONE`
- **状态**：`DONE`
- **Owner**: Integrator

## Evidence
- **Coverage Task IDs**: `T8-999`
- **Date**: `2026-02-10`
- **Verifier**: `tester`
- **Command(s)**:
  - `bash scripts/run_governance_checks.sh --mode strict --changed-file bin/multipowers`
- **Exit Code**: `0`
- **Key Output**:
  - `governance artifact: outputs/governance/demo.json`
PLAN_EOF

if python3 scripts/check_plan_evidence.py --require-governance-evidence "$tmp_plan" >/tmp/test_plan_gov_out.txt 2>/tmp/test_plan_gov_err.txt; then
    echo "  [PASS] governance evidence requirement satisfied"
else
    echo "  [FAIL] governance evidence should pass"
    cat /tmp/test_plan_gov_err.txt
    exit 1
fi

echo ""
echo "Plan governance evidence tests PASSED"
