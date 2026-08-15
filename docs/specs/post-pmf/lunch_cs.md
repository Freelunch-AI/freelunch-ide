# Lunch CS | Great way to learn Computer Science  & Engineering | Computer Emulation Visual Playground Platform: understand low-level Modern Cloud Computing visually

Interactive & Visual Playground running in the browser, where users can understand how Computers, Clusters & Compilers work by visualizing low-level details involving hardware (compute & communication), OS, drivers, protocols and compilers.

Things that are not emulated:

- Digital electronics (e.g., digital circuits of CPU cores & Memory are not emulated, but the behaviour is emulated using their abstractions) never emulated
- Power electronics (e.g, battery) never emulated
- Process-level details of kernel-space code (OS/driver/EBPF code) not emulated by default

Basically super **lighweight VMs (small memory/cores) & VNs running in the browser**, but not only can they be used, they are also fully observable and editable.

Comes with a “Modern Low-level Computing” **Video Course** that teaches low-level computing (from Processor/Memory/Storage/OS/Networking fundamentals, to Compilers/Interpreters, and Distributed Systems) using the Playground. The course is a game where you are guided to build mini implementations of popular open source tools following 2 steps: (1) naive implementation; (2) optimized implementation.

**Similar Projects:** **

- *OpenROAD* (playground, low-level, visual, very detailed (RTL) emulation)
- *gem5* (playground, low-level, cloud machines, not visual, detailed (Microarchitecture) emulation)
- *QEMU* (not playground, low-level, cloud machines, not visual, simplified emulation)
- *GIRUS* (playground, high-level, cloud clusters, visual)
- *Tinkercad* (playground, high-level, arduinos, visual)
- *Computer Science Crash Course* (youtube cs intro course, high-level, visual)

**Playground Features:**

- source code/libraries-executable/shared objects visual mapping to understand the compilation & linking process
- debug/replay/record in slow speed
- real-world observability logs and traces on a sidebar
- QEMU-backed fast-forwarding
- emulate components in isolation, avoiding having to emulate entire computer all the time
- disable some components to make visualization cleaner (e.g., disable drivers source code line-by-line execution visualization, disable k3s visualization)
- expand components to see their internal components/workings & minimize components to see just interface (e.g., Virtual System RAM can be expanded to MMU, Memory Controller, Physical DRAM, Peripheral Registers)
- IDE p/node where user can use terminal, edit files, install and run low-memory programs (cpu-only, gpu and tpu programs) with debuggers
- switch components of the same type, using a component library (e.g., switch OS, switch CPU, switch GPU, switch storage, swith ethernet, switch intra-node GPU interconnect, switch inter-node GPU interconnect, etc)
- profile and get calibrated real world estimates of the playground profiling
- stop and inspect values in storage/memory/cache/registers mid-execution
- chaos engineering: overheat a piece of hardware, kill a piece of hardware or even kill an entire node

Also comes with an **AI Assistant** that can answer questions based on text/highlight/image capture/video capture.

Example learning scenario:
