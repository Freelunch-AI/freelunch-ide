# Bruno's AI-assisted Coding Setup

Run everyting inside a local dev docker container for safety reasons. VSCode connects easily via Remote Connection.

## Project Rules (AGENTS.md)

### Things you should NEVER do

- never acess directories above the project's root directory
- never read sensitive files (e.g., .env)
- never run destructive commands without my permission. Always present me a dry run alternative if possible.
- never try to extract extreme performance (speed, throughput) at the expense of making code complexity significantly higher, unless specifically prompted for extreme performance optimization
- never try to bypass pre-commit hooks or ny other quality gate, unless explicitely asked to do so by me.
- never modify opencode-related folders/files such as /.agents, AGENTS.md and skills, without my explicit approval

### Things you should always do

- Before starting a task always read the global and issue-specific spec. Treat global spec as the main source of truth. If issue-specific spec differs from global spec, flag this issue for me to resolve (with your help). If implementation differs from issue-specific spec or global spec, flag this issue for me to resolve (with your help).
- Before starting something new, check the last uncommitted and commited changes made with git. Only start this new thing if nothing seems suspicious (e.g., new code was written without correspinding tests, weird code changes, etc)
- Before using an unfamiliar dependency/API, consult its official documentation relevant to the operation being performed. Do not reread documentation already understood in the current session.
- every review document (in ./.agent/session-persistence-candidate/reviews/) created should contain in its initial metadata a reference to the exact version of what was reviewed which can be a file of specific commit (e.g., spec review and security spec review) or an entire commit (e.g., code review and security code review).
- log all mistakes you made in ./.agent/persistent/knowledge/mistakes.jsonl file, each entry in the form {"what_was_done": "placeholder", "what was wrong": "placeholder", "why it was wrong": "placeholder", "how the mistake was corrected": placeholder}. 
    - What counts as mistakes?
        - Anything you realize you did wrong before, having evidence to support why its wrong and explanation of why its wrong
        - Anything the I (aka your user) had to intervene to change something you already did becomes it had serious problems. I might say explicitely that you did something wrong (e.g., "change di code you wrote because its not readable", "change these tests you wrote becaue they dont reflect the spec", "change your implementation plan to more fine-grained end-to-end steps, where you start by") or just ask you if you are shure something is correct. Note: before counting it as a mistake and changing it, you must confirm the problem by talking to me with arguments. 
    - Whats doesnt count as mistakes? (1) User intervetions where the user requests spec changes/updates are never mistakes; (2) User interverntions for non-critical low-level changes that to satisfy his preferences (e.g., "extend this class to also support this capability that is currently outisde of the class")
- Log all codebase-related (e.g., how does this function work?) questions I ask to you in a ./.agent/persistent/user-codebase-questions.jsonl, each entry in the form "question": "placeholder", "answer": "placeholder", "branch": "placeholder", "commit_hash": placeholder"} where branch is the branch in which the user is in when he asked the question and commit_hash is the last commit at the time the user asked the question.
- Log all new usefull nuggets of codebase knowledge to ./.agent/session-persistant-candidate/non_obvious_conjectures.md where you keep track of non-obvious conjectures you make for the session. Each conjecture has the following data: (1) description of the conjecture, current evidence of the conjecture, already a fact? (yes or no) and risk if wrong (high, medium or low). Always ask question to user before you are about to act on a risky conjecture you dont have much evidence.
- keep your code clean and organized, refactoring might be needed
- moduarization: the codebase should have a few big modules with clear boundaries and relationships, and each big module is composed of many little modules. Dont let fils become too big, prefer breaking into multiple files where each one has a clear meaning/job.
- if using bash commands for file/content search: prefer `fd` (fdfind) and `rg` (ripgrep) over standard `find` and `grep` for better performance and git-awareness.
- always make a plan before doing stuff
- before concluding a task, critically re-evaluate your reasoning, assumptions, and implementation. Verify that the solution satisfies the user's objective, that no possibly affected areas have been overlooked, and that no unnecessary regressions have been introduced.
- when E2E tesitng a product: be picky about the UI you see and be obsessed with pixle perfection. 
- If something unrelated looks wrong, record it as a todo in ./.agent/session/todos.md it creates a correctness, security, build, or test failure affecting the current task. do not do one todo item at a time, batch them into related todos, and implement one batch at a time.
- If you realize that you are stuck in a loop where you can't solve a specific problem, you need to change your approach. Save the problem description in ./.agent/session/problem_stuck.md (should have problem summary, how to verify solution of problem, things already tried and the output gotten in each attempt) and then try these in order:
  1. Just reset the opencode context window and try again (the problem description will be at ./.agent/session/problem_stuck.md). 
  2. Reset the opencode context window and Change the OpenCode Coding Model (e.g., change from Kimi K3 to Opus 4.8) and increase it's reasoning if possible
  3. Reset the opencode context window and Take a fresh look at the problem. Try a different approach (if necessary you can update plan.md or core-implementaiton-tasks.md)
  4. Spawn the most approppriate sub-agent to help
  5. Reset to the latest commit as your last resource, then reset the OpenCode context window, then try agian to solve the problem.
- If you solve the problem you are stuck, delete ./.agent/session/problem_stuck.md immediatelly.
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
- Wehn doing any kind of artifact optimization (e.g., function performance optimization): never switch the current implementation for a candidate one before comparing both on the relevant evaluations. IF the candidate beats the current in the final eval score, than you can change, and the candidate then becomes the current implementation.

### Project Directories and Terminology

- Explanation of the /.agent directory structure is in ./agent/directory_structure.md. It explains the directory and files you will be using in sessions and across sessions for doing work effectively.
- Explanation of the terminology used in this project is in ./docs/terminology.md, always check it out when confused about terminalogy I use or which terminology you should use.

### Who are you (the AI agent)

You are a rigorous platform engineer working on the Freelunch IDE project, specifically working towards the Demo. You follow modern software engineering best practices such as: SOLID principles, tests-first-development, test coverage, dependecy injection, etc. You should always aim for the simplest approach by default, not the perfectly optimized/scalable/fastest one. Always review that you've done every task i asked you to do before saying you have finished.

### How you should treat me (the human user thats using you to code)

You should treat the me as the CEO thats sets objectives for you to build and also reviews your work. You should always explain to me everything you want to do/did the most step by step way. I may be wrong sometimes, therefore you should always reason about what I say and provide your take before a final decision. I may sometimes ask for things there are to vague/broad/ambiguous that require more specification to implement, or for things that might be incosistent with spec (globla spec: founding_doc.md, tech_stack.md, roadmap.md; issue-specific spec: prd.md, tech_stack.md, architecture.md) in this case you should ask for clarifying questions.

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
- if want to change core-implementation-tasks-plan.md: you should not overwrite the file, you should append to it the reason of the replanning and a summary of the current state of the codebase, then append a new core implementation tasks plan graph aling with an explanation of what was changed and why. So the resulting file will actually contain (in order) the previous core tasks plan and the new one.

### Big Refactoring mid-coding: rewriting tests and/or modifying scaffold (directories, files, interfaces, data models, etc; the skeleton in which logic gets written inside)

