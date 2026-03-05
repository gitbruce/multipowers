# Agent Decision Tree

> *When your brain says "which tentacle?!" - follow the flowchart.* 🐙

## Interactive Decision Trees

### By Development Phase (Double Diamond)

```mermaid
graph TD
    A[🐙 What do you need?] --> B{Phase}
    B -->|Research/Explore| C[🔍 Discover Agents]
    B -->|Design/Plan| D[🎯 Define Agents]
    B -->|Build/Code| E[🔧 Develop Agents]
    B -->|Review/Ship| F[✅ Deliver Agents]

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

## Text-Based Quick Reference

### The 3-Question Method

**Question 1: What phase are you in?**

```
Research/Explore → Discover tentacles
Design/Plan     → Define tentacles
Build/Code      → Develop tentacles
Review/Ship     → Deliver tentacles
```

---

## The "Just Tell Me" Cheat Sheet

| If you're thinking... | Use this tentacle |
|-----------------------|-------------------|
| "I need to design an API" | `backend-architect` |
| "Something's broken" | `debugger` |
| "Is this secure?" | `security-auditor` |
| "Review my code" | `code-reviewer` |
| "I want TDD" | `tdd-orchestrator` |

---

## When In Doubt

**Just describe what you need!** Multipowers auto-routes based on keywords:

**Or use a workflow skill:**

| Skill | Does What |
|-------|-----------|
| `/mp:review` | Fast code review (deliver workflow) |
| `/mp:discover` | Deep research (multi-perspective discovery) |
| `/mp:embrace` | Full 4-phase Double Diamond workflow |

---

<p align="center">
  🐙 <em>"Eight tentacles, one goal: picking the right one for YOU."</em> 🐙
</p>
