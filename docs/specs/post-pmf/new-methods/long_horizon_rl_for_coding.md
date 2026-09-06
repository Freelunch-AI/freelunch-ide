# Learning Long-Horizon Coding Agents Through Search, Verification, and Hierarchical Recovery

## Abstract

We propose a training framework for long-horizon coding agents based on four complementary ideas:

1. **Task-conditioned branch-and-verify search**
2. **Recursive counterfactual trajectory expansion**
3. **Task-conditioned value-guided search allocation**
4. **Hierarchical recovery learning**

Long-horizon software engineering is difficult for reinforcement learning because a task may require hundreds or thousands of heterogeneous actions before its final correctness can be evaluated. Terminal rewards therefore create severe exploration and credit-assignment problems, while the space of possible code-edit trajectories is effectively intractable.

Our central idea is to transform long-horizon coding into structured search over **repository states** conditioned on a specific task or PRD specification $T$.

A **Planner** decomposes $T$ into semantic subtasks and determines the overall strategy. An **Implementor** executes those subtasks. During training, the system explores alternative implementations in isolated repository branches. Candidate branches are externally verified, and only successful branches are promoted into canonical happy-path trajectories.

These trajectories are then recursively expanded at selected states. Instead of uniformly exploring the enormous trajectory space, search is concentrated on states where additional exploration is expected to be valuable. A learned task-conditioned value function estimates whether a repository state is promising for the specific task $T$, while an adaptive search controller considers uncertainty, novelty, semantic boundaries, downstream impact, and compute budget.

Because multiple rollouts can originate from the same state-task pair, the resulting data provides **counterfactual supervision**: the learner can compare alternative decisions taken from the same state for the same task rather than simply learning from isolated successful trajectories. These comparisons are used to improve both Planner and Implementor policies through grouped policy optimization.

Search, however, creates an important deployment mismatch. During training, a failed branch can simply be discarded, whereas a deployed agent must often recover from mistakes in its current repository state. We therefore train two distinct recovery capabilities:

* **Local Implementor recovery:** recover the current semantic subtask from a broken implementation state.
* **Global Planner recovery:** recognize that the overall strategy is flawed and revise the remaining plan.

The two recovery objectives operate at different temporal scales. Implementor recovery terminates when the current subtask is repaired. Planner recovery receives credit according to the eventual outcome of the revised global trajectory.

Planner and Implementor are trained through alternating optimization. The Planner defines intermediate semantic milestones, which provide dense supervision to the Implementor. However, the Planner itself is ultimately evaluated only against the fixed final task objective, preventing it from gaming training by inventing trivial intermediate goals.

After Planner and Implementor training, both policies are frozen. A separate **Test Maker** is then trained to transform the task specification $T$ into an executable final test suite. The Test Maker does not observe coding trajectories, preventing it from overfitting to the particular behavior of the coding agent. Mutation testing provides an automatic training signal, while specification fidelity is treated as the ultimate criterion.

> **Central Hypothesis**
>
> Long-horizon coding policies can be learned more effectively by selectively searching and verifying alternative semantic trajectories for a task, using task-conditioned value estimation to allocate search and separate recovery objectives to handle local implementation failures and global planning failures, than by relying on single-path imitation or uniformly explored long-horizon reinforcement learning.

---

# 1. The Core Problem

A realistic coding task $T$ can require a long sequence of interdependent decisions:

