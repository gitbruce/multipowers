# Agent Decision Tree

> *When your brain says "which tentacle?!" - follow the flowchart.* 🐙

## Interactive Decision Trees

### By Development Phase

```mermaid
graph TD
    A[🐙 What do you need?] --> B{Phase}
    B -->|Research/Explore| C[🔍 PROBE Agents]
    B -->|Design/Plan| D[🎯 GRASP Agents]
    B -->|Build/Code| E[🔧 TANGLE Agents]
    B -->|Review/Ship| F[✅ INK Agents]

    C --> C1{What domain?}
    C1 -->|AI/LLM| C2[ai-engineer]
    C1 -->|Business/Metrics| C3[business-analyst]
    C1 -->|Context/Coordination| C4[context-manager]

    D --> D1{What are you designing?}
    D1 -->|Backend API| D2[backend-architect]
    D1 -->|Frontend UI| D3[frontend-developer]
    D1 -->|Database| D4[database-architect]
    D1 -->|Cloud/Infra| D5[cloud-architect]
    D1 -->|GraphQL| D6[graphql-architect]

    E --> E1{What are you building?}
    E1 -->|With tests| E2[tdd-orchestrator]
    E1 -->|Python code| E3[python-pro]
    E1 -->|TypeScript code| E4[typescript-pro]
    E1 -->|Debugging| E5[debugger]
    E1 -->|DevOps| E6[devops-troubleshooter]

    F --> F1{What kind of review?}
    F1 -->|Code quality| F2[code-reviewer]
    F1 -->|Security| F3[security-auditor]
    F1 -->|Performance| F4[performance-engineer]
    F1 -->|Testing| F5[test-automator]
    F1 -->|Deployment| F6[deployment-engineer]

    style A fill:#9b59b6,stroke:#8e44ad,color:#fff
    style C fill:#3498db,stroke:#2980b9,color:#fff
    style D fill:#e74c3c,stroke:#c0392b,color:#fff
    style E fill:#f39c12,stroke:#d68910,color:#fff
    style F fill:#27ae60,stroke:#1e8449,color:#fff
```

---

### By Task Type

```mermaid
graph LR
    A[🦑 Task Type] --> B{What's the primary task?}

    B -->|Debug/Troubleshoot| C[debugger]
    B -->|Review Code| D[code-reviewer]
    B -->|Security Audit| E[security-auditor]
    B -->|Optimize Performance| F[performance-engineer]
    B -->|Write Tests| G[test-automator]
    B -->|Design Architecture| H{Architecture Type}
    B -->|Deploy/CI/CD| I[deployment-engineer]
    B -->|Document| J[docs-architect]

    H -->|Backend| K[backend-architect]
    H -->|Frontend| L[frontend-developer]
    H -->|Database| M[database-architect]
    H -->|Cloud| N[cloud-architect]
    H -->|GraphQL| O[graphql-architect]

    style A fill:#9b59b6,stroke:#8e44ad,color:#fff
    style H fill:#e74c3c,stroke:#c0392b,color:#fff
```

---

### By Technology Stack

```mermaid
graph TD
    A[💻 Technology] --> B{Primary Language/Framework}

    B -->|Python| C[python-pro]
    B -->|TypeScript/JavaScript| D[typescript-pro]
    B -->|React/Next.js| E[frontend-developer]
    B -->|GraphQL| F[graphql-architect]
    B -->|AWS/GCP/Azure| G[cloud-architect]
    B -->|Kubernetes/Docker| H[deployment-engineer]
    B -->|Any/Multiple| I{What's the task?}

    I -->|Architecture| J[backend-architect]
    I -->|Security| K[security-auditor]
    I -->|Testing| L[test-automator]
    I -->|Performance| M[performance-engineer]

    style A fill:#9b59b6,stroke:#8e44ad,color:#fff
    style I fill:#f39c12,stroke:#d68910,color:#fff
```

---

