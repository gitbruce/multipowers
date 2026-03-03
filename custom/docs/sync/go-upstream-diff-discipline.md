# Go Branch Upstream Diff Discipline

## Rule Set

- Keep sync direction fixed: `upstream/main -> main -> go`.
- Keep `main` as mirror; keep implementation deltas on `go`.
- No legacy layering mechanism; all shared sync uses rules contracts.

## Structure Policy (`.claude` -> `.claude-plugin/.claude`)

`MUST_HOMOMORPHIC` (must remain structurally aligned):
- `.claude/commands` -> `.claude-plugin/.claude/commands` (shared subset)
- `.claude/skills` -> `.claude-plugin/.claude/skills` (shared subset)
- `.claude/references` -> `.claude-plugin/.claude/references`
- `.claude/state` -> `.claude-plugin/.claude/state`

`ALLOW_FORK` (explicitly allowed divergence):
- `.claude-plugin/.claude/commands/init.md`
- `.claude-plugin/.claude/commands/mp.md`
- `.claude-plugin/.claude/commands/persona.md`
- `.claude-plugin/.claude/skills/skill-persona.md`
- plus explicit ignore lists in `config/sync/claude-structure-rules.json`

## Execution Discipline

- Run `./scripts/sync-all.sh -dry-run` before applying sync.
- Run `./scripts/validate-claude-structure.sh -dry-run` before commit/push.
- Never resolve drift by reverting user local uncommitted files.
