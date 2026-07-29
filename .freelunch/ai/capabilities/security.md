# Security Review

Use this capability when a change touches authentication, authorization, tenants, secrets, user-controlled input, network or file access, serialization, dependencies, LLM tools, CI/CD, infrastructure, deployment, or another trust boundary. It combines Agency AppSec and AI-code auditing with the layered review model from Anthropic security guidance.

## Three Review Layers

1. **Pattern pass:** inspect changed lines for known-dangerous APIs and obvious local mistakes such as hardcoded secrets, unsafe deserialization, raw command construction, unsafe HTML, permissive policies, or path use.
2. **Diff pass:** reason about the complete change, intended behavior, privilege changes, and security regression tests.
3. **Context pass:** trace attacker-controlled data and authority across related files to find authorization bypass, IDOR, SSRF, prompt injection, secret exposure, and cross-file validation gaps.

These are review procedures, not installed hooks. Do not add a plugin, start a model call, send code to another endpoint, or change Git behavior.

## Threat Model

For a material boundary, identify:

- assets and sensitivity;
- actors and credentials;
- entry points and attacker-controlled values;
- components, data flows, and trust boundaries;
- authentication and authorization decisions;
- persistence, logging, external calls, and privileged actions;
- spoofing, tampering, repudiation, disclosure, denial-of-service, and privilege-escalation paths that actually apply.

Convert material threats into testable requirements. Prefer a small data-flow or attack-path description over a generic OWASP checklist.

## Review Checklist

### Identity And Access

- Authenticate at the intended boundary and authorize every protected object/action server-side.
- Check ownership and tenant scope, not only role presence.
- Do not trust client-editable metadata, headers, form fields, or model output for privilege decisions.
- Use least-privilege service identities and avoid privilege propagation across background jobs or tools.

### Input, Execution, And Output

- Validate type, shape, length, range, encoding, and allowed targets at trust boundaries.
- Keep data separate from SQL, shell, templates, paths, URLs, and model instructions.
- Use parameterized APIs, allowlists, canonical paths, output encoding, and safe parsers.
- Bound pagination, uploads, decompression, recursion, regex work, retries, and concurrency.
- Avoid exposing stack traces, internal identifiers, sensitive configuration, or cross-tenant data.

### Secrets And Cryptography

- Never print or repeat a raw secret in findings. Report type, redacted location, and exposure path.
- Treat a committed or client-reachable secret as compromised; removal must be paired with provider-side rotation and history/impact assessment.
- Distinguish intentionally public identifiers from credentials to avoid false positives.
- Use maintained cryptographic libraries and correct nonce, randomness, comparison, storage, and key-lifecycle behavior. Never design a custom primitive.

### Dependencies And Delivery

- Verify canonical source, immutable version or digest where appropriate, license, transitive impact, install scripts, and maintainer posture.
- Review workflow permissions, untrusted PR execution, artifact provenance, secret availability, cache poisoning, and deploy credentials.
- Require explicit rollback and health evidence for risky delivery changes.
- Treat scanners as evidence sources, not proof of security; validate exploitable findings and record scanner gaps.

### AI-Assisted And LLM Code

- Keep untrusted content in a data/user role rather than concatenating it into system or developer instructions.
- Treat tool-enabled model calls as privileged execution: validate arguments, scope tools, enforce authorization outside the model, and require confirmation for consequential actions.
- Check client bundles and public environment prefixes for secret leakage.
- Review generated database policies for missing enforcement, blanket-allow predicates, and client-controlled role checks.
- Mark prompt-injection and taint conclusions with honest confidence; demonstrate the source-to-sink path.

## Finding Standard

For each confirmed issue provide:

- source and sink or violated boundary;
- concrete exploit or failure path;
- affected asset and impact;
- severity and confidence;
- smallest secure fix in the repository's language and framework;
- regression test or rescan needed to prove closure;
- rotation, migration, or incident step when exposure already occurred.

Separate confirmed vulnerabilities from hardening suggestions. Do not claim compliance, completeness, or a security percentage.

This review supplements, but does not replace, human review, secret scanning, SAST, dependency scanning, image scanning, DAST, or penetration testing required by project policy.
