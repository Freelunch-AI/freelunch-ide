# FreeLunch AI Workflow Bundle

This repository includes a recommended, model-agnostic workflow for Codex, Claude Code, Cursor, and OpenCode. It exposes five entry points and keeps source-derived procedures, specialist lenses, and optional tool adapters behind them.

The bundle is optional. Developers remain responsible for their changes, and CI/CD remains responsible for enforced quality gates.

## Entry Points

| Workflow                      | Codex               | Claude Code         | Cursor              | OpenCode            |
| ----------------------------- | ------------------- | ------------------- | ------------------- | ------------------- |
| Plan an issue without changes | `$plan <issue-url>` | `/plan <issue-url>` | `/plan <issue-url>` | `/plan <issue-url>` |
| Implement an approved plan    | `$implement`        | `/implement`        | `/implement`        | `/implement`        |
| Add and run relevant tests    | `$test`             | `/test`             | `/test`             | `/test`             |
| Run an independent review     | `$review`           | `/review`           | `/review`           | `/review`           |
| Open a linked draft PR        | `$pr`               | `/pr`               | `/pr`               | `/pr`               |

Codex uses `$name` for explicit Agent Skill invocation. The other supported clients expose slash invocation. Names, inputs, safety rules, state, and outputs remain the same.

## Setup

Use a supported client from the repository root. Project configuration is checked in:

- Codex reads `AGENTS.md` and `.agents/skills/`.
- Claude Code reads `CLAUDE.md` and `.claude/skills/`.
- Cursor reads `.cursor/rules/` and `.cursor/commands/`.
- OpenCode reads `AGENTS.md` and `.opencode/commands/`.

The workflow needs issue and pull request context. Use an existing GitHub connector or MCP integration, or authenticate the GitHub CLI with `gh auth login`. No credentials or provider-specific MCP configuration are committed to this repository.

## Layout

```text
.agents/skills/          canonical five public workflows
.claude/skills/          Claude Code Markdown entry points
.cursor/commands/        Cursor Markdown entry points
.opencode/commands/      OpenCode Markdown entry points
.freelunch/ai/           internal capabilities, source policy, and progress template
.freelunch/progress/     local, ignored per-issue progress records
ai-stack.lock.json       pinned methodology and capability sources
```

Each client entry point is checked in directly as Markdown and points to the canonical workflow. There is no generator, transpiler, build step, or language runtime. Internal capability files are outside client discovery paths, so they do not appear as extra commands.

## Progress And Safety

`plan` is read-only. After the developer approves branch creation, `implement` creates `.freelunch/progress/issue-<number>.md` from the checked-in template. Later workflows update the same record so work can resume in another session or supported client.

Branch creation, commits, pushes, rebases, merges, and remote GitHub changes require explicit approval. `pr` creates a draft pull request and stops; it never merges.

## Dependencies

`ai-stack.lock.json` pins the exact BMAD Method, Agency Agents, GBrain, Graphify, RTK, and official Anthropic plugin revisions reviewed while designing the internal capabilities. It records licenses, selected files, adopted behavior, and deliberate exclusions. `.freelunch/ai/sources.md` explains how each source maps into the five workflows.

No third-party prompt package is added to a public discovery directory. Only original FreeLunch adaptations are checked in. The Anthropic plugin repository is reference-only because its content is all rights reserved under Anthropic's commercial terms.

GBrain, Graphify, and RTK are optional. When already available, the workflows may use read-only knowledge lookup, an existing code graph, or selective output compression. The repository does not install their runtimes, enable hooks, create graphs, configure accounts, or make them quality-gate dependencies.
