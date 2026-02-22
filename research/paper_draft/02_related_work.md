# 2. Related Work

## 2.1 Sloshing Simulation: State of Practice

Sloshing analysis has historically relied on three complementary approaches: physical model testing, mesh-based CFD, and particle-based methods.

**Physical model testing** remains the gold standard for LNG tank certification. GTT (Gaztransport & Technigaz) and Seoul National University operate the world's largest sloshing test facilities, with SNU's database comprising over 20,000 hours of scale-model experiments and 540 TB of raw data across multiple tank geometries, fill levels, and excitation conditions \cite{kim2019snu}. The ISOPE comparative study \cite{isope2012} coordinated nine institutions to benchmark sloshing predictions against a Mark III tank geometry (946$\times$118$\times$670 mm), revealing that 9 of 16 CFD pressure peak predictions fell outside the $\pm 2\sigma$ experimental envelope. While authoritative, physical testing costs millions of dollars per campaign and requires months per evaluation cycle.

**Mesh-based CFD** is the dominant computational approach. Commercial solvers (ANSYS Fluent, STAR-CCM+) and the open-source OpenFOAM (interDyMFoam with VOF) are widely used for sloshing analysis \cite{chen2018}. However, mesh-based methods face fundamental limitations for violent sloshing: free-surface fragmentation causes mesh distortion, wave breaking requires interface reconstruction, and splashing demands adaptive mesh refinement --- all of which increase computational cost and reduce robustness for the most physically interesting (and dangerous) sloshing regimes.

**Smoothed Particle Hydrodynamics (SPH)** addresses these limitations through a meshless Lagrangian formulation. DualSPHysics is the leading open-source GPU-accelerated SPH solver, with 687 GitHub stars, 52,000+ downloads, and 746 citations \cite{dominguez2022}. It achieves 40--100$\times$ speedup over CPU implementations through CUDA parallelism \cite{crespo2011}. Commercial alternatives include Altair nanoFluidX (\$10K+/year) and AVL PreonLab. English et al. \cite{english2021} validated DualSPHysics with the modified Dynamic Boundary Condition (mDBC) against SPHERIC Test 10 sloshing benchmarks, demonstrating quantitative agreement with experimental pressure measurements. Liu et al. \cite{liu2024} conducted a comprehensive parametric study (3 fill levels $\times$ 5 frequencies $\times$ 3 amplitudes) for large-tank pitch sloshing.

