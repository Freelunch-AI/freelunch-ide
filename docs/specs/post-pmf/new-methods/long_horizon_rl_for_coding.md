# Learning Long-Horizon Coding Agents Through Search, Verification, and Hierarchical Recovery

## Abstract

We propose a training framework for long-horizon coding agents based on agent-controlled branch-and-verify search, recursive counterfactual trajectory expansion, adaptive search allocation, and hierarchical recovery learning.

Long-horizon software engineering is difficult for reinforcement learning because a task may require hundreds or thousands of heterogeneous actions before its final correctness can be evaluated. Terminal rewards therefore create severe exploration and credit-assignment problems, while the space of possible code-edit trajectories is effectively intractable.

Our central idea is to transform long-horizon coding into structured search over semantic state transitions. A **Planner** decomposes a PRD into semantic subtasks, while an **Implementor** executes them. During training, the agent explores alternative implementations in isolated repository branches. Only verified successful branches are merged into the canonical repository history, producing clean happy-path trajectories that exclude failed experiments and rollback behavior.

These happy trajectories are recursively expanded at selected states. Rather than uniformly exploring the enormous trajectory space, the system preferentially samples states with high uncertainty, semantic subtask boundaries, high downstream impact, or low previous sampling frequency. Expensive LLM-powered search is used selectively, with its state-selection decisions potentially distilled into cheaper search heuristics.

The resulting search tree provides multiple independently verified continuations from the same state. This creates counterfactual supervision that allows policies to learn not only which actions succeed, but which successful decisions are preferable under task-quality and compute constraints.

Search alone, however, creates a deployment mismatch: during search, failed branches can simply be discarded, whereas a deployed coding agent must often recover from mistakes in its current repository state. We therefore introduce hierarchical recovery learning with different objectives for the Planner and Implementor:

* **Implementor recovery is local:** From a broken state, the Implementor learns to diagnose and repair the current semantic subtask and terminates when that subtask is complete.
* **Planner recovery is global:** When evidence reveals that the overall decomposition, ordering, or architecture is poor, the Planner learns to revise the remaining task plan and receives credit according to the eventual task outcome.

Planner and Implementor are trained jointly through alternating optimization, repeatedly freezing one policy while improving the other. Intermediate Planner-defined verification goals provide dense signals to the Implementor only when they represent genuine semantic milestones. The Planner itself is optimized against the fixed final evaluation objective, preventing it from gaming its reward by creating trivial intermediate tests.

Finally, after Planner and Implementor training, both policies are frozen and a separate **Test Maker** policy is trained to transform a PRD into an executable final test suite. The Test Maker never observes coding trajectories. Its objective is to produce tests that faithfully encode the specification and efficiently discriminate correct implementations from flawed ones, using mutation testing as an automatic anchor and human judgment as the ultimate specification-fidelity evaluation.

> **Central Hypothesis:**
> *Long-horizon coding policies can be learned more effectively by selectively searching and verifying alternative semantic trajectories, while separately learning local implementation recovery and global planning recovery, than by relying on single-path imitation or uniformly explored long-horizon reinforcement learning.*

---

## 1. The Core Problem

A realistic coding task can require a complex sequence of dependencies:

```
PRD 
 ↓ 
Understand repository 
 ↓ 
Architectural decision 
 ↓ 
Decompose task 
 ↓ 
Modify interfaces 
 ↓ 
Implement feature 
 ↓ 
Update callers 
 ↓ 
Debug compilation 
 ↓ 
Run tests 
 ↓ 
Discover architectural problem 
 ↓ 
Replan 
 ↓ 
Refactor 
 ↓ 
Integration testing 
 ↓ 
Performance validation 
 ↓ 
E2E success

```

The final reward might be:

$$R = \begin{cases} 1 & \text{if the task satisfies the specification} \\ 0 & \text{otherwise} \end{cases}$$

This creates four fundamental problems:

