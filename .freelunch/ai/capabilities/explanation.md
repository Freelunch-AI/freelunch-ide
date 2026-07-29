# Explanation And Developer Handoff

Use this capability when the developer asks how a change works, when an approval depends on understanding risk, or when a non-obvious repository pattern should survive handoff. It adapts explanatory output style without injecting a mandatory block into every response.

## Explain What Matters

1. Start from the issue outcome and acceptance criterion.
2. Trace the affected components, state or data flow, and important boundary.
3. Explain why the implementation follows an existing repository pattern or why a new decision was necessary.
4. Name the tradeoff, failure path, security consequence, and rollback implication when relevant.
5. Show how tests or observed evidence demonstrate the behavior.
6. Link the actual files, decisions, and commands rather than teaching generic programming concepts.

Keep straightforward changes brief. Use a diagram or longer walkthrough only when relationships are difficult to explain linearly.

## Accuracy Rules

- Distinguish observed behavior, documented intent, and inference.
- Do not claim a check ran when it was skipped or inferred.
- Verify code comments and examples against current code before using them as explanation.
- Avoid repeated preambles, decorative insight boxes, lectures, quizzes, and token-heavy narration.
- Correct misunderstandings directly and ask a knowledge-check question only when the developer requests one or a risky approval genuinely depends on it.

A useful handoff leaves the developer able to locate the change, understand the decision, reproduce the evidence, and recognize the main failure mode.
