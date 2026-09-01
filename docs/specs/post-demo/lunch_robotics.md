# Lunch Robotics: Universal brain for any robot

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

```text id="j6xv5n"
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
                      ▼
              Digital Twin
                      │
                      ▼
             System Identification
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

The agent is a **purposefully RL-trained LLM agent** whose task is to convert real-world interaction data into a sensible executable simulation.

```text id="f5a5e0"
                    Video
                      │
                      ▼
             ┌─────────────────┐
             │ Real-to-Sim LLM │
             │      Agent      │
             └────────┬────────┘
                      │
             ┌────────┼─────────┐
             ▼        ▼         ▼
         geometry   objects   agents
             │        │         │
             ▼        ▼         ▼
          physics   materials  world model
             │        │         │
             └────────┼─────────┘
                      ▼
               Digital Twin
```

The agent is responsible for **constructing a plausible simulator**, not necessarily for estimating every parameter perfectly.

It can:

* inspect the video
* identify objects and agents
* determine relevant geometry
* select appropriate simulation primitives
* decide which components require learned models
* incorporate environment hints
* create simulation code/configuration
* run simulations
* inspect failures
* revise the digital twin
* determine what remains uncertain

The agent therefore acts as an **AI simulation engineer**.

---

# 2. System Identification

The initial digital twin produced by the agent will inevitably be imperfect.

Rather than requiring the LLM to estimate every physical parameter, the engine separates **semantic simulator construction** from **quantitative system identification**.

The agent might determine:

> "This object is a wooden table."

But the simulator still needs to determine parameters such as:

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

These parameters are estimated using optimization against observed real-world trajectories.

$$
\theta^*
=
\arg\min_\theta
D\left(
\tau_{\text{real}},
\tau_{\text{sim}}(\theta)
\right)
$$

The process therefore becomes:

```text id="f4m2ha"
        Real-world data
               │
               ▼
        LLM constructs
       plausible simulator
               │
               ▼
          Digital Twin
               │
               ▼
      optimization-based
     system identification
               │
               ▼
       calibrated twin
```

This division is deliberate:

> **The LLM determines the structure of the simulator; optimization determines its numerical parameters.**

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

The same general model can therefore simulate different types of autonomous entities.

Its predictions implicitly combine behavioral decisions with unobserved low-level dynamics.

The simulator becomes:

$$
\boxed{
\text{Digital Twin}
=
\text{Explicit Physics}
+
\text{General World Model}
}
$$

---

# 4. The Data Collection Protocol

The input video should not consist only of successful robot demonstrations.

The robot should be deliberately operated under the behavioral regimes it is expected to encounter.

For example:

```text id="8hwxql"
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

The task is therefore an explicit conditioning variable.

A single trained policy can potentially learn many tasks rather than requiring a separate policy for every task.

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

```text id="1ihh0p"
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
Complex real-world task
```

The curriculum can be automatically generated and adjusted based on the robot's performance.

---

# 7. Real-World RL Closes the Sim-to-Real Gap

Simulation is not expected to perfectly reproduce reality.

Before deployment, the simulation-trained policy therefore undergoes a final stage of **real-world reinforcement learning**.

This stage also uses curriculum learning.

The robot begins with simple, low-risk tasks and progressively moves toward more complex tasks.

```text id="6k6c8s"
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

```text id="e5ibk8"
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

This eliminates the need for continuous human supervision during the final adaptation stage.

---

# 8. Task Interface

The deployed robot should accept tasks through natural language.

```text id="7h4h8s"
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

This makes the task interface independent of the underlying control system.

The SDK provides the bridge between the learned policy and the robot's physical action interface.

---

# 9. Inference-Time Search

The digital twin is not discarded after training.

It remains available during deployment as a **counterfactual simulator**.

Before executing an uncertain action, the robot can simulate alternatives:

```text id="q2f1zq"
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

---

# 10. End-to-End System

The complete system can be viewed as four stages:

```text id="s2w7q0"
┌────────────────────────────────────────────────────────────┐
│                     1. OBSERVE                             │
│                                                            │
│ Robot interaction video + robot specification + metadata  │
└───────────────────────────┬────────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────────┐
│                     2. CONSTRUCT                           │
│                                                            │
│ RL-trained Real-to-Sim Agent → Digital Twin                │
│                                                            │
│ Optimization → System Identification                       │
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

> **A previosly trained AI agent, equipped with system identification tools, can convert a small amount of real-world robot video data into an executable digital twin that is accurate enough to support large-scale policy learning and counterfactual planning at inference-time.**

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
\text{Digital Twin}
\rightarrow
\text{Millions of simulated interactions}
\rightarrow
\text{Real-World RL}
\rightarrow
\text{Robot Policy}
}
$$

The ultimate product abstraction is therefore:

> **Give us the robot, interaction video, and action interface. We build its digital twin, train its policy, adapt it to reality, and return a robot that can learn to perform tasks specified in natural language.**

This positions the system not merely as a world model or simulator, but as an **end-to-end robot learning engine powered by automatically generated digital twins**.