1. **Sparse reward:** Correctness may only become observable hundreds of actions later.
2. **Huge action space:** There are enormous numbers of possible code modifications and tool interactions.
3. **Long-range dependencies:** An architectural decision early in the task can determine whether a later subtask is easy or impossible.
4. **Expensive exploration:** Exhaustive search over possible coding trajectories is computationally infeasible.

The proposed framework treats this as a hierarchical search-and-learning problem, rather than attempting to directly solve the entire trajectory space with conventional RL.

---

## 2. Hierarchical Agent Structure

The coding agent consists primarily of two interacting policies.

### Planner

Responsible for:

* Understanding the PRD
* Decomposing the task into semantic subtasks
* Ordering those subtasks
* Deciding when the current plan is no longer appropriate
* Replanning remaining work
* Defining meaningful intermediate verification goals

### Implementor

Responsible for:

* Executing the current semantic objective
* Modifying the repository
* Using development tools
* Running relevant checks
* Diagnosing implementation failures
* Completing the current subtask

The distinction is fundamental:

> **Planner** decides *what* should happen.
> **Implementor** decides *how* to execute it.

---

## 3. Semantic State Transitions

The framework does not treat every token or tool call as an equally meaningful decision. Instead, it organizes trajectories around semantic transitions:

$$S_t \xrightarrow{\text{subtask}} S_{t+1}$$

where $S_{t+1}$ represents a repository state in which the current semantic objective has been completed.

For example:

```
S0
 │
 ├── Implement authentication abstraction
 ▼
S1
 │
 ├── Implement OAuth provider
 ▼
S2
 │
 ├── Migrate API callers
 ▼
S3

```

This provides natural points for verification, branching, counterfactual exploration, recovery, and credit assignment.

---

## 4. Agent-Controlled Happy-Path Search

During training, the agent can explore multiple candidate implementations using isolated repository branches.

```
                  S
          ┌───────┼───────┐
          │       │       │
       Branch A Branch B Branch C
          │       │       │
        FAIL   SUCCESS   FAIL
          │       │       │
       discard  MERGE   discard
                  │
                  ▼
                 S'

```

> **Critical Property:** The agent itself constructs the canonical trajectory by promoting only verified successful branches.

Thus the canonical history contains clean progressions ($S_0 \rightarrow S_1 \rightarrow S_2 \rightarrow S_3 \rightarrow \text{SUCCESS}$) rather than sequences cluttered by bad attempts, rollbacks, and retries.

The failed branches are still valuable training material for recovery learning, but they are not represented as part of the canonical happy path. This produces a clean distinction between **solution trajectories** and **recovery trajectories**.

---

## 5. Verification Is External to the Policy

The environment determines whether proposed candidate repository states are correct:

$$\text{agent action} \rightarrow \text{repository state} \rightarrow \text{tests / checkers / static analysis} \rightarrow \text{objective result} \rightarrow \text{reward}$$

* The **agent** controls which candidate to promote.
* The **environment** controls whether that candidate is actually successful.

This prevents the search procedure from becoming self-confirming.

---

## 6. Recursive Counterfactual Expansion

A single successful trajectory is not necessarily optimal. Suppose the system discovers $S_0 \rightarrow S_1 \rightarrow S_2 \rightarrow S_3$. It can revisit $S_1$ and independently explore alternative continuations:

```
                  S1
          ┌───────┼───────┐
          │       │       │
          A       B       C
          │       │       │
        FAIL   SUCCESS SUCCESS
                  │       │
                 .72     .94

```

The successful alternatives become additional training examples. This process is recursive:

$$\text{happy trajectory} \rightarrow \text{select valuable state} \rightarrow \text{branch} \rightarrow \text{verify} \rightarrow \text{retain successful continuations} \rightarrow \text{repeat}$$

The result is a tree of independently verified solution trajectories.

---

## 7. Counterfactual Supervision

At a state $s$, we can observe:

