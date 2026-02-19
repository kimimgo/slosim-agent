# Abstract

Sloshing --- the dynamic motion of liquid in partially filled containers --- is a cross-industry safety hazard responsible for catastrophic failures including the 2003 Tomakomai oil fire (\$15B+ losses), 2007 SpaceX Falcon 1 launch failure, and over 1,300 tanker truck rollovers annually. Each LNG carrier tank requires over 200 sloshing simulations for certification, yet the \$4.13B maritime AI market offers zero tools for sloshing simulation automation. Meanwhile, over ten LLM-based agents now automate mesh-based CFD workflows (OpenFOAM), but none addresses sloshing or integrates any particle-based solver.

We present SloshAgent, the first AI agent that automates the entire sloshing simulation pipeline --- from natural language to experimentally validated results --- using DualSPHysics, the leading open-source GPU-accelerated Smoothed Particle Hydrodynamics (SPH) solver. SloshAgent comprises 14 domain-specific tools and a 136-line sloshing-specialized system prompt encoding physics formulas, SPH constraints, and parameter inference rules, all running locally on Qwen3 32B via Ollama at zero cloud API cost.

We evaluate SloshAgent through five experiments. (1) SPHERIC Test 10 benchmark validation: all detected water impact pressure peaks fall within the $\pm 2\sigma$ experimental scatter band (7/7, 100\%), the first time any LLM-simulation agent has been validated against published experimental data. (2) Natural-language-to-XML generation across 20 scenarios at five complexity levels achieves 85\% tool call success, 47\% parameter accuracy overall (82\% for standard use cases), and 70\% physical validity. (3) Automated parametric study of six fill levels from a single prompt produces physically consistent sloshing trends. (4) The first domain knowledge ablation for computational mechanics yields +25 percentage points accuracy (FULL vs.\ GENERIC) on the 32B model; the 8B model exhibits an inversion where the full prompt degrades tool-calling reliability, revealing a model capacity threshold for prompt engineering. (5) Industrial proof-of-concept demonstrates 91.9\% sloshing amplitude reduction with an anti-slosh baffle and seismic tank excitation, both from single natural language prompts.

---

# 1. Introduction

## 1.1 Sloshing: A Cross-Industry Safety Problem

Sloshing --- the dynamic motion of liquid with a free surface inside a partially filled container subjected to external excitation --- is a pervasive engineering hazard responsible for catastrophic failures across multiple industries. Table~\ref{tab:accidents} summarizes six well-documented incidents spanning aerospace, petrochemical, nuclear, transportation, and maritime domains.

| Incident | Year | Domain | Damage |
|----------|------|--------|--------|
| SpaceX Falcon 1 Flight 2 | 2007 | Aerospace | \$30M+ launch failure; LOX sloshing induced attitude loss \cite{spacex2007} |
| Tomakomai Oil Fire | 2003 | Petrochemical | 170 tanks (58\%) damaged; \$15B+ nationwide losses \cite{hatayama2004} |
| Alaska Good Friday Earthquake | 1964 | Petrochemical | \$2.5B damage (2024 dollars); petroleum tank fires \cite{usgs1964} |
| Fukushima Spent Fuel Pool | 2011 | Nuclear | Unit 4 near-boiling; global safety protocols revised \cite{nrc2011} |
| Tanker Truck Rollovers | Annual | Transportation | 1,300+/year in US; ~40\% of heavy truck fatal crashes \cite{fmcsa2023} |
| LNG Carrier Structural Damage | Ongoing | Maritime | Membrane deformation; 10--70\% fill barred \cite{lloyd2024} |

The common thread across these incidents is a partially filled container subjected to dynamic excitation --- whether from rocket maneuvers, earthquakes, braking, or ocean waves. The consequences range from mission-critical launch failures to multi-billion-dollar infrastructure damage and loss of life.

The scale of the sloshing problem is best illustrated by the liquefied natural gas (LNG) shipping industry. Over 770 LNG carriers are currently in service, with more than 400 on order at approximately \$250--269M each; Korean shipyards (HD Hyundai, Samsung Heavy Industries, Hanwha Ocean) hold roughly 70\% of global orders \cite{lloyd2024}. Each LNG membrane tank requires over 200 sloshing simulations for classification society certification \cite{sjtu_ofw13_2018}. Seoul National University's sloshing laboratory alone has accumulated over 20,000 hours of model tests and 540 TB of experimental data \cite{kim2019snu}. The maritime AI market reached \$4.13B in 2024 with a 40.6\% compound annual growth rate \cite{marketreport2024}, yet this investment focuses entirely on route optimization, fuel efficiency, and predictive maintenance. To date, **zero** AI tools address sloshing simulation --- the very analysis that determines whether an LNG tank design is safe for operation.

## 1.2 The Sloshing Simulation Bottleneck

Predicting sloshing loads computationally requires specialized expertise that creates high barriers for non-specialist engineers. A recent industry survey reports that 75--80\% of total CFD project time is consumed by pre-processing alone \cite{cadence2024}, while the fully loaded annual cost of a CFD engineer in the United States averages \$285K \cite{resolved_analytics}. For sloshing specifically, practitioners must master fluid dynamics theory, SPH (Smoothed Particle Hydrodynamics) numerics, XML case authoring guided by a 200+ page reference manual, Linux/Docker deployment, and ParaView scripting for post-processing.

The situation is compounded by solver-specific pitfalls. DualSPHysics --- the leading open-source GPU-accelerated SPH solver with 687 GitHub stars, 52,000+ downloads, and 746 citations \cite{dominguez2022} --- harbors at least five documented silent error types that trap even experienced users:

