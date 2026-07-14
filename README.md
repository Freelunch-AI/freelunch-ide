# Freelunch Founding Document

## Vision & Business Model

Freelunch is the next-gen venture-studio that goes beyond spawning startups.

We partner with founders and provide scaleup-infra-as-a-service for our scaleups, via our ai-native platform stack that allows small teams to seamlessly scale with low head count, after they have already reached bootstrapped Product-Market Fit (PMF). We create and support startups building apps (excludes dev tool companies).

The core platform we provide is the devops/mlops platform, but we also will provide ERP (adapted for scaleup needs), internal tech upskilling platform and offices in the future, so that founders can focus all their efforts on their specific business needs. All these platforms are free, ejectable (leave with artifacts + gitflow) and forkable for our startups (but not publicly forkable) creating a strong inner source (but not open source) ecosystem & community.

**Business Model**: We take one of these 3 approaches: (1) bring in founders after validating the idea, giving them 40% equity + salary; (2) bring in pre-scale companies via y combinator-like form application, taking 20% equity of them; (3) bring in founders for an edge city-like idea exploration phase (or with an idea already), validating the idea together with the founder and starting the company together with the founder taking 30% equity. In all cases we help with eventual fundraising through our VC network. Since we operate essentially as a cofounder, we also get diluted as much as the founders in funding rounds. 