```text
PRD (T)
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

The final reward is determined by whether the repository satisfies the task specification.

```math
R(T) =
\begin{cases}
1 & \text{if the repository satisfies } T \\
0 & \text{otherwise}
\end{cases}
```

This creates four fundamental problems:

1. **Sparse reward:** Correctness with respect to $T$ may only become observable hundreds of actions later.
2. **Huge action space:** There are enormous numbers of possible code modifications and tool interactions.
3. **Long-range dependencies:** An architectural decision early in the task can determine whether a later subtask is easy or impossible.
4. **Expensive exploration:** Exhaustive search over coding trajectories is computationally infeasible.

The proposed framework treats coding as a **task-conditioned hierarchical search problem**, rather than as one undifferentiated sequence of low-level actions.

---

# 2. Hierarchical Agent Structure

The coding agent consists primarily of two interacting policies.

## Planner

The Planner is responsible for the global strategy for task $T$.

Its policy can be written as:

```math
\pi_P(a_P \mid s,T)
```

The Planner is responsible for:

* Understanding the PRD specification $T$
* Decomposing $T$ into semantic subtasks
* Ordering subtasks
* Making architectural decisions
* Defining meaningful intermediate verification milestones
* Detecting when the current plan is no longer appropriate
* Replanning the remaining task

## Implementor

The Implementor is responsible for executing the current semantic goal.

Its policy can be written as:

```math
\pi_I(a_I \mid s,T,g)
```

where $g$ is the current semantic subtask.

The Implementor is responsible for:

* Executing the current subtask
* Modifying the repository
* Using development tools
* Running relevant checks
* Diagnosing implementation failures
* Completing the current semantic goal

The distinction is fundamental:

> **Planner:** decides *what should happen* to satisfy $T$.
>
> **Implementor:** decides *how to execute it* given task $T$ and goal $g$.

---

# 3. Semantic State Transitions

The framework does not treat every token or tool call as an equally meaningful decision.

Instead, trajectories are organized around **semantic repository states**.

A semantic transition can be represented as:

```math
S_t \xrightarrow{\mathrm{subtask}} S_{t+1}
```

where $S_{t+1}$ represents a repository state in which the current semantic objective has been completed.

For example:

```text
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

This provides natural points for:

* Verification
* Branching
* Counterfactual exploration
* Recovery
* Task-conditioned credit assignment

The key abstraction is therefore not the individual token or shell command, but the **semantic transition produced by a sequence of low-level actions**.

---

# 4. Agent-Controlled Happy-Path Search

During training on task $T$, the agent explores multiple candidate implementations using isolated repository branches.

```text
                  State S
                    │
          ┌─────────┼─────────┐
          │         │         │
       Branch A  Branch B  Branch C
          │         │         │
         FAIL    SUCCESS     FAIL
          │         │         │
       discard    merge     discard
                    │
                    ▼
                   S'
```

The canonical trajectory is constructed by promoting only externally verified successful branches.

```math
S_0 \rightarrow S_1 \rightarrow S_2 \rightarrow S_3 \rightarrow \mathrm{SUCCESS}(T)
```

This has an important consequence:

The canonical history contains **clean successful progress**, rather than a mixture of:

* failed experiments
* rollbacks
* retries
* abandoned approaches

Failed branches are not discarded completely. They remain valuable as **recovery training data**.

This creates a clean distinction between:

* **Solution trajectories**
* **Recovery trajectories**

---

# 5. Verification Is External to the Policy

The environment determines whether a candidate repository state satisfies task requirements.

```text
Agent action
     ↓
Repository state
     ↓
Tests / checkers / static analysis
     ↓
Objective result
     ↓
Reward for task T
```

The separation is deliberate:

* The **agent** controls which candidates to explore and promote.
* The **environment** determines whether the candidate actually satisfies $T$.

This prevents the search procedure from becoming self-confirming.

The policy cannot simply declare a branch successful. Success is determined by an external evaluation process.

---

# 6. Recursive Counterfactual Expansion

A single successful trajectory is not necessarily the best trajectory.

Suppose search finds:

```math
S_0 \rightarrow S_1 \rightarrow S_2 \rightarrow S_3
```

The system can return to $S_1$ and independently explore alternative continuations.

```text
                  S1
           ┌──────┼──────┐
           │      │      │
           A      B      C
           │      │      │
          FAIL SUCCESS SUCCESS
                  │       │
                0.72     0.94
```

The important property is that the alternatives originate from the **same state under the same task**.

The search procedure therefore recursively performs:

