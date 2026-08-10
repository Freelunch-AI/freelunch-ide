# FreeLunch — Demo Features: Implementation Order & Roadmap

## Mock

[Freelunch Mock](https://github.com/Freelunch-AI/freelunch-ide/blob/main/docs/mock.html)

## OSS Design References

We will probably borrow patterns from:

- **Kubero** — a Kubernetes-based, developer-friendly platform
- **OKD** — open source edition of Red Hat OpenShift, a k8s-based complete platform focused on enterprises
- **Tilt** — a strong dev/experimentation experience for Kubernetes
- **Backstage** — a plugin-based internal developer platform interface
- **Ray** — modern distributed programming framework for Python (*inspiration for the lunch-lang distributed programming framework idea, to be used within freelunch-ide, though ray works as runtime and lunch-lang would be at compile time*)

The goal is to take ideas from these tools to simplify the developer experience without copying their implementation wholesale.

---

## Glossary

| Term | Definition |
|---|---|
| **Monorepo** | The single Git repository for a customer organisation. Contains all products, services, workflows, platform config, generated L2 artifacts, and CI/CD configuration. One monorepo per customer. |
| **Workload** | Any single deployable unit inside the monorepo. In the Demo this is a Service; the Workflow directory remains a placeholder for post-Demo DAG workflows. |
| **Service** | A long-running, always-on, request-driven Workload (HTTP, gRPC, event consumer, etc.). |
| **Workflow (DAG)** | A trigger-driven Workload that runs to completion and has per-run execution state. App DAG workflows such as Airflow are post-Demo. |
| **Layer 1 (L1)** | FreeLunch's abstraction API. Developers model Workloads on the canvas, which maintains deterministic CUE files as the Git-backed authoring representation. |
| **Layer 2 (L2)** | K8s/Helm artifacts generated from L1. L2 is the deployment source of truth. Customers can view and edit L2 directly; conflicts with L1 are surfaced during compilation and must be resolved explicitly. |
| **Canvas** | The primary visual authoring surface for Services. Canvas blocks maintain deterministic L1 CUE files in the monorepo. |
| **Persona** | One of the three core personas—Platform Admin, Platform Engineer, and Developer. Temporary approval grants may be assigned to personas for staging/production review. |
| **Hotfix** | A permission grantable by Platform Admin to any Persona. Allows an explicitly audited merge to `main` without normal CI gates. |
| **Ephemeral Staging Environment** | A short-lived isolated environment spun up per PR, torn down after merge/close. Its desired state is managed through the same GitOps path as other environments. |
| **Blue-Green Deployment** | Argo Rollouts-managed promotion between active (blue) and preview (green) versions, with the previous revision retained for a configured rollback window. |
| **Agent Integration Layer** | A read-only REST API with an OpenAPI contract and a first-party skill through which coding agents query workload, pipeline, error, cost, and observability data. Repository changes remain the agent's normal write path. |
| **Local Dev/Experimentation Environment** | A fully local Kubernetes environment for validating workloads before shared CI/Staging. It uses ProxMox + Talos Linux and is inspired by the developer experience of Tilt. |
| **CI Executor** | Dagger is the reusable, container-native execution layer for CI steps. GitHub Actions is the remote CI orchestrator; Act emulates GitHub Actions locally. |

---

## Deferred Decisions and Scope

- **Kubernetes as the primary target** — the Demo is centered on vanilla Kubernetes rather than a managed cloud-specific stack.
- **Local environment** — the Demo uses ProxMox + Talos Linux to provide a local VM-based Kubernetes environment. This is local virtualization, not a simulation of a specific public cloud.
- **Deployment targets** — the Demo focuses on the local environment; deployment to existing EC2-based clusters (and later multi-cloud) is a later extension rather than a core Demo requirement.
- **Post-Demo scope** — real cloud deployment, support for app workflows (e.g. Airflow DAGs), stateful services, block store, internal package registry, and other items from the founding doc.
- **Node autoscaling** — cloud-specific node autoscaling such as Karpenter is post-Demo. The local Demo uses a fixed node pool.
- **New L1 abstraction types** — creating entirely new abstraction types from scratch is post-Demo. The Demo supports configuration of the existing abstractions.

---

# Implementation Strategy

The roadmap is intentionally organized around **vertical dependency slices**, rather than installing infrastructure components merely because they are available.

The critical path is:

```text
Foundation (Group 1)
  ↓
L1/L2 Model Core (Group 2) ───→ [IDE & CLI Shell Created]
  ↓
CI Execution (Group 3)
  ↓
GitOps Deployment Foundation (Group 4)
  ↓
★ FIRST USABLE SLICE (Sample Application Created)
  ↓
Progressive Delivery + Staging (Group 5)
  ↓
Observability + Cost (Group 6)
  ↓
Agent Integration Layer (Group 7) + Advanced Platform Features (Group 8)
