# FAQ

## Why does mp block before execution?
Because required context under `.multipowers/` is missing or incomplete.

## How to recover quickly?
Run `/mp:init` with explicit JSON answers, then verify required files exist:
- `product.md`
- `product-guidelines.md`
- `tech-stack.md`
- `workflow.md`
- `tracks.md`
- `CLAUDE.md`

## Where should context live?
Always under this repository's `.multipowers/` directory.

## Why keep main untouched?
This migration is governed on `go` branch. Keeping `main` untouched prevents accidental scope bleed and reduces integration risk.

## Why sync upstream continuously?
The mapping baseline depends on upstream versions and tags (for example `v8.31.1`). Syncing prevents stale parity calculations and mapping drift.
