# Lunch Robotics: Universal Brain for Any Robot

## The Vision

We propose an end-to-end engine that turns a useless robot into a useful and self-improving robot, requiring only:

- **Tactile Environment Probing Dataset**: Video and synchronized contact data from a human operating a handheld gripper to interact with, tap, deform, push, and manipulate environment objects and surfaces.
- **Task Demonstration Dataset**: One video per task of a human performing the task.
- **Robot Specifications**: The robot's SDK and hardware specification, including kinematics, joint limits, and end-effector dynamics.

The user should not need to manually build a simulator, design an RL environment, engineer a training curriculum, or train the robot policy.

### Required Setup

- Interactive Environment + Task Data collected via the tactile-sensing handheld gripper.
- Robot SDK installed on the robot.
- Robot Hardware Specification exposed through the SDK.
- Optional environment metadata and hints.

The engine automatically constructs the simulation environment, calibrates its physics, trains the robot policy, and adapts it to the physical hardware.

```mermaid
flowchart TD
    A["Interactive Tactile Probing + Task Data"] --> D["Real-to-Sim Agent"]
    B["Robot SDK + Specs"] --> D
    C["Environment Metadata"] --> D

    D --> E["Calibrated Digital Twin"]
    E --> F["Surrogate Model"]
    F --> G["Simulation RL"]
    G --> H["Real-World RL"]
    H --> I["Trained Robot Policy"]
    I --> J["Deployment"]
```

At deployment time, the user simply gives the robot a task via natural language.

**Example:**  
*"Clean the table."*

The task is captured through the robot's audio sensor, converted to text using speech recognition, and passed to the robot policy through the SDK.

*(Note: To safeguard initial operations, homes are provided with non-breaking cups and dishes to use during the first month of robot usage.)*

### The Goal

$$
\text{Robot} + \text{Tactile Video/Contact Data} + \text{SDK} + \text{Task} \rightarrow \text{Autonomous Robot}
$$

---

## 0. Foundation: Pre-Existing Base Models and General SFT Phase

Before initiating any environment-specific procedures, calibration, or deployment workflows, the entire engine is built upon a robust foundational training phase. This phase establishes broad, cross-domain capabilities across diverse scenarios before any targeted adaptation takes place.

### 0.1 Pre-Existing Base Models
The architecture rests upon a suite of powerful pre-existing foundational models that provide core priors across modalities:
- **Base Robot Action Model**: Supplies foundational motor primitives and generalized manipulation skills.
- **Base 3D Editing Model**: Powers spatial reasoning, scene manipulation, and generative latent state modifications.
- **Base 2D World Model**: Captures visual dynamics, temporal transitions, and environmental continuity.
- **Base Agent Coding LLM**: Drives tool use, reasoning, simulation code generation, and error diagnosis.

### 0.2 General SFT Phase
On top of these pre-existing base models, the system undergoes an extensive **Supervised Fine-Tuning (SFT) training phase**. This phase exposes the models to vast collections of:
- **General Environments**: Diverse synthetic and recorded spatial layouts, material types, and physics profiles.
- **Task Variations**: A comprehensive spectrum of manipulation, locomotion, and interaction primitives.
- **Scenario Variations**: Edge cases, dynamic obstacles, multi-agent interactions, and varying lighting or spatial geometries.

This general SFT phase ensures that the models possess deep semantic understanding, robust tool-use proficiency, and broad spatial-temporal reasoning. Consequently, when a user provides specific tactile probing data and robot specifications for a deployment environment, the system does not learn from scratch; instead, it leverages these generalized priors to rapidly specialize into a high-fidelity digital twin and policy.

---

## 1. The Real-to-Sim Agent

The digital twin is constructed by an autonomous agent rather than a fixed reconstruction pipeline.

The agent is a purposefully RL-trained LLM agent whose task is to convert real-world multimodal interaction data into a sensible, executable, and calibrated simulation.

The agent is equipped with a toolbox for:

- Multimodal video and tactile/contact signal parsing
- 3D scene and object reconstruction
- Simulator construction
- Physics simulation and contact modeling
- Trajectory and wrench/force extraction
- System identification and material property estimation
- Parameter optimization
- Simulation evaluation
- Experiment design
- World-model integration
- Code generation and execution

