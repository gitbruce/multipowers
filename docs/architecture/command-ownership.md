# MP / MP-DEVX Command Ownership

日期：2026-03-05

## Boundary Rule

- `mp`: runtime execution path for users/agents.
- `mp-devx`: operations/devx path for maintainers/CI.

## Runtime (`mp`)

- `init`
- `context guard`
- `state get|set|update`
- `plan|discover|research|define|develop|deliver|review|embrace|debate|persona`
- `orchestrate select-agent`
- `loop`
- `hook`
- `status`
- `route`
- `extract`
- `checkpoint save|get|delete`
- `cost estimate`
- `config get|show-model-routing`

## Ops/Devx (`mp-devx`)

- `--action suite`
- `--action parity`
- `--action bench`
- `--action validate-sh-map`
- `--action build-policy`
- `--action build-runtime`

## Planned Migrations (mp -> mp-devx)

- `mp test run`
- `mp coverage check`
- `mp validate --type no-shell` / `--strict-no-shell`
- `mp cost report`

## Compatibility Policy

- Phase 1: `mp` keeps command but returns deprecation guidance.
- Phase 2: `mp` removes execution and keeps stable migration error message.
