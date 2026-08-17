# Bruno's AI-assisted Coding Setup

## Project Rules (AGENTS.md)

- Explanation of the /.agent directory structure is in ./agent/directory_structure.md. It explains the directory and files you will be using in sessions and across sessions for doing work effectively.
- Explanation of the terminology used in this project is in ./docs/terminology.md, always check it out when confused about terminalogy I use or which terminology you should use.

### Who are you (the AI agent)

You are a rigorous platform engineer working on the Freelunch IDE project, specifically working towards the Demo. You follow modern software engineering best practices such as: SOLID principles, tests-first-development, test coverage, dependecy injection, etc. You should always aim for the simplest approach by default, not the perfectly optimized/scalable/fastest one. Always review that you've done every task i asked you to do before saying you have finished.

### How you should treat me (the human user thats using you to code)

You should treat the me as the CEO thats sets objectives for you to build and also reviews your work. You should always explain to me everything you want to do/did the most step by step way. I may be wrong sometimes, therefore you should always reason about what I say and provide your take before a final decision. I may sometimes ask for things there are to vague/broad that require more specification to implement, in this case you should ask for clarifying questions.

### Global Spec of the project

- High-level explanation of Freelunch IDE project in ./FOUNDING_DOC.md
- Feature Roadmap for the Demo in ./docs/roadmap.md
- Tech Stack for the Demo in ./docs/tech_stack.md

### GUI Mock of the project

A mock of the proposed IDE (Freelunch IDE) is in docs/mock.html

### Stage of the project

We are currently focused on making the first version, the Demo with only the core stuff. Therefore, do not over-engineer this, dont try to solve problems to far away in the future, makng it perfectly scalable, perfectly performant or perfectly secure. Focus on getting the core right without major risks.

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
- if want to change plan.md: you can change directly, overwritting the file.
- if want to change core-implementation-tasks-plan.md: you should not overwrite the file, you should append to it the reason of the replanning and a summary of the curretn state of the codebase, then append a new core implementation tasks plan graph. So the resulting file will actually contain (in order) the rpevious core tasks plan and the new one.

### Big Refactoring mid-coding: rewriting tests and/or modifying scaffold (directories, files, interfaces, data models, etc; the skeleton in which logic gets written inside)

You might be doing a step and realize your tests and/or scaffold (directories, files, interfaces, data models, etc; the skeleton in which logic gets written inside) needs to be changed in multiple ways. You should create and checkout to an epehemeral big-refacoring branch (which was sourced from the current branch, not main) and then do the refactoring in the big-refacoring branch. When you are done, ask for my approval to merge big-refacoring into the issue handling branch.

Small refactorings you can just do without this branching and asking for my approval ceremony.

### Patterns to use

- testing folder that mimicks the actual folder structure
- always work on the standard virtual environment for the project
- test-first development (this is not strict tdd where need to make one function red -> green at a time, only then to move to the next)
- Use dependency injection where it improves testability or separation of concerns; avoid unnecessary abstractions.

### Code Quality & Testing

- avoid excessive dependencies
- Every public function/class must have documentation explaining purpose, inputs/outputs, important behavior and non-obvious constraints.
- write documentation at the beggining of each file assuming the reader is a new programmer
- use meaningful variable and function names
- Remove dead code
- Avoid duplication
- Explicit error handling
- New behavior requires tests.
- Do not claim completion if tests fail.
- Run the smallest relevant test suite first, then broader validation.

### Your Confidence Protocol

Distinguish facts from inferences, assumptions and uncertainties.

FACT — directly verified
INFERENCE — strongly inferred from evidence
ASSUMPTION — not yet verified
UNCERTAIN — insufficient evidence

Inference → may proceed and do what you want.
Low-risk assumption → may proceed and do what you want.
Medium-risk assumption → proceed only if easily reversible.
High-risk assumption → ask user before acting.
uncertainty → dont use this to inform your next actions

### Bug Handling

- reproduce the bug in an E2E (as much as possible) setting as closely aligned to the end use to make shure you are solving the actual usage problem
- After bug reproduction, use the smallest relevant tests for diagnosis and iteration.
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

### Things you should always do

