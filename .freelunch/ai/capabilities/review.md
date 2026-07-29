# Independent Review

Use this capability for `review`.

Review the branch as an independent senior engineer. Start from the issue, approved plan, repository instructions, base branch, and actual diff rather than implementation claims.

Prioritize findings in this order:

1. incorrect behavior, data loss, unsafe Git or operational actions;
2. security and authorization defects;
3. violated acceptance criteria or architectural contracts;
4. missing tests for material changed behavior;
5. maintainability problems likely to cause defects.

Require concrete evidence for every finding. Cite the file and location, explain the failure scenario, and avoid style-only comments already handled by formatters or linters. Recheck a suspected issue before reporting it and suppress duplicates.

Keep reviewer independence where the client supports isolated context or subagents. Fix only findings the developer confirms, then rerun the smallest checks that prove the fix and any affected broader suite.

End with findings ordered by severity, open questions, verification evidence, and residual risk. Say clearly when no issues were found.
