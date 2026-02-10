#!/usr/bin/env bash
# Workflow engine tests
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

echo "Testing workflow engine..."

workspace=$(create_workspace)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cd "$workspace"

# Use a stub ask-role to isolate workflow execution behavior.
cat > bin/ask-role <<'ASK_EOF'
#!/usr/bin/env bash
set -euo pipefail
role=${1:-}
prompt=${2:-}
log_file=${WORKFLOW_TEST_LOG:-}
if [ -n "$log_file" ]; then
  echo "${role}|${prompt}" >> "$log_file"
fi
if [ "${WORKFLOW_FAIL_ROLE:-}" = "$role" ]; then
  echo "stub-fail-$role" >&2
  exit 23
fi
echo "stub-ok-$role"
ASK_EOF
chmod +x bin/ask-role

# Test 1: workflows are configured and node-level role override executes.
echo "[TEST 1] Node-level role override works"
export WORKFLOW_TEST_LOG="$workspace/workflow.log"
: > "$WORKFLOW_TEST_LOG"

if out=$(./bin/multipowers workflow run subagent-driven-development --task "Implement API validation" --json 2>/tmp/test_workflow_engine_err.txt); then
    if python3 - "$out" "$WORKFLOW_TEST_LOG" <<'PY'
import json
import sys

payload = json.loads(sys.argv[1])
with open(sys.argv[2], "r", encoding="utf-8") as handle:
    lines = [line.strip() for line in handle if line.strip()]

if len(payload.get("nodes", [])) < 2:
    raise SystemExit("expected at least two nodes")
roles = [node.get("role") for node in payload["nodes"]]
if roles[:2] != ["coder", "architect"]:
    raise SystemExit(f"unexpected node roles: {roles}")

logged_roles = [line.split("|", 1)[0] for line in lines]
if logged_roles[:2] != ["coder", "architect"]:
    raise SystemExit(f"unexpected execution order: {logged_roles}")

print("ok")
PY
    then
        echo "  [PASS] Workflow uses configured roles and overrides"
    else
        echo "  [FAIL] Workflow output/roles mismatch"
        cat /tmp/test_workflow_engine_err.txt
        exit 1
    fi
else
    echo "  [FAIL] Workflow run should succeed"
    cat /tmp/test_workflow_engine_err.txt
    exit 1
fi

# Test 2: project workflow config overrides default.
echo "[TEST 2] Project workflow config takes precedence"
mkdir -p conductor/config
cat > conductor/config/workflows.json <<'JSON_EOF'
{
  "workflows": {
    "custom-review": {
      "default_role": "librarian",
      "nodes": [
        {
          "id": "research",
          "prompt_template": "Research task: {task}"
        }
      ]
    }
  }
}
JSON_EOF

: > "$WORKFLOW_TEST_LOG"
if out=$(./bin/multipowers workflow run custom-review --task "find docs" --json 2>/tmp/test_workflow_engine_err.txt); then
    if python3 - "$out" "$WORKFLOW_TEST_LOG" <<'PY'
import json
import sys

payload = json.loads(sys.argv[1])
if payload.get("workflow") != "custom-review":
    raise SystemExit("wrong workflow name")
nodes = payload.get("nodes", [])
if len(nodes) != 1 or nodes[0].get("role") != "librarian":
    raise SystemExit(f"unexpected nodes: {nodes}")

with open(sys.argv[2], "r", encoding="utf-8") as handle:
    lines = [line.strip() for line in handle if line.strip()]
if not lines or not lines[0].startswith("librarian|"):
    raise SystemExit("project override workflow did not execute librarian role")

print("ok")
PY
    then
        echo "  [PASS] Project workflow override is honored"
    else
        echo "  [FAIL] Project workflow override validation failed"
        cat /tmp/test_workflow_engine_err.txt
        exit 1
    fi
else
    echo "  [FAIL] Custom workflow should succeed"
    cat /tmp/test_workflow_engine_err.txt
    exit 1
fi

# Test 3: failure path includes actionable diagnostics.
rm -f conductor/config/workflows.json
echo "[TEST 3] Node failure has clear diagnostics"
export WORKFLOW_FAIL_ROLE="architect"
set +e
output=$(./bin/multipowers workflow run subagent-driven-development --task "trigger fail" 2>&1)
status=$?
set -e
unset WORKFLOW_FAIL_ROLE

if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected workflow failure when architect node fails"
    exit 1
fi
if echo "$output" | grep -q "Node 'review' failed (role=architect, exit_code=23)"; then
    echo "  [PASS] Failure diagnostics are actionable"
else
    echo "  [FAIL] Missing actionable failure diagnostics"
    echo "  Output: $output"
    exit 1
fi

echo ""
echo "Workflow engine tests PASSED"
