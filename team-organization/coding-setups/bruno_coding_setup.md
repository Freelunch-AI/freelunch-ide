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
- when E2E tesitng a product: be picky about the UI you see and be obsessed with pixle perfection. 
- If something looks off (even if tis not directly related to the things youa re doing) try to get it fixed along
- If you realize that you are stuck in a loop where you did them same action multiple times, you need to change your approach or even reset to the latest commit if everyhting is chaotic
- Always write code filled with debug logs to help you debug in dev phase. Dont worry, there is a scheduled step later that is dedicated for you to remove these excessive logs which arent good for production.
- Only call sub-agents if you are having difficulty doing it on your own and need a fresh view point froma specialist (e.g., stuck in a feature bug -> call debugging specialist; stuck in some testing error -> call testing specialist, etc)
- Before creating a skill from scratch for a common thing (not project-specific, e.g., frontend design) search for existing skills in skills.sh which can be installed via npx skills add
- For UI work (e.g., making a button or siebar), always create multiple mocks in the same html file before implementation. So that we can know precisely what we want to build. Each mock should only show what we interested in (not the whole UI which is already mocked in docs/mock.html) have an ID. I will say which ID I chose or give feedback for you to regenerate.
- if you encounter code-spec mismatch you should explain the mismatch, initiate a discussion with the user, which wil culminate in either code or spec change (or both). Spec should always be the goldern standard we look up to, so it can never be outdated.

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
- can go back from a step to a previous step if necessary to fix issue created earlier (but step jumps need to be tracked in an feature_flow.md file which has the name and number of issue on its title)
- when a new feature starts, first need search for any completed feature_flow.md (inside .gitignore) and store it in completed_feature_flow folder (inside .gitignore) in the form feature_flow[i].md
- session summary hook: when a session ends store a summary of key things done/key problems encoutered/tips/learnings/todos in the session in the respective section of feature_flow.md that agent was in (e.g., under step 3 or step 12) in this json form {"key things done": "placeholder", "key problems encoutered": {"problem":" placeholder", "solved_or_not": placeholder, "tips for next agent working on this": "placeholder"}, "learnings": "", "todos": "placeholder"}
- approval gates mean that either the user (developer) or a specific AI agent needs to give approval to continue the flow
- after every approval (human or ai), a git commit is made
- multiple features can be implemented in parallel by having separate terminals, each one in its respective worktree and branch.
- "AI Review" menas the same AI thats coding reviews its own work
- "Independent AI Reviewer" means that a different model with fresh context must be used

Flow:

A: Issue-specific Spec

1. **/start Start Feature Building: point to github issue, agent will read the issue and create feature branch with appropriate name according to the branching strategy file**
2. **/clarify Ask User Clarifying Questions & do web search if necessary** [User Approval Gate with AI Review Suggestions]
3. Loop until 2 is succesfull [User Approval Gate with AI Security Reviewer Suggestions]
    1. **/spec Build issue-specific Spec (prd.md + architecture.md + tech_stack.md under ./issue folder which should be mentioned inside .gitignore)**

    ---- New Session (reset context) ----
   
    2. **/reviewspec** Review Spec with Indepdente AI Reviewer, make shure to also check consistency with Global Spec (Founding Doc + Roadmap + Tech Stack), possibly catching things not specified in spec and problems in spec that needed to be overruled [User Approval Gate]

    ---- New Session (reset context) ----
    
    3. **/specsecreview Specialized Spec Security Review** flagging critical problems & warnings 
4. [Make plan first & keep updating the plan at every step] **/fspecgrillme Understand the feature spec, then grill User with questions to see if he really understands the feature spec, user reviews feature spec and asks questions until he has full understanding** [AI Approval Gate]

---- New Session (reset context) ----

B: Boilerplate

5. **/boilerdep Define Allowed boilerplate dependencies** (e.g., programming language, build tool, testing tools, package manager, etc) [User Approval Gate with AI Review Suggestions]
6. [Make plan first & keep updating the plan at every step] [User Approval Gate with Indepedent AI Plan Reviewer Suggestions] **/boiler Setup/Modify the code boilerplate (stucture/skeleton)** (directories, files, functions, classes, types, docstrings, test build command, final packaging build command), install boilerplate depedencies & Review against Issue-specific & Global Spec (PRD + Architecture + Tech Stack) catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled [User Approval Gate with AI Review Suggestions]
7. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands the commit, user review code and asks questions until he has full understanding** [AI Approval Gate]

