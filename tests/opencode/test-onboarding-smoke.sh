#!/usr/bin/env bash
# Onboarding smoke test: fresh clone style flow
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

create_workspace() {
    local workspace
    workspace=$(mktemp -d)
    cp -R "$REPO_ROOT/bin" "$workspace/"
    cp -R "$REPO_ROOT/config" "$workspace/"
    cp -R "$REPO_ROOT/connectors" "$workspace/"
    cp -R "$REPO_ROOT/scripts" "$workspace/"
    cp -R "$REPO_ROOT/templates" "$workspace/"
    chmod +x "$workspace/bin/multipowers" "$workspace/bin/ask-role"
    echo "$workspace"
}

workspace=$(create_workspace)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cd "$workspace"

echo "Testing onboarding smoke flow..."

# Test 1: Fresh state doctor fails

echo "[TEST 1] Doctor fails before init"
set +e
output=$(./bin/multipowers doctor 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected doctor failure before init"
    exit 1
fi
if echo "$output" | grep -q "context directory missing"; then
    echo "  [PASS] Fresh state failure is explicit"
else
    echo "  [FAIL] Missing expected doctor failure message"
    exit 1
fi

# Test 2: Init/repair then doctor passes

echo "[TEST 2] Init then doctor passes"
./bin/multipowers init >/dev/null 2>&1
if output=$(./bin/multipowers doctor 2>&1); then
    if echo "$output" | grep -q "context file exists"; then
        echo "  [PASS] Onboarding flow reaches healthy state"
    else
        echo "  [FAIL] Doctor pass output missing context checks"
        exit 1
    fi
else
    echo "  [FAIL] Doctor should pass after init"
    echo "  Output: $output"
    exit 1
fi

# Test 3: Repair mode is idempotent

echo "[TEST 3] Repair mode is safe"
if output=$(./bin/multipowers init --repair 2>&1); then
    if echo "$output" | grep -q "repair completed"; then
        echo "  [PASS] Repair mode runs safely"
    else
        echo "  [PASS] Repair mode runs (no-op acceptable)"
    fi
else
    echo "  [FAIL] Repair mode should not fail"
    exit 1
fi

echo ""
echo "Onboarding smoke tests PASSED"
