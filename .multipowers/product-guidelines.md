# Product Guidelines

## Decision Rules
- Prefer atomic-first runtime design over opaque high-level black boxes.
- Keep deterministic behavior in Go; keep reasoning and orchestration in Markdown skills.
- Treat docs and contracts as executable governance, not optional notes.
- Resolve ambiguity by explicit mapping (source file -> target file -> target symbol).

## Conductor-Inspired Principles
- Plan is source of truth.
- Tech-stack and architecture changes must be explicit and reviewable.
- Prefer test-first for new behavior.
- Non-interactive and CI-safe commands only.
- Security and privacy impact must be documented.

## UX and Content
- Document behavior with concrete command contracts and absolute constraints.
- Avoid vague status outputs; every blocked/error state includes remediation.
- Keep migration language specific: strategy, target path, target symbol, status.

## Delivery Guardrails
- No direct feature work on `main`; all delivery lands on `go`.
- Sync upstream references before large edits and before release-ready pushes.
- No completion claims without verification output.
- Keep commits scoped: mapping updates, copy updates, verification evidence updates.
