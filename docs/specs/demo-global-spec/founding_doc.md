# FreeLunch — Company Intro + Product Requirements Document

# 1. Company Intro

## 1.1 What is FreeLunch?

**FreeLunch is a venture studio/accelerator that provides scaleup infrastructure-as-a-service to the companies it builds and supports.**

We partner with founders and pre-scale companies after they have reached, or are working toward, bootstrapped Product-Market Fit. Our goal is to let small teams scale software-intensive businesses without having to build large internal platform, DevOps, MLOps, and eventually other operational teams.

Our core product is an **AI-native developer platform for building and operating scalable applications**. Over time, FreeLunch may also provide ERP, technical upskilling, offices, and other shared infrastructure.

The platforms are designed for **portability**:

* Customers retain standard artifacts and GitOps workflows.
* Infrastructure remains inspectable and directly operable.
* Platforms can eventually be detached/ejected automatically.
* Portfolio companies can fork the platform internally, creating an inner-source ecosystem.

FreeLunch does **not** initially build its own coding agent or LLM. Instead, it turns existing coding agents such as Claude Code and OpenCode into complete software/AI engineering agents by giving them access to platform context, documentation, observability, and operational state.

## 1.2 Business Model

### Near Term

FreeLunch monetizes through **platform engineering consulting** and custom developer-platform implementations for scaleups.

### Long Term

FreeLunch becomes a **venture studio/accelerator**, using three main models:

1. **Idea → company:** FreeLunch validates an idea with a founder, provides salary and infrastructure, and takes significant founding equity.
2. **Pre-scale company:** Existing companies apply to join the accelerator in exchange for equity.
3. **Co-building:** FreeLunch and a founder explore and validate an opportunity together before forming the company.

FreeLunch also supports portfolio companies with fundraising through its VC network and participates in dilution alongside founders during future financing rounds.

## 1.3 Product Vision

FreeLunch is a **cloud-development IDE and internal developer platform** that combines application development, infrastructure, deployment, experimentation, and operations into one environment.

The core analogy is:

> **IDE + Internal Developer Platform as a single product**

Developers should be able to compose services and infrastructure visually, write application code, deploy through GitOps, observe the running system, and diagnose production issues without needing to become experts in the underlying Kubernetes/cloud-native ecosystem.

The platform provides a higher-level abstraction while preserving access to the underlying systems.

### Core principles

* **Developer-first:** infrastructure complexity should not dominate application development.
* **Git-native:** Git remains the source of truth for desired state and changes.
* **Kubernetes-native:** use Kubernetes and its ecosystem rather than replacing them.
* **Portable:** customers retain standard deployment artifacts and GitOps workflows.
* **Layered:** simple abstractions for developers, direct access to lower-level artifacts for platform engineers.
* **AI-native:** existing coding agents receive platform context and operational capabilities.
* **Composable:** integrate existing cloud-native tools rather than rebuilding them.

## 1.4 Competitive Position

FreeLunch combines ideas from or uses under the hood:

* **Kubero / Kubefirst** — developer-friendly Kubernetes and GitOps.
* **Backstage** — unified developer platform.
* **Tilt** — local development and experimentation.
* **Ray** — higher-level distributed programming abstractions.
* **Karmada** — multi-cluster abstraction.
* **Terraform / Crossplane** — infrastructure management.
* **SigNoz / OpenTelemetry** — observability.
* **OpenCost** — cost visibility.
* **Infisical / Vault** — secrets management.
* **SkyPilot** — ephemeral workloads.
* **MLflow / Langfuse** — AI/ML experimentation and observability.
* **n8n / ReactFlow** — visual composition.

FreeLunch's differentiation is combining these capabilities into a **single application-centric development and operations experience**, rather than another collection of independent platform tools.

---

# 2. Product Requirements Document

## 2.1 Product

**FreeLunch IDE / Developer Platform**

A Kubernetes-native, GitOps-based development environment for building, deploying, observing, and operating distributed applications.

The IDE and Dev Portal are the **same product surface**. The CLI provides setup and inspection capabilities.

## 2.2 Initial Product Goal

The first goal is to prove that FreeLunch can make scaling a stateless application significantly easier than using the underlying Kubernetes/cloud-native tooling directly.

The Demo must demonstrate:

1. Modeling and building multiple services.
2. Deploying them through GitOps.
3. Automatically scaling them under increased traffic.
4. Observing application behavior, deployments, and costs.
5. Diagnosing a failure using platform context and an AI coding agent.
6. Rolling back safely.

The Demo is a validation artifact, **not yet the MVP**.

---

# 3. Demo Product Experience

## 3.1 Build

A developer can create a stateless application consisting of:

* A Service built from source.
* A Service imported from an existing container image.
* Connections between Services.
* A virtual Service representing an externally managed dependency.

Source-based Services use **Cloud Native Buildpacks**, avoiding the need to manually create Dockerfiles.

