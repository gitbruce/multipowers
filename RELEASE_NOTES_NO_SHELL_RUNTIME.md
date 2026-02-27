# Release Notes: Strict No-Shell Runtime

## Summary

This release migrates runtime execution to Go-only entrypoints.

## Highlights

- All runtime shell scripts removed from repository.
- `/mp:*` command paths execute through `.claude-plugin/bin/mp`.
- New dev/CI helper runtime: `cmd/mp-devx`.
- Strict validator added: `octo validate --strict-no-shell --json`.
- Shell-to-Go mapping evidence captured before deletion:
  - `docs/plans/evidence/no-shell-runtime/mapping/sh-to-go-map.csv`

## Verification

- `rg --files | rg '\\.sh$'` => `0`
- `go test ./...` pass
- `go vet ./...` pass
- `go run ./cmd/mp validate --strict-no-shell --dir . --json` pass

## Notes

- If older docs mention shell paths, treat them as historical references only.
- Runtime and CI should use `.claude-plugin/bin/mp` and `go run ./cmd/mp-devx`.
