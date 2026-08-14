# FreeLunch Terminology Standard

## Purpose

FreeLunch uses a deliberately constrained vocabulary so that developers, documentation, and AI agents interpret architectural terms consistently.

When a term has multiple meanings in software engineering, the definitions below establish its canonical FreeLunch meaning.

---

# 1. Core Terms

| Term                       | Canonical FreeLunch meaning                                                                                                                                                                                                                                          | Use when referring to                                                       |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------- |
| **Application**            | A customer software product composed of one or more Workloads and their associated source code and configuration.                                                                                                                                                    | The customer's software as a whole.                                         |
| **Workload**               | A single deployable unit. In the Demo, every Workload is a Service; application Workflows are not supported yet.                                                                                                                                                     | Services, and future Workflows.                                             |
| **Service**                | A long-running, request-driven Workload that normally remains running and handles requests or events.                                                                                                                                                                | Customer microservices or FreeLunch internal services.                      |
| **Pipeline**               | An automated CI/CD process that builds, tests, validates, and deploys software.                                                                                                                                                                                      | FreeLunch CI/CD.                                                            |
| **Environment**            | An isolated runtime context in which Workloads are developed, tested, or run.                                                                                                                                                                                        | Local Dev, Remote Dev, CI, Staging, and Production environments.            |
| **Platform**               | The FreeLunch-managed set of abstractions, policies, infrastructure, and tooling used to build and operate customer Workloads.                                                                                                                                       | FreeLunch as a developer/platform foundation.                               |
| **Infrastructure**         | The complete set of physical, virtual, and managed resources that provide the foundation for running the Platform and customer Workloads. It includes Compute Infrastructure and other infrastructure such as networking, storage, and externally managed resources. | Infrastructure as a whole.                                                  |
| **Compute Infrastructure** | The compute-oriented infrastructure provisioned and managed through infrastructure tooling. It consists primarily of one or more Compute Clusters and may also include supporting external compute-related resources.                                                | The infrastructure layer provisioned by IaC that provides compute capacity. |
| **Compute Cluster**        | A set of machines connected through a network and operated together as a compute resource pool. A Compute Cluster does not necessarily imply Kubernetes.                                                                                                             | Proxmox clusters and other machine clusters.                                |
| **Cluster**                | A generic grouping of machines or resources operated together. Prefer a more specific term such as Compute Cluster or Kubernetes Cluster when the type is known.                                                                                                     | Only when the cluster type is intentionally generic.                        |
| **Kubernetes Cluster**     | A Kubernetes control plane and its associated worker compute resources that provide a Kubernetes runtime environment.                                                                                                                                                | Kubernetes runtime environments.                                            |
| **Node**                   | A machine or VM participating in a Kubernetes cluster and capable of running Pods.                                                                                                                                                                                   | Kubernetes compute machines.                                                |
| **Replica**                | One running copy of a Workload.                                                                                                                                                                                                                                      | Individual running copies of a Service.                                     |
| **Deployment**             | The process of making a specific Workload version run in an Environment.                                                                                                                                                                                             | Delivering a Workload version to an Environment.                            |
| **Release**                | A specific version of a Workload or Application intended to be deployed or promoted.                                                                                                                                                                                 | Versioned software being released.                                          |

---

# 2. L1 / L2 and Artifacts

