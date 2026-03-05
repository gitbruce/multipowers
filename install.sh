#!/bin/bash
# Multipowers Installer
#
# Recommended installation method (preferred over curl|bash):
#   claude plugin marketplace add https://github.com/gitbruce/multipowers
#   claude plugin install multipowers@nyldn-plugins --scope user
#   claude plugin enable multipowers --scope user
#   claude plugin update multipowers --scope user
#
# This script exists for convenience and uses the Claude Code plugin manager
# under the hood when available.

set -euo pipefail

echo "🐙 Installing Multipowers..."

if ! command -v claude >/dev/null 2>&1; then
  echo "❌ Error: Claude Code CLI ('claude') not found in PATH."
  echo ""
  echo "Install Claude Code first, then install the plugin with:"
  echo "  claude plugin marketplace add https://github.com/gitbruce/multipowers"
  echo "  claude plugin install multipowers@nyldn-plugins --scope user"
  echo "  claude plugin enable multipowers --scope user"
  echo "  claude plugin update multipowers --scope user"
  exit 1
fi

echo "📦 Using Claude Code plugin manager (recommended)..."

# Ensure marketplace exists and is fresh (idempotent).
claude plugin marketplace add https://github.com/gitbruce/multipowers >/dev/null 2>&1 || true
claude plugin marketplace update nyldn-plugins >/dev/null 2>&1 || true

# Install/enable/update (idempotent).
claude plugin install multipowers@nyldn-plugins --scope user >/dev/null 2>&1 || true
claude plugin enable multipowers --scope user >/dev/null 2>&1 || true
claude plugin update multipowers --scope user >/dev/null 2>&1 || true

echo ""
echo "✅ Installation complete!"
echo ""
echo "Next steps:"
echo "1. Fully restart Claude Code (Cmd+Q, then reopen)"
echo "2. Run: /octo:setup"
echo ""
echo "Troubleshooting:"
echo "- If commands don't appear, check: ~/.claude/debug/*.txt"
echo "- Verify install: claude plugin list | grep multipowers"
echo ""