```mermaid
flowchart TD
    A["Tactile Video + Contact Data"] --> B["Real-to-Sim LLM Agent"]

    B --> C["Geometry"]
    B --> D["Objects"]
    B --> E["Agents"]
    B --> F["Physics & Contact Parameters"]
    B --> G["Materials & Compliance"]
    B --> H["World Model"]

    C --> I["Initial Digital Twin"]
    D --> I
    E --> I
    F --> I
    G --> I
    H --> I

    I --> J["System Identification Tools"]
    J --> K["Calibrated Digital Twin"]
    K --> L["Validation"]

    L -->|"Failure"| B
    L -->|"Success"| M["Finalize"]
```

The agent is responsible for constructing and validating the entire digital twin. Integrating concepts from generative simulation (e.g., RoboGen), the agent utilizes a self-guided propose-generate-learn cycle to automatically scale diverse task creation and scene generation. However, unlike pure LLM-based generators that lack physical accuracy, our agent grounds this generation in multimodal tactile data. It can:

- Inspect the multi-view video and synchronized contact feedback.
- Autonomously propose new sub-tasks and populate pertinent assets to create diversified training environments.
- Identify objects, geometry, and surface contact dynamics to anchor the generative simulation.
- Select simulation primitives and contact models, such as rigid, soft-body, or elastoplastic models.
- Determine which components require learned world-model dynamics.
- Create simulator code and physical property configurations.
- Identify uncertain physical parameters, such as mass, friction, compliance, and damping.
- Invoke specialized system-identification tools against contact trajectories.
- Compare simulated contact forces and kinematics with physical data.
- Modify and tune the simulator until the twin reaches the required fidelity.

The agent acts as an AI simulation engineer and system-identification engineer, yielding infinite, yet physically grounded, training data.

---

## 2. Agentic System Identification

Passive video alone leaves physical properties fundamentally under-constrained. Synchronized tactile and force data collected during human exploration removes these ambiguities.

The Real-to-Sim Agent uses specialized system-identification tools to optimize the simulator against multimodal real-world trajectories.

The agent first infers approximate physical bounds:  
*"This object appears to be a wooden block."*

It then estimates the numerical parameter vector:

$$
\theta = \begin{bmatrix} m \\ \mu_s \\ \mu_d \\ e \\ k \\ c \\ \vdots \end{bmatrix}
$$

where the parameters include:

- $m$: mass  
- $\mu_s$: static friction coefficient  
- $\mu_d$: dynamic friction coefficient  
- $e$: coefficient of restitution  
- $k$: contact stiffness / compliance  
- $c$: damping factor  

```mermaid
flowchart TD
    A["Real-World Contact Trajectories<br/>(Video + Force/Tactile Data)"] --> B["Real-to-Sim Agent"]

    B --> C["Identify Uncertain Physics Parameter"]
    C --> D["System ID Tool"]
    D --> E["Parameter Optimization"]
    E --> F["Simulator Rollout"]
    F --> G["Compare Kinematics + Contact Wrench"]
    G --> H["Multimodal Loss/Error"]
    H --> I["Agent Diagnoses Mismatch"]

    I -->|"Revise"| B
    I -->|"Sufficient Fidelity"| J["Finalize"]
```

The optimization objective minimizes both state kinematic trajectories $\tau$ and contact wrench trajectories $W$:

$$
\theta^* = \arg\min_{\theta} \left( D_{\text{kin}}\left(\tau_{\mathrm{real}}, \tau_{\mathrm{sim}}(\theta)\right) + \lambda D_{\text{force}}\left(W_{\mathrm{real}}, W_{\mathrm{sim}}(\theta)\right) \right)
$$

The LLM agent reasons about physical discrepancies at a high level, selects parameters to unfreeze, and uses optimization tools to solve for physical parameters.

---

## 3. Agentic Entities Are Simulated by a General World Model

Entities that exhibit complex, non-scripted behaviors, such as humans, pets, or other dynamic agents, are simulated using a general world model.

Rather than requiring internal state access, the model predicts future states directly from observation histories:

$$
S^{\mathrm{agent}}_{t+1:t+H} \sim P_{\phi}\left(S^{\mathrm{agent}}_{t+1:t+H} \mid S_{\leq t}, E\right)
$$

