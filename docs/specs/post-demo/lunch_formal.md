# Lunch Formal: Formal Verirication of Freelunch Codebases

Use Cases:
- guaranteeing that major safety and cost issues are not present in codebases that are heavily written by AI Agents, not humans.

Steps:
1. Codebase is mapped to a standard system specification (SSS) document containing code entitites, relationships, user types, user access points, external communication... everything needed to reason about the system. Where the system is platform + applications.
2. The SSS is converted to LEAN to enable formal reasoning
3. User invariants (things the user wants to verify is always true or not about the system) are expressed in natural language
4. User Invariants are converted to LEAN via a specialized LLM Formal Verification Agent
5. LEAN runs to try to prove the invaraint and spits a runnable test counter example if the invariant is not true
6. The test is run. If the test passes, then the LM Formal Verification Agent reasons about what it did wrong in the autoformalization steps, and corrects. It does this until the spitted test outputs what was expected after the LEAN run.
