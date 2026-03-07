---
name: {{.Name}}
description: {{printf "%q" .Description}}
---

# {{.Name}}

{{.RuntimeSection}}
{{- if .UpstreamBody }}

## Upstream Workflow

{{.UpstreamBody}}
{{- end }}
