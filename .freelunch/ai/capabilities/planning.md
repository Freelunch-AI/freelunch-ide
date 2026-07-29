# Planning And Architecture

Use this capability for `plan` and when implementation uncovers a material design gap. It combines BMAD's stakes-scaled product and architecture discipline with Agency's pragmatic architecture lens.

## Intake

1. Read the issue, its decision-bearing comments, repository instructions, relevant specifications, and recent related changes.
2. Apply `.freelunch/ai/capabilities/context.md` before asking the developer to repeat a past decision.
3. Separate requirements, observed current behavior, developer decisions, and agent assumptions.
4. Calibrate depth to stakes. A small internal change may need a few acceptance criteria; a cross-module platform decision needs explicit boundaries, failure modes, and operational consequences.
5. Identify the real concern set: compatibility, security, data ownership, integration density, public contracts, migration, reliability, operations, UX, or deployment. Do not force every plan through a fixed document template.

## Requirements Discipline

- Describe capabilities and outcomes before implementation mechanisms.
- Give each acceptance criterion one testable behavior and include important negative or error behavior.
- Preserve explicit exclusions and non-goals.
- Use user journeys only when a multi-step human workflow is materially involved. Keep simple internal tooling concise.
- Mark unresolved inference as an assumption and state what would invalidate it.
- Keep technical detail that constrains implementation; omit detail that merely demonstrates knowledge.

## Architecture Spine

Record an architectural decision only when independent implementations could otherwise choose incompatibly, the choice is not obvious from existing code, and the tradeoff is real.

For each load-bearing decision state:

- `Binds`: components or future work that must follow it;
- `Prevents`: the specific divergence or failure it avoids;
- `Rule`: the concise invariant to follow;
- `Tradeoff`: what becomes harder or less flexible;
- `Revisit when`: the measurable condition that justifies reopening it.

Prefer the repository's current paradigm and dependency direction. For new boundaries, compare the smallest realistic options, name reversibility and operational cost, and avoid introducing DDD, ports-and-adapters, CQRS, eventing, or microservices unless their constraints solve an observed problem.

## Plan Construction

1. State the problem, scope, non-goals, and decision sources.
2. List conflicts or questions that materially change the result.
3. Produce acceptance criteria traceable to the issue.
4. Describe the smallest architecture, boundaries, data flow, and important alternatives.
5. Identify security, compatibility, migration, rollback, observability, and documentation needs that actually apply.
6. Break work into ordered, independently verifiable tasks. Keep unrelated cleanup as follow-up.
7. Map each criterion and risky boundary to a test level and evidence type.
8. Decide whether any dependency is allowed, including source, version, license, and why existing code is insufficient.
9. Name relevant internal specialist lenses and why each is triggered.

For a small, well-specified change, keep the result short. Do not manufacture a PRD, architecture document, ADR, or meeting workflow when acceptance criteria and a few tasks are enough.
