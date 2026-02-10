#!/usr/bin/env bash
# Routing lane selection tests
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
    chmod +x "$workspace/bin/multipowers"
    echo "$workspace"
}

echo "Testing routing lane command..."

workspace=$(create_workspace)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cd "$workspace"

# Test 1: Fast lane for bounded low-risk task

echo "[TEST 1] Low-risk task routes to fast lane"
out=$(./bin/multipowers route --task "Rename one variable in a single file" --json)
if python3 - "$out" <<'PY'
import json
import sys
payload = json.loads(sys.argv[1])
if payload.get("lane") != "fast":
    raise SystemExit(1)
if not payload.get("reason"):
    raise SystemExit(1)
print("ok")
PY
then
    echo "  [PASS] Fast lane selected"
else
    echo "  [FAIL] Expected fast lane"
    echo "  Output: $out"
    exit 1
fi

# Test 2: High-risk task routes to standard lane

echo "[TEST 2] High-risk task routes to standard lane"
out=$(./bin/multipowers route --task "Perform architecture migration across services" --risk-hint high --json)
if python3 - "$out" <<'PY'
import json
import sys
payload = json.loads(sys.argv[1])
if payload.get("lane") != "standard":
    raise SystemExit(1)
if payload.get("suggested_workflow") not in {"writing-plans", "subagent-driven-development"}:
    raise SystemExit(1)
print("ok")
PY
then
    echo "  [PASS] Standard lane selected"
else
    echo "  [FAIL] Expected standard lane"
    echo "  Output: $out"
    exit 1
fi

# Test 3: Deterministic behavior for same input

echo "[TEST 3] Routing is deterministic"
out1=$(./bin/multipowers route --task "Review architecture and workflow migration" --risk-hint medium --json)
out2=$(./bin/multipowers route --task "Review architecture and workflow migration" --risk-hint medium --json)
if python3 - "$out1" "$out2" <<'PY'
import json
import sys
p1 = json.loads(sys.argv[1])
p2 = json.loads(sys.argv[2])
keys = ["lane", "reason", "suggested_workflow", "suggested_role"]
for key in keys:
    if p1.get(key) != p2.get(key):
        raise SystemExit(1)
print("ok")
PY
then
    echo "  [PASS] Deterministic route output"
else
    echo "  [FAIL] Routing output is not deterministic"
    echo "  Out1: $out1"
    echo "  Out2: $out2"
    exit 1
fi

# Test 4: Force lane override

echo "[TEST 4] Force lane override works"
out=$(./bin/multipowers route --task "Major migration" --risk-hint high --force-lane fast --json)
if python3 - "$out" <<'PY'
import json
import sys
payload = json.loads(sys.argv[1])
if payload.get("lane") != "fast":
    raise SystemExit(1)
if "force_lane" not in payload:
    raise SystemExit(1)
print("ok")
PY
then
    echo "  [PASS] Force lane override honored"
else
    echo "  [FAIL] Force lane override not honored"
    echo "  Output: $out"
    exit 1
fi

echo ""
echo "Routing lane tests PASSED"
