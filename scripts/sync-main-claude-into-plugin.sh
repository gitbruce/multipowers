#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# Refresh source assets from main, then keep go branch layout strict.
git checkout main -- .claude

rm -rf .claude-plugin/.claude
mv .claude .claude-plugin/.claude

# Keep bundled persona defaults aligned with main source-of-truth.
git checkout main -- config/agents.yaml
git checkout main -- config/orchestration.yaml
mkdir -p .claude-plugin/config .claude-plugin/agents
cp config/agents.yaml .claude-plugin/config/agents.yaml
cp config/orchestration.yaml .claude-plugin/config/orchestration.yaml
# Legacy mirror path for compatibility with older tooling.
cp config/agents.yaml .claude-plugin/agents/config.yaml

test ! -d .claude
echo "Synced main:.claude -> .claude-plugin/.claude"
