#!/usr/bin/env bash
# verify-architecture-diff-docs.sh
# Validates consistency across architecture diff documents
# Usage: scripts/verify-architecture-diff-docs.sh

set -euo pipefail

DOC1="docs/architecture/commands_skills_difference.md"
DOC2="docs/architecture/script-differences.md"
DOC3="docs/architecture/other-differences.md"
TRACKER="docs/architecture/gap-remediation-tracker.md"

err() { echo "ERROR: $*" >&2; exit 1; }
warn() { echo "WARN: $*" >&2; }

# Check all docs exist
for f in "$DOC1" "$DOC2" "$DOC3" "$TRACKER"; do
  [ -f "$f" ] || err "missing file: $f"
done

# Verify go baseline hash consistency
echo "Checking baseline hash consistency..."
vals=$(rg -o 'go=[0-9a-f]{7,}' "$DOC1" "$DOC2" "$DOC3" 2>/dev/null | awk -F: '{print $NF}' | sort -u | wc -l | tr -d ' ')
[ "$vals" = "1" ] || err "go baseline hash mismatch across docs (found $vals unique hashes)"

# Verify E0-E3 legend in each doc
echo "Checking evidence level legends..."
for f in "$DOC1" "$DOC2" "$DOC3"; do
  for level in E0 E1 E2 E3; do
    rg -q "$level" "$f" || err "missing $level legend in $f"
  done
done

# Verify decision tokens
echo "Checking decision tokens..."
for f in "$DOC1" "$DOC2" "$DOC3"; do
  rg -qi 'decision' "$f" || err "missing decision token in $f"
done

# Verify hook lifecycle events in script-differences.md
echo "Checking hook lifecycle events..."
for evt in SessionStart UserPromptSubmit PreToolUse PostToolUse Stop SubagentStop; do
  rg -q "$evt" "$DOC2" || err "missing lifecycle event $evt in script-differences.md"
done

# Verify mcp-server/openclaw decision tags in other-differences.md
echo "Checking mcp/openclaw decision tags..."
rg -n 'mcp-server/|openclaw/' "$DOC3" | rg -q 'DEFER_WITH_CONDITION|EXCLUDE_WITH_REASON|MIGRATE_TO_GO' \
  || warn "mcp/openclaw rows may be missing explicit decision (verify manually)"

# Verify tracker has required sections
echo "Checking tracker structure..."
rg -q "Commands/Skills High-Risk" "$TRACKER" || err "missing Commands/Skills High-Risk section in tracker"
rg -q "Script Missing Decision Classification" "$TRACKER" || err "missing Script Missing Decision Classification section in tracker"
rg -q "Other-Differences Partial/Missing Contracts" "$TRACKER" || err "missing Other-Differences section in tracker"
rg -q "E0 Upgrade Queue" "$TRACKER" || err "missing E0 Upgrade Queue section in tracker"

echo ""
echo "verify-architecture-diff-docs: PASS"
