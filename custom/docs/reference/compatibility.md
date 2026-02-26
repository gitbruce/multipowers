# Compatibility Matrix

| Upstream Base | Multipowers Overlay | Required Action | Known Caveats |
|---|---|---|---|
| v8.23.1 (`f6a815a`) | overlay-v1 | Apply overlay scripts and run custom tests | Requires Node runtime for JSON parsing helpers |
| future upstream patch/minor | overlay-v1+ | Run `custom/scripts/octo-devx sync`, then validate | Conflicts possible in high-churn command/docs files |

## Verification Scope
- Command registration includes `/octo:persona`
- Overlay config files parse correctly
- Proxy hook applies only to external providers
