#!/usr/bin/env bash
# Governance checks pipeline tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing governance checks..."

# Test 1: docs-only changes pass.
echo "[TEST 1] docs-only change passes"
if output=$(bash scripts/run_governance_checks.sh --changed-file README.md 2>&1); then
    if echo "$output" | grep -q "PASS"; then
        echo "  [PASS] docs-only governance check passed"
    else
        echo "  [FAIL] Missing PASS signal"
        echo "  Output: $output"
        exit 1
    fi
else
    echo "  [FAIL] docs-only change should pass"
    echo "  Output: $output"
    exit 1
fi

# Test 2: code changes without docs should fail docs-sync.
echo "[TEST 2] missing docs update is detected"
set +e
output=$(bash scripts/run_governance_checks.sh --changed-file scripts/route_task.py 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected governance to fail for docs-sync violation"
    exit 1
fi
if echo "$output" | grep -q "docs update required" && echo "$output" | grep -q "Changed code files"; then
    echo "  [PASS] Docs-sync violation is actionable"
else
    echo "  [FAIL] Missing actionable docs-sync message"
    echo "  Output: $output"
    exit 1
fi

# Test 3: code + docs together pass docs-sync gate.
echo "[TEST 3] code+docs changes satisfy docs-sync"
if output=$(bash scripts/run_governance_checks.sh --changed-file scripts/route_task.py --changed-file README.md 2>&1); then
    if echo "$output" | grep -q "PASS"; then
        echo "  [PASS] code+docs changes pass governance"
    else
        echo "  [FAIL] Missing PASS signal for code+docs"
        echo "  Output: $output"
        exit 1
    fi
else
    echo "  [FAIL] code+docs should pass docs-sync gate"
    echo "  Output: $output"
    exit 1
fi

# Test 4: npm governance script is wired.
echo "[TEST 4] npm governance:check script exists"
if npm run -s governance:check -- --changed-file README.md >/tmp/test_governance_npm_out.txt 2>/tmp/test_governance_npm_err.txt; then
    echo "  [PASS] npm governance:check is executable"
else
    echo "  [FAIL] npm governance:check should run"
    cat /tmp/test_governance_npm_err.txt
    exit 1
fi

echo ""
echo "Governance checks tests PASSED"
