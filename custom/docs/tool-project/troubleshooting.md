# Troubleshooting (Tool Project Maintainers)

- If overlay commands are missing, run `./custom/scripts/apply-custom-overlay.sh`.
- If sync fails due to dirty tree, commit/stash and retry `./custom/scripts/sync-upstream.sh`.
- If proxy behavior is wrong, verify `custom/config/proxy.json`.
- If command docs drift, reapply overlay and verify `.claude/commands/*` against `custom/commands/*`.
