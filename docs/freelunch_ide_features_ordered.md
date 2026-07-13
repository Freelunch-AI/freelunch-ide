# FreeLunch — Demo/MVP Features: Implementation Order

## Glossary (Revised)

| Term | Definition |
|---|---|
| **Monorepo** | The single Git repository for a customer organisation. Contains all products, services, workflows, and platform config. One monorepo per customer. Replaces "Product Repo." |
| **Workload** | Any single deployable unit inside the monorepo. Either a Service or a Workflow. |
| **Service** | A long-running, always-on, request-driven Workload (HTTP, gRPC, event consumer, etc.). |
| **Workflow** | A trigger-driven Workload that runs to completion. Has per-run execution state. |
| **Layer 1 (L1)** | FreeLunch's abstraction API. What Developers write. Manually compiled into L2. |
| **Layer 2 (L2)** | K8s/Helm artifacts. Source of truth for deployment. Customers can view and edit directly. Raw edits are preserved on recompile via conflict detection. |
| **Persona** | One of four pre-defined roles: Platform Admin, Platform Engineer, Developer, Tech Lead. |
| **Hotfix** | A permission grantable by Platform Admin to any Persona. Allows merging directly to main without CI gates. |
| **Ephemeral Staging Environment** | A short-lived isolated environment spun up per PR, torn down after merge/close. |
| **Blue-Green Deployment** | Traffic switch between old (blue) and new (green) version. Rollback = switch back to blue. |
| **Eject** | Customer leaves FreeLunch. Receives L2 artifacts + GitOps pipeline intact. FreeLunch-installed components remain in the cluster, owned by the customer. |
| **Coding Agent API** | Read-only HTTP API for coding agents to query workload statuses, pipeline state, errors, costs, and observability data. Authenticated via Keycloak client credentials. |
| **Coding/Experimenting Environment** | A girus-powered production simulation sandbox. Reads live SigNoz/OTel data to seed simulations. Used to anticipate production problems. |

---

## Deferred Decisions (Pending Team Input)

- **Public/Private Hub for reusable blocks** — whether it is in the Demo or post-MVP. If included, it is a Theia-integrated registry for sharing Workload templates across teams.
- **Vanilla K8s vs LocalStack EKS** — both are supported; which the Demo primarily targets is TBD.
- **Specific tool choices** (Kargo, Bazel, etc.) — conceptual roles are defined; final tool selection happens during implementation of each feature.

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
│       ├── services/      ← L1 Service definitions (Developers)
│       └── workflows/     ← L1 Workflow definitions (Developers)
└── .github/
    └── workflows/         ← CI/CD pipeline (managed by FreeLunch)
