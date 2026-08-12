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

## Another major key difference: calibrated & fine-grained confidence-backed reviews

Today’s coding agents can find increasingly subtle bugs, but they still struggle to know how much to trust their own conclusions. A review agent may identify a real race condition and a speculative performance issue with the same confident tone. The goal is to train a model with a genuine, calibrated sense of uncertainty at the claim level, so it can distinguish what it knows from what it is merely hypothesizing.

The core training approach would use reinforcement learning with confidence-sensitive rewards. Each review finding is evaluated against ground truth, and the reward is tied to both correctness and confidence: the model is rewarded more for being highly confident when it is right, and punished more heavily for being highly confident when it is wrong. In other words, confidently wrong findings should be much more costly than cautious mistakes, while confidently correct findings receive the strongest reward. This directly trains the model to develop a useful internal model of its own reliability rather than simply learning to sound uncertain.

There are two promising ways to expose that learned signal. Token-based confidence would train the model to delimit spans with special tokens such as <VERY_CONFIDENT>...</VERY_CONFIDENT>, <CONFIDENT>...</CONFIDENT>, and <NOT_CONFIDENT>...</NOT_CONFIDENT>, allowing different parts of a review to carry different confidence levels. A separate confidence head would instead keep confidence outside the generated text, predicting something like P(claim is correct) for each individual finding and letting the review system convert that signal into whatever interface or behavior is appropriate.

The result is a code-review agent that doesn't just try to find more bugs—it learns to know which of its own findings deserve to be trusted. Very high-confidence findings can automatically block a PR, medium-confidence findings can be surfaced to the developer, and low-confidence findings can trigger further investigation instead of being presented as facts.

## Recall is the north-star metric

Traditional review systems have to balance precision and recall. We intentionally move the tradeoff toward recall because the costs are asymmetric. A false positive costs an engineer a few minutes; a missed production bug can cost a company days, money, customers, or reputation.

The objective is to push automated defect recall high enough that reviewing every change manually is no longer worthwhile.

**Our Verification Layer** acts as a hard filter. High-recall candidate generation is fine internally within the swarm, but the swarm should only surface a warning to the human if it can produce an **execution proof** (e.g., a failing test case, a reproducible trace, or an explicit AST path breaking a contract). If the swarm *thinks* there is a bug but cannot prove it, it should automatically construct a runtime harness to test it before pinging the human.

Whata about complex cloud problems? The environment will give the code review agent access to local cloud-emulation tools such as Proxmox, LocalStack, Act, and Testcontainers, allowing it to reproduce complex cloud scenarios locally and provide failing-test proofs even for problems that would otherwise be difficult to reproduce.

## Architecture

To make this execution-emulation layer viable in production, the swarm must use a Tiered Escalation Cascade rather than jumping straight to heavy infrastructure simulation:

1. Level 1 (Static AST Analysis): Candidate agents generate hypotheses using static code/diff analysis.

2. Level 2 (In-Memory Verification): If a candidate bug can be proven with a standard mock or unit test, do it here (fast, low compute).

3. Level 3 (Containerized Verification): If the bug requires stateful dependencies (e.g., database, cache, message bus), invoke Testcontainers and/or Floci to run a localized integration proof.

4. Level 4 (Full Topology Simulation): Only for cross-service, workflow, or infrastructure bugs, spin up Act and/or lightweight cluster topologies (e.g., via Kind, Proxmox, OpenStack or Floci) as a final verification step.

## Humans become the escalation layer

The ideal workflow is simple: the AI reviews everything, and humans investigate the exceptions. A normal change can be investigated automatically and proceed when no credible defect is found. When the swarm discovers a potential problem, it produces the evidence and calls a human to investigate.

This changes the role of the engineer from **reviewer of every change** to **investigator of the changes that actually need human judgment**. Instead of spending hours reviewing code that is probably correct, engineers spend their time on the small fraction of changes where the system has found credible evidence of a problem.

## Review during development, not just PRs

The same system can run during development via cli or triggered by your main coding agent harness (e.g., claude code), to catch problems ebefore PR phase. The earlier a problem is encounterd, the cheaper it is to fix it.

## Training Pipeline

1. train each model indidually via SFT on open bug/vulnerability finding datasets + synthehtic mutation testing dataset (note: this dataset is constructed based on know real-world bug/vulnerability modes, not randomly) specific to that models role (e.g., maintainability, security, performance, scalability)
2. train each model indidually via RL (described in step 3) with calibration punishment.
3. train the entire system together (models, prompts, hyperparameters) via RL with calibration punishment where:
    1. its rewarded a little if:
        - giving high confidence on successes
    2. its rewarded a lot if:
        - for bugs/vulnerabilities it can: produce a test that passes in main but not on the mutated satellite branch 
        - for code quality: rewarded via llm-as-a-judge proportionally to how clean the test-passing diff it produces to a changes in product spec (quality code should produce cleaner diffs)
        3. giving low confidence on mistakes (bad test)
    3. punished a little for:
        - giving low confidence on successes (good test)
        - each step it takes
    4. punished more for:
        - producing a mistake
    5. punsihed a lot for:
        - giving high confidence on a mistake
4. train the entire system together (models, prompts, hyperparameters) via RL on real PRs that solved Bug issues

## First Users

The product will win fastest in more mature repos with strong existing unit-test coverage, clear module boundaries, and high business risk (e.g., financial logic, smart contracts, core API services) where compute costs are easily justified by defect prevention.
