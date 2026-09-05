# Running the Local Demo Cluster

The local Kubernetes cluster from roadmap 1.2. Everything below was executed exactly as
written — this is a transcript, not a design sketch.

> **Verified on:** macOS 26 / Apple Silicon (arm64), Docker 29.4.0.
> **Not yet verified on:** Linux or WSL2. The commands are OS-independent and the pinned
> tools cover `linux-64`, `linux-aarch64` and `osx-arm64`, but nobody has run this on Linux
> yet. If you are the first, please correct anything that is wrong here rather than working
> around it.

---

## What you need first

**Docker must be installed and running.** That is the only prerequisite the environment
does not provide for you — roadmap 1.2 says "a fresh machine with Docker installed", and
everything else is provisioned by the repository.

```bash
docker version --format '{{.Server.Version}}'
```

If that errors, start Docker Desktop (macOS/Windows) or `sudo systemctl start docker`
(Linux) before continuing.

Then, from a clean checkout:

```bash
pixi install       # Go, Node, linters, helm
pixi run setup     # Go modules, npm packages, and the cluster tools
```

`pixi run setup` installs **k3d** and **kubectl** into `~/.freelunch/bin`, at versions
pinned in `src/cli/internal/toolchain/scripts/versions.env`. They are deliberately **not**
added to your `PATH` and **not** installed system-wide, so they cannot collide with a k3d or
kubectl you already use for something else.

To use them directly:

```bash
export PATH="$HOME/.freelunch/bin:$PATH"
```

The `pixi run task cluster:*` commands below address them by full path, so you do not need
to do this unless you want to run `kubectl` yourself.

---

## Bring the cluster up

```bash
pixi run task cluster:up
```

This creates the cluster described by `src/cli/internal/cluster/k3d-cluster.yaml` and then
prints the nodes. Expect roughly 40 seconds on first run, less afterwards. You should see:

```
NAME                     STATUS   ROLES           AGE   VERSION
k3d-freelunch-agent-0    Ready    <none>          10s   v1.35.5+k3s1
k3d-freelunch-server-0   Ready    control-plane   18s   v1.35.5+k3s1
```

Two nodes, both `Ready`. The Kubernetes version is pinned in the config, so it will not
drift when k3d is upgraded.

The cluster is also written into your default kubeconfig and made the current context, so
plain `kubectl` works immediately. `cluster:down` removes that entry again.

---

## Verify it actually works

Nodes being `Ready` only proves the cluster exists. The property that matters — and the one
that killed the earlier ProxMox/Talos plan — is whether you can reach a service **from your
host machine**. Test it:

```bash
export PATH="$HOME/.freelunch/bin:$PATH"

kubectl create deployment probe --image=traefik/whoami --replicas=2
kubectl expose deployment probe --port=80 --target-port=80

kubectl apply -f - <<'EOF'
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: probe
spec:
  rules:
    - http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: probe
                port:
                  number: 80
EOF
```

**Wait for the ingress controller before curling.** k3s installs Traefik through a Helm job
that takes a further 30–60 seconds after the nodes report `Ready`. Requests sent before it
finishes simply fail, which looks exactly like a broken cluster:

```bash
kubectl -n kube-system rollout status deploy/traefik --timeout=180s
```

Then:

```bash
curl -s http://localhost:8080/ | head -4
```

Expected:

```
Hostname: probe-586f846674-d2bhj
IP: 127.0.0.1
IP: ::1
IP: 10.42.1.4
```

Confirm it is load-balancing across both pods:

```bash
for i in $(seq 1 12); do curl -s http://localhost:8080/ | grep -i '^Hostname:'; done | sort | uniq -c
```

Expected: two hostnames, roughly six requests each.

Clean up the probe:

```bash
kubectl delete ingress probe && kubectl delete svc probe && kubectl delete deployment probe
```

---

## Ports

| Host port | Goes to | Why not the obvious port |
|---|---|---|
| `8080` | ingress HTTP (`:80` in-cluster) | ports below 1024 are privileged on Linux and would require root |
| `8443` | ingress HTTPS (`:443` in-cluster) | same |
| `5050` | the cluster's local image registry | recent macOS binds `5000` for AirPlay Receiver |
| random | the Kubernetes API, via the k3d load balancer | k3d picks it; the kubeconfig records it |

