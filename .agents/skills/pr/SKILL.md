---
name: pr
description: Prepare and open a draft FreeLunch pull request for the current issue branch. Use when the user invokes pr and wants a linked draft PR containing implementation, test, and review evidence without automatic merging.
---

# Pull Request

Read `.freelunch/ai/workflow.md`, `.freelunch/ai/capabilities/context.md`, `.freelunch/ai/capabilities/handoff.md`, and `.freelunch/ai/capabilities/command-output.md` before starting.

Require an issue, non-main branch, actual diff, and review/test evidence. Read `.github/pull_request_template.md` and `branching_strategy.md`.

1. Inspect Git status, branch ancestry, commits, the exact raw diff, progress record, decision sources, and applicable check results.
2. Identify uncommitted work, an unpushed branch, failures, unresolved review findings, or divergence from `main`.
3. Ask for separate explicit approval before each required commit, rebase/update, or push. Never force push.
4. Draft the PR title and body from the repository template. Link the issue, use `Resolves #<number>` only when the PR fully completes it, include exact validation evidence and remaining risks, and link material source inspiration when external methodology shaped the change.
5. Show the draft content and ask for approval before creating remote GitHub state.
6. Create a draft PR targeting `main` through the available GitHub integration.
7. Record the PR URL and set the progress stage to `pr-opened`.

Stop after reporting the draft PR URL and any pending checks or reviews. Never mark it ready, approve it, enable auto-merge, or merge it unless the developer makes a separate, explicit request allowed by repository policy.
