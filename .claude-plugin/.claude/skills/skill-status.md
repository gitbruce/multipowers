# skill-status

Runtime status and health check skill.

## Overview

This skill provides comprehensive runtime status using the atomic mp status command
to check context, providers, validation, and hook readiness.

## Status Command

**Action:** Get comprehensive runtime status
```bash
mp status --dir . --json
```

## Response Fields

The status command returns comprehensive health information:

### Overall Status

- `status`: Overall runtime status ("ready" | "context_incomplete" | "no_providers" | "degraded")
- `data.ready`: Boolean indicating if runtime is ready for work

### Context Status

- `data.context_complete`: Boolean indicating if .multipowers context is complete
- `data.context_missing`: Array of missing context files
- `data.context_path`: Path to the project directory

### Provider Status

- `data.providers_available`: Array of available provider names
- `data.providers_count`: Number of available providers

### Validation Status

- `data.validation_status`: Last validation result ("passed" | "failed: <reason>")
- `data.last_validation`: Type of last validation performed

### Hook Status

- `data.hook_ready`: Boolean indicating if hooks are ready
- `data.hook_events`: Array of supported hook events

## Interpreting Status

### Ready State

```json
{
  "status": "ready",
  "data": {
    "ready": true,
    "context_complete": true,
    "providers_count": 2,
    "validation_status": "passed",
    "hook_ready": true
  }
}
```

**Action:** Runtime is ready for all workflows.

### Context Incomplete

```json
{
  "status": "context_incomplete",
  "data": {
    "ready": false,
    "context_complete": false,
    "context_missing": ["product.md", "tech.md"]
  }
}
```

**Action:** Run `/mp:init` to complete context setup.

### No Providers

```json
{
  "status": "no_providers",
  "data": {
    "ready": false,
    "providers_count": 0,
    "providers_available": []
  }
}
```

**Action:** Configure at least one AI provider (Claude, Codex, Gemini).

### Degraded State

```json
{
  "status": "degraded",
  "data": {
    "ready": false,
    "context_complete": true,
    "providers_count": 1,
    "validation_status": "failed: missing workspace"
  }
}
```

**Action:** Check specific failure and remediate.

## State Integration

Combine status with state for workflow tracking:

```bash
# Get current status
mp status --dir . --json

# Get workflow state
mp state get --dir . --json

# Update state with status
mp state update --data '{"last_status_check":"<timestamp>","runtime_ready":true}' --dir . --json
```

## Response Contract

The status command returns a JSON response with:
- `status`: Overall runtime status
- `message`: Human-readable status
- `data`: Structured status data (see above)

## Example Usage

User: "/mp:status"

1. Run status command
2. Parse response
3. Report health status to user
4. Provide remediation if not ready

User: "Check if runtime is ready for development"

1. Run status command
2. Check `data.ready` and `data.context_complete`
3. Check `data.providers_count` > 0
4. Report readiness status
