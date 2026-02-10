#!/usr/bin/env bash
# Template sync candidate checker tests
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

echo "Testing template sync candidates checker..."

workspace=$(create_workspace)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cd "$workspace"
./bin/multipowers init --repair >/tmp/test_template_sync_init_out.txt 2>/tmp/test_template_sync_init_err.txt

# Test 1: no-drift scenario passes.
echo "[TEST 1] no drift passes"
if python3 scripts/check_template_sync_candidates.py >/tmp/test_template_sync_nodrift_out.txt 2>/tmp/test_template_sync_nodrift_err.txt; then
    echo "  [PASS] no-drift check passes"
else
    echo "  [FAIL] no-drift check should pass"
    cat /tmp/test_template_sync_nodrift_err.txt
    exit 1
fi

# Test 2: drift scenario is detected.
echo "[TEST 2] drift is reported"
echo "\nMaintainer delta" >> conductor/context/workflow.md
set +e
output=$(python3 scripts/check_template_sync_candidates.py 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] expected drift checker to fail when drift exists"
    exit 1
fi
if echo "$output" | grep -q "drift" && echo "$output" | grep -q "workflow.md"; then
    echo "  [PASS] drift checker reports changed mapping"
else
    echo "  [FAIL] drift diagnostics missing"
    echo "  Output: $output"
    exit 1
fi

# Test 3: doctor surfaces drift as warning only.
echo "[TEST 3] doctor emits warning but remains non-blocking"
set +e
doctor_output=$(MULTIPOWERS_CONTEXT_MODE=lenient ./bin/multipowers doctor 2>&1)
doctor_status=$?
set -e
if [ $doctor_status -ne 0 ]; then
    echo "  [FAIL] doctor should remain non-blocking for template drift"
    echo "  Output: $doctor_output"
    exit 1
fi
if echo "$doctor_output" | grep -qi "template" && echo "$doctor_output" | grep -qi "warn"; then
    echo "  [PASS] doctor surfaces template drift warning"
else
    echo "  [FAIL] doctor missing template drift warning"
    echo "  Output: $doctor_output"
    exit 1
fi

echo ""
echo "Template sync candidate tests PASSED"
