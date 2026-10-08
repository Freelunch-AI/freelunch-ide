# Lunch Debug

System debugger that lets devs/engineers/agents visually or programatically:
- set distributed break points
- inspect "complete traces" (the causal history of external data exchanges, internal data esxcahnges, events, state changes and workloads triggered; profiled and wtith corresponding metrics and logs attached to them; with also sub traces of ai agents and DAGs embedded inside)
- travel to any point-in-time, replay and inspect traaffic or storage data of your system running on k8s (prod, staging or dev/debug environments)
- modify code or data mid-runs

If debugging prod or staging, a new nearly identical safe isolated environment is created for the debugger to work on.
