#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

[[ -f "$ROOT/custom/templates/CLAUDE.md" ]]
[[ -f "$ROOT/custom/templates/FAQ.md" ]]
rg -n 'write_from_template "\$custom_templates_root/CLAUDE\.md" "\$croot/CLAUDE\.md"' "$ROOT/scripts/orchestrate.sh" >/dev/null
rg -n 'write_from_template "\$custom_templates_root/FAQ\.md" "\$croot/FAQ\.md"' "$ROOT/scripts/orchestrate.sh" >/dev/null

echo "PASS test-claude-faq-init-render"
