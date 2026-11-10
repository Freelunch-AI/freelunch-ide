# Lunch Robotics: Universal Brain for Any Robot

## The Vision

We are building a **universal robot brain** that can turn a general-purpose robot into an autonomous worker inside an arbitrary physical environment. The fundamental problem in robotics is not simply making robots capable of moving or manipulating objects. It is giving them enough general intelligence to understand an unfamiliar world, reason about what should be done, predict the consequences of possible actions, adapt to the physics and embodiment of that world, dynamically decompose and replan long-horizon tasks, and continuously improve as they encounter situations that were not represented in their original training data.

Lunch Robotics separates this problem into **general physical intelligence, world simulation, adaptive high-level planning, hierarchical task execution, environment-specific intelligence, pre-deployment robustness, and continual fleet learning**, with each layer contributing a distinct capability.

General robot intelligence is learned offline through a **Foundation Model Data Funnel** that combines massive amounts of human video with progressively more robot-relevant forms of supervision. Massive human video provides broad knowledge about the physical world. A generative World Model can then produce large quantities of synthetic egocentric human manipulation experience that is filtered and action-annotated. Human action supervision connects world knowledge to purposeful behavior, manipulation data collected with a data-collecting gripper introduces contact and embodiment information, and a smaller amount of teleoperation data provides direct grounding in robot control.

The World Model itself is trained through three complementary phases. First, it learns the structure and evolution of the physical world from massive passive video through future and spatial prediction. Second, it is fine-tuned on synthetic multi-environment data to learn **3D-conditioned temporal prediction**, allowing it to use explicit scene geometry and noisy 3D observations when available. Third, it is further fine-tuned with **action-conditioned training**, where actions are explicitly provided as inputs and the model learns to predict what will happen after those actions. This progressively transforms the World Model from a passive predictor, into a 3D-conditioned temporal model, and finally into a learned **action-conditioned simulator** capable of evaluating hypothetical robot behaviors.

The synthetic egocentric-data stage is deliberately positioned between passive foundation-model pretraining and dense action-conditioned training. Once the foundation World Model is sufficiently capable, it can generate large numbers of hypothetical human manipulation episodes covering tasks, objects, viewpoints, and physical situations that are expensive to collect in the real world. It can generate completely new egocentric trajectories and can also take real human videos and modify their initial frames with an image-editing model to create controlled variations of the scene, object configuration, or task. Humans and learned critics filter these generations for plausibility and usefulness, while pose-estimation models recover approximate hand and body actions from the surviving videos. These generated examples become an additional source of action-conditioned supervision, allowing the model to learn controllable physical dynamics at much greater scale than direct real-world collection alone would permit.

This is not treated as a closed synthetic-data loop: **real egocentric data remains the grounding source**. The synthetic generation process is also used to discover where the foundation model is weak and determine which new real-world experiences should be collected to correct those weaknesses.

These datasets are not treated as isolated sequential stages in which each dataset replaces the previous one. They jointly train a shared **World Model + VLA foundation brain**, with the training distribution becoming increasingly robot-relevant while retaining the broad coverage of the lower-fidelity datasets. Once this foundation has learned to perceive, predict, and generate grounded robot actions, **Stage 7 trains the high-level VLM planner, not the VLA**, through supervised fine-tuning on high-quality planning and reasoning trajectories followed by reinforcement learning. The planner learns to reason about the task, current world state, execution history, constraints, and existing plan in order to select the right next subgoal, identify missing prerequisites, dynamically create useful intermediate subtasks, and revise the hierarchical plan when necessary. Its reward is tied to the quality of these planning decisions and their downstream task outcomes, rather than to the generation of low-level actions.

The VLA has a separate responsibility: it learns to generate robot action chunks from the current observation, active subtask, and selected target future. Candidate actions are evaluated through the action-conditioned World Model and VLM comparison loop. The VLA is further adapted to a particular robot and environment through Stage 8 agentic simulation RL. At inference time, the **planner's reasoning budget** can be selected according to task complexity, uncertainty, and available compute.

After the global reasoning-trained planner and action-capable foundation system have been built, they are adapted to a particular robot and environment through a final **agentic simulation RL stage** before deployment. A VLM training agent acts as an autonomous curriculum and scenario-generation manager, deciding which tasks and subtasks to train, how difficult they should be, how the environment should be varied, and from which initial states the robot should start. The curriculum progresses from small, simple tasks toward increasingly complex and long-horizon tasks while simultaneously varying the environment and initial state. This teaches the robot not only to solve expected tasks, but to remain robust when objects, geometry, physics, clutter, or other agents differ from the nominal environment, and when the robot encounters a partially completed, degraded, or otherwise unexpected task state that requires recovery.

At runtime, the system is **hierarchical and dynamically replanning**. A high-level VLM planner maintains a structured plan for the user's task, identifies the current subtask, and continuously updates the plan as new observations arrive. It can keep the current subtask, modify it, reorder subtasks, remove obsolete steps, or **create entirely new subtasks dynamically** when execution reveals that an intermediate action or recovery step is required. This planner operates at a slower semantic timescale than the physical control loop. Under the current subtask, the World Model imagines desirable target futures, the VLM selects one, and the VLA repeatedly proposes action chunks whose consequences are simulated and scored against that target.

The final action-selection mechanism is therefore not a simple VLA-to-robot mapping. The system operates on **three nested timescales**:

1. **Highest-level semantic planning:** the high-level VLM planner maintains and dynamically replans the hierarchical task structure and can create new subtasks online.
2. **Lower-frequency target planning:** the World Model imagines three possible medium-horizon future trajectories representing different ways the current subtask could successfully progress. A VLM scores those futures and selects the most desirable one as the current **target future**.
3. **Higher-frequency action selection:** the VLA generates three candidate action chunks. Each candidate is passed once through the action-conditioned World Model to produce one predicted future trajectory. A VLM compares those three predicted futures against the selected target future and chooses the candidate whose consequence best matches the target. The selected action chunk is executed, and the high-frequency loop repeats while the target remains active.

Conceptually, the architecture is:

```text
HIGHEST-LEVEL: DYNAMIC TASK PLANNING
─────────────────────────────────────
User Task
   ↓
High-Level VLM Planner
   ↓
Hierarchical Task / Subtask Plan
   ↓
Current Subtask
   ↓

LOW-FREQUENCY TARGET PLANNING
──────────────────────────────
World Model
   ↓
3 Imagined Target Futures
   ↓
VLM Target Scoring
   ↓
Best Target Future
   │
   │ remains active for multiple
   │ high-frequency action decisions
   ▼

HIGH-FREQUENCY ACTION SELECTION
────────────────────────────────
VLA
 ↓
3 Candidate Action Chunks
 ↓
World Model
 ↓
3 Predicted Future Trajectories
 ↓
VLM Comparison Against Target
 ↓
Best Action Chunk
 ↓
Robot
 ↓
New Observation
 ↓
Repeat

       │
       │ periodically or when execution
       │ invalidates the current plan/target
       ▼
Dynamic Replanning
```

The target-selection stage is:

```math
\left\{z_{\mathrm{target}}^{(1)}(t:t+H),z_{\mathrm{target}}^{(2)}(t:t+H),z_{\mathrm{target}}^{(3)}(t:t+H)\right\}
\sim
\mathrm{WM}(s_t,T_{\mathrm{subtask}})
```

The VLM scores each target:

```math
r_j^{\mathrm{target}}=\mathrm{VLM}_{\mathrm{score}}\left(z_{\mathrm{target}}^{(j)},s_t,T_{\mathrm{subtask}}\right)
```

and selects:

```math
j^*=\arg\max_{j\in\{1,2,3\}}r_j^{\mathrm{target}}
```

The selected target is:

```math
z_{\mathrm{target}}=z_{\mathrm{target}}^{(j^*)}
```

The VLA then samples three candidate action chunks:

```math
\left\{a_t^{(1)},a_t^{(2)},a_t^{(3)}\right\}\sim\pi_{\theta}\left(a_t\mid o_t,T_{\mathrm{subtask}},r_{1:k},z_{\mathrm{target}},S_{\mathrm{tutorial}}\right)
```

Each candidate is passed through the action-conditioned World Model once:

```math
z_{\mathrm{pred}}^{(i)}(t:t+H)=\mathrm{WM}_{\mathrm{action}}(s_t,a_t^{(i)})
```

The VLM scores the three predicted futures against the selected target:

```math
r_i=\mathrm{VLM}_{\mathrm{score}}\left(z_{\mathrm{pred}}^{(i)},z_{\mathrm{target}},s_t,T_{\mathrm{subtask}}\right)
```

and selects:

```math
i^*=\arg\max_{i\in\{1,2,3\}}r_i
```

The executed action is:

```math
a_t=a_t^{(i^*)}
```

The important distinction is therefore:

```text
Dynamic high-level plan
       ↓
Current subtask
       ↓
3 target futures
       ↓
VLM
       ↓
1 selected target future
       ↓
3 candidate actions
       ↓
3 predicted futures
       ↓
VLM
       ↓
1 selected action
```

There are **three candidate futures at the action-selection stage, not nine**. The VLA proposes what could be done, the World Model predicts what would happen, and the VLM determines both which imagined future is desirable and which proposed action most closely produces that future. The high-level planner adds another semantic layer: it determines **what should be accomplished next**, can dynamically decompose a task into subtasks, and can change that decomposition whenever execution reveals new information.

The reasoning-trained planner and action-capable foundation system are then adapted to a specific robot and physical environment by an autonomous **Real-to-Sim Agent**. Given a walkthrough of the environment, tactile probing data, task context, demonstrations, and the robot's hardware specification, the agent constructs and calibrates a digital twin, performs system identification, determines which aspects of the environment should be represented through explicit physics and which should be handled by learned dynamics, and trains a fast surrogate simulator. The resulting calibrated simulator becomes the training environment for Stage 8 agentic RL, where a VLM training agent constructs the task curriculum and generates environment and initial-state variations. The output is an **environment-specific robot system**, including an environment-adapted VLA and the supporting planner and World Model components.

The system does not stop learning after deployment. Each deployed robot continuously improves locally through simulation, while the platform monitors real-world execution for mistakes. A failure may be identified explicitly by the human user or automatically by a VLM observing the robot through cameras installed in the environment. Every detected mistake is captured as a rich episode containing the task and subtask being executed, the robot state and action trajectories, the high-level planner's reasoning trace and plan revisions, the active hierarchical plan and planner state, any dynamically generated subtasks, the three target trajectories sampled during the most recent target update, the VLM scores used to select the target, the three candidate actions considered by the VLA, the World Model's predicted future for each candidate, the VLM scores used to rank those predicted futures against the target, the selected action, the actual execution, the corresponding decoded imagined videos, and the relevant simulator state and environment parameters.

These failures are **not automatically used to update the global model**. Instead, they are sent back to the Lunch Robotics team, where we analyze error modes across deployments and determine which failures are genuinely generalizable. Environment-specific quirks remain local, while problems that reveal broader capability gaps become the basis for carefully curated training datasets. Those curated datasets are then used to improve a new version of the **global lab system**, with the central training process applying additional supervised fine-tuning, planner reasoning supervision, reinforcement learning, World Model training, VLM evaluation training, or a combination of them depending on the identified capability gap. The updated lab system is then combined with knowledge accumulated in the environment-specific systems through continual weight mixing, producing a new mixed initialization that is redistributed to every deployment, after which each environment performs its own local RL again starting from this stronger initialization.

This creates a compounding global-local learning system:

```math
\boxed{\text{Foundation Data}\rightarrow\text{World Model}\rightarrow\text{Synthetic Egocentric Action Data}\rightarrow\text{3D-Conditioned World Model}\rightarrow\text{Action-Conditioned World Model}\rightarrow\text{Global VLA + Planner}\rightarrow\text{Planner Reasoning SFT + RL}\rightarrow\text{Environment-Specific Agentic RL}\rightarrow\text{Deployment}\rightarrow\text{Failure Discovery}\rightarrow\text{Error-Mode Analysis}\rightarrow\text{Curated Data}\rightarrow\text{New Lab System}\rightarrow\text{Weight Mixing}\rightarrow\text{All Deployments}}
```

The key product idea is therefore simple:

> **The robots specialize locally, while Lunch Robotics learns globally.**

---

# 1. The Architecture

The Lunch Robotics architecture is built around four nested learning systems operating on top of a shared foundation model and a reasoning-trained high-level planner. The first is the **global foundation model**, which learns broad physical, manipulation, predictive, and action-generation intelligence. The second is the **global planner-training system**, which teaches the high-level VLM to decompose tasks, choose subgoals, and dynamically revise plans. The third is the **environment-specific adaptation system**, which takes the shared system and specializes it to a particular robot and physical environment. The fourth is the **fleet learning system**, which turns deployment failures into curated global training data and feeds the resulting improvements back into the shared brain. At the center of the architecture is the interaction between the VLA, the World Model, the VLM evaluator, and the persistent high-level planner.

The **high-level VLM planner** is a persistent agent rather than a one-shot task decomposer. It maintains a hierarchical representation of the current task, chooses which subtask should be attempted next, and continually replans when new observations invalidate assumptions. The planner can **create new subtasks dynamically**, not just execute a predefined list. This is essential for physical autonomy because important intermediate objectives often become apparent only after the robot interacts with the world. Stage 7's Reasoning SFT + RL explicitly trains this planner behavior.

For example, a request such as `Clean the table` may initially produce:

```text
1. Identify objects on table
2. Collect objects
3. Wipe table
4. Return objects
```

During execution, the robot may discover that an object blocks the cleaning area. The planner can dynamically insert:

```text
3a. Move obstructing object
```

Or it may discover that the object is fragile:

```text
3a. Identify safe grasp region
3b. Move fragile object to safe location
3c. Resume cleaning
```

Or it may discover that the original route is impossible:

```text
Original:
   Pick object → place in drawer

Observation:
   Drawer blocked

Replan:
   Move blocker
       ↓
   Open drawer
       ↓
   Pick object
       ↓
   Place object
```

The planner therefore operates as an **online hierarchical controller over semantic state**, while the lower-level loops operate over physical state.

The World Model provides two distinct capabilities. First, it can imagine multiple possible desirable futures for the current subtask. Second, because it is action-conditioned, it can predict the consequence of each proposed action. The architecture operates at two physical-control frequencies plus one semantic-planning frequency: the **high-level planner** updates when the task structure needs revision, the **target-future generation and target-future VLM scoring loop runs at lower frequency**, and the **VLA action generation, World Model prediction, and candidate VLM scoring loop runs at higher frequency**.