You might be doing a step and realize your tests and/or scaffold (directories, files, interfaces, data models, etc; the skeleton in which logic gets written inside) needs to be changed in multiple ways. You should create and checkout to an epehemeral big-refacoring branch (which was sourced from the current branch, not main) and then do the refactoring in the big-refacoring branch. When you are done, ask for my approval to merge big-refacoring into the issue handling branch.

Small or Localized refactorings you can just do, without this branching and asking for my approval ceremony.

### Patterns to use

- testing folder that mimicks the actual folder structure
- always work on the standard virtual environment for the project
- test-first development where tests are wrritne before code (this is not: strict tdd where need to make one function red -> green at a time and only then to move to the next)
- Use dependency injection where it improves testability or separation of concerns; 

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
- Only when doing code review: Use the tool lizard (https://github.com/terryyin/lizard) for measuring clyclomatic complexity of functions. The results will help you analyze if some functions are too complex or not.

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

- reproduce the bug in an E2E (as much as possible) setting as closely aligned to the end use to make shure you are solving the actual usage problem
- After bug reproduction, use the smallest relevant tests for diagnosis and iteration.
- Find the Root Cause of the Bug. Make hypothessis and test your hypotheses until you find the real hypothesis that is the root cause of the bug.
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

- Implent first (Always Required): Unit tests 
- Implement second (Always Required): Integration tests 
- Implement third (only required if the project already provides an external API or GUI): End-to-end tests

Skipping any required level = Not Complete

> This is the end of the AGENTS.md file

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
- "AI Review" means the same AI thats coding reviews its own work
- "Independent AI Reviewer" means that a different model with fresh context must be used
- When a session ends, a hook must be called to: (1) empty all files inside .agent/session and (2) analyze all files recursevly in .agent/session-persistant-candidate to check if there is knowledge that is still true or probabably true for the current state of the codebase, then transfer the usefull knowlegde to .agent/persistant/knowledge/non_obvious_conjectures_and_facts.md, then finally empty all files inside .agent/session-persistant-candidate; (3) assess if the current harness-related files (AGENTS.md, skills, model configs, etc) need improvement (or bloat removal) based on your analysis, traces (stored in ~/.opencode-trace/), past mistakes, feedbacks received, post-mortems, new knowledge obtained and make appapriate harness improvements if necessary using harness-eval skill for ensuring the proposed change actually improves the harness and remember to git commit if changes were made after evals. But before making changes to harness-related files you need my approval (user approval), so ask for my approval explaining the changes made and the eval results that are evidence the the change will improve the harness.
- At the start of any slash command: a hook must be called to git add & commit if there were not commited changes made. If in a new session without context fo what was done: use .agent/persistent/current-issue/flow/issue_flow.md's last progress data to infer a good commit message.
- If stuck in a loop, developers should try changing the implementer model from the default to other 
- Ensure language-specific (in our case here: go and typescript) linters and formatters are running continuosly on every file edit, using opencode hooks
- have a git pre-commit hook that: builds (must be sucessfull without warnings), runs tests (all must pass) and test coverage (> 90% required).

### Issue Flow (note: bug, refactoring or performance issue handling allow skipping multiple steps that are required for a new feature, skip when you deem the step unnecessary):

A: Issue-specific Spec & Core Implementaion Tasks Plan

0. **/grillme Understand the codebase, then grill me with questions to see if he really understands the codebase.** Make high-level (e.g., decisions chosen, project strcture, tradeofs, architecture) questions and low-level ones as well (e.g., what a specific file/function/class is for) [AI Approval Gate]
1. **/start Start Issue Handling**: point to github issue, agent will read the issue and create satelite branch with appropriate name according to the branching strategy file. Will then study the repo and ask user clarifying questions about probem and solution. This step ends when a common problem & solution understanding is reached with the user.
2. {only if issue of type == bug} **/reproducebug** Reproduce the reported bug in the issue by writing and running a failing test, and confirm it fails for the expected reason described in the bug issue description.
3. Loop until 2 and 3 are succesfull [User Approval Gate with AI Security Reviewer Suggestions]
    1. **/spec Build issue-specific Spec (prd.md + architecture.md + tech_stack.md under .agent/issue-spec folder**

    ----<<separate terminal block (reset context) >>----
    
    2. **/reviewspec** Review Spec with Indepedente AI Reviewer, make shure to also check consistency with Global Spec (Founding Doc + Roadmap + Tech Stack), possibly catching things not specified in global spec and problems in global spec. Incosistencies between both specs should be flagged to the user with recommendations. If you detect a problem in global, spec first modify global spec, then update issue-specific spec. Stores final spec review in .agent/session-persistence-candidate/reviews/final_sec_review_[timestamp].md [User Approval Gate]
    3. **/specsecreview Specialized Spec Security Review** flagging critical problems & warnings. Stores final sec spec review in .agent/session/reviews/final_specsec_review_[timestamp].md

    ----<</separate terminal block >>----

---- New Session (reset context) ----

4. **/make-tasks-plan Convert the spec into a graph of tasks (core-implementation-tasks-plan.md), each one depending on previous tasks or not depending on any** [User Approval Gate with Indepedent AI Plan Reviewer Suggestions]
5. **/tasksplangrillme Grill User with questions to see if he really understands the tasks plan until he has full understanding** [AI Approval Gate]

6: Core Implementation (p/task of the core-implementation-tasks-plan.md). Loop until core-implementation-tasks-plan.md is finished

---- New Session (reset context) ----

B1: Common scaffold

7. **/boilerdep Define Allowed scaffold dependencies (e.g., programming language, build tool, testing tools, package manager, etc)** You cannot use dependencies that are not actively maintained (search github to see last PRs). A specific dependency should can chosen over its alternatives if it has important capabilities other dont, if its signfificantly more popular, if its signficantly easier to use or if its signficantly older. [User Approval Gate with AI Review Suggestions]
8. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate with Indepedent AI Plan Reviewer Suggestions] **/boiler Setup/Modify the common scaffold (stucture/skeleton/foundation)** (directories, files, functions, classes, types, docstrings, data models (if statefull stuff is required), dev/test/build/package/publish command automations, etc) that are needed before core implementation, install scaffold depedencies & Review against Issue-specific (PRD + Architecture + Tech Stack) & Global Spec (PRD + Roadmap + Tech Stac) catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled [User Approval Gate with AI Review Suggestions]

B2: Tests & Logic

9. [Make/Remake plan.md first & keep updating the plan as you progress]  [User Approval Gate with Indepedent AI Plan Reviewer Suggestions] **/writetests Write/Modify the functional tests** (unit tests, integration tests) & Review the tests against spec (Issue-specific & Global Spec) to see if they are consistent. Incosistencies between tests/issue-specific-spec/global-spec should be flagged to the user with recommendations. [User Approval Gate with AI Independent Reviewer Suggestions] 
10. **/testtests Test the functional tests with placeholder feature code (all tests must fail in this phase) and guarantee test coverage of all core logic ((should be over 90% and ideally near 100% coverage))**

11. Loop until 3 is sucessfull [User Approval Gate with AI Independent Reviewer Suggestions]
    1. **/featdep Define Allowed feature code dependecies with explanation why use each dependency**. You cannot use dependencies that are not actively maintained (search github to see last PRs). A specific dependency should can chosen over its alternatives if it has important capabilities other dont, if its signfificantly more popular, if its signficantly easier to use or if its signficantly older [User Approval Gate with AI Review Suggestions]
    2. [Make/Remake plan.md first & keep updating the plan as you progress]  [User Approval Gate with Indepedent AI Plan Reviewer Suggestions]**/feat Write feature code using only the allowed feature code dependencies & Review against Spec and Deisgn System if doing GUI work (issue-specific spec and global spec) catching inconsistencies with spec/design system, things not specified in specs/design system and problems in spec/design system that needed to be overruled**. Incosistencies between code/tests;issue-specific-spec/global-spec should be flagged to the user with recommendations. Important: should try to make multiple related tests pass at a time, always aim for the smallest coherent behavioral slice, as end-to-end as possible across the components) that produces a useful feedback signal  [User Approval Gate with AI Independent Reviewer Suggestions]
    3. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate] **/test Build and Test feature code with the functional tests, generate testing & test coverage reports** & Review against Issue-specific & Global Spec, repeat this step until all tests pass