The simulator composes deterministic physics with learned behavioral dynamics:

$$
\text{Digital Twin} = \text{Explicit Physics (Calibrated via Tactile Data)} + \text{General World Model}
$$

This allows the digital twin to combine highly calibrated physical interactions with learned dynamics for entities whose behavior cannot be efficiently described through explicit physics alone.

---

## 4. Surrogate Models Accelerate Simulation

Once the digital twin has been calibrated using physical interaction data, it serves as the teacher for a learned neural surrogate simulator.

$$
f_{\mathrm{sim}}(s_t, a_t) \approx f_{\mathrm{surrogate}}(s_t, a_t)
$$

```mermaid
flowchart TD
    A["Calibrated Digital Twin"] --> B["Generate Simulation Data"]
    B --> C["Synthetic Trajectories"]
    C --> D["Train Surrogate Model"]
    D --> E["Fast Approximate Simulator"]

    E --> F["Large-Scale RL"]
    F --> G["Candidate Policy"]

    G --> H["Validate in High-Fidelity Twin"]
    H -->|"Mismatch"| D
    H -->|"Sufficient Accuracy"| I["Continue Training"]
```

The surrogate is not intended to replace the high-fidelity simulator completely.  
Instead, it provides a much faster approximation that allows orders of magnitude more policy rollouts during RL.

The hierarchy operates as:

$$
\text{Real World (Tactile + Video)} \rightarrow \text{Calibrated Digital Twin} \rightarrow \text{Surrogate Model} \rightarrow \text{Fast Simulation RL}
$$

---

## 5. Multimodal Data Collection Protocol

To eliminate system-identification ambiguities while keeping setup overhead minimal, data collection uses a specialized, sensor-equipped handheld gripper.

The human operator manually interacts with the environment to explicitly ground physical properties.

```mermaid
flowchart TD
    A["Handheld Tactile Gripper"] --> B["Environment Interaction Data"]
    A --> C["Task Demonstration Data"]

    B --> D["Physical Surface Sweeps<br/>(Video + Contact/Force/Friction Data)"]
    B --> E["Object Dynamics Exploration<br/>(Pushes, Taps, Lifts, Deformations)"]

    C --> F["Human Task Execution Video"]

    D --> G["Real-to-Sim SysID Engine"]
    E --> G
    F --> H["Goal & Reward Formulation"]
```

### Data Streams Collected

#### 1. Interactive Environment Probing

A human holds a sensorized gripper and systematically "plays" with the target workspace.

The operator:

- Taps surfaces
- Slides across tables
- Lifts objects
- Pushes objects
- Squeezes objects
- Deforms compliant materials
- Manipulates objects under different contact conditions

Collected signals include:

- RGB-D video streams
- High-frequency contact forces
- Tactile pressure maps
- Slipping indicators
- End-effector pose
- Interaction forces and torques

#### 2. Task Demonstration Video

A clean, uninterrupted recording of a human performing the goal task.

Examples include:

- Sorting objects
- Wiping a counter
- Picking up fragile items
- Manipulating articulated objects

This structured exploration grounds object mass, friction profiles, joint resistance, surface stiffness, and compliance directly from real-world contact events before policy training begins.

---

## 6. On-the-Fly Learning from Human Feedback

Deployment is continuous and adaptive.

When the robot makes an error in the physical world, the human operator provides on-the-fly corrections via natural language, voice, or spatial gestures.

```mermaid
flowchart TD
    A["Human Feedback on Real-World Mistake"] --> B["Robot Brain Parses Failure Mode"]
    B --> C["Instantiate Sub-Task RL Env in Calibrated Twin"]
    C --> D["Targeted RL Policy Fine-Tuning"]
    D --> E["VLM Judge Evaluates Rollout"]
    E -->|"Mistake Unresolved"| C
    E -->|"Mistake Resolved"| F["Deploy Updated Policy"]
```

### Self-Correction Loop