The VLA therefore does not need to directly solve:

```math
(o_t,T)\rightarrow a_t^*
```

Instead, the overall system decomposes the problem into:

```math
(T,s_t,h_t)\rightarrow\text{Dynamic Plan}\rightarrow T_{\mathrm{subtask}}\rightarrow\{z_{\mathrm{target}}^{(1)},z_{\mathrm{target}}^{(2)},z_{\mathrm{target}}^{(3)}\}\rightarrow z_{\mathrm{target}}\rightarrow\{a_t^{(1)},a_t^{(2)},a_t^{(3)}\}\rightarrow\{z_{\mathrm{pred}}^{(1)},z_{\mathrm{pred}}^{(2)},z_{\mathrm{pred}}^{(3)}\}\rightarrow a_t^*
```

The complete nested control loop is:

```text
                    HIGHEST-LEVEL LOOP
       ┌────────────────────────────────────────┐
       │                                        │
User Task → High-Level VLM Planner             │
                 ↓                             │
      Dynamic Hierarchical Plan                │
                 ↓                             │
           Current Subtask                     │
                 │                             │
                 ▼                             │
              LOW-FREQUENCY LOOP               │
                 │                             │
            World Model                        │
                 ↓                             │
         3 Target Futures                      │
                 ↓                             │
         VLM Target Scoring                    │
                 ↓                             │
        Selected Target Future                 │
                 │                             │
                 ▼                             │
           HIGH-FREQUENCY LOOP                 │
                 │                             │
                 VLA                            │
                 ↓                             │
        3 Candidate Actions                    │
                 ↓                             │
         World Model × 3                       │
                 ↓                             │
        3 Predicted Futures                    │
                 ↓                             │
         VLM Comparison                        │
                 ↓                             │
           Best Action                         │
                 ↓                             │
         Execute Action Chunk                  │
                 ↓                             │
            Observation                        │
                 │                             │
                 └─────────────────────────────┘
                        │
                        │ replan when needed
                        ▼
               Updated Task / Subtask Plan
```

A target can remain fixed over several action-chunk decisions:

```math
z_{\mathrm{target}}^{(j)}\qquad\text{for}\qquad k_j\leq k<k_{j+1}
```

where `j` indexes target updates and `k` indexes the higher-frequency action decisions. This prevents the system from unnecessarily regenerating and semantically rescoring a new target future after every small action chunk. The target is recomputed when the robot has progressed sufficiently, when the current target becomes inconsistent with the observed state, when the environment changes materially, when uncertainty rises enough to justify replanning, or when the high-level planner changes the active subtask.

---

# 2. Foundation Model Data Funnel

The foundation model and pre-deployment training process is organized as an **eight-stage data and learning funnel** in which the amount of available data decreases as fidelity, action grounding, reasoning, and environment specificity increase. The first six stages build a globally capable World Model + VLA foundation from massive human video through increasingly direct action supervision. Stage 7 trains the **high-level VLM planner** through reasoning SFT and outcome-driven RL. Stage 8 takes the resulting global robot system into a calibrated simulation of the specific environment and robot in which it is expected to operate, where **agentic reinforcement learning exposes the system to a curriculum of increasingly complex tasks, systematic environment variations, and varied initial states that teach robustness and recovery before real-world deployment**.

The eight stages are:

```text
                        ┌────────────────────────────┐
                        │ STAGE 8                    │
                        │ AGENTIC ENVIRONMENT RL     │
                        │ Curriculum + robustness   │
                        │ + recovery before reality  │
                        └─────────────┬──────────────┘
                                      │
                        ┌─────────────┴──────────────┐
                        │ STAGE 7                    │
                        │ PLANNER REASONING SFT + RL│
                        │ Dynamic subgoals + replans │
                        └─────────────┬──────────────┘
                                      │
                        ┌─────────────┴──────────────┐
                        │ STAGE 6                    │
                        │ TELEOPERATION              │
                        │ Direct robot actions       │
                        └─────────────┬──────────────┘
                                      │
                        ┌─────────────┴──────────────┐
                        │ STAGE 5                    │
                        │ HUMAN + GRIPPER DATA       │
                        │ Manipulation + contact     │
                        └─────────────┬──────────────┘
                                      │
                        ┌─────────────┴──────────────┐
                        │ STAGE 4                    │
                        │ 3D-CONDITIONED WORLD MODEL │
                        │ Synthetic multi-env data   │
                        └─────────────┬──────────────┘
                                      │
                        ┌─────────────┴──────────────┐
                        │ STAGE 3                    │
                        │ SYNTHETIC HUMAN ACTION     │
                        │ Scalable action supervision│
                        └─────────────┬──────────────┘
                                      │
                        ┌─────────────┴──────────────┐
                        │ STAGE 2                    │
                        │ HUMAN VIDEO + ACTION       │
                        │ Approximate action data    │
                        └─────────────┬──────────────┘
                                      │
                 ┌────────────────────┴──────────────────┐
                 │ STAGE 1                               │
                 │ MASSIVE HUMAN EGOCENTRIC VIDEO        │
                 │ Broadest world coverage               │
                 └───────────────────────────────────────┘

                        ↓    ↓    ↓    ↓    ↓    ↓    ↓

                   GLOBAL FOUNDATION TRAINING

                             ↓

                   WORLD MODEL + VLA BRAIN

                             ↓

                STAGE 7: HIGH-LEVEL PLANNER SFT + RL

                             ↓

                STAGE 8: ENVIRONMENT SPECIALIZATION

                             ↓

                   REAL-WORLD DEPLOYMENT
```

The first seven stages form the **global foundation and planner-training system**:

| Stage | Data | Scale | Supervision | Primary Capability |
| --- | --- | ---: | --- | --- |
| **1** | Massive human egocentric video | Massive | Next-video-latent prediction | World representation and physical prediction |
| **2** | Human video + VLM pose waypoints | Large | Approximate action supervision | Human action understanding |
| **3** | Synthetic human egocentric data | Very large | Synthetic approximate action supervision | Scalable action-conditioned coverage |
| **4** | Synthetic multi-environment video + explicit 3D scene data + noise conditions | Large | 3D-conditioned next-video-latent prediction | 3D-conditioned world dynamics |
| **5** | Human manipulation + data-collecting gripper | Medium | Human manipulation action supervision | Contact and manipulation |
| **6** | Teleoperation data | Small | Direct robot action supervision | Robot control and embodiment |
| **7** | High-level planner reasoning traces + RL rollouts | Targeted | Planner SFT + outcome-driven RL | Dynamic hierarchical planning, subgoal selection, and subtask creation |

Stage 8 is different because its data is **generated specifically for the expected deployment**:

| Stage | Data | Scale | Supervision | Primary Capability |
| --- | --- | ---: | --- | --- |
| **8** | Agentic RL across curriculum + environment and initial-state variations | Massive | Environment reward + interaction | Environment-specific competence, robustness, and recovery |

The distinction is fundamental:

```text
Stages 1–6
       ↓
GLOBAL WORLD MODEL + VLA FOUNDATION

Stage 7
       ↓
REASONING-TRAINED HIGH-LEVEL VLM PLANNER

Stage 8
       ↓
ENVIRONMENT-SPECIFIC ROBOT SYSTEM

Real Deployment
       ↓
CONTINUAL LOCAL + GLOBAL LEARNING
```

Stage 8 is therefore not another foundation-data layer in the same sense as human video or teleoperation. It is the **final pre-deployment specialization stage** in which the global system is turned into a policy and planner stack that is competent and robust in the particular robot and physical world where it is expected to operate.

---

# 2.1 The Data Pyramid

The data funnel has both a **scale dimension** and a **grounding dimension**. Data becomes smaller and more expensive as it moves toward direct robot action, but the lower stages are intentionally retained because they contain broad information that the highest-fidelity datasets cannot economically reproduce.

```text
                       HIGH FIDELITY / LOW SCALE
                                 ▲
                                 │
                    Planner VLM Reasoning SFT + RL
                                 │
                          Teleoperation
                                 │
                      Human + Gripper
                                 │
                   Synthetic Human Actions
                                 │
                     Human + Waypoints
                                 │
                    Massive Human Video
                                 │
                                 ▼
                   LOW FIDELITY / HIGH SCALE
```

The critical addition is that the synthetic human stage is **generated using both the foundation World Model and real-video seeds**, while the next stage adds explicit 3D conditioning through synthetic multi-environment temporal prediction. Together these create a controlled bridge from broad passive world knowledge to increasingly structured physical prediction and dense action supervision without pretending that generated data contains novel real-world grounding.

Stage 8 then sits orthogonally above this pyramid as the **environment-specific specialization layer**:

```text
GLOBAL DATA PYRAMID
       ↓
Global World Model + VLA
       ↓
Planner VLM Reasoning SFT + RL
       ↓
Global Robot System with Reasoning-Trained Planner
       ↓
Calibrated Deployment Environment
       ↓
VLM Curriculum Agent
       ↓
Agentic Environment RL
       ↓
Robust Environment-Specific Robot System
```

---

# 2.2 Stage 1 — Massive Human Egocentric Video

Massive human egocentric video provides the broadest source of physical-world information. The starting point is the best available **JEPA-like World Model**, which is trained or fine-tuned to predict future observations while also learning to predict spatially masked regions of the environment in latent space. The goal is not simply to reproduce pixels, but to learn a useful internal representation of objects, spatial relationships, motion, interactions, contact-relevant structure, and temporal evolution that can later support both future imagination and action-conditioned prediction.

```mermaid
flowchart LR
    A["Massive Human Egocentric Video"] --> B["JEPA-like World Model"]
    B --> C["Future Prediction"]
    B --> D["Spatial Prediction"]
    C --> E["Latent World Representation"]
    D --> E
```

At this point, the World Model is primarily learning the structure of the world from passive observations. It answers questions such as what is likely to happen next, what the hidden part of the scene looks like, how objects are related in space and time, and how the physical state evolves. This creates the latent world representation that later becomes the basis for action-conditioned simulation.

---

# 2.3 Stage 2 — Human Video + VLM Pose Waypoints

The second data layer connects broad world understanding to purposeful action. A VLM processes large-scale human egocentric videos and extracts approximate pose and action waypoints, providing weak but scalable action supervision without requiring robot data. At the same time, the World Model produces latent future imagination from the observed video, allowing the VLA to learn from both the observed action trajectory and the predicted evolution of the environment.

```mermaid
flowchart TD
    A["Human Egocentric Video"] --> B["VLM"]
    B --> C["Human Pose / Action Waypoints"]
    A --> D["World Model"]
    D --> E["Latent Future Imagination"]
    A --> F["VLA"]
    C --> F
    E --> F
```

The model begins learning a relationship of the form:

```math
\text{Observation}
+
\text{Task}
+
\text{Predicted Future}
\rightarrow
\text{Action}
```

This creates a bridge between passive physical-world understanding and purposeful interaction.

---

# 2.4 Stage 3 — Synthetic Human Egocentric Action Data

Once the foundation World Model has learned a sufficiently broad model of human behavior, objects, scenes, and physical evolution, it becomes a **generative source of additional egocentric manipulation experience**. The system generates this experience through two complementary mechanisms. First, the World Model can generate entirely new human egocentric trajectories from a task, scene, or manipulation specification. Second, the system can take the initial frames of real human egocentric videos and use a powerful **image-editing model** to create controlled visual variations of those real scenes before handing the modified initial conditions to the World Model. These edits can vary object identity, object appearance, object position, scene configuration, clutter, background conditions, and other visually controllable properties while preserving the basic structure of the observed real-world situation.

```mermaid
flowchart TD
    A["Real Human Egocentric Video"] --> B["Initial Frames"]
    B --> C["Image Editing Model"]
    C --> D["Variation 1"]
    C --> E["Variation 2"]
    C --> F["Variation N"]
    G["Task / Scene Specification"] --> H["Foundation World Model"]
    D --> H
    E --> H
    F --> H
    H --> I["Synthetic Egocentric Trajectories"]
    J["Foundation World Model"] --> K["Fully Generated Trajectories"]
    I --> L["Human / VLM Filtering"]
    K --> L
    L --> M["Pose Estimator"]
    M --> N["3D Hand / Body Actions"]
    N --> O["Synthetic Action Dataset"]
    O --> P["Action-Conditioned World Model"]
    O --> Q["VLA Training"]
    L --> R["Failure / Quality Analysis"]
    R --> S["Capability Gap Map"]
    S --> T["Targeted Real Data Collection"]
    T --> J
```

The real-initial-frame branch is particularly valuable because it anchors synthetic generation to **actual camera statistics, real object appearances, real spatial configurations, real hand morphology, and real environmental layouts**. Instead of asking the World Model to hallucinate an entire scene from scratch, the system can start from a real observation and make small, controlled changes before generating the subsequent trajectory.

The image-editing branch can also create **small variations of the task itself**. A video of a person placing a red cup on a shelf could, for example, be transformed into variants involving a different cup, a slightly different shelf position, a different nearby object, or a slightly changed placement objective. The objective is not to generate arbitrary tasks that drift far from reality, but to create a dense neighborhood of related examples around genuine real-world observations.

The two synthetic branches complement each other:

```text
                    SYNTHETIC HUMAN DATA
                             │
               ┌─────────────┴─────────────┐
               │                           │
               ▼                           ▼
      Fully generated trajectories   Real-video-seeded
                                     image variations
               │                           │
               └─────────────┬─────────────┘
                             ▼
                    World Model rollout
                             ↓
                  Human / VLM filtering
                             ↓
                     Pose estimation
                             ↓
                  Action annotations
```

The first filtering layer is deliberately human-led. Humans inspect generated trajectories and reject examples containing obvious physical inconsistencies, implausible contacts, temporal artifacts, bad hand geometry, contradictory object behavior, or other failures that would otherwise reinforce errors in the model. Over time, these human judgments can be distilled into learned critics and VLM-based quality filters, allowing the system to scale while maintaining a human-controlled quality bar.

The surviving videos are then processed with pose-estimation and reconstruction models to recover approximate hand and body trajectories:

```math
V_{\mathrm{synthetic}}
\rightarrow
\mathrm{PoseEstimator}
\rightarrow
A_{\mathrm{human}}
```

where `A_human` can contain hand pose, wrist motion, body pose, temporal waypoints, and other available action representations. The resulting synthetic demonstrations are mixed with real human action data and increasingly grounded manipulation data to train the World Model and VLA.

