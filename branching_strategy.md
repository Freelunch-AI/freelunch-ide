# Branching Strategy

## Overview

This repository follows a lightweight GitHub Flow with tag-based production releases for the Freelunch repository itself.

### Goals

* Keep `main` always releasable.
* Develop features in short-lived branches.
* Require Pull Requests for all changes.
* Automatically deploy every merge to **Staging**.
* Deploy to **Production** only from a version tag.

---

# Branches

## `main`

* `main` is the single long-lived branch.
* It must always be in a deployable state.
* Direct commits are prohibited.
* All changes must be merged through Pull Requests.

---

## Feature Branches

Create branches from the latest `main`.

Naming conventions:

```
feature/<description>
bugfix/<description>
hotfix/<description>
refactor/<description>
docs/<description>
test/<description>
chore/<description>
```

Examples:

```
feature/oauth-login
bugfix/login-timeout
hotfix/payment-webhook
refactor/cache-layer
```

Feature branches should be short-lived and deleted after merging.

---

# Pull Requests

Every Pull Request must:

* Target `main`
* Pass all CI checks
* Be up to date with `main`
* Receive the required approvals
* Resolve all review comments before merging

---

# Merge Strategy

Use **Squash and Merge** for all Pull Requests.

Benefits:

* One commit per feature
* Clean project history
* Easier rollback
* Simpler changelog generation

---

# Deployment Strategy

## Staging

Every successful merge into `main` automatically:

1. Builds the project
2. Runs the full CI pipeline
3. Deploys the latest commit to the **Staging** environment

This ensures that `main` is continuously validated in a production-like environment.

---

## Production

Production (freelunch package + template) deployments are triggered **only by Git tags**.

Create a release tag:

```bash
git checkout main
git pull
git tag v1.4.0
git push origin v1.4.0
```

Pushing a semantic version tag automatically triggers the Production deployment pipeline.

This provides:

* Explicit production releases
* Immutable version history
* Easy rollback to previous versions
* Clear mapping between deployments and source code

---

# Versioning

Use Semantic Versioning.

Examples:

```
v1.0.0
v1.1.0
v1.1.1
v2.0.0
```

---

# Hotfixes

Critical fixes follow the normal Pull Request process.

```
main
  ↓
hotfix/payment-timeout
  ↓
Pull Request
  ↓
main
  ↓
Automatic Staging deployment
  ↓
Create new version tag
  ↓
Production deployment (freelunch package + template)
```

Hotfixes should **not** be pushed directly to `main`.

---

# Protected Branch Rules

Protect the `main` branch with:

* No direct pushes
* Required Pull Requests
* Required status checks
* Required approvals
* No force pushes
* No branch deletion

---

# Commit Messages

Follow the Conventional Commits specification.

Examples:

```
feat(auth): add OAuth login

fix(api): handle timeout

docs: update installation guide

refactor(cli): simplify parser

test(cache): add Redis integration tests
```

---

# Release Flow

```
feature/*
      │
      ▼
Pull Request
      │
      ▼
Review + CI
      │
      ▼
Squash Merge
      │
      ▼
main
      │
      ▼
Automatic Staging Deployment
      │
      ▼
Create Tag (vX.Y.Z)
      │
      ▼
Automatic Production Deployment (freelunch package + template)
```