$$D_s = \{(a_1, R_1), (a_2, R_2), \dots, (a_n, R_n)\}$$

rather than merely observing one action. This allows the learner to distinguish:

* *"This action eventually worked."*

from:

* *"From this exact state, this alternative produced a better solution than the other explored alternatives."*

| Action | Approach | Outcome | Resource Cost |
| --- | --- | --- | --- |
| **Action A** | Extend existing interface | Works | 500k tokens |
| **Action B** | Introduce dedicated abstraction | Works | **200k tokens** |
| **Action C** | Rewrite authentication layer | Works | 1M tokens |

With an explicit compute-aware objective, **Action B** becomes the preferred decision even though all three satisfy functional correctness.

---

## 8. Compute-Aware Search

Search computation is itself part of the optimization problem. The conceptual objective is:

$$R(\tau) = R_{\text{task}} - \lambda C(\tau)$$

where:

$$C = C_{\text{tokens}} + C_{\text{model calls}} + C_{\text{branches}} + C_{\text{tests}} + C_{\text{execution}}$$

The goal is not to search as much as possible, but to **spend additional computation when its expected value exceeds its cost**. Routine states should be solved directly, while ambiguous or high-impact states receive additional search.

---

## 9. Adaptive Search Controller

Exhaustively expanding every state is infeasible. The system learns to prioritize states for expensive search based on key signals:

```
                    state
                      │
                      ▼
                 cheap scorer
                      │
          ┌───────────┴───────────┐
          │                       │
      low value               high value
          │                       │
     single pass           expensive search
                                  │
                                  ▼
                          branch-and-verify

```

* **Signals used:** Model uncertainty, candidate disagreement, semantic subtask boundaries, expected downstream impact, novelty, previous sampling frequency, historical search value, and remaining compute budget.

---

## 10. LLM-Powered Search

An LLM can occasionally act as a search controller by evaluating state uncertainty and potential downstream consequences. It provides expensive guidance only when useful, which can subsequently be distilled into a cheaper state-prioritization model:

$$f_\theta(S) \rightarrow P(\text{search is valuable})$$

This turns the LLM into a teacher for a learned, light-weight search heuristic.

---

## 11. Non-Uniform State Sampling

The system deliberately oversamples:

* **High-uncertainty states:** States where multiple actions appear plausible.
* **Semantic subtask boundaries:** States where a major semantic decision begins or ends.
* **Undersampled states:** States that have received relatively little exploration.
* **High-impact states:** States where an early decision can substantially affect downstream complexity.

The resulting sampling distribution is approximately:

$$p(s) \propto f(U(s), B(s), I(s), N(s), C(s))$$

where $U$ is uncertainty, $B$ is subtask boundaries, $I$ is expected impact, $N$ is novelty/sampling deficiency, and $C$ represents computational considerations.

---

## 12. The Critical Search/Deployment Mismatch

Happy-path search has an inherent flaw:

* **During training:** $\text{bad branch} \rightarrow \text{discard} \rightarrow \text{continue from clean state}$
* **At deployment:** $\text{bad implementation} \rightarrow \text{agent must recover}$

The deployed agent cannot always rewind the repository to the last known-good state. Therefore, successful-path search alone may produce a policy that is brittle when reality deviates. This motivates explicit recovery learning.

---

## 13. Hierarchical Recovery

Planner and Implementor recover at different scales:

### Local Implementor Recovery

The Implementor is responsible for recovering from an incorrect implementation of the current semantic subtask:

$$G_I: S_{\text{bad}} \rightarrow S_{\text{subtask-complete}}$$

The recovery episode terminates as soon as the current subtask is successfully completed.

### Global Planner Recovery

The Planner is responsible for recognizing when the overall plan is flawed:

$$G_P: (S_t, P_{\text{remaining}}) \rightarrow P'_{\text{remaining}}$$

The Planner can revise remaining subtasks, subtask ordering, architecture, interfaces, decomposition, or implementation strategy across the rest of the task.

