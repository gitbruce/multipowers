#!/usr/bin/env bash
# Routing/workflow observability tests
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

echo "Testing routing/workflow observability..."

workspace=$(create_workspace)
stub_bin=$(mktemp -d)
cleanup() {
    rm -rf "$workspace" "$stub_bin"
}
trap cleanup EXIT

cd "$workspace"
./bin/multipowers init --repair >/tmp/test_obs_repair_out.txt 2>/tmp/test_obs_repair_err.txt

# Stub model CLIs used by ask-role during workflow execution.
cat > "$stub_bin/codex" <<'STUB'
#!/usr/bin/env python3
print("stub-codex")
STUB
cat > "$stub_bin/gemini" <<'STUB'
#!/usr/bin/env python3
print("stub-gemini")
STUB
chmod +x "$stub_bin/codex" "$stub_bin/gemini"

request_id="req-obsv-001"
track_id="track-obsv-001"

PATH="$stub_bin:$PATH" ./bin/multipowers route --task "Migrate architecture for auth" --risk-hint high --request-id "$request_id" --track-id "$track_id" --json >/tmp/test_obs_route_out.json
PATH="$stub_bin:$PATH" ./bin/multipowers workflow run subagent-driven-development --task "Implement auth migration" --request-id "$request_id" --track-id "$track_id" --json >/tmp/test_obs_workflow_out.json

log_file="outputs/runs/$(date +%Y-%m-%d).jsonl"
if [ ! -f "$log_file" ]; then
    echo "  [FAIL] Expected observability log file: $log_file"
    exit 1
fi

if python3 - "$log_file" "$request_id" "$track_id" <<'PY'
import json
import sys

log_file = sys.argv[1]
request_id = sys.argv[2]
track_id = sys.argv[3]

entries = []
with open(log_file, "r", encoding="utf-8") as handle:
    for line in handle:
        line = line.strip()
        if not line:
            continue
        payload = json.loads(line)
        if payload.get("request_id") == request_id:
            entries.append(payload)

if not entries:
    raise SystemExit("no entries for request_id")

required_events = ["lane_selected", "workflow_started", "workflow_node_executed", "workflow_finished"]
events = [entry.get("event") for entry in entries]
for event in required_events:
    if event not in events:
        raise SystemExit(f"missing event: {event}")

lane_index = events.index("lane_selected")
workflow_start_index = events.index("workflow_started")
workflow_finish_index = len(events) - 1 - events[::-1].index("workflow_finished")
if not (lane_index < workflow_start_index < workflow_finish_index):
    raise SystemExit(f"invalid event order: {events}")

lane_entry = next(entry for entry in entries if entry.get("event") == "lane_selected")
if lane_entry.get("lane") != "standard":
    raise SystemExit(f"expected standard lane, got {lane_entry.get('lane')}")
if lane_entry.get("track_id") != track_id:
    raise SystemExit("lane_selected missing track_id")
if not lane_entry.get("reason"):
    raise SystemExit("lane_selected missing reason")

node_entries = [entry for entry in entries if entry.get("event") == "workflow_node_executed"]
if len(node_entries) < 2:
    raise SystemExit("expected at least two workflow_node_executed events")
for node_entry in node_entries:
    if not node_entry.get("node"):
        raise SystemExit("workflow_node_executed missing node")
    if node_entry.get("workflow") != "subagent-driven-development":
        raise SystemExit("workflow_node_executed missing workflow")

print("ok")
PY
then
    echo "  [PASS] Structured observability events are complete and ordered"
else
    echo "  [FAIL] Observability schema/sequence validation failed"
    exit 1
fi

echo ""
echo "Routing observability tests PASSED"