```

- Platform versioning declared in `platform/` (e.g., `freelunch.yaml: platform-version: x.y.z`)
- Separation between platform config (Platform Engineers) and Workload definitions (Developers)

### 1.2 Local K8s Infrastructure
The fully-local runtime environment for the Demo.

> **Story:** As a Platform Admin, I run `freelunch install` on a fresh machine with Docker installed and a running Kind cluster appears locally with LocalStack emulating AWS services alongside it. All FreeLunch components deploy into it without touching the internet. I can also run `freelunch install --adopt` against an already-running cluster and FreeLunch installs without disrupting existing workloads.

- **Kind** — vanilla Kubernetes-in-Docker as primary local K8s target
- **LocalStack** — emulates AWS services (S3, IAM, SQS, RDS, etc.)
- Two install modes supported:
  - **Fresh** — `freelunch install` spins up a new Kind cluster
  - **Adopt existing** — `freelunch install --adopt` installs FreeLunch components onto an already-running cluster without disrupting existing workloads
- EKS (via LocalStack) is also a supported target; final Demo target (Kind vs LocalStack EKS) is TBD

### 1.3 Keycloak (Local Instance)
OIDC identity provider for all human authentication within FreeLunch.

> **Story:** As a Platform Admin, I open the local Keycloak admin console, create accounts for my team, and assign each user to the correct Persona group. Those users can then log into the Theia IDE with their credentials on the next login — no separate account setup required in FreeLunch itself.

- Self-hosted local Keycloak instance
- Used by: Theia IDE (developer login), Coding Agent API (machine-to-machine client credentials)
- The 4 Personas map to Keycloak groups
- Does **not** handle end-user auth for customer applications (out of Demo scope)

### 1.4 Vault (Local Instance)
Application secrets store for all Workload credentials.

> **Story:** As a Platform Admin, I store a database password in Vault at path `secret/my-service/db-password`. A pod running `my-service` starts with a `DB_PASSWORD` environment variable already populated — no Vault SDK in the application code, no Kubernetes Secret created manually.

- Self-hosted local Vault instance (not an AWS service)
- Stores: DB connection strings, API keys, service credentials provisioned by the customer
- External-secrets-operator (Group 2) syncs secrets from Vault → K8s native Secrets
- Pods consume secrets as env vars and never communicate with Vault directly

---

## Group 2 — Core Infrastructure Components
*Depends on: Group 1 (K8s cluster, Keycloak, Vault)*

### 2.1 External-secrets-operator
Bridge between Vault and K8s native secrets.

> **Story:** As a Platform Admin, I define an `ExternalSecret` resource that points to a path in Vault. A Kubernetes Secret appears in the target namespace automatically and stays in sync — whenever the value in Vault changes, the K8s Secret updates without any manual intervention.

- Installed in the K8s cluster
- Watches Vault for secret changes and syncs them into K8s Secrets
- Pods consume K8s Secrets as env vars — no direct Vault dependency in application code

### 2.2 ArgoCD
GitOps sync engine. Watches the monorepo's L2 artifacts and syncs them to the K8s cluster.

> **Story:** As a Platform Engineer, I push a change to the L2 artifacts in the monorepo. Within seconds, ArgoCD detects the change and syncs it to the K8s cluster automatically — no `kubectl apply` is ever run manually. The cluster state always reflects what is in Git.

- Installed in the K8s cluster via Helm
- Syncs L2 artifacts (K8s manifests, Helm charts) from the monorepo to the target cluster
- Is the deployment actuator for all Workloads — all deploys go through ArgoCD, not `kubectl apply` directly

### 2.3 4 Personas + Permission Model
IAM-inspired capability-based permission system.

> **Story:** As a Platform Admin, I assign Alice to the Developer Persona in Keycloak. Alice can write L1 Service definitions and open PRs, but cannot edit L2 artifacts or approve her own PRs. When I grant Alice the hotfix permission, she can merge a branch directly to `main` without going through the CI pipeline. Bob, a Platform Engineer, can edit L2 artifacts directly and tune L1 abstractions — things Alice cannot do.

- **4 Personas:**
  - **Platform Admin** — full access; manages Roles, permissions, setup
  - **Platform Engineer** — configures platform, edits L2 directly, tunes L1 abstractions via CUE
  - **Developer** — writes Services and Workflows in L1
  - **Tech Lead** — Developer with PR merge rights (granted by Platform Admin; not a separate system role)
- **Hotfix** — a permission grantable by Platform Admin to any Persona; allows merging directly to `main` with no CI gates
- Personas map to Keycloak groups
- Permission enforcement spans: GitHub, Theia IDE, and K8s simultaneously
- When permissions conflict across systems, most restrictive applies until Platform Admin resolves it

---

## Group 3 — L1/L2 Model Core
*Depends on: Group 1 (monorepo structure), Group 2 (ArgoCD for deployment validation)*

### 3.1 L1 Abstraction Schema
The CUE-based schema defining what Developers write.

> **Story:** As a Developer, I create a new file `products/my-product/services/api-server.cue` with a simple L1 Service definition specifying name, port, and language. No Dockerfile and no K8s YAML — just the FreeLunch abstraction. The file validates against the L1 schema and is ready to compile.

- Defines the L1 types: `Service`, `Workflow`, and platform config
- Platform Engineers can **tune/override existing** L1 abstractions (change defaults, adjust policies, resource limits, HPA thresholds)
- Creating entirely new L1 abstraction types from scratch is **post-Demo**
- A Service definition in L1 requires no Dockerfile and no K8s YAML from the Developer

### 3.2 L1→L2 Compilation Engine
Manually triggered engine that compiles L1 abstractions into K8s/Helm artifacts.

> **Story:** As a Developer, I click "Compile" in the Theia IDE after editing my L1 Service definition. FreeLunch produces K8s manifests and Helm charts in the `platform/` directory. If I had previously edited a Deployment manifest in L2 directly, FreeLunch shows me a conflict diff and waits for my resolution — it never silently overwrites my manual changes.

- **Manually triggered** by the developer (via Theia command or `freelunch compile`)
- Compiles L1 → L2 (K8s manifests, Helm charts)
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

> **Story:** As a Developer, I open a PR that changes two services and appends a new column to a database schema. FreeLunch automatically classifies it as **High** criticality (multi-workload + schema change). My reviewers see the classification badge and know to apply thorough review and validate staging before approving — no one has to manually decide how risky this PR is.

| Level | Trigger | Implication |
|---|---|---|
| **Low** | Single-Workload change; interfaces and schema unchanged | Standard review |
| **Medium** | Multi-Workload change; interfaces and schema unchanged | Broader review scope |
| **High** | Any interface or schema append change | Careful review, staging analysis |
| **Critical** | Storage engine changes, resource usage above threshold, manual infra changes, or Hotfixes | Maximum scrutiny |

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

## Group 5 — Deployment Engine
*Depends on: Group 2 (ArgoCD), Group 3 (L2 artifacts), Group 4 (CI/CD pipeline)*

### 5.1 Blue-Green Deployment Engine
Zero-downtime deployment via traffic switching.

> **Story:** As a Developer, my PR is merged to `main`. The new version (green) is deployed alongside the running version (blue), and once green passes health checks, the load balancer switches all traffic to it with zero downtime. If I run `freelunch rollback`, traffic switches back to blue instantly. Green started with the same replica count blue was running at the moment of cutover — not the L1-configured minimum.

- New version (green) deployed alongside old version (blue)
- Traffic switch via K8s load balancer / LocalStack ALB target group switching
- **Manual rollback:** switch traffic back to blue (previous `main` commit before last merge)
- **Scaling level sync:** green inherits blue's live autoscaling state at cutover (replica count, HPA settings) — not the L1-defined minimum
- **DB sync (schema):** FreeLunch enforces backward-compatible schema migrations only (append-only). Blue can still read the DB after rollback. This is enforcement, not data copying.

### 5.2 Ephemeral Staging Environments
One isolated environment per PR for parallel testing.

> **Story:** As a Developer, I open a PR and FreeLunch automatically provisions an isolated staging environment scoped to my PR. I run smoke tests against it without touching any other developer's staging work. When my PR is merged, the environment tears itself down automatically — I never manually create or delete a staging environment.

- Spun up when a PR enters the Deploy Day stage
- Torn down automatically after PR is merged or closed
- Multiple PRs can be tested in parallel without environment contamination
- Uses the same K8s cluster (different namespace) or a dedicated Kind cluster per PR — implementation detail TBD

### 5.3 Autoscaling
Always-on autoscaling for all Workloads.

> **Story:** As a Platform Engineer, I set `min: 2, max: 10` in the platform config for `service-a`. When traffic spikes, K8s scales the pods up to 10 automatically. A Developer who tries to change the HPA settings in their L1 definition without explicit permission sees their change rejected — autoscaling policy stays in platform hands.

- Pod vertical and horizontal autoscaling (VPA + HPA)
- Node autoscaling (Karpenter or equivalent)
- Platform Engineers configure scaling policies; Developers cannot change them unless explicitly granted that permission by the Platform Admin

---

## Group 6 — Secrets + Stateful Service Wiring
*Depends on: Group 1 (Vault), Group 2 (external-secrets-operator), Group 3 (L1 schema + compilation engine)*

### 6.1 Stateful Service Wiring via Annotations
FreeLunch wires customer-provisioned stateful services to Workloads automatically.

> **Story:** As a Developer, I add `db: postgres://my-prod-db` as an annotation in my L1 Service definition. After running compile, my pod starts with a `DATABASE_URL` environment variable already populated — no IAM policy written, no Vault entry created, no K8s Secret defined by me. FreeLunch generated all of that from the single annotation.