1. **Feedback Parsing**: The robot's central reasoning engine translates natural-language critique, such as *"You're squeezing the paper cup too hard and crushing it,"* into concrete physical parameter bounds and failure predicates.
2. **Automated RL Sub-Environment Generation**: The digital twin instantiates a specialized micro-simulation environment replicating the precise spatial state, contact properties, and dynamic parameters where the failure occurred.
3. **Simulated Policy Refinement**: The policy undergoes focused RL updates inside the sub-environment to learn corrected wrench and trajectory profiles.
4. **VLM Evaluation**: A Vision-Language Model (VLM) evaluates simulated execution videos to verify that the robot successfully achieves the objective without triggering the identified failure mode.
5. **Redeployment**: The updated policy parameters are hot-swapped onto the physical robot.

## 7. Training the Robot Policy

Whether the system is learning a pre-known task from initial demonstrations or fine-tuning to correct mistakes on-the-fly, policy optimization proceeds using a unified, large-scale RL approach.

The robot policy is represented as:

$$
\pi_{\theta}(a_t \mid o_t, T)
$$

where:

- $a_t$ is the action at time $t$
- $o_t$ is the robot's observation
- $T$ is the task specification (or natural language correction)
- $\theta$ represents the policy parameters

The surrogate handles massive rollout volume, while the calibrated digital twin verifies contact-heavy manipulation steps to ensure physical grounding.

### 7.1 Unified Reward Generation via Imagined Motions

Strict kinematic reference trajectories heavily overfit to specific initial scenario settings and are far too brittle for generalized learning. Instead, for robust learning across both new and pre-known tasks, the reward is formulated based on achieving specific **imagined target motions** (representing latent state trajectories over a very small temporal delta $\Delta t$), rather than static snapshots. 

Crucially, generating these latent motion trajectories is not a novel paradigm; latent motions are already the native output of most modern video generation models out-of-the-box.

The pipeline is as follows:

1. **Imagined Motion Generation:** A 3D Vision-Language Model (VLM) diffusion model acts as an editing model. It processes the initial state and the task description (or language correction) to output a target **embedded latent motion trajectory** ($z_{\text{target}(t:t+\delta)}$) representing the successful temporal sequence, leveraging off-the-shelf video generation model outputs.
2. **Latent Space Embedding:** During simulated rollouts, the current state trajectory segment $s_{t:t+\delta}$ is continuously embedded into the same latent space to produce $z_{t:t+\delta}$.
3. **Latent Similarity Reward:** The system issues a reward based on the similarity between the current embedded motion trajectory $z_{t:t+\delta}$ and the imagined target latent motion $z_{\text{target}(t:t+\delta)}$. The rawrd is only give, if and only if, that specifc reward was not given before for that trajectory.
4. **Temporal Maintenance (Delta Time):** For tasks requiring sustained action (e.g., holding or carrying), the reward requires the latent similarity to remain above a threshold across the window $\Delta t$.
5. **Destruction Penalty:** A VLM judge continuously observes the simulated rollouts to detect severe or destructive results—such as knocking a glass over or crushing a paper cup. If an unsafe event occurs, a massive penalty is applied to strictly forbid that behavior.

The step reward $r_t$ is formulated in the latent space as:

$$
r_t = w_{\text{goal}} \cdot \text{sim}(z_{t:t+\delta}, z_{\text{target}(t:t+\delta)}) - w_{\text{penalty}} \mathbb{I}_{\text{fail}}(s_t)
$$

where:
* $z_{t:t+\delta} = \text{Encoder}(s_{t:t+\delta})$ is the latent embedding of the current simulator state trajectory.
* $z_{\text{target}(t:t+\delta)}$ is the target latent motion generated by the 3D VLM diffusion editing model.
* $\text{sim}(\cdot, \cdot)$ is a similarity metric (such as cosine similarity) in the embedded latent space.
* $\mathbb{I}_{\text{fail}}(s_t) \in \{0, 1\}$ is a Boolean indicator evaluated by the VLM for destructive/catastrophic failures.
* $w_{\text{goal}}$ and $w_{\text{penalty}}$ are weighting constants.

---

## 8. Curriculum Learning

Training scales from basic object manipulation to complex tasks across increasingly difficult environmental variables.

```mermaid
flowchart TD
    A["Simple Single-Object Contact"] --> B["Multi-Object Clutter"]
    B --> C["Precise / Compliant Manipulation"]
    C --> D["Long-Horizon Sequences"]
    D --> E["Dynamic Environments"]
    E --> F["Full Task Autonomy"]
```