```math
\mathrm{HappyTrajectory}(T)
\rightarrow
\mathrm{SelectState}(s,T)
\rightarrow
\mathrm{Branch}
\rightarrow
\mathrm{Verify}
\rightarrow
\mathrm{RetainSuccessfulBranches}
\rightarrow
\mathrm{Repeat}
```

The result is a search tree containing multiple independently verified continuations.

---

# 7. Counterfactual Supervision

At a state $s$ for task $T$, the system can observe multiple candidate actions and their downstream returns.

```math
D_{s,T}
=
\left\{
(a_1,R_1),
(a_2,R_2),
\ldots,
(a_K,R_K)
\right\}
```

This provides a fundamentally richer signal than a single successful trajectory.

Without counterfactual search, the learner sees:

> "Action A worked."

With counterfactual search, the learner can observe:

> "From this exact state for this exact task, Action A worked better than Actions B and C."

For example:

| **Action** | **Approach for Task $T$**       | **Outcome** | **Resource Cost** |
| ---------- | ------------------------------- | ----------- | ----------------: |
| **A**      | Extend existing interface       | Works       |       500k tokens |
| **B**      | Introduce dedicated abstraction | Works       |   **200k tokens** |
| **C**      | Rewrite authentication layer    | Works       |         1M tokens |

If the training objective includes compute efficiency, Action B is preferable even though all three solutions satisfy functional correctness.

This turns search into a source of **relative decision supervision**.

---

# 8. Search-to-Policy Reinforcement Learning

Counterfactual search produces multiple rollouts from the same state-task pair.

The search tree therefore provides samples of the form:

```math
(s,T,a,R)
```

where:

* $s$ is the semantic repository state.
* $T$ is the task specification.
* $a$ is the Planner or Implementor action.
* $R$ is the discounted return of the resulting rollout.

For example:

```text
State S, Task T

├── Action A → Return 0.92
├── Action B → Return 0.81
├── Action C → Return 0.37
└── Action D → Return 0.95
```

This naturally produces a **relative-ranking problem** conditioned on the exact state and task.

## Task-Grouped GRPO Over Shared-State Rollouts

For a batch of $K$ rollouts from the same state $s$ under the same task $T$:

```math
\mathcal{B}(s,T)
=
\left\{
(a_i,R_i)
\right\}_{i=1}^{K}
```

The group mean return is:

```math
\bar{R}(s,T)
=
\frac{1}{K}
\sum_{i=1}^{K}
R_i
```

The advantage for each rollout is:

```math
A_i
=
R_i-\bar{R}(s,T)
```

The policy objective is:

```math
\mathcal{L}_{\mathrm{GRPO}}
=
-\frac{1}{K}
\sum_{i=1}^{K}
\min
\left(
r_i A_i,
\mathrm{clip}(r_i,1-\epsilon,1+\epsilon)A_i
\right)
```

where the policy ratio is:

```math
r_i
=
\frac{
\pi_\theta(a_i \mid s,T)
}{
\pi_{\theta_{\mathrm{old}}}(a_i \mid s,T)
}
```

The important point is not merely that GRPO is used, but **where the groups come from**.

The groups are deliberately constructed so that multiple actions originate from the same:

```text
state s
task T
```

This gives the optimization a direct counterfactual comparison at the decision point being learned.

## Why Task-Grouped GRPO Fits Search

| **Conventional RL**                                            | **Task-Grouped Search GRPO**                                      |
| -------------------------------------------------------------- | ----------------------------------------------------------------- |
| Trajectories start from different states and tasks.            | Multiple trajectories intentionally start from the same $(s,T)$.  |
| Advantage depends on a global or learned baseline.             | Advantage is computed directly from sibling rollouts.             |
| Exploration is primarily generated by the policy.              | Exploration is generated by branch-and-verify search.             |
| Alternative decisions are rarely controlled at the same state. | Alternative decisions are explicitly generated at the same state. |

Because search explicitly constructs counterfactual alternatives, policy optimization receives the comparisons that matter for local decision making.