All of them are bound to **`127.0.0.1` only**. Docker's default is every interface, which
would put a Keycloak with committed dev credentials and an unauthenticated, writable
registry on whatever network the laptop is connected to. The bind address is set per port
in `k3d-cluster.yaml`, and the reasoning is in the comment there — do not widen it to make
a demo reachable from another device; that wants an explicit, opt-in mechanism, not a
default.

The registry answers on the host:

```bash
curl -s http://localhost:5050/v2/_catalog
# {"repositories":[]}
```

Push images to `localhost:5050/...` to make them available in-cluster without a public
registry. This is also the basis of the offline requirement in 1.2 — images are pulled once
and served locally afterwards.

---

## Everyday commands

| Command | Does |
|---|---|
| `pixi run task cluster:up` | Create the cluster and show nodes |
| `pixi run task cluster:down` | Delete the cluster |
| `pixi run task cluster:status` | Cluster list plus `kubectl get nodes -o wide` |
| `pixi run task setup:cluster-tools` | Reinstall/repair k3d and kubectl |
| `pixi run task setup:airgap-images` | Cache the k3s image bundle for offline creation (~220MB) |
| `pixi run task pin:tools` | Regenerate `checksums.txt` after a version bump |

The same cluster lifecycle is also on the CLI, which is the path a customer gets — they
have a `freelunch` binary and no pixi, no checkout and no Taskfile:

| Command | Does |
|---|---|
| `freelunch install` | Create the cluster from the config embedded in the binary |
| `freelunch uninstall` | Delete the cluster |
| `freelunch status` | Report whether the cluster exists, whether every node is Ready, and each node's state |

Both drive the same pinned k3d and kubectl in `~/.freelunch/bin`, so they are
interchangeable — use whichever is at hand. The tasks print more, which is usually what
you want while working on the cluster itself; the CLI is what the runbook's audience will
eventually ship against.

The cluster is **disposable by design** — roadmap 1.2 specifies "fresh start only". If
anything is strange, delete and recreate rather than repairing in place:

```bash
pixi run task cluster:down && pixi run task cluster:up
```

---

## Offline / airgapped cluster creation

Roadmap 1.2 requires the environment to come up "without touching the internet". A k3d
cluster does **not** manage that on its own: the `rancher/k3s` image ships no preloaded
workload images, so traefik, coredns, local-path-provisioner, klipper-lb, klipper-helm,
metrics-server, pause and busybox are all pulled from the network as the cluster boots.

Cache them once:

```bash
pixi run task setup:airgap-images
```

That downloads the official `k3s-airgap-images-<arch>.tar.gz` for the pinned k3s release
into `~/.freelunch/images`, verifying it against the digest in `checksums.txt`. It is
about **220 MB**, which is why it is not part of `pixi run setup` — run it when you need
offline capability.

After that, both `freelunch install` and `pixi run task cluster:up` mount that directory
into every node at `/var/lib/rancher/k3s/agent/images`, and k3s imports the bundle at
startup before scheduling anything. Cluster creation takes the same ~25s as it does
online. Without the cache, both fall back to pulling from the network — the CLI logs
which mode it chose at `Debug`.

Use the **official bundle**, never a hand-written image list. Enumerating a running
cluster with `crictl images` returns six images; the release's own `k3s-images.txt` lists
eight. `busybox` and `metrics-server` are pulled later than a fresh inspection catches, so
a hand-rolled list looks complete and fails only once the network is actually gone.

### Verifying an airgap honestly

Two traps make this easy to get wrong.

**`k3d cluster create` reporting success proves very little.** It waits for the k3s server
process, not for workloads. With registries unreachable and no cache, k3d still prints
"Cluster 'freelunch' created successfully!" while the node holds nothing but `pause` and
every component is stuck. Always check what actually runs:

```bash
kubectl -n kube-system rollout status deploy/traefik
kubectl -n kube-system get pods
```

