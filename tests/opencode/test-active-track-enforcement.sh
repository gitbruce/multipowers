#!/usr/bin/env bash
# Active track enforcement tests
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

echo "Testing active track enforcement..."

workspace=$(create_workspace)
stub_bin=$(mktemp -d)
cleanup() {
    rm -rf "$workspace" "$stub_bin"
}
trap cleanup EXIT

cd "$workspace"
./bin/multipowers init --repair >/tmp/test_track_enforce_init_out.txt 2>/tmp/test_track_enforce_init_err.txt

cat > "$stub_bin/codex" <<'STUB'
#!/usr/bin/env python3
print("stub-codex")
STUB
cat > "$stub_bin/gemini" <<'STUB'
#!/usr/bin/env python3
print("stub-gemini")
STUB
chmod +x "$stub_bin/codex" "$stub_bin/gemini"

# Test 1: workflow run blocked without active track.
echo "[TEST 1] workflow run requires active track"
set +e
output=$(PATH="$stub_bin:$PATH" ./bin/multipowers workflow run brainstorming --task "draft plan" --json 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected workflow run to fail without active track"
    exit 1
fi
if echo "$output" | grep -q "active track" && echo "$output" | grep -q "track start"; then
    echo "  [PASS] workflow run blocked with actionable message"
else
    echo "  [FAIL] Missing actionable active track message"
    echo "  Output: $output"
    exit 1
fi

# Test 2: run command blocked without active track.
echo "[TEST 2] run requires active track"
set +e
output=$(PATH="$stub_bin:$PATH" ./bin/multipowers run --task "Fix typo" --json 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected run to fail without active track"
    exit 1
fi
if echo "$output" | grep -q "active track"; then
    echo "  [PASS] run command blocked without active track"
else
    echo "  [FAIL] Missing active track diagnostic"
    echo "  Output: $output"
    exit 1
fi

# Test 3: start track writes marker and allows execution.
echo "[TEST 3] track start sets active marker"
./bin/multipowers track new "active-track-enforcement" >/tmp/test_track_enforce_new_out.txt
track_basename=$(basename conductor/tracks/track-*-active-track-enforcement.md .md)
./bin/multipowers track start "$track_basename" >/tmp/test_track_enforce_start_out.txt

if [ ! -f conductor/.active_track ]; then
    echo "  [FAIL] active track marker not created"
    exit 1
fi
if ! grep -q "$track_basename" conductor/.active_track; then
    echo "  [FAIL] active track marker does not store track id"
    exit 1
fi

if out=$(PATH="$stub_bin:$PATH" ./bin/multipowers workflow run brainstorming --task "draft plan" --json 2>/tmp/test_track_enforce_run_err.txt); then
    if python3 - "$out" "$track_basename" <<'PY'
import json
import sys
payload = json.loads(sys.argv[1])
if payload.get("track_id") != sys.argv[2]:
    raise SystemExit(f"track_id mismatch: {payload.get('track_id')}")
print("ok")
PY
    then
        echo "  [PASS] active track marker enables execution with track_id"
    else
        echo "  [FAIL] workflow output missing active track id"
        cat /tmp/test_track_enforce_run_err.txt
        exit 1
    fi
else
    echo "  [FAIL] workflow run should succeed with active track"
    cat /tmp/test_track_enforce_run_err.txt
    exit 1
fi

# Test 4: explicit allow-untracked bypass works.
echo "[TEST 4] allow-untracked bypass works"
rm -f conductor/.active_track
if PATH="$stub_bin:$PATH" ./bin/multipowers workflow run brainstorming --task "draft plan" --allow-untracked --json >/tmp/test_track_enforce_allow_out.json 2>/tmp/test_track_enforce_allow_err.txt; then
    echo "  [PASS] allow-untracked bypass works"
else
    echo "  [FAIL] allow-untracked should bypass enforcement"
    cat /tmp/test_track_enforce_allow_err.txt
    exit 1
fi

echo ""
echo "Active track enforcement tests PASSED"