| Term               | Canonical FreeLunch meaning                                                                                         | Use when referring to                                                                                               |
| ------------------ | ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| **Configuration**  | Declarative data that controls how a Workload or FreeLunch component behaves.                                       | Configuration files, platform settings, and Workload settings.                                                      |
| **L1 Schema**      | The CUE-based schema defining the structure, constraints, and supported abstractions of FreeLunch L1 configuration. | The formal schema governing L1. Do not shorten this to "schema" when referring to the FreeLunch abstraction schema. |
| **Schema**         | A formal data-model definition describing the structure and constraints of data.                                    | Data models, database schemas, API/data schemas, and similar concepts.                                              |
| **Manifest**       | A declarative document describing a runtime object, especially a Kubernetes resource.                               | Kubernetes manifests.                                                                                               |
| **Artifact**       | A concrete generated or built output consumed by another part of the software lifecycle.                            | Helm charts, Kubernetes manifests, container images, and similar outputs.                                           |
| **Build**          | The process of transforming source code into a runnable artifact.                                                   | Source-to-image and executable builds.                                                                              |
| **Compile**        | The process of transforming L1 configuration into L2 artifacts.                                                     | The L1→L2 CUE-based compilation process.                                                                            |
| **L1**             | FreeLunch's high-level abstraction layer used by developers to model Workloads.                                     | The FreeLunch abstraction model.                                                                                    |
| **L2**             | The lower-level Kubernetes/Helm representation produced from L1 and used as the deployment source of truth.         | Generated Kubernetes and Helm artifacts.                                                                            |
| **Desired State**  | The state declared by Git/FreeLunch that the runtime should converge toward.                                        | GitOps and declarative deployment.                                                                                  |
| **Observed State** | The state currently reported by the runtime, including relevant externally managed resources where applicable.      | Cluster operational state and external operational state.                                                           |
| **State**          | The current condition of a Workload, Environment, Pipeline, deployment, or other managed object.                    | Runtime and operational condition.                                                                                  |

---

# 3. GitOps and Repository Terms

| Term              | Canonical FreeLunch meaning                                                                                             | Use when referring to                          |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| **Repository**    | A Git repository containing source code and/or FreeLunch configuration.                                                 | Git repositories generally.                    |
| **User Monorepo** | The single Git repository containing a customer's Products, Workloads, platform configuration, and CI/CD configuration. | The customer's canonical repository.           |
| **Product**       | A logical customer software product containing one or more Workloads, configurations, and related resources.            | The organizational grouping under `products/`. |
| **GitOps**        | The operating model in which Git contains declarative desired state and repository changes drive system changes.        | FreeLunch's deployment model.                  |

---

# 4. CI/CD Terms

| Term               | Canonical FreeLunch meaning                                                                      | Use when referring to                           |
| ------------------ | ------------------------------------------------------------------------------------------------ | ----------------------------------------------- |
| **Pipeline**       | The complete automated CI/CD process for a change.                                               | Build → test → validate → staging → deployment. |
| **Pipeline Stage** | A logical phase within a Pipeline.                                                               | Build, test, scan, staging, etc.                |
| **Job**            | A finite automated execution unit within a Pipeline.                                             | A concrete CI task/execution.                   |
| **Test Run**       | One execution of a test suite or test stage.                                                     | A particular test execution.                    |
| **CI/CD**          | Continuous integration and continuous delivery/deployment practices implemented by the Pipeline. | The overall software delivery process.          |

### Workflow

**Workflow is not a Demo FreeLunch concept.**

The Demo does not support application Workflows or DAG Workloads.

The word may appear in vendor-specific terminology such as:

```text
.github/workflows/
GitHub Actions workflow
```

These refer to GitHub terminology and are not FreeLunch application Workflows.

Post-Demo, **Workflow** may be introduced as a separate Workload type for finite, trigger-driven application computations such as Airflow-style DAGs.

---

# 5. Platform, Permissions, and Security Terms

| Term                | Canonical FreeLunch meaning                                                                      | Use when referring to                         |
| ------------------- | ------------------------------------------------------------------------------------------------ | --------------------------------------------- |
| **Persona**         | One of FreeLunch's named user categories with associated responsibilities and permissions.       | Platform Admin, Platform Engineer, Developer. |
| **Permission**      | Authorization to perform a specific action.                                                      | Editing L2, approving PRs, using Hotfix, etc. |
| **Role**            | A named collection of permissions/responsibilities when a role abstraction is actually required. | RBAC or permission bundles.                   |
| **Policy**          | A rule defining what is allowed, required, constrained, or enforced.                             | Platform policies and governance.             |
| **Platform Policy** | A FreeLunch policy controlling supported Workload configuration and behavior.                    | Defaults, constraints, and enforced policies. |
| **Hotfix**          | A permission allowing a user to merge directly to `main` without normal CI gates.                | The FreeLunch Hotfix capability.              |
| **Secret**          | Sensitive credential or value used by a Workload or platform component.                          | API keys, credentials, tokens, etc.           |
| **Identity**        | An entity recognized by the authentication/authorization system.                                 | Human and machine identities.                 |

---

# 6. Observability and Operations Terms

