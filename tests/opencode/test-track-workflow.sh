#!/usr/bin/env bash
# Track workflow tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "Testing track workflow commands..."

workspace=$(mktemp -d)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cp -R "$REPO_ROOT/bin" "$workspace/"
cp -R "$REPO_ROOT/templates" "$workspace/"
chmod +x "$workspace/bin/multipowers"

cd "$workspace"
mkdir -p conductor/tracks

# Test 1: Create new track

echo "[TEST 1] Create new track"
track_name="test-feature-$(date +%s)"
./bin/multipowers track new "$track_name" >/dev/null

shopt -s nullglob
matches=(conductor/tracks/track-*-${track_name}.md)
shopt -u nullglob

if [ ${#matches[@]} -ne 1 ]; then
    echo "  [FAIL] Track file not created"
    exit 1
fi
track_file=${matches[0]}

if grep -q '^\*\*Status:\*\* Proposed' "$track_file"; then
    echo "  [PASS] Status set to Proposed"
else
    echo "  [FAIL] Status not set correctly"
    exit 1
fi
if grep -q '^\*\*Updated At:\*\*' "$track_file"; then
    echo "  [PASS] Updated At field exists"
else
    echo "  [FAIL] Updated At field missing"
    exit 1
fi
if grep -q '^\*\*Owner:\*\*' "$track_file"; then
    echo "  [PASS] Owner field exists"
else
    echo "  [FAIL] Owner field missing"
    exit 1
fi

# Test 2: Slug fallback works for non-ASCII names

echo "[TEST 2] Non-ASCII track name fallback"
non_ascii_name="中文 功能"
./bin/multipowers track new "$non_ascii_name" >/dev/null

shopt -s nullglob
non_ascii_matches=(conductor/tracks/track-*-task-*.md)
shopt -u nullglob
if [ ${#non_ascii_matches[@]} -lt 1 ]; then
    echo "  [FAIL] Expected fallback task slug for non-ASCII name"
    exit 1
fi
echo "  [PASS] Non-ASCII name generated fallback slug"

# Test 3: List tracks

echo "[TEST 3] List tracks"
echo "# Test Track 1" > conductor/tracks/track-001-test1.md
echo "# Test Track 2" > conductor/tracks/track-002-test2.md
output=$(./bin/multipowers track list 2>&1)

if echo "$output" | grep -q "track-001-test1.md" && echo "$output" | grep -q "track-002-test2.md"; then
    echo "  [PASS] Track list includes expected files"
else
    echo "  [FAIL] Track list output invalid"
    echo "  Output: $output"
    exit 1
fi

# Test 4: Ambiguous track selector shows actionable error

echo "[TEST 4] Ambiguous track selector"
cat > conductor/tracks/track-101-api.md <<'TRACK_EOF'
# Track: API 101
**Status:** Proposed
**Updated At:** 2026-01-01
**Owner:** Test
TRACK_EOF
cat > conductor/tracks/track-102-api.md <<'TRACK_EOF'
# Track: API 102
**Status:** Proposed
**Updated At:** 2026-01-01
**Owner:** Test
TRACK_EOF

set +e
output=$(./bin/multipowers track start "api" 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected ambiguous track selector to fail"
    exit 1
fi
if echo "$output" | grep -q "ambiguous" && echo "$output" | grep -q "Use full filename"; then
    echo "  [PASS] Ambiguous selector gives actionable guidance"
else
    echo "  [FAIL] Missing actionable ambiguous selector message"
    echo "  Output: $output"
    exit 1
fi

# Test 5: Start track (Proposed -> In Progress)

echo "[TEST 5] Start track transition"
cat > conductor/tracks/track-003-test-start.md <<'TRACK_EOF'
# Track: Start Test
**Status:** Proposed
**Updated At:** 2026-01-01
**Owner:** Test
**Goal:** Test status transition
TRACK_EOF

./bin/multipowers track start "test-start" >/dev/null
if grep -q '^\*\*Status:\*\* In Progress' conductor/tracks/track-003-test-start.md; then
    echo "  [PASS] Status changed to In Progress"
else
    echo "  [FAIL] Start transition failed"
    exit 1
fi

# Test 6: Prevent restart of completed track

echo "[TEST 6] Prevent restarting completed track"
cat > conductor/tracks/track-004-test-completed.md <<'TRACK_EOF'
# Track: Completed Test
**Status:** Completed
**Updated At:** 2026-01-01
**Owner:** Test
TRACK_EOF

set +e
output=$(./bin/multipowers track start "test-completed" 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Should not restart completed track"
    exit 1
fi
if echo "$output" | grep -q "Cannot restart a completed track"; then
    echo "  [PASS] Restart prevented"
else
    echo "  [FAIL] Missing restart prevention message"
    exit 1
fi

# Test 7: Complete track (In Progress -> Completed)

echo "[TEST 7] Complete track transition"
cat > conductor/tracks/track-005-test-complete.md <<'TRACK_EOF'
# Track: Complete Test
**Status:** In Progress
**Updated At:** 2026-01-01
**Owner:** Test
TRACK_EOF

./bin/multipowers track complete "test-complete" >/dev/null
if grep -q '^\*\*Status:\*\* Completed' conductor/tracks/track-005-test-complete.md; then
    echo "  [PASS] Status changed to Completed"
else
    echo "  [FAIL] Complete transition failed"
    exit 1
fi

# Test 8: Prevent completing non in-progress track

echo "[TEST 8] Prevent invalid complete transition"
cat > conductor/tracks/track-006-test-proposed.md <<'TRACK_EOF'
# Track: Proposed Test
**Status:** Proposed
**Updated At:** 2026-01-01
**Owner:** Test
TRACK_EOF

set +e
output=$(./bin/multipowers track complete "test-proposed" 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Should not complete proposed track"
    exit 1
fi
if echo "$output" | grep -q "Can only complete tracks in progress"; then
    echo "  [PASS] Invalid completion prevented"
else
    echo "  [FAIL] Missing invalid completion message"
    exit 1
fi

# Test 9: Status output

echo "[TEST 9] Show track status"
cat > conductor/tracks/track-007-test-status.md <<'TRACK_EOF'
# Track: Status Test
**Status:** In Progress
**Updated At:** 2026-02-09
**Owner:** Test User
**Goal:** Test status display
TRACK_EOF

output=$(./bin/multipowers track status "test-status" 2>&1)
if echo "$output" | grep -q '^\*\*Status:\*\* In Progress' && echo "$output" | grep -q '^\*\*Owner:\*\* Test User'; then
    echo "  [PASS] Status output is correct"
else
    echo "  [FAIL] Status output mismatch"
    echo "  Output: $output"
    exit 1
fi

echo ""
echo "All track workflow tests PASSED"