---- New Session (reset context) ----

C: Tests & Logic

8. Loop until 2 is sucessfull [User Approval Gate with AI Independent Reviewer Suggestions]
    1. [Make plan first & keep updating the plan at every step]  [User Approval Gate with Indepedent AI Plan Reviewer Suggestions] **/writetests Write/Modify the functional tests** (unit tests, integration tests) & Review against Issue-specific & Global Spec [User Approval Gate] with AI QA Review Suggestions] 
    2. **/testtests Test the functional tests with placeholder feature code and guarantee 100% test coverage**
9. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]
10. Loop until 4 is sucessfull [User Approval Gate with AI Independent Reviewer Suggestions]
    1. **/featdep Define Allowed feature code dependecies** [User Approval Gate with AI Review Suggestions]
    2. [Make plan first & keep updating the plan at every step]  [User Approval Gate with Indepedent AI Plan Reviewer Suggestions]**/feat Write feature code using only the allowed feature code dependencies & Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled**. Important: should make one test pass at a time, do not try to code multiple things in parallel and make multiple tests pass at once. [User Approval Gate with AI Independent Reviewer Suggestions]
    3. [Make plan first & keep updating the plan at every step]  [User Approval Gate] **/fixstatic Fix Linting errors, Static Analysis & Simplify Code** [User Approval Gate with AI Independent Reviewer Suggestions]
    4. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/test Build and Test feature code with the functional tests, generate testing & test coverage reports** & Review against Issue-specific & Global Spec, repeat this step until all tests pass
11. [Make plan first & keep updating the plan at every step] **/evalrefact evaluate refactoring opportunities that would improve code quality and maintanability, [Make plan first & keep updating the plan at every step] then implement the chossen refactoring bits one by one, after each one is done, evaluate if it actually is better than before (if not, just keep ho it was before), only then move to the next** [AI Griller Approval Gate]
12. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]
13. [Make plan first & keep updating the plan at every step] **/stripdebuglogs** Remove debug logs from the code, only leave essential logs [User Approval Gate with AI Reviewer Suggestions]

---- New Session (reset context) ----

D: Code Review & Documentation

14. Loop until 1 is sucessfull or go back to a previous step [User Approval Gate with AI Independent Reviewer Suggestions]
    1. [Make plan first & keep updating the plan at every step] **/review Independent Code Review** (including Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled)
    2.  [Make plan first & keep updating the plan at every step] [User Approval Gate]**/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec [User Approval Gate with AI Reviewer Suggestions]
15. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands the changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]
16. **/document Document**:  Final User Documentation (how to install & use the product) & Controbutor Documentation (how to understanding the codebase), both in the form of step by step tutorial. [User Approval Gate with Independent AI Reviewer Approval Suggestionse]

---- New Session (reset context) ----

E: Security Review & PR

17. Loop until 1 is sucessfull or go back to a previous step [User Approval Gate with AI Reviewer Suggestions]
    1. **/secreview Specialized Security Review** flagging critical problems & warnings
    2. [Make plan first & keep updating the plan at every step] [User Approval Gate]**/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec [User Approval Gate with AI Reviewer Suggestions]
18. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands the changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]
19. **/pr Push & Open PR with Summary of Changes, PR has to link the Issue it solves. Never merge automatically.**

---- New Session (reset context) ----

F: Make fixes based on PR Reviews and/or CI failures until PR is merged

