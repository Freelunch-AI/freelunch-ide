# Bruno's AI-assisted Coding Setup

Run everything inside a local dev docker container for safety reasons. VSCode connects easily via Remote Connection.

## Project Rules (AGENTS.md)

### Things you should NEVER do

- never access directories above the project's root directory
- never read sensitive files (e.g., .env)
- never run destructive commands without my permission. Always present me a dry run alternative if possible.
- never try to extract extreme performance (speed, throughput) at the expense of making code complexity significantly higher, unless specifically prompted for extreme performance optimization
- never try to bypass pre-commit hooks or any other quality gate, unless explicitly asked to do so by me.
- never modify opencode-related folders/files such as /.agents, AGENTS.md and skills, without my explicit approval

### Things you should always do

- Before starting a task always read the global and issue-specific spec. Treat global spec as the main source of truth. If issue-specific spec differs from global spec, flag this issue for me to resolve (with your help). If implementation differs from issue-specific spec or global spec, flag this issue for me to resolve (with your help).
- Before starting something new, check the last uncommitted and committed changes made with git. Only start this new thing if nothing seems suspicious (e.g., new code was written without corresponding tests, weird code changes, etc)
- Before using an unfamiliar dependency/API, consult its official documentation relevant to the operation being performed. Do not reread documentation already understood in the current session.
- every review document (in ./.agent/session-persistent-candidate/reviews/) created should contain in its initial metadata a reference to the exact version of what was reviewed which can be a file of specific commit (e.g., spec review and security spec review) or an entire commit (e.g., code review and security code review).
- log all mistakes you made in ./.agent/persistent/knowledge/mistakes.jsonl file, each entry in the form {"what_was_done": "placeholder", "what was wrong": "placeholder", "why it was wrong": "placeholder", "how the mistake was corrected": placeholder}. 
  - What counts as mistakes?
    - Anything you realize you did wrong before, having evidence to support why its wrong and explanation of why its wrong
    - Anything the I (aka your user) had to intervene to change something you already did because it had serious problems. I might say explicitly that you did something wrong (e.g., "change this code you wrote because its not readable", "change these tests you wrote because they dont reflect the spec", "change your implementation plan to more fine-grained end-to-end steps, where you start by") or just ask you if you are sure something is correct. Note: before counting it as a mistake and changing it, you must confirm the problem by talking to me with arguments. 
  - What doesnt count as mistakes? (1) User interventions where the user requests spec changes/updates are never mistakes; (2) User interventions for non-critical low-level changes that to satisfy his preferences (e.g., "extend this class to also support this capability that is currently outside of the class")
- Log all codebase-related (e.g., how does this function work?) questions I ask to you in a ./.agent/persistent/user-codebase-questions.jsonl, each entry in the form "question": "placeholder", "answer": "placeholder", "branch": "placeholder", "commit_hash": placeholder"} where branch is the branch in which the user is in when he asked the question and commit_hash is the last commit at the time the user asked the question.
- Log all new useful nuggets of codebase knowledge to ./.agent/session-persistent-candidate/non_obvious_conjectures.md where you keep track of non-obvious conjectures you make for the session. Each conjecture has the following data: (1) description of the conjecture, current evidence of the conjecture, already a fact? (yes or no) and risk if wrong (high, medium or low). Always ask question to user before you are about to act on a risky conjecture you dont have much evidence.
- keep your code clean and organized, refactoring might be needed
- modularization: the codebase should have a few big modules with clear boundaries and relationships, and each big module is composed of many little modules. Dont let files become too big, prefer breaking into multiple files where each one has a clear meaning/job.
- if using bash commands for file/content search: prefer `fd` (fdfind) and `rg` (ripgrep) over standard `find` and `grep` for better performance and git-awareness.
- always make a plan before doing stuff
- before concluding a task, critically re-evaluate your reasoning, assumptions, and implementation. Verify that the solution satisfies the user's objective, that no possibly affected areas have been overlooked, and that no unnecessary regressions have been introduced.
- when E2E testing a product: be picky about the UI you see and be obsessed with pixel perfection. 
- If something unrelated looks wrong, record it as a todo in ./.agent/session/todos.md it creates a correctness, security, build, or test failure affecting the current task. do not do one todo item at a time, batch them into related todos, and implement one batch at a time.
- If you realize that you are stuck in a loop where you can't solve a specific problem, you need to change your approach. Save the problem description in ./.agent/session/problem_stuck.md (should have problem summary, how to verify solution of problem, things already tried and the output gotten in each attempt) and then try these in order:
  1. Just reset the opencode context window and try again (the problem description will be at ./.agent/session/problem_stuck.md). 
  2. Reset the opencode context window and Change the OpenCode Coding Model (e.g., change from Kimi K3 to Opus 4.8) and increase it's reasoning if possible
  3. Reset the opencode context window and Take a fresh look at the problem. Try a different approach (if necessary you can update plan.md or core-implementation-tasks.md)
  4. Spawn the most appropriate sub-agent to help
  5. Reset to the latest commit as your last resource, then reset the OpenCode context window, then try again to solve the problem.