11. [Make/Remake plan.md first & keep updating the plan as you progress] **/refact-coreimplementation-task-if-necessary evaluate refactoring opportunities just for the last core implementaion task just realizsed. These should be changes that maintain same functionality but improve code clarity, quality and maintanability, [Make plan first & keep updating the plan as you progress]  [User Approval Gate] then implement the chossen refactoring bits one by one, after each one is done, evaluate if it actually is better than before (if not, just keep how it was before), only then move to the next** [AI Approval Gate]

 ----<<separate terminal block (reset context) >>----

12. **/grillme Understand the codebase, then grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]

----<</separate terminal block >>----

13. [Make/Remake plan.md first & keep updating the plan as you progress] **/stripdebuglogs** Remove debug logs from the code if thre are any [User Approval Gate with AI Reviewer Suggestions]

14. **/determine-if-core-implementation-task-is-done** should evaluate if the core implementaiotn tasks objectives were in fact achieved and if the core-implementaiotn-tasks.md still holds. If all the tests actually reflect the spec and all the tests pass, mark as done the the current core implementaiton tasks inside core-implementaiotn-tasks.md. Shouldnt rely on historic test passing outputs or reports, should run the entire test suite for the current core implementaiotn task and for previously completed core implementation tasks (to catch regressions if any).

---- New Session (reset context) ----

**Requirement to continue the flow: core-implementation-tasks-plan.md needs to be fully complete, i.e., all core implementation tasks finished**

15. **/end-to-end-testing:** Should test end-to-end to make shure the issue was completely handled. If code involved GUI, should have GUI tester to test its usability and how good it looks, emulating real user behaviours. Also should check if end-to-end tests actually reflect the issue's prd requirements. Also verify is the issue spec is still consistent with global spec.

---- New Session (reset context) ----

C: Code Review

16. Loop until 1 is sucessfull or go back to a previous step [User Approval Gate with AI Independent Reviewer Suggestions]
    1. [Make/Remake plan.md first & keep updating the plan as you progress] **/review Independent Code Review: review maintanability, modularity, test coverage, file sizes, tests compliance to spec, possibile GUI compliance to Design Doc. The last thing you must do: run mutation testing** (also identify things not specified in spec/design system (if doing GUI work) and problems in spec/design system (if doing GUI work)). This step will generate a .agent/session-persistence-candidate/reviews/final_code_review_[timestamp].md file
    2.  [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate]**/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec. Note: code review is stored in .agent/session-persistence-candidate/reviews/final_code_review_[timestamp].md [User Approval Gate with AI Reviewer Suggestions]

---- New Session (reset context) ----

D: Security Review, Documentation & PR

17. Loop until 1 is sucessfull or go back to a previous step [User Approval Gate with AI Reviewer Suggestions]
    1. **/secreview Specialized Security Review** flagging critical problems & warnings. Stores review in .agent/session/reviews/final_code_sec_review_[timestamp].md file
    2. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate]**/redo Make necessary ccode/test/docs changes, build & test** & Review against Issue-specific & Global Spec [User Approval Gate with AI Reviewer Suggestions]

18. Loop until 1 is sucesfull
        1. **/end-to-end-testing:** Should test end-to-end to make shure the issue was completely handled. If code involved GUI, should have GUI tester to test its usability and how good it looks, emulating real user behaviours. Also should check if end-to-end tests actually reflect the issue's prd requirements. Also verify is the issue spec is still consistent with global spec.
        2. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate] **/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec -- code, tests, specs should all be consistent with each other, if not flagg insconsistencies for the user to resolve [User Approval Gate with AI Indepentes Reviewer Suggestions]

 ----<<separate terminal block (reset context) >>----

19. **/grillme Understand the codebase, then grill User with questions to see if he really understands the changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]

----<</separate terminal block >>----

20. **/write-pr (1) Check the .agent/persistent/pre-pr-checklist.md to see if everthing was done/updated. (2) Write the PR to ./.agent/current-pr/pr.md

21. **/document Document**: (A) Contributor Documentation (how to understand the codebase, in the form of step by step tutorial from high-level to low-level going through pieces of files). (B) Final User Documentation only if its already usable (how to install & use the product) [User Approval Gate with Independent AI Reviewer Approval Suggestionse]

22. **/send-pr (1) Pushes and Opens Draft PR using ./.agent/current-pr/pr.md

---- New Session (reset context) ----

E: Make fixes based on PR Reviews and/or CI failures until PR is merged

22. Loop until 1 is sucessfull or go back to a previous step 
    1. (On PR Review or CI failure Notification manually checked by user) **/prreviews Read PR Reviews & CI Run from Github and write them locally on a dedicated folder**
    2. Loop until 1 is succesfull
        1. [Make/Remake plan.md first & keep updating the plan as you progress] **/review Independent Code Review: review maintanability, modularity, test coverage, file sizes, tests compliance to spec, possibile GUI compliance to Design Doc. The last thing you must do: run mutation testing** (also identify things not specified in spec/design system (if doing GUI work) and problems in spec/design system (if doing GUI work)). This step will generate a .agent/session-persistence-candidate/reviews/final_code_review_[timestamp].md file
        2. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate] **/redo Make necessary code/test/docs changes** & Review against Issue-specific & Global Spec -- code, tests, specs should all be consistent with each other, if not flagg insconsistencies for the user to resolve [User Approval Gate with AI Indepentes Reviewer Suggestions]

    ---- New Session (reset context) ----

    3. Loop until 1 is succesfull
        1. **/secreview Specialized Independent Security Review** flagging critical problems & warnings. Stores final sec review in .agent/session/reviews/final_code_sec_review_[timestamp].md
        2. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate] **/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec -- code, tests, specs should all be consistent with each other, if not flagg insconsistencies for the user to resolve [User Approval Gate with AI Indepentes Reviewer Suggestions]

    4. Loop until 1 is sucesfull
        1. **/end-to-end-testing:** Should test end-to-end to make shure the issue was completely handled. If code involved GUI, should have GUI tester to test its usability and how good it looks, emulating real user behaviours. Also should check if end-to-end tests actually reflect the issue's prd requirements. Also verify is the issue spec is still consistent with global spec.
        2. [Make/Remake plan.md first & keep updating the plan as you progress] [User Approval Gate] **/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec -- code, tests, specs should all be consistent with each other, if not flagg insconsistencies for the user to resolve [User Approval Gate with AI Indepentes Reviewer Suggestions]

     ----<<separate terminal block (reset context) >>----
    
    5. **/grillme Grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Approval Gate]

    ----<</separate terminal block >>----

    6. **/write-pr (1) Check the .agent/persistent/pre-pr-checklist.md to see if everthing was done/updated. (2) Write the PR to ./.agent/current-pr/pr.md

    7. **/document Document**: (A) Contributor Documentation (how to understand the codebase, in the form of step by step tutorial from high-level to low-level going through pieces of files). (B) Final User Documentation only if its already usable (how to install & use the product) [User Approval Gate with Independent AI Reviewer Approval Suggestionse]

    8. **/send-pr (1) Pushes and Opens Draft PR using ./.agent/current-pr/pr.md

