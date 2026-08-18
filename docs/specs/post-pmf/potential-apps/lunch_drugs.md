# Lunch Drugs | Drug Discovery App: digital twin of human body for drug discovery at scale

Disclaimer: brain is excluded from the simulation

## Overview

Digital Twin of the human body that can be used for:

- testing new drugs
- finding optimal drug dosage and diet.
- understading the effects of a drug, infection, desease, lesion.

It’s a foundation model with human body-centric priors based on scientific models, which is trained on multimodal biological data.

It will be a SaaS platform for pharma R&D teams to simulate drug ADMET (absorption, distribution, metabolism, excretion, toxicity) across virtual populations.

## **Major challenges**

- **Data:** Getting enough and timely biodata (genomic, proteomic, metabolomic, imaging, sensor, clinical, wearable, etc)  to make it possible to learn latent variables and parameters specific to that person
- **Evals:** Digital Twin Evaluation: controlled clinical studies comparing predicted vs. actual pharmacokinetics/dynamics.
- **Modelling**
    - Finding the rights abstraction/solution levels to ensure proper modelling/solving and the simulation being computationally feasble. Some parts/moments need more fine-grained modelling/solving than others.
    - Doing ML on top of prior biological knowledge of the human body and governing biophysics.
    - Mapping fisiological behaviours to biological markers to be able to map simulation observations to human behavioural change
    - Modular modelling/learning of body parts (lungs, heart, vascular system, etc) separately and then a final unified learning phase
    - Multi-scale modelling in parallel (from causal PBK models with learnable variables, to models with ML models al the way to a foundation model modelling the entire subsystem) and informing each other for convergence in distribution.
    - Leveraging multimodal data.
- **Computing:** High Performance Computing for cutting down learning/inference compute costs and learning/simulation time

## Major fields and methods

- Foundation Models/LLMs
- Scientific ML / Physics-informed ML
- Causal Inference
- Computational Chemistry (low-level)
- Computational Biology (mid-level)
- Medicine (high-level)

## **MVP**

- **Start with one organ and one use-case** (e.g., e.g., liver digital twin for personalized oncology drug dosage).
- **Leverage open biomedical datasets** (e.g., UK Biobank, NIH OpenData, Physiome Project).
- **Build multi-scale models** — couple mechanistic PBPK models with lower level ML-learned layers.
- **Partner with clinics/pharma** for training and validation data, and co-development.
- **Build as a non-critical helper tool** to find good candidates for real world clinical drug trials on animals.
