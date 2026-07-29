# Internal Skill Authoring

Use this capability when a FreeLunch workflow repeatedly needs the same non-obvious procedure. It adapts skill-creator guidance and GBrain's thin-router/fat-skill pattern while preserving the five-command product decision.

## Default Decision

Keep exactly five public skills: `plan`, `implement`, `test`, `review`, and `pr`.

- Add a reusable procedure under `.freelunch/ai/capabilities/` when it supports an existing workflow.
- Add detail as a directly linked internal reference when it is needed only for a specific surface.
- Add a new discoverable skill only through a separately approved change to the public workflow contract.
- Do not expose upstream agents, tools, or role names as additional commands.

## Creation Process

1. Run the procedure successfully on a concrete repository task before extracting it.
2. Collect representative usage: normal invocation, incomplete input, non-trigger, risky mutation, and recovery from failure.
3. Separate latent work that needs judgment from deterministic work that needs a tool. Do not add executable code when clear Markdown and existing tools are sufficient.
4. Define inputs, outputs, safety boundaries, approvals, evidence, and stop conditions.
5. Put core imperative procedure in the canonical `SKILL.md`; move surface-specific detail to a directly linked internal capability.
6. Write a specific description that states what the skill does and when it applies. Avoid a description so broad that it triggers on unrelated work.
7. Remove duplicated guidance and context the model already knows. Keep repository-specific constraints and failure lessons.
8. Test the workflow in each supported client shape without adding a generator or runtime.

## Progressive Disclosure

Use three levels:

1. public name and description for discovery;
2. canonical skill body for the stage procedure;
3. internal capability files loaded only for the changed surface.

Keep references shallow and explicit. A workflow should not need to crawl a tree to discover a mandatory rule, and the same instruction should have one authoritative home.

## Conformance Review

Verify:

- valid frontmatter and matching directory/name;
- exact public set remains five for Codex, Claude Code, Cursor, and OpenCode;
- every client adapter resolves to its canonical `.agents/skills/<name>/SKILL.md`;
- every referenced capability exists and is internal to discovery;
- triggering examples include positive, negative, incomplete-input, and safety cases;
- Git and remote mutations remain separately approved;
- source provenance, revision, license, adopted behavior, and exclusions are updated in `ai-stack.lock.json`;
- no secrets, machine-specific paths, unavailable scripts, or undeclared runtime assumptions were added.

After real use, update the reusable procedure that caused friction rather than adding a one-off workaround to a client adapter.
