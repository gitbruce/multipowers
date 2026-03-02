# Product

## Summary
Claude Octopus on `go` branch adopts a no-shell hybrid runtime:
- Go runtime provides deterministic atomic capabilities.
- Markdown skills provide staged reasoning and orchestration.
- Script migration is tracked against upstream `v8.31.1` with per-file ownership mapping.

## Target Users
- Maintainers and contributors working on `go` branch runtime and docs.
- Engineers migrating shell-era behaviors into Go packages without regressing behavior.

## Primary Goal
Deliver a stable hybrid architecture with explicit contracts:
1. Atomic command surfaces for `state`, `validate`, `hook`, `route`, `test`, `coverage`, `status`.
2. Normalized JSON response contract usable by skills for branching decisions.
3. Full `v8.31.1` shell inventory mapped as `COPY_FROM_MAIN` or `MIGRATE_TO_GO`.

## Non-Goals
- Reintroducing shell runtime logic as the control plane.
- One-to-one mechanical translation of every shell script.
- Modifying `main` branch during `go` implementation delivery.

## Constraints
- Keep `main` branch untouched.
- Always sync with `upstream` before major mapping or implementation updates.
- Preserve contract fields: `status`, `action`, `error_code`, `message`, `data`, `remediation`.
- Keep compatibility facades where immediate command continuity is required.

## Initial Success Signals
- Atomic commands are callable and return normalized contract output.
- `docs/architecture/script-differences.md` stays parity-accurate with upstream baseline inventory.
- `COPY_FROM_MAIN` entries are copied and present on `go`.
- `MIGRATE_TO_GO` entries have concrete Go file/method targets.
- Verification gates pass on `go` branch before push.