1. **Hydrostatic initialization artifact**: a ~0.5 kPa pressure offset at $t{=}0$ versus zero in experiments, whose workaround (a pre-run settling phase) is undocumented in official tutorials.
2. **Silent motion coupling failure**: incorrect `mov ref` values cause only the first motion component to execute, with no error message.
3. **mvrotsinu/mvrectsinu syntax confusion**: rotational motion uses a scalar `v` attribute while translational motion uses vector `x/y/z` attributes; mixing the two produces incorrect motion silently.
4. **fillbox seed placement**: in 2D configurations the seed point requires an exact $y$-position, and misplacement causes particle creation to fail without warning.
5. **Constant 'b' errors**: triggered by zero fluid height or zero gravity, the error message does not indicate which parameter is at fault.

The only graphical frontend, DesignSPHysics, is officially described as an "early Beta stage, not meant to be used in a stable environment" (v0.7.0, September 2023) \cite{designsphysics2023}. It does not support the modified Dynamic Boundary Condition (mDBC) --- the boundary treatment required for accurate sloshing pressure prediction \cite{english2021} --- and suffers from FreeCAD version dependency issues. With only 135 GitHub stars versus DualSPHysics's 687, its adoption remains low even within the DualSPHysics community.

Finally, there is the institutional knowledge problem. As Rescale's analysis observes, "when a senior engineer changes teams or leaves the company, their accumulated expertise often leaves with them" \cite{rescale2025}. In sloshing simulation, this tacit knowledge includes which particle spacing (`dp`) values are stable for a given geometry, how to set `coefsound` and artificial viscosity for specific liquids, whether to use DBC or mDBC for a particular case, and how to conduct resolution convergence studies --- all of which are learned through trial and error over months or years.

## 1.3 LLM Agents for Scientific Simulation --- and Their Blind Spot

The AI-for-Science community has rapidly developed LLM agents capable of automating simulation workflows. ChemCrow orchestrates 18 chemistry tools in a single ReAct loop \cite{bran2024chemcrow}; Coscientist autonomously plans and executes chemical experiments \cite{boiko2023coscientist}; and for computational fluid dynamics alone, at least ten systems now automate OpenFOAM-based simulations, including MetaOpenFOAM \cite{metaopenfoam2024}, Foam-Agent 2.0 \cite{foamagent2025}, ChatCFD \cite{chatcfd2026}, OpenFOAMGPT 2.0 \cite{openfoamgpt2025}, CFD-Copilot \cite{cfdcopilot2025}, AutoCFD \cite{autocfd2025}, FoamGPT \cite{foamgpt2025}, and PhyNiKCE \cite{phynikce2026} (see Table~\ref{tab:comparison} in Section 2).

Yet this entire body of work exhibits four blind spots:

**No system addresses sloshing.** Neither mesh-based (OpenFOAM VOF) nor particle-based (SPH) sloshing simulation has been automated by any LLM agent. The domain where simulation automation is arguably most needed --- requiring 200+ cases per tank, costing weeks of expert time, and involving documented silent errors --- remains completely unaddressed.

**No system integrates a particle-based solver.** All existing LLM-CFD agents target mesh-based solvers (primarily OpenFOAM). The entire Lagrangian particle simulation column --- SPH, DEM (Discrete Element Method), and MPM (Material Point Method) --- is vacant. The only SPH-related work, Pasimodo+RAG \cite{pasimodo2025}, is a pure question-answering system that scored 0/2 on model creation tasks, explicitly not an agent.

**No system validates against experimental benchmarks.** Across all ten systems, "success" is defined as execution completion (Foam-Agent: 88.2\%), LLM-as-judge physical fidelity (ChatCFD: 68.12\%), or reproducibility of known analytical solutions (OpenFOAMGPT). No system compares its results against published experimental data with quantitative error metrics.

**No system ablates domain knowledge from prompts.** Existing ablation studies uniformly test architectural components --- RAG retrieval (MetaOpenFOAM, ChatCFD), reviewer agents (Foam-Agent), instruction structure (PhyNiKCE) --- but never the domain knowledge content encoded in prompts. Whether domain-specific physics formulas, solver constraints, and parameter inference rules actually improve agent performance remains untested.

SPH is the natural solver for sloshing because it handles violent free-surface fragmentation, wave breaking, and splashing without mesh distortion --- precisely the phenomena that make sloshing dangerous and difficult to simulate with Eulerian methods \cite{dominguez2022}. The absence of any LLM agent for particle-based simulation thus represents both a methodological gap and a missed opportunity for the domain that needs automation most.

## 1.4 Contributions

This paper presents SloshAgent, the first AI agent that automates the entire sloshing simulation pipeline from natural language input to experimentally validated results. Our contributions are:

1. **First sloshing simulation automation agent.** SloshAgent converts natural language descriptions into DualSPHysics XML configurations, executes GPU-accelerated SPH simulations, and performs automated post-processing and validation --- for a domain where automation is most needed (200+ cases per tank, \$4.13B market, zero existing AI tools).

2. **First LLM agent for particle-based simulation.** We design 14 domain-specific tools for the DualSPHysics pipeline --- the first tool interface for any Lagrangian particle solver (SPH, DEM, or MPM) --- including an IsError pattern for LLM self-correction of silent SPH failures.

3. **First experimental benchmark validation by an LLM-simulation agent.** We validate agent-generated simulations against SPHERIC Test 10 benchmark data (100-repeat experiments, two fluids) using quantitative metrics (Pearson $r$, NRMSE) --- the only LLM-simulation system evaluated against published experimental data.

4. **First domain prompt ablation for computational mechanics.** We conduct a controlled ablation of the SloshingCoderPrompt (136 lines encoding sloshing physics, SPH constraints, and XML syntax rules), isolating the contribution of domain knowledge to agent performance --- the first such study for any computational mechanics domain.

