---
name: implement
description: Implement an approved FreeLunch issue plan. Use when the user invokes implement after planning and wants the agent to create the approved issue branch, make scoped code or documentation changes, use relevant internal specialists, and track progress.
---

# Implement

Read `.freelunch/ai/workflow.md`, `.freelunch/ai/capabilities/context.md`, `.freelunch/ai/capabilities/implementation.md`, and `.freelunch/ai/capabilities/specialists.md`. Read the applicable files under `.freelunch/ai/capabilities/` named `code-intelligence.md`, `command-output.md`, `planning.md`, `testing.md`, `security.md`, `explanation.md`, or `skill-authoring.md` only when the changed surface requires that guidance.

Require an approved plan associated with one GitHub issue. If the plan is missing, stale, materially ambiguous, or conflicts with the repository, stop and direct the developer to `plan`.

1. Show the issue, base branch, proposed branch, scope, and first task.
2. Ask for explicit approval before creating or switching the branch.
3. After approval, create `.freelunch/progress/issue-<number>.md` from `.freelunch/ai/progress-template.md` and set the stage to `implementing`.
4. Implement the approved tasks in repository-sized increments. Follow existing patterns, keep dependencies explicit, load only relevant specialist lenses, and finish each increment with a behavior-preserving simplification pass.
5. Run focused checks as work proceeds and record observed evidence. An existing RTK installation may reduce routine noisy output, but use raw output for omitted failure detail and exact diff evidence.
6. Update acceptance-criterion status, decisions, changed files, validations, and remaining work in the progress record.
7. Do not commit or push unless the developer separately approves the exact Git action.

When implementation is complete, summarize the diff and evidence, update the progress record, and recommend `test`. Do not invoke another public workflow automatically.
