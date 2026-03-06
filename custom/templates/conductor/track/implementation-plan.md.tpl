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

## Task Group {{.CurrentGroup}}
- Why: {{.Objective}}
- What: Render and maintain the canonical track artifacts for `{{.TrackID}}`.
- How: Use validated templates and write them to `.multipowers/tracks/{{.TrackID}}/`.
- Key Design: Treat the coordinator and metadata as the single source of truth for track state.
- Verification: {{.VerificationCommand}}
- Done When: {{.DoneWhen}}
