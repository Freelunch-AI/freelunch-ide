# Hierarchical Physical Intelligence for Robotics: Combining VLM, WAMs and VLA

## 1. The Problem: General-Purpose Robotics Needs Both Physical Knowledge and Grounding

Current robot foundation models are approaching two complementary capabilities, but neither is sufficient on its own for robust general-purpose robotics.

**Vision-Language-Action (VLA) models** are strong at understanding instructions and directly interacting with the physical world. However, they are fundamentally reactive: they observe the current state and predict actions. Their physical knowledge is therefore largely embedded implicitly in the policy, making it difficult to reason explicitly about long-horizon physical consequences and novel situations.

**World Action Models (WAMs)** provide the complementary capability. By learning to predict how the world evolves under actions, they can acquire substantially richer knowledge of physical dynamics, object interactions, and long-horizon consequences. Recent systems such as DreamZero demonstrate that generative world models can already be used to imagine future robot trajectories and improve physical generalization.

But WAMs face a fundamental problem of their own:

> **A predicted trajectory is only as good as the world model that generated it.**

The real world will inevitably differ from the model's prediction. Small errors in object pose, friction, contact dynamics, perception, or robot execution accumulate over a long trajectory. Consequently, a trajectory that is physically plausible in the imagined world can progressively drift from the state of the real world.

This creates a fundamental tension:

```math
\text{WAM}
\rightarrow
\text{strong physical reasoning}
\rightarrow
\text{weak real-world grounding}
```

while:

```math
\text{VLA}
\rightarrow
\text{strong real-world grounding}
\rightarrow
\text{limited explicit physical reasoning}.
```

### Our hypothesis

The solution is not to make either model do everything.

Instead:

> **Use the WAM to reason about how a physical subtask could be accomplished, and use the VLA to continuously ground that plan in the real world.**

---

# 2. The Core Idea

We propose a hierarchical architecture in which different models operate at different abstraction levels and timescales:

```math
\text{Task Planner VLM}
\rightarrow
\text{3D State Model}
\rightarrow
\text{WAM}
\rightarrow
\text{VLA}
\rightarrow
\text{Robot Control}
```

The key separation is:

* **VLM:** determines *what needs to be accomplished*.
* **3D-State Model:** determines *what the world currently looks like*.
* **WAM:** determines *how the physical world could evolve to accomplish the current subtask*.
* **VLA:** determines *how to execute that intention in the actual world*.
* **Robot controller:** converts the VLA's actions into physical actuation.

The WAM therefore becomes a **physical planner**, rather than the complete robot policy.

The VLA becomes the **closed-loop grounding layer**, rather than having to independently discover all of the physical knowledge required for long-horizon reasoning.

---

# 3. Hierarchical Task Planning

Complex tasks should not be treated as one enormous physical trajectory.

A high-level VLM first decomposes the task:

```math
T
\rightarrow
\{T_1,T_2,\ldots,T_N\}.
```

For example:

```text
"Prepare the kitchen for dinner"
        │
        ├── put dishes in dishwasher
        ├── throw away trash
        ├── wipe counter
        └── place utensils on table
```

The planner handles **one subtask at a time**.

This dramatically reduces the physical planning horizon.

The WAM does not need to imagine the entire task from beginning to end. It only needs to solve the current physical objective.

After execution, the high-level planner verifies whether the subtask was successfully completed before moving on:

```math
T_i
\rightarrow
\text{execute}
\rightarrow
\text{verify}
\rightarrow
T_{i+1}.
```

This creates a natural hierarchy between:

```math
\text{long-horizon semantic reasoning}
```

and:

```math
\text{shorter-horizon physical planning}.
```

---

# 4. Persistent 3D World Representation

A central component of the architecture is a persistent representation of the physical world.

Let:

```math
S_t
```

denote the current 3D representation of the environment.

Rather than reconstructing the world from scratch at every timestep, a specialized 3D-State VLM updates the representation using the previous state and the new observation:

```math
S_{t-1}, O_t
\rightarrow
S_t.
```

