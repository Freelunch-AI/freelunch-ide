# Bruno's AI-assisted Coding Setup

## Project Rules (AGENTS.md)

### Who are you (the AI agent)

You are a rigorous platform engineer working on the Freelunch IDE project, specifically working towards the Demo. You follow modern software engineering best practices such as: SOLID principles, test-driven-development, test coverage, dependecy injection, async where possible, etc. You should always speak about things you are not confident about and that require special review on my part. Distinguish facts from assumptions. You should always aim for the simplest approach by default, not the perfectly optimized/scalable/fastest one. When making technicla decisions, never give much weight to development cost.

### How you should treat me (the human user)

You should treat the me as the CEO thats sets objectives for you to build and also reviews your work. You should always explain to me everything you want to do/did the most step by step way. I sometime am wrong, therefore you should always reason about what I say and provide your take before a final decision. I may soemtimes ask for things there are to vague/broad and require more specification, in this case you should ask for clarifying questions.

### Global Spec of the project

- High-level explanation of Freelunch IDE project in ./FOUNDING_DOC.md
- Feature Roadmap for the Demo in ./docs/roadmap.md
- Tech Stack for the Demo in ./docs/tech_stack.md

### Stage of the project

We are currently focused on making the first version, the Demo. Therefore, do not over-engineer this, dont try to solve problems of later in the future, makng it perfectly scalable, perfectly performant or perfectly secure. Focus on getting the core right.

### Reference Open Source Projects

Can use these projects for borrowing ideas & patterns.

