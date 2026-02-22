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

The scale of the sloshing problem is best illustrated by the liquefied natural gas (LNG) shipping industry. Over 770 LNG carriers are currently in service, with more than 400 on order at approximately \$250--269M each; Korean shipyards (HD Hyundai, Samsung Heavy Industries, Hanwha Ocean) hold roughly 70\% of global orders \cite{lloyd2024}. Each LNG membrane tank requires over 200 sloshing simulations for classification society certification \cite{sjtu_ofw13_2018}. Seoul National University's sloshing laboratory alone has accumulated over 20,000 hours of model tests and 540 TB of experimental data \cite{kim2019snu}. The maritime AI market reached \$4.13B in 2024 with a 40.6\% compound annual growth rate \cite{marketreport2024}, yet this investment focuses entirely on route optimization, fuel efficiency, and predictive maintenance. To date, to the best of our knowledge, no AI tools address sloshing simulation --- the very analysis that determines whether an LNG tank design is safe for operation.

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

**No system addresses sloshing.** To the best of our knowledge, neither mesh-based (OpenFOAM VOF) nor particle-based (SPH) sloshing simulation has been automated by any LLM agent. The domain where simulation automation is arguably most needed --- requiring 200+ cases per tank, costing weeks of expert time, and involving documented silent errors --- remains unaddressed.

**No system integrates a particle-based solver.** All existing LLM-CFD agents target mesh-based solvers (primarily OpenFOAM). The entire Lagrangian particle simulation column --- SPH, DEM (Discrete Element Method), and MPM (Material Point Method) --- appears vacant. The only SPH-related work, Pasimodo+RAG \cite{pasimodo2025}, is a pure question-answering system that scored 0/2 on model creation tasks, explicitly not an agent.

**No system validates against experimental benchmarks.** Across all ten systems surveyed, "success" is defined as execution completion (Foam-Agent: 88.2\%), LLM-as-judge physical fidelity (ChatCFD: 68.12\%), or reproducibility of known analytical solutions (OpenFOAMGPT). We find no system that compares its results against published experimental data with quantitative error metrics.

**No system ablates domain knowledge from prompts.** Existing ablation studies uniformly test architectural components --- RAG retrieval (MetaOpenFOAM, ChatCFD), reviewer agents (Foam-Agent), instruction structure (PhyNiKCE) --- but not the domain knowledge content encoded in prompts. Whether domain-specific physics formulas, solver constraints, and parameter inference rules actually improve agent performance remains untested.

SPH is the natural solver for sloshing because it handles violent free-surface fragmentation, wave breaking, and splashing without mesh distortion --- precisely the phenomena that make sloshing dangerous and difficult to simulate with Eulerian methods \cite{dominguez2022}. The absence of any LLM agent for particle-based simulation thus represents both a methodological gap and a missed opportunity for the domain that needs automation most.

## 1.4 Contributions

This paper presents SloshAgent, the first AI agent that automates the entire sloshing simulation pipeline from natural language input to experimentally validated results. Our contributions are:

1. **First sloshing simulation automation agent.** SloshAgent converts natural language descriptions into DualSPHysics XML configurations, executes GPU-accelerated SPH simulations, and performs automated post-processing and validation. For standard use cases (explicit parameters and domain-inference inputs), parameter accuracy reaches 82\% (95\% CI [0.64, 0.95]), with accuracy dropping significantly for paper-specific and edge-case scenarios ($p = 0.001$), revealing a clear complexity--accuracy boundary.

2. **First LLM agent for particle-based simulation.** We design 14 domain-specific tools for the DualSPHysics pipeline --- to our knowledge, the first tool interface for any Lagrangian particle solver (SPH, DEM, or MPM) --- including an IsError pattern for LLM self-correction of silent SPH failures.

3. **First experimental benchmark validation by an LLM-simulation agent.** We validate agent-generated simulations against SPHERIC Test 10 benchmark data (102-repeat experiments, two fluids): all seven detected water impact pressure peaks fall within the experimental $\pm 2\sigma$ scatter band, whose non-normal distribution (Shapiro-Wilk $p < 0.005$, CV = 20--36\%) confirms the appropriateness of band-based validation for stochastic impact flows.

4. **First domain prompt ablation for computational mechanics.** We conduct a controlled ablation of the SloshingCoderPrompt (136 lines encoding sloshing physics, SPH constraints, and XML syntax rules). The full prompt yields a medium effect size ($d = 0.58$) over a generic baseline on the 32B model; the 8B model exhibits an inversion where the full prompt degrades tool-calling reliability, revealing a model capacity threshold for domain prompt engineering.

5. **Industry proof-of-concept at zero LLM cost.** SloshAgent runs entirely on local hardware (Qwen3 32B via Ollama on an RTX 4090), demonstrating sloshing simulation automation without cloud API dependencies, including parametric study orchestration and industrial baffle design comparison (91.9\% amplitude reduction) from single natural language prompts.
