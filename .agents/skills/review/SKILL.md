---
name: review
description: Independently review a FreeLunch issue branch for correctness, acceptance-criteria drift, security, tests, and maintainability. Use when the user invokes review or asks for a high-signal review before opening a pull request.
---

# Review

Read `.freelunch/ai/workflow.md`, `.freelunch/ai/capabilities/review.md`, `.freelunch/ai/capabilities/security.md`, and `.freelunch/ai/capabilities/specialists.md` before starting. Read the testing or explanation capability when verification evidence or developer handoff is incomplete.

Require the issue, approved plan, base branch, and current diff. Set the progress stage to `reviewing` when a progress record exists.

1. Review independently from the implementation narrative. Use an isolated reviewer context when the client supports it without weakening repository safety rules.
2. Check the actual diff against repository instructions, acceptance criteria, architecture, and tests.
3. Apply security review proportionally to changed trust boundaries.
4. Validate suspected defects and report only high-signal findings with severity, location, failure scenario, and rationale.
5. Ask for approval before changing code to address findings.
6. After approved fixes, rerun the smallest proving checks and affected broader checks.
7. Record findings and their disposition in the progress record. Set `ready-for-pr` only when no blocking finding remains.

End with findings first, followed by verification evidence and residual risk. Recommend `pr` when ready. Never approve or merge the pull request on the developer's behalf.
