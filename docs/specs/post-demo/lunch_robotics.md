# Lunch Robotics: Universal Brain for Any Robot

## The Vision

We propose an end-to-end engine that turns a robot into a **task-capable autonomous system** using video of the robot interacting with its environment.

The user should not need to manually build a simulator, design an RL environment, engineer a training curriculum, or train the robot policy.

The required inputs are:

1. **Robot interaction video** collected in the target environment.
2. **Robot SDK** installed on the robot.
3. **Action interface** exposed through the SDK.
4. **Audio sensor** for receiving natural-language tasks.
5. Optional environment metadata and hints.

The engine automatically:

```text
                  USER INPUTS
                      │
          ┌───────────┼───────────┐
          ▼           ▼           ▼
      robot video   robot SDK   environment
                                 metadata
          │           │
          └───────────┼───────────┘
                      ▼
             Real-to-Sim Agent
                      │
                      │ builds + tests + identifies
                      ▼
              Calibrated Digital Twin
                      │
                      ▼
             Simulation RL
                      │
                      ▼
            Real-World RL
                      │
                      ▼
             Trained Robot Policy
                      │
                      ▼
                 DEPLOYMENT
```

At deployment time, the user simply gives the robot a task.

For example:

> "Clean the table."

The task is captured through the robot's audio sensor, converted to text using speech recognition, and passed to the robot policy through the SDK.

The goal is:

$$
\boxed{
\text{Robot}
+
\text{Video}
+
\text{SDK}
+
\text{Task}
\rightarrow
\text{Autonomous Robot}
}
$$

---

# 1. The Real-to-Sim Agent

The digital twin is constructed by an autonomous agent rather than a fixed reconstruction pipeline.

The agent is a **purposefully RL-trained LLM agent** whose task is to convert real-world interaction data into a sensible, executable, and calibrated simulation.

The agent is equipped with a toolbox for:

* video understanding
* 3D reconstruction
* simulator construction
* physics simulation
* trajectory extraction
* system identification
* parameter optimization
* simulation evaluation
* experiment design
* world-model integration
* code generation and execution

```text
                    Video
                      │
                      ▼
             ┌─────────────────┐
             │ Real-to-Sim LLM │
             │      Agent      │
             └────────┬────────┘
                      │
          ┌───────────┼────────────┐
          ▼           ▼            ▼
      geometry     objects       agents
          │           │            │
          ▼           ▼            ▼
       physics     materials   world model
          │           │            │
          └───────────┼────────────┘
                      ▼
               Initial Twin
                      │
                      ▼
             System ID Tools
                      │
                      ▼
              Calibrated Twin
                      │
                      ▼
                 validation
                      │
                 ┌────┴────┐
                 │         │
              failure    success
                 │         │
                 ▼         ▼
              revise    finalize
```

The agent is responsible for **constructing and validating the entire digital twin**.

It can:

* inspect the video
* identify objects and agents
* reconstruct relevant geometry
* select simulation primitives
* determine which components should use explicit physics
* determine which components require learned world-model dynamics
* incorporate environment hints
* create simulator code and configuration
* identify uncertain parameters
* invoke system-identification tools
* run optimization and simulation experiments
* compare simulated and real trajectories
* diagnose mismatches
* modify the simulator
* repeat the process until the twin reaches the required fidelity

The agent therefore acts as an **AI simulation engineer and system-identification engineer**.

---

# 2. Agentic System Identification

The initial digital twin produced from video will inevitably be imperfect.

The critical insight is that the Real-to-Sim Agent does not need to solve system identification directly.

Instead, it is equipped with **specialized system-identification tools** that allow it to measure and optimize the simulator against real-world observations.

The agent may initially infer:

> "This object is a wooden table."

and construct a corresponding simulation model.

But many numerical parameters remain unknown:

$$
\theta =
[
m,
\mu,
e,
k,
c,
\ldots
]
$$

where the parameters may include mass, friction, elasticity, stiffness, damping, and other dynamics.

The agent can invoke system-identification tools to estimate them.

For example:

```text
             Real-world trajectories
                       │
                       ▼
              Real-to-Sim Agent
                       │
                 "friction is
                  uncertain"
                       │
                       ▼
             System ID Tool
                       │
                       ▼
             Parameter search
                       │
                       ▼
              Simulator rollout
                       │
                       ▼
            Compare with reality
                       │
                       ▼
                  error
                       │
                       ▼
              Agent diagnoses
                    mismatch
                       │
                       ▼
             modify / re-run
```

A generic optimization objective is:

$$
\theta^*
=
\arg\min_\theta
D\left(
\tau_{\text{real}},
\tau_{\text{sim}}(\theta)
\right)
$$

But optimization is not the responsibility of the agent itself.

