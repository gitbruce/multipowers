#!/usr/bin/env bash
# Full routing flow test: standard + fast lanes with observability
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

echo "Testing full routing flow..."

workspace=$(create_workspace)
stub_bin=$(mktemp -d)
cleanup() {
    rm -rf "$workspace" "$stub_bin"
}
trap cleanup EXIT

cd "$workspace"
./bin/multipowers init --repair >/tmp/test_full_route_repair_out.txt 2>/tmp/test_full_route_repair_err.txt

# Stub external CLIs used by connectors.
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

request_std="req-full-standard-001"
request_fast="req-full-fast-001"
track_id="track-full-001"

# Standard lane route + workflow execution path.
PATH="$stub_bin:$PATH" ./bin/multipowers route --task "Migrate architecture across services" --risk-hint high --request-id "$request_std" --track-id "$track_id" --json >/tmp/test_full_route_std.json
PATH="$stub_bin:$PATH" ./bin/multipowers workflow run subagent-driven-development --task "Implement auth migration" --request-id "$request_std" --track-id "$track_id" --json >/tmp/test_full_route_workflow.json

# Fast lane route path.
PATH="$stub_bin:$PATH" ./bin/multipowers route --task "Fix typo in README" --request-id "$request_fast" --track-id "$track_id" --json >/tmp/test_full_route_fast.json

log_file="outputs/runs/$(date +%Y-%m-%d).jsonl"
if [ ! -f "$log_file" ]; then
    echo "  [FAIL] Expected structured log file: $log_file"
    exit 1
fi

if python3 - "$log_file" "$request_std" "$request_fast" "$track_id" /tmp/test_full_route_std.json /tmp/test_full_route_fast.json <<'PY'
import json
import sys

log_file = sys.argv[1]
request_std = sys.argv[2]
request_fast = sys.argv[3]
track_id = sys.argv[4]
std_json_path = sys.argv[5]
fast_json_path = sys.argv[6]

with open(std_json_path, "r", encoding="utf-8") as handle:
    std_route = json.load(handle)
with open(fast_json_path, "r", encoding="utf-8") as handle:
    fast_route = json.load(handle)

if std_route.get("lane") != "standard":
    raise SystemExit(f"expected standard lane, got {std_route.get('lane')}")
if fast_route.get("lane") != "fast":
    raise SystemExit(f"expected fast lane, got {fast_route.get('lane')}")

entries = []
with open(log_file, "r", encoding="utf-8") as handle:
    for line in handle:
        line = line.strip()
        if not line:
            continue
        payload = json.loads(line)
        if payload.get("request_id") in {request_std, request_fast}:
            entries.append(payload)

std_entries = [e for e in entries if e.get("request_id") == request_std]
fast_entries = [e for e in entries if e.get("request_id") == request_fast]

if not std_entries:
    raise SystemExit("missing standard request entries")
if not fast_entries:
    raise SystemExit("missing fast request entries")

std_events = [e.get("event") for e in std_entries if e.get("event")]
for required in ["lane_selected", "workflow_started", "workflow_node_executed", "workflow_finished"]:
    if required not in std_events:
        raise SystemExit(f"missing standard-lane event: {required}")

fast_events = [e.get("event") for e in fast_entries if e.get("event")]
if "lane_selected" not in fast_events:
    raise SystemExit("missing fast-lane lane_selected event")
if "workflow_started" in fast_events:
    raise SystemExit("fast lane should not emit workflow_started")

lane_std = next(e for e in std_entries if e.get("event") == "lane_selected")
lane_fast = next(e for e in fast_entries if e.get("event") == "lane_selected")
if lane_std.get("track_id") != track_id or lane_fast.get("track_id") != track_id:
    raise SystemExit("track_id not preserved in lane events")

print("ok")
PY
then
    echo "  [PASS] Full flow covers standard+fast lanes with event logs"
else
    echo "  [FAIL] Full flow validation failed"
    exit 1
fi

echo ""
echo "Full routing flow tests PASSED"