- Customer **provisions** their own stateful services (RDS, DynamoDB, SQS, PostgreSQL, etc.) independently
- Customer declares the dependency in L1 via annotations/tags (e.g., `db: postgres://my-database`)
- On L1→L2 compilation, FreeLunch generates:
  - IAM role for pod access (IRSA pattern)
  - Vault secret entry for credentials
  - K8s Secret (synced by external-secrets-operator from Vault)
  - Env var injection into the pod spec
- The Service's application code reads the env var — it never knows about Vault, IAM, or the provisioning layer

---

## Group 7 — Observability
*Depends on: Group 1 (K8s), Group 5 (deployed Workloads)*

### 7.1 SigNoz Setup
Platform-level observability backend for customer Workload telemetry.

> **Story:** As a Platform Admin, after running `freelunch install`, I open the SigNoz UI in my browser and see metrics, logs, and traces flowing from the sample application's pods. FreeLunch's own internal components (Keycloak, ArgoCD, Vault) do not appear in SigNoz — only customer Workloads are visible.

- Installed in the K8s cluster
- Collects: metrics, logs, traces from customer-deployed Workloads
- **Scope:** Customer workloads only. FreeLunch's own internal components (SigNoz itself, Keycloak, ArgoCD, Vault) are **not monitored** in the Demo — FreeLunch internal observability is out of Demo scope.