The agent **decides when and how to use the optimization tools**, what parameters to expose, which observations to compare, and whether the resulting simulator is sufficiently accurate.

This creates an agentic system-identification loop:

$$
\boxed{
\text{Observe}
\rightarrow
\text{Hypothesize}
\rightarrow
\text{Simulate}
\rightarrow
\text{Measure Error}
\rightarrow
\text{Optimize}
\rightarrow
\text{Revise}
}
$$

The division of labor is therefore:

> **The LLM reasons about the structure and debugging of the simulator; specialized system-identification tools perform the numerical estimation.**

This is analogous to giving a software engineer a compiler, profiler, debugger, and test suite rather than expecting the engineer to perform those operations manually.

---

# 3. Agentic Entities Are Simulated by a General World Model

Not every entity in the environment can be adequately represented by traditional physics.

Humans, cars, pedestrians, animals, and other robots make decisions.

The digital twin therefore uses a **general world model as the behavioral baseline for agentic entities**.

Rather than requiring access to their internal actions, the model predicts their future states directly from observations.

$$
S^{agent}_{t+1:t+H}
\sim
P_\phi
\left(
S^{agent}_{t+1:t+H}
\mid
S_{\leq t},E
\right)
$$

The same general world model can simulate different types of agentic entities.

Its predictions implicitly combine behavioral decisions with unobserved low-level dynamics.

The simulator therefore becomes:

$$
\boxed{
\text{Digital Twin}
=
\text{Explicit Physics}
+
\text{General World Model}
}
$$

The Real-to-Sim Agent determines how these components should be composed and can use observed interaction data to validate the resulting behavior.

---

# 4. The Data Collection Protocol

The input video should not consist only of successful robot demonstrations.

The robot should be deliberately operated under the behavioral regimes it is expected to encounter.

For example:

```text
              Robot interaction data
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
       passive       normal      aggressive
          │            │            │
          └────────────┼────────────┘
                       │
                    mistakes
                       │
                       ▼
              environment reactions
```

This is particularly important because other agents react to the robot.

A human may behave differently when the robot:

* waits
* approaches
* moves aggressively
* unexpectedly changes direction
* blocks a path
* makes a mistake

Therefore:

> **Robot mistakes are valuable interventions for learning the environment's response to the robot.**

The video dataset consequently provides evidence not only about the robot and environment, but also about how the surrounding world reacts to different robot behaviors.

---

# 5. Training the Robot Policy

Once the digital twin has been constructed and calibrated, the engine trains the robot's policy inside the simulator.

The policy learns:

$$
\pi_\theta(a_t\mid o_t,T)
$$

where:

* \(o_t\) is the robot's observation,
* \(T\) is the task,
* \(a_t\) is the robot action.

The task is an explicit conditioning variable.

A single trained policy can therefore learn to perform many tasks rather than requiring a separate policy for every task.

Simulation provides massive amounts of experience:

$$
S_0
\rightarrow
A_0
\rightarrow
S_1
\rightarrow
A_1
\rightarrow
\cdots
$$

allowing reinforcement learning to explore policies without consuming physical robot time.

---

# 6. Curriculum Learning

Training begins with simple tasks and progressively increases difficulty.

The curriculum can vary:

* task length
* number of objects
* environmental clutter
* precision requirements
* number of interactions
* disturbances
* uncertainty
* presence of other agents
* consequences of failure

```text
Simple task
     │
     ▼
More objects
     │
     ▼
More complex manipulation
     │
     ▼
Longer horizon
     │
     ▼
Dynamic environment
     │
     ▼
Other agents
     │
     ▼
Complex task
```

The curriculum can be automatically generated and adjusted based on the robot's performance.

---

# 7. Real-World RL Closes the Sim-to-Real Gap

Simulation is not expected to perfectly reproduce reality.

Before deployment, the simulation-trained policy therefore undergoes a final stage of **real-world reinforcement learning**.

This stage also uses curriculum learning.

The robot begins with simple, low-risk tasks and progressively moves toward more complex tasks.

```text
              Simulation RL
                    │
                    ▼
              Initial policy
                    │
                    ▼
             Real-world RL
                    │
          ┌─────────┴─────────┐
          ▼                   ▼
      simple tasks       complex tasks
          │                   │
          └─────────┬─────────┘
                    ▼
                deployment
```

The purpose is not to relearn the entire policy.

It is to adapt the simulation-trained policy to:

* residual physics errors
* actuator imperfections
* perception errors
* unmodeled dynamics
* real-world agent behavior
* discrepancies between simulation and reality

A VLM provides automated task generation and reward evaluation.

```text
                    VLM
               ┌──────┴──────┐
               ▼             ▼
        task generation   reward evaluation
               │             ▲
               ▼             │
            Robot ───────────┘
               │
               ▼
             RL update
               │
               ▼
        harder curriculum
```