20. Loop until 1 is sucessfull or go back to a previous step 
    1. (On PR Review or CI failure Notification manually checked by user) **/prreviews Read PR Reviews & CI Run from Github and write them locally on a dedicated folder**
    3. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/redo Make necessary code/test/docs changes, build & test** & Review against Issue-specific & Global Spec  [User Approval Gate with AI Review Suggestions]
    4. Loop until 1 is sucessfull or go back to a previous step
        1. [Make plan first & keep updating the plan at every step] **/review Independent Code Review** (including Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled)
        2. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/redo Make necessary code/test/docs changes, build & test.** & Review against Issue-specific & Global Spec [User Approval Gate with AI Review Suggestions]
    5. [Make plan first & keep updating the plan at every step] **/grillme Grill User with questions to see if he really understands changes since lst grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]
    6. **/document Document**:  Final User Documentation (how to install & use the product) & Contributor Documentation (how to understanding the codebase), both in the form of step by step tutorial. [User Approval Gate with AI Review Suggestions]

    ---- New Session (reset context) ----

    6. Loop until 1 is succesfull
        1. **/secreview Specialized Independent Security Review** flagging critical problems & warnings 
        2. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/redo Make necessary code/test/docs changes, build & test.** & Review against Issue-specific & Global Spec [User Approval Gate with AI Review Suggestions]
    7. [Make plan first & keep updating the plan at every step] **/grillme Grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]
    8. **/pr Push & Open PR with Summary of Changes, PR has to link the Issue it solves. Never merge automatically.**

## AI-assisted Coding Tech Stack
- Agent-native Editor/Terminal: Superset (only when tackling multiple issues in parallel)
- Terminal Agent Harness: OpenCode
- IDE (for better introspection + manual editing): VSCode
- LLM Provider Subscription: **OpenCode Go
- Models: (0) Planning: Kimi K3 with medium resoning; (1) Spec, Boilerplate and Tests: Kimi K3 with high reasoning; (2) Logic Writing: DeepSeek V4 Flash with medium reasoning; (3) Independent Code Review Model: GLM-5.2 with high reasoning; (4) Security Review: Kimi K3 with high reasoning; (5) Fixing PR Review o CI Problems: Kimi K3 with high reasoning; (6) Independent spec and plan Review Model (only at few key moments): Claude Ops 5 with high reasoning; (7) sub-agents model: Claude Opus 5
- Feature Flow: Start with just making each step a slash command, and leave to the developer to follow the steps (note: this doesnt enforce step execution, se requires developer commitment). Note: each slash commands should remember to update the feature building progress at the end (feature_flow.md file) or, if its the first step, create the file if its still not created). Each slash command should have a simple name & also use the sub-agent that is most appropriate for the step.
- OpenCode Plugins: graphify, rtk, lavish-axi and summarize-session, cross-platform-screenshot-capture.
- OpenCode MCPs: Github MCP, Playwright MCP (only when making the ide).
- Custom Freelunch OpenCode sub-agents: security-specialist, code-review-and-refactoring-specialist, testing-specialist, debugging-specialist. Agency-agents repo provides some agents out-of-the-box.
- Custom Freelunch OpenCode Slash Commands: one for each unique step of the feature building flow
- Dependency Docs: a dependency_docs.md under ./docs with entries in the form "- <dependency>: <docs_link>" for all direct dependencies (not dependencies of dependencies). Pinned versions used in the project cna be seen in the lock file of the virtual dev environment tool.
- Skills: 
    - Custom Skills: created on-demand, under a `custom_skills` folder, via manual creation or via `skill-creator` to avoid having to repeat the same solution process over and over. Make these: codebase-grillme, fspec-grillme, make-spec, understand-external-codebase, bash-cheatsheet.
    - Use existing skills: skill-creator, grillme, lavish-axi
- Ralph Loop Engine (For agents to without supervision to achieve a goal, usefull for when you sleep/eat/or are just living life): **good nigh, have fun (gnhf)**. Note: this will burn tokens, only use if you are pretty confortable token-wise.

## Custom Plugin PRD: Usage Guard Plugin

### Overview

Usage Guard is an OpenCode plugin that continuously monitors local LLM token consumption and proactively warns developers when their current usage pattern is likely to trigger subscription throttling.

Unlike cost dashboards, Usage Guard is subscription-aware. Users configure their LLM subscription (starting with **OpenCode Go**) and the plugin estimates throttling risk using known characteristics of that subscription model together with the user's recent usage patterns.

The plugin also analyzes the user's workflow and provides actionable recommendations to reduce token consumption before limits are reached.

---

### Problem

Developers using subscription-based LLM plans receive little or no warning before hitting provider rate limits.

Once throttled, they must interrupt their workflow, switch models, or wait for their rolling usage window to recover.

