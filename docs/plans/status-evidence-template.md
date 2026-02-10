# Status Evidence Template

Use this template whenever a task status changes to `DONE`.

## Task Status Fields

For each task section, keep both fields aligned:

- `- **Status**: \`TODO|IN_PROGRESS|BLOCKED|DONE\``
- `- **状态**：\`TODO|IN_PROGRESS|BLOCKED|DONE\``

## Evidence Block

- **Task ID**: `T?-???`
- **Coverage Task IDs**: `T?-???[, T?-??? ...]`
- **Date**: `YYYY-MM-DD`
- **Verifier**: `name/role`
- **Command(s)**:
  ```bash
  # exact verification command(s)
  ```
- **Exit Code**: `0`
- **Key Output**:
  - `short excerpt proving success`
- **Notes**:
  - `environment assumptions / known limitations`

## Governance Proof (for major changes)

When plan verification uses governance mode, include at least one of:

- `bash scripts/run_governance_checks.sh --mode strict ...`
- `semgrep`, `biome`, `ruff` execution evidence
- Governance artifact path (e.g. `outputs/governance/<artifact>.json`)

## Rule

A task marked `DONE` (by `Status` or `状态`) without evidence coverage is treated as `IN_PROGRESS`.

When strict governance evidence is required, missing governance proof fails evidence validation.