## Text-Based Quick Reference

### The 3-Question Method

**Question 1: What phase are you in?**

```
Research/Explore → PROBE tentacles
Design/Plan     → GRASP tentacles
Build/Code      → TANGLE tentacles
Review/Ship     → INK tentacles
```

**Question 2: What's your domain?**

```
Backend API     → backend-architect
Frontend UI     → frontend-developer
Database        → database-architect
Cloud/Infra     → cloud-architect
Security        → security-auditor
Testing         → tdd-orchestrator / test-automator
Performance     → performance-engineer
```

**Question 3: What's your language?**

```
Python          → python-pro
TypeScript/JS   → typescript-pro
Multiple/Any    → Use domain-specific tentacle
```

---

## Common Scenarios

### "I'm building a new feature"

```
Is it backend?
├─ Yes → Does it involve database schema?
│        ├─ Yes → database-architect FIRST, then backend-architect
│        └─ No  → backend-architect
└─ No → Is it frontend?
         ├─ Yes → frontend-developer
         └─ No  → What is it? (Describe and auto-route)
```

### "I need to fix something"

```
Is it a bug/error?
├─ Yes → debugger
│        └─ Still stuck? → devops-troubleshooter (if infra)
└─ No → Is it slow?
         ├─ Yes → performance-engineer
         │        └─ Database slow? → database-architect
         └─ No  → Is it insecure?
                  ├─ Yes → security-auditor
                  └─ No  → code-reviewer (quality issues)
```

### "I need a review"

```
What kind of review?
├─ General quality     → code-reviewer
├─ Security focused    → security-auditor
├─ Performance focused → performance-engineer
└─ Multiple concerns   → Use /mp:review skill
```

---

## The "Just Tell Me" Cheat Sheet

| If you're thinking... | Use this tentacle |
|-----------------------|-------------------|
| "I need to design an API" | `backend-architect` |
| "Something's broken" | `debugger` |
| "Is this secure?" | `security-auditor` |
| "Why is it slow?" | `performance-engineer` |
| "Review my code" | `code-reviewer` |
| "I want TDD" | `tdd-orchestrator` |
| "Design the database" | `database-architect` |
| "React/Next.js work" | `frontend-developer` |
| "Python code" | `python-pro` |
| "TypeScript types" | `typescript-pro` |
| "Cloud infrastructure" | `cloud-architect` |
| "GraphQL schema" | `graphql-architect` |
| "CI/CD pipeline" | `deployment-engineer` |
| "Production incident" | `incident-responder` |
| "Write documentation" | `docs-architect` |
| "Create a diagram" | `mermaid-expert` |

---

## When In Doubt

**Just describe what you need!** Claude Octopus auto-routes based on keywords:

```
"Build user authentication with OAuth"
→ Auto-routes to: backend-architect + database-architect

"Review this code for security issues"
→ Auto-routes to: security-auditor

"My API is slow"
→ Auto-routes to: performance-engineer
```

**Or use a workflow skill:**

| Skill | Does What |
|-------|-----------|
| `/mp:review` | Fast code review (define + develop) |
| `/mp:research` | Deep research (4-perspective discover) |
| `/mp:security` | Security audit (red team/blue team) |
| `/mp:embrace` | Full 4-phase Double Diamond workflow |

---

## Anti-Patterns: Don't Do This

| Bad Idea | Why | Do This Instead |
|----------|-----|-----------------|
| Use `security-auditor` for code style | Wrong domain | Use `code-reviewer` |
| Use `debugger` for architecture | Wrong phase | Use `backend-architect` |
| Use `database-architect` for API design | Wrong order | Design schema first, then API |
| Invoke principles agents directly | Internal use only | Use the persona that applies them |
| Use opus for simple tasks | Wastes tokens | Let auto-routing pick appropriate tier |

---

<p align="center">
  🐙 <em>"Eight tentacles, one goal: picking the right one for YOU."</em> 🐙
</p>
