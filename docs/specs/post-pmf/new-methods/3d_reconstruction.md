# World-Model-Guided 3D Reconstruction from Long-Horizon Egocentric Video

## The Problem

Modern video world models are remarkably capable of predicting plausible future observations, but their temporal context remains fundamentally different from the context required for accurate 3D reconstruction. Video generation can often operate from only a small set of recent frames because its objective is local prediction: given what was just observed, generate what comes next. In contrast, reconstructing a complete 3D environment from egocentric video requires aggregating observations distributed across a much longer trajectory. The camera may only see the back of an object several minutes after first observing its front, revisit a room from a new viewpoint, or reveal previously occluded surfaces much later in the sequence. A reconstruction system therefore needs to preserve and combine scene-relevant information across potentially thousands of frames rather than relying on the short temporal windows commonly used for video generation.

This creates a fundamental separation between **video prediction** and **scene reconstruction**. Increasing the context length of a large video world model is an expensive and inefficient way to solve the reconstruction problem because most historical frames are redundant from the perspective of next-frame prediction. At the same time, purely geometric reconstruction methods are limited by the information contained in their observations and often struggle with incomplete visibility, ambiguous correspondences, occlusions, imperfect camera estimates, dynamic objects, and surfaces that are never directly observed. The central hypothesis of this work is that these two limitations can be addressed by combining the complementary strengths of the two systems: a specialized 3D reconstruction model should aggregate the long history, while a pretrained world model should provide the learned prior required to resolve ambiguity and refine the resulting scene.

## Core Idea

We propose a two-stage architecture in which **long-context geometric reconstruction and learned world knowledge are explicitly separated**.

Given a long egocentric video trajectory,

```math
V = \{I_1, I_2, \ldots, I_T\},
```

A specialized real-time, monocular-based 3D reconstruction system (e.g., LingBot-Map) incrementally updates a persistent representation, to produce an initial reconstruction

```math
S_0 = R(V).
```

The reconstruction may contain camera poses, depth, point clouds, Gaussian splats, meshes, object-level representations, semantic features, and uncertainty estimates. Importantly, this stage is allowed to use arbitrarily long temporal context. It is optimized primarily for geometric consistency rather than video generation.

The resulting reconstruction is then passed to a pretrained video world model, which serves as a **learned prior over physical s**. Rather than asking the world model to reconstruct the scene directly from thousands of frames, we condition it on the compact output of the reconstruction system and ask it to identify inconsistencies, infer missing structure, and refine the scene toward a representation that is more consistent with the distribution of real-world environments learned during large-scale video pretraining:

```math
S^* = W(S_0, V_{\mathrm{local}}),
```

where `W` is the world-model-based refinement operator and `V_local` can provide a small number of high-resolution observations when necessary.

The key design principle is therefore:

**Use specialized algorithms/models/geometry to make an intial estimate, and use the world model to refine the estimate using its prior spatial knowledge.**

## Why a World Model Should Help Reconstruction

A reconstruction model is fundamentally constrained by direct evidence. If an object is only partially observed, multiple 3D structures may explain the available images. If a surface is completely occluded, there may be no direct geometric evidence for its shape. If the camera trajectory is imperfectly estimated, correspondences across distant observations may also become ambiguous.

A large video world model has access to a different source of information: **statistical knowledge about how the physical world looks and behaves**. Through massive video pretraining, it can learn object structure, appearance continuity, common scene layouts, viewpoint transformations, object persistence, physical relationships, and correlations between visible and unobserved regions. The proposed system therefore treats the world model as a learned prior that resolves ambiguities left by geometric reconstruction rather than replacing the reconstruction process itself.

For example, the geometric model might determine with high confidence that several observations correspond to the same partially observed object, while leaving part of its geometry uncertain. The world model can use the visible evidence together with its learned prior to infer which completion is physically and semantically plausible. In this formulation, the world model is not hallucinating an arbitrary scene: it is refining a geometry-constrained hypothesis produced from real observations.

This distinction is crucial. The reconstruction model supplies the **evidence**, while the world model supplies the **prior**.

## Long-Horizon Memory Without Long World-Model Context

The central engineering advantage is that the proposed architecture does not require the video world model itself to attend to the entire trajectory.

A long egocentric sequence can instead be processed incrementally:

```math
M_{t+1} = \mathrm{Update}(M_t, I_t),
```

