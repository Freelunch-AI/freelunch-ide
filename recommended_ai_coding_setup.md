# AI-assisted Coding Recommended Setup

## Suggested Feature Building Flow (tool-agnostic) (assuming repo foundation (group 1) is already layed out)

- each unique step is a slash command
- each step allows single-agent or multi-agent collab. If multi-agent: them manin agent manages worker in their separate sub-branch which has its own git worktree. One terminal for each. But you only approve the final commit presented by the main agent (main agent reviews/approves sub-agent sub-branch under the hood)
- steps can (and probably should) use specilist sub-agents at the type of task (e.g., specislist security agent for security review)
- can go back from a step to a previous step if necessary to fix issue created earlier (but step progress needs to be tracked in an feature_flow.md file)
- after every approval, a git commit is made
- multiple features can be implemented in parallel by having separate worktrees for each working branch. One terminal for each.
- approval gates mean that either the user (developer) or a specific AI agent needs to give approval to continue the flow

1. **/Start Feature Building: point to github issue, agent will read the issue and create feature branch with appropriate name according to the branching strategy file**
3. **/Ask User Clarifying Questions & do web search if necessary** [User Approval Gate]
4. **/Build issue-specific Spec (PRD + Architecture + Tech Stack)** & Review against Global Spec (PRD + Architecture + Tech Stack) catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled [User Approval Gate] 
5. **/Define Allowed boilerplate dependencies** (e.g., programming language, build tool, testing tools, package manager, etc) [User Approval Gate]
6. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/Setup/Modify the code stucture/skeleton** (directories, files, functions, classes, types, docstrings, test build command, final packaging build command), install boilerplate depedencies & Review against Issue-specific & Global Spec (PRD + Architecture + Tech Stack) catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled [User Approval Gate] 
7. [Make plan first & keep updating the plan at every step]  [User Approval Gate] **/Write/Modify the functional tests** (unit tests, integration tests) & Review against Issue-specific & Global Spec [User and AI QA Approval Gate] 
8. **/Test the functional tests with placeholder feature code**
9. **/Define Allowed feature code dependecies** [User Approval Gate]
10. [Make plan first & keep updating the plan at every step]  [User or Independent AI Reviewer Approval Gate]**/Write feature code using only the allowed feature code dependencies & Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled** [User or Independent AI Reviewer Approval Gate]
11. [Make plan first & keep updating the plan at every step]  [User or Independent AI Reviewer Approval Gate] **/Fix Linting errors, Static Analysis & Simplify Code** [User or Independent AI Reviewer Approval Gate]
12. [Make plan first & keep updating the plan at every step] [User Approval Gate] **/Build and Test feature code with the functional tests, generate testing & test coverage reports** & Review against Issue-specific & Global Spec [User or Independent AI Reviewer Approval Gate], repeat this step until all tests pass 
13. [Make plan first & keep updating the plan at every step] **/Independent Code Review** (including Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled)
14.  [Make plan first & keep updating the plan at every step] [User or Independent AI Reviewer Approval Gate]**/Make necessary code/test changes, build & test** & Review against Issue-specific & Global Spec [User Approval Gate]
15. [Make plan first & keep updating the plan at every step] **/Grill User with questions to see if he really understands the commit, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]
16. **/Document**:  Final User Documentation (how to install & use the product) & Controbutor Documentation (how to understanding the codebase), both in the form of step by step tutorial. [User or Independent AI Reviewer Approval Gate]
17. **/Specialized Security Review** flagging critical problems & warnings [AI Security Reviewer Approval Gate]
18. **/Push & Open PR**
19. [On PR Rreview Notification] **Read PR Reviews from Github**
20. [Make plan first & keep updating the plan at every step] [User or Independent AI Reviewer Approval Gate] **Make necessary code/test/docs changes, build & test** & Review against Issue-specific & Global Spec  [User Approval Gate]
21. [Make plan first & keep updating the plan at every step] **/Independent Code Review** (including Review against Spec catching inconsistencies with spec, things not specified in spec and problems in spec that needed to be overruled)
22. [Make plan first & keep updating the plan at every step] [User or Independent AI Reviewer Approval Gate] **/Make necessary code/test/docs changes, build & test.** & Review against Issue-specific & Global Spec [User Approval Gate]
23. [Make plan first & keep updating the plan at every step] **/Grill User with questions to see if he really understands the commit, user review code and asks questions until he has full understanding** [AI Griller Approval Gate]
24. **/Document**:  Final User Documentation (how to install & use the product) & Contributor Documentation (how to understanding the codebase), both in the form of step by step tutorial. [User or Independent AI Reviewer Approval Gate]
25. **/Specialized Security Review** flagging critical problems & warnings [AI Security Reviewer Approval Gate]
26. **/Push to PR with Summary of Changes, PR has to link the Issue it solves. Never merge automatically.**

## Tools Suggestion
- Agent-native Editor/Terminal: **Superset**
- Terminal Agent Harness: Use **OpenCode Go** Subscription.
- IDE (for better introspection + manual editing): **VSCode**
- Model: current best oss coding model (because of the OpenCode Go subscription)
- Feature Flow: Start with just making each step a slash command, and leave to the developer to follow the steps (note: this doesnt enforce step execution, se requires our commitment). Note: each slash commands should remember to update the feature building progress at the end (feature_flow.md file) or create the file if its still not created). Each slash command should have a simple name & also use the sub-agents that are most appropriate for the step.
- OpenCode Plugins (package skills + slash commands + sub-agents ... until a unit): **signoz plugin, opencost plugin, headlamp plugin, graphify, rtk, security-guidance, code-review, code-simplifier, explanatory-output-style, skill-creator and summarize-session.** 
- OpenCode MCPs: Github MCP.
- Custom Freelunch OpenCode Sub-agents: **spec-specialist, security-specialist, refactoring-specialist, cloud-architect, platform-engineer, testing-specialist, debugging-specialist, grilling-specialist, proxmox-specislist, talos-linux-specialist, linux-specialist, k8s-specialist, sre, cicd-specialist, golang-specialist, typescript-specialist.**
- Custom Freelunch OpenCode Slash Commands: one for each unique step of the feature building flow
- Dependency Docs: a `dependency_docs` folder holding the skill or documentation (if no skill is available) of pinned version of every first-level repo dependency (by default) with two sub-folders: `system_dependencies` (e.g., proxmox, kubetcl, etc) and `libraries` (e.g., go libraries, node libraries). Can also add lower level repo depdendencies when debugging (adding a docs of a dependency of a dependency to help debugging)
- Custom Skills: created on-demand, under a `custom_skills` folder, via `skill-creator` to avoid having to repeat the same solution process over and over
