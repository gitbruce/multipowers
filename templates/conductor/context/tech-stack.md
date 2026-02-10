# Tech Stack

## Runtime & Languages
- Shell: Bash
- Python: 3.x (connector and validation scripts)
- Node.js: npm test runner / plugin checks

## Key Components
- `bin/multipowers`: project CLI (init/doctor/track)
- `bin/ask-role`: role dispatch bridge
- `connectors/*.py`: external CLI adapters and structured logging
- `scripts/*.py`: validation and governance checks

## Tooling Requirements
- Required: `python3`
- Optional: `jq`
- Testing: `npm test` (core suite), optional integration tests

## Operational Conventions
- Structured logs in `outputs/runs/YYYY-MM-DD.jsonl`
- Context priority files:
  1. `conductor/context/product.md`
  2. `conductor/context/product-guidelines.md`
  3. `conductor/context/workflow.md`
  4. `conductor/context/tech-stack.md`
