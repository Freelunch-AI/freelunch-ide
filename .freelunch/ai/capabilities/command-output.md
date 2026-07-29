# Command Output Discipline

Use this capability when shell output may be large. It adapts RTK's output-filtering approach without making RTK a dependency or hiding evidence.

## Optional RTK Use

1. Check whether `rtk` is already available before considering it.
2. Do not install RTK, initialize client hooks, rewrite global configuration, or require every command to use it.
3. Use it selectively for supported, noisy commands such as test, lint, build, status, log, or routine infrastructure output.
4. Prefer structured command output such as `--json` plus a real parser when the underlying tool provides it.
5. Record the command actually executed, including an `rtk` prefix when used.

RTK's savings measure reduced command-output bytes. Do not present that percentage as total token, latency, or cost savings.

## Raw Evidence Fallback

Use raw output instead of a filter when:

- reviewing the exact patch or surrounding diff context;
- diagnosing a failure whose details were omitted;
- verifying warnings, skipped checks, counts, exit codes, or ordering;
- inspecting security scan evidence, stack traces, generated files, or binary/tool identity;
- the command is unsupported or filtered output looks incomplete;
- a result will be cited as a review finding or release gate.

When compressed output reports a failure, rerun the smallest failing target raw. When it reports success, still verify the process exit status and any required counts or skipped checks.

## Output Handling

- Start with the narrowest command that can prove the behavior.
- Preserve error context needed to reproduce and debug the failure.
- Summarize repetitive success output, but never suppress a distinct failure.
- Keep secrets and tokens out of commands, logs, progress records, and reports.
- Label a check `passed`, `failed`, `skipped`, or `unavailable`; never collapse the last three into success.