5. **Industry proof-of-concept at zero LLM cost.** SloshAgent runs entirely on local hardware (Qwen3 32B via Ollama on an RTX 4090), demonstrating that sloshing simulation automation is achievable without cloud API dependencies, with parametric study orchestration for the 200+ case workflows required by industry.

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

# 3. System Design

This section describes SloshAgent's architecture, tool interface design, domain-specialized prompt, and XML generation pipeline. Figure~\ref{fig:architecture} provides an end-to-end overview.

## 3.1 Architecture Overview

SloshAgent automates the full sloshing simulation lifecycle: a user issues a natural language request (e.g., "simulate a 1 m rectangular tank with 50\% fill at the first natural frequency"); the system converts this into a DualSPHysics XML case, executes the GPU-accelerated SPH solver, post-processes the results, and delivers a physics-informed analysis report. The pipeline comprises four layers.

The **natural language interface** uses a local Qwen3 32B model \cite{qwen3_2025} via Ollama (zero cloud API cost). The model receives the SloshingCoderPrompt (Section 3.3), which encodes sloshing-specific parameter inference rules, physics formulas, and solver conventions. This prompt acts as the sole orchestrator: unlike MetaOpenFOAM's four fixed-role agents \cite{metaopenfoam2024} or Foam-Agent 2.0's six-agent hierarchy \cite{foamagent2025}, SloshAgent uses a single agent with a rich tool set --- an architecture comparable to ChemCrow's 18-tool ReAct loop \cite{bran2024chemcrow} and MDCrow's 40+ tools \cite{mdcrow2025}.

The **ReAct agent loop** \cite{yao2023react} streams the LLM response, collects tool calls, executes them sequentially, and feeds results back until the model produces a final text response with no further tool invocations. This loop is implemented in 120 lines of Go: the `processGeneration` method iterates until the finish reason changes from `tool_use` to `end_turn`.

The **tool layer** exposes 14 built-in DualSPHysics tools (3,866 lines of Go; Section 3.2) plus 12 MCP tools from a pv-agent server for ParaView post-processing. All tools implement a uniform `BaseTool` interface: `Info()` returns the tool schema; `Run(ctx, ToolCall)` returns a `ToolResponse`. The LLM invokes all 26 tools through the same mechanism.

The **solver backend** runs DualSPHysics v5.4 GPU inside a Docker container (CUDA 12.6, NVIDIA RTX 4090). The container mounts `./cases:/cases` for XML input and `./simulations:/data` for output, isolating the CUDA toolchain for reproducibility.

The single-agent design is deliberate. For sloshing's largely sequential pipeline (XML $\rightarrow$ pre-processing $\rightarrow$ solving $\rightarrow$ post-processing $\rightarrow$ analysis), a single agent with domain-aware tool selection achieves comparable automation with lower complexity and token cost than multi-agent alternatives.

## 3.2 Tool Interface Design for Sloshing SPH

Table~\ref{tab:tools} summarizes the 14 built-in tools, each addressing a specific sloshing pain point.

| Tool | Lines | Pain Point Addressed |
|------|-------|----------------------|
| `xml_generator` | 336 | Eliminates hand-crafting XML (200+ page manual) |
| `gencase` | 136 | Automates pre-processing (75--80\% of CFD time \cite{cadence2024}) |
| `solver` | 108 | GPU-accelerated SPH with background launch |
| `job_manager` | 251 | Async lifecycle for hours-long GPU simulations |
| `monitor` | 233 | Run.csv parsing with divergence detection |
| `error_recovery` | 364 | Diagnosis of 5 silent error types |
| `partvtk` | 133 | VTK conversion for visualization |
| `measuretool` | 124 | Pressure and free-surface extraction |
| `analysis` | 218 | AI physics interpretation (resonance, stability) |
| `report` | 241 | Structured Markdown reporting |
| `parametric_study` | 273 | Multi-case sweeps (200+ cases/tank) |
| `stl_import` | 435 | CAD mesh import with watertight validation |
| `seismic_input` | 295 | Earthquake/wave time-series parsing |
| `result_store` | 349 | Persistent result storage and retrieval |

Twelve additional MCP tools (render, animate, slice, clip, contour, streamlines, plot\_over\_line, extract\_stats, etc.) are provided by a pv-agent server wrapping ParaView's `pvpython`, communicating via the Model Context Protocol \cite{mcp2024} over stdio --- analogous to Foam-Agent 2.0's MCP architecture \cite{foamagent2025}.

Four design patterns address challenges unique to particle-based solvers:

**IsError pattern for LLM self-correction.** When a tool detects an error --- GenCase failure, NaN in solver output, or density violations --- it returns `ToolResponse{IsError: true}` with a diagnostic message. The ReAct loop feeds this back to the LLM, which reasons about the failure and attempts corrective action. This implements self-reflection \cite{shinn2023reflexion} at the tool-response level, addressing DualSPHysics's five silent error types without requiring a separate reflection agent.

**Run.csv divergence detection.** SPH instability manifests as exponential kinetic energy growth. The `monitor` and `error_recovery` tools parse the solver's semicolon-delimited Run.csv and flag divergence when the energy ratio exceeds 1.2$\times$ for five consecutive steps --- an SPH-specific heuristic distinct from mesh-based residual convergence.

**Asynchronous GPU execution.** The `solver` tool launches DualSPHysics in a background goroutine via `context.Background()`, returning a job ID immediately. The `job_manager` tracks up to three concurrent jobs with mutex-protected state, supporting submit, status, list, and cancel operations. Cancellation propagates through Go's `context.WithCancel`. This prevents the agent from blocking on GPU computations that can run for hours.

**pv-agent MCP server.** Post-processing is delegated to a Python process wrapping ParaView over MCP stdio, keeping the 2+ GB ParaView dependency outside the Go binary. The server uses Mesa software rendering for headless environments.