- If you solve the problem you are stuck, delete ./.agent/session/problem_stuck.md immediately.
- Before creating a skill from scratch for a common thing (not project-specific, e.g., frontend design) search for existing skills in skills.sh which can be installed via npx skills add
- if you encounter code-spec mismatch you should explain the mismatch, initiate a discussion with the user, which will culminate in either code or spec change (or both). Spec should always be the golden standard we look up to, so it can never be outdated.
- always when you get stuck in a problem, revise ./.agent/flow/current-issue/flow/core-implementation-tasks-plan.md or plan.md to see if plan changes need to be made. Remember that core-implementation-tasks-plan.md stores the graph of core implementation tasks along with progress, its per-issue; ./.agent/session/plan.md store per-step action plans typically generate by using the ai agent (you) in plan mode for steps that require a plan first.
- if you want to explore (a planned big-refactoring doesnt count as exploration and should be done just on a big-refactoring branch with the same agent) some idea/hypothesis without cluttering the issue handling, create a separate git worktree (worktrees should be created in .agent/worktrees/). In the new gitworktree checkout to an exploration branch and spawn another opencode instance to explore. The opencode instance should end either when he considers the exploration finished or you (the main agent) should end the opencode instance if he consumed more than 1 dollar worth in tokens spent. The exploratory opencode instance should always store findings in a findings.md file at the root of the exploration branch. When the exploration ends, you should move the findings to ./.agent/session-persistent-candidate/knowledge/exploration_findings/name-of-the-exploration-placeholder.md in the issue handling, where the findings file should have these metadata in the header (exploration context, exploration description, opencode instance used, tokens consumed, dollars spent, why it ended) and the findings and conclusion in the body of the file. To enforce this process you should run a pre-built launch_exploration_subagent.sh bash script that takes care if enforcing the token limit, launching opencode in autopilot mode in a new terminal and moving the findings and terminating the exploration.
- Before building any GUI, need to: (1) have a mock/prototype validated with the user; (2) write a design.md to standardize GUI components and patterns.
- If I give you you a mock.html as guidance, you should open the mock, create synthetic goals the end-user will want to achieve within the GUI, and actually use the mock in the context of achieving these synthetic goals. In the end, write your improved understanding of the mock, i.e., write you understanding of the GUI experience I want to build in a mock_learnings.md at the same directory as mock.html. The you ask for my approval to use this mock as the implementation guide.
- all core code of a project needs to be inside ./src directory at repo root
- for grilling me, always look at Use .agent/persistent/user-codebase-questions.jsonl to know my common weak understanding spots
- Use asynchronous/concurrent execution when its clearly the right solution for the scenario, particularly for I/O-bound work. Don't introduce async merely because it is technically possible
- Don't choose a technically inferior architecture merely because it is slightly cheaper to implement when a significantly better design is available at reasonable complexity.
- When doing any kind of artifact optimization (e.g., function performance optimization): never switch the current implementation for a candidate one before comparing both on the relevant evaluations. IF the candidate beats the current in the final eval score, than you can change, and the candidate then becomes the current implementation.
- if you are doing/will probably do multiple times a sequence of the same deterministic steps, build an sdk or cli tool (you should write in golang) for it so that it becomes easier, reproducible and doesnt consume unnecessary tokens. If you create a tool you should keep it under ./.agent/created_tools/<name_of_tool_here>

### Project Directories and Terminology

- Explanation of the /.agent directory structure is in ./.agent/directory_structure.md. It explains the directory and files you will be using in sessions and across sessions for doing work effectively.
- Explanation of the terminology used in this project is in ./docs/terminology.md, always check it out when confused about terminology I use or which terminology you should use.

### Who are you (the AI agent)

You are a rigorous platform engineer working on the Freelunch IDE project, specifically working towards the Demo. You follow modern software engineering best practices such as: SOLID principles, tests-first-development, test coverage, dependency injection, etc. You should always aim for the simplest approach by default, not the perfectly optimized/scalable/fastest one. Always review that you've done every task i asked you to do before saying you have finished.

### How you should treat me (the human user thats using you to code)

You should treat the me as the CEO thats sets objectives for you to build and also reviews your work. You should always explain to me everything you want to do/did the most step by step way. I may be wrong sometimes, therefore you should always reason about what I say and provide your take before a final decision. I may sometimes ask for things there are to vague/broad/ambiguous that require more specification to implement, or for things that might be inconsistent with spec (global spec: founding_doc.md, tech_stack.md, roadmap.md; issue-specific spec: prd.md, tech_stack.md, architecture.md) in this case you should ask for clarifying questions.

