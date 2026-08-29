# The Claude Code for Code Review

In 2026, AI coding agents can produce software faster than engineering teams can reliably inspect it. **Code review is becoming the bottleneck.** Greptile, CodeRabbit, and similar AI PR reviewers have shown that LLMs can already automate meaningful parts of review, but they are fundamentally limited by using existing models as reviewers rather than training the review system itself toward a measurable objective. We want to build the next step: **an autonomous code-review agent swarm trained as a whole specifically to maximize the recall of real software defects.**

## A Claude Code built for review

The product is an autonomous engineering swarm that can operate throughout the development lifecycle—not just on pull requests. Like Claude Code, it can explore the repository, use tools, run tests, inspect dependencies, modify code to test hypotheses, and reason across the codebase. The difference is its objective: **Claude Code tries to make the code work; our swarm tries to prove that it has problems.**

This means the product isn't simply a GitHub bot that comments on pull requests. It can operate alongside developers and coding agents while they are building software, continuously investigating changes, challenging implementations, searching for regressions, and testing suspicious behavior. Pull-request review is simply one place where the same underlying review agent can be deployed.

## The Technology we use

We use an agetn swarm to perform the review with a problem candidate generation phase and a problem filtering phase.

### The key difference: we train the swarm that was initialized as a set of indepedently-built generators and verififers

Each initial verifier and generator is built independelty by our sub-teams without any knowledge of the other teams, with the sole objective of hill climbing a benchmark.

We leverage swarm-aware multi-model optimization infrastructure such as AgentJet (Multi-agent RL) for inner loop model optimizaiton and DSPY (LLM-based program optimization) for outer loop system optmization to train the agents toward the outcome of the entire review process. The optimization target isn't "did the model write a convincing review?" It is **"did the swarm discover a real defect?"**. Wieghts, prompts and hyperarameters can be optimized.

This allows us to optimize the entire system:
- how each cr doubles down on its specific strenghts
- how findings are filtered
- how each verififers opnion is taken into account at final decision
- which generators/verifier participate and when
- which models they use
- what tools they invoke
- how deeply they investigate

**The swarm is the thing being optimized.**

### Another major key difference: calibrated & fine-grained confidence-backed reviews

Today’s coding agents can find increasingly subtle bugs, but they still struggle to know how much to trust their own conclusions. A review agent may identify a real race condition and a speculative performance issue with the same confident tone. The goal is to train a model with a genuine, calibrated sense of uncertainty at the claim level, so it can distinguish what it knows from what it is merely hypothesizing.

The core training approach would use reinforcement learning with confidence-sensitive rewards. Each review finding is evaluated against ground truth, and the reward is tied to both correctness and confidence: the model is rewarded more for being highly confident when it is right, and punished more heavily for being highly confident when it is wrong. In other words, confidently wrong findings should be much more costly than cautious mistakes, while confidently correct findings receive the strongest reward. This directly trains the model to develop a useful internal model of its own reliability rather than simply learning to sound uncertain.

There are two promising ways to expose that learned signal. Token-based confidence would train the model to delimit spans with special tokens such as <VERY_CONFIDENT>...</VERY_CONFIDENT>, <CONFIDENT>...</CONFIDENT>, and <NOT_CONFIDENT>...</NOT_CONFIDENT>, allowing different parts of a review to carry different confidence levels. A separate confidence head would instead keep confidence outside the generated text, predicting something like P(claim is correct) for each individual finding and letting the review system convert that signal into whatever interface or behavior is appropriate.

The result is a code-review agent that doesn't just try to find more bugs—it learns to know which of its own findings deserve to be trusted. High-confidence problem findings can automatically block a PR, medium-confidence findings can be surfaced to the developer, and low-confidence findings can trigger further investigation instead of being presented as facts.

## Two Offerings of Lunch Review: Ralph Version & Human Version

### lunch-review-human: Maximize Defect Recall

When a human engineer is the final decision-maker, the system should optimize primarily for extremely high defect recall. False positives are acceptable because the human can discard speculative findings, while missed defects remain invisible.

