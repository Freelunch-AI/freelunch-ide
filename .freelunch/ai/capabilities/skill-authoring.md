# Internal Skill Authoring

Use this capability when a workflow repeatedly needs the same non-obvious procedure.

1. Execute the procedure successfully on a concrete repository task before extracting it. Do not create a skill from an imagined process.
2. Identify the reusable instructions, references, scripts, inputs, outputs, and safety boundaries.
3. Prefer a concise internal capability under `.freelunch/ai/capabilities/` when it supports one of the five public workflows.
4. Add a new discoverable skill only through a separately approved workflow-surface change. The default is to keep exactly five public skills.
5. Keep procedural instructions imperative and move detail into directly linked references.
6. Add representative positive, incomplete-input, negative-trigger, and safety tests.
7. Update `ai-stack.lock.json` before adapting any external source and verify its immutable revision and license.

Do not copy upstream content whose license does not permit redistribution. Attribute reference-only design inputs in the lock manifest.
