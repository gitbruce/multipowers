# {{PROJECT_NAME}} - Claude Code Working Agreement

## 0. Scope and Priority
- Scope: this file applies to this target project only.
- Priority: system instructions > command contracts > this file > task notes.
- Output boundary: all generated artifacts must stay inside this target project, primarily under `.multipowers/`.

## 1. Project Overview
- Project: {{PROJECT_NAME}}
- Summary: {{PRODUCT_SUMMARY}}
- Target users: {{TARGET_USERS}}
- Primary goal: {{PRIMARY_GOAL}}
- Non-goals: {{NON_GOALS}}
- Constraints: {{CONSTRAINTS}}

## 2. Language-Agnostic Engineering Rules
- Prefer explicit typing where the language supports it.
- If the language is dynamic, validate all external inputs at boundaries.
- Keep functions focused and small.
- Public interfaces require concise usage docs.
- Avoid hardcoded environment paths, ports, or credentials.

## 3. Security Baseline
- Never hardcode secrets, tokens, or private keys.
- Never log sensitive user or business data.
- Validate and sanitize all untrusted inputs.
- Use safe query/command execution patterns to prevent injection.

## 4. Testing and Quality Gates
- New behavior requires tests.
- Critical flows require integration coverage.
- Before completion: format, lint, test, and build/type-check must pass.

## 5. Runtime Preconditions
- Runtime profile: {{RUNTIME}}
- Framework: {{FRAMEWORK}}
- Database: {{DATABASE}}
- Deployment target: {{DEPLOYMENT}}
- Runtime precommands are defined in `.multipowers/context/runtime.json`.
- Policy: fail-fast. If any precommand fails, stop execution immediately.

## 6. Workflow
{{WORKFLOW_FLOW}}

## 7. Learned Failure Avoidance
Read `.multipowers/FAQ.md` before executing high-impact commands.
`FAQ.md` is auto-generated and auto-refined from observed failures.

## 8. Project-Specific Notes
- Add project-specific constraints here.
