# Troubleshooting (Target Project Users)

- If `/octo` still uses old behavior, reinstall plugin in user scope:
  - `/plugin uninstall octo@nyldn-plugins --scope user`
  - `/plugin install octo@nyldn-plugins --scope user`
- If `conductor/` is created in wrong location, run `/octo:init` from your target project root.
- If persona routing looks wrong, run `/octo:persona list` to confirm installed version and model mapping.
- If marketplace uninstall fails, confirm scope and plugin id (`octo@nyldn-plugins`).
