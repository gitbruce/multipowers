# Go Migration Release Notes (9.0.0-go)

- Runtime switched to Go single binary (`bin/octo`).
- `scripts/orchestrate.sh` now wrapper with fallback (`OCTO_RUNTIME=legacy`).
- Hooks route through `octo hook --event ...`.
- Spec-driven guard enforced via Go pipeline.
- Artifacts constrained to target `/.multipowers/*`.
