# Getting Started (Tool Project Maintainers)

## Branch Discipline

1. Keep `main` synced to `upstream/main` only.
2. Do all custom development on `multipowers`.
3. Never merge `multipowers` back into `main`.
4. Periodically merge `main` into `multipowers` and reapply overlay.

## Quick Start

```bash
git switch main
git fetch upstream
git merge --ff-only upstream/main
git switch multipowers
./custom/scripts/mp-devx overlay
./bin/mp persona list
```

## Daily Maintainer Workflow

1. Sync upstream with `./custom/scripts/mp-devx sync`
2. Reapply overlay
3. Validate with `bash tests/integration/test-sync-overlay.sh`
