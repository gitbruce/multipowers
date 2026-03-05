# Doctor Verification Evidence (2026-03-06)

## Test Verification

Command:

```bash
go test ./...
```

Result:
- Exit code: `0`
- Updated packages covered in this wave:
  - `internal/doctor`
  - `internal/cli`
  - `cmd/mp-devx`
  - `internal/hooks`
  - `internal/decisions`
  - `internal/issues`
  - `internal/faq`

## Smoke Verification

### 1) `mp-devx doctor --list`

Command:

```bash
go run ./cmd/mp-devx --action doctor --list
```

Result highlights:
- Printed table columns: `check_id`, `purpose`, `fail_capable`
- Listed all 16 checks in fixed alphabetical order by `check_id`

### 2) `mp-devx doctor --check-id config --json`

Command:

```bash
go run ./cmd/mp-devx --action doctor --check-id config --json
```

Result highlights:
- Exit code: `0`
- JSON check object includes required fields:
  - `check_id`
  - `status`
  - `message`
  - `detail`
  - `timed_out`
  - `elapsed_ms`
  - `timeout_ms`
  - `fail_capable`
- Observed `timeout_ms=45000` (single-check default)

### 3) `mp doctor --check-id config --timeout 10s --json`

Command:

```bash
tmpbin=$(mktemp -u /tmp/mp-devx.XXXXXX)
go build -o "$tmpbin" ./cmd/mp-devx
MP_DEVX_BIN="$tmpbin" go run ./cmd/mp doctor --check-id config --timeout 10s --json
rm -f "$tmpbin"
```

Result highlights:
- Exit code: `0`
- Output schema matches `mp-devx --action doctor` output
- Observed `timeout_ms=10000` (explicit override applied)
- Confirms `mp doctor` proxy path is using shared doctor engine

## Governance Artifacts Verified

- `.coderabbit.yaml` exists at repo root.
- `/mp:doctor` command asset exists at:
  - `.claude-plugin/.claude/commands/doctor.md`
- Plugin manifest includes doctor command path:
  - `./.claude/commands/doctor.md`
- Hook config includes:
  - `EnterPlanMode`
  - `WorktreeCreate`
  - `WorktreeRemove`
