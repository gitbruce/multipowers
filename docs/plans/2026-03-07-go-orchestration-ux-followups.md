# Go Orchestration UX Follow-Ups

Date: 2026-03-07  
Status: Deferred after re-baseline

## Purpose

This note captures orchestration items that remain valuable, but are not required to close the current runtime-hardening gap. The current closure wave focuses on retry reliability, traceability, structured logs, regression goldens, and verification evidence.

## Deferred Items

### O06: Merge explainability CLI
- Build explain resolver output with value + source provenance.
- Expose `mp orchestrate explain`.
- Add negative tests and CLI reference examples.

### O07: Result caching
- Add cache key strategy and filesystem-backed cache store.
- Integrate read/write cache path into orchestration entrypoints.
- Add TTL and invalidation tests.

### O08: Planner visualization
- Export execution plan DAG.
- Render Mermaid output.
- Document visualization workflow and examples.

## Re-entry Criteria

Move these items back into an active implementation plan only when one of the following becomes true:

1. users need config precedence debugging in daily runtime work,
2. orchestration token/runtime cost becomes a dominant bottleneck,
3. planning complexity makes visual graph output a support requirement.
