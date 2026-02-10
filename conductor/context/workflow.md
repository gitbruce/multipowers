# Workflow: Multipowers Delivery and Execution Routing

## 0. Setup + Track Baseline

1. Initialize/repair conductor scaffold.
2. Ensure context files exist and are meaningful.
3. Create track via `new track` and start execution from that track.

## 1. Router Routing Gate

For every incoming task, Router chooses one lane:

### Fast Lane (No Skill)

Use when task is small, bounded, and low coupling.

- Typical: tiny bug fix, config adjustment, narrow script update.
- Path: directly dispatch an execution role via `ask-role`.
- Goal: minimize overhead and complete quickly.

### Standard Lane (Workflow-Driven)

Use when change is significant or needs explicit quality gates.

- Select workflow first (e.g., brainstorming → writing-plans → subagent-driven-development).
- Execute workflow node-by-node.
- Switch roles only where workflow requires specialist nodes.

## 2. Standard Lane Node Pattern

Example (`subagent-driven-development`):

1. Default executor handles main implementation nodes.
2. Spec compliance node is reviewed by Architect role.
3. Code quality node is reviewed by Architect role.
4. Findings loop back to implementer until approved.

This enforces "flow first, role second".

## 3. Role Dispatch Mechanics

- All external role execution goes through `bin/ask-role`.
- `ask-role` injects context and uses role config to select connector.
- Connector invokes non-interactive model CLI (`codex`/`gemini`/others by config).

## 4. Major Change Governance Flow

When task is major:

1. Capture changed files (`git diff --name-only`).
2. Run checks:
   - `semgrep`
   - `biome` (TS/JS)
   - `ruff` (Python)
3. Fix issues and rerun checks.
4. Update docs mapped from changed files.
5. Only then mark task/track complete.

## 5. Template Sync Rule

When maintainers improve reusable workflow/context conventions in `conductor/`, evaluate and sync safe parts into `templates/conductor/` for downstream user projects.