The swarm can therefore operate as a deliberately high-recall candidate generator: investigate aggressively, surface weak hypotheses, and trade precision for coverage. The objective is to maximize the number of real defects discovered while providing enough evidence and calibrated confidence for the engineer to distinguish strong findings from speculative ones.

### lunch-review-ralph: Minimize High-Confidence Regression-Inducing False Positives

When review findings are consumed by an autonomous Ralph loop, the objective changes fundamentally. A false positive is no longer merely an annoying review comment: the agent may actually modify the code in response to it.

The critical metric therefore becomes high-confidence regression-inducing false positives: findings that are incorrect and whose acceptance causes the implementation to become less compliant with the actual product specification or introduces a behavioral regression.

The Ralph configuration should therefore optimize for very low high-confidence regression-inducing false-positive rate, while maintining high recall.

Non-regression-inducing false positive examples: unnecessary modularization, stylistic changes, overengineering security, overengineering performance, or other improvements that leave behavior correct. These induce non-critical time & cost problems. 

Regression-inducing false positive examples: findings that cause tests to be rewritten around incorrect behavior or findngs that change build/test/package/deploy or publish commands in a way that breaks intended behaviour & validation.

## Architecture

1. [Read Mode] [Multi-agent] [Test-time scaling] Static Codebase Understanding, Making hypotheses and predicitons for hypothesis, also sending questions to human engineer
2. [Read Mode] [Multi-agent] Runs experiments to test hypothesis and improve understanding
3. [Read Mode] [Multi-agent] [High-recall] Static Problem Candidate Generation with causal evidence, specification/contract grounding and problem execution proof (Failing Test) building type per problem (Level 1, Level 2, Level 3 or Level 4)
4. [Execution Mode] [Multi-agent] Independent Build/Test/Package/Deploy or Publish Command Runner, Problem Execution Proof (Failing Test) Builder and Problem Report Builder
    1. Level 1 (Static AST Analysis): Candidate agents generate hypotheses using static code/diff analysis.
    2. Level 2 (In-Memory Verification): If a candidate bug can be proven with a standard mock or unit test, do it here (fast, low compute).
    3. Level 3 (Containerized Verification): If the bug requires stateful dependencies (e.g., database, cache, message bus), invoke Testcontainers and/or Floci to run a localized integration proof.
    4. Level 4 (Full Topology Simulation): Only for cross-service, workflow, or infrastructure bugs, spin up Act and/or lightweight cluster topologies (e.g., via Kind, Proxmox, OpenStack or Floci) as a final verification step.
5. [Execution Mode] [Multi-agent] [Test-time scaling] [Confidence Calibrated Statements] [High-precision] [Very Low regression-inducing FP rate] Verifier that Review the Report to Produce the Final Problem Report

That produces a report in the form:

### Functional Problems
- Critical Problems: {"problems_a": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "spec_grounding", str, "confidence": float}; ...}
- Important Problems: {"problems_b": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "spec_grounding", str, "confidence": float}; ...}
- Enhancement Opportunities: {"problems_g": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "spec_grounding", str, "confidence": float}; ...}
- Non-critical score: float