> This is the end of the .agent/persistent/pre-pr-checklist.md file

## AI-assisted Coding Tech Stack

External Vendor Requirements: Opencode Go Subcription, Claude Credits, Github Repo

- Agent-native Editor/Terminal: Orca (only when tackling multiple issues in parallel)
- Terminal Harnesses: OpenCode + failproofai + claude-tap, open-code-review and openwiki
- IDE (for better introspection + manual editing): VSCode (with language-specific linters and formatters plugins running continuosly on every file edit)
- LLM Provider Subscription: **OpenCode Go
- Models: (0) Planning: Current great coding model thats not so expensive with medium resoning; (1) Spec, scaffold and Tests: Current great coding model thats not so expensive with high reasoning; (2) Core coding: Best coding qualirty/price model with medium reasoning; (3) Independent Spec Review: Best Coding Model with high reasoning; (4) Plan Review: Current great coding model thats not so expensive with high reasoning; (5) Security Review: Best Review Model with high reasoning; (6) end-to-end testing: Current great coding model thats not so expensive with high reasoning; (7) sub-agents model: Best Coding Model with high reasoning; (8) Code Review: 2 not so correlated great Review Models that are not so expensive with high reasoning as cadidate generators and open-code-review using best review model with high reasoning as final decision maker; (9) Grillme: Best coding qualirty/price model with medium reasoning
- Local Routing: opencode-model-router (opencode plugin)
    - Fast Model: qwen2.5-coder:7b
    - Medium Model: current coder model
    - Heavy Model: current coder model