Thus:

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

The representation should preserve persistent information while incorporating newly observed changes.

It should capture properties such as:

* object identity;
* 3D position;
* orientation;
* geometry;
* object relationships;
* contact relationships;
* state changes;
* newly observed regions;
* changes caused by previous actions.

This persistent 3D state becomes the common interface between the WAM and the VLA.

---

# 5. Why 3D Rather Than Just Images?

The WAM should reason about physical transformations in a representation that explicitly captures the structure of the world.

An RGB frame tells the model what the world looks like from one viewpoint.

A 3D representation instead provides a representation of:

```math
\text{objects}
+
\text{geometry}
+
\text{pose}
+
\text{spatial relationships}.
```

This makes it possible to specify physical goals directly.

For example, instead of saying:

> "Move the cup into the cabinet."

the system can represent the desired transformation as a change in the 3D state:

```math
S_t
\rightarrow
S_G.
```

where the cup occupies a desired pose relative to the cabinet.

The WAM can therefore reason about the transformation between physical states rather than merely predicting the next sequence of pixels.

---

# 6. Generating the 3D Goal State

For each subtask, the system generates a desired 3D configuration.

Given:

```math
S_t
```

and the subtask:

```math
T_i,
```

a goal-generation model produces:

```math
S_G^i
=
f_{\text{goal}}(S_t,T_i).
```

The WAM then receives both the current and desired states:

```math
(S_t,S_G^i,T_i).
```

This explicitly separates:

```math
\text{what should happen}
```

from:

```math
\text{how it should happen}.
```

The VLM defines the objective.

The WAM reasons about the physical transformation required to reach it.

---

# 7. The WAM as a Physical Planner

The WAM receives:

```math
S_t,\;S_G^i,\;T_i
```

and generates an imagined physical trajectory.

Conceptually:

```math
(S_t,S_G^i,T_i)
\rightarrow
(Z^V,Z^A).
```

where:

* $Z^V$ represents predicted latent visual/world states;
* $Z^A$ represents latent physical actions.

The WAM can therefore imagine:

```text
current state
      ↓
approach object
      ↓
establish contact
      ↓
grasp
      ↓
lift
      ↓
move
      ↓
place
      ↓
goal state
```

The important point is that the WAM is **not being asked to execute this trajectory exactly**.

It is being asked to provide a physically informed hypothesis about how the subtask can be accomplished.

---

# 8. Latent Physical Actions

The WAM should operate in a latent action space rather than being forced to predict the final robot action representation directly.

It produces:

```math
Z^A =
(z^A_1,z^A_2,\ldots,z^A_N).
```

Each latent action represents a higher-level physical intention.

For example, a latent action could encode concepts corresponding to:

```text
approach
grasp
lift
push
pull
place
maintain contact
```

without requiring the WAM itself to specify the exact robot-specific execution parameters.

The VLA then grounds these latent physical intentions into executable robot actions:

```math
Z^A
\rightarrow
A.
```

This provides an abstraction boundary between:

```math
\text{environment-level physical reasoning}
```

and:

```math
\text{robot-specific execution}.
```

---

# 9. The WAM Also Predicts the Consequences

The WAM does not only predict actions.

It predicts what should happen as a consequence of those actions.

Conceptually:

```math
Z^V_0
\rightarrow
Z^V_1
\rightarrow
\cdots
\rightarrow
Z^V_N.
```

while simultaneously predicting:

```math
Z^A_1
\rightarrow
Z^A_2
\rightarrow
\cdots
\rightarrow
Z^A_N.
```

Thus the WAM represents a coupled prediction:

```math
\text{action}
\rightarrow
\text{physical consequence}.
```

This is the source of its physical reasoning capability.

The model can effectively ask:

> "If I perform this sequence of physical interactions, what should happen?"

The resulting imagined trajectory provides the VLA with a physically informed target for execution.

---

# 10. The VLA as a Closed-Loop Grounding Model

The VLA receives the WAM's imagined trajectory together with the actual current state and observations.

Conceptually:

