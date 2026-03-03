---
command: setup
description: "Shortcut for /mp:sys-setup - Check Claude Octopus setup status"
redirect: sys-setup
---

# Setup (Shortcut)

This is a shortcut alias for `/mp:sys-setup`.

Running setup detection...

```bash
MP_BIN=""
if [[ -n "${CLAUDE_PLUGIN_ROOT:-}" ]] && [[ -x "${CLAUDE_PLUGIN_ROOT}/bin/mp" ]]; then
  MP_BIN="${CLAUDE_PLUGIN_ROOT}/bin/mp"
elif [[ -x "$PWD/.claude-plugin/bin/mp" ]]; then
  MP_BIN="$PWD/.claude-plugin/bin/mp"
elif [[ -x "./.claude-plugin/bin/mp" ]]; then
  MP_BIN="./.claude-plugin/bin/mp"
else
  echo "mp binary not found. Restart Claude Code and use /mp:* commands, or build via scripts/build.sh." >&2
  exit 1
fi

"$MP_BIN" detect-providers
```

For full setup documentation, see `/mp:sys-setup`.
