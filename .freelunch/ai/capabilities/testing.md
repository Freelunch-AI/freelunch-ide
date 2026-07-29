# Testing And Evidence

Use this capability for `test` and for verification inside `implement` or `review`. It combines behavioral coverage, deterministic automation, failure classification, and evidence receipts.

## Discover The Real Test Surface

Read checked-in configuration, task definitions, CI workflows, and nearby tests. Do not guess commands from the language alone. Apply `.freelunch/ai/capabilities/command-output.md` when output is noisy.

Map changed behavior to the cheapest reliable level:

| Level | Use for |
| --- | --- |
| Static/format | Syntax, formatting, types, lint rules, declarative structure |
| Unit | Pure behavior, invariants, boundary cases, regression protection |
| Integration | Component, process, database, API, filesystem, or tool boundaries |
| End-to-end | Critical human journeys where the integration itself is the risk |
| Build/package | Compilation, bundles, artifacts, installation, release shape |
| Operational | Health, migration, rollback, observability, infrastructure contracts |

Do not chase line coverage or browser tests when a lower level proves the contract better.

## Behavioral Coverage

For each acceptance criterion and changed contract, check:

- expected success behavior;
- important invalid, empty, limit, and failure behavior;
- authorization or tenant boundaries when applicable;
- concurrency, ordering, retry, cancellation, or idempotency when applicable;
- compatibility or migration behavior;
- user-visible state and accessibility for UI changes;
- error propagation and observable diagnostics.

Tests should assert behavior and contracts, not implementation trivia. A proposed test gap must name the regression it would catch and why existing tests do not already cover it.

## Deterministic Automation

- Let each test own or isolate its data. Avoid hidden order dependencies and shared mutable fixtures.
- Wait for observable conditions, responses, or state transitions rather than fixed sleeps.
- For UI automation, prefer accessible roles and labels; use stable test IDs only when semantics cannot identify the element.
- Set up repeated prerequisites through an appropriate lower-level fixture when the setup journey is not under test.
- Attach logs, traces, screenshots, or other artifacts needed to diagnose a CI-only failure.
- Treat retries as a signal of flakiness, not proof that a test is healthy.

## Execution Order

1. Run the smallest new or affected test.
2. Run the affected package or module suite.
3. Run applicable static, format, type, and lint checks.
4. Run integration, end-to-end, build, or operational checks only when the changed boundary requires them.
5. Run the broader repository suite when configured and feasible.
6. Re-run after a fix and inspect the result; do not infer success from the edit.

## Failure Classification

Classify each failure before changing code or tests:

| Class | Meaning |
| --- | --- |
| Product regression | Changed behavior violates the intended contract |
| Test defect or stale expectation | The test is wrong, brittle, or expects intentionally replaced behavior |
| Flake | The same unchanged target is nondeterministic; retry only to diagnose |
| Environment/infra | A prerequisite, service, platform, or resource is unavailable |
| New incomplete work | A newly added test correctly exposes unfinished implementation |

Never weaken assertions, remove a security test, accept a retry-only pass, or update an expected value merely to make the suite green.

## Evidence Receipt

For every applicable check record:

- exact command;
- scope and environment when material;
- exit status and meaningful counts;
- `passed`, `failed`, `skipped`, or `unavailable`;
- concise observed output or artifact path;
- limitation or remaining untested path.

Local success does not prove CI success. Say clearly which CI-only surfaces remain.

For Markdown workflow changes, also verify that exactly five public names exist for each client, every adapter resolves to its canonical skill, all referenced internal files exist, frontmatter is valid, the lockfile parses, and no unexpected executable or public command was added.