Similar Orgs: 
- [Result](https://www.ycombinator.com/companies/result): from one of Y Combinator’s recent batches. But their focus is small-scale businesses, we focus on companies scaling to infinity. They also aren't focused on the product side that much (the actual devops/mlops required for building scalable/safe/fast/cost-efficient apps beyond Demos), instead on the ERP side (finance/marketing/hr/etc). Another important difference is the business model: they are a PaaS, while we are a venture-studio.
- [Shiva](https://www.omshiva.ai/): a new brazilian vc/accelerator that is backing small hyper productive teams with salaries + AI/Cloud credits instead of large investments.

## Strategy

First, we need to focus on the core devops platform and validate that it actually makes company scaling easy. We need to build the Demo (freelunch devops platform), migrate an existing SaaS app to it and artificially scale it. We also need to compare it to existing open source options and popular PaaS vendors (in terms of devex, capabilities, lock-in & cost).

### Demo Goal
1. Docs: A developer has access to documentation explaining how an example app  was migrated from a pre-existing k8s setup, deployed and operated with freelunch. 
2. Building: The same developer then creates his hello world app: creates two services visually, writes source code inside one of them and imports a container for the other, connects them, adds an external Postgres dependency, telemetry configs, clicks Deploy, and the platform generates standard Terraform, Kubernetes manifests, CI/CD statuses & logs, with observability, secrets, API Gateway, networking all built-in.
3. Scaling: The developer increases inbound traffic artificially to emulate real world product adoption. The infra & services scale automatically to match the increased traffic in a cost-efficient manner.
4. Incident Resolution: Later, when production fails, he inspect traces, topology, logs, deployments, and costs in the same workspace, with the help of Claude Code that has access to the complete picture through the unified agent API. He then does a rollback to the previous prod commit.
5. Ejection: Finally, the developer ejects form the platform maintaining his artifacts and gitflow.

## Demo (showing the idea, but not an MVP yet)

While traditional PaaS' like Heroku/Railway/Render are nice for 0 -> 1, startups move away from them when scaling (because of high cost, strong limitations and high lock-in). These startups suffer with the complexities of Terraform/K8s, along with all the overwhelming ecosystem of cloud-native tools around them, just to get simple things done; safely, fast and cost-efficient. This distracts them from the actual product they are making, requires expensive DevOps/MLOps/Data hires and work is still duplicated across teams. Freelunch IDE is a new kind of IDE for cloud development (e.g., services and source code are first-class citizens and you build inside composable building blocks of application code and infra) and operations, effectively working as an internal developer platform. It streamlines the entire SDLC (dev/experimentation –> ci –> staging –> prod –> observability), with AI assistance in every stage (from coding, to end-to-end testing, to incident solving). But we don’t develop our own coding agent or LLM (at first), we transform coding agents (e.g., Claude Code or Open Code) into complete Software/AI engineering agents.

Freelunch IDE adds developer-friendly abstraction & visual layer on top of the complex Cloud/DevOps/K8s/Data/AI ecosystems, while still letting you deal directly with the underlying tooling (e.g., use kubectl or bring your own cluster) when needed via our 2-layer API model.
A good analogy is with game engines (e.g., Unity), which are specialized & visual IDEs which streamline game development. In this case, it's for general cloud/k8s development and operations.
The most similar open source projects out there are currently Kubero (easy source-code-to-k8s deploy), Kubefirst (modern k8s gitops template) Meshery (visual cloud/k8s/services management & observability) & Backstage (central touching point for developers to have everything they need).

[Freelunch Mock](https://github.com/Freelunch-AI/freelunch-ide/blob/main/docs/freelunch_ide_mock.html)

Demo Innovative Features: (1) K8s-based, with layer 1 API (platform abstractions) that gets compiled to layer 2 (k8s artifacts) which triggers GitOps Deploy Flow; (2) Visual Cloud-native IDE (Backwards-compatible with VSCode) where you build scalable apps by writing, composing & visually debugging building blocks working together; (3) API for Coding Agents to query everything outside of code (think of it as the Dev Portal, but for agents) statuses, state, errors, costs, infra observability data, app observability data, tickets; and create tickets & send notifications to human engineers/devs. (We will be implementing Observality for the customers only in the demo/mvp).


### Demo Expected Features 

- Infra: (1) Blue-green deploys/rollbacks (with db sync & scaling level sync); (5) Auth: CI/CD uses Github as OIDC Provider, IDE uses Keycloak as OIDC Provider, Pods use K8s as OIDC Provider and Vault as Application Secrets Store; (5) IaC knowledge not necessary for getting services running and observing them; (6) Only superficial K8s knowledge necessary; (7) Autoscaling: pod, node vertical and horizontal autoscaling; (7) Local Dev Environment for developing/validating services powered by K8s-in-Docker (Kind); (8) Backups & Restores already set up; (9) Support for using existing EKS clusters.
Interfaces: (1) IDE + Dev Portal being the same thing; (2) Minimal CLI for set up and inspection; (3) GitOps for actual modifications to the systems (unit tests -> integration tests (ephemeral) -> staging (ephemeral) -> prod); (4) Personas: platform admin, platform engineer, developer, tech lead (developer that can merge PRs)
Docs: (1) comes with an example app deployed and documentation explaining how it was migrated from a pre-existing k8s setup
ch.
- Observability: (1) Infra & App Observability; (2) Cost observability
- Application: (1) Easily build and deploy stateless services without dealing necessarily with containers (powered by cloud-native buildpacks)
- Lock-in: (1) option to eject and maintain artifacts + gitops flow.
- Migration to Freelunch: (1) blue-green stateless migration
- Extensions: any Open VSX extension should work
- Language support: any language since its based on containers (later, freelunch’s distributed programming framework will require language-specific work, but this is not for the Demo)

Demo Limitations: (1) fully local & aws-only (local aws cloud emulation); (2) no support for hosting stateful services; (3) no gpu, a/b testing, frontend, data engineering or mlops support; (4) not tool-agnostic yet (e.g., relies on terraform instead of allowing any IaC tool); (5) no distributed programming framework yet; (6) no visual slow-motion replay of traces & time-travel yet; (7) No Project Management & AI Agent Management yet; (8) no polyrepo support yet; (9) no on-premise cluster neither embedded device support yet; (10) no emphasis on auth & security; (11) no remote k8s development/experimentation environment support; (12) Public/Private Hub for reusable blocks; (13) No support for DAGs; (14) No native AI-assisted Engineering; (15) No budget-enforcement & deployment cost gates; (16) No IDE action logging.

What's missing to become an MVP? Support for hosting stateful services, actual deployment to the cloud, a/b testing, proper auth, proper security. And of course, validation with a real-world scaleup using it.

**Demo Estimated Stack**: theia + typescript, golang, git + pre-commit, Github + Act + Dagger + Bazel, Docker, Trivy, AWS ECR, Kargo, CUE, ArgoCD, Terraform, AWS EKS, Helm, KEDA, Karpenter, Backstage, Keycloak, external-secrets-operator, Vault, K6, testcontainers, kubetest, wiremock, Open Telemetry, headlamp, SigNoz, OpenCost, MKDocs, Cloudflare.

**Our own Coding/Experimenting Environment Estimated Stack**: linux/wsl, pixi, git, girus

Basic setup for us to start development: github repo & access control, virtual environment & package management; linting/formatting; building; testing; publishing IDE binary and CLI Golang package; updating docs website.

Monetization while building our product: platform consulting. Being hired to be the outsourced platform engineering team of existing scaleups that don’t have a team for this. Hired to build a custom platform for the company (caveat: will have some specific needs that don't generalize to other companies).


## Key Technical Decisions & FAQ

- IDE vs Platform Separate from IDE: building a developer platform separate from the IDE is easier and makes adoption easier, however, it can’t provide the seamless developer experience as a single integrated environment where you code, debug and observe your systems in the same place, with holistic AI Assistance. This caps the potential of becoming a groundbreaking tool that changes the way in which teams develop cloud software and scale companies.
- Theia vs VScode Fork: Use Theia as a replaceable workbench shell, ReactFlow as the center of the user experience, Monaco/xterm/LSP/DAP as the building blocks of service-centric mini IDEs, steal catalog/template/auth/plugin concepts from Backstage, steal remote-development patterns from Gitpod/Che/Coder/Tilt, and build your own Kubernetes-native application model instead of inheriting someone else's workspace or portal model.
We did not choose to fork VS Code because of Freelunch's long-term center of gravity is a distributed-systems runtime workbench (services, agents, deployments, observability, debugging, topology graphs) rather than a code editor, and a VS Code fork would increasingly force the product to conform to editor-centric assumptions instead of allowing the runtime model to become the primary user experience.
- Primary Programming Language: We chose Golang as our primary language because of: (1) its cloud-native popularity and use in similar projects, (2) its good performance; (3) experience in our team + small learning curve compared to other languages like Rust for example. However, Typescript will be used for the Theia Frontend.
- Why Venture Studio and not PaaS Business Model? Because PaaS inherently has significant lock-in for the customers and doesn’t tend to become unicorn-level businesses. Though we acknowledge that starting a venture studio is inherently more complex, because 2 things need to work to prove its worth: the freelunch platform and one portfolio company.
Stateless only vs Stateless & Statefull: we chose to support stateless workload only because it will make the Demo much simpler and will be sufficient to show the core idea of the IDE. But of course, it won't be very useful for scaleups yet, only after we implement support for stateful workloads.
- Monorepo vs Polyrepo: the goal is to support both, but for the Demo we want to keep it as simple as possible. Therefore, supporting just monorepo makes more sense since it’s a simpler option that is also being highly adopted by platform engineering/devops teams around the world.
- Multi-cloud vs just AWS: the goal is to be multi-cloud eventually, but for the Demo we want to keep it as simple as possible, as just AWS already does the job of validating things for us.
- OSS vs AWS Services (e.g., Vault vs AWS Secrets Manager): the goal now is to just use AWS as a compute provider, as use OSS tools for the rest. Because we don’t want the MVP to be locked-in AWS and we don’t want unnecessary costs. A big part of our pitch is avoiding the lock-in that PaaS platforms like Railway have and being cost-efficient. But we should eventually also support users picking AWS services instead (because there are companies which have AWS credits they need to use)



