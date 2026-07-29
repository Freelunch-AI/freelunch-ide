---
name: plan
description: Plan a FreeLunch repository change from a GitHub issue without modifying files or Git. Use when the user invokes plan, provides an issue URL for analysis, or asks for acceptance criteria, architecture, implementation tasks, test strategy, and a branch name before coding.
---

# Plan

Read `.freelunch/ai/workflow.md`, `.freelunch/ai/capabilities/context.md`, and `.freelunch/ai/capabilities/planning.md` before starting. Read the applicable files under `.freelunch/ai/capabilities/` named `code-intelligence.md`, `security.md`, `testing.md`, or `specialists.md` only when the issue surface triggers that guidance.

Require one GitHub issue URL. Resolve the issue and decision-bearing comments through the available GitHub integration, then inspect the repository context needed to validate it. Search durable context before asking the developer to repeat a prior decision.

Remain read-only. Do not edit files, create progress state, create or switch branches, commit, push, or change GitHub data.

Produce:

1. a concise problem and scope summary;
2. ambiguities or conflicts that need a decision;
3. testable acceptance criteria;
4. the smallest suitable architecture and important tradeoffs;
5. ordered implementation tasks with relevant internal specialist lenses;
6. a risk-based test strategy;
7. dependency decisions;
8. a branch name that follows `branching_strategy.md`;
9. a proposed progress-state summary for `implement` to create after approval.

For cross-module questions, use `.freelunch/ai/capabilities/code-intelligence.md`; an existing Graphify graph is optional and source inspection remains authoritative. Ask focused clarifying questions only when repository evidence and primary-source research cannot resolve a material decision. End by asking the developer to approve or revise the plan. Do not invoke `implement` automatically.
