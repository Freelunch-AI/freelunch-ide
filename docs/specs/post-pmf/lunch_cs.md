# Lunch CS — Playground for Humans & AI Agents to Learn Performance Optimization

## The Vision

**Lunch CS** is an interactive, visual playground for learning how modern computers work—and how to make them faster.

It brings the computing stack into the browser through lightweight virtual machines and networks, while making the normally invisible behavior of computers observable and interactive:

**CPU/GPU → Memory → OS/Drivers → Compilers → Networking → Distributed Systems**

Users can run programs, inspect their execution, modify them, slow down and replay execution, and see how software decisions propagate through the underlying system.

> **Make computers a playground instead of a black box.**

---

## Learn by Building & Optimizing

Lunch CS comes with a game-like course in modern low-level computing.

Instead of learning only from lectures and diagrams, users build simplified versions of real systems and then optimize them.

Every lesson follows:

**Build a naive implementation → Run it → Profile it → Understand the bottleneck → Optimize it → Observe the result**

For example:

* Build a memory allocator → optimize allocation and cache behavior.
* Build a TCP server → optimize concurrency and I/O.
* Build a compiler/interpreter → optimize execution.
* Build a distributed key-value store → optimize communication and coordination.
* Build GPU programs → optimize memory access and parallelism.

The playground lets learners inspect things such as caches, memory, registers, scheduling, network traffic, compiler transformations, and GPU execution while their code is running.

Components can also be **expanded or collapsed** to move between abstraction levels—for example, from a simple "RAM" component down to the MMU, memory controller, and DRAM.

---

## A Training Ground for AI Agents

The same environment is designed for AI agents to learn performance engineering.

Today, an optimization agent often operates in a sparse loop:

**Change code → run benchmark → observe latency/throughput → try again**

Lunch CS provides a much richer environment where agents can observe the consequences of their changes throughout the system.

### Offline learning

Agents can train with reinforcement learning on thousands or millions of controlled optimization tasks without repeatedly consuming expensive real hardware.

They can first master individual optimization techniques:

* Cache locality
* SIMD/vectorization
* Memory tiling
* Kernel fusion
* Synchronization
* GPU optimization
* Scheduling
* Communication optimization

Once these skills are learned, agents can be trained on open-ended tasks where the objective is simply:

**Improve performance while preserving correctness.**

### Online optimization

The same environment can then be used by agents when optimizing real workloads.

An agent can:

**Inspect → Profile → Modify → Experiment → Validate → Deploy**

For example, it could optimize a CUDA kernel, Triton program, database, distributed service, or LLM serving workload.

Lunch CS gives the agent a place to experiment and understand the implications of an optimization before applying it to expensive production infrastructure.

---

## A Fully Observable Computer

The key idea is that performance is not represented by a single benchmark number.

Instead of simply seeing:

```text
Latency: 42 ms
GPU utilization: 67%
```

users and agents can see what produced those numbers:

* Cache hits and misses
* Memory traffic
* CPU/GPU execution
* Kernel launches
* GPU occupancy
* Scheduling
* Synchronization
* Network communication
* Storage access
* Compiler transformations

Execution can be paused, replayed, and inspected at any point.

This turns performance optimization from **guess-and-check** into an interactive process of understanding, experimenting, and improving.

---

## The Platform

Lunch CS combines ideas that are currently fragmented across different tools:

**Interactive simulator + IDE + debugger + profiler + systems course + RL environment**

It is inspired by projects such as **gem5, QEMU, OpenROAD, CompilerGym, and GIRUS**, but focuses on bringing the entire software/hardware stack together in an environment designed for both humans and AI agents.

The long-term goal is to create a common playground where:

**Humans learn how computers work.
Engineers learn how to optimize them.
AI agents learn performance engineering offline.
And those agents use the same environment to optimize real systems online.**

> **Lunch CS — Learn the machine. Experiment with it. Make it faster.**
