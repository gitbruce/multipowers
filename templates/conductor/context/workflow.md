# Workflow

## Standard Delivery Loop
1. Clarify requirement and scope.
2. Create/Update plan under `docs/plans/`.
3. Implement in small increments with tests.
4. Run verification (`bash -n`, `py_compile`, `npm test`).
5. Update evidence section before marking `DONE`.

## Track Lifecycle
1. `./bin/multipowers track new <feature-name>`
2. `./bin/multipowers track start <track-name>`
3. Deliver plan tasks with evidence
4. `./bin/multipowers track complete <track-name>`

## Artifacts
- Design docs: `docs/design/`
- Plans: `docs/plans/`
- Track records: `conductor/tracks/`
