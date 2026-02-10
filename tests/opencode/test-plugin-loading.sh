#!/usr/bin/env bash
# Test: Plugin Loading
# Verifies that the superpowers plugin loads correctly in OpenCode
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "=== Test: Plugin Loading ==="

# Source setup to create isolated environment
source "$SCRIPT_DIR/setup.sh"

# Trap to cleanup on exit
trap cleanup_test_env EXIT

# Test 1: Verify plugin file exists and is registered
echo "Test 1: Checking plugin registration..."
if [ -L "$HOME/.config/opencode/plugins/superpowers.js" ]; then
    echo "  [PASS] Plugin symlink exists"
else
    echo "  [FAIL] Plugin symlink not found at $HOME/.config/opencode/plugins/superpowers.js"
    exit 1
fi

# Verify symlink target exists
if [ -f "$(readlink -f "$HOME/.config/opencode/plugins/superpowers.js")" ]; then
    echo "  [PASS] Plugin symlink target exists"
else
    echo "  [FAIL] Plugin symlink target does not exist"
    exit 1
fi

# Test 2: Verify lib/skills-core.js is in place
echo "Test 2: Checking skills-core.js..."
if [ -f "$HOME/.config/opencode/superpowers/lib/skills-core.js" ]; then
    echo "  [PASS] skills-core.js exists"
else
    echo "  [FAIL] skills-core.js not found"
    exit 1
fi

# Test 3: Verify skills directory is populated
echo "Test 3: Checking skills directory..."
skill_count=$(find "$HOME/.config/opencode/superpowers/skills" -name "SKILL.md" | wc -l)
if [ "$skill_count" -gt 0 ]; then
    echo "  [PASS] Found $skill_count skills installed"
else
    echo "  [FAIL] No skills found in installed location"
    exit 1
fi

# Test 4: Check using-superpowers skill exists (critical for bootstrap)
echo "Test 4: Checking using-superpowers skill (required for bootstrap)..."
if [ -f "$HOME/.config/opencode/superpowers/skills/using-superpowers/SKILL.md" ]; then
    echo "  [PASS] using-superpowers skill exists"
else
    echo "  [FAIL] using-superpowers skill not found (required for bootstrap)"
    exit 1
fi

# Test 5: Verify plugin JavaScript syntax (basic check)
echo "Test 5: Checking plugin JavaScript syntax..."
plugin_file="$HOME/.config/opencode/superpowers/.opencode/plugins/superpowers.js"
if node --check "$plugin_file" 2>/dev/null; then
    echo "  [PASS] Plugin JavaScript syntax is valid"
else
    echo "  [FAIL] Plugin has JavaScript syntax errors"
    exit 1
fi

# Test 6: Verify runtime capabilities exported by plugin
echo "Test 6: Checking plugin runtime capabilities..."
if PLUGIN_FILE="$plugin_file" SUPERPOWERS_ROOT="$HOME/.config/opencode/superpowers" PERSONAL_ROOT="$HOME/.config/opencode" node --input-type=module - <<'NODE'
import process from 'process';

const pluginPath = process.env.PLUGIN_FILE;
const mod = await import(`file://${pluginPath}`);
const plugin = mod.default;

if (!plugin || plugin.name !== 'superpowers') {
  throw new Error('plugin export missing or invalid name');
}
if (typeof plugin.setup !== 'function') {
  throw new Error('plugin.setup must be a function');
}

const runtime = plugin.setup({
  superpowersRoot: process.env.SUPERPOWERS_ROOT,
  personalRoot: process.env.PERSONAL_ROOT,
});

if (!runtime || typeof runtime.listSkills !== 'function' || typeof runtime.resolveSkill !== 'function') {
  throw new Error('runtime methods listSkills/resolveSkill missing');
}

const skills = runtime.listSkills();
if (!Array.isArray(skills) || skills.length < 1) {
  throw new Error('listSkills returned no skills');
}

const resolved = runtime.resolveSkill('using-superpowers');
if (!resolved || typeof resolved.skillFile !== 'string') {
  throw new Error('resolveSkill did not resolve using-superpowers');
}

console.log('runtime-ok');
NODE
then
    echo "  [PASS] Plugin exposes meaningful runtime behavior"
else
    echo "  [FAIL] Plugin runtime behavior check failed"
    exit 1
fi

# Test 7: Verify personal test skill was created
echo "Test 7: Checking test fixtures..."
if [ -f "$HOME/.config/opencode/skills/personal-test/SKILL.md" ]; then
    echo "  [PASS] Personal test skill fixture created"
else
    echo "  [FAIL] Personal test skill fixture not found"
    exit 1
fi

echo ""
echo "=== All plugin loading tests passed ==="
