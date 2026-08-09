# [dev mode] [declarative] [self-contained] Build Front — *lunch-lang*: Polyglot Microservice build-time framework (distributed app compiler) for Modern Programming.


It is Ray for microservices. It is a distributed cloud programming framework that doesn’t require engineers to know distributed systems neither kubernetes in depth. More specifcally, it is a build-time decorator-based extension to multiple programming languages that makes building distributed systems easier (distributed app compiler). You use it’s CLI to build your code with lunch-lang decorators on top of any of the supported programming languages; the build generates final source code (code gen).

It also gives the user freedom to choose the tools used under the hood for the actual implementation (e.g., fastAPI/Uvicorn or gRPC for making HTTP server).
from your app code + lunch-lang decorators.

Helps with:
- aysnc programming
- concurrent programming
- distributed communication
- durable programs
- real-time media streaming
- communication with external services made my other comapny teams
- data-driven programming
- artifact usage
- artifact a/b testing
- parallel programming
- type-safe programming
- thread-safe programming
- observability with lineage
- documentation
- secrets
- functional testing
- profiling
- microservice configuration
- building the api gateway with authn/authz and DNS setup
- state re-initialization after change
- webserver build: write html builders, get webserver with DNS setup + SSL setup
- stable, tool-agnostic SDKs for supporting services (e.g., database, messae queue, background job)
- polyglot services (makes the necessary bindings and mappings) and frontends (compiles high performance non-javascript to webassembly that gets used within javascript)
making mcp servers
- task-level parallelism in DAGs
- lm-native programs

Use only the features you want (e.g., concurrent programming + testing + software infra).

Comes with IDE extension for highlighting, linting, etc

