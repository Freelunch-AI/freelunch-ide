# [dev + ops mode] [declarative] [self-contained] Project Management & AI-assisted Engineering Front — lunch-vibes: for project management & ai-assisted engineering

AI-assisted Development and Project management, with github/jira/linear/etc integrations, best practices enforced and team-based development built-in. 

Relevant Benchmark to optimize agaist: Horizon-SWE.

Reference projects: vibe-kanban.

It’s vibe-kanban on steroids. Added features:
- automatic setup & autoscaling of self-hosted Coding LLM API with the best current coding model & model routing (local fast/cheap SLM, cloud ok LLM, cloud slow/constly gigantic LLM).
- helps with spec-generation (vision spec and PR spec), with available spec templates
- built-in cost management per agent and per LLM API, with side-by-side agent comparisons for the same task
- agent imporovement: easily create skills and assign them to agents, finetune (if using open model), evaluate, track experiments and promote agents. Agetns also self-improve by leveraging its interactions with you and its failures via prompt optimization, skill optimization and/or automatic finetuning.
- alocate max hardware resources of your machine that cna be used for running agents, the rest of the agents are run on your VPC.
- external context: agents can access knowledge from your project management tools (e.g., Notion, Jira, Github, Slack etc) and from monitoring/observability tools (e.g., Datadog, Prometheus, etc)
- structured development: software engineering-backed standard procedures for handling issues of these types: feature request, optimization and bug.
- Drift monitoring & resolution for issue-ticket-prspec-plan-tests-code-docs by identifying source of trith
- qa: Multi-model Quality assurance jobs & agents triggered by diffs, commits and PRs.
- agent cnfiguration: Granular configuration of what each agent can and cannot do, by defining types of agents and then instantiating them.
- PR battle: Have two seprate agents work on same issue, then after some time pick the best and borrow ideas from the other implementation
- PR iteration: If PR isnt accepted right away, agent reads the feedback and iterates on the PR until it gets accepted.
- suepervisor agent: supervisor agent that monitors worker agents to break them out of loops and reset to last stable commit when things go astray. Supervisor agent also can ask human engineer clarifying questions and pass them to the respective worker agent.
- create and assign github issues directly to agents that can use lunch-vibekill themselves.
- ai-assitance not just for spec making & coding, also:
  - PR Reviews
  - Frontend Testing
  - Infra & App Monitoring/Observability
  - Ticket creation
  - Issue creation
  - Tests & Evals additions based on problems detected

Development Journey:

AGENTS.md: organization wide, rules for guiding agents, independent of the project

GUARDRAILS.md: whitelist or blacklist of tools, proihibited commands p/tool p/agent, commands that always need approval, command that never need approval

Assignments:
- Each branch to a developer/engineer (human or ai) with one other engineer as his/her helper for giving support in debugging/discussions. Ideally the helper engineer is working on another PR that is closely related to what the main engineer is working on
- Also assign each feature optimization PR to a horizontal engineer. Note: PRs can also be also be assigned to an AI agent working with git worktrees with a team of sub-agents (each sub-agent working on a git worktree).
- Assign Main branch to a human Lead/Staff Engineer
- If branch B opens from branch A (assuming branch A is also a feature branch), then the develper/swe responsible for branch A will also be responsible for reviewing the PR branch b makes to branch at some point
- 
Within a PR branch: AI-assited Coding Guidance:
- Every building/debugging step is done by n diverse agents in parallel and a hypervisor agent chooses the best implementation to show the user (if the user doesnt like, shows other according to user feedback, etc). Leveraging git-worktree.
- Every project creation routine must follow this sequence of steps:
   - Braching Policy creation and testing. With AI-assitance (scaffolding + improvement)
  - Pre-commit Hooks creation and testing. With AI-assitance (scaffolding + improvement)
  - CI/CD Pipeline creation and testing. With AI-assitance (scaffolding + improvement)
- Every feature building/removal routine must follow this sequence of steps (With AI-assitance):
  - Brach creation and checkout
  - (Optional - only if new plugin is needed) IDE plugin install
  - Setup/Modify the code stucture/skeleton (directories, files, functions, classes, types, docstrings, test build command, final packaging build command)
  - Write/Modify the functional tests
    - unit tests
    - integration tests
    - end-to-end tests
  - Test the functional tests with dummy feature code
  - Allowed feature code dependecies
  - Write feature code using only the allowed feature code dependencies
  - Fix Linting errors
  - Statically Review changes made
  - Build and Test feature code with the functional tests
  - Define possibly affected regression tests
  - Build necessary code and run possibly affected regression tests
  - Move tests to regression tests
  - Remove legacy regression tests that dont make sense anymore
  - Update documentation
- Every feature (performance, security, compliance, scalability, refactoring or documentation) optimization routine must follow this sequence (With AI-assitance):
  - Brach creation and checkout
  - (Optional - only if new plugin is needed) IDE plugin install
  - (Optional - only if sturcture needs change) Setup/Modify the code stucture (directories, files, functions, classes, types, docstrings, test build command, final packaging build command)
  - Remove Regression tests that dont make sense anymore
  - (Optional - only if code structure was changed) Write the feature optimization tests
    - unit tests
    - integration tests
    - end-to-end tests
  - (Optional - only if tests were changed) Test the tests
  - Allowed feature code dependecies
  - Write optimized feature code using only the allowed feature code dependencies
  - Statically Review changes made
  - Build and Test optimized feature code with the functional and performance tests
  - (Optional - only if code structure was changed) Define possibly affected regression tests
  - (Optional - only if code structure was changed) Build necessary code and run possibly affected regression tests
  - Move tests to the respective horizontal optimization regression tests
  - (Optional - only if code structure was changed) Remove legacy regression tests that dont make sense anymore
  - Update documentation
