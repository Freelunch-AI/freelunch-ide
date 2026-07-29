# Independent Review

Use this capability for `review`. It adapts Agency's reviewer perspective and Anthropic's code-review toolkit into one high-signal, client-neutral process.

## Establish Scope

1. Read the issue and decision-bearing comments, approved plan, progress record, repository instructions, base branch, PR title/body when present, and exact diff.
2. Resolve instruction scope by directory: only apply a nested instruction file to changes inside its subtree.
3. Review the actual change rather than trusting the implementation summary.
4. Use raw diff output for final evidence. Graph or compressed output may help navigate but cannot be the sole basis for a finding.

## Review Lenses

Apply only relevant lenses, but do not skip correctness:

1. **Contract and scope:** acceptance criteria, non-goals, repository rules, unintended files, public API or data compatibility.
2. **Correctness:** compilation, logic, boundary values, state transitions, races, resource lifecycle, and failure scenarios.
3. **Security:** trust-boundary and data-flow review from `.freelunch/ai/capabilities/security.md`.
4. **Tests:** behavior-to-test mapping, negative paths, brittleness, and untested material risk.
5. **Error behavior:** broad catches, swallowed errors, silent defaults, fallback chains, exhausted retries, cleanup, and actionable diagnostics.
6. **Types and invariants:** invalid representable states, construction validation, mutation guards, ownership, and unsafe conversion.
7. **Comments and documentation:** factual accuracy, non-obvious rationale, stale examples, parameters, side effects, and error claims.
8. **Maintainability and performance:** only concrete defect risk, avoidable hot-path cost, or complexity likely to cause errors.
9. **Simplification:** changed-code opportunities that preserve behavior and do not broaden scope.

## Finding Validation

For every suspected issue:

1. Locate the introduced or changed line and trace the relevant context.
2. Describe a concrete input, state, or execution path that triggers the problem.
3. Check whether tests, callers, framework behavior, or repository rules invalidate the concern.
4. Distinguish introduced defects from pre-existing conditions.
5. Reproduce with a focused test or command when practical.
6. Suppress duplicates, formatter/linter-only notes, subjective preferences, and findings that remain speculative.

Use an isolated reviewer context or multiple independent passes when the client supports it and the risk warrants the cost. The main workflow must validate and reconcile their output. Subagents are optional, not part of the contract.

## Finding Format

Report findings first, ordered by severity:

- `Critical`: exploitable security issue, data loss, unsafe operational action, or unusable core path;
- `High`: clear wrong behavior, authorization failure, broken contract, or likely production regression;
- `Medium`: material edge case, test gap, or maintainability defect with a concrete failure path;
- `Low`: include sparingly when it prevents real future confusion; omit cosmetic nits.

Each finding includes severity, confidence, file and line, violated contract or failure scenario, impact, evidence, and the smallest complete fix. A suggestion block is appropriate only when that small suggestion fully resolves the issue.

If no validated issues remain, say so directly and list residual test or environment gaps. Do not invent findings to satisfy a quota.

## Mutation Boundaries

- Review is advisory and read-only until the developer approves fixes.
- Do not post comments, approve, request changes, mark ready, or merge without separate explicit approval.
- After an approved fix, rerun the focused proving check and any affected broader check, then update the finding disposition.