The goal is to detect aggressive usage **before** throttling occurs and help users adjust their behavior.

---

### Goals

* Continuously monitor local token usage.
* Understand the user's subscription plan.
* Estimate throttling risk using subscription-specific heuristics.
* Display warnings directly inside the terminal.
* Recommend practical ways to reduce token usage.

---

### Non-Goals

* Predict the exact remaining provider quota.
* Reverse engineer provider rate-limiting algorithms.
* Display API billing or costs.
* Synchronize usage across multiple machines (MVP).

---

### Supported Subscriptions

#### MVP

* OpenCode Go

#### Future

* OpenCode Max
* Claude Pro
* Claude Max
* ChatGPT Plus/Pro
* Gemini subscriptions
* API usage profiles

Each subscription profile defines heuristics appropriate for its known usage model.

---

### Terminal Scope

#### MVP

Usage Guard monitors a **single OpenCode terminal session**.

Risk estimation, rolling statistics, and warnings are computed only from activity originating in the current terminal. This provides immediate value while keeping the implementation simple and avoiding cross-process coordination.

#### Future

Support multiple concurrently running OpenCode terminals.

The plugin should distinguish between:

* **Terminal Risk** — Usage and throttling risk for the current terminal.
* **Global Risk** — Combined usage across all active OpenCode terminals.

Example:

```text
Current Terminal
🟡 Moderate Risk

Global Usage
🔴 High Risk

Reason:
Three OpenCode terminals are actively consuming tokens.
```

The architecture should be designed around a per-terminal abstraction so global aggregation can be added without changing the core risk estimation logic.

---

### Data Source

Read usage data from OpenCode's local SQLite database.

Available metrics include:

* Timestamp
* Provider
* Model
* Input tokens
* Output tokens
* Cached tokens
* Request duration

---

### Core Metrics

Continuously compute rolling statistics:

* Tokens/minute
* Requests/minute
* Last 5 minutes
* Last 10 minutes
* Last 30 minutes
* Exponential Moving Average (EMA)
* Context growth rate
* Cache hit ratio

---

### Risk Estimation

Estimate throttling probability by combining:

* Subscription profile
* Rolling token velocity
* Request frequency
* Historical usage patterns

Risk levels:

* 🟢 Low
* 🟡 Moderate
* 🔴 High

The heuristics should be configurable and evolve as community knowledge improves.

---

### Intelligent Recommendations

When high usage is detected, Usage Guard should explain **why** and recommend ways to reduce token consumption based on observed patterns.

Examples:

* Large contexts are being sent repeatedly.
* Long-running agent sessions are consuming most tokens.
* Context is growing continuously without resets.
* Cache hit rate is low.
* An expensive model is being used for simple tasks.
* Multiple agents are running simultaneously (future multi-terminal support).
* Frequent retries are increasing token usage.

Possible recommendations:

* Start a fresh conversation.
* Switch to a cheaper model for routine tasks.
* Reduce the number of concurrent agents.
* Enable or improve prompt caching.
* Break large tasks into smaller batches.
* Reduce unnecessary context included in prompts.

Recommendations should be personalized based on actual usage patterns rather than generic advice.

---

### Terminal Warnings

Warnings should appear directly inside the active OpenCode terminal. Users should not need to open a dashboard or separate interface.

Example:

```text
⚠ Usage Guard

High usage detected.

Subscription:
OpenCode Go

10-minute token rate is significantly above your expected sustainable usage.

Risk:
🔴 High

Suggestions:
• Start a new conversation
• Switch to Gemini Flash for routine tasks
• Reduce concurrent agents
```

Warnings should be lightweight, infrequent, configurable, and should not interrupt the normal development workflow.

---

### Configuration

```yaml
usageGuard:
  enabled: true

  subscription: opencode-go

  windows:
    - 5m
    - 10m
    - 30m

  warningThreshold: medium

  notifyEvery: 5m
```

---

### Future Enhancements

* Multi-terminal monitoring
* Global usage aggregation
* Desktop notifications
* Historical usage graphs
* Burn-rate forecasting
* Provider-specific optimization recommendations
* Community-maintained subscription profiles
* Prometheus/OpenTelemetry metrics export
* Automatic model-switch suggestions
* Team-wide usage analytics

