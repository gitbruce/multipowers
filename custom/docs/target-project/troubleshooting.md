# Troubleshooting (Target Project Users)

- If `/octo` still uses old behavior, reinstall plugin in user scope:
  - `/plugin uninstall octo@nyldn-plugins --scope user`
  - `/plugin install octo@nyldn-plugins --scope user`
- If `.multipowers/` is created in wrong location, run `/octo:init` from your target project root.
- If pre-run commands are not applied, check `.multipowers/context/runtime.json`:
  - `pre_run.enabled` must be `true`
  - at least one `entries[].match` must match your execution context (`all`, provider name, phase, role, or runtime tag)
  - `entries[].commands` must be non-empty
- If commands fail repeatedly, review `.multipowers/FAQ.md` for deduplicated root causes and fixes.
- If persona routing looks wrong, run `/octo:persona list` to confirm installed version and model mapping.
- If marketplace uninstall fails, confirm scope and plugin id (`octo@nyldn-plugins`).
