# CLI Reference - `mp`

This document covers the Go-native `mp` CLI as it exists in the current `go` branch.

## Global Flags

Most commands accept:

- `--dir <path>`: target project directory (default `.`)
- `--prompt <text>`: free-form input or JSON payload, depending on the command
- `--json`: machine-readable response

## Initialization

### `mp init`

```bash
mp init --dir <project> --prompt '{...wizard answers json...}' --json
```

Behavior:

- bootstraps `.multipowers/` in the target project
- writes `.multipowers/context/runtime.json`
- initializes the canonical track registry at `.multipowers/tracks/tracks.md`
- creates track artifacts under `.multipowers/tracks/<track_id>/`
- sets `pre_run.enabled=false` by default in `runtime.json`

Important:

- `.multipowers/tracks.md` is legacy-only and is not read by the runtime
- missing context blocks spec execution with `run_init` guidance; files are never generated silently

## Track Runtime

### `mp track group-start`

```bash
mp track group-start \
  --dir <project> \
  --track-id <track_id> \
  --group g1 \
  --execution-mode workspace|worktree \
  --json
```

Starts an explicit implementation group for an existing spec track.

Rules:

- updates `current_group` / `group_status=in_progress`
- clears stale commit / verification evidence for the new group
- updates `.multipowers/tracks/tracks.md`
- blocks if the track requires linked-worktree execution and the current checkout is not a linked git worktree

### `mp track group-complete`

```bash
mp track group-complete \
  --dir <project> \
  --track-id <track_id> \
  --group g1 \
  --commit-sha <sha> \
  --json
```

Completes the active implementation group.

Rules:

- requires a matching active group
- requires `--commit-sha`
- records `last_commit_sha`, `last_verified_at`, and `completed_groups`
- updates `.multipowers/tracks/tracks.md`

## Spec-Driven Commands

```bash
mp plan "..."
mp discover "..."
mp define "..."
mp develop "..."
mp deliver "..."
mp embrace "..."
mp debate "..."
```

Runtime behavior:

- all spec artifacts land under `.multipowers/tracks/<track_id>/`
- the active track is coordinated by `internal/tracks.TrackCoordinator`
- spec command touches update `last_command` / `last_command_at`
- spec commands do **not** fake implementation progress by mutating `current_group`
- if an implementation group is still `in_progress` without commit / verification evidence, the next spec pipeline step is blocked

## Validation and Context

```bash
mp context validate --dir <project> --json
mp validate workspace --dir <project> --json
mp validate no-shell --dir <project> --json
mp validate tdd-env --dir <project> --json
mp validate test-run --dir <project> --json
mp validate coverage --dir <project> --json
```

## State and Status

```bash
mp state get --key <key> --json
mp state set --key <key> --value <value> --json
mp state update --data '{"k":"v"}' --json
mp status --dir <project> --json
```

## Orchestration Commands

```bash
mp orchestrate --phase <discover|define|develop|deliver> --prompt "..." --json
mp persona --agent <agent> --prompt "..." --json
mp loop --agent <agent> --prompt "..." --max-iterations <n> --json
mp route --intent <discover|define|develop|deliver> --json
mp config show-model-routing --json
```

Current orchestration hardening includes:

- step-level retry policy fields (`retry.idempotent`, `max_retries`, `backoff_ms`, `jitter_ratio`, `retryable_codes`)
- deterministic retry execution for eligible idempotent steps
- per-run `trace_id` propagation across plans, events, and step results
- structured lifecycle JSONL logs under `.multipowers/<logs_subdir>/orchestration-<trace_id>.jsonl`

## Hooks, Doctor, Policy, and Utilities

```bash
mp hook --event <event> --prompt "..." --json
mp doctor [--list] [--check-id <id>] [--timeout <duration>] [--save] [--verbose] [--json]
mp policy sync [--apply] [--ignore-id <id>] [--rollback-id <id>] [--revoke-id <id>] --json
mp policy stats --json
mp policy gc --json
mp policy tune --mode <balanced|accuracy|storage> --json
mp extract --url <url> --json
mp cost estimate --prompt "..." --json
mp cost report --metrics-dir <dir> --json
mp checkpoint save --checkpoint-id <id> --json
mp checkpoint get --checkpoint-id <id> --json
mp checkpoint delete --checkpoint-id <id> --json
```

## See Also

- `docs/ARCHITECTURE.md`
- `docs/plans/2026-03-07-go-orchestration-ux-followups.md`
