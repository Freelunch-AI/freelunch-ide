---
name: test
description: Test a FreeLunch issue implementation and record trustworthy evidence. Use when the user invokes test or asks to add missing tests and run the applicable unit, integration, end-to-end, typecheck, lint, formatting, and build checks for the current issue branch.
---

# Test

Read `.freelunch/ai/workflow.md`, `.freelunch/ai/capabilities/testing.md`, and `.freelunch/ai/capabilities/command-output.md` before starting. Read `.freelunch/ai/capabilities/security.md` or the test-automation lens in `.freelunch/ai/capabilities/specialists.md` when the changed boundary requires it.

Require an issue progress record or enough issue and branch context to identify the approved acceptance criteria. Do not treat missing state as permission to guess scope.

1. Set the progress stage to `testing`.
2. Discover applicable checks from repository configuration and documentation.
3. Map acceptance criteria and changed behavior to tests. Add missing tests when they are within the approved scope.
4. Run focused checks first and broader applicable checks afterward. Use an existing RTK installation only for supported noisy output and rerun the smallest failure raw when filtering omitted diagnostic evidence.
5. Classify and diagnose failures without weakening meaningful assertions, treating a retry-only pass as healthy, or hiding errors.
6. Record every command as passed, failed, skipped, or unavailable with concise observed evidence.
7. Update acceptance criteria and remaining work. Set the stage to `ready-for-pr` only when implementation and applicable checks are complete; otherwise use `blocked` or retain `testing`.

Update the progress record and recommend `review`. Do not claim CI will pass solely because local checks pass.
