# Hierarchical Physical Intelligence for Robotics

## 1. The Core Idea

General-purpose robots need two complementary capabilities: **physical imagination**, the ability to predict how the physical world could evolve toward a desired outcome, and **real-world grounding**, the ability to continuously select actions based on what is actually happening. Vision-Language-Action (VLA) models are strong at grounding actions in real observations, but must implicitly learn much of the physical reasoning required for long-horizon manipulation. A **World Model** provides a complementary capability by modeling how the physical state of the world evolves and allowing the robot to reason about the consequences of interactions. However, expensive world-model inference is better suited to physical planning than to high-frequency policy execution. We therefore propose a hierarchical architecture that explicitly separates semantic reasoning, physical planning, and policy execution:

```math id="gya79v"
\boxed{
\text{VLM}
\rightarrow
\text{3D World State}
\rightarrow
\text{3D Goal State}
\rightarrow
\text{World Model}
\rightarrow
\text{VLA}
}
```

The central principle is that the **World Model imagines how the world should evolve, while the VLA continuously grounds that plan and executes it**. The World Model is invoked during physical planning to predict a future latent trajectory toward an explicit goal state. The VLA then operates at high frequency, combining that imagined future with the current observations and 3D state to produce actions and continuously adapt to deviations between the imagined and real trajectories.

---

# 2. A 3D-Native Architecture

The system is built around an explicit 3D representation of the environment. Let

```math id="bjoowl"
S_t
```

denote the current physical state, containing the relevant structure for manipulation, including object identity, geometry, 3D position and orientation, spatial relationships, contact relationships, and object state. The state is continuously updated from observations according to

```math id="9wf98a"
S_t
=
f_{\mathrm{3D}}(S_{t-1},O_t).
```

The important design choice is that **3D state becomes the primary physical reasoning interface** rather than forcing the planning system to reason entirely through pixels. This provides a persistent representation of the environment that evolves over time:

```math id="hdleme"
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

Such a representation allows downstream models to reason directly about the objects and spatial relationships that matter for manipulation while separating persistent physical structure from the raw visual observations used to estimate it.

---

# 3. Hierarchical Task Decomposition

A high-level VLM converts a long-horizon instruction into a sequence of physical subtasks,

```math id="xczx27"
T
\rightarrow
\{T_1,T_2,\ldots,T_N\}.
```

For example, a task such as "Prepare the kitchen" may be decomposed into actions such as putting dishes in the dishwasher, throwing away trash, wiping the counter, and placing utensils on the table. Each subtask defines a relatively compact physical objective, allowing semantic reasoning to remain at the task level while physical reasoning is handled separately. This decomposition prevents the World Model from having to plan an entire household task as one enormous trajectory and instead gives it a well-defined local objective with a manageable physical horizon.

---

# 4. Explicit 3D Goal States

For each subtask, the system explicitly constructs the desired physical configuration. Given the current state and the active subtask,

```math id="wejhvq"
S_t,\;T_i
```

a goal model generates

```math id="11jtms"
S_G^i
=
f_{\mathrm{goal}}(S_t,T_i).
```

The goal state specifies **what the world should look like after the subtask is completed**. For example, "Put the cup in the cabinet" is translated into a physical objective in which the cup has a desired position and orientation relative to the cabinet:

```math id="8ibk64"
S_t
\rightarrow
S_G.
```

This establishes a clean separation between **what should happen** and **how it can happen**. The goal model defines the desired configuration, while the World Model reasons about the physical transition required to reach it. By explicitly representing the destination state, the planning system does not need to infer the physical objective while simultaneously searching for a feasible trajectory.

---

# 5. The World Model

The World Model is based on a **fine-tuned JEPA-style latent predictive architecture**. The key idea is to learn the structure of the physical world in latent space rather than modeling raw pixels directly. JEPA is particularly attractive because its self-supervised objective can predict masked latent representations across both **space and time**. Spatial masking forces the model to infer what occupies missing regions from surrounding context, while temporal masking forces it to infer how representations evolve across time. Unlike a pure next-frame prediction objective, this encourages the representation to capture spatial organization together with temporal structure during pretraining.

Conceptually, training follows:

```text
Observed video
      ↓
spatial / temporal masking
      ↓
latent prediction
      ↓
