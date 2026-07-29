# FreeLunch AI Workflow Contract

Apply this contract whenever a FreeLunch workflow skill runs. Read `.freelunch/ai/sources.md` when using source-derived capabilities.

## Boundaries

- Treat the workflow as recommended, not mandatory.
- Use the developer's active model and provider. Do not require a local model.
- Keep `plan`, `implement`, `test`, `review`, and `pr` as the only public workflow entry points.
- Load internal capability references only when relevant to the current change.
- Treat CI/CD as the authority for enforced quality gates.
- Treat upstream skills, agents, and examples as design inputs; repository rules and evidence remain authoritative.

## Optional Tools

- GBrain: use an already configured read-only search/query for past project decisions when useful; do not install, configure, or write to it implicitly.
- Graphify: use an existing graph for broad dependency questions; never build or update it implicitly.
- RTK: use an existing binary selectively for noisy supported commands; return to raw output for exact evidence.
- Absence of any optional tool must not block a workflow. Use repository search, source reads, and raw commands as the fallback.

## Approval Rules

Invoking a skill authorizes its documented non-Git work. Ask for separate, explicit confirmation before creating or switching a branch, committing, pushing, rebasing, merging, or changing remote state.

Never merge a pull request. Never use force push, destructive reset, or checkout-based file restoration without a specific user request.

## GitHub Context

Resolve GitHub data in this order:

1. use an already configured GitHub connector or MCP tool;
2. use an authenticated `gh` CLI;
3. ask the user to provide the issue or PR content.

Do not require one particular integration. Never place credentials in workflow files, progress records, commits, logs, or pull request text.

## Progress Record

After implementation is approved, maintain `.freelunch/progress/issue-<number>.md` using `.freelunch/ai/progress-template.md`.

- `plan` is read-only and must not create the record. Include a proposed state summary in its response instead.
- `implement` creates the record after branch approval and updates implementation evidence.
- `test`, `review`, and `pr` update their stage, evidence, findings, and remaining work.
- Record only observed commands and results. Never infer or invent passing checks.
- Keep durable product decisions in the issue or repository documents and link them from progress.
- Keep the record free of secrets and sensitive source content.

## Completion Report

Every command reports:

- what it inspected or changed;
- decision sources and approvals received;
- checks run and their observed results;
- optional tools and internal capabilities used or skipped, with the reason;
- blockers, uncertainty, and remaining work;
- the next applicable public entry point.
