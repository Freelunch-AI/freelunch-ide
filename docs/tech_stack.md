# Demo Tech Stack

This is a first-pass mapping of capability to tool choice, derived from the FreeLunch feature order and founding docs.

Notes:

- we might realize during implementation that some tools are worth more borrowing ideas/code than to actually just use them out-of-the-box
- embedded tool UIs doesnt mean we wont provide our code-centric user experience for it, just means we can provide to the user, easily via plugin, a good UI that is already battle tested, as our baseline.
- Most influential decisions: **Theia**, **ProxMox**, **Talos Linux**, **Go**, **Bazel**

## 1. Stack for the FreeLunch repo itself

This covers the tooling used to build, test, document, and ship the FreeLunch package, template, and IDE experience itself.

- Dev OS: **Linux** or **WSL2**
- Experimentation (optional): **Girus**
- Package management / virtual environments for the repo itself → **Pixi**
- Version control and collaboration for the repo itself → **Git** + **GitHub**
- CI/CD for the FreeLunch repo and templates → **GitHub Actions**
- Build orchestration → **Task**
- Static checks, formatting, and pre-commit enforcement → **pre-commit hooks**
- Testing for the repo and platform engine → **Go tests**, with **Testcontainers**/**WireMock** for integration-style coverage and **Loft** using **vCluster** for controlled cluster environments
- Primary language for the core engine / CLI → **Go**
- Primary language for the IDE frontend → **TypeScript**
- Versioned Go package publishing for the FreeLunch CLI/core engine → **GitHub Releases** + **Go module proxy**
- Versioned monorepo template publishing for the FreeLunch CLI bootstrap flow → **GitHub Releases** / **GitHub repository template**
- Documentation site for the repo and product narrative → **MkDocs**

## 2. Stack for the product (the devops platform demo)

This covers the runtime, delivery, and operations stack used to demonstrate the FreeLunch platform to customers.

- Programming Languages Supported: **Any**
- Local cloud-service emulation Cloud Compute Cluster → **ProxMox** running **Talos Linux** with **MetalLB**
- Cloud Workloads Runtime → **Kubernetes** with **Cilium Mesh/Gateway**
- Container registry for built images → **Harbor**
- GitOps deployment engine for L2 artifacts → **ArgoCD**
- Kubernetes manifest packaging for deployed workloads → **Helm**
- Infra Provisioning → **Terraform**
- Kubernetes Installation → comes for free with **Talos Linux**
- Identity provider for human authentication in the IDE → **Keycloak**
- Secrets management for application credentials → **Vault**
- Secret synchronization from Vault into Kubernetes → **external-secrets-operator**
- Abstraction/schema language for platform configuration and L1 definitions → **CUE**
- CI/CD pipeline for customer workloads → **GitHub Actions** + **Dagger** (embedded UI via plugin)
- Container Image Build from source without Dockerfiles → **Cloud Native Buildpacks**
- Ephemeral Environment Management → **Loft** (embedded UI via plugin) using k8s namespaces 
- Progressive delivery / deployment automation → **Kargo** (embedded UI via plugin)
- Autoscaling for workloads → **VPA** + **KEDA** + **Karpenter**
- Observability backend for metrics/logs/traces → **SigNoz** (embedded UI via plugin)
- Distributed tracing / telemetry instrumentation → **OpenTelemetry**
- Kubernetes observability → **Headlamp** (embedded UI via plugin)
- Cost observabilityfor workloads → **OpenCost** + **Prometheus** (embedded UI via plugin)
- Load testing → **K6**
- Test doubles / API mocking for integration tests → **WireMock**
- Kubernetes integration testing → **kubetest**
- Security scanning for images and artifacts → **Trivy**
- IDE/workbench frameowrk for making the IDE → **Eclipse Theia**
- Plugins -> **OpenVSX** (for **VSCode** extensions) + **Backstage** (for embedded control planes) plugins
