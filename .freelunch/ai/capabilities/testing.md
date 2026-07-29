# Testing And Evidence

Use this capability for `test` and for verification inside `implement` or `review`.

1. Discover the repository's real test, typecheck, lint, formatting, and build commands from checked-in configuration and documentation.
2. Map each acceptance criterion and changed behavior to the cheapest reliable check.
3. Add missing tests for changed behavior when they are in scope. Favor deterministic tests with meaningful failure messages.
4. Run focused tests first, then the broader applicable suite. Include integration and end-to-end checks only when the changed boundary requires them.
5. Distinguish product failures, test defects, environment failures, and checks that were not available.
6. Record the exact command, status, and concise evidence. Never report a check as passing when it was skipped, unavailable, or only inferred.
7. Avoid chasing coverage percentages without behavior value. Call out important untested paths explicitly.

Do not weaken assertions, remove tests, or hide failures merely to produce a green result.
