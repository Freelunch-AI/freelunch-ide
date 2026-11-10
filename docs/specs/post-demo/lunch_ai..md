# [dev + ops mode] [declarative] [self-contained] Project Management & AI-assisted Engineering Front — lunch-ai-coding: for project management & ai-assisted engineering

AI-assisted SDLC, with github/jira/linear/etc integrations, best practices enforced and team-based development built-in. 

Features:
- kanban for your agents
- automatic type labelling of Issues & PRs and tracking sucess rate of agents on theses types. When sucess rate is above threshold, requires no more human approval (e.g., issue approval or pr review)
- assign issues to local and remote agents (when local machine is not enough for running all the coding agents you want to run, because 1 sandbox/agent starts becoming heavy)
- get notified on important remote agent events
- agent cost manageemnt
- self-hosted, open source-based LLM API for agents to use instead of external LLM providers
- your custom coding model always beating major LLM APIs: proactive monitoring of where new LLM APIs beat the custom model, then finetuning the cusotm model on distilled data from the best performing API for the specific lagging place. This is done until custom model is again the best model.
- first mate: have a manager agent that you talk to, the manager deals with all the other agents. Manager can explain, make changes, monitor and provide help to the other agents.
- built-in autonomous cloud agents (with easy evals creation for them + improver agent that tries to improve them continuosly on evals) for: app spec-code-docs disalignment flagging & suggestions, app feature request creation from user feedback, app tech debt removal, app performance enhancement, app pentesting, app pr review, platform pr review, app staging testing, app deploy babysitting, app operations (problem identification, root cause analysis, issue creation), platform operations (problem identification, root cause analysis, issue creation); that ping human engineer when stuck in a loop or uncertain about an important thing.
- automatic pr merge policy that can merge prs automatically depending on the type of pr (ai-made pr acceptance rate is tracked per type of pr, if its very high, might be good option to just auto merge taht type of pr)

## On self-improvement of agents and automatic discovery of the best company/repo-specific coding setup for human developers

Coding agents differ in how well they perform depending on the codebase, task, and engineering environment. The choice of agent, model, tools, instructions, and workflows can significantly affect the quality, cost, and speed of software development. However, most organizations select these configurations based on public benchmarks, vendor claims, or manual experimentation, which rarely reflect the organization's actual repository, architecture, conventions, and development workflows. Our approach is to automatically evaluate and optimize coding agent configurations against the organization's own codebase, using real engineering tasks and measurable results rather than relying on generic benchmarks.

### Generating candidate configurations

The platform generates candidate setups by combining different coding harnesses (Claude Code, Codex, etc.), models, global rules, Skills, plugins, tools, and development workflows. Each candidate represents a complete, reproducible configuration of the coding environment. The search process is configurable: organizations can define an experimentation budget, the number of candidates to evaluate, and prior knowledge that constrains or guides the search. For example, they can exclude agents that have consistently performed poorly, prioritize configurations that previously produced good results, or require certain tools and coding rules to be present. This allows the platform to build on previous experiments rather than repeatedly exploring configurations that are already known to be ineffective.

### Evaluating candidates on real repository tasks

Candidates are evaluated using issues from the organization's own repositories. Each candidate receives the same task specification and initial repository state, and operates independently in an isolated environment, ensuring that candidates cannot access each other's changes or intermediate results. Each candidate must implement the requested changes and submit a pull request. The resulting implementations are evaluated against the same criteria:

- **Correctness:** Does the implementation satisfy the issue specification? Do the relevant tests pass? Are there regressions?
- **Code quality:** Is the implementation maintainable, understandable, modular, and consistent with the repository's architecture and conventions?
- **Performance:** Does the resulting implementation meet relevant runtime performance requirements, such as latency and memory consumption? Where applicable, this is measured through profiling and representative deployment scenarios.
- **Cost:** What was the total cost of producing the implementation, including model tokens, inference charges, and compute resources?
- **Development time:** How long did the agent take to produce the result, including tool execution, testing, pre-commit hooks, and other workflow steps?

Evaluation combines automated checks, such as tests, static analysis, and profiling, with structured code review. Failed mandatory checks can disqualify a candidate, while the remaining dimensions are scored according to the organization's priorities. Where results are variable, candidates can be evaluated across multiple tasks and repeated trials to distinguish consistent improvements from isolated successes.

### Selecting and deploying the best configuration

The platform compares candidates using the evaluation results and the organization's scoring criteria. The objective is not necessarily to find the configuration that produces the highest-quality code regardless of cost or time. Instead, organizations can define their own trade-offs, such as prioritizing correctness and maintainability while imposing cost and execution-time constraints. The best-performing configuration is saved as a reproducible engineering environment, including its agent, model, tools, instructions, and workflow configuration. It can then be deployed as an autonomous coding agent that continuously picks up issues from the organization's GitHub repositories, implements changes, runs the required checks, and opens pull requests using the validated configuration.

### Continuous evaluation and promotion

The selected configuration becomes the production baseline, but it is not assumed to remain optimal indefinitely. New models, coding agents, tools, and workflow strategies are continuously evaluated against the incumbent configuration using comparable repository tasks and evaluation criteria. A challenger is promoted only when the results demonstrate a meaningful improvement according to the organization's priorities, without violating mandatory quality or reliability requirements. This creates a continuous optimization process in which the platform generates candidate configurations, evaluates them on repository-specific tasks, compares their results against the production configuration, and promotes those that demonstrate measurable improvements as new models and tools become available. The platform also retains historical evaluation results, allowing organizations to track how configurations perform over time and understand which changes contributed to improvements or regressions.

### Learning from evaluation data

Over time, the platform accumulates a dataset connecting engineering environment configurations, repository characteristics, task specifications, implementation results, code reviews, and execution costs. This data can be used to improve the configuration search process itself, helping the platform predict which combinations of agents, models, tools, and workflows are worth evaluating for a particular repository. Our longer-term research direction is to use this data to fine-tune open-source foundation models for software engineering. Rather than relying exclusively on external prompts, Skills, and repository instructions to communicate engineering practices to an agent, specialized models could learn some of these practices directly during training. This could reduce the amount of context required at inference time and potentially improve latency and inference costs, although whether it produces measurable improvements over context-based approaches would itself need to be evaluated.

The result is a system that does more than execute software engineering tasks: it systematically measures, selects, deploys, and improves the engineering environments responsible for executing them.
