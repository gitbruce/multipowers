# Go Branch Upstream Diff Discipline

- Keep upstream-heavy files minimal on go branch where possible.
- Move governance behavior into Go runtime (`cmd/`, `internal/`, `pkg/`).
- Keep markdown commands/skills as thin wrappers only.
- Record any unavoidable upstream conflict files in PR notes.
