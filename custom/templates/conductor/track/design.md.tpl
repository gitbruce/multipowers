# Design

## Track Overview
- Track ID: `{{.TrackID}}`
- Title: {{.TrackTitle}}
- Status: {{.Status}}
- Last Command: {{if .LastCommand}}{{.LastCommand}}{{else}}(none){{end}}
- Current Group: {{if .CurrentGroup}}{{.CurrentGroup}}{{else}}(not started){{end}}
- Execution mode: {{.ExecutionMode}}
- Completed groups: {{if .CompletedGroups}}{{join .CompletedGroups ", "}}{{else}}(none){{end}}

## Scope Mapping
- Intent Objective: {{.Objective}}
- Key Design: Centralize track lifecycle state in coordinator + metadata.
