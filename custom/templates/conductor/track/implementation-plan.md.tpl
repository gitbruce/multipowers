# Implementation Plan

- Track ID: `{{.TrackID}}`
- Title: {{.TrackTitle}}
- Status: {{.Status}}

## Execution Mode Decision
- Execution Mode: {{.ExecutionMode}}
- Complexity Score: {{.ComplexityScore}}
- Worktree Required: {{.WorktreeRequired}}
- Rationale: {{.ExecutionRationale}}
- Enforcement: update metadata + update tracks registry + commit after each group

## Execution State
- Last Command: {{if .LastCommand}}{{.LastCommand}}{{else}}(none){{end}}
- Last Command At: {{if .LastCommandAt}}{{.LastCommandAt}}{{else}}(unknown){{end}}
- Current Group: {{if .CurrentGroup}}{{.CurrentGroup}}{{else}}(not started){{end}}
- Group Status: {{if .GroupStatus}}{{.GroupStatus}}{{else}}(idle){{end}}
- Completed Groups: {{if .CompletedGroups}}{{join .CompletedGroups ", "}}{{else}}(none){{end}}

## Task Group {{if .CurrentGroup}}{{.CurrentGroup}}{{else}}(not started){{end}}
- Why: {{.Objective}}
- What: Render and maintain the canonical track artifacts for `{{.TrackID}}`.
- How: Use validated templates and write them to `.multipowers/tracks/{{.TrackID}}/`.
- Key Design: Treat the coordinator and metadata as the single source of truth for track state.
- Verification: {{.VerificationCommand}}
- Done When: {{.DoneWhen}}