The critical principle is that **synthetic data provides scale but not independent grounding**. The foundation World Model acquired its understanding from real observations, so its generations primarily sample and recombine knowledge already present in the model. Real-video-seeded editing adds another useful property: because synthetic generation starts from actual observations, the resulting examples remain anchored to real scene structure while varying selected factors.

The synthetic generation process also becomes a **diagnostic instrument for the foundation model**. The team can track which kinds of generations humans consistently reject, where pose reconstruction fails, which manipulation categories produce low-quality trajectories, and which task or object combinations lead to repeated inconsistencies.

This creates an active data-acquisition loop:

```text
Real Egocentric Data
       ↓
Foundation World Model
       ↓
Synthetic Egocentric Generation
       │
       ├── New generations from the World Model
       │
       └── Real-video initial frames
               ↓
         Image editing
               ↓
         Controlled variations
       │
       └──────────────┐
                      ▼
               Human / Learned Filtering
                      ↓
                Failure Analysis
                      ↓
             Identify Weaknesses
                      ↓
         Targeted Real Data Collection
                      ↓
             Improved Foundation WM
                      ↺
```

The purpose of this loop is broader than synthetic-data augmentation. The World Model becomes both a **data generator and a diagnostic model of its own uncertainty and blind spots**.

---

# 2.5 Stage 4 — 3D-Conditioned World Model Fine-Tuning

Stage 4 adds an explicit **3D-conditioned learning capability** to the foundation World Model. The model is fine-tuned on large amounts of synthetic video generated from multiple types of simulated environments, where every training example can provide both the recent visual history and an explicit 3D representation of the environment. The objective remains temporal prediction: given previous frames, predict the next frame, but now the World Model can also condition its prediction on structured 3D information when that information is available.

The synthetic environments should cover diverse geometries, object configurations, interaction layouts, viewpoints, and physical conditions so that the model learns a general relationship between visual observations, explicit 3D scene state, and temporal evolution rather than memorizing one simulator.

```text
Previous Frames
     +
3D Environment Representation
     +
Noise / Observation Condition
     ↓
3D-Conditioned World Model
     ↓
Predicted Next Frame
```

The 3D input can contain whatever explicit scene information is available from the simulator and later from real-world perception, such as geometry, depth, point clouds, object poses, spatial relationships, or other structured representations of the scene. The important capability is not dependence on a single 3D format, but learning that **explicit geometric state can be used as an additional conditioning signal for temporal prediction**.

The training objective can be expressed conceptually as:

```math
\hat{o}_{t+1}\sim P_{\phi}\left(o_{t-k:t},G_t,n_t\right)
```

where `o_{t-k:t}` denotes the recent sequence of observed frames, `G_t` is the available 3D representation of the environment, and `n_t` is a noise or observation-condition variable describing the rendering or sensing conditions under which the sequence was generated.

The **noise factor** is deliberately varied during training to prevent the model from learning only clean simulator dynamics. The simulation can introduce controlled visual and sensing imperfections such as sensor noise, depth corruption, missing observations, motion blur, lighting variation, rendering artifacts, imperfect reconstruction, or other disturbances that approximate the conditions encountered by a real perception stack. The World Model therefore learns to use 3D information without assuming that either the video or the explicit scene representation is perfectly clean.

A crucial property of this stage is that **3D conditioning is optional**. Training can mix examples where the 3D representation is available with examples where it is absent or degraded. This teaches the World Model to exploit explicit geometry when it exists while retaining the ability to predict from visual history alone.

```mermaid
flowchart TD
    A["Synthetic Simulated Environments"] --> B["Rendered Frame Sequences"]
    A --> C["Explicit 3D Scene State"]
    A --> D["Noise / Sensor Conditions"]
    B --> E["3D-Conditioned World Model"]
    C --> E
    D --> E
    E --> F["Predicted Next Frame / Future Latents"]
    F --> G["Temporal Prediction Loss"]
    H["3D Available"] --> E
    I["3D Missing / Degraded"] --> E
```

This stage is intentionally **not yet action-conditioned**. The model is still learning to predict how the observed world evolves from one frame to the next. The new capability is that it can use explicit 3D scene state to disambiguate geometry, occlusion, spatial relationships, and other aspects of the physical world that may be difficult to infer from a short visual history alone.

The role of Stage 4 is therefore to establish the progression:

```text
Passive Video Understanding
       ↓
Synthetic Human Action Coverage
       ↓
3D-Conditioned Temporal Prediction
       ↓
Action-Conditioned Prediction
       ↓
Contact + Robot Action Grounding
```

The same capability can later be used by the action-conditioned World Model. Once actions become inputs, the model can combine **visual history + explicit 3D state + action** when available, giving the simulator a stronger representation of the physical configuration in which an action is being evaluated.

---

# 2.6 Action-Conditioned World Model Training

The World Model, now strengthened by the Stage 4 3D-conditioned temporal prediction objective, is exposed to increasingly large collections of action-labelled trajectories from real human video, synthetic human video, manipulation datasets, and teleoperation. The action-conditioned model can continue to use explicit 3D scene state when it is available, while remaining capable of operating from visual state alone. Instead of learning only passive prediction, the model learns to predict how the world changes after an action.

```math
z_{t+1}\sim P_{\phi}(z_{t+1}\mid s_t,G_t,a_t)
```

and, when 3D state is unavailable:

```math
z_{t+1}\sim P_{\phi}(z_{t+1}\mid s_t,a_t)
```

Over longer horizons:

```math
z_{t:t+H}\sim P_{\phi}(s_t,a_{t:t+H-1})
```

The key conceptual transition is:

```text
Passive Observation
       ↓
"What normally happens next?"
       ↓
Action-Conditioned Prediction
       ↓
"What happens if I do this?"
```

The model can now be queried with hypothetical action sequences and asked to predict their consequences. Synthetic human action trajectories are particularly useful here because they provide large-scale coverage of action variations, while the real-world stages remain essential for correcting the model where its predictions diverge from reality.

This transforms the World Model into a learned **action-conditioned simulator**. The model does not need to produce a perfect pixel-level simulation of the future; it needs to accurately predict the latent aspects of the future that matter for decision-making, including object motion, contact state, task progress, spatial relationships, physical consequences, and other task-relevant changes.

---

# 2.7 Stage 5 — Human Manipulation + Data-Collecting Gripper

The fifth layer introduces substantially more direct information about manipulation and contact. Humans perform manipulation tasks using a **data-collecting gripper** that records end-effector motion and interaction signals, giving the model access to information that is difficult to recover from ordinary video alone, including grasping, pushing, pulling, contact transitions, deformation, slip, and other manipulation dynamics.

```mermaid
flowchart LR
    A["Human Manipulation + Data-Collecting Gripper"] --> B["Contact + Motion + Action Data"]
    B --> C["Shared World Model / VLA Training"]
```

This dataset is useful for both sides of the system: for the VLA, it provides better action supervision; for the World Model, it provides better supervision for learning how actions transform physical states. The same action data therefore helps the system learn both what action should be taken and what will happen if that action is taken.

---

# 2.8 Stage 6 — Teleoperation

Teleoperation provides the highest-fidelity foundation-model supervision because demonstrations are generated directly by real robots. It provides true robot action distributions, embodiment-specific kinematics, actuator constraints, interaction dynamics, temporal structure, and realistic observation-action correlations, making it the most direct source of robot grounding in the foundation model.

```mermaid
flowchart TD
    A["Teleoperation"] --> B["Direct Robot Action Supervision"]
    B --> C["Shared VLA Training"]
    B --> D["Action-Conditioned World Model"]
```

The same teleoperation trajectories therefore serve two complementary purposes: the VLA learns which actions are associated with successful behavior, while the World Model learns what happens after those actions are executed.

---

# 2.9 Stage 7 — High-Level Planner Reasoning SFT + RL

The earlier stages teach the shared foundation system **what the world looks like, how it evolves, how humans manipulate it, how robot actions affect it, and how to predict the consequences of those actions**. Stage 7 trains a different capability: **how the high-level VLM planner should reason about a task and determine what the robot should accomplish next**.

The target of Stage 7 is the persistent **high-level VLM planner**, not the VLA. The planner operates at the semantic task level. It receives the user's task, the current observed or estimated world state, the execution history, the current hierarchical plan, available skill knowledge, and relevant constraints. It produces or updates the hierarchical plan, selects the active subtask, and can create new intermediate subgoals when the current plan is incomplete or no longer valid.

The first part of Stage 7 is supervised fine-tuning on high-quality planning trajectories. These examples teach the planner to:

- Decompose long-horizon tasks into useful hierarchical subtasks.
- Select the most appropriate next subgoal given the current state.
- Identify prerequisites that were not apparent when the task was first planned.
- Recognize when an existing subtask is impossible, obsolete, unsafe, or already complete.
- Insert new subtasks and recovery steps when observations reveal that they are necessary.
- Reorder, modify, or replace the remaining plan as the environment changes.
- Use available skills, task context, and execution history when deciding what should happen next.

The planner's output is a **plan and a semantic subgoal**, not a robot action chunk:

```math
(P_t,g_t)\sim\pi_{\psi}\left(\cdot\mid T,s_t,H_t,P_{t-1},\mathcal{K}\right)
```

where `P_t` is the updated hierarchical plan, `g_t` is the selected next subgoal, `T` is the user's task, `s_t` is the current world state, `H_t` is execution history, `P_{t-1}` is the previous plan, and `K` is available skill and task knowledge.

After SFT, the planner is further optimized with reinforcement learning using full or partial task-execution rollouts. It is rewarded for reasoning that leads to appropriate subgoal choices and useful plan updates, including discovering missing prerequisites, creating necessary intermediate subtasks, choosing effective recovery plans, and making progress toward successful task completion. Rewards should reflect the downstream quality of the plan—not merely whether the planner produced a plausible-looking reasoning trace.

Conceptually, the RL objective is:

```math
\max_{\psi}\mathbb{E}_{\tau\sim\pi_{\psi}}\left[R_{\mathrm{planning}}(\tau;T,s_0)\right]
```

where the planning reward evaluates the quality of the planner's subgoal decisions and plan revisions in the context of the resulting task execution. Relevant signals include task progress and completion, whether selected subgoals were necessary and achievable, whether important prerequisites were identified, whether new subtasks were created when needed, and whether replanning led to safe and effective recovery.

The VLA is **not** the Stage 7 reasoning-RL target. Once the planner selects the current subgoal, the VLA performs the lower-level action-generation role: it receives the observation, active subtask, selected target future, and retrieved tutorial context, and proposes candidate action chunks. The action-conditioned World Model predicts their consequences, and the VLM evaluation loop selects among them. This preserves a clear division of responsibilities:

```text
HIGH-LEVEL VLM PLANNER — TRAINED IN STAGE 7
Task + World State + Execution History + Current Plan
                         ↓
              Reasoning SFT + RL
                         ↓
       Updated Hierarchical Plan / Next Subgoal
                         ↓

VLA — ACTION-GENERATION COMPONENT
Current Observation + Active Subtask + Target Future
                         ↓
             Candidate Action Chunks
                         ↓

WORLD MODEL + VLM ACTION-EVALUATION LOOP
Predict Candidate Consequences → Compare → Select Action
```

Stage 7 therefore produces a **reasoning-trained high-level VLM planner** that can dynamically decompose and replan tasks. It does not turn the VLA into a chain-of-thought action proposer. The VLA continues to learn action generation from action supervision and environment-specific RL, while the planner learns which objectives and subgoals the robot should pursue.

---

# 2.10 Stage 8 — Agentic RL in the Expected Environment

After the global robot system has been trained—including the action-capable World Model + VLA foundation and the reasoning-trained high-level VLM planner—the model is specialized to the particular robot and environment through **massive agentic reinforcement learning inside the calibrated deployment simulator**.

Stage 8 follows the same fundamental mechanism as the local online RL loop used after deployment. The robot learns through interaction, task complexity is increased through curriculum learning, and the training environment is continuously varied. The key difference is that Stage 8 happens **before real-world deployment**, allowing the robot to acquire environment-specific competence, robustness, and recovery ability in simulation.

The central idea is:

> **The VLM training agent constructs the worlds and situations the robot learns in; the robot system learns how to operate within them.**

The Stage 8 architecture is:

```text
Global Robot System
(World Model + VLA + Reasoning-Trained High-Level VLM Planner)
                         ↓
                 Calibrated Digital Twin
                         ↓
                    VLM Training Agent
                         ↓
       Task / Subtask Selection + Environment Variation
                         +
                  Initial-State Variation
                         ↓
                    Simulation Episode
                         ↓
                       Agentic RL
                         ↓
                 Performance Evaluation
                         ↓
                    VLM Training Agent
                         ↺
```

The VLM training agent acts as an **autonomous simulation curriculum manager**. It receives the environment description, calibrated simulator, robot capabilities, available task set, current policy performance, and known failure modes. It then decides what the robot should experience next.

Conceptually, the VLM training agent chooses:

```math
(T,E,s_0,h_0)\sim\mathcal{G}_{\mathrm{VLM}}\left(E_{\mathrm{real}},\mathcal{T},\Pi,\mathcal{F}\right)
```

where:

- `T` is the selected task or task family;
- `E` is the environment configuration;
- `s_0` is the initial physical state;
- `h_0` represents the initial hierarchical task state or plan context;
- `E_real` is the calibrated deployment environment;
- `𝒯` is the available task set;
- `Π` represents the current robot system and its performance;
- `𝓕` represents known failure modes.

The VLM training agent therefore controls four dimensions of the training distribution:

```text
1. WHAT should the robot learn?
       ↓
   Task / capability

2. WHAT should it accomplish next?
       ↓
   Current subtask / decomposition opportunity

3. WHAT WORLD should it learn it in?
       ↓
   Environment configuration / variation

4. FROM WHAT STATE should it start?
       ↓
   Initial-state variation / recovery scenario
```

The high-level planner is also trained and evaluated to handle **dynamic decomposition**. It does not need the entire subtask tree to be specified in advance. During a training episode it may discover that a new intermediate objective is necessary, and the scenario is scored partly on whether the planner can recognize the need, create the appropriate subgoal, and guide the robot toward successful execution.

---

# 2.11 VLM-Generated Curriculum

The Stage 8 curriculum follows the same basic principle as online local RL: the robot progresses from **small tasks to increasingly larger and more complex tasks**.

