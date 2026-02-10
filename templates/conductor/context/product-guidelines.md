# Product Guidelines (User Project)

## Decision Rules

1. Prioritize end-user value over technical novelty.
2. Keep delivery workflow explicit (plan, implement, verify, document).
3. Keep scope tight: build what is requested, avoid overbuilding.
4. Require evidence before claiming completion.

## Change Policies

- Use tracks for non-trivial work and keep acceptance criteria visible.
- Prefer workflow-driven execution for major changes.
- Route specialist checkpoints (e.g., code review/security) to dedicated roles.

## Quality Gates

- Requirements are concrete and testable.
- Functional and error paths are validated.
- Backward compatibility and migrations are considered where relevant.
- Major changes run post-change scans/lint/tests and fix findings.

## Documentation Rules

- Update user-facing docs when behavior changes.
- Update developer docs when architecture/contracts change.
- Ensure changed-file scope maps to doc updates.
