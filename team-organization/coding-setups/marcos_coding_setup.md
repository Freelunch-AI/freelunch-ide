# Marcos' AI Coding Setup

This is my personal setup and feature-building flow. Like the other setups in this folder it is descriptive, not a required team standard — `branching_strategy.md` and the repository CI define the shared requirements regardless of which tools a developer uses.

The whole setup rests on one principle: **I own the architecture, the AI owns the implementation inside it.** Everything below follows from that.

## Tools

- **AI client:** Claude Code, with the model chosen per task — a stronger model for complex or architectural work, a cheaper one for mechanical work. No second agent client.
- **IDE:** Cursor for Go and everything else; IntelliJ IDEA for JVM languages.
- **Version control:** plain `git` from the terminal. I create branches manually, following `branching_strategy.md`.
- **Slash commands:** I don't maintain a fixed command library. My slash commands *are* the skills I write. When I find myself explaining the same kind of task to the AI twice, that explanation becomes a skill, and from then on the task is one command.
- **Quality gate:** unit tests locally, CI/CD after the commit. Nothing else.

### Current skills

Skills live in `.claude/skills/` in the workspace and are invoked as `/<name>`:

| Skill | Purpose |
| --- | --- |
| `add-service` | Add a new service component to the `jusadvisor-web-go` framework — interfaces, no-ops, concrete implementation, HTTP handlers, lifecycle wiring, `main.go` registration. |
| `tdd` | Red-green-refactor loop for building a feature or fixing a bug test-first. |
| `fpt` | First Principles Thinking session — deconstruct assumptions and rebuild a solution from the ground up. |
| `grill-me` | The AI interviews me relentlessly about a plan or design until every branch of the decision tree is resolved. |
| `grill-with-docs` | Same grilling, but challenged against the existing domain model and documented decisions, updating `CONTEXT.md`/ADRs inline. |
| `handoff` | Compact the current session into a handoff document for the next agent session. |
| `write-a-skill` | Create a new skill with proper structure and progressive disclosure — this is how the list above grows. |

`add-service` is the clearest example of what a skill is for me: it is not a generic workflow, it is *my* framework's extension procedure written down once so the AI can repeat it exactly.

## Feature-Building Flow

1. **Pick the task.**
2. **Create the feature branch by hand**, named after the feature, following `branching_strategy.md`.
3. **Read the existing code myself** and decide how the feature is going to be implemented. This step is mine, not the AI's. I want to have the design in my head before any agent touches the repository.
4. **Judge whether the existing code is good enough to build on.**
   - **If yes** → go straight to step 6 and hand the implementation to the AI.
   - **If no, or if this is a new project** → step 5 first.
5. **Write the foundation by hand.** For a new project, or whenever the current code isn't a solid base, I always implement the framework, the first encapsulations and abstractions, the base classes and objects myself. I do this because AI is bad at defining code architecture: told to build something from the ground up it produces code with no encapsulation, no order and no organization, and once the codebase grows that becomes very expensive to undo. Creating the code model and the framework up front gives the AI a shape to work *inside* instead of a blank page.
6. **Teach the AI the framework** — what I built, how it is organized, and how I want it used.
7. **Turn the recurring work into skills.** For each specialized thing I want implemented, I write a skill that captures how it must be done in this framework. From then on that work is a single command.
8. **Implement by invoking the skills.** No orchestration, no gates, no multi-agent setup. Just the right skill for the task.
9. **Always require unit tests.** Every implementation the AI produces comes with unit tests. Passing tests are how I measure whether the implementation is actually correct.
10. **Commit the feature branch and trigger CI/CD** once the feature is implemented and all unit tests pass. CI is my external quality gate.

## Why it stays this simple

I deliberately don't do a few things the heavier setups do:

- **No approval-gate pipeline.** The gate is the architecture I wrote by hand plus the unit tests. If the AI works inside a good framework and the tests pass, the risky part was already handled at step 5.
- **No multi-agent orchestration, worktrees or sub-branch supervision.** One agent, one branch, one terminal.
- **No per-feature spec documents.** The design lives in my head at step 3 and in the code model at step 5. If I need to stress-test it, I use `/grill-me` or `/grill-with-docs` instead of writing a spec.
- **No fixed command library to maintain.** Skills are created on demand and only when a task repeats.
- **No local review, security or simplification stage.** That's what CI and the PR review are for. Keeping those out of the local loop is what keeps the environment simple and flexible.

The cost of this flow is that it depends on me doing the architecture work honestly. The benefit is that there is almost nothing to maintain, and the AI is fast and reliable precisely because it never has to invent structure.
