#!/usr/bin/env bash
# Update command behavior tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "Testing update command..."

workspace=$(mktemp -d)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

mkdir -p "$workspace"
cp -R "$REPO_ROOT/bin" "$workspace/"
cp -R "$REPO_ROOT/scripts" "$workspace/"
cp -R "$REPO_ROOT/config" "$workspace/"
chmod +x "$workspace/bin/multipowers"

cd "$workspace"

git init >/tmp/test_update_git_init_out.txt 2>/tmp/test_update_git_init_err.txt
git config user.email "test@example.com"
git config user.name "Test User"
git add .
git commit -m "init" >/tmp/test_update_git_commit_out.txt 2>/tmp/test_update_git_commit_err.txt

# Test 1: update --check returns actionable state.
echo "[TEST 1] update --check returns status"
if out=$(./bin/multipowers update --check --json 2>/tmp/test_update_check_err.txt); then
    if python3 - "$out" <<'PY'
import json
import sys
payload = json.loads(sys.argv[1])
for key in ["branch", "dirty", "has_upstream", "ahead", "behind"]:
    if key not in payload:
        raise SystemExit(f"missing key: {key}")
print("ok")
PY
    then
        echo "  [PASS] update --check payload is valid"
    else
        echo "  [FAIL] update --check payload invalid"
        cat /tmp/test_update_check_err.txt
        exit 1
    fi
else
    echo "  [FAIL] update --check should succeed"
    cat /tmp/test_update_check_err.txt
    exit 1
fi

# Test 2: apply requires explicit confirmation.
echo "[TEST 2] update --apply requires confirmation"
set +e
output=$(./bin/multipowers update --apply 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] update --apply without confirmation should fail"
    exit 1
fi
if echo "$output" | grep -q -- "--yes"; then
    echo "  [PASS] apply confirmation requirement enforced"
else
    echo "  [FAIL] missing apply confirmation diagnostic"
    echo "  Output: $output"
    exit 1
fi

# Test 3: dirty tree blocks apply safely.
echo "[TEST 3] dirty tree blocks update apply"
echo "dirty" >> bin/multipowers
set +e
output=$(./bin/multipowers update --apply --yes 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] dirty tree should block apply"
    exit 1
fi
if echo "$output" | grep -q "dirty"; then
    echo "  [PASS] dirty tree is handled safely"
else
    echo "  [FAIL] missing dirty tree diagnostic"
    echo "  Output: $output"
    exit 1
fi

echo ""
echo "Update command tests PASSED"