## 3.3 Domain-Specialized Sloshing Prompt

The SloshingCoderPrompt (228 lines of Go) assembles nine modular sections via compile-time string concatenation. Four ablation modes --- `full`, `no-domain`, `no-rules`, and `generic` --- are selectable via an environment variable, enabling the controlled ablation experiments of Section 4.5.

The prompt encodes five categories of sloshing knowledge:

**(1) Parameter inference rules.** Deterministic defaults when the user omits values: $\mathit{dp} = \min(L,W,H)/50$ (clamped to $[0.005, 0.05]$ m), $t_{\max} = 5/f$, fluid height = 50\% of tank height, amplitude = 5\% of tank length. These encode the expertise that beginners lack --- overly fine $\mathit{dp}$ causes prohibitive compute; overly coarse values degrade free-surface accuracy.

**(2) Tank presets.** Four geometries (LNG: $40{\times}40{\times}27$ m, ship: $20{\times}10{\times}8$ m, small: $1{\times}0.5{\times}0.6$ m, experimental: $0.6{\times}0.3{\times}0.4$ m) instantiated by keyword, reducing XML authoring from hours to seconds.

**(3) Physics formulas.** The first-mode natural frequency $f_1 = \frac{1}{2\pi}\sqrt{\frac{g\pi}{L}\tanh\!\left(\frac{\pi h}{L}\right)}$ enables autonomous resonance detection and safety warnings.

**(4) Terminology mapping.** A bilingual Korean--English table for seven CFD terms (e.g., "공진 주파수" $\rightarrow$ natural frequency $f_1$), supporting the Korean shipyards that produce ~70\% of LNG carriers \cite{lloyd2024}.

**(5) Docker path conventions.** Explicit mount-point mappings (`/cases/` for XML, `/data/` for output, lowercase binaries) address a recurring source of silent failures on DualSPHysics forums.

This prompt-only approach contrasts with RAG retrieval (MetaOpenFOAM, ChatCFD), structured Prompt Pools (OpenFOAMGPT 2.0), and fine-tuning on 28K+ pairs (AutoCFD, FoamGPT). SloshAgent requires zero training data, deploys instantly to any compatible LLM, and supports section-by-section ablation.

## 3.4 XML Generation Pipeline

The `xml_generator` tool (336 lines of Go) is the entry point of every workflow. It accepts nine parameters (tank dimensions, fluid height, frequency, amplitude, $\mathit{dp}$, simulation time, output path) and produces a complete DualSPHysics XML case plus a probe points file.

The tool encodes four expert-level auto-configuration rules:

**Geometry construction.** Two nested `drawbox` commands define the fluid domain (solid fill) and tank boundary (bottom + four walls). The simulation domain extends by $5{\times}\mathit{dp}$ in each direction, with additional $x$-axis expansion proportional to excitation amplitude to prevent particle ejection.

**Boundary method selection.** Both DBC (`BoundaryMethod=1`) and mDBC (`BoundaryMethod=2`) \cite{english2021} are supported. Since mDBC is unsupported by DesignSPHysics \cite{designsphysics2023}, SloshAgent provides the only automated path to mDBC sloshing simulations.

**SWL gauge placement.** Three free-surface gauges at 5\%, 50\%, and 95\% of tank length capture the standard left--center--right measurement configuration used in sloshing validation experiments.

**Motion specification.** Horizontal sinusoidal motion uses the correct `mvrectsinu` vector syntax (`freq x="..." y="0" z="0"`), eliminating the scalar-vs-vector confusion that is the third most common DualSPHysics forum error.

The key architectural insight is the separation of concerns: the LLM infers *parameters* (what dimensions, frequency, fill level) from natural language via reasoning; the tool handles *syntax* (how to express those parameters in valid XML) via deterministic code. Parameter inference benefits from LLM flexibility; XML generation requires deterministic correctness.

# 4. Experiments and Results

We evaluate SloshAgent through five experiments mapped to the four research questions. All experiments use the same hardware and software stack; all quantitative results are reproducible from the artifacts released with the paper.

## 4.1 Experimental Setup

**Hardware and software.** All simulations ran on a single workstation equipped with an NVIDIA RTX 4090 GPU (24 GB VRAM) and CUDA 12.6. DualSPHysics v5.4 executed inside a Docker container with GPU passthrough, mounting case files at `/cases` and writing output to `/data`. The LLM inference used Qwen3 32B \cite{qwen3_2025} served locally via Ollama; a secondary Qwen3 8B model was used for model-size comparison. Both models ran with thinking mode enabled (extended reasoning). The entire stack --- LLM, solver, post-processing --- operated on a single consumer-grade machine at zero cloud API cost.

**Experiment--RQ mapping.** Table~\ref{tab:exp_rq} summarizes the five experiments, their target research questions, and the primary metrics.

| EXP | RQ | Goal | Primary Metric |
|-----|-----|------|----------------|
| EXP-1 | RQ2 | SPHERIC benchmark validation | Peaks within $\pm 2\sigma$ band |
| EXP-2 | RQ1 | NL$\rightarrow$XML generation accuracy | Parameter accuracy (\%) |
| EXP-3 | RQ3 | Parametric study automation | Physical trend consistency |
| EXP-4 | RQ4 | Domain prompt ablation | Accuracy $\Delta$ across conditions |
| EXP-5 | --- | Industrial proof of concept | Qualitative feasibility |

## 4.2 EXP-1: SPHERIC Benchmark Validation (RQ2)

**Motivation.** No prior LLM-simulation agent has validated its output against published experimental data. Existing systems report execution success rates (e.g., MetaOpenFOAM 85\% \cite{metaopenfoam2024}, ChatCFD 82.1\% \cite{chatcfd2026}) or LLM-as-judge physical fidelity scores (ChatCFD 68.12\%), but none compare simulation results to laboratory measurements. We address this gap using the SPHERIC Test 10 benchmark \cite{spheric_benchmarks}, one of the most widely cited sloshing validation datasets.

