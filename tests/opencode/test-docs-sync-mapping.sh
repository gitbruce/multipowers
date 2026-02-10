#!/usr/bin/env bash
# Docs sync mapping checker tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing docs-sync mapping checks..."

# Test 1: mapped code path without mapped docs fails.
echo "[TEST 1] mapped code change requires mapped docs"
set +e
output=$(python3 scripts/check_docs_sync.py --changed-file bin/multipowers 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] expected mapped docs-sync failure"
    exit 1
fi
if echo "$output" | grep -q "required docs"; then
    echo "  [PASS] mapped docs failure is actionable"
else
    echo "  [FAIL] mapped docs diagnostics missing"
    echo "  Output: $output"
    exit 1
fi

# Test 2: mapped docs change satisfies rule.
echo "[TEST 2] mapped docs update passes"
if python3 scripts/check_docs_sync.py --changed-file bin/multipowers --changed-file README.md >/tmp/test_docs_sync_map_out.txt 2>/tmp/test_docs_sync_map_err.txt; then
    echo "  [PASS] mapped docs rule passes with README update"
else
    echo "  [FAIL] mapped docs should pass when README changed"
    cat /tmp/test_docs_sync_map_err.txt
    exit 1
fi

# Test 3: ignored test-only change passes without docs.
echo "[TEST 3] test-only changes are ignored"
if python3 scripts/check_docs_sync.py --changed-file tests/opencode/test-docs-sync-mapping.sh >/tmp/test_docs_sync_testonly_out.txt 2>/tmp/test_docs_sync_testonly_err.txt; then
    echo "  [PASS] test-only changes do not require docs"
else
    echo "  [FAIL] test-only changes should not require docs"
    cat /tmp/test_docs_sync_testonly_err.txt
    exit 1
fi

echo ""
echo "Docs-sync mapping tests PASSED"