The developer interacts primarily with FreeLunch abstractions rather than Kubernetes manifests.

## 3.2 Deploy

The normal deployment flow is:

```text
Developer change
    ↓
Git
    ↓
CI
    ↓
Generate deployment artifacts
    ↓
Git commit
    ↓
ArgoCD
    ↓
Kubernetes
```

FreeLunch does not normally mutate the cluster directly.

## 3.3 Scale

The developer can generate increasing traffic against the sample application.

Kubernetes automatically scales workloads according to their configured policies.

The Demo focuses on workload autoscaling rather than cloud-node autoscaling.

## 3.4 Observe

The developer can inspect:

* Application logs, metrics, and traces.
* Workload health.
* Deployment and rollout state.
* CI/CD status.
* Resource usage.
* Estimated workload costs.

Observability is provided through **OpenTelemetry + SigNoz**, while **OpenCost** provides cost allocation.

## 3.5 Diagnose

A coding agent can use the FreeLunch platform context to investigate failures.

The agent receives:

* Workload state.
* Deployment state.
* CI/test results.
* Recent errors.
* Observability data.
* Cost information.
* Platform documentation.

The agent can modify the repository using its normal Git/GitHub workflow, but the FreeLunch Agent API itself is **read-only**.

## 3.6 Roll Back

Deployments use **Argo Rollouts** for blue-green delivery.

The previous revision remains available for rollback.

Rollback follows the GitOps model rather than allowing the CLI or agent to directly mutate the cluster.

---

# 4. Architecture

## 4.1 Layered Application Model

FreeLunch provides two conceptual layers.

### Layer 1 — Developer Abstraction

The developer declares application intent using FreeLunch's platform model.

The model is represented using **CUE**.

CUE is used for:

* schema validation,
* constraints and defaults,
* deterministic evaluation,
* compilation of the higher-level application model.

CUE is **not the deployment artifact** and should not be treated as a second Kubernetes manifest format.

### Layer 2 — Deployment Artifacts

The platform model is compiled into customer-visible Kubernetes/Helm deployment artifacts.

```text
FreeLunch Application Model
        ↓
       CUE
        ↓
  L2 Kubernetes/Helm
        ↓
       Git
        ↓
     ArgoCD
        ↓
   Kubernetes
```

L2 remains inspectable and directly operable by platform engineers.

The generated artifacts are versioned in Git, allowing developers and platform engineers to see exactly what will be deployed.

## 4.2 GitOps

Git is the central mutation mechanism.

The platform should avoid direct imperative deployment operations wherever possible.

```text
IDE / CLI / Agent
       ↓
     Git
       ↓
      CI
       ↓
  Generated L2
       ↓
    ArgoCD
       ↓
 Kubernetes
```

## 4.3 Platform Logic

FreeLunch primarily acts as a **build-time/compiler layer**, not a Kubernetes controller.

Existing Kubernetes controllers should own domain-specific operational behavior:

* ArgoCD → GitOps reconciliation.
* Argo Rollouts → progressive delivery.
* HPA → workload scaling.
* External Secrets → secret synchronization.
* Other specialized controllers → their respective domains.

FreeLunch composes these capabilities into a coherent developer experience.

---

# 5. Demo Scope

## 5.1 Required

### Developer Experience

* FreeLunch IDE built on Eclipse Theia.
* VS Code-compatible extension ecosystem through Open VSX.
* Visual application composition.
* Minimal CLI.
* Git-native workflows.
* Platform documentation.
* Local development/experimentation environment.

### Application Platform

* Stateless Services.
* Source-to-image builds with Cloud Native Buildpacks.
* Container-image Services.
* Virtual Services.
* Kubernetes deployment.
* Helm/L2 artifacts.
* Workload autoscaling.

### CI/CD

* GitHub Actions.
* Act for local CI emulation.
* Dagger for reusable CI execution.
* Pre-commit hooks.
* Unit, integration, functional, and load testing.
* Security scanning.

### Deployment

* ArgoCD.
* Argo Rollouts.
* Blue-green deployments.
* Rollback.
* Ephemeral staging environments.
* GitHub and Kubernetes permissions.

### Operations

* OpenTelemetry.
* SigNoz.
* OpenCost.
* Workload health.
* Deployment state.
* CI/CD state.
* Logs, metrics, and traces.
* Cost visibility.

### AI Integration

* Read-only Coding Agent API.
* OpenAPI contract.
* First-party FreeLunch skill.
* Documentation-backed agent workflows.

---

# 6. Explicitly Out of Demo Scope

The Demo does **not** attempt to solve the entire platform vision.

### Application capabilities

* Stateful services.
* Database/queue lifecycle management.
* DAG workflows.
* GPU workloads.
* Frontend application management.
* Data engineering.
* Full MLOps.
* Confidential computing.
* Distributed programming framework.

### Infrastructure

