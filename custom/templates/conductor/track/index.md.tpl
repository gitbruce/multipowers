# Track Index

- Intent: `intent.md`
- Design: `design.md`
- Implementation Plan: `implementation-plan.md`
- Metadata: `metadata.json`

## Track Summary
- Track ID: `{{.TrackID}}`
- Title: {{.TrackTitle}}
- Status: {{.Status}}
- Last Command: {{if .LastCommand}}{{.LastCommand}}{{else}}(none){{end}}
- Current Group: {{if .CurrentGroup}}{{.CurrentGroup}}{{else}}(not started){{end}}
