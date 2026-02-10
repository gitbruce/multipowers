#!/usr/bin/env bash
# Core ask-role tests
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

echo "Testing ask-role core functionality..."

workspace=$(create_workspace)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cd "$workspace"

# Initialize context baseline
./bin/multipowers init >/dev/null 2>&1

# Test 1: Non-existent role should be rejected

echo "[TEST 1] Non-existent role validation"
set +e
output=$(./bin/ask-role non_existent_role "test" 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected non-zero exit for non-existent role"
    exit 1
fi
if echo "$output" | grep -q "Role 'non_existent_role' not found"; then
    echo "  [PASS] Non-existent role detected"
else
    echo "  [FAIL] Expected role not found message"
    echo "  Output: $output"
    exit 1
fi

# Test 2: Project config should take precedence over default config

echo "[TEST 2] Project config priority"
mkdir -p conductor/config
cat > conductor/config/roles.json <<'JSON_EOF'
{
  "roles": {
    "project-specific-role": {
      "description": "Project-only role",
      "tool": "system",
      "system_prompt": "Project specific"
    }
  }
}
JSON_EOF

set +e
output=$(./bin/ask-role project-specific-role "test" 2>&1)
status=$?
set -e

if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected non-zero for system role dispatch"
    exit 1
fi

if echo "$output" | grep -q "Using roles config: conductor/config/roles.json"; then
    echo "  [PASS] Project config takes precedence"
else
    echo "  [FAIL] Project config priority not observed"
    echo "  Output: $output"
    exit 1
fi

rm -f conductor/config/roles.json

# Test 3: Missing context strategy is consistent across strict/lenient

echo "[TEST 3] Context strategy consistency"
rm -rf conductor/context

set +e
output=$(./bin/ask-role prometheus "test" 2>&1)
status=$?
set -e
if [ $status -eq 0 ]; then
    echo "  [FAIL] Strict mode should fail when context is missing"
    exit 1
fi
if echo "$output" | grep -q "Missing required context files"; then
    echo "  [PASS] Strict mode blocks missing context"
else
    echo "  [FAIL] Strict mode missing expected context error"
    echo "  Output: $output"
    exit 1
fi

stub_dir=$(mktemp -d)
trap 'rm -rf "$workspace" "$stub_dir"' EXIT
cat > "$stub_dir/gemini" <<'PY'
#!/usr/bin/env python3
print("stub-gemini-ok")
PY
chmod +x "$stub_dir/gemini"

set +e
output=$(PATH="$stub_dir:$PATH" MULTIPOWERS_CONTEXT_MODE=lenient ./bin/ask-role prometheus "test" 2>&1)
status=$?
set -e
if [ $status -ne 0 ]; then
    echo "  [FAIL] Lenient mode should continue with warnings"
    echo "  Output: $output"
    exit 1
fi
if echo "$output" | grep -q "\[ASK-ROLE WARNING\] Missing required context files"; then
    echo "  [PASS] Lenient mode warns and continues"
else
    echo "  [FAIL] Lenient mode missing warning output"
    echo "  Output: $output"
    exit 1
fi

echo ""
echo "All ask-role core tests PASSED"