**DesignSPHysics** \cite{designsphysics2023} is the only graphical frontend for DualSPHysics, built as a FreeCAD plugin. However, it remains self-described as "early Beta stage, not meant to be used in a stable environment" (v0.7.0). Critically, it does not support mDBC (Issue \#171), the boundary treatment required for accurate sloshing pressure prediction, and suffers from FreeCAD version compatibility issues. Its 135 GitHub stars (versus DualSPHysics's 687) reflect limited adoption.

A distinct line of work uses neural networks for sloshing load **prediction** (not simulation). Kim et al. \cite{kim2019snu} trained an ANN on SNU's experimental database to predict sloshing pressures without running new simulations. While valuable for rapid screening, this approach cannot generate flow fields, validate new geometries, or replace physics-based simulation for certification.

## 2.2 LLM Agents for CFD and Simulation

The rapid development of LLM agents for scientific simulation has produced over ten systems targeting computational fluid dynamics and related domains. Table~\ref{tab:comparison} provides a systematic comparison verified from full-text review of each publication.

**Table 1.** Comparison of LLM-based simulation agents. All entries verified from primary publications.

| System | Year | Solver | Architecture | LLM | Benchmark | Success Metric | Cost |
|--------|------|--------|-------------|-----|-----------|---------------|------|
| MetaOpenFOAM 1.0/2.0 | 2024--25 | OpenFOAM | 4-agent MetaGPT | GPT-4o ($T{=}0.01$) | 8$\to$13 cases | 85--86.9\% Pass@1 | \$0.15--0.22/case |
| OpenFOAMGPT 2.0 | 2025 | OpenFOAM v2406 | 4-agent Prompt Pool | Claude-3.7-Sonnet ($T{=}0$) | 455 cases | 100\% reproducibility | Cloud |
| Foam-Agent 2.0 | 2025 | OpenFOAM | 6-agent + 11 MCP | Claude 3.5 Sonnet | 110 tasks | 88.2\% exec success | ~334K tok/case |
| ChatCFD | 2026 | OpenFOAM | 4-stage + structured KB | DeepSeek-R1+V3 | 315 cases | 82.1\% exec / 68.12\% phys | \$0.208/case |
| PhyNiKCE | 2026 | OpenFOAM | Neurosymbolic (CSP) | Neural+symbolic | 13 configs, 340 runs | 96\% over ChatCFD | N/A |
| CFD-Copilot | 2025 | OpenFOAM v2406 | MetaGPT + 100+ MCP | Qwen3-8B LoRA + 32B | NACA 0012 | $U$ 96.4\%, $p$ 93.2\% | Local |
| AutoCFD | 2025 | OpenFOAM | Fine-tune + multi-agent | Qwen2.5-7B (28.7K pairs) | 21 cases | 88.7\% accuracy | \$0.020/case |
| FoamGPT | 2025 | OpenFOAM | LoRA fine-tune | Qwen3-8B | CFDLLMBench | 26.36\% exec | Local |
| MooseAgent | 2025 | MOOSE (FEM) | LangGraph, ~5 agents | DeepSeek-R1+V3 | 9 cases | 93\% avg | <\$0.14/case |
| MCP-SIM | 2025 | FEniCS (FEM) | 6 agents + shared mem | Multi-agent | 12 tasks | 100\% | Cloud |
| Pasimodo+RAG | 2025 | Pasimodo (SPH, closed) | Pure RAG Q\&A | Llama/Gemma 3--27B | 28 prompts | **0/2 model creation** | Local |
| **SloshAgent (Ours)** | **2026** | **DualSPHysics v5.4** | **14 tools + ReAct + MCP** | **Qwen3 32B (local)** | **SPHERIC Test 10** | **$r{>}0.9$ (exp. valid.)** | **\$0 LLM** |

Several critical observations emerge from this comparison:

**All systems target mesh-based solvers.** Ten of eleven systems use OpenFOAM; the remaining two use finite element solvers (MOOSE, FEniCS). The entire column of Lagrangian particle methods --- SPH, DEM, MPM --- is vacant. Pasimodo+RAG \cite{pasimodo2025} is the only SPH-adjacent work, but it is a pure retrieval-augmented Q\&A system that explicitly scored 0/2 on model creation tasks and does not qualify as a simulation agent.

**Success metrics are not comparable across systems.** MetaOpenFOAM reports 85\% Pass@1 on its own 8-case benchmark but achieves only 55.5\% on Foam-Agent's benchmark and 6.2\% on ChatCFD's 315-case benchmark \cite{chatcfd2026}. ChatCFD's "physical fidelity" score (68.12\%) uses LLM-as-judge evaluation, not comparison to experimental measurements. OpenFOAMGPT's 100\% reproducibility measures identical re-runs, verified only against analytical solutions (e.g., Poiseuille flow). These disparities highlight the absence of a unified validation methodology.

**No system validates against experimental benchmark data.** This is the most significant gap. All existing success metrics measure whether generated code compiles and runs (execution success), whether outputs match other simulations or analytical solutions (reference comparison), or whether an LLM judge deems results physically plausible. None compares agent-generated simulation results against published experimental data with quantitative error metrics such as Pearson correlation or NRMSE. The ISOPE benchmark study \cite{isope2012} demonstrated that even expert CFD predictions frequently fall outside experimental uncertainty bounds, underscoring that execution success alone is an insufficient measure of simulation quality.

**No system ablates domain knowledge content.** Ablation studies in existing work consistently test architectural decisions: MetaOpenFOAM ablates RAG retrieval and reviewer agents; Foam-Agent tests the effect of removing its reviewer node; ChatCFD evaluates RAG module contributions; PhyNiKCE ablates instruction structure; and MDCrow compares prompt styles \cite{mdcrow2025}. None isolates the effect of domain-specific physics knowledge encoded in prompts. Whether sloshing resonance formulas, SPH parameter inference rules, or solver-specific syntax constraints actually improve agent performance over generic prompting has never been tested.

## 2.3 LLM Agents for Scientific Discovery

Beyond simulation-specific agents, a broader ecosystem of LLM agents for scientific discovery provides architectural context for SloshAgent.

**Chemistry and molecular dynamics.** ChemCrow \cite{bran2024chemcrow} demonstrated that a single ReAct agent with 18 chemistry tools can plan and execute multi-step chemical syntheses, achieving expert-level performance on drug discovery and materials design tasks (Nature Machine Intelligence, 2024). Coscientist \cite{boiko2023coscientist} autonomously designed, planned, and executed chemical experiments using robotic lab equipment (Nature, 2023). MDCrow \cite{mdcrow2025} extended this paradigm to molecular dynamics with 40+ tools for LAMMPS and GROMACS simulation setup, analysis, and debugging (ICLR 2025). These systems established the single-agent-with-many-tools pattern that SloshAgent adopts.

**Materials and mechanics.** MatSciAgent \cite{matsciagent2025} automates materials science workflows including literature review, hypothesis generation, and experimental design. MechAgents \cite{mechagents2024} applies LLM agents to finite element analysis, generating FEA code for structural mechanics problems. Neither addresses fluid dynamics or sloshing.

**General scientific agents.** The AI Scientist \cite{aiscientist2024} (Sakana AI) automates the full scientific research cycle from hypothesis to paper writing. ScienceAgentBench \cite{scienceagentbench2025} provides a standardized evaluation of 102 data-driven scientific discovery tasks, on which the best agent achieves only 32.4\% success (ICLR 2025). Recent surveys \cite{survey_agentic_ai2025, survey_agentic_science2025} comprehensively catalog the rapidly growing landscape but identify computational mechanics --- and sloshing in particular --- as an unaddressed domain.

SloshAgent draws architectural inspiration from ChemCrow's single-agent-with-many-tools paradigm and MDCrow's domain-specialized tool design, while addressing a fundamentally different physics domain (violent free-surface fluid dynamics) with a fundamentally different solver type (GPU-accelerated particle methods) and a validation methodology (experimental benchmark comparison) that no prior system has attempted.
