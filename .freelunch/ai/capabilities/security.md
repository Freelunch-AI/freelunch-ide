# Security Review

Use this capability when a change touches trust boundaries, authentication, authorization, secrets, user-controlled input, network access, serialization, file paths, dependencies, CI/CD, or infrastructure.

1. Identify assets, actors, entry points, privilege boundaries, and attacker-controlled data in the changed surface.
2. Trace input through validation, authorization, storage, execution, and output.
3. Check for injection, cross-site scripting, request forgery, path traversal, unsafe deserialization, secret exposure, insecure randomness, authorization bypass, IDOR, and unsafe defaults where applicable.
4. Review dependency and build changes for provenance, pinning, install scripts, and least privilege.
5. Check logs, errors, telemetry, and progress records for leaked credentials or sensitive data.
6. Prefer a reproducible exploit path or failing test over a speculative warning.
7. Separate confirmed vulnerabilities from hardening suggestions and false-positive-prone patterns.

This review is assistive. It does not replace human review, SAST, dependency scanning, secret scanning, DAST, or penetration testing required by CI/CD or release policy.
