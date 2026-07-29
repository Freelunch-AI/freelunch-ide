# Implementation And Refactoring

Use this capability for `implement` and for approved fixes from `review`.

1. Reconfirm the approved acceptance criteria and current task before editing.
2. Inspect existing patterns and ownership boundaries. Prefer established helpers, frameworks, and conventions.
3. Add a dependency only when the approved plan requires it. Verify source, version, license, maintenance status, and why existing code is insufficient.
4. Use test-first development when it clarifies behavior or protects a regression. Do not force placeholder tests or ceremony onto documentation and configuration changes.
5. Implement in small, coherent increments. Keep the diff within the issue scope.
6. Simplify recently changed code after behavior is correct: remove accidental duplication and nesting while preserving observable behavior.
7. Run focused checks during implementation and record only observed results.

Select specialist guidance by changed surface, not by a fixed checklist. Architecture, Go, TypeScript, platform, Kubernetes, SRE, CI/CD, documentation, or debugging perspectives are useful only when that surface is present.
