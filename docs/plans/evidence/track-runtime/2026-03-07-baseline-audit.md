# 2026-03-07 Baseline Audit

- Baseline commit: `844315b6da14913288e1cafe6aa5b51ed5066765`
- Audit timestamp (UTC): `2026-03-07T13:07:12Z`

## Commands

### `go test ./... -count=1`
- Exit code: `0`
- Result: PASS
- First lines:

```text
?    github.com/gitbruce/multipowers/cmd/mp [no test files]
ok   github.com/gitbruce/multipowers/cmd/mp-devx 0.146s
ok   github.com/gitbruce/multipowers/internal/app 0.099s
ok   github.com/gitbruce/multipowers/internal/cli 0.087s
ok   github.com/gitbruce/multipowers/internal/orchestration 0.832s
ok   github.com/gitbruce/multipowers/internal/tracks 0.028s
```

### Focused implementation grep
- Command:

```bash
rg -n 'current.CurrentGroup = command|CurrentGroup = command|current.CurrentGroup = "post-tool"|trace_id|backoff|jitter|cache_hit|plan_graph_mermaid' internal custom docs -S
```

- Exit code: `0`
- Key findings:
  - `internal/cli/root.go` still writes `current.CurrentGroup = command`
  - `internal/hooks/post_tool_use.go` still writes `current.CurrentGroup = "post-tool"`
  - `internal/app/spec_track_lifecycle_test.go` mirrors the same command-as-group behavior
  - orchestration hardening items such as `trace_id`, retry backoff / jitter, cache hit metadata, and plan mermaid rendering appear in docs but not in runtime implementation files

## Re-baseline conclusion

### Spec-track runtime
- Implemented already:
  - `.multipowers/context/runtime.json` generation in init
  - hard migration to `.multipowers/tracks/tracks.md`
  - `TrackCoordinator`-based artifact creation and active track reuse
  - template rendering, complexity scoring, and lifecycle smoke coverage
- Still open:
  - `current_group` is overloaded with command names instead of real implementation groups
  - `last_commit_sha` is not written by the runtime path
  - group enforcement only becomes real for synthetic `gN` metadata states, not for normal CLI-driven progress

### Go orchestration optimization
- Implemented already:
  - routing hardcode guard subtasks in O01
  - `mp-devx lint-config` and lint gate wiring in O02-S02 / O02-S03
  - typed orchestration error code follow-through in O03 subtasks
  - retry policy fields noted in O04-S01
- Still open:
  - retry loop execution and deterministic retry matrix (`O04-S02`, `O04-S03`)
  - trace propagation and structured runtime logs (`O05-*`)
  - planner / report / degraded golden regressions (`O09-*`)
  - verification evidence and rollout docs (`O10-*`)
- Deferred from this closure wave unless code inspection later proves otherwise:
  - explainability CLI (`O06-*`)
  - result caching (`O07-*`)
  - plan visualization (`O08-*`)
