#!/usr/bin/env bash
set -euo pipefail

SOURCE_REF="${1:-main}"
TARGET_REF="${2:-go}"

echo "content-diff source=${SOURCE_REF} target=${TARGET_REF}"
git diff --name-status "${SOURCE_REF}...${TARGET_REF}" -- \
  .claude \
  ':!**/*.go' \
  ':!custom/**' \
  ':!tests/**' \
  ':!hooks/**' \
  ':!openclaw/**' \
  ':!mcp-server/**' \
  ':!scripts/**' \
  ':!docs/plans/**' \
  ':!.multipowers/**' \
  ':!.github/workflows/**' \
  ':!config/sync/**' \
  ':!.dependencies/**' \
  ':!coverage.out' \
  ':!**/coverage.out' \
  ':!RELEASE_NOTES*.md' \
  ':!MIGRATION-*.md' \
  ':!IMPLEMENTATION_SUMMARY.md' || true
