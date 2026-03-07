---
command: {{.Name}}
description: {{printf "%q" .Description}}
{{- if .SkillName }}
skill: {{.SkillName}}
{{- end }}
---

# /mp:{{.Name}}

{{.RuntimeSection}}
{{- if .UpstreamBody }}

## Upstream Workflow

{{.UpstreamBody}}
{{- end }}