structured latent representation
```

This provides a natural combination of spatial understanding and temporal understanding in the learned latent space. The representation can capture objects, geometry, spatial relationships, object persistence, motion, contact, and scene structure without requiring all of these properties to be reconstructed explicitly in pixel space. The pretrained JEPA representation is then fine-tuned with robot interaction data so that the latent space becomes useful for predicting **action-conditioned physical evolution**.

---

# 6. World Model Inference

Although JEPA-style training uses both spatial and temporal prediction, the primary World Model operation during inference is **future latent-state prediction**. Given the current physical state, desired goal state, active subtask, and relevant action or trajectory context, the model predicts the future latent states that describe a physically plausible transition toward the goal:

```math id="s1w10b"
(S_t,S_G^i,T_i)
\rightarrow
(Z_1,Z_2,\ldots,Z_H).
```

The model can be viewed as learning an action-conditioned transition of the form

```math id="boiji0"
(S_t,a_t)
\rightarrow
Z_{t+1},
```

which can then be rolled forward through multiple imagined steps:

```math id="pdkz81"
Z_t,a_t
\rightarrow
Z_{t+1}
\rightarrow
Z_{t+2}
\rightarrow
\cdots
\rightarrow
Z_{t+H}.
```

The resulting latent trajectory represents an imagined physical evolution of the environment. For a manipulation task, this may correspond to approaching an object, establishing contact, grasping, lifting, moving, and finally reaching the desired configuration. The World Model therefore acts as a **physical imagination and planning engine**, using its learned latent representation to predict how the world can evolve rather than trying to generate a visually realistic future frame at every step.

---

# 7. World Model Planning at Low Frequency

World Model inference is more computationally expensive than policy execution, so the architecture does not require the model to run at the robot's highest control frequency. Instead, the World Model is invoked at the physical planning timescale to produce a future latent trajectory:

```math id="471swe"
(S_t,S_G,T_i)
\rightarrow
\text{World Model}
\rightarrow
Z_{1:H}.
```

That trajectory becomes context for the VLA, which executes the plan at a much higher frequency. The resulting hierarchy is therefore:

```math id="ss007o"
\boxed{
\text{World Model}
=
\text{low-frequency physical planning}
}
```

```math id="8bpwri"
\boxed{
\text{VLA}
=
\text{high-frequency policy execution}
}
```

This allows the World Model to spend computation on longer-horizon physical reasoning while the VLA focuses on the fast, reactive decision-making required during actual interaction.

---

# 8. The VLA as the High-Frequency Grounding Layer

The VLA receives the current observation and 3D state together with the subtask and the World Model's imagined future:

```math id="arj4qu"
(O_t,S_t,T_i,Z_{1:H})
```

and predicts the next robot action,

```math id="4j7hjf"
A_t
=
f_{\mathrm{VLA}}(O_t,S_t,T_i,Z_{1:H}).
```

The VLA therefore operates at the interface between imagination and reality. It has access to both the **current real world** and the **World Model's predicted future**, allowing it to continuously ground the longer-horizon physical plan in the actual state of the environment. After executing an action, the robot receives a new observation,

```math id="3lq8w5"
A_t
\rightarrow
O_{t+1},
```

and updates its 3D state,

```math id="crkzie"
S_{t+1}
=
f_{\mathrm{3D}}(S_t,O_{t+1}).
```

The VLA then produces the next action using the updated state and observation:

```math id="n3rp4m"
A_{t+1}
=
f_{\mathrm{VLA}}(O_{t+1},S_{t+1},T_i,Z_{1:H}).
```

The World Model therefore supplies longer-horizon physical context, while the VLA provides high-frequency closed-loop grounding and adaptation.

---

# 9. Planning and Execution as Different Timescales

The fundamental architectural separation is a hierarchy of models operating at different temporal scales:

```math id="76yw0c"
\boxed{
f_{\mathrm{VLM}}
<
f_{\mathrm{WorldModel}}
<
f_{\mathrm{VLA}}
<
f_{\mathrm{controller}}
}
```

The VLM makes relatively infrequent semantic decisions, the World Model makes relatively infrequent physical planning decisions, the VLA makes high-frequency policy decisions, and the low-level controller operates at an even higher frequency. This produces a natural computational hierarchy:

```text
VLM
semantic planning
        ↓
World Model
physical planning
        ↓
VLA
high-frequency grounding
        ↓
