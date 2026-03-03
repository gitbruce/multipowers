# Go Branch Upstream Diff Discipline

## Rule Set

- Keep sync direction fixed: `upstream/main -> main -> go`.
- Keep `main` as mirror; keep implementation deltas on `go`.
- No legacy layering mechanism; all shared sync uses rules contracts.

## Structure Policy (`.claude` -> `.claude`)

`MUST_HOMOMORPHIC` (must remain structurally aligned):
- `.claude/commands` -> `.claude/commands` (shared subset)
- `.claude/skills` -> `.claude/skills` (shared subset)
- `.claude/references` -> `.claude/references`
- `.claude/state` -> `.claude/state`

`ALLOW_FORK` (explicitly allowed divergence):
- `.claude/commands/init.md`
- `.claude/commands/mp.md`
- `.claude/commands/persona.md`
- `.claude/skills/skill-persona.md`
- plus explicit ignore lists in `config/sync/claude-structure-rules.json`

## Execution Discipline

- Run `./scripts/sync-all.sh -dry-run` before applying sync.
- Run `./scripts/validate-claude-structure.sh -dry-run` before commit/push.
- Never resolve drift by reverting user local uncommitted files.
