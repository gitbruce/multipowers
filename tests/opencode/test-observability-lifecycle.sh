#!/usr/bin/env bash
# Observability lifecycle coverage tests
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

echo "Testing observability lifecycle events..."

workspace=$(create_workspace)
stub_bin=$(mktemp -d)
cleanup() {
    rm -rf "$workspace" "$stub_bin"
}
trap cleanup EXIT

cd "$workspace"
./bin/multipowers init --repair >/tmp/test_obs_lifecycle_init_out.txt 2>/tmp/test_obs_lifecycle_init_err.txt

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

./bin/multipowers track new "obs-lifecycle" >/tmp/test_obs_lifecycle_new_out.txt
track_basename=$(basename conductor/tracks/track-*-obs-lifecycle.md .md)
./bin/multipowers track start "$track_basename" >/tmp/test_obs_lifecycle_start_out.txt

request_run="req-obslife-run-001"
request_complete="req-obslife-complete-001"

PATH="$stub_bin:$PATH" ./bin/multipowers run --task "Fix typo in README" --request-id "$request_run" --json >/tmp/test_obs_lifecycle_run.json

# no code changes => governance script should pass, but lifecycle events must still emit
PATH="$stub_bin:$PATH" ./bin/multipowers track complete "$track_basename" --request-id "$request_complete" >/tmp/test_obs_lifecycle_complete_out.txt

log_file="outputs/runs/$(date +%Y-%m-%d).jsonl"
if [ ! -f "$log_file" ]; then
    echo "  [FAIL] Expected structured log file: $log_file"
    exit 1
fi

if python3 - "$log_file" "$request_run" "$request_complete" <<'PY'
import json
import sys

log_file = sys.argv[1]
request_run = sys.argv[2]
request_complete = sys.argv[3]

run_events = []
complete_events = []
with open(log_file, 'r', encoding='utf-8') as handle:
    for line in handle:
        line = line.strip()
        if not line:
            continue
        payload = json.loads(line)
        if payload.get('request_id') == request_run:
            run_events.append(payload)
        if payload.get('request_id') == request_complete:
            complete_events.append(payload)

run_names = [e.get('event') for e in run_events]
for required in ['lane_selected', 'fast_lane_dispatched', 'fast_lane_finished']:
    if required not in run_names:
        raise SystemExit(f"missing run lifecycle event: {required}")
if run_names.index('fast_lane_dispatched') > run_names.index('fast_lane_finished'):
    raise SystemExit(f"invalid fast lane event order: {run_names}")

complete_names = [e.get('event') for e in complete_events]
for required in ['governance_started', 'governance_finished']:
    if required not in complete_names:
        raise SystemExit(f"missing governance lifecycle event: {required}")
if complete_names.index('governance_started') > complete_names.index('governance_finished'):
    raise SystemExit(f"invalid governance event order: {complete_names}")

print('ok')
PY
then
    echo "  [PASS] lifecycle observability events complete and ordered"
else
    echo "  [FAIL] lifecycle observability validation failed"
    exit 1
fi

echo ""
echo "Observability lifecycle tests PASSED"