---

## 14. Local Implementor Recovery

A failed implementation creates a dead-end state ($S_0 \rightarrow S_1 \rightarrow S_{\text{dead}}$). Recovery training begins directly at $S_{\text{dead}}$:

$$S_{\text{dead}} \rightarrow \text{inspect failure} \rightarrow \text{diagnose root cause} \rightarrow \text{modify implementation} \rightarrow \text{run checks} \rightarrow \text{subtask complete}$$

The episode ends as soon as the subtask is complete, eliminating the need to execute all remaining downstream subtasks and making recovery data inexpensive to generate.

---

## 15. Generating Implementor Recovery Data

Recovery starting states are constructed using multiple sources:

* Real failed search branches
* Synthetic AST/code mutations
* Incorrect intermediate implementations
* Dependency/interface corruption
* Deliberately injected logical bugs
* Intermediate verification failures

One happy trajectory can yield many local recovery training episodes:

$$\text{known-good state} \rightarrow \text{inject subtle bug} \rightarrow \text{broken state} \rightarrow \text{Implementor recovery} \rightarrow \text{subtask complete}$$

---

## 16. Global Planner Recovery

When an underlying plan makes downstream work unnecessarily difficult, a strong Planner must recognize that the approach itself is flawed rather than forcing the Implementor to patch a bad design.