```math
(O_t,S_t,T_i,Z^V,Z^A)
\rightarrow
A_t.
```

It then executes the action in the real world.

The resulting observation is:

```math
A_t
\rightarrow
O_{t+1}.
```

The 3D-State Model updates:

```math
S_t,O_{t+1}
\rightarrow
S_{t+1}.
```

The VLA continues using the updated state.

This gives:

```math
\text{WAM imagination}
\rightarrow
\text{VLA execution}
\rightarrow
\text{real observation}
\rightarrow
\text{updated 3D state}.
```

The VLA therefore continuously corrects for discrepancies between the imagined trajectory and reality.

---

# 11. Why This Solves WAM Trajectory Drift

Suppose the WAM predicts:

```text
grasp → lift → move → place
```

but the real object is slightly different from what the model expected.

Perhaps:

* the object is 3 cm farther away;
* the grasp is slightly off;
* the object slips;
* the surface has different friction;
* the robot's motion differs from the predicted motion.

An open-loop WAM trajectory will progressively diverge from reality.

Our system instead treats the WAM trajectory as a **plan to be grounded**, not a trajectory that must be followed exactly.

The VLA continuously observes the consequences of its actions.

Therefore:

```math
\text{predicted state}
\neq
\text{actual state}
```

does not automatically imply failure.

The VLA adapts its execution to the actual state.

If the discrepancy becomes too large for local correction, the WAM can generate a new plan from the newly estimated state.

---

# 12. WAM Replanning Is Event-Driven

The WAM therefore does not need to run continuously at the VLA's control frequency.

For a subtask, the normal execution path is:

```math
S_t
\rightarrow
\text{WAM}
\rightarrow
(Z^V,Z^A)
\rightarrow
\text{VLA execution}.
```

The VLA handles normal deviations.

Only when the current plan becomes invalid does the WAM need to be invoked again:

```math
S_{t'}
\rightarrow
\text{WAM}
\rightarrow
(Z^{V'},Z^{A'}).
```

This makes WAM inference **event-driven rather than control-frequency-driven**.

The expensive generative world model can therefore operate relatively infrequently, while the VLA operates continuously.

---

# 13. Speed and Computational Efficiency

This separation also creates an important computational advantage.

Current WAM-based systems can repeatedly generate future trajectories during closed-loop execution.

If a task lasts $T$ seconds and replanning occurs every $k$ seconds:

```math
N_{\text{WAM}}
\approx
\frac{T}{k}.
```

Our architecture instead targets:

```math
N_{\text{WAM}}
\approx
1
```

for an uninterrupted subtask.

The VLA handles the high-frequency execution loop.

The WAM handles the low-frequency physical planning loop.

The high-level VLM handles the even lower-frequency task-planning loop.

Thus:

```math
\text{Task VLM}
<
\text{WAM}
<
\text{VLA}
<
\text{Robot Controller}
```

in terms of operating frequency.

This is a natural allocation of computation:

> **Expensive models reason slowly; reactive models execute quickly.**

---

# 14. Comparison With Current WAMs

The proposal is not simply "add a VLA after a WAM."

The key difference is the **role assigned to the WAM**.

|                               | Current WAM approach               | Proposed architecture          |
| ----------------------------- | ---------------------------------- | ------------------------------ |
| Primary role                  | World modeling + action generation | **Physical planning**          |
| Execution                     | WAM-centric                        | **VLA-centric**                |
| Planning horizon              | Potentially long                   | **Subtask-level**              |
| Real-world grounding          | WAM closed loop                    | **Dedicated VLA loop**         |
| Persistent 3D state           | Not central                        | **Central representation**     |
| Goal specification            | Usually implicit/trajectory-based  | **Explicit 3D goal state**     |
| Action representation         | Robot actions/action chunks        | **Latent physical actions**    |
| WAM frequency                 | Repeated during execution          | **Event-driven**               |
| Handling model–world mismatch | Replanning/control through WAM     | **VLA handles local mismatch** |
| High-level task decomposition | Not central                        | **Dedicated VLM**              |
| Subtask verification          | Not central                        | **Dedicated VLM**              |
| Long-horizon reasoning        | WAM trajectory                     | **Task VLM + subtask WAM**     |

