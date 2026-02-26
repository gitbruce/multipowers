# Go Big-Bang Migration Design (multipowers -> go branch)

Date: 2026-02-26
Status: Approved
Scope: Replace shell execution core with Go, maximize Claude Code hooks, keep command/skill layer thin.

## 1. Decisions Locked

- Migration mode: big-bang (not gradual), with freeze window.
- Branching: create `go` from `multipowers` and execute migration there.
- Go version: 1.22.
- File-size policy: soft warning when source file exceeds 500 lines.
- Spec context guard required files: `product.md`, `product-guidelines.md`, `tech-stack.md`, `workflow.md`, `tracks.md`, `CLAUDE.md`.
- `FAQ.md` and `context/runtime.json` are non-required for context readiness.
- Runtime pre-run contract: if `context/runtime.json` exists, execute pre-run and fail-fast.

## 2. Architecture

## 2.1 Single execution kernel

- Go CLI is the only execution kernel.
- `.claude/commands/*.md` and `.claude/skills/*.md` become thin wrappers that call `octo ... --json`.
- Common constraints are removed from Markdown and enforced in Go pipeline.

## 2.2 Unified pipeline

1. Resolve target root (`--dir` or `PWD`).
2. Context guard (5 core files + `CLAUDE.md`).
3. Auto-init when missing, then re-check.
4. Hard-stop if init fails or context still incomplete.
5. Runtime pre-run loading/execution (optional file, mandatory enforcement when present).
6. Command execution.
7. Post-process: track update + FAQ synthesis + event logging.

## 2.3 Hooks-first governance

Claude Code hooks are policy entry points. Hooks call Go subcommands only.

- `SessionStart`: inject stable summaries from:
  - `.multipowers/product.md`
  - `.multipowers/product-guidelines.md`
  - `.multipowers/tech-stack.md`
  - `.multipowers/workflow.md`
  - `.multipowers/CLAUDE.md`
  - current track status
- Summary limit: each file summary <= 20 lines.
- `UserPromptSubmit`: preflight guard for spec-driven `/mp:*`.
- `PreToolUse`: write-boundary and dangerous command governance.
- `PostToolUse`: FAQ/event/track post-processing.
- `Stop`/`SubagentStop`: prevent premature termination when mandatory workflow state is incomplete.

## 2.4 Provider layer

- Provider interface + registry (Codex/Gemini/Claude).
- Unified proxy routing through provider router.
- Debate/multi-LLM quorum: attempt up to 3 providers, require >=2 to continue.

## 2.5 Filesystem boundary

- During target-project execution, outputs must stay in target project.
- `.multipowers/*` is canonical workspace.
- No business artifact writes to `$HOME` or tool project paths.

## 3. Directory Layout

```text
cmd/mp/main.go
internal/app/*
internal/cli/*
internal/context/*
internal/runtime/*
internal/providers/*
internal/workflows/*
internal/tracks/*
internal/faq/*
internal/hooks/*
internal/fsboundary/*
internal/execx/*
internal/render/*
internal/util/*
pkg/api/*
.claude-plugin/hooks.json
```

## 4. Adoption from tmp/gemini.md and tmp/cc.md

Adopted:
- Go-only execution kernel.
- Unified JSON contract.
- `context guard` as first-class subcommand.
- Provider interface + registry.
- Middleware/pipeline enforcement.
- Domain-based package split.

Partially adopted:
- Multi-binary proposal converted to single binary subcommands.
- JSON IPC simplified to stdout JSON contract.
- Dynamic plugin loading deferred; static extension points first.

Rejected:
- Any `~/.claude-octopus/*` operational artifact paths.
- Gradual migration plan for this effort (big-bang chosen).

## 5. Acceptance Criteria

- All `/mp:*` execution paths enforced in Go pipeline.
- Context guard hard-stop behavior is deterministic and non-bypassable.
- Hook events route into Go handlers and are test-covered.
- Runtime pre-run enforced when runtime config exists.
- Debate quorum behavior enforced: 3->2 continue, <2 fail.
- Target project boundary policy enforced end-to-end.
- Key behavior parity verified for `plan`, `debate`, `embrace` (structure and workflow semantics), not strict byte-for-byte markdown identity.
- Performance benchmark targets met:
  - `octo hook preflight` hot path p95 <= 50ms
  - cold start p95 <= 120ms

## 6. Common Patterns Extraction Summary

### Pattern 1: Context Guard
- Current: repeated shell checks in many skills/commands.
- Go replacement: `octo context guard --dir <target> --json --auto-init`
- Contract:
  - validate required files (5 + `CLAUDE.md`)
  - if missing and `--auto-init`, run init once and re-check
  - return JSON: `status`, `missing_files`, `init_triggered`, `error_code`

### Pattern 2: Provider Detection
- Current: repeated `command -v` shell snippets.
- Go replacement: `internal/providers/detector.go` + registry status APIs.
- Contract: `DetectAll()` reports provider availability consistently for commands/hooks.

### Pattern 3: State Management
- Current: shell + inline Python JSON handling.
- Go replacement: typed state manager (`internal/state/*` equivalent under `tracks`/`faq`/`app`).
- Contract: atomic read/write with schema validation and deterministic error codes.

### Pattern 4: Visual Indicators
- Current: markdown templates and shell printf duplication.
- Go replacement: centralized render package (`internal/render/*`).
- Contract: provider/status banners are generated by Go, not duplicated in markdown.

### Pattern 5: Validation Gates
- Current: shell file-glob checks and ad-hoc exits.
- Go replacement: governance gates in pipeline/hook handlers.
- Contract: reject invalid execution states with machine-readable remediation.

Important path policy:
- No operational artifacts under `~/.claude-octopus/*`.
- Runtime artifacts must live under target project `/.multipowers/*`.

## 7. Core Interface Design (Reference)

```go
// internal/core/interfaces.go (reference shape)
type ContextGuard interface {
    Validate(projectDir string) (*ValidationResult, error)
    RequiredFiles() []string
    AutoInit(projectDir string) error
}

type ProviderDetector interface {
    DetectAll() (*ProviderStatus, error)
    IsAvailable(provider string) bool
}

type GovernanceGate interface {
    Check(gateType string) (*GateResult, error)
    Block(reason string) error
}

type HookHandler interface {
    Handle(event HookEvent) (*HookResult, error)
}

type BoundaryPolicy interface {
    ValidateWritePath(targetPath string, projectRoot string) error
}
```