where `M_t` is a persistent 3D reconstruction memory. The reconstruction system absorbs observations from arbitrarily far back in time while maintaining a compact spatial representation. Once sufficient evidence has accumulated, the world model operates over this compact representation rather than over the complete pixel history.

This changes the scaling problem from **temporal attention over thousands of frames** to **reasoning over a persistent spatial state**. The world model can therefore retain its normal short-context video interface while gaining access to information that originated much earlier in the trajectory.

A further extension is to maintain a small episodic memory alongside the 3D state. The reconstruction system can retrieve the most informative historical frames associated with uncertain regions and present only those frames to the world model:

```text
Long egocentric video
        │
        ▼
Long-context 3D reconstruction
        │
        ├──────► persistent 3D scene
        │
        └──────► uncertainty / frame retrieval
                         │
                         ▼
                  World-model prior
                         │
                         ▼
                 refined 3D scene
```

This gives the world model access to both **global spatial context** and **targeted visual evidence** without requiring global pixel-level attention over the entire video.

## Iterative Refinement

The strongest version of the system should not necessarily be one-way. The initial reconstruction can be treated as a hypothesis that the world model critiques and refines, after which the reconstruction system can use the refined state to improve its own geometry.

```math
S_0 = R(V)
```

```math
S_{k+1} = W(S_k, V_{\mathrm{retrieved}})
```

followed by another reconstruction or optimization step:

```math
S_{k+1}' = R(V, S_{k+1}).
```

This creates a geometry-prior refinement loop. The reconstruction system enforces consistency with the observed trajectory, while the world model resolves ambiguities and produces a more plausible scene hypothesis. The two systems therefore act as complementary constraints rather than competing reconstruction methods.

## The Research Question

The central research question is:

> **Can the learned physical-world knowledge of a large video world model improve long-horizon 3D reconstruction beyond what can be achieved using geometric reconstruction alone?**

This can be decomposed into three subquestions.

First, can a specialized reconstruction system reliably compress arbitrarily long egocentric trajectories into a persistent spatial representation without retaining the entire pixel history?

Second, can a pretrained world model use that representation to resolve geometric ambiguities, complete partially observed structures, and correct inconsistencies that are difficult to solve from geometry alone?

Third, does the refined representation become a better state for downstream world modeling and robotics than either raw video or conventional reconstruction?

## Research Program

The first stage would establish a strong long-context reconstruction baseline. A modern video-to-3D system would process long trajectories and maintain persistent camera, depth, and scene representations. Experiments would explicitly stress long-range effects such as revisits, loop closure, large viewpoint changes, occlusion, and delayed observations of previously unseen surfaces.

The second stage would introduce the world-model refinement module. Rather than fine-tuning the entire world model immediately, several interfaces should be tested: conditioning on rendered views of the reconstruction, conditioning on 3D Gaussian or point representations, conditioning on depth and camera trajectories, and conditioning on learned reconstruction features. The objective is to determine which representation allows the pretrained world model to contribute the most useful prior while remaining grounded in the observed scene.

The third stage would introduce uncertainty-guided refinement. The reconstruction model would estimate where geometry is unreliable, and the world model would selectively operate on those regions rather than refining the entire scene. Historical frames associated with uncertain geometry could be retrieved and supplied as additional evidence. This turns the world model into a targeted **3D reconstruction reasoner** instead of an expensive global post-processor.

Finally, the system should be tested on dynamic egocentric environments. The representation can be extended from static 3D geometry to a 4D scene state containing persistent objects, poses, and motion. This would allow the same framework to distinguish between true geometric uncertainty and changes caused by moving objects or camera motion.

## Training Strategy

A particularly attractive property of this research direction is that the world-model refinement stage does not necessarily require large amounts of newly labeled 3D data. A world model pretrained on massive video corpora already contains the prior we want to exploit. The new training signal can instead focus on teaching the interface between reconstruction and generation.

Training examples can consist of real long trajectories with ground-truth or high-quality reconstructed scenes. Corrupted versions of those reconstructions can then be generated by introducing realistic failure modes: missing surfaces, incorrect depth, incomplete objects, misaligned fragments, noisy camera poses, inconsistent geometry, and partial observations. The world model is trained to transform these imperfect reconstructions into scene states that better explain the original observations.

Synthetic corruption is particularly useful because it allows control over the difficulty and type of ambiguity without requiring manually annotated reconstruction failures.

