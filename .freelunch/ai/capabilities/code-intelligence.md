# Codebase Intelligence

Use this capability for cross-module architecture, dependency paths, data flow, impact analysis, or a review whose scope is too broad for one-file inspection. It adapts Graphify's graph navigation while preserving a source-first fallback.

## Select The Narrowest Reliable Method

1. Use direct file reads and repository search for a known symbol, file, or small change.
2. Use language-aware references, call hierarchy, or existing repository tooling when available.
3. Use Graphify only when `graphify-out/graph.json` already exists and a working `graphify` command or configured integration is already available.
4. Never install Graphify, build a graph, configure an API key, enable a hook, or update generated graph files as a side effect of a FreeLunch workflow.
5. If the graph is absent, stale, unhealthy, or insufficient, fall back to source search and say which limitation applied.

## Graph Operations

Choose the operation from the question:

| Question | Operation |
| --- | --- |
| Broad neighbors or "what connects to X?" | `graphify query` with breadth-first traversal |
| Specific dependency or route from A to B | `graphify path` or depth-first query |
| Focused concept and its immediate relationships | `graphify explain` |
| Broad architecture only after scoped operations are insufficient | `GRAPH_REPORT.md` or the generated wiki |

Use terms found in the graph's vocabulary. If no graph label plausibly matches the question, stop using the graph rather than manufacturing synonyms or edges.

## Evidence Rules

Graph relationships have different evidentiary weight:

- `EXTRACTED`: deterministically present in source. Confirm freshness and cite the source location.
- `INFERRED`: supported by contextual evidence but not a direct structural fact. Validate it in source before using it for an implementation or review conclusion.
- `AMBIGUOUS`: a lead for investigation, never a finding or architectural fact by itself.

Never invent a relationship. Cite source files and locations from the graph, and distinguish what the graph reported from what source inspection confirmed. Generated graph output is navigation help, not the repository's source of truth.

## Impact Analysis

For a proposed or actual change, trace:

1. entry points and callers;
2. state and data ownership;
3. public interfaces and consumers;
4. persistence, network, process, and privilege boundaries;
5. error propagation and fallbacks;
6. tests that exercise the path;
7. deployment, observability, and documentation consumers when applicable.

Report the confirmed path, uncertain edges, affected modules, and the smallest verification set. Do not expand the implementation solely because the graph surfaced unrelated neighbors.