**Do not test with a Docker `--internal` network.** It looks like the obvious way to
remove connectivity, but k3s refuses to start on one:

```
level=fatal msg="Error: no default routes found in /proc/net/route or /proc/net/ipv6_route"
```

k3s needs a default route to select its node IP, so an internal network fails for a reason
that has nothing to do with images — it tells you nothing about your airgap. Blackhole the
registries instead, which leaves the network intact but makes every pull fail:

```yaml
# /tmp/airgap-registries.yaml
mirrors:
  docker.io:
    endpoint: ["http://127.0.0.1:1"]
  registry.k8s.io:
    endpoint: ["http://127.0.0.1:1"]
  ghcr.io:
    endpoint: ["http://127.0.0.1:1"]
  quay.io:
    endpoint: ["http://127.0.0.1:1"]
```

```bash
k3d cluster create --config <config> \
  --volume "$HOME/.freelunch/images:/var/lib/rancher/k3s/agent/images@server:*;agent:*" \
  --registry-config /tmp/airgap-registries.yaml
```

Run it once **without** the `--volume` first. If that still yields a working cluster, your
test has no teeth and the real one is not proving anything either.

One thing the cache does not cover: the three images Docker itself needs on the host —
`rancher/k3s`, `ghcr.io/k3d-io/k3d-proxy` and `ghcr.io/k3d-io/k3d-tools`. Those are pulled
by the Docker daemon, not by k3s, so a truly offline first run needs them already present
in the local Docker image store.

---

## Changing versions

k3d, kubectl and the Kubernetes version are pinned in two files:

- `src/cli/internal/toolchain/scripts/versions.env` — k3d and kubectl
- `src/cli/internal/cluster/k3d-cluster.yaml` — the k3s image

Keep the kubectl and k3s **minor** versions in step; kubectl supports one minor of skew in
either direction. After changing `versions.env`:

```bash
pixi run task pin:tools     # regenerates checksums.txt
```

Commit `versions.env` and `checksums.txt` together. The installers refuse to run against a
stale `checksums.txt` rather than accepting a binary you have not pinned.

---

## Troubleshooting

**`curl: (52) Empty reply` or connection refused on `:8080`**
Traefik is probably still installing. Run
`kubectl -n kube-system rollout status deploy/traefik --timeout=180s` and try again. This is
the single most common false alarm — the cluster looks broken and is merely not finished.

**`Cannot connect to the Docker daemon`**
Docker is not running. Everything here is containers; there is no fallback.

**`port is already allocated` on 8080, 8443 or 5050**
Something else on your machine holds the port. Find it with `lsof -i :8080`, or change the
port number in `k3d-cluster.yaml` — the ports are arbitrary and only need to be free. Keep
the `127.0.0.1:` prefix when you do.

**The cluster is not reachable from another machine**
Intended: every host port is bound to loopback (see [Ports](#ports)).

**`checksums.txt is stale — run 'pixi run task pin:tools'`**
`versions.env` was changed without regenerating the checksums. Run the command it suggests
and commit both files.

**`no pinned checksum for <tool> <version> <os>/<arch>`**
Your platform is not in `checksums.txt`. If it is a platform we intend to support, add it to
`PLATFORMS` in `pin-tools.sh` and to `platforms` in `pixi.toml`, then regenerate.

**k3d or kubectl behaves oddly**
Both are verified by checksum on every run, so a corrupted binary repairs itself:
`pixi run task setup:cluster-tools`.

**`kubectl` talks to the wrong cluster**
`cluster:up` switches your current context. Check with `kubectl config current-context` —
it should be `k3d-freelunch`.

---

## What this does not cover

1.2 ends at a working cluster. ArgoCD, Argo Rollouts and external-secrets are Group 2;
Keycloak is 1.3; Vault is 1.4; SigNoz and OpenCost are Group 6. `freelunch install` will
eventually install all of them; today it stops at the cluster, which is exactly what 1.2
asks for.

`freelunch start` and `freelunch stop` do not exist yet. They are reserved for resuming
and suspending a cluster that already exists — k3d's own `cluster start`/`stop` — rather
than being second names for install and uninstall.
