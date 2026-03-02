# Release Notes: No-Shell Hybrid Runtime

## Summary

This release implements a hybrid architecture: Go atomic engine + Markdown reasoning skills.

## Architecture

- **Go Runtime (Engine)**: Provides deterministic atomic commands
  - `mp state get/set/update` - State management
  - `mp validate --type <type>` - Typed validation (workspace, no-shell, tdd-env, test-run, coverage)
  - `mp route --intent <intent>` - Provider routing
  - `mp test run` - Test execution
  - `mp coverage check` - Coverage analysis
  - `mp status` - Runtime health check
  - `mp hook --event <event>` - Hook handling

- **Markdown Skills (Brain)**: Provide stepwise reasoning and orchestration
  - Call atomic `mp` commands instead of shell scripts
  - Consume JSON outputs to decide next actions
  - Preserve LLM decision space for non-deterministic reasoning

## JSON Contract

All atomic commands return a normalized response:
```json
{
  "status": "ok|error|blocked",
  "action": "continue|ask_user_questions",
  "error_code": "E_XXX",
  "message": "Human-readable status",
  "data": { ... },
  "remediation": "Suggested fix if blocked"
}
```

## Highlights

- All runtime shell scripts removed from repository
- `/mp:*` commands execute through `.claude-plugin/bin/mp`
- Core skills (flow-*.md, skill-tdd/validate/status.md) restored with reasoning content
- Legacy high-level commands (discover, define, develop, deliver) are compatibility facades

## Verification

- `go test ./...` pass
- `go vet ./...` pass
- `mp validate --type no-shell --dir . --json` pass
- `mp status --dir . --json` returns real runtime health

## Notes

- If older docs mention shell paths, treat them as historical references only
- Runtime and CI should use `.claude-plugin/bin/mp` and `go run ./cmd/mp-devx`
- Skills should call atomic `mp` commands and branch based on JSON responses