### 7.2 Customer Workload Observability
End-to-end observability for services the customer deploys via FreeLunch.

> **Story:** As a Developer, I deploy `service-a` and open the Theia IDE's observability panel. I see CPU/memory usage for its pods, the current CI/CD pipeline stage for my latest PR, and the distributed traces my service is emitting via its existing OpenTelemetry SDK — all in one place, without changing a single line of instrumentation code.

- **K8s infra metrics** — pod/node health, resource usage for customer Workloads
- **CI/CD pipeline visibility** — GitHub Actions + Dagger pipeline state
- **Application-level telemetry** — customer instruments their code via OpenTelemetry SDKs; FreeLunch provides the OTel instrumentation layer and routes to SigNoz
- No re-instrumentation required on migration: existing OTel-compatible telemetry continues to work
- All observability data scoped to the customer's own Workloads

### 7.3 Cost Observability
Per-Workload cost breakdown (visibility only — no budget enforcement in Demo).

> **Story:** As a Platform Admin, I open the cost panel in the Theia IDE and see a breakdown of compute costs per Workload for the current month. I can tell the team that `service-b` is consuming 60% of the cluster cost. No alerts fire and no deployments are blocked — it is visibility only.

- Cost profiling per Workload and per product
- Powered by OpenCost or equivalent tooling
- Surfaces in Theia IDE and Coding Agent API

---

## Group 8 — Theia IDE + CLI
*Depends on: Group 1 (Keycloak), Group 3 (L1/L2 model), Group 4 (CI/CD pipeline), Group 7 (observability)*

### 8.1 Theia IDE
Primary user interface — IDE and Dev Portal are the same unified surface.

> **Story:** As a Developer, I open the Theia IDE in my browser and log in with my Keycloak credentials. In one unified interface I can edit L1 Service definitions, trigger compilation, view my pipeline status, inspect pod health, read logs, browse cost data, and install any Open VSX extension — without switching tools or opening a separate Dev Portal.

- Built on Eclipse Theia (backwards-compatible with VS Code extensions)
- **Open VSX extensions:** developers can install any extension from the Open VSX registry freely; this is a Theia configuration concern, not a feature FreeLunch builds from scratch
- FreeLunch ships with pre-installed extensions for: L1 editor, L1→L2 compilation trigger, pipeline viewer, observability panels (metrics, logs, traces, costs), Workload status
- All modifications to the system happen via Git (GitOps) or CLI — the IDE is the observation and editing surface, not an action dispatcher
- Authenticated via Keycloak SSO