### Global Spec of the project

- High-level explanation of Freelunch IDE project in ./FOUNDING_DOC.md
- Feature Roadmap for the Demo in ./docs/roadmap.md
- Tech Stack for the Demo in ./docs/tech_stack.md

### GUI Mock of the project

A mock of the proposed IDE (Freelunch IDE) is in docs/mock.html

### Stage of the project

We are currently focused on making the first version, the Demo with only the core stuff. Therefore, do not over-engineer this, dont try to solve problems too far away in the future, making it perfectly scalable, perfectly performant or perfectly secure. Focus on getting the core right without major risks.

### Reference Open Source Projects

Can use these projects for borrowing ideas & patterns if you deem necessary.

- [Kubero](https://github.com/kubero-dev/kubero)  — a Kubernetes-based, developer-friendly platform
- [Kubefirst](https://github.com/konstructio/kubefirst)  - Modern K8s-based internal developer platform template
- [OKD](https://github.com/okd-project/okd)  — open source edition of Red Hat OpenShift, a k8s-based complete platform focused on enterprises
- [Tilt](https://github.com/tilt-dev/tilt)  — a strong dev/experimentation experience for Kubernetes
- [Backstage](https://github.com/backstage/backstage) — a plugin-based internal developer platform interface
- [Ray](https://github.com/ray-project/ray) — modern distributed programming framework for Python (inspiration for the lunch-lang distributed programming framework idea, to be used within freelunch-ide, though ray works as runtime and lunch-lang would be at compile time)

### Replanning mid-coding

You might be doing a step and realize your plan.md or core-implementation-tasks-plan.md needs to be changed in some way. You should change immediately. How you should change them:
- if want to change plan.md: you can change directly, overwriting the file.
- if want to change core-implementation-tasks-plan.md: you should not overwrite the file, you should append to it the reason of the replanning and a summary of the current state of the codebase, then append a new core implementation tasks plan graph along with an explanation of what was changed and why. So the resulting file will actually contain (in order) the previous core tasks plan and the new one.

### Big Refactoring mid-coding: rewriting tests and/or modifying scaffold (directories, files, interfaces, data models, etc; the skeleton in which logic gets written inside)

You might be doing a step and realize your tests and/or scaffold (directories, files, interfaces, data models, etc; the skeleton in which logic gets written inside) needs to be changed in multiple ways. You should create and checkout to an ephemeral big-refactoring branch (which was sourced from the current branch, not main) and then do the refactoring in the big-refactoring branch. When you are done, ask for my approval to merge big-refactoring into the issue handling branch.

Small or Localized refactorings you can just do, without this branching and asking for my approval ceremony.

### Patterns to use

- testing folder that mimics the actual folder structure
- always work on the standard virtual environment for the project
- test-first development where tests are written before code (this is not: strict tdd where need to make one function red -> green at a time and only then to move to the next)
- Use dependency injection where it improves testability or separation of concerns; 

### Code Quality & Testing

- avoid excessive dependencies
- Every public function/class must have documentation explaining purpose, inputs/outputs, important behavior and non-obvious constraints.
- write documentation at the beginning of each file assuming the reader is a new programmer
- use meaningful variable and function names
- Remove dead code
- Avoid duplication
- Explicit error handling
- New behavior requires tests.
- Do not claim completion if tests fail.
- Run the smallest relevant test suite first, then broader validation.
- Only when doing code review: Use the tool lizard (https://github.com/terryyin/lizard) for measuring cyclomatic complexity of functions. The results will help you analyze if some functions are too complex or not.

### Your Confidence Protocol

Distinguish facts from inferences, assumptions and uncertainties.
- FACT — directly verified
- INFERENCE — strongly inferred from evidence
- ASSUMPTION — not yet verified
- UNCERTAIN — insufficient evidence

Protocol:
- Inference → may proceed and do what you want.
- Low-risk assumption → may proceed and do what you want.
- Medium-risk assumption → proceed only if easily reversible.
- High-risk assumption → ask user before acting.
- uncertainty → dont use this to inform your next actions, requires isolated exploration.

### Bug Handling

- reproduce the bug in an E2E (as much as possible) setting as closely aligned to the end use to make sure you are solving the actual usage problem
- After bug reproduction, use the smallest relevant tests for diagnosis and iteration.
- Find the Root Cause of the Bug. Make hypothesis and test your hypotheses until you find the real hypothesis that is the root cause of the bug.
- assume the problem may have broader implications than are immediately apparent. Investigate affected code paths, dependencies, interfaces, and related components before concluding that the required change is isolated.
- understand why something broke before changing it. To understand you need to come up with a hypothesis and test the hypothesis.
- use debugger if possible

### Linting

Always keep a look at linting warnings and errors. Fix them immediately.

### Evidence-based Statements

Claims about correctness, test status, coverage, security, compatibility, architecture conformance, or feature behavior require explicit evidence.

For example:

Bad:

“The API is backwards compatible.”

Good:

“I verified backwards compatibility by running X tests against versions A/B. Results: ...”

Similarly:

- “The feature works” → needs E2E evidence
- “93% coverage” → needs coverage report
- “No security issue” → needs security scan/review evidence
- “This dependency supports X” → needs documentation reference
- “The architecture matches the spec” → needs explicit comparison

### Security

- Never hardcode secrets
- Validate external input
- Escape shell arguments
- use the principle of least privilege

### Testing Order

- Implement first (Always Required): Unit tests 
- Implement second (Always Required): Integration tests 
- Implement third (only required if the project already provides an external API or GUI): End-to-end tests

Skipping any required level = Not Complete

> This is the end of the AGENTS.md file

---

## Issue Flow (tool-agnostic)

Notes for implementation:

- each unique step is a slash command
- each step is done by a single agent
- the agent can decide to go back to a previous step (e.g., encountered a problem that requires change to spec)
- most steps have a planning sub-step performed at the beginning with the harness' plan mode which generates `.agent/session/plan.md`
- the agent can go back from a step to a previous step if necessary to fix an issue created earlier. Step jumps must be tracked in `issue_flow.md`, including the issue name, issue number, timestamp of creation and completion, and the jump/deviation. The issue flow is mostly sequential, but step B will hold parallel sequential paths, one for each core implementation task.
- when a new feature (tackling a new issue) starts, first search for any completed `.agent/flow/completed-issues/completed_issue_flows/issue_flow_[i].md` and use it when relevant as historical context for the issue being handled
- session summary hook: when a session ends, store a summary of key things done, key problems encountered, tips/learnings/todos in the current issue's `issue_flow.md` under the step the agent was in, using this JSON form:
  `{"key things done": "placeholder", "key problems encountered": {"problem": "placeholder", "solved_or_not": placeholder, "tips for next agent working on this": "placeholder"}, "learnings": "", "todos": "placeholder"}`
- approval gates mean that either the user (developer) or a specific AI agent needs to give approval to continue the flow
- "AI Review" means the same AI that is coding reviews its own work
- "Independent AI Reviewer" means that a different model with fresh context must be used
- when a session ends, a hook must be called to: (1) empty all files inside `.agent/session`; (2) recursively analyze `.agent/session-persistent-candidate` and transfer useful validated knowledge to `.agent/persistent/knowledge/non_obvious_conjectures_and_facts.md`, then empty `.agent/session-persistent-candidate`; and (3) assess whether harness-related files need improvement or bloat removal using the `harness-eval` skill and available traces/past mistakes/feedback. Before modifying harness-related files, ask for user approval and provide the evaluation evidence. If changes are made after evaluation, git commit them.
- at the start of any slash command, a hook must git add & commit if there are uncommitted changes that should not be left outstanding. **Exception: when the slash command is starting an inferred harness A/B evaluation of the current uncommitted working tree, do not auto-commit the harness changes before the evaluation snapshot is created.** In a new session without context, use `.agent/flow/current-issue/flow/issue_flow.md`'s last progress data to infer a suitable commit message.
- if stuck in a loop, developers should try changing the implementer model from the default to another model
- ensure language-specific (in our case Go and TypeScript) linters and formatters are running continuously on every file edit, using OpenCode hooks
- have a git pre-commit hook that builds (successfully without warnings), runs tests (all must pass), and verifies test coverage is above 90%

### Issue Flow (note: bug, refactoring or performance issue handling may skip multiple steps that are unnecessary)

A: Issue-specific Spec & Core Implementation Tasks Plan

0. **/ensure-my-understanding Understand the codebase, then grill me with questions to see if I really understand the codebase.** Make high-level questions (e.g., decisions chosen, project structure, tradeoffs, architecture) and low-level questions (e.g., what a specific file/function/class is for). [AI Approval Gate]

1. **/start Start Issue Handling**: point to the GitHub issue; create a satellite branch with an appropriate name according to the branching strategy file; study the repo; ask user clarifying questions about the problem and solution. This step ends when a common problem/solution understanding is reached with the user.

2. **/reproducebug** only for bug issues: reproduce the reported bug by writing and running a failing test and confirm that it fails for the expected reason described in the bug issue.

3. Loop until **3.ii and 3.iii** are successful [User Approval Gate with AI Security Reviewer Suggestions]
   i. **/spec Build Issue-specific Spec**: create `prd.md`, `architecture.md`, and `tech_stack.md` under `.agent/flow/current-issue/issue-spec_[i]/`.
   ii. **/reviewspec Review Spec**: use an Independent AI Reviewer to check consistency with the Global Spec (Founding Doc + Roadmap + Tech Stack), catch omissions, and flag problems in the Global Spec. If a Global Spec problem is found, modify the Global Spec first, then update the issue-specific spec. Store the final spec review in `.agent/session-persistent-candidate/reviews/final_spec_review_[timestamp].md`. [User Approval Gate]
   iii. **/specsecreview Specialized Spec Security Review**: flag critical security problems and warnings. Store the final security spec review in `.agent/session-persistent-candidate/reviews/final_spec_sec_review_[timestamp].md`.

----<<separate terminal block (reset context)>>----

B: Core Implementation

4. **/make-tasks-plan Convert the spec into a dependency graph of core implementation tasks** in `.agent/flow/current-issue/flow/core-implementation-tasks-plan.md`. Each node should represent a coherent behavioral/functional slice rather than a component-only task. [User Approval Gate with Independent AI Plan Reviewer Suggestions]

5. **/tasksplangrillme Grill User on the task plan** until the user has a full understanding of it. [AI Approval Gate]

6. **Core Implementation Loop**: process each pending core implementation task in the graph until the graph is fully complete. Each task follows steps 7–12 below.

---- New Session (reset context) ----

B1: Common Scaffold

7. **/boilerdep Define Allowed Scaffold Dependencies**: choose dependencies for the common scaffold. Dependencies must be actively maintained. A dependency may be preferred when it has important capabilities the alternatives lack, is significantly more popular, is significantly easier to use, or is significantly older. [User Approval Gate with AI Review Suggestions]

8. **/boiler Setup/Modify the Common Scaffold**: make/remake `plan.md` first; create or modify the required project skeleton, dependencies, build/test/package automations, and other foundation pieces. Review against Issue-specific and Global Spec. [User Approval Gate with AI Review Suggestions]

B2: Tests & Logic

9. **/writetests Write/Modify Functional Tests**: make/remake `plan.md` first. Write/modify unit and integration tests and review them against the Issue-specific and Global Spec. Flag inconsistencies rather than silently resolving spec conflicts. [User Approval Gate with Independent AI Reviewer Suggestions]

10. **/testtests Test the Functional Tests with Placeholder Feature Code**: all functional tests must fail in this phase, while test files themselves must be valid and runnable. Ensure strict test coverage of the intended core logic (above 90% of the entire codebase once implemented, ideally near 100%).

11. Loop through 11.i–iii until the feature implementation satisfies the functional tests and validation requirements. [User Approval Gate with Independent AI Reviewer Suggestions]
   i. **/featdep Define Allowed Feature-code Dependencies** with explanation of why each is used. Apply the same maintenance/popularity/ease/capability criteria as scaffold dependencies. [User Approval Gate with AI Review Suggestions]
   ii. **/feat Write Feature Code** using only the allowed feature-code dependencies; make/remake `plan.md` first. Implement the smallest coherent behavioral slice and try to make multiple related tests pass at a time. Review against Issue-specific and Global Spec and Design System when GUI work is involved. [User Approval Gate with Independent AI Reviewer Suggestions]
   iii. **/test Build and Test Feature Code**: make/remake `plan.md` first; build and run functional tests; generate test and coverage reports; review against Issue-specific and Global Spec; repeat until all required tests pass.

12. **/refact-coreimplementation-task-if-necessary**: make/remake `plan.md` first; evaluate refactoring opportunities limited to the last completed core implementation task. Refactor only when it preserves behavior and improves clarity, quality, or maintainability. After each refactoring, evaluate whether it is actually better; otherwise keep the previous version. [User Approval Gate] [AI Approval Gate]

----<<separate terminal block (reset context)>>----

13. **/grillme Grill User on the Latest Changes**: user reviews the code and asks questions; the agent grills the user until the user has full understanding of the changes. [AI Approval Gate]

----<</separate terminal block>>----

14. **/stripdebuglogs**: make/remake `plan.md` first; remove debug logs if there are any. [User Approval Gate with AI Reviewer Suggestions]

15. **/determine-if-core-implementation-task-is-done**: determine whether the current core implementation task objectives were actually achieved and whether the task graph still holds. Do not rely on historic test output; rerun the entire suite for the current task and previously completed core tasks to catch regressions. Only mark the current task done when tests reflect the spec and pass.

---- New Session (reset context) ----

16. **/end-to-end-testing**: test end-to-end to make sure the issue was completely handled. If GUI is involved, include GUI usability/visual testing that emulates real user behavior. Check that E2E tests reflect PRD requirements and that the issue spec remains consistent with the Global Spec.

---- New Session (reset context) ----

C: Code Review

17. Loop until code review succeeds [User Approval Gate with Independent AI Reviewer Suggestions]
   i. **/review Independent Code Review**: make/remake `plan.md` first; review maintainability, modularity, test coverage, file sizes, test/spec compliance, and GUI compliance where relevant. The last thing to do is run mutation testing. Also identify omissions and spec/design problems. Store the final code review in `.agent/session-persistent-candidate/reviews/final_code_review_[timestamp].md`.
   ii. **/redo Make Necessary Code/Test Changes**: make/remake `plan.md` first; implement justified changes, build and test; review again against Issue-specific and Global Spec. [User Approval Gate with AI Reviewer Suggestions]

---- New Session (reset context) ----

D: Security Review, Documentation & PR

18. Loop until security review succeeds [User Approval Gate with AI Reviewer Suggestions]
   i. **/secreview Specialized Security Review**: flag critical problems and warnings. Store the final security review in `.agent/session-persistent-candidate/reviews/final_code_sec_review_[timestamp].md`.
   ii. **/redo Make Necessary Code/Test/Docs Changes**: make/remake `plan.md` first; implement changes, build and test, then review against Issue-specific and Global Spec. [User Approval Gate with AI Reviewer Suggestions]

19. Loop until E2E verification succeeds
   i. **/end-to-end-testing**: run the final E2E verification, including GUI testing when relevant, and verify E2E tests reflect PRD requirements and Global Spec consistency.
   ii. **/redo Make Necessary Code/Test/Docs Changes** as required. [User Approval Gate with Independent Reviewer Suggestions]

----<<separate terminal block (reset context)>>----

20. **/grillme Grill User on the Final Changes** until the user has full understanding. [AI Approval Gate]

----<</separate terminal block>>----

21. **/write-pr**: check `.agent/flow/current-pr/pre-pr-checklist.md`, confirm it is satisfied, then write the PR to `.agent/flow/current-pr/pr.md`. The PR is not considered final until the documentation step is complete.

22. **/document Document**: (A) Contributor Documentation: explain how to understand the codebase using a step-by-step tutorial from high level to low level; (B) Final User Documentation only if the product is already usable. [User Approval Gate with Independent AI Reviewer Approval Suggestions]

23. **/send-pr**: push and open the draft PR using `.agent/flow/current-pr/pr.md`.

---- New Session (reset context) ----

E: Make Fixes Based on PR Reviews and/or CI Failures Until PR Is Merged

24. Loop until all required external review/CI failures are addressed [User Approval Gate with AI Reviewer Suggestions]
   i. **/prreviews**: after the user manually checks a PR review or CI failure notification, read the PR reviews and CI runs from GitHub and write them locally in the dedicated current-issue PR-review area.
   ii. **/review Independent Code Review**: make/remake `plan.md` first; review maintainability, modularity, test coverage, file sizes, tests/spec compliance, and run mutation testing last. Store the final code review in `.agent/session-persistent-candidate/reviews/final_code_review_[timestamp].md`.
   iii. **/redo Make Necessary Code/Test/Docs Changes**: apply justified fixes, build and test, and review against Issue-specific and Global Spec. [User Approval Gate with Independent AI Reviewer Suggestions]

---- New Session (reset context) ----

25. Loop until security review succeeds
   i. **/secreview Specialized Independent Security Review**: store final security review in `.agent/session-persistent-candidate/reviews/final_code_sec_review_[timestamp].md`.
   ii. **/redo Make Necessary Code/Test/Docs Changes** and validate. [User Approval Gate with Independent AI Reviewer Suggestions]

26. Loop until final E2E verification succeeds
   i. **/end-to-end-testing**: rerun final E2E verification, including GUI testing when relevant.
   ii. **/redo Make Necessary Code/Test/Docs Changes** and validate. [User Approval Gate with Independent AI Reviewer Suggestions]

----<<separate terminal block (reset context)>>----

27. **/grillme Grill User on Changes Since Last Grillme** until the user has full understanding. [AI Approval Gate]

----<</separate terminal block>>----

28. **/write-pr**: recheck `.agent/flow/current-pr/pre-pr-checklist.md`, then regenerate `.agent/flow/current-pr/pr.md` so it reflects the final state.

29. **/document**: update contributor documentation and usable final user documentation. [User Approval Gate with Independent AI Reviewer Approval Suggestions]

30. **/send-pr**: push/open/update the draft PR using `.agent/flow/current-pr/pr.md`.

The human user is responsible for checking when the PR is ultimately merged.

> This is the end of the Issue Flow section.

## AI-assisted Coding Tech Stack

External Vendor Requirements: Opencode Go Subscription, Claude Credits, Github Repo

- Agent-native Editor/Terminal: Orca (only when tackling multiple issues in parallel)
- Terminal Harnesses: OpenCode + failproofai + claude-tap, open-code-review and openwiki
- IDE (for better introspection + manual editing): VSCode (with language-specific linters and formatters plugins running continuously on every file edit)
- LLM Provider Subscription: **OpenCode Go
- Models: (0) Planning: Current great coding model thats not so expensive with medium reasoning; (1) Spec, scaffold and Tests: Current great coding model thats not so expensive with high reasoning; (2) Core coding: Best coding quality/price model with medium reasoning; (3) Independent Spec Review: Best Coding Model with high reasoning; (4) Plan Review: Current great coding model thats not so expensive with high reasoning; (5) Security Review: Best Review Model with high reasoning; (6) end-to-end testing: Current great coding model thats not so expensive with high reasoning; (7) sub-agents model: Best Coding Model with high reasoning; (8) Code Review: 2 not so correlated great Review Models that are not so expensive with high reasoning as candidate generators and open-code-review using best review model with high reasoning as final decision maker; (9) Grillme: Best coding quality/price model with medium reasoning
- Local Routing: opencode-model-router (opencode plugin)
    - Fast Model: qwen2.5-coder:7b
    - Medium Model: current coder model
    - Heavy Model: current coder model
- Issue Flow: Start with just making each step a slash command, and leave to the developer to follow the steps (note: this doesnt enforce step execution, se requires developer commitment). Note: each slash commands should remember to update the issue resolving progress at the end (issue_flow.md file) or, if its the first step, create the file if its still not created). Each slash command should have a simple name.
- OpenCode Plugins: graphify, rtk, opencode-quota, opencode-model-router, opencode-trace
- OpenCode MCPs: Github MCP
- Custom Freelunch OpenCode sub-agents: security-specialist, code-review-and-refactoring-specialist, testing-specialist, debugging-specialist. Tip: Agency-agents repo provides some sub-agents out-of-the-box.
- Custom Freelunch OpenCode Slash Commands: one for each unique step of the issue building flow
- Dependency Docs: a dependency_docs.md under ./docs with entries in the form "- <dependency>: <docs_link>" for all direct dependencies (not dependencies of dependencies). Pinned versions used in the project can be seen in the lock file of the virtual dev environment tool.
- Skills: 
    - Custom Skills: created on-demand, under a `custom_skills` folder, via manual creation or via `skill-creator` to avoid having to repeat the same solution process over and over. Make these custom skills:
        - harness-eval (evaluates proposed harness changes on the current codebase, explained better in the end of this doc)
        - ensure-my-understanding (continually ask questions of the latest changes to codebase to me, to see if i understand the codebase. Always give score my answers and give feedback to it. Only stop when you feel i understand the codebase. The user can also specify specific files for you to grill him about instead of the entire codebase)
        - understand-external-codebase (1. Build a doc explaining in detail the characteristics and internals of an external github codebase; 2. Add to this doc an explanation of where and why this codebase can be helpful as a reference for ideas/patterns for the current project being built)
        - update-fixed-context (1. Infers new useful knowledge from ./.agent/persistent/knowledge/mistakes.jsonl and ./.agent/flow/completed-issues/completed_issue_flows; 2. Add this new useful knowledge to AGENTS.md if its not already there)
        - make-core-implementation-tasks-plan (transforms the spec into a graph of tasks, where: (1) each task can depend on other tasks being already done or not depend on any; (2) the tasks should not be of the form "one task implements each component that will be needed in this feature, e.g., one task for the backend, another for database, another for gateway and another for frontend", the tasks should be done in the form of "one task implements a slice (governed by on or more integration/end-to-end tests) of multiple components, e.g., this task implements a functionality slice of frontend, backend, gateway and frontend that together brings us one step closer to our end goal and guarantee rich cross-component feedback along development". The core implementation tasks plan needs to be stored in .agent/flow/current-issue/flow/core-implementation-tasks-plan.md. The core-implementation-tasks-plan.md. file should have the graph structure of the plan, where each node is a task. For each node there is also a pending/in-progress/done checkbox. Do not confuse with plan.md which is a per-step ephemeral small plan for step execution.)
        - document using openwiki and with the following guidelines in ./.openwiki/INSTRUCTIONS.md: should first check inconsistency (global spec is king), staleness and incompleteness of existing documentation (if any) and then update/create: (1) Contributor Documentation: visualization of repository structure explaining succinctly each directory and file, (1.2) Step by step contributor tutorial to help a newcomer understand the codebase; (2) User Documentation (only do after first version 0.1.0 is released): (2.1) User API Reference. (2.2) User step by step tutorial starting from scratch; (2.3) User guides to do common stuff; (2.4) FAQ. 
        make sure the documentation explains well things that I usually have a hard-time understanding.
        - ui-taste (UI Taste gives Claude a visual sense of taste. Instead of relying only on abstract design principles, the skill provides curated examples of bad, good, and stellar GUIs across different application categories and problem modes, including screenshots and their underlying HTML/CSS. This gives the agent an understanding of what makes GUIs look good. The agent should launch the current GUI, identify the biggest visual shortcomings, and iteratively improve them. The goal isn't to force a particular design style—it is to help Claude distinguish "functional but mediocre" from "genuinely beautiful and easy to use", giving coding agents a practical visual benchmark for judging their own work.)
    - Use existing skills: visual-recap (ideally also running in github actions for every PR), skill-creator, i-have-adhd, quick-recap, mutation-testing, chrome-devtools-cli, optimize-anything, grillme (every grillme run should log all the questions, answers and feedback gave to the user into a .agent/persistent/user-grills/grill[i].md where i is the id of the grill and the file should have timestamp, commit, what the grill was about and grill score in the beginning of it. Ever grill should start by looking at the commit, what was grilled in the last grill and user-codebase-questions.jsonl file), lavish-axi, code-review-and-quality, api-and-interface-design, browser-testing-with-devtools (only when working with frontend part), security-and-hardening, cc-skills-golang, maintainable-typescript (only when working with frontend part), improve-codebase-architecture, screenshot (only when working with frontend part), extract-design-system (only when working with frontend part), frontend-design (only when working with frontend part).

## How Code Review is Done

1. Start a new session with opencode: do code review with Model A and store the review in .agent/session-persistent-candidate/reviews/[A]_code_review_[timestamp].md, where A is a placeholder for the actual mode name, and timestamp is placeholder for the actual timestamp
2. Start a new session with opencode: do code review with Model B and store the review in .agent/session-persistent-candidate/reviews/[B]_code_review_[timestamp].md, where B is a placeholder for the actual mode name, and timestamp is placeholder for the actual timestamp
3. Start a new session with open-code-review: do code review with Model C explicitly telling it to look at the candidate problems flagged inside .agent/session-persistent-candidate/reviews/ folder and store the resulting code review inside .agent/session-persistent-candidate/reviews/final_code_review_[timestamp].md

## Token Efficiency Laws

- Avoid small actions → batch related small tasks into larger coherent work units (use a todo_buffer.md to store all todos and then batch them before prompting the harness)
- Don't derail the agent from its main goal → context spent on unrelated work is expensive and increases context pollution.
- Use a graph/codebase-understanding tool → avoid repeatedly spending LLM tokens rediscovering repository structure and relationships.
- Related big tasks in the same session, unrelated new stuff gets its new session
- Start a new session when context becomes sufficiently polluted → carrying a huge amount of irrelevant history can become more expensive than rebuilding a clean context.
- Avoid rambling/random studying with the agent, do all of this in ChatGPT/Gemini/Grok webpages
- Stay in the same session while the context is still useful → preserving cached/reusable context avoids paying to rebuild understanding.
- Until you reach scaffold code with tests, dont switch models
- Always mention files with @ for the agent to look/modify instead of letting the agent wander the repo for that file
- Use something like rtk to compact tool outputs
- Have a router + local model for doing simple stuff

## .Agent Directory Structure Guide Doc (./.agent/directory_structure.md)

## `.agent/` Directory Structure

The `.agent/` directory contains the AI agent's workflow state, persistent knowledge, session state, issue-specific process artifacts and runtime-created tools. It is an internal directory used by the coding agent and should not contain product source code.

### Directory Structure (./.agent/directory_structure.md file)

```text
.agent/
├── flow/
│   ├── current-issue/
│   │   ├── raw_github_issue.md
│   │   ├── flow/
│   │   │   ├── issue_flow_[i].md where i is the issue number
│   │   │   └── core-implementation-tasks-plan.md
│   │   └── issue-spec_[i]/ where i is the issue number
│   │       ├── prd.md
│   │       ├── architecture.md
│   │       └── tech_stack.md
│   ├── current-pr/
│   │   ├── pre-pr-checklist.md
│   │   └── pr.md
│   └── completed-issues/
│       ├── completed_issue_flows/
│       │   └── issue_flow_[i].md where i is the issue number
│       ├── completed_issue_specs/
│       │   └── issue-spec_[i]/
│       │       ├── prd.md
│       │       ├── architecture.md
│       │       └── tech_stack.md
│       └── completed_prs/
│           └── pr_[i].md where i is the issue number
│
├── persistent/
│   ├── knowledge/
│   │   ├── non_obvious_conjectures_and_facts.md
│   │   └── mistakes.jsonl
│   └── user-understanding/
│       ├── user-codebase-questions.jsonl
│       └── user-grills/
│           └── grill[i].md
│
├── session/
│   ├── plan.md
│   ├── todos.md
│   ├── problem_stuck.md
│   └── debug-logs/
│
├── session-persistent-candidate/
│   ├── knowledge/
│   │   ├── non_obvious_conjectures.md
│   │   └── exploration_findings/
│   │       └── <name-of-exploration>_[timestamp].md
│   └── reviews/
│       ├── final_spec_review_[timestamp].md
│       ├── final_spec_sec_review_[timestamp].md
│       ├── final_code_review_[timestamp].md
│       └── final_code_spec_review_[timestamp].md
│
├── created_tools/
│
├── directory_structure.md
└── harness-evals/
    └── <historical harness evaluation results are stored here if configured>
