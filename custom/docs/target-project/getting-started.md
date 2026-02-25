# Getting Started (Target Project Users)

## Install Plugin (User Scope)

```text
/plugin marketplace add /mnt/f/src/ai/claude-octopus
/plugin install octo@nyldn-plugins --scope user
```

## Initialize in Your Project

In your target project directory:

```text
/octo:init
```

Expected:
- creates `conductor/` in your target project
- initializes project context files and tracks registry

## Run Spec-Driven Commands

- `/octo:plan`
- `/octo:discover`, `/octo:define`, `/octo:develop`, `/octo:deliver`
- `/octo:embrace`, `/octo:review`, `/octo:debate`, `/octo:research`

If context is missing, `/octo:init` is auto-triggered.
