# Hierarchical Physical Intelligence for Robotics

## 1. The Core Idea

General-purpose robots need two complementary capabilities:

**Physical imagination** — predicting how the physical world could evolve toward a desired outcome.

**Real-world grounding** — continuously selecting actions based on what is actually happening.

Vision-Language-Action (VLA) models are strong at grounding actions in real observations, but must implicitly learn much of the physical reasoning required for long-horizon manipulation.

World Action Models (WAMs) provide a complementary capability: they can imagine future physical states and reason about the consequences of interactions. However, their expensive generative inference is better suited to **physical planning** than to high-frequency policy execution.

We propose a hierarchical architecture that separates these roles:

```math
\boxed{
\text{VLM}
\rightarrow
\text{3D World State}
\rightarrow
\text{3D Goal State}
\rightarrow
\text{WAM}
\rightarrow
\text{VLA}
}
```

The central principle is:

> **The WAM imagines how the world should evolve; the VLA continuously grounds that plan and executes it.**

The WAM is invoked at the beginning of a physical motion segment to generate a future trajectory toward an explicit goal state. The VLA then runs at high frequency, using that trajectory together with the current 3D state and observations to produce robot actions.

---

# 2. A 3D-Native Architecture

The system is built around an explicit 3D representation of the environment.

Let

```math
S_t
```

denote the current physical state.

It represents the relevant structure of the environment, including:

* object identity;
* geometry;
* 3D position and orientation;
* spatial relationships;
* contact relationships;
* object state.

The state is continuously updated from observations:

```math
S_t
=
f_{\mathrm{3D}}(S_{t-1},O_t).
```

The important design choice is that **3D state is the primary physical reasoning interface**.

Rather than forcing the planning system to reason entirely through pixels, the architecture explicitly represents the objects and spatial relationships that matter for manipulation.

This gives the system a persistent physical state:

```math
S_0
\rightarrow
S_1
\rightarrow
S_2
\rightarrow
\cdots
\rightarrow
S_t.
```

---

# 3. Hierarchical Task Decomposition

A high-level VLM converts a long-horizon instruction into a sequence of physical subtasks:

```math
T
\rightarrow
\{T_1,T_2,\ldots,T_N\}.
```

For example:

```text
Prepare the kitchen
        ↓
Put dishes in dishwasher
        ↓
Throw away trash
        ↓
Wipe the counter
        ↓
Place utensils on the table
```

Each subtask defines a relatively compact physical objective.

This separates:

```text
semantic reasoning
```

from:

```text
physical reasoning
```

and prevents the WAM from having to plan an entire household task as one enormous trajectory.

---

# 4. Explicit 3D Goal States

For each subtask, the system explicitly constructs the desired physical configuration.

Given the current state and subtask:

```math
S_t,\;T_i
```

a goal model generates:

```math
S_G^i
=
f_{\mathrm{goal}}(S_t,T_i).
```

The goal state describes **what the world should look like after the subtask is completed**.

For example:

```text
"Put the cup in the cabinet"
```

becomes a physical objective such as:

```math
S_t
\rightarrow
S_G
```

where the cup has a desired position and orientation relative to the cabinet.

This establishes an explicit separation:

```text
WHAT should happen?
        ↓
   Goal State

HOW can it happen?
        ↓
       WAM
```

The WAM therefore does not need to infer the physical objective while simultaneously planning the motion.

---

# 5. The WAM as a Physical Imagination Engine

The WAM receives the current state, the desired goal state, and the active subtask:

```math
(S_t,S_G^i,T_i).
```

It predicts a physically plausible sequence of future states:

```math
(S_t,S_G^i,T_i)
\rightarrow
(Z_1,Z_2,\ldots,Z_H).
```

where each

```math
Z_h
```

represents an imagined future world state.

Conceptually:

```text
Current state
      ↓
Approach object
      ↓
Establish contact
      ↓
Grasp
      ↓
Lift
      ↓
Move
      ↓
Place
      ↓
Goal state
```

The WAM is therefore acting as a **physical planner**.

Its output is not intended to directly replace the robot policy. It provides a physically informed prediction of how the subtask can unfold.

---

# 6. WAM Planning at Low Frequency

The WAM is computationally expensive compared with a policy model.

Rather than invoking it for every control step, the architecture uses it at the beginning of a motion-planning segment:

```math
(S_t,S_G,T_i)
\rightarrow
\text{WAM}
\rightarrow
Z_{1:H}.
```

