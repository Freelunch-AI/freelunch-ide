# Bruno's AI-assisted Coding Setup

## Doc-specific section markers

**[for Humans or Outer Loop Agents]**: these mark sections of this document which are intended for humans or agents setting up and operating the issue --> PR factory, meaning that who should act on these sections is not doing the actual issue -> PR work, but setting up or operating the process itself that enables the inner loop agents to work properly, effectively and securely inside the process.

**[for Inner Loop Agents]**: these mark sections of this document which are intended for agents who are working in the process of issue --> PR, effectevely doing the final intended work, working within a strcutured process. These sections should be treated as harness artifacts that will be put in the right place in the harness, and provided to these agents at runtime (as agents.md, skills, or something else)

## Fundamental Setup [for Humans or Outer Loop Agents]

Requirements for the Fundamental Setup:
- OpenCode Go Subscription API Key
- Windows with WSL2 and a Linux distribution such as Ubuntu
- Docker Engine running inside the WSL2 distribution (recommended), or one deliberately selected Docker Desktop daemon exposed through WSL integration
- Docker Compose

Use **ordinary Docker containers with Docker's default `runc` runtime** for all OpenSandbox sandboxes. Do not install or configure an alternative OCI runtime or virtualization layer. No Kubernetes cluster, Kubernetes node, RuntimeClass, Helm or Kubernetes control plane is required.

For the simplest and most predictable WSL setup, install Docker Engine directly inside the Ubuntu WSL2 distribution and run the OpenSandbox server against that daemon. Use one Docker daemon consistently: do not install Docker Engine inside WSL and accidentally point some services at the separate Docker Desktop daemon. Verify which daemon is active with `docker context show` and `docker info`. Keep any host-side checkout, setup files or bind-mounted data in the WSL Linux filesystem (for example, under `~/projects`) rather than `/mnt/c`, so Linux permissions, symlinks, filesystem events and file performance behave as expected. The actual issue working copies should live in OpenSandbox-managed Docker volumes inside their assigned sandbox; do not bind-mount the project repository into the OpenCode container. Before installing OpenSandbox, verify the daemon with `docker context show`, `docker info`, and `docker run --rm hello-world`; all Compose services and OpenSandbox must target this same daemon.

Disable OpenCode's native filesystem and terminal tools. Give each issue's OpenCode process only its assigned OpenSandbox MCP, which provides file editing and shell execution in that issue's ordinary Docker sandbox. The working repository and its durable storage belong to the issue-specific sandbox and its persistent Docker-managed storage; the OpenCode container must not mount the project repository or the Docker socket. VS Code may use separate remote windows: one connected to the OpenCode container for interacting with the agent, and another connected to the relevant coding environment for inspecting or manually editing the repository.

**Use the OpenSandbox Docker backend with Docker's default `runc` runtime.** Do not configure a `[secure_runtime]` section and do not register an alternative runtime with Docker. Configure the Docker backend with bridge networking because OpenSandbox's Docker egress sidecar and per-sandbox `networkPolicy` require it. A minimal server configuration is:

```toml
[server]
# Required when the lifecycle server runs inside Docker Compose.
host = "0.0.0.0"
port = 8080
# Supply the actual value through OPENSANDBOX_SERVER_API_KEY, not this file.

[proxy]
# The server runs on a Compose network while sandbox containers use Docker's
# default bridge; proxy through host-published sandbox ports.
resolve_internal = false

[runtime]
type = "docker"
execd_image = "opensandbox/execd:<version-compatible-with-server>"

[egress]
image = "opensandbox/egress:<version-compatible-with-server>"
mode = "dns+nft"
readiness_timeout_seconds = 30.0

[docker]
network_mode = "bridge"
# Must resolve from the OpenSandbox server and OpenCode client containers.
host_ip = "host.docker.internal"
port_range_min = 40000
port_range_max = 60000
no_new_privileges = true
drop_capabilities = ["AUDIT_WRITE", "MKNOD", "NET_ADMIN", "NET_RAW", "SYS_ADMIN", "SYS_MODULE", "SYS_PTRACE", "SYS_TIME", "SYS_TTY_CONFIG"]
pids_limit = 4096
seccomp_profile = "" # Empty selects Docker's default seccomp profile.
```

Pin OpenSandbox server, `execd` and egress images to a mutually compatible release set; pin the SDK version too, and do not use `latest`. Meet the minimum versions in the Credential Vault guide (server >= 0.2.0, egress >= 1.1.1, and the minimum version for the language SDK you actually use). The egress sidecar is attached only to sandboxes created with a `networkPolicy`.

The `[proxy].resolve_internal = false` and `[docker].host_ip` values above are for the recommended Compose topology: the OpenSandbox server is attached to a private Compose network, while sandboxes are created on Docker's default bridge and expose their `execd`/service ports on the Docker host. Set `[docker].publish_host` to the actual default-bridge gateway IPv4 obtained from `docker network inspect bridge`, rather than leaving it at the default `0.0.0.0`; this limits the host interfaces on which sandbox endpoint ports are published. Add `extra_hosts: ["host.docker.internal:host-gateway"]` to the OpenSandbox server and to every client container (including OpenCode) that must connect to returned sandbox endpoint URLs. Reserve the configured Docker port range in local firewall policy. Confirm from inside both containers that `host.docker.internal` resolves to the Docker host gateway and that an issued sandbox endpoint is reachable. If instead the OpenSandbox server runs directly as a process inside WSL and can route directly to sandbox bridge IPs, use the direct/internal proxy mode documented by OpenSandbox and choose a host address reachable by the clients; do not blindly reuse the Compose-specific settings.

**Credential Vault has a hard `dns+nft` requirement.** For credential injection, create the sandbox with both an explicit outbound `networkPolicy` and `credentialProxy.enabled = true`, then create narrowly scoped Credential Vault bindings through the egress-sidecar API. Use `defaultAction = "deny"` and explicitly allow only the destinations required by the task. The Credential Vault guide states that it refuses to activate in DNS-only mode because direct-IP connections can bypass DNS policy. Therefore, **do not fall back to `mode = "dns"` and assume secret injection still works**. DNS-only mode is an option only for sandboxes that do not use Credential Vault and only need domain-based egress filtering; it does not enforce CIDR/IP rules. Before relying on this architecture in WSL2, run an end-to-end preflight on the actual WSL2 Docker daemon: create a bridge-networked sandbox with `networkPolicy` and `credentialProxy.enabled = true`, verify the egress sidecar becomes healthy, create a test Credential Vault binding, and verify that a scoped HTTPS request receives the injected test credential while a denied destination remains blocked. Merely having the `nft` userspace command installed is not sufficient proof that the WSL2 kernel and Docker networking support the required NAT interception. If this preflight fails, fix the WSL2 kernel/network prerequisites or use a separate credential-broker design; do not silently downgrade the Credential Vault to DNS-only mode.

