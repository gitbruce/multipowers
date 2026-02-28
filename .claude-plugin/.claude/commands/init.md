---
command: init
description: Guided wizard setup for .multipowers context
---

# /mp:init

Run a guided wizard first, then call Go runtime.

Actions:
1. Ask one batched `AskUserQuestion` with these fields:
   - `project_name`
   - `summary`
   - `target_users`
   - `primary_goal`
   - `non_goals`
   - `constraints`
   - `runtime`
   - `framework`
   - `database`
   - `deployment`
   - `workflow`
   - `track_name`
   - `track_objective`
2. Build a single-line JSON object from answers and pass it as `--prompt`.
3. Ensure `${CLAUDE_PLUGIN_ROOT}/bin/mp` exists.
4. Execute:
```bash
"${CLAUDE_PLUGIN_ROOT}/bin/mp" init --dir "$PWD" --prompt "$INIT_WIZARD_JSON" --json
```
5. Parse JSON response.
6. If `status` is `error` or `blocked`, stop immediately.

Do not implement command logic in markdown.
