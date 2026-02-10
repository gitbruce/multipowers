#!/usr/bin/env bash
# Roles schema validation tests
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$REPO_ROOT"

echo "Testing roles schema validation..."

validator="scripts/validate_roles.py"
schema="config/roles.schema.json"

if [ ! -f "$validator" ]; then
    echo "  [FAIL] Missing validator: $validator"
    exit 1
fi
if [ ! -f "$schema" ]; then
    echo "  [FAIL] Missing schema: $schema"
    exit 1
fi

# Test 1: Default config should validate

echo "[TEST 1] Default config validates"
if python3 "$validator" --config config/roles.default.json --schema "$schema" --quiet >/dev/null 2>&1; then
    echo "  [PASS] Default config is valid"
else
    echo "  [FAIL] Default config should be valid"
    exit 1
fi

# Test 2: Missing required field should fail

echo "[TEST 2] Missing required field fails"
tmp_cfg=$(mktemp)
cat > "$tmp_cfg" <<'JSON_EOF'
{
  "roles": {
    "test-role": {
      "description": "Test role",
      "tool": "gemini"
    }
  }
}
JSON_EOF

if python3 "$validator" --config "$tmp_cfg" --schema "$schema" --quiet >/dev/null 2>&1; then
    echo "  [FAIL] Invalid config should fail validation"
    rm -f "$tmp_cfg"
    exit 1
else
    echo "  [PASS] Missing required field detected"
fi
rm -f "$tmp_cfg"

# Test 3: Invalid tool enum should fail

echo "[TEST 3] Invalid tool enum fails"
tmp_cfg=$(mktemp)
cat > "$tmp_cfg" <<'JSON_EOF'
{
  "roles": {
    "test-role": {
      "description": "Test role",
      "tool": "invalid-tool",
      "system_prompt": "Test"
    }
  }
}
JSON_EOF

if python3 "$validator" --config "$tmp_cfg" --schema "$schema" --quiet >/dev/null 2>&1; then
    echo "  [FAIL] Invalid tool should fail validation"
    rm -f "$tmp_cfg"
    exit 1
else
    echo "  [PASS] Invalid tool detected"
fi
rm -f "$tmp_cfg"

# Test 4: Invalid args type should fail

echo "[TEST 4] Invalid args type fails"
tmp_cfg=$(mktemp)
cat > "$tmp_cfg" <<'JSON_EOF'
{
  "roles": {
    "test-role": {
      "description": "Test role",
      "tool": "gemini",
      "system_prompt": "Test",
      "args": "not-an-array"
    }
  }
}
JSON_EOF

if python3 "$validator" --config "$tmp_cfg" --schema "$schema" --quiet >/dev/null 2>&1; then
    echo "  [FAIL] Invalid args type should fail validation"
    rm -f "$tmp_cfg"
    exit 1
else
    echo "  [PASS] Invalid args type detected"
fi
rm -f "$tmp_cfg"

echo ""
echo "All roles schema tests PASSED"