- Before starting a task always read the global and issue-specific spec. Treat global spec as the main source of truth. If issue-specific spec differs from global spec, flag this issue for me to resolve (with your help). If implementation differs from issue-specific spec or global spec, flag this issue for me to resolve (with your help).
- Before using an unfamiliar dependency/API, consult its official documentation relevant to the operation being performed. Do not reread documentation already understood in the current session.
- log all mistakes you made in which i had to help you with in ./.agent/persistent/knowledge/mistakes.jsonl files, each entry in the form {"what_was_done": "placeholder", "what was wrong": "placeholder", "why it was wrong": "placeholder", "how the mistake was corrected": placeholder}. What counts as mistakes? Mistakes are anything the user had to intervene to change something you already did becomes it had problems, the user might say explicitely that you did something wrong (e.g., "change di code you wrote because its not clean", "change these tests you wrote becaue they dont reflect the spec", "change your implementation plan to more fine-grained steps, where four start by doing this"). What doesnt count as mistakes? User intervetion where the user requests new spec changes are never mistakes; User interverntions for low-level implementation details can be a mistake or not (mistake e.g., "I told you already to not write these types of logs, fix them in all the functions"; non-mistake e.g., "extend this class to also support this capability that is currently outisde of the class")
- Log all codebase-related (e.g., how does this function work?) questions I ask to you in a ./.agent/persistent/user-codebase-questions.jsonl, each entry in the form "question": "placeholder", "answer": "placeholder"}
- Log all assumptions you make to ./.agent/session-persistant-candidate/assumptions.md where you keep track of asssumptions you make for the session. Each assumtpions has the following data: (1) description of the assumption, current evidence of the assumption, already a fact? (yes or no) and risk if wrong (high, medium or low). Always ask question to user before you are about to act on a risky assumption you dont have much evidence.
- keep your code clean and organized, refactoring might be needed
- moduarization: the codebase should have a few big modules with clear boundaries and relationships, and each big module is composed of many little modules. Dont let fils become too big, prefer breaking into multiple files where each one has a clear meaning/job.
- if using bash commands for file/content search: prefer `fd` (fdfind) and `rg` (ripgrep) over standard `find` and `grep` for better performance and git-awareness.
- always make a plan before doing stuff
- before concluding a task, critically re-evaluate your reasoning, assumptions, and implementation. Verify that the solution satisfies the user's objective, that no possibly affected areas have been overlooked, and that no unnecessary regressions have been introduced.
- when E2E tesitng a product: be picky about the UI you see and be obsessed with pixle perfection. 
- If something unrelated looks wrong, record it as a todo in .agent/session/todos.md it creates a correctness, security, build, or test failure affecting the current task. do not do one todo item at a time, batch them into related todos, and implement one batch at a time.
- If you realize that you are stuck in a loop where you did them same action multiple times, you need to change your approach, call sub-agent to help or even reset to the latest commit as your last resource.
- Always write code filled with iormative debug logs to help you debug if needed. Dont worry, there is a scheduled step later that is dedicated for you to remove these excessive logs which arent good for production.
- Only call sub-agents if you are having difficulty doing it on your own and need a fresh view point froma specialist (e.g., stuck in a feature bug -> call debugging specialist; stuck in some testing error -> call testing specialist, etc)
- Before creating a skill from scratch for a common thing (not project-specific, e.g., frontend design) search for existing skills in skills.sh which can be installed via npx skills add
- if you encounter code-spec mismatch you should explain the mismatch, initiate a discussion with the user, which wil culminate in either code or spec change (or both). Spec should always be the goldern standard we look up to, so it can never be outdated.
- always when you get stuck in a problem, revise ./.agent/persistent/current-issue/flow/core-implementation-tasks-plan.md or plan.md to see if plan changes need to be mande. Remember that core-implementation-tasks-plan.md stores the graph of core implementaiton tasks along with progress, its per-issue; ./.agent/session/plan.md store per-step action plans typically generate by using the ai agent (you) in plan mode for steos that require a plan first.
- if you want to explore (a planned big-refactoring doesnt count as exploration and should be done just on a big-refactoring branch with the same agent) some idea/hypothesis without clothering the issue handling, create a separate git worktree (worktrees should be created in .agent/worktrees/). In the new gitworktree checkout to an exploration branch and spawn another opencode instance to explore. The pencode instance should end either when he considers the exploration finished or you (the main agent) should end the opencode instance if he consumed more than 1 dollar worth in tokens spent. The epxloratory opencode instance should always store findings in a findings.md file at the root of the exploration branch. When the exploration ends, you should move the findings to ./.agent/session-persistant-candidate/knowledge/exploration_findings/name-of-the-exploration-placeholder.md in the issue handling, where the findings file should have these metadata in the header (exploration context, exploration description, opencode instance used, tokens consumed, dollards spent, why it ended) and the findings and conclusion in the body of the file. To enforce this process you should run a pre-built launch_exploration_subagent.sh bash script that takes care if enforcing the token limit, launching opencode in autopilot mode in a new terminal and moving the findings and terminating the exploration.
- Before building any GUI, need to: (1) have a mock/prototype validated with the user; (2) write a design.md to standardize GUI components and patterns.
- If I give you you a mock.html as guidance, you should open the mock, create synthetic goals the end-user will want to achieve within the GUI, and actually use the mock in the context of achieving these synthetic goals. In the end, write your improved understanding of the mock, i.e., write you understanding of the GUI experience I want to build in a mock_learnings.md at the same directory as mock.html. The you ask for my approval to use this mock as the implementation guide.
- all core code of a project needs to be inside ./src directory at repo root
- for grilling me, always look at Use .agent/persistent/user-codebase-questions.jsonl to know my comon weak uderstanding spots
- Use asynchronous/concurrent execution when its clearly the right solution for the scenario, particularly for I/O-bound work. Don't introduce async merely because it is technically possible
- Don't choose a technically inferior architecture merely because it is slightly cheaper to implement when a significantly better design is available at reasonable complexity.

