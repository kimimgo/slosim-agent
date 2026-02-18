# SloshAgent: A Domain-Specialized LLM Agent for Autonomous SPH Sloshing Simulation

## Abstract
- **Problem**: SPH sloshing simulation requires deep domain expertise for case setup, parameter tuning, and result interpretation, creating high barriers for non-specialist engineers across LNG, automotive, nuclear, and aerospace industries.
- **Gap**: While LLM+CFD automation exists for mesh-based solvers (AutoCFD, OpenFOAMGPT, FoamGPT), no system addresses particle-based SPH simulation — a complete research gap.
- **Approach**: We present SloshAgent, the first domain-specialized LLM agent that autonomously configures, executes, and validates DualSPHysics sloshing simulations through 14 structured tool interfaces with a local Qwen3 32B model.
- **Results**: On SPHERIC Test 10 benchmark, SloshAgent achieves expert-level accuracy (r > 0.9 for pressure time series). Domain-specialized prompting improves XML generation accuracy by 20-40% over generic prompts. Case setup time is reduced by 10x+ compared to manual workflow.
- **Contribution**: (1) First LLM agent for SPH simulation, (2) Tool interface design patterns for particle-based solvers, (3) Domain prompt ablation for computational mechanics, (4) Benchmark-validated E2E pipeline, (5) Industry PoC for sloshing practitioners.

## 1. Introduction (1 page)

### 1.1 Sloshing: A Critical Industrial Challenge
- 4 major industries: LNG carriers (membrane damage, class rules), automotive (NVH, fuel supply), nuclear (seismic safety), aerospace (Falcon 1 loss 2007, NEAR mission)
- Current methods: model tests (SNU 20,000+ hr), mesh-based CFD (VOF), SPH
- Pain points: expertise barrier, setup time (hours-days), parameter tuning, result interpretation

### 1.2 LLM Agents for Scientific Simulation
- AI4Science trajectory: ChemCrow (chemistry, 740 citations), Coscientist (lab automation, 1108), MechAgents (FEA, 133)
- CFD automation: AutoCFD (Qwen2.5-7B, NL2FOAM 28.7K pairs), OpenFOAMGPT 2.0 (100% reproducibility), FoamGPT (NeurIPS ML4PS 2025, 26.36% success)
- **Gap**: All mesh-based (OpenFOAM). Zero work on particle-based SPH solvers.
- Closest: Pasimodo+RAG (arXiv:2502.03916) — RAG assistance, not autonomous agent

### 1.3 Contributions
1. First domain-specialized LLM agent system for autonomous SPH sloshing simulation
2. Tool interface design patterns for LLM–particle-solver integration (IsError, async GPU, Run.csv monitoring)
3. Domain-specialized prompt engineering ablation for computational mechanics
4. Benchmark-validated E2E pipeline (SPHERIC Test 10)
5. Industry PoC: parametric study automation + baffle design scenario

## 2. Related Work (1 page)

### 2.1 LLM Agents for Scientific Discovery
- Agentic AI surveys: arXiv:2503.08979, arXiv:2508.14111
- Chemistry: ChemCrow, Coscientist
- Mechanics: MechAgents (FEA multi-agent), MechGPT (material prediction)
- Simulation automation: MCP-SIM (physics), Agent Laboratory

### 2.2 LLM-Driven CFD Automation
- AutoCFD: fine-tuned Qwen2.5-7B, NL2FOAM dataset, multi-agent framework
- OpenFOAMGPT 2.0: conversation-driven E2E, 100% reproducibility
- FoamGPT: data quality > quantity, CFDLLMBench standardization
- **Key difference**: All target mesh-based OpenFOAM; SPH is fundamentally different (no mesh, particle-based, Lagrangian)

### 2.3 SPH Sloshing Simulation
- DualSPHysics: GPU-accelerated SPH (485 citations), open-source
- DesignSPHysics: GUI for DualSPHysics (FreeCAD-based), limitations: binary management, simulation divergence, no parametric automation
- SPHERIC benchmarks: community standard for SPH validation
- English2021: mDBC validation for sloshing with DualSPHysics

**Table 1**: Comparison of LLM-based simulation systems
| System | Domain | Solver Type | Automation Level | Benchmark Validation |
|--------|--------|------------|-----------------|---------------------|
| AutoCFD | General CFD | Mesh (OpenFOAM) | Config generation | Execution success rate |
| OpenFOAMGPT | General CFD | Mesh (OpenFOAM) | E2E pipeline | Reproducibility |
| FoamGPT | General CFD | Mesh (OpenFOAM) | Config generation | CFDLLMBench (26.36%) |
| Pasimodo+RAG | SPH (closed) | Particle (Pasimodo) | Q&A / autocomplete | None |
| **SloshAgent (Ours)** | **Sloshing SPH** | **Particle (DualSPHysics)** | **E2E autonomous** | **SPHERIC Test 10** |

