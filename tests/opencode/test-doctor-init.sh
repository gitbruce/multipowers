#!/usr/bin/env bash
# Doctor and init tests
# Verifies context strategy consistency and safe init modes without polluting repo workspace

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

snapshot_conductor_state() {
    if [ ! -d "$REPO_ROOT/conductor" ]; then
        echo "absent"
        return
    fi

    tar -cf - -C "$REPO_ROOT" conductor 2>/dev/null | sha256sum | awk '{print $1}'
}

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

echo "Testing doctor and init commands..."

before_snapshot=$(snapshot_conductor_state)
workspace=$(create_workspace)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cd "$workspace"

# Test 1: Doctor should fail when context is missing in strict mode

echo "[TEST 1] Doctor fails with missing context (strict)"
set +e
output=$(./bin/multipowers doctor 2>&1)
status=$?
set -e

if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected doctor to fail before init"
    exit 1
fi
if echo "$output" | grep -q "context directory missing"; then
    echo "  [PASS] Doctor reports missing context"
else
    echo "  [FAIL] Doctor did not report missing context"
    echo "  Output: $output"
    exit 1
fi

# Test 2: --repair should non-destructively fill missing context files

echo "[TEST 2] Init --repair is non-destructive"
mkdir -p conductor/config conductor/tracks
echo "custom" > conductor/config/custom.cfg
echo "# existing track" > conductor/tracks/track-000-existing.md

./bin/multipowers init --repair >/tmp/test_doctor_init_repair_out.txt 2>/tmp/test_doctor_init_repair_err.txt

required_files=(
    "conductor/context/product.md"
    "conductor/context/product-guidelines.md"
    "conductor/context/workflow.md"
    "conductor/context/tech-stack.md"
)
for required in "${required_files[@]}"; do
    if [ -f "$required" ]; then
        echo "  [PASS] $required created"
    else
        echo "  [FAIL] Missing required context file: $required"
        exit 1
    fi
done

if [ -f "conductor/config/custom.cfg" ] && [ -f "conductor/tracks/track-000-existing.md" ]; then
    echo "  [PASS] Existing config/track files preserved"
else
    echo "  [FAIL] Existing files were unexpectedly removed"
    exit 1
fi

# Test 3: Doctor passes after repair

echo "[TEST 3] Doctor passes after repair"
if output=$(./bin/multipowers doctor 2>&1); then
    if echo "$output" | grep -q "context file exists"; then
        echo "  [PASS] Doctor validates context files"
    else
        echo "  [FAIL] Doctor missing context validation output"
        exit 1
    fi
else
    echo "  [FAIL] Doctor should pass after repair"
    echo "  Output: $output"
    exit 1
fi

# Test 4: --force requires --yes

echo "[TEST 4] Init --force requires --yes"
set +e
output=$(./bin/multipowers init --force 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected init --force without --yes to fail"
    exit 1
fi
if echo "$output" | grep -q "requires --yes"; then
    echo "  [PASS] --force confirmation required"
else
    echo "  [FAIL] Missing confirmation requirement message"
    exit 1
fi

# Test 5: --force --yes recreates structure

echo "[TEST 5] Init --force --yes recreates structure"
echo "temporary" > conductor/context/to-remove.md
./bin/multipowers init --force --yes >/tmp/test_doctor_init_force_out.txt 2>/tmp/test_doctor_init_force_err.txt
if [ ! -f "conductor/context/to-remove.md" ]; then
    echo "  [PASS] --force --yes recreates structure"
else
    echo "  [FAIL] --force --yes did not recreate structure"
    exit 1
fi

# Test 6: Lenient mode downgrades missing context to warnings

echo "[TEST 6] Lenient mode allows doctor to continue"
rm -rf conductor/context
if output=$(MULTIPOWERS_CONTEXT_MODE=lenient ./bin/multipowers doctor 2>&1); then
    if echo "$output" | grep -q "\[WARN\] context directory missing"; then
        echo "  [PASS] Lenient mode uses warnings"
    else
        echo "  [FAIL] Expected warning in lenient mode"
        exit 1
    fi
else
    echo "  [FAIL] Lenient mode doctor should not fail"
    echo "  Output: $output"
    exit 1
fi

# Test 7: Repository workspace must remain unchanged

echo "[TEST 7] Repository workspace remains unchanged"
after_snapshot=$(snapshot_conductor_state)
if [ "$before_snapshot" = "$after_snapshot" ]; then
    echo "  [PASS] No repository conductor pollution"
else
    echo "  [FAIL] Repository conductor state was modified"
    echo "  Before: $before_snapshot"
    echo "  After : $after_snapshot"
    exit 1
fi

echo ""
echo "All doctor and init tests PASSED"
