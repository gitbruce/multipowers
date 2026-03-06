# Design

- Track ID: `{{.TrackID}}`
- Title: {{.TrackTitle}}

## Architecture
- Centralize track lifecycle state in shared coordinator and metadata helpers.
- Keep canonical artifacts under `.multipowers/tracks/{{.TrackID}}/`.

## Scope Mapping
- Objective: {{.Objective}}
- Current Group: {{.CurrentGroup}}
- Status: {{.Status}}

## Operational Notes
- Execution mode: {{.ExecutionMode}}
- Completed groups: {{join .CompletedGroups ", "}}