| Term                   | Canonical FreeLunch meaning                                                                         | Use when referring to                        |
| ---------------------- | --------------------------------------------------------------------------------------------------- | -------------------------------------------- |
| **Observability**      | The ability to understand runtime behavior through telemetry such as metrics, logs, and traces.     | SigNoz, OpenTelemetry, Workload health, etc. |
| **Telemetry**          | Runtime data used for observability.                                                                | Metrics, logs, and traces.                   |
| **Metric**             | A numeric measurement recorded over time.                                                           | CPU, memory, request rate, latency, etc.     |
| **Log**                | A timestamped record of an event produced by software.                                              | Application and infrastructure logs.         |
| **Trace**              | A representation of end-to-end execution of a request across components.                            | Distributed tracing.                         |
| **Health**             | The operational condition of a Workload or runtime object.                                          | Running, degraded, failing, etc.             |
| **Rollback**           | Restoring a previously deployed Workload version and reconciling the declarative state accordingly. | Argo Rollouts and GitOps rollback.           |
| **Cost Observability** | Visibility into estimated resource consumption and associated infrastructure cost.                  | OpenCost and FreeLunch cost views.           |

---

# 7. Developer Experience Terms

| Term                                      | Canonical FreeLunch meaning                                                                                                        | Use when referring to                        |
| ----------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------- |
| **IDE**                                   | The primary interactive development environment provided by FreeLunch.                                                             | The Eclipse Theia-based FreeLunch IDE.       |
| **CLI**                                   | FreeLunch's command-line interface for setup and inspection.                                                                       | `freelunch init`, `install`, `status`, etc.  |
| **Canvas**                                | The visual authoring surface used to model Workloads.                                                                              | Service Blocks and their relationships.      |
| **Block**                                 | A visual authoring unit on the Canvas representing a concrete Workload or dependency.                                              | Service Blocks and Virtual Blocks.           |
| **Virtual Block**                         | A Block representing an externally managed dependency that FreeLunch does not host.                                                | External databases, queues, etc.             |
| **Local Dev/Experimentation Environment** | The local runtime environment used to develop, test, observe, and experiment with Workloads before shared CI/Staging environments. | The local FreeLunch development environment. |
| **Ephemeral Staging Environment**         | A temporary isolated Environment created for a PR and destroyed after the PR is merged or closed.                                  | Per-PR staging.                              |

---

# 8. Agent and AI Terms

| Term                        | Canonical FreeLunch meaning                                                                                               | Use when referring to                                       |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------- |
| **Agent**                   | Software that autonomously reasons about a goal, observes results, and takes actions through available capabilities.      | Coding agents interacting with FreeLunch.                   |
| **Tool**                    | A callable capability an Agent can invoke to perform a specific operation.                                                | API calls, shell commands, repository operations, etc.      |
| **Skill**                   | Reusable procedural knowledge that teaches an Agent how to accomplish a class of tasks using available tools and context. | The first-party FreeLunch Skill.                            |
| **Agent Run**               | One concrete execution of an Agent against a task or goal.                                                                | A particular agent execution.                               |
| **Tool Call**               | One invocation of a Tool by an Agent.                                                                                     | Individual agent actions.                                   |
| **Context**                 | Information made available to an Agent during an execution.                                                               | Repository information, documentation, platform state, etc. |
| **Capability**              | An action or class of actions an actor is able and authorized to perform.                                                 | Agent and user capabilities.                                |
| **Agent Integration Layer** | FreeLunch's read-only programmatic interface exposing platform and operational state to coding agents.                    | The REST API and related agent integration.                 |

---

# 9. Words to Avoid Entirely

These words should **not** be used as FreeLunch architectural terminology.

They may still appear when quoting external documentation, referring to vendor terminology, or speaking in ordinary English.

| Word          | Reason                                                                                                                           |
| ------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| **System**    | Too broad. It can refer to nearly any level of the architecture. Name the actual entity instead.                                 |
| **Server**    | Ambiguous between a machine, process, backend application, or infrastructure object. Use the precise underlying concept instead. |
| **Thing**     | Completely non-specific.                                                                                                         |
| **Stuff**     | Non-specific and unsuitable for architectural terminology.                                                                       |
| **Component** | Too broad to function as a stable FreeLunch architectural primitive. Use the actual object being described.                      |