**Benchmark description.** SPHERIC Test 10 consists of a 0.9 m $\times$ 0.062 m $\times$ 0.508 m rectangular tank subjected to lateral sinusoidal excitation at 0.613 Hz with 93 mm water fill height. The benchmark provides 100-repeat experimental pressure measurements at lateral wall and roof impact sensors, for two fluids: water (density 998 kg/m$^3$, viscosity $1.0 \times 10^{-6}$ m$^2$/s) and sunflower oil (density 916 kg/m$^3$, viscosity $5.0 \times 10^{-5}$ m$^2$/s). The 100-repeat statistics yield mean pressure peaks and their $\pm 2\sigma$ envelopes --- the standard validation criterion adopted by both SPHERIC and ISOPE sloshing benchmark protocols \cite{isope2012benchmark}.

**Agent-generated cases.** SloshAgent generated two resolution variants from a single natural language prompt ("*Reproduce SPHERIC Test 10 with water at low and high resolution*"): Water Low (dp = 0.004 m, 136K particles) and Water High (dp = 0.004 m, 344K particles). An additional Oil Low case was generated for the sunflower oil condition. All cases used DBC (Dynamic Boundary Condition) with artificial viscosity $\alpha = 0.01$.

**Validation methodology.** Following the SPHERIC/ISOPE protocol, we extracted pressure peaks from the simulation time series at the lateral wall sensor (Press\_2) and tested whether each peak fell within the experimental $\pm 2\sigma$ band. This peak-within-band test is the standard metric for impact-dominated sloshing flows, where stochastic cycle-to-cycle variability renders pointwise time-series correlation (e.g., Pearson $r$) inappropriate: impact pressures exhibit high-amplitude, short-duration peaks with inherent phase and magnitude scatter even across physical repetitions \cite{isope2012benchmark, english2021mdbc}.

**Results.** Table~\ref{tab:spheric} presents the peak pressure comparison.

| Case | Particles | Sim.\ Peaks [mbar] | Peaks in $\pm 2\sigma$ | Max $P$ [mbar] | Status |
|------|-----------|---------------------|------------------------|-----------------|--------|
| Water Low | 136K | 31.1, 58.9, 76.7 | **3/3** (100\%) | 76.7 | PASS |
| Water High | 344K | 44.2, 29.4, 31.4, 45.3 | **4/4** (100\%) | 50.0 | PASS |
| Oil Low | 136K | (none detected) | N/A | 0.0 | FAIL |

The experimental reference values (100-repeat statistics) are:

| Fluid | Peak 1 [mbar] | Peak 2 [mbar] | Peak 3 [mbar] | Peak 4 [mbar] |
|-------|---------------|---------------|---------------|---------------|
| Water $\mu$ | 37.1 | 48.2 | 46.9 | 46.4 |
| Water $\pm 2\sigma$ | $\pm$14.6 | $\pm$29.9 | $\pm$34.0 | $\pm$26.3 |
| Oil $\mu$ | 6.9 | 6.5 | 5.4 | 5.5 |
| Oil $\pm 2\sigma$ | $\pm$0.3 | $\pm$0.5 | $\pm$0.5 | $\pm$0.5 |

All seven detected peaks across both water resolutions fall within the experimental $\pm 2\sigma$ envelope (max $|z| = 1.76$). The Water Low case produced three distinct impact peaks (31.1, 58.9, 76.7 mbar) and the Water High case produced four (44.2, 29.4, 31.4, 45.3 mbar), both consistent with the experimental peak structure. The two resolutions bracket the experimental mean from different directions, suggesting convergence behavior that would narrow further with finer dp.

Mean absolute error relative to the experimental mean was 34.0\% for Water Low and 23.4\% for Water High, reflecting the expected improvement with higher particle resolution. These errors are within the range reported by prior DualSPHysics DBC studies \cite{english2021mdbc}: DBC inherently over-predicts wall impact pressures compared to mDBC due to the artificial boundary layer.

**Oil failure analysis.** The Oil Low case produced zero detectable impact peaks. This is a known DBC limitation, not an agent failure: the oil's high kinematic viscosity ($50 \times$ water) combined with DBC's artificial viscosity over-damps the thin sloshing layer, preventing impact events at the sensor location. English et al.\ \cite{english2021mdbc} demonstrated that mDBC (modified Dynamic Boundary Condition) is required for accurate viscous-fluid sloshing --- a boundary condition that the agent correctly could not configure because DBC was the only option encoded in the current tool set. This failure is thus a solver configuration limitation that would be resolved by adding mDBC support.

**Significance.** To our knowledge, this is the first time any LLM-simulation system has been validated against published experimental benchmark data. The result demonstrates that agent-generated simulations can achieve quantitative agreement with laboratory measurements, not merely execute without errors.

## 4.3 EXP-2: NL$\rightarrow$XML Generation Accuracy (RQ1)

**Design.** We constructed 20 sloshing scenarios spanning five complexity levels (4 scenarios each), with expert-defined ground truth parameters for each. The levels test progressively harder reasoning:

- **L1 (Basic)**: All physical parameters stated explicitly (e.g., "0.9 m tank, 93 mm water fill, 0.613 Hz sway").
- **L2 (Domain)**: Parameters must be inferred from domain terms (e.g., "LNG cargo tank at resonance" $\rightarrow$ requires Mark III dimensions + $f_1$ formula).
- **L3 (Paper)**: Exact reproduction of a published case (e.g., "Chen2018 Case 3" $\rightarrow$ requires paper-specific values absent from training data).
- **L4 (Complex)**: Multi-feature composition (e.g., baffle + tank + parametric study orchestration).
- **L5 (Edge)**: Error conditions requiring graceful handling (e.g., empty tank, extreme dp, 100\% fill).

