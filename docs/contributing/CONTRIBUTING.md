# Freelunch IDE - Contributor Guide

## Codebase Explanation

### Folder & File Structure

### Step-by-step Codebase Tutorial

## Contributing Rules

### Global Spec (founding_doc.md, roadmap.md, tech_stack.md)

- is the source of truth for that work we need to do
- remains high-level: does not include implementaiton details that are subject to freuqnt changes such as data models, code structure, etc
- is updated after every Code PR merge: previous steps are removed and remaining docs are updated with new requirements surfaced.

### Issues

- Follow Issue Template
- Issues must be self-contained and highly descriptive of the problem and solution that should be implemented
- Every Bug Issue closed must create a ./docs/post-mortems/post-mortem_[i].md where i is the number of the issue, which describes the issue, when it happenned, the data evidence for the issue (telemetry, screenshots, user reviews, etc), desription of the solution, data evidence of the solution, who was involved in solving it, time it took to solve it, the process that was used to arrive at the solution, how to avoid and reduce the impact of similar problems in the future, the changes in standard operating protocols imposed.
- If Issue is incosistent with Global Spec in any way: (1) Flag the incosistency as an issue comment; (2) assuem the issue's version is the right version, unsless
stated otherwise or the issue gets edited.

### PRs

- Follow PR Template
- Code PRs need to resolve a specific issue
- PRs must be self-contained and highly descriptive of the solution and solution implementation
- Prefer multiple straighforward PRs over a single big PR with multiple different stuff
- Prefer batching multiple changes to the same file in a single PR
- Never commit on top of someone else's PR, just give review and make comments.
