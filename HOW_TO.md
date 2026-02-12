# How to Use Multipowers

## Initial Setup

```bash
npm install
./bin/multipowers init
./bin/multipowers doctor
```

## Start Work on a Feature

```bash
# Create a new track
./bin/multipowers track new feature-name

# Start the track (sets active track)
./bin/multipowers track start feature-name
```

## Execute Tasks

### Fast Lane (Small Tasks)

Use for tiny fixes, config changes, narrow updates:

```bash
./bin/multipowers run --task "Fix typo in README"
```

### Standard Lane (Significant Work)

Use for features, refactors, multi-step changes:

```bash
./bin/multipowers run --task "Add user authentication API" --risk-hint high
```

### Manual Workflow Execution

```bash
# List available workflows
./bin/multipowers workflow list

# Run a specific workflow
./bin/multipowers workflow run subagent-driven-development --task "Implement feature X"
```

## Direct Role Dispatch

Skip router and dispatch directly to a role:

```bash
./bin/ask-role architect "Review the auth module design"
./bin/ask-role coder "Implement the login function with tests"
./bin/ask-role librarian "Find examples of rate limiting patterns"
```

## Finish Work

```bash
# Complete the track (runs governance checks if configured)
./bin/multipowers track complete feature-name
```

## Governance for Major Changes

When making significant modifications:

```bash
# Run governance checks
bash scripts/run_governance_checks.sh --mode strict --changed-file src/auth.py

# Or use npm script
npm run governance -- --mode strict --changed-file src/auth.py --changed-file docs/api.md
```

## Inspect State

```bash
# Check health and configuration
./bin/multipowers doctor

# View available workflows
./bin/multipowers workflow list

# Validate workflow configuration
./bin/multipowers workflow validate

# Check for updates
./bin/multipowers update --check
```

## View Execution Logs

Structured logs are written to:

```bash
outputs/runs/YYYY-MM-DD.jsonl
```

Filter by `request_id` or `track_id` to trace one execution timeline.

## Common Flags

| Flag | Purpose |
|------|---------|
| `--allow-untracked` | Execute without active track |
| `--risk-hint low\|medium\|high` | Influence lane selection |
| `--json` | Machine-readable output |
| `--track-id` | Associate execution with a track |

## Workflow Examples

```bash
# Brainstorm a feature design
./bin/multipowers workflow run brainstorming --task "Design payment flow"

# Write an implementation plan
./bin/multipowers workflow run writing-plans --task "Plan API refactor"

# Execute with code review gates
./bin/multipowers workflow run subagent-driven-development --task "Build user service"
```
