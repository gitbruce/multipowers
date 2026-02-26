# Config Schema

- `models.json`: providers, role_routing, fallback_lane
- `proxy.json`: enabled, providers, host, port, no_proxy
- `persona-lanes.json`: personas, fallback_lane
- `setup.toml`: `/octo:init` setup protocol and template source paths

## Generated Target-Project Artifacts

- `.multipowers/CLAUDE.md`: project-level working agreement generated from `custom/templates/CLAUDE.md`
- `.multipowers/FAQ.md`: auto-generated, deduplicated failure-avoidance knowledge
- `.multipowers/context/runtime.json`: runtime pre-run contract (`fail-fast` policy)