- Every fix routine must follow this sequence
  - Make the problem reproducible and formalized
  - Ask for more information & do web search if necessary
  - Make a few hypotheses about the Root Problem
  - Test Hypotheses about the Root Problem & Identify the Root Problem
  - Make the fix
  - Create new tests if necessary
  - Run tests to ensure it was fixed
  - Write debugging post-mortem & store it in the post-portem registry

Block owners are repsonsible for monitoring their blocks in production

Creation of issues (problems identified during dev that dont require monitoring data to understand) or self-explainable (cards with reference to monitoring data and codebase, also with reference to replayable system for a window around that time), time-freezed (values referenced are the values that were showing at the time of the ticket, not current time) tickets (problems identified during ops) trigger as hook that finds the root cause, when the engineer aproves the root cause, it causes the plan to be rebuilt (changing it still requires approval from PD/Lead)

Each Spec is a different system that is developed on it’s respective isolated workspace (each dev team works in their separate workspace) inside the platform. Workspaces can also me merged. The default workspace has all the implementation modes inside of them (unless configured otherwise), new workspaces only have the implementation modes needed (stated in the spec).

Spec, Plan, Tests, Code, Docs, Issues and Tickets are always monitored for drifts between them. If a drift is detected, engineer is presented with a dirft report explaining the differences, where the engineer must choose the actual source of thruth. If the engineer doesnt clarify within a specified time period, then the souce of truth defaults to being the spec, issues and tickets and the other ones will be fixed to sync with them.

When the AI Assistant finds a problem with the spec and/or a problem with the system during development, it creates an issue and present a problem report to the engineer explaining the problems and with reccomendded acitons paths. The engineers chooses a recommended action ppth or makes a custom one, which will cause the plan and/or spec to be fixed and development can continue. If the engineer doesnt chosse within a specified time period, then the AI Assistant automatically decides the best action path and makes the modifications itself.

lunch-vibes will not allow the Human engineer or the AI Assistant to bypass (will block the feature branch commit or PR) the right development structure:

- On Commit:
  - Agent that analyzes the commit, to see if its consistent with the current step
  - Not allowing feature code to be commited if the latest commit brakes regression tests or passes less feature tests than the previous one
  - Assuring development is done in the right branch
  - Runs stanadard pre-comit hooks (building, testing, security, complience, docstrings, etc)
  - Ensure developer understand the code (probably AI-generated code): synthetically generates a quizz with questions like “if X changes to X’, how does it affect Y” with answers derived from static or dynamic analysis. If developer score below trshold on the quizz, blocks the commit.
- On PR:
  - Agent analyzes the PR’s readme to see if the PR can be understood from it
  - Agent that analyzes end-to-end feature tests to see if its consistent current step and everthing is properly documented
  — Not AI-related: Standard CI jobs are run (building, testing, security, complience, docstrings, etc) —
  - Agent ensures developer understands the code (probably AI-generated code): synthetically generates a quizz with questions like “if X changes to X’, how does it affect Y”. Each questions spawns a chat thread with the AI reviewer. The AI reviewer will evaluated all quiz chat threads and set a final code-undertanding score for the developer. If developer score below trshold on the quizz, blocks the PR.
  - AI PR Reviewer helps Human PR reviewer

Voice & Highlight mode: interac twith the ai assistant by talking and highlighting code wile you talk. You say a standard your like “Fish” to signal you ended your monologue and the Assistant should start working on it.

AI specialization support:
- Vibes data mode: where you cna organize materials (or links to materials) for the AI assistant to leverage. Organization can happen per topic or per task the ai assistant needs to do. This data can be used both for training or at inference time (as context). Screen record tutorials you've made go here directy. Materials are categorized as knowledge or examples.
- Vibes eval mode: where you can build eval datasets and evaluate different versions of sub AI Assistants on them. Can also evaluate agetns side by side working on the same task, arena-style.
- Vibes experiment tracker mode: where you can track all versions of the AI Assistant you've "vibe finetuned", select specific versions to "vibe finetune" or to set as default
- Vibes finetuning mode: "vibe finetune" an AI assistant via iterative cycles of pointing it to "vibes data" (examples or knowledge), explaining the goal of the finetune in natural language, evaluating it on "vibes evals", and giving feedback
- Vibes config mode: config AI Assitant agents by: changing their type (Platform Engineer, Developer, AI Scientist, Data Engineer) changing their descriptions, tool set, confidence treshold, number of agents in each multi-agent swarm
- Vibes action scope mode: define which write-type actions each AI assistant agent can do autonomously (e.g., monitoring action), which ones require human engineer approval (e.g., PR to main) and which ones AI cant even touch (e.g., k8s low-level access). By default AI Assistant cannot do any write-type actions and can do all read-only type actions.
AI also supports getting context from external tools like Jira, Slack, DataDog, etc via MCP/skills.
