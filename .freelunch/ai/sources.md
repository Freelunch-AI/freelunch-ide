# Source Integration Policy

`ai-stack.lock.json` is the machine-readable source of truth for every external methodology reviewed for this bundle. It pins an immutable revision, license, selected files, adopted behavior, and intentionally excluded behavior.

The source projects are design inputs. They are not installed dependencies and their agents are not additional public skills. FreeLunch keeps five public entry points and loads the following original adaptations only when the changed surface warrants them.

## Coverage

| Source | What FreeLunch retains | How it is used |
| --- | --- | --- |
| [BMAD Method](https://github.com/bmad-code-org/BMAD-METHOD) | Stakes-scaled requirements, explicit assumptions, architecture invariants, ordered implementation, definition-of-done evidence, independent review | Core procedure across all five workflows |
| [Agency Agents](https://github.com/msitarzewski/agency-agents) | Architecture, minimal-change, review, DevOps, SRE, technical-writing, AppSec, AI-code security, test-automation, and evidence perspectives | Routed internally by `.freelunch/ai/capabilities/specialists.md` |
| [GBrain](https://github.com/garrytan/gbrain) | Thin public router, progressive context loading, durable-decision lookup, source precedence, conflict/gap reporting, friction capture, test-before-bulk | Repository-native context protocol; optional read-only lookup when GBrain is already available |
| [Graphify](https://github.com/Graphify-Labs/graphify) | Graph-first navigation for broad architecture and dependency questions, source-backed edges, confidence labels, query/path/explain selection | Optional adapter when an existing graph and tool are already present; ordinary repository search remains the fallback |
| [RTK](https://github.com/rtk-ai/rtk) | Selective compression of noisy command output, honest savings language, and raw-output fallback | Optional adapter when `rtk` is already installed; never a required command prefix |
| [Anthropic Claude Code plugins](https://github.com/anthropics/claude-code/tree/main/plugins) | High-signal review, second-pass validation, changed-code simplification, comment/test/error/type lenses, layered security review, concise explanation, and progressive skill design | Original FreeLunch procedures only; upstream text is reference-only because redistribution is not licensed |

## Named Inputs

The names discussed in issue planning map into the internal layer as follows:

| Named input | FreeLunch location |
| --- | --- |
| `agency-agents` | `capabilities/specialists.md` |
| `gbrain` | `capabilities/context.md`, plus progress and source-precedence rules in `workflow.md` |
| `graphify` | `capabilities/code-intelligence.md` |
| `rtk` | `capabilities/command-output.md` |
| `code-review` and PR review agents | `capabilities/review.md` |
| `code-simplifier` | the changed-code simplification section in `capabilities/implementation.md` and the simplification lens in `capabilities/review.md` |
| `security-guidance` | the three review layers and detailed checks in `capabilities/security.md` |
| `explanatory-output-style` | `capabilities/explanation.md`, activated only when explanation adds value |
| `skill-creator` | `capabilities/skill-authoring.md` |

These paths are deliberately outside every client's discovery directory. A user sees five choices; the selected workflow sees the deeper procedures it needs.

## Loading Rules

1. Start with the invoked public skill and `.freelunch/ai/workflow.md`.
2. Load only the capability files named by that skill or triggered by the changed surface.
3. Apply specialist lenses inside the current workflow. Do not advertise or invoke them as extra commands.
4. Treat optional tools as accelerators, never prerequisites. Do not install, configure, or update them as a side effect of a workflow.
5. Prefer repository evidence over an upstream example when the two differ. Upstream stack-specific commands and numeric targets are not FreeLunch requirements.
6. Follow the current repository and client safety rules even when an upstream source suggests autonomous mutation, posting, merging, or auto-fixing.

## Deliberate Exclusions

- No source repository is vendored wholesale. Large upstream prompt packages would make context routing worse and would drift from their owners.
- No Python, JavaScript, Go, Rust, shell, generator, transpiler, or hook is added to run these Markdown workflows.
- No GBrain, Graphify, or RTK installation or account is required.
- No Anthropic plugin content is copied. The pinned files are traceable design references only.
- No reviewer must invent a quota of findings, fail a first iteration by default, or produce a score without evidence.
- No compressed or generated representation is accepted as the sole evidence for a finding.

## Updating A Source

Treat a revision change as a dependency update:

1. inspect the upstream diff between the old and proposed revisions;
2. verify the new revision and license from the canonical repository;
3. review every selected file that changed;
4. update the relevant FreeLunch capability only for behavior that still fits this repository;
5. update `ai-stack.lock.json` and explain adopted and rejected changes in the pull request;
6. verify all five public workflows still resolve and no new public command appeared.