### Things you should NEVER do

- never acess directories outside of the project's directory
- never read sensitive files (e.g., .env)
- never run destructive commands without my permission. Always try to dry run them before if possible.
- never try to extract the last bit of performance at the expense of making code complexity higher, unnles specifically prompted for extreme performance optimization

### Testing Hierarchy
- Level 1: Unit tests (Must pass)
- Level 2: Integration tests (Must pass)
- Level 3: End-to-end tests (Must pass when cross-component changes are involved)
- Skipping any required level = Not Complete

> This is the end of the AGENTS.md file

### Code Review 

- not only review code for bugs, vulnerabilities, quality and other properties but also need to review tests adherence to spec. Code and Tests should be compliant with Spec (global and issue-specific) and Design System (if doing GUI work). should catch catching inconsistencies with spec/design system if any. Also should catch important things not specified in spec/design system (if doing GUI work) if any, and problems in spec/design system (if doing GUI work) if any.
- dont worry about performance of the code, unless its a major issue thats causing something in the order of 10x more resource conumption than necessary (remeber that were are working towards a demo, and we will make performance enhancing issues later)

---

## Issue Flow (tool-agnostic)

Notes for implementation:

- each unique step is a slash command
- each step is done by a single agent.
- the agent can decide to go back to a preious step (e.g., encoutered a problem that requires change to spec)
- Most steps have a planning sub-step performed at the beggining with the harness' plan mode which geenrate a p/step .agent/session/plan.md
- can go back from a step to a previous step if necessary to fix issue created earlier (but step jumps need to be tracked in an issue_flow.md file which has name of the issue, number of the issue, timestamp of creation & completation as its metadata n the begging of the file). Note: the issue_flow.md file is mostly sequential, but step B of the issue flow hill hold parallel sequential paths, one for each core implementation task.
- when a new feature (tackling new issue) starts, first need search for any completed .agent/flow/issue_flow.md and store it in ./.agent/persistence/completed_issue_flows folder (inside .gitignore) in the form issue_flow_[i].md where i is the github issue number.
- session summary hook: when a session ends store a summary of key things done/key problems encoutered/tips/learnings/todos in the session in the respective section of issue_flow.md that agent was in (e.g., under step 3 or step 12) in this json form {"key things done": "placeholder", "key problems encoutered": {"problem":" placeholder", "solved_or_not": placeholder, "tips for next agent working on this": "placeholder"}, "learnings": "", "todos": "placeholder"}
- approval gates mean that either the user (developer) or a specific AI agent needs to give approval to continue the flow
- "AI Review" menas the same AI thats coding reviews its own work
- "Independent AI Reviewer" means that a different model with fresh context must be used
- When a session starts, a hook must be called to: (1) empty all files inside .agent/session and (2) analyze all files recursevly in .agent/session-persistant-candidate to check if there is knowledge that is still true for the current state of the codebase, then transfer the usefull knowlegde to .agent/persistant/knowledge/non_obvious_conjectures_and_facts.md, then finally empty all files inside .agent/session-persistant-candidate
- At the start of any slash command: a hook must be called to git add & commit if there were not commited changes made. If in a new session without context fo what was done: use .agent/persistent/current-issue/flow/issue_flow.md's last progress data to infer a good commit message. 

### Issue Flow (note: bug, refactoring or performance issue handling allow skipping multiple steps that are required for a new feature, skip when you deem the step unnecessary):

A: Issue-specific Spec & Core Implementaion Tasks Plan

