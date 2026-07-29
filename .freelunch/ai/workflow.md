# FreeLunch AI Workflow Contract

Apply this contract whenever a FreeLunch workflow skill runs.

## Boundaries

- Treat the workflow as recommended, not mandatory.
- Use the developer's active model and provider. Do not require a local model.
- Keep `plan`, `implement`, `test`, `review`, and `pr` as the only public workflow entry points.
- Load internal capability references only when they are relevant to the current change.
- Treat CI/CD as the authority for enforced quality gates.

## Approval Rules

Invoking a skill authorizes its documented non-Git work. Ask for a separate, explicit confirmation before creating or switching a branch, committing, pushing, rebasing, merging, or changing remote state.

Never merge a pull request. Never use force push, destructive reset, or checkout-based file restoration without a specific user request.

## GitHub Context

Resolve GitHub data in this order:

1. Use an already configured GitHub connector or MCP tool.
2. Use an authenticated `gh` CLI.
3. Ask the user to provide the issue or PR content.

Do not require one particular integration. Never place credentials in workflow files, progress records, commits, logs, or pull request text.

## Progress Record

After implementation is approved, maintain `.freelunch/progress/issue-<number>.md` using `.freelunch/ai/progress-template.md`.

- `plan` is read-only and must not create the record. Include a proposed state summary in its response instead.
- `implement` creates the record after branch approval and updates implementation evidence.
- `test`, `review`, and `pr` update their stage, evidence, findings, and remaining work.
- Record only observed commands and results. Never infer or invent passing checks.
- Keep the record free of secrets and sensitive source content.

## Completion Report

Every command reports:

- what it inspected or changed;
- decisions and approvals received;
- checks run and their observed results;
- skipped internal capabilities and why;
- blockers or remaining work;
- the next applicable public entry point.
