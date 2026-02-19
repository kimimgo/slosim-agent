# 5. Discussion

## 5.1 Limitations

We identify eight limitations that bound the generalizability of our findings.

**Single GPU hardware.** All experiments ran on a single NVIDIA RTX 4090 (24 GB VRAM), which limits practical particle counts to approximately 500K. Industrial LNG tank assessments at full scale may require millions of particles for converged pressure predictions \cite{dominguez2022dualsphysics}. Multi-GPU or cluster-based execution --- supported by DualSPHysics but not yet integrated into SloshAgent --- would be required for production-scale deployments.

**Single LLM family.** All experiments used Qwen3 variants (32B and 8B) served via Ollama \cite{qwen2025, ollama2024}. We did not evaluate proprietary models (GPT-4, Claude) or other open-weight families (Llama, Mistral). While the zero-cost local deployment is a design feature, the generalizability of our accuracy and ablation findings to other model architectures remains untested. CFD-Copilot \cite{cfdcopilot2025} also uses Qwen3-32B for its general agent, suggesting this model class is competitive for simulation tasks, but direct cross-model comparison is needed.

**Solver-specific integration.** SloshAgent targets DualSPHysics v5.4 exclusively. The 14-tool interface assumes DualSPHysics-specific I/O patterns: XML case definition, GenCase preprocessing, binary particle output, and Run.csv telemetry. Other SPH solvers (SPHinXsys, GPUSPH) and mesh-based sloshing codes (OpenFOAM VOF) use different input formats and execution models. While we argue in Section 5.2 that the tool design patterns are transferable, this claim is not experimentally validated.

**DBC boundary limitation.** The oil sloshing failure in EXP-1 (zero detected impact peaks) is a documented limitation of the Dynamic Boundary Condition (DBC): its artificial viscosity over-damps low-energy viscous fluid impacts. English et al.\ \cite{english2021mdbc} demonstrated that mDBC (modified DBC) resolves this issue. The current tool set does not support mDBC configuration, making viscous-fluid sloshing an open limitation.

**Moderate overall accuracy.** The 47\% parameter accuracy in EXP-2 requires contextualization. This aggregate is driven down by L3 (paper-specific parameters absent from training data, 15\%) and L5 (edge cases requiring graceful refusal, 0\%). For realistic usage scenarios (L1 + L2), where users provide explicit parameters or domain-level descriptions, accuracy reaches 82\% (32B). Nevertheless, the overall figure indicates that SloshAgent is not yet a drop-in replacement for expert configuration at all complexity levels.

**8B ablation inversion.** The domain prompt ablation (EXP-4) only produces the expected hierarchy (FULL $>$ NO-RULES $>$ NO-DOMAIN $>$ GENERIC) on the 32B model. On 8B, the full 136-line prompt degrades tool-calling reliability (7/10 vs 10/10), causing the FULL condition to underperform NO-RULES. This indicates that domain prompt engineering requires a minimum model capacity threshold --- a practical constraint for resource-limited deployments.

**No user study.** All evaluations are automated or expert-judged. The claim that SloshAgent makes sloshing simulation accessible to non-specialists is qualitative. A formal user study with naval architects or mechanical engineers who lack CFD expertise would be required to substantiate this claim.

**Incomplete quantitative validation.** The Chen2018 parametric study comparison (EXP-3) demonstrates physically consistent trends but lacks point-by-point quantitative validation against published figures, which would require figure digitization. The SPHERIC validation (EXP-1) uses the community-standard peak-within-band metric but does not report time-series correlation, which is inappropriate for sparse impact signals (see Section 4.2) but is sometimes expected by reviewers outside the sloshing community.

## 5.2 From Sloshing to General Particle Simulation

While SloshAgent is intentionally sloshing-specific, three of its design patterns generalize to the broader Lagrangian particle simulation ecosystem:

**IsError pattern for silent failure detection.** Particle solvers share a class of failures that produce output without error messages: divergent energy, particle leakage through boundaries, insufficient boundary layer thickness. SloshAgent's `error_recovery` tool monitors Run.csv telemetry for anomalous kinetic energy growth and solver warnings. This pattern applies directly to DEM (Discrete Element Method) simulations, where particle overlap or explosive contact forces manifest as energy spikes, and to MPM (Material Point Method), where cell-crossing instabilities produce similar signatures.

**Asynchronous GPU job management.** The `job_manager` and `monitor` tools implement a polling-based architecture that decouples the LLM reasoning loop from the GPU solver execution. This pattern is necessary for any GPU-accelerated particle solver where simulation times range from minutes to hours --- far exceeding LLM API timeout windows. The same architecture would serve SPHinXsys (multi-physics SPH), LIGGGHTS (DEM for granular flow), and MPM solvers such as Taichi-MPM.

**Domain-encoded system prompts.** The EXP-4 ablation demonstrates that encoding domain physics (natural frequency formulas, fill-level sensitivity, parameter inference rules) into the system prompt yields +25 percentage points accuracy gain over a generic prompt (32B). This finding suggests that other particle simulation domains --- powder processing (DEM), geomechanics (MPM), astrophysical SPH --- would similarly benefit from domain-specific prompt engineering, provided the target LLM has sufficient capacity (32B+).

The broader implication is a "particle-solver + LLM agent" paradigm: any Lagrangian solver with a text-based input format, GPU execution, and domain-specific parameter constraints is a candidate for LLM agent automation using the patterns established here.

## 5.3 Industry Impact Projection

SloshAgent's practical value is best understood against the economics of current sloshing analysis workflows.

**Time reduction.** A single sloshing case requires expert-level XML authoring (1--4 hours depending on complexity), solver execution (minutes on GPU), and result interpretation. SloshAgent reduces the authoring phase to a single natural language prompt (seconds to minutes of LLM inference). For the 200+ case parametric studies required for LNG tank certification \cite{isope2012sloshing}, this compression from months of expert effort to days of automated execution represents a qualitative shift in workflow efficiency.

**Cost structure.** SloshAgent operates at zero marginal LLM cost: Qwen3 32B runs locally via Ollama on consumer-grade hardware. Combined with DualSPHysics (open-source, GPU-accelerated), the entire stack avoids both cloud API fees (\$0.01--0.06 per 1K tokens for proprietary models) and commercial solver licenses (\$10K+/year). By contrast, CFD consulting rates of \$80--120/hour and project costs of \$2K--20K per study \cite{cadcrowd2024} create significant barriers for iterative design exploration. BCG projects that AI-assisted R\&D can deliver 10--20\% time-to-market reduction and up to 20\% cost reduction \cite{bcg2025}.

**Reproducibility advantage.** All components of the SloshAgent stack are open: DualSPHysics (LGPL), Qwen3 (Apache 2.0), and SloshAgent itself. This contrasts with systems that depend on proprietary APIs (MetaOpenFOAM on GPT-4o \cite{metaopenfoam2024}, ChatCFD on DeepSeek \cite{chatcfd2026}), where reproducibility is contingent on API availability and pricing stability.