## Discounted Returns Over Semantic Decisions

Returns are computed over semantic transitions rather than individual tool calls.

```math
R_t
=
\sum_{k=t}^{T_{\max}}
\gamma^{k-t}
r_k(T)
```

The reward terms can include:

* Final correctness for task $T$
* Intermediate semantic verification
* Compute penalties
* Quality bonuses
* Performance bonuses

## Planner Updates

Planner training uses tuples of the form:

```math
(s,T,a_P,R)
```

where $a_P$ represents a planning decision such as:

* Subtask decomposition
* Subtask ordering
* Architectural choice
* Replanning

## Implementor Updates

Implementor training uses tuples of the form:

```math
(s,T,g,a_I,R)
```

where $a_I$ represents:

* Code edits
* Tool usage
* Debugging actions
* Implementation choices

## Search Distills Into the Policy

The resulting learning pipeline is:

```text
Branch-and-Verify Search
          │
          ▼
Multiple rollouts from the same (s,T)
          │
          ▼
Counterfactual training tuples
          │
          ▼
Task-grouped policy optimization
          │
          ▼
Planner and Implementor internalize search decisions
```

Expensive search is therefore used primarily during training.

The long-term goal is for the learned policies to reproduce good search decisions directly during inference.

---

# 9. Learning a Task-Conditioned Value Function

Search-generated trajectories provide supervision not only for the policy, but also for a **task-conditioned value function**.

The value function is:

```math
V_\psi(s,T)
```

It estimates the expected future return from repository state $s$ for task $T$.

```math
V_\psi(s,T)
=
\mathbb{E}_\pi
\left[
R_t
\mid
S_t=s,\mathrm{Task}=T
\right]
```

The policy answers:

> "What action should I take next?"

The value function answers:

> "How promising is this state for this task?"

This distinction is important because a repository state cannot be evaluated independently of the task it is supposed to satisfy.

For example, the same state could have:

```math
V_\psi(s,T_{\mathrm{REST}}) \approx 0.95
```

and:

```math
V_\psi(s,T_{\mathrm{GraphQL}}) \approx 0.05
```

The value is therefore explicitly **task-conditioned**.

---

# 10. Value-Guided Search Allocation

The learned value function has two primary uses.

## Early Dead-End Detection

Instead of rolling every branch out to complete task success, search can stop branches whose predicted value is sufficiently low.

```text
               State S, Task T
                      │
       ┌──────────────┼──────────────┐
       │              │              │
    Branch A       Branch B       Branch C
       │              │              │
     0.92            0.84           0.07
   continue        continue          prune
```

Conceptually:

```math
V_\psi(s,T) \leq \tau
```

indicates that the state is sufficiently unlikely to be worth further exploration.

Branches above the threshold continue to receive search.

This allows search to avoid expensive full rollouts for states that already appear to be dead ends.

## Search Prioritization

The value function can also determine where to spend additional search budget.

High-value states can receive:

* More rollouts
* More candidate actions
* Deeper counterfactual expansion

Low-value states can receive less computation.

```text
                 Root State
                /     |     \
             0.88    0.31    0.79
               │       │       │
             expand   prune   expand
```

This turns the value model into a **search allocation mechanism**, not merely a predictor used after search has already happened.

---

# 11. AlphaZero Analogy

The interaction between policy, value, and search is conceptually similar to AlphaZero.

| **AlphaZero**                          | **Long-Horizon Coding**                                 |
| -------------------------------------- | ------------------------------------------------------- |
| Policy proposes moves.                 | Planner and Implementor propose semantic actions.       |
| Value estimates future game outcome.   | Value estimates future task return.                     |
| MCTS expands promising states.         | Branch-and-verify expands promising repository states.  |
| Search generates training information. | Verified search generates policy and value supervision. |

The analogy is conceptual rather than literal.

The proposed coding system is not simply MCTS applied to repositories. Its search operates over semantic repository transitions, uses external software verification, and includes hierarchical Planner/Implementor recovery.

