# Troubleshooting (Tool Project Maintainers)

- If overlay commands are missing, run `./custom/scripts/mp-devx overlay`.
- If sync fails due to dirty tree, commit/stash and retry `./custom/scripts/mp-devx sync`.
- If proxy behavior is wrong, verify `custom/config/proxy.json`.
- If command docs drift, reapply overlay and verify `.claude/commands/*` against `custom/commands/*`.
