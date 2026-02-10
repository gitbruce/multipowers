#!/usr/bin/env bash
# ask-role argument handling tests
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

echo "Testing ask-role argument handling..."

./bin/multipowers init >/dev/null 2>&1

# Test 1: Non-existent role should produce readable error

echo "[TEST 1] Non-existent role error"
set +e
output=$(./bin/ask-role nonexistent_role "test prompt" 2>&1)
status=$?
set -e

if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected non-zero exit"
    exit 1
fi
if echo "$output" | grep -q "Role 'nonexistent_role' not found"; then
    echo "  [PASS] Error message is readable"
else
    echo "  [FAIL] Expected readable role error"
    echo "  Output: $output"
    exit 1
fi

# Test 2: Roles with multiple args should parse without shell errors

echo "[TEST 2] Multi-arg role parsing"
set +e
output=$(./bin/ask-role hephaestus "test prompt" 2>&1)
status=$?
set -e

if [ $status -eq 0 ]; then
    echo "  [PASS] hephaestus call completed"
else
    if echo "$output" | grep -qE "Failed to call codex|Command failed|not found"; then
        echo "  [PASS] Args parsed; failure occurred at connector execution stage"
    else
        echo "  [FAIL] Unexpected failure before connector stage"
        echo "  Output: $output"
        exit 1
    fi
fi

# Test 3: Single-arg role should parse without errors

echo "[TEST 3] Single-arg role parsing"
set +e
output=$(./bin/ask-role prometheus "test prompt" 2>&1)
status=$?
set -e

if [ $status -eq 0 ]; then
    echo "  [PASS] prometheus call completed"
else
    if echo "$output" | grep -qE "Failed to call gemini|Command failed|not found"; then
        echo "  [PASS] Args parsed; failure occurred at connector execution stage"
    else
        echo "  [FAIL] Unexpected failure before connector stage"
        echo "  Output: $output"
        exit 1
    fi
fi

echo ""
echo "All ask-role args tests PASSED"
