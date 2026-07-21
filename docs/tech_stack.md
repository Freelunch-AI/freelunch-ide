# Demo Tech Stack

This is a first-pass mapping of capability to tool choice, derived from the FreeLunch feature order and founding docs.

## 1. Stack for the FreeLunch repo itself

This covers the tooling used to build, test, document, and ship the FreeLunch package, template, and IDE experience itself.

- Dev OS: **Linux** or **WSL2**
- Experimentation (optional): **Girus**
- Package management / virtual environments for the repo itself → **Pixi**
- Version control and collaboration for the repo itself → **Git** + **GitHub**
- CI/CD for the FreeLunch repo and templates → **GitHub Actions**
- Build orchestration → **Task**
- Static checks, formatting, and pre-commit enforcement → **pre-commit hooks**
- Testing for the repo and platform engine → **Go tests**, with **Testcontainers**/**WireMock** for integration-style coverage
- Primary language for the core engine / CLI → **Go**
- Primary language for the IDE frontend → **TypeScript**
- Versioned Go package publishing for the FreeLunch CLI/core engine → **GitHub Releases** + **Go module proxy**
- Versioned monorepo template publishing for the FreeLunch CLI bootstrap flow → **GitHub Releases** / **GitHub repository template**
- Documentation site for the repo and product narrative → **MkDocs**

## 2. Stack for the product (the devops platform demo)

This covers the runtime, delivery, and operations stack used to demonstrate the FreeLunch platform to customers.

- Programming Languages Supported: **Any**
- Local cloud-service emulation of AWS → **Floci**
- Compute Clusters → **AWS EC2** running **Talos Linux**
- Cloud Workloads Runtime → **Kubernetes** with **Cilium Mesh/Gateway**
- Container registry for built images → **AWS ECR**
- GitOps deployment engine for L2 artifacts → **ArgoCD**
- Kubernetes manifest packaging for deployed workloads → **Helm**
- Platform Infra Provisioning → **Terraform**
- Application Infra Provisioning -> **Crossplane**
- Kubernetes Installation → **Talos Linux**
- Identity provider for human authentication in the IDE → **Keycloak**
- Secrets storage for application credentials → **Vault**
- Secret synchronization from Vault into Kubernetes → **external-secrets-operator**
- Abstraction/schema language for platform configuration and L1 definitions → **CUE**
- Monorepo polyglot build system -> **Bazel**
- CI/CD pipeline for customer workloads → **GitHub Actions** + **Dagger**
- Container Image Build from source without Dockerfiles → **Cloud Native Buildpacks**
- Progressive delivery / deployment automation → **Kargo**
- Autoscaling for workloads → **KEDA** + **Karpenter**
- Observability backend for metrics/logs/traces → **SigNoz**
- Distributed tracing / telemetry instrumentation → **OpenTelemetry**
- Kubernetes cluster administration → **Rancher**
- Cost observabilityfor workloads → **OpenCost** + **Prometheus**
- Load testing → **K6**
- Test doubles / API mocking for integration tests → **WireMock**
- Kubernetes integration testing → **kubetest**
- Security scanning for images and artifacts → **Trivy**
- IDE/workbench frameowrk for making the IDE → **Eclipse Theia**
- Plugins -> **OpenVSX** (for **VSCode** extensions) + **Backstage** (for embedded control planes) plugins