---

# 10. Words to Avoid in Specific Uses

These words are valid English and may be used, but **only with the meaning specified below**.

| Word                       | Do not use it for                                                                                 | Use it only for                                                                                                                               |
| -------------------------- | ------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| **Workflow**               | FreeLunch CI/CD, generic development processes, or Demo Workloads                                 | Vendor terminology such as GitHub Actions workflows, or the future post-Demo application Workload type                                        |
| **Schema**                 | The FreeLunch L1 abstraction schema                                                               | Data models, database schemas, API/data schemas, and similar data-structure definitions                                                       |
| **Cluster**                | A Kubernetes cluster when the intended meaning is actually a broader machine cluster              | A generic cluster when the specific cluster type is intentionally unspecified; otherwise prefer **Compute Cluster** or **Kubernetes Cluster** |
| **Infrastructure**         | Compute infrastructure when the intended scope is only machines and provisioned compute resources | The complete infrastructure layer, including compute, networking, storage, and other supporting infrastructure                                |
| **Compute Infrastructure** | All infrastructure indiscriminately                                                               | Infrastructure provisioned to provide compute capacity, including Compute Clusters and supporting external compute resources                  |
| **Service**                | A Kubernetes networking object                                                                    | A FreeLunch customer Workload; use **Kubernetes Service** for the Kubernetes object                                                           |
| **Deployment**             | A Kubernetes `Deployment` object unless explicitly qualified                                      | The act/process of deploying a Workload; use **Kubernetes Deployment** for the Kubernetes object                                              |
| **Resource**               | An unspecified piece of infrastructure                                                            | A specific resource class when qualified, such as **Kubernetes Resource** or **Compute Resource**                                             |
| **Instance**               | A generic synonym for any running software                                                        | A specific instantiated object when the underlying technology actually uses the term                                                          |
| **Module**                 | A generic synonym for an architectural component                                                  | A code/software module                                                                                                                        |
| **Process**                | A generic synonym for a workflow or pipeline                                                      | An executing OS/process-level program                                                                                                         |
| **Platform**               | Kubernetes, cloud infrastructure, or an arbitrary subsystem                                       | The FreeLunch Platform                                                                                                                        |
| **App**                    | A Workload or Service                                                                             | Only when speaking informally about an Application                                                                                            |
| **Object**                 | A generic synonym for any FreeLunch entity                                                        | An actual object in a technology that defines objects, such as a Kubernetes object                                                            |
| **Environment**            | A shell environment, virtualenv, or arbitrary configuration context                               | An isolated runtime context such as Local Dev, Remote Dev, CI, Staging, or Production                                                         |
| **Job**                    | A generic synonym for any task                                                                    | A finite automated execution unit in a Pipeline or a technology-specific Job                                                                  |
| **Run**                    | A generic synonym for state or deployment                                                         | One concrete execution of a Pipeline, Agent, test, or other executable unit                                                                   |

---

# 11. Naming Rules

### Prefer the most specific noun available

Do not write:

> "The system updates the resource."

Write:

> "ArgoCD syncs the L2 Helm artifact to the Kubernetes cluster."

### Qualify overloaded vendor concepts

Write:

* **FreeLunch Service**
* **Kubernetes Service**
* **CI/CD Pipeline**
* **GitHub Actions workflow**
* **Kubernetes Deployment**
* **Kubernetes Resource**
* **Compute Resource**
* **Compute Cluster**
* **Kubernetes Cluster**

when context does not make the distinction obvious.

### Do not overload a term across abstraction levels

The intended hierarchy is:

```text
Application
└── Workload
    └── Service
```

```text
Infrastructure
├── Compute Infrastructure
│   └── Compute Cluster
│       └── Machines / Nodes
├── Networking
├── Storage
└── External / Managed Resources
```

```text
Platform
├── IDE
├── CLI
├── Compiler
└── Platform Services
```

```text
Environment
└── Kubernetes Cluster
    └── Kubernetes Resources
```

```text
CI/CD
└── Pipeline
    ├── Stages
    └── Jobs
```

The general rule is:

> **One FreeLunch term should correspond to one conceptual level and one primary meaning.**
