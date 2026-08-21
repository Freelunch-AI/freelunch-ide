# Running the Local Auth Service

The local OIDC identity provider from roadmap 1.3 — a Keycloak instance in the Demo
cluster. Everything below was executed exactly as written.

> **Verified on:** macOS 26 / Apple Silicon (arm64), Docker 29.4.0, Keycloak 26.7.2.
> **Not yet verified on:** Linux or WSL2 — same status as the cluster runbook.

> ⚠️ **This instance is not secure, and holds nothing real.** It runs Keycloak in dev
> mode: in-memory database, known admin password, HTTP without TLS. Every credential in
> it is a demo credential. Do not put anything real in it.

---

## Bring it up

The auth service is part of the default install:

```bash
./bin/freelunch install            # cluster + auth service
./bin/freelunch install --only auth    # just the auth service, cluster already up
./bin/freelunch install --skip auth    # cluster only
```

Keycloak takes **about 60 seconds** to come up after install returns. `freelunch status`
is honest about this — it probes the OIDC discovery endpoint over HTTP rather than
trusting pod readiness, and reports "not ready" until discovery actually answers:

```
Auth service is ready: realm "freelunch", issuer http://keycloak.localhost:8080/realms/freelunch
```

The issuer being printed is deliberate: it is the value that breaks first when hostname
configuration is wrong, and seeing it here beats diagnosing a mysterious 401 later.

## The console, and the demo users

Admin console: **http://keycloak.localhost:8080/admin/** — user `admin`, password
`admin`. (`*.localhost` resolves to loopback on macOS and Linux; no `/etc/hosts` entry
needed.)

The `freelunch` realm is imported at startup from
`src/cli/internal/cluster/../auth/freelunch-realm.json`, embedded in the binary. It
defines:

| | | password |
|---|---|---|
| `carol` | Platform Admin | `demo` |
| `bob` | Platform Engineer | `demo` |
| `alice` | Developer | `demo` |

Groups: `platform-admin`, `platform-engineer`, `developer` (the three Personas), plus
`developer-tech-lead` and `platform-tech-lead` (roadmap 2.4's *temporary grants* — not
Personas) and a `hotfix` realm role (a *permission*).

Two clients:

- **`freelunch-ide`** — public, Authorization Code + PKCE, for the Theia IDE (Group 7).
  Its redirect URIs are placeholders until the IDE exists.
- **`freelunch-agent`** — confidential, client credentials, for the Agent Integration
  Layer (roadmap 8.1). Secret: `freelunch-agent-demo-secret` (not a secret, see above).

Try 8.1's exact authentication path:

```bash
curl -s \
  -d client_id=freelunch-agent \
  -d client_secret=freelunch-agent-demo-secret \
  -d grant_type=client_credentials \
  http://keycloak.localhost:8080/realms/freelunch/protocol/openid-connect/token
```

## The one rule: the realm JSON is the source of truth

The database is in-memory. **Anything you create by hand in the console is lost on the
next pod restart** — and the realm is re-imported from the committed JSON on every boot,
which is exactly what makes the committed state durable. To change the realm (a user, a
client, a redirect URI), edit `freelunch-realm.json` and run
`freelunch install --only auth`; it re-applies the ConfigMap and restarts the pod.

## Troubleshooting

**`status` says not ready long after 60s.** Look at the pod:
`kubectl -n freelunch-system get pods`, then
`kubectl -n freelunch-system logs deploy/keycloak | tail -20`. A realm JSON syntax error
shows up here as an import failure at startup.

**Redirects or token validation fail with the right credentials.** Check the issuer:
`curl -s http://keycloak.localhost:8080/realms/freelunch/.well-known/openid-configuration | grep issuer`.
If it is not `http://keycloak.localhost:8080/realms/freelunch`, the `KC_HOSTNAME`
configuration in `keycloak.yaml` has drifted — Keycloak bakes the hostname into every
endpoint URL it advertises, and clients follow those, not the URL you typed.

**503 for a few seconds after a pod restart.** `kubectl rollout status` returns when the
new pod is ready, but Traefik switches endpoints a few seconds later. This window is why
`freelunch status` reads discovery instead of pod state.

## What this does not cover

1.3 registers clients and stops. IDE login is Group 7; the permission model that makes
the Personas mean something is 2.4; machine-to-machine consumers are Group 8. The realm's
groups exist so those can build on them, not because anything enforces them today.
