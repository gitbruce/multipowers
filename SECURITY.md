# Security Policy for Multipowers

## Threat Model (Go Native Engine)

Multipowers orchestrates external AI CLI tools using a Go-native atomic engine. The threat landscape has shifted from shell-scripting vulnerabilities to runtime execution boundaries and secret management.

### Trust Boundaries

| Boundary | Description | Risk Level | Mitigation |
|----------|-------------|------------|------------|
| User Prompts | Natural language inputs | Medium | LLM output validation and prompt shielding. |
| File System | Parallel execution writes | High | **FSBoundary** & **Worktree Isolation**. |
| API Keys | Provider credentials | High | Environment-only storage & mask-on-render. |
| External CLIs | Responses from AI models | Low | JSON schema validation and retry quotas. |

## Security Controls

### 1. Physical Isolation (FSBoundary)
Multipowers enforces strict file system boundaries within the Go runtime (`internal/fsboundary`).
- **Path Guard**: Commands are physically blocked from reading or writing files outside the explicitly authorized Project Root.
- **Symlink Protection**: Resolves all paths to absolute physical locations to prevent symlink-based escapes.

### 2. Execution Safety
- **Parameterized Execution**: legacy Shell string concatenation is replaced with `os/exec` slice-based arguments to eliminate Shell Injection.
- **Bounded Concurrency**: `Worktree Slots` prevent resource exhaustion by capping active parallel agents.

### 3. Secrets Management
- **Key formats**: Automatic regex validation for OpenAI (`sk-...`) and Google API keys.
- **Masking**: All logs and verbose outputs automatically mask sensitive strings with `***`.
- **Zero-Storage Policy**: API keys are never written to disk or session state files.

### 4. Audit & Logging
Security-relevant events are tracked via standard runtime logging. Advanced audit logging to `~/.multipowers/logs/audit.log` is currently in the architectural roadmap.

---

## Reporting Vulnerabilities

**Please DO NOT create public GitHub issues for security vulnerabilities.**

To report a vulnerability, please use **GitHub Security Advisories** to create a private report. We aim to acknowledge all reports within 48 hours.

---

**Last Updated:** March 2026
**Status:** Multipowers v8.x Security Standard
