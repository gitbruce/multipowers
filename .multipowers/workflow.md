# Workflow

## Delivery Loop
1. Discover: verify intent, constraints, and upstream baseline.
2. Define: lock acceptance criteria, mapping policy, and branch governance.
3. Develop: implement incrementally with verification gates.
4. Deliver: validate contracts/parity and publish evidence.

## Working Agreements
- Always keep `main` untouched for this migration track.
- Always sync `upstream` before major mapping updates and before final push.
- Update `script-differences.md` whenever strategy/path/symbol ownership changes.
- Keep `COPY_FROM_MAIN` and `MIGRATE_TO_GO` counts explicit and verifiable.
- If blocked, record remediation rather than silently skipping requirements.

## Quality Gates
- `go fmt ./...`
- `go test ./... -short`
- `go vet ./...`
- no-shell runtime validation
- script matrix parity (`TOTAL == ROWS`)
- mapping completeness (no blank strategy/target/symbol/status fields)
