# Freelunch Founding Document

For the optional cross-client development workflow, see [FreeLunch AI Workflow Bundle](docs/ai_workflow.md).

## Vision & Business Model

Freelunch is the next-gen venture-studio that goes beyond spawning startups.

We partner with founders and provide scaleup-infra-as-a-service for our scaleups, via our ai-native platform stack that allows small teams to seamlessly scale with low head count, after they have already reached bootstrapped Product-Market Fit (PMF). We create and support startups building apps (excludes dev tool companies).
The core platform we provide is the devops/mlops platform, but we also will provide ERP (adapted for scaleup needs), internal tech upskilling platform and offices in the future, so that founders can focus all their efforts on their specific business needs. All these platforms are free and designed for portability: customers retain the standard artifacts and GitOps flow. Dedicated detach or eject automation is post-MVP. The platforms are forkable for our startups (but not publicly forkable), creating a strong inner source (but not open source) ecosystem & community.

**Business Model**: 
- Intermediate Business Model: Developer Platform Consulting or Selling Monthly/Annual Licences of our IDE/Platform.
- Final Business Model: Venture Studio/Accelerator. We take one of these 3 approaches: (1) bring in founders after validating the idea, giving them 40% equity + salary; (2) bring in pre-scale companies via y combinator-like form application, taking 20% equity of them; (3) bring in founders for an edge city-like idea exploration phase (or with an idea already), validating the idea together with the founder and starting the company together with the founder taking 30% equity. In all cases we help with eventual fundraising through our VC network. Since we operate essentially as a cofounder, we also get diluted as much as the founders in funding rounds. 

