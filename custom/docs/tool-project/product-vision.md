# Product Vision: Multipowers (Go Branch)

## Vision

Multipowers is a workflow-first orchestration plugin for Claude Code that makes multi-provider AI execution reliable, transparent, and repeatable via a **no-shell hybrid runtime**.

The goal is to provide consistent, high-quality decisions through structured Go-based orchestration, explicit quality gates, and role-specialized execution, while keeping reasoning in high-level Markdown skills.

## Product Positioning

- Project type: Go-powered plugin project for Claude Code.
- Primary audience: Maintainers and users running advanced `/mp:*` workflows.
- Core use cases:
  - Multi-AI software development (Double Diamond phases).
  - Multi-AI knowledge work (research, PRD, debate, etc.).
  - Deterministic state tracking and session resume across runs.

## Strategic Direction (Go Branch)

1. **Deterministic Core**: Move all state, validation, and routing logic from shell scripts to Go atomic commands and packages (`internal/`).
2. **First-Class Hooks**: Use a Go-based hook lifecycle (`mp hook run`) as the primary control plane for policy enforcement.
3. **Reasoning-Orchestration Split**: Maintain high-level reasoning in Markdown skills while delegating deterministic work to Go CLI surfaces.
4. **Provider Flexibility**: Support Codex, Gemini, and Claude-native execution with Go-level adapters and role-based routing policies.
5. **No-Shell Discipline**: Enforce strict no-shell runtime checks to prevent regression into fragile shell-based control flows.

## Product Principles

- **Go for Determinism, Markdown for Reasoning**: Logic and state management belong in Go; strategy and synthesis belong in Markdown.
- **Contract-Driven Communication**: Every Go command returns normalized JSON fields (`status`, `action`, etc.).
- **Evidence-First Completion**: Work is not "done" until verified by Go-based tests and validation gates.
- **Provider Transparency**: Explicit provider role selection based on task intent (via `router_intent.go`).
- **State Continuity**: Persistence is handled by `internal/tracks` to ensure reliability across sessions.

## Model Selection Principles (via `models.json`)

- **Planning, Architecture, and High-Importance Decisions**: Default to Codex (`gpt-5.3-codex`).
- **Heavy Implementation and Code Authoring**: Default to Claude Opus (`claude-opus`).
- **Documentation and Test-Case Generation**: Default to Claude Sonnet (`claude-sonnet`).
- **External-World Research**: Default to Gemini (`gemini-3-pro-preview`).
- **Quality Gates**:
  - Heavy/High-token audits: Claude Opus.
  - Lighter/Lower-token review passes: Codex.

## Success Signals

- Successful execution of full `embrace` workflows with Go-based state and gate checkpoints.
- `mp route --intent` accurately maps intents to provider lanes and skills.
- Zero reliance on legacy shell scripts for core decision-making or state management.
- Complete documentation of the migration path from upstream shell-era scripts to Go packages.
