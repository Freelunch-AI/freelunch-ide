# [Outputs Artifacts] [dev mode] [imperative] [self-contained]  [mostly code] Data & AI Artifacts Front — lunch-artifact

ide plugin, notebook-based, for experimenting, building, tracking, publishing to Artifact Registry and promoting artifacts to PROD. Artifacts: datasets, data transform pipelines (can be implemented as ETL, ELT or Streaming), prompts, models, inference req/res pipelines, agent definition packs. Contains these sub-GUIs: dev-mode blocks view, data storage discovery (viz & query), artifact registry, visual data transform pipeline (ETL/ELT/Streaming) creation, notebooks with code vizualization, dataset visualization & annotation, no-code automl and playground. No code GUI-supported AutoML features: dataset vibe-creation/improvement (based on dataset natural language description/feedback or envisioned finetuned model), dataset vibe-eda (with genrative UI), model vibe-training/tuning, dataset vibe-evaluation and model vibe-evaluation. Also, visualization of your code as computational graph with tensor visualizations. Also model expert GUI: model conversion, profiling, interpretability, comparison, compression, compilation, scaling law estimation.

Note: models can only be used by developers once they are evaluated and vetted by an Data/AI Scientist.

Notebook run on an optionally GPU-powered remote dev environment with access to a ready-to-use ephemeral worklads cluster (SkyPilot + KubeRay/slurm/NVIDIA NeMo/hf accelerate/jax/dask/spark/deepspeed/kubeflow trainer/kubetorch) running on the data/ai experimentation k8s+kueue cluster. In which cloud the ephemeral workloads cluster runs on? That will depend on the configuration set (e.g., TPU → GCP; High-end GPU → Nebius; Simple GPU → AWS; CPU → Azure) and on cost analysis.

Note: this front works on top of a separate artifacts-specific repo folder (but user doesnt interact directly with the repo) use git-based experiment tracking where an experiment is defined as a commit with: command to run to generate artifacts, artifacts generated schema and locations, command to run to evaluate each artifact. Artifacts per se are not tracked by git, they have pointer file (that points a file that has the verison tree tracked with cloud pointers to each artifact version in the tree) that is tracked by git and lakefs is used to do git-like versioning on the cloud

Artifacts can be stored with the following info:

- reproducibility: git-commit experiement that produces and evaluates it. Its a git commit with: conatiner image that things run inside, artifact build script, artifact evaluation script
- ownership: who built that artifact
- versioning: the version of the artifact
- evaluation: the eval metrics of the artifact
- lineage: dependency graph showing artifacts and libraries used to create it
- status: in production, lifeboat or dev
- size
- last modified
- evalauted or not
- vetted for prod or not
- precondition data tests: required condition of the data scehama and data of the data source which enables the transformation to be applied
- preconditin data transformation: features that need to be avaialble in the describe data source, also if these fatures need to e contoiously built via batch and/or streaming. Features are described via engine-agnostic (using Ibis) data transformations with support for multiple data sources/targets, transformation engines & transformations patterns (batch, streaming or both)