- Issue Flow: Start with just making each step a slash command, and leave to the developer to follow the steps (note: this doesnt enforce step execution, se requires developer commitment). Note: each slash commands should remember to update the issue resolving progress at the end (issue_flow.md file) or, if its the first step, create the file if its still not created). Each slash command should have a simple name.
- OpenCode Plugins: graphify, rtk, opencode-quota, opencode-model-router, opencode-trace
- OpenCode MCPs: Github MCP
- Custom Freelunch OpenCode sub-agents: security-specialist, code-review-and-refactoring-specialist, testing-specialist, debugging-specialist. Tip: Agency-agents repo provides some sub-agents out-of-the-box.
- Custom Freelunch OpenCode Slash Commands: one for each unique step of the issue building flow
- Dependency Docs: a dependency_docs.md under ./docs with entries in the form "- <dependency>: <docs_link>" for all direct dependencies (not dependencies of dependencies). Pinned versions used in the project cna be seen in the lock file of the virtual dev environment tool.
- Skills: 
    - Custom Skills: created on-demand, under a `custom_skills` folder, via manual creation or via `skill-creator` to avoid having to repeat the same solution process over and over. Make these custom skills:
        - harness-eval (evalautes proposed harness changes on the current codebase, explained better in the end of this doc)
        - grill-my-understanding (continually ask questions of the latest changes to codebase to me, to see if i understand the codebase. Always give score my answers and give feedback to it. Only stop when you feel i understand the codebase. The user can also specifify specific files for you to grill him about instead of the entire codebase)
        - understand-external-codebase (1. Build a doc eplxianng in detail the characteristics and internals of an external github codebase; 2. Add to this doc an explanation of where and why this codebase can be helpfull as a reference for ideias/patterns for the current project being built)
        - update-fixed-context (1. Infers new usefull knowledge from ./.agent/persistant/knowledge/mistakes.jsonl and ./.agent/persistant/completed_issue_flows; 2. Add this new usefull knowledge to AGENTS.md if its not already there)
        - make-core-implementation-tasks-plan (transforms the spec into a graph of tasks, where: (1) each task can depend on other tasks being already done or not depend on any; (2) the tasks should not be of the form "one task implements each component that will be needed in this feature, e.g., oen task for the backend, another for databas,e another for gateway and another for frontend", the tasks should be done in the form of "one task implments a slice (governed by on or more integration/end-to-end tests) of multiple components, e.g., this task implements a funcitonality slice of frontend, backend, gateway and frontend that together brings us one step closer to our end goal and guarantee rich cross-component feedback along development". The core implementation tasks plan needs to be stored in .agent/persistent/current-issue/flow/core-implementation-tasks-plan.md. The core-implementation-tasks-plan.md. file should have the graph structure of the plan, where each node is a task. For each node there is also a pending/in-progress/done checkbox. Do not confuse with plan.md which is a per-step ephmeral small plan for step execution.)
        - document using openwiki and with the following guidelines in ./.openwiki/INSTRUCTIONS.md: should first check incosistency (global spec is king), staleness and incompleteness of existing documentation (if any) and then update/create: (1) Contributor Documentation: visualization o repositoty strcture explaining succintly each directory and file, (1.2) Step by step contributor tutorial to help a newcomer understand the codebase; (2) User Documentation (only do after first version 0.1.0 is released): (2.1) User API Reference. (2.2) User step by step tutorial starting from sratch; (2.3) User guides to do common stuff; (2.4) FAQ. 
        make shure the documentaiton explains well things that I usually have a hard-time understanding.
        - ui-taste (UI Taste gives Claude a visual sense of taste. Instead of relying only on abstract design principles, the skill provides curated examples of bad, good, and stellar GUIs across different application categories and problem modes, including screenshots and their underlying HTML/CSS. This gives the agent an understanding of what makes GUIs look good. The agent should launch the current GUI, identify the biggest visual shortcomings, and iteratively improve them. The goal isn't to force a particular design style—it is to help Claude distinguish "functional but mediocre" from "genuinely beatifull and easy to use", giving coding agents a practical visual benchmark for judging their own work.)
    - Use existing skills: visual-recap (ideally also running in github actions for every PR), skill-creator, i-have-adhd, quick-recap, mutation-testing, chrome-devtools-cli, optimize-anything, grillme (every grillmre run should log all the questions, answers and feedback gave to the user inot a .agent/persistent/user-grills/grill[i].md where i is the id of the grill and the file should have timestamp, commit, what the grill was about and grill score in the beggining of it. Ever grill should start by looking at the commit, what was grilled in the last grill and user-codebase-questions.jsonl file), lavish-axi, code-review-and-quality, api-and-interface-design, browser-testing-with-devtools (only when working with frontend part), security-and-hardening, cc-skills-golang, maintainable-typescript (only when working with frontend part), improve-codebase-architecture, screenshot (only when working with frontend part), extract-design-system (only when working with frontend part), frontend-design (only when working with frontend part).

## How Code Review is Done

1. Start a new session with opencode: do code review with Model A and store the review in .agent/session-persistence-candidate/reviews/[A]_code_review_[timestamp].md, where A is a placeholder for the actual mode name, and timestamp is placeholder for the actual timestamp
2. Start a new session with opencode: do code review with Model B and store the review in .agent/session-persistence-candidate/reviews/[B]_code_review_[timestamp].md, where B is a placeholder for the actual mode name, and timestamp is placeholder for the actual timestamp
3. Start a new session with open-code-review: do code review with Model C explicitely telling it to look at the candidate problems flagged inside .agent/session-persistence-candidate/reviews/ folder and store the resulting code review inside .agent/session-persistence-candidate/reviews/final_code_review_[timestamp].md

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
|   └── current-pr/pr.md
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
|   └── problem_stuck.md
|   └── debug-logs/
│
├── session-persistence-candidate/
│   └── knowledge/
│       └── non_obvious_conjectures.md
│       └── exploration_findings/
│           └── <name-of-exploration>_[timestamp].md
|   └── reviews/
|       └── final_spec_review_[timestamp].md
|       └── final_spec_sec_review_[timestamp].md
|       └── final_code_review_[timestamp].md
|       └── final_code_spec_review_[timestamp].md
│
```

### `persistent/`

Contains durable project knowledge and historical state that should survive across sessions and remain useful to future agents.

* `knowledge/non_obvious_conjectures_and_facts.md` — durable, non-obvious facts or evidence-backed insights about the codebase that are useful to future agents.

* `mistakes.jsonl` — records mistakes made by the agent that required user intervention (e.g. of user intervention: user asks "are you shure about this implementation, seems it can risk eposing credentials)". Each entry records what was done, what was wrong, why it was wrong, and how it was corrected.

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

* `current-pr/` — contains current pr made to solve the current issue.

* `completed-issues/` — archives durable artifacts from issues that have been completed.

  * `completed_issue_flows/` — contains archived `issue_flow_[i].md` files documenting how completed issues were implemented.
  * `completed-issue-specs/` — contains archived `issue-spec_[i]/` directories, including the PRD, architecture, and technology-stack documents for completed issues.

Persistent knowledge should contain only information that is expected to remain useful beyond the current session. Current-issue state belongs under `current-issue/` while historical issue state belongs under `completed-issues/`.

### `session/`

Contains temporary state for the current agent session.

* `plan.md` — the ephemeral step-by-step plan for the current slash-command execution. Plan progress needs to be tracked.
* `todos.md` — the ephemeral todo list used to store next things that need to be done. Need to check if an item is already done before actually doing it.

Everything in `session/` is disposable. At the beginning of a new session, the contents of `.agent/session/` are cleared. These files must not be treated as historical records of the project or issue. Durable workflow state belongs in `persistent/current-issue/`.

### `session-persistence-candidate/`

Contains information discovered during the current session that may be useful beyond the session but has not yet been promoted to persistent knowledge.

* `non_obvious_conjectures.md` — inferences or assumptions made during the current session, including their evidence, and their risk if wrong.
* `knowledge/exploration_findings/` — findings produced by exploratory sub-agents. Each exploration should produce a findings file containing the exploration context, description, model/agent used, token/cost information, reason it ended, findings, and conclusion.

At the beginning of a new session, useful information from this directory is reviewed and promoted into `.agent/persistent/knowledge/non_obvious_conjectures_and_facts.md` when appropriate. After promotion, the candidate directory is cleared.

### `directory_structure.md`

Documents the purpose and organization of the `.agent/` directory itself. It should be updated whenever the directory structure or the responsibilities of its files change.

### Important distinction: session state vs. persistent state

The `.agent/` directory deliberately separates **ephemeral working state** from **durable project state**:

* `session/` contains information needed only to continue the current session.
* `session-persistence-candidate/` contains potentially reusable discoveries/feedback that have not yet been validated or promoted.
* `persistent/` contains validated, durable knowledge and issue history.

Within `persistent/`, `current-issue/` is the active issue's durable workspace, while `completed-issues/` preserves the historical record of issues that have already been completed.


### General Rules

1. **Do not put source code in ****`.agent/`****.** Product code belongs under `./src/`.
2. **Do not treat ****`.agent/session/`**** as persistent storage.** Its contents may be deleted when a session starts.
3. **Do not promote assumptions or inferences automatically.** An assumption or infrence should only become persistent knowledge after it has been sufficiently verified.
4. **Do not silently modify historical issue records.** They are part of the project's implementation history.

> This is the end of the .agent/directory_structure.md file

## Custom PR Skill

Mandatory pre-pr checklist (ready to push and open PR):

- global spec, issue-specific spec and tests all consistent with each other
- all core (at least 90% total code coverage) logic covered by tests, with test coverage report evidence for it
- unit, integration and end-to-end (if possible) tests present and passing
- no big changes were made after the last independent code review was run
- no significant changes were made after the last security review
- spec, documentation and code are consistent + contrbutor docs cover all the code + user docs cover all Aexternal functionality exposed
- user understands the PR that will be made at the function interface/class interface/file/directory level

Making PRs:

- write PR to .agent/current-pr/pr.md (if there is already a file there, it should be deeleted first)
- ensure the Mandatory Pre-PR checklist is checked before making a PR
- follow the project's PR template
- in the pr you write:
        - highlight key decisions that might be questioned by a reviewer (e.g., why wasnt another distributed pattern chosen? why created a helper function instead of using a dependency? etc), problems encoutered, solutions and tradeoffs chosen
        - highlight what was tested and link to evidence that shows your test: (1) for standard tests: terminal output first shows the current commit hash and then runs the test suite and shows all tests are green, together with above 90% coverage; (2) if UI work: need before fix and after fix videocapture;
        - make a risk assesment of the PR (Low, Medium, High) based on how many changes it makes, the type of changes it makes, test coverage, etc

Handling PR Comments and Reviews:

- you should respond PR comments addresing the problems raised by them (PR comments are not always 100% correct, so you have to judge the comment's good and bad points (if any), then respond with your posistion on the concern raised by the comment) or just respond PR questions made by other developers
- you should make appropriate changes according to what you think should be changed after reading the comments
- if not shure what to change becasue of subjective/incosistent/ambiguous PR comments, then ask me clarifying questions.

---

---

## Harness Evaluation Custom Skill (harness-eval.md)

Use this skill to run controlled A/B evaluations of AI coding setups.

## When to Use This Skill

Use this skill whenever you want empirical evidence about whether a proposed change to an AI coding setup (harness change) is better or not than the baseline for the specific repo being worked on.

Typical uses:

* compare **Model A vs Model B**;
* compare **harness setup A vs harness setup B**;
* compare **a skill vs no skill**;
* compare **Skill A vs Skill B**;
* compare **Harness A vs Harness B**;
* evaluate harness changes currently present in the working tree;
* evaluate a harness change **before** making it;
* investigate whether an `AGENTS.md`, skill, hook, rule, subagent, model, or other harness change actually helps on the actual repoository being worked on.

The evaluation supports four task types:

```text
coding
spec / architecture / design
code review
doc review
```

The caller must explicitly specify which eval type(s) to run.

When `doc review` is included, the caller must specify the **exact document path(s)** to review.

Do not invent missing eval types or document paths.

---

# 1. Determine the Two Options

The evaluation always compares:

```text
Option A
Option B
```

The options can be provided explicitly **or inferred from the current working tree**.

## 1.1 Explicit Options

When the user provides two configurations, use them directly.

Examples:

```text
Model A vs Model B

Skill A vs Skill B

Harness setup before vs after

Harness A vs Harness B
```

Determine exactly what differs between the options.

---

## 1.2 Infer Options From Current Changes

When the user asks to evaluate their current harness changes without explicitly providing A and B, inspect the current working tree.

Look for changes to:

```text
AGENTS.md
CLAUDE.md
.agent/
.claude/
skills/
rules/
hooks/
subagents/
commands/
prompts/
harness configuration
model configuration
tool configuration
planning/review configuration
```

Also detect model changes represented by:

* configuration files;
* command arguments;
* environment-driven configuration;
* harness configuration.

Construct:

```text
Option A = current configuration with the evaluated change reverted
Option B = current configuration including the evaluated change
```

Preserve all unrelated user work in both options.

For example, if the working tree contains:

```text
modified application code
new skill
modified AGENTS.md
```

and the current evaluation is the new skill:

```text
Option A:
  current application changes
  old skill
  old AGENTS.md

Option B:
  current application changes
  new skill
  old AGENTS.md
```

The application changes are preserved in both options.

---

## 1.3 Infer Multiple Independent Evaluations

If the user asks to evaluate the current harness changes generally, identify independent relevant changes and create a separate experiment for each when possible.

For example:

```text
Change 1:
new code-review skill

Change 2:
modified AGENTS.md

Change 3:
Model A → Model B
```

Construct:

```text
Evaluation 1:
  A = no skill
  B = new skill

Evaluation 2:
  A = old AGENTS.md
  B = modified AGENTS.md

Evaluation 3:
  A = Model A
  B = Model B
```

Do not combine independent changes into a single experiment.

If one change depends on another and they cannot be meaningfully separated, evaluate them together and explicitly record the combined experimental variable.

---

# 2. Required Eval Input

Determine:

```text
EVAL_TYPES =
  one or more of:
    coding
    spec
    review
    doc-review
```

Also determine:

```text
DOCS =
  exact document paths
  required if doc-review is included
```

Option configuration may be:

```text
explicitly supplied
or inferred from current changes
```

If the eval type is missing, stop and request it.

If `doc-review` is selected without exact document paths, stop and request them.

---

# 3. FIRST: Isolate the Evaluation

**Before modifying anything, create the isolated evaluation environment.**

The user's current repository is the **main checkout**.

Treat it as read-only.

Do not:

* create the mutation there;
* install an evaluated harness there;
* run an evaluated agent there;
* modify harness files there;
* create evaluation artifacts there;
* modify `.agent/harness-evals/` while agents are running.

The evaluation must execute in **two separate Docker containers**, with **one independent repository per option**.

The main checkout is only the source from which the evaluation snapshot is created.

---

# 4. Snapshot the Current Working Tree

Capture the user's actual current state.

Do not use only `HEAD`.

The snapshot must include:

* committed files;
* tracked modifications;
* relevant untracked files;
* current application changes;
* current harness changes.

Do not discard user changes.

Exclude:

```text
.git/
.agent/harness-evals/
temporary evaluation state
temporary caches
credentials
secrets
```

Never copy secrets into the evaluation repository.

Record:

```text
HEAD commit
working-tree state
included untracked files
excluded paths
snapshot hash
```

---

# 5. Create an Isolated Evaluation Repository

Do not create evaluated repositories using a normal worktree from the original repository.

The original repository's Git history may reveal the ground truth.

Instead:

1. copy the evaluation snapshot;
2. create a new temporary Git repository;
3. initialize it;
4. add the snapshot;
5. create a single root commit.

The resulting repository must contain the exact snapshot state but **none of the original Git history**.

The evaluated agents must not be able to recover hidden ground truth with commands such as:

```bash
git show <original-commit>
git diff <original-history>
git log <original-history>
```

---

# 6. Determine Validation

From the isolated repository:

1. identify the normal test/build commands;
2. identify repository-specific validation;
3. run the known-good baseline validation;
4. confirm the baseline passes.

Record the exact commands.

Do not invent a new correctness oracle unless the eval definition explicitly requires it.

---

# 7. Create One Mutation Per Trial

For each trial, create **one mutation**.

The mutation is shared by both options.

Never independently generate a mutation for A and B.

Use:

```text
known-good evaluation repository
        ↓
create mutation once
        ↓
validate mutation
        ↓
freeze mutation
        ↓
clone exact frozen state
      ↙       ↘
 Option A   Option B
```

Record:

```text
mutation_id
seed
eval_type
target path(s)
target region(s)
mutation description
ground-truth content
known-good tree hash
mutated tree hash
```

---

# 8. Mutation Rules

Mutations must never alter dependencies or evaluation infrastructure.

## Never mutate dependencies

Do not:

* delete dependencies;
* add dependencies;
* remove dependencies;
* downgrade dependencies;
* upgrade dependencies;
* modify lockfiles;
* alter package manifests;
* alter vendored dependencies.

The dependency environment must remain identical for both options.

If a mutation requires changing dependencies, reject that mutation.

## Never mutate protected paths

Protect:

```text
tests/
.github/
CI/CD infrastructure
build infrastructure
dependency configuration
lockfiles
pre-commit configuration
evaluation infrastructure
harness-evaluation infrastructure
```

Add repository-specific protected paths as required.

---

# 9. Coding Eval

For a coding eval, delete application logic that is **covered by tests**.

This is mandatory.

Eligible targets include:

* functions;
* methods;
* classes;
* services;
* controllers;
* handlers;
* algorithms;
* business logic;
* orchestration logic;
* state transitions;
* modules;
* cross-module application logic.

Do not delete:

* tests;
* dependencies;
* dependency configuration;
* lockfiles;
* CI/CD;
* build infrastructure;
* evaluation infrastructure;
* harness infrastructure;
* generated/vendor code;
* uncovered application code;
* comments only;
* formatting only;
* trivial imports.

Prefer AST/semantic region selection.

Before accepting a mutation:

```text
known-good repository
    → relevant tests PASS

mutated repository
    → relevant tests FAIL

ground truth restored
    → relevant tests PASS
```

If the deleted code is not covered by tests or the mutation does not cause the relevant tests to fail, reject the mutation.

---

# 10. Spec / Architecture / Design Eval

Choose an existing technical specification, architecture, or design document.

Remove a meaningful region.

Possible targets:

* architecture documentation;
* design documents;
* ADRs;
* subsystem descriptions;
* component descriptions;
* interface descriptions;
* data flows;
* deployment architecture;
* technical rationale.

Do not modify application code or tests.

The original removed content is the hidden ground truth.

Give both agents the incomplete document.

After they finish, compare the reconstructed region with the original removed region.

The comparison is based on the original codebase document, not on a subjective judgment by the evaluator.

Do not give either agent access to the original removed content.

---

# 11. Code Review Eval

Start from the known-good application code.

Inject one or more subtle bugs.

A valid bug must:

* compile;
* be plausible;
* change behavior;
* be covered by tests;
* cause deterministic validation failure.

Examples:

* wrong comparison;
* boundary condition;
* off-by-one;
* wrong branch;
* incorrect default;
* stale state;
* incorrect error propagation;
* ordering;
* state transition;
* cleanup;
* timeout;
* mapping;
* concurrency.

Verify:

```text
known-good → PASS
mutated    → FAIL
```

Do not modify tests.

Do not tell the agents:

* how many bugs exist;
* where they are;
* which files changed;
* what type of bugs were injected.

---

# 12. Doc Review Eval

The caller must provide the exact document path(s).

For each specified document:

1. verify it exists;
2. read the document;
3. inject concrete repository-grounded defects;
4. preserve the original document as hidden ground truth;
5. freeze the mutated document.

Possible defects:

* incorrect technical statement;
* incorrect API description;
* outdated command;
* incorrect configuration;
* contradiction with implementation;
* incorrect example;
* broken reference;
* missing requirement;
* incorrect architecture statement.

Avoid subjective style-only issues.

Do not tell the agent what was changed.

After the agent finishes, compare the resulting document against the original.

---

# 13. Freeze the Mutation

After mutation creation:

1. run the required validation;
2. confirm the intended mutation is present;
3. confirm protected files are unchanged;
4. confirm dependencies are unchanged;
5. compute the mutated tree hash;
6. freeze the mutation.

The mutation must never change during the A/B comparison.

---

# 14. Create Two Separate Docker Containers

Create exactly two independent execution environments for the trial:

```text
Container A
Container B
```

They must be **separate Docker containers**.

Each container must have its own:

```text
filesystem
repository
HOME
agent state
harness state
temporary directory
writable cache
build directory
```

Do not use one container with two sequential agent sessions.

Do not run both options in the same container.

Do not mount one writable repository into both containers.

Do not share writable agent state or caches.

Use the same Docker image for A and B.

Record the immutable Docker image digest.

---

# 15. Create Two Independent Repositories

Create:

```text
option-a/repo
option-b/repo
```

Initialize both from the **exact same frozen mutation tree**.

Before applying the A/B difference:

```text
tree_hash(A) == tree_hash(B)
```

must be true.

Also verify:

```text
mutation_id(A) == mutation_id(B)
mutation_seed(A) == mutation_seed(B)
dependencies(A) == dependencies(B)
tests(A) == tests(B)
```

Do not start the agents until these checks pass.

---

# 16. Configure Option A and Option B

Apply only the intended experimental difference.

Examples:

### Model

```text
A = Model A
B = Model B
```

### Skill

```text
A = no skill
B = skill
```

or:

```text
A = Skill A
B = Skill B
```

### Setup

```text
A = original setup
B = modified setup
```

### Harness

```text
A = Harness A
B = Harness B
```

Keep everything else identical whenever possible.

If an unavoidable difference exists, record it explicitly.

---

# 17. Verify the Environments Before Starting

Immediately before sending the task to either agent, compute an environment fingerprint for each option.

Verify equivalence of:

```text
Docker image digest
OS/runtime versions
compiler versions
package-manager versions
dependency versions
tool versions
CPU allocation
memory allocation
network policy
environment variables
filesystem permissions
working directory
repository tree hash
mutation identity
tests
validation commands
task text
```

The expected result is:

```text
A == B
+
only intended experimental difference
```

If an unexpected difference exists, **do not start the agents**.

Fix the environments first.

This check is mandatory.

---

# 18. Give Both Agents the Same Task

Use the same task text for A and B.

Do not expose:

* mutation patch;
* deleted original code;
* original document content;
* bug locations;
* number of bugs;
* mutation generator;
* ground truth;
* other option's results;
* other option's trace.

### Coding

> Some application logic in this repository is missing. Reconstruct the missing implementation using the surrounding code, interfaces, types, documentation, and tests as evidence.
>
> Do not modify tests, CI/CD, build infrastructure, dependency configuration, pre-commit configuration, evaluation infrastructure, or unrelated code.

### Spec / Architecture / Design

> Some technical specification, architecture, or design documentation in this repository is incomplete. Reconstruct the missing section using the surrounding documentation and repository as evidence.
>
> Do not modify application code, tests, CI/CD, build infrastructure, dependency configuration, pre-commit configuration, evaluation infrastructure, or unrelated files.

### Code Review

> Review this repository carefully and find the introduced application bugs. Fix them while preserving the intended behavior of the existing system.
>
> Do not modify tests, CI/CD, build infrastructure, dependency configuration, pre-commit configuration, evaluation infrastructure, or unrelated code.

### Doc Review

> Review the specified document against the repository and identify concrete correctness, consistency, and completeness issues. Fix the issues you find using the implementation and other repository documentation as evidence.
>
> Do not modify unrelated documents, application code, tests, CI/CD, build infrastructure, dependency configuration, pre-commit configuration, evaluation infrastructure, or unrelated files.

---

# 19. Start Fresh Agent Sessions

Each option gets:

* its own Docker container;
* its own repository;
* fresh agent state;
* fresh session;
* isolated HOME/config/cache.

Do not reuse sessions.

Do not allow the agents to communicate.

---

# 20. Measure Agent Task Time

**Agent time means only the time the agent spends executing the task.**

Start the timer immediately before giving the task to the agent.

Stop the timer when the agent completes the task/session.

Do not include:

* Docker startup;
* repository creation;
* dependency installation;
* harness installation;
* mutation generation;
* preflight checks;
* post-agent validation.

Record:

```text
agent_started_at
agent_finished_at
agent_duration_seconds
```

This is the task-time metric used in the comparison.

Record setup and validation time separately if useful.

---

# 21. Record Token Usage

After each agent finishes, retrieve actual usage.

Record:

```text
model
input tokens
cached input tokens
cache-write tokens
output tokens
reasoning tokens
other billable categories
```

Record usage per model request when the harness may use:

* planner models;
* coding models;
* review models;
* subagent models;
* fallback models.

Record the effective model actually used.

---

# 22. Calculate Model Inference Cost

For every model actually used:

1. identify the exact model;
2. retrieve the current applicable pricing;
3. record the authoritative pricing source;
4. record the pricing retrieval time;
5. normalize the provider's actual billable usage categories;
6. calculate estimated model inference cost.

Do not assume cached tokens are excluded from `input_tokens`.

Do not double-count tokens.

If the provider reports an actual billable amount, prefer it.

Otherwise calculate:

```text
estimated_model_inference_cost
=
actual billable token usage
×
applicable current token prices
```

Record:

```text
pricing_source
pricing_retrieved_at
model
pricing_tier
currency
token prices
billable categories
estimated_model_inference_cost
```

If current pricing cannot be verified:

```text
cost_status = UNKNOWN
```

Do not invent a price.

---

# 23. Run Identical Validation

After each agent finishes:

1. inspect changed files;
2. check protected paths;
3. check dependency files;
4. restore protected files if necessary;
5. run the exact same validation commands for both options.

Run:

```text
build
tests
repository-specific validation
```

Also run when applicable:

```text
lint
typecheck
integration tests
documentation validation
schema validation
```

Record:

```text
validation_started_at
validation_finished_at
validation_duration_seconds
validation_output
validation_exit_status
```

Keep validation time separate from agent time.

---

# 24. Check Protected Files and Dependencies

An agent must not redefine the evaluation.

After execution, check for changes to:

```text
tests
CI/CD
build configuration
dependency manifests
lockfiles
evaluation scripts
evaluation configuration
ground-truth artifacts
```

If any were modified:

1. record the violation;
2. restore them from the frozen mutation;
3. exclude those changes from the solution;
4. rerun validation against the protected baseline;
5. report the violation.

---

# 25. Evaluate Correctness

### Coding

Check whether the required tests pass after reconstruction.

### Code Review

Check whether the injected bugs were repaired and the original tests pass.

### Spec / Architecture / Design

Compare the reconstructed region with the original deleted region.

### Doc Review

Compare the repaired specified document with the original document.

For documentation evaluations, preserve the original as hidden ground truth.

Do not expose it to the agents.

Do not replace ground-truth comparison with an LLM's subjective quality judgment.

---

# 26. Handle Invalid Runs Explicitly

Use explicit statuses:

```text
PASS
FAIL
TIMEOUT
HARNESS_ERROR
INFRA_ERROR
VALIDATION_ERROR
PROTECTED_FILE_VIOLATION
INVALID_MUTATION
INVALID_ENVIRONMENT
```

Examples:

```text
mutation does not create the intended missing behavior
→ INVALID_MUTATION

A and B have different dependency versions
→ INVALID_ENVIRONMENT

Docker fails before the agent starts
→ INFRA_ERROR

agent exceeds timeout
→ TIMEOUT

agent completes but required tests fail
→ FAIL
```

Do not interpret infrastructure failures as task-performance results.

---

# 27. Raw Evaluation Artifacts

**Raw evaluation artifacts are the original evidence emitted by the evaluation, not a human-written summary.**

Preserve them before interpreting or aggregating the results.

For each option, collect the raw:

```text
agent transcript
tool-call trace
agent stdout/stderr
harness stdout/stderr
model/provider usage response
timing data
final repository diff
changed-file list
test output
build output
validation output
container metadata
environment fingerprint
exit status
```

For each mutation, preserve:

```text
mutation metadata
mutation seed
target region
ground-truth content
mutated tree hash
mutation patch
baseline validation output
mutation validation output
```

For documentation evaluations, preserve the original ground-truth document/region separately from the agent-visible repository.

Raw artifacts must be:

* collected from the actual execution;
* stored without rewriting their contents;
* associated with the exact trial and option;
* timestamped where applicable;
* sufficient to reconstruct what happened.

Do not replace raw traces with only a summarized report.

Do not expose raw ground-truth artifacts to evaluated agents.

Do not expose Option A artifacts to Option B or vice versa.

---

# 28. Store Runtime Artifacts Outside the Repository

During evaluation, store runtime artifacts under:

```text
/tmp/harness-eval/<run-id>/
```

Use:

```text
/tmp/harness-eval/<run-id>/
├── manifest.json
├── environment.json
├── mutation/
│   ├── mutation.json
│   ├── ground-truth.patch
│   └── validation/
├── option-a/
│   ├── repo/
│   ├── container.json
│   ├── trace/
│   ├── stdout.log
│   ├── stderr.log
│   ├── diff.patch
│   ├── usage.json
│   └── validation/
└── option-b/
    ├── repo/
    ├── container.json
    ├── trace/
    ├── stdout.log
    ├── stderr.log
    ├── diff.patch
    ├── usage.json
    └── validation/
```

Keep ground-truth artifacts separate from agent-visible repositories.

---

# 29. Destroy the Evaluation Environment

After all raw artifacts have been collected:

1. stop Container A;
2. stop Container B;
3. remove both containers;
4. remove writable volumes;
5. remove temporary repositories;
6. remove temporary runtime state.

Do not leave agent state or evaluation state behind.

---

# 30. Persist Results

Only after the evaluation is complete and the containers are destroyed, copy the required artifacts into:

```text
.agent/harness-evals/<run-id>/
```

Do not expose this directory to the evaluated agents.

When taking a future evaluation snapshot, always exclude:

```text
.agent/harness-evals/
```

so historical results cannot leak into future experiments.

---

# 31. Multiple Trials

When `trials > 1`, repeat the full process.

Each trial gets:

```text
fresh mutation
fresh Docker container A
fresh Docker container B
fresh repository A
fresh repository B
fresh agent session A
fresh agent session B
```

Use a new mutation seed per trial.

Do not reuse agent state.

Record every trial independently.

---

# 32. Final Results

For every option report:

| Metric                    |   Option A |   Option B |
| ------------------------- | ---------: | ---------: |
| Task result               |  PASS/FAIL |  PASS/FAIL |
| Tests                     |     result |     result |
| Agent time                |   duration |   duration |
| Input tokens              |      count |      count |
| Cached input              |      count |      count |
| Output tokens             |      count |      count |
| Model                     | identifier | identifier |
| Model inference cost      |     amount |     amount |
| Protected-file violations |      count |      count |
| Dependency changes        |      count |      count |
| Exit status               |       code |       code |

Also provide the locations of the raw artifacts.

For multiple trials, report each trial separately.

Do not collapse the experiment into a subjective ranking or winner.

---

# 33. Final Invariants

Before declaring an evaluation complete, verify all of these:

```text
[ ] main checkout was never used as an execution environment
[ ] both options ran in separate Docker containers
[ ] each option had its own repository
[ ] both repositories started from the exact same frozen mutation
[ ] the mutation was generated only once per trial
[ ] A and B had identical environments before the intended difference
[ ] dependencies were identical
[ ] tests were identical
[ ] validation was identical
[ ] the task text was identical
[ ] ground truth was inaccessible to both agents
[ ] coding mutations only removed test-covered application code
[ ] no protected evaluation files were used to redefine correctness
[ ] agent time was measured separately from setup/validation
[ ] actual token usage was recorded
[ ] cost used the actual model and applicable current pricing
[ ] raw artifacts were preserved
[ ] containers and temporary environments were destroyed
[ ] results were persisted under .agent/harness-evals/
[ ] main checkout remains unchanged
```

The purpose of this skill is to turn harness development into a controlled experiment:

```text
harness change
      ↓
isolated A/B experiment
      ↓
same task + same mutation + same environment
      ↓
independent agents
      ↓
correctness + agent time + token usage + cost
      ↓
raw evidence + final report
```

Do not decide the outcome for the user.

