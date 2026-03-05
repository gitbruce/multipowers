# Config-Driven Model Routing Verification

**Date:** 2026-03-04
**Status:** VERIFIED
**Implementation Plan:** `docs/plans/2026-03-04-config-driven-model-routing-implementation.md`

## Summary

The config-driven model routing implementation has been completed and verified. All 9 tasks (T01-T09) have been implemented with passing tests.

## Verification Commands

### 1. Build Policy and Runtime

```bash
./scripts/build.sh
```

Output:
```
Building policy...
loading source config from config
validating and compiling policy
wrote .claude-plugin/runtime/policy.json (checksum: 0c7cdc41762dbceb)
policy build ok
Building binaries...

Build complete:
  - /path/to/.claude-plugin/runtime/policy.json
  - /path/to/.claude-plugin/bin/mp
  - /path/to/.claude-plugin/bin/mp-devx
```

### 2. Run All Tests

```bash
go test ./internal/policy ./internal/hooks ./internal/cli ./internal/workflows ./internal/validation ./cmd/mp-devx -count=1
```

Output:
```
ok      github.com/gitbruce/multipowers/internal/policy      0.069s
ok      github.com/gitbruce/multipowers/internal/hooks       0.010s
ok      github.com/gitbruce/multipowers/internal/cli         0.150s
ok      github.com/gitbruce/multipowers/internal/workflows   0.912s
ok      github.com/gitbruce/multipowers/internal/validation  0.067s
ok      github.com/gitbruce/multipowers/cmd/mp-devx          0.006s
```

### 3. Verify Policy Compilation

```bash
./mp-devx -action build-policy
```

### 4. Verify Config Show/Hide Toggle

```bash
./mp config show-model-routing --dir . --json
./mp config show-model-routing --dir . --value off --json
```

## Architecture

### Development-Time (mutable)

- `config/workflows.yaml` - Workflow model policies with task-level overrides
- `config/agents.yaml` - Agent/persona model policies
- `config/executors.yaml` - Executor definitions and fallback chains

### Build-Time (deterministic)

```
config/*.yaml -> [compile] -> .claude-plugin/runtime/policy.json
```

### Run-Time (read-only)

- `.claude-plugin/runtime/policy.json` - Compiled policy
- `.claude-plugin/bin/mp` - Main binary
- `.claude-plugin/bin/mp-devx` - Dev tools binary

## Key Features Verified

1. **Workflow Task-Specific Model Selection**
   - Workflows support default + task-level overrides
   - Precedence: task -> default

2. **External Hard Enforcement**
   - External executors (codex_cli, gemini_cli) use hard enforcement
   - Model arg is injected into command template

3. **Claude Code Hint Enforcement**
   - claude_code executor uses hint enforcement
   - Model is advisory, not enforced

4. **One-Hop Automatic Fallback**
   - Fallback policies define single-hop chains
   - Cross-provider fallback supported

5. **Config Visibility Toggle**
   - `mp config show-model-routing` controls visibility
   - Default: visible (true)

6. **Hardcoded Model Guardrails**
   - Test scans for model strings outside allowed paths
   - Allowed: config/, testdata/, test files, docs

## Files Changed

### Created

- `config/workflows.yaml`
- `config/agents.yaml`
- `config/executors.yaml`
- `internal/policy/types.go`
- `internal/policy/load.go`
- `internal/policy/validate.go`
- `internal/policy/compile.go`
- `internal/policy/runtime_policy.go`
- `internal/policy/resolve.go`
- `internal/policy/dispatch.go`
- `internal/settings/runtime_settings.go`
- `internal/validation/model_hardcode_guard_test.go`
- `internal/validation/runtime_artifact_test.go`
- `internal/policy/e2e_resolve_dispatch_test.go`

### Modified

- `internal/hooks/handler.go`
- `internal/hooks/session_start.go`
- `internal/cli/root.go`
- `scripts/build.sh`

## Commit History

1. `feat(policy): add source config schemas for workflows agents executors`
2. `feat(policy): add source config loader and semantic validator`
3. `feat(policy): compile runtime policy and expose mp-devx build actions`
4. `feat(runtime): resolve workflow and agent execution contracts from compiled policy`
5. `feat(dispatch): enforce external model constraints and one-hop automatic fallback`
6. `feat(config): add show_model_routing toggle with default visible`
7. `test(policy): enforce no hardcoded runtime models outside config`
8. `build(runtime): generate policy and binaries for runtime-only plugin bundle`

## Remaining Risks

1. **Legacy modelroute compatibility**: The old `internal/modelroute` package is retained for backward compatibility but marked deprecated. Migration is gradual.

2. **External executor availability**: The system gracefully degrades when external CLIs (codex, gemini) are not installed, falling back through the chain.

3. **Policy drift**: The checksum in policy.json enables drift detection but enforcement is manual.

## Acceptance Criteria Status

| Criteria | Status |
|----------|--------|
| Workflow/agent routing decisions come from config/*.yaml via compiled policy | PASS |
| Non-config hardcoded model references are removed from runtime decision paths | PASS |
| External executor paths enforce model hard constraints with one-hop automatic fallback | PASS |
| Claude Code paths use hint-only routing | PASS |
| /mp:config controls routing/fallback visibility and defaults to visible | PASS |
| Build outputs include .claude-plugin/runtime/policy.json and .claude-plugin/bin/* | PASS |
