# Product Guidelines: Routing, Workflows, and Roles

## 1) Context Governance

- Read `conductor/context/*.md` before starting any track.
- Treat context files as stable background; update them only when product direction changes.
- Keep task-specific decisions in tracks/plans, not in stable context files.

## 2) Router Decision Model

Router is the global coordinator and must choose one lane for each task:

- **Fast Lane**: direct role dispatch without a skill workflow.
- **Standard Lane**: select a workflow first (for example, `subagent-driven-development`) and execute its nodes.

Selection factors include:

- change size and risk
- coupling across components
- need for explicit review checkpoints

## 3) Standard Lane Rule (Important)

In the Standard Lane, Router selects a **workflow**, not one fixed role for the entire task.

- Most nodes use the workflow’s default executor role.
- Specialized nodes must switch to dedicated specialist roles.
- Example: implementation by Coder, review nodes by Architect.

## 4) Role Contract

- **Router**: route work by task type and workflow needs.
- **Architect**: produce architecture/design/planning outputs and handle review/verification.
- **Coder**: implement changes with TDD discipline.
- **Librarian**: perform focused research and evidence gathering.

Role-to-model mapping is defined in role config and executed through CLI connectors.

## 5) Major Change Policy

For significant modifications, Router (or the delegated workflow) must:

1. Record the changed file list.
2. Run post-change checks: `semgrep`, `biome` (TS/JS), and `ruff` (Python).
3. Fix findings and rerun checks.
4. Update documentation affected by changed files.

No completion claim is allowed without these governance artifacts.
