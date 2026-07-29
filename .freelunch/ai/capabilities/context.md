# Context And Durable Decisions

Use this capability whenever work depends on past decisions, issue discussion, repository-wide intent, or facts that may live outside the current session. It adapts GBrain's thin-router and source-aware retrieval ideas to repository-native state.

## Context Layers

Keep three kinds of state separate:

| Layer | Contents | FreeLunch source |
| --- | --- | --- |
| Durable project knowledge | Approved product decisions, issue discussion, specifications, architecture, contracts | GitHub issues and checked-in documentation |
| Operational progress | Current branch, accepted plan, completed tasks, validations, findings, remaining work | `.freelunch/progress/issue-<number>.md` |
| Session context | The current request, recent tool output, temporary reasoning | Current client conversation |

Do not turn the progress record into a second specification. Link durable decisions rather than paraphrasing them loosely, and do not treat session memory as evidence for a decision that should survive the session.

## Source Precedence

For the current change, resolve conflicts in this order unless the developer explicitly changes a decision:

1. the developer's current explicit instruction;
2. the latest clear decision in the linked issue discussion;
3. approved acceptance criteria and repository specifications;
4. the current code, tests, and configuration for observed behavior;
5. the issue progress record for operational state;
6. an optional connected knowledge tool;
7. external sources.

Do not silently choose when two authoritative sources conflict. Quote or link both, explain the practical consequence, and ask for the smallest decision needed to continue.

## Retrieval Protocol

1. Decompose the question into exact terms, concepts, relationships, and time-sensitive facts.
2. Search the linked issue, repository instructions, specifications, and relevant code before asking the developer to repeat context.
3. Read the most relevant sources deeply enough to confirm scope; do not load every document merely because it exists.
4. For architecture or dependency questions, apply `.freelunch/ai/capabilities/code-intelligence.md`.
5. Use an already configured GBrain search or query tool only when it is available and the question concerns past decisions or project context. Start with keyword search, then semantic query if exact search is thin.
6. Cite the issue comment, document path, code location, graph source location, or knowledge-page identifier behind every material decision.
7. State missing, stale, ambiguous, and conflicting context explicitly.

## Write Boundaries

- Record observed workflow state in the progress file after implementation begins.
- Do not write to GBrain, GitHub issues, or another external knowledge store without explicit approval for that remote mutation.
- Preserve a developer's exact decision wording when recording it. Separate it from agent inference.
- Never persist secrets, credentials, sensitive source excerpts, raw model prompts, or speculative conclusions.
- Test any approved bulk context migration on a small representative sample before processing the full set.

## Friction Record

When a workflow is blocked or confusing, record a concise progress entry with:

- stage or command;
- observed symptom;
- classification: `blocker`, `error`, `confused`, or `nit`;
- evidence or exact error;
- workaround used, if any;
- follow-up that would prevent recurrence.

Do not describe a tool as healthy or a workaround as successful until the proving check has been rerun.