---

# 12. Training the Value Function

Every verified rollout produces a target return for the states it visited.

For example:

```text
S0 → S1 → S2 → S3 → SUCCESS(T)
```

Each state receives a discounted future return.

The value model can be trained with:

```math
\mathcal{L}_{\mathrm{value}}
=
\frac{1}{N}
\sum_{i=1}^{N}
\left(
V_\psi(S_i,T_i)-R_i
\right)^2
```

The value model therefore learns from the same verified search data that trains the policies.

This creates a feedback loop:

```text
Search
  ↓
Verified trajectories
  ↓
Policy learning + value learning
  ↓
Better actions + better search allocation
  ↓
Improved search
```

---

# 13. Adaptive Search Controller

The full search budget should not be spent uniformly.

A cheap search controller can determine whether a state deserves expensive branch-and-verify search.

```text
              State S, Task T
                     │
                     ▼
             Cheap Search Scorer
                     │
          ┌──────────┴──────────┐
          │                     │
       low value             high value
          │                     │
      single pass          expensive search
                                │
                                ▼
                        branch-and-verify
```

The controller can combine several signals:

* Task-conditioned uncertainty
* Candidate disagreement
* Semantic subtask boundaries
* Expected downstream impact
* Novelty
* Previous sampling frequency
* Task-conditioned value
* Remaining compute budget

The goal is not merely to search less.

The goal is to **search selectively where search is most informative**.

---

# 14. LLM-Powered Search Control

An LLM can act as an expensive search controller during training.

Its role is to estimate whether exploring a state is likely to produce valuable information for task $T$.

Conceptually:

```math
f_\theta(S,T)
\rightarrow
P(\mathrm{search\ is\ valuable})
```

The expensive controller can then be distilled into a cheaper model or heuristic.

The resulting sampling distribution can depend on:

```math
p(s\mid T)
\propto
f
\left(
U(s,T),
B(s),
I(s,T),
N(s),
C(s)
\right)
```

where:

* $U(s,T)$ is task-conditioned uncertainty.
* $B(s)$ indicates semantic boundaries.
* $I(s,T)$ estimates downstream impact.
* $N(s)$ captures novelty or under-sampling.
* $C(s)$ represents computational cost.

The learned value function can be incorporated as another signal into this allocation mechanism.

---

# 15. The Critical Search/Deployment Mismatch

Happy-path search creates a fundamental mismatch between training and deployment.

During training:

```text
bad branch
   ↓
discard
   ↓
continue from clean state
```

At deployment:

```text
bad implementation
   ↓
agent must recover
   ↓
continue from current repository state
```

A deployed coding agent cannot always rewind to the last known-good state.

Therefore, learning only successful forward trajectories is insufficient.

The agent must also learn **how to recover when its current state is already wrong**.

---

# 16. Hierarchical Recovery

We therefore introduce two different recovery objectives.

### Local Implementor Recovery

Recover the current semantic subtask from an incorrect implementation.

### Global Planner Recovery

Recognize that the overall strategy is flawed and revise the remainder of the task.

These are fundamentally different problems.

The first is **local corrective control**.

The second is **global strategic recovery**.

---

# 17. Local Implementor Recovery

Suppose the agent reaches a broken state:

```text
S0 → S1 → S_dead
```

The Implementor is trained to recover the current subtask directly from $S_{\mathrm{dead}}$.

Conceptually:

```math
G_I(S_{\mathrm{dead}},T,g)
\rightarrow
S_{\mathrm{subtask\ complete}}
```

The recovery episode looks like:

```text
Broken state
    ↓
Inspect failure
    ↓
Diagnose
    ↓
Modify code
    ↓
Run checks
    ↓
Subtask complete
```

The crucial property is the termination condition:

> **Implementor recovery stops when the current subtask is repaired.**

There is no need to execute all remaining subtasks.

This makes recovery data much cheaper to generate than full-task recovery trajectories.

---

# 18. Generating Implementor Recovery Data

