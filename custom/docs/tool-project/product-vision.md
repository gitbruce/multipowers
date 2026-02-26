# Product Vision: Claude Octopus

## Vision

Claude Octopus is a workflow-first orchestration plugin for Claude Code that makes multi-provider AI execution reliable, transparent, and repeatable for both software delivery and knowledge work.

The target outcome is not "more prompts". The target outcome is consistent, high-quality decisions and deliverables through structured phases, explicit quality gates, and role-specialized execution.

## Product Positioning

- Project type: tool/plugin project for Claude Code.
- Primary audience: maintainers and advanced users running `/octo:*` workflows.
- Core use cases:
  - Multi-AI software workflows (discover, define, develop, deliver, embrace).
  - Multi-AI knowledge workflows (research, PRD, debate, docs/deck output).
  - Cross-session orchestration with resumable state and governance checks.

## Strategic Direction

1. Keep `bin/octo` as the single operational engine for routing, provider execution, and workflow control.
2. Keep command UX simple (`/octo:*` in Claude Code, CLI fallback via `./bin/octo`) while expanding capability behind stable interfaces.
3. Maintain provider flexibility (Codex, Gemini, Claude-native) with graceful degradation when one or more providers are unavailable.
4. Continue strengthening quality gates, validation hooks, and review loops so major changes cannot silently bypass verification.
5. Preserve context hygiene: stable project context in `conductor/context/*`, task-specific execution state in workflow artifacts.

## Product Principles

- Workflow before improvisation: major work follows explicit phases and checkpoints.
- Evidence before completion: recommendations require verifiable outputs (tests, review notes, validation artifacts).
- Role specialization: personas and skills should be selected by task intent, not by habit.
- Provider transparency: users should see which providers are active and what each contributes.
- Safe defaults: supervised operation and clear fallbacks over opaque autonomous behavior.

## Model Selection Principles

- Planning, architecture, and other high-importance decisions default to Codex (`gpt-5.3-codex`).
- Heavy coding and implementation default to Claude Opus (mapped in this environment to GLM-5).
- Documentation and test-case authoring default to Claude Sonnet (mapped in this environment to GLM-4.7).
- External-world research tasks (market, ecosystem, literature, competitive research) default to Gemini (`gemini-3-pro-preview`).
- Quality checks use a split policy:
  - Heavy, high-token audits use Claude Opus (GLM-5).
  - Lighter, lower-token review passes use Codex.

## Success Signals

- Users can reliably run full Double Diamond flows with clear phase transitions and gate outcomes.
- `auto` routing chooses the right workflow mode with minimal correction.
- Multi-provider runs produce measurably better synthesis than single-provider baselines.
- Session resume and state tracking reduce rework between runs.
- Documentation, commands, and tests remain aligned with shipped behavior.