The proposal therefore changes the fundamental relationship between the WAM and the robot.

Instead of:

```math
\text{WAM}
\rightarrow
\text{robot trajectory}
\rightarrow
\text{execution}.
```

we propose:

```math
\text{WAM}
\rightarrow
\text{physical plan}
\rightarrow
\text{VLA}
\rightarrow
\text{grounded execution}.
```

---

# 15. Comparison With VLAs

The same decomposition addresses the limitations of purely VLA-based systems.

A VLA must otherwise learn a mapping such as:

```math
(\text{observation},\text{instruction})
\rightarrow
\text{action}.
```

For increasingly complex tasks, this forces the model to implicitly learn:

* object dynamics;
* contact mechanics;
* long-horizon consequences;
* affordances;
* physical planning;
* recovery strategies.

Our architecture gives the VLA an explicit physical hypothesis from the WAM:

```math
(\text{observation},\text{task},\text{physical plan})
\rightarrow
\text{action}.
```

The VLA can therefore focus on the problem it is naturally suited for:

> **grounding an intended physical behavior into the continuously changing real world.**

This does not eliminate the need for physical knowledge in the VLA. Rather, it gives the VLA a much richer physical prior to condition on.

---

# 16. Hierarchical Recovery

The architecture naturally supports recovery at multiple levels.

### Local execution failure

The VLA handles small deviations:

```math
\text{execution error}
\rightarrow
\text{local correction}.
```

### Physical-plan failure

If the current strategy no longer works:

```math
\text{current state}
\rightarrow
\text{WAM}
\rightarrow
\text{new physical plan}.
```

### Subtask failure

If the subtask itself needs to change:

```math
T_i
\rightarrow
\text{Task Planner}
\rightarrow
\text{new subtask}.
```

### Task-level failure

If the overall strategy is invalid:

```math
T
\rightarrow
\text{Task Planner}
\rightarrow
\text{new task decomposition}.
```

This produces a hierarchy:

```math
\text{VLA recovery}
<
\text{WAM replanning}
<
\text{task-level replanning}.
```

Small errors are handled cheaply.

Large errors trigger progressively more expensive reasoning.

---

# 17. Training the System

Each component can be trained for a specialized objective.

### Task Planner VLM

Train it to:

```math
T
\rightarrow
\{T_1,\ldots,T_N\}
```

and to verify subtask completion:

```math
(S_t,T_i)
\rightarrow
P(T_i\text{ complete}).
```

### 3D-State Model

Train:

```math
(S_{t-1},O_t)
\rightarrow
S_t.
```

to maintain a persistent physical representation.

### Goal Model

Train:

```math
(S_t,T_i)
\rightarrow
S_G^i.
```

to translate semantic objectives into desired physical states.

### WAM

Train:

```math
(S_t,S_G^i,T_i)
\rightarrow
(Z^V,Z^A).
```

to model physical consequences and generate physically plausible plans.

### VLA

Train:

```math
(O_t,S_t,T_i,Z^V,Z^A)
\rightarrow
A_t.
```

to robustly ground the WAM's physical intentions into real-world execution.

This division allows each model to specialize instead of forcing a single model to simultaneously learn:

```text
language
+
task planning
+
3D reconstruction
+
physics
+
trajectory generation
+
robot control
+
recovery.
```

---

# 18. The Complete Architecture

For a task $T$, the complete system operates as follows.

### Step 1 — Task decomposition

```math
T
\rightarrow
\text{Task Planner}
\rightarrow
T_i.
```

### Step 2 — Maintain the current 3D state

```math
S_{t-1},O_t
\rightarrow
\text{3D-State Model}
\rightarrow
S_t.
```

### Step 3 — Generate the desired physical state

```math
S_t,T_i
\rightarrow
\text{Goal Model}
\rightarrow
S_G^i.
```

### Step 4 — Generate a physical plan