* Public-cloud deployment.
* Multi-cloud.
* On-premise clusters.
* Embedded devices.
* Remote Kubernetes development environments.
* Full multi-cluster management.

### Platform

* Arbitrary new L1 abstraction types.
* Fully tool-agnostic IaC.
* Public/private reusable-block marketplace.
* Automated detach/eject.
* Agent-triggered platform mutations.
* Project management.
* Agent management.
* IDE action logging.
* DORA/IDE analytics.
* Budget enforcement and deployment cost gates.

### Advanced product capabilities

* A/B testing.
* Time-travel/slow-motion trace replay.
* System-wide experiment tracking.
* DataOps promotion workflows.
* Full OpenLineage integration.
* Monitoring actions directly from IDE widgets.

---

# 7. Demo vs MVP

The Demo proves the **core developer experience**.

The MVP requires the platform to become useful to a real scaleup.

Major additions include:

* Stateful workloads.
* Real cloud deployment.
* Proper authentication.
* Proper security.
* A/B testing.
* Validation with a real scaleup.

The most important validation is not feature completeness: it is demonstrating that a real scaleup can operate and grow applications with significantly less platform complexity and headcount.

---

# 8. Personas

### Platform Admin

Owns the platform, organization configuration, identity, and emergency permissions.

### Platform Engineer

Defines platform capabilities, policies, infrastructure, and deployment configuration.

### Developer

Builds, deploys, observes, and operates applications through the platform.

### Tech Lead

A developer with additional responsibility for reviewing and merging changes.

---

# 9. Key Technical Decisions

## Theia instead of a VS Code fork

FreeLunch is ultimately a distributed-systems workbench rather than a code editor.

Theia provides a replaceable IDE shell while allowing FreeLunch to make services, deployments, topology, observability, and runtime state first-class concepts.

Monaco, xterm, LSP, and DAP can provide code-centric experiences where needed.

## Go + TypeScript

Go is the primary backend/CLI language because of its cloud-native ecosystem, performance, and team familiarity.

TypeScript is used for the IDE frontend.

## CUE-based compilation

CUE provides the declarative modeling, validation, constraints, and deterministic evaluation required for the higher-level application model.

FreeLunch compiles that model into standard deployment artifacts rather than implementing a custom Kubernetes runtime controller.

## GitOps over direct deployment

Git provides auditability, review, reproducibility, rollback, and portability.

ArgoCD is responsible for reconciling approved desired state into Kubernetes.

## Kubernetes-native rather than Kubernetes replacement

FreeLunch should make Kubernetes easier to use rather than hide it permanently.

Platform engineers retain access to the underlying Kubernetes and generated artifacts.

## Local-first Demo

The Demo uses a fully local Kubernetes environment to reduce iteration time and infrastructure cost.

The target local environment is:

* Proxmox
* Talos Linux
* Kubernetes
* Terraform

The local environment is intended to exercise a realistic infrastructure path rather than merely emulate Kubernetes.

---

# 10. Demo Validation

The Demo should be compared against:

* Kubernetes + ecosystem tooling directly.
* Open-source developer platforms such as Kubero and Kubefirst.
* PaaS products such as Heroku, Railway, and Render.

Comparison criteria:

| Dimension            | Question                                                           |
| -------------------- | ------------------------------------------------------------------ |
| Developer Experience | How quickly can a developer build and operate the application?     |
| Capabilities         | Can the platform support realistic scaleup requirements?           |
| Complexity           | How much Kubernetes/cloud-native knowledge is required?            |
| Portability          | Can the customer operate the underlying artifacts independently?   |
| Lock-in              | How difficult is it to leave the platform?                         |
| Cost                 | What infrastructure and engineering cost does the approach create? |
| Headcount            | How much platform/DevOps expertise is required?                    |

The strongest validation is adoption by a real scaleup.

---

# 11. Development Environment

FreeLunch's own development environment should be reproducible and lightweight:

* Linux / WSL2.
* Pixi.
* Git + GitHub.
* Girus for Kubernetes experimentation where useful.

The development repository should provide:

* Reproducible development environment.
* Linting and formatting.
* Automated tests.
* Build tasks.
* Release/publishing workflows.
* IDE binary publishing.
* Go CLI/module publishing.
* Documentation website generation.

---

# 12. Initial Commercial Strategy

Before the venture-studio/accelerator model is fully operational, FreeLunch can generate revenue through **platform engineering consulting**.

The initial service is effectively an outsourced platform engineering team for scaleups that need:

* Internal developer platforms.
* Kubernetes/cloud infrastructure.
* CI/CD.
* Observability.
* Platform automation.

Custom implementations should be used selectively because company-specific requirements do not always generalize into the core FreeLunch product.

---

# 13. Product North Star

FreeLunch succeeds when a small engineering team can build, scale, and operate a sophisticated distributed application **without needing to become experts in every cloud-native tool and how they should work together.**

The platform should make low-level access available, but without making the complexity mandatory.
