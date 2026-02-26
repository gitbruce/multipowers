---
name: octopus-security-audit
description: |
  Comprehensive security auditing skill leveraging the security-auditor persona.
  Use for vulnerability scanning, OWASP compliance checks, and security reviews.
---

# Security Audit Skill

Invokes the security-auditor persona for thorough security analysis during the `ink` (deliver) phase.

## Usage

```bash
# Via mp runtime
./scripts/mp spawn security-auditor "Scan for SQL injection vulnerabilities"

# Via auto-routing (detects security intent)
./scripts/mp auto "security audit the payment processing module"
```

## Capabilities

- OWASP Top 10 vulnerability detection
- SQL injection and XSS scanning
- Authentication/authorization review
- Secrets and credential detection
- Dependency vulnerability assessment
- Security configuration review

## Persona Reference

This skill wraps the `security-auditor` persona defined in:
- `agents/personas/security-auditor.md`
- CLI: `claude-opus`
- Model: `claude-opus-4.6`
- Phases: `ink`
- Expertise: `owasp`, `vulnerability-scanning`, `security-review`

## Example Prompts

```
"Scan for hardcoded credentials in the codebase"
"Check for CSRF vulnerabilities in form handlers"
"Review the API authentication implementation"
"Analyze the encryption at rest configuration"
```
