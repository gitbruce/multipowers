# Superpowers vs. Multipowers Functional Mapping & Analysis

This document provides a comparative analysis between the [Superpowers](https://github.com/obra/superpowers) framework and the **Multipowers** (Go-native branch) project. It maps functionalities and evaluates the architectural trade-offs of each approach.

## 1. Core Architectural Philosophy

| Feature | Superpowers | Multipowers (Go Branch) |
| :--- | :--- | :--- |
| **Primary Engine** | **Markdown Behavioral Programming**: The agent reads `SKILL.md` files to understand its own constraints and state machine. | **Go-Native Hybrid Engine**: A compiled Go binary (`mp`) handles orchestration, isolation, and state, while Markdown provides high-level reasoning. |
| **Logic Implementation** | Natural language instructions within Markdown. | Type-safe Go code under `internal/` with structured JSON contracts. |
| **State Management** | Implicit (in-context) and file-based (Markdown checkmarks). | Explicit (Persistent KV store, `metadata.json`, and JSONL logs). |
| **Enforcement** | Cognitive (Agent self-policing based on instructions). | Runtime (Go engine physically blocks execution if gates fail). |

---

## 2. Command & Skill Mapping

### 2.1 Planning & Design
| Superpowers Skill | Multipowers Command/Logic | Mapping Notes |
| :--- | :--- | :--- |
| `brainstorming` | `/mp:define` / `/mp:research` | Superpowers uses a Socratic "one question at a time" approach. Multipowers uses specialized Personas (`backend-architect`). |
| `writing-plans` | `/mp:plan` | Both generate structured Markdown implementation plans. Multipowers enforces a `complexity_score` during this phase. |

### 2.2 implementation & Execution
| Superpowers Skill | Multipowers Command/Logic | Mapping Notes |
| :--- | :--- | :--- |
| `executing-plans` | `/mp:develop` | Multipowers adds a concurrent execution engine with progress tracking. |
| `subagent-driven-dev` | `internal/orchestration` | Multipowers uses a Planner-Executor-Synthesizer pattern to manage sub-agents. |
| `tdd` | `internal/validation` | Multipowers enforces TDD via pre-run hooks and `doctor` checks. |
| `using-git-worktrees` | `internal/isolation` | Multipowers automatically calculates if a worktree is required based on complexity. |

### 2.3 Quality & Delivery
| Superpowers Skill | Multipowers Command/Logic | Mapping Notes |
| :--- | :--- | :--- |
| `requesting-code-review` | `/mp:review` | Both check code against the plan. Multipowers supports "Reviewer-Led Abort" in parallel flows. |
| `finishing-a-dev-branch` | `/mp:deliver` / `/mp:embrace` | Both handle merging and cleanup. Multipowers includes automated `Synthesis` logic. |
| `systematic-debugging` | `/mp:doctor` | Superpowers focuses on RCA steps; Multipowers focuses on automated environment/code health checks. |

---

## 3. Pros and Cons Analysis

### **Superpowers (Markdown-Centric)**

#### **Pros:**
- **Extreme Flexibility**: Easy to modify the agent's behavior by simply editing text. No compilation required.
- **Portability**: Works in any environment that can read Markdown (Claude Code, Cursor, etc.).
- **Human-Readable Logic**: The "source code" of the process is the documentation itself.
- **Agent Cognition**: Forces the agent to "think" about the process by keeping it in its prompt context.

#### **Cons:**
- **Context Pressure**: Large `SKILL.md` files consume significant token space in every turn.
- **Non-Deterministic**: The agent might occasionally ignore instructions ("hallucinate") or bypass gates.
- **Performance**: Sequential reasoning is slower than compiled logic for structural tasks (like worktree management).
- **Weak State**: Hard to maintain complex state across long-running parallel tasks without explicit database/KV support.

### **Multipowers (Go-Native Hybrid)**

#### **Pros:**
- **Deterministic Enforcement**: Hard gates (like `Complexity Scoring`) are physical blocks in Go, making the process tamper-proof.
- **Context Efficiency**: Heavy logic is offloaded to the Go binary; the agent only gets concise JSON results, saving tokens.
- **Concurrency**: Built-in support for parallel agents, mailbox-based IPC, and slot-limited worktrees.
- **Scalability**: Capable of handling enterprise-scale refactors with thousands of files using compiled performance.

#### **Cons:**
- **Higher Friction for Customization**: Changing the core engine requires Go knowledge and a rebuild/deploy cycle.
- **Deployment Overhead**: Requires the user to have the Go binary installed or bundled in the plugin.
- **Abstracted Logic**: The "why" of a block might be hidden inside Go code rather than being visible in the Markdown skill.

---

## 4. Conclusion & Evolution

**Superpowers** is an excellent framework for **cognitive guidance**—ensuring the agent follows a high-quality human process. 

**Multipowers** is an evolution into a **production-grade AI OS**. It takes the "best practices" defined by frameworks like Superpowers and encodes them into a robust, high-performance runtime. 

**Recommended Path:** Use Superpowers-style Markdown for **strategy and reasoning** (the "Brain") and Multipowers-style Go logic for **orchestration and safety** (the "Muscles").
