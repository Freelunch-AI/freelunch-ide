# Claude Code Guidance

Read and follow `AGENTS.md` for repository-wide context, safety boundaries, and review rules.

The five project skills under `.claude/skills/` are thin Claude Code adapters for the canonical workflows under `.agents/skills/`. Keep the canonical workflow authoritative and do not create additional public commands for internal capabilities.

Use `plan`, `implement`, `test`, `review`, and `pr` only when the developer explicitly invokes them. Never advance stages or perform Git mutations automatically.