```text
Simple Primitive
       ↓
Small Task
       ↓
Multi-Step Task
       ↓
Multi-Object Task
       ↓
Long-Horizon Task
       ↓
Full Deployment Task
       ↓
Complex Multi-Stage Autonomy
```

Early training may begin with isolated primitives such as grasping, placing, pushing, opening, or navigating to a specific location. Once those behaviors become reliable, the curriculum combines them into increasingly complex sequences, eventually requiring the robot to complete the full tasks it is expected to perform in deployment.

The VLM training agent monitors performance and decides when to increase difficulty:

```text
Current Robot System
       ↓
VLM Training Agent evaluates performance
       ↓
Reliable on current curriculum?
       ├── No → generate more examples at current difficulty
       └── Yes
              ↓
        Increase difficulty
              ↓
        Generate next curriculum
```

The curriculum is therefore adaptive rather than a fixed script. The VLM training agent repeatedly generates scenarios around the current capability boundary, concentrating simulation compute where it is most useful.

---

# 2.12 Dynamic Hierarchical Planning and Subtask Creation

The high-level VLM planner maintains a hierarchical task state rather than a frozen sequence of instructions. Given the task objective and the current world state, it proposes an initial plan, selects the active subtask, and monitors execution against the expected state transitions. **This is the capability explicitly trained in Stage 7: planner SFT teaches the VLM how to construct and revise useful hierarchical plans, while planner RL rewards decisions that select appropriate subgoals and create effective new subtasks as the task unfolds.**

```math
P_t=\mathrm{Planner}\left(T,S_t,H_t,\mathcal{K},F_t\right)
```

where `P_t` is the current plan, `T` is the task, `S_t` is the current world state, `H_t` is execution history, `K` is available skill knowledge, and `F_t` is the set of known failures or constraints.

The planner does not merely choose the next pre-authored subtask. It can update the plan:

```text
Current Plan
   ↓
New Observation / Outcome
   ↓
VLM Planner
   ↓
┌──────────────────────────────┐
│ Keep current subtask         │
│ Modify current subtask       │
│ Reorder subtasks             │
│ Remove obsolete subtask      │
│ Insert a new subtask         │
│ Create a new recovery plan   │
│ Replace the remaining plan   │
└──────────────────────────────┘
   ↓
Updated Plan
```

A dynamically created subtask can itself become a complete mini-task with its own target-future and action-selection loops. For example:

```text
Main Task: Prepare dinner

Initial plan:
1. Gather ingredients
2. Prepare food
3. Cook food
4. Plate food

New observation:
Ingredient is inside sealed package

Dynamic subtask inserted:
2a. Find safe opening point
2b. Open package
2c. Transfer ingredient

Continue with:
2d. Prepare food
```

This makes the architecture explicitly **hierarchical and closed-loop through continuous replanning**. The robot can discover what needs to be done rather than assuming that every useful intermediate action was known beforehand.

Replanning can be triggered by:

```text
├── New object or obstacle discovered
├── Expected object state not observed
├── Action failure
├── Unexpected contact
├── Another agent changes the environment
├── Current subtask becomes impossible
├── Current subtask is already completed
├── A new prerequisite is discovered
├── Safety constraint becomes active
└── Uncertainty rises beyond a threshold
```

The high-level planner therefore changes **what the robot is trying to accomplish**, while the lower-frequency target planner changes **which desirable future it is trying to reach**, and the high-frequency action loop changes **how it moves toward that future right now**.

---

# 2.13 Environment Variation

For every curriculum level, the VLM training agent creates multiple variations of the expected deployment environment. The objective is to prevent the environment-specific robot system from simply memorizing the exact configuration reconstructed by the Real-to-Sim system.

```text
                   EXPECTED ENVIRONMENT
                          │
                 ┌────────┴────────┐
                 │                 │
                 ▼                 ▼
            Nominal Scene      Variations
                                   │
               ┌───────────────────┼───────────────────┐
               │                   │                   │
               ▼                   ▼                   ▼
            Geometry            Objects             Physics
               │                   │                   │
         positions             identity              mass
         orientation            appearance           friction
         clutter                availability         compliance
                                                     contact
               │                   │                   │
               └───────────────────┼───────────────────┘
                                   │
                       ┌───────────┴───────────┐
                       │                       │
                       ▼                       ▼
                   Sensing                 Dynamic Events
                       │                       │
                 lighting                human movement
                 occlusion               object movement
                 sensor noise             interruptions
```

Relevant variations can include:

```text
STATIC CONDITIONS
├── Object positions
├── Object orientation
├── Object identity
├── Object appearance
├── Object quantity
├── Clutter
├── Lighting
├── Camera viewpoint
└── Scene configuration

PHYSICAL CONDITIONS
├── Mass
├── Friction
├── Compliance
├── Damping
├── Contact parameters
├── Actuator response
└── Sensor noise

TASK CONDITIONS
├── Goal specification
├── Initial state
├── Object availability
├── Task ordering
├── Partial completion
└── Time constraints

DYNAMIC CONDITIONS
├── Object displacement
├── Object occlusion
├── Human interruption
├── Moving agents
├── Unexpected contact
├── Tool failure
├── Grasp failure
└── Environmental disturbances
```

The system does not randomize everything uniformly. The VLM training agent reasons about what variations are relevant to the actual deployment environment and generates them accordingly.

---

# 2.14 Initial-State Variation and Recovery Training

Environment variation alone is not sufficient for robust autonomy. A robot can be trained across many different scenes while still implicitly assuming that every task begins from a clean, ideal state. Stage 8 therefore also varies the **initial task state**.

The VLM training agent can intentionally initialize the environment and robot in states that are partially completed, unusual, degraded, or inconsistent with the nominal demonstration.

```text
Nominal Initial State

Cup on table
   ↓
Pick up cup
   ↓
Place cup
```

versus:

```text
Recovery Initial State

Cup already partially lifted
   ↓
Robot must infer current state
   ↓
Continue manipulation
```

or:

```text
Cup has slipped
   ↓
Robot detects unexpected state
   ↓
Planner replans
   ↓
Create recovery subgoal
   ↓
Recover
   ↓
Complete task
```

or:

```text
Task is already 60% complete
   ↓
Robot enters midway
   ↓
Infer current state
   ↓
Continue from there
```

The purpose is to teach the robot that **it is not always starting from scratch**.

---

# 2.15 Curriculum Over Task Complexity, Environment, and Initial State

The complete Stage 8 curriculum has **four simultaneous axes**:

```text
                        STAGE 8 CURRICULUM
                               │
      ┌────────────────────────┼────────────────────────┐
      │                        │                        │
      ▼                        ▼                        ▼
 TASK COMPLEXITY         ENVIRONMENT STATE       INITIAL STATE
      │                        │                        │
small → large            nominal → varied       nominal → degraded
      │                        │                        │
      └────────────────────────┼────────────────────────┘
                               │
                        DYNAMIC PLANNING
                               │
                       static → adaptive
                               │
                               ▼
                          AGENTIC RL
```

The planner dimension teaches the robot that the task structure itself can change:

```text
Fixed decomposition
       ↓
State-dependent decomposition
       ↓
Dynamic subtask insertion
       ↓
Dynamic recovery planning
       ↓
Open-ended hierarchical autonomy
```

This produces training scenarios such as:

```text
Level 1
Simple task
+ nominal environment
+ clean initial state
+ fixed/simple decomposition

Level 2
Simple task
+ environment variation
+ clean initial state
+ modest replanning

Level 3
Multi-step task
+ environment variation
+ mildly perturbed initial state
+ dynamic subtask creation

Level 4
Long-horizon task
+ broad environment variation
+ partially completed task
+ recovery planning

Level 5
Full deployment task
+ broad environment variation
+ unexpected states
+ open-ended subtask generation
```

The robot is progressively taught:

> **Solve the task.**

then:

> **Solve the task under different environmental conditions.**

then:

> **Solve the task even when the world is not in the state you expected.**

and ultimately:

> **Figure out what needs to be done next when the original plan is no longer sufficient.**

---

# 2.16 VLM Agent as Scenario Generator

The VLM training agent is responsible for constructing the complete training scenario rather than merely providing a reward.

It observes:

```text
Environment description
Robot capabilities
Task requirements
Current curriculum
Policy performance
Simulation outcomes
Known failure modes
```

and decides:

```text
What task next?
What subtask structure is appropriate?
What environment configuration?
What initial state?
How difficult should it be?
What variation should be introduced?
Should a new subtask be inserted during execution?
Should the curriculum advance or remain at the current level?
```

For example:

```text
Current robot system is excellent at:
   picking objects from clear tables

Current weakness:
   recovering after object displacement

VLM Training Agent decides:
   increase task complexity slightly
   introduce object displacement
   initialize from partially completed states
   create recovery-subtask opportunities during execution
   maintain realistic environment variation
```

The next batch of simulation episodes is then specifically designed around that weakness. The VLM training agent that constructs scenarios is conceptually separate from the **high-level runtime planner** trained in Stage 7: the training agent chooses what experiences the system should learn from, while the runtime planner chooses what the robot should accomplish next during an actual task.

---

# 2.17 Agentic RL

Once the VLM training agent constructs a scenario, the robot operates autonomously inside the simulator using the same complete architecture that will be used during deployment.

```text
Scenario Generated by VLM Training Agent
       ↓
Observation
       ↓
High-Level VLM Planning
       ↓
Dynamic Subtask Selection
       ↓
Target-Future Generation
       ↓
VLM Target Selection
       ↓
Candidate Action Generation
       ↓
World Model Prediction
       ↓
VLM Evaluation
       ↓
Action Selection
       ↓
Execution
       ↓
New Observation
       ↓
Continue / Recover / Replan / Create New Subtask
```

The agentic RL objective can be expressed conceptually as:

```math
\max_{\theta}\mathbb{E}_{T,E,s_0}\left[R_{\mathrm{task}}\left(\tau;T,E,s_0\right)\right]
```

where `theta` represents the environment-specific action-generation policy parameters. The objective is evaluated over a distribution of tasks, environment configurations, and initial states rather than a single fixed scenario. The high-level planner remains active during every rollout and its decisions affect the task outcome, but the dedicated reasoning SFT + RL that teaches hierarchical decomposition, subgoal selection, and dynamic subtask creation is Stage 7 training of the planner. Stage 8 specializes the robot's execution policy to its deployment environment; it should not be confused with Stage 7 planner reasoning training.

The learning process therefore trains the complete decision system:

```text
Planning
   ↓
Subtask Creation
   ↓
Reasoning
   ↓
Target Selection
   ↓
Action Proposal
   ↓
World Model Prediction
   ↓
VLM Evaluation
   ↓
Action Selection
   ↓
Execution
   ↓
Outcome
   ↓
Replanning / Recovery
```

This is important because an eventual deployment failure may originate in the high-level plan, dynamic subtask decomposition, VLA candidate generation, World Model prediction, VLM ranking, target replanning, or recovery behavior rather than in the low-level policy alone.

---

# 2.18 Simulation Curriculum Feedback

The VLM training agent receives feedback from the outcome of simulation episodes and uses that feedback to continuously update the curriculum.

```text
Generated Scenario
       ↓
Agentic RL
       ↓
Outcome
       ↓
Success / Failure / Partial Success
       ↓
VLM Curriculum Agent
       ↓
Next Scenario
```

When the policy becomes highly reliable on a task distribution, the VLM training agent increases task complexity or variation. When performance drops, it generates more examples at the current level until the system recovers.

```text
Current Level
     ↓
Performance
     ↓
 ┌────┴────┐
 │         │
Good      Poor
 │         │
 ↓         ↓
Harder    More practice
Tasks     at current level
 │         │
 └────┬────┘
      ↓
Next Episodes
```

At the same time, the agent tracks where dynamic replanning is required. If failures cluster around incomplete plans or missing intermediate objectives, it can generate episodes specifically designed to require **online subgoal and subtask creation**.

---

# 2.19 Learning Robustness, Not Memorization

The purpose of environment-specific RL is not to memorize the exact geometry observed during calibration. It is to learn the invariances and recovery strategies that remain valid across the expected deployment distribution.

For example, a robot should not learn:

> "The cup is always 22 cm to the left of the plate."

It should learn:

> "The cup is usually located somewhere on this workspace, and I should visually locate it and approach it safely."

The simulator therefore deliberately perturbs variables that should **not** change the underlying task strategy. This makes Stage 8 a form of **environment-specific robustness learning**: the robot is trained against the variations expected around its deployment environment rather than being overfit to one deterministic simulator state.

---

# 2.20 Simulation-to-Real Validation Gate

Before deployment, the environment-specific robot system must pass a final validation process across both the nominal digital twin and the generated robustness distribution.

```text
Calibrated Twin
     ↓
Nominal Evaluation
     ↓
Variation Evaluation
     ↓
Recovery Evaluation
     ↓
Replanning Evaluation
     ↓
Long-Horizon Evaluation
     ↓
Safety Evaluation
     ↓
High-Fidelity Validation
     ↓
Deployment
```

The robot should only move toward real deployment once performance is robust across the relevant distribution rather than merely strong on the original reconstructed scene.

---

# 2.21 Why Stage 8 Is Separate From Real Deployment

The purpose of Stage 8 is to move as much of the **learning burden as possible into simulation before the robot encounters the real environment**.

The system should therefore enter deployment having already experienced:

```text
Nominal behavior
+
Task progression
+
Object variation
+
Scene variation
+
Physics variation
+
Sensor noise
+
Task variation
+
Dynamic replanning
+
Unexpected disturbances
+
Partial task states
+
Dynamic subtask creation
+
Recovery
+
Long-horizon execution
```

Real deployment then becomes primarily a process of **validation, calibration, and discovering residual gaps**, rather than the first time the robot experiences meaningful variation.

---

# 2.22 Why Co-Training Instead of Sequential Fine-Tuning?

The different global data layers contain complementary information, and sequential fine-tuning risks allowing the final, smallest dataset to dominate the model. Massive human video provides diversity and broad physical knowledge, synthetic human demonstrations provide large-scale action-conditioned coverage and a mechanism for probing current model weaknesses, the 3D-conditioned stage teaches the World Model to exploit explicit geometric scene state under realistic observation conditions, teleoperation provides highly accurate grounding in robot embodiment, action-conditioned training teaches the World Model the relationship between actions and physical consequences, and **planner reasoning training teaches the high-level VLM how to use world knowledge, execution history, and available skills to select appropriate subgoals and dynamically revise hierarchical plans**.

