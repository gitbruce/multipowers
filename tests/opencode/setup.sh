#!/usr/bin/env bash
# Setup script for OpenCode plugin tests
# Creates an isolated test environment with proper plugin installation
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

export TEST_HOME
TEST_HOME=$(mktemp -d)
export HOME="$TEST_HOME"
export XDG_CONFIG_HOME="$TEST_HOME/.config"
export OPENCODE_CONFIG_DIR="$TEST_HOME/.config/opencode"

install_root="$HOME/.config/opencode/superpowers"
mkdir -p "$install_root"

cp -r "$REPO_ROOT/lib" "$install_root/"
cp -r "$REPO_ROOT/skills" "$install_root/"

mkdir -p "$install_root/.opencode/plugins"
plugin_target="$install_root/.opencode/plugins/superpowers.js"
plugin_source="$REPO_ROOT/.opencode/plugins/superpowers.js"

if [ ! -f "$plugin_source" ]; then
    echo "[SETUP FAIL] Missing plugin runtime source: $plugin_source" >&2
    exit 1
fi

cp "$plugin_source" "$plugin_target"

mkdir -p "$HOME/.config/opencode/plugins"
ln -sf "$plugin_target" "$HOME/.config/opencode/plugins/superpowers.js"

mkdir -p "$HOME/.config/opencode/skills/personal-test"
cat > "$HOME/.config/opencode/skills/personal-test/SKILL.md" <<'EOF_SKILL'
---
name: personal-test
description: Test personal skill for verification
---
# Personal Test Skill

This is a personal skill used for testing.

PERSONAL_SKILL_MARKER_12345
EOF_SKILL

mkdir -p "$TEST_HOME/test-project/.opencode/skills/project-test"
cat > "$TEST_HOME/test-project/.opencode/skills/project-test/SKILL.md" <<'EOF_PROJECT'
---
name: project-test
description: Test project skill for verification
---
# Project Test Skill

This is a project skill used for testing.

PROJECT_SKILL_MARKER_67890
EOF_PROJECT

echo "Setup complete: $TEST_HOME"
echo "Plugin installed to: $plugin_target"
echo "Plugin registered at: $HOME/.config/opencode/plugins/superpowers.js"
echo "Test project at: $TEST_HOME/test-project"

cleanup_test_env() {
    if [ -n "${TEST_HOME:-}" ] && [ -d "$TEST_HOME" ]; then
        rm -rf "$TEST_HOME"
    fi
}

create_test_workspace() {
    local workspace
    workspace=$(mktemp -d)

    cp -R "$REPO_ROOT/bin" "$workspace/"
    cp -R "$REPO_ROOT/config" "$workspace/"
    cp -R "$REPO_ROOT/connectors" "$workspace/"
    cp -R "$REPO_ROOT/scripts" "$workspace/"
    cp -R "$REPO_ROOT/templates" "$workspace/"

    echo "$workspace"
}

export -f cleanup_test_env
export -f create_test_workspace
export REPO_ROOT
