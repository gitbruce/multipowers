# Go Physical `.claude` Path Migration Design

## Context

Current `go` branch stores Claude workspace assets under `.claude-plugin/.claude`, while `main` uses root `.claude`.
The target is physical path parity with `main` for workspace assets.

## Decision

Use one-shot physical migration:

- Move `.claude-plugin/.claude` to root `.claude`
- Remove `.claude-plugin/.claude` entirely
- No compatibility layer, no symlink, no backup, no transition window

Keep `.claude-plugin` only for plugin packaging/runtime assets:

- `.claude-plugin/bin/*`
- `.claude-plugin/plugin.json`
- `.claude-plugin/marketplace.json`
- `.claude-plugin/custom/config/*`

## Goals

- `go` and `main` share the same physical Claude workspace path: `.claude/*`
- Eliminate path-mapping maintenance overhead for workspace assets
- Keep packaging/runtime pipeline stable

## Non-Goals

- No migration of plugin packaging assets out of `.claude-plugin`
- No dual-path support
- No backward compatibility for `.claude-plugin/.claude`

## Architecture

### 1) Path Model

- Workspace source of truth on `go`: `.claude/*`
- Plugin runtime/package artifacts remain in `.claude-plugin/*` (excluding `.claude-plugin/.claude`)

### 2) Policy and Validation

- Structure parity rules become same-path checks:
  - `main:.claude/commands -> go:.claude/commands`
  - `main:.claude/skills -> go:.claude/skills`
  - `main:.claude/references -> go:.claude/references`
  - `main:.claude/state -> go:.claude/state`
- Existing naming-policy exclusions remain:
  - `*.go`
  - `custom/**`
  - other non-parity scopes already approved in policy docs

### 3) Script and Tooling Alignment

- All scripts/tools reading workspace docs must use `.claude/*`
- Build/release wrappers for binaries keep `.claude-plugin/bin/*` paths

## Execution Order

1. Physical move `.claude-plugin/.claude` -> `.claude` and delete old path
2. Update sync rules, validators, tests, testdata to new roots
3. Update scripts/tooling path references for workspace docs
4. Update architecture/sync docs to same-path narrative
5. Run CI-equivalent checks and collect evidence
6. Merge only when no residual old-path references remain in active code/docs

## Risks and Handling

### Residual hardcoded path risk

- Detection: repo-wide search for `.claude-plugin/.claude`
- Action: fix all active references; allow only intentionally historical docs if explicitly accepted

### Validator mismatch risk

- Detection: `./scripts/validate-claude-structure.sh -dry-run`
- Action: fix rules and roots, avoid masking with broad ignores

### Runtime break risk in wrappers/scripts

- Detection: run smoke commands for `scripts/mp`, `scripts/mp-devx`, and relevant tooling scripts
- Action: normalize workspace read paths to `.claude/*`

## Verification Matrix

1. `go test ./internal/devx ./cmd/mp-devx -v`
2. `./scripts/validate-claude-structure.sh -dry-run`
3. `./scripts/mp-devx -action validate-structure-parity -dry-run`
4. `./scripts/report-claude-content-diff.sh`
5. `rg -n "\\.claude-plugin/\\.claude" scripts internal cmd config .github docs custom/docs`
6. Evidence doc with command, exit code, UTC timestamp, and first 60 output lines