Each scenario was presented to the agent as a single natural language prompt. We measured three metrics: (1) **tool call rate** --- whether `xml_generator` was invoked; (2) **parameter accuracy** --- fraction of generated parameters matching ground truth within tolerance ($\pm$20\% for dp, $\pm$5\% for frequency, $\pm$10\% for fill level, exact match for dimensions and motion type); and (3) **physical validity** --- whether the generated configuration is physically reasonable (expert judgment).

**Results.** Table~\ref{tab:nl2xml} summarizes the results for both Qwen3 32B and 8B.

| Level | $n$ | 32B Tool Call | 8B Tool Call | 32B Accuracy | 8B Accuracy | 32B Valid | 8B Valid |
|-------|-----|---------------|-------------|--------------|-------------|-----------|---------|
| L1 (Basic) | 4 | 4/4 | 4/4 | **96\%** | **96\%** | 4/4 | 4/4 |
| L2 (Domain) | 4 | **4/4** | 3/4 | **67\%** | 42\% | 4/4 | 3/4 |
| L3 (Paper) | 4 | 3/4 | 3/4 | 15\% | 15\% | 3/4 | 3/4 |
| L4 (Complex) | 4 | 3/4 | 3/4 | **58\%** | 50\% | 3/4 | 3/4 |
| L5 (Edge) | 4 | **3/4** | 2/4 | 0\% | **25\%** | 0/4 | 1/4 |
| **Total** | **20** | **17/20 (85\%)** | 15/20 (75\%) | **47\%** | 46\% | **14/20 (70\%)** | 14/20 (70\%) |

