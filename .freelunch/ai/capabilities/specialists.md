# Specialist Routing

Agency Agents provides rich role prompts. FreeLunch retains their useful review questions as internal lenses instead of exposing a large agent catalog. Select a lens because the changed surface triggers it, load only the relevant section, and return one reconciled workflow result.

## Routing Table

| Trigger | Lens | Primary workflow |
| --- | --- | --- |
| New module boundary, shared interface, state ownership, or durable dependency rule | Software architect | `plan`, `review` |
| Any scoped implementation or bug fix | Minimal-change engineer | `implement`, `review` |
| Correctness, maintainability, performance, or contract review | Code reviewer | `review` |
| GitHub Actions, build, packaging, deployment, IaC, or environment promotion | DevOps automator | `plan`, `implement`, `review` |
| Reliability target, observability, capacity, incident, or rollout behavior | SRE | `plan`, `test`, `review` |
| User, contributor, architecture, runbook, or operational documentation | Technical writer | `implement`, `pr` |
| Trust boundary, identity, secrets, untrusted input, dependency, CI, or infrastructure | Application security engineer | `plan`, `implement`, `review` |
| AI-generated integration, client/server secret boundary, database policy, or tool-enabled LLM | AI-code security auditor | `review` |
| Browser journey or cross-system end-to-end behavior | Test automation engineer | `test`, `review` |
| Readiness or completion claim | Evidence/reality checker | `test`, `review`, `pr` |

## Software Architect

- Start from the domain problem, existing boundaries, and team constraints before selecting technology.
- Decide whether simple layered code is enough before adding rich domain models, ports, events, CQRS, or services.
- Protect dependency direction: core policy must not depend on delivery, storage, vendor, or framework details without an explicit reason.
- Identify data/state ownership, transactional boundaries, integration contracts, failure containment, and evolution path.
- Compare realistic options by simplicity, reversibility, coupling, reliability, operability, and team cost.
- Record only decisions that prevent incompatible downstream work; name what each decision gives up.

Deliver: the smallest architecture, binding invariants, alternatives considered, deferred decisions, and revisit conditions.

## Minimal-Change Engineer

- Read the task literally and map each changed file to a requirement.
- Prefer the boring local change when it fully satisfies the contract.
- Avoid speculative flags, wrappers, compatibility shims, defensive branches for impossible internal states, and refactors of untouched neighbors.
- Validate real external boundaries even when the patch is small.
- Walk the final diff line by line; remove additions justified only by "while here" or hypothetical future use.
- List useful unrelated observations as follow-up instead of changing them.

Deliver: scoped diff, file-by-file reason, proving checks, and explicit follow-ups not included.

## Code Reviewer

- Prioritize correctness, security, contracts, tests, maintainability, and demonstrated hot-path cost.
- Make each finding specific: location, scenario, impact, evidence, and complete fix.
- Validate before reporting and provide the complete high-signal review in one pass when possible.
- Omit style preferences handled by automation and avoid praise or criticism that does not help the decision.

Deliver: severity-ordered findings, evidence, residual risk, and a clear no-findings statement when appropriate.

## DevOps Automator

- Inspect current CI/CD, task runners, artifact flow, environments, and repository platform choices before proposing tools.
- Review build reproducibility, least-privilege permissions, secrets, artifact provenance, promotion, health checks, rollback, and failure visibility.
- Prefer versioned declarative automation and idempotent operations.
- Require monitoring and recovery only in proportion to the deployed surface; do not add enterprise machinery to a local-only change.
- Verify current action/tool versions from primary sources when changing delivery configuration.

Deliver: pipeline or infrastructure change, security/rollback consequences, environment matrix, and exact validation evidence.

## SRE

- Define the user-visible reliability behavior before selecting metrics.
- Use relevant golden signals: latency, traffic, errors, and saturation.
- Distinguish service objectives, alert thresholds, and implementation metrics.
- Review timeouts, retries, backoff, idempotency, overload, degradation, rollback, and capacity.
- Prefer progressive rollout and measurable recovery over big-bang deployment.
- Reduce repeat operational work, but do not automate an unproven manual recovery path.

Deliver: failure modes, observable signals, recovery/rollback plan, and remaining operational risk.

## Technical Writer

- Write for the actual reader and task: contributor setup, operator recovery, architecture decision, API contract, or reviewer handoff.
- Verify every command, path, option, example, and stated behavior against the repository.
- Explain prerequisites, expected result, failure path, and recovery where needed.
- Preserve terminology and decision wording; remove filler, duplicated prose, stale promises, and implementation detail from product requirements.
- Keep durable rationale in the relevant decision document rather than in obvious code comments.

Deliver: concise accurate documentation and a note of any unverified example or environment-specific assumption.

## Application Security Engineer

- Build a small threat model around assets, actors, entry points, and trust boundaries.
- Trace attacker-controlled data to authorization, storage, execution, network, path, rendering, logging, and model-tool sinks.
- Separate exploitable findings from defense-in-depth suggestions and validate proposed fixes.
- Review dependencies and delivery paths as carefully as first-party code.
- Add a regression test for a fixed vulnerability when feasible and never expose raw secrets in evidence.

Deliver: testable security requirements or evidence-backed findings with exploit path, fix, and verification.

## AI-Code Security Auditor

- Check generated code for demo-friendly unsafe defaults: hardcoded/client-exposed secrets, permissive data policies, missing object authorization, and trusted client metadata.
- Keep untrusted text separate from higher-priority model instructions and independently authorize every tool action.
- Treat leaked credentials as already compromised and include rotation.
- Use conservative source-to-sink analysis, label heuristic conclusions, and prefer silence over a noisy false positive.
- Rescan or rerun the focused security test after an approved fix.

Deliver: worst-first findings with source, sink, exploit, confidence, remediation, and closure evidence.

## Test Automation Engineer

- Reserve end-to-end tests for critical journeys and prove lower-level behavior below the browser when possible.
- Give tests isolated data, condition-based waits, user-facing selectors, and failure artifacts.
- Set up shared prerequisites through stable fixtures or APIs when that setup is not the subject under test.
- Diagnose flakes at their root; retries reveal instability but do not cure it.
- Keep CI duration and parallel safety visible without weakening coverage.

Deliver: journey-to-test map, deterministic tests, artifact locations, and explicit flake or environment limitations.

## Evidence And Reality Checker

- Cross-check completion claims against the issue, actual diff, command results, and user-visible artifacts.
- For UI changes, inspect responsive screenshots and interaction states when a runnable app exists.
- For non-UI changes, prefer tests, structured output, build artifacts, configuration reads, or a reproducible failure path.
- State what was not tested and avoid unsupported readiness scores.
- Do not assume a first iteration must contain issues and never manufacture a minimum finding count.

Deliver: claim-to-evidence map, confirmed gaps, and a readiness statement bounded by observed evidence.

## Reconciliation

When lenses disagree, the active public workflow resolves the conflict using the issue, repository rules, source evidence, risk, and developer decisions. An upstream agent's persona, metric, command example, or default verdict is never binding.
