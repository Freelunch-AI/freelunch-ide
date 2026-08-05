# Bruno's AI-assisted Coding Setup

## Project Rules (agents.md)

### Who are you (the AI agent)
You are a rigorous platform engineer working on the Freelunch IDE project, specifically working towards the Demo. You follow modern software engineering best practices such as: SOLID principles, test-driven-development, test coverage, dependecy injection, async where possible, etc. You should always speak about things you are not confident about and that require special review on my part. Distinguish facts from assumptions. You should always aim for the simplest approach by default, not the perfectly optimized/scalable/fastest one.

### How you should treat me (the human user)
You should treat the me as the CEO thats sets objectives for you to build and also reviews your work. You should always explain to me everything you want to do/did the most step by step way. I sometime am wrong, therefore you should always reason about what I say and provide your take before a final decision. I may soemtimes ask for things there are to vague/broad and require more specification, in this case you should ask for clarifying questions.

### Global Spec of the project
- High-level explanation of Freelunch IDE project in ./FOUNDING_DOC.md
- Feature Roadmap for the Demo in ./docs/roadmap.md
- Tech Stack for the Demo in ./docs/tech_stack.md

### Reference Open Source Projects
Can use these projects for borrowing ideas & patterns.
- Kubero — a Kubernetes-based, developer-friendly platform
- Kubefirst - Modern K8s-based internal developer platform template
- OKD — open source edition of Red Hat OpenShift, a k8s-based complete platform focused on enterprises
- Tilt — a strong dev/experimentation experience for Kubernetes
- Backstage — a plugin-based internal developer platform interface
- Ray — modern distributed programming framework for Python (inspiration for the lunch-lang distributed programming framework idea, to be used within freelunch-ide, though ray works as runtime and lunch-lang would be at compile time)

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
- log all mistakes you made in which i had to help you with in ./mistakes.jsonl files, each entry in the form {"what_was_done": "placeholder", "what was wrong": "placeholder", "why it was wrong": "placeholder", "how the mistake was corrected": placeholder}
- keep your code clean and organized, refactoring might be needed
- if using bash commands for file/content search: prefer `fd` (fdfind) and `rg` (ripgrep) over standard `find` and `grep` for better performance and git-awareness.
- always make a plan before doing stuff
- before concluding a task, critically re-evaluate your reasoning, assumptions, and implementation. Verify that the solution satisfies the user's objective, that no possibly affected areas have been overlooked, and that no unnecessary regressions have been introduced.

### Things you should never do
- never acess directories outside of the project's directory
- never read sensitive files (e.g., .env)
- never run destructive commands
- never rewrite architecture unless specifically asked

### Bug Handling
- assume the problem may have broader implications than are immediately apparent. Investigate affected code paths, dependencies, interfaces, and related components before concluding that the required change is isolated.
- understand why something broke before changing it. to understand you need to come up with a hypothesis and test the hypothesis.
- use debugger if possible

## Feature Building Flow (tool-agnostic) (assuming repo foundation (group 1) is already layed out)

Notes:
- each unique step is a slash command
- each step is done by a single agent.
- each step has a planning sub-step performed at the beggining
- steps can (and probably should) use a specialist sub-agent for that type of task (e.g., specislist security agent for security review)
- can go back from a step to a previous step if necessary to fix issue created earlier (but step jumps need to be tracked in an feature_flow.md file which has the name and number of issue on its title)
- when a new feature starts, first need search for any completed feature_flow.md (inside .gitignore) and store it in completed_feature_flow folder (inside .gitignore) in the form feature_flow[i].md
- approval gates mean that either the user (developer) or a specific AI agent needs to give approval to continue the flow
- after every approval (human or ai), a git commit is made
- multiple features can be implemented in parallel by having separate terminals/worktrees for each working feature branch.
- "AI Review" menas the same AI thats coding reviews its own work
- "Independent AI Reviewer" means that a different model with fresh context must be used

Flow:

1. **/start Start Feature Building: point to github issue, agent will read the issue and create feature branch with appropriate name according to the branching strategy file**
2. **/clarify Ask User Clarifying Questions & do web search if necessary** [User Approval Gate with AI Review Suggestions]
3. Loop until 2 is succesfull [User Approval Gate with AI Security Reviewer Suggestions]
    1. **/spec Build issue-specific Spec (PRD + Architecture + Tech Stack)** & Review against Global Spec (Founding Doc + Roadmap + Tech Stack) catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled [User Approval Gate with AI Review Suggestions]
    2. **/specsecreview Specialized Spec Security Review** flagging critical problems & warnings 
4. [Make plan first & keep updating the plan at every step] **/fspecgrillme Understand the feature spec, then grill User with questions to see if he really understands the feature spec, user reviews feature spec and asks questions until he has full understanding** [AI Griller Approval Gate]

---- New Session ----

5. **/boilerdep Define Allowed boilerplate dependencies** (e.g., programming language, build tool, testing tools, package manager, etc) [User Approval Gate with AI Review Suggestions]
6. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/boiler Setup/Modify the code boilerplate (stucture/skeleton)** (directories, files, functions, classes, types, docstrings, test build command, final packaging build command), install boilerplate depedencies & Review against Issue-specific & Global Spec (PRD + Architecture + Tech Stack) catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled [User Approval Gate with AI Review Suggestions]
7. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands the commit, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]

---- New Session ----

8. Loop until 2 is sucessfull [User Approval Gate with AI Independent Reviewer Suggestions]
    1. [Make plan first & keep updating the plan at every step]  [User Approval Gate] **/writetests Write/Modify the functional tests** (unit tests, integration tests) & Review against Issue-specific & Global Spec [User Approval Gate] with AI QA Review Suggestions] 
    2. **/testtests Test the functional tests with placeholder feature code**
9. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]
10. Loop until 4 is sucessfull [User Approval Gate with AI Independent Reviewer Suggestions]
    1. **/featdep Define Allowed feature code dependecies** [User Approval Gate with AI Review Suggestions]
    2. [Make plan first & keep updating the plan at every step]  [User Approval Gate]**/feat Write feature code using only the allowed feature code dependencies & Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled** [User Approval Gate with AI Independent Reviewer Suggestions]
    3. [Make plan first & keep updating the plan at every step]  [User Approval Gate] **/fixstatic Fix Linting errors, Static Analysis & Simplify Code** [User Approval Gate with AI Independent Reviewer Suggestions]
    4. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/test Build and Test feature code with the functional tests, generate testing & test coverage reports** & Review against Issue-specific & Global Spec, repeat this step until all tests pass 
11. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]

---- New Session ----

12. Loop until 1 is sucessfull or go back to a previous step [User Approval Gate with AI Independent Reviewer Suggestions]
    1. [Make plan first & keep updating the plan at every step] **/review Independent Code Review** (including Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled)
    2.  [Make plan first & keep updating the plan at every step] [User Approval Gate]**/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec [User Approval Gate with AI Reviewer Suggestions]
13. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands the changes since last grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]
14. **/document Document**:  Final User Documentation (how to install & use the product) & Controbutor Documentation (how to understanding the codebase), both in the form of step by step tutorial. [User Approval Gate with Independent AI Reviewer Approval Suggestionse]

---- New Session ----

15. Loop until 1 is sucessfull or go back to a previous step [User Approval Gate with AI Reviewer Suggestions]
    1. **/secreview Specialized Security Review** flagging critical problems & warnings
    2. [Make plan first & keep updating the plan at every step] [User Approval Gate]**/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec [User Approval Gate with AI Reviewer Suggestions]
16. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands the changes since last grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]
17. **/pr Push & Open PR with Summary of Changes, PR has to link the Issue it solves. Never merge automatically.**

---- New Session ----

18. Loop until 1 is sucessfull or go back to a previous step 
    1. [On PR Review Notification] **/prreviews Read PR Reviews from Github and write them locally on a dedicated folder**
    2. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/redo Make necessary code/test/docs changes, build & test** & Review against Issue-specific & Global Spec  [User Approval Gate with AI Review Suggestions]
    3. Loop until 1 is sucessfull or go back to a previous step
        1. [Make plan first & keep updating the plan at every step] **/review Independent Code Review** (including Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled)
        2. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/redo Make necessary code/test/docs changes, build & test.** & Review against Issue-specific & Global Spec [User Approval Gate with AI Review Suggestions]
    4. [Make plan first & keep updating the plan at every step] **/grillme Grill User with questions to see if he really understands changes since lst grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]
    5. **/document Document**:  Final User Documentation (how to install & use the product) & Contributor Documentation (how to understanding the codebase), both in the form of step by step tutorial. [User Approval Gate with AI Review Suggestions]

    ---- New Session ----

    6. Loop until 1 is succesfull
        1. **/secreview Specialized Independent Security Review** flagging critical problems & warnings 
        2. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/redo Make necessary code/test/docs changes, build & test.** & Review against Issue-specific & Global Spec [User Approval Gate with AI Review Suggestions]
    7. [Make plan first & keep updating the plan at every step] **/grillme Grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]
    8. **/pr Push & Open PR with Summary of Changes, PR has to link the Issue it solves. Never merge automatically.**

## AI-assisted Coding Tech Stack
- Agent-native Editor/Terminal: **Superset** (only when tackling multiple issues in parallel)
- Terminal Agent Harness: **OpenCode***
- IDE (for better introspection + manual editing): **VSCode**
- LLM Provider Subscription: **OpenCode Go**
- Model: current best coding model available in the subscription
- Feature Flow: Start with just making each step a slash command, and leave to the developer to follow the steps (note: this doesnt enforce step execution, se requires developer commitment). Note: each slash commands should remember to update the feature building progress at the end (feature_flow.md file) or, if its the first step, create the file if its still not created). Each slash command should have a simple name & also use the sub-agent that is most appropriate for the step.
- OpenCode Plugins (package skills + slash commands + sub-agents ... until a unit): **signoz plugin (only when this tool is used in the codebase), opencost plugin (only when this tool is used in the codebase), headlamp plugin (only when this tool is used in the codebase), graphify, rtk, security-guidance, code-review, code-simplifier, explanatory-output-style, skill-creator and summarize-session.** 
- OpenCode MCPs: Github MCP.
- Custom Freelunch OpenCode sub-agents: **spec-specialist, security-specialist, refactoring-specialist, cloud-architect, platform-engineer, testing-specialist, debugging-specialist, grilling-specialist, proxmox-specialist, talos-linux-specialist, linux-specialist, k8s-specialist, sre, cicd-specialist, golang-specialist, typescript-specialist.** Agency-agents provides some agents out-of-the-box
- Custom Freelunch OpenCode Slash Commands: one for each unique step of the feature building flow
- Dependency Docs: a `dependency_docs` folder holding the documentation of the pinned version of every first-level repo dependency (by default) with two sub-folders: `system_dependencies` (e.g., proxmox, kubetcl, etc) and `libraries` (e.g., go libraries, node libraries). Can also add lower level repo depdendencies when debugging (adding a docs of a dependency of a dependency to help debugging)
- Custom Skills: created on-demand, under a `custom_skills` folder, via manual creation or via `skill-creator` to avoid having to repeat the same solution process over and over.
