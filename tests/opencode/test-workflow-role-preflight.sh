#!/usr/bin/env bash
# Workflow role preflight validation tests
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

echo "Testing workflow role preflight..."

workspace=$(create_workspace)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cd "$workspace"
./bin/multipowers init --repair >/tmp/test_role_preflight_init_out.txt 2>/tmp/test_role_preflight_init_err.txt

# prepare active track for execution command constraints
./bin/multipowers track new "role-preflight" >/tmp/test_role_preflight_new_out.txt
track_basename=$(basename conductor/tracks/track-*-role-preflight.md .md)
./bin/multipowers track start "$track_basename" >/tmp/test_role_preflight_start_out.txt

mkdir -p conductor/config
cat > conductor/config/workflows.json <<'JSON_EOF'
{
  "workflows": {
    "broken-flow": {
      "default_role": "coder",
      "nodes": [
        {
          "id": "step-1",
          "role": "not-a-real-role",
          "prompt_template": "Do thing: {task}"
        }
      ]
    }
  }
}
JSON_EOF

# Test 1: workflow run fails before node execution for unknown role.
echo "[TEST 1] workflow run preflight catches unknown role"
set +e
output=$(./bin/multipowers workflow run broken-flow --task "do work" --json 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected workflow run to fail on unknown role"
    exit 1
fi
if echo "$output" | grep -q "unknown role" && echo "$output" | grep -q "broken-flow" && echo "$output" | grep -q "step-1"; then
    echo "  [PASS] run preflight reports workflow + node + role"
else
    echo "  [FAIL] run preflight diagnostics incomplete"
    echo "  Output: $output"
    exit 1
fi

# Test 2: workflow validate catches unknown role too.
echo "[TEST 2] workflow validate catches unknown role"
set +e
output=$(./bin/multipowers workflow validate 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected workflow validate to fail"
    exit 1
fi
if echo "$output" | grep -q "unknown role" && echo "$output" | grep -q "broken-flow"; then
    echo "  [PASS] validate preflight reports unknown role"
else
    echo "  [FAIL] validate output missing unknown role diagnostics"
    echo "  Output: $output"
    exit 1
fi

# Test 3: valid workflow still passes.
echo "[TEST 3] valid workflow still passes"
cat > conductor/config/workflows.json <<'JSON_EOF'
{
  "workflows": {
    "ok-flow": {
      "default_role": "coder",
      "nodes": [
        {
          "id": "step-1",
          "prompt_template": "Do thing: {task}"
        }
      ]
    }
  }
}
JSON_EOF

if ./bin/multipowers workflow validate >/tmp/test_role_preflight_validate_ok_out.txt 2>/tmp/test_role_preflight_validate_ok_err.txt; then
    echo "  [PASS] valid workflow validation succeeds"
else
    echo "  [FAIL] valid workflow should validate"
    cat /tmp/test_role_preflight_validate_ok_err.txt
    exit 1
fi

echo ""
echo "Workflow role preflight tests PASSED"