The curriculum progressively increases:

- Number of interacting objects
- Contact complexity
- Precision requirements
- Material variability
- Environmental dynamics
- Task horizon
- Degree of uncertainty

---

## 9. Real-World RL Closes the Sim-to-Real Gap

A final adaptation stage grounds the policy directly on physical hardware.

```mermaid
flowchart TD
    A["Surrogate / Twin RL Policy"] --> B["Real-World RL Adaptation"]
    B --> C["VLM Automated Task Generation & Reward"]
    C --> D["Fine-Tuned Hardware Policy"]
```

Real-world RL adaptation to close the sim-to-real gap is executed **just once every 3 months**. During this periodic phase, the human asks the robot to perform a task and provides natural language feedback explaining what it did wrong, driving the adaptation loop:

$$
\text{Task Generation} \rightarrow \text{Real Interaction} \rightarrow \text{VLM Evaluation} \rightarrow \text{Reward} \rightarrow \text{RL Update}
$$

This infrequent stage addresses residual discrepancies between the calibrated simulation and the deployed hardware.

---

## 10. Safety Guardrails and Adversarial Risk Mitigation

To ensure that the robot never performs catastrophic or strictly forbidden behaviors—such as attempting to retrieve items from a cabinet full of delicate glassware—the system implements a rigorous safety alignment pipeline combining Supervised Fine-Tuning (SFT) and specialized Reinforcement Learning.

### 10.1 Safety SFT and Adversarial Twin RL
Before deployment, the system undergoes extensive safety Supervised Fine-Tuning followed by targeted Reinforcement Learning inside the specific digital twin:
- **Adversarial Training**: The robot is intentionally placed in simulation environments where it is adversarially instructed to perform unsafe or destructive tasks.
- **Penalty Optimization**: Similar to how it learns desirable behaviors, the policy receives massive penalties for hazardous actions. Its explicit optimization objective is to achieve zero penalties even when coaxed or instructed to execute risky routines.
- **Resulting Behavioral Bounds**: This conditioning ensures the robot develops robust behavioral boundaries against executing hazardous or forbidden actions in real homes.

### 10.2 Planner LLM Rejection of Risky Tasks
Beyond low-level policy safeguards, high-level task execution relies on a central planner LLM:
- **Intent Screening**: The high-level planner LLM responsible for parsing user intent and defining sub-tasks is explicitly trained and fine-tuned to recognize and **reject** dangerous or risky tasks.
- **Safe Fallback**: When presented with a harmful or unsafe request (e.g., handling delicate items in unstable storage zones), the planner LLM refuses the instruction and prompts the user for clarification or a safe alternative, preventing the downstream policy from ever generating execution trajectories.

---

## 11. Robot SDK

The deployed policy maps high-level human speech to robot control through the SDK. The user interacts with the robot at the level of tasks rather than low-level robot commands, letting the SDK bridge the gap between intention and physical actuation.

### 11.1 SDK Motion Smoothing

Raw neural network policies often produce high-frequency, noisy, or abrupt action sequences. If sent directly to the hardware, these commands can cause jittery movement or damage the physical joints. To ensure safe execution, the Robot SDK acts as a protective translation layer incorporating a real-time Motion Smoother.

To ensure low-latency decision-making on resource-constrained physical hardware, the deployed diffusion policy undergoes **Consistency Distillation**. By enforcing self-consistency along the diffusion policy's learned trajectories, the system creates a *Consistency Policy* that can generate an action sequence in a single inference step. This speeds up action generation by an order of magnitude compared to standard iterative denoising.

Furthermore, every deployed policy automatically benefits from motion smoothing.

The smoother takes the high-speed intended action input ($a_t$) and context from the action history ($a_{<t}$) to perform online interpolation and jerk-limited trajectory generation.

$$
q_{t+1}, \dot{q}_{t+1}, \ddot{q}_{t+1} = f_{\text{smooth}}(a_t, a_{<t}, q_t, \dot{q}_t, \ddot{q}_t)
$$

To maintain a completely hardware-agnostic implementation, the motion smoother queries the robot's specific physical limits directly from the SDK. It then enforces these constraints continuously, bounding the output trajectory by maximum velocity, acceleration, and jerk:

