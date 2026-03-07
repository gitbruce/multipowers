# Intent

- Track ID: `{{.TrackID}}`
- Title: {{.TrackTitle}}
- Objective: {{.Objective}}
- Status: {{.Status}}
- Last Command: {{if .LastCommand}}{{.LastCommand}}{{else}}(none){{end}}
- Current Group: {{if .CurrentGroup}}{{.CurrentGroup}}{{else}}(not started){{end}}
- Completed groups so far: {{if .CompletedGroups}}{{join .CompletedGroups ", "}}{{else}}(none){{end}}

## Success Criteria
- Keep all canonical artifacts under `.multipowers/tracks/{{.TrackID}}/`.
- Preserve canonical registry updates in `.multipowers/tracks/tracks.md`.
- Complete the current group {{if .CurrentGroup}}`{{.CurrentGroup}}`{{else}}after it is explicitly started{{end}}.