Similar Orgs: 
- [Result](https://www.ycombinator.com/companies/result): from one of Y Combinator’s recent batches. But their focus is small-scale businesses, we focus on companies scaling to infinity. They also aren't focused on the product side that much (the actual devops/mlops required for building scalable/safe/fast/cost-efficient apps beyond Demos), instead on the ERP side (finance/marketing/hr/etc). Another important difference is the business model: they are a PaaS, while we are a venture-studio.
- [Wildlife Studios](https://wildlifestudios.com/): Venture Studio focused on mobile games with mature platform engineering.
- [Shiva](https://www.omshiva.ai/): a new brazilian vc/accelerator that is backing small hyper productive teams with salaries + AI/Cloud credits instead of large investments.

## Vision of the Core Product via Analogies with existing Tools

Analogy: Freelunch: platform = \
for_distributed_apps_and_ai_powered(Unreal Engine: visual & all-in-one IDE) +

 for_distributed_apps_and_all_in_one(Cursor: ai-powered IDE) + 
 
 lunch_platform_native(Paperclip: agent orchestration plane) +
 
 just_similar_visuals_and_marketplace(N8N: node-based IDE) +
 
 just_similar_visuals(Meshery: mesh-agnostic service & infra topology design & observability) +
 
multi_language_and_build_time_and_k8s_native(Spring Boot: microservices framework) +

 multi_language_and_build_time_and_made_for_services_on_top_of_k8s(Ray: distributed programming framework) +
 
developer_friendly_and_with_finops(Kubefirst: gitops k8s developer platform template) +

 Karmada: single-cluster abstraction for multi-cluster
 
 build_time_platform_with_infra_support_and_with_environment_progression_and_layer_conflict_resolution(Kubero: k8s developer platform)
 
 inferred_from_code(Infisical: secrets & config management) +
 
 Rancher: k8s cluster management +
 
 Terraform/Crossplane: IaC +
 
 Talscale: Modern VPNs +
 
 decoupled_from_programming_language(ZenML: tool-agnostic DAGs) + integrated_within_devops(Langfuse: LM Systems Observability) +
 
 Marimo: structured notebooks +

 SkyPilot: ephemeral workloads layer +
 
 integrated_within_devops(MLFlow: Artifact Experiment Tracking) +

 declarative_with_visuals_write_audit_publish_pattern_and_data_transformation_management(Ibis tool-agnostic data transformation) +
 
 Beam: cross-storage data transformation engine + 
 
 with_copy_on_write_and_anonymization_support(Kubeblocks: k8s structured data storage management) +
 
 Rook: k8s raw data storage management +
 
 ready_made_and_webmcp_powered_portal(Backstage: unified dashbaords & control planes in a single view) +

## Strategy

First, we need to focus on the core devops platform and validate that it actually makes company scaling easy. We need to build the Demo (freelunch devops platform), deploy a representative stateless app, and artificially scale it. We also need to compare it to existing open source options and popular PaaS vendors (in terms of devex, capabilities, lock-in & cost).

### Demo Goal
1. Docs: A developer has access to documentation explaining how the stateless sample app is modeled, deployed, and operated with Freelunch.
2. Building: The same developer creates multiple services on the canvas: writes source code for one, imports a container image for another, connects them, and adds a virtual Service block for an externally managed database dependency. Cloud Native Buildpacks turn the source into an image, while the canvas-maintained CUE captures configuration, metadata, and image references that compile into versioned Helm charts. The workspace surfaces CI/CD status and logs with observability, secrets, API Gateway, and networking built in.
3. Scaling: The developer increases inbound traffic artificially to emulate real world product adoption. The infra & services scale automatically to match the increased traffic in a cost-efficient manner.
4. Incident Resolution: Later, when production fails, the developer inspects traces, topology, logs, deployments, and costs in the same workspace. Claude Code uses the documentation-backed FreeLunch skill and read-only Coding Agent API to correlate the evidence, and the developer rolls back to the retained previous revision through Argo Rollouts.

## Demo (showing the idea, but not an MVP yet)

While traditional PaaS' like Heroku/Railway/Render are nice for 0 -> 1, startups move away from them when scaling (because of high cost, strong limitations and high lock-in). These startups suffer with the complexities of Terraform/K8s, along with all the overwhelming ecosystem of cloud-native tools around them, just to get simple things done; safely, fast and cost-efficient. This distracts them from the actual product they are making, requires expensive DevOps/MLOps/Data hires and work is still duplicated across teams. Freelunch IDE is a new kind of IDE for cloud development (e.g., services and source code are first-class citizens and you build inside composable building blocks of application code and infra) and operations, effectively working as an internal developer platform. It streamlines the entire SDLC (dev/experimentation –> ci –> staging –> prod –> observability), with AI assistance in every stage (from coding, to end-to-end testing, to incident solving). But we don’t develop our own coding agent or LLM (at first), we transform coding agents (e.g., Claude Code or Open Code) into complete Software/AI engineering agents.

Freelunch IDE adds developer-friendly abstraction & visual layer on top of the complex Cloud/DevOps/K8s/Data/AI ecosystems, while still letting you deal directly with the underlying tooling (e.g., use kubectl or bring your own cluster) when needed via our 2-layer API model.
A good analogy is with game engines (e.g., Unity), which are specialized & visual IDEs which streamline game development. In this case, it's for general cloud/k8s development and operations.
The most similar open source projects out there are currently Kubero (easy source-code-to-k8s deploy), Kubefirst (modern k8s gitops template) Meshery (visual cloud/k8s/services management & observability) & Backstage (central touching point for developers to have everything they need).

[Freelunch Mock](https://github.com/Freelunch-AI/freelunch-ide/blob/main/docs/freelunch_ide_mock.html)

The mock is an exploratory interface prototype. The [ordered feature specification](docs/freelunch_ide_features_ordered.md) is the source of truth for Demo/MVP implementation scope.

Demo Innovative Features: (1) K8s-based, with a canvas-maintained layer 1 API (platform abstractions) that deterministically compiles configuration, metadata, and Buildpacks-produced image references into versioned Helm charts and triggers the GitOps deployment flow; (2) a visual cloud-native IDE (backwards-compatible with VS Code) where developers compose and debug scalable building blocks; (3) a read-only Coding Agent API with an OpenAPI contract and a documentation-backed first-party FreeLunch skill for interactive and headless platform workflows. Agent-triggered mutations, ticket creation, and notifications are post-MVP.


### Demo Expected Features 

- Infra: (1) Argo Rollouts-managed blue-green deploys and rollbacks for stateless applications; (2) Auth: CI/CD uses GitHub as OIDC Provider, IDE uses Keycloak as OIDC Provider, Pods use K8s as OIDC Provider and Vault as Application Secrets Store; (3) IaC knowledge not necessary for getting services running and observing them; (4) Only superficial K8s knowledge necessary; (5) Autoscaling: pod, node vertical and horizontal autoscaling; (6) Local Dev Environment for developing/validating services powered by K8s-in-Docker (Kind); (7) Backups & Restores already set up; (8) Support for using existing EKS clusters.
Interfaces: (1) IDE + Dev Portal being the same thing; (2) Minimal CLI for set up and inspection; (3) GitOps for actual modifications to the systems (unit tests -> integration tests (ephemeral) -> staging (ephemeral) -> prod); (4) Personas: platform admin, platform engineer, developer, tech lead (developer that can merge PRs)
Docs: (1) comes with a stateless sample app and documentation explaining how it is modeled, deployed, and operated
- Observability: (1) Infra & App Observability through SigNoz; (2) OpenCost with its Prometheus backend for Demo cost observability; (3) SigNoz and OpenCost UIs embedded in FreeLunch as plugins, with a unified FreeLunch-owned experience planned later
- Application: (1) Easily build and deploy stateless services without dealing necessarily with containers (powered by cloud-native buildpacks)
- Platform lifecycle: (1) Declarative platform versioning with compatibility checks and guided resolution for breaking schema changes
- Lock-in: (1) standard, customer-owned L2 artifacts and GitOps flow remain inspectable and directly operable; automated detach/eject workflows are post-MVP
- Extensions: any Open VSX extension should work
- Language support: any language since its based on containers (later, freelunch’s distributed programming framework will require language-specific work, but this is not for the Demo)

**Demo Limitations**: (1) fully local & aws-only (local aws cloud emulation); (2) no support for hosting stateful services or orchestrating database/queue migration; (3) no gpu, a/b testing, frontend, data engineering or mlops support; (4) not tool-agnostic yet (e.g., relies on terraform instead of allowing any IaC tool); (5) no distributed programming framework yet; (6) no visual slow-motion replay of traces & time-travel yet; (7) no Project Management & AI Agent Management yet; (8) no polyrepo support yet; (9) no on-premise cluster neither embedded device support yet; (10) no emphasis on auth & security; (11) no remote k8s development/experimentation environment support; (12) no Public/Private Hub for reusable blocks; (13) no support for DAGs; (14) no detach/eject automation or agent-triggered platform mutations; (15) no budget enforcement or deployment cost gates; (16) no IDE action logging; (17) no DORA metrics or IDE usage analytics; (18) no system-wide experiment tracking; (19) no widgets for performing monitoring actions; (20) no DataOps promotion of data into production storage systems through GitOps; (21) no lineage-full observability through OpenLineage.
What's missing to become an MVP? Support for hosting stateful services, actual deployment to the cloud, a/b testing, proper auth, proper security. And of course, validation with a real-world scaleup using it.

**Demo Estimated Stack**: theia + typescript, golang, git + pre-commit, Github + Act + Dagger + Bazel, Docker, Trivy, AWS ECR, CUE, ArgoCD, Argo Rollouts, Terraform, AWS EKS, Helm, KEDA, Karpenter, Backstage, Keycloak, external-secrets-operator, Vault, K6, testcontainers, kubetest, wiremock, Open Telemetry, headlamp, SigNoz, OpenCost, Prometheus, MKDocs, Cloudflare.
**Our own Coding/Experimenting Environment Estimated Stack**: linux/wsl, pixi, git, girus

Basic setup for us to start development: github repo & access control, virtual environment & package management; linting/formatting; building; testing; publishing IDE binary and CLI Golang package; updating docs website.

Monetization while building our product: platform consulting. Being hired to be the outsourced platform engineering team of existing scaleups that don’t have a team for this. Hired to build a custom platform for the company (caveat: will have some specific needs that don't generalize to other companies).


## Key Technical Decisions & FAQ

- IDE vs Platform Separate from IDE: building a developer platform separate from the IDE is easier and makes adoption easier, however, it can’t provide the seamless developer experience as a single integrated environment where you code, debug and observe your systems in the same place, with holistic AI Assistance. This caps the potential of becoming a groundbreaking tool that changes the way in which teams develop cloud software and scale companies.
- Theia vs VScode Fork: Use Theia as a replaceable workbench shell, ReactFlow as the center of the user experience, Monaco/xterm/LSP/DAP as the building blocks of service-centric mini IDEs, steal catalog/template/auth/plugin concepts from Backstage, steal remote-development patterns from Gitpod/Che/Coder/Tilt, and build your own Kubernetes-native application model instead of inheriting someone else's workspace or portal model.
We did not choose to fork VS Code because of Freelunch's long-term center of gravity is a distributed-systems runtime workbench (services, agents, deployments, observability, debugging, topology graphs) rather than a code editor, and a VS Code fork would increasingly force the product to conform to editor-centric assumptions instead of allowing the runtime model to become the primary user experience.
- Primary Programming Language: We chose Golang as our primary language because of: (1) its cloud-native popularity and use in similar projects, (2) its good performance; (3) experience in our team + small learning curve compared to other languages like Rust for example. However, Typescript will be used for the Theia Frontend.
- Why Venture Studio and not PaaS Business Model? Because PaaS inherently has significant lock-in for the customers and doesn’t tend to become unicorn-level businesses. Though we acknowledge that starting a venture studio is inherently more complex, because 2 things need to work to prove its worth: the freelunch platform and one portfolio company.
- Stateless only vs Stateless & Statefull: we chose to support stateless workload only because it will make the Demo much simpler and will be sufficient to show the core idea of the IDE. But of course, it won't be very useful for scaleups yet, only after we implement support for stateful workloads.
- Monorepo vs Polyrepo for the User's gitops: the goal is to support both, but for the Demo we want to keep it as simple as possible. Therefore, supporting just monorepo makes more sense since it’s a simpler option that is also being highly adopted by platform engineering/devops teams around the world.
- Fully local vs Deployment to the Cloud: we choose the deo to be fully local to speed up iteration speed for the demo, and also avoid cost risks. Though we recognize that being local is just enough for a Demo not an MVP that will require cloud deployment.
