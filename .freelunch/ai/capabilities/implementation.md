# Implementation And Simplification

Use this capability for `implement` and for developer-approved fixes from `review`. It combines BMAD's acceptance-criterion loop, Agency's minimal-change discipline, and the code-simplifier's behavior-preserving cleanup pass.

## Scope Control

1. Reconfirm the approved issue, acceptance criteria, task, base branch, and current worktree before editing.
2. Trace the smallest set of files and boundaries required for the task. Use `.freelunch/ai/capabilities/code-intelligence.md` when the path is cross-module or unclear.
3. Preserve unrelated user changes and avoid cleanup that is not required for the criterion.
4. Every changed line must support the task, its test, or a necessary contract/documentation update.
5. Surface worthwhile out-of-scope work as a follow-up with evidence. Do not smuggle it into the diff.

Minimal does not mean fragile. Validate untrusted input at system boundaries, preserve required error handling, and update every consumer when a contract intentionally changes.

## Implementation Loop

For each approved task:

1. Identify the acceptance criterion and proving check.
2. Add or adjust a failing test first when it clarifies behavior or protects a regression. Documentation and declarative configuration changes do not need ceremonial test-first work.
3. Implement the smallest coherent behavior that satisfies the criterion using existing patterns.
4. Run the focused proving check and inspect the observed result.
5. Handle specified error and boundary cases.
6. Update the progress record with changed files, decisions, and actual evidence.
7. Continue only after the current task is coherent and its focused check is understood.

If implementation evidence contradicts the approved plan, stop expanding the solution. Document the conflict and return for a developer decision.

## Behavior-Preserving Simplification

After behavior is correct, review only recently changed code:

- reduce avoidable nesting and branching;
- remove redundant variables, copies, comments, and abstractions;
- consolidate duplication only when the abstraction has a clear stable responsibility;
- choose readable control flow over dense expressions or clever one-liners;
- improve names when the existing name obscures changed behavior;
- preserve useful domain boundaries, test seams, error context, and explicit invariants;
- never remove a fallback, validation, log, or compatibility path without proving it is unnecessary.

Rerun the focused behavior checks after simplification. A cleanup is incomplete if observable behavior changed unintentionally.

## Dependencies And Generated Work

- Add a dependency only when the approved plan permits it. Verify canonical source, exact version, license, maintenance posture, install behavior, and transitive risk.
- Do not introduce a language runtime merely to execute Markdown skills.
- Test an approved bulk or generated operation on a small representative sample, inspect the output, then obtain any required mutation approval before the full run.
- Do not hand-edit generated output unless the repository explicitly owns it as source.