An even more powerful objective is to require the refined scene to support **novel-view prediction**. Rather than supervising only the final geometry, the system can render the refined reconstruction from previously unseen camera trajectories and compare those views against held-out observations. This directly tests whether the world-model prior recovered the true underlying scene rather than merely producing visually plausible geometry.

## Relation to Existing World-Model Research

Recent work has already established that persistent spatial memory can dramatically improve world models. *Video World Models with Long-term Spatial Memory* augments short-term frame history with geometry-grounded spatial and episodic memories, explicitly addressing the problem of forgetting previously observed environments. *Learning 3D Persistent Embodied World Models* similarly aggregates generated RGB-D observations into a persistent 3D map and uses that map to improve long-horizon simulation. More recent approaches such as PERSIST make the 3D scene itself an explicit persistent latent state of the generative world model.

There is also increasingly direct work connecting world models and reconstruction. WorldStereo introduces geometric memories inside a video generation system and uses them to obtain more geometrically consistent multi-view generation and 3D reconstruction. WorldMirror explores a unified reconstruction model that can incorporate geometric priors such as camera poses and depth to resolve ambiguities in 3D prediction. FR3D goes further toward persistent future 3D representations and uses foundation-model knowledge to improve zero-shot geometric prediction.

However, these directions largely frame the problem as **building 3D memory into the world model, jointly generating geometry and video, or using geometric information to control generation**. The proposed research asks a different question: **can an already powerful generative world model be used as a learned prior to refine the output of an independent long-context 3D reconstruction system?** This separation makes it possible to exploit mature reconstruction algorithms for arbitrarily long observations while using the world model specifically where learned world knowledge provides value.

This distinction is also important for systems such as Cosmos. NVIDIA's current Cosmos family supports substantially longer video contexts than earlier generations, including hundreds of frames depending on resolution, but its primary abstraction remains video generation rather than persistent reconstruction from an arbitrarily long observational trajectory. The proposed approach does not attempt to turn Cosmos into a giant-context structure-from-motion system. Instead, it uses a dedicated reconstruction module to solve the long-memory problem and exposes the resulting compact spatial state to the world model.

## Evaluation

The project should be evaluated on three levels.

**Geometric accuracy:** depth, camera trajectory, point-cloud alignment, surface reconstruction, Gaussian or mesh quality, and consistency across revisited areas.

**World consistency:** whether the reconstruction remains consistent when viewed from novel viewpoints, whether objects retain their identity and structure across long trajectories, and whether previously occluded or weakly observed regions are reconstructed correctly.

**Downstream physical reasoning:** whether the refined scene produces better future-view prediction, action-conditioned simulation, and robot planning than the initial reconstruction.

The most important ablation is straightforward:

```text
Long video
   ↓
3D reconstruction
```

versus

```text
Long video
   ↓
3D reconstruction
   ↓
World-model refinement
```

Additional ablations can remove the long-term memory, the world-model prior, uncertainty-guided retrieval, or iterative refinement. This will establish whether improvements actually come from the learned world prior rather than simply from additional computation.

## Why This Matters

A successful system would establish a new division of labor between perception and generative world modeling. Instead of forcing a single video model to simultaneously remember arbitrarily long histories, recover geometry, infer occluded structure, and predict future observations, each component would operate in the regime where it is strongest.

The 3D reconstruction model becomes the **persistent memory mechanism**, capable of integrating observations over minutes or hours of egocentric video. The world model becomes the **learned physical prior**, capable of interpreting ambiguous or incomplete geometry using knowledge acquired from massive video pretraining. Together they produce a persistent world representation that is both grounded in real observations and informed by learned world knowledge.

The ultimate goal is not merely better 3D reconstruction. It is to establish **3D scene state as the bridge between long-horizon visual memory and generative world modeling**:

```text
                 Long Egocentric Video
                          │
                          ▼
              Long-Context 3D Reconstruction
                          │
                          ▼
                Persistent 3D / 4D State
                          │
                  ┌───────┴───────┐
                  ▼               ▼
            World-Model       Memory / Retrieval
                Prior               │
                  └───────┬─────────┘
                          ▼
                 Refined World State
                          │
                ┌─────────┴─────────┐
                ▼                   ▼
        Future Simulation       Robot Planning
```

If successful, this would provide a practical path from massive egocentric video to persistent, physically meaningful world representations without requiring the world model itself to process the entire history at every prediction step.