### Security Problems
- Critical Problems: {"problems_c": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Important Problems: {"problems_d": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Enhancement Opportunities: {"problems_h": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Non-critical score: float

### Code Quality & Documentation Problems
- Critical Problems: {"problems_e": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Important Problems: {"problems_f": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Enhancement Opportunities: {"problems_i": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Non-critical score: float

### Performance Problems
- Critical Problems: {"problems_g": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Important Problems: {"problems_f": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Enhancement Opportunities: {"problems_j": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Non-critical score: float

### UI Problems (If project has UI)
- Critical Problems: {"problems_k": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Important Problems: {"problems_l": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Enhancement Opportunities: {"problems_m": {"problem_name": str, "problem_description": str, "problem_executable_proof": str, "confidence": float}; ...}
- Non-critical score: float

Where:
- non-critical score is a score of the diff excluding criticla problems. In ralph loop, you set a score treshold which determines when you pass review or not. Critical problems are not considered because all of them must be solved. This allows humans to be called only for: (1) answering questions; (2) solving stuck loops; (3) reviewing low-confidence critical problems.
- critical high-confidence regression-inducing false positives (generally functional false positives) extremely close to 0
- non-critical high-confidence regression-inducing false positives (generally functional false positives) are close to 0
- code quality & documentation problems proofs are cosntructed by makig a test that calls LLMs to check
- the whole system is trained end-to-end with RL to optimize for extemely low high-confidence regression-inducing false positives, high-recall, decent precision and calibrated confidence.
- spec_grounding is a causal justification how why the code is not prducing the intended behaviour described in the spec, with spec citations
- can be configured to use a different random_seed at each run (usefull for getting ralph loops unstuck) where it uses different models, and slighly different prompts and hyperparameters
- awarm has specific agents for each type of problem
- lunch review can be configured for 3 different levels of compute: light (lowest cost, good results), medium (costs more, better results) and heavy (costs the most, best results)

## Humans become the escalation layer

The ideal workflow is simple: claude code is generating code, lunch review is reviwing code and human are only called for these things:
- aswering questions
- solving stuck loops
- analyzing low-confidence critical problems surfaced

This changes the role of the engineer from **reviewer of every change** to **investigator of the changes that actually need human judgment**. Instead of spending hours reviewing code that is probably correct, engineers spend their time on the small fraction of changes.

## Review during development via CLI or Review PRs via Gitub Action

- During development: via cli or triggered by your main coding agent harness (e.g., claude code), to catch problems before PR phase. The earlier a problem is encounterd, the cheaper it is to fix it.
- For PR: at every PR or triggered by PR comment. Requests changes for critical problems.

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
        - token consumption
    4. punished more for:
        - producing a mistake
    5. punsihed a lot for:
        - giving high confidence on a mistake
4. train the entire system together (models, prompts, hyperparameters) via RL on real PRs that got reviewed and patched.
5. distill the traces of the big system and finetune the best open course coding model on it.
6. Result: A single model/Single Agent for cidate generation + A single model/Single Agent for candidate verification

## Future Development Towards Developing of our own Coding Agent: Why Review Must Go Beyond Final Diff Review for the purpose of SFTing a Coding Agent (alongside RL)

A final diff is only the **end state** of a coding agent's work. For training coding agents, reviewing only the final diff loses the most valuable information: **how the agent arrived there**.

A final diff review answers:

> **"Is this change correct?"**

A trajectory review answers the more important training question:

> **"Where did the agent's reasoning or implementation process go wrong, and how should it have proceeded?"**

### Final Diff Review

The fundamental review task remains:

```text
Specification + Code Diff → Review
```

This provides high-quality supervision for correctness, regressions, missing requirements, and architectural violations.

However, it cannot reliably identify:

* incorrect decisions that were later corrected,
* wasted exploration,
* dead-end approaches,
* the first point where the agent diverged from the plan,
* why a particular change was made,
* whether the agent is converging toward the intended solution.

### Trajectory Review

For agent training, the model should additionally review the **sequence of changes**:

```text
Specification + Plan + Diff₁ → Diff₂ → ... → Diffₙ
                         ↓
                   Review FM
```

This allows the reviewer to provide **process-level supervision**:

* identify the first incorrect step,
* explain why it diverged from the specification or plan,
* distinguish productive exploration from wasted work,
* identify missing investigations or tests,
* recommend the next corrective action,
* assess whether the trajectory is converging.

This generates a batch of new training data to SFT the model on.

This is particularly important because coding-agent RL has sparse objective feedback: a trajectory that receives `0 new tests passed` could represent either a completely misguided attempt or a nearly-correct solution that failed on one final detail.

SFT happens alongside RL which rewards states where more tests pass and penalizes a bit avery step.

All the training should ahppen under the lens of Curriculum Learning:

```text
SFT & RL for very simples coding tasks -> SFT & RL for a bit harder coding tasks -> SFT & RL for a bit harder coding tasks ....
```

## First Users

The product will win fastest in more mature repos with strong existing unit-test coverage, clear module boundaries, high AI-coding usage, and high business risk (e.g., financial logic, smart contracts, core API services) where compute costs are easily justified by defect prevention.
