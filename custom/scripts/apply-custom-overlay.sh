#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

mkdir -p "$REPO_ROOT/.claude/commands"
cp "$REPO_ROOT/custom/commands/persona.md" "$REPO_ROOT/.claude/commands/persona.md"

node -e 'JSON.parse(require("fs").readFileSync(process.argv[1],"utf8"));' "$REPO_ROOT/custom/config/models.json"
node -e 'JSON.parse(require("fs").readFileSync(process.argv[1],"utf8"));' "$REPO_ROOT/custom/config/proxy.json"
node -e 'JSON.parse(require("fs").readFileSync(process.argv[1],"utf8"));' "$REPO_ROOT/custom/config/persona-lanes.json"

echo "Overlay applied successfully"
