# Claude Octopus Plugin Configuration

## ⚠️ CRITICAL: Plugin Name

**The plugin name in `plugin.json` MUST remain `"mp"`**

```json
{
  "name": "mp"  // ⚠️ DO NOT CHANGE
}
```

### Why?

- Command prefix: `/mp:discover`, `/mp:debate`, etc.
- Changing this breaks all existing commands and user workflows
- Package name (`multipowers` in `package.json`) is different and correct

### More Information

- **Detailed explanation:** `PLUGIN_NAME_LOCK.md`
- **All safeguards:** `../SAFEGUARDS.md`
- **Validate:** Run `make test-plugin-name`

---

This directory contains the Claude Code plugin configuration for Claude Octopus.
