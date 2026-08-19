# Freelunch IDE - Contributor Guide

## Codebase Explanation

### Folder & File Structure

### Step-by-step Codebase Tutorial

## Contributing Rules

### Issues

- Follow Issue Template
- Issues must be self-contained and highly descriptive of the problem and solution that should be implemented
- Every Bug Issue closed must create a ./docs/post-mortems/post-mortem_[i].md where i is the number of the issue, which describes the issue, when it happenned, the data evidence for the issue (telemetry, screenshots, user reviews, etc), desription of the solution, data evidence of the solution, who was involved in solving it, time it took to solve it, the process that was used to arrive at the solution, how to avoid and reduce the impact of similar problems in the future, the changes in standard operating protocols imposed.

### PRs

- Follow PR Template
- Code PRs need to resolve a specific issue
- PRs must be self-contained and highly descriptive of the solution and solution implementation
- Prefer multiple straighforward PRs over a single big PR with multiple different stuff
- Prefer batching multiple changes to the same file in a single PR
- Never commit on top of someone else's PR, just give review and make comments.
