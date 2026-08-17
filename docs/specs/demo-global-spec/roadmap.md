# FreeLunch — Demo Features: Implementation Order & Roadmap

## Mock

[Freelunch Mock](https://github.com/Freelunch-AI/freelunch-ide/blob/main/docs/mock.html)

## OSS Design References

We will probably borrow some ideas/patterns from:

* **Kubero** — a Kubernetes-based, developer-friendly platform
* **OKD** — open source edition of Red Hat OpenShift, a Kubernetes-based complete platform focused on enterprises
* **Tilt** — a strong developer/experimentation experience for Kubernetes
* **Backstage** — a plugin-based internal developer platform interface
* **Ray** — a modern distributed programming framework for Python; inspiration for the `lunch-lang` distributed programming framework idea to be used within `freelunch-ide`, although Ray operates as a runtime while `lunch-lang` would operate at compile time                                                                                                                |
---

## Deferred Decisions and Scope

* **Kubernetes as the primary target** — the Demo is centered on vanilla Kubernetes rather than a managed cloud-specific stack.
* **Local environment** — the Demo uses ProxMox + Talos Linux to provide a local VM-based Kubernetes environment. This is local virtualization, not a simulation of a specific public cloud.
* **Deployment targets** — the Demo focuses on the local environment. Deployment to existing EC2-based clusters, and later multi-cloud support, are extensions rather than core Demo requirements.
* **Post-Demo scope** — real cloud deployment, support for app workflows such as Airflow DAGs, stateful services, block store, internal package registry, and other items from the founding document.
* **Node autoscaling** — cloud-specific node autoscaling such as Karpenter is post-Demo. The local Demo uses a fixed node pool.
* **New L1 abstraction types** — creating entirely new abstraction types from scratch is post-Demo. The Demo supports configuration of existing abstractions.

---

## Implementation Strategy

The roadmap is organized around **vertical dependency slices**, rather than installing infrastructure components merely because they are available.

The Sample Application is introduced early, at the model layer, because the correctness of the core FreeLunch abstraction can be verified without deploying anything. The same application then follows the roadmap as a continuous integration and runtime fixture.

## Critical Path

```text
Foundation (Group 1)
  ↓
★ L1/L2 Model Core (Group 2)
  │
  └── ★ Sample Application introduced as L1→L2 golden fixture
  ↓
IDE & CLI Shell
  ↓
CI Execution (Group 3)
  ↓
★ GitOps Deployment Foundation (Group 4)
  │
  └── ★ FIRST USABLE RUNTIME SLICE
       Sample Application deployed to Kubernetes
  ↓
Progressive Delivery + Staging (Group 5)
  ↓
Observability + Cost (Group 6)
  ↓
Agent Integration Layer (Group 7)
  +
Advanced Platform Features (Group 8)
```

## Two Important Milestones

### ★ Group 2 — First Model Usable

The Sample Application proves that a developer can declare a real application using the FreeLunch abstractions and that the system deterministically produces the expected deployment artifacts.

```text
Sample Application
  → FreeLunch L1 configuration
  → CUE evaluation / validation
  → L2 Kubernetes/Helm artifacts
  → artifact validation
```

This milestone answers:

> **Can FreeLunch express a real application and reliably turn that declaration into the deployment artifacts we expect?**

No Kubernetes deployment or ArgoCD is required yet.

### ★ Group 4 — First Usable Runtime Slice

The same Sample Application is taken through the complete deployment path:

```text
Sample Application
  → FreeLunch L1 configuration
  → CUE evaluation
  → L2 artifacts
  → Dagger CI
  → GitHub Actions
  → Git
  → ArgoCD
  → Kubernetes
  → basic workload health
```

This milestone answers:

> **Can FreeLunch take a declared application all the way from source and desired state to a running Kubernetes workload through GitOps?**

Everything after this builds on that working slice.

---

# Group 1 — Foundation

*No external FreeLunch dependencies. These establish the local runtime, repository structure, identity, and secrets primitives.*

## 1.1 Customer Monorepo Structure

Define the canonical directory topology for the customer monorepo.

> **Story:** As a Platform Admin, I run `freelunch init my-company` and a local Git repository is scaffolded with the canonical directory structure. GitHub remote creation/push is handled by the GitHub integration rather than being implicitly part of repository scaffolding.

```text
monorepo/
├── platform/            ← Platform configuration
├── products/
│   └── <product-name>/
│       ├── services/    ← FreeLunch L1 Service definitions
│       └── workflows/   ← Placeholder for post-Demo DAG Workloads
├── l2/                  ← Generated/customer-visible deployment artifacts
│   └── <environment>/   ← Environment-specific desired state
└── .github/
    └── workflows/       ← Thin GitHub Actions orchestration
```

* `products/<product-name>/workflows/` remains a placeholder because app DAG workflows are post-Demo.
* `.github/workflows/` contains thin orchestration that invokes Dagger rather than duplicating build/test logic.
* `l2/` is generated from L1 but remains customer-visible and editable by authorized Platform Engineers.
* The platform version is declared in `platform/freelunch.yaml`.
* The FreeLunch product repository's own TypeScript + Go tooling is an implementation concern and is not mixed into the customer monorepo definition.
* The customer monorepo is multi-language friendly and includes linting that validates user edits against FreeLunch-supported structures.

## 1.2 Local Dev/Experimentation Environment

The fully local runtime environment for the Demo.

> **Story:** As a Platform Admin, I run `freelunch start` on a supported machine and FreeLunch provisions the local Kubernetes environment and installs the platform components. The environment does not require a connection to a customer cloud cluster.

* **Kubernetes** — primary runtime target for the Demo.
* **ProxMox + Talos Linux** — local VM-based Kubernetes cluster.
* **Provisioning** — Terraform provisions the local infrastructure; Talos configures the Kubernetes nodes.
* Avoid a Kind-only implementation because the Demo is intended to exercise a real Kubernetes deployment path.
* `freelunch start` is the canonical command for creating/starting the local Demo environment.

## 1.3 Auth Service (Local Instance)

OIDC identity provider for the IDE and platform APIs.

> **Story:** As a Platform Admin, I open the local auth service admin console, create accounts for my team, and assign users to the correct Persona groups. Those users can log into the FreeLunch IDE and obtain credentials for the Agent Integration Layer.

* Self-hosted local OIDC provider.
* Used by the FreeLunch IDE and Agent Integration Layer.
* Three core Personas map to identity-provider groups:

  * Platform Admin
  * Platform Engineer
  * Developer
* Temporary approval grants are authorization concepts, not additional core Personas.
* Does not handle end-user auth for customer applications.

## 1.4 Secrets Store (Local Instance)

Application secrets store for Workload credentials.

> **Story:** As a Platform Admin, I store a third-party API credential in the secrets store. A pod running a Workload receives the value through a Kubernetes Secret without application code communicating directly with the secrets store.

* Self-hosted local secrets store.
* External Secrets Operator syncs secrets into Kubernetes Secrets.
* Pods consume Kubernetes Secrets as environment variables or mounted volumes.
* The secrets store is not an application SDK dependency.

---

# ★ Group 2 — L1/L2 Model Core

*Depends on: Group 1.*

This group establishes the FreeLunch configuration model and its deterministic transformation into deployment artifacts.

It must be independent of the deployment controller. The compiler can produce and validate L2 artifacts before ArgoCD exists.

## 2.1 L1 Workload Model

The high-level FreeLunch model used to represent Workloads.

> **Story:** As a Developer, I add a Service to the FreeLunch Canvas and configure its supported properties. FreeLunch maintains the corresponding L1 configuration deterministically in the monorepo.

* The Canvas is a block-based authoring surface.
* Blocks may represent a single pod or a supported multi-pod composition.
* Virtual blocks represent externally managed dependencies without hosting them.
* One deterministic L1 definition exists per Workload.
* L1 contains the high-level concepts developers interact with rather than exposing Kubernetes resources directly.
* Platform constraints and defaults are applied through the underlying CUE model.
* Creating entirely new L1 abstraction types from scratch is post-Demo.

## 2.2 CUE Evaluation and L1 → L2 Compilation

The deterministic evaluation and compilation layer that transforms FreeLunch configuration into deployment artifacts.

> **Story:** As a Developer, I compile a Service definition. FreeLunch validates the configuration and evaluates the applicable CUE constraints and policies to produce the corresponding Kubernetes/Helm artifacts in `l2/`.

* CUE is the engine behind validation and evaluation.
* The compiler is a standalone library/CLI capability and does not depend on ArgoCD.
* Local compilation performs validation and can be triggered explicitly with `freelunch compile`.
* CI later invokes the same compilation capability automatically.
* Validates L1 configuration against the platform model and applicable policies.
* Produces deterministic Kubernetes/Helm artifacts.
* Generated artifacts may include Argo Rollouts resources, but compilation does not require the Argo Rollouts controller to be installed.

## 2.3 L2 Artifact Management

The customer-visible deployment layer.

> **Story:** As a Platform Engineer, I edit an L2 artifact directly to add a supported customization. A later L1 compilation detects the divergence and shows a conflict instead of silently deleting the manual change.

* L2 lives in the monorepo and is versioned through Git.
* Platform Engineers may edit L2 directly.
* L2 is the deployment source of truth.
* L1 is the higher-level authoring model; L2 is the rendered deployment state.
* L1/L2 conflicts are explicit and must be resolved through Git/FreeLunch workflows.
* `l2/<environment>/` provides a place for environment-specific desired state.
* Deployment does not happen from the compiler directly; deployment happens later through GitOps.

## 2.4 ★ Sample Application — L1/L2 Golden Fixture

Introduce the Sample Application as soon as the model exists.

> **Story:** As a Developer, I define the Sample Application using FreeLunch abstractions. FreeLunch produces deterministic L2 artifacts, and automated tests verify that those artifacts represent the intended application correctly.

The initial Sample Application should be deliberately small and deterministic.

It serves as a **golden fixture** for:

* L1 configuration correctness.
* CUE evaluation correctness.
* L1 validation.
* L2 artifact generation.
* L2 artifact structure and semantics.
* Regression testing of the FreeLunch model.

The initial fixture does **not** need to run in Kubernetes yet.

The Sample Application is expanded later as the platform gains additional capabilities.

---

# Group 3 — CI/CD Execution

*Depends on: Groups 1–2.*

This group establishes the complete CI execution model before GitOps deployment is considered complete.

## 3.1 Pre-commit Hooks

Local enforcement layer, re-enforced by remote CI.

> **Story:** As a Developer, I run `git commit` with a formatting violation. The pre-commit hook blocks the commit. If I bypass the hook, remote CI independently runs the same required checks.

* Static checks — linting, type checking, security smell detection.
* Format enforcement.
* L1 validation.
* Customer-repository structure validation.
* Pipeline configuration validation where applicable.
* Remote CI, not Act, is the authoritative server-side enforcement mechanism.

## 3.2 Dagger CI Execution

Reusable, container-native CI execution layer.

> **Story:** As a Developer, I push a Go Service with no Dockerfile. The Dagger pipeline builds the image, executes tests and scans, and produces the CI outputs required by the GitHub Actions workflow.

* Dagger owns reusable build/test/scan execution logic.
* Buildpacks provide source-to-image builds without requiring a Dockerfile.
* The same Dagger pipeline can run locally and remotely.
* Dagger is not the GitHub Actions orchestrator.

## 3.3 GitHub Actions + Act Pipeline

Remote GitHub Actions pipeline with Act-based local emulation.

> **Story:** As a Developer, I open a PR. GitHub Actions runs the same Dagger-based pipeline that can be emulated locally with Act. The remote workflow is the authoritative CI result.

```text
Local
  └─ pre-commit
  └─ Act
      └─ Dagger
          ├─ compile / validate L1 → L2
          ├─ build images
          ├─ unit tests
          ├─ contract tests
          ├─ functional tests
          ├─ load tests
          └─ security scans

Remote PR
  └─ GitHub Actions
      └─ Dagger
          ├─ compile / validate L1 → L2
          ├─ build images
          ├─ unit tests
          ├─ contract tests
          ├─ functional tests
          ├─ load tests
          └─ security scans
```

* **Act is only a local emulator.** It is not the server-side CI system.
* **GitHub Actions is the remote CI orchestrator.**
* **Dagger is the reusable execution layer used by both.**
* CI validates the Sample Application as part of the model/compiler regression suite.
* CI publishes/commits generated L2 artifacts through the defined Git flow.
* Remote CI is authoritative even when the same workflow was tested locally with Act.

---

# ★ Group 4 — GitOps Deployment Foundation

*Depends on: Groups 1–3.*

**ArgoCD belongs here, not in the initial infrastructure group.** It becomes meaningful only once there is a complete path producing versioned L2 artifacts through CI.

This group turns the Sample Application from a **model-level fixture** into the **first running FreeLunch workload**.

## 4.1 ArgoCD

GitOps sync engine.

> **Story:** As a Platform Engineer, CI commits an approved L2 change to Git. ArgoCD detects the desired-state change and reconciles it into the Kubernetes cluster. No normal deployment path uses `kubectl apply` directly.

* Installed in the Kubernetes cluster via Helm.
* Watches the appropriate environment-specific L2 path.
* Requires repository credentials/access to the customer monorepo.
* Is the deployment actuator for normal Workloads.
* CI produces and versions desired state; ArgoCD reconciles that state.
* ArgoCD must not be treated as a prerequisite for L1/L2 compilation.

## 4.2 Basic GitOps Deployment

The first complete deployment vertical slice.

> **Story:** A Sample Application Service is defined, CI validates and builds it, the resulting L2 is committed to the correct environment, and ArgoCD reconciles it into Kubernetes.

* Establish the minimum Git → ArgoCD → Kubernetes path before adding blue-green delivery.
* No direct cluster mutation from the FreeLunch IDE/CLI.
* Deployment status is derived from Git and cluster reconciliation state.
* Environment-specific desired state must be explicit.
* The Sample Application becomes the canonical first deployed workload.

### ★ First Usable Runtime Slice

At the completion of this item:

```text
Sample Application
  → L1
  → CUE evaluation
  → L2
  → Dagger
  → GitHub Actions
  → Git
  → ArgoCD
  → Kubernetes
  → basic workload health
```

This is the first complete FreeLunch runtime workflow.

## 4.3 GitHub + Permission Enforcement

Connect identity, repository permissions, CI gates, and Kubernetes authorization.

> **Story:** A Developer can modify L1 and open PRs but cannot approve their own protected deployment path. A Platform Engineer can edit L2. Platform Admin controls emergency permissions.

* Auth-service identity is the source of FreeLunch persona identity.
* GitHub branch protection/repository permissions enforce PR and merge rules.
* Kubernetes RBAC enforces cluster-side permissions.
* Temporary approval grants map to explicit GitHub/deployment gates.
* A Persona group in the OIDC provider alone does **not** enforce permissions in GitHub or Kubernetes.
* Hotfix bypasses must be explicit, auditable, and limited to authorized users.

---

# Group 5 — Progressive Delivery + Staging

*Depends on Group 4 and the deployed Sample Application.*

## 5.1 Argo Rollouts Blue-Green Integration

Controlled application promotion through standard Argo Rollouts resources.

> **Story:** A change reaches the production environment through the GitOps flow. ArgoCD syncs the generated Rollout. Argo Rollouts manages the preview and active versions and promotes the new revision only after the configured checks succeed.

* FreeLunch generates a standard `Rollout`, active Service, preview Service, and analysis resources in L2.
* ArgoCD syncs resources; Argo Rollouts manages ReplicaSets, selectors, pauses, and analysis state.
* Demo traffic switching is scoped to active/preview Services in the same namespace.
* Promotion requires readiness and configured pre-promotion checks.
* Rollback changes the Git desired state and lets ArgoCD reconcile it; the CLI must not directly mutate the cluster.
* Stateful migrations and database/queue synchronization are outside Demo scope.

## 5.2 Ephemeral Staging Environments

One isolated environment per PR.

> **Story:** A PR that passes CI gets an isolated staging environment. The environment is deployed through the same GitOps machinery, smoke/e2e tests run against it, and it is removed after the PR is merged or closed.

* Created only after CI succeeds.
* Desired state is represented in Git and reconciled by ArgoCD.
* Environment isolation uses namespaces or lightweight virtual clusters such as vcluster.
* Multiple PRs can be tested in parallel.
* Teardown removes the environment's GitOps resources as well as the runtime resources.

## 5.3 Demo Autoscaling

Basic workload autoscaling.

> **Story:** A Platform Engineer configures supported HPA settings for a Service. Kubernetes scales the Workload within the configured limits.

* HPA is in Demo scope.
* Scaling configuration is part of the existing L1/platform-policy model.
* Platform Engineers control scaling policy; Developers cannot override enforced policy.
* VPA and cloud-specific node autoscaling such as Karpenter are post-Demo unless the local environment gains a concrete supported implementation.

---

# Group 6 — Observability + Cost

*Depends on Groups 4–5 and a deployed Sample Application.*

## 6.1 SigNoz Setup

Observability backend for customer Workload telemetry.

> **Story:** After the Sample Application is deployed, the FreeLunch observability surface shows its metrics, logs, and traces.

* Installed in the Kubernetes cluster.
* Collects customer Workload telemetry.
* FreeLunch provides the OpenTelemetry collection/routing layer; it does not magically instrument arbitrary application code.
* Existing OTel-compatible telemetry continues to work.
* Internal FreeLunch components are outside the customer observability view in the Demo.
* CI/CD telemetry requires an explicit exporter/integration; it should not be described as appearing in SigNoz automatically.

## 6.2 Customer Workload Observability

End-to-end operational visibility.

* Kubernetes workload/resource health.
* Application logs, metrics, and traces.
* CI/CD status from GitHub Actions and Dagger.
* Deployment/rollout state from ArgoCD and Argo Rollouts.
* All observability data is scoped to customer Workloads.
* The IDE aggregates these sources rather than implying Act is a remote telemetry source.

## 6.3 Cost Observability

Per-Workload cost breakdown.

> **Story:** As a Platform Admin, I open the cost panel and see estimated cost allocation for each product and Workload.

* OpenCost is installed with the Prometheus backend required by the Demo implementation.
* FreeLunch-generated resources carry stable product/Workload labels.
* The local Demo uses explicit custom prices and label values as estimates.
* Idle and unallocated costs remain visible.
* Budget enforcement, alerts, and deployment gates are out of scope.
* The long-term unified SigNoz-backed cost experience is post-Demo.

---

# Group 7 — Agent Integration Layer

*Depends on Groups 1, 4, 5, 6, and the repository/GitHub integration.*

## 7.1 Coding Agent Observability API

Read-only API for coding agents to query operational state outside the codebase.

> **Story:** As a coding agent, I obtain an auth-service client credential and query Workload, pipeline, deployment, error, cost, and observability state. The API cannot mutate the platform.

* Auth uses machine-to-machine OAuth2 client credentials.
* Publishes an OpenAPI contract.
* Exposes:

  * Workload status
  * Pipeline state and test results
  * Deployment/rollout state
  * Recent errors
  * Cost data
  * Infrastructure health
  * Application telemetry
  * Platform/customer documentation metadata
* Scope the API to structured operational state rather than claiming that literally every IDE feature is exposed.
* No write endpoints.
* Repository modifications happen through the agent's normal local Git workflow and GitHub access.

## 7.2 First-Party FreeLunch Skill

Enable coding agents to use FreeLunch correctly.

> **Story:** As a Developer, I install the FreeLunch skill in my coding agent. The agent can modify repository files and Git branches normally, while using the FreeLunch API to inspect deployment and operational state.

* Packages FreeLunch documentation, repository conventions, and operational runbooks.
* Uses the OpenAPI-described API.
* Covers Workload modeling, source-to-image builds, L1 validation/compilation, CI, deployment, rollout inspection, and diagnosis.
* Credentials are supplied by the agent runtime and never embedded in the skill.
* GitHub access is an explicit dependency of the agent workflow.
* MCP is not used in the Demo.
* Platform mutations continue to follow documented Git/GitOps paths.

---

# Group 8 — Advanced Platform Features

*Depends on Group 2 and the IDE Policy Editor.*

## 8.1 Platform Defaults and Policy Editor

Graphical editor for platform defaults, constraints, and enforced policies.

> **Story:** As a Platform Engineer, I change the default Service memory request. FreeLunch previews affected Workloads and the resulting L2 diff, then writes the platform configuration change to the Platform Engineer's branch.

* The underlying policy model already exists in Group 2; this group adds the graphical editor UI within the IDE shell.
* Distinguishes defaults, allowed constraints, and enforced policies.
* Supports settings already represented by the existing L1 abstractions.
* Shows affected Workloads and generated L2 changes before writing Git changes.
* Platform configuration remains inspectable and usable by CLI/CI.
* The underlying CUE model evaluates the platform configuration and produces the resulting L2 artifacts.
* Raw CUE authoring is not part of the normal user workflow.
* New L1 abstraction types remain post-Demo.
* Changes are tracked through Git and the approval flow.

## 8.2 Platform Versioning

Declarative platform version management.

> **Story:** As a Platform Engineer, I run `freelunch upgrade` after a new FreeLunch version is published. The CLI compares the target platform version with the current version and reports incompatible configuration.

* Versioned CLI releases and monorepo templates are published.
* `freelunch upgrade` compares target and current platform versions.
* Compatible schema/model upgrades can be applied after review.
* Breaking changes are reported as actionable conflicts.
* Platform versioning is deliberately after the initial vertical slice because it is not required to prove the core deployment workflow.

---

# Cross-Cutting Workstreams

Unlike the numbered groups, which represent structural dependencies and backend implementation order, these workstreams represent continuous user-facing and validation activities spanning multiple groups.

## FreeLunch IDE + CLI

* **Starts:** Immediately after Group 2.
* **Continues:** Through Groups 3–8.
* **Role:** Primary user interface. The IDE and Dev Portal are the same surface, with the CLI used for setup and inspection.
* The shell is established as soon as the L1/L2 model and compilation path exist.
* Subsequent platform capabilities are integrated into the shell incrementally as backend groups complete.

### IDE/CLI Integration by Group

| Group          | IDE / CLI Capabilities                                                                                              |
| -------------- | ------------------------------------------------------------------------------------------------------------------- |
| **Group 2**    | Visual Workload editor, Workload configuration views, L1/L2 preview and diff, `freelunch init`, `freelunch compile` |
| **Group 3**    | Pre-commit integration, local Act execution, CI status views                                                        |
| **Groups 4–5** | ArgoCD/Rollouts state panels, environment status, `freelunch status`, `freelunch rollback`                          |
| **Group 6**    | Workload health telemetry, SigNoz views, OpenCost panels, `freelunch logs`                                          |
| **Group 8**    | Policy Editor graphical interface, `freelunch upgrade`                                                              |

### Key Specifications

* Built on Eclipse Theia.
* Open VSX extensions available.
* Authenticated via the local OIDC auth service.
* Local Dev/Experimentation experience inspired by Tilt, including predefined scenarios such as load tests.
* Every platform mutation strictly follows the GitOps flow.
* CUE remains an implementation detail of the configuration/evaluation pipeline rather than the primary IDE abstraction.

---

## Sample Application

* **Starts:** Group 2.
* **First deployed:** Group 4.
* **Continues:** Through Groups 4–8.
* **Role:** Purpose-built application and golden fixture used to validate the FreeLunch model, CI/CD pipeline, deployment workflow, and subsequent platform capabilities.

### Group 2 — Model Fixture

The initial Sample Application is deliberately small.

It validates:

```text
L1
→ CUE evaluation
→ L2
→ deterministic artifact validation
```

### Group 4 Onward — Runtime Fixture

The same application is deployed through:

```text
L1
→ L2
→ CI
→ GitOps
→ Kubernetes
```

It then expands continuously as new platform capabilities are added.

### Composition

* One Service built from source using Cloud Native Buildpacks.
* One imported container-image Service.
* A virtual Service block for an externally managed database dependency.
* Health, echo, controlled-error, and latency test paths.

---

## Documentation

* **Starts:** Group 1.
* **Continues:** Through Groups 1–8.
* **Role:** Documentation is maintained continuously alongside feature implementation rather than deferred to the end.
* Built from Markdown via MkDocs.
* Covers:

  * Canvas/Workload modeling
  * L1 configuration
  * CUE evaluation and compilation
  * CI/CD
  * GitOps
  * Rollouts
  * Observability
  * Cost
  * Permissions
  * Agent integration
  * Operational runbooks
* Exportable as a unified `.md` file for LLM/offline consumption.

---

## Integration Testing & End-to-End Validation

* **Starts:** Group 2.
* **Continues:** Through Groups 3–8.
* **Role:** Validation starts with deterministic L1→L2 tests for the Sample Application and expands into full end-to-end deployment tests once Group 4 is available.

### Validation Progression

```text
Group 2
  → L1 validation
  → CUE evaluation
  → L2 artifact validation
  → Sample Application golden tests

Group 3
  → CI executes model/compiler tests
  → build/test/scan validation

Group 4
  → GitOps deployment
  → Kubernetes runtime validation

Groups 5–8
  → progressive delivery
  → staging
  → observability
  → cost
  → agent integration
  → platform management
```

Every newly completed feature in Groups 5–8 should be validated against the Sample Application using automated Dagger/GitHub Actions workflows.

---

# Low Priority / Post-Demo

*Deferred until the core vertical slice and Sample Application are stable.*

Items in this section should only be added after the end-to-end path is reliable:

* Stateful workload abstractions
* Block storage management
* Internal package registry
* VPA
* Karpenter / cloud-specific node autoscaling
* Real public cloud targets
* Support for app DAG workflows
* Completely new/custom L1 abstraction types
* Graphical policy editing for completely custom L1 types
* **PR criticality classification**
* **Selective test execution based on dependency analysis**

PR criticality and selective test execution are optimization and governance features rather than prerequisites for the core Demo workflow.

---

# Summary Table

| Group / Workstream | Focus                                                                                                        | Key Dependency            |
| ------------------ | ------------------------------------------------------------------------------------------------------------ | ------------------------- |
| **Group 1**        | Monorepo structure, local Kubernetes runtime, OIDC auth, secrets store                                       | None                      |
| **★ Group 2**      | L1 Workload model, CUE evaluation/compiler, L2 artifact model, Sample Application golden fixture             | Group 1                   |
| **★ Workstream**   | **FreeLunch IDE & CLI Shell**                                                                                | **Group 2 — L1/L2 Model** |
| **Group 3**        | Pre-commit, Dagger, GitHub Actions, Act                                                                      | Groups 1–2                |
| **★ Group 4**      | ArgoCD, basic GitOps deployment, GitHub/Kubernetes permission enforcement, first deployed Sample Application | Groups 1–3                |
| **★ Milestone**    | **First Usable Runtime Slice**                                                                               | **Group 4**               |
| **Group 5**        | Argo Rollouts blue-green, ephemeral staging, Demo HPA                                                        | Group 4                   |
| **Group 6**        | SigNoz telemetry, customer workload observability, OpenCost                                                  | Groups 4–5                |
| **Group 7**        | Coding Agent API, OpenAPI contract, first-party Agent Skill                                                  | Groups 4–6                |
| **Group 8**        | Policy Editor UI, declarative platform versioning                                                            | Group 2 + IDE Shell       |
| **Documentation**  | MkDocs documentation and operational runbooks                                                                | Continuous — Groups 1–8   |
| **Low Priority**   | Post-core features and extensions                                                                            | Core groups complete      |

---

# Recommended Implementation Milestones

## Milestone A — Foundation

```text
Monorepo
+ Local Kubernetes
+ OIDC
+ Secrets
```

---

## ★ Milestone B — Model & Developer Shell Core

```text
Canvas
  → FreeLunch L1 configuration
  → L1 validation
  → CUE evaluation
  → L2 generation
  → deterministic artifact validation
  → Sample Application golden fixture
  → FreeLunch IDE & CLI Shell initial launch
```

### Milestone B Exit Criteria

The Sample Application can be declared using FreeLunch abstractions and produces the expected L2 artifacts deterministically.

---

## Milestone C — CI Integration

```text
Act (local)
└── Dagger
     ├── compile / validate
     ├── build
     ├── test
     └── scan

GitHub Actions (remote)
└── same Dagger pipeline
```

The Sample Application's model and build/test validation run through the same CI execution path.

---

## ★ Milestone D — GitOps Vertical Slice

### First Usable Runtime Slice

```text
Sample Application
→ FreeLunch L1
→ CUE evaluation
→ generated/versioned L2
→ GitHub Actions
→ Git
→ ArgoCD
→ Kubernetes
→ basic workload health
```

> **First Usable Runtime Slice Reached:** The Sample Application is now a real running Workload. From this point onward, it becomes the primary end-to-end integration fixture for the rest of the Demo.

---

## Milestone E — Production-Like Delivery

```text
PR
→ CI
→ Ephemeral Staging
→ Smoke/E2E
→ Approval
→ main
→ ArgoCD
→ Argo Rollouts
→ Blue/Green
```

---

## Milestone F — Operational Experience

```text
Kubernetes + CI + Rollouts
→ SigNoz
→ OpenCost
→ IDE views
   ├── Observability
   └── Cost
→ CLI status & logs
```

---

## Milestone G — Agent Experience

```text
Operational state
→ Read-only API
→ OpenAPI
→ First-party Agent Skill
```

---

## Milestone H — Platform Management

```text
Existing platform configuration/policy model
→ IDE Policy Editor UI
→ Platform Versioning
→ freelunch upgrade
```

---

# Roadmap Dependency Overview

```text
                         ┌─────────────────────────┐
                         │       GROUP 1            │
                         │       Foundation         │
                         └────────────┬────────────┘
                                      │
                                      ▼
                    ★ ┌─────────────────────────┐
                      │       GROUP 2            │
                      │      L1 / L2 Core        │
                      │                          │
                      │  L1 model                │
                      │  CUE evaluation          │
                      │  L2 generation           │
                      │  Sample Application      │
                      │  golden fixture           │
                      └────────────┬────────────┘
                                   │
                      ┌────────────┴────────────┐
                      │                         │
                      ▼                         ▼
             ┌─────────────────┐      ┌────────────────────┐
             │  IDE + CLI      │      │   GROUP 3          │
             │  Shell          │      │   CI/CD Execution   │
             └─────────────────┘      └─────────┬──────────┘
                                                │
                                                ▼
                               ★ ┌────────────────────┐
                                 │     GROUP 4        │
                                 │ GitOps Foundation   │
                                 │                    │
                                 │ Sample Application  │
                                 │ becomes deployed   │
                                 └─────────┬──────────┘
                                           │
                                           ▼
                                  ★ FIRST USABLE
                                    RUNTIME SLICE
                                           │
                                           ▼
                                  ┌─────────────────┐
                                  │ Sample          │
                                  │ Application     │
                                  │ Runtime Fixture │
                                  └────────┬────────┘
                                           │
                 ┌─────────────────────────┼──────────────────────┐
                 │                         │                      │
                 ▼                         ▼                      ▼
          ┌─────────────┐           ┌─────────────┐       ┌─────────────┐
          │  GROUP 5    │           │  GROUP 6    │       │  GROUP 8    │
          │ Delivery +  │           │ Observability│       │  Platform   │
          │ Staging     │           │ + Cost      │       │ Management  │
          └──────┬──────┘           └──────┬──────┘       └─────────────┘
                 │                         │
                 └────────────┬────────────┘
                              │
                              ▼
                     ┌───────────────────┐
                     │     GROUP 7       │
                     │ Agent Integration │
                     └───────────────────┘

Documentation:        Groups 1 → 8 continuously
Integration Testing:  Groups 2 → 8 continuously
IDE + CLI:            Groups 2 → 8 continuously
Sample Application:   Group 2 → 8 continuously
```

## Core Principle

The roadmap deliberately establishes **two levels of proof** before expanding the platform:

1. **Model proof — Group 2:** FreeLunch can express a real application and deterministically produce the correct deployment artifacts.
2. **Runtime proof — Group 4:** Those artifacts can travel through CI and GitOps to produce a healthy running Kubernetes workload.

The Sample Application is the bridge between these two proofs and remains the continuous fixture for the rest of the Demo.
