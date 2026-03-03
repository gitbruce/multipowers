# Main-Go Naming Parity Design

## Context

Current sync direction is fixed as `upstream/main -> main -> go`.
The `go` branch keeps product-specific customization, but long-term maintenance requires reducing avoidable structural drift from `main`.

## Goals

- Keep naming conventions aligned with `main` for:
  - Markdown files (`*.md`)
  - Claude Code related files under `.claude*` policy/command/skill/reference/state scopes
- Do not constrain Go implementation files (`*.go`), which remain independently evolved in `go`.
- Allow controlled customization in `init/mp/persona/skill-persona` while preserving naming discipline.

## Non-Goals

- No physical migration from `.claude-plugin/.claude` to `.claude`.
- No requirement for full content parity across all Markdown or Claude Code files.
- No naming parity enforcement for Go source files.

## Design

### 1) Policy Model

Use three policy classes:

- `MUST_HOMOMORPHIC`
  - Path and file-set parity required for shared static knowledge surfaces.
  - Targets shared command/skill/reference/state subsets.
- `ALLOW_FORK_WITH_NAME_PARITY`
  - Content divergence allowed.
  - File naming must follow `main` naming conventions.
  - Add/remove operations must be explicitly registered in the rules contract.
  - Applied to `init/mp/persona/skill-persona`.
- `ALLOW_FORK`
  - Explicitly go-only artifacts not expected to match `main`.

### 2) Scope Rules

- Enforce naming parity for:
  - `*.md` in mapped parity scopes
  - Claude Code related files in `.claude` / `.claude-plugin/.claude` mapped scopes
- Exclude from parity:
  - `*.go` (all Go files remain unconstrained by naming parity)

### 3) Validation Flow

Validation command stays read-only (`-dry-run`) and evaluates:

1. Source (`main`) and target (`go`) tree names under configured roots.
2. Rule decision per root pair.
3. For `ALLOW_FORK_WITH_NAME_PARITY`, ensure:
   - naming-format compliance vs `main` conventions;
   - new/deleted filenames are present in explicit allowlist entries.
4. Emit machine-readable summary and non-zero exit on rule violations.

### 4) CI Integration

Two-lane checks:

- Required gate:
  - naming/structure parity check for `*.md` and Claude Code-related scoped files
- Informational report:
  - content-level diff summary without blocking merge

This keeps drift pressure on structure/naming while preserving customization flexibility.

## Error Handling

- Missing rules file or malformed rule entry: fail fast with actionable error path.
- Unknown rule decision: fail validation.
- Unregistered add/delete in `ALLOW_FORK_WITH_NAME_PARITY`: fail validation and list offending files.
- Go file mismatch reports are suppressed by scope filter to avoid false policy failures.

## Testing Strategy

- Unit tests:
  - rule loading and schema validation
  - filename normalization and naming comparison
  - allowlist registration enforcement for add/delete
- Integration tests:
  - git tree resolution for `main` vs `go`
  - `validate-structure-parity -dry-run` return code and key output
- CI assertions:
  - required job fails on naming/structure violations
  - informational job always reports but does not block

## Rollout

1. Extend rules contract with `ALLOW_FORK_WITH_NAME_PARITY`.
2. Update validator logic and tests.
3. Wire CI required/informational lanes.
4. Update sync docs with final policy language and exception process.
5. Capture verification evidence with timestamps and command outputs.