A simplified objective is:

```math
\mathcal{L}
=
\lambda_{\mathrm{WM}}\mathcal{L}_{\mathrm{WM}}
+
\lambda_{\mathrm{WM-3D}}\mathcal{L}_{\mathrm{WM-3D}}
+
\lambda_{\mathrm{WM-action}}\mathcal{L}_{\mathrm{WM-action}}
+
\lambda_{\mathrm{human}}\mathcal{L}_{\mathrm{human}}
+
\lambda_{\mathrm{synthetic}}\mathcal{L}_{\mathrm{synthetic}}
+
\lambda_{\mathrm{gripper}}\mathcal{L}_{\mathrm{gripper}}
+
\lambda_{\mathrm{teleop}}\mathcal{L}_{\mathrm{teleop}}
+
\lambda_{\mathrm{planner}}\mathcal{L}_{\mathrm{planner}}
```

The planner loss is applied to the high-level VLM planner. It is distinct from the VLA action-generation losses and from any separate VLM evaluator or safety-checker losses.

The fundamental principle is:

> **Low-fidelity real data provides scale and broad world knowledge; synthetic data expands action-conditioned coverage and reveals model weaknesses; 3D-conditioned training teaches the World Model to use explicit geometry under realistic observation conditions; high-fidelity data provides grounding in manipulation and robot control; action-conditioned training turns prediction into simulation; planner reasoning training teaches the high-level VLM to select and create useful subgoals; environment-specific agentic RL turns the global system into robust deployment behavior.**

---

# 2.23 Active Real-World Data Collection

The synthetic egocentric-data stage creates an explicit mechanism for deciding what additional real-world data should be acquired. Instead of assuming that more random video is always beneficial, Lunch Robotics can compare the foundation model's performance across generated tasks, measure generation quality, and identify the categories in which the model consistently produces implausible or incomplete behavior.

A simplified priority function can be thought of as:

```math
p(c)\propto\mathrm{Frequency}(c)\cdot\mathrm{FailureRate}(c)\cdot\mathrm{TransferValue}(c)
```

where `c` denotes a capability or interaction category. This creates a continually updated **data acquisition frontier** describing where additional reality-grounded video is most valuable.

```text
Foundation Model
     ↓
Synthetic Generation
     ↓
Failure Distribution
     ↓
Capability Gap Map
     ↓
Prioritized Real Data Collection
     ↓
New Real Egocentric Data
     ↓
Foundation Model Update
```

---

# 3. Deployment Inputs

Once the global robot system exists, a new deployment requires only a relatively small amount of environment-specific information. The user records a **30–90 second walkthrough video** of the workspace, providing coarse geometry, furniture, large objects, spatial layout, camera scale, and an initial scene representation from which the Real-to-Sim Agent can construct the initial digital twin.

The user also provides an **Environment & Task Context Description** explaining what the environment is used for, which objects matter, what tasks the robot is expected to perform, what constitutes success, unusual or non-obvious procedures, environmental constraints, and the behavior of other agents such as humans or autonomous systems.

The user additionally performs **tactile environment probing** with a sensorized handheld gripper. The gripper can tap, slide, push, lift, squeeze, deform compliant materials, and probe contact conditions while recording RGB-D video, pose, forces, tactile signals, slip information, and interaction trajectories. This information is primarily used for physical system identification rather than simply visual reconstruction.

Each task has two videos: an **execution demonstration**, where the human performs the task naturally without explaining it and which is primarily used for evaluation, and a **tutorial video**, where the human explains the task and its important details while performing it. The tutorial is transformed into a modular skill representation that can later be dynamically retrieved by the robot.

Finally, the robot vendor provides the **robot SDK and hardware specification**, including kinematics, joint limits, actuator information, end-effector specifications, and hardware control constraints. The vendor also provides approximately one minute of random-policy execution containing synchronized observations and actions. This recording helps the system identify robot dynamics, actuator behavior, latency, joint response, and low-level control characteristics.

---

# 4. The Real-to-Sim Agent

The system does not rely on a fixed, hand-engineered simulator-generation pipeline. Instead, an **RL-trained LLM agent acts as an autonomous simulation engineer and system-identification engineer**.

Given the environment videos, tactile measurements, task descriptions, robot specifications, and World Model predictions, the agent uses specialized tools to construct the digital twin. The main tools are:

1. A 3D reconstruction tool (for example, a streaming reconstruction or mapping system).
2. A learnable simulator framework, such as a differentiable physics engine.

The agent updates robot physics using the robot specification, initial random-robot-action data, task videos, and streaming egocentric video during production. It updates environment physics using the simulator framework's differentiable-learning support, task videos, and streaming egocentric video. It updates initial conditions using the initial walkthrough video and the output of the streaming 3D reconstruction tool. It updates materials, weight, and weight distribution using initial tactile data from humans interacting with objects, task videos, and later production video of humans manipulating objects.

```mermaid
flowchart TD
    A["Environment Walkthrough"] --> B["Real-to-Sim Agent"]
    C["Tactile + Contact Data"] --> B
    D["Environment & Task Context"] --> B
    E["Task Videos"] --> B
    F["Robot Specs + Random Policy"] --> B
    G["World Model"] --> B
    B --> H["Geometry"]
    B --> I["Objects"]
    B --> J["Agents"]
    B --> K["Physics"]
    B --> L["Materials"]
    B --> M["Dynamics"]
    H --> N["Initial Digital Twin"]
    I --> N
    J --> N
    K --> N
    L --> N
    M --> N
    N --> O["System Identification"]
    O --> P["Calibrated Digital Twin"]
    P --> Q["Validation"]
    Q -->|"Mismatch"| B
    Q -->|"Validated"| R["Finalize"]
```

The goal is not to create a visually perfect replica of the environment. It is to create a **useful executable model of the world** that is sufficiently accurate for policy training, candidate-action evaluation, and counterfactual reasoning.

---

# 5. Agentic System Identification

Visual reconstruction alone cannot reveal many physical quantities that matter for manipulation. The system may need to estimate object mass, friction, compliance, damping, restitution, actuator response, and other parameters.

```math
\theta=
\begin{bmatrix}
m \\
\mu_s \\
\mu_d \\
e \\
k \\
c \\
\vdots
\end{bmatrix}
```

The system then minimizes the mismatch between real and simulated behavior:

```math
\theta^*
=
\arg\min_{\theta}
\left(
D_{\mathrm{kin}}\left(\tau_{\mathrm{real}},\tau_{\mathrm{sim}}(\theta)\right)
+
\lambda D_{\mathrm{force}}\left(W_{\mathrm{real}},W_{\mathrm{sim}}(\theta)\right)
\right)
```

The important distinction is that the **agent reasons about what must be identified and which experiment will be informative**, while specialized numerical tools perform the parameter estimation itself.

---

# 6. Explicit Physics + General World Model

Not every entity in the environment should be represented with hand-designed physics. Rigid objects and contact interactions can often be modeled explicitly, while humans, pets, and other complex non-scripted entities are better represented through the general World Model.

```math
S^{\mathrm{agent}}_{t+1:t+H}
\sim
P_{\phi}(S_{\leq t},E)
```

The digital twin therefore becomes a hybrid system:

```math
\text{Digital Twin}
=
\text{Explicit Calibrated Physics}
+
\text{Learned World Dynamics}
```

The action-conditioned World Model adds another layer: explicit physics can handle quantities that require physical precision, while the learned model can predict complex, difficult-to-model interactions and other-agent behavior.

---

# 7. Learned Surrogate Simulator

High-fidelity simulation is necessary for calibration and validation but is too expensive to run at the scale required for large-scale RL. The calibrated digital twin is therefore used to generate experience from which the system trains a learned **surrogate simulator**.

```math
f_{\mathrm{sim}}(s_t,a_t)
\approx
f_{\mathrm{surrogate}}(s_t,a_t)
```

The surrogate trades some physical fidelity for enormous speed and becomes the main engine for large-scale policy optimization and environment-specific training.

```mermaid
flowchart LR
    A["Calibrated Digital Twin"] --> B["Synthetic Experience"]
    B --> C["Train Surrogate"]
    C --> D["Fast Simulator"]
    D --> E["Massive Policy Rollouts"]
    E --> F["Candidate Policies"]
    F --> G["High-Fidelity Twin Validation"]
```

The resulting hierarchy is:

```math
\text{Real World}
\rightarrow
\text{Calibrated Twin}
\rightarrow
\text{Surrogate}
\rightarrow
\text{Massive Agentic RL}
```

There are therefore two simulation mechanisms in the architecture. The high-fidelity digital twin provides accurate physical validation, while the learned World Model provides extremely fast latent prediction for candidate actions.

---

# 8. Training the Environment-Specific VLA

The environment-specific VLA begins from the current **global or mixed VLA**, rather than being trained from scratch. It is conditioned on the current observation, the current subtask selected by the high-level planner, the selected World Model target imagination, and retrieved tutorial context:

```math
\pi_{\theta}
\left(a_t\mid o_t,T_{\mathrm{subtask}},z_{\mathrm{target}},S_{\mathrm{tutorial}}\right)
```

The VLA generates three candidate action chunks:

```math
\{a_t^{(1)},a_t^{(2)},a_t^{(3)}\}
\sim
\pi_{\theta}
\left(a_t\mid o_t,T_{\mathrm{subtask}},z_{\mathrm{target}},S_{\mathrm{tutorial}}\right)
```

The World Model then predicts one future for each candidate, and the VLM selects the candidate whose predicted future best matches the selected target.

The **environment-specific training** of this VLA occurs through Stage 8 agentic RL:

```text
Global / Mixed Robot System
       ↓
Calibrated Environment
       ↓
VLM Curriculum Agent
       ↓
Task + Subtask + Environment + Initial-State Generation
       ↓
Agentic RL
       ↓
Environment-Specific VLA / Robot System
```

The resulting environment-specific VLA therefore learns not just the nominal environment, but the expected distribution of physical conditions and recovery states around that environment. The high-level planner remains responsible for semantic decomposition and selecting the next subgoal; the VLA remains responsible for proposing physical action chunks.

---

# 9. Imagined Goal States and Target Future Trajectories

Exact trajectory matching is brittle because many different trajectories can successfully solve the same task. Instead, the runtime system uses the World Model to construct multiple possible latent representations of successful future states and short-horizon motions:

```math
\left\{
z_{\mathrm{target}}^{(1)}(t:t+H),
z_{\mathrm{target}}^{(2)}(t:t+H),
z_{\mathrm{target}}^{(3)}(t:t+H)
\right\}
```

The three target trajectories represent different plausible ways the world could be heading if the current subtask is progressing successfully. They do not need to prescribe an exact robot trajectory, because two very different physical motions may lead to essentially the same successful state.

The VLM evaluates these possibilities:

```math
r_j^{\mathrm{target}}
=
\mathrm{VLM}_{\mathrm{score}}
\left(z_{\mathrm{target}}^{(j)},s_t,T_{\mathrm{subtask}}\right)
```

and selects:

```math
j^*=\arg\max_{j\in\{1,2,3\}}r_j^{\mathrm{target}}
```

The selected target is:

```math
z_{\mathrm{target}}=z_{\mathrm{target}}^{(j^*)}
```

This target remains active across multiple action chunks.

---

# 10. Three-Action World Model Selection

At each lower-frequency target update, the World Model first samples three possible target futures:

```math
\left\{z_{\mathrm{target}}^{(1)},z_{\mathrm{target}}^{(2)},z_{\mathrm{target}}^{(3)}\right\}
\sim
\mathrm{WM}(s_t,T_{\mathrm{subtask}})
```

The VLM scores them and selects:

```math
j^*=\arg\max_{j\in\{1,2,3\}}r_j^{\mathrm{target}}
```

The selected target future is:

```math
z_{\mathrm{target}}=z_{\mathrm{target}}^{(j^*)}
```

This target remains fixed while the high-frequency action-selection loop executes multiple action chunks.

At each high-frequency action decision, the VLA generates three candidate action chunks:

```math
\{a_t^{(1)},a_t^{(2)},a_t^{(3)}\}
```

Each candidate is passed through the action-conditioned World Model once:

```math
z_{\mathrm{pred}}^{(i)}(t:t+H)
=
\mathrm{WM}_{\mathrm{action}}(s_t,a_t^{(i)})
```

The VLM compares them against the selected target:

```math
r_i=
\mathrm{VLM}_{\mathrm{score}}
\left(z_{\mathrm{pred}}^{(i)},z_{\mathrm{target}},s_t,T_{\mathrm{subtask}}\right)
```

and selects:

```math
i^*=\arg\max_{i\in\{1,2,3\}}r_i
```

The architecture therefore has exactly:

```text
3 target futures
       ↓
1 selected target
       ↓
3 candidate action chunks
       ↓
3 predicted candidate futures
       ↓
1 selected action
```

The **target-imagination process happens at lower frequency**, while the **action-chunk generation, World Model simulation, and VLM candidate-scoring process happens at higher frequency**. The high-level VLM planner sits above both and can invalidate the current target by changing the active subtask.

---

# 11. Environment-Specific Agentic RL

Stage 8 is the core **pre-deployment learning stage**. The VLM curriculum agent continuously constructs simulation episodes, while the environment-specific robot system learns through RL.

The stage follows the same learning philosophy as online local RL:

```text
Task Curriculum
+
Dynamic Subtask Structure
+
Environment Variation
+
Initial-State Variation
       ↓
Agentic RL
       ↓
Improved Robot System
       ↓
New Evaluation
       ↓
Updated Curriculum
```

The difference is that Stage 8 can perform this learning at massive scale before exposing the policy to physical hardware.

The VLM training agent creates the scenario, while the robot system learns inside it:

```text
OUTER AGENT
VLM Training Agent
   ↓
chooses task + subtask opportunities + environment + initial state

INNER ROBOT SYSTEM
High-Level VLM Planner + VLA + World Model
   ↓
plans + reasons + imagines + acts + evaluates + recovers
```

This creates a nested optimization loop:

```text
                VLM TRAINING AGENT
                        │
           ┌────────────┼────────────┐
           │            │            │
        Task /      Environment   Initial State
        Subtasks     Variations    Variations
           │            │            │
           └────────────┼────────────┘
                        ↓
                 Simulation Episode
                        ↓
                    Agentic RL
                        ↓
                 System Evaluation
                        ↓
                VLM Curriculum Update
                        ↺
```

