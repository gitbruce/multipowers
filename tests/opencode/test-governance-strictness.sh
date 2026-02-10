#!/usr/bin/env bash
# Governance strict/advisory policy tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing governance strictness policy..."

fake_path=$(mktemp -d)
cleanup() {
    rm -rf "$fake_path"
}
trap cleanup EXIT

ln -s "$(command -v python3)" "$fake_path/python3"
ln -s "$(command -v git)" "$fake_path/git"
ln -s "$(command -v bash)" "$fake_path/bash"
ln -s "$(command -v dirname)" "$fake_path/dirname"

# Test 1: strict mode fails if required tools unavailable for code files.
echo "[TEST 1] strict mode fails on missing governance tool"
set +e
output=$(PATH="$fake_path" bash scripts/run_governance_checks.sh --mode strict --changed-file scripts/route_task.py --changed-file README.md 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] strict mode should fail without tooling"
    exit 1
fi
if echo "$output" | grep -Eq "required tool missing|FAIL"; then
    echo "  [PASS] strict mode blocks on missing tooling"
else
    echo "  [FAIL] strict mode diagnostics missing"
    echo "  Output: $output"
    exit 1
fi

# Test 2: advisory mode allows missing tool and passes.
echo "[TEST 2] advisory mode allows missing governance tool"
if output=$(PATH="$fake_path" bash scripts/run_governance_checks.sh --mode advisory --changed-file scripts/route_task.py --changed-file README.md 2>&1); then
    if echo "$output" | grep -q "PASS"; then
        echo "  [PASS] advisory mode keeps workflow moving"
    else
        echo "  [FAIL] advisory mode missing PASS signal"
        echo "  Output: $output"
        exit 1
    fi
else
    echo "  [FAIL] advisory mode should not hard fail"
    echo "  Output: $output"
    exit 1
fi

# Test 3: npm governance defaults to strict mode.
echo "[TEST 3] npm governance script is strict-ready"
if npm run -s governance -- --help >/tmp/test_gov_strict_help_out.txt 2>/tmp/test_gov_strict_help_err.txt; then
    if grep -q -- "--mode" /tmp/test_gov_strict_help_out.txt || grep -q -- "--mode" /tmp/test_gov_strict_help_err.txt; then
        echo "  [PASS] governance help exposes mode flag"
    else
        echo "  [FAIL] governance mode help missing"
        exit 1
    fi
else
    echo "  [FAIL] npm governance --help should run"
    cat /tmp/test_gov_strict_help_err.txt
    exit 1
fi

echo ""
echo "Governance strictness tests PASSED"