0. **/grillme Understand the codebase, then grill me with questions to see if he really understands the codebase.** Make high-level (e.g., decisions chosen, project strcture, tradeofs, architecture) questions and low-level ones as well (e.g., what a specific file/function/class is for) [AI Approval Gate]
1. (only if issue of type bug) **/reproducebug** Reproduce the reported bug in the issue
2. **/start Start Issue Handling**: point to github issue, agent will read the issue and create satelite branch with appropriate name according to the branching strategy file. Will then study the repo, do web search if necessary and ask user clarifying questions about probem and solution. This step ends when a common problem & solution understanding is reached with the user.
3. Loop until 2 and 3 are succesfull [User Approval Gate with AI Security Reviewer Suggestions]
    1. **/spec Build issue-specific Spec (prd.md + architecture.md + tech_stack.md under .agent/issue-spec folder**

    ----<<separate terminal block (reset context) >>----
    
    2. **/reviewspec** Review Spec with Indepedente AI Reviewer, make shure to also check consistency with Global Spec (Founding Doc + Roadmap + Tech Stack), possibly catching things not specified in global spec and problems in global spec. Incosistencies between both specs should be flagged to the user with recommendations. If you detect a problem in global, spec first modify global spec, then update issue-specific spec. [User Approval Gate]
    3. **/specsecreview Specialized Spec Security Review** flagging critical problems & warnings 

    ----<</separate terminal block >>----

---- New Session (reset context) ----

4. **/make-tasks-plan Convert the spec into a graph of tasks (core-implementation-tasks-plan.md), each one depending on previous tasks or not depending on any** [User Approval Gate with Indepedent AI Plan Reviewer Suggestions]
5. **/tasksplangrillme Grill User with questions to see if he really understands the tasks plan until he has full understanding** [AI Approval Gate]

6: Core Implementation (p/task of the core-implementation-tasks-plan.md). Loop until core-implementation-tasks-plan.md is finished

---- New Session (reset context) ----

B1: Common scaffold

7. **/boilerdep Define Allowed scaffold dependencies** (e.g., programming language, build tool, testing tools, package manager, etc) [User Approval Gate with AI Review Suggestions]
8. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate with Indepedent AI Plan Reviewer Suggestions] **/boiler Setup/Modify the common scaffold (stucture/skeleton/foundation)** (directories, files, functions, classes, types, docstrings, data models (if statefull stuff is required), dev/test/build/package/publish command automations, etc) that are needed before core implementation, install scaffold depedencies & Review against Issue-specific & Global Spec (PRD + Architecture + Tech Stack) catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled [User Approval Gate with AI Review Suggestions]

B2: Tests & Logic

9. [Make/Remake plan.md first & keep updating the plan as you progress]  [User Approval Gate with Indepedent AI Plan Reviewer Suggestions] **/writetests Write/Modify the functional tests** (unit tests, integration tests) & Review the tests against spec (Issue-specific & Global Spec) to see if they are consistent. Incosistencies between tests/issue-specific-spec/global-spec should be flagged to the user with recommendations. [User Approval Gate with AI Independent Reviewer Suggestions] 
10. **/testtests Test the functional tests with placeholder feature code (all tests must fail in this phase) and guarantee test coverage of all core logic ((should be near 100% coverage))**

11. Loop until 3 is sucessfull [User Approval Gate with AI Independent Reviewer Suggestions]
    1. **/featdep Define Allowed feature code dependecies** [User Approval Gate with AI Review Suggestions]
    2. [Make/Remake plan.md first & keep updating the plan as you progress]  [User Approval Gate with Indepedent AI Plan Reviewer Suggestions]**/feat Write feature code using only the allowed feature code dependencies & Review against Spec and Deisgn System if doing GUI work (issue-specific spec and global spec) catching inconsistencies with spec/design system, things not specified in specs/design system and problems in spec/design system that needed to be overruled**. Incosistencies between code/tests;issue-specific-spec/global-spec should be flagged to the user with recommendations. Important: should try to make multiple related tests pass at a time, always aim for the smallest coherent behavioral slice, as end-to-end as possible across the components) that produces a useful feedback signal  [User Approval Gate with AI Independent Reviewer Suggestions]
    3. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate] **/test Build and Test feature code with the functional tests, generate testing & test coverage reports** & Review against Issue-specific & Global Spec, repeat this step until all tests pass

11. [Make/Remake plan.md first & keep updating the plan as you progress] **/refact-coreimplementation-task-if-necessary evaluate refactoring opportunities just for the last core implementaion task just realizsed. These should be changes that maintian same functionality but improve code clarity, quality and maintanability, [Make plan first & keep updating the plan as you progress]  [User Approval Gate] then implement the chossen refactoring bits one by one, after each one is done, evaluate if it actually is better than before (if not, just keep how it was before), only then move to the next** [AI Approval Gate]

 ----<<separate terminal block (reset context) >>----

12. **/grillme Understand the codebase, then grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]

----<</separate terminal block >>----

13. [Make/Remake plan.md first & keep updating the plan as you progress] **/stripdebuglogs** Remove debug logs from the code, only leave essential logs [User Approval Gate with AI Reviewer Suggestions]

