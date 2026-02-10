#!/usr/bin/env bash
# Track complete governance gate tests
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

echo "Testing track complete governance gate..."

workspace=$(create_workspace)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cd "$workspace"
./bin/multipowers init --repair >/tmp/test_track_gate_init_out.txt 2>/tmp/test_track_gate_init_err.txt

git init >/tmp/test_track_gate_git_init_out.txt 2>/tmp/test_track_gate_git_init_err.txt
git config user.email "test@example.com"
git config user.name "Test User"
git add .
git commit -m "init" >/tmp/test_track_gate_git_commit_out.txt 2>/tmp/test_track_gate_git_commit_err.txt

./bin/multipowers track new "governance-gate" >/tmp/test_track_gate_new_out.txt
track_basename=$(basename conductor/tracks/track-*-governance-gate.md .md)
./bin/multipowers track start "$track_basename" >/tmp/test_track_gate_start_out.txt

# make a code-only change to trigger docs-sync governance failure
echo "# governance gate test" >> scripts/route_task.py

# Test 1: complete blocked by governance by default.
echo "[TEST 1] track complete blocks on governance failure"
set +e
output=$(./bin/multipowers track complete "$track_basename" 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] track complete should fail without governance evidence"
    exit 1
fi
if echo "$output" | grep -q "governance" && echo "$output" | grep -Eq "failed|FAIL"; then
    echo "  [PASS] governance gate blocks completion"
else
    echo "  [FAIL] missing governance gate diagnostics"
    echo "  Output: $output"
    exit 1
fi

# Test 2: explicit bypass works and remains auditable.
echo "[TEST 2] skip-governance bypass is explicit and auditable"
if ./bin/multipowers track complete "$track_basename" --skip-governance >/tmp/test_track_gate_complete_out.txt 2>/tmp/test_track_gate_complete_err.txt; then
    track_file="conductor/tracks/${track_basename}.md"
    if grep -q '^\*\*Status:\*\* Completed' "$track_file" && grep -q '^\*\*Governance:\*\* skipped (explicit bypass)' "$track_file"; then
        echo "  [PASS] bypass completes with governance metadata"
    else
        echo "  [FAIL] track metadata missing governance bypass record"
        cat "$track_file"
        exit 1
    fi
else
    echo "  [FAIL] skip-governance should allow completion"
    cat /tmp/test_track_gate_complete_err.txt
    exit 1
fi

echo ""
echo "Track complete governance gate tests PASSED"