The simulator therefore becomes an **automatically generated training distribution** rather than a single fixed environment.

---

---

# 12. Real Deployment and Local Continual Learning

Simulation will inevitably differ from reality, so deployment creates a continual local learning loop. The environment-specific robot system operates on the real robot while its calibrated simulator continues to generate additional training scenarios and perform local simulation RL.

```math
\text{Real Deployment}
\rightarrow
\text{Observed Outcome}
\rightarrow
\text{Failure Identification}
\rightarrow
\text{VLM Scenario Generation}
\rightarrow
\text{Targeted Simulation}
\rightarrow
\text{Online Local RL}
\rightarrow
\text{Updated Environment-Specific Robot System}
```

The local post-deployment loop uses the same mechanism as Stage 8:

```text
VLM Training Agent
       ↓
Task Curriculum
+
Dynamic Subtask Generation
+
Environment Variations
+
Initial-State Variations
       ↓
Local Simulation RL
```

The key difference is that **real-world failures now provide the evidence used to decide what new scenarios should be generated**.

---

# 13. Post-Deployment Failure Data

Every deployed robot is a source of valuable real-world experience. The system monitors execution continuously and identifies situations in which the robot makes a mistake. These mistakes can be explicitly reported by the human user or automatically detected by a VLM observing the robot through cameras installed in the environment.

The system captures the complete decision context:

```text
Failure Episode
├── Task being executed
├── Hierarchical plan at time of failure
├── Current subtask / selected subgoal
├── Dynamically generated subtasks
├── Environment and task context
├── Robot state trajectory
├── High-level planner reasoning trace, plan revisions, and subgoal decisions
├── Target trajectory 1
├── Target trajectory 2
├── Target trajectory 3
├── VLM score for target 1
├── VLM score for target 2
├── VLM score for target 3
├── Selected target trajectory
├── Candidate action 1
├── Candidate action 2
├── Candidate action 3
├── World Model prediction for action 1
├── World Model prediction for action 2
├── World Model prediction for action 3
├── VLM score for predicted future 1
├── VLM score for predicted future 2
├── VLM score for predicted future 3
├── Selected action
├── Actual robot action trajectory
├── Real video of the robot execution
├── Failure timestamp / event
├── Human feedback or VLM failure judgment
└── Relevant simulator state + environment parameters
```

---

# 14. Failure-Conditioned Data Generation

A single failure should not remain a single training example. Once a failure is identified, the system reconstructs the relevant state inside the calibrated digital twin and uses the VLM training agent to generate a local curriculum around the failure.

```mermaid
flowchart TD
    A["Real-World Failure"] --> B["Failure Episode"]
    B --> C["Task + Plan + Subtask"]
    B --> D["Planner Reasoning Context"]
    B --> E["3 Target Futures"]
    B --> F["VLM Target Scores"]
    B --> G["3 Candidate Actions"]
    B --> H["3 World Model Predictions"]
    B --> I["VLM Candidate Scores"]
    B --> J["Real Execution"]
    C --> K["Failure Analysis"]
    D --> K
    E --> K
    F --> K
    G --> K
    H --> K
    I --> K
    J --> K
    K --> L["Failure Condition"]
    L --> M["VLM Curriculum Agent"]
    M --> N["Task Variation"]
    M --> O["Subtask / Plan Variation"]
    M --> P["Environment Variation"]
    M --> Q["Initial-State Variation"]
    N --> R["Local Agentic RL"]
    O --> R
    P --> R
    Q --> R
```

The objective is not to memorize the original mistake. It is to train the system on the surrounding state distribution so that similar failures become less likely.

For example:

```text
Real failure:
object displaced unexpectedly

       ↓

VLM training agent generates:
├── small displacement
├── medium displacement
├── large displacement
├── different object orientations
├── partially completed manipulation
├── alternative recovery states
├── recovery subtask creation opportunities
└── nearby environmental variations

       ↓

Local curriculum RL
       ↓
Improved recovery behavior
```

---

# 15. Sending Deployment Failures Back to the Lunch Robotics Lab

The raw failure stream is also sent back to the Lunch Robotics lab, but **raw deployment failures are not directly used to fine-tune the global system**.

The central dataset pipeline begins with human-led error-mode analysis because different deployments produce different kinds of mistakes and not all failures represent deficiencies in the universal brain.

Some failures are caused by local geometry, local objects, local task conventions, local planner behavior, or local environment calibration. Others reveal genuine limitations in general robot intelligence, World Model prediction, target imagination, candidate generation, planner reasoning, VLM evaluation, dynamic planning, or action selection.

The Lunch Robotics team analyzes failure episodes across deployments to identify **systematic, recurring, and generalizable error modes**.

---

# 16. Curating the Global Training Dataset

Once a generalizable error mode has been identified, the Lunch Robotics team creates a **curated training dataset** around that capability gap.

```math
\mathcal{D}_{\mathrm{curated}}
=
\mathrm{Curate}\left(
\mathcal{D}_{\mathrm{deployment}},
\mathcal{D}_{\mathrm{simulation}},
\mathcal{D}_{\mathrm{planner}},
\mathcal{D}_{\mathrm{WM}},
\mathcal{D}_{\mathrm{VLM}},
\mathcal{D}_{\mathrm{foundation}}
\right)
```

The curated dataset can combine selected real-world failures with successful examples, counterfactual successful trajectories, targeted simulation rollouts, adversarial near-failure scenarios, World Model imagined trajectories, decoded imagined videos, planner reasoning traces, candidate action sets, target-future alternatives, VLM rankings, dynamic subtask plans, replan decisions, and relevant examples from the original foundation datasets.

The central question is:

> **What capability is actually missing, where in the decision loop does the failure originate, and what data will teach the global system that capability in a way that transfers beyond the environment where the failure was observed?**

---

# 17. Fine-Tuning a New Global Lab System

The curated dataset is used to improve the **global lab system**:

```math
W_{\mathrm{lab}}'
=
\mathrm{FineTune}\left(W_{\mathrm{lab}},\mathcal{D}_{\mathrm{curated}}\right)
```

Different capability gaps may require updating different components.

**VLA:**

```math
W_{\mathrm{VLA}}'
=
\mathrm{FineTune}\left(W_{\mathrm{VLA}},\mathcal{D}_{\mathrm{VLA,curated}}\right)
```

**World Model:**

```math
W_{\mathrm{WM}}'
=
\mathrm{FineTune}\left(W_{\mathrm{WM}},\mathcal{D}_{\mathrm{WM,curated}}\right)
```

**High-Level Planner VLM:**

```math
W_{\mathrm{planner}}'
=
\mathrm{FineTune}\left(W_{\mathrm{planner}},\mathcal{D}_{\mathrm{planner,curated}}\right)
```

Planner-specific data contains successful and failed hierarchical plans, subgoal choices, dynamically inserted subtasks, prerequisite discoveries, replan decisions, and the resulting task outcomes. This is the component directly aligned with Stage 7's planner reasoning objective.

**VLM Evaluator / Safety Checker:**

```math
W_{\mathrm{evaluator}}'
=
\mathrm{FineTune}\left(W_{\mathrm{evaluator}},\mathcal{D}_{\mathrm{evaluator,curated}}\right)
```

The VLM evaluator or safety checker may also be updated when the identified capability gap concerns future ranking, task alignment, or safety. These are separate responsibilities from high-level planning.

The resulting global system therefore evolves as a coordinated stack:

```text
Curated Failure Data
       ↓
┌────────┼───────────────┬────────────────┐
↓        ↓               ↓                ↓
VLA   World Model   Planner VLM    VLM Evaluator
↓        ↓               ↓                ↓
└────────┼───────────────┴────────────────┘
         ↓
     New Lab System
```

---

# 18. Environment-Specific Models and the Global Lab System

The architecture maintains two distinct classes of model system.

The **global lab system** is the shared general-purpose model stack maintained by Lunch Robotics, capturing broadly reusable capabilities learned from the foundation data funnel and from curated deployment-derived training. It includes the World Model, VLA, high-level planner VLM, and any VLM evaluator or safety components.

Each deployment maintains its own **environment-specific system**, which is optimized for its physical environment, robot embodiment, objects, task distribution, local dynamics, and local operating conventions.

```text
                         GLOBAL LAB SYSTEM
              (World Model + VLA + Planner + Evaluators)
                              │
                              ▼
                      Shared Initialization
                              │
                ┌─────────────┼─────────────┐
                │             │             │
                ▼             ▼             ▼
          Environment A  Environment B  Environment C
          Specific Stack Specific Stack Specific Stack
                │             │             │
                ▼             ▼             ▼
       Local Agentic Simulation RL
                │
                ▼
       Real Deployment + Local Learning
```

---

# 19. Continual Weight Mixing

The next step is to combine the information accumulated by the global lab system and the environment-specific systems.

Suppose:

```math
W_{\mathrm{lab}}
```

is the current global model stack and:

```math
W_1,W_2,\ldots,W_N
```

are environment-specific model systems.

After the lab system has been improved using curated deployment data:

```math
W_{\mathrm{lab}}'
=
\mathrm{FineTune}\left(W_{\mathrm{lab}},\mathcal{D}_{\mathrm{curated}}\right)
```

the system performs a continual model-merging step:

```math
W_{\mathrm{mixed}}
=
\mathrm{Mix}\left(W_{\mathrm{lab}}',W_1,W_2,\ldots,W_N\right)
```

The exact implementation may eventually use parameter-space merging, model deltas, adapter composition, selective merging, or another technique. The architectural principle is that useful information accumulated centrally and useful knowledge discovered through deployment should be consolidated into a stronger common initialization.

The resulting **mixed system** is redistributed to all environments. Each deployment then starts a new round of environment-specific adaptation from this shared state:

```math
W_i^{(t+1)}
=
\mathrm{Adapt}\left(W_{\mathrm{mixed}},E_i\right)
```

Every global update therefore becomes available to the entire fleet.

---

# 20. The Global-Local Learning Flywheel

Lunch Robotics consequently has four interacting learning loops.

The **foundation-model loop** discovers and scales general physical knowledge:

```math
\boxed{
\text{Real Human Egocentric Video}
\rightarrow \text{Foundation World Model}
\rightarrow \text{Synthetic Egocentric Data}
\rightarrow \text{Failure / Weakness Analysis}
\rightarrow \text{Targeted Real Data Collection}
\rightarrow \text{Stronger Foundation World Model}
}
```

The **planner-learning loop** improves high-level task reasoning:

```math
\boxed{
\text{Planner Demonstrations + Execution Rollouts}
\rightarrow \text{Planner Reasoning SFT}
\rightarrow \text{Planner RL}
\rightarrow \text{Better Subgoal Selection + Dynamic Planning}
\rightarrow \text{Planner Failure Analysis}
\rightarrow \text{Curated Planner Data}
\rightarrow \text{Improved Global Planner}
}
```

The **pre-deployment environment loop** specializes the global system:

```math
\boxed{
\text{Global Robot System}
\rightarrow \text{Real-to-Sim}
\rightarrow \text{Calibrated Environment}
\rightarrow \text{VLM Curriculum Agent}
\rightarrow \text{Task + Dynamic Subtasks + Environment + Initial-State Variation}
\rightarrow \text{Agentic RL}
\rightarrow \text{Robust Environment-Specific System}
}
```

The **post-deployment local loop** continues improving the specific environment:

```math
\boxed{
\text{Environment-Specific System}
\rightarrow \text{Real Deployment}
\rightarrow \text{Failure}
\rightarrow \text{VLM Scenario Generation}
\rightarrow \text{Targeted Simulation}
\rightarrow \text{Local Agentic RL}
\rightarrow \text{Improved Environment-Specific System}
}
```

The **global fleet loop** turns collective experience into improvements to the shared brain:

```math
\boxed{
\text{Many Deployments}
\rightarrow \text{Failure Episodes}
\rightarrow \text{Team Error-Mode Analysis}
\rightarrow \text{Curated Data}
\rightarrow \text{New Lab System}
\rightarrow \text{Weight Mixing}
\rightarrow \text{All Deployments}
}
```

Together, these loops create a system in which **reality grounds the model, the model generates hypotheses, synthetic data expands training, simulation builds robustness, VLM training agents construct curricula, the high-level planner dynamically decomposes and replans tasks, deployment discovers residual failures, and curated experience improves the shared brain**.

---

# 21. Safety

Safety is handled through two main mechanisms.

First, the **high-level planner VLM is supervised fine-tuned specifically for safety**, so that it learns to identify unsafe tasks, situations, and intended subgoals before the robot begins acting. Safety constraints are maintained during replanning, not just at initial task decomposition.

Second, at inference time, an additional **VLM safety checker** evaluates the candidate trajectories produced by the system.

The target futures can be safety-checked before one is selected as the desired future, while the three predicted candidate futures can also be safety-checked before the final candidate is selected:

```text
HIGHEST LEVEL
High-Level VLM Planner
 ↓
Task / Subtask Plan
 ↓
Replanning + Safety Constraints

LOW FREQUENCY
World Model
 ↓
3 Target Futures
 ↓
VLM Target Evaluation + Safety
 ↓
Safe Target Future
 ↓

HIGH FREQUENCY
VLA
 ↓
3 Candidate Actions
 ↓
World Model
 ↓
3 Predicted Future Trajectories
 ↓
Decode to Visual Trajectories
 ↓
VLM Safety Checker
 ↓
Reject Unsafe Futures
 ↓
VLM Task-Alignment Scoring
 ↓
Select Safe Best Candidate
 ↓
Robot
```

The safety mechanism therefore sits directly inside the planning, target-selection, and action-selection loops rather than relying only on the policy to have learned safe behavior.

---

# 22. Robot SDK

The Robot SDK is the final hardware abstraction layer between the policy and the physical robot.

The robot should interact with the user through natural language while the underlying system translates those instructions into task specifications, builds and revises a hierarchical task plan, retrieves the appropriate skill context, determines the required **planner reasoning budget**, generates target futures for the current subtask, evaluates them, generates VLA candidate actions, evaluates their possible consequences through the World Model and VLM, and converts the selected action into a hardware-safe trajectory.

