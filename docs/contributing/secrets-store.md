# Running the Local Secrets Store

The local secrets store from roadmap 1.4 — an OpenBao instance in the Demo cluster.
Everything below was executed exactly as written.

> **Verified on:** macOS 26 / Apple Silicon (arm64), Docker 29.4.0, OpenBao 2.6.2.
> **Not yet verified on:** Linux or WSL2 — same status as the other runbooks.

> ⚠️ **This store is not secure, holds nothing real, and loses its contents on every
> restart.** It runs OpenBao in dev mode: in-memory, auto-unsealed, root token
> `freelunch-dev-root` (a constant in the code — not a secret). Do not put anything real
> in it.

## Why OpenBao and not HashiCorp Vault

Decision D1: Vault has been **BUSL-1.1 since August 2023**, which restricts offering it
as an embedded service in a commercial product. OpenBao is the Linux Foundation fork of
the last MPL-2.0 codebase, with the same HTTP API. Two renames to know: the CLI is
**`bao`** and the env vars are **`BAO_ADDR`/`BAO_TOKEN`** — Vault tutorials otherwise
apply. Beware conda-forge's `vault` package: it advertises MPL-2.0 but builds BUSL-1.1
HashiCorp Vault; the metadata is stale.

## Bring it up

```bash
./bin/freelunch install               # cluster + auth + secrets
./bin/freelunch install --only secrets    # store only, cluster already up
```

Install waits for the store (it comes up in seconds, unlike Keycloak) and **seeds the
demo credential** `secret/example_service · api-key` on every run — dev mode loses
contents on restart, so seeding is repeatable by design. `freelunch status` reports the
mounted engine:

```
Secrets store is ready: secret/ (kv v2)
```

"Ready" means exactly that line: the store answers, is unsealed, **and** a KV **v2** engine
is mounted at `secret/`. Dev mode mounts it on every start, so if `status` instead says the
store is `SEALED`, or `up but NOT READY` with a missing or different engine, the deployment
has been changed out from under you and waiting will not fix it — inspect it, or recreate
the environment. A store that is still coming up says `not ready` with no detail.

## Using it

There is deliberately **no Ingress** — the consumer is external-secrets-operator (2.1),
which runs in-cluster. From the host, go through the pod:

```bash
export PATH="$HOME/.freelunch/bin:$PATH"
POD=$(kubectl -n freelunch-system get pod -l app.kubernetes.io/name=openbao \
  --field-selector=status.phase=Running -o jsonpath='{.items[0].metadata.name}')

# write / read with the in-pod CLI
kubectl -n freelunch-system exec $POD -- env BAO_ADDR=http://127.0.0.1:8200 \
  BAO_TOKEN=freelunch-dev-root bao kv put secret/example_service api-key=new-value
kubectl -n freelunch-system exec $POD -- env BAO_ADDR=http://127.0.0.1:8200 \
  BAO_TOKEN=freelunch-dev-root bao kv get secret/example_service
```

## ⚠️ The KV v2 path rewrite — read this before wiring 2.1

The CLI addresses the secret as `secret/example_service`. The **HTTP API serves it at
`secret/data/example_service`** — KV v2 inserts `data/` after the mount. The API is what
external-secrets-operator uses, and a `SecretStore` configured with the CLI-shaped path
finds nothing and reports an *empty secret, not an error*.

This is the `SecretStore` shape 2.1 should start from (token auth against the dev store;
Kubernetes auth was **deliberately deferred to 2.1** rather than half-configured here):

```yaml
apiVersion: external-secrets.io/v1
kind: SecretStore
metadata:
  name: freelunch-openbao
  namespace: freelunch-system
spec:
  provider:
    vault:                                  # the vault provider speaks OpenBao's API
      server: http://openbao.freelunch-system.svc:8200
      path: secret                          # the mount — ESO adds data/ itself for v2
      version: v2                           # THE load-bearing line
      auth:
        tokenSecretRef:
          name: openbao-token               # a K8s Secret holding freelunch-dev-root
          key: token
```

With `version: v2`, ESO handles the `data/` insertion; an `ExternalSecret` then refers to
`key: example_service` and property `api-key`. Getting `version` wrong reproduces the
silent-empty failure described above.

## What 1.4 does not deliver — do not file it as a bug

The roadmap story ends with a pod getting an `API_KEY` env var. **That last leg is 2.1's**
(external-secrets-operator): 1.4 builds the store; 2.1 builds the bridge. Until 2.1, no
pod receives anything — the demonstrable claim is that the secret is stored, seeded on
install, and readable at the documented API path.

## Troubleshooting

**`status` says not ready.** `kubectl -n freelunch-system get pods`, then
`kubectl -n freelunch-system logs deploy/openbao | tail`.

**`status` says SEALED.** Dev mode never seals. Someone changed the deployment args away
from `-dev`; re-run `freelunch install --only secrets` to restore the committed shape.

**Wrote a secret, API returns nothing.** You read the CLI path. Insert `data/` after the
mount: `/v1/secret/data/<path>`, not `/v1/secret/<path>`.

**Everything vanished.** The pod restarted; that is the dev-mode contract. Re-run
`freelunch install --only secrets` to re-seed.
