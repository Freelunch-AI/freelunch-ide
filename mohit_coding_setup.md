# Mohit's AI Coding Setup

This is my current personal setup and feature-building flow. It is descriptive, not a required team standard. Repository checks and CI define shared quality requirements regardless of which AI tools a developer uses.

## Tools

- **Agent clients:** Codex as the primary client and Claude Code as a secondary client.
- **Editor:** The workflow is editor-independent; I use the agent beside the repository and make manual edits in a VS Code-compatible editor when useful.
- **Models:** The best suitable cloud coding model available through my existing Codex or Claude subscription. The workflow is not tied to one model.
- **GitHub access:** GitHub CLI or the client's GitHub integration for reading issues, reviewing history, and creating pull requests.
- **Personal skills:** Five local Markdown workflows named `plan`, `implement`, `test`, `review`, and `pr`. They are available only in my local checkout and are not part of the project toolchain.

The local workflows adapt useful parts of [BMAD Method v6.10.0](https://github.com/bmad-code-org/BMAD-METHOD), [Agency Agents](https://github.com/msitarzewski/agency-agents), [GBrain](https://github.com/garrytan/gbrain), [Graphify](https://github.com/Graphify-Labs/graphify), [RTK](https://github.com/rtk-ai/rtk), and the review, security, simplification, explanation, and skill-authoring guidance in [Claude Code](https://github.com/anthropics/claude-code). These sources provide internal guidance rather than additional public commands, required plugins, or project dependencies.

## Feature-Building Flow

1. **Understand the issue:** Read the issue, decision-bearing comments, repository guidance, and relevant specifications. Surface missing or conflicting requirements before changing code.
2. **Plan:** Produce testable acceptance criteria, a scoped implementation plan, relevant risks, a test strategy, and a branch name. Planning is read-only.
3. **Create the branch:** Update `main`, create an issue-specific branch that follows `branching_strategy.md`, and confirm the approved scope.
4. **Implement:** Make small, scoped changes using existing repository patterns. Use specialist, security, or dependency analysis only when the changed surface warrants it.
5. **Test:** Run the smallest relevant checks first, then applicable unit, integration, end-to-end, static-analysis, formatting, and build checks configured by the repository.
6. **Review:** Independently compare the diff with the issue and plan, validate correctness and security concerns, simplify changed code where behavior is preserved, and rerun affected checks after fixes.
7. **Open the PR:** Commit and push only after reviewing the diff. Open a draft PR linked to the issue with a concise implementation and validation summary. Never merge automatically.
8. **Address feedback:** Reviewers provide findings or comments; I update my own PR branch, rerun relevant checks, and record the result.

## PR Ownership

Only the PR author changes or force-pushes that PR's branch. Reviewers inspect the branch and leave comments or findings for the author to apply. If a separate implementation is needed, it should use a new branch and PR rather than modifying another developer's work.
