#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# Refresh source assets from main, then keep go branch layout strict.
git checkout main -- .claude

rm -rf .claude-plugin/.claude
mv .claude .claude-plugin/.claude

# Keep bundled persona defaults aligned with main.
git checkout main -- agents/config.yaml
mkdir -p .claude-plugin/agents
cp agents/config.yaml .claude-plugin/agents/config.yaml

test ! -d .claude
echo "Synced main:.claude -> .claude-plugin/.claude"
