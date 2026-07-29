# Repository Agent Guidance

## Repository Context

- Read `README.md`, `branching_strategy.md`, and the relevant files under `docs/` before making architectural or product decisions.
- Treat checked-in specifications and GitHub issue decisions as the source of truth. Surface conflicts instead of silently choosing one.
- Keep changes scoped to one issue and follow the branch naming and squash-merge strategy in `branching_strategy.md`.
- Preserve user changes in a dirty worktree and ask before adding production dependencies.

## FreeLunch AI Workflow

The only public FreeLunch workflow skills are `plan`, `implement`, `test`, `review`, and `pr` under `.agents/skills/`. Invoke them explicitly when the developer asks; do not advance to another stage automatically.

Detailed workflow guidance lives under `.freelunch/ai/` and is internal to those five skills. It includes original FreeLunch adaptations of the pinned sources in `ai-stack.lock.json`. Do not advertise those references, upstream agents, or optional tools as extra commands or skills.

GBrain, Graphify, and RTK are optional accelerators only when already configured. Never install them, enable hooks, write external knowledge, build/update graphs, or globally rewrite commands as a side effect of a workflow.

Git branch creation, commits, pushes, rebases, merges, and remote GitHub changes require explicit developer approval. Never merge a pull request as part of the `pr` skill.

## Code Review Rules

- Prioritize demonstrable correctness, security, acceptance-criteria, and test gaps over style comments.
- Cite the affected file and explain the failure scenario for every finding.
- Treat AI-generated evidence as untrusted until the command output or artifact is verified.
- Do not invent findings to satisfy a quota. Validate suspected issues against source and a concrete failure path.