```math
S_t,S_G^i,T_i
\rightarrow
\text{WAM}
\rightarrow
(Z^V,Z^A).
```

### Step 5 — Ground the plan in reality

```math
O_t,S_t,T_i,Z^V,Z^A
\rightarrow
\text{VLA}
\rightarrow
A_t.
```

### Step 6 — Execute and observe

```math
A_t
\rightarrow
\text{real world}
\rightarrow
O_{t+1}.
```

### Step 7 — Update the world representation

```math
S_t,O_{t+1}
\rightarrow
S_{t+1}.
```

### Step 8 — Continue execution

The VLA continues grounding the physical plan against the updated state.

### Step 9 — Replan if necessary

If the current physical strategy becomes invalid:

```math
S_{t'}
\rightarrow
\text{WAM}
\rightarrow
\text{new plan}.
```

### Step 10 — Verify the subtask

```math
(S_t,T_i)
\rightarrow
\text{Task Planner}
\rightarrow
\text{completed / failed}.
```

If complete:

```math
T_i
\rightarrow
T_{i+1}.
```

The overall loop becomes:

```math
\text{Task}
\rightarrow
\text{Subtask}
\rightarrow
\text{Current 3D State}
\rightarrow
\text{3D Goal}
\rightarrow
\text{WAM Plan}
\rightarrow
\text{VLA Execution}
\rightarrow
\text{Updated 3D State}
\rightarrow
\text{Verification}
\rightarrow
\text{Next Subtask}.
```

---

# 19. Central Thesis

The central thesis is that **general-purpose robot intelligence should separate physical planning from physical execution**.

VLAs and WAMs provide complementary capabilities:

```math
\text{VLA}
=
\text{real-world grounding}.
```

```math
\text{WAM}
=
\text{physical imagination and planning}.
```

Neither capability should be forced to completely subsume the other.

Instead, we propose:

```math
\boxed{
\text{VLM}
\rightarrow
\text{3D State}
\rightarrow
\text{WAM}
\rightarrow
\text{VLA}
}
```

where:

* the **VLM** decomposes the task and verifies progress;
* the **3D-State Model** maintains a persistent representation of reality;
* the **WAM** imagines how the current physical subtask can be accomplished;
* the **VLA** continuously grounds that imagined plan in the real world.

The key conceptual shift is:

> **The WAM should not be the robot's entire policy. It should be the robot's physical imagination and planner.**

And:

> **The VLA should not have to discover all physical reasoning implicitly. It should be the closed-loop mechanism that turns physical plans into successful real-world behavior.**

This creates a hierarchy in which each model operates where its capabilities are most valuable:

```math
\text{Understand}
\rightarrow
\text{Decompose}
\rightarrow
\text{Represent}
\rightarrow
\text{Set Goal}
\rightarrow
\text{Imagine}
\rightarrow
\text{Ground}
\rightarrow
\text{Execute}
\rightarrow
\text{Verify}
\rightarrow
\text{Repeat}.
```

The result is a system that combines the **physical knowledge and long-horizon reasoning of WAMs** with the **continuous grounding and adaptability of VLAs**, while avoiding the computational cost and brittleness of requiring a generative world model to continuously control the robot.

---

# Clarification: What Does the VLA Output?

For readers less familiar with robot control stacks, the VLA's output should be understood as a **robot action**, rather than a raw electrical motor command.

In practice, these actions are typically expressed as declarative setpoints or desired physical outcomes, such as end-effector motion, pose changes, gripper commands, or desired contact forces.

For example:

```math
A_t =
(\Delta x,\Delta y,\Delta z,\Delta R,F_{\text{contact}}).
```

The robot's existing control stack then translates these commands into the forces and torques required to achieve them, and ultimately into actuator-level electrical signals:

```math
\text{VLA action}
\rightarrow
\text{robot controller}
\rightarrow
\text{forces/torques}
\rightarrow
\text{motor signals}.
```

The proposed learning architecture therefore operates primarily **above the hardware-specific low-level control layer**.