```mermaid
flowchart LR
    A["Human Speech"] --> B["Speech-to-Text"]
    B --> C["Task Specification"]
    C --> D["Tutorial Skill Retrieval"]
    D --> E["High-Level VLM Planner"]
    E --> F["Current Subtask"]
    F --> G["World Model: 3 Target Futures"]
    G --> H["VLM Target Selection"]
    H --> I["Selected Target Future"]
    I --> J["Environment-Specific VLA: 3 Candidate Actions"]
    J --> K["World Model: 3 Predicted Futures"]
    K --> L["VLM Future Comparison"]
    L --> M["Selected Action"]
    M --> N["Robot SDK"]
    N --> O["Motion Smoother"]
    O --> P["Hardware Controller"]
    P --> Q["Robot"]
    Q --> R["New Observation"]
    R --> E
```

---

# 22.1 Motion Smoothing

A neural action policy can produce noisy or abrupt outputs, so the SDK applies online smoothing and jerk-limited trajectory generation:

```math
(q_{t+1},\dot{q}_{t+1},\ddot{q}_{t+1})
=
f_{\mathrm{smooth}}(a_t,a_{<t},q_t,\dot{q}_t,\ddot{q}_t)
```

The resulting trajectory is constrained by the robot's hardware limits:

```math
\begin{aligned}
|\dot{q}(t)| &\leq v_{\max} \\
|\ddot{q}(t)| &\leq a_{\max} \\
|\dddot{q}(t)| &\leq j_{\max}
\end{aligned}
```

These constraints are read directly from the robot specification, keeping the policy hardware-agnostic. For low-latency deployment, a diffusion-based action policy can additionally undergo **Consistency Distillation**, allowing an iterative diffusion policy to be transformed into a much faster inference process.

---

# 23. Runtime Task Execution

At runtime, the experience should be intentionally simple. The user can say:

> **"Clean the table."**

The system converts the request into a task specification, retrieves the relevant tutorial, and hands the problem to the high-level VLM planner. The planner constructs an initial hierarchical plan, selects the first subtask, and continuously revises the plan as execution proceeds.

The runtime system then operates on three nested timescales.

```mermaid
flowchart TD
    A["User: Clean the table"] --> B["Speech Recognition"]
    B --> C["Task Specification"]
    C --> D["Retrieve Table-Cleaning Tutorial"]
    D --> E["High-Level VLM Planner"]
    E --> F["Hierarchical Plan + Current Subtask"]
    F --> G["LOW FREQUENCY: World Model"]
    G --> H["3 Target Future Trajectories"]
    H --> I["VLM Target Selection"]
    I --> J["Selected Target Future"]
    J --> K["HIGH-FREQUENCY LOOP"]
    K --> L["VLA Samples 3 Candidate Actions"]
    L --> M["World Model Predicts 3 Futures"]
    M --> N["VLM Compares 3 Futures to Target"]
    N --> O["Select Best Candidate"]
    O --> P["Safety Validation"]
    P --> Q["Robot SDK"]
    Q --> R["Robot"]
    R --> S["New Observation"]
    S --> T["High-Level Planner Re-evaluates"]
    T --> U{"Plan Still Valid?"}
    U -->|"Yes"| K
    U -->|"No"| V["Replan / Create New Subtask"]
    V --> F
```

The high-frequency loop repeatedly executes:

```text
VLA
 ↓
3 Candidate Actions
 ↓
3 World Model Futures
 ↓
VLM Comparison
 ↓
Best Action
 ↓
Execute Action Chunk
 ↓
New Observation
 ↓
Repeat
```

The lower-frequency loop periodically executes:

```text
Current Observation + Current Subtask
 ↓
World Model
 ↓
3 Target Futures
 ↓
VLM Scoring
 ↓
Selected Target
 ↓
Return to High-Frequency Loop
```

The highest-level loop periodically or eventfully executes:

```text
Current Observation + Execution History
 ↓
High-Level VLM Planner
 ↓
Check Current Plan
 ↓
Keep / Modify / Reorder / Remove / Create Subtask
 ↓
Updated Plan
 ↓
Select Current Subtask
```

The target therefore remains active across multiple action chunks. It is recomputed when the robot has made sufficient progress, when the current target has become stale, when the observed state diverges from the expected state, when the environment changes materially, when uncertainty increases sufficiently to justify replanning, or when the high-level planner switches to a different subtask.

---

# 24. End-to-End System

The full architecture is a continuous pipeline from foundation-model training to deployment and back again.

```mermaid
flowchart TD
    A["FOUNDATION MODEL DATA FUNNEL"]
    A1["MASSIVE Human Egocentric Video"]
    A2["LARGE Human Video + VLM Pose Waypoints"]
    A3["VERY LARGE Synthetic Human Egocentric Data"]
    A4["MEDIUM Human + Data-Collecting Gripper"]
    A5["SMALL Teleoperation"]
    A1 --> B["PASSIVE WORLD MODEL"]
    A2 --> C["ACTION / VLA TRAINING"]
    A3 --> C
    A4 --> C
    A5 --> C
    B --> D["ACTION-CONDITIONED WORLD MODEL"]
    C --> D
    E["Real Video Initial Frames"] --> F["Image Editing Variations"]
    F --> A3
    D --> G["GLOBAL WORLD MODEL + VLA"]
    A3 --> H["Synthetic Failure Analysis"]
    H --> I["Targeted Real Egocentric Collection"]
    I --> B
    G --> J["HIGH-LEVEL PLANNER REASONING SFT"]
    J --> K["HIGH-LEVEL PLANNER REASONING RL"]
    K --> L["GLOBAL SYSTEM WITH REASONING-TRAINED PLANNER"]
    L --> M["REAL-TO-SIM AGENT"]
    M --> N["CALIBRATED DIGITAL TWIN"]
    N --> O["SURROGATE SIMULATOR"]
    O --> P["VLM TRAINING AGENT"]
    P --> Q["TASK CURRICULUM"]
    P --> R["DYNAMIC SUBTASK SCENARIOS"]
    P --> S["ENVIRONMENT VARIATIONS"]
    P --> T["INITIAL-STATE VARIATIONS"]
    Q --> U["AGENTIC SIMULATION RL"]
    R --> U
    S --> U
    T --> U
    U --> V["ENVIRONMENT-SPECIFIC ROBOT SYSTEM"]
    V --> W["HIGH-LEVEL DYNAMIC PLANNER"]
    W --> X["CURRENT SUBTASK"]
    X --> Y["LOW-FREQUENCY TARGET IMAGINATION"]
    Y --> Z["3 TARGET FUTURES"]
    Z --> AA["VLM TARGET SELECTION"]
    AA --> AB["SELECTED TARGET"]
    AB --> AC["HIGH-FREQUENCY ACTION LOOP"]
    AC --> AD["3 CANDIDATE ACTIONS"]
    AD --> AE["WORLD MODEL"]
    AE --> AF["3 PREDICTED FUTURES"]
    AF --> AG["VLM FUTURE COMPARISON"]
    AG --> AH["SAFETY CHECK"]
    AH --> AI["SELECT ACTION"]
    AI --> AJ["REAL DEPLOYMENT"]
    AJ --> AK["RUNTIME MONITORING"]
    AK --> AL["SUCCESSFUL EXPERIENCE"]
    AK --> AM["FAILURE EPISODES"]
    AM --> AN["VLM TARGETED LOCAL CURRICULUM"]
    AN --> AO["LOCAL SIMULATION RL"]
    AO --> V
    AM --> AP["LUNCH ROBOTICS ERROR-MODE ANALYSIS"]
    AP --> AQ["CURATED GLOBAL DATA"]
    AQ --> AR["NEW LAB VLA / WORLD MODEL / PLANNER / VLM"]
    AR --> AS["WEIGHT MIXING"]
    V --> AS
    AS --> AT["MIXED GLOBAL SYSTEM"]
    AT --> M
    AK --> W
```

The local environment lifecycle is:

```math
\text{Mixed Global System}
\rightarrow \text{Environment Reconstruction}
\rightarrow \text{System Identification}
\rightarrow \text{Calibrated Simulation}
\rightarrow \text{VLM Curriculum Agent}
\rightarrow
\left[
\begin{array}{c}
\text{Task Curriculum} \\
+ \\
\text{Dynamic Subtask Generation} \\
+ \\
\text{Environment Variations} \\
+ \\
\text{Initial-State Variations}
\end{array}
\right]
\rightarrow \text{Agentic RL}
\rightarrow \text{Environment-Specific Robot System}
\rightarrow \text{Dynamic Task Planning}
\rightarrow \text{3 Target Futures}
\rightarrow \text{VLM Target Selection}
\rightarrow \text{Selected Target}
\rightarrow \text{3 Candidate Actions}
\rightarrow \text{3 World Model Predictions}
\rightarrow \text{VLM Candidate Selection}
\rightarrow \text{Repeated High-Frequency Control}
\rightarrow \text{Dynamic Replanning}
\rightarrow \text{Real Deployment}
\rightarrow \text{Failure}
\rightarrow \text{VLM Targeted Local Curriculum}
\rightarrow \text{Local Simulation RL}
```

The global fleet lifecycle is:

```math
\text{Many Deployment Failures}
\rightarrow \text{Lunch Robotics Error-Mode Analysis}
\rightarrow \text{Curated Global Dataset}
\rightarrow \text{New Lab VLA / World Model / Planner / VLM}
\rightarrow \text{Weight Mixing}
\rightarrow \text{Mixed Global System}
\rightarrow \text{All Deployments}
```

The foundation-model lifecycle is:

```math
\text{Massive Real Egocentric Video}
\rightarrow \text{Passive Foundation World Model}
\rightarrow
\left[
\begin{array}{c}
\text{New Synthetic Egocentric Generation} \\
+ \\
\text{Real Initial Frames} \rightarrow \text{Image Editing} \rightarrow \text{Controlled Variations}
\end{array}
\right]
\rightarrow \text{Filtering + Action Extraction}
\rightarrow \text{Action-Conditioned Training}
\rightarrow \text{Synthetic Failure Analysis}
\rightarrow \text{Targeted Real Data Collection}
\rightarrow \text{Foundation World Model Update}
\rightarrow \text{Repeat}
```

The planner-training lifecycle is:

```math
\text{High-Level Planner Demonstrations}
\rightarrow \text{Planner Reasoning SFT}
\rightarrow \text{Planner RL Rollouts}
\rightarrow \text{Rewarded Subgoal Selection + Plan Revision}
\rightarrow \text{Reasoning-Trained Global Planner}
\rightarrow \text{Deployment Failure Analysis}
\rightarrow \text{Curated Planner Data}
\rightarrow \text{Improved Global Planner}
```

The pre-deployment environment-training lifecycle is:

```math
\text{Global Robot System with Reasoning-Trained Planner}
\rightarrow \text{Real-to-Sim}
\rightarrow \text{Calibrated Expected Environment}
\rightarrow \text{VLM Curriculum Agent}
\rightarrow
\left[
\begin{array}{c}
\text{Task Curriculum} \\
+ \\
\text{Dynamic Subtasks} \\
+ \\
\text{Environment Variations} \\
+ \\
\text{Initial-State Variations}
\end{array}
\right]
\rightarrow \text{Agentic RL}
\rightarrow \text{Robust Environment-Specific Robot System}
\rightarrow \text{High-Fidelity Validation}
\rightarrow \text{Real Deployment}
```

The resulting architecture is not a one-way pipeline. It is a **closed learning system** in which the foundation model learns from reality, the World Model generates synthetic hypotheses about reality, real-video-seeded image editing creates controlled variations around genuine observations, synthetic data expands action-conditioned training, model failures identify missing knowledge, targeted real-world collection corrects those weaknesses, planner reasoning training teaches the high-level VLM how to choose and create useful subgoals, the Real-to-Sim system constructs the expected deployment world, a VLM training agent constructs the curriculum and variation distribution, the high-level runtime planner dynamically decomposes and replans the task, agentic RL trains competence and recovery, target imagination defines the medium-horizon objective, VLM target selection determines which future is desirable, the VLA generates multiple local action alternatives, the World Model predicts their consequences, VLM comparison selects the best local action, deployments expose the residual weaknesses of those policies, and the company converts the most generalizable weaknesses into improvements to the shared brain.

---

# 25. The Core Research Thesis

The central thesis of Lunch Robotics is that universal robot intelligence should be built from complementary components:

```math
\boxed{
\text{Foundation Model Data Funnel}
+
\text{Generative Egocentric Data Flywheel}
+
\text{Action-Conditioned World Model}
+
\text{High-Level Planner Reasoning SFT + RL}
+
\text{Agentic Real-to-Sim}
+
\text{VLM-Agentic Environment RL}
+
\text{Dynamic Hierarchical Planning}
+
\text{Curated Fleet Learning}
}
```

The Foundation Model Data Funnel solves the first problem: **how do we learn broad physical and manipulation intelligence without requiring enormous quantities of expensive robot data?**

The answer is to combine data sources with radically different scale and fidelity inside one continuously co-trained foundation model:

```math
\text{Massive Human Video}
+
\text{Human Action Learning}
+
\text{Synthetic Human Action Data}
+
\text{Robot-Relevant Manipulation}
+
\text{Teleoperation}
\rightarrow
\text{Global World Model + VLA}
```

The generative egocentric-data flywheel solves the next problem: **how do we scale action-conditioned experience and determine which missing parts of reality are most important to collect?**

```math
\text{Foundation World Model}
\rightarrow
\left[
\begin{array}{c}
\text{New Human Video Generation} \\
+ \\
\text{Real Initial Frames} \rightarrow \text{Image-Edited Variations}
\end{array}
\right]
\rightarrow \text{Synthetic Human Videos}
\rightarrow \text{Filtering}
\rightarrow \text{Pose / Action Extraction}
\rightarrow \text{Action-Conditioned Training}
```

and:

```math
\text{Synthetic Failures}
\rightarrow \text{Weakness Analysis}
\rightarrow \text{Targeted Real Data Collection}
\rightarrow \text{Better Foundation World Model}
```

The essential principle is:

> **Synthetic data scales the model's existing beliefs; real-world data determines whether those beliefs are correct.**

The Action-Conditioned World Model solves the next problem: **how does the robot predict the consequences of an action before committing to it?**

```math
\text{Current State}
+
\text{Action}
\rightarrow
\text{Predicted Future Trajectory}
```

Target future imagination solves the next problem: **how does the robot decide what future it should be moving toward when there are multiple valid ways to solve the current subtask?**

```math
\text{Current State}
+
\text{Current Subtask}
\rightarrow
\text{3 Target Futures}
\rightarrow
\text{VLM Ranking}
\rightarrow
\text{Selected Target Future}
```

