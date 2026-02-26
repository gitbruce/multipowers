# flow-deliver

This skill is a thin Go wrapper.

Execute:
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/octo" deliver --dir "$PWD" --prompt "<user-prompt>" --json
```

Rules:
- Parse JSON response only.
- If `status` is `error` or `blocked`, stop immediately.
- Do not perform direct implementation logic in this markdown skill.