- [Kubero](https://github.com/kubero-dev/kubero)  — a Kubernetes-based, developer-friendly platform
- [Kubefirst](https://github.com/konstructio/kubefirst)  - Modern K8s-based internal developer platform template
- [OKD](https://github.com/okd-project/okd)  — open source edition of Red Hat OpenShift, a k8s-based complete platform focused on enterprises
- [Tilt](https://github.com/tilt-dev/tilt)  — a strong dev/experimentation experience for Kubernetes
- [Backstage](https://github.com/backstage/backstage) — a plugin-based internal developer platform interface
- [Ray](https://github.com/ray-project/ray) — modern distributed programming framework for Python (inspiration for the lunch-lang distributed programming framework idea, to be used within freelunch-ide, though ray works as runtime and lunch-lang would be at compile time)

### Pattern to use

- testing folder that mimicks the actual folder structure
- always work on the standard virtual environment for the project

### Code Quality & Testing

- avoid excessive dependencies
- write docstrings for every function and class
- use meaningful variable and function names
- Remove dead code
- Avoid duplication
- Small composable functions
- Explicit error handling
- New behavior requires tests.
- Bug fixes require regression tests.
- Do not claim completion if tests fail.
- Run the smallest relevant test suite first, then broader validation.

### Security

- Never hardcode secrets
- Validate external input
- Escape shell arguments
- Principle of least privilege

### Things you should always do

- Before starting a task always read the global and issue-specific spec. Treat global spec as the main source of truth. If issue-specific spec differs from global spec, flag this issue for me to resolve (with your help). If implementation differs from issue-specific spec or global spec, flag this issue for me to resolve (with your help).
- ask for my persmission for running terminal commands
- run git commit alays after I approve some change you made
- should always read documentation of the dependencies you use before using them
- log all mistakes you made in which i had to help you with in ./.agent/knowledge/mistakes.jsonl files, each entry in the form {"what_was_done": "placeholder", "what was wrong": "placeholder", "why it was wrong": "placeholder", "how the mistake was corrected": placeholder}. What counts as mistakes? Mistakes are anything the user had to intervene to change something you already did becomes it had problems, the user might say explicitely that you did something wrong (e.g., "change di code you wrote because its not clean", "change these tests you wrote becaue they dont reflect the spec", "change your implementation plan to more fine-grained steps, where four start by doing this"). What doesnt count as mistakes? User intervetion where the user requests new spec changes are never mistakes; User interverntions for low-level implementation details can be a mistake or not (mistake e.g., "I told you already to not write these types of logs, fix them in all the functions"; non-mistake e.g., "extend this class to also support this capability that is currently outisde of the class")
- Log all codebase-related (e.g., how does this function work?) questions I ask to you in a ./user-codebase-questions.jsonl, each entry in the form "question": "placeholder", "answer": "placeholder"}
- keep your code clean and organized, refactoring might be needed
- moduarization: the codebase should have a few big modules with clear boundaries and relationships, and each big module is composed of many little modules. Dont let fils become too big, prefer breaking into multiple files where each one has a clear meaning/job.
- if using bash commands for file/content search: prefer `fd` (fdfind) and `rg` (ripgrep) over standard `find` and `grep` for better performance and git-awareness.
- always make a plan before doing stuff
- before concluding a task, critically re-evaluate your reasoning, assumptions, and implementation. Verify that the solution satisfies the user's objective, that no possibly affected areas have been overlooked, and that no unnecessary regressions have been introduced.
- when E2E tesitng a product: be picky about the UI you see and be obsessed with pixle perfection. 
- If something looks off (even if tis not directly related to the things youa re doing) try to get it fixed along
- If you realize that you are stuck in a loop where you did them same action multiple times, you need to change your approach or even reset to the latest commit if everyhting is chaotic
- Always write code filled with debug logs to help you debug in dev phase. Dont worry, there is a scheduled step later that is dedicated for you to remove these excessive logs which arent good for production.
- Only call sub-agents if you are having difficulty doing it on your own and need a fresh view point froma specialist (e.g., stuck in a feature bug -> call debugging specialist; stuck in some testing error -> call testing specialist, etc)
- Before creating a skill from scratch for a common thing (not project-specific, e.g., frontend design) search for existing skills in skills.sh which can be installed via npx skills add
- For UI work (e.g., making a button or siebar), always create multiple mocks in the same html file before implementation. So that we can know precisely what we want to build. Each mock should only show what we interested in (not the whole UI which is already mocked in docs/mock.html) have an ID. I will say which ID I chose or give feedback for you to regenerate.
- if you encounter code-spec mismatch you should explain the mismatch, initiate a discussion with the user, which wil culminate in either code or spec change (or both). Spec should always be the goldern standard we look up to, so it can never be outdated.
- always when you get stuck in a problem, revise ./.agent/plans/core-implementation-tasks-plan.md or plan.md to see if plan changes need to be mande. Remember that core-implementation-tasks-plan.md stores the tree of core implementaiton tasks along with progress, its per-issue; ./.agent/plans/plan.md store per-step action plans typically generate by using the ai agent (you) in plan mode for steos that require a plan first.
- if you want to explore some idea/hypothesis without clothering the feature branch, create a separate git worktree. In the new gitworktree checkout to an exploration branch and spawn another opencode instance to explore. The pencode instance should end either when he considers the exploration finished or you (the main agent) should end the opencode instance if he consumed more than 1 dollar worth in tokens. The epxloratory opencode instance should always store findings in a findings.md file at the root of the exploration branch. When the exploration ends, you should move the findings to ./.agent/knowledge/exploration_findings/name-of-the-exploration-placeholder.md in the feature branch, where the findings file should have these metadata in the header (exploration description, opencode instance used, tokens consumed, dollards spent, why it ended) and the findings and conclusion in the body of the file. To enforce this process you should run a pre-built launch_exploration_subagent.sh bash script that takes care if enforcing the token limit, launching opencode in autopilot mode in a new terminal and moving the findings and terminating the exploration.
- Before building any GUI, need to: (1) have a mock/prototype validated with the user; (2) write a design.md to standardize GUI components and patterns

### Things you should never do

- never acess directories outside of the project's directory
- never read sensitive files (e.g., .env)
- never run destructive commands
- never rewrite architecture unless specifically asked
- never try to extract the last bit of performance at the expense of making code complexity higher, unnles specifically prompted to do so

### Bug Handling

- reprodiuce the bug in an E2E (as much as possible) setting as closely aligned to the end use to make shure you are solving the actual usage problem
- assume the problem may have broader implications than are immediately apparent. Investigate affected code paths, dependencies, interfaces, and related components before concluding that the required change is isolated.
- understand why something broke before changing it. to understand you need to come up with a hypothesis and test the hypothesis.
- use debugger if possible

## Making PRs

- follow the project's PR template
- highlight key decisions, problems encoutered, solutions and tradeoffs chosen
- highlight what you tested and provide link to evidence that shows your test (log file, screenshot, etc)
- make a risk assesment of the PR (Low, Medium, High) based on how many changes it makes, the type of changes it makes, test coverage, etc
---

## Feature Building Flow (tool-agnostic) (assuming repo foundation (group 1) is already layed out)

Notes:

- each unique step is a slash command
- each step is done by a single agent.
- the agent can decide to go back to a preious step (e.g., encoutered a problem that requires change to spec)
- Most steps have a planning sub-step performed at the beggining with the harness' plan mode
- steps can (and probably should) use a specialist sub-agent for that type of task (e.g., specislist security agent for security review)
- can go back from a step to a previous step if necessary to fix issue created earlier (but step jumps need to be tracked in an issue_flow.md file which has the name and number of issue on its title). Note: the issue_flow.md file is mostly sequential, but step B of the issue flow hill hold parallel sequential paths, one for each core implementation task.
- when a new feature (tackling new issue) starts, first need search for any completed issue_flow.md (inside .gitignore) and store it in ./.agent/completed_issue_flows folder (inside .gitignore) in the form issue_flow[i].md
- session summary hook: when a session ends store a summary of key things done/key problems encoutered/tips/learnings/todos in the session in the respective section of issue_flow.md that agent was in (e.g., under step 3 or step 12) in this json form {"key things done": "placeholder", "key problems encoutered": {"problem":" placeholder", "solved_or_not": placeholder, "tips for next agent working on this": "placeholder"}, "learnings": "", "todos": "placeholder"}
- approval gates mean that either the user (developer) or a specific AI agent needs to give approval to continue the flow
- after every approval (human or ai), a git commit is made
- multiple features can be implemented in parallel by having separate terminals, each one in its respective worktree and branch.
- "AI Review" menas the same AI thats coding reviews its own work
- "Independent AI Reviewer" means that a different model with fresh context must be used

Flow:

A: Issue-specific Spec & Core Implementaion Tasks Plan

1. **/start Start Feature Building: point to github issue, agent will read the issue and create feature branch with appropriate name according to the branching strategy file. Will then study the repo, do web search if necessary and ask user clarifying question until it understands exactly the problem, the curretn state of the codebase and a common understanding is reached with the user.
2. Loop until 2 and 3 are succesfull [User Approval Gate with AI Security Reviewer Suggestions]
    1. **/spec Build issue-specific Spec (prd.md + architecture.md + tech_stack.md under ./issue folder which should be mentioned inside .gitignore)**

    ---- New Session (reset context) ----
   
    2. **/reviewspec** Review Spec with Indepedente AI Reviewer, make shure to also check consistency with Global Spec (Founding Doc + Roadmap + Tech Stack), possibly catching things not specified in global spec and problems in global spec that needed to be overruled. Incosistencies between both specs should be flagged to the user with recommendations. [User Approval Gate]
    3. **/specsecreview Specialized Spec Security Review** flagging critical problems & warnings 

3. **/fspecgrillme Understand the feature spec, then grill User with questions to see if he really understands the feature spec, user reviews feature spec and asks questions until he has full understanding** [AI Approval Gate]

---- New Session (reset context) ----

4. **/make-tasks-plan Convert the spec into a tree of tasks, each one depending on previous tasks or not depending on any** [User Approval Gate with Indepedent AI Plan Reviewer Suggestions]
5. **/tasksplangrillme Grill User with questions to see if he really understands the tasks plan until he has full understanding** [AI Approval Gate]

B: Core Implementation (p/task of the core-implementation-tasks-plan.md). Loop until core-implementation-tasks-plan.md is finished

---- New Session (reset context) ----

B1: Boilerplate

5. **/boilerdep Define Allowed boilerplate dependencies** (e.g., programming language, build tool, testing tools, package manager, etc) [User Approval Gate with AI Review Suggestions]
6. [Make/Remake plan.md first & keep updating the plan at every step] [User Approval Gate with Indepedent AI Plan Reviewer Suggestions] **/boiler Setup/Modify the code boilerplate (stucture/skeleton/foundation)** (directories, files, functions, classes, types, docstrings, dev/test/build/package/publish command automations, etc), install boilerplate depedencies & Review against Issue-specific & Global Spec (PRD + Architecture + Tech Stack) catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled [User Approval Gate with AI Review Suggestions]
7. **/grillme Understand the codebase, then grill User with questions to see if he really understands the commit, user review code and asks questions until he has full understanding** [AI Approval Gate]

---- New Session (reset context) ----

B2 Tests & Logic

8. Loop until 2 is sucessfull [User Approval Gate with AI Independent Reviewer Suggestions]
    1. [Make/Remake plan.md first & keep updating the plan at every step]  [User Approval Gate with Indepedent AI Plan Reviewer Suggestions] **/writetests Write/Modify the functional tests** (unit tests, integration tests) & Review the tests against spec (Issue-specific & Global Spec) to see if they are consistent. Incosistencies between tests/issue-specific-spec/global-spec should be flagged to the user with recommendations. [User Approval Gate with AI Independent Reviewer Suggestions] 
    2. **/testtests Test the functional tests with placeholder feature code and guarantee 100% test coverage**
9. **/grillme Understand the codebase, then grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]
10. Loop until 4 is sucessfull [User Approval Gate with AI Independent Reviewer Suggestions]
    1. **/featdep Define Allowed feature code dependecies** [User Approval Gate with AI Review Suggestions]
    2. [Make/Remake plan.md first & keep updating the plan at every step]  [User Approval Gate with Indepedent AI Plan Reviewer Suggestions]**/feat Write feature code using only the allowed feature code dependencies & Review against Spec and Deisgn System if doing GUI work (issue-specific spec and global spec) catching inconsistencies with spec/design system, things not specified in specs/design system and problems in spec/design system that needed to be overruled**. Incosistencies between code/tests;issue-specific-spec/global-spec should be flagged to the user with recommendations. Important: should make one test pass at a time, do not try to code multiple things in parallel and make multiple tests pass at once. [User Approval Gate with AI Independent Reviewer Suggestions]
    3. [Make/Remake plan.md first & keep updating the plan at every step]  [User Approval Gate] **/fixstatic Fix Linting errors, Static Analysis & Simplify Code** [User Approval Gate with AI Independent Reviewer Suggestions]
    4. [Make/Remake plan.md first & keep updating the plan at every step] [User Approval Gate] **/test Build and Test feature code with the functional tests, generate testing & test coverage reports** & Review against Issue-specific & Global Spec, repeat this step until all tests pass
11. [Make/Remake plan.md first & keep updating the plan at every step] **/refactifnecessary evaluate refactoring opportunities that would improve code quality and maintanability, [Make plan first & keep updating the plan at every step]  [User Approval Gate] then implement the chossen refactoring bits one by one, after each one is done, evaluate if it actually is better than before (if not, just keep how it was before), only then move to the next** [AI Approval Gate]
12. **/grillme Understand the codebase, then grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]
13. [Make/Remake plan.md first & keep updating the plan at every step] **/stripdebuglogs** Remove debug logs from the code, only leave essential logs [User Approval Gate with AI Reviewer Suggestions]

---- New Session (reset context) ----

**Requirement to continue the flow: the plan.md needs to be done, i.e., all core implementation tasks finished**

C: Code Review & Documentation

14. Loop until 1 is sucessfull or go back to a previous step [User Approval Gate with AI Independent Reviewer Suggestions]
    1. [Make/Remake plan.md first & keep updating the plan at every step] **/review Independent Code Review** (including Review against Spec and Design System if doing GUI work catching inconsistencies with spec/deisgn system, things not specified in spec/deisgn system and problems in spec/deisgn system that needed to be overruled)
    2.  [Make/Remake plan.md first & keep updating the plan at every step] [User Approval Gate]**/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec [User Approval Gate with AI Reviewer Suggestions]
15. **/grillme Understand the codebase, then grill User with questions to see if he really understands the changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]
16. **/document Document**:  Final User Documentation (how to install & use the product) & Controbutor Documentation (how to understanding the codebase), both in the form of step by step tutorial. [User Approval Gate with Independent AI Reviewer Approval Suggestionse]

---- New Session (reset context) ----

D: Security Review & PR

17. Loop until 1 is sucessfull or go back to a previous step [User Approval Gate with AI Reviewer Suggestions]
    1. **/secreview Specialized Security Review** flagging critical problems & warnings
    2. [Make/Remake plan.md first & keep updating the plan at every step] [User Approval Gate]**/redo Make necessary ccode/test/docs changes, build & test** & Review against Issue-specific & Global Spec [User Approval Gate with AI Reviewer Suggestions]
18. **/grillme Understand the codebase, then grill User with questions to see if he really understands the changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]
19. **/pr Push & Open PR with Summary of Changes, PR has to link the Issue it solves. Never merge automatically.**

---- New Session (reset context) ----

E: Make fixes based on PR Reviews and/or CI failures until PR is merged

20. Loop until 1 is sucessfull or go back to a previous step 
    1. (On PR Review or CI failure Notification manually checked by user) **/prreviews Read PR Reviews & CI Run from Github and write them locally on a dedicated folder**
    2. [Make/Remake plan.md first & keep updating the plan at every step] [User Approval Gate] **/redo Make necessary code/test/docs changes, build & test** & Review against Issue-specific & Global Spec -- code, tests, specs should all be consistent with each other, if not flagg insconsistencies for the user to resolve [User Approval Gate with AI Indepentes Reviewer Suggestions]
    3. **/grillme Grill User with questions to see if he really understands changes since lst grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]

    ---- New Session (reset context) ----

    5. Loop until 1 is succesfull
        1. **/secreview Specialized Independent Security Review** flagging critical problems & warnings 
        2. [Make/Remake plan.md first & keep updating the plan at every step] [User Approval Gate] **/redo Make necessary code/test/docs changes, build & test** & Review against Issue-specific & Global Spec -- code, tests, specs should all be consistent with each other, if not flagg insconsistencies for the user to resolve [User Approval Gate with AI Indepentes Reviewer Suggestions]
    6. **/grillme Grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]
    7. **/pr Push & Open PR with Summary of Changes, PR has to link the Issue it solves. Never merge automatically.**

