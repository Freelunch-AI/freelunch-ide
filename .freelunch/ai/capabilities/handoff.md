# Handoff And Pull Request Evidence

Use this capability when updating progress state and preparing `pr` output.

## Reviewer-Facing Content

Include:

- linked issue and exact approved scope;
- branch and base branch;
- acceptance criteria status;
- important technical decisions, alternatives, and tradeoffs;
- modules or contracts materially changed;
- exact validation commands and observed results;
- confirmed review findings and their disposition;
- security, migration, compatibility, rollout, or documentation impact when applicable;
- known limitations, skipped checks, blockers, and remaining work;
- source inspiration links when methodology or third-party design input materially shaped the change.

Use `Resolves #N` only when the pull request fully completes the issue. Otherwise use a non-closing link.

## Evidence Quality

- Read the actual diff, status, commits, and progress record rather than copying a session summary.
- Use raw output to confirm final diff and gate evidence when filtering could omit detail.
- Keep screenshots or recordings under the Implementation section when they help reviewers inspect UI or workflow behavior.
- Distinguish local checks from CI and state pending checks clearly.
- Do not paste full logs, source files, secret values, model prompts, or generated reports into the PR body.
- Do not claim production readiness, compliance, or complete security from an agent review.

The progress record is resumable working state. The PR description is the durable, concise review contract.