### 8.2 FreeLunch CLI
Minimal CLI for setup and day-to-day inspection.

> **Story:** As a Platform Admin, I run `freelunch init my-company` to scaffold the monorepo, then `freelunch install` to bring up the full FreeLunch stack locally (Kind, Keycloak, Vault, ArgoCD, SigNoz, Theia). Later, from my terminal, I run `freelunch status` and see the health of every Workload and environment at a glance — without opening the IDE.

**Setup commands** (run by Platform Admin):
- `freelunch init` — bootstrap a new monorepo with FreeLunch structure
- `freelunch install` — install FreeLunch components into the K8s cluster (Theia, ArgoCD, SigNoz, Keycloak, Vault, etc.)
- `freelunch install --adopt` — install onto an existing cluster without disruption
- `freelunch migrate` — migrate an existing repo/cluster into FreeLunch
- `freelunch configure` — set IP whitelists, cluster targets, rollback policies, Role permission overrides
- `freelunch upgrade` — apply a new platform version; flags breaking changes for Platform Engineer resolution

**Inspection commands** (any authorized Persona):
- `freelunch status` — health of all Workloads, environments, and pipeline
- `freelunch logs` — tail logs for a Workload
- `freelunch rollback` — trigger a manual rollback to the previous blue
- `freelunch eject` — export L2 artifacts and exit FreeLunch

---

## Group 9 — Coding Agent API
*Depends on: Group 1 (Keycloak), Group 7 (observability), Group 5 (deployment state)*

### 9.1 Coding Agent API
Read-only HTTP API for coding agents to query all platform state outside of code.

> **Story:** As a coding agent integrated into a developer's workflow, I obtain a Keycloak client credential token and call `GET /api/v1/workloads` to see the status of all Services and Workflows. I call `GET /api/v1/pipeline/pr/42` to get the test results and current stage of PR #42. I call `GET /api/v1/costs` to get a per-Workload cost breakdown. I cannot trigger deploys, create records, or modify anything — every endpoint is read-only.

- **Auth:** Keycloak client credentials (machine-to-machine OAuth2 flow)
- **Data exposed:**
  - Workload statuses (running, degraded, failing)
  - Pipeline state (current PR stage, test results)
  - Errors (recent errors from logs/traces)
  - Cost data (per Workload cost breakdown)
  - Infra observability (pod/node health, resource usage)
  - App observability (custom metrics/traces the customer has instrumented)
- **Out of Demo scope:** ticket creation, notifications, agent management
- Think of it as the Dev Portal, but for agents — all read, no write

---

## Group 10 — Advanced Platform Features
*Depends on: Group 3 (L1 schema + compilation engine)*

### 10.1 CUE Engine for L1 Abstraction Overrides
Allows Platform Engineers to tune and override existing FreeLunch L1 abstractions.

> **Story:** As a Platform Engineer, I open `platform/overrides.cue` and change the default memory limit for all Services from `256Mi` to `512Mi`. After the next compile, all existing L1 Service definitions inherit the new default in their generated L2 output — Developers do not need to change any of their L1 files. I cannot create a new L1 abstraction type from scratch; that is post-Demo.

- Override defaults: change default ports, resource limits, HPA thresholds, retry policies, etc.
- All overrides stay within the existing abstraction structure FreeLunch defines
- **Not in Demo:** creating entirely new L1 abstraction types from scratch (post-Demo)
- Changes are tracked via Git (who changed what, who approved)

### 10.2 Platform Versioning
Declarative platform version management with automated upgrade detection.

> **Story:** As a Platform Engineer, I run `freelunch upgrade` after a new FreeLunch version is published. The CLI diffs the new L1 schema against our current `platform/freelunch.yaml`, reports zero breaking changes, and applies the upgrade with zero downtime. In a separate scenario with a breaking change, the CLI prints exactly what needs to change in our L1 config and blocks the upgrade until I resolve it.