The VLM can generate increasingly difficult tasks and evaluate the robot's behavior without requiring a human to manually label trajectories.

This creates an automated loop:

$$
\text{Task Generation}
\rightarrow
\text{Real Interaction}
\rightarrow
\text{VLM Evaluation}
\rightarrow
\text{Reward}
\rightarrow
\text{RL}
\rightarrow
\text{Harder Task}
$$

The final real-world RL stage therefore serves as a **closed-loop adaptation phase between the calibrated digital twin and reality**.

---

# 8. Task Interface

The deployed robot accepts tasks through natural language.

```text
Human speech
     │
     ▼
 Audio sensor
     │
     ▼
 Speech-to-text
     │
     ▼
     Task
     │
     ▼
 Robot SDK
     │
     ▼
πθ(a | observation, task)
     │
     ▼
   Robot
```

For example:

> "Put all the cups in the dishwasher."

becomes a structured task representation consumed by the policy.

The SDK provides the bridge between the learned policy and the robot's physical action interface.

This means the user does not need to understand the underlying policy, simulator, or RL infrastructure.

They provide the robot with a task.

---

# 9. Inference-Time Search

The digital twin is not discarded after training.

It remains available during deployment as a **counterfactual simulator**.

Before executing an uncertain action, the robot can simulate alternatives:

```text
                 Current state
                      │
          ┌───────────┼───────────┐
          ▼           ▼           ▼
        action 1    action 2    action 3
          │           │           │
       simulate    simulate    simulate
          │           │           │
          ▼           ▼           ▼
       future 1    future 2    future 3
          │           │           │
          └───────────┼───────────┘
                      ▼
                select action
```

The robot can therefore use the same digital twin for both:

**learning what to do** and **reasoning about what to do next**.

This creates a closed relationship between training and inference:

$$
\text{Digital Twin}
\rightarrow
\begin{cases}
\text{Policy Training}\\
\text{Counterfactual Search}
\end{cases}
$$

---

# 10. End-to-End System

The complete system can be viewed as four stages:

```text
┌────────────────────────────────────────────────────────────┐
│                     1. OBSERVE                             │
│                                                            │
│ Robot interaction video + robot specification + metadata  │
└───────────────────────────┬────────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────────┐
│                     2. CONSTRUCT                           │
│                                                            │
│ RL-trained Real-to-Sim Agent                               │
│                                                            │
│   Video understanding                                      │
│        ↓                                                   │
│   Simulator construction                                   │
│        ↓                                                   │
│   System identification tools                              │
│        ↓                                                   │
│   Simulation / measurement / optimization                  │
│        ↓                                                   │
│   Iterative validation                                     │
│        ↓                                                   │
│   Calibrated Digital Twin                                  │
└───────────────────────────┬────────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────────┐
│                     3. LEARN                               │
│                                                            │
│ Curriculum RL in Digital Twin                              │
│                         ↓                                  │
│ Curriculum RL in Real Environment                          │
└───────────────────────────┬────────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────────┐
│                     4. DEPLOY                              │
│                                                            │
│ Audio → Speech-to-Text → Task → Robot Policy → Action      │
│                                                            │
│ Digital Twin remains available for inference-time search   │
└────────────────────────────────────────────────────────────┘
```

## The Core Research Thesis

The central hypothesis is:

> **A purposefully RL-trained AI agent, equipped with specialized simulation and system-identification tools, can convert a relatively small amount of real-world robot interaction data into a calibrated executable digital twin that is accurate enough to support large-scale robot policy learning and counterfactual planning.**

This changes the economics of robot learning.

Instead of:

$$
\text{Millions of real-world interactions}
\rightarrow
\text{robot policy}
$$

the goal becomes:

$$
\boxed{
\text{Small real-world dataset}
\rightarrow
\text{AI-generated Digital Twin}
\rightarrow
\text{System Identification}
\rightarrow
\text{Millions of simulated interactions}
\rightarrow
\text{Real-World RL}
\rightarrow
\text{Robot Policy}
}
$$

The important architectural distinction is:

> **The Real-to-Sim Agent is not itself the system-identification algorithm. It is an intelligent operator of system-identification tools.**

It decides what needs to be identified, chooses appropriate experiments and tools, interprets the resulting errors, and iteratively improves the simulator.

The ultimate product abstraction is therefore:

> **Give us the robot, interaction video, and action interface. An AI agent builds and calibrates its digital twin, trains its policy, adapts it to reality, and returns a robot capable of learning tasks specified in natural language.**

This positions Lunch Robotics not merely as a world model or simulator, but as an **end-to-end robot learning engine powered by an autonomous Real-to-Sim agent**.