Recovery states can come from:

* Failed search branches
* Intermediate verification failures
* Incorrect intermediate implementations
* Synthetic AST or code mutations
* Dependency corruption
* Interface corruption
* Injected logical bugs

A single successful trajectory can therefore generate many recovery episodes.

```text
Known-good state
      ↓
Inject failure
      ↓
Broken state
      ↓
Implementor recovery
      ↓
Subtask complete
```

The same trajectory becomes useful for both forward learning and recovery learning.

---

# 19. Global Planner Recovery

Not every failure should be repaired locally.

Sometimes the underlying plan itself is the problem.

For example, the original strategy for task $T$ might be:

```text
A: Refactor DB
→ B: Add Caching
→ C: Migrate API
→ D: Observability
```

After downstream evidence reveals that the architecture is wrong, the Planner might instead choose:

```text
A': Introduce Abstraction
→ C: Migrate API
→ B: Add Caching
→ D: Observability
```

The key distinction is that the Planner changes the **future trajectory**, rather than merely repairing the current code.

The global recovery operator can be represented as:

```math
G_P(S_t,T,P_{\mathrm{remaining}})
\rightarrow
P'_{\mathrm{remaining}}
```

---

# 20. Why Planner Recovery Needs Global Credit

Planning decisions can have consequences many semantic transitions later.

For example:

```text
Step 10
Choose Architecture A
        ↓
Step 30
Implement Feature
        ↓
Step 60
Update Callers
        ↓
Step 80
Discover Incompatibility
        ↓
Step 85
Replan
        ↓
Step 120
SUCCESS(T)
```

A Planner recovery decision at step 85 cannot be evaluated solely by whether the immediate subtask succeeds.

Its quality depends on whether the revised global trajectory eventually produces a better outcome for task $T$.

Therefore, Planner recovery receives credit according to the **eventual global outcome** of the revised plan.

---

# 21. Two Recovery Operators

The framework explicitly learns two different operators.

## Local Corrective Control

```math
G_I(S_{\mathrm{bad}},T,g)
\rightarrow
S_{\mathrm{subtask\ complete}}
```

The Implementor repairs the current subtask.

## Global Strategic Recovery

```math
G_P(S_t,T,P_{\mathrm{remaining}})
\rightarrow
P'_{\mathrm{remaining}}
```

The Planner revises the remaining strategy.

| **Failure Type**                   | **Recovery Mechanism** |
| ---------------------------------- | ---------------------- |
| Incorrect code or implementation   | Implementor recovery   |
| Poor architecture or decomposition | Planner recovery       |
| Local debugging failure            | Implementor recovery   |
| Global strategy failure            | Planner recovery       |

---

# 22. Planner/Implementor Alternating Optimization

Planner and Implementor are coupled policies.

Training both simultaneously can make the learning problem unstable because each policy continually changes the environment seen by the other.

We therefore use alternating optimization.

The training process is:

```math
(P_0,I_0)
\rightarrow
(P_1,I_0)
\rightarrow
(P_1,I_1)
\rightarrow
(P_2,I_1)
\rightarrow
\cdots
```

Conceptually:

```text
Train Planner with Implementor frozen
                  ↓
Train Implementor with Planner frozen
                  ↓
Train Planner with Implementor frozen
                  ↓
Train Implementor with Planner frozen
                  ↓
Repeat
```

Each iteration uses:

* Happy-path search data
* Counterfactual trajectories
* Value-function targets
* Implementor recovery episodes
* Planner recovery episodes

This lets each component improve against a relatively stable version of the other.

---

# 23. Intermediate Verification vs. Final Evaluation

The Planner defines intermediate verification milestones.

These milestones are useful because they provide dense feedback to the Implementor.

For example:

```text
Goal g1 → Authentication abstraction complete
Goal g2 → OAuth provider working
Goal g3 → API callers migrated
```

However, allowing the Planner to define its own reward would create a reward-hacking problem.

The solution is an asymmetric objective:

