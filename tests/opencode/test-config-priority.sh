#!/usr/bin/env bash
# Configuration priority tests
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

workspace=$(create_workspace)
cleanup() {
    rm -rf "$workspace"
}
trap cleanup EXIT

cd "$workspace"

echo "Testing configuration priority..."

./bin/multipowers init >/dev/null 2>&1

# Test 1: Default config is used when project config does not exist

echo "[TEST 1] Default config fallback"
rm -f conductor/config/roles.json
set +e
output=$(./bin/ask-role prometheus "test" 2>&1)
status=$?
set -e

if echo "$output" | grep -q "Using roles config: config/roles.default.json"; then
    echo "  [PASS] Default config selected"
else
    echo "  [FAIL] Expected default config path"
    echo "  Output: $output"
    exit 1
fi

if [ $status -ne 0 ]; then
    echo "  [PASS] Connector stage failure is acceptable in test environment"
fi

# Test 2: Project config takes precedence

echo "[TEST 2] Project config precedence"
mkdir -p conductor/config
cat > conductor/config/roles.json <<'JSON_EOF'
{
  "roles": {
    "project-role": {
      "description": "Project-only role",
      "tool": "system",
      "system_prompt": "Project config role"
    }
  }
}
JSON_EOF

set +e
output=$(./bin/ask-role project-role "test" 2>&1)
status=$?
set -e

if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected non-zero for system role dispatch"
    exit 1
fi
if echo "$output" | grep -q "Using roles config: conductor/config/roles.json"; then
    echo "  [PASS] Project config selected"
else
    echo "  [FAIL] Expected project config path"
    echo "  Output: $output"
    exit 1
fi

# Test 3: Invalid project config should fail fast

echo "[TEST 3] Invalid project config fails fast"
echo "invalid-json" > conductor/config/roles.json

set +e
output=$(./bin/ask-role prometheus "test" 2>&1)
status=$?
set -e

if [ $status -eq 0 ]; then
    echo "  [FAIL] Expected non-zero for invalid project config"
    exit 1
fi
if echo "$output" | grep -q "Invalid JSON"; then
    echo "  [PASS] Invalid project config detected"
else
    echo "  [FAIL] Expected invalid JSON error"
    echo "  Output: $output"
    exit 1
fi

echo ""
echo "Configuration priority tests PASSED"
