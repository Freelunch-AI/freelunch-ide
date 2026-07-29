---
name: plan
description: Plan a FreeLunch repository change from a GitHub issue without modifying files or Git. Use when the user invokes plan, provides an issue URL for analysis, or asks for acceptance criteria, architecture, implementation tasks, test strategy, and a branch name before coding.
---

# Plan

Read `.freelunch/ai/workflow.md` and `.freelunch/ai/capabilities/planning.md` before starting.

Require one GitHub issue URL. Resolve the issue through the available GitHub integration, then inspect the repository context needed to validate it.

Remain read-only. Do not edit files, create progress state, create or switch branches, commit, push, or change GitHub data.

Produce:

1. a concise problem and scope summary;
2. ambiguities or conflicts that need a decision;
3. testable acceptance criteria;
4. the smallest suitable architecture and important tradeoffs;
5. ordered implementation tasks with relevant internal capabilities;
6. a risk-based test strategy;
7. dependency decisions;
8. a branch name that follows `branching_strategy.md`;
9. a proposed progress-state summary for `implement` to create after approval.

Ask focused clarifying questions only when repository evidence and primary-source research cannot resolve a material decision. End by asking the developer to approve or revise the plan. Do not invoke `implement` automatically.
