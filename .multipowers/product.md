# Product

## Summary
Claude Octopus `go` branch is not a version-only fork. It is a standards-preserving evolution of the upstream product baseline (`https://github.com/nyldn/claude-octopus`) into a no-shell hybrid runtime:
- Keep upstream product intent (multi-provider orchestration, structured workflows, quality gates).
- Move runtime-critical logic to deterministic Go atomic commands.
- Keep reasoning/orchestration in Markdown skills.
- Maintain explicit script migration traceability from upstream `v8.31.1`.

## Upstream Base Requirements
The upstream project baseline remains the product foundation and must be reflected in `go` branch behavior/design:
- Multi-provider orchestration model: Codex (implementation depth), Gemini (ecosystem breadth), Claude (synthesis/orchestration).
- Structured methodology: Double Diamond phases (`discover`, `define`, `develop`, `deliver`) with phase quality gates.
- Consensus-oriented delivery: disagreement should be surfaced and gated, not silently ignored.
- Production workflows: lifecycle orchestration, review/security gates, research/debate workflows, and persona-driven specialization.
- Transparency and safety: explicit provider visibility, namespace isolation, deterministic outputs, and clear operational boundaries.

## Invocation and Call Process (Claude Code Features)
`go` branch execution must explicitly leverage Claude Code-native capabilities rather than ad-hoc shell control flow.

### 1. Command Entry and Intent Routing
- User invokes `/mp:*` commands.
- Runtime resolves intent and calls atomic surfaces:
  - `mp state get|set|update`
  - `mp validate --type <workspace|no-shell|tdd-env|test-run|coverage>`
  - `mp hook run --event <event>`
  - `mp route --intent <intent> --provider-policy <policy>`
  - `mp test run`
  - `mp coverage check`
  - `mp status`

### 2. Hook Lifecycle as First-Class Control Plane
The flow must rely on hook events for deterministic policy enforcement and context continuity:
- `SessionStart`: initialize session state and context synchronization.
- `UserPromptSubmit`: check context contract and block/guide when incomplete.
- `PreToolUse`: validate routing/security/policy before potentially expensive actions.
- `PostToolUse`: enforce quality checkpoints and telemetry updates.
- `Stop/SubagentStop`: finalize state and capture completion signals.
- Additional operational events (where available): `ConfigChange`, `WorktreeCreate`, `WorktreeRemove`.

### 3. Agent-Teams Orchestration
- Use teammate lifecycle hooks (for example idle dispatch and task-completion transitions) as orchestration signals.
- Maintain explicit dependency/transition checks between task units.
- Keep dispatch/governance decisions auditable via structured state and migration docs.

### 4. Reasoning + Deterministic Contract Boundary
- Skills decide branching and questioning.
- Go commands return normalized JSON contract fields:
  - `status`, `action`, `error_code`, `message`, `data`, `remediation`.
- No ambiguous free-form-only failures for runtime-critical paths.

## Target Users
- Maintainers and contributors working on `go` branch runtime, hooks, and migration docs.
- Engineers translating upstream shell-era behavior into Go package ownership without parity drift.
- Reviewers validating that architectural changes preserve upstream product goals.

## Primary Goal
Deliver a stable hybrid architecture that is both upstream-aligned and operationally strict:
1. Atomic command contracts are callable and deterministic.
2. Hook-driven control flow is explicit and policy-enforced.
3. Agent-teams lifecycle integration is documented and operationally usable.
4. `v8.31.1` script inventory is fully classified (`COPY_FROM_MAIN` vs `MIGRATE_TO_GO`) with concrete ownership targets.

## Non-Goals
- Reintroducing shell runtime logic as primary control plane.
- One-to-one mechanical translation of every shell script.
- Modifying `main` branch during `go` implementation delivery.
- Replacing reasoning-layer decisions with opaque monolithic runtime behavior.

## Constraints
- Keep `main` branch untouched.
- Always sync with `upstream` before major mapping/implementation updates.
- Preserve normalized response contract (`status`, `action`, `error_code`, `message`, `data`, `remediation`).
- Keep compatibility facades where continuity is required, but do not hide atomic ownership.
- Push implementation work to remote `go` branch only after verification.

## Initial Success Signals
- Upstream baseline requirements are explicitly represented in `go` branch product docs and runtime contracts.
- Hook lifecycle is part of the documented invocation process, not implied behavior.
- Agent-teams orchestration signals are tracked and policy-gated.
- `docs/architecture/script-differences.md` remains parity-accurate with upstream baseline inventory.
- `COPY_FROM_MAIN` artifacts exist where intended; `MIGRATE_TO_GO` rows map to concrete file/method targets.
- Verification gates pass before pushing to `origin/go`.
