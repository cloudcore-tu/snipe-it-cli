# doc-sync

Check whether code changes require documentation updates.

## Check Targets

- `docs/agents/project.md`
- `.claude/CLAUDE.md`
- `AGENTS.md`
- `docs/adr/`
- `CHANGELOG.md`
- `README.md`
- `docs/api-coverage.md`

## Workflow

1. Inspect the diff.
2. Classify changes: dependency, configuration, architecture, command surface, user-visible behavior.
3. Verify matching docs were updated.

## Rules

- Update `docs/agents/` first for shared agent guidance.
- Keep `.claude/CLAUDE.md`, `AGENTS.md`, and tool-specific skill wrappers thin.
- If command surface changes, check both `README.md` and `docs/api-coverage.md`.
