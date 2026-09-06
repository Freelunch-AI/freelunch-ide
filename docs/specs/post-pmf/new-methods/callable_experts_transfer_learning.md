# Callable Skill Experts for Foundation Models

Large foundation models are increasingly expected to solve complex robotics and scientific problems, yet many of the capabilities required for these tasks already exist in highly specialized models, algorithms, optimizers, and simulators. The bottleneck is that foundation models typically have to **relearn these capabilities from end-to-end data**, which is often scarce or expensive.

We propose **Callable Skill Experts (CSEs)**: a mechanism for incorporating arbitrary specialized computational skills directly into the architecture of a foundation model.

A skill expert is exposed through a simple interface:

```math
y_k = E_k(x_k)
```

The expert receives an input prepared by the foundation model and returns an output that can be integrated back into its hidden state. Internally, `E_k` can be anything: a neural network, optimization procedure, simulator, classical algorithm, or composition of multiple components.

For example, a Dex-Net expert could internally solve:

```math
g^* = \arg\max_g Q_{\mathrm{DexNet}}(o, g)
```

while exposing only:

```math
g^* = E_{\mathrm{DexNet}}(o)
```

to the foundation model. The optimization procedure is therefore completely encapsulated inside the expert.

## Skill Invocation

At selected layers, the foundation model learns whether to invoke a skill and how to interface with it. Given hidden state `h_l`, a learned input adapter produces the expert input:

```math
x_k = P_k(h_l)
```

A learned gate determines whether the expert should be invoked:

```math
g_k = \sigma(W_k h_l)
```

The expert output is projected back into the foundation model's representation space:

```math
h_l' =
h_l +
g_k A_k\left(E_k\left(P_k(h_l)\right)\right)
```

The same interface can therefore expose fundamentally different types of expertise. The expert may internally perform a direct prediction, an `argmax`, an optimization procedure, simulation, or an arbitrary program; to the foundation model, it is simply a callable function that maps an input to an output.

## Learning to Use Skills

Skill invocation becomes part of the model's policy rather than a manually designed pipeline. The expert-augmented model is trained with reinforcement learning to learn:

* when a skill should be invoked,
* which skill is useful for the current state,
* what information should be provided to the expert,
* and how the expert's output should influence subsequent decisions.

The resulting policy can be written as:

```math
\pi_{\theta}^{\mathrm{expert}}
=
\pi_{\theta}
\left(
a_t
\mid
h_t,
E_{1:K}
\right)
```

The experts therefore act as additional computational capabilities available to the policy during training.

## Expert-Free Distillation

The key step is to use the expert-augmented model to generate new training data.

After reinforcement learning, the model can use its specialized skills to produce high-quality trajectories, solutions, or demonstrations:

```math
\mathcal{D}_{\mathrm{expert}}
=
\left\{
\tau_i
\sim
\pi_{\theta}^{\mathrm{expert}}
\right\}
```

We then train a new model without access to the experts:

```math
\pi_{\phi}
\leftarrow
\mathrm{Train}
\left(
\mathcal{D}_{\mathrm{expert}}
\right)
```

The goal is therefore not to build a foundation model that permanently depends on external skills, but to **transfer specialized capabilities into the model itself**:

```math
\text{Specialized Expertise}
\rightarrow
\text{Expert-Augmented Model}
\rightarrow
\text{Synthetic Experience}
\rightarrow
\text{Expert-Free Model}
```

This provides a mechanism for **bootstrapping foundation models from existing domain knowledge** when large-scale end-to-end training data is unavailable.

## Why Robotics?

Robotics is a particularly natural domain because it already contains a large ecosystem of specialized capabilities that are difficult to acquire purely from demonstrations:

* grasping and manipulation models such as Dex-Net,
* inverse kinematics solvers,
* motion planners,
* collision checkers,
* trajectory optimizers,
* physics simulators,
* model-based controllers,
* specialized perception models.

Instead of forcing a foundation model to rediscover these capabilities from scratch, Callable Skill Experts allow the model to exploit them during training and subsequently **internalize their useful behavior through expert-free distillation**.

The broader hypothesis is:

> **Existing specialized computational knowledge can be converted into end-to-end foundation-model capabilities by temporarily exposing the model to callable skill experts and distilling the resulting behavior back into the model.**