Reasoning SFT and RL solve another problem: **how do we teach the system to reason about long-horizon tasks, identify missing prerequisites, and decide what the robot should accomplish next?** The training target is the **high-level VLM planner**. Planner SFT teaches it to construct and revise hierarchical plans; planner RL rewards decisions that produce appropriate subgoals, discover necessary intermediate objectives, and lead to successful downstream task outcomes.

```math
\text{Task}
+
\text{World State}
+
\text{Execution History}
\rightarrow
\text{Planner Reasoning}
\rightarrow
\text{Updated Hierarchical Plan}
\rightarrow
\text{Next Subgoal}
```

This is distinct from the VLA's action-generation problem:

```math
\text{Observation}
+
\text{Current Subtask}
+
\text{Selected Target Future}
\rightarrow
\text{Candidate Action Chunks}
```

The planner decides **what should be accomplished next**. The VLA proposes **how to physically act toward the selected subgoal**. The World Model predicts the consequences of those actions, and the VLM evaluation loop selects among them.

Dynamic hierarchical planning solves another key problem: **how do we make long-horizon task execution adaptive when the necessary subtask sequence is not fully known in advance?**

```math
\text{Task}
+
\text{Current State}
\rightarrow
\text{Plan}
\rightarrow
\text{Subtask}
\rightarrow
\text{Observation}
\rightarrow
\text{Replan}
\rightarrow
\text{New / Modified Subtask}
```

The planner can therefore discover intermediate objectives at runtime rather than relying on a static plan generated once at the beginning.

The Real-to-Sim system solves the next problem: **how do we adapt that general intelligence to an arbitrary robot operating in an arbitrary physical environment?**

```math
\text{Real Environment}
\rightarrow
\text{Digital Twin}
\rightarrow
\text{System Identification}
\rightarrow
\text{Calibrated Simulation}
```

Stage 8 then solves a separate but crucial problem:

> **How do we make the environment-specific robot robust before it ever enters the real world?**

The answer is to put the model into a calibrated simulation and have a VLM training agent continually generate:

```text
Task Curriculum
+
Dynamic Subtasks
+
Environment Variations
+
Initial-State Variations
```

and optimize the robot through agentic RL:

```math
\text{Calibrated Environment}
+
\text{Curriculum}
+
\text{Dynamic Planning}
+
\text{Variations}
+
\text{Initial States}
\rightarrow
\text{Agentic RL}
\rightarrow
\text{Robust Environment-Specific Robot System}
```

The crucial distinction is:

> **Stage 8 is the pre-deployment version of local online RL.**

It teaches the robot progressively harder tasks while continuously varying the physical environment and initial task state. It also trains and evaluates the planner's ability to adapt the task decomposition when unexpected intermediate conditions arise. This gives the system **robustness to environmental variation**, **the ability to recover when it does not start from the expected state**, and **the ability to dynamically create intermediate subtasks when the original plan is insufficient**.

The VLM training agent therefore operates above the training episodes, while the runtime planner operates inside each episode:

```text
VLM Training Agent
       ↓
chooses what experiences the robot should learn from

High-Level VLM Planner
       ↓
chooses what the robot should accomplish next

Robot VLA
       ↓
learns how to execute the selected subgoal physically
```

The candidate-action mechanism solves the next problem: **how do we avoid relying on a single VLA prediction to make every decision perfectly?**

```math
\{a_t^{(1)},a_t^{(2)},a_t^{(3)}\}
\rightarrow
\{z_{\mathrm{pred}}^{(1)},z_{\mathrm{pred}}^{(2)},z_{\mathrm{pred}}^{(3)}\}
\rightarrow
\text{VLM Comparison}
\rightarrow
a_t^*
```

The action whose predicted future best matches the selected target future is executed:

```math
a_t^*
=
a_t^{\left(
\arg\max_i
\mathrm{VLM}_{\mathrm{score}}
\left(z_{\mathrm{pred}}^{(i)},z_{\mathrm{target}}\right)
\right)}
```

The deployment learning system solves the final problem: **how do we continue improving after the robot is operating in the real world?**

The answer is to let each deployment learn locally while sending rich failure data back to the central lab. The Lunch Robotics team analyzes those failures across environments, identifies generalizable error modes, and curates the datasets required to teach those capabilities to the global model and planner.

```math
\text{Deployment Failures}
\rightarrow
\text{Team Error-Mode Analysis}
\rightarrow
\text{Curated Data}
\rightarrow
\text{New Lab VLA / World Model / Planner / VLM}
```

Finally, the global and local systems are combined and redistributed:

```math
\text{Global Lab System}
+
\text{Environment-Specific Systems}
\rightarrow
\text{Mixed System}
\rightarrow
\text{All Deployments}
\rightarrow
\text{Environment-Specific Agentic RL}
```

The most important architectural principle is therefore:

> **The robots specialize locally, while Lunch Robotics learns globally.**

A deployed robot is not merely a consumer of a fixed model. It is an autonomous learning agent operating inside a particular physical environment and a sensor collecting valuable evidence about what the shared brain still does not understand.

At the foundation level, the World Model itself also acts as a data-generation and diagnosis engine: it produces hypothetical human experiences, uses real observations as seeds for controlled variations, exposes its own blind spots, and directs the organization toward the real-world data needed to improve its understanding of physical reality.

At the environment-training level, the VLM training agent performs the analogous role for deployment-specific learning: it constructs the curriculum, varies the environment, varies the initial state, creates scenarios requiring dynamic replanning, and continually adapts the training distribution to the current robot system.

At runtime, the high-level planner performs the analogous role for execution: it turns an abstract goal into a living hierarchy of subtasks, observes what actually happens, and rewrites the plan whenever reality invalidates the current assumptions.

This creates a compounding learning system:

```math
\boxed{
\text{More Real Data}
\rightarrow \text{Better Foundation Model}
\rightarrow \text{More Useful Synthetic Experience}
\rightarrow \text{Better Action Conditioning}
\rightarrow \text{Better World Model}
\rightarrow \text{Better Planner Reasoning}
\rightarrow \text{Better Planning}
\rightarrow \text{Better Target Imagination}
\rightarrow \text{Better VLM Evaluation}
\rightarrow \text{Better Simulation Training}
\rightarrow \text{More Robust Robot}
\rightarrow \text{Better Real Deployment}
\rightarrow \text{More Real-World Experience}
}
```

and across the fleet:

```math
\boxed{
\text{More Deployments}
\rightarrow \text{More Real-World Experience}
\rightarrow \text{More Error Modes Discovered}
\rightarrow \text{Better Curated Data}
\rightarrow \text{Better Lab VLA}
\rightarrow \text{Better World Model}
\rightarrow \text{Better Planner}
\rightarrow \text{Better Target Imagination}
\rightarrow \text{Better VLM Evaluation}
\rightarrow \text{Better Agentic Simulation RL}
\rightarrow \text{Better Initialization}
\rightarrow \text{Better Environment Adaptation}
\rightarrow \text{Better Robots}
\rightarrow \text{More Deployments}
}
```

---

# 26. The Product Vision

The end state is **zero-to-hero robot autonomy**. The user provides a robot, an environment, the tasks it needs to perform, a small amount of tactile probing, and task demonstrations. The rest of the system operates behind the scenes.

```math
\begin{aligned}
&\text{Massive Human Video}
+\text{Human Action Data}
+\text{Synthetic Human Action Data}
+\text{Robot-Relevant Data}
+\text{Teleoperation} \\
&\qquad\downarrow \\
&\text{Global World Model + VLA} \\
&\qquad\downarrow \\
&\text{Action-Conditioned World Model} \\
&\qquad\downarrow \\
&\text{High-Level Planner Reasoning SFT + RL} \\
&\qquad\downarrow \\
&\text{Reasoning-Trained Global Planner} \\
&\qquad\downarrow \\
&\text{Real-to-Sim Agent} \\
&\qquad\downarrow \\
&\text{Calibrated Digital Twin} \\
&\qquad\downarrow \\
&\text{VLM Curriculum Agent} \\
&\qquad\downarrow \\
&\text{Task + Dynamic Subtasks + Environment + Initial-State Variation} \\
&\qquad\downarrow \\
&\text{Massive Agentic Simulation RL} \\
&\qquad\downarrow \\
&\text{Robust Environment-Specific Robot System} \\
&\qquad\downarrow \\
&\text{Dynamic Hierarchical Planning} \\
&\qquad\downarrow \\
&\text{Current Subtask} \\
&\qquad\downarrow \\
&\text{LOW-FREQUENCY TARGET IMAGINATION} \\
&\qquad\downarrow \\
&\text{3 Target Futures} \\
&\qquad\downarrow \\
&\text{VLM Target Selection} \\
&\qquad\downarrow \\
&\text{Selected Target Future} \\
&\qquad\downarrow \\
&\text{HIGH-FREQUENCY ACTION SELECTION} \\
&\qquad\downarrow \\
&\text{3 Candidate Action Chunks} \\
&\qquad\downarrow \\
&\text{3 World Model Predicted Futures} \\
&\qquad\downarrow \\
&\text{VLM Future Comparison} \\
&\qquad\downarrow \\
&\text{Selected Action} \\
&\qquad\downarrow \\
&\text{Real Deployment} \\
&\qquad\downarrow \\
&\text{Failure Discovery} \\
&\qquad\downarrow \\
&\text{Targeted Local Learning}
\end{aligned}
```

Across the fleet, deployment failures are converted into curated improvements to the global model system:

```math
\begin{aligned}
&\text{Deployment Failures} \\
&\qquad\downarrow \\
&\text{Lunch Robotics Error-Mode Analysis} \\
&\qquad\downarrow \\
&\text{Curated Global Training Data} \\
&\qquad\downarrow \\
&\text{New Lab VLA + World Model + Planner + VLM} \\
&\qquad\downarrow \\
&\text{Weight Mixing} \\
&\qquad\downarrow \\
&\text{Mixed Global System} \\
&\qquad\downarrow \\
&\text{Every Deployment}
\end{aligned}
```

At deployment, the user can simply say:

> **"Clean the table."**

The robot understands the task, retrieves the relevant skill, and hands the task to the high-level VLM planner. The planner creates an initial hierarchical plan, determines the next subgoal, and reasons for as long as the task's complexity and uncertainty justify. The runtime system then operates through the nested planning and control loops. It periodically imagines three possible medium-horizon target futures, uses a VLM to select the best target, generates three plausible action chunks, predicts one future for each through the World Model, compares those three imagined futures against the target using the VLM, selects the action associated with the best future, executes it through a hardware-safe control stack, and monitors the result.

Meanwhile, the high-level planner watches the execution state and can dynamically replan. It can discover an unexpected prerequisite, create a new subtask, abandon an invalid branch, reorder remaining subtasks, or switch to a recovery sequence without requiring a new human instruction. **The planner's Stage 7 training is what teaches it to make these semantic planning decisions; the VLA's role remains action generation.**

The selected target is not regenerated for every individual action chunk. The robot repeatedly performs the fast action-selection loop against that target and periodically generates a new set of target futures as the subtask progresses. The high-level planner sits above this loop and can invalidate the current target by changing the active subtask.

Before reaching this deployment point, however, the robot has already undergone a dedicated environment-specific curriculum. A VLM training agent has generated progressively more difficult tasks, varied the physical environment around the calibrated deployment world, and initialized tasks from different states, including partially completed and degraded states. The policy therefore enters reality having already learned not only how to perform the expected tasks, but how to operate under the kinds of variation and recovery conditions it is expected to encounter.

When the robot encounters something it does not understand, the system does not simply record a failure and move on. It captures the entire event, uses the local VLM training agent to generate targeted simulation scenarios, and learns locally from the failure. The information also goes back to the Lunch Robotics lab, where the team determines whether the failure reveals a broader capability, planner reasoning, planning, target-imagination, prediction, VLM-ranking, or action-selection gap.

The team then curates the appropriate training data, improves the global VLA and/or World Model and/or high-level planner and/or VLM, mixes the new global knowledge with the knowledge accumulated by specialized deployments, and redistributes the resulting system across the fleet.

At the foundation level, the same philosophy applies before deployment. Real egocentric video teaches the World Model how the physical world behaves; the World Model then generates large quantities of hypothetical human manipulation experience and can also use real video initial frames as seeds for an image-editing model that creates controlled scene variations and slight task variations; humans and learned critics identify where those generations are implausible; pose-estimation models turn the surviving generations into action supervision; and the resulting failure patterns determine which real-world experiences should be collected next.

The long-term objective is not to build another robot-specific policy. It is to build a **universal brain that can be installed into any robot, adapted to any physical environment, reason at the appropriate level for any task, dynamically decompose and replan long-horizon tasks, learn broad physical intelligence from massive real-world experience, use its World Model to generate and interrogate synthetic human experience, use real video as a seed for controlled synthetic variation, identify the limits of its own understanding, guide targeted real-world data collection, adapt itself to an expected deployment environment through VLM-guided curriculum learning, become robust through large-scale simulation RL across environmental and initial-state variations, recover from unexpected states, periodically imagine multiple possible desirable futures, select the best target future, generate multiple candidate action chunks at high frequency, simulate their consequences through the World Model, use a VLM to choose the action whose predicted outcome best matches the currently selected future, continuously revise its high-level plan and create new subtasks when needed, and continuously improve through the collective experience of every robot running it.**

The complete philosophy can be summarized as:

```text
REALITY
   ↓
Foundation Learning
   ↓
World Model
   ↓
Synthetic Hypotheses
   ↓
Targeted Real Data
   ↓
Global World Model + VLA
   ↓
Planner Reasoning SFT + RL
   ↓
Reasoning-Trained High-Level Planner
   ↓
Real-to-Sim
   ↓
VLM-Generated Curriculum
   ↓
Dynamic Subtask Training
   ↓
Environment Variations
   +
Initial-State Variations
   ↓
Agentic Simulation RL
   ↓
Robust Environment-Specific Robot
   ↓
Dynamic Hierarchical Planning
   ↓
Target-Future Selection
   ↓
Action Selection
   ↓
Real Deployment
   ↓
Residual Failures
   ↓
Local Recovery Learning
   ↓
Global Error-Mode Analysis
   ↓
Curated Global Improvement
   ↓
All Robots
```

The fundamental insight is:

> **Reality provides the grounding. World Models generate hypotheses. VLM training agents decide what to learn. The high-level VLM planner decides what to accomplish. The VLA proposes how to act. RL improves decisions through outcomes. Simulation provides scale. Deployment provides new evidence. Lunch Robotics turns the most generalizable evidence into a better universal brain.**
