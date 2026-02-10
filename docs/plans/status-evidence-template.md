# Status Evidence Template

Use this template whenever a task status changes to `DONE`.

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

## Rule

A task marked `DONE` without evidence coverage is considered `IN_PROGRESS`.