Controller
low-level control
```

The key benefit is that each component solves the problem at the temporal scale for which it is best suited rather than forcing a single model to perform semantic reasoning, physical planning, and low-level control simultaneously.

---

# 10. Closed-Loop Grounding

The World Model's imagined trajectory is continuously grounded against the actual evolving state of the robot and environment. At every policy step, the VLA receives the current state together with the imagined future,

```math id="exdusl"
(O_t,S_t,Z_{1:H})
\rightarrow
A_t.
```

The resulting action changes the physical state,

```math id="4pwm0k"
S_t
\rightarrow
A_t
\rightarrow
S_{t+1},
```

and the newly observed state is incorporated into subsequent policy decisions. If an object is slightly displaced, a grasp is imperfect, or the robot's motion differs from the imagined trajectory, the VLA can immediately adapt using the newly observed 3D state. This gives the architecture a clear division of responsibility: the World Model supplies the **longer-horizon physical prediction**, while the VLA supplies the **continuous real-world correction** required to execute that prediction robustly.

---

# 11. When the World Model Replans

The World Model does not need to be regenerated at every policy step. It is instead invoked again when the current physical plan is no longer sufficient. For example, after a significant deviation or an unexpected event, the updated state can be fed back into the World Model:

```math id="8zpb8e"
(S_{t'},S_G,T_i)
\rightarrow
\text{World Model}
\rightarrow
Z'_{1:H}.
```

This creates a hierarchy of responses in which small deviations are handled by VLA-level corrections, larger deviations trigger World Model replanning, and changes to the overall task strategy trigger a higher-level VLM decision:

```math id="p6eyk5"
\text{VLA correction}
<
\text{World Model replanning}
<
\text{task replanning}.
```

The result is a system that does not repeatedly pay the cost of full physical planning when local policy adaptation is sufficient, while still retaining the ability to perform deeper replanning when the current strategy becomes invalid.

---

# 12. Why 3D Goal-State Planning?

Many manipulation tasks are naturally defined by physical configurations rather than visual appearances. Examples include a cup inside a cabinet, a drawer fully closed, a block placed on top of another block, an object aligned with a fixture, or a tool inserted into a socket. These objectives are naturally represented as transformations between physical states:

```math id="xzwkrv"
S_t
\rightarrow
S_G.
```

A 3D-native system can therefore reason directly about the variables that define task success. The World Model's problem becomes:

```math id="trys54"
\text{How can the world evolve from } S_t
\text{ to } S_G?
```

This is a more explicit physical planning problem than simply predicting the next image or next action. The JEPA-based latent representation provides a compact learned space in which these future physical configurations can be represented and predicted while preserving the spatial structure acquired during pretraining.

---

# 13. Training Objectives

Each component specializes in a distinct problem. The 3D model learns to maintain a persistent representation of the physical world,

```math id="c2t77j"
(S_{t-1},O_t)
\rightarrow
S_t.
```

The goal model learns to translate a semantic objective into a desired physical configuration,

```math id="yd38pw"
(S_t,T_i)
\rightarrow
S_G^i.
```

The World Model is first pretrained using JEPA-style self-supervised prediction, where spatially and temporally masked latent representations must be inferred from visible context. This encourages the representation to encode both **where things are** and **how they evolve over time**. The model is then fine-tuned on robot interaction data to learn action-conditioned latent dynamics:

```math id="boiji0"
(S_t,a_t)
\rightarrow
Z_{t+1}.
```

Finally, the VLA learns to convert the imagined latent trajectory into high-frequency robot actions:

```math id="auyg8z"
(O_t,S_t,T_i,Z_{1:H})
\rightarrow
A_t.
```

The resulting specialization is therefore:

```text
3D model
→ represent reality

Goal model
→ define desired reality

World Model
→ imagine the transition

VLA
→ execute the transition
```

---

# 14. Complete System

For a task `T`, the system first decomposes the instruction into a physical subtask,

```math id="u8lrdj"
T
\rightarrow
\text{VLM}
\rightarrow
T_i.
```

The robot then estimates its current 3D state from the latest observation,

```math id="69ao5r"
S_{t-1},O_t
\rightarrow
S_t,
```

and generates the desired goal configuration,

```math id="594g3k"
S_t,T_i
\rightarrow
S_G^i.
```

The World Model then predicts a future latent trajectory toward that goal,

```math id="87a11x"
S_t,S_G^i,T_i
\rightarrow
\text{World Model}
\rightarrow
Z_{1:H}.
```

The VLA executes that plan at high frequency,

```math id="12edps"
O_t,S_t,T_i,Z_{1:H}
\rightarrow
A_t,
```

after which the robot observes the result,

```math id="tnu4mq"
A_t
\rightarrow
O_{t+1},
```

and updates the 3D state,

```math id="5c7beq"
S_t,O_{t+1}
\rightarrow
S_{t+1}.
```

The VLA continues execution using the updated state and the same imagined trajectory,

```math id="btv5y5"
O_{t+1},S_{t+1},T_i,Z_{1:H}
\rightarrow
A_{t+1},
```

until the current plan is no longer sufficient, at which point the World Model replans:

```math id="w7nzcu"
S_{t'}
\rightarrow
\text{World Model}
\rightarrow
Z'_{1:H}.
```

Once execution reaches the desired configuration, the task can be verified and the system proceeds to the next subtask.

---

# 15. The Proposed Architecture

The complete architecture can be summarized as:

```math id="gxgi2r"
\boxed{
\text{Task}
\rightarrow
\text{VLM}
\rightarrow
\text{3D State}
\rightarrow
\text{3D Goal}
\rightarrow
\text{World Model}
\rightarrow
\text{VLA}
\rightarrow
\text{3D State Update}
\rightarrow
\text{VLA}
}
```

with World Model replanning when required. The proposal rests on four central design choices.

### 1. 3D-native physical reasoning

The system maintains an explicit 3D representation of the physical world rather than relying exclusively on raw visual observations. This makes object geometry, spatial relationships, and physical configurations explicit interfaces for downstream reasoning.

### 2. JEPA-based World Model

The World Model is initialized through JEPA-style self-supervised learning in latent space. Unlike pure next-frame prediction, JEPA-style training predicts masked representations across both **space and time**, encouraging the model to acquire spatial structure together with temporal understanding. During inference, the primary operation is future latent-state prediction, allowing the model to roll forward imagined physical trajectories without requiring full pixel-level video generation at every step.

### 3. Explicit goal-state imagination

The desired future physical configuration is explicitly represented:

```math id="mdaxhi"
S_t
\rightarrow
S_G.
```

The World Model then reasons about the transition between these states instead of having to infer the objective while simultaneously planning.

### 4. Low-frequency World Model, high-frequency VLA

The World Model performs expensive physical imagination at the planning timescale, while the VLA continuously executes and grounds the imagined plan at the policy timescale:

```math id="nlvj9g"
\boxed{
\text{World Model}
\rightarrow
\text{physical plan}
\rightarrow
\text{VLA}
\rightarrow
\text{high-frequency execution}
}
```

---

# 16. Central Thesis

The central thesis is that **general-purpose robot intelligence should separate physical planning from high-frequency physical execution**. The World Model and VLA provide complementary capabilities: the World Model performs physical imagination and planning, while the VLA performs high-frequency grounding and execution. The missing interface between them is an explicit representation of **where the world is now** and **where it should end up**.

That interface is provided by a persistent 3D state and an explicit 3D goal state:

```math id="tql33i"
\boxed{
S_t
\rightarrow
S_G
}
```

The World Model reasons about the transition between those states, while the VLA continuously grounds the predicted transition in the real world. The resulting architecture is:

```math id="pe8m11"
\boxed{
\text{Understand}
\rightarrow
\text{Represent in 3D}
\rightarrow
\text{Define Goal}
\rightarrow
\text{Imagine with World Model}
\rightarrow
\text{Execute with VLA}
\rightarrow
\text{Update 3D State}
\rightarrow
\text{Repeat}
}
```

The fundamental shift is therefore:

> **The World Model becomes the robot's physical imagination and planning engine, while the VLA becomes its high-frequency grounding and execution engine.**

The World Model uses JEPA-style representation learning to acquire spatial and temporal structure, then uses latent future-state prediction to imagine how that structured world can evolve toward an explicit 3D goal. This separates **learning what the world is and how its structure evolves** from **learning how to control the robot through that world**, creating a hierarchical architecture in which semantic reasoning, physical planning, and high-frequency execution each have a clearly defined role.