Examples of the things it will be able to do: 
- distributed communication: a build-time decorator that makes it clear that a list is a shared list between blocks (with updates being made to the list under the hood via durable (e.g. kafka) or non-durable (e.g., redis pubsub) streaming). Shared objects can be global or service specific.
- distributed communication: a build-time decorator that makes it clear that a function is an RPC function that can be used by another service (in sync or async manner);
- distributed communication: only allowed to call services explicitly defined as service dependencies.
- durable programs: you decorate side-effect awaits with @confirmation which will cause the @durable state (default is for all global objects to be durable) to be checkpointed to disk under the tag the that side-effect was indeed completed. If the program crashes after, it replays from the last checkpoint automcatically via k8s operator, and automcatically tries the next side-effect (via Idempotent execution). This also automatically makes the services that have the side effect rpc function to support idempotent execution and expose a helper query endpoint to see the status.
- real-time media streaming: make video & audio streaming declarative via decorator on top of video and audio objects.
- secrets: developers put any placeholder with @@ decorator. When building the code: will check if that team has acess indeeed to that service; (2) if has access: code will be changed to get the secret that was injected into the container at start-time, the secrets is gotten from the vault.
- data-driven programming: a build-time decorator that marks variables, functions/methods and/or models as learnable (with optional default value/distribution and optional text description in English of what you want to learn). These are learned via multiple analytical + numerical + AI methods working together to learn from your signal (loss function, natural language feedback and/or reward). Similar do DSPy, but works at the source code level and across distributed services/tasks, it is not a library, it is a systems optimizer working more like a compiler. Also get optimization time and cost estimates based on the size of your data/models, the complexity of your codebase and how many learnable object you have. You can also: (1) Impose symbolic regression or neural compression on learned functions to make tem simpler and interpretable; (2) learn control-theory controllers provided dynamicla system, observer and setpoint; (3) use supported ML frameworks & numerical solver libraries within your code; (4) impose equality/inequality constraints for the optimization; (5) define differential privacy requirements; (6) define sensitive variables.
- concurrent programming: a build-time decorator that make it clear that a block of code operates under the use of a lock + lock condition; 
- concurrent and easy async programming: build-time decorator tha tmake it clear that a function call is badly blocking other code and therefore should be submited as a Compute-bound routine or IO-bound routine to a few background threads (mixing IO-bound and Compute-bound task interleaving in the optimal way) that handles task execution via green threading (down the road, maybe this can even be detected automatically by static analysis + running the program in debug mode), like Goroutines (note: if using python =! 3.13 free-threaded then yuses background processes instead of threads because of the GIL); 
- parallel programming: a build-time decorator that makes it clear that a compute-heavy function should be compiled to exploit parallel hardware (CPU SIMD and/or GPU), of course it cant compile any function, has limitations on the libraries you use inside the function, like Numba).
- type-safe programming: type checking of programming languages with no enforced type checking (e.g., python). 
- type-safe programming: explicit http (gRPC or OpenAPI) interfaces for ready-made container services that generate client code.
- observability: build-time-decorator-based system-debugger. Unifies as a single visual story (that is simple at first view but can be expanded): profiling, events, traces, external communication, evals, binary problems, stack traces, logs (unstructured and structured), metrics, artifacts generated (artifacts use lineage-aware observability), state changes (the change and who caused the change) and user feedback. This enables a unified view of the entire system functioning, for root causa analysis and bootleneck identification. You can also set breakpoints on lines of code while you watch state, like in a traditional (local single-process) debugger. Note: observability decorators can be applied to data contracts as well, in case of ready-made supporting service (e.g., redis) that you dont want to analyze/make changes to its internals. Note 2: observability decorators can be part of a semantic observavility object (pipeline, AI agent, etc). E.g., a pipeline, where you define a pipeline of semantic steps and their associated code blocks; when observing the system, the developer will see the entire pipeline and the current sep being done. Exception telemtry is added by default, no need for developer instrumentation.
- testing: Codegen tests with local mocks by just decorating with examples (inuput-output examples for the code you are testing and mockers).
- testing: integration tests organized by different initial sample data
- testing: automatic load tests
- evals and offline experimentation on top of dev infra: try differents approches, services, methods, data, scale, etc. Evaluate them in parallel locally or in a remote cluster, and systematically choose the best combination. All in a git-native way (like dvc does).
- microservice configuration: define global and per-service configuration in beatifull ui fields, that can be changed any time and implemented on the microservices.
- building the api gateway with authn/authz + DNS setup: you setup the authorization rules (roles and accesse rights) automatically builds an API gateway for external clients to acess your externally facing functions. AI Gateway comes with built-in authentication/authorization server for authorization flows configured by the platform engineer (Private API Key, JWT, OAuth/OIDC, Public+Private Key).
- state re-initialization after change without re-deployment: set hooks for certain objects of fixed type. When this object changes, the state of the service is recomputed using the hook script which receives automatically the last object and the new object as arguments. E.g., when changing embedding model and the storage of embedding data needs to be re-done.
- thread-safe programming: compile-time thread safety checker (like in Rust) on top of existing languages without imposing new syntax on any language. This safety chcker will flag programs that cannot be proven to be thread-safe at compile time.
- artifact usage: refernce artifacts from the block Hub as if they were local files. Two types of ereferences: “to specific artifact” and “to best artifact”. If choose “to best artifact” the service will contiously ping the artifact hub for the best artifact or will subcribe for a new best artifact event.
- artifact versioning: annotate artifacts that will be generated dynamically in produciton, so that they are stored and versioned in the artifact registry with pointer to the exact commit and file/line where it gets generated.
- communication with external services made my other comapny teams: app developer dont want to know where the service is running jsut adress it my its domain name
stable, tool-agnostic SDKs for supporting services (e.g., database, messae queue, background job): you use stable SDKs that a build-time get transformed into client code for the specific tool configured (e.g., messae queue SDK becomes Kafka client code, because Kafka was configured as the Message Queue tool)
- making mcp servers: simply decorate a service block functions with @mcp
- for DAGs: task-level parallel programming: you might need to run 1000 containers at the sme time, each working on a piece of the data. Each container gets an instrisic ID like CUDA programming thread IDs. In this way you can easilt write only once, and during excution each container will ge its piece of data.
- lm-native: write tradtional functions but also now write llm functions which are functions with @@lm in the body. Which at build time is converted to the configured lm call using the fucntion signature and docdstring as prompts.

Ok, but what’s the first deliverable?

lunch-rpc: a subset of lunch-lang, focusing just on the distributed communication feature and just for the dev stage (excluding ci/cd and ops).