* **Implementor:** evaluated against Planner-defined intermediate goals.
* **Planner:** evaluated against the fixed final task objective.

The final task reward can be written as:

```math
R_{\mathrm{final}}(T)
=
f
\left(
\mathrm{correctness}(T),
\mathrm{performance}(T),
\mathrm{quality}(T),
\mathrm{compute}
\right)
```

The Planner therefore cannot redefine what success means.

The final evaluation remains externally specified by task $T$.

---

# 24. Test Maker as a Separate Final Stage

After Planner and Implementor training, both policies are frozen.

```math
(P,I)
\rightarrow
\mathrm{FREEZE}
```

A separate **Test Maker** policy is then trained.

Its objective is:

```math
\pi_M:
T
\rightarrow
\mathcal{E}_T
```

where $\mathcal{E}_T$ is the executable final test suite for task $T$.

The Test Maker does **not** observe coding trajectories.

It does not see:

* Happy paths
* Failed branches
* Recovery trajectories
* Planner decisions
* Implementor actions

It receives the task specification and relevant repository context.

This prevents the Test Maker from simply learning the coding style or failure patterns of the particular agent being trained.

---

# 25. Test Maker Objective

Mutation testing provides an automatic anchor.

```text
                              PRD T
                                │
                                ▼
                        Generated Tests E_T
                                │
                 ┌──────────────┴──────────────┐
                 │                             │
                 ▼                             ▼
        Correct Implementation        Mutated Implementation
                 │                             │
                 ▼                             ▼
               PASS                          FAIL
```

A conceptual test-generation reward is:

```math
R_T
=
R_{\mathrm{specification}}(\mathcal{E}_T,T)
+
\lambda R_{\mathrm{mutation}}(\mathcal{E}_T,T)
-
\mu C_{\mathrm{tests}}(\mathcal{E}_T)
```

The objective balances:

* Specification fidelity
* Mutation score
* Test quality
* Test execution cost

Mutation testing is useful because it provides an automatic signal for whether tests discriminate between correct and incorrect implementations.

However, mutation score alone is insufficient.

A test suite can achieve a high mutation score while still failing to encode important requirements in the PRD.

Therefore, **specification fidelity remains the ultimate criterion**.

---

# 26. Complete Training Loop

```text
                               PRD (T)
                                  │
                                  ▼
                            ┌───────────┐
                            │  PLANNER  │
                            └─────┬─────┘
                                  │
                          semantic subtask (g,T)
                                  │
                                  ▼
                            ┌──────────────┐
                            │ IMPLEMENTOR  │
                            └──────┬───────┘
                                   │
                           repository state (s)
                                   │
                                   ▼
                         adaptive search gate
                                   │
                    ┌──────────────┴──────────────┐
                    │                             │
               single pass                branch-and-verify
                    │                             │
                    │                    verified alternatives
                    │                             │
                    └──────────────┬──────────────┘
                                   │
                                   ▼
                         happy trajectory for T
                                   │
                                   ▼
                       counterfactual expansion
                           at selected (s,T)
                                   │
                                   ▼
                   policy update + value update
                                   │
                                   │
                ───────── RECOVERY TRAINING ─────────
                                   │
                    ┌──────────────┴──────────────┐
                    │                             │
         implementation dead ends           plan failures
                    │                             │
                    ▼                             ▼
             Implementor recovery           Planner recovery
                    │                             │
                    ▼                             ▼
            subtask complete                revised plan
                    │                             │
                    └──────────────┬──────────────┘
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
                          PRD T → final tests
```

---

# 27. The Four Core Learning Mechanisms

The framework combines four core mechanisms.

### 1. Task-Conditioned Verified Search

Discover successful implementations through isolated branch-and-verify exploration for task $T$.

### 2. Counterfactual Expansion

Explore multiple alternative continuations from the same semantic state for the same task.

### 3. Hierarchical Recovery

Train:

* Local Implementor recovery for implementation failures.
* Global Planner recovery for strategic failures.