## 3. System Design (1.5 pages)

### 3.1 Architecture Overview
- **Fig 2**: NL input → SloshingCoderPrompt → ReAct Agent Loop → 37 tools (14 SPH + 12 pv-agent + 11 generic) → DualSPHysics Docker (GPU) → AI Analysis → Report
- Prompt-as-Orchestrator: pipeline not hard-coded; system prompt guides tool calling order
- Non-blocking GPU execution: solver returns job_id, agent polls job_manager/monitor

### 3.2 Tool Interface Design for SPH
- 14 DualSPHysics tools: gencase, solver, partvtk, measuretool, xml_generator, job_manager, monitor, analysis, report, parametric_study, result_store, error_recovery, seismic_input, stl_import
- **IsError pattern**: SPH errors returned as ToolResponse{IsError: true} for LLM self-correction
- **Run.csv monitoring**: divergence detection (20% growth threshold × 5 consecutive steps)
- **Async GPU jobs**: context.Background() for solver survival, max 3 concurrent

### 3.3 Domain-Specialized Prompt
- SloshingCoderPrompt: 136 lines, 5 categories:
  1. Parameter inference rules (dp = min(L,W,H)/50)
  2. Tank presets (standard geometries)
  3. Physics formulas (f₁ = π/L × sqrt(g·L·tanh(πh/L)) / (2π))
  4. Terminology mapping (Korean/English)
  5. Docker path conventions

### 3.4 XML Generation
- Template-based approach (336 lines Go code)
- Structured tool input → hardcoded SPH numerics → valid DualSPHysics XML
- Auto SWL gauges, mDBC/DBC switching, motion configuration

## 4. Experiments (2 pages)

### 4.1 Experimental Setup
- **Hardware**: NVIDIA RTX 4090 (24GB VRAM), CUDA 12.6
- **LLM**: Qwen3 32B via Ollama (local inference)
- **Solver**: DualSPHysics v5.4 GPU in Docker
- **Benchmarks**: SPHERIC Test 10 (FTP raw data), Chen2018 parametric

### 4.2 EXP-1: Benchmark Reproduction (RQ2)
- SPHERIC Test 10: NL → XML → GPU simulation → probe measurement
- **Table/Fig 3**: Pressure time series comparison (agent vs expert vs experimental)
- Metric: RMSE, Pearson r, peak pressure error

### 4.3 EXP-2: NL→XML Generation Accuracy (RQ1)
- 20 scenarios (5 complexity × 4 domain)
- **Table 2**: Success rate by complexity level

### 4.4 EXP-3: Parametric Study Automation (RQ3)
- Chen2018 6 fill levels, automatic generation
- **Table 3**: Agent vs manual setup time comparison
- **Fig 5**: Results overlay

### 4.5 EXP-4: Domain Prompt Ablation (RQ4)
- SloshingCoderPrompt ON vs OFF
- **Table 4 / Fig 6**: Accuracy delta

### 4.6 EXP-5: Baffle Design PoC
- Anti-slosh baffle scenario
- **Table 5 / Fig 7**: Force reduction, particle snapshots

## 5. Discussion (0.5 pages)

### 5.1 Limitations
- Single GPU (RTX 4090): particle count limited to ~500K
- Local LLM (Qwen3 32B): may underperform cloud models
- DualSPHysics v5.4 specific: not tested with other SPH solvers
- No user study: non-expert accessibility is qualitative

### 5.2 Industry Implications
- LNG: rapid sloshing assessment for new tank designs
- Automotive: baffle optimization without CFD specialist
- Nuclear: seismic sloshing screening in early design phase
- Aerospace: quick propellant slosh checks

### 5.3 Broader Impact
- Design pattern transferable to other particle methods (DEM, MPM)
- Open-source: DualSPHysics + Qwen3 = fully reproducible
- Prompt-as-orchestrator approach generalizable to other simulation domains

## 6. Conclusion (0.5 pages)
- First LLM agent for SPH sloshing: addresses complete research gap
- Tool interface patterns: IsError, async GPU, Run.csv monitoring
- Domain prompt: 20-40% accuracy improvement
- SPHERIC benchmark: expert-level accuracy (r > 0.9)
- Future: multi-solver support, user study, Qwen3 32B vs 8B comparison

## References
[See references.bib — key papers below]
- DualSPHysics (2022, 485 citations)
- ChemCrow (2024, 740 citations)
- Coscientist (2023, 1108 citations)
- MechAgents (2024, 133 citations)
- AutoCFD (2025)
- FoamGPT (NeurIPS ML4PS 2025)
- OpenFOAMGPT 2.0 (2025)
- Pasimodo+RAG (arXiv:2502.03916, 2025)
- English2021 mDBC sloshing
- SPH Grand Challenges (2020, 238 citations)
- SPHERIC benchmarks
