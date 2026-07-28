# Demo Tech Stack

This is a first-pass mapping of capability to tool choice, derived from the FreeLunch feature order and founding docs.

Notes:

- these are not strict requirements, think of it as an implementaion guide
- we might realize during implementation that some tools are worth more borrowing ideas/code than to actually just use them out-of-the-box
- embedded tool UIs doesnt mean we wont provide our code-centric user experience for it, just means we can provide to the user, easily via plugin, a good UI that is already battle tested, as our baseline.
- Most influential decisions: **Theia**, **ProxMox**, **Talos Linux**, **Go**

## 1. Stack for the FreeLunch repo itself

This covers the tooling used to build, test, document, and ship the FreeLunch package, template, and IDE experience itself.

- Dev OS: **Linux** or **WSL2**
- Experimentation (optional): **Girus** (locked version in pixi)
- Virtual Dev Environment for freelunch development itself & Task runner → **Pixi**
- Version control and collaboration for the repo itself → **Git** + **GitHub**
- CI/CD for the FreeLunch repo and templates → **GitHub Actions**
- Testing for the repo and platform engine → **Go tests**, with **Testcontainers**/**WireMock** (locked version in pixi) for integration-style coverage and **Loft** (locked version in pixi) using **vCluster** for controlled cluster environments
- Primary language for the core engine / CLI → **Go**
- Primary language for the IDE frontend → **TypeScript**
- Versioned Go package publishing for the FreeLunch CLI/core engine → **GitHub Releases** + **Go module proxy**
- Versioned monorepo template publishing for the FreeLunch CLI bootstrap flow → **GitHub Releases** / **GitHub repository template**
- Documentation site for the repo and product narrative → **MkDocs** (locked version in pixi)

## 2. Stack for the product (the devops platform demo)

This covers the runtime, delivery, and operations stack used to demonstrate the FreeLunch platform to customers.

- Programming Languages Supported: **Any**
- Local cloud-service emulation Cloud Compute Cluster → **ProxMox** (locked version in pixi) running **Talos Linux** with **MetalLB**
- Cloud Workloads Runtime → **Kubernetes** (locked CLI version in pixi) with **Cilium Mesh/Gateway**
- Container registry for built images → **Harbor** (locked version in pixi)
- GitOps deployment engine for L2 artifacts → **ArgoCD** (locked version in pixi)
- Kubernetes manifest packaging for deployed workloads → **Helm** (locked version in pixi)
- Infra Provisioning → **Terraform** (locked version in pixi)
- Kubernetes Installation → comes for free with **Talos Linux**
- Identity provider for human authentication in the IDE → **Keycloak** (locked version in pixi)
- Secrets management for application credentials → **Vault** (locked version in pixi)
- Secret synchronization from Vault into Kubernetes → **external-secrets-operator**
- DNS -> **ExternalDNS** with **Cloudflare** as the DNS Provider
- TLS -> **cert-manager**
- Abstraction/schema language for platform configuration and L1 definitions → **CUE** (locked version in pixi)
- CI/CD pipeline for customer workloads → **GitHub Actions** + **Dagger** (embedded UI via plugin)
- Container Image Build from source without Dockerfiles → **Cloud Native Buildpacks** (locked version in pixi)
- Ephemeral Environment Management → **Loft** (embedded UI via plugin) using k8s namespaces 
- Progressive delivery / deployment automation → **Kargo** (locked version in pixi) (embedded UI via plugin)
- Autoscaling for workloads → **VPA** + **KEDA** + **Karpenter**
- Observability backend for metrics/logs/traces → **SigNoz** (locked version in pixi) (embedded UI via plugin)
- Distributed tracing / telemetry instrumentation → **OpenTelemetry**
- Kubernetes observability → **Headlamp** (embedded UI via plugin)
- Cost observabilityfor workloads → **OpenCost** + **Prometheus** (embedded UI via plugin)
- Load testing → **K6** (locked CLI version in pixi)
- Test doubles / API mocking for integration tests → **WireMock**
- Kubernetes integration testing → **kubetest** (locked version in pixi)
- Security scanning for images and artifacts → **Trivy** (locked version in pixi)
- IDE/workbench frameowrk for making the IDE → **Eclipse Theia** (locked version in pixi)
- Plugins -> **OpenVSX** (for **VSCode** extensions) + **Backstage** (for embedded control planes) plugins