$$
\begin{align*}
|\dot{q}(t)| &\leq v_{\max} \\
|\ddot{q}(t)| &\leq a_{\max} \\
|\dddot{q}(t)| &\leq j_{\max}
\end{align*}
$$

```mermaid
flowchart TD
    A["Human Speech"] --> B["Audio Sensor"]
    B --> C["Speech-to-Text"]
    C --> D["Task Spec"]
    D --> E["Tactile-Aware Policy"]
    
    E -->|Intended Action a_t| F["Robot SDK"]
    
    subgraph SDK["Robot SDK"]
        G["Action History a_<t"] --> H
        F --> H["Motion Smoother<br/>(Online Interpolation)"]
        H --> I["Jerk-Limited<br/>Trajectory Generator"]
        J["Hardware Limits<br/>v_max, a_max, j_max"] --> I
    end
    
    I --> K["Hardware Controller"]
    K --> L["Robot"]
```

---

## 12. Inference-Time Search

The digital twin and surrogate remain active during runtime for counterfactual safety checks before executing high-uncertainty trajectories.

```mermaid
flowchart TD
    A["Current State"] --> B["Candidate Trajectories"]
    B --> C["Fast Surrogate Rollouts"]
    C --> D["Twin High-Fidelity Contact Check"]
    D --> E["Select Optimal Action"]
```

The surrogate provides cheap evaluation of many candidate trajectories.  
Promising or high-risk trajectories can then be evaluated using the higher-fidelity calibrated twin before execution on the physical robot.

---

## 13. End-to-End System

```mermaid
flowchart TD
    A["1. MULTIMODAL PROBING<br/><br/>Handheld Gripper Tactile Video + Contact Data<br/>+ Task Video<br/>+ Robot Hardware Specs"]

    A --> B["2. CONSTRUCT & CALIBRATE<br/><br/>Real-to-Sim Agent<br/>↓<br/>Multimodal SysID (Mass, Friction, Stiffness)<br/>↓<br/>Calibrated Digital Twin<br/>↓<br/>Learned Surrogate Simulator"]

    B --> C["3. LEARN<br/><br/>Massive RL on Surrogate<br/>↓<br/>High-Fidelity Twin Validation<br/>↓<br/>Real-World RL Fine-Tuning (Every 3 Months)"]

    C --> D["4. DEPLOY & ADAPT<br/><br/>Audio Task Input → Planner LLM Guardrails → Policy → Robot SDK → Motion Smoother → Hardware Controller → Robot<br/>↓<br/>Inference-Time Counterfactual Search<br/>↓<br/>On-the-Fly Human Feedback Correction Loop (VLM Judged)"]
```

The complete system therefore forms a continuous loop:

$$
\text{Probe} \rightarrow \text{Construct} \rightarrow \text{Calibrate} \rightarrow \text{Learn} \rightarrow \text{Deploy} \rightarrow \text{Observe} \rightarrow \text{Correct} \rightarrow \text{Learn Again}
$$

> **Note**: Every deployed policy automatically benefits from motion smoothing via the Robot SDK, ensuring fluid and hardware-safe operations out-of-the-box. Additionally, homes are supplied with non-breaking cups and dishes for the first month of usage.

---

## The Core Research Thesis

The central hypothesis is:

> Providing a Real-to-Sim agent with synchronized video and contact data from human environment probing allows it to build fully calibrated, contact-accurate digital twins. These twins generate fast surrogate models that make large-scale RL, real-world adaptation, and live human-feedback self-correction practical.

Instead of relying on pure video estimation or uncalibrated physics:

$$
\text{Interactive Tactile/Video Data} \rightarrow \text{Agentic Multimodal SysID} \rightarrow \text{Calibrated Twin} \rightarrow 	ext{Fast Surrogate RL} \rightarrow \text{Real Adaptation} \rightarrow \text{VLM-Guided Self-Correction}
$$

### The Product Vision

Delivers zero-to-hero autonomy:

> Provide the tactile probing data, task video, and robot SDK. The AI engine constructs the calibrated physics twin, trains the surrogate and policy, fine-tunes on hardware, and delivers a self-correcting autonomous robot.