The resulting imagined trajectory becomes context for the VLA.

The WAM therefore operates at a relatively low frequency, while the VLA operates at high frequency.

```math
\boxed{
\text{WAM}
=
\text{low-frequency physical planning}
}
```

```math
\boxed{
\text{VLA}
=
\text{high-frequency policy execution}
}
```

This separation allows each model to operate at the timescale appropriate to its role.

---

# 7. The VLA as the High-Frequency Grounding Layer

The VLA receives:

```math
(O_t,S_t,T_i,Z_{1:H})
```

and predicts the next robot action:

```math
A_t
=
f_{\mathrm{VLA}}(O_t,S_t,T_i,Z_{1:H}).
```

The VLA therefore has access to both:

```text
the current real world
```

and:

```text
the WAM's imagined future.
```

Its task is to continuously ground the imagined physical trajectory in the actual environment.

After execution:

```math
A_t
\rightarrow
O_{t+1}.
```

The 3D state is updated:

```math
S_{t+1}
=
f_{\mathrm{3D}}(S_t,O_{t+1}).
```

The VLA then produces the next action:

```math
A_{t+1}
=
f_{\mathrm{VLA}}(O_{t+1},S_{t+1},T_i,Z_{1:H}).
```

Thus the execution loop is:

```math
\boxed{
\text{WAM plan}
\rightarrow
\text{VLA}
\rightarrow
\text{real world}
\rightarrow
\text{3D state update}
\rightarrow
\text{VLA}
\rightarrow
\cdots
}
```

---

# 8. Planning and Execution as Different Timescales

The fundamental architectural separation is therefore:

```math
\boxed{
f_{\mathrm{VLM}}
<
f_{\mathrm{WAM}}
<
f_{\mathrm{VLA}}
<
f_{\mathrm{controller}}
}
```

The VLM makes relatively infrequent semantic decisions.

The WAM makes relatively infrequent physical planning decisions.

The VLA makes high-frequency policy decisions.

The robot controller operates at an even higher frequency.

This creates a natural computational hierarchy:

```text
VLM
semantic planning
        ↓
WAM
physical planning
        ↓
VLA
high-frequency grounding
        ↓
Controller
low-level control
```

---

# 9. Closed-Loop Grounding

The WAM's imagined trajectory is continuously grounded against the actual evolving state.

At every policy step:

```math
(O_t,S_t,Z_{1:H})
\rightarrow
A_t.
```

The resulting observation changes the state:

```math
S_t
\rightarrow
A_t
\rightarrow
S_{t+1}.
```

The VLA therefore continuously adapts its action to the current physical configuration.

For example, if an object is slightly displaced, a grasp is imperfect, or the robot's motion differs from the imagined trajectory, the VLA can adjust its next action using the newly observed 3D state.

The WAM provides the **longer-horizon physical context**.

The VLA provides the **high-frequency closed-loop adaptation**.

---

# 10. When the WAM Replans

The WAM does not need to be continuously regenerated during normal execution.

It can instead be invoked when the current physical plan is no longer sufficient.

For example:

```math
(S_{t'},S_G,T_i)
\rightarrow
\text{WAM}
\rightarrow
Z'_{1:H}.
```

This creates a hierarchy of responses:

```math
\text{VLA correction}
<
\text{WAM replanning}
<
\text{task replanning}.
```

Small deviations are handled by the VLA.

Changes that invalidate the current physical strategy trigger the WAM.

Changes to the overall task strategy trigger the high-level planner.

---

# 11. Why 3D Goal-State Planning?

Many manipulation tasks are naturally defined by physical configurations rather than visual appearances.

Examples include:

```text
cup inside cabinet
drawer fully closed
block on top of another block
object aligned with fixture
tool inserted into socket
```

These are naturally represented as transformations:

```math
S_t
\rightarrow
S_G.
```

A 3D-native system can therefore reason directly about the physical variables that define task success.

The WAM's problem becomes:

```math
\text{How can the world evolve from } S_t
\text{ to } S_G?
```

This is a much more explicit physical planning problem than simply predicting the next image or action.

---

# 12. Training Objectives

Each component specializes in a distinct problem.

### 3D World State

Learn to maintain a persistent physical representation:

```math
(S_{t-1},O_t)
\rightarrow
S_t.
```

### Goal Model

Learn to translate a semantic objective into a desired physical configuration:

```math
(S_t,T_i)
\rightarrow
S_G^i.
```

### WAM

Learn to predict physically plausible future states conditioned on the current and desired state:

```math
(S_t,S_G^i,T_i)
\rightarrow
Z_{1:H}.
```

### VLA

Learn to ground the imagined trajectory into high-frequency robot actions:

```math
(O_t,S_t,T_i,Z_{1:H})
\rightarrow
A_t.
```

The resulting specialization is:

```text
3D model
→ represent reality

Goal model
→ define desired reality

WAM
→ imagine the transition

VLA
→ execute the transition
```

---

# 13. Complete System

For a task `T`, the system operates as follows.

### Step 1 — Task decomposition

```math
T
\rightarrow
\text{VLM}
\rightarrow
T_i.
```

### Step 2 — Estimate the current 3D state

```math
S_{t-1},O_t
\rightarrow
S_t.
```

### Step 3 — Explicitly generate the goal state

```math
S_t,T_i
\rightarrow
S_G^i.
```

### Step 4 — Imagine the physical transition

```math
S_t,S_G^i,T_i
\rightarrow
\text{WAM}
\rightarrow
Z_{1:H}.
```

### Step 5 — Run the VLA at high frequency

```math
O_t,S_t,T_i,Z_{1:H}
\rightarrow
A_t.
```

### Step 6 — Observe the world

```math
A_t
\rightarrow
O_{t+1}.
```

### Step 7 — Update the 3D state

```math
S_t,O_{t+1}
\rightarrow
S_{t+1}.
```

### Step 8 — Continue high-frequency VLA execution

```math
O_{t+1},S_{t+1},T_i,Z_{1:H}
\rightarrow
A_{t+1}.
```

### Step 9 — Replan when needed

```math
S_{t'}
\rightarrow
\text{WAM}
\rightarrow
Z'_{1:H}.
```

### Step 10 — Verify the subtask

```math
(S_t,T_i)
\rightarrow
\text{VLM}
\rightarrow
\text{success / failure}.
```

Then:

```math
T_i
\rightarrow
T_{i+1}.
```

---

# 14. The Proposed Architecture

The complete architecture can be summarized as:

```math
\boxed{
\text{Task}
\rightarrow
\text{VLM}
\rightarrow
\text{3D State}
\rightarrow
\text{3D Goal}
\rightarrow
\text{WAM}
\rightarrow
\text{VLA}
\rightarrow
\text{3D State Update}
\rightarrow
\text{VLA}
}
```

with WAM replanning when required.

The three central design choices are:

### 1. 3D-native physical reasoning

The system maintains an explicit 3D representation of the physical world rather than relying exclusively on raw visual observations.

### 2. Explicit goal-state imagination

The system explicitly represents the desired future physical configuration:

```math
S_t
\rightarrow
S_G.
```

The WAM then reasons about the transition between these states.

### 3. Low-frequency WAM, high-frequency VLA

The WAM performs expensive physical imagination at the planning timescale.

The VLA continuously executes and grounds that plan at the policy timescale.

```math
\boxed{
\text{WAM}
\rightarrow
\text{physical plan}
\rightarrow
\text{VLA}
\rightarrow
\text{high-frequency execution}
}
```

---

# 15. Central Thesis

The central thesis is:

> **General-purpose robot intelligence should separate physical planning from high-frequency physical execution.**

WAMs and VLAs provide complementary capabilities.

```math
\text{WAM}
=
\text{physical imagination}
```

```math
\text{VLA}
=
\text{high-frequency grounding and execution}
```

The missing interface between them is an explicit representation of **where the world is now** and **where it should end up**.

That interface is a persistent 3D state and an explicit 3D goal state:

```math
\boxed{
S_t
\rightarrow
S_G
}
```

The WAM reasons about the transition.

The VLA continuously grounds that transition in the real world.

The resulting architecture is:

```math
\boxed{
\text{Understand}
\rightarrow
\text{Represent in 3D}
\rightarrow
\text{Define Goal}
\rightarrow
\text{Imagine with WAM}
\rightarrow
\text{Execute with VLA}
\rightarrow
\text{Update 3D State}
\rightarrow
\text{Repeat}
}
```

The fundamental shift is therefore:

> **The WAM becomes the robot's physical imagination and planning engine, while the VLA becomes its high-frequency grounding and execution engine.**

This allows expensive generative physical reasoning and fast reactive control to coexist in a single hierarchical architecture.
