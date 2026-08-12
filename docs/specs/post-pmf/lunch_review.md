# The Claude Code for Code Review

In 2026, AI coding agents can produce software faster than engineering teams can reliably inspect it. **Code review is becoming the bottleneck.** Greptile, CodeRabbit, and similar AI PR reviewers have shown that LLMs can already automate meaningful parts of review, but they are fundamentally limited by using existing models as reviewers rather than training the review system itself toward a measurable objective. We want to build the next step: **an autonomous code-review agent swarm trained as a whole specifically to maximize the recall of real software defects.**

## A Claude Code built for review

The product is an autonomous engineering swarm that can operate throughout the development lifecycle—not just on pull requests. Like Claude Code, it can explore the repository, use tools, run tests, inspect dependencies, modify code to test hypotheses, and reason across the codebase. The difference is its objective: **Claude Code tries to make the code work; our swarm tries to prove that it has problems.**

This means the product isn't simply a GitHub bot that comments on pull requests. It can operate alongside developers and coding agents while they are building software, continuously investigating changes, challenging implementations, searching for regressions, and testing suspicious behavior. Pull-request review is simply one place where the same underlying review agent can be deployed.

## The core technology: a trainable review swarm

Rather than relying on one reviewer model, the system creates a swarm of specialized agents. An investigator might explore the codebase and identify suspicious behavior, a test agent might construct counterexamples and execute targeted tests, a security agent might investigate vulnerabilities, and an adversarial reviewer might try to disprove the findings of the other agents. A verifier then evaluates the evidence before the swarm reaches a final decision.

The important property is that these agents don't have to use the same model. Different models can specialize in different failure modes, and the swarm can combine them according to whatever configuration produces the best results.

## The key difference: we train the swarm

This is the fundamental distinction from today's AI code-review products. With conventional AI review, the workflow is essentially **PR → existing model → review**. We instead want **code → specialized review swarm → investigation → verification → outcome → reward → RL training → better swarm**.

We leverage swarm-aware multi-model RL infrastructure such as AgentJet to train the agents toward the outcome of the entire review process. The optimization target isn't "did the model write a convincing review?" It is **"did the swarm discover a real defect?"**

This allows us to optimize the entire system: which agents participate, which models they use, what tools they invoke, how deeply they investigate, how they debate, how they verify findings, when they escalate to humans, how results are combined, and how findings are filtered. The models are components. **The swarm is the thing being optimized.**

## Recall is the north-star metric

Traditional review systems have to balance precision and recall. We intentionally move the tradeoff toward recall because the costs are asymmetric. A false positive costs an engineer a few minutes; a missed production bug can cost a company days, money, customers, or reputation.

The objective is to push automated defect recall high enough that reviewing every change manually is no longer worthwhile.

**Our Verification Layer** acts as a hard filter. High-recall candidate generation is fine internally within the swarm, but the swarm should only surface a warning to the human if it can produce an **execution proof** (e.g., a failing test case, a reproducible trace, or an explicit AST path breaking a contract). If the swarm *thinks* there is a bug but cannot prove it, it should automatically construct a runtime harness to test it before pinging the human.

## Humans become the escalation layer

The ideal workflow is simple: the AI reviews everything, and humans investigate the exceptions. A normal change can be investigated automatically and proceed when no credible defect is found. When the swarm discovers a potential problem, it produces the evidence and calls a human to investigate.

This changes the role of the engineer from **reviewer of every change** to **investigator of the changes that actually need human judgment**. Instead of spending hours reviewing code that is probably correct, engineers spend their time on the small fraction of changes where the system has found credible evidence of a problem.

## Review during development, not just PRs

The same system can run continuously alongside developers and coding agents. A coding agent might implement a feature while the review swarm independently investigates the changes as they are produced. It can identify suspicious assumptions, construct tests, search for regressions, and challenge the implementation before the work ever reaches a pull request.

## Training Pipeline

1. train each model indidually via SFT on open bug/vulnerability finding datasets + synthehtic mutation testing dataset (note: this dataset is constructed based on know real-world bug/vulnerability modes, not randomly) specific to that models role (e.g., maintainability, security, performance, scalability)
2. train the entire system together (models, prompts, hyperparameters) via RL where:
    1. its rewarded a lot if:
        1. for bugs/vulnerabilities it can: produce a test that passes in main but not on the mutated satellite branch 
        2. for code quality: rewarded via llm-as-a-judge proportionally to how clean the test-passing diff it produces to a changes in product spec (quality code should produce cleaner diffs)
    2. punished a little for each step it takes
    3. punished more for producing a bad test.
3. train the entire system together (models, prompts, hyperparameters) via RL on real PRs that solved Bug issues
