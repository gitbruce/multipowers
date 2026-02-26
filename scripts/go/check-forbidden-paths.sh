#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")/../.." && pwd)"
forbidden='~/.claude-octopus|/home/.*/\.claude-octopus'
if rg -n "$forbidden" "$root/internal" "$root/cmd" "$root/pkg" >/tmp/forbidden.txt 2>/dev/null; then
  echo "ERROR: forbidden path reference found"
  cat /tmp/forbidden.txt
  exit 1
fi
echo "OK: no forbidden runtime paths"
