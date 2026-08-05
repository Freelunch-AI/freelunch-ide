# Bruno's AI-assisted Coding Setup

## Feature Building Flow (tool-agnostic) (assuming repo foundation (group 1) is already layed out)

- each unique step is a slash command
- each step is done by a single agent.
- each step has a planning sub-step performed at the beggining
- steps can (and probably should) use a specialist sub-agent for that type of task (e.g., specislist security agent for security review)
- can go back from a step to a previous step if necessary to fix issue created earlier (but step jumps need to be tracked in an feature_flow.md file which has the name and number of issue on its title)
- when a new feature starts, first need search for any completed feature_flow.md (inside .gitignore) and store it in completed_feature_flow folder (inside .gitignore) in the form feature_flow[i].md
- approval gates mean that either the user (developer) or a specific AI agent needs to give approval to continue the flow
- after every approval (human or ai), a git commit is made
- multiple features can be implemented in parallel by having separate terminals/worktrees for each working feature branch.

1. **/start Start Feature Building: point to github issue, agent will read the issue and create feature branch with appropriate name according to the branching strategy file**
3. **/clarify Ask User Clarifying Questions & do web search if necessary** [User Approval Gate]

Loop until 4 is succesfull
4. **/spec Build issue-specific Spec (PRD + Architecture + Tech Stack)** & Review against Global Spec (Founding Doc + Roadmap + Tech Stack) catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled [User Approval Gate]
5. **/specsecreview Specialized Spec Security Review** flagging critical problems & warnings [AI Security Reviewer Approval Gate]

6. [Make plan first & keep updating the plan at every step] **/fspecgrillme Understand the feature spec, then grill User with questions to see if he really understands the feature spec, user reviews feature spec and asks questions until he has full understanding** [AI Griller Approval Gate]

---- New Session ----

5. **/boilerdep Define Allowed boilerplate dependencies** (e.g., programming language, build tool, testing tools, package manager, etc) [User Approval Gate]
6. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/boiler Setup/Modify the code boilerplate (stucture/skeleton)** (directories, files, functions, classes, types, docstrings, test build command, final packaging build command), install boilerplate depedencies & Review against Issue-specific & Global Spec (PRD + Architecture + Tech Stack) catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled [User Approval Gate]
7. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands the commit, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]

---- New Session ----

Loop until 9 is sucessfull
8. [Make plan first & keep updating the plan at every step]  [User Approval Gate] **/writetests Write/Modify the functional tests** (unit tests, integration tests) & Review against Issue-specific & Global Spec [User and AI QA Approval Gate] 
9. **/testtests Test the functional tests with placeholder feature code**

10. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]

Loop until 14 is sucessfull
11. **/featdep Define Allowed feature code dependecies** [User Approval Gate]
12. [Make plan first & keep updating the plan at every step]  [User or Independent AI Reviewer Approval Gate]**/feat Write feature code using only the allowed feature code dependencies & Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled** [User or Independent AI Reviewer Approval Gate]
13. [Make plan first & keep updating the plan at every step]  [User or Independent AI Reviewer Approval Gate] **/fixstatic Fix Linting errors, Static Analysis & Simplify Code** [User or Independent AI Reviewer Approval Gate]
14. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/test Build and Test feature code with the functional tests, generate testing & test coverage reports** & Review against Issue-specific & Global Spec [User or Independent AI Reviewer Approval Gate], repeat this step until all tests pass 

15. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands changes since last grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]

---- New Session ----

Loop until 17 is sucessfull or go back to a previous step
16. [Make plan first & keep updating the plan at every step] **/review Independent Code Review** (including Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled)
17.  [Make plan first & keep updating the plan at every step] [User or Independent AI Reviewer Approval Gate]**/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec [User Approval Gate]

18. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands the changes since last grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]
19. **/document Document**:  Final User Documentation (how to install & use the product) & Controbutor Documentation (how to understanding the codebase), both in the form of step by step tutorial. [User or Independent AI Reviewer Approval Gate]

---- New Session ----

Loop until 17 is sucessfull or go back to a previous step
17. **/secreview Specialized Security Review** flagging critical problems & warnings [AI Security Reviewer Approval Gate]
18. [Make plan first & keep updating the plan at every step] [User or Independent AI Reviewer Approval Gate]**/redo Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec [User Approval Gate]

19. [Make plan first & keep updating the plan at every step] **/grillme Understand the codebase, then grill User with questions to see if he really understands the changes since last grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]

20. **/pr Push & Open PR with Summary of Changes, PR has to link the Issue it solves. Never merge automatically.**

---- New Session ----

Loop until 19 is sucessfull or go back to a previous step
{
19. [On PR Review Notification] **/prreviews Read PR Reviews from Github and write them locally on a dedicated folder**
20. [Make plan first & keep updating the plan at every step] [User or Independent AI Reviewer Approval Gate] **/redo Make necessary code/test/docs changes, build & test** & Review against Issue-specific & Global Spec  [User Approval Gate]

Loop until 21 is sucessfull or go back to a previous step
21. [Make plan first & keep updating the plan at every step] **/review Independent Code Review** (including Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled)
22. [Make plan first & keep updating the plan at every step] [User or Independent AI Reviewer Approval Gate] **/redo Make necessary code/test/docs changes, build & test.** & Review against Issue-specific & Global Spec [User Approval Gate]

23. [Make plan first & keep updating the plan at every step] **/grillme Grill User with questions to see if he really understands changes since lst grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]
24. **/document Document**:  Final User Documentation (how to install & use the product) & Contributor Documentation (how to understanding the codebase), both in the form of step by step tutorial. [User or Independent AI Reviewer Approval Gate]

---- New Session ----

Loop until 25 is succesfull
25. **/secreview Specialized Independent Security Review** flagging critical problems & warnings [AI Security Reviewer Approval Gate]
26. [Make plan first & keep updating the plan at every step] [User or Independent AI Reviewer Approval Gate] **/redo Make necessary code/test/docs changes, build & test.** & Review against Issue-specific & Global Spec [User Approval Gate]
27. [Make plan first & keep updating the plan at every step] **/grillme Grill User with questions to see if he really understands changes since lst grillme, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]

---- New Session ----

28. **/pr Push & Open PR with Summary of Changes, PR has to link the Issue it solves. Never merge automatically.**
}

## Tech Stack
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