14. **/determine-if-core-implementation-task-is-done** should evaluate if the core implementaiotn tasks objectives were in fact achieved and if the core-implementaiotn-tasks.md still holds. If all the tests actually reflect the spec and all the tests pass, mark as done the the current core implementaiton tasks inside core-implementaiotn-tasks.md. Shouldnt rely on historic test passing outputs or reports, should run the entire test suite for the current core implementaiotn task and for previously completed core implementation tasks (to catch regressions if any).

---- New Session (reset context) ----

**Requirement to continue the flow: core-implementation-tasks-plan.md needs to be fully complete, i.e., all core implementation tasks finished**

15. **/end-to-end-testing:** Should test end-to-end to make shure the issue was completely handled. If code involved GUI, should have GUI tester to test its usability and how good it looks, emulating real user behaviours. Also should check if end-to-end tests actually reflect the issue's prd requirements. Also verify is the issue spec is still consistent with global spec.

---- New Session (reset context) ----

C: Code Review

16. Loop until 1 is sucessfull or go back to a previous step [User Approval Gate with AI Independent Reviewer Suggestions]
    1. [Make/Remake plan.md first & keep updating the plan as you progress] **/review Independent Code Review** (including Code Review, including review against compliance to Spec (global and issue-specific) and Design System (if doing GUI work) catching inconsistencies with spec/design system, things not specified in spec/design system (if doing GUI work) and problems in spec/design system (if doing GUI work)). This step will generate a .agent/session/code_reviews/final_code_review file
    2.  [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate]**/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec. Note: code review is stored in .agent/session/code_reviews/final_code_review [User Approval Gate with AI Reviewer Suggestions]

---- New Session (reset context) ----

D: Security Review, Documentation & PR

17. Loop until 1 is sucessfull or go back to a previous step [User Approval Gate with AI Reviewer Suggestions]
    1. **/secreview Specialized Security Review** flagging critical problems & warnings
    2. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate]**/redo Make necessary ccode/test/docs changes, build & test** & Review against Issue-specific & Global Spec [User Approval Gate with AI Reviewer Suggestions]

18. **/end-to-end-testing:** Should test end-to-end to make shure the issue was completely handled. If code involved GUI, should have GUI tester to test its usability and how good it looks, emulating real user behaviours. Also should check if end-to-end tests actually reflect the issue's prd requirements. Also verify is the issue spec is still consistent with global spec.

19. **/document Document**: Final User Documentation if its already usable (how to install & use the product) & Controbutor Documentation (how to understand the codebase), both in the form of step by step tutorial. [User Approval Gate with Independent AI Reviewer Approval Suggestionse]

 ----<<separate terminal block (reset context) >>----

20. **/grillme Understand the codebase, then grill User with questions to see if he really understands the changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]

----<</separate terminal block >>----

21. **/pr (1) Check the pre-pr-checklist.md to see if everthing was done/updated. (2) Push & Open PR with Summary of Changes, PR has to link the Issue it solves. Never merge automatically.**

---- New Session (reset context) ----

E: Make fixes based on PR Reviews and/or CI failures until PR is merged

22. Loop until 1 is sucessfull or go back to a previous step 
    1. (On PR Review or CI failure Notification manually checked by user) **/prreviews Read PR Reviews & CI Run from Github and write them locally on a dedicated folder**
    2. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate] **/redo Make necessary code/test/docs changes, build & test** & Review against Issue-specific & Global Spec -- code, tests, specs should all be consistent with each other, if not flagg insconsistencies for the user to resolve [User Approval Gate with AI Indepentes Reviewer Suggestions]

    ---- New Session (reset context) ----

    3. Loop until 1 is succesfull
        1. **/secreview Specialized Independent Security Review** flagging critical problems & warnings 
        2. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate] **/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec -- code, tests, specs should all be consistent with each other, if not flagg insconsistencies for the user to resolve [User Approval Gate with AI Indepentes Reviewer Suggestions]

    4. **/end-to-end-testing:** Should test end-to-end to make shure the issue was completely handled. If code involved GUI, should have GUI tester to test its usability and how good it looks, emulating real user behaviours. Also should check if end-to-end tests actually reflect the issue's prd requirements. Also verify is the issue spec is still consistent with global spec.

     ----<<separate terminal block (reset context) >>----
    
    5. **/grillme Grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]

    ----<</separate terminal block >>----

    6. **/document Document**: Final User Documentation if its already usable (how to install & use the product) & Controbutor Documentation (how to understand the codebase), both in the form of step by step tutorial. [User Approval Gate with Independent AI Reviewer Approval Suggestionse]
    
    7. **/pr (1) Check the .agent/persistent/pre-pr-checklist.md to see if everthing was done/updated. (2) Push to already open PR (need to check if PR is already opened) with Summary of Changes, PR has to link the Issue it solves. Never merge automatically.****

> This is the end of the .agent/persistent/pre-pr-checklist.md file

## AI-assisted Coding Tech Stack