The Docker (`runc`) + egress configuration is explicitly listed as supported by OpenSandbox and requires `[docker].network_mode = "bridge"`. See the [OpenSandbox Docker Compose example](https://github.com/opensandbox-group/OpenSandbox/blob/main/server/docker-compose.example.yaml), [server configuration](https://github.com/opensandbox-group/OpenSandbox/blob/main/server/configuration.md), [egress compatibility notes](https://github.com/opensandbox-group/OpenSandbox/blob/main/docs/guides/egress-upstream-proxy.md), [egress architecture](https://open-sandbox.ai/architecture/network/egress), and [Credential Vault guide](https://github.com/opensandbox-group/OpenSandbox/blob/main/docs/guides/credential-vault.md).

**Security trade-off:** ordinary Docker containers share the WSL2 Linux host kernel. They provide process, filesystem, namespace, capability and network isolation, but they do not provide a separate-kernel/VM isolation boundary. Treat sandboxed agent code as untrusted, but do not describe this setup as VM-grade isolation or as a strong boundary against a kernel escape. Reduce risk with OpenSandbox's default seccomp profile, `no_new_privileges`, dropped Linux capabilities, per-issue containers and workspaces, default-deny egress policies, resource limits, and strict separation of the Docker socket. Never mount `/var/run/docker.sock` into an OpenCode or sandbox container. Only the trusted OpenSandbox control-plane service may access the Docker socket because it must create and manage sandboxes; anyone who can control that socket effectively has control over the Docker host. The WSL user/root account and anyone with access to the Docker daemon are therefore trusted administrators and can inspect sandbox files, sidecar state and container environment variables; this architecture protects secrets from the sandboxed agent, not from the local machine administrator.

If the OpenSandbox control plane runs as a Docker Compose service, mount the WSL Docker socket into that service only. Set `OPENSANDBOX_SERVER_API_KEY` in the control-plane container environment; this overrides `[server].api_key`. When the server runs in a container, bind its HTTP listener to `0.0.0.0` inside the container, but do not publish the lifecycle API port to the LAN: expose it only on the private Compose network and, when needed for local administration, bind a host port to `127.0.0.1`. Add `extra_hosts: ["host.docker.internal:host-gateway"]` to both the OpenSandbox server and OpenCode client containers because the server returns host-published sandbox endpoint URLs. This must match `[proxy].resolve_internal = false` and `[docker].host_ip = "host.docker.internal"` in the server config; do not configure only one of these pieces. Keep the API authenticated and private. A user with access to the WSL Docker daemon/socket can inspect or control the OpenSandbox control plane and its containers; treat the WSL Docker administrator as trusted. OpenCode receives only its assigned issue-specific OpenSandbox MCP endpoint, never the Docker API or socket.

Secrets needed by the agent to use tools should be stored through OpenSandbox's Credential Vault and injected at tool-use time through the supported Docker (`runc`) egress path. Do not put real external credentials in sandbox environment variables, files, command arguments, agent prompts or logs. Scope every credential binding by destination host, HTTP method and path, and ensure its destination is explicitly allowed by the sandbox's egress policy. Keep the OpenCode Go subscription credential separate: only the APISIX gateway may access it.

Use **Apache APISIX 3.19.0 or newer** as the LLM credential gateway, configured as a generic HTTP reverse proxy rather than an LLM request/response translator. APISIX is an open-source Apache-2.0 API gateway. Pin the image version. The 3.19.0 minimum is deliberate: this release makes `upstream.tls.verify` effective for HTTPS/GRPCS upstreams. Compose its OSS plugins to implement the required behavior:

1. `key-auth` validates a locally generated gateway key from a dedicated request header such as `X-Gateway-API-Key`. Create an APISIX Consumer and a `key-auth` Consumer credential for each issue's OpenCode process; its `key` must be a secret reference such as `$ENV://OPENCODE_GATEWAY_KEY_ISSUE_<id>`, not the literal key in the declarative route/consumer file. Configure the route's `header` to `X-Gateway-API-Key` and set `hide_credentials: true` so this client credential is removed before forwarding. APISIX documents environment-variable and secret-manager references for Consumer authentication credentials.
2. Use `serverless-pre-function` in the `rewrite` phase to enforce an explicit **request-header allowlist**. Preserve only the headers required by the API families used by OpenCode Go, such as `accept`, `accept-encoding`, `content-type`, `content-length`, `user-agent`, `x-gateway-api-key`, `x-opencode-session`, and the provider/version headers observed in actual OpenCode requests (for example `anthropic-version`, `anthropic-beta`, `openai-beta`, `openai-organization`, and `openai-project`). Add a provider-specific header only after confirming it is required. Remove all other client-supplied headers, especially client-supplied `Authorization`, `Proxy-Authorization`, cookies and spoofable forwarding headers. APISIX should generate its own forwarding/trace headers.
3. Use `proxy-rewrite.headers.set.Authorization` to **overwrite** the upstream `Authorization` header with the full value held by APISIX, referenced as `$ENV://OPENCODE_GO_AUTHORIZATION` (the environment variable contains the complete value, including the `Bearer ` prefix). Do not append to or preserve the client-supplied authorization value.
4. Store each issue's gateway credential (for example, `OPENCODE_GATEWAY_KEY_ISSUE_<id>`) and the complete upstream authorization value (`OPENCODE_GO_AUTHORIZATION`) in APISIX runtime environment variables or a supported secret manager. Use the corresponding `$ENV://OPENCODE_GATEWAY_KEY_ISSUE_<id>` reference in that Consumer's `key-auth.key` and `$ENV://OPENCODE_GO_AUTHORIZATION` in `proxy-rewrite.headers.set.Authorization`. APISIX supports `$ENV://...` references in Consumer authentication credentials and plugin string fields; ensure those variables are exposed to APISIX's Nginx worker through the documented `env` configuration. Because an unset APISIX environment reference can be treated as a literal string, the Compose entrypoint/preflight must explicitly test that every configured issue key and the upstream authorization variable are non-empty and exit before APISIX starts if any is missing. Do not put either real secret in route configuration, OpenCode configuration, repository files or agent-visible logs. Remember that a WSL Docker administrator can inspect container environment variables; the host Docker administrator is inside the trusted computing base.
5. Create routes for `/inference/go` and `/inference/go/*` only, with a fixed upstream node `opencode.ai:443` and `scheme = "https"`. Preserve the request method, body, query string and path; do not rewrite the path. Set the APISIX upstream `pass_host` to `node` so the upstream HTTP Host and TLS SNI are derived from the fixed `opencode.ai` node (do not try to set `Host` through `proxy-rewrite.headers.set`). Set `upstream.tls.verify = true`. Do not let the caller choose an upstream host, port or arbitrary destination.
6. Enable APISIX's `proxy-buffering` plugin with `disable_proxy_buffering: true` on inference routes so streamed responses (including SSE) reach OpenCode without response buffering.

OpenCode supports custom provider `options.headers`. Configure the OpenCode custom provider to send `X-Gateway-API-Key: {env:OPENCODE_LLM_GATEWAY_KEY}` to APISIX. Orca's issue launcher must inject the matching issue-specific local gateway key into that OpenCode process as `OPENCODE_LLM_GATEWAY_KEY`; it is an intentionally agent-visible, narrowly scoped credential for the fixed APISIX inference route, not the OpenCode Go subscription credential. Keep it process-scoped, do not write it into repository/config files, and revoke it when the issue process ends. If the selected SDK adapter requires an `apiKey` value to construct a request, use a non-secret placeholder; APISIX must overwrite `Authorization` before the request leaves the gateway. Do not configure the actual OpenCode Go key in `auth.json`, `options.apiKey`, or any OpenCode process environment.

The current OpenCode Go base URL is `https://opencode.ai/inference/go`. Configure separate OpenCode custom-provider entries for the API dialects you actually use, rather than assuming one shared `baseURL` can be combined safely with several SDK adapters. Each entry should use the matching SDK package and a `baseURL` prefix that makes the installed adapter generate exactly the documented complete path below when it appends its own endpoint suffix:
- OpenAI Chat Completions: `/inference/go/openai/v1/chat/completions`
- OpenAI Responses: `/inference/go/openai/v1/responses`
- Anthropic Messages: `/inference/go/anthropic/v1/messages`
- Gemini: `/inference/go/google/v1beta/models/<model>:generateContent` (and the corresponding `:streamGenerateContent` request where required)

For example, the OpenAI adapters can use the `/inference/go/openai/v1` base prefix, the Anthropic adapter `/inference/go/anthropic`, and the Google adapter `/inference/go/google/v1beta`, provided the installed SDK versions append the suffixes shown above. Use `@ai-sdk/openai-compatible` for Chat Completions, `@ai-sdk/openai` for Responses, `@ai-sdk/anthropic` for Messages, and `@ai-sdk/google` for Gemini. APISIX must pass requests through; it must not translate the request body or API dialect. Before using the real provider, run each configured provider against a mock upstream and assert the exact final method, path, query string, body, headers and streaming behavior. Do not invent model-catalogue or usage paths: allow those only when confirmed from the actual requests made by the installed OpenCode Go integration, and keep them under the fixed `/inference/go` prefix.

OpenCode Go requires a stable `x-opencode-session` header for every conversation. Preserve it unchanged through APISIX and include it in the header allowlist. The launcher must provide a unique, opaque `OPENCODE_GO_SESSION_ID` for each OpenCode Go conversation, and the provider configuration must send it as `x-opencode-session` using OpenCode's environment-variable substitution. Keep that identifier stable for all requests in one conversation; generate a new one for a new conversation. Do not reuse one static identifier across conversations. Verify the actual header at the APISIX ingress and at a mock upstream before using the real provider. The Go service also asks clients to identify themselves with a client-specific User-Agent; configure and preserve an identifiable value such as `freelunch-opencode/<version>` rather than a generic HTTP-library name. See the [OpenCode Go API documentation](https://dev.opencode.ai/v2/docs/console/go/), [OpenCode provider configuration](https://docs.opencode.ai/docs/providers/), and [OpenCode configuration/environment substitution](https://docs.opencode.ai/docs/config/).

For the initial single-node APISIX deployment, prefer file-driven standalone mode so etcd is not required solely for local routing. Keep the APISIX Admin API/control interface inaccessible to the OpenCode agent; expose only the inference data-plane listener. APISIX references: [key-auth](https://apisix.apache.org/docs/apisix/plugins/key-auth/), [proxy-rewrite](https://apisix.apache.org/docs/apisix/plugins/proxy-rewrite/), [serverless functions](https://apisix.apache.org/docs/apisix/plugins/serverless/), [proxy-buffering](https://apisix.apache.org/docs/apisix/plugins/proxy-buffering/), [APISIX Secrets](https://apisix.apache.org/docs/apisix/terminology/secret/), and [deployment modes](https://apisix.apache.org/docs/apisix/deployment-modes/).

Before considering the gateway integration complete, test it against a harmless/mock upstream and then OpenCode Go: (1) missing and invalid `X-Gateway-API-Key` values are rejected; (2) a valid local key is accepted; (3) the local key is not forwarded; (4) arbitrary incoming headers are removed unless allowlisted, including client-supplied `Host`, `Authorization`, `Proxy-Authorization`, cookies and forwarding headers; (5) any client-supplied `Authorization` value is removed/overwritten and the APISIX-held Go credential is the only upstream credential; (6) the real upstream credential never appears in responses, OpenCode configuration, logs or traces; (7) all four API-family paths are forwarded with their original methods, bodies, query strings and paths; (8) upstream requests use `Host: opencode.ai`, SNI `opencode.ai` and successful certificate verification; (9) streaming responses work without buffering; (10) `x-opencode-session` reaches the upstream unchanged and is stable per conversation; (11) the expected client-specific User-Agent is forwarded; and (12) requests to arbitrary hosts or paths outside the configured OpenCode Go route prefix are rejected.

In summary, local issue handling uses these components:
- One ordinary Docker container runs the OpenCode processes. It has no project-repository mount, no Docker socket and no OpenCode Go subscription credential. Each issue process receives only its own local APISIX gateway key and its assigned OpenSandbox MCP endpoint.
- The OpenSandbox server uses its local Docker backend with Docker's default `runc` runtime. Each issue gets a separate ordinary Docker sandbox container, repository/worktree and durable storage. The sandbox executes file and shell actions and provides observations through its issue-specific MCP endpoint. Transparent egress credential injection uses OpenSandbox Credential Vault through the documented Docker (`runc`) egress integration. It must be enabled only with an explicit per-sandbox `networkPolicy` and `credentialProxy.enabled = true`.
- One ordinary Docker container runs Apache APISIX. It validates the local gateway key, applies the request-header allowlist, overwrites upstream `Authorization`, and forwards to the fixed OpenCode Go endpoint without LLM protocol translation.

Docker Compose may orchestrate OpenCode, APISIX and the OpenSandbox control-plane service. The OpenSandbox control plane is trusted and may need access to the host Docker Engine; if containerized, mount the Docker socket into that service only. OpenSandbox then creates ordinary Docker sandbox containers using Docker’s default `runc` runtime. No Kubernetes resources are involved.

Notes for handling multiple issues concurrently:
- Issue workflows may share the ordinary OpenCode container and APISIX gateway, but each issue/conversation runs in its own OpenCode process/Orca terminal with its own process environment.
- Each issue workflow uses its own git worktree/branch, ordinary Docker OpenSandbox workload, persistent volume and MCP endpoint.
- `SANDBOX_ENDPOINT` is a **process-scoped environment variable**: Orca's issue launcher sets it in the environment of that issue's OpenCode process (child tool processes inherit it). The OpenCode MCP configuration references it with OpenCode's `{env:SANDBOX_ENDPOINT}` substitution. `OPENCODE_LLM_GATEWAY_KEY` is also process-scoped and contains only the local APISIX credential for that issue's Consumer. Never set either value as a single container-global value for concurrent issues.
- `SANDBOX_ENDPOINT` selects the MCP endpoint; it is routing configuration, not by itself a security boundary. Give each issue's MCP endpoint its own authentication/secure-access credential and expose that credential only to the corresponding OpenCode process. The assigned MCP server must authenticate and authorize every request and bind the authenticated identity to the assigned sandbox. Do not assume that the environment variable alone prevents a process from contacting another reachable endpoint.
- Because multiple OpenCode processes share one container and network namespace, process-scoped environment variables do not provide network isolation. Run issue processes under distinct unprivileged OS user IDs, with private `HOME`/OpenCode configuration directories and per-issue MCP credentials readable only by the corresponding user. Restrict cross-user process inspection. If Orca cannot provide and preserve this OS-level process separation, use separate OpenCode containers per issue rather than claiming that `SANDBOX_ENDPOINT` alone isolates them.
- `OPENCODE_GO_SESSION_ID` is conversation-scoped, not issue-scoped: generate a new value for every Go conversation. A child OpenCode process that starts a new conversation must get a new session ID while retaining the issue's `SANDBOX_ENDPOINT` and `OPENCODE_LLM_GATEWAY_KEY`.

Important: use Orca's terminal/session management for issue-specific terminal sessions. Each issue's OpenCode process must receive its own working directory and process environment, including `SANDBOX_ENDPOINT` and `OPENCODE_LLM_GATEWAY_KEY`; each conversation also receives its own `OPENCODE_GO_SESSION_ID`. Child sessions spawned for that issue retain the issue mapping and local gateway credential but receive a new Go session ID when they represent a new conversation. Do not manually create a terminal session that loses the issue mapping.

## Project Rules (AGENTS.md) [for Inner Loop Agents]

### Things you should NEVER do

- never access directories above the parent of the project's root directory. The project root's immediate parent is the highest directory level an agent may access; use it only when a documented workflow explicitly requires it (for example, the harness-evaluation workspace under `../tmp/`). Never traverse to `../../` or higher.
- never read sensitive files (e.g., .env)
- never run destructive commands without my permission. Always present me a dry run alternative if possible.
- never try to extract extreme performance (speed, throughput) at the expense of making code complexity significantly higher, unless specifically prompted for extreme performance optimization
- never try to bypass pre-commit hooks or ny other quality gate, unless explicitely asked to do so by me.
- never modify opencode-related folders/files such as /.agents, AGENTS.md and skills, without my explicit approval

### Things you should always do

- Before starting a task always read the global and issue-specific spec. Treat global spec as the main source of truth. If issue-specific spec differs from global spec, flag this issue for me to resolve (with your help). If implementation differs from issue-specific spec or global spec, flag this issue for me to resolve (with your help).
- Before starting something new, check the last uncommitted and commited changes made with git. Only start this new thing if nothing seems suspicious (e.g., new code was written without correspinding tests, weird code changes, etc)
- Before acting upon a review document, check if its claims are up-to-date with respect to the current codebase. There may be claims or entire documents that are stale, meaning that they refer to old things of the codebase that are now different in the latest commit.
- Before using an unfamiliar dependency/API, consult its official documentation relevant to the operation being performed. Do not reread documentation already understood in the current session.
- every review document (in ./.agent/session-persistent-candidate/reviews/) created should contain in its initial metadata a reference to the exact version of what was reviewed which can be a file of specific commit (e.g., spec review and security spec review) or an entire commit (e.g., code review and security code review).
- log all mistakes you made in ./.agent/persistant/knowledge/mistakes.jsonl file, each entry in the form {"what_was_done": "placeholder", "what was wrong": "placeholder", "why it was wrong": "placeholder", "how the mistake was corrected": placeholder, "commit hash of the codebase at the time of mistake": placeholder}. 
    - What counts as mistakes?
        - Anything you realize you did wrong before, having evidence to support why its wrong and explanation of why its wrong
        - Anything the I (aka your user) had to intervene to change something you already did becomes it had serious problems. I might say explicitely that you did something wrong (e.g., "change this code you wrote because its not readable", "change these tests you wrote becaue they dont reflect the spec", "change your implementation plan to more fine-grained end-to-end steps, where you start by ...") or just ask you if you are shure something is correct. Note: before counting it as a mistake and changing it, you must consider that I might be wrong, therefore you must reason about it and come to a conclusion based on arguments and evidence.
    - Whats doesnt count as mistakes? (1) User intervetions where the user requests spec changes/updates are never mistakes; (2) User interverntions for non-critical low-level changes that to satisfy his preferences (e.g., "extend this class to also support this capability that is currently outisde of the class")
- Log all codebase-related (e.g., how does this function work?) questions I ask to you in a ./.agent/persistant/user-understanding/user-codebase-questions.jsonl, each entry in the form "question": "placeholder", "answer": "placeholder", "branch": "placeholder", "commit_hash": placeholder"} where branch is the branch in which the user is in when he asked the question and commit_hash is the last commit of the branch that user was in, at the time the user asked the question.
- Log all new usefull nuggets of codebase knowledge to ./.agent/session-persistent-candidate/knowledge/non_obvious_conjectures.md where you keep track of non-obvious conjectures you make for the session. Each conjecture has the following data: (1) description of the conjecture, current evidence of the conjecture, already a fact? (yes or no) and risk if wrong (high, medium or low). Always ask question to user before you are about to act on a risky conjecture you dont have much evidence.
- keep your code clean and organized, refactoring might be needed
- modularization: the codebase should have a few big modules with clear boundaries and relationships, and each big module is composed of many little modules. Dont let files become too big, prefer breaking into multiple files where each one has a clear meaning/job.
- if using bash commands for file/content search: prefer `fd` (fdfind) and `rg` (ripgrep) over standard `find` and `grep` for better performance and git-awareness.
- always make a plan before doing stuff
- before concluding a task, critically re-evaluate your reasoning, assumptions, and implementation. Verify that the solution satisfies the user's objective, that no possibly affected areas have been overlooked, and that no unnecessary regressions have been introduced.
- when E2E testing a product: be picky about the UI you see and be obsessed with pixel perfection. 
- If something unrelated looks wrong, record it as a todo in ./.agent/session/todos.md it creates a correctness, security, build, or test failure affecting the current task. do not do one todo item at a time, batch them into related todos, and implement one batch at a time.
- If you realize that you are stuck in a loop where you can't solve a specific problem, you need to do something different. Save the problem description in ./.agent/session/problem_stuck.md (should have: problem summary, how to verify solution of problem, things already tried and the output gotten in each attempt) and then try these in order:
  1. Just reset the opencode context window and try again (the problem description will be at ./.agent/session/problem_stuck.md). 
  2. Reset the opencode context window and Change the OpenCode Coding Model (e.g., change from Kimi K3 to Opus 4.8) and increase it's reasoning if possible
  3. Reset the opencode context window and Take a fresh look at the problem. Try a different approach (if necessary you can update plan.md or core-implementaiton-tasks.md)
  4. Spawn the most approppriate sub-agent to help
  5. Spawn a codex child terminal session and prompt to to solve the problem (needs codex installed)
  6. Reset to the latest commit as your last resource, then reset the OpenCode context window, then try agian to solve the problem.
- If you solve the problem you are stuck, delete ./.agent/session/problem_stuck.md immediatelly.
- if you encounter code-spec-docs mismatch you should explain the mismatch, initiate a discussion with the user, which wil culminate in either code, spec or docs change (or multiple changes).
- always when you get stuck in a problem, revise ./.agent/flow/current-issue/flow/core-implementation-tasks-plan.md or plan.md to see if plan changes need to be mande. Remember that core-implementation-tasks-plan.md stores the graph of core implementaiton tasks along with progress, its per-issue; ./.agent/session/plan.md store per-step action plans typically generate by using the ai agent (you) in plan mode for steos that require a plan first.
- if you want to explore (a planned big-refactoring doesnt count as exploration and should be done just on a big-refactoring branch with the same agent) some idea/hypothesis without clothering the issue handling branch, create a separate git worktree (worktrees should be created in .agent/worktrees/). In the new gitworktree checkout to an exploration branch and spawn another opencode instance to explore. The opencode instance should end either when he considers the exploration finished or you (the main agent) should end the opencode instance if he consumed more than 1 dollar worth in tokens spent. The epxloratory opencode instance should always store findings in a findings.md file at the root of the exploration branch. When the exploration ends, you should move the findings to ./.agent/session-persistent-candidate/knowledge/exploration_findings/name-of-the-exploration-placeholder.md in the issue handling, where the findings file should have these metadata in the header (exploration context, exploration description, opencode instance used, tokens consumed, dollards spent, why it ended) and the findings and conclusion in the body of the file. To enforce this process you should run a pre-built launch_exploration_subagent.sh bash script that takes care if enforcing the token limit, launching opencode in autopilot mode in a new terminal and moving the findings and terminating the exploration.
- Before building any GUI, need to: (1) have a mock/prototype validated with the user; (2) write a design.md to standardize GUI components and patterns.
- If I give you you a mock.html as guidance, you should open the mock, create synthetic goals the end-user will want to achieve within the GUI, and actually use the mock in the context of achieving these synthetic goals. In the end, write your improved understanding of the mock, i.e., write you understanding of the GUI experience I want to build in a mock_learnings.md at the same directory as mock.html. The you ask for my approval to use this mock as the implementation guide.
- all core code of a project needs to be inside ./src directory at project repo root
- for quizzing me, always look at Use .agent/persistant/user-understanding/user-codebase-questions.jsonl to know my comon weak uderstanding spots
- Use asynchronous/concurrent execution when its clearly the right solution for the scenario, particularly for I/O-bound work. Don't introduce async merely because it is technically possible and more elegant.
- Don't choose a technically inferior architecture merely because it is slightly cheaper to implement when a significantly better design is available at reasonable complexity.
- Wehn doing any kind of artifact optimization (e.g., function performance optimization): never switch the current implementation for a candidate one before comparing both on the relevant evaluations. IF the candidate beats the current in the final eval score, than you can change, and the candidate then becomes the current implementation.
- if you are doing/will probably do multiple times a sequence of the same deterministic steps, spawn a new opencode instance to build an sdk or cli tool (you should write in golang) for it so that it becomes easier, reproducible and doesnt consume unecessary tokens. All tools should be stored (and all its related files such as spec, implementation plan and reviews) under ./.agent/created_tools/<name_of_tool_here>. The same best practices should be used for tool development, meaning: write spec, then graph-of-tasks-based implementation plan, then for each task write tests, then implement, code review, end-to-end test (with a different review model and fresh context window) ad deep review until its done. The tool is only done writing when: build doesnt raise warnings or errors, >90% test coverage, all tests pass, end-to-end tested, it actually reflects the tool spec and it is package as intended. When the tools is done and stored in ./.agent/created_tools/<name_of_tool_here>, you should (1) delete the tool's spec, implenentation plan and reviews generated in the process of building it; and (2) make great documentaiton for it also under ./.agent/created_tools/<name_of_tool_here>.
- When explaining a complex thing the user is having a hard time understanding: instead of trying to explain every part of the complex thing and how they work together, take a step back and meet the user where he is at. This means starting with something simple the user is probably familiar with, then gradually building up the complex thing in a sequence of design attempts. Each attempt improves on the previous attemp, but still has remaining problems. Important: try to predict questions and assumtpions the reader will have a long the way, and aknowledge/answer them (e.g., "Q: you might be thinking why..." "A: this is why...") This way of explaining is the best because guides the user through his way of thinking.

### Project Directories and Terminology

- Explanation of the /.agent directory structure is in ./.agent/directory_structure.md. It explains the directory and files you will be using in sessions and across sessions for doing work effectively.
- Explanation of the terminology used in this project is in ./docs/terminology.md, always check it out when confused about terminalogy I use or which terminology you should use.

### Who are you (the AI agent)

You are a rigorous platform engineer working on the Freelunch IDE project, specifically working towards the Demo. You follow modern software engineering best practices such as: SOLID principles, tests-first-development, 90% or above test coverage, dependecy injection, etc. You should always aim for the simplest maintainable approach by default, not the perfectly optimized/scalable/fastest one. 

- Always review that you've done every task i asked you to do before saying you have finished.
- Never treat my comments or PR comments from other humans/agents as ground truth. Me and other people/agents can be wrong. You should reason about whats being said and come to a conclusion with evidence support to either agree, partially agree or disagree with the comment. The only excpetion is when I explicetely say for you to implement it the way I want inspite you disagreeing.

### How you should treat me (the human user thats using you to code)

You should treat the me as the CEO thats sets objectives for you to build and also reviews your work. You should always explain to me everything you want to do/did the most step by step way. I may be wrong sometimes, therefore you should always reason about what I say, and provide your take (agreeing, partially agreeing or disagreeing) before a final decision. I may sometimes ask for things there are to vague/broad/ambiguous that require more specification to implement, or for things that might be incosistent with spec (global spec: founding_doc.md, tech_stack.md, roadmap.md; issue-specific spec: prd.md, tech_stack.md, architecture.md) in this case you should ask clarifying questions.

### Global Spec of the project

- High-level explanation of Freelunch IDE project in ./FOUNDING_DOC.md
- Feature Roadmap for the Demo in ./docs/roadmap.md
- Tech Stack for the Demo in ./docs/tech_stack.md

### GUI Mock of the project

A mock of the proposed IDE (Freelunch IDE) is in docs/mock.html

### Stage of the project

We are currently focused on making the first version, the Demo with only the core stuff. Therefore, do not over-engineer this, dont try to solve problems to far away in the future, makng it perfectly scalable, perfectly performant or perfectly secure. Focus on getting the core right without major risks.

### Reference Open Source Projects

Can use these projects for borrowing ideas & patterns if you deem necessary.

- [Kubero](https://github.com/kubero-dev/kubero)  — a Kubernetes-based, developer-friendly platform
- [Kubefirst](https://github.com/konstructio/kubefirst)  - Modern K8s-based internal developer platform template
- [OKD](https://github.com/okd-project/okd)  — open source edition of Red Hat OpenShift, a k8s-based complete platform focused on enterprises
- [Tilt](https://github.com/tilt-dev/tilt)  — a strong dev/experimentation experience for Kubernetes
- [Backstage](https://github.com/backstage/backstage) — a plugin-based internal developer platform interface
- [Ray](https://github.com/ray-project/ray) — modern distributed programming framework for Python (inspiration for the lunch-lang distributed programming framework idea, to be used within freelunch-ide, though ray works as runtime and lunch-lang would be at compile time)

### Replanning mid-coding

You might be doing a step and realize your plan.md or core-implementation-tasks-plan.md needs to be changed in some way. You should change immediately. How you should change them:
- if want to change plan.md: you can change directly, overwritting the file.
- if want to change core-implementation-tasks-plan.md: you should not overwrite the file, you should append to it the reason of the replanning and a summary of the current state of the codebase, then append a new core implementation tasks plan graph aling with an explanation of what was changed and why. So the resulting file will actually contain (in order) the previous core tasks plan and the new one.

### Big Refactoring mid-coding: rewriting tests and/or modifying scaffold (directories, files, interfaces, data models, etc; the skeleton in which logic gets written inside)

You might be doing a step and realize your tests and/or scaffold (directories, files, interfaces, data models, etc; the skeleton in which logic gets written inside) needs to be changed in multiple ways. You should create and checkout to an epehemeral big-refacoring branch (which was sourced from the current branch, not main) and then do the refactoring in the big-refacoring branch. When you are done, ask for my approval to merge big-refacoring into the issue handling branch.

Small or Localized refactorings you can just do, without this branching and asking for my approval ceremony.

### Patterns to use

- testing folder that mimicks the actual folder structure
- always work on the standard virtual environment for the project
- test-first development where tests are wrritne before code (this is not: strict tdd where need to make one function red -> green at a time and only then to move to the next)
- Use dependency injection where it improves testability or separation of concerns; 

### Code Quality & Testing

- avoid excessive dependencies
- Every public function/class must have documentation explaining purpose, inputs/outputs, important behavior and non-obvious constraints.
- write documentation at the beggining of each file assuming the reader is a new programmer
- use meaningful variable and function names
- Remove dead code
- Avoid duplication
- Explicit error handling
- New behavior requires tests.
- Do not claim completion if tests fail.
- Run the smallest relevant test suite first, then broader validation.
- Only when doing code review: Use the tool lizard (https://github.com/terryyin/lizard) for measuring clyclomatic complexity of functions. The results will help you analyze if some functions are too complex or not.

### Your Confidence Protocol

Distinguish facts from inferences, assumptions and uncertainties.
- FACT — directly verified
- INFERENCE — strongly inferred from evidence
- ASSUMPTION — not yet verified
- UNCERTAIN — insufficient evidence

Protocol:
- Inference → may proceed and do what you want.
- Low-risk assumption → may proceed and do what you want.
- Medium-risk assumption → proceed only if easily reversible.
- High-risk assumption → ask user before acting.
- uncertainty → dont use this to inform your next actions, requires isolated exploration.

### Bug Handling

- reproduce the bug in an E2E (as much as possible) setting as closely aligned to the end use to make shure you are solving the actual usage problem
- After bug reproduction, use the smallest relevant tests for diagnosis and iteration.
- Find the Root Cause of the Bug. Make hypothessis and test your hypotheses until you find the real hypothesis that is the root cause of the bug.
- assume the problem may have broader implications than are immediately apparent. Investigate affected code paths, dependencies, interfaces, and related components before concluding that the required change is isolated.
- understand why something broke before changing it. To understand you need to come up with a hypothesis and test the hypothesis.
- use debugger if possible

### Linting

Always keep a look at linting warnings and errors. Fix them immediately.

### Evidence-based Statements

Claims about correctness, test status, coverage, security, compatibility, architecture conformance, or feature behavior require explicit evidence.

For example:

Bad:

“The API is backwards compatible.”

Good:

“I verified backwards compatibility by running X tests against versions A/B. Results: ...”

Similarly:

- “The feature works” → needs E2E evidence
- “93% coverage” → needs coverage report
- “No security issue” → needs security scan/review evidence
- “This dependency supports X” → needs documentation reference
- “The architecture matches the spec” → needs explicit comparison

### Security

- Never hardcode secrets
- Validate external input
- Escape shell arguments
- use the principle of least privilege

### Testing Order

- Implent first (Always Required): Unit tests 
- Implement second (Always Required): Integration tests 
- Implement third (only required if the project contains a user-facing API, CLI or GUI/Frontend): End-to-end tests

Skipping any required level = Not Complete

> This is the end of the AGENTS.md file

---

## Issue Flow (tool-agnostic) [for Humans or Outer Loop Agents]

Notes for implementation:

- each unique step is a slash command
- each step is done by a single agent
- the agent can decide to go back to a previous step (e.g., encountered a problem that requires change to spec)
- most steps have a planning sub-step performed at the beginning with the harness' plan mode which generates `.agent/session/plan.md`
- the agent can go back from a step to a previous step if necessary to fix an issue created earlier. Step jumps must be tracked in `issue_flow_[i].md`, including the issue name, issue number, timestamp of creation and completion, and the jump/deviation. The issue flow is mostly sequential, but step B will hold parallel sequential paths, one for each core implementation task.
- when a new feature (tackling a new issue) starts, first search for any completed `.agent/flow/completed-issues/completed_issue_flows/issue_flow_[i].md` and use it when relevant as historical context for the issue being handled
- session summary hook: when a session ends, store a summary of key things done, key problems encountered, tips/learnings/todos in the current issue's `issue_flow_[i].md` under the step the agent was in, using this JSON form:
  `{"key things done": "placeholder", "key problems encountered": {"problem": "placeholder", "solved_or_not": placeholder, "tips for next agent working on this": "placeholder"}, "learnings": "", "todos": "placeholder"}`
- approval gates mean that either the user (developer) or a specific AI agent needs to give approval to continue the flow
- An author self-check may be performed, but it is not a review and cannot satisfy a review or verification gate. Every LLM-based review or verification must use a different effective model from the model that authored the artifact being reviewed/verified. When multiple reviewers or verifiers are involved in the same work unit, use pairwise-distinct effective models. Prefer different model families/providers, not merely different aliases or reasoning settings. Deterministic build/test/static-analysis tools remain independent evidence and do not need an LLM. If model availability prevents distinct models, explicitly mark the check as non-independent and request approval before treating a critical gate as complete.
- The workflow is initially human-directed; slash commands expose the steps but the human chooses when to invoke them. Do not add an automated workflow state machine yet.
- The main issue-flow steps must be numbered sequentially from 1, with no gaps or duplicates; child steps must use stable hierarchical identifiers (for example, `12.a`). Implement `validate-issue-flow` as an automated validator for this authoritative Issue Flow section. It must check sequential main-step numbering, unique child identifiers, that every parent exists, that every explicit step/cross-reference/loop target resolves to an existing identifier, and that referenced step ranges are valid. Run it whenever this document's Issue Flow section or the generated slash-command definitions change, and run it before relying on the flow. Store the validator under `.agent/created_tools/validate-issue-flow/` and include tests for malformed numbering, duplicate identifiers, dangling references and invalid loop ranges.
- when a session ends, a hook must be called to: (1) empty all files inside `.agent/session`; (2) recursively analyze `.agent/session-persistent-candidate` and transfer useful validated knowledge to `.agent/persistant/knowledge/non_obvious_conjectures_and_facts.md`, then empty `.agent/session-persistent-candidate`; and (3) assess whether harness-related files need improvement or bloat removal using the `harness-eval` skill, by first proposing harness improvement ideas, then evalauting them. Do 3 rounds of harness imrpovement idea geenration and running evals so that you can test combinations of what works best. Some harness changes might not be directly evaluatable (e.g., harness-eval skill doesnt support directly evaluating a planning skill) however there is always a downstream evaluation that can be done (a good planning skill should make the agent write spec and code better, provided that everythig also remains the same, so you can eval spec writing and coding downstream effects with harness-eval skill). Before actually modifying harness-related files, ask for user approval and provide the evaluation results gotten. If changes are approved by me (user), git commit them.
- At the start of every slash command, retain the automatic commit hook: if there are uncommitted changes that the workflow says should not be left outstanding, the hook must commit them before the command begins. Before generating the commit message, study the relevant workflow artifacts under `.agent/` to establish exactly what was done, which implementation-plan tasks were completed, what problems were encountered, and how they were resolved. At minimum inspect the current issue's `issue_flow_[i].md`, `core-implementation-tasks-plan.md`, relevant issue-specific spec, the latest `.agent/session/plan.md`, session summaries, and any relevant PR-review/review artifacts; compare these with `git status` and the actual diff so the message describes the real changes rather than merely repeating the last step name. Use a concise conventional commit subject and a body when needed to explain the scope, rationale, and important problem resolution. Do not invent work not evidenced by the diff/workflow artifacts, and do not include unrelated user-owned changes in the commit. **Exception: when the slash command is starting an inferred harness A/B evaluation of the current uncommitted working tree, do not auto-commit the harness changes before the evaluation snapshot is created.** In a new session without context, first reconstruct the commit context from `.agent/flow/current-issue/flow/issue_flow_[i].md` and the other relevant workflow artifacts.
- ensure language-specific (in our case Go and TypeScript) linters and formatters are running continuously on every file edit, using OpenCode hooks
- have a git pre-commit hook that builds (successfully without warnings), runs tests (all must pass to consider a unit of work done), and verifies test coverage is above 90%

### Issue Flow (note: bug, refactoring or performance issue handling may skip multiple steps that are unnecessary)

The workflow is human-directed for now: each unique command is exposed as a slash command, and the developer decides when to invoke it. The flow is sequential by default, but an agent may return to an earlier step when new evidence requires a specification, test, design or implementation change. Every such deviation must be recorded in `issue_flow_[i].md` with the reason and affected step identifiers.

A: Issue-specific Spec & Core Implementation Tasks Plan

1. **/start Start Issue Handling**: point to the GitHub issue. Before doing anything: (1) if there is not already a dedicated sandbox for this issue, create a new ordinary Docker OpenSandbox environment for it, with its issue-specific repository/worktree and persistent storage; (2) ensure the issue has its own MCP endpoint and mapping in the issue-sandbox configuration; (3) have Orca launch the issue's terminal/OpenCode process with `SANDBOX_ENDPOINT` set in that process environment only (child processes inherit it); (4) configure that OpenCode process to use this endpoint as its only filesystem/terminal MCP and disable native filesystem and terminal tools; (5) ensure the process has access to APISIX but not the actual OpenCode Go credential. Study the repository and branching strategy, inspect the current git state, and ask the user any clarifying questions needed to reach a shared understanding of the problem and intended solution. This step ends when the problem/solution understanding is agreed.

2. **/quizme Understand the Codebase**: quiz the developer on high-level architecture, design decisions, tradeoffs and low-level codebase details until the developer understands the codebase. [AI Approval Gate]

3. **/reproducebug** (bug issues only): reproduce the reported bug by writing and running a failing test, and confirm that it fails for the reason described in the issue.

4. Loop until **4.b and 4.c** are successful [User Approval Gate with Independent AI Reviewer Suggestions]
   - **4.a /spec Build Issue-specific Spec**: create `prd.md`, `architecture.md` and `tech_stack.md` under `.agent/flow/current-issue/issue-spec_[i]/`.
   - **4.b /reviewspec Review Spec**: use an Independent AI Reviewer to check consistency with the Global Spec (Founding Doc + Roadmap + Tech Stack), catch omissions and flag problems in the Global Spec. If a Global Spec problem is found, modify the Global Spec first, then update the issue-specific spec. Store the final review in `.agent/session-persistent-candidate/reviews/final_spec_review_[timestamp].md`. [User Approval Gate]
   - **4.c /specsecreview Specialized Spec Security Review**: flag critical security problems and warnings. Use a model different from the spec author and the independent spec reviewer. Store the final security review in `.agent/session-persistent-candidate/reviews/final_spec_sec_review_[timestamp].md`.

----<<separate terminal block (reset context)>>----

B: Core Implementation

5. **/make-tasks-plan Convert the Spec into a Dependency Graph** in `.agent/flow/current-issue/flow/core-implementation-tasks-plan.md`. Each node must represent a coherent behavioral/functional slice rather than a component-only task. [User Approval Gate with Independent AI Plan Reviewer Suggestions]

6. **/tasksplangrillme Grill User on the Task Plan** until the developer understands the plan. [AI Approval Gate]

7. **Core Implementation Loop**: process each pending core implementation task in the graph until the graph is complete. Each task follows steps 8–13 below.

---- New Session (reset context) ----

B1: Common Scaffold

8. **/scaffolddep Define Allowed Scaffold Dependencies**: choose dependencies for the common scaffold. Have a different model review the dependency decision. Dependencies must be actively maintained. A dependency may be preferred when it has important capabilities the alternatives lack, is significantly more popular, is significantly easier to use, or is significantly older. [User Approval Gate with Independent AI Reviewer Suggestions]

9. **/scaffold Set Up/Modify the Common Scaffold**: make/remake `.agent/session/plan.md` first; create or modify the required project skeleton, signatures, dependencies, build/test/package automations and other foundation pieces. Review against the Issue-specific and Global Specs. Leave function/class bodies as high-level pseudocode rather than implementing feature logic. [User Approval Gate with Independent AI Reviewer Suggestions]

B2: Tests & Logic

10. **/writetests Write/Modify Functional Tests**: make/remake `.agent/session/plan.md` first; write or modify unit and integration tests and review them against the Issue-specific and Global Specs. Flag inconsistencies rather than silently resolving spec conflicts. [User Approval Gate with an Independent AI Reviewer]

11. **/testtests Test the Functional Tests with Placeholder Feature Code**: all functional tests must fail in this phase, while the test files themselves must be valid and runnable. Ensure strict test coverage of the intended core logic (above 90% of the entire codebase once implemented, ideally near 100%).

12. Loop through **12.a–12.c** until the feature implementation satisfies the functional tests and validation requirements. [User Approval Gate with Independent AI Reviewer Suggestions]
   - **12.a /featdep Define Allowed Feature-code Dependencies** with an explanation of why each is used. Apply the same maintenance/popularity/ease/capability criteria as scaffold dependencies. [User Approval Gate with Independent AI Reviewer Suggestions]
   - **12.b /feat Write Feature Code**: make/remake `.agent/session/plan.md` first; implement using only the allowed feature-code dependencies. Implement the smallest coherent behavioral slice and try to make multiple related tests pass at a time. Review against the Issue-specific and Global Specs and Design System when GUI work is involved. [User Approval Gate with an Independent AI Reviewer]
   - **12.c /test Build and Test Feature Code**: make/remake `.agent/session/plan.md` first; build and run functional tests; generate test and coverage reports; review against the Issue-specific and Global Specs; repeat until the build has no errors/warnings and all required tests pass. The verification model must differ from the implementation model.

13. **/refact-coreimplementation-task-if-necessary**: make/remake `.agent/session/plan.md` first; evaluate refactoring opportunities limited to the last completed core implementation task. Refactor only when behavior is preserved and clarity, quality or maintainability improves. After each refactoring, evaluate whether it is actually better using an independent model and the relevant tests/evaluations; otherwise keep the previous version. [User Approval Gate] [Independent AI Approval Gate]

----<<separate terminal block (reset context)>>----

14. **/quizme Quiz User on the Latest Changes**: the developer reviews the code and asks questions; the agent grills the developer until the developer understands the latest changes. [AI Approval Gate]

15. **/stripdebuglogs**: make/remake `.agent/session/plan.md` first; remove debug logs if any remain. [User Approval Gate with an Independent AI Reviewer]

16. **/determine-if-core-implementation-task-is-done**: determine whether the current core implementation task objectives were achieved and whether the task graph still holds. Do not rely on historical test output; rerun the entire suite for the current task and previously completed core tasks to catch regressions. Use an independent verification model. Only mark the task done when tests reflect the spec and pass.

---- New Session (reset context) ----

17. **/end-to-end-testing**: test end-to-end to make sure the issue was completely handled. For GUIs, include usability/visual testing that emulates real user behavior. Also perform risk-appropriate stress testing: for APIs, load/availability/latency tests; for GUIs, aggressive scrolling, random group browsing, rapid navigation and repeated media loading. The E2E verification model must differ from the implementation and earlier review models.

---- New Session (reset context) ----

C: Code Review

18. Loop until code review succeeds [User Approval Gate with Independent AI Reviewer Suggestions]
   - **18.a /review Independent Code Review**: make/remake `.agent/session/plan.md` first; review maintainability, modularity, test coverage, file sizes, test/spec compliance and GUI compliance where relevant. Run mutation testing last. Identify omissions and spec/design problems. Use a model distinct from the implementer and other reviewers. Store the final code review in `.agent/session-persistent-candidate/reviews/final_code_review_[timestamp].md`.
   - **18.b /redo Make Necessary Code/Test Changes**: make/remake `.agent/session/plan.md` first; implement justified changes, build and test, and review against Issue-specific and Global Specs. [User Approval Gate with an Independent AI Reviewer]

---- New Session (reset context) ----

D: Security Review, Documentation & PR

19. Loop until security review succeeds [User Approval Gate with Independent AI Reviewer Suggestions]
   - **19.a /secreview Specialized Security Review**: flag critical security problems and warnings. Use a model different from the implementer and code reviewer. Store the final security review in `.agent/session-persistent-candidate/reviews/final_code_sec_review_[timestamp].md`.
   - **19.b /redo Make Necessary Code/Test/Docs Changes**: make/remake `.agent/session/plan.md` first; implement justified changes, build and test, then review against Issue-specific and Global Specs. The verification/review model must differ from the implementer and security reviewer. [User Approval Gate]

20. Loop until End-to-end testing succeeds
   - **20.a /end-to-end-testing**: verify the complete issue behavior in the end-user environment. Use an independent model distinct from the implementer, code reviewer and security reviewer.
   - **20.b /redo Make Necessary Code/Test/Docs Changes** as required. Rebuild and rerun required tests; use an independent verification model.

----<<separate terminal block (reset context)>>----

21. **/quizme Quiz User on the Final Changes** until the developer understands the final changes. [AI Approval Gate]

22. **/write-pr**: check `.agent/flow/current-pr/pre-pr-checklist.md`, confirm it is satisfied, then write the PR to `.agent/flow/current-pr/pr.md`. The PR is not final until documentation and required reviews are complete.

23. **/document Document**: (A) Contributor Documentation: explain how to understand the codebase through a step-by-step tutorial from high to low level; (B) Final User Documentation only if the product is already usable. Use an independent documentation reviewer. [User Approval Gate]

24. **/deep-review**: perform two phases. Phase 1 runs multiple problem-finding rounds with fresh contexts and different models, until no materially new problems are found. Phase 2 uses three verifiers to assess each candidate problem; a problem is validated only when at least two verifiers agree that it can cause a production failure under a defined condition. Use different models for the verifiers. Do not reuse a model with merely different prompt wording and describe it as independent. If there are not enough distinct models, report the reduced independence and request approval before treating the review as complete. Code styling/formatting preferences and file splitting alone are not production problems.

25. **/send-pr**: push and open the draft PR using `.agent/flow/current-pr/pr.md`.

---- New Session (reset context) ----

E: Make Fixes Based on PR Reviews and/or CI Failures Until PR Is Merged

26. Loop until all required external review/CI failures are addressed [User Approval Gate with Independent AI Reviewer Suggestions]
   - **26.a /prreviews**: after the user checks a PR review or CI failure notification, read the PR reviews and CI runs from GitHub and write them locally in the dedicated current-issue PR-review area.
   - **26.b /review Independent Code Review**: make/remake `.agent/session/plan.md` first; review maintainability, modularity, test coverage, file sizes, test/spec compliance and run mutation testing last. Use a model different from the model that made the changes and from the earlier code reviewer. Store the final code review in `.agent/session-persistent-candidate/reviews/final_code_review_[timestamp].md`.
   - **26.c /redo Make Necessary Code/Test/Docs Changes**: apply justified fixes, build and test, and review against Issue-specific and Global Specs. Use an independent verification model. [User Approval Gate]

---- New Session (reset context) ----

27. Loop until security review succeeds
   - **27.a /secreview Specialized Independent Security Review**: use a distinct security-review model and store the final security review in `.agent/session-persistent-candidate/reviews/final_code_sec_review_[timestamp].md`.
   - **27.b /redo Make Necessary Code/Test/Docs Changes** and validate with a model distinct from the implementer and security reviewer. [User Approval Gate]

28. Loop until final End-to-end verification succeeds
   - **28.a /end-to-end-testing**: test end-to-end to make sure the issue was completely handled, including GUI usability/visual testing and risk-appropriate stress testing. Use a model distinct from the implementer and all preceding reviewers.
   - **28.b /redo Make Necessary Code/Test/Docs Changes** and validate with an independent model. [User Approval Gate]

----<<separate terminal block (reset context)>>----

29. **/quizme Quiz User on Changes Since Last Quiz** until the developer understands the changes made since the previous quiz. [AI Approval Gate]

30. **/write-pr**: recheck `.agent/flow/current-pr/pre-pr-checklist.md`, then regenerate `.agent/flow/current-pr/pr.md` so it reflects the final state.

31. **/document**: update contributor documentation and usable final user documentation. Use an independent documentation reviewer. [User Approval Gate]

32. **/deep-review**: repeat the two-phase review described in step 24 for the post-review changes. Use fresh contexts and models distinct from the authors of those changes and from the earlier reviewers. Store the final validated problems list in `.agent/session-persistent-candidate/reviews/problems_list_[timestamp].md`.

33. **/send-pr**: push/open/update the draft PR using `.agent/flow/current-pr/pr.md`.

The human user is responsible for checking when the PR is ultimately merged.

> This is the end of the Issue Flow section.

## AI-assisted Coding Tech Stack [for Humans or Outer Loop Agents]

External Vendor Requirements: Opencode Go Subcription, Claude Credits, Github Repo

- Agent-native Editor/Terminal: Orca (only when tackling multiple issues in parallel)
- Terminal Harnesses: OpenCode + failproofai + claude-tap, open-code-review and openwiki
- IDE (for better introspection + manual editing): VSCode (with language-specific linters and formatters plugins running continuosly on every file edit)
- LLM Provider Subscription: **OpenCode Go**
- LLM credential gateway / transparent HTTP reverse proxy: **Apache APISIX (Apache-2.0)**, using `key-auth` for local gateway-key validation and `proxy-rewrite` to replace the upstream `Authorization` header. Use file-driven standalone mode for the initial local single-node setup; the initial configuration should not require etcd. This is a generic API gateway, not an LLM protocol translator.
- OpenCode Go upstream: `https://opencode.ai/inference/go`; preserve the API-family path and the stable `x-opencode-session` header for each conversation.
- Models: (0) Planning: good coding model with medium reasoning; (1) Spec/scaffold/tests: good coding model with high reasoning; (2) Core coding: best current coding quality/price model with medium reasoning; (3) Independent Spec Review: a different effective model from the spec author; (4) Plan Review: a different effective model from the plan author; (5) Security Review: a distinct security-capable model; (6) End-to-end testing/verification: a distinct model from the implementation and prior reviewers; (7) Sub-agents: suitable coding/research models; (8) Code Review: two candidate-generating models from different model families/providers and a third distinct model as final adjudicator; (9) quizme: coding-quality model with medium reasoning; (10) UI Testing: best current computer-use model. For every review and AI verification, use a model different from the implementer and from the other reviewers/verifiers for that work unit; prefer distinct model families/providers. Record the actual model identifiers in the workflow artifacts.
- Local Routing: opencode-model-router (opencode plugin)
    - Fast Model: qwen2.5-coder:7b
    - Medium Model: current coder model
    - Heavy Model: current coder model
- Issue Flow: Start with just making each step a slash command, and leave to the developer to follow the steps (note: this doesnt enforce step execution, se requires developer commitment). Note: each slash commands should remember to update the issue resolving progress at the end (issue_flow_[i].md file) or, if its the first step, create the file if its still not created). Each slash command should have a simple name.
- OpenCode Plugins: graphify, rtk, opencode-quota, opencode-model-router, opencode-trace
- OpenCode External Tools: Github CLI, Context 7 CLI
- Custom Freelunch OpenCode sub-agents: security-specialist, code-review-and-refactoring-specialist, testing-specialist, debugging-specialist. Tip: Agency-agents repo provides some sub-agents out-of-the-box.
- Custom Freelunch OpenCode Slash Commands: one for each unique step of the issue building flow
- Custom workflow validator: `validate-issue-flow`, stored under `.agent/created_tools/validate-issue-flow/`, automatically checks step numbering and all cross-references whenever the flow or slash-command definitions change
- Dependency Docs: a dependency_docs.md under ./docs with entries in the form "- <dependency>: <docs_link>" for all direct dependencies (not dependencies of dependencies). Pinned versions used in the project cna be seen in the lock file of the virtual dev environment tool.
- Skills: 
    - Custom Skills: created on-demand, under a `custom_skills` folder, via manual creation or via `skill-creator` to avoid having to repeat the same solution process over and over. Make these custom skills:
        - harness-eval (evalautes proposed harness changes on the current codebase, explained better in the end of this doc)
        - quizme (continually ask questions of the latest changes to codebase to me, to see if i understand the codebase. Always give score my answers and give feedback to it. Only stop when you feel i understand the codebase. The user can also specifify specific files for you to grill him about instead of the entire codebase) (every quizme run should log all the questions, answers and feedback gave to the user inot a .agent/persistant/user-understanding/user-quizzes/quiz[i].md where i is the id of the quiz and the file should have timestamp, commit, what the quiz was about and grill score in the beggining of it. Ever grill should start by looking at the commit, what was quizzed in the last grill and user-codebase-questions.jsonl file)
        - understand-external-codebase (1. Build a doc eplxianng in detail the characteristics and internals of an external github codebase; 2. Add to this doc an explanation of where and why this codebase can be helpfull as a reference for ideias/patterns for the current project being built) 
        - update-fixed-context (1. Infers new usefull knowledge from ./.agent/persistant/knowledge/mistakes.jsonl and ./.agent/flow/completed-issues/completed_issue_flows; 2. Add this new usefull knowledge to AGENTS.md if its not already there)
        - make-core-implementation-tasks-plan (transforms the spec into a graph of tasks, where: (1) each task can depend on other tasks being already done or not depend on any; (2) the tasks should not be of the form "one task implements each component that will be needed in this feature, e.g., oen task for the backend, another for databas,e another for gateway and another for frontend", the tasks should be done in the form of "one task implments a slice (governed by on or more integration/end-to-end tests) of multiple components, e.g., this task implements a funcitonality slice of frontend, backend, gateway and frontend that together brings us one step closer to our end goal and guarantee rich cross-component feedback along development". The core implementation tasks plan needs to be stored in .agent/flow/current-issue/flow/core-implementation-tasks-plan.md. The core-implementation-tasks-plan.md. file should have the graph structure of the plan, where each node is a task. For each node there is also a pending/in-progress/done checkbox. Do not confuse with plan.md which is a per-step ephmeral small plan for step execution.)
        - document using openwiki and with the following guidelines in ./.openwiki/INSTRUCTIONS.md: should first check incosistency (global spec is king), staleness and incompleteness of existing documentation (if any) and then update/create: (1) Contributor Documentation: visualization o repositoty strcture explaining succintly each directory and file, (1.2) Step by step contributor tutorial to help a newcomer understand the codebase; (2) User Documentation (only do after first version 0.1.0 is released): (2.1) User API Reference. (2.2) User step by step tutorial starting from sratch; (2.3) User guides to do common stuff; (2.4) FAQ. 
        make shure the documentaiton explains well things that I usually have a hard-time understanding.
        - ui-taste (UI Taste gives Claude a visual sense of taste. Instead of relying only on abstract design principles, the skill provides curated examples of bad, good, and stellar GUIs across different application categories and problem modes, including screenshots and their underlying HTML/CSS. This gives the agent an understanding of what makes GUIs look good. The agent should launch the current GUI, identify the biggest visual shortcomings, and iteratively improve them. The goal isn't to force a particular design style—it is to help Claude distinguish "functional but mediocre" from "genuinely beatifull and easy to use", giving coding agents a practical visual benchmark for judging their own work.)
    - Use existing skills: ctx7, unslop, visual-recap + explain-diff (ideally also running in github actions for every PR), skill-creator, i-have-adhd, quick-recap, mutation-testing, chrome-devtools-cli, optimize-anything, grillme, lavish-axi, code-review-and-quality, api-and-interface-design, browser-testing-with-devtools (only when working with frontend part), security-and-hardening, cc-skills-golang, maintainable-typescript (only when working with frontend part), improve-codebase-architecture, screenshot (only when working with frontend part), extract-design-system (only when working with frontend part), frontend-design (only when working with frontend part), teach, ponytail

## How Code Review is Done [for Humans or Outer Loop Agents]

1. Start a new session with opencode: do code review with Model A and store the review in .agent/session-persistent-candidate/reviews/[A]_code_review_[timestamp].md, where A is a placeholder for the actual mode name, and timestamp is placeholder for the actual timestamp
2. Start a new session with opencode: do code review with Model B and store the review in .agent/session-persistent-candidate/reviews/[B]_code_review_[timestamp].md, where B is a placeholder for the actual mode name, and timestamp is placeholder for the actual timestamp
3. Start a new session with open-code-review using Model C, which must be different from Models A and B, explicitly telling it to inspect the candidate problems flagged inside `.agent/session-persistent-candidate/reviews/` and store the final review in `.agent/session-persistent-candidate/reviews/final_code_review_[timestamp].md`. Use a fresh context for every reviewer.

## Token Efficency Laws [for Humans or Outer Loop Agents]

- Avoid small actions → batch related small tasks into larger coherent work units (use a todo_buffer.md to store all todos and then batch them before prompting the harness)
- Don't derail the agent from its main goal → context spent on unrelated work is expensive and increases context pollution.
- Use a graph/codebase-understanding tool → avoid repeatedly spending LLM tokens rediscovering repository structure and relationships.
- Related big tasks in the same session, unrelated new stuff gets its new session
- Start a new session when context becomes sufficiently polluted → carrying a huge amount of irrelevant history can become more expensive than rebuilding a clean context.
- Avoid rambling/random studying with the agent, do all of this in ChatGPT/Gemini/Grok webpages
- Stay in the same session while the context is still useful → preserving cached/reusable context avoids paying to rebuild understanding.
- Keep the core implementation model consistent through the common scaffold and initial functional-test suite when continuity is useful; this does not apply to reviews or verifications, which must use distinct models under the review-independence rule above.
- Always mention files with @ for the agent to look/modify instead o letting the agent winder the repo for that file
- Use something like rtk to compact tool outputs
- Have a router + local model for doing simple stuff

## .Agent Directory Structure Guide Doc (./.agent/directory_structure.md)

## `.agent/` Directory Structure

The `.agent/` directory contains the AI agent's workflow state, persistent knowledge, session state, issue-specific process artifacts and runtime-created tools. It is an internal directory used by the coding agent and should not contain product source code.

### Directory Structure (./.agent/directory_structure.md file)

```text
.agent/
├── flow/
│   ├── current-issue/
│   │   ├── raw_github_issue.md
│   │   ├── flow/
│   │   │   ├── issue_flow_[i].md where i is the issue number
│   │   │   └── core-implementation-tasks-plan.md
│   │   └── issue-spec_[i]/ where i is the issue number
│   │       ├── prd.md
│   │       ├── architecture.md
│   │       └── tech_stack.md
│   ├── current-pr/
│   │   ├── pre-pr-checklist.md
│   │   └── pr.md
│   └── completed-issues/
│       ├── completed_issue_flows/
│       │   └── issue_flow_[i].md where i is the issue number
│       ├── completed_issue_specs/
│       │   └── issue-spec_[i]/
│       │       ├── prd.md
│       │       ├── architecture.md
│       │       └── tech_stack.md
│       └── completed_prs/
│           └── pr_[i].md where i is the issue number
│
├── persistant/
│   ├── knowledge/
│   │   ├── non_obvious_conjectures_and_facts.md
│   │   └── mistakes.jsonl
│   └── user-understanding/
│       ├── user-codebase-questions.jsonl
│       └── user-quizzes/
│           └── quiz[i].md
│
├── session/
│   ├── plan.md
│   ├── todos.md
│   ├── problem_stuck.md
│   └── debug-logs/
│
├── session-persistent-candidate/
│   ├── knowledge/
│   │   ├── non_obvious_conjectures.md
│   │   └── exploration_findings/
│   │       └── [i]_<name-of-exploration>_[timestamp].md
│   └── reviews/
│       ├── final_spec_review_[timestamp].md
│       ├── final_spec_sec_review_[timestamp].md
│       ├── final_code_review_[timestamp].md
│       └── final_code_spec_review_[timestamp].md
│
├── created_tools/
│
├── directory_structure.md
├── supervised_evals/
│   └── <mistake_name>_[mistake_id]/
│       ├── metadata.json
│       ├── task.md
│       └── expected.diff
└── harness_evals_results/
    └── [i]_<experiment_name>_[timestamp]/
```

#### `flow/`

Contains the active issue implementation flow and PR state.

* `current-issue/` — contains the durable state of the issue currently being implemented.
  * `raw_github_issue.md` — the original GitHub issue, preserved as the source of truth for the issue being worked on.
  * `flow/` — contains the durable implementation workflow state for the current issue.
    * `issue_flow_[i].md` — records sequential progress through the Issue Flow for issue `i`, including completed steps, timestamps, approvals, session summaries, and deviations or jumps between steps.
    * `core-implementation-tasks-plan.md` — contains the dependency graph of the core implementation tasks. Each task represents a coherent behavioral or functional slice and records dependencies and progress.
  * `issue-spec_[i]/` — contains the specification produced for issue `i`.
    * `prd.md` — product requirements and expected behavior.
    * `architecture.md` — architectural design and implementation boundaries.
    * `tech_stack.md` — technologies, dependencies and relevant technical choices.

* `current-pr/` — contains the current PR artifacts for the active issue.
  * `pre-pr-checklist.md` — mandatory checklist that must be satisfied before creating or updating the PR.
  * `pr.md` — the current generated PR description.

* `completed-issues/` — archives durable artifacts from issues that have been completed.
  * `completed_issue_flows/` — archived issue flow records.
  * `completed_issue_specs/` — archived issue-specific specifications.
  * `completed_prs/` — archived PR descriptions for completed issues.

The `completed-issues/` archive should preserve the complete historical implementation record: issue flow, issue specification and final PR.

#### `persistant/`

Contains durable project knowledge that should survive across sessions and remain useful to future agents.

* `knowledge/non_obvious_conjectures_and_facts.md` — durable, non-obvious facts or evidence-backed insights about the codebase that are useful to future agents.
* `knowledge/mistakes.jsonl` — records mistakes made by the agent that required user intervention or were otherwise confirmed as real mistakes.
* `user-understanding/user-codebase-questions.jsonl` — records codebase-related questions asked by the user and their answers, branch and commit context.
* `user-understanding/user-quizzes/` — contains the history of user-understanding grill sessions. Each `quiz[i].md` records the timestamp, commit, topic, questions, user answers, feedback and score.

#### `session/`

Contains temporary state for the current agent session.

* `plan.md` — the ephemeral step-by-step plan for the current slash-command execution.
* `todos.md` — the ephemeral todo list used to batch unrelated correctness/security/build/test issues.
* `problem_stuck.md` — the problem description used when the agent is stuck in a loop.
* `debug-logs/` — temporary debug logs.

Everything in `.agent/session/` is disposable. At the beginning of a new session, its contents are cleared. These files must not be treated as historical records of the project or issue.

### `session-persistent-candidate/`

Contains information discovered during the current session that may be useful beyond the session but has not yet been promoted to durable knowledge.

* `knowledge/non_obvious_conjectures.md` — inferences or assumptions made during the current session, including current evidence and risk if wrong.
* `knowledge/exploration_findings/` — findings produced by exploratory sub-agents. Each findings file should contain the exploration context, description, agent/model used, token/cost information, reason it ended, findings and conclusion.
* `reviews/` — review artifacts generated during the current workflow session. Every review stored here must contain initial metadata identifying the exact reviewed version.

At the beginning of a new session, useful information from this directory is reviewed and, when sufficiently validated, promoted into `.agent/persistant/knowledge/non_obvious_conjectures_and_facts.md`; then the candidate directory is cleared.

### `created_tools/`

Where all tools created by agents at runtime are stored.

### `directory_structure.md`

Documents the purpose and organization of `.agent/`. It should be updated whenever the directory structure or responsibilities of its files change.

### `supervised_evals/`

Contains the persistent supervised evaluation suite derived from confirmed entries in `./.agent/persistant/knowledge/mistakes.jsonl`. Each cached supervised eval contains the mistake ID, the exact mistake commit, the agent-facing task, and the expected correction diff. The expected diff is hidden from evaluated agents. Existing cached evals are reused on later `harness-eval` invocations; only previously unseen mistake entries are converted into new evals.

### `harness_evals_results/`

Contains historical results from completed harness evaluations. It is runtime evidence, not benchmark ground truth, and is excluded from future evaluation snapshots so previous results cannot leak into evaluated agents.

### Important distinction: session state vs. durable state

The `.agent/` directory deliberately separates temporary state from durable project state:

* `session/` contains information needed only to continue the current session.
* `session-persistent-candidate/` contains potentially reusable discoveries that have not yet been validated or promoted.
* `persistant/` contains validated, durable knowledge.
* `flow/current-issue/` contains active issue state.
* `flow/completed-issues/` contains historical issue state.

### General Rules

1. **Do not put source code in `.agent/`.** Product code belongs under `./src/`.
2. **Do not treat `.agent/session/` as persistent storage.** Its contents may be deleted when a session starts.
3. **Do not promote assumptions or inferences automatically.** An assumption or inference should only become persistent knowledge after it has been sufficiently verified.
4. **Do not silently modify historical issue records.** They are part of the project's implementation history.

> This is the end of the .agent/directory_structure.md file

## Custom PR Skill [for Inner Loop Agents]

Mandatory pre-pr checklist (ready to push and open/update a PR):

- global spec, issue-specific spec and tests are consistent with each other
- strict 90% minimum coverage of the entire codebase, with a coverage report as evidence
- unit and integration tests are present and passing; end-to-end tests are required when the project exposes an external API or GUI
- no big changes were made after the last independent code review that have not been re-reviewed
- no significant changes were made after the last security review that have not been re-reviewed
- spec, documentation, tests and code are consistent; contributor docs cover the codebase; user docs cover all exposed external functionality when applicable
- user understands the PR at the function/interface/file/directory level

### Making PRs

- write the PR to `.agent/flow/current-pr/pr.md`; if the file already exists, delete/recreate it so stale content is not preserved
- ensure `.agent/flow/current-pr/pre-pr-checklist.md` is satisfied before making/updating a PR
- follow the project's PR template
- highlight key architectural/design decisions that a reviewer might question, problems encountered, solutions and tradeoffs
- document what was tested and link to evidence. Standard test evidence should show the current commit hash, test suite result and strict 90%+ coverage. GUI work should include before/after video capture where appropriate.
- make a risk assessment of the PR (Low, Medium, High) based on change scope, type of change, test coverage and other relevant evidence

In the summary section, pick the smallest view that makes the key point clear.

- Show logic or an algorithm as pseudocode:

```text
on(save)
  if content is unchanged
    return cached result
  write new content
  return fresh result
```

- Show runtime control flow as a call tree:

```text
submitForm
  createSession
    persistPrompt
    launchAgent
  navigateToSession
```

- Show UI structure as a component tree, including state and module boundaries that matter:

```text
<SessionPage> (apps/example/src/routes/session.tsx)
  useSessionEvents()
  <SessionToolbar>
    <RunSkillButton> (packages/ui)
```

- Show file responsibility or a broad refactor as a shallow file tree:

```text
src/
├── commands/       # parses user actions
├── sessions/       # owns session state
└── transport/      # sends API requests
```

- Show component interaction, control flow, or data flow with Mermaid:

```mermaid
sequenceDiagram
    participant User
    participant UI
    participant Daemon
    User->>UI: choose command
    UI->>Daemon: send expanded prompt
    Daemon-->>UI: stream result
```

- Use `diff` when the point is what changes and the surrounding shape already exists. Match the diff shape to the topic.

For a component change:

```diff
 <SessionPage>
   useSessionEvents()
   <SessionToolbar>
+    <RunSkillButton />
   <SessionTimeline>
+    <SkillResultCard />
```

For a file-layout change:

```diff
 src/
 ├── commands/
+│   └── show-me.ts       # expands the slash command
 ├── sessions/
-└── transport.ts
+└── transport/
+    ├── client.ts
+    └── stream.ts
```

For a call-tree or call-stack change:

```diff
 submitForm
   createSession
     persistPrompt
+    expandSkillMention
     launchAgent
-  navigateToSession
+  navigateToSession
+    subscribeToEvents
```

For a state or control-flow change:

```diff
 on(save)
-  write content
+  if content is unchanged
+    return cached result
+  write new content
+  invalidate cache
```

- Show the whole block when most of it is new, when omitted context would hide ownership or order, or when the user needs a copyable target shape:

```ts
function expandSkill(command: string): string {
  const skillName = command.slice(1);
  return `use the ${skillName} skill`;
}
```

### Handling PR Comments and Reviews

- respond to PR comments by addressing the problems raised, while judging whether each comment is correct and explaining the position when it is not followed
- make appropriate changes when the review identifies a real problem
- when it is unclear what to change because feedback is subjective, inconsistent or ambiguous, ask the user to resolve the ambiguity

---

## Harness Evaluation Custom Skill (harness-eval.md) [for Inner Loop Agents]

Use this skill to run controlled A/B evaluations of AI coding setups.

## When to Use This Skill

Use this skill whenever empirical evidence is needed to determine whether a proposed AI coding setup change is better or worse than a baseline for the specific repository.

Typical uses:

* compare Model A vs Model B
* compare a harness setup A vs B
* compare a skill vs no skill
* compare Skill A vs Skill B
* evaluate harness changes currently present in the working tree
* evaluate a harness change before making it
* investigate whether an AGENTS.md rule, skill, hook, subagent, model, or other harness change helps on the actual repository

The evaluation supports these task types:

```text
coding
test-writing
spec / architecture / design
code review
doc review
```

The caller must explicitly specify which eval type(s) to run.

When `doc-review` is included, the caller must specify the **exact document path(s)** to review.

Do not invent missing eval types or document paths.

---

## 1. Determine the Two Options

The evaluation always compares:

```text
Option A
Option B
```

The options can be explicitly supplied or inferred from the working tree.

## 1.1 Explicit A/B Evaluation

When the user explicitly provides the two configurations, use them directly.

Examples:

```text
Model A vs Model B
Skill A vs Skill B
Harness setup A vs B
Harness A vs Harness B
```

For an explicitly defined A/B evaluation, **uncommitted working-tree changes are ignored** unless the user explicitly includes them in the option definition.

Before executing, print a terminal preflight summary showing:

```text
EVAL MODE: explicit A/B
Option A: <description>
Option B: <description>
Uncommitted working-tree changes: ignored unless explicitly included
```

Determine exactly what differs between the options.

## 1.2 Inferred A/B Evaluation From Current Working-Tree Harness Changes

When the user asks to evaluate the current harness changes without explicitly providing A and B, inspect the current working tree.

The intended comparison is:

```text
Option A = last commit
Option B = last commit + selected uncommitted harness-related changes
```

Uncommitted application changes are excluded from both options.

The harness should identify candidate harness-related changes using the current repository structure and changed files. Do **not** hard-code an exhaustive definition of "harness-related"; harness structure can change over time.

Before execution, print the exact detected files and directories proposed for B, for example:

```text
EVAL MODE: inferred from current working tree

Option A:
  last commit: <commit>
  uncommitted changes: excluded

Option B:
  last commit: <commit>
  uncommitted harness-related changes:
    AGENTS.md
    skills/foo/SKILL.md
    .agent/...

Uncommitted application changes excluded from both options:
  src/...

Please confirm the detected A/B change set before the experiment starts.
```

The user may add or remove files/directories from the detected change set. Do not start the experiment until the user confirms the final A/B change set.

The inferred experiment must use exactly:

```text
Option A = last committed repository state
Option B = last committed repository state + user-confirmed uncommitted harness changes
```

## 1.3 Infer Multiple Independent Evaluations

If the user asks to evaluate current harness changes generally, identify independent relevant changes and create a separate experiment for each when possible.

For example:

```text
Change 1:
new code-review skill

Change 2:
modified AGENTS.md

Change 3:
model configuration change
```

Construct separate evaluations when the changes can be meaningfully isolated.

Do not combine independent changes into a single experiment.

If one change depends on another and they cannot be meaningfully separated, evaluate them together and explicitly record the combined experimental variable.

---

# 2. Required Eval Input

Determine:

```text
EVAL_TYPES =
  one or more of:
    coding
    test-writing
    spec
    review
    doc-review
```

Also determine:

```text
DOCS =
  exact document paths
  required if doc-review is included
```

Option configuration may be explicitly supplied or inferred from current changes.

If the eval type is missing, stop and request it.

If `doc-review` is selected without exact document paths, stop and request them.

The caller-selected `EVAL_TYPES` governs reconstruction evals only. The cached supervised eval suite described below is always run in addition to those selected reconstruction evals.

---

# 2.1 Required A/B Trial Count

Every A/B evaluation must begin with **3 independent paired screening trials** because model behavior is probabilistic. Three is the minimum screening sample, not a hard cap.

One trial consists of one independent execution of Option A and one independent execution of Option B against the same frozen mutation and task for that trial.

The initial screening therefore produces:

```text
3 paired trials
×
2 options
=
6 agent executions
```

Before launching an experiment, declare the primary quality metric, the smallest practically meaningful A/B difference for that metric, and the maximum confirmation budget. After the initial screening, run additional paired confirmation trials when the results are close or unstable. Treat the result as close/unstable if (a) the direction of the paired A-minus-B difference changes across the three screening trials, (b) the aggregate difference is smaller than the predeclared practically meaningful difference, or (c) a paired confidence interval still includes both no difference and a practically meaningful advantage for either option. Add confirmation trials in batches of 3, keeping A and B matched within every trial and using independent fresh mutations for reconstruction trials. Recalculate the paired effect and interval after each batch. If the predeclared trial budget is exhausted before the result is stable, report the comparison as inconclusive rather than declaring a winner.

Each trial must use fresh agent sessions and isolated execution environments. Trials must not reuse agent state, writable caches, repositories, or other execution state from previous trials.

Every screening and confirmation trial must be independently recorded. Never stop before the initial 3 paired trials are complete. Do not declare a winner after the initial screening when the results are close or unstable; run the confirmation batches described above.

Each reconstruction trial must use an independent mutation seed and therefore a fresh mutation, while A and B within that trial receive the exact same frozen mutation.

For supervised evals, there is no generated mutation to reseed. All screening and confirmation trials reuse the same cached historical mistake state, task and expected diff, while still using fresh agent sessions, isolated containers, repositories and writable execution state.

---

# 2.2 Automatic Cached Supervised Eval Suite

In addition to the caller-selected reconstruction eval types, **every `harness-eval` invocation must run the complete cached supervised eval suite** for each A/B experiment. Supervised evals are not an `EVAL_TYPES` value and do not require the caller to request them explicitly.

A supervised eval is a replay of a confirmed real mistake recorded in `./.agent/persistant/knowledge/mistakes.jsonl`. Each supervised eval contains exactly these benchmark ingredients:

```text
task
codebase commit at the time of the mistake
expected correction diff
```

The expected diff is hidden from the evaluated agents. The main evaluation agent receives the expected diff only after the agent execution and compares the agent-produced diff against it to produce the supervised evaluation score.

Supervised evals run with the same isolation requirements as reconstruction evals:

* Option A and Option B each run in their own Docker Compose environment.
* Each supervised eval gets fresh repositories, containers, agent sessions, HOME directories, harness state, temporary directories, writable caches, and build directories.
* A supervised eval must not share writable execution state with another supervised eval or with a reconstruction eval.
* A and B receive the exact same historical mistake state and task within each trial; only the intended A/B harness difference may vary.
* The historical mistake commit is the application/codebase baseline. Apply the current Option A or Option B harness configuration on top of that historical baseline without importing later application changes.

Every cached supervised eval is executed on every invocation. Existing cached eval definitions are reused; they are not regenerated merely because the repository or current harness has changed.

---

# 2.3 Construct and Cache Supervised Evals From New Mistakes

At the beginning of every `harness-eval` invocation, synchronize the supervised-eval cache with `./.agent/persistant/knowledge/mistakes.jsonl`. Only mistake entries that are not already represented in `./.agent/supervised_evals/` are transformed.

Use a stable `mistake_id` derived from the canonical JSON representation of the complete mistake entry, for example:

```text
mistake_id = SHA-256(canonical_json(mistake_entry))
```

Do not use the JSONL line number as the identity because entries may be moved or reformatted. The cache directory itself is the lookup/index: if `./.agent/supervised_evals/<mistake-id>/` exists, that mistake entry is already cached. No separate cache index is required.

### Construction must happen outside the main repository

The construction process is **not allowed to create a Git repository, checkout, worktree, or other mutable historical-code workspace anywhere inside the main repository**, including under `.agent/supervised_evals/`. The main repository may only be read to obtain source material and benchmark-construction inputs. Any filtering, redaction, deletion, truncation, or other sanitization must happen only in a copied external workspace; never in the main repository itself.

The **project folder** is the Git checkout containing the project's `.git` directory. The agent may work outside the project folder only **one level above the project folder**. Any such path must be expressed relative to the project root using exactly one `../` component. The agent must not use absolute paths for external workspaces and must not traverse two or more levels above the project root (for example, `../../...`).

All construction work must happen in an external temporary workspace such as:

```text
../tmp/harness-eval/<run-id>/supervised-eval-construction/<mistake-id>/
```

The construction workspace may contain a temporary clone/checkout with the original Git history because construction is trusted benchmark-building work. This historical repository must never be exposed to the evaluated agent. Only after construction is complete may the final `task.md`, `expected.diff`, and minimal cache metadata be copied into `./.agent/supervised_evals/<mistake-id>/`.

This separation is mandatory because the supervised-eval agent must start from a historical codebase state, while putting a Git repository inside the user's working repository creates a nested-repository and accidental-scope risk.

For each new mistake entry:

1. Parse and validate the required mistake fields, especially the commit hash of the codebase at the time of the mistake.
2. Create the external construction workspace and reconstruct the codebase at exactly that mistake commit there.
3. Use the recorded mistake description and the trusted repository history available in the construction workspace to identify the confirmed correction that actually resolved the mistake.
4. Construct the **expected diff** as the smallest semantically complete correction patch from the mistake-commit state to the confirmed corrected state. The expected diff must be grounded in the real correction; do not invent a hypothetical patch.
5. Construct the agent-facing **task** from the mistake description and evidence available to the agent at that historical state. The task may describe the problem that needs fixing, but must not reveal the logged correction procedure, hidden expected diff, or post-mistake ground truth.
6. Validate that applying the expected diff to the mistake-commit state produces the confirmed correction.
7. Only after the construction is validated, copy the completed eval into `./.agent/supervised_evals/<mistake-id>/`.

The cached supervised eval contains only:

```text
metadata.json
  mistake_id
  mistake_commit

task.md
expected.diff
```

The two metadata fields above have direct operational purposes:

* `mistake_id` identifies which `mistakes.jsonl` entry the cached eval represents and is also encoded by the cache directory name.
* `mistake_commit` tells the harness which historical application state to reconstruct for execution.

Do **not** store redundant construction metadata such as task paths, expected-diff paths, baseline tree identifiers, expected-diff hashes, construction timestamps, or correction-source references. The paths are fixed by the cache layout, the mistake commit already identifies the historical baseline, the expected diff itself is the ground truth, and construction bookkeeping is irrelevant once the cached benchmark is validated.

Do not expose `./.agent/supervised_evals/`, `expected.diff`, or any construction workspace/evidence to evaluated agents. These paths must be excluded from every evaluated repository snapshot.

A cached supervised eval is immutable during normal `harness-eval` invocations. If a mistake entry already has a cached eval, reuse that eval rather than silently rebuilding it. An intentional rebuild must be an explicit maintenance action outside the normal evaluation flow.

If a new mistake cannot be converted into a trustworthy expected diff, do not invent one. Do not create a partial cache entry; report the construction failure for that invocation and leave the mistake uncached so it can be retried on a later invocation.

---

# 3. FIRST: Isolate the Evaluation

**Before running any evaluated agent, create the isolated evaluation environment.**

The user's current repository is the **main checkout**.

Treat it as read-only for evaluated execution. The harness may update the persistent supervised-eval cache before the experiment and persist results after the experiment, but evaluated agents must never execute against the main checkout.

Do not:

* create a reconstruction mutation in the main checkout
* install an evaluated harness in the main checkout
* run an evaluated agent in the main checkout
* modify harness files in the main checkout as part of the actual experiment
* create runtime evaluation artifacts in the main checkout while agents are running
* modify `.agent/supervised_evals/` while agents are running
* modify `.agent/harness_evals_results/` while agents are running

The evaluation uses two separate Docker Compose environments, one per option. The main checkout is only the source from which the evaluation snapshots and historical supervised-eval repositories are created.

---

# 4. Snapshot the Correct Source State

The snapshot semantics depend on the evaluation mode.

## Explicit A/B

Create each option directly from the explicitly defined configuration. Ignore unrelated uncommitted working-tree changes unless the user explicitly included them.

## Inferred A/B From Current Harness Changes

Create the source states as follows:

```text
Option A source = HEAD / last commit
Option B source = HEAD / last commit + user-confirmed uncommitted harness changes
```

Do not include uncommitted application changes in either option.

For both modes:

Exclude from every agent-visible evaluation snapshot:

```text
.git/
.agent/persistant/knowledge/mistakes.jsonl
.agent/supervised_evals/
.agent/harness_evals_results/
temporary evaluation state
temporary caches
credentials
secrets
```

**Excluding a file means omitting it from the copied evaluation snapshot/container. Never delete, truncate, modify, or otherwise change the corresponding file in the main repository.** In particular, the main repository's `./.agent/persistant/knowledge/mistakes.jsonl` must remain untouched because it is the persistent source for future supervised-eval construction.

Never copy secret values into the evaluation repositories.

Do not extract secrets from A/B experiments into metadata, traces, fingerprints or reports. Secret values are not relevant to the decision-making metrics.

Record:

```text
source commit
working-tree state / explicit option definition
included files
excluded files
snapshot tree identifier
```

---

# 5. Create Isolated Evaluation Repositories

Do not create evaluated repositories using normal worktrees from the original repository.

All temporary evaluated repositories must be created under a sibling directory one level above the project folder, using a path relative to the project root such as `../tmp/harness-eval/<run-id>/...`. Do not use an absolute filesystem path or traverse above that parent directory.

The original repository's Git history may reveal hidden ground truth.

Instead:

1. create the agent-visible sanitized snapshot described in Section 5.1
2. copy that sanitized snapshot into the external one-level-up workspace
3. create a new temporary Git repository in the external one-level-up workspace
4. initialize it
5. add the snapshot
6. create a single root commit

The resulting repository must contain the exact intended snapshot but none of the original Git history.

## 5.1 Sanitize Everything Before Creating Agent Containers

Before either A/B repository is copied into a Docker container, create an **agent-visible sanitized snapshot**. The sanitized snapshot is a derived copy used only for the evaluation; the main repository is never modified as part of sanitization.

The sanitization process must inspect **all files and directories that would otherwise be copied into the agent-visible repository**, not only `.agent/`. The goal is to prevent the evaluated agent from recovering information about the mistake, expected correction, prior agent behavior, hidden mutation, or evaluation outcome from repository context.

At minimum, scan for and remove or redact information such as:

```text
confirmed mistakes and mistake descriptions
expected corrections or expected diffs
post-mistake corrected code that reveals the solution
agent transcripts, tool traces, or previous agent reasoning
harness-evaluation results or rankings
mutation metadata or mutation patches
supervised-eval task/ground-truth construction evidence
prior reviews that reveal the specific defect or correction
issue-flow/session history that reveals the solution, mistake, or hidden benchmark state
non-obvious conjectures/facts that reveal hidden benchmark information
created tools or their artifacts that reveal the benchmark construction or expected solution
user quizzes or codebase-question answers that directly reveal the hidden benchmark information
other generated or historical artifacts whose contents could reveal what the agent is supposed to reconstruct or fix
```

Do not assume a file is safe merely because its path is not obviously related to evaluations. **Read and assess every copied text-bearing artifact for benchmark leakage.** Binary/generated artifacts must also be excluded when they contain hidden ground truth or other benchmark-specific information that an evaluated agent could inspect.

If a file contains useful repository context plus benchmark-specific information, sanitize only the benchmark-leaking portions while preserving the unrelated useful context. If the file cannot be safely sanitized without risking leakage, exclude it from the agent-visible snapshot.

The following two directories receive special treatment:

```text
.agent/session/
.agent/session-persistent-candidate/
```

In every A/B agent container, preserve only their directory/file structure. Every file inside them must be **zero bytes**, and no file may contain historical, task, mistake, review, plan, transcript, conjecture, or other content. Directories may remain present, but they must contain only these empty files/directories.

The final agent-visible snapshot must therefore satisfy:

```text
main repository
      │
      ├── remains completely unchanged
      │
      └── sanitized copy
              ↓
       Option A / Option B containers
```

The main repository's `mistakes.jsonl`, supervised-eval cache, historical results, session state, and session-persistent-candidate contents are never deleted or sanitized in place. Only the snapshot is filtered.

The evaluated agents must not be able to recover hidden ground truth with commands such as:

```bash
git show <original-commit>
git diff <original-history>
git log <original-history>
```

---

# 6. Determine Validation

From the isolated repository:

1. identify the normal test/build commands
2. identify repository-specific validation
3. run the known-good baseline validation
4. confirm the baseline passes

Record the exact commands.

Do not invent a new correctness oracle unless the selected eval definition requires it.

---

# 7. Create One Mutation Per Trial

For each reconstruction trial, create **one mutation**.

The mutation is shared by both options.

Never independently generate a mutation for A and B.

For supervised evals, **do not create an additional mutation**. The cached historical mistake state is the fixed benchmark input; use the cached mistake commit and expected diff instead.

Use:

```text
known-good evaluation repository
        ↓
create mutation once
        ↓
validate mutation
        ↓
freeze mutation
        ↓
clone exact frozen state
      ↙       ↘
 Option A   Option B
```

Record:

```text
mutation_id
seed
eval_type
target path(s)
target region(s)
mutation description
ground-truth content or canonical ground-truth reference
known-good tree identifier
mutated tree identifier
```

---

# 8. Mutation Rules

Mutations must not alter application dependencies or evaluation infrastructure unless the selected evaluation definition explicitly calls for it.

## Application Dependencies

Application dependencies must be identical in A and B.

Do not use mutations to:

* add application dependencies
* remove application dependencies
* upgrade application dependencies
* downgrade application dependencies
* modify application lockfiles
* modify application package manifests
* modify vendored application dependencies

If a mutation requires an application dependency change, reject it.

## Harness Dependencies

Harness-level dependencies are different: a harness change being evaluated may intentionally add, remove, upgrade or downgrade a **harness-level dependency**. Such a dependency difference is valid when it is the intended experimental variable.

The harness must record harness-level dependency differences as part of the experimental difference and must still keep application dependencies identical.

## Protected Paths

Protect at minimum:

```text
tests/
.github/
CI/CD infrastructure
build infrastructure
application dependency configuration
application lockfiles
pre-commit configuration
evaluation infrastructure
harness-evaluation infrastructure
```

Add repository-specific protected paths as required.

Protected-path rules are evaluation-specific: a path may be part of the intended harness A/B difference when it is explicitly classified as harness configuration, while application tests and application correctness artifacts remain protected during mutations.

---

# 9. Coding Eval

For a coding eval, remove selected application implementation logic while preserving the existing tests.

Eligible targets include:

* functions
* methods
* classes
* services
* controllers
* handlers
* algorithms
* business logic
* orchestration logic
* state transitions
* modules
* cross-module application logic

Do not delete:

* tests
* application dependencies
* application dependency configuration
* application lockfiles
* CI/CD
* build infrastructure
* evaluation infrastructure
* harness infrastructure
* generated/vendor code
* comments only
* formatting only
* trivial imports

Prefer AST/semantic region selection.

The selected missing implementation must be sufficiently represented by existing tests and surrounding repository evidence for the coding reconstruction task to be meaningful.

The coding mutation is valid when:

1. the known-good repository passes its relevant validation
2. the mutation creates the intended missing implementation
3. the ground truth can restore the original implementation

The mutated repository does **not** need to fail deterministic validation. The coding agent's task is to reconstruct the missing implementation from repository evidence.

---

# 10. Test-writing Eval

The test-writing eval measures the ability to practice test-first development by reconstructing **tests plus only the minimal project skeleton/stubs required for those tests to run**, using the repository's specifications and surrounding technical context.

Start from a known-good repository state in which its last merged PR is the one we want to recostruct the skeleton/tests of the selected new behaviour being merged.

For the mutation:

* define commit hash of the repo before the merge of the PR we want to recostruct, define the PR number we want to recostruct the tests, define the issue number this PR resolves
* setup the repo to be exactly as the commit hash of the repo before the PR
* put the issue description in the raw issue directory inside .agent

Note: when running this eval need to make shure the agent doesnt have access to github, because if it has access it can cheat by looking at the actual PR made or things in the future.

The agent must:

* recreate/modify tests from the issue description
* recreate/modify the project skeleton/stubs required for those tests to compile and run
* use the issue description as the source of behavioral intent
* **not implement the actual application behavior**

All reconstructed tests must be runnable and all of them must fail for the expected reason: the behavior under test is intentionally not implemented. Test failures caused by invalid test code, compile errors, missing imports, broken test setup, or unrelated infrastructure problems are not acceptable.

The evaluation therefore has two distinct requirements:

```text
test validity (determisistic):
    all tests compile/run
    all tests fail because the intended behavior is unimplemented

semantic test quality (agentic evaluation):
    tests correctly represent the requirements in the issue
    and cover the intended behavior, constraints, and important cases
```

The actual merged PR on github is used by the agentic evaluation as a baseline for comparisons of generate skeleton/tsts vs actual merged skeleton/tests.

Do not require the agent to reconstruct the real implementation.

---

# 11. Spec / Architecture / Design Eval

Choose an existing technical specification, architecture or design document.

Remove a meaningful region.

Possible targets:

* architecture documentation
* design documents
* ADRs
* subsystem descriptions
* component descriptions
* interface descriptions
* data flows
* deployment architecture
* technical rationale

Do not modify application code or tests.

The original removed content is hidden ground truth.

Give both agents the incomplete document.

After they finish, a main evaluation agent performs semantic comparison between the reconstructed content and the original hidden content, using the surrounding document and repository as context.

Literal textual identity is not required when the candidate preserves the intended technical meaning and required details.

---

# 12. Code Review Eval

Start from the known-good application code and inject one or more concrete defects.

Code-review defects must cover the full range of issues that a strong code review should detect, not only functional bugs.

Defect categories should include, where applicable:

* functional/behavioral bugs
* security vulnerabilities
* incorrect error handling
* state or concurrency bugs
* defects that make the implementation fail to conform to the specification
* architectural/design violations grounded in repository/spec requirements
* code-quality and maintainability problems
* functions that are too large
* classes that are too large
* files that are too large
* misleading or poor names
* inconsistent conventions
* violations of established project best practices
* problematic duplication, coupling or structure

A valid injected defect must be plausible and grounded in the known-good repository, specification, architecture, design system or established repository conventions.

Across the full evaluation suite, include all applicable defect categories listed above. A single trial does not need to contain every category.

Avoid purely subjective style preferences. For code-quality defects, provide concrete repository/spec evidence in the hidden mutation metadata explaining why the original implementation is preferable.

A defect does **not** need to:

* be covered by tests
* cause a test failure
* cause deterministic validation failure
* change runtime behavior

Some defects should be deliberately test-detectable, while others should require reasoning about security, specification conformance, maintainability or code quality.

Do not modify tests to create code-review defects.

Do not tell the agent:

* how many defects exist
* where they are
* which files changed
* what defect categories were injected

The mutation metadata must preserve for each defect:

```text
defect_id
category
target region
original/correct content or canonical reference
defect description
repository/spec evidence
```

---

# 13. Doc Review Eval

The caller must provide the exact document path(s).

For each specified document:

1. verify it exists
2. read the document
3. inject concrete repository-grounded defects
4. preserve the original document as hidden ground truth
5. freeze the mutated document

Possible defects:

* incorrect technical statement
* incorrect API description
* outdated command
* incorrect configuration
* contradiction with implementation
* incorrect example
* broken reference
* missing requirement
* incorrect architecture statement

Avoid subjective style-only issues.

Do not tell the agent what was changed.

After the agent finishes, compare the repaired document against the hidden original using semantic evaluation from the main evaluation agent.

---

# 14. Freeze the Mutation

After mutation creation:

1. run the required baseline validation
2. confirm the intended mutation is present
3. confirm application dependency files are unchanged
4. confirm protected evaluation files are unchanged
5. compute the mutated tree identifier
6. freeze the mutation

The mutation must never change during the A/B comparison.

---

# 15. Create Exactly Two Docker Compose Environments

Each A/B trial uses exactly two Docker Compose environments, one for each option. Each Compose environment contains exactly one primary agent container. Additional service containers are not permitted unless explicitly required by the repository's evaluation environment.

```text
Option A
  └── Docker Compose environment A
       └── primary agent container A

Option B
  └── Docker Compose environment B
       └── primary agent container B
```

The two Compose environments must be independent.

Each primary agent container must have its own **sanitized** copy of the evaluation repository and its own:

```text
filesystem
repository
HOME
agent state
harness state
temporary directory
writable cache
build directory
```

The repository copied into these containers is the sanitized agent-visible snapshot. It is never the live main repository and never contains `.agent/persistant/knowledge/mistakes.jsonl`.

Do not run both options in the same container.

Do not mount one writable repository into both environments.

Do not share writable agent state or caches.

Use the same Docker image for A and B whenever the intended experimental difference does not require otherwise. Record the immutable Docker image digest.

If the intended A/B variable itself changes the agent runtime image or harness dependency environment, record that as the intended experimental difference rather than treating it as an unexpected environment mismatch.

---

# 16. Create Two Independent Repositories

Create:

```text
option-a/repo
option-b/repo
```

Initialize both from the exact same frozen mutation tree before applying the intended A/B harness difference.

Before applying the A/B difference:

```text
tree_identifier(A) == tree_identifier(B)
```

must be true.

Also verify:

```text
mutation_id(A) == mutation_id(B)
mutation_seed(A) == mutation_seed(B)
application_dependencies(A) == application_dependencies(B)
tests(A) == tests(B)
```

Harness-level dependency differences are allowed only when they are part of the intended A/B difference.

Do not start the agents until these checks pass.

---

# 17. Configure Option A and Option B

Apply only the intended experimental difference.

Examples:

### Model

```text
A = Model A
B = Model B
```

### Skill

```text
A = no skill
B = skill
```

or:

```text
A = Skill A
B = Skill B
```

### Setup

```text
A = original setup
B = modified setup
```

### Harness

```text
A = Harness A
B = Harness B
```

Keep everything else identical whenever possible.

If an unavoidable difference exists, record it explicitly.

For inferred working-tree harness evaluations, this is the point where the user-confirmed harness file set from Section 1.2 is applied to B.

For a supervised eval, first reconstruct the historical application/codebase state at the cached `mistake_commit`, then apply the current Option A harness configuration to the A repository and the current Option B harness configuration to the B repository. Do not import application changes made after the mistake commit.

---

# 18. Verify the Environments Before Starting

Immediately before sending the task to either agent, compute an environment fingerprint for each option.

Verify equivalence of all non-experimental execution properties, including:

```text
Docker image digest
OS/runtime versions
compiler versions
package-manager versions
application dependency versions
tool versions
CPU allocation
memory allocation
network policy
non-sensitive environment variables
filesystem permissions
working directory
repository tree identifier
mutation identity
tests
validation commands
task text
```

Do not collect or expose secret values from environment variables.

The expected result is:

```text
A == B
+
only intended experimental difference(s)
```

Expected differences may include harness configuration, harness-level dependencies, model configuration, skills, prompts, commands or other explicitly intended harness changes.

If an unexpected difference exists, **do not start the agents**. Fix the environments first.

This check is mandatory.

---

# 19. Give Both Agents the Same Task

Use the same task text for A and B for the selected eval type.

Do not expose:

* mutation patch
* deleted original code
* original document content
* defect locations
* number of defects
* mutation generator
* ground truth
* other option's results
* other option's trace

### Coding

> Some application logic in this repository is missing. Reconstruct the missing implementation using the surrounding code, interfaces, types, documentation, and tests as evidence.
>
> Do not modify tests, CI/CD, build infrastructure, application dependency configuration, application lockfiles, pre-commit configuration, evaluation infrastructure, or unrelated code.

### Test-writing

> Some application behavior and its tests are missing from this repository. Reconstruct the tests and only the minimal project skeleton/stubs required for those tests to run, using the specifications, interfaces, types, documentation and surrounding repository as evidence.
>
> Do not implement the actual application behavior. All reconstructed tests must compile/run and must fail because the behavior they test is not implemented.
>
> Do not modify unrelated application code, dependency configuration, CI/CD, build infrastructure, pre-commit configuration, evaluation infrastructure, or unrelated files.

### Spec / Architecture / Design

> Some technical specification, architecture, or design documentation in this repository is incomplete. Reconstruct the missing section using the surrounding documentation and repository as evidence.
>
> Do not modify application code, tests, CI/CD, build infrastructure, dependency configuration, pre-commit configuration, evaluation infrastructure, or unrelated files.

### Code Review

> Review this repository carefully and identify the introduced defects. For every defect you address, explain the problem and make the correct fix. Defects may involve functionality, security, specification conformance, architecture/design, maintainability, code quality, naming, conventions, or other concrete engineering problems.
>
> Do not modify tests, CI/CD, build infrastructure, application dependency configuration, application lockfiles, pre-commit configuration, evaluation infrastructure, or unrelated code.

### Doc Review

> Review the specified document against the repository and identify concrete correctness, consistency, and completeness issues. Fix the issues you find using the implementation and other repository documentation as evidence.
>
> Do not modify unrelated documents, application code, tests, CI/CD, build infrastructure, dependency configuration, pre-commit configuration, evaluation infrastructure, or unrelated files.

### Supervised Eval

Use the cached task for the supervised eval. The repository starts at the exact historical commit recorded for that mistake. Ask the agent to solve the task normally. Do not expose the expected correction diff, the post-correction repository state, the mistake's correction procedure, or the cache metadata.

The task given to the agent must be identical for Option A and Option B within a trial.

---

# 20. Start Fresh Agent Sessions

Each option gets:

* its own Docker Compose environment
* its own repository
* fresh agent state
* fresh session
* isolated HOME/config/cache

Do not reuse sessions.

Do not allow the agents to communicate.

The agent may always run the repository's tests when the selected task permits it. There is no separate read-only code-review run.

---

# 21. Measure Agent Task Time

**Agent time means only the time the agent spends executing the task.**

Start the timer immediately before giving the task to the agent.

Stop the timer when the harness execution reports that the task is finished or when the total agent-harness execution reaches the configured timeout limit.

Do not include:

* Docker startup
* repository creation
* dependency installation
* harness installation
* mutation generation
* preflight checks
* post-agent validation

The timeout applies to the **total agent harness execution time for the task**, including model reasoning, tool calls and sub-agents used during that execution.

Record:

```text
agent_started_at
agent_finished_at
agent_duration_seconds
timeout_seconds
```

Record setup and validation time separately if useful.

---

# 22. Record Token Usage

After each agent finishes, retrieve actual usage.

Record:

```text
model
input tokens
cached input tokens
cache-write tokens
output tokens
reasoning tokens
other billable categories
```

Record usage per model request when the harness may use:

* planner models
* coding models
* review models
* subagent models
* fallback models

Record the effective model actually used.

Model names in this document may remain descriptive (for example, "best coding quality/price model") because model choices change frequently. The exact effective model identifier must still be recorded in experiment results.

---

# 23. Calculate Model Inference Cost

For every model actually used:

1. identify the exact model
2. retrieve the current applicable pricing
3. record the authoritative pricing source
4. record the pricing retrieval time
5. normalize the provider's actual billable usage categories
6. calculate estimated model inference cost

Do not assume cached tokens are excluded from `input_tokens`.

Do not double-count tokens.

If the provider reports an actual billable amount, prefer it.

Otherwise calculate:

```text
estimated_model_inference_cost
=
actual billable token usage
×
applicable current token prices
```

Record:

```text
pricing_source
pricing_retrieved_at
model
pricing_tier
currency
token prices
billable categories
estimated_model_inference_cost
```

If current pricing cannot be verified:

```text
cost_status = UNKNOWN
```

Do not invent a price.

---

# 24. Run Identical Validation

After each agent finishes:

1. inspect changed files
2. check protected paths
3. check application dependency files
4. restore protected files if necessary
5. run the exact same validation commands for both options

Run:

```text
build
tests
repository-specific validation
```

Also run when applicable:

```text
lint
typecheck
integration tests
documentation validation
schema validation
```

Record:

```text
validation_started_at
validation_finished_at
validation_duration_seconds
validation_output
validation_exit_status
```

Keep validation time separate from agent time.

For the test-writing eval, validation should specifically verify that all reconstructed tests are runnable and all of them fail for the expected unimplemented-behavior reason.

---

# 25. Check Protected Files and Dependencies

An agent must not redefine the evaluation.

After execution, check for changes to:

```text
tests
CI/CD
build configuration
application dependency manifests
application lockfiles
evaluation scripts
evaluation configuration
ground-truth artifacts
```

Harness-level files and dependencies may legitimately differ between A and B when they are the intended experimental variable.

If the agent modified protected files:

1. record the violation
2. restore them from the frozen mutation / intended option baseline
3. exclude those changes from the solution
4. rerun validation against the protected baseline
5. report the violation

For the test-writing eval, tests and the minimal project skeleton are intentionally writable and are therefore not protected for that eval. Application behavior implementation remains out of scope for the agent.

---

# 26. Evaluate Correctness

### Coding

Check whether the reconstructed implementation satisfies the required tests and repository validation.

### Test-writing

Check:

```text
all reconstructed tests run
all reconstructed tests fail for the expected missing-behavior reason
project skeleton/stubs are sufficient for the tests to run
no actual application behavior was implemented
```

Then use the main evaluation agent to semantically compare the reconstructed tests and skeleton against the specifications and surrounding repository evidence.

### Code Review

A code-review defect is considered **found** only when the agent both:

1. presents/identifies the defect; and
2. makes a fix that semantically matches the correct behavior/content in the original known-good repository.

The main evaluation agent performs the semantic comparison between:

```text
agent's finding + agent's fix
                vs
original hidden repository + specification/architecture/conventions
```

The primary metrics are precision, recall and F1. Count a finding at most once per injected defect; duplicate reports of the same defect do not increase true positives. A finding is a true positive only when the evaluator confirms that it matches an injected defect and the agent both identifies the defect and makes a semantically correct fix. An injected defect not found is a false negative. A reported finding that is not a real defect is a false positive.

```text
precision = true_positives / (true_positives + false_positives)
recall = true_positives / (true_positives + false_negatives)
f1_score = 2 * precision * recall / (precision + recall)
detection_score = recall  # retained as a backward-compatible alias
```

Handle zero denominators as `N/A`, not as an invented zero score. Report at minimum:

```text
defects_injected
true_positives
false_positives
false_negatives
precision
recall
f1_score
detection_score
```

A valid additional defect that was not injected should be reported separately as `validated_additional_findings`; do not count it as a false positive merely because it was outside the injected benchmark set.

Also report the categories of the injected defects.

A defect that was only discussed but not fixed does not count as found.

A defect that was fixed without being identified/presented to the user does not count as found.

### Spec / Architecture / Design

Compare the reconstructed region with the original hidden region using semantic evaluation by the main evaluation agent. Preserve intended meaning and required technical details; literal wording need not match.

### Doc Review

Compare the repaired specified document with the original hidden document using semantic evaluation by the main evaluation agent.

### Supervised Eval

The main evaluation agent compares the agent-produced final diff with the cached `expected.diff`. Score semantic correctness rather than literal patch identity so equivalent implementations can receive full credit. The comparison must determine whether the mistake was completely corrected, whether the changes are unnecessarily broad or introduce regressions, and how closely the resulting behavior/content matches the confirmed correction.

Use a normalized `supervised_diff_score` in the range `0.0` to `1.0`, where `1.0` means the intended correction is semantically complete and the diff contains no material unjustified changes. Partial credit is allowed for partial corrections. A score of `0.0` means the proposed change does not materially correct the mistake.

Record both the produced diff and the hidden expected diff as separate ground-truth/evidence artifacts. Never expose the expected diff to the evaluated agent.

---

# 27. Handle Invalid Runs Explicitly

Use explicit statuses:

```text
PASS
FAIL
TIMEOUT
HARNESS_ERROR
INFRA_ERROR
VALIDATION_ERROR
PROTECTED_FILE_VIOLATION
INVALID_MUTATION
INVALID_SUPERVISED_EVAL
INVALID_ENVIRONMENT
```

Examples:

```text
mutation does not create the intended missing behavior
→ INVALID_MUTATION

A and B have different application dependency versions
→ INVALID_ENVIRONMENT

Docker/Compose fails before the agent starts
→ INFRA_ERROR

agent reaches the configured total agent-harness timeout
→ TIMEOUT

agent completes but required validation fails
→ FAIL
```

Do not interpret infrastructure failures as task-performance results.

---

# 28. Raw Evaluation Artifacts

**Raw evaluation artifacts are the original evidence emitted by the evaluation, not a human-written summary.**

Preserve them before interpreting or aggregating the results.

For each option, collect the raw:

```text
agent transcript
tool-call trace
agent stdout/stderr
harness stdout/stderr
model/provider usage response
timing data
final repository diff
changed-file list
test output
build output
validation output
container metadata
environment fingerprint
exit status
```

Do not collect secret values as raw evaluation artifacts.

For each mutation, preserve:

```text
mutation metadata
mutation seed
target region
ground-truth content or canonical ground-truth reference
mutated tree identifier
mutation patch
baseline validation output
mutation validation output
```

For each supervised eval, preserve separately from the agent-visible repository:

```text
supervised_eval_id
mistake_id
mistake commit
task
expected correction diff
baseline validation output
```

Do not add redundant hashes, path metadata, timestamps, or construction references to the persisted benchmark definition.

For documentation evaluations, preserve the original ground-truth document/region separately from the agent-visible repository.

Raw artifacts must be:

* collected from the actual execution
* stored without rewriting their contents
* associated with the exact trial and option
* timestamped where applicable
* sufficient to reconstruct what happened

Do not replace raw traces with only a summarized report.

Do not expose raw ground-truth artifacts to evaluated agents.

Do not expose Option A artifacts to Option B or vice versa.

---

# 29. Store Runtime Artifacts Outside the Repository

During evaluation, store runtime artifacts under a sibling directory one level above the project folder. Paths must remain relative to the project root:

```text
../tmp/harness-eval/<run-id>/
```

Use:

```text
../tmp/harness-eval/<run-id>/
├── manifest.json
├── environment.json
├── mutation/
│   ├── mutation.json
│   ├── ground-truth.patch
│   └── validation/
├── option-a/
│   ├── repo/
│   ├── container.json
│   ├── trace/
│   ├── stdout.log
│   ├── stderr.log
│   ├── diff.patch
│   ├── usage.json
│   └── validation/
└── option-b/
    ├── repo/
    ├── container.json
    ├── trace/
    ├── stdout.log
    ├── stderr.log
    ├── diff.patch
    ├── usage.json
    └── validation/
└── supervised-eval/
    └── <supervised-eval-id>/
        ├── task.md
        └── expected.diff
```

The `supervised-eval/` ground-truth directory is never mounted into either agent container.

Keep ground-truth artifacts separate from agent-visible repositories.

Keep Option A and Option B artifacts isolated from each other during execution.

---

# 30. Destroy the Evaluation Environment

After all raw artifacts have been collected:

1. stop the Option A Compose environment
2. stop the Option B Compose environment
3. remove both primary agent containers
4. remove writable volumes
5. remove temporary repositories
6. remove temporary runtime state

Do not leave agent state or evaluation state behind.

---

# 31. Persist Results

Only after the evaluation is complete and the containers are destroyed, copy the required artifacts into:

```text
.agent/harness_evals_results/<run-id>/
```

Do not expose this directory to the evaluated agents.

The persistent supervised-eval cache remains separate:

```text
.agent/supervised_evals/
```

When taking a future evaluation snapshot, always exclude both directories:

```text
.agent/harness_evals_results/
.agent/supervised_evals/
```

Historical results and supervised ground truth must never leak into future evaluated repositories.

---

# 32. Three Required Trials

Every A/B evaluation runs **at least 3 paired trials**: 3 initial screening trials, followed by additional batches of 3 confirmation trials when results are close or unstable.

Repeat the full evaluation process independently for each trial.

For reconstruction evals, each screening or confirmation trial gets:

```text
fresh mutation
fresh Docker Compose environment A
fresh Docker Compose environment B
fresh repository A
fresh repository B
fresh agent session A
fresh agent session B
```

Use a new mutation seed per reconstruction trial.

For supervised evals, each screening or confirmation trial gets the same cached benchmark inputs:

```text
same supervised_eval_id
same mistake commit
same task
same expected diff
```

but still gets:

```text
fresh Docker Compose environment A
fresh Docker Compose environment B
fresh repository A
fresh repository B
fresh agent session A
fresh agent session B
```

Do not reuse agent state, repositories, writable caches, or execution state between trials.

Within each reconstruction trial, the mutation is created exactly once and the resulting frozen mutation is shared by Option A and Option B. Do not independently generate mutations for A and B. Within each supervised trial, the cached historical benchmark is shared by A and B and no additional mutation is created.

The initial screening therefore has:

```text
Trial 1: Option A + Option B
Trial 2: Option A + Option B
Trial 3: Option A + Option B
```

When confirmation is required, append additional paired trials in batches of 3. Record every screening and confirmation trial independently, including all raw evaluation artifacts and metrics. Do not stop before the initial three trials are complete, and do not declare a winner from an unstable screening result.

---

# 33. Final Results

For every option report:

| Metric | Option A | Option B |
| --- | ---: | ---: |
| Task result | Score | Score |
| Tests | result | result |
| Agent time | duration | duration |
| Input tokens p/model used| count | count |
| Cached input tokens p/model used| count | count |
| Output tokens p/model used | count | count |
| Reasoning tokens p/model used | count | count |
| Model inference cost p/model used | amount | amount |
| Protected-file violations | count | count |
| Application dependency changes | count | count |
| Harness dependency changes | count | count |
| Exit status | code | code |

For supervised evals also report:

```text
supervised_eval_id
mistake_id
mistake_commit
supervised_diff_score
produced_diff_path
```

The expected diff remains available as hidden ground truth through the supervised-eval cache; its path is fixed by the cache layout and does not need to be duplicated in the result metadata.

For code review also report:

```text
defects_injected
true_positives
false_positives
false_negatives
precision
recall
f1_score
detection_score
validated_additional_findings
defect categories
```

For test-writing also report:

```text
all_tests_runnable = yes/no
all_tests_fail_for_expected_reason = yes/no
skeleton_valid = yes/no
semantic_test_quality = evaluator result
```

Also provide the locations of the raw artifacts.

Report every screening and confirmation trial separately. Identify the initial three screening trials and any additional confirmation batches.

Also provide an aggregate view across all screening and confirmation trials actually run where the metric supports aggregation, while preserving the initial screening results, each confirmation batch, individual paired results and raw evidence.

Do not collapse the experiment into a subjective ranking or winner.

---

# 34. Final Invariants

Before declaring an evaluation complete, verify all of these:

```text
[ ] main checkout was never used as an execution environment
[ ] main checkout was never modified for snapshot sanitization or benchmark construction
[ ] `.agent/persistant/knowledge/mistakes.jsonl` was omitted from agent-visible snapshots rather than deleted or modified in the main checkout
[ ] every copied file/dir was scanned for benchmark-leaking information before container creation
[ ] `.agent/session/` and `.agent/session-persistent-candidate/` contained only empty files/directories inside A and B containers
[ ] any agent-accessible path outside the project folder stayed within the project root's immediate parent and used exactly one ../ prefix where required
[ ] no agent-accessible path used an absolute external filesystem path or ../../ traversal
[ ] exactly two Docker Compose environments were created for the A/B trial
[ ] each option had its own primary agent container
[ ] each option had its own repository
[ ] both repositories started from the exact same frozen mutation before the intended A/B difference
[ ] the mutation was generated only once per trial
[ ] A and B had identical non-experimental environments before the intended difference
[ ] application dependencies were identical
[ ] any harness-level dependency differences were intentional and recorded
[ ] tests were identical before the intended eval mutation
[ ] validation commands were identical
[ ] task text was identical
[ ] ground truth was inaccessible to both agents
[ ] coding mutations removed application implementation logic and did not require test deletion
[ ] test-writing mutations deleted selected application implementation, tests and/or required skeleton so the agent reconstructed tests plus minimal stubs/skeleton only
[ ] test-writing agent did not implement actual application behavior
[ ] code-review mutations covered functional, security, spec-conformance and code-quality defect classes across the evaluation suite
[ ] code-review defects were not required to fail tests
[ ] no protected evaluation files were used to redefine correctness
[ ] agent time was measured separately from setup/validation
[ ] total agent-harness timeout was enforced
[ ] actual token usage was recorded
[ ] cost used the actual model and applicable current pricing
[ ] secret values were not extracted into evaluation metadata or artifacts
[ ] raw artifacts were preserved
[ ] supervised eval expected diffs remained hidden from evaluated agents
[ ] every cached supervised eval ran for all initial 3 A/B screening trials and any required confirmation trials
[ ] new mistakes were added to the supervised-eval cache only once and existing cached evals were reused
[ ] containers and temporary environments were destroyed
[ ] results were persisted under .agent/harness_evals_results/
[ ] supervised eval definitions were persisted under .agent/supervised_evals/
[ ] main checkout was not used as an evaluated execution environment
```

The purpose of this skill is to improve the harness itself systematically after very session:

```text
harness change
      ↓
isolated A/B experiment
      ↓
reconstruction evals: same frozen mutation + same task + same non-experimental environment
supervised evals: same cached historical mistake state + same task + same non-experimental environment
      ↓
independent agents
      ↓
correctness / supervised diff score + task time + token usage + cost
      ↓
raw evidence + final report + final harness improvement suggestions
      ↓
human makes the appropriate harness changes
```

Do not decide the outcome for the user, just give actioanable suggestions (maximum 3 suggestions).
