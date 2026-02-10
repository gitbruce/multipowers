# Workflow (User Project)

## 0) Setup and Context Anchoring

1. Initialize project conductor scaffold.
2. Fill `conductor/context/*.md` with app-specific constraints.
3. Validate baseline with project health checks.

## 1) Track Lifecycle

1. Create a track for meaningful work.
2. Start the track and define acceptance criteria.
3. Execute through fast lane or standard lane.
4. Complete track only after verification and doc sync.

## 2) Routing Policy

### Fast Lane

- Use for small, low-risk, bounded fixes.
- Dispatch role directly for quick turnaround.

### Standard Lane

- Use for significant, cross-cutting, or high-risk changes.
- Choose workflow first (`brainstorming` → `writing-plans` → execution workflow).
- In workflow nodes, switch to specialist reviewer roles when required.

## 3) Major Change Governance

For significant modifications:

1. Record changed files.
2. Run quality/security checks for the app stack.
3. Fix issues and rerun checks.
4. Update relevant docs (user + developer).
5. Attach evidence to track before completion.

## 4) Artifacts

- Design docs: `docs/design/`
- Implementation plans: `docs/plans/`
- Track records: `conductor/tracks/`