- Platform version declared in `platform/freelunch.yaml` in the monorepo
- FreeLunch team publishes new versioned releases
- `freelunch upgrade` diffs the new version's L1 schema against the current monorepo's L1 config
- Non-breaking upgrades: applied automatically with zero downtime
- Breaking changes: flagged as conflicts; Platform Engineer edits L1 config to be compatible, then re-runs upgrade

### 10.3 Eject Capability
Customer leaves FreeLunch and operates independently.

> **Story:** As a Platform Admin who wants to leave FreeLunch, I run `freelunch eject`. The CLI produces an `./ejected/` directory containing all L2 K8s manifests, Helm charts, ArgoCD config, GitHub Actions workflows, and a `README.md` explaining each file. I take this directory to a plain K8s setup and everything runs without any FreeLunch dependency. Keycloak, Vault, SigNoz, and ArgoCD remain running in my cluster — they are mine now.

- `freelunch eject` produces:
  - All L2 artifacts (K8s manifests, Helm charts) as-is — fully functional without FreeLunch
  - The GitOps pipeline (ArgoCD config + GitHub Actions workflows) intact and runnable
  - A README explaining each component and how to operate independently
- After eject: no L1 compilation, no CLI, no Theia IDE
- FreeLunch-installed components (Keycloak, Vault, SigNoz, ArgoCD) **remain in the customer's cluster** — customer owns and operates them going forward

---

## Group 11 — Migration
*Depends on: Group 5 (blue-green deployment), Group 6 (wiring), Group 7 (observability), Group 8 (Theia IDE)*

### 11.1 Blue-Green Migration Flow
The process for a customer to migrate from their existing K8s setup to FreeLunch.

> **Story:** As a Platform Admin migrating from our old K8s setup, I create the monorepo and define our services in L1. FreeLunch deploys them to a green namespace while our old prod system keeps running as blue. I open the Theia IDE and observe both blue and green side by side in the observability panels. After validating green with smoke tests, I switch the load balancer to green. During the entire migration window I can switch back to blue at any time with zero data loss.

Step-by-step:
1. Customer has existing production system (blue) running on their own K8s setup
2. Customer creates a new monorepo with FreeLunch structure; defines Services in L1
3. FreeLunch compiles L1→L2, deploys Services to a green environment (separate namespace or cluster)
4. FreeLunch sets up DB/queue sync between old DB (blue) and new DB (green) — customer provisioned both
5. Customer uses Theia IDE to observe **both blue and green in parallel**
6. Customer runs smoke tests + observability checks on green
7. Customer switches traffic from blue to green (load balancer switch)
8. After validation period, DB sync is torn down, blue is decommissioned
9. Rollback during migration window = switch load balancer back to blue (sync is still running, no data loss)

### 11.2 DB/Queue Sync During Migration
Temporary bidirectional data sync for safe blue-green migration.

> **Story:** As a Platform Admin, after the green environment is deployed, I run `freelunch migrate --sync db my-postgres` and FreeLunch establishes CDC replication from the old database (blue) to the new database (green). Both stay in sync during the validation window. After I switch traffic to green and confirm it, I run `freelunch migrate --sync teardown` and the replication is cleanly removed — no stale sync processes left behind.

- FreeLunch sets up CDC (Change Data Capture) or replication between old DB (blue) and new DB (green) during the migration window
- Supported engines (to be finalized during implementation): PostgreSQL, MySQL, Redis, SQS/Kafka as minimum set
- Sync is active only during the migration window and is explicitly torn down after traffic cutover is validated
- Bidirectional sync ensures rollback to blue has zero data loss

---

## Group 12 — Documentation + Sample App
*Depends on: Groups 1–11 (demonstrates all major features)*

### 12.1 Sample Application
A purpose-built Demo application that exercises the full FreeLunch workflow.

> **Story:** As a new customer evaluating FreeLunch, I clone the Demo monorepo and run `freelunch install`. The sample HTTP CRUD service is already deployed and visible in the Theia IDE — I send a `POST /items` request to create a record and `GET /items` to retrieve it. SigNoz is showing its metrics and traces. The repo includes a pointer to the legacy repo the sample was migrated from, so I can see exactly what changed.

