# Lunch Robotics | Intelligent Robotics App: universal brain for any robot

The last piece of intelligent robotics. Works out-of-the-box with any robot in any environment (simulated or real). 

A software that you put on any robot and it makes the robot inteligent and self-improving. Support deployment on multiple embodiements, multiple embedded OSs and to multiple devices with suficient compute, memory and internet connection.

The software does online (powered by foundation models & logic engines, and with an initialization phase):

1. Vibe tuning Stuff
  - model vibe tuning
  - agent vibe tuning
2. simulator stuff
  - agentic generation
    - simulator assets generation
    - simulator (driver + dynamical system + observer) generation
  - learning
    - simulator (driver + dynamical system + observer) finetuning (makes use of SciML methods)
    - simulator (driver + dynamical system + observer) symbolic regression (makes use of SciML methods)
    - fast surrogate simulator learning (makes use of SciML methods)
3. (control theory) controller stuff
  - controller synthesis
  - controller finetuning
4. other agents stuff
  - modelling other agents (their dynamical systems + policies)
  - communication with other AI agentes and communication with humans via speech generation/processing
5. sensor fusion & state estimation (observation with requires and optional sensors, simulator -> state)
6. remote safe sandbox usage (internet browsing + file editing + terminal + storage APIs + human communication APIs) with attached gpu ray cluster (that scales to zero) + gpu genesis (physics engine) cluster + storage services.
7. curriculum synthesis (pipeline of learning setups in order to get to a model that do well at the goal task)
8. reward/feedback engineering (reward function can change over time, ideal should go from very dens to very sparse as the agente gets better)
9. hypothesis testing and causal experiments
10. policy learning/finetuning (via RL or imitation learning)
11. action smoothing (steps -> curves)
12. emergency detector and emergency policy learning/finetuning
13. planning/replanning & HIL plan reviewing
14. data collection (including active learning) & dataset building
15. logic engine run & autoformalization/deformalization

Where heavy (learning and synthesis) stuff can be delegated to happend async at the cloud and observation->action stuff + data collection stuff happens always (async or sync) at the robot.
	
Each robot needs to provide an MCP server declaring its hardware spec (actuators, sensors and structure/mobility),

Foundation Models (powering lunch-robo) iniitialized with prior knowledge and finetuned on a lot of different virtual {robot dynamics (real world: provided by MCP and inferred), environment dynamics (real world: inferred), context (real world: provided by the user), goal (real world: provided by user), constraints (real world: inferred and/or provided by user)} setups generated usng the Genesis Generative Physics Engine.

Business model: OSS, but offer it as a service (Robotic-Intelligence-as-a-service: RIaaS) where users just need to use our SDK in their robots.
