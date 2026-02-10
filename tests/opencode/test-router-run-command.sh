#!/usr/bin/env bash
# Router one-command execution tests
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

echo "Testing router one-command run..."

workspace=$(create_workspace)
stub_bin=$(mktemp -d)
cleanup() {
    rm -rf "$workspace" "$stub_bin"
}
trap cleanup EXIT

cd "$workspace"
./bin/multipowers init --repair >/tmp/test_router_run_init_out.txt 2>/tmp/test_router_run_init_err.txt

cat > "$stub_bin/codex" <<'STUB'
#!/usr/bin/env python3
print("stub-codex")
STUB
cat > "$stub_bin/gemini" <<'STUB'
#!/usr/bin/env python3
print("stub-gemini")
STUB
cat > "$stub_bin/claude" <<'STUB'
#!/usr/bin/env python3
print("stub-claude")
STUB
chmod +x "$stub_bin/codex" "$stub_bin/gemini" "$stub_bin/claude"

# Test 1: fast lane route+execute in one command.
echo "[TEST 1] run command executes fast lane end-to-end"
request_fast="req-router-fast-001"
if out=$(PATH="$stub_bin:$PATH" ./bin/multipowers run --task "Fix typo in README" --request-id "$request_fast" --allow-untracked --json 2>/tmp/test_router_run_err1.txt); then
    if python3 - "$out" <<'PY'
import json
import sys
payload = json.loads(sys.argv[1])
if payload.get("lane") != "fast":
    raise SystemExit("expected fast lane")
exec_payload = payload.get("execution", {})
if exec_payload.get("status") != "ok":
    raise SystemExit(f"expected fast execution ok, got {exec_payload}")
if payload.get("request_id") != "req-router-fast-001":
    raise SystemExit("request_id mismatch")
print("ok")
PY
    then
        echo "  [PASS] fast lane run works"
    else
        echo "  [FAIL] fast lane payload invalid"
        cat /tmp/test_router_run_err1.txt
        exit 1
    fi
else
    echo "  [FAIL] fast lane run should succeed"
    cat /tmp/test_router_run_err1.txt
    exit 1
fi

# Test 2: standard lane route+workflow in one command.
echo "[TEST 2] run command executes standard lane workflow"
request_std="req-router-std-001"
if out=$(PATH="$stub_bin:$PATH" ./bin/multipowers run --task "Migrate architecture for auth layer" --risk-hint high --request-id "$request_std" --allow-untracked --json 2>/tmp/test_router_run_err2.txt); then
    if python3 - "$out" <<'PY'
import json
import sys
payload = json.loads(sys.argv[1])
if payload.get("lane") != "standard":
    raise SystemExit("expected standard lane")
exec_payload = payload.get("execution", {})
if exec_payload.get("status") != "ok":
    raise SystemExit(f"expected standard execution ok, got {exec_payload}")
workflow_payload = exec_payload.get("workflow_output", {})
if not workflow_payload.get("nodes"):
    raise SystemExit("workflow output missing nodes")
print("ok")
PY
    then
        echo "  [PASS] standard lane run works"
    else
        echo "  [FAIL] standard lane payload invalid"
        cat /tmp/test_router_run_err2.txt
        exit 1
    fi
else
    echo "  [FAIL] standard lane run should succeed"
    cat /tmp/test_router_run_err2.txt
    exit 1
fi

# Test 3: request_id is propagated through events.
echo "[TEST 3] run keeps request_id across lifecycle events"
log_file="outputs/runs/$(date +%Y-%m-%d).jsonl"
if python3 - "$log_file" "$request_fast" "$request_std" <<'PY'
import json
import sys
from collections import defaultdict

log_file = sys.argv[1]
requests = {sys.argv[2], sys.argv[3]}

by_req = defaultdict(list)
with open(log_file, "r", encoding="utf-8") as handle:
    for line in handle:
        line = line.strip()
        if not line:
            continue
        payload = json.loads(line)
        request_id = payload.get("request_id")
        if request_id in requests:
            by_req[request_id].append(payload.get("event"))

if sys.argv[2] not in by_req or "fast_lane_finished" not in by_req[sys.argv[2]]:
    raise SystemExit("missing fast lane lifecycle events")
if sys.argv[3] not in by_req or "workflow_finished" not in by_req[sys.argv[3]]:
    raise SystemExit("missing standard workflow lifecycle events")

print("ok")
PY
then
    echo "  [PASS] request_id lifecycle propagation verified"
else
    echo "  [FAIL] request_id lifecycle propagation failed"
    exit 1
fi

echo ""
echo "Router run command tests PASSED"
