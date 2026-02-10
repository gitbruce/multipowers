#!/usr/bin/env bash
# multipowers route/workflow command integration tests
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

echo "Testing multipowers route/workflow commands..."

workspace=$(create_workspace)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cd "$workspace"

# Test 1: route command works and returns JSON lane data.
echo "[TEST 1] route command outputs lane JSON"
out=$(./bin/multipowers route --task "Fix typo in readme" --json)
if python3 - "$out" <<'PY'
import json
import sys
payload = json.loads(sys.argv[1])
if payload.get("lane") not in {"fast", "standard"}:
    raise SystemExit(1)
if "reason" not in payload:
    raise SystemExit(1)
print("ok")
PY
then
    echo "  [PASS] route command returns valid JSON"
else
    echo "  [FAIL] route command JSON invalid"
    echo "  Output: $out"
    exit 1
fi

# Test 2: workflow list and validate commands work.
echo "[TEST 2] workflow list/validate commands"
if list_out=$(./bin/multipowers workflow list 2>&1); then
    if echo "$list_out" | grep -q "subagent-driven-development"; then
        echo "  [PASS] workflow list shows configured workflows"
    else
        echo "  [FAIL] workflow list missing expected workflow"
        echo "  Output: $list_out"
        exit 1
    fi
else
    echo "  [FAIL] workflow list should succeed"
    echo "  Output: $list_out"
    exit 1
fi

if validate_out=$(./bin/multipowers workflow validate 2>&1); then
    if echo "$validate_out" | grep -q "\[WORKFLOW-VALIDATE\] PASS"; then
        echo "  [PASS] workflow validate succeeds"
    else
        echo "  [FAIL] workflow validate missing pass output"
        echo "  Output: $validate_out"
        exit 1
    fi
else
    echo "  [FAIL] workflow validate should succeed"
    echo "  Output: $validate_out"
    exit 1
fi

# Test 3: workflow run dry-run works and resolves nodes.
echo "[TEST 3] workflow run dry-run"
out=$(./bin/multipowers workflow run brainstorming --task "Explore options" --allow-untracked --dry-run --json)
if python3 - "$out" <<'PY'
import json
import sys
payload = json.loads(sys.argv[1])
if payload.get("workflow") != "brainstorming":
    raise SystemExit(1)
nodes = payload.get("nodes", [])
if not nodes:
    raise SystemExit(1)
if nodes[0].get("status") != "skipped":
    raise SystemExit(1)
print("ok")
PY
then
    echo "  [PASS] workflow dry-run returns node plan"
else
    echo "  [FAIL] workflow dry-run payload invalid"
    echo "  Output: $out"
    exit 1
fi

# Test 4: workflow run missing task has actionable error.
echo "[TEST 4] workflow error message is actionable"
set +e
out=$(./bin/multipowers workflow run brainstorming 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] workflow run without --task should fail"
    exit 1
fi
if echo "$out" | grep -q "requires --task" && echo "$out" | grep -q "Usage:"; then
    echo "  [PASS] workflow error message is actionable"
else
    echo "  [FAIL] workflow error output not actionable"
    echo "  Output: $out"
    exit 1
fi

echo ""
echo "Multipowers route/workflow tests PASSED"
