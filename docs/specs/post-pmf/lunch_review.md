# Lunch Review: A Foundation Model for Autonomous Code Review

## The Thesis

AI coding agents are rapidly increasing software production. The bottleneck is no longer writing code — it is verifying that code is correct.

Today's AI reviewers (Greptile, CodeRabbit, Claude Code review mode, etc.) use general-purpose LLMs as reviewers. They optimize for generating useful review comments.

**Lunch Review is different: we are building a code review foundation model whose objective is discovering real software defects.**

Instead of treating code review as prompt engineering, we train specialized review models and jointly optimize an autonomous reviewer swarm around a measurable objective:

> **maximize discovery of real defects while minimizing high-confidence false positives.**

---

# Building a Review Foundation Model

The long-term goal is not a GitHub bot.

It is a **review foundation model**: a family of models trained specifically for software verification rather than code generation.

The foundation model learns capabilities that existing coding models are not explicitly trained for:

- finding functional defects,
- finding security vulnerabilities,
- detecting regressions,
- verifying behavior against specifications,
- constructing executable proofs,
- estimating calibrated confidence for every finding.

The final product is distilled into lightweight reviewer models that can run continuously during development or inside pull request workflows.

---

# Claude Code Writes Code. Lunch Review Tries to Break It.

Claude Code is optimized to produce working implementations.

**Lunch Review is optimized to prove implementations are wrong.**

The reviewer continuously investigates code by:

- reading specifications,
- exploring the repository,
- executing tools,
- running tests,
- generating new tests,
- reproducing failures,
- verifying hypotheses.

Pull request review is simply one deployment mode. The same reviewer can operate while code is being written.

---

# What Makes Lunch Review Different?

## 1. Specialized Reviewer Swarm

The system is composed of specialized reviewer agents.

Examples include:

- functional correctness,
- security,
- performance,
- architecture,
- documentation,
- UI behavior.

Each specialist is responsible for generating and verifying findings inside its own domain.

This specialization allows each model to become significantly stronger than a single general-purpose reviewer.

## 2. Executable Evidence, Not Opinions

Every finding is treated as a hypothesis that should be proven whenever possible.

The reviewer escalates verification through increasingly expensive execution environments.

| Level | Verification |
|-------|--------------|
| **L1** | Static code and AST analysis |
| **L2** | Unit tests or mocked execution |
| **L3** | Containerized integration tests |
| **L4** | Full multi-service topology simulation |

Whenever possible, every finding includes:

- specification grounding,
- causal explanation,
- executable failing test,
- reproduction instructions.

The output is evidence rather than speculation.

## 3. Learned Confidence Calibration

Current AI reviewers sound equally confident whether they are correct or hallucinating.

Lunch Review explicitly learns confidence as part of training.

Every finding predicts its probability of being correct.

The model is rewarded for:

- high confidence when correct,
- low confidence when uncertain.

It is punished heavily for:

- high confidence when incorrect.

The result is a reviewer that learns what it actually knows.

Confidence becomes a machine-actionable signal rather than a writing style.

---

# Two Operating Modes

## Lunch Review Human

The human reviewer is the final decision-maker.

Objective:

**maximize defect recall.**

The swarm aggressively surfaces hypotheses, evidence, and confidence scores. Humans filter speculative findings.

## Lunch Review Ralph

The reviewer operates inside an autonomous coding-agent loop.

Objective:

**minimize high-confidence regression-inducing false positives while maintaining high recall.**

Incorrect high-confidence findings are extremely costly because another agent may modify code in response to them.

---

# Review Architecture

## Stage 1 — Understand the Codebase

Multiple agents build an understanding of:

- repository structure,
- specifications,
- documentation,
- dependency graph,
- recent changes.

Agents generate hypotheses and identify areas worth investigating.

## Stage 2 — Investigate and Generate Candidates

Specialists independently search for problems inside their domain.

This stage intentionally prioritizes **high recall**.

Agents may:

- inspect execution paths,
- generate tests,
- mutate inputs,
- run tools,
- inspect runtime behavior.

The output is a large pool of candidate findings.

## Stage 3 — Verify and Filter Candidates

Independent verifier agents evaluate every candidate.

Verification includes:

- reproducing failures,
- validating executable proofs,
- checking specification grounding,
- rejecting unsupported hypotheses,
- assigning calibrated confidence.

This stage converts a noisy high-recall candidate pool into a high-precision review report.

---

# Review Output

Every finding contains structured evidence.

```yaml
category: Functional

severity: Critical

confidence: 0.98

spec_grounding:
  - Checkout specification §3.2

evidence:
  - failing integration test
  - execution trace
  - reproduction command
```

Findings are grouped into:

- Functional
- Security
- Performance
- Code Quality
- Documentation
- UI (optional)

Each category contains:

- Critical Problems
- Important Problems
- Enhancement Opportunities

The reviewer also produces an overall non-critical quality score for autonomous review loops.

---

# Training Strategy

## Stage 1 — Train Specialist Reviewers

Each specialist is trained independently.

Training combines supervised learning and reinforcement learning on datasets specific to its domain.

Examples include:

- real-world bugs,
- vulnerability datasets,
- mutation-testing datasets,
- performance regressions,
- maintainability improvements.

Each specialist learns to maximize performance on its own review task.

## Stage 2 — Train Confidence Calibration

Every specialist is further trained with reinforcement learning to produce calibrated confidence.

Rewards encourage:

- high confidence on correct findings,
- low confidence on uncertain findings.

Penalties increase dramatically for confidently incorrect findings.

Calibration is learned before specialists interact with each other.

## Stage 3 — End-to-End Joint Optimization

Once specialists are individually strong, the entire reviewer swarm is optimized jointly.

This optimization adjusts:

- prompts,
- routing,
- verifier aggregation,
- model selection,
- tool usage,
- investigation depth,
- compute allocation,
- hyperparameters.

The optimization target is **system-level review quality**, not individual model quality.

The swarm learns how specialists collaborate to maximize overall defect discovery while minimizing costly false positives.

---

# Humans Become the Escalation Layer

The ideal workflow is simple.

```text
Claude Code writes code.

Lunch Review investigates code.

Humans resolve ambiguity.
```

Humans are only involved when:

1. the reviewer needs specification clarification,
2. autonomous loops become stuck,
3. a critical finding has insufficient confidence.

Engineers stop reviewing every change and instead investigate the small fraction of changes that require human judgment.

---

# Why This Matters

As AI agents generate increasingly large volumes of software, trustworthy verification becomes the scarce resource.

Lunch Review is building the review equivalent of a coding foundation model: a system trained specifically to understand, investigate, verify, and prioritize software defects.

The objective is simple:

> **maximize discovery of real defects while making confidence reliable enough for autonomous software engineering.**
