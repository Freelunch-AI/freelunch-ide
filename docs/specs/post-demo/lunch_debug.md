# Lunch Debug

System debugger that lets devs/engineers/agents visually or programatically:
- reply production traffic wit only a subset of necessary production workloads deployed in the isolated system debugging environment
- set distributed break points
- inspect "complete traces" (the causal history of external data exchanges, internal data esxcahnges, events, state changes and workloads triggered; profiled and wtith corresponding metrics and logs attached to them; with also sub traces of ai agents and DAGs embedded inside)
- travel to any point-in-time, replay and inspect traffic or storage data of any part of your system
- modify code or data mid-runs

Notes:
- If debugging prod or staging, a new nearly identical safe isolated environment is created for the debugger to work on. Data-heavy services are replicated with point-in-time synced data. However the default scale is very small, the goal is to inspect individual bugs, not to replicate scaling issues.