### 4. Adaptive Task-Value Search

Use a task-conditioned value function and search controller to allocate expensive exploration selectively.

The overall learning process is:

```math
\mathrm{Search}(T)
\rightarrow
\mathrm{Verify}(T)
\rightarrow
\mathrm{Compare}(T)
\rightarrow
\mathrm{Learn}(T)
\rightarrow
\mathrm{Recover}(T)
```

---

# 28. What Makes This Different From Simply Doing RL?

| **Dimension**                 | **Conventional RL**              | **Proposed Formulation**                      |
| ----------------------------- | -------------------------------- | --------------------------------------------- |
| **Reward**                    | Often sparse terminal reward     | Final task reward plus local rewards for implementor when avaiable  |
| **State conditioning**        | State-centric                    | State + task specification                    |
| **Counterfactuals**           | Usually incidental               | Explicitly generated from the same $(s,T)$    |
| **Value function**            | Task-agnostic $V(s)$             | Task-conditioned $V_\psi(s,T)$                |
| **Search allocation**         | Uniform or policy-driven         | Value- and uncertainty-guided                 |
| **Policy learning**           | Learns from sampled trajectories | Learns from sibling counterfactual rollouts   |
| **Recovery**                  | Mixed into global trajectories   | Separate local and global recovery objectives |
| **Planning**                  | Often entangled with execution   | Dedicated Planner policy                      |
| **Implementation**            | Same policy handles everything   | Dedicated Implementor policy                  |
| **Deployment mismatch**       | Often implicit                   | Explicitly trained recovery capability        |

The central difference is therefore not simply the use of reinforcement learning.

The proposed framework changes:

* How trajectories are represented
* How training data is generated
* How alternatives are compared
* How search computation is allocated
* How planning and implementation are separated
* How recovery is learned

---

# 29. Experimental Hypotheses

### H1 — Task-Conditioned Value Function

$V_\psi(s,T)$ predicts downstream task success more accurately than a task-agnostic value estimator $V(s)$.

### H2 — Counterfactual Search

Multiple verified continuations from the same $(s,T)$ provide better policy-learning signals than isolated successful trajectories.

### H3 — Adaptive Search

Value-guided selective search achieves a better performance-to-compute tradeoff than uniform branching.

### H4 — Hierarchical Recovery

Combining local Implementor recovery with global Planner recovery outperforms either recovery mechanism alone.

### H5 — Search-to-Policy Transfer

Search-discovered behavior can be internalized into policy parameters, reducing inference-time search requirements.

---

# 30. Critical Ablations

1. SFT on successful trajectories
2. Standard RL / RLVR
3. RL with successful-trajectory replay
4. Branch-and-verify without task conditioning
5. Task-conditioned branch-and-verify
6. Task-conditioned branch-and-verify + counterfactual expansion
7. Counterfactual expansion + value-guided search
8. Search + Implementor recovery
9. Search + Planner recovery
10. Search + both recovery mechanisms
11. **Full proposed system**

---

# 31. The Central Research Claim

Long-horizon coding should be treated as **structured, task-conditioned search over repository states**, rather than an undifferentiated sequence of thousands of low-level actions.

The framework separates the problem into two major learning dimensions:

> **Forward Decisions:** Learned through verified task-conditioned counterfactual search.
>
> **Failure Recovery:** Learned through hierarchical recovery, with global strategy handled by the Planner and local implementation repair handled by the Implementor.

The complete learning process is:

```math
\mathrm{Search}(T)
\rightarrow
\mathrm{Verification}(T)
\rightarrow
\mathrm{Counterfactual\ Expansion}
\rightarrow
\mathrm{Policy\ and\ Value\ Learning}
\rightarrow
\mathrm{Recovery\ Learning}
```

The ultimate objective is:

> **Use expensive, task-conditioned, value-guided search during training to teach coding agents how to solve long-horizon software engineering tasks, so that inference requires substantially fewer rollouts and remains robust when execution deviates from the expected plan.**