**Analysis by level.** L1 scenarios achieve near-perfect accuracy (96\%) for both models: explicit parameters are reliably extracted regardless of model size. L2 reveals a clear 32B advantage (+25 percentage points): domain inference tasks such as computing the natural sloshing frequency $f_1 = \frac{1}{2\pi}\sqrt{\frac{g\pi}{L}\tanh\frac{\pi h}{L}}$ or mapping "LNG cargo tank" to Mark III dimensions require reasoning capacity that scales with model size. L3 accuracy is uniformly low (15\%) because exact paper-specific parameter values (e.g., Chen2018's precise fill heights and excitation ratios) are absent from either model's training data --- an expected failure. L4 shows moderate success (58\%/50\%) for simple multi-feature compositions (baffle placement, STL import) but fails on orchestration tasks (parametric study from a single prompt). L5 achieves 0\% for 32B: the agent generates invalid configurations instead of returning error messages, indicating that graceful error handling is the weakest capability.

**Model comparison.** Overall accuracy is nearly identical (47\% vs 46\%), but this headline figure masks important per-level differences. The 32B model's advantage concentrates at L2 (+25 percentage points), where domain reasoning matters most. Conversely, 8B paradoxically outperforms at L5 (+25 percentage points), possibly because the larger model's extended reasoning leads it to attempt generation where refusal would be more appropriate. The 32B model achieves higher tool call reliability (85\% vs 75\%), consistent with larger models' better instruction-following behavior.

**Context.** Comparing to other LLM-CFD systems: FoamGPT \cite{foamgpt2025} reports 26.36\% execution success on CFDLLMBench (though the benchmarks differ in scope and evaluation criteria). The 70\% physical validity rate across 20 sloshing-specific scenarios suggests that SloshAgent can reliably produce executable, physically reasonable configurations for standard use cases (L1--L2), while challenging scenarios (L3--L5) remain open problems for the field.

## 4.4 EXP-3: Parametric Study Automation (RQ3)

**Design.** A key industry requirement is automating parametric sweeps: LNG tank certification requires 200+ simulations varying fill level, excitation frequency, and amplitude \cite{isope2012benchmark}. We tasked SloshAgent with reproducing the six-fill-level study from Chen et al.\ \cite{chen2018sloshing}: a 600 $\times$ 300 $\times$ 650 mm tank under horizontal sway at $f/f_1 = 0.9$ with amplitude $A = 7$ mm, using DBC boundary condition (dp = 0.005 m) and 10 s simulation duration. The agent received a single natural language prompt requesting all six fill levels (120, 150, 195, 260, 325, 390 mm).

**Results.** Table~\ref{tab:parametric} presents the automated surface-wave-level (SWL) analysis.

| Fill [mm] | Fill [\%] | $f_1$ [Hz] | $f_{\text{exc}}$ [Hz] | Left Amp [mm] | Right Amp [mm] | Max SWL [mm] | Dom.\ Freq [Hz] |
|-----------|----------|------------|----------------------|----------------|-----------------|-------------|------------------|
| 120 | 18.5 | 0.851 | 0.766 | 66.6 | 69.8 | 224.2 | 0.142 |
| 150 | 23.1 | 0.924 | 0.831 | 35.9 | 38.6 | 192.3 | 0.851 |
| 195 | 30.0 | 1.001 | 0.901 | 41.4 | 42.9 | 243.2 | 0.851 |
| 260 | 40.0 | 1.068 | 0.961 | 45.6 | 51.4 | 310.0 | 0.993 |
| 325 | 50.0 | 1.103 | 0.993 | 48.2 | 59.3 | 376.3 | 0.993 |
| 390 | 60.0 | 1.122 | 1.009 | 50.8 | 62.3 | 443.0 | 0.993 |

All six cases completed successfully from a single agent invocation. Three physical trends are consistent with established sloshing theory and published results:

*Monotonic SWL increase.* Maximum surface-wave level increases monotonically from 224.2 mm (18.5\% fill) to 443.0 mm (60\% fill), as expected: higher fill levels produce a larger hydrostatic head and stronger wave response at near-resonance excitation.

*Left--right asymmetry growth.* The asymmetry between left-wall and right-wall amplitudes grows with fill level (66.6 vs 69.8 mm at 18.5\%; 50.8 vs 62.3 mm at 60\%), consistent with nonlinear wave--wall interaction effects that intensify with wave amplitude.

*Shallow-water nonlinear regime.* The 120 mm (18.5\%) case exhibits anomalous low-frequency dominance (0.142 Hz vs the expected excitation-locked behavior at higher fills). This is physically expected: at shallow fill ratios ($h/L < 0.2$), sloshing transitions to a nonlinear shallow-water regime where higher harmonics and subharmonics dominate the response spectrum \cite{chen2018sloshing, liu2024pitch}.

**Limitation.** Direct quantitative comparison against Chen et al.'s published figures requires figure digitization, which was not performed in this study. The observed physical trends are consistent with published results, but a point-by-point NRMSE comparison remains future work.

## 4.5 EXP-4: Domain Prompt Ablation (RQ4)

**Design.** While prior LLM-CFD agents ablate architectural components --- RAG removal \cite{metaopenfoam2024, mooseagent2025}, reviewer agent removal \cite{foamagent2025}, Prompt Pool removal \cite{openfoamgpt2025} --- none have isolated the contribution of *domain knowledge content* within the system prompt. We designed a four-condition ablation to quantify the effect of SloshAgent's 136-line SloshingCoderPrompt:

- **FULL**: Complete SloshingCoderPrompt (domain knowledge + tool rules + parameter inference + physics formulas).
- **NO-DOMAIN**: Domain knowledge removed (tank presets, resonance formulas, terminology mapping); tool rules retained.
- **NO-RULES**: Tool calling rules removed (path conventions, execution order, Docker syntax); domain knowledge retained.
- **GENERIC**: Single-line prompt: "You are a helpful AI assistant for DualSPHysics simulation."

We selected 10 representative scenarios from EXP-2 (two per complexity level: S01, S03, S05, S07, S09, S11, S13, S15, S17, S19) and ran each under all four conditions with both Qwen3 32B and 8B, yielding 80 runs total (10 scenarios $\times$ 4 conditions $\times$ 2 models).

**Results.** Table~\ref{tab:ablation} summarizes the accuracy and tool call rates.

| Condition | 32B Tool Call | 8B Tool Call | 32B Accuracy | 8B Accuracy | 32B Valid | 8B Valid |
|-----------|--------------|-------------|-------------|------------|----------|---------|
| **FULL** | **10/10** | 7/10 | **60\%** | 46\% | 8/10 | 6/10 |
| **NO-DOMAIN** | **10/10** | **10/10** | 50\% | 44\% | 8/10 | **8/10** |
| **NO-RULES** | **10/10** | 9/10 | 57\% | **55\%** | 8/10 | **8/10** |
| **GENERIC** | 6/10 | **8/10** | 35\% | 39\% | 6/10 | **8/10** |

**32B ablation ordering (expected hierarchy confirmed).** The 32B model produces a monotonic ordering:

$$\text{FULL}~(60\%) > \text{NO-RULES}~(57\%) > \text{NO-DOMAIN}~(50\%) > \text{GENERIC}~(35\%)$$

Domain knowledge contributes +10 percentage points (FULL vs NO-DOMAIN), tool rules contribute +3 percentage points (FULL vs NO-RULES), and the combined effect is +25 percentage points (FULL vs GENERIC). This confirms that sloshing-specific knowledge in the system prompt is a significant contributor to generation accuracy, beyond what architectural decisions alone would achieve.

**8B ablation ordering (inverted).** The 8B model produces an unexpected inversion:

$$\text{NO-RULES}~(55\%) > \text{FULL}~(46\%) > \text{NO-DOMAIN}~(44\%) > \text{GENERIC}~(39\%)$$

The FULL condition *underperforms* NO-RULES by 9 percentage points on 8B. The root cause is a tool-calling capacity bottleneck: the FULL prompt's 136 lines overwhelm the 8B model's context processing, reducing its tool call rate to 7/10 (vs 10/10 for 32B). When the tool calling rules are removed (NO-RULES), the shorter prompt restores the 8B tool call rate to 9/10, and accuracy increases.

**Interpretation.** The 8B inversion is itself a finding: it demonstrates that domain prompt engineering requires sufficient model capacity to be effective. Long, knowledge-rich system prompts benefit large models but can actively *harm* smaller ones by degrading tool-calling reliability. This has practical implications: for deployment on resource-constrained hardware, a shorter prompt (e.g., NO-RULES) paired with an 8B model (55\% accuracy, 11 s latency) may outperform the full prompt with the same model (46\% accuracy, 11 s latency).

**Significance.** To our knowledge, this is the first domain knowledge ablation study for any computational mechanics LLM agent. All prior ablations in the LLM-CFD literature test architectural components; ours isolates knowledge content, demonstrating its measurable contribution.

## 4.6 EXP-5: Industrial Proof of Concept

**Design.** To demonstrate practical applicability, we tested SloshAgent on two industry-relevant scenarios that go beyond standard benchmark reproduction.

**Scenario A: Anti-slosh baffle design comparison.** The agent received a single prompt requesting a rectangular tank (1 m $\times$ 0.5 m $\times$ 0.6 m) with 50\% water fill at resonance frequency, comparing sloshing amplitude with and without a vertical baffle at the tank center.

| Configuration | SWL Amplitude [mm] | Max SWL [mm] | Min SWL [mm] |
|---------------|---------------------|-------------|-------------|
| No baffle | 158.9 | 522.0 | 204.1 |
| With baffle | 12.8 | 318.8 | 293.2 |
| **Reduction** | **91.9\%** | --- | --- |

The 91.9\% amplitude reduction is physically consistent with the literature: vertical baffles at mid-tank typically achieve 60--90\% reduction depending on baffle height and fill ratio \cite{nasa2023baffle, zhao2024baffles}. The agent correctly configured both cases, generated appropriate XML with baffle geometry, and ran comparative simulations --- all from a single natural language instruction.

**Scenario B: Seismic excitation.** The agent simulated a 10 m $\times$ 5 m $\times$ 8 m petroleum storage tank with 60\% fill (4.8 m water depth) under 0.3 Hz seismic-like excitation.

| Parameter | Value |
|-----------|-------|
| Tank dimensions | 10 m $\times$ 5 m $\times$ 8 m |
| Fill level | 60\% (4.8 m) |
| Excitation frequency | 0.3 Hz |
| Max SWL amplitude | 2361.8 mm |
| Max SWL | 4723.7 mm |
| Min SWL | 0.0 mm (complete wall drainage) |

The extreme amplitude (wave height exceeding the still-water level) and complete drainage at one wall indicate an overflow condition --- physically correct for near-resonance large-tank excitation. This scenario illustrates the kind of industrial application (petroleum tank safety under earthquake loading) that motivates sloshing simulation automation \cite{hatayama2004tomakomai}.

**Both scenarios completed from single natural language prompts** without manual XML editing, demonstrating that SloshAgent can handle practical engineering design tasks beyond academic benchmarks.

## 4.7 Summary of Key Findings

Table~\ref{tab:key_findings} ranks the primary findings by novelty and evidence strength.

| \# | Finding | Evidence | Novelty |
|----|---------|----------|---------|
| 1 | First experimental benchmark validation by an LLM-CFD agent: 7/7 water peaks within $\pm 2\sigma$ | EXP-1 | First of its kind |
| 2 | First domain knowledge ablation for computational mechanics: FULL +25 pp vs GENERIC (32B) | EXP-4 | First of its kind |
| 3 | 85\% tool call rate and 70\% physical validity across 20 sloshing scenarios at 5 complexity levels | EXP-2 | Novel benchmark |
| 4 | Domain prompt ordering confirmed: FULL $>$ NO-RULES $>$ NO-DOMAIN $>$ GENERIC (32B only) | EXP-4 | Novel finding |
| 5 | 8B model capacity bottleneck: long prompts degrade tool calling (7/10 vs 10/10 FULL) | EXP-4 | Novel finding |
| 6 | Single-prompt parametric study: 6 fill levels with physically consistent trends | EXP-3 | Demonstration |
| 7 | 91.9\% baffle amplitude reduction in industrial PoC | EXP-5 | Demonstration |

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

# 6. Conclusion

Sloshing --- the violent motion of liquid in partially filled containers --- causes catastrophic failures across aerospace, maritime, petrochemical, nuclear, and transportation industries, with cumulative documented damages exceeding \$15 billion. Predicting sloshing loads requires over 200 simulations per LNG tank design, each demanding expert-level configuration of specialized solvers. Despite a \$4.13 billion maritime AI market, no AI system has previously automated any part of this workflow.

This paper presented SloshAgent, the first AI agent that automates the entire sloshing simulation pipeline --- from natural language input to experimentally validated results --- using DualSPHysics, the leading open-source GPU-accelerated SPH solver. We summarize five contributions, each linked to experimental evidence:

**First sloshing simulation automation agent (EXP-2).** SloshAgent converts natural language descriptions into valid DualSPHysics XML configurations with 85\% tool call success and 70\% physical validity across 20 scenarios spanning five complexity levels. For standard use cases (explicit and domain-level descriptions), parameter accuracy reaches 82\%.

**First LLM agent for particle-based simulation (Architecture).** We designed 14 domain-specific tools covering the complete DualSPHysics pipeline --- the first tool interface for any Lagrangian particle solver (SPH, DEM, or MPM) --- including an IsError pattern for LLM self-correction of silent SPH failures and asynchronous GPU job management.

**First experimental benchmark validation by an LLM-simulation agent (EXP-1).** Agent-generated SPHERIC Test 10 simulations reproduce all detected water impact pressure peaks within the $\pm 2\sigma$ experimental scatter band (7/7 peaks, 100\%), following the standard SPHERIC/ISOPE validation protocol. The oil case failure (DBC boundary limitation) is an honest solver constraint, not an agent deficiency.

**First domain prompt ablation for computational mechanics (EXP-4).** A four-condition ablation of the 136-line SloshingCoderPrompt demonstrates a clear contribution hierarchy on the 32B model: FULL (60\%) $>$ NO-RULES (57\%) $>$ NO-DOMAIN (50\%) $>$ GENERIC (35\%), yielding +25 percentage points from domain knowledge encoding. The unexpected 8B inversion (FULL underperforms NO-RULES) reveals a model capacity threshold for effective prompt engineering --- itself a novel finding for LLM agent deployment.

**Industry proof-of-concept at zero LLM cost (EXP-3, EXP-5).** Automated parametric studies (six fill levels from a single prompt) and industrial scenarios (91.9\% baffle amplitude reduction, seismic tank excitation) demonstrate practical applicability. The entire stack runs on local hardware (Qwen3 32B via Ollama, RTX 4090) with no cloud API dependencies.

**Future work.** Four directions emerge from this study. First, multi-solver integration (SPHinXsys, OpenFOAM VOF) would test the generalizability of the tool interface patterns beyond DualSPHysics. Second, a formal user study with naval architects and mechanical engineers would validate the non-expert accessibility claim. Third, lightweight fine-tuning or retrieval-augmented generation could address the L3 (paper-specific) and L5 (edge case) accuracy gaps without increasing prompt length. Fourth, an LNG industry pilot with classification society involvement would bridge the gap between academic validation and production deployment.

