---
command: doctor
description: "Run runtime governance diagnostics (16 checks) via mp-devx doctor engine"
---

# Doctor

Run Go-native diagnostics for runtime governance health.

Usage examples:

```bash
/mp:doctor
/mp:doctor --list
/mp:doctor --check-id config --timeout 10s --json
/mp:doctor --save --verbose
```

Execution:

```bash
MP_BIN=""
if [[ -n "${CLAUDE_PLUGIN_ROOT:-}" ]] && [[ -x "${CLAUDE_PLUGIN_ROOT}/bin/mp" ]]; then
  MP_BIN="${CLAUDE_PLUGIN_ROOT}/bin/mp"
elif [[ -x "$PWD/.claude-plugin/bin/mp" ]]; then
  MP_BIN="$PWD/.claude-plugin/bin/mp"
elif [[ -x "./.claude-plugin/bin/mp" ]]; then
  MP_BIN="./.claude-plugin/bin/mp"
else
  echo "mp binary not found. Build via mp-devx --action build-runtime." >&2
  exit 1
fi

# Forward options exactly as entered after /mp:doctor.
"$MP_BIN" doctor --dir "$PWD" $ARGUMENTS
```

Notes:
- `/mp:doctor` reuses `mp-devx --action doctor` via `mp doctor` proxy.
- Non-zero exit only when at least one check is `fail`.