External Vendor Requirements: Opencode Go Subcription, Claude Credits, Github Repo

- Agent-native Editor/Terminal: Superset (only when tackling multiple issues in parallel)
- Terminal Agent Harnesses: OpenCode & open-code-review
- IDE (for better introspection + manual editing): VSCode
- LLM Provider Subscription: **OpenCode Go
- Coder Models: (0) Planning: Kimi K3 with medium resoning; (1) Spec, scaffold and Tests: Kimi K3 with high reasoning; (2) Core coding: DeepSeek V4 Flash with medium reasoning; (3) Independent Spec Review: Claude Opus 4.8 with high reasoning; (4) Plan Review: Kimi K3 with high reasoning; (5) Security Review: Claude Opus 4.8 with high reasoning; (6) end-to-end testing: Kimi K3 with high reasoning; (7) sub-agents model: Claude Opus 4.8; (8) Code Review: Opus4.6 (Model A) & Qwen3.7-Max (Model B) as cadidate generators and open-code-review using Qwen3.8-Max (Model C).
- Grillme Model: DeepSeek V4 Flash
- Local Routing: opencode-model-router (opencode plugin)
    - Fast Model: qwen2.5-coder:7b
    - Medium Model: current coder model
    - Heavy Model: current coder model
- Issue Flow: Start with just making each step a slash command, and leave to the developer to follow the steps (note: this doesnt enforce step execution, se requires developer commitment). Note: each slash commands should remember to update the issue resolving progress at the end (issue_flow.md file) or, if its the first step, create the file if its still not created). Each slash command should have a simple name.
- OpenCode Plugins: graphify, rtk, opencode-quota, opencode-model-router
- OpenCode MCPs: Github MCP
- Custom Freelunch OpenCode sub-agents: security-specialist, code-review-and-refactoring-specialist, testing-specialist, debugging-specialist. Tip: Agency-agents repo provides some sub-agents out-of-the-box.
- Custom Freelunch OpenCode Slash Commands: one for each unique step of the issue building flow
- Dependency Docs: a dependency_docs.md under ./docs with entries in the form "- <dependency>: <docs_link>" for all direct dependencies (not dependencies of dependencies). Pinned versions used in the project cna be seen in the lock file of the virtual dev environment tool.
- Skills: 
    - Custom Skills: created on-demand, under a `custom_skills` folder, via manual creation or via `skill-creator` to avoid having to repeat the same solution process over and over. Make these custom skills: 
        - grill-my-understanding (continually ask questions of the latest changes to codebase to me, to see if i understand the codebase. Always give score my answers and give feedback to it. Only stop when you feel i understand the codebase. The user can also specifify specific files for you to grill him about instead of the entire codebase)
        - understand-external-codebase (1. Build a doc eplxianng in detail the characteristics and internals of an external github codebase; 2. Add to this doc an explanation of where and why this codebase can be helpfull as a reference for ideias/patterns for the current project being built)
        - update-fixed-context (1. Infers new usefull knowledge from ./.agent/persistant/knowledge/mistakes.jsonl and ./.agent/persistant/completed_issue_flows; 2. Add this new usefull knowledge to AGENTS.md if its not already there)
        - make-core-implementation-tasks-plan (transforms the spec into a graph of tasks, where: (1) each task can depend on other tasks being already done or not depend on any; (2) the tasks should not be of the form "one task implements each component that will be needed in this feature, e.g., oen task for the backend, another for databas,e another for gateway and another for frontend", the tasks should be done in the form of "one task implments a slice (governed by on or more integration/end-to-end tests) of multiple components, e.g., this task implements a funcitonality slice of frontend, backend, gateway and frontend that together brings us one step closer to our end goal and guarantee rich cross-component feedback along development". The core implementation tasks plan needs to be stored in .agent/persistent/current-issue/flow/core-implementation-tasks-plan.md. The core-implementation-tasks-plan.md. file should have the graph structure of the plan, where each node is a task. For each node there is also a pending/in-progress/done checkbox. Do not confuse with plan.md which is a per-step ephmeral small plan for step execution.)
        - document (should first staleness and incompleteness of existing documentation (if any) and then update/create: (1) Contributor Documentation: visualization o repositoty strcture explaining succintly each directory and file, (1.2) Step by step contributor tutorial to help a newcomer understand the codebase; (2) User Documentation (only do after first version 0.1.0 is released): (2.1) User API Reference. (2.2) User step by step tutorial starting from sratch; (2.3) User guides to do common stuff; (2.4) FAQ. 
        make shure the documentaiton explains well things that I usually have a hard-time understanding.
        - ui-taste (UI Taste gives Claude a visual sense of taste. Instead of relying only on abstract design principles, the skill provides curated examples of bad, good, and stellar GUIs across different application categories and problem modes, including screenshots and their underlying HTML/CSS. This gives the agent an understanding of what makes GUIs look good. The agent should launch the current GUI, identify the biggest visual shortcomings, and iteratively improve them. The goal isn't to force a particular design style—it is to help Claude distinguish "functional but mediocre" from "genuinely beatifull and easy to use", giving coding agents a practical visual benchmark for judging their own work.)
    - Use existing skills: skill-creator, i-have-adhd, chrome-devtools-cli, grillme (every grillmre run should log all the questions, answers and feedback gave to the user inot a .agent/persistent/user-grills/grill[i].md where i is the id of the grill and the file should have timestamp, commit, what the grill was about and grill score in the beggining of it. Ever grill should start by looking at the commit, what was grilled in the last grill and user-codebase-questions.jsonl file), lavish-axi, code-review-and-quality, api-and-interface-design, browser-testing-with-devtools (only when working with frontend part), security-and-hardening, cc-skills-golang, maintainable-typescript (only when working with frontend part), improve-codebase-architecture, screenshot (only when working with frontend part), extract-design-system (only when working with frontend part), frontend-design (only when working with frontend part).

## How Code Review is Done

1. Start a new session with opencode: do code review with Model A and store the review in .agent/session/code_reviews/code_review_[A].md, where A is a placeholder for the actual mode name 
2. Start a new session with opencode: do code review with Model B and store the review in .agent/session/code_reviews/code_review_[B].md, where B is a placeholder for the actual mode name
3. Start a new session with open-code-review: do code review with Model C explicitely telling it to look at the candidate problems flagged inside .agent/session/coe-reviews/ folder and store the resulting code review inside .agent/session/code_reviews/final_code_review

## Token Efficency Laws

- Avoid small actions → batch related small tasks into larger coherent work units (use a todo_buffer.md to store all todos and then batch them before prompting the harness)
- Don't derail the agent from its main goal → context spent on unrelated work is expensive and increases context pollution.
- Use a graph/codebase-understanding tool → avoid repeatedly spending LLM tokens rediscovering repository structure and relationships.
- Related big tasks in the same session, unrelated new stuff gets its new session
- Start a new session when context becomes sufficiently polluted → carrying a huge amount of irrelevant history can become more expensive than rebuilding a clean context.
- Avoid rambling/random studying with the agent, do all of this in ChatGPT/Gemini/Grok webages
- Stay in the same session while the context is still useful → preserving cached/reusable context avoids paying to rebuild understanding.
- Until you reach scaffold code with tests, dont swittch models
- Always mention files with @ for the agent to look/modify instead o letting the agent winder the repo for that file
- Use something like rtk to compact tool outputs
- Have a router + local model for doing simple stuff

## .Agent Directory Strcture Guide Doc (./agent/directory_structure.md)

## `.agent/` Directory Structure

The `.agent/` directory contains the AI agent's workflow state, persistent knowledge, session state, and issue-specific process artifacts. It is an internal directory used by the coding agent and should not contain product source code.

### Directory Structure (./.agent/directory_structure.md file)

```text
.agent/
├── persistent/
│   ├── knowledge/
│       └── non_obvious_conjectures_and_facts.md
│       └── mistakes.jsonl
|   |── user-understanding/
|       └── user-codebase-questions.jsonl
│       └── user-grills/
│           └── grill[i].md
|   |── current-issue/
|        └── raw_github_issue.md
         └── flow
|            └── issue_flow_[i].md where i is the issue number
|            └── core-implementation-tasks-plan.md
|        └── issue-spec_[i]/ where i is the issue number
|               └── prd.md
|               └── architecture.md
|               └── tech_stack.md
│   └── completed-issues/
|       └── completed_issue_flows/
│       |   └── issue_flow_[i].md where i is the issue number
|       └──completed-issue-specs/
|           └── issue-spec_[i]/ where i is the issue number
|               └── prd.md
|               └── architecture.md
|               └── tech_stack.md
│
│── directory_structure.md
│
├── session/
│   └── plan.md
|   └── todos.md
│
├── session-persistence-candidate/
│   ├── assumptions.md
│   └── knowledge/
│       └── exploration_findings/
│           └── <name-of-exploration>.md
│
```

### `persistent/`

Contains durable project knowledge and historical state that should survive across sessions and remain useful to future agents.

* `knowledge/non_obvious_conjectures_and_facts.md` — durable, non-obvious facts, conclusions, and insights about the codebase that are useful to future agents.

* `mistakes.jsonl` — records mistakes made by the agent that required user intervention. Each entry records what was done, what was wrong, why it was wrong, and how it was corrected.

* `user-understanding/user-codebase-questions.jsonl` — records codebase-related questions asked by the user and the corresponding answers. This helps identify areas where the user may need additional explanation during future `/grillme` sessions.

* `user-understanding/user-grills/` — contains the history of `/grillme` sessions. Each `grill[i].md` records the commit being reviewed, what was tested, the questions asked, the user's answers and feedback, and the final grill score.

* `current-issue/` — contains the durable state of the issue currently being implemented.

  * `raw_github_issue.md` — the original GitHub issue, preserved as the source of truth for the issue being worked on.
  * `flow/` — contains the durable implementation workflow state for the current issue.

    * `issue_flow_[i].md` — records the sequential progress through the Issue Flow for issue `i`, including completed steps, timestamps, approvals, session summaries, and deviations or jumps between steps.
    * `core-implementation-tasks-plan.md` — contains the dependency graph of the core implementation tasks for the current issue. Each task represents a coherent behavioral or functional slice and records its dependencies and progress.
  * `issue-spec_[i]/` — contains the specification produced for issue `i`.

    * `prd.md` — product requirements and expected behavior.
    * `architecture.md` — architectural design and implementation boundaries.
    * `tech_stack.md` — technologies, dependencies, and relevant technical choices.

* `completed-issues/` — archives durable artifacts from issues that have been completed.

  * `completed_issue_flows/` — contains archived `issue_flow_[i].md` files documenting how completed issues were implemented.
  * `completed-issue-specs/` — contains archived `issue-spec_[i]/` directories, including the PRD, architecture, and technology-stack documents for completed issues.

Persistent knowledge should contain only information that is expected to remain useful beyond the current session. Current-issue state belongs under `current-issue/` while historical issue state belongs under `completed-issues/`.

### `session/`

Contains temporary state for the current agent session.

* `plan.md` — the ephemeral step-by-step plan for the current slash-command execution.
* `todos.md` — the temporary task list used to track work during the current session.

Everything in `session/` is disposable. At the beginning of a new session, the contents of `.agent/session/` are cleared. These files must not be treated as historical records of the project or issue. Durable workflow state belongs in `persistent/current-issue/`.

### `session-persistence-candidate/`

Contains information discovered during the current session that may be useful beyond the session but has not yet been promoted to persistent knowledge.

* `assumptions.md` — assumptions made during the current session, including their evidence, whether they have been verified, and their risk if wrong.
* `knowledge/exploration_findings/` — findings produced by exploratory sub-agents. Each exploration should produce a findings file containing the exploration context, description, model/agent used, token/cost information, reason it ended, findings, and conclusion.

At the beginning of a new session, useful information from this directory is reviewed and promoted into `.agent/persistent/knowledge/` when appropriate. After promotion, the candidate directory is cleared.

### `directory_structure.md`

Documents the purpose and organization of the `.agent/` directory itself. It should be updated whenever the directory structure or the responsibilities of its files change.

### Important distinction: session state vs. persistent state

The `.agent/` directory deliberately separates **ephemeral working state** from **durable project state**:

* `session/` contains information needed only to continue the current session.
* `session-persistence-candidate/` contains potentially reusable discoveries that have not yet been validated or promoted.
* `persistent/` contains validated, durable knowledge and issue history.

Within `persistent/`, `current-issue/` is the active issue's durable workspace, while `completed-issues/` preserves the historical record of issues that have already been completed.


### General Rules

1. **Do not put source code in ****`.agent/`****.** Product code belongs under `./src/`.
2. **Do not treat ****`.agent/session/`**** as persistent storage.** Its contents may be deleted when a session starts.
3. **Do not promote assumptions or inferences automatically.** An assumption or infrence should only become persistent knowledge after it has been sufficiently verified.
4. **Do not silently modify historical issue records.** They are part of the project's implementation history.

> This is the end of the .agent/directory_structure.md file

## Custom PR Skill

Mandatory Pre-PR checklist (ready to push and open PR):

- global spec, issue-specific spec and tests all consistent with each other
- all core (at least 90% total code coverage) logic covered by tests, with test coverage report evidence for it
- unit, integration and end-to-end tests present and passing
- no big changes were made after the last independent code review was run
- no significant changes were made after the last security review
- documentation is up to date with the code
- user understand the PR that will be made at the function interface/class interface/file/directory level

Making PRs:

- ensure the pre-pr-checklist.md is checked before making a PR
- follow the project's PR template
- in the pr you write: highlight key decisions, problems encoutered, solutions and tradeoffs chosen
- in the pr you write: highlight what you tested and provide link to evidence that shows your test (log file, screenshot, etc)
- in the pr you write: make a risk assesment of the PR (Low, Medium, High) based on how many changes it makes, the type of changes it makes, test coverage, etc

Handling PR Comments and Reviews:

- you should respond PR comments adresing the issues raised or just respond PR questions made by other developers
- you should make apprpriate changes according to the feedback received


