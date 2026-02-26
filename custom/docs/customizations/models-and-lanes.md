# Models and Lanes

## What Changed From Upstream
Added role-based lane mapping in `custom/config/models.json`.

## Why This Exists
To keep workload-specific model selection without deep core edits.

## How To Use
Edit `custom/config/models.json` and run `./scripts/mp-devx overlay`.

## Operational Impact
Predictable model routing for heavy coding, docs/tests, architecture, and business analysis.

## Rollback Path
Remove/disable custom lane config and rely on upstream defaults.
