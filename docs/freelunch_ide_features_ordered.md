# FreeLunch — Demo/MVP Features: Implementation Order

## Mock

[Freelunch Mock](https://github.com/Freelunch-AI/freelunch-ide/blob/main/docs/freelunch_ide_mock.html)

## Glossary (Revised)

| Term | Definition |
|---|---|
| **Monorepo** | The single Git repository for a customer organisation. Contains all products, services, workflows, and platform config. One monorepo per customer. Replaces "Product Repo." |
| **Workload** | Any single deployable unit inside the monorepo. Either a Service or a Workflow. |
| **Service** | A long-running, always-on, request-driven Workload (HTTP, gRPC, event consumer, etc.). |
| **Workflow** | A trigger-driven Workload that runs to completion. Has per-run execution state. |
| **Layer 1 (L1)** | FreeLunch's abstraction API. Developers model Workloads on the canvas, which maintains CUE files as the Git-backed representation. Manually compiled into L2. |
| **Layer 2 (L2)** | K8s/Helm artifacts. Source of truth for deployment. Customers can view and edit directly. Raw edits are preserved on recompile via conflict detection. |
| **Canvas** | The primary visual authoring surface for Services and Workflows. Canvas blocks maintain deterministic L1 CUE files in the monorepo. |
| **Persona** | One of four pre-defined roles: Platform Admin, Platform Engineer, Developer, Tech Lead. |
| **Hotfix** | A permission grantable by Platform Admin to any Persona. Allows merging directly to main without CI gates. |
| **Ephemeral Staging Environment** | A short-lived isolated environment spun up per PR, torn down after merge/close. |
| **Blue-Green Deployment** | Argo Rollouts-managed promotion between active (blue) and preview (green) versions, with the previous revision retained for a configured rollback window. |
| **Agent Integration Layer** | A read-only REST API with an OpenAPI contract and a first-party skill through which coding agents query workload, pipeline, error, cost, and observability data. |
| **Coding/Experimenting Environment** | A girus-powered production simulation sandbox. Reads live SigNoz/OTel data to seed simulations. Used to anticipate production problems. |

---

## Deferred Decisions and Scope

- **Vanilla K8s vs LocalStack EKS** — both are supported; which the Demo primarily targets is TBD.
- **Remaining tool choices** (Kargo, Bazel, etc.) — Argo Rollouts and OpenCost are selected below; other conceptual roles are finalized during implementation of each feature.
- **Post-MVP scope** — stateful-service wiring and datastore or queue migration are excluded from the Demo/MVP and require separate engine-specific designs.

---

## Group 1 — Foundation
*No dependencies. Start here.*

### 1.1 Monorepo Structure
Define the canonical directory topology for the customer monorepo.

> **Story:** As a Platform Admin, I run `freelunch init my-company` and a new Git repository is created with the canonical directory structure (`platform/`, `products/<product-name>/services/`, `products/<product-name>/workflows/`, `.github/workflows/`) — ready to commit and push to GitHub.

```
monorepo/
├── platform/              ← L1 platform config (Platform Engineers)
├── products/
│   └── <product-name>/
│       ├── services/      ← Canvas-maintained L1 Service definitions
│       └── workflows/     ← Canvas-maintained L1 Workflow definitions
└── .github/
    └── workflows/         ← CI/CD pipeline (managed by FreeLunch)
```

- Separation between platform config (Platform Engineers) and canvas-maintained Workload definitions (Developers)

### 1.2 Local K8s Infrastructure
The fully-local runtime environment for the Demo.

> **Story:** As a Platform Admin, I run `freelunch install` on a fresh machine with Docker installed and a running Kind cluster appears locally with LocalStack emulating AWS services alongside it. All FreeLunch components deploy into it without touching the internet. I can also run `freelunch install --adopt` against an already-running cluster and FreeLunch installs without disrupting existing workloads.

- **Kind** — vanilla Kubernetes-in-Docker as primary local K8s target
- **LocalStack** — emulates the selected AWS APIs needed by the local Demo
- Two install modes supported:
  - **Fresh** — `freelunch install` spins up a new Kind cluster
  - **Adopt existing** — `freelunch install --adopt` installs FreeLunch components onto an already-running cluster without disrupting existing workloads
- EKS (via LocalStack) is also a supported target; final Demo target (Kind vs LocalStack EKS) is TBD

### 1.3 Keycloak (Local Instance)
OIDC identity provider for all human authentication within FreeLunch.

> **Story:** As a Platform Admin, I open the local Keycloak admin console, create accounts for my team, and assign each user to the correct Persona group. Those users can then log into the Theia IDE with their credentials on the next login — no separate account setup required in FreeLunch itself.

- Self-hosted local Keycloak instance
- Used by: Theia IDE (developer login), Agent Integration Layer (machine-to-machine client credentials)
- The 4 Personas map to Keycloak groups
- Does **not** handle end-user auth for customer applications (out of Demo scope)

### 1.4 Vault (Local Instance)
Application secrets store for all Workload credentials.

> **Story:** As a Platform Admin, I store a third-party API credential in Vault at path `secret/my-service/api-key`. A pod running `my-service` starts with an `API_KEY` environment variable already populated — no Vault SDK in the application code and no Kubernetes Secret created manually.

- Self-hosted local Vault instance (not an AWS service)
- Stores API keys and service credentials provisioned by the customer
- External-secrets-operator (Group 2) syncs secrets from Vault → K8s native Secrets
- Pods consume secrets as env vars and never communicate with Vault directly

---

## Group 2 — Core Infrastructure Components
*Depends on: Group 1 (K8s cluster, Keycloak, Vault)*

### 2.1 External-secrets-operator
Bridge between Vault and K8s native secrets.

> **Story:** As a Platform Admin, I define an `ExternalSecret` resource for a third-party API credential stored in Vault. A Kubernetes Secret appears in the target namespace automatically and stays in sync — whenever the value in Vault changes, the K8s Secret updates without any manual intervention.

- Installed in the K8s cluster
- Watches Vault for secret changes and syncs them into K8s Secrets
- Pods consume K8s Secrets as env vars — no direct Vault dependency in application code

### 2.2 ArgoCD
GitOps sync engine. Watches the monorepo's L2 artifacts and syncs them to the K8s cluster.

> **Story:** As a Platform Engineer, I push a change to the L2 artifacts in the monorepo. Within seconds, ArgoCD detects the change and syncs it to the K8s cluster automatically — no `kubectl apply` is ever run manually. The cluster state always reflects what is in Git.

- Installed in the K8s cluster via Helm
- Syncs L2 artifacts (K8s manifests, Helm charts) from the monorepo to the target cluster
- Is the deployment actuator for all Workloads — all deploys go through ArgoCD, not `kubectl apply` directly

### 2.3 Argo Rollouts
Progressive delivery controller used for blue-green application deployments.

> **Story:** As a Platform Admin, I run `freelunch install` and the Argo Rollouts controller and CRDs are installed alongside ArgoCD. When ArgoCD syncs a generated `Rollout`, the controller maintains its active and preview versions and reports promotion state back to FreeLunch.

- Installed in the K8s cluster via Helm
- Separate responsibility from ArgoCD: ArgoCD syncs desired state; Argo Rollouts manages active and preview ReplicaSets and Service selectors
- FreeLunch generates standard `Rollout`, active Service, preview Service, and analysis resources in L2
- Group 5 defines promotion, validation, and rollback behavior

### 2.4 4 Personas + Permission Model
IAM-inspired capability-based permission system.

> **Story:** As a Platform Admin, I assign Alice to the Developer Persona in Keycloak. Alice can model Services and Workflows on the canvas and open PRs, but cannot edit L2 artifacts or approve her own PRs. When I grant Alice the hotfix permission, she can merge a branch directly to `main` without going through the CI pipeline. Bob, a Platform Engineer, can edit L2 artifacts directly and configure platform policies — things Alice cannot do.

- **4 Personas:**
  - **Platform Admin** — full access; manages Roles, permissions, setup
  - **Platform Engineer** — configures platform policies, edits L2 directly, and may use advanced CUE editing
  - **Developer** — models Services and Workflows on the canvas
  - **Tech Lead** — Developer with PR merge rights (granted by Platform Admin; not a separate system role)
- **Hotfix** — a permission grantable by Platform Admin to any Persona; allows merging directly to `main` with no CI gates
- Personas map to Keycloak groups
- Permission enforcement spans: GitHub, Theia IDE, and K8s simultaneously
- When permissions conflict across systems, most restrictive applies until Platform Admin resolves it

---

## Group 3 — L1/L2 Model Core
*Depends on: Group 1 (monorepo structure), Group 2 (ArgoCD and Argo Rollouts for deployment validation)*

### 3.1 L1 Abstraction Schema
The CUE-based schema validating what the canvas maintains for Developers.

> **Story:** As a Developer, I add an API Service block to the canvas and specify its name, port, and language. FreeLunch deterministically maintains `products/my-product/services/api-server.cue` and validates it against the L1 schema. Cloud Native Buildpacks turn my source code into a container image without a Dockerfile; the CUE records the Workload configuration, metadata, and image reference that FreeLunch compiles into a versioned Helm chart. I do not write K8s YAML or CUE by hand.

- Defines the L1 types: `Service`, `Workflow`, and platform config
- One deterministic CUE file per Workload is committed to Git and remains usable by the CLI and CI without the canvas
- The canvas supports the FreeLunch-defined CUE subset; raw CUE is inspectable and available through an explicit advanced mode
- Unsupported advanced CUE expressions are surfaced and never silently overwritten by the canvas
- Platform Engineers configure defaults, constraints, and enforced policies through the Platform Policy Editor
- Creating entirely new L1 abstraction types from scratch is **post-Demo**
- Cloud Native Buildpacks own source-to-image builds; L1 CUE owns Workload configuration, metadata, and the resulting image reference
- A Service definition in L1 requires no Dockerfile and no K8s YAML from the Developer

### 3.2 L1→L2 Compilation Engine
Manually triggered engine that compiles L1 configuration, metadata, and image references into versioned Helm artifacts.

> **Story:** As a Developer, I click "Compile" in the Theia IDE after saving my Service on the canvas. FreeLunch validates the canvas-maintained CUE and produces K8s manifests and Helm charts in the `platform/` directory. If a Platform Engineer previously edited a Deployment manifest in L2 directly, FreeLunch shows a conflict diff and waits for resolution — it never silently overwrites the manual change.

- **Manually triggered** by the developer (via Theia command or `freelunch compile`)
- Validates canvas-maintained L1 CUE against the versioned schema and platform policies
- Packages the generated K8s manifests and Argo Rollouts resources as a versioned Helm chart
- Compiles L1 → L2 (K8s manifests, Helm charts, and Argo Rollouts resources)
- **Conflict detection:** If the customer has directly edited L2, the engine detects conflicts between the new L1 output and existing L2 manual overrides. Customer sees a diff and resolves conflicts before accepting the new L2 state
- Raw L2 edits are **never silently discarded** — they are always surfaced for review
- L2 is the source of truth for deployment; ArgoCD deploys whatever is in L2

### 3.3 L2 Artifact Management
The customer-visible layer of K8s/Helm artifacts.

> **Story:** As a Platform Engineer, I open the `platform/` directory in the Theia IDE, edit a Deployment manifest directly to add a custom sidecar container, and commit it. The next time a Developer compiles their L1 changes, the sidecar appears in the conflict diff as a manual override — it is never silently removed, and its git history shows exactly who added it and when.

- Customer can inspect and edit L2 artifacts directly at any time
- L2 lives in the monorepo (versioned via Git) — all changes are auditable
- Platform Engineers may edit L2 directly; those edits persist across L1 recompilations
- L2 → ArgoCD → K8s cluster is the deploy path

---

## Group 4 — CI/CD Pipeline
*Depends on: Group 1 (monorepo structure), Group 3 (L1/L2 model), Group 2 (permissions)*

### 4.1 GitHub Actions Pipeline Structure
The managed CI/CD pipeline installed in the monorepo's `.github/workflows/`.

> **Story:** As a Developer, I open a PR against `develop`. Without any manual setup, GitHub Actions runs image build, security scan, unit tests, contract tests, functional tests, load tests, and a PR compliance check in sequence. After code review and remote integration tests pass, a PR to `main` is opened automatically. Once smoke and e2e load tests pass on the ephemeral staging environment, the PR to `main` is accepted and a blue-green deploy to prod is triggered — all without me touching a pipeline file.

Full pipeline for non-hotfix PRs:
```
Local
  └─ unit tests + integration tests (dev environment, testcontainers / WireMock)
  └─ pre-commit hooks

PR opened → develop branch
  └─ build executables & container images (Buildpacks — no Dockerfile required)
  └─ image static security scan
  └─ remote unit tests (mocks)
  └─ API contract tests
  └─ API functional tests (mocks)
  └─ API load tests
  └─ PR compliance check (PR type annotated, commit message format, PR text format)

PR Review

→ remote integration tests
  └─ if passed: PR to develop accepted + PR to main opened

Deploy Day (Ephemeral Staging Environment)
  └─ smoke tests
  └─ e2e load tests
  └─ production-like observability
  └─ if passed: PR to main accepted

→ Blue-Green Deployment to K8s prod
→ Manual rollback available
```

Hotfix path: merge directly to `main` → deploy to prod. No CI gates. Only Personas with the Hotfix permission can do this.

### 4.2 Pre-commit Hooks
Local enforcement layer, re-enforced server-side.

> **Story:** As a Developer, I run `git commit` with a formatting violation in my code. The pre-commit hook catches it and blocks the commit with a clear error message. I try bypassing with `--no-verify` and push anyway — GitHub Actions detects the bypass server-side and fails the PR with the same error. My code never reaches `develop`.

- **Static checks** — linting, type checking, security smell detection
- **Format enforcement** — code formatter; fails if code changes after format
- **L1 validation** — canvas-maintained CUE must conform to the versioned schema and platform policies
- **Pipeline structure validation** — when a Platform Engineer inserts a new module into the CI/CD pipeline, the hook validates input/output type-signature compatibility with neighboring modules
- **Bypass enforcement** — `--no-verify` bypass is caught by GitHub Actions server-side and blocks the PR. Local bypass never reaches production

### 4.3 Dagger CI Execution
Executes CI pipeline steps (build, test, scan) in a reproducible, container-native way.

> **Story:** As a Developer, I push a Go service with no Dockerfile. Dagger detects the language via Buildpacks, builds a container image automatically, runs the test suite inside the container, scans the image for vulnerabilities, and pushes the image to the registry. I never write or maintain a Dockerfile or a build script.

- Called by GitHub Actions
- Handles: image builds (via Buildpacks), test execution, security scanning
- Buildpacks support: Developers can deploy Services from source code with no Dockerfile

### 4.4 PR Criticality Scale
Classification system applied to every PR.

> **Story:** As a Developer, I open a PR that changes two services and appends a field to a public API contract. FreeLunch automatically classifies it as **High** criticality (multi-workload + interface change). My reviewers see the classification badge and know to apply thorough review and validate staging before approving — no one has to manually decide how risky this PR is.

| Level | Trigger | Implication |
|---|---|---|
| **Low** | Single-Workload change; public interfaces unchanged | Standard review |
| **Medium** | Multi-Workload change; public interfaces unchanged | Broader review scope |
| **High** | Any public API or event-contract change | Careful review, staging analysis |
| **Critical** | Resource usage above threshold, manual infra changes, or Hotfixes | Maximum scrutiny |

Two structural types outside the criticality scale:
- **Telemetry additions** — low blast radius, specialized review
- **CI/CD structural additions** — reviewed by Platform Engineers for pipeline type compatibility

### 4.5 Selective Test Execution
Allows PR Authors to skip integration tests for unaffected Workloads.

> **Story:** As a Developer, I open a PR that only touches `service-a`. When the integration test stage starts, I mark `service-b` and `service-c` integration tests as skipped. The pipeline runs only `service-a`'s integration tests and finishes significantly faster. The build graph confirms those services are unaffected by my change, and the skip is recorded in the PR audit trail.

- PR Author decides which integration tests to skip for Workloads unaffected by the change
- Default: run all tests
- Build graph (via chosen build tool) informs which Workloads are actually affected

---

## Group 5 — Deployment
*Depends on: Group 2 (ArgoCD and Argo Rollouts), Group 3 (L2 artifacts), Group 4 (CI/CD pipeline)*

### 5.1 Argo Rollouts Blue-Green Integration
Controlled application promotion through standard Argo Rollouts resources.

> **Story:** As a Developer, my PR is merged to `main` and ArgoCD syncs the generated `Rollout`. Argo Rollouts deploys green behind a preview Service while blue continues serving through the active Service. After green passes its configured analysis, it is promoted to active. During the configured rollback window, `freelunch rollback` restores the retained previous revision and reconciles the selected revision in Git.

- FreeLunch generates a standard `Rollout`, active Service, preview Service, and optional analysis templates in L2
- ArgoCD syncs the resources; Argo Rollouts manages ReplicaSets, Service selectors, promotion pauses, and analysis state
- Demo traffic switching is scoped to active and preview Services in the same namespace
- Promotion occurs only after readiness and configured pre-promotion analysis succeed
- **Manual rollback:** restore a retained previous revision within the configured rollback window and reconcile Git with cluster state
- Provider propagation delays and rollback eligibility are surfaced rather than described as universally instant or zero-downtime
- Stateful migration, database or queue synchronization, and schema migration orchestration are outside the MVP; this feature covers application rollout only

### 5.2 Ephemeral Staging Environments
One isolated environment per PR for parallel testing.

> **Story:** As a Developer, I open a PR and FreeLunch automatically provisions an isolated staging environment scoped to my PR. I run smoke tests against it without touching any other developer's staging work. When my PR is merged, the environment tears itself down automatically — I never manually create or delete a staging environment.

- Spun up when a PR enters the Deploy Day stage
- Torn down automatically after PR is merged or closed
- Multiple PRs can be tested in parallel without environment contamination
- Uses the same K8s cluster (different namespace) or a dedicated Kind cluster per PR — implementation detail TBD

### 5.3 Autoscaling
Always-on autoscaling for all Workloads.

> **Story:** As a Platform Engineer, I set `min: 2, max: 10` for `service-a` in the Platform Policy Editor. When traffic spikes, K8s scales the pods up to 10 automatically. A Developer who attempts to change the HPA policy on the canvas without explicit permission sees the change rejected — autoscaling policy stays in platform hands.

- Pod vertical and horizontal autoscaling (VPA + HPA)
- Node autoscaling (Karpenter or equivalent)
- Platform Engineers configure scaling policies; Developers cannot change them unless explicitly granted that permission by the Platform Admin

---

## Group 6 — Observability
*Depends on: Group 1 (K8s), Group 5 (deployed Workloads)*

### 6.1 SigNoz Setup
Platform-level observability backend for customer Workload telemetry.

> **Story:** As a Platform Admin, after running `freelunch install`, I open the SigNoz UI in my browser and see metrics, logs, and traces flowing from the sample application's pods. FreeLunch's own internal components (Keycloak, ArgoCD, Vault) do not appear in SigNoz — only customer Workloads are visible.

- Installed in the K8s cluster
- Collects: metrics, logs, traces from customer-deployed Workloads
- **Scope:** Customer workloads only. FreeLunch's own internal components (SigNoz itself, Keycloak, ArgoCD, Vault) are **not monitored** in the Demo — FreeLunch internal observability is out of Demo scope.

### 6.2 Customer Workload Observability
End-to-end observability for services the customer deploys via FreeLunch.

> **Story:** As a Developer, I deploy `service-a` and open the Theia IDE's observability panel. I see CPU/memory usage for its pods, the current CI/CD pipeline stage for my latest PR, and the distributed traces my service is emitting via its existing OpenTelemetry SDK — all in one place, without changing a single line of instrumentation code.

- **K8s infra metrics** — pod/node health, resource usage for customer Workloads
- **CI/CD pipeline visibility** — GitHub Actions + Dagger pipeline state
- **Application-level telemetry** — customer instruments their code via OpenTelemetry SDKs; FreeLunch provides the OTel instrumentation layer and routes to SigNoz
- Existing OTel-compatible telemetry continues to work without re-instrumentation
- All observability data scoped to the customer's own Workloads

### 6.3 Cost Observability
Per-Workload cost breakdown (visibility only — no budget enforcement in Demo).

> **Story:** As a Platform Admin, I open the cost panel in the Theia IDE and see OpenCost's estimated month-to-date compute cost for every product and Workload, including idle and unallocated cluster cost. I can identify that `service-b` is the largest active consumer and inspect its CPU, memory, storage, and network allocation. No alerts fire and no deployments are blocked — it is visibility only.

- For the Demo, OpenCost is installed via Helm with its Prometheus backend while SigNoz remains the platform observability backend
- FreeLunch-generated resources carry stable `freelunch.io/product`, `freelunch.io/workload`, and `freelunch.io/workload-id` labels
- The Demo embeds the OpenCost and SigNoz UIs in FreeLunch through Theia-integrated plugins so users can inspect cost and observability without leaving the IDE
- The local Kind Demo uses explicit custom CPU, memory, storage, and network prices; the UI labels all values as estimates
- Idle and unallocated costs remain visible so Workload percentages have a defined cluster-cost denominator
- Cost profiling surfaces through the Theia IDE and Agent Integration Layer; budget enforcement, alerts, and deployment gates are out of Demo scope
- **Post-Demo:** send OpenCost metrics to the SigNoz backend, remove the dedicated Prometheus backend and OpenCost UI, and replace them with a unified FreeLunch cost experience that queries SigNoz

---

## Group 7 — Theia IDE + CLI
*Depends on: Group 1 (Keycloak), Group 3 (L1/L2 model), Group 4 (CI/CD pipeline), Group 6 (observability)*

### 7.1 Theia IDE
Primary user interface — IDE and Dev Portal are the same unified surface.

> **Story:** As a Developer, I open the Theia IDE in my browser and log in with my Keycloak credentials. On the canvas I configure a Service without writing CUE, trigger compilation, view my pipeline and Argo Rollouts state, inspect pod health, read logs, browse estimated cost data, and install Open VSX extensions — without switching tools or opening a separate Dev Portal.

- Built on Eclipse Theia (backwards-compatible with VS Code extensions)
- **Open VSX extensions:** developers can install any extension from the Open VSX registry freely; this is a Theia configuration concern, not a feature FreeLunch builds from scratch
- FreeLunch ships with a canvas for L1 Workload authoring and hosts the Platform Policy Editor defined in Group 9 for authorized Platform Engineers
- Canvas changes deterministically maintain the Git-backed CUE files; supported fields round-trip between the canvas and CUE
- Raw CUE remains inspectable and can be unlocked through an explicit advanced mode; unsupported expressions are never silently overwritten
- Additional integrated surfaces include the L1→L2 compilation trigger, Git conflict diff, pipeline and Argo Rollouts viewers, observability panels, cost panel, and Workload status
- All modifications to the system happen via Git (GitOps) or CLI — the IDE is the observation and editing surface, not an action dispatcher
- Authenticated via Keycloak SSO

### 7.2 FreeLunch CLI
Minimal CLI for setup and day-to-day inspection.

> **Story:** As a Platform Admin, I run `freelunch init my-company` to scaffold the monorepo, then `freelunch install` to bring up the full FreeLunch stack locally. Later, from my terminal, I run `freelunch status` and see the health of every Workload and environment at a glance.

**Setup commands** (run by Platform Admin):
- `freelunch init` — bootstrap a new monorepo with FreeLunch structure
- `freelunch install` — install FreeLunch components into the K8s cluster, including Theia, ArgoCD, Argo Rollouts, SigNoz, OpenCost with its Prometheus backend, Keycloak, and Vault
- `freelunch install --adopt` — install onto an existing cluster without disruption
- `freelunch configure` — set IP whitelists, cluster targets, rollback policies, Role permission overrides

**Inspection commands** (any authorized Persona):
- `freelunch status` — health of all Workloads, environments, and pipeline
- `freelunch logs` — tail logs for a Workload
- `freelunch rollback` — trigger a manual rollback to the previous blue

---

## Group 8 — Agent Integration Layer
*Depends on: Group 1 (Keycloak), Group 5 (deployment state), Group 6 (observability)*

### 8.1 Coding Agent API
Canonical read-only HTTP API for coding agents to query platform state outside of code.

> **Story:** As a coding agent integrated into a developer's workflow, I obtain a Keycloak client credential token and call `GET /api/v1/workloads` to see the status of all Services and Workflows. I call `GET /api/v1/pipeline/pr/42` to get the test results and current stage of PR #42. I call `GET /api/v1/costs` to get a per-Workload cost breakdown. I cannot trigger deploys, create records, or modify anything — every endpoint is read-only.

- **Auth:** Keycloak client credentials (machine-to-machine OAuth2 flow)
- Publishes an OpenAPI contract for endpoint discovery, typed clients, and compatibility checks
- **Data exposed:**
  - Workload statuses (running, degraded, failing)
  - Pipeline state (current PR stage, test results)
  - Errors (recent errors from logs/traces)
  - Cost data (per Workload cost breakdown)
  - Infra observability (pod/node health, resource usage)
  - App observability (custom metrics/traces the customer has instrumented)
- Read-only authorization is enforced by Keycloak scopes and server routes; no write endpoints are registered
- **Out of Demo scope:** ticket creation, notifications, agent management, and every platform mutation

### 8.2 First-Party FreeLunch Skill
Documentation-backed workflows for using FreeLunch interactively or headlessly.

> **Story:** As a Developer, I install the FreeLunch skill in my coding agent. I can ask how to model and compile a Workload, inspect the same repository and platform state exposed by the IDE, or diagnose why `service-b` is degraded. The skill combines packaged FreeLunch documentation with the authenticated, read-only Coding Agent API, structured CLI output, and Git-backed project files so the workflows also work headlessly.

- Packages the FreeLunch documentation, concepts, repository conventions, and operational runbooks for agent use
- Uses the OpenAPI-described REST API in 8.1 for authorized platform reads and structured CLI or Git-backed files for headless workflows
- Covers Workload modeling, source-to-image builds, L1 validation and compilation, pipeline and rollout inspection, and links back to equivalent Theia views
- The skill correlates Workload health, pipeline state, recent errors, cost data, and observability summaries through the authorized API
- The first-party skill defines diagnostic workflows such as degraded-Workload analysis, PR readiness, deployment regression analysis, and monthly cost investigation
- Credentials are provided by the agent runtime or secure environment and are never embedded in the skill
- The skill registers no MCP server; platform mutations continue to follow the documented GitOps and CLI paths

---

## Group 9 — Advanced Platform Features
*Depends on: Group 3 (L1 schema + compilation engine), Group 7 (Theia IDE + CLI)*

### 9.1 Platform Defaults and Policy Editor
Allows Platform Engineers to configure supported defaults, constraints, and enforced policies without normally editing CUE.

> **Story:** As a Platform Engineer, I open the Platform Policy Editor and change the default Service memory request from `256Mi` to `512Mi`. FreeLunch shows which Workloads omit an explicit value, previews their resulting L2 changes, and deterministically updates `platform/overrides.cue` on my branch. After approval and compilation, only those Workloads inherit the new default.

- The editor distinguishes defaults used when values are omitted, allowed constraints, and enforced policies that Developers cannot override
- Supported settings include default ports, resource requests and limits, HPA thresholds, and retry policies within existing FreeLunch abstractions
- Before writing a Git change, FreeLunch shows affected Workloads, unchanged explicit configurations, validation failures, and the generated L2 diff
- `platform/overrides.cue` remains standard, deterministic CUE that is inspectable by customers and usable by CLI and CI
- Raw CUE editing is an explicit advanced mode; unsupported expressions are preserved and surfaced rather than silently rewritten
- **Not in Demo:** creating entirely new L1 abstraction types from scratch (post-Demo)
- Changes are tracked via Git (who changed what, who approved)

---

## Group 10 — Documentation + Sample App
*Depends on: Groups 1–9 (demonstrates all major MVP features)*

### 10.1 Sample Application
A purpose-built Demo application that exercises the full FreeLunch workflow.

> **Story:** As a new customer evaluating FreeLunch, I clone the Demo monorepo and run `freelunch install`. On the canvas I create one stateless Service from source, import a container image as a second Service, connect them, and add a virtual Service block for an externally managed database dependency. I compile the generated CUE and exercise the services while SigNoz shows their telemetry, OpenCost shows their estimated allocation, and the FreeLunch skill can summarize their health through the read-only Coding Agent API.

- Multiple stateless Services: one built from source with Cloud Native Buildpacks and one using an imported container image
- A virtual Service block represents the external database dependency without FreeLunch hosting, migrating, or synchronizing the database
- Includes health, echo, controlled-error, and latency paths for exercising deployment and observability
- Deployed and observable in the Demo monorepo
- Demonstrates: canvas authoring, CUE validation, L1→L2 compilation, GitOps deployment, Argo Rollouts blue-green promotion, observability, cost allocation, and agent diagnosis
- Requires no FreeLunch-hosted database, queue, persistent volume, or stateful-service wiring

### 10.2 Documentation
Auto-generated documentation website + LLM-friendly export.

> **Story:** As a Platform Engineer onboarding to FreeLunch, I open the MKDocs documentation website and find complete guides to the canvas, generated CUE, platform policies, GitOps deployment, Argo Rollouts, observability, cost allocation, and agent integration. I can also download a single `freelunch-docs.md` file for accurate offline or LLM-assisted reference.

- Built from Markdown files via MKDocs
- Includes the sample app, canvas-to-CUE workflow, platform policy model, and operational runbooks
- Full documentation exportable as a single `.md` file for LLM consumption

---

## Group 11 — Low Priority
*Depends on: Group 6 (observability), Group 5 (deployed workloads). Deferred until core groups are complete.*

### 11.1 Coding/Experimenting Environment
A girus-powered sandbox for anticipating production problems.

> **Story:** As a Developer, I open the Coding/Experimenting Environment and select "simulate 10x traffic on `service-a` based on last week's p95 patterns." FreeLunch seeds the girus simulation with real SigNoz metrics from last week and I observe how `service-a` behaves under that load before it happens in production — without touching the live cluster or writing any simulation config manually.

- Uses girus to create interactive, scenario-based K8s simulation environments
- **Connected to live observability data** — reads SigNoz/OTel metrics from the customer's deployed Workloads to seed simulation state
- Developer can say "simulate 10x traffic based on last week's patterns" and observe system behavior before it hits production
- Running environment stack: `linux/wsl, pixi, git, girus`
- In Demo scope but **explicitly low priority** — implement only after all other groups are complete

---

## Summary Table

| Group | Features | Key Dependency |
|---|---|---|
| **1** | Monorepo structure, Kind/LocalStack, Keycloak, Vault | None |
| **2** | External-secrets-operator, ArgoCD, Argo Rollouts, 4 Personas + permissions | Group 1 |
| **3** | Canvas-maintained L1 schema, L1→L2 compilation engine, L2 artifact management | Groups 1–2 |
| **4** | GitHub Actions pipeline, pre-commit hooks, Dagger, PR Criticality, selective test execution | Groups 1–3 |
| **5** | Argo Rollouts blue-green integration, ephemeral staging, autoscaling | Groups 2–4 |
| **6** | SigNoz, customer workload observability, OpenCost cost observability | Groups 1, 5 |
| **7** | Theia canvas and integrated platform surfaces, FreeLunch CLI | Groups 1, 3, 4, 6 |
| **8** | Coding Agent API, OpenAPI, first-party skill | Groups 1, 5, 6 |
| **9** | Platform Policy Editor | Groups 3, 7 |
| **10** | Stateless sample app, documentation | Groups 1–9 |
| **11** | Coding/Experimenting Environment (low priority) | Groups 5, 6 |