- A simple HTTP CRUD service (receives requests to be persisted and queried) backed by a database
- Deployed and observable in the Demo monorepo
- The database is customer-provisioned (simulated); FreeLunch wires it via L1 annotations
- Demonstrates: L1 authoring, L1→L2 compilation, GitOps deploy, blue-green deploy, observability, migration path
- Has a corresponding "legacy repo" it was migrated from

### 12.2 Documentation
Auto-generated documentation website + LLM-friendly export.

> **Story:** As a Platform Engineer onboarding to FreeLunch, I open the MKDocs documentation website and find a complete guide — what each monorepo directory does, how to write a Service in L1, the full migration walkthrough with a link to the legacy sample repo. I also download a single `freelunch-docs.md` file and paste it into an AI assistant to get accurate answers about the platform without searching through docs.

- Built from Markdown files via MKDocs
- Includes: sample app explanation, migration walkthrough (with pointer to legacy repo)
- Full documentation exportable as a single `.md` file for LLM consumption

---

## Group 13 — Low Priority
*Depends on: Group 7 (observability), Group 5 (deployed workloads). Deferred until core groups are complete.*

### 13.1 Coding/Experimenting Environment
A girus-powered sandbox for anticipating production problems.

> **Story:** As a Developer, I open the Coding/Experimenting Environment and select "simulate 10x traffic on `service-a` based on last week's p95 patterns." FreeLunch seeds the girus simulation with real SigNoz metrics from last week and I observe how `service-a` behaves under that load before it happens in production — without touching the live cluster or writing any simulation config manually.

- Uses girus to create interactive, scenario-based K8s simulation environments
- **Connected to live observability data** — reads SigNoz/OTel metrics from the customer's deployed Workloads to seed simulation state
- Developer can say "simulate 10x traffic based on last week's patterns" and observe system behavior before it hits production
- Running environment stack: `linux/wsl, pixi, git, girus`
- In Demo scope but **explicitly low priority** — implement only after all other groups are complete

---

## Group 14 — Pending Team Decision

### 14.1 Public/Private Hub for Reusable Blocks *(decision pending)*
A registry/marketplace for sharing reusable Workload templates.

> **Story:** As a Platform Engineer, I publish a reusable `background-worker` Workload template to the Private Hub directly from the Theia IDE. A Developer on another product team browses the Hub, selects the `background-worker` template, and instantiates it in their `products/my-product/workflows/` directory in minutes — no need to write the L1 definition from scratch.

- Listed as part of the IDE's Innovative Features in `features_overview.md`
- Would allow teams to publish and consume reusable Service and Workflow templates
- Team input needed to confirm whether this is Demo scope or post-Demo
- If in Demo: scoped to the Theia IDE, private sharing only within the customer organisation; public Hub is post-Demo

---

## Summary Table

| Group | Features | Key Dependency |
|---|---|---|
| **1** | Monorepo structure, Kind/LocalStack, Keycloak, Vault | None |
| **2** | External-secrets-operator, ArgoCD, 4 Personas + permissions | Group 1 |
| **3** | L1 schema, L1→L2 compilation engine, L2 artifact management | Groups 1–2 |
| **4** | GitHub Actions pipeline, pre-commit hooks, Dagger, PR Criticality, selective test execution | Groups 1–3 |
| **5** | Blue-green deployment, ephemeral staging, autoscaling | Groups 2–4 |
| **6** | Stateful service wiring via annotations | Groups 1–3 |
| **7** | SigNoz, customer workload observability, cost observability | Groups 1, 5 |
| **8** | Theia IDE, FreeLunch CLI | Groups 1, 3, 4, 7 |
| **9** | Coding Agent API | Groups 1, 5, 7 |
| **10** | CUE overrides, platform versioning, eject | Group 3 |
| **11** | Blue-green migration, DB/queue sync | Groups 5, 6, 7, 8 |
| **12** | Sample app, documentation | Groups 1–11 |
| **13** | Coding/Experimenting Environment (low priority) | Groups 5, 7 |
| **14** | Public/Private Hub (pending team decision) | Group 8 |
