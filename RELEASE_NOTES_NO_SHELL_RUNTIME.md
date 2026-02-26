# Release Notes: Strict No-Shell Runtime

## Summary

This release migrates runtime execution to Go-only entrypoints.

## Highlights

- All runtime shell scripts removed from repository.
- `/octo:*` command paths execute through `bin/octo`.
- New dev/CI helper runtime: `cmd/octo-devx`.
- Strict validator added: `octo validate --strict-no-shell --json`.
- Shell-to-Go mapping evidence captured before deletion:
  - `docs/plans/evidence/no-shell-runtime/mapping/sh-to-go-map.csv`

## Verification

- `rg --files | rg '\\.sh$'` => `0`
- `go test ./...` pass
- `go vet ./...` pass
- `go run ./cmd/octo validate --strict-no-shell --dir . --json` pass

## Notes

- If older docs mention shell paths, treat them as historical references only.
- Runtime and CI should use `bin/octo` and `go run ./cmd/octo-devx`.
