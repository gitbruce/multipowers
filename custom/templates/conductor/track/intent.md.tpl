# Intent

- Track ID: `{{.TrackID}}`
- Title: {{.TrackTitle}}
- Status: {{.Status}}

## Objective
{{.Objective}}

## Success Criteria
- Complete the current group `{{.CurrentGroup}}`.
- Preserve canonical artifacts under `.multipowers/tracks/{{.TrackID}}/`.

## Constraints
- Completed groups so far: {{join .CompletedGroups ", "}}
- Execution mode stays aligned with the generated implementation plan.