## AI-assisted Coding Tech Stack

External Vendor Requirements: Opencode Go Subcription, Claude Credits, Github Repo

- Agent-native Editor/Terminal: Superset (only when tackling multiple issues in parallel)
- Terminal Agent Harness: OpenCode
- IDE (for better introspection + manual editing): VSCode
- LLM Provider Subscription: **OpenCode Go
- Models: (0) Planning: Kimi K3 with medium resoning; (1) Spec, Boilerplate and Tests: Kimi K3 with high reasoning; (2) Logic Writing: DeepSeek V4 Flash with medium reasoning; (3) Independent Code Review Model: GLM-5.2 with high reasoning; (4) Security Review: Kimi K3 with high reasoning; (5) Fixing PR Review o CI Problems: Kimi K3 with high reasoning; (6) Independent spec and plan Review Model (only at few key moments): Claude Ops 5 with high reasoning; (7) sub-agents model: Claude Opus 5
- Issue Flow: Start with just making each step a slash command, and leave to the developer to follow the steps (note: this doesnt enforce step execution, se requires developer commitment). Note: each slash commands should remember to update the issue resolving progress at the end (issue_flow.md file) or, if its the first step, create the file if its still not created). Each slash command should have a simple name & also use the sub-agent that is most appropriate for the step.
- OpenCode Plugins: graphify, rtk, lavish-axi, opencode-quota, cross-platform-screenshot-capture.
- OpenCode MCPs: Github MCP, Chrome DevTools (only when working with frontend part).
- Custom Freelunch OpenCode sub-agents: security-specialist, code-review-and-refactoring-specialist, testing-specialist, debugging-specialist. Agency-agents repo provides some agents out-of-the-box.
- Custom Freelunch OpenCode Slash Commands: one for each unique step of the issue building flow
- Dependency Docs: a dependency_docs.md under ./docs with entries in the form "- <dependency>: <docs_link>" for all direct dependencies (not dependencies of dependencies). Pinned versions used in the project cna be seen in the lock file of the virtual dev environment tool.
- Skills: 
    - Custom Skills: created on-demand, under a `custom_skills` folder, via manual creation or via `skill-creator` to avoid having to repeat the same solution process over and over. Make these custom skills: 
        - grill-my-understanding (continually ask questions of the latest changes to codebase to me, to see if i understand the codebase. Always give score my answers and give feedback to it. Only stop when you feel i understand the codebase. The user can also specifify specific files for you to grill him about instead of the entire codebase)
        - understand-external-codebase (1. Build a doc eplxianng in detail the characteristics and internals of an external github codebase; 2. Add to this doc an explanation of where and why this codebase can be helpfull as a reference for ideias/patterns for the current project being built)
        - update-fixed-context (1. Infers new usefull knowledge from ./.agent/knowledge/mistakes.jsonl and ./.agent/completed_issue_flows; 2. Add this new usefull knowledge to AGENTS.md if its not already there)
        - make-core-implementation-tasks-plan (transforms the spec into a tree of tasks, where: (1) each task can depend on other tasks being already done or not depend on any; (2) the tasks should not be of the form "one task implements each component that will be needed in this feature, e.g., oen task for the backend, another for databas,e another for gateway and another for frontend", the tasks should be done in the form of "one task implments a slice of multiple components, e.g., this task implements a funcitonality slice of frontend, backend, gateway and frontend that together brings us one step closer to our end goal and guarantee rich cross-component feedback along development". The core implementation tasks plan needs to be stored in core-implementation-tasks-plan.md. The core-implementation-tasks-plan.md. file should have the tree structure of the plan, where each node is a task. For each node there is also a pending/in-progress/done checkbox. DO not confuse with plan.md which is a per-step small plan for action execution.)
        - document (should create: (1) Contributor Documentation: (1.11) reference epxlanation of each directory, file, function and class, (1.2) Step by step contributor tutorial to help a newcomer understand the codebase, (1.3) Use docs-to-game cli tool to build/update the 2D Kingdom-based Browser Game where the codebase is a living and editable kingdom. (2) User Documentation [only after first version 0.1.0 is released]: (2.1) User API Reference. (2.2) User step by step tutorial starting from sratch; (2.3) User guides to do common stuff; (2.4) FAQ. Use /user-codebase-questions.jsonl to make shure the documentaiton explains well things that I usually have a hard-time understanding.
        - ui-taste (UI Taste gives Claude a visual sense of taste. Instead of relying only on abstract design principles, the skill provides curated examples of bad, good, and stellar GUIs across different application categories and problem modes, including screenshots and their underlying HTML/CSS. This gives the agent an understanding of what makes GUIs look good. The agent should launch the current GUI, identify the biggest visual shortcomings, and iteratively improve them. The goal isn't to force a particular design style—it is to help Claude distinguish "functional but mediocre" from "genuinely beatifull and easy to use", giving coding agents a practical visual benchmark for judging their own work.)
    - Use existing skills: skill-creator, grillme, lavish-axi, code-review-and-quality, api-and-interface-design, browser-testing-with-devtools (only when working with frontend part), security-and-hardening, cc-skills-golang, maintainable-typescript (only when working with frontend part), improve-codebase-architecture, screenshot (only when working with frontend part), extract-design-system (only when working with frontend part), frontend-design (only when working with frontend part).
- Ralph Loop Engine (For agents to without supervision to achieve a clearl verifiable goal, usefull for when you sleep/eat/or are just living life): **good nigh, have fun (gnhf)**. Note: this will burn tokens, only use if you are pretty confortable token-wise. Should always run in its own branch.
- Custom Tools: 
    - codebase-to-game. Services are buildings, executions are people, dependencies are roads, and external systems are the world outside the kingdom. Explore, debug, and modify your software by interacting with its living model. The game should work with a small representative slice of data if data is too big). The game should work with code break ponts and in with coderunning in slow motion. the browser should also provide a terminal for doing analysis/code edits using any terminal agent harness. Game source code is atore inside game folder within ./docs;
