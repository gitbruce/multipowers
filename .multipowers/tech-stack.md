# Tech Stack

## Runtime
- Runtime: Go 1.21+
- Framework: Go standard library + internal package architecture
- State/Docs: repository files under `.multipowers/`, docs, and internal state artifacts
- Deployment target: `origin/go` branch (after verification)

## Engineering Defaults
- Atomic CLI surfaces in `internal/cli` with structured responses.
- Domain logic in `internal/{tracks,context,validation,hooks,providers,workflows,...}`.
- No-shell runtime checks must remain active.
- Migration matrix in `docs/architecture/script-differences.md` is source of truth for shell parity.

## Core Command Surfaces
- `mp state get|set|update`
- `mp validate --type <workspace|no-shell|tdd-env|test-run|coverage>`
- `mp hook run --event <event>` (or normalized equivalent)
- `mp route --intent <intent> --provider-policy <policy>`
- `mp test run`
- `mp coverage check`
- `mp status`

## Contract Requirements
All atomic results should remain compatible with:
- `status`
- `action`
- `error_code`
- `message`
- `data`
- `remediation`
