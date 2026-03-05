---
event: Setup
description: Auto-initialize Multipowers workspace and verify providers on --init
---

# Setup Hook (Claude Code v2.1.10)

This hook runs automatically when Claude Code is started with `--init`, `--init-only`, or `--maintenance` flags.

## What It Does

1. Creates session-aware workspace directories
2. Verifies provider availability (Codex/Gemini)
3. Validates Claude Code version compatibility
4. Initializes analytics tracking

## Trigger

Runs on Setup event (v2.1.10 feature):
- `claude --init`
- `claude --init-only`
- `claude --maintenance`

## Implementation

```bash
${CLAUDE_PLUGIN_ROOT}/bin/mp init --quiet
${CLAUDE_PLUGIN_ROOT}/bin/mp detect-providers
```

## Output

On successful setup:
```
Multipowers workspace initialized
Session: ${CLAUDE_SESSION_ID}
Results: ~/.multipowers/results/${CLAUDE_SESSION_ID}/
Plans: ~/.multipowers/plans/${CLAUDE_SESSION_ID}/
```

## Workspace Structure

```
~/.multipowers/
├── results/
│   └── ${SESSION_ID}/           # Session-specific results
│       ├── .session-id
│       ├── .created-at
│       ├── probe-synthesis-*.md
│       └── ...
├── plans/
│   └── ${SESSION_ID}/           # Session-specific plans
├── logs/
│   └── ${SESSION_ID}/           # Session-specific logs
└── analytics/
    └── agent-usage.csv
```

## Benefits

- Automatic workspace setup on first use
- Provider issues detected early
- Session isolation for cleaner organization
