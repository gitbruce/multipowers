# Command Snapshot (2026-03-05)

## mp (current)

- checkpoint save|get|delete
- extract
- cost estimate|report
- init
- state get|set|update
- context guard
- validate (workspace|no-shell|tdd-env|test-run|coverage)
- plan
- discover|research
- define
- develop
- deliver|review
- embrace
- debate
- persona
- orchestrate select-agent
- loop
- hook
- status
- route
- test run
- coverage check
- config show-model-routing|get

## mp-devx (current)

- --action suite
- --action parity
- --action bench
- --action validate-sh-map
- --action build-policy
- --action build-runtime

## Target Split

- Keep runtime on `mp` (workflow/orchestration/state/context/hook).
- Move ops/testing/validation/reporting commands to `mp-devx`.
