---
name: review
description: Independently review a FreeLunch issue branch for correctness, acceptance-criteria drift, security, tests, and maintainability. Use when the user invokes review or asks for a high-signal review before opening a pull request.
---

# Review

Read `.freelunch/ai/workflow.md`, `.freelunch/ai/capabilities/context.md`, `.freelunch/ai/capabilities/review.md`, `.freelunch/ai/capabilities/security.md`, `.freelunch/ai/capabilities/specialists.md`, and `.freelunch/ai/capabilities/command-output.md` before starting. Read the applicable files under `.freelunch/ai/capabilities/` named `code-intelligence.md`, `testing.md`, or `explanation.md` when the scope or evidence requires that guidance.

Require the issue, approved plan, base branch, and current diff. Set the progress stage to `reviewing` when a progress record exists.

1. Review independently from the implementation narrative. Use isolated or multiple reviewer passes when the client supports them and risk warrants it; the main workflow must validate and reconcile every finding.
2. Check the actual diff against repository instructions, acceptance criteria, architecture, and tests.
3. Apply security review proportionally to changed trust boundaries.
4. Apply relevant correctness, test, error-path, type-invariant, comment-accuracy, security, and changed-code simplification lenses. Validate suspected defects and report only high-signal findings with severity, confidence, location, failure scenario, evidence, and a complete fix.
5. Ask for approval before changing code to address findings.
6. After approved fixes, rerun the smallest proving checks and affected broader checks.
7. Record findings and their disposition in the progress record. Set `ready-for-pr` only when no blocking finding remains.

Use raw source and diff context for final findings even when Graphify or RTK helped navigate. End with findings first, followed by verification evidence and residual risk. Recommend `pr` when ready. Never comment, approve, request changes, or merge on the developer's behalf without separate explicit approval.