* **Original Plan:** $\text{A: Refactor DB} \rightarrow \text{B: Add Caching} \rightarrow \text{C: Migrate API} \rightarrow \text{D: Observability}$
* **Revised Plan:** $\text{A': Introduce Repository Abstraction} \rightarrow \text{C: Migrate API} \rightarrow \text{B: Add Caching} \rightarrow \text{D: Observability}$

The Planner is explicitly empowered to alter the remaining future trajectory.

---

## 17. Why Planner Recovery Needs Global Credit

A planning decision at step 10 might only reveal its consequences at step 80.

$$\text{Step 10: Choose Arch A} \rightarrow \text{Step 30: Feature} \rightarrow \text{Step 60: Callers} \rightarrow \text{Step 80: Incompatibility} \rightarrow \text{Step 85: Replan} \rightarrow \text{Step 120: Success}$$

The value of the Planner's recovery decision cannot be measured at a single subtask boundary. Therefore, Planner recovery receives credit based on the **eventual outcome of the revised global trajectory**.

---

## 18. Two Recovery Operators

The system explicitly learns two distinct forms of recovery:

* **Local Corrective Control:**

$$G_I: \text{bad implementation} \rightarrow \text{current subtask complete}$$


* **Global Strategic Recovery:**

$$G_P: \text{bad global plan} \rightarrow \text{better remaining plan}$$



This mirrors real software engineering: *"The code is wrong"* versus *"The code is fine, but the approach was wrong."*

---

## 19. Planner/Implementor Alternating Optimization

Planner and Implementor are coupled policies trained through alternating optimization:

$$\left(P_0, I_0\right) \xrightarrow{\text{Train } P} \left(P_1, I_0\right) \xrightarrow{\text{Train } I} \left(P_1, I_1\right) \xrightarrow{\text{Train } P} \left(P_2, I_1\right) \xrightarrow{\text{Train } I} \left(P_2, I_2\right) \dots$$

Each iteration leverages happy-path search data, counterfactual trajectories, Implementor recovery episodes, and Planner global-replanning episodes.

---

## 20. Intermediate Verification

The Planner defines intermediate verification milestones for subtasks. These checks provide dense learning signals to the Implementor.

To prevent reward-hacking by the Planner:

* **Implementor:** Evaluated against Planner-defined intermediate verification.
* **Planner:** Evaluated strictly against the fixed final evaluation suite, task success, solution quality, and compute efficiency.

---

## 21. Fixed Final Evaluation

Training uses a fixed final evaluation objective ($T_0$) consisting of unit tests, integration tests, end-to-end tests, static analysis, and performance checks:

$$R_{\text{final}} = f(\text{correctness}, \text{performance}, \text{quality}, \text{compute})$$

This maintains an uncorrupted, stable ground truth while the Planner and Implementor co-adapt.

---

## 22. Test Maker as a Separate Final Stage

Once Planner and Implementor training is complete, both policies are frozen:

$$(P, I) \rightarrow \text{FREEZE}$$

A separate **Test Maker** is then trained on the objective: $\text{PRD} \rightarrow \text{final executable test suite}$.

The Test Maker **never observes coding trajectories** (happy paths, failed branches, or recovery actions). It only receives the specification and repository context, preventing it from over-fitting to the specific coding patterns of the agent.

---

## 23. Test Maker Objective

The Test Maker uses mutation testing as an automatic anchor:

```
                  PRD
                   │
                   ▼
            Generated Tests
                   │
         ┌─────────┴─────────┐
         ▼                   ▼
Correct Implementation   Mutated Implementation
         │                   │
         ▼                   ▼
       PASS                FAIL

```

The conceptual objective is:

$$R_T = R_{\text{specification}} + \lambda R_{\text{mutation}} - \mu C_{\text{tests}}$$

Direct specification fidelity (including human judgment) is evaluated alongside mutation scores to ensure comprehensive requirement coverage.

---

## 24. Why the Test Maker Is Trained Last

The Test Maker learns to construct tests tailored to evaluate a mature, fixed coding system:

$$\text{PRD} \rightarrow \text{Planner + Implementor (Alternating Training)} \rightarrow \text{Mature Agent (Frozen)} \rightarrow \text{Train Test Maker} \rightarrow (\text{PRD} \rightarrow \text{Test Suite})$$

Its core goal is to transform vague requirements into the most useful executable specification for evaluating code correctness.

---

## 25. Complete Training Loop

```
                              PRD
                               │
                               ▼
                         ┌───────────┐
                         │  PLANNER  │
                         └─────┬─────┘
                               │
                        semantic subtask
                               │
                               ▼
                        ┌──────────────┐
                        │ IMPLEMENTOR  │
                        └──────┬───────┘
                               │
                               ▼
                        repository state
                               │
                               ▼
                      adaptive search gate
                               │
                 ┌─────────────┴─────────────┐
                 │                           │
            single pass              branch-and-verify
                 │                           │
                 │                 verified alternatives
                 │                           │
                 └─────────────┬─────────────┘
                               │
                               ▼
                        happy trajectory
                               │
                               ▼
                    recursive counterfactual
                           expansion
                               │
                               ▼
                       policy improvement


                    PARALLEL RECOVERY

                    happy trajectories
                               │
                 ┌─────────────┴─────────────┐
                 │                           │
      implementation dead ends         plan failures
                 │                           │
                 ▼                           ▼
          local Implementor           global Planner
             recovery                    recovery
                 │                           │
                 ▼                           ▼
          subtask complete             revised plan
                 │                           │
                 └─────────────┬─────────────┘
                               │
                               ▼
                      alternating training
                               │
                               ▼
                            repeat
                               │
                               ▼
                     Planner + Implementor
                            FREEZE
                               │
                               ▼
                          Test Maker
                               │
                               ▼
                      PRD → final tests

```

---

## 26. The Four Core Learning Mechanisms

1. **Verified Trajectory Search:** Discover successful implementations through isolated branch-and-verify exploration.
2. **Counterfactual Expansion:** Explore multiple alternative successful continuations from high-value states.
3. **Hierarchical Recovery:** Train local Implementor recovery and global Planner replanning.
4. **Adaptive Computation:** Allocate expensive search exclusively to states with high expected value.

$$\text{Search} \rightarrow \text{Verify} \rightarrow \text{Compare} \rightarrow \text{Learn} \rightarrow \text{Recover}$$

---

## 27. What Makes This Different From Simply Doing RL

| Dimension | Conventional RL | Proposed Formulation |
| --- | --- | --- |
| **Trajectory View** | $s_0 \rightarrow a_0 \rightarrow a_1 \dots \rightarrow a_{1000} \rightarrow R$ | Branching tree over semantic states ($S_t \xrightarrow{\text{subtask}} S_{t+1}$) |
| **Search & Verify** | Evaluates long, noisy, single-path trajectories | Isolates subtasks in repository branches; retains verified merges |
| **Data Generation** | Relies on sparse terminal success/failure | Generates counterfactual state pairs and targeted recovery pairs |
| **Recovery Learning** | Implicitly mixed into global trajectory | Separated into local (subtask) and global (plan) recovery operators |

---

## 28. Experimental Hypotheses

* **H1 (Counterfactual Search):** Multiple verified continuations from the same state provide better learning signals than isolated successful trajectories.
* **H2 (Recursive Expansion):** Repeatedly expanding valuable states produces better policies than one-generation trajectory search.
* **H3 (Adaptive Search):** Selective search achieves a better performance/compute tradeoff than uniform branching.
* **H4 (Local Implementor Recovery):** Training on dead-end $\rightarrow$ subtask-completion trajectories substantially improves robustness to implementation failures.
* **H5 (Global Planner Recovery):** Training the Planner to recognize and revise flawed plans improves performance as task horizon increases.
* **H6 (Hierarchical Recovery):** Combining local Implementor recovery with global Planner recovery outperforms either mechanism alone.
* **H7 (Search-to-Policy Transfer):** Search-discovered behaviors can be internalized into policy weights, reducing inference-time compute requirements.

---

## 29. Critical Ablations

A comprehensive evaluation should compare the full system against isolated components and baselines:

1. SFT on successful trajectories
2. Standard RL / RLVR
3. RL with successful-trajectory replay
4. Branch-and-verify
5. Branch-and-verify + counterfactual expansion
6. Counterfactual expansion + adaptive search
7. Search + Implementor recovery
8. Search + Planner recovery
9. Search + both recovery mechanisms
10. **Full proposed system**

---

## 30. Evaluation

### Primary Axes

* $\text{Success Rate} \quad\text{vs.}\quad \text{Task Horizon}$
* $\text{Success Rate} \quad\text{vs.}\quad \text{Compute per Task}$

### Key Metric Categories

* **Coding Performance:** Pass@1, task success, success vs. horizon, unseen repositories, solution quality.
* **Search Efficiency:** Tokens/task, model calls/task, branches/task, test executions, wall-clock time, compute/successful task.
* **Counterfactual Search Quality:** Improvement from additional branching, solution diversity, value of selected states, advantage over single-path training.
* **Planner Recovery:** Plan-error detection rate, replanning success, downstream task improvement, architectural mistake recovery.
* **Implementor Recovery:** Dead-end recovery rate, recovery tokens, recovery latency, performance by failure type.
* **Test Maker:** Requirement coverage, mutation score, specification fidelity, false-positive/negative rates, test execution cost, human-rated PRD compliance.

---

## 31. The Central Research Claim

Long-horizon coding should be treated as structured search over semantic state transitions, rather than an undifferentiated sequence of thousands of low-level actions.

> **Forward Decisions:** Learned through verified counterfactual search.
> **Failure Recovery:** Learned through hierarchical recovery (global strategy via Planner, local code via Implementor).

$$\underbrace{\text{Search}}_{\text{discover}} \rightarrow \underbrace{\text{Verification}}_{\text{evaluate}} \rightarrow \underbrace{\text{Counterfactual expansion}}_{\text{compare}} \rightarrow \underbrace{\text{Policy learning}}_{\text{internalize}} \rightarrow \underbrace{\text{Recovery learning}}_{\text{robustify}}$$

**Ultimate Objective:** Use expensive search during training to teach the coding agent how to solve long-horizon tasks without requiring that search at inference time.
